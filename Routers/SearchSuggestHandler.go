package Routers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"encoding/json"
	"fmt"
	"strings"
)

func SearchSuggestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	searchTerm := ps.ByName("term");

	var torrents []Models.Torrent




	words:=strings.Split(searchTerm, " ")
	lastWord := words[len(words) - 1]
	App.Db.Debug().
		Table("zdb_termlist('torrents', 'name', '" + lastWord + "', NULL, 10) ").
		Select("term as name").
		Order("totalfreq desc").
		Find(&torrents)
	var response Models.SearchSuggestResponse;
	response.Suggestions = make([]string, len(torrents))
	for i, torrent:= range torrents{
		response.Suggestions[i] = torrent.Name
	}
	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
