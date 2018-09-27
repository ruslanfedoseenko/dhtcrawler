package Routers

import (
	"github.com/abbot/go-http-auth"
	"github.com/julienschmidt/httprouter"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"net/http"
	"path"
	"github.com/ruslanfedoseenko/dhtcrawler/Middleware"
	"github.com/op/go-logging"
)

var App *Config.App

type HostSwitch map[string]http.Handler
var RouterLog = logging.MustGetLogger("Router")
// Implement the ServerHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	//log.Println("Trying to serve host", r.Host)
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		// Handle host names for wich no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func Secret(user, realm string) string {
	users := map[string]string{
		"admin": "$apr1$upxZKqIX$nHMIFp/bAjG8VuZpIzKEP1",
	}

	if a, ok := users[user]; ok {
		return a
	}
	return ""
}
func setupAdminHandlers(router *httprouter.Router) {
	authentificactor := auth.NewBasicAuthenticator("Search Admin", Secret)
	adminFolderPath := path.Join(App.Config.HttpConfig.StaticDataFolder, "admin")
	strippedHandler := http.StripPrefix(
		"/admin",
		http.FileServer(http.Dir(adminFolderPath)),
	)

	router.HandlerFunc("GET", "/admin/*filepath",
		auth.JustCheck(
			authentificactor,
			strippedHandler.ServeHTTP),
	)

}
func Setup(app *Config.App) {
	App = app
	go func() {
		router := httprouter.New()
		router.GET("/search/suggest/:term", SearchSuggestHandler)
		router.GET("/torrents/count/", TorrentCountHandler)
		router.GET("/torrent/info/:infohash", TorrentInfoHandler)
		router.GET("/torrents/", TorrentsListHandler)
		router.GET("/torrents/stats/", TorrentStatsHandler)
		router.GET("/torrents/page/:pageNumber", TorrentsListHandler)
		router.GET("/torrents/search/:term", TorrentSearchHandler)
		router.GET("/torrents/search/:term/page/:pageNumber", TorrentSearchHandler)
		RouterLog.Info("Starting Http Backend")
		setupAdminHandlers(router)
		setupAuthHandlers(router)
		router.ServeFiles("/index/*filepath", http.Dir(App.Config.HttpConfig.StaticDataFolder))

		hs := make(HostSwitch)
		hs["btoogle.com"] = router
		hs["www.btoogle.com"] = router
		hs["localhost:6060"] = router

		http.ListenAndServe(":6060", hs)

	}()
}
func setupAuthHandlers(router *httprouter.Router) {
	router.POST("/auth/login", AuthLoginHandler)
	router.POST("/auth/register", AuthRegistrationHandler)
	router.GET("/auth/logout", Middleware.AuthRequest(AuthLogoutHandler))
	router.GET("/auth/currentUser", Middleware.AuthRequest(AuthCurrentUserInfo))
}
