package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sklevenz/cf-api-broker/config"
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerContentType, contentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{"buildVersion": buildVersion, "buildCommit": buildCommit})
}

func basicAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg, err := config.New(configPath)

		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		// handle basic auth
		u, p, ok := r.BasicAuth()
		if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
			unauthorised(w)
			return
		}

		if u != cfg.Server.BasicAuth.UserName || p != cfg.Server.BasicAuth.Passowrd {
			unauthorised(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)
}
