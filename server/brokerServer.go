package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	supportedAPIVersionValue string = "2.14"

	headerAPIVersion            string = "X-Broker-API-Version"
	headerAPIOrginatingIdentity string = "X-Broker-API-Originating-Identity"
	headerAPIRequestIdentity    string = "X-Broker-API-Request-Identity"

	headerContentType string = "Content-Type"

	contentTypeCSS  string = "text/css; charset=utf-8"
	contentTypeHTML string = "text/html; charset=utf-8"
	contentTypeTEXT string = "text/plain; charset=utf-8"
	contentTypeJSON string = "application/json; charset=utf-8"
)

// NewBrokerServer implements static routes for serving a home page and the routes
// defined by OSB v2.0 API
func NewBrokerServer(staticDir string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	v2Router := router.PathPrefix("/v2/").Subrouter()
	v2Router.HandleFunc("/catalog/", catalogHandler).Name("v2.catalog").Methods(http.MethodGet)
	v2Router.Use(apiVersionHandler)

	router.HandleFunc("/version/", versionHandler).Name("version").Methods(http.MethodGet)
	router.HandleFunc("/health/", healthHandler).Name("health").Methods(http.MethodGet)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))).Name("static").Methods(http.MethodGet)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir))).Name("home").Methods(http.MethodGet)
	router.Use(logHandler)

	return router
}

func apiVersionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supportedAPIVersion := strings.Split(supportedAPIVersionValue, ".")[0]
		requestedAPIVersionValue := strings.Split(r.Header.Get(headerAPIVersion), ".")[0]

		if requestedAPIVersionValue == "" {
			err := fmt.Errorf("HTTP Status: (%v) - mandatory request header %v not set", http.StatusPreconditionFailed, headerAPIVersion)
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Header().Set(headerContentType, contentTypeTEXT)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		requestedAPIVersion := strings.Split(requestedAPIVersionValue, ".")[0]
		if supportedAPIVersion != requestedAPIVersion {
			err := fmt.Errorf("HTTP Stauts: (%v) - requested API version is %v but supported API version is %v", http.StatusPreconditionFailed, requestedAPIVersionValue, supportedAPIVersionValue)
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Header().Set(headerContentType, contentTypeTEXT)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	fmt.Fprint(w, "{\"catalog\": true}")
}
