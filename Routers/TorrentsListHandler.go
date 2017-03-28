package Routers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func TorrentsListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var page uint64 = 1
	var err error
	pageNumberStr := ps.ByName("pageNumber")
	if len(pageNumberStr) > 0 {

		page, err = strconv.ParseUint(pageNumberStr, 10, 64)
		if err != nil {
			log.Println("Parsing Error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var itemsPerPage uint64 = App.Config.ItemsPerPage
	var torrents []Models.Torrent = make([]Models.Torrent, itemsPerPage, itemsPerPage)
	var itemsCount uint64

	row := App.Db.Debug().Table("realtime_counters").Select("torrent_count").Row()
	row.Scan(&itemsCount)
	err = App.Db.Debug().
		Preload("Files").Preload("Titles").Preload("Tags").
		Model(&Models.Torrent{}).
		Limit(itemsPerPage).
		Offset((page - uint64(1)) * itemsPerPage).
		Find(&torrents).Error

	for i := range torrents {
		torrents[i].FilesTree = Models.BuildTree(torrents[i].Files)
	}

	if err != nil {
		log.Println("Error:", err.Error())
	}
	log.Println("Total items count:", itemsCount)

	var pageCountFix uint64 = 0
	if itemsCount % itemsPerPage != 0 {
		pageCountFix = 1
	}
	paginatedResponse := Models.PaginatedTorrentsResponse{
		Torrents:     torrents,
		PageCount:    itemsCount / itemsPerPage + pageCountFix,
		ItemsCount:   itemsCount,
		Page:         page,
		ItemsPerPage: itemsPerPage,
	}
	data, err := json.Marshal(paginatedResponse)

	if err != nil {
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
