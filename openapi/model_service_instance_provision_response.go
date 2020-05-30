/*
 * Open Service Broker API
 *
 * The Open Service Broker API defines an HTTP(S) interface between Platforms and Service Brokers.
 *
 * API version: master - might contain changes that are not yet released
 * Contact: open-service-broker-api@googlegroups.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ServiceInstanceProvisionResponse struct {

	DashboardUrl string `json:"dashboard_url,omitempty"`

	Metadata ServiceInstanceMetadata `json:"metadata,omitempty"`
}
