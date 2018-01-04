package Routers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"log"
	"net/http"
)

func TorrentCountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var counters Models.Counters

	err := App.Db.Model(&Models.Counters{}).First(&counters).Error

	if err != nil {
		log.Println(err.Error())
	}
	data, err := json.Marshal(counters)
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
