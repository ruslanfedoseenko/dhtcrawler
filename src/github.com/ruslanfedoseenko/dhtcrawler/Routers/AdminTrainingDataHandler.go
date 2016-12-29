package Routers

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

func AdminTrainingDataHandler(w http.ResponseWriter, r *http.Request) {
	var trainingTokens []Models.TrainingPortion

	App.Db.Find(&trainingTokens)
	data, err := json.Marshal(trainingTokens)
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func AdminTrainingDataHandlerAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var addtrainingToken Models.AddTrainingToken

	err := json.NewDecoder(r.Body).Decode(&addtrainingToken)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = App.Db.Create(&(addtrainingToken.Token)).Error
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Token Recived:", addtrainingToken.Token)
	data, err := json.Marshal(addtrainingToken.Token)
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func AdminRefreshClassifier(w http.ResponseWriter, r *http.Request) {
	App.Classifier.Refresh(App.Db)
}
