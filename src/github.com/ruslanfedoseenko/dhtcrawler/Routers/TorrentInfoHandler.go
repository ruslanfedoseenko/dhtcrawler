package Routers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func TorrentInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var infohash string = ps.ByName("infohash")
	if len(infohash) != 40 {
		log.Println("Invalid infohash specified:", infohash)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var torrent Models.Torrent
	App.Db.First(&Models.Torrent{}, map[string]interface{}{"Infohash": infohash}).Scan(&torrent)
	App.Db.Model(&torrent).Association("Files").Find(&torrent.Files)
	App.Db.Model(&torrent).Association("Group").Find(&torrent.Group)
	data, err := json.Marshal(torrent)

	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
