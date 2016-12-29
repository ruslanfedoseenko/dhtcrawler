package Services

import (
	"sync"
	"sync/atomic"
	"time"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
	"github.com/op/go-logging"
)

var scrapeerLog = logging.MustGetLogger("Scraper")

type result struct {
	response   tracker.ScrapeResponse
	trackerUrl string
	startId    int32
	job        work
}

type Scraper struct {
	tasks                   chan work
	results                 chan result
	workersDone             sync.WaitGroup
	quit 			chan bool
	numThreads              int
	torrentsCount           int32
	startId                 int32
	startIdGuard            sync.Mutex
	dbGuard                 sync.Mutex
	trackerUrls             []string
	workerId                int32
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

	scrapeerLog.Info("Scrape will repeat every", uint64(App.Config.ScrapeConfig.ScrapeTimeout), "hours")
	App.Scheduler.Every(uint64(App.Config.ScrapeConfig.ScrapeTimeout)).Hours().Do(scraper.Start)
	App.AddService(scraper)

}

func (s *Scraper) Start() {
	scrapeerLog.Info("waiting for previos task")
	s.workersDone.Wait()
	s.initChannels()
	s.workersDone.Add(s.numThreads)
	for i := 0; i < s.numThreads; i++ {
		go s.scrapeThreadWorker()
	}
	go s.workManagerThread()
	go s.resultsManagementThread()
	scrapeerLog.Info("Scraping started")
}

func minInt32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func (s *Scraper) initChannels() {
	s.tasks = make(chan work, s.numThreads*4)
	s.results = make(chan result, s.numThreads*16)
	s.quit = make(chan bool)
}

func (s *Scraper) workManagerThread() {

	var counters Models.Counters
	err := App.Db.First(&counters).Error
	if err != nil {
		scrapeerLog.Error("Error reading counters from db:",err)
	}
	scrapeerLog.Info("Last scraped", counters.LastScrapedId, "Torrent Count", counters.TorrentCount)
	s.torrentsCount = counters.TorrentCount

	var torrents []Models.Torrent
	err = App.Db.Model(&Models.Torrent{}).
		Where("extract( epoch from (now() - last_scrape))/3600 > ? OR last_scrape IS NULL",
			App.Config.ScrapeConfig.ScrapeTimeout).Limit(DefaultWorkSize).Find(&torrents).Error
	if err != nil {
		scrapeerLog.Error("Find valid torrents for scrape err:",err)
	}
	if len(torrents) == 0 {
		scrapeerLog.Info("Scrape not needed all is up to date... Exiting now")
		s.quit <- true
		close(s.tasks)
		s.workersDone.Wait()
		close(s.results)
		return
	}
	s.startId = minInt32(counters.LastScrapedId, torrents[0].Id)

	for i := 0; i < s.numThreads*2; i++ {
		if s.torrentsCount <= s.startId {
			break
		}
		workSize := minInt32(DefaultWorkSize, s.torrentsCount-s.startId)
		s.tasks <- work{
			count:      workSize,
			startId:    s.startId,
		}

		s.startId = atomic.AddInt32(&s.startId, workSize)
	}

	for s.torrentsCount > s.startId {
		scrapeerLog.Debug("Scrape Step:", s.startId)
		workSize := minInt32(DefaultWorkSize, s.torrentsCount-s.startId)
		time.Sleep(200 * time.Millisecond)
		s.tasks <- work{
			count:      workSize,
			startId:    s.startId,
		}
		s.startId = atomic.AddInt32(&s.startId, workSize)
	}
	scrapeerLog.Info("Waiting for workers")
	s.quit <- true;
	s.workersDone.Wait()
	scrapeerLog.Info("Exit workManagerThread")
}

