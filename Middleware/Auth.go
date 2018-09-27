package Middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/jwt"
)

func AuthRequest(handleFunc httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		log.Println("In auth restricted section")

		// read cookies
		AuthCookie, authErr := r.Cookie("AuthToken")
		if authErr == http.ErrNoCookie {
			log.Println("Unauthorized attempt! No auth cookie")
			jwt.NullifyTokenCookies(&w, r)
			// http.Redirect(w, r, "/login", 302)
			http.Error(w, http.StatusText(401), 401)
			return
		} else if authErr != nil {
			log.Panic("panic: %+v", authErr)
			jwt.NullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		RefreshCookie, refreshErr := r.Cookie("RefreshToken")
		if refreshErr == http.ErrNoCookie {
			log.Println("Unauthorized attempt! No refresh cookie")
			jwt.NullifyTokenCookies(&w, r)
			http.Redirect(w, r, "/login", 302)
			return
		} else if refreshErr != nil {
			log.Panic("panic: %+v", refreshErr)
			jwt.NullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// grab the csrf token
		requestCsrfToken := grabCsrfFromReq(r)
		log.Println(requestCsrfToken)

		// check the jwt's for validity
		authTokenString, refreshTokenString, csrfSecret, err := jwt.CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value, requestCsrfToken)
		if err != nil {
			if err.Error() == "Unauthorized" {
				log.Println("Unauthorized attempt! JWT's not valid!")
				// Utils.NullifyTokenCookies(&w, r)
				// http.Redirect(w, r, "/login", 302)
				http.Error(w, http.StatusText(401), 401)
				return
			} else {
				// @adam-hanna: do we 401 or 500, here?
				// it could be 401 bc the token they provided was messed up
				// or it could be 500 bc there was some error on our end
				log.Println("err not nil")
				log.Panic("panic: %+v", err)
				// Utils.NullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}
		}
		log.Println("Successfully recreated jwts")

		// @adam-hanna: Change this. Only allow whitelisted origins! Also check referer header
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// if we've made it this far, everything is valid!
		// And tokens have been refreshed if need-be
		jwt.SetAuthAndRefreshCookies(&w, r, authTokenString, refreshTokenString)
		w.Header().Set("X-CSRF-Token", csrfSecret)

		handleFunc(w, r, params)
	}
}

func grabCsrfFromReq(r *http.Request) string {
	csrfFromFrom := r.FormValue("X-CSRF-Token")

	if csrfFromFrom != "" {
		return csrfFromFrom
	} else {
		return r.Header.Get("X-CSRF-Token")
	}
}
