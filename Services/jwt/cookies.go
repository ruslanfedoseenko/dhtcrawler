package jwt

import (
	"log"
	"time"
	"net/http"
)

func NullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {
	authCookie := http.Cookie{
		Name: "AuthToken",
		Value: "",
		Expires: time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name: "RefreshToken",
		Value: "",
		Expires: time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &refreshCookie)

	// if present, revoke the refresh cookie from our db
	RefreshCookie, refreshErr := r.Cookie("RefreshToken")
	if refreshErr == http.ErrNoCookie {
		// do nothing, there is no refresh cookie present
		return
	} else if refreshErr != nil {
		log.Panic("panic: %+v", refreshErr)
		http.Error(*w, http.StatusText(500), 500)
	}

	RevokeRefreshToken(RefreshCookie.Value)
}

func SetAuthAndRefreshCookies(w *http.ResponseWriter, r *http.Request, authTokenString string, refreshTokenString string) {

	authCookie := http.Cookie{
		Name: "AuthToken",
		Value: authTokenString,
		HttpOnly: true,
		Domain: r.Host,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name: "RefreshToken",
		Value: refreshTokenString,
		HttpOnly: true,
		Domain: r.Host,
	}

	http.SetCookie(*w, &refreshCookie)
}