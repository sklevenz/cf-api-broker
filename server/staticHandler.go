package server

import (
	"encoding/json"
	"fmt"
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

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if config.Get().Server.AuthType == config.AuthTypeBasic {

			// handle basic auth
			u, p, ok := r.BasicAuth()
			if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
				unauthorised(w)
				return
			}

			if u != config.Get().Server.BasicAuth.UserName || p != config.Get().Server.BasicAuth.Password {
				unauthorised(w)
				return
			}
		} else {
			err := fmt.Errorf("Config error: unsupported AuthType: \"%v\"", config.Get().Server.AuthType)
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)
}
