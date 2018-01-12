package Routers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"log"
	"net/http"
	"strconv"
)

var searchLog = logging.MustGetLogger("SearchHandler")

func TorrentSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	searchTerm := ps.ByName("term")
	searchLog.Info("Search Term:", searchTerm)
	pageNumberStr := ps.ByName("pageNumber")
	var page uint64 = 1
	var err error
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
	if len(searchTerm) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var torrents []Models.Torrent
	var itemsCount uint64 = 0
	itemsPerPage := App.Config.ItemsPerPage

	nameQuery := Utils.ZdbBuildQuery("name", searchTerm, (page - 1) * itemsPerPage, itemsPerPage)

	var countHolder Models.ZdbEstimateCountHolder
	App.Db.Debug().Raw("select zdb_estimate_count as count from zdb_estimate_count('torrents'," + nameQuery + ")").Scan(&countHolder)
	itemsCount = countHolder.Count
	App.Db.Debug().
		Preload("ScraperResults").
		Preload("Files").
		Preload("Titles").
		Preload("Tags").
		Model(Models.Torrent{}).
		Where("zdb('torrents', ctid) ==> " + nameQuery).
		Order("zdb_score('torrents', torrents.ctid) desc").
		Find(&torrents)

	for i := range torrents {

		torrents[i].FilesTree = Models.BuildTree(torrents[i].Files)

	}

	var pageCountFix uint64 = 0
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
