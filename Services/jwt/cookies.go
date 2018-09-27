package jwt

import (
	"log"
	"time"
	"net/http"
)

func NullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {

	var isHttps = r.Header.Get("X-Forwarded-Proto") == "https"

	authCookie := http.Cookie{
		Name: "AuthToken",
		Value: "",
		Expires: time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
		Domain: r.Host,
		Path: "/api",
		Secure: isHttps,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name: "RefreshToken",
		Value: "",
		Expires: time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
		Domain: r.Host,
		Path: "/api",
		Secure: isHttps,
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

	var isHttps = r.Header.Get("X-Forwarded-Proto") == "https"
	authCookie := http.Cookie{
		Name: "AuthToken",
		Value: authTokenString,
		HttpOnly: true,
		Domain: r.Host,
		Path: "/api",
		Secure: isHttps,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name: "RefreshToken",
		Value: refreshTokenString,
		HttpOnly: true,
		Domain: r.Host,
		Path: "/api",
		Secure: isHttps,
	}

	http.SetCookie(*w, &refreshCookie)
}