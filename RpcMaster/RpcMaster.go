package main

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
)



func main() {

	app := Config.NewApp()

	Services.SetupRpc(app)
	app.Run()
}

