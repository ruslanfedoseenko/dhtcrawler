package Routers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func CategoriesTreeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var categories []Models.GeneralCategory
	App.Db.Model(&Models.GeneralCategory{}).Where("parent_id IS NULL").Find(&categories)
	if len(categories) > 0 {
		for i, category := range categories {
			categories[i].Children = make([]Models.GeneralCategory, 10)
			App.Db.Model(categories[i]).Association("Children").Find(&categories[i].Children)
			for j, subCategory := range categories[i].Children {
				categories[i].Children[j].Name = category.Name + "\\" + subCategory.Name

			}
		}
	}
	data, err := json.Marshal(categories)
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
