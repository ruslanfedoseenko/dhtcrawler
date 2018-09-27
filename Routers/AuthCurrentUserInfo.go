package Routers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"strconv"
	"encoding/json"
	"github.com/ruslanfedoseenko/dhtcrawler/Errors"
)

func AuthCurrentUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	authUserIdStr := params.ByName("AuthUserId")

	authUserId, err := strconv.Atoi(authUserIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var data,err = json.Marshal(Models.NewError(Errors.InvalidUsername))

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
	var user Models.User

	err = App.Db.Where(&Models.User{
		Id: authUserId,
	}).First(&user).Error

	data, err := json.Marshal(&user)

	if err != nil {
		httpLog.Errorf(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
