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

func TestHeaderHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, brokerAPIVersion, response.Header().Get("X-Broker-API-Version"))
}

func TestRedirect(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, http.StatusMovedPermanently, response.Result().StatusCode)
}

func TestHomeHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Contains(t, response.Body.String(), "Cloud Foundry API - OSB Broker")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestStaticCSS(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/static/css/broker.css", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "CSS for Cloud Foundry API - OSB Broker")
}

func TestHealth(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/health/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.JSONEq(t, `{"ok":true}`, response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	//	catalogHandler(response, request)
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, "catalog", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
