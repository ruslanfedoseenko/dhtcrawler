package Services

import (
	"github.com/jasonlvhit/gocron"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"time"
)

var statsUpdLog = logging.MustGetLogger("StatsUpdater")

type StatsUpdater struct {
}

func (updater *StatsUpdater) updateStatistics() {

	statsUpdLog.Debug("Starting Update Torrent Statistic")
	var lastInserted Models.TorrentStatsPart

	err := App.Db.Model(&Models.TorrentStatsPart{}).Where("id = ( SELECT MAX( id )	FROM  torrent_stats )").Scan(&lastInserted).Error

	if err != nil {
		statsUpdLog.Error("Select last id: Error occured", err)
	}

	if lastInserted.Date.Hour() == time.Now().Hour() {
		_, time := gocron.NextRun()
		statsUpdLog.Info("Next Update Statistics task run in:", time)
		return
	}

	err = App.Db.Exec("INSERT INTO  torrent_stats (  torrents_count ,  files_count ,  files_size ) SELECT  torrent_count ,  files_count ,  total_file_size FROM  realtime_counters").Error

	if err != nil {
		statsUpdLog.Error("App.Db.Create: Error occured", err)
	}
	_, time := gocron.NextRun()
	statsUpdLog.Info("Next Update Statistics task run in:", time)
}

func (updater StatsUpdater) Start() {
	go updater.updateStatistics()
}

func SetupTorrentCountStatsUpdater(app *Config.App) {
	App = app
	statsUpdater := StatsUpdater{}
	App.AddService(statsUpdater)
	App.Scheduler.Every(1).Hour().Do(statsUpdater.Start)

}
