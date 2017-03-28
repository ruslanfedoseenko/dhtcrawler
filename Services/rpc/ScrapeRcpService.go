package Rpc

import (
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/jinzhu/gorm"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
	"sync/atomic"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"time"
	"github.com/jasonlvhit/gocron"
	"github.com/lib/pq"
	"sync"
)

var scrapeRpcServiceLog = logging.MustGetLogger("ScrapeRpcService")

type ScrapeRpcService struct {
	lastScrapedId     int32
	lastTorrentId     int32
	hasAvailableTasks bool
	scrapeTimeOut     int
	db                *gorm.DB
	scheduler         *gocron.Scheduler
	results           chan *ScrapeResult
	lastScrapeMutex   sync.Mutex
}

func NewScrapeRpcService(app *Config.App) *ScrapeRpcService {
	service := ScrapeRpcService{
		db: app.Db,
		scheduler: app.Scheduler,
		lastScrapedId: -1,
		lastTorrentId: -1,
		hasAvailableTasks: true,
		scrapeTimeOut: app.Config.ScrapeConfig.ScrapeTimeout,
		results: make(chan *ScrapeResult, 65536),
	}
	for i := 0; i < 8; i++ {
		go service.resultsWriterThread()
	}

	return &service
}

const DefaultWorkSize int = 74

type ScrapeTask struct {
	InfoHashes []string
	LastID     uint32
}

type ScrapeResult struct {
	TrackerUrl string
	Response   *tracker.ScrapeResponse
	LastID     uint32
}

func (s *ScrapeRpcService) getLastScrapeID() int32 {
	return atomic.LoadInt32(&s.lastScrapedId)
}
func (s *ScrapeRpcService) setLastScrapeId(value int32) {
	atomic.StoreInt32(&s.lastScrapedId, value)
	s.db.Exec("UPDATE realtime_counters SET last_scraped_id = ?", value)
}

func (s *ScrapeRpcService) HasAvailableTasks() bool {
	return s.hasAvailableTasks
}

func (s *ScrapeRpcService) GetNextScrapeTask() *ScrapeTask {
	s.lastScrapeMutex.Lock()
	defer s.lastScrapeMutex.Unlock()
	lastScrapeID := s.getLastScrapeID()

	var torrent Models.Torrent
	s.db.Model(&Models.Torrent{}).Order("id desc").First(&torrent)
	s.lastTorrentId = torrent.Id

	var counters Models.Counters
	err := s.db.First(&counters).Error
	if err != nil {
		scrapeRpcServiceLog.Error("Error reading counters from db:", err)
	}
	if (lastScrapeID == -1) {
		scrapeRpcServiceLog.Info("Last scraped", counters.LastScrapedId, "Torrent Count", counters.TorrentCount)
		s.setLastScrapeId(counters.LastScrapedId)
		lastScrapeID = counters.LastScrapedId
	} else {
		scrapeRpcServiceLog.Info("Db Last scraped", counters.LastScrapedId, "Torrent Count", counters.TorrentCount, "Local LastScraped ID", lastScrapeID)
		if (s.lastTorrentId - lastScrapeID < 100) {
			lastScrapeID = 1
		}
	}

	var torrents []Models.Torrent
	err = s.db.Model(&Models.Torrent{}).
		Where("((extract( epoch from (now() - last_scrape)))/3600 > ? OR last_scrape IS NULL)" +
		" AND id >= ?", s.scrapeTimeOut, lastScrapeID).Limit(DefaultWorkSize).Scan(&torrents).Error
	if (err != nil) {
		scrapeRpcServiceLog.Error("Failed to get torrents", err)
	}
	if len(torrents) == 0 {
		scrapeRpcServiceLog.Error("Failed to get more torretns from ", lastScrapeID)
	}
	var task ScrapeTask
	task.InfoHashes = make([]string, len(torrents))
	task.LastID = uint32(torrents[len(torrents) - 1].Id)
	s.setLastScrapeId(int32(task.LastID))
	for i, torrent := range torrents {
		task.InfoHashes[i] = torrent.Infohash
	}
	return &task
}

func (s *ScrapeRpcService) resultsWriterThread() {
	for {
		result := <-s.results

		var infohases []string;
		for key, value := range result.Response.ScrapeDatas {
			if value.Completed + value.Leechers + value.Seeders != 0 {
				infohases = append(infohases, key)
			}

		}
		if (len(infohases) > 0) {
			var torrentsToHash map[string]Models.Torrent = make(map[string]Models.Torrent, len(infohases))
			var torrents []Models.Torrent;
			s.db.Debug().Preload("ScraperResults").Where("infohash in (?)", infohases).Find(&torrents);
			for i := 0; i < len(torrents); i++ {
				torrentsToHash[torrents[i].Infohash] = torrents[i]
			}

			for i := 0; i < len(torrents); i++ {
				torrent := torrents[i]
				info := result.Response.ScrapeDatas[torrent.Infohash]

				if info.Completed + info.Leechers + info.Seeders != 0 {
					var found bool = false;
					for j := 0; j < len(torrent.ScraperResults); j++ {
						if torrent.ScraperResults[j].TrackerUrl == result.TrackerUrl {
							found = true;
							torrent.ScraperResults[j].LastUpdate = pq.NullTime{
								Time:time.Now(),
								Valid: true,
							}
							torrent.ScraperResults[j].Leaches = info.Leechers
							torrent.ScraperResults[j].Seeds = info.Seeders
							s.db.Debug().Model(&torrent.ScraperResults[j]).Update(torrent.ScraperResults[j]);
						}
					}
					if !found {
						s.db.Debug().Model(&torrent).
							Association("ScraperResults").
							Append(&Models.ScrapeTorrentResult{
							Leaches:info.Leechers,
							Seeds: info.Seeders,
							TrackerUrl:result.TrackerUrl,
							LastUpdate: pq.NullTime{
								Time:time.Now(),
								Valid: true,
							},
						});
					}
				}
			}
		}


	}
}

func (s *ScrapeRpcService) ReportScrapeResults(result *ScrapeResult) {

	if (len(result.Response.ScrapeDatas) != DefaultWorkSize) {
		var torrent Models.Torrent
		err := s.db.Model(&Models.Torrent{}).
			Where("extract( epoch from (now() - last_scrape))/3600 > ? OR last_scrape IS NULL",
			s.scrapeTimeOut).First(&torrent).Error
		if (err != nil) {
			scrapeRpcServiceLog.Error("Failed to find lastscrapeId", err)
		}
		if (torrent.Id < s.lastTorrentId) {
			s.setLastScrapeId(torrent.Id)
		}
	} else {
		if (result.LastID > uint32(s.getLastScrapeID())) {
			s.setLastScrapeId(int32(result.LastID))

		}
	}
	scrapeRpcServiceLog.Debug("Results Queue Len", len(s.results))
	s.results <- result

}
