package main

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
	_"net/http/pprof"
	"log"
	"net/http"
)

func main() {

	app := Config.NewApp()

	Services.SetupScrape(app)
	Services.SetupDhtCrawling(app)
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()
	app.Run()
}