func (s *Scraper) resultsManagementThread() {

	var canQuit bool = false
	for {

		scrapeerLog.Info("resultsManagementThread waiting for result")
		select {
			case result := <-s.results:
			{
				resultLen := len(result.response.ScrapeDatas)
				if resultLen == 0 {
					scrapeerLog.Info("Empty scrape result received.", result.trackerUrl,  result.job.startId)

					continue
				}
				s.dbGuard.Lock()
				/*if (int32(resultLen) != result.job.count){
					if (s.startId > result.job.startId){
						trackerIndex := Utils.IndexOf(s.trackerUrls, result.job.trackerUrl) + 1
						if (trackerIndex < len(s.trackerUrls)) {
							scrapeerLog.Println("Recived not full results from tracker", result.job.trackerUrl, "restarting with new tracker", s.trackerUrls[trackerIndex])
							workSize := minInt32(DefaultWorkSize, s.torrentsCount - s.startId)
							heap.Push(s.workPriorityQueue, &Item{
								priority:High,
								value:work{
									count:workSize,
									startId: result.startId,
									trackerUrl: s.trackerUrls[trackerIndex],
								},
							})
							var torrents []Models.Torrent
							App.Db.Debug().Model(&Models.Torrent{}).
								Where("(HOUR(TIMEDIFF(CURRENT_TIMESTAMP,`last_scrape`)) > ? OR `last_scrape` IS NULL)" +
								" AND `id` >= ?", App.Config.ScrapeConfig.ScrapeTimeout, s.startId).Limit(workSize + 1).Find(&torrents)
							if (len(torrents) > 0){
								s.startId = torrents[len(torrents) - 1].Id
							}

						} else {
							scrapeerLog.Println("Recived not full results from tracker", result.job.trackerUrl, "but no more trackers left")
						}


					}
				}*/

				App.Db.Exec("UPDATE `realtime_counters` SET last_scraped_id = ?", result.startId+result.job.count)
				for key, value := range result.response.ScrapeDatas {
					scrapeerLog.Info("Writing to db", key, "Leechers", value.Leechers, "Seeds", value.Seeders, "Completed" ,value.Completed)
					App.Db.Debug().Model(&Models.Torrent{}).
						Where(map[string]interface{}{"infohash": key}).
						Update(map[string]interface{}{"Leechers": value.Leechers, "Seeds": value.Seeders + value.Completed, "LastScrape": time.Now()})

				}
				s.dbGuard.Unlock()
				taskQueued := len(s.tasks)
				resultsQueued := len(s.results)
				scrapeerLog.Debug("Task queued", taskQueued, "Results Queued", resultsQueued)
				if taskQueued == 0 && resultsQueued == 0 && canQuit {

					break
				}
			}
			case <- s.quit:
			{
				canQuit = true
				scrapeerLog.Info("resultsManagementThread may exit")

			}
		}
	}
	close(s.tasks)
	s.workersDone.Wait()
	close(s.results)
	scrapeerLog.Info("Exit resultsManagementThread")
}

func (s *Scraper) scrapeThreadWorker() {
	workerId := atomic.LoadInt32(&s.workerId)
	atomic.AddInt32(&s.workerId, 1)
	var canQuit = false
	scrapeTimeout := App.Config.ScrapeConfig.ScrapeTimeout
	for {
		select{
			case task := <-s.tasks :
			{
				scrapeerLog.Info(workerId, "recived work", task)

				result := result{
					job:        task,
					startId:    task.startId,
				}
				var torrents []Models.Torrent

				s.dbGuard.Lock()
				App.Db.Model(&Models.Torrent{}).
					Where("((extract( epoch from (now() - last_scrape)))/3600 > ? OR last_scrape IS NULL)"+
					" AND id >= ?", scrapeTimeout, task.startId).Limit(task.count).Scan(&torrents)

				s.dbGuard.Unlock()
				infoHashes := make([]string, len(torrents), len(torrents))

				for index, torrent := range torrents {
					infoHashes[index] = torrent.Infohash
				}
				if len(infoHashes) > 0 {
					for index, trackerUrl := range s.trackerUrls {
						scrapeResponse, err := tracker.Scrape(trackerUrl, infoHashes)

						if err != nil {
							scrapeerLog.Error("failed scraping tracker", trackerUrl ,err.Error())
							continue
						}
						if len(scrapeResponse.ScrapeDatas) == 0 {
							scrapeerLog.Info("Empty scrape response from tracker", trackerUrl)
							if index + 1 != len(s.trackerUrls) {
								continue;
							}

						}
						result.response = scrapeResponse
						result.trackerUrl = trackerUrl
						break;
					}


				}

				s.results <- result

				taskQueued := len(s.tasks)
				resultsQueued := len(s.results)
				scrapeerLog.Debug("Task queued", taskQueued, "Results Queued", resultsQueued)
				if taskQueued == 0 && resultsQueued == 0 && canQuit {

					break
				}
			}
			case <- s.quit:
			{
				canQuit = true;
				scrapeerLog.Info("scrapeThreadWorker may quit")
			}
		}
	}
	s.workersDone.Done()
	scrapeerLog.Info("Exit scrapeThreadWorker", workerId)
}
