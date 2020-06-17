package server

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sklevenz/cf-api-broker/openapi"
	"github.com/stretchr/testify/assert"
)

func TestNoApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestWrongApiVersionFormat(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "abc")
	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestWrongApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "1.2")
	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusPreconditionFailed, response.Result().StatusCode)
}

func TestCorrectApiVersion(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestRedirect(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, contentTypeHTML, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusMovedPermanently, response.Result().StatusCode)
}

func TestCatalogHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	request.Header.Set(headerAPIVersion, "2.2")
	NewRouter(staticDir, testConfigPath).ServeHTTP(response, request)

	assert.Equal(t, "{\"catalog\": true}", response.Body.String())
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestAPIOriginatingIdentity(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	request.Header.Set(headerAPIOrginatingIdentity, "cloudfoundry eyANCiAgInVzZXJfaWQiOiAiNjgzZWE3NDgtMzA5Mi00ZmY0LWI2NTYtMzljYWNjNGQ1MzYwIg0KfQ==")

	handler := originatingIdentityLogHandler(testHandler)
	handler.ServeHTTP(response, request)

	assert.Contains(t, buf.String(), "cloudfoundry")
	assert.Contains(t, buf.String(), "683ea748-3092-4ff4-b656-39cacc4d5360")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestAPIRequestIdentity(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	request.Header.Set(headerAPIRequestIdentity, "e26cee84-6b38-4456-b34e-d1a9f002c956")

	handler := requestIdentityLogHandler(testHandler)
	handler.ServeHTTP(response, request)

	assert.Contains(t, buf.String(), "e26cee84-6b38-4456-b34e-d1a9f002c956")
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
func TestHttpErrorHandler(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleHttpError(w, http.StatusInternalServerError, errors.New("blabla"))
	})

	handler := requestIdentityLogHandler(testHandler)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Contains(t, response.Body.String(), "blabla")
}

func TestOSBErrorHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/v2/catalog/", nil)
	response := httptest.NewRecorder()

	err := &openapi.Error{
		Error:            "AsyncRequired",
		Description:      "blabla",
		InstanceUsable:   true,
		UpdateRepeatable: true,
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleOSBError(w, http.StatusInternalServerError, *err)
	})

	handler := requestIdentityLogHandler(testHandler)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
	assert.Equal(t, contentTypeJSON, response.Header().Get(headerContentType))
	assert.Contains(t, response.Body.String(), "AsyncRequired")
	assert.Contains(t, response.Body.String(), "blabla")
	assert.Contains(t, response.Body.String(), "instance_usable")
	assert.Contains(t, response.Body.String(), "instance_usable")
}
