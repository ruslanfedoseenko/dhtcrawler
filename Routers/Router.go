package Routers

import (
	"github.com/abbot/go-http-auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"path"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
)

var App *Config.App

type HostSwitch map[string]http.Handler

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
func IndexRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/index/", 301)
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
	router.HandlerFunc("GET", "/adm/groups/", auth.JustCheck(authentificactor, AdminCategoriesTreeHandler))
	router.HandlerFunc("GET", "/adm/relearn/", auth.JustCheck(authentificactor, AdminRefreshClassifier))
	router.HandlerFunc("POST", "/adm/group/addtoken/", auth.JustCheck(authentificactor, AdminCategoriesTreeHandlerAddToken))
	router.HandlerFunc("GET", "/adm/trainingdata/", auth.JustCheck(authentificactor, AdminTrainingDataHandler))
	router.HandlerFunc("POST", "/adm/trainingdata/", auth.JustCheck(authentificactor, AdminTrainingDataHandlerAdd))

}
func Setup(app *Config.App) {
	App = app
	go func() {
		router := httprouter.New()
		router.GET("/", IndexRedirect)
		router.GET("/categories/", CategoriesTreeHandler)
		router.GET("/torrents/count/", TorrentCountHandler)
		router.GET("/torrent/info/:infohash", TorrentInfoHandler)
		router.GET("/torrents/", TorrentsListHandler)
		router.GET("/torrents/stats/", TorrentStatsHandler)
		router.GET("/torrents/page/:pageNumber", TorrentsListHandler)
		router.GET("/torrents/group/:groupId/page/:pageNumber", TorrentsListHandler)
		router.GET("/torrents/search/:term", TorrentSearchHandler)
		router.GET("/torrents/search/:term/page/:pageNumber", TorrentSearchHandler)
		log.Println("Starting Http Backend")
		setupAdminHandlers(router)

		router.ServeFiles("/index/*filepath", http.Dir(App.Config.HttpConfig.StaticDataFolder))

		hs := make(HostSwitch)
		hs["search.cutetorrent.info"] = router
		hs["localhost:6060"] = router

		http.ListenAndServe(":6060", hs)

	}()
}
