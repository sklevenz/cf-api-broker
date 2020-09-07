package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sklevenz/cf-api-broker/data"
	"github.com/sklevenz/cf-api-broker/openapi"
)

const (
	supportedAPIVersionValue string = "2.14"

	headerAPIVersion            string = "X-Broker-API-Version"
	headerAPIOrginatingIdentity string = "X-Broker-API-Originating-Identity"
	headerAPIRequestIdentity    string = "X-Broker-API-Request-Identity"
)

type userIDType struct {
	UserID string `json:"user_id"`
}

type originatingIdentityType struct {
	Platform string     `json:"platform"`
	UserID   userIDType `json:"user_id_object"`
}

func handleOSBError(w http.ResponseWriter, code int, err openapi.Error) {
	output, _ := json.Marshal(err)

	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(code)
	w.Write(output)
}

func handleHTTPError(w http.ResponseWriter, code int, err error) {

	output, _ := json.Marshal(&openapi.Error{
		Error:       http.StatusText(code),
		Description: err.Error(),
	})

	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(code)
	w.Write(output)
}

func parseOriginatingIdentityHeader(value string) (*originatingIdentityType, error) {
	value = strings.TrimSpace(value)
	values := strings.Split(value, " ")

	originatingIdentity := &originatingIdentityType{}
	originatingIdentity.Platform = values[0]

	encoded, err := base64.StdEncoding.DecodeString(values[1])

	if err != nil {
		log.Printf("Error in Originating Identity Header, user_id not base64 encoded: %v", value)
		return nil, err
	}

	json.Unmarshal([]byte(encoded), &originatingIdentity.UserID)

	return originatingIdentity, nil
}

func originatingIdentityLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(headerAPIOrginatingIdentity)
		if headerValue != "" {
			originatingIdentity, err := parseOriginatingIdentityHeader(headerValue)
			if err == nil {
				log.Printf("Originating Identity: %v", originatingIdentity)
			}
		} else {
			log.Printf("Header %v not set", headerAPIOrginatingIdentity)
		}

		next.ServeHTTP(w, r)
	})
}

func requestIdentityLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(headerAPIRequestIdentity)
		if headerValue != "" {
			log.Printf("Request Identity: %v", headerValue)
		} else {
			log.Printf("Header %v not set", headerAPIRequestIdentity)
		}

		next.ServeHTTP(w, r)
	})
}

func apiVersionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supportedAPIVersion := strings.Split(supportedAPIVersionValue, ".")[0]
		requestedAPIVersionValue := strings.Split(r.Header.Get(headerAPIVersion), ".")[0]

		if requestedAPIVersionValue == "" {
			err := fmt.Errorf("HTTP Status: (%v) - mandatory request header %v not set", http.StatusPreconditionFailed, headerAPIVersion)
			log.Printf("Error: %v", err)
			handleHTTPError(w, http.StatusPreconditionFailed, err)
			return
		}

		requestedAPIVersion := strings.Split(requestedAPIVersionValue, ".")[0]
		if supportedAPIVersion != requestedAPIVersion {
			err := fmt.Errorf("HTTP Status: (%v) - requested API version is %v but supported API version is %v", http.StatusPreconditionFailed, r.Header.Get(headerAPIVersion), supportedAPIVersionValue)
			log.Printf("Error: %v", err)
			handleHTTPError(w, http.StatusPreconditionFailed, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func etagHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dat, err := data.NewCloudFoundryMetaData(configPath)

		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(headerETag, fmt.Sprintf("W/\"%v\"", dat.GetLastModifiedHash()))

		next.ServeHTTP(w, r)
	})
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	dat, err := data.NewCloudFoundryMetaData(configPath)

	if err != nil {
		handleHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	js, err := json.Marshal(buildCatalog(dat))
	if err != nil {
		handleHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set(headerContentType, contentTypeJSON)
	reader := bytes.NewReader(js)
	http.ServeContent(w, r, "xxx", dat.GetLastModified(), reader)
}

func buildCatalog(dat *data.CloudFoundriesType) *openapi.Catalog {
	catalog := openapi.Catalog{}
	var services []openapi.Service
	var service openapi.Service

	service.Id = "cf"
	service.Name = "cloudfoundry"
	service.Description = "cloudfoundry api endpoint service"
	service.Tags = append(service.Tags, "cf", "api", "cloudfoundry")
	service.Requires = []string{}
	service.Bindable = true
	service.InstancesRetrievable = true
	service.BindingsRetrievable = true
	service.AllowContextUpdates = true
	service.Metadata = map[string]interface{}{}
	service.DashboardClient = openapi.DashboardClient{}
	service.PlanUpdateable = true

	plans := []openapi.Plan{}

	for cfName, cf := range dat.CloudFoundries {

		plan := openapi.Plan{}
		plan.Id = cfName
		plan.Name = cfName
		plan.Description = "cloud foundry api instance"
		plan.Metadata = make(map[string]interface{})
		plan.Metadata["labels"] = cf.Labels
		plan.Free = true
		plan.Bindable = true
		plan.PlanUpdateable = true
		plan.Schemas = openapi.SchemasObject{}
		plan.MaximumPollingDuration = 10
		plan.MaintenanceInfo = openapi.MaintenanceInfo{}

		plans = append(plans, plan)
	}

	service.Plans = plans
	services = append(services, service)
	catalog.Services = services

	log.Printf("Catalog: %v", catalog)

	return &catalog
}
