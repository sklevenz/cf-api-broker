package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNoApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeTEXT, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestWrongApiVersionFormat(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "abc")
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeTEXT, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestWrongApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "1.2")
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeTEXT, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestCorrectApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestRedirect(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, contentTypeHTML, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusMovedPermanently, response.Result().StatusCode)
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewBrokerServer(staticDir).ServeHTTP(response, request)

	assert.Equal(t, "{\"catalog\": true}", response.Body.String())
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

