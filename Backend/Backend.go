package main

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Routers"
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
)

// @APIVersion 1.0.0
// @APITitle Btoogle Torrent Search API
// @APIDescription API to search for Torrents
// @Contact ruslan.fedoseenko.91@gmail.com
// @TermsOfServiceUrl http://google.com/
// @License BSD
// @LicenseUrl http://opensource.org/licenses/BSD-2-Clause
// @BasePath http://btoogle.com/api/


func main() {

	app := Config.NewApp()

	Routers.Setup(app)

	Services.SetupTorrentCountStatsUpdater(app)
	app.Run()
}


