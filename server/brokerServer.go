package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	brokerAPIVersion string = "2.14"
)

// NewBrokerServer creates a new OSB broker server and configures handlers and routes
func NewBrokerServer() *mux.Router {
	broker := mux.NewRouter().StrictSlash(true)

	broker.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))).Name("static")
	//broker.Handle("/", http.FileServer(http.Dir(staticPath))).Name("home").Methods(http.MethodGet)
	broker.HandleFunc("/", homeHandler).Name("home").Methods(http.MethodGet)
	broker.HandleFunc("/v2/catalog/", catalogHandler).Name("catalog").Methods(http.MethodGet)

	broker.Use(logHandler)
	broker.Use(headerHandler)

	return broker
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./static"))
	fs.ServeHTTP(w, r)
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "catalog")
}
