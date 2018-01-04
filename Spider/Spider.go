package main

import (
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
	"log"
	"os"
)

var logger = logging.MustGetLogger("Spider")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module} %{shortfunc} - %{level:.4s} %{color:reset}: %{message}`,
)

func main() {

	setupLog()

	app := Config.NewApp()

	Services.SetupScrape(app)
	Services.SetupDhtCrawling(app)
	app.Run()
}

func setupLog() {
	file, err := os.OpenFile("spider.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	fileBackend := logging.NewLogBackend(file, "", log.LstdFlags)
	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetFormatter(format)
	logging.SetBackend(fileBackend, stdOutBackend)
}
