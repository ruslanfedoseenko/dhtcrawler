package Routers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func TorrentSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	searchTerm := ps.ByName("term")
	searchTerm = strings.Replace(searchTerm, " ", "%", -1)
	pageNumberStr := ps.ByName("pageNumber")
	var page int = 1
	var err error
	if len(pageNumberStr) > 0 {
		page, err = strconv.Atoi(pageNumberStr)
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
	if len(searchTerm) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var torrents []Models.Torrent
	itemsCount := 0
	itemsPerPage := App.Config.ItemsPerPage


	App.Db.Debug().Table("torrents").Select("id, group_id, leechers, seeds, infohash, name").Where("zdb('torrents', ctid) ==> 'name:(" + strings.Replace(searchTerm, "%", " ", 500) + ")'").Count(&itemsCount).Limit(itemsPerPage).Offset((page - 1) * itemsPerPage).Scan(&torrents)

	for i := range torrents {
		var files []Models.File
		App.Db.Model(&torrents[i]).Association("Files").Find(&files)
		torrents[i].FilesTree = Models.BuildTree(files)
		App.Db.Model(&torrents[i]).Related(&torrents[i].Titles, "Titles")
		App.Db.Model(&torrents[i]).Association("Group").Find(&torrents[i].Group)
	}

	pageCountFix := 0
	if itemsCount%itemsPerPage != 0 {
		pageCountFix = 1
	}
	paginatedResponse := Models.PaginatedTorrentsResponse{
		Torrents:     torrents,
		PageCount:    itemsCount/itemsPerPage + pageCountFix,
		ItemsCount:   itemsCount,
		Page:         page,
		ItemsPerPage: itemsPerPage,
	}
	data, err := json.Marshal(paginatedResponse)

	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

