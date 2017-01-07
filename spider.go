package main

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
	"github.com/ruslanfedoseenko/dhtcrawler/Routers"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/op/go-logging"
	"log"
	"os"
	//"github.com/ruslanfedoseenko/dhtcrawler/Services/tracker"
)

var logger = logging.MustGetLogger("Main")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module} %{shortfunc} - %{level:.4s} %{color:reset}: %{message}`,
)
func main() {

	setupLog();

	/*var infoHashes = []string {"00006a32d5ccef525cbc82669ac2dadb2a9994ff"};
	response, err := tracker.Scrape("http://p4p.arenabg.com:1337/announce", infoHashes)
	if err != nil {
		logger.Error(err.Error());
	}
	for k,v := range response.ScrapeDatas {
		logger.Debug("Key", k, "Value", v)
	}

	return;*/
	app := Config.NewApp()

	Routers.Setup(app)

	Services.SetupScrape(app)
	Services.SetupTorrentCountStatsUpdater(app)
	Services.SetupDhtCrawling(app)
	app.Run()
}


func setupLog(){
	file,err := os.OpenFile("crawler.log", os.O_APPEND|os.O_WRONLY| os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	fileBackend := logging.NewLogBackend(file, "", log.LstdFlags)
	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetFormatter(format)
	logging.SetBackend(fileBackend, stdOutBackend)
}