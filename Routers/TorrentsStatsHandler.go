package Routers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"log"
	"net/http"
)

func TorrentStatsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var stats []Models.TorrentStatsPart

	App.Db.Model(&Models.TorrentStatsPart{}).Order("id desc").Limit(24).Find(&stats)

	len := len(stats)

	for i, j := 0, len-1; i < j; i, j = i+1, j-1 {
		stats[i], stats[j] = stats[j], stats[i]
	}

	data, err := json.Marshal(stats)

	if err != nil {
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
