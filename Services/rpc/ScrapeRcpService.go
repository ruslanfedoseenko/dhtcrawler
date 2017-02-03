package Rpc

import (
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/jinzhu/gorm"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
	"sync/atomic"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"time"
)

var scrapeRpcServiceLog = logging.MustGetLogger("ScrapeRpcService")


type ScrapeRpcService struct{
	lastScrapedId int32
	lastTorrentId int32
	hasAvailableTasks bool
	scrapeTimeOut int
	db *gorm.DB
	results chan *ScrapeResult

}

func NewScrapeRpcService(app *Config.App) *ScrapeRpcService {
	service := ScrapeRpcService{
		db: app.Db,
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

type ScrapeTask struct{
	InfoHashes []string
	LastID uint32
}

type ScrapeResult struct{
	Response *tracker.ScrapeResponse
	LastID uint32
}
func (s *ScrapeRpcService) getLastScrapeID() int32{
	return atomic.LoadInt32(&s.lastScrapedId)
}
func (s *ScrapeRpcService) setLastScrapeId(value int32){
	atomic.StoreInt32(&s.lastScrapedId, value)
	s.db.Exec("UPDATE `realtime_counters` SET last_scraped_id = ?", value)
}

func (s *ScrapeRpcService) HasAvailableTasks() bool {
	return s.hasAvailableTasks
}

func (s *ScrapeRpcService) GetNextScrapeTask() *ScrapeTask {
	lastScrapeID := s.getLastScrapeID()

	var torrent Models.Torrent
	s.db.Model(&Models.Torrent{}).Order("id desc").First(&torrent)
	s.lastTorrentId = torrent.Id

	var counters Models.Counters
	err := s.db.First(&counters).Error
	if err != nil {
		scrapeRpcServiceLog.Error("Error reading counters from db:",err)
	}
	if (lastScrapeID == -1) {
		scrapeRpcServiceLog.Info("Last scraped", counters.LastScrapedId, "Torrent Count", counters.TorrentCount)
		s.setLastScrapeId(counters.LastScrapedId)
		lastScrapeID = counters.LastScrapedId
	} else {
		scrapeRpcServiceLog.Info("Db Last scraped", counters.LastScrapedId, "Torrent Count", counters.TorrentCount, "Local LastScraped ID", lastScrapeID)
		if (s.lastTorrentId - lastScrapeID < 100) {
			lastScrapeID = 1
			s.hasAvailableTasks = false
		}
	}

	var torrents []Models.Torrent
	err = s.db.Model(&Models.Torrent{}).
		Where("((extract( epoch from (now() - last_scrape)))/3600 > ? OR last_scrape IS NULL)"+
		" AND id >= ?", s.scrapeTimeOut, lastScrapeID).Limit(DefaultWorkSize).Scan(&torrents).Error
	if (err != nil) {
		scrapeRpcServiceLog.Error("Failed to get torrents",err)
	}
	if len(torrents) == 0 {
		scrapeRpcServiceLog.Error("Failed to get more torretns from ",lastScrapeID)
	}
	var task ScrapeTask
	task.InfoHashes = make([]string, len(torrents))
	task.LastID = uint32(torrents[len(torrents)-1].Id)
	defer s.setLastScrapeId(int32(task.LastID))
	for i, torrent := range torrents{
		task.InfoHashes[i] = torrent.Infohash
	}
	return &task
}

func (s *ScrapeRpcService) resultsWriterThread(){
	for {
		result := <- s.results
		tx := s.db.Begin()
		tx.Exec("SET zombodb.batch_mode = true;")
		for key, value := range result.Response.ScrapeDatas {
			scrapeRpcServiceLog.Info("Writing to db", key, "Leechers", value.Leechers, "Seeds", value.Seeders, "Completed" ,value.Completed)
			err := tx.Debug().Model(&Models.Torrent{}).
				Where(map[string]interface{}{"infohash": key}).
				Update(map[string]interface{}{"Leechers": value.Leechers, "Seeds": value.Seeders + value.Completed, "LastScrape": time.Now()}).Error
			if err != nil {
				scrapeRpcServiceLog.Error("Error updating S L C", err)
				tx.Commit()
				tx = s.db.Begin()
			}
		}
		tx.Commit()
	}
}

func (s *ScrapeRpcService) ReportScrapeResults(result *ScrapeResult){

	if (len(result.Response.ScrapeDatas) != DefaultWorkSize) {
		var torrent Models.Torrent
		err := s.db.Model(&Models.Torrent{}).
			Where("extract( epoch from (now() - last_scrape))/3600 > ? OR last_scrape IS NULL",
			s.scrapeTimeOut).First(&torrent).Error
		if (err != nil){
			scrapeRpcServiceLog.Error("Failed to find lastscrapeId",err)
		}
		if (torrent.Id < s.lastTorrentId){
			s.setLastScrapeId(torrent.Id)
		}
	} else {
		if (result.LastID > uint32(s.getLastScrapeID())){
			s.setLastScrapeId(int32(result.LastID))

		}
	}
	scrapeRpcServiceLog.Debug("Results Queue Len", len(s.results))
	s.results <- result


}
