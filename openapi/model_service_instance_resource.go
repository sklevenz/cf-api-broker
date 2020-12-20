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

type ServiceInstanceResource struct {
	ServiceId string `json:"service_id,omitempty"`

	PlanId string `json:"plan_id,omitempty"`

	DashboardUrl string `json:"dashboard_url,omitempty"`

	Parameters map[string]interface{} `json:"parameters,omitempty"`

	MaintenanceInfo MaintenanceInfo `json:"maintenance_info,omitempty"`

	Metadata ServiceInstanceMetadata `json:"metadata,omitempty"`
}
