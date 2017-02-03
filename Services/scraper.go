package Services

import (
	"sync"
	"sync/atomic"

	"github.com/ruslanfedoseenko/dhtcrawler/Config"

	"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/rpc"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
)

var scrapeerLog = logging.MustGetLogger("Scraper")



type Scraper struct {
	workersDone             sync.WaitGroup
	quit 			chan bool
	numThreads              int
	trackerUrls             []string
	workerId                int32
	rpcClient 		*Rpc.RpcClient
}

const DefaultWorkSize int32 = 74

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
	s.workersDone.Wait()
	scrapeerLog.Info("waiting for previos task completed")
	s.initChannels()
	s.workersDone.Add(s.numThreads)
	for i := 0; i < s.numThreads; i++ {
		go s.scrapeThreadWorker()
	}
	scrapeerLog.Info("Scraping started")
}

func minInt32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func (s *Scraper) initChannels() {
	s.quit = make(chan bool, s.numThreads + 1)
}

func (s *Scraper) scrapeThreadWorker() {
	workerId := atomic.LoadInt32(&s.workerId)
	atomic.AddInt32(&s.workerId, 1)

	for {

		if !s.rpcClient.HasAvailableTasks() {
			break
		}
		task := s.rpcClient.GetNextScrapeTask()
		scrapeerLog.Info("ScrapeWorker",workerId, "recived work", task)


		if len(task.InfoHashes) > 0 {
			var result Rpc.ScrapeResult
			result.Response = new(tracker.ScrapeResponse)
			result.Response.ScrapeDatas = make(map[string]tracker.ScrapeTorrentInfo)
			for index, trackerUrl := range s.trackerUrls {
				scrapeResponse, err := tracker.Scrape(trackerUrl, task.InfoHashes)
				scrapeerLog.Info("ScrapeResult from",trackerUrl, "Len", len(scrapeResponse.ScrapeDatas), scrapeResponse.ScrapeDatas)
				if err != nil {
					scrapeerLog.Error("failed scraping tracker", trackerUrl, err.Error())
					continue
				}
				if len(scrapeResponse.ScrapeDatas) == 0 {
					scrapeerLog.Info("Empty scrape response from tracker", trackerUrl)
					if index + 1 != len(s.trackerUrls) {
						continue;
					}

				}
				for key, value := range scrapeResponse.ScrapeDatas {
					result.Response.ScrapeDatas[key] = value
					itemToRemove := Utils.IndexOf(task.InfoHashes, key)
					if (itemToRemove > -1) {
						task.InfoHashes = append(task.InfoHashes[:itemToRemove], task.InfoHashes[itemToRemove + 1:]...)
					} else {
						scrapeerLog.Error("Failed to find", key, "Task", task)
					}

				}
				infoHashesLen := len(task.InfoHashes)
				scrapeerLog.Info("InfoHashesLeft", infoHashesLen)
				if infoHashesLen == 0 {
					break;
				}

			}
			s.rpcClient.ReportScrapeResults(&result)

		}

		select {
			case <-s.quit:{
				scrapeerLog.Info("Recieved Quit Signal")
				break;
			}
			default:{
				scrapeerLog.Info("Not recieved Quit Signal")
			}
		}



	}
	s.workersDone.Done()
	scrapeerLog.Info("Exit scrapeThreadWorker", workerId)
}
