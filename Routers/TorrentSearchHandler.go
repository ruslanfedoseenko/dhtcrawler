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
	"github.com/lestrrat/go-ngram"
	"github.com/op/go-logging"
)
var searchLog = logging.MustGetLogger("SearchHandler")
func TorrentSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	searchTerm := ps.ByName("term")
	searchLog.Info("Search Term:", searchTerm)
	gramsTerm := ""
	words := strings.Split(searchTerm, " ")
	for _, word := range words{
		wordLen := len(word)
		if wordLen == 0 {
			continue
		}
		oneGram := ngram.NewTokenize(1, word)
		triGram := ngram.NewTokenize(3, word)
		twoGran := ngram.NewTokenize(2, word)
		for _, s := range oneGram.Tokens() {

			gramsTerm += s.String()
			gramsTerm += " "
		}
		if wordLen > 1 {
			for _, s := range twoGran.Tokens() {

				gramsTerm += s.String()
				gramsTerm += " "
			}
			if wordLen > 2 {
				for _, s := range triGram.Tokens() {

					gramsTerm += s.String()
					gramsTerm += " "
				}
			}
		}

	}

	searchLog.Info("Tokenized:", gramsTerm)
	pageNumberStr := ps.ByName("pageNumber")
	var page uint64 = 1
	var err error
	if len(pageNumberStr) > 0 {
		page, err = strconv.ParseUint(pageNumberStr,10,64)
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

	nameQuery := " 'name:("+ strings.Replace(gramsTerm, "%", " ", 500) + ")'"

	var countHolder Models.ZdbEstimateCountHolder
	App.Db.Debug().Raw("select zdb_estimate_count as count from zdb_estimate_count('torrents',"+nameQuery + ")").Scan(&countHolder)
	itemsCount = countHolder.Count
	App.Db.Debug().
		Table("torrents").
		Select("zdb_score('torrents', torrents.ctid) as score, id, group_id, leechers, seeds, infohash, name").
		Where("zdb('torrents', ctid) ==> " + nameQuery ).
		Order("zdb_score('torrents', torrents.ctid) desc").
		Limit(itemsPerPage).
		Offset((page - 1) * itemsPerPage).
		Scan(&torrents)

	for i := range torrents {
		var files []Models.File
		App.Db.Model(&torrents[i]).Association("Files").Find(&files)
		torrents[i].FilesTree = Models.BuildTree(files)
		App.Db.Model(&torrents[i]).Related(&torrents[i].Titles, "Titles")
		App.Db.Model(&torrents[i]).Association("Group").Find(&torrents[i].Group)
	}

	var pageCountFix uint64= 0
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

