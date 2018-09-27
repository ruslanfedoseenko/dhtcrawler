package Routers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"net/http"
)

func TorrentInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var infohash string = ps.ByName("infohash")
	if len(infohash) != 40 {
		httpLog.Errorf("Invalid infohash specified: %s", infohash)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var torrent Models.Torrent
	App.Db.
		Preload("Files").
		Preload("Titles").
		Preload("Tags").
		Preload("ScraperResults").
		First(&torrent, map[string]interface{}{"infohash": infohash})

	torrent.FilesTree = Models.BuildTree(torrent.Files)

	data, err := json.Marshal(torrent)

	if err != nil {
		httpLog.Error(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
