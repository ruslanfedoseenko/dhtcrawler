package Routers

import (

	"encoding/json"
	"log"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func AdminCategoriesTreeHandler(w http.ResponseWriter, r *http.Request) {
	var categories []Models.GeneralCategory
	App.Db.Model(&Models.GeneralCategory{}).Where("parent_id IS NULL").Find(&categories)
	if len(categories) > 0 {
		for i, category := range categories {
			App.Db.Model(categories[i]).Related(&categories[i].TrainingData, "TrainingData")
			categories[i].Children = make([]Models.GeneralCategory, 10)
			App.Db.Model(categories[i]).Association("Children").Find(&categories[i].Children)
			for j, subCategory := range categories[i].Children {
				App.Db.Model(categories[i].Children[j]).Related(&categories[i].Children[j].TrainingData, "TrainingData")
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

func AdminCategoriesTreeHandlerAddToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var addTrainingTokenRequest Models.AddTrainingToken

	err := json.NewDecoder(r.Body).Decode(&addTrainingTokenRequest)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println("Adding token", addTrainingTokenRequest.Token, "to category:", addTrainingTokenRequest.GroupID)
	err = App.Db.Exec("INSERT INTO `training_data_to_gemeral_groups` (`training_portion_id`, `general_category_id`) VALUES (?, ?)", addTrainingTokenRequest.Token.Id, addTrainingTokenRequest.GroupID).Error

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB error"))
	}

}
