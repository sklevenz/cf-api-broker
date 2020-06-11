package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	headerContentType string = "Content-Type"

	contentTypeCSS  string = "text/css; charset=utf-8"
	contentTypeHTML string = "text/html; charset=utf-8"
	contentTypeTEXT string = "text/plain; charset=utf-8"
	contentTypeJSON string = "application/json; charset=utf-8"
)

// NewRouter implements static routes for serving a home page and the routes
// defined by OSB v2.0 API
func NewRouter(staticDir string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	v2Router := router.PathPrefix("/v2/").Subrouter()
	v2Router.Use(apiVersionHandler)
	v2Router.Use(originatingIdentityLogHandler)
	v2Router.HandleFunc("/catalog/", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)

	router.HandleFunc("/version/", versionHandler).Name("version").Methods(http.MethodGet)
	router.HandleFunc("/health/", healthHandler).Name("health").Methods(http.MethodGet)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))).Name("static").Methods(http.MethodGet)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir))).Name("home").Methods(http.MethodGet)

	router.Use(logHandler)

	return router
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RAW Request Object: %v", r)

		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("-%v -%v -%v", r.Method, r.RequestURI, time.Since(start))
	})
}
