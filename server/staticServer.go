package server

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	buildVersion string = "n/a"
	buildCommit  string = "n/a"
)

// SetBuildVersion get build information from calling application
func SetBuildVersion(version string, commit string) {
	buildVersion = version
	buildCommit = commit
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %v", r)
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{"buildVersion": buildVersion, "buildCommit": buildCommit})
}
