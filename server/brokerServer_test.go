package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	broker := NewBrokerServer()

	routeHome := broker.Get("home")
	routeStatic := broker.Get("static")
	routeCatalog := broker.Get("catalog")

	assert.Equal(t, "home", routeHome.GetName())
	assert.Equal(t, "static", routeStatic.GetName())
	assert.Equal(t, "catalog", routeCatalog.GetName())
}

func TestHeaderHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, brokerAPIVersion, response.Header().Get("X-Broker-API-Version"))
}

func TestHomeHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer().ServeHTTP(response, request)

	assert.Contains(t, "Cloud Foundry API Broker", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	//	catalogHandler(response, request)
	NewBrokerServer().ServeHTTP(response, request)

	assert.Equal(t, "catalog", response.Body.String())
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
