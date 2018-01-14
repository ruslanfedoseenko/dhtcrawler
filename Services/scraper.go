package Services

import (
	"sync"
	"sync/atomic"

	"github.com/ruslanfedoseenko/dhtcrawler/Config"

	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/rpc"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
	"time"
)

var scrapeerLog = logging.MustGetLogger("Scraper")

type Scraper struct {
	workersDone      sync.WaitGroup
	isAlreadyWaiting bool
	quit             chan bool
	numThreads       int
	trackerUrls      []string
	workerId         int32
	rpcClient        *Rpc.RpcClient
}

func NewScraper() (s *Scraper) {
	s = new(Scraper)

	s.trackerUrls = App.Config.ScrapeConfig.Trackers
	s.numThreads = App.Config.ScrapeConfig.WorkerThreads
	s.workerId = 0
	return
}

var scraper *Scraper

func SetupScrape(app *Config.App) {
	App = app

	scraper = NewScraper()
	scraper.rpcClient = Rpc.GetRpcCLientInstance(app)
	scrapeerLog.Info("Scrape will repeat every", uint64(App.Config.ScrapeConfig.ScrapeTimeout), "hours")
	App.Scheduler.Every(uint64(App.Config.ScrapeConfig.ScrapeTimeout)).Hours().Do(scraper.Start)
	App.AddService(scraper)

}
func (s *Scraper) Start() {
	go s.startInternal()
}
func (s *Scraper) startInternal() {
	scrapeerLog.Info("waiting for previos task")
	if s.isAlreadyWaiting {
		return
	}
	s.isAlreadyWaiting = true
	s.workersDone.Wait()
	s.isAlreadyWaiting = false
	scrapeerLog.Info("waiting for previos task completed")
	s.initChannels()
	s.workersDone.Add(s.numThreads)
	for i := 0; i < s.numThreads; i++ {
		go s.scrapeThreadWorker()
	}
	scrapeerLog.Info("Scraping started")
}

func (s *Scraper) initChannels() {
	s.quit = make(chan bool, s.numThreads+1)
}

func (s *Scraper) scrapeThreadWorker() {
	workerId := atomic.LoadInt32(&s.workerId)
	atomic.AddInt32(&s.workerId, 1)

	for {

		if !s.rpcClient.HasAvailableTasks() {
			break
		}
		task := s.rpcClient.GetNextScrapeTask()
		if task == nil {
			<-time.After(5 * time.Second)
			continue
		}
		scrapeerLog.Info("ScrapeWorker", workerId, "recived work", task)

		if len(task.InfoHashes) > 0 {

			for _, trackerUrl := range s.trackerUrls {
				var result Rpc.ScrapeResult
				result.TrackerUrl = trackerUrl
				var err error
				var response tracker.ScrapeResponse
				response, err = tracker.Scrape(trackerUrl, task.InfoHashes)
				result.Response = &response
				scrapeerLog.Info("ScrapeResult from", trackerUrl, "Len", len(result.Response.ScrapeDatas), result.Response.ScrapeDatas)
				if err != nil {
					scrapeerLog.Error("failed scraping tracker", trackerUrl, err.Error())
					continue
				}
				if len(result.Response.ScrapeDatas) == 0 {
					scrapeerLog.Info("Empty scrape response from tracker", trackerUrl)
					continue
				}
				scrapeerLog.Info("reporting scrape result for", trackerUrl, "result.TrackerUrl = ", result.TrackerUrl, "whole result", result)
				s.rpcClient.ReportScrapeResults(&result)
			}

		}

		select {
		case <-s.quit:
			{
				scrapeerLog.Info("Recieved Quit Signal")
				break
			}
		default:
			{
				scrapeerLog.Info("Not recieved Quit Signal")
			}
		}

	}
	s.workersDone.Done()
	scrapeerLog.Info("Exit scrapeThreadWorker", workerId)
}
