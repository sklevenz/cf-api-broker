package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	staticDir string = "./../static"
)

func TestHomeHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeHTML, response.Header().Get(headerContentType))
	assert.Contains(t, response.Body.String(), "Cloud Foundry API - OSB Broker")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestStaticCSS(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/static/css/broker.css", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeCSS, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "CSS for Cloud Foundry API - OSB Broker")
}

func TestVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/version/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.JSONEq(t, `{"buildVersion":"n/a", "buildCommit":"n/a"}`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
func TestHealth(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/health/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.JSONEq(t, `{"ok":true}`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
