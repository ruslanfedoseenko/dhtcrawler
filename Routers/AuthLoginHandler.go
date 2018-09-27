package Routers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/jwt"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"github.com/ruslanfedoseenko/dhtcrawler/Errors"
)

var UserNotFound = errors.New("User not found");

func AuthLogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jwt.NullifyTokenCookies(&w, r)
	// use 302 to force browser to do GET request
	http.Redirect(w, r, "/#/login", 302)
}

func AuthLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var credentials Models.Credentials
	bytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}
	json.Unmarshal(bytes, &credentials)
	log.Println(credentials)

	user, uuid, loginErr := LogUserIn(credentials.Username, credentials.Password)
	log.Println(user, uuid, loginErr)
	if loginErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if loginErr == UserNotFound {
			var data,err = json.Marshal(Models.NewError(Errors.InvalidUsername))

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}
	} else {
		// no login err
		// now generate cookies for this user
		authTokenString, refreshTokenString, csrfSecret, err := jwt.CreateNewTokens(uuid, user.Role)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		// set the cookies to these newly created jwt's
		jwt.SetAuthAndRefreshCookies(&w, r, authTokenString, refreshTokenString)
		w.Header().Set("X-CSRF-Token", csrfSecret)

		w.WriteHeader(http.StatusOK)

		var data, err2 = json.Marshal(user)
		if err2 != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	}
}
func LogUserIn(username string, password string) (user *Models.User,uuid string,err error) {
	userCount := 0
	user = &Models.User{}

	err = App.Db.Debug().Model(&Models.User{}).Where(&Models.User{
		Username: username,
	}).First(user).Count(&userCount).Error
	if err != nil {
		return
	}

	if userCount == 0 {
		err = UserNotFound
		return nil, "", err
	}
	log.Println(user)

	return user, strconv.Itoa(user.Id), checkPasswordAgainstHash(user.PasswordHash, password)
}

func checkPasswordAgainstHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
