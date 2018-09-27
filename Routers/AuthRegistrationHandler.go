package Routers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Errors"
	"net/mail"
	"golang.org/x/crypto/bcrypt"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/jwt"
	"strconv"
)

func AuthRegistrationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var user Models.User
	decoder.Decode(&user)
	validateUser(w, user)
	user.Role = "user"
	e := saveUser(&user)
	if e != nil {
		err := Models.NewErrorAddText(Errors.InvalidMail, e)
		bytes, e := json.Marshal(err)
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
		return
	}

	authTokenString, refreshTokenString, csrfSecret, err := jwt.CreateNewTokens(strconv.Itoa(user.Id), user.Role)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	// set the cookies to these newly created jwt's
	jwt.SetAuthAndRefreshCookies(&w, r, authTokenString, refreshTokenString)
	w.Header().Set("X-CSRF-Token", csrfSecret)

	w.WriteHeader(http.StatusOK)
}
func saveUser(user *Models.User) (err error) {

	user.PasswordHash, err = generateBcryptHash(user.PasswordHash)
	if err != nil {
		return
	}
	err = App.Db.Create(user).Error
	return
}

func generateBcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash[:]), err
}

func validateUser(w http.ResponseWriter, user Models.User) {
	userCount := 0
	App.Db.Where(&Models.User{
		Username: user.Username,
	}).Count(&userCount)

	if userCount > 0 {
		w.WriteHeader(http.StatusBadRequest)
		bytes, e := json.Marshal(Models.NewError(Errors.UserNameAlreadyExists))
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
		return
	}

	App.Db.Where(&Models.User{
		Mail: user.Mail,
	}).Count(&userCount)

	if userCount > 0 {
		w.WriteHeader(http.StatusBadRequest)
		bytes, e := json.Marshal(Models.NewError(Errors.UserEmailIsUsed))
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
		return
	}
	_, e := mail.ParseAddress(user.Mail)
	if e != nil {
		err := Models.NewErrorAddText(Errors.InvalidMail, e)
		bytes, e := json.Marshal(err)
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
		return
	}
}
