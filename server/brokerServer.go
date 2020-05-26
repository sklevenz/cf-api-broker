package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	brokerAPIVersion string = "2.14"
)

// NewBrokerServer implements static routes for serving a home page and the routes
// defined by OSB v2.0 API
func NewBrokerServer(staticDir string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	v2Router := router.PathPrefix("/v2/").Subrouter()
	v2Router.HandleFunc("/catalog/", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)

	v2Router.Use(headerHandler)

	router.HandleFunc("/health/", healthHandler).Name("health").Methods(http.MethodGet)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))).Name("static").Methods(http.MethodGet)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir))).Name("home").Methods(http.MethodGet)

	router.Use(logHandler)

	return router
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func headerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Broker-API-Version", brokerAPIVersion)
		next.ServeHTTP(w, r)
	})
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "catalog")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
