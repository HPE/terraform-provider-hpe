// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Integration event models for the generic integration resource

// IntegrationEventRequest is the body used to create or update an outbound integration event.
type IntegrationEventRequest struct {
	Name                   string           `json:"name"`
	Entity                 string           `json:"entity"`
	EventType              string           `json:"eventType"`
	UseBaseNotifier        bool             `json:"useBaseNotifier"`
	Notifier               *NotifierRequest `json:"notifier,omitempty"`
	ThirdPartyEventType    string           `json:"thirdPartyEventType,omitempty"`
	Headers                []KeyValuePair   `json:"headers,omitempty"`
	EventPayload           string           `json:"eventPayload,omitempty"`
	EndPointURI            string           `json:"endPointURI,omitempty"`
	ResponseHeaders        []KeyValuePair   `json:"responseHeaders,omitempty"`
	ResourceGroupAllowed   bool             `json:"resourceGroupAllowed,omitempty"`
	CustomAttributeAllowed bool             `json:"customAttributeAllowed,omitempty"`
}

// IntegrationEventResponse is the API response for a created/fetched outbound integration event.
type IntegrationEventResponse struct {
	ID                     string           `json:"id"`
	Name                   string           `json:"name"`
	Entity                 string           `json:"entity"`
	EventType              string           `json:"eventType"`
	UseBaseNotifier        bool             `json:"useBaseNotifier"`
	Notifier               *NotifierRequest `json:"notifier,omitempty"`
	ThirdPartyEventType    string           `json:"thirdPartyEventType,omitempty"`
	Headers                []KeyValuePair   `json:"headers,omitempty"`
	EventPayload           string           `json:"eventPayload,omitempty"`
	EndPointURI            string           `json:"endPointURI,omitempty"`
	ResponseHeaders        []KeyValuePair   `json:"responseHeaders,omitempty"`
	ResourceGroupAllowed   bool             `json:"resourceGroupAllowed"`
	CustomAttributeAllowed bool             `json:"customAttributeAllowed"`
	Active                 bool             `json:"active"`
	EventLevel             string           `json:"eventLevel,omitempty"`
	ModifiedTime           string           `json:"modifiedTime,omitempty"`
	ModifiedBy             string           `json:"modifiedBy,omitempty"`
}
