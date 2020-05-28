package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	supportedAPIVersionValue string = "2.14"

	headerAPIVersion            string = "X-Broker-API-Version"
	headerAPIOrginatingIdentity string = "X-Broker-API-Originating-Identity"
	headerAPIRequestIdentity    string = "X-Broker-API-Request-Identity"
)

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
