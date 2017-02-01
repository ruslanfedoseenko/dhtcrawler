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
	var groupId int = -1
	var err error
	pageNumberStr := ps.ByName("pageNumber")
	groupIdStr := ps.ByName("groupId")
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

	if len(groupIdStr) > 0 {
		groupId, err = strconv.Atoi(groupIdStr)
		if err != nil {
			log.Println("Parsing Error", err.Error())
			return
		}
		if groupId < 1 {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	var itemsPerPage uint64 = App.Config.ItemsPerPage
	var torrents []Models.Torrent = make([]Models.Torrent, itemsPerPage, itemsPerPage)
	var itemsCount uint64

	if groupId > 0 {

		err = App.Db.Debug().
			Model(&Models.Torrent{}).
			Where(&Models.Torrent{GroupId: int32(groupId)}).
			Count(&itemsCount).Order("seeds desc,leechers desc").
			Limit(itemsPerPage).
			Offset((page - uint64(1)) * itemsPerPage).
			Find(&torrents).Error
	} else {
		row := App.Db.Debug().Table("realtime_counters").Select("torrent_count").Row()
		row.Scan(&itemsCount)
		err = App.Db.Debug().Model(&Models.Torrent{}).Limit(itemsPerPage).Offset((page - uint64(1)) * itemsPerPage).Find(&torrents).Error
	}
	if err != nil {
		log.Println("Error:", err.Error())
	}
	log.Println("Total items count:", itemsCount)

	for i := range torrents {
		var files []Models.File
		App.Db.Model(&torrents[i]).Association("Files").Find(&files)
		torrents[i].FilesTree = Models.BuildTree(files)
		App.Db.Model(&torrents[i]).Related(&torrents[i].Titles, "Titles")
		App.Db.Model(&torrents[i]).Association("Group").Find(&torrents[i].Group)
		if torrents[i].Group.ParentId != 0 {
			var parent Models.GeneralCategory
			App.Db.Model(&torrents[i].Group).Related(&parent, "ParentId")
			torrents[i].Group.Name = parent.Name + "\\" + torrents[i].Group.Name
		}

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
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
