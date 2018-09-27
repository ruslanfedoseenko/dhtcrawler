package Middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/jwt"
	"github.com/op/go-logging"
)

var AuthMiddlewareLog = logging.MustGetLogger("AuthMiddleware")

func AuthRequest(handleFunc httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		AuthMiddlewareLog.Info("In auth restricted section")

		// read cookies
		AuthCookie, authErr := r.Cookie("AuthToken")
		if authErr == http.ErrNoCookie {
			AuthMiddlewareLog.Error("Unauthorized attempt! No auth cookie")
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
			AuthMiddlewareLog.Error("Unauthorized attempt! No refresh cookie")
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
		AuthMiddlewareLog.Infof("Request Refresh Token: %s",requestCsrfToken)


		// check the jwt's for validity
		authTokenString, refreshTokenString, csrfSecret, err := jwt.CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value, requestCsrfToken)
		if err != nil {
			if err.Error() == "Unauthorized" {
				AuthMiddlewareLog.Error("Unauthorized attempt! JWT's not valid!")
				// Utils.NullifyTokenCookies(&w, r)
				// http.Redirect(w, r, "/login", 302)
				http.Error(w, http.StatusText(401), 401)
				return
			} else {
				// @adam-hanna: do we 401 or 500, here?
				// it could be 401 bc the token they provided was messed up
				// or it could be 500 bc there was some error on our end
				AuthMiddlewareLog.Error("err not nil")
				log.Panic("panic: %+v", err)
				// Utils.NullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}
		}
		claims, err := jwt.GetClaims(AuthCookie.Value)

		params = append(params, httprouter.Param{
			Key: "AuthUserId",
			Value: claims.Subject,
		})

		AuthMiddlewareLog.Error("Successfully recreated jwts")

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
	csrfFromCookie, err := r.Cookie("csrf")
	if csrfFromFrom != "" {
		return csrfFromFrom
	} else  if err == nil && csrfFromCookie.Value != "" {
		return csrfFromCookie.Value
	} else {
		return r.Header.Get("X-CSRF-Token")
	}
}
