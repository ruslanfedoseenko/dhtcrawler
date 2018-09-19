// @SubApi Search Helper API [/search]
package Routers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"net/http"
	"strings"
)

// @Title SearchSuggestHandler
// @Description retrieves suggestion for input search string
// @Accept  json
// @Produce  json
// @Param   term     path    string     true        "Value of Search imput"
// @Success 200 {array}  dhtcrawler.Model.SearchSuggestResponse
// @Failure 400 {object} error    "Customer ID must be specified"
// @Router /search/suggest/{term} [get]

func SearchSuggestHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	searchTerm := ps.ByName("term")

	var torrents []Models.Torrent

	words := strings.Split(strings.Trim(searchTerm, " "), " ")
	lastWord := words[len(words)-1]
	App.Db.Debug().
		Table("zdb_termlist('torrents', 'name_autocomplete', '" + lastWord + "', NULL, 5000) ").
		Select("term as name").
		Order("totalfreq desc").
		Limit(25).
		Find(&torrents)
	var response Models.SearchSuggestResponse
	response.Suggestions = make([]string, len(torrents))
	for i, torrent := range torrents {
		response.Suggestions[i] = torrent.Name
	}
	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
