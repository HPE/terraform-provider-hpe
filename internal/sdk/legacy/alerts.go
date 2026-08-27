package morpheus

import (
	"fmt"
	"time"
)

var (
	// AlertsPath is the API endpoint for alerts
	AlertsPath = "/api/monitoring/alerts"
)

// Alert structures for use in request and response payloads
type Alert struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	AllApps     bool          `json:"allApps"`
	AllChecks   bool          `json:"allChecks"`
	AllGroups   bool          `json:"allGroups"`
	Active      bool          `json:"active"`
	MinSeverity string        `json:"minSeverity"`
	MinDuration int64         `json:"minDuration"`
	DateCreated time.Time     `json:"dateCreated"`
	LastUpdated time.Time     `json:"lastUpdated"`
	Checks      []int64       `json:"checks"`
	CheckGroups []interface{} `json:"checkGroups"`
	Apps        []interface{} `json:"apps"`
	Contacts    []struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Method string `json:"method"`
		Notify bool   `json:"notify"`
		Close  bool   `json:"close"`
	} `json:"contacts"`
}

// ListAlertsResult structure parses the list alerts response payload
type ListAlertsResult struct {
	Alerts *[]Alert    `json:"alerts"`
	Meta   *MetaResult `json:"meta"`
}

// GetAlertResult structure parses the get alert response payload
type GetAlertResult struct {
	Alert *Alert `json:"alert"`
}

// CreateAlertResult structure parses the create alert response payload
type CreateAlertResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Alert   *Alert            `json:"alert"`
}

// UpdateAlertResult structure parses the update alert response payload
type UpdateAlertResult struct {
	CreateAlertResult
}

// DeleteAlertResult structure parses the delete alert response payload
type DeleteAlertResult struct {
	DeleteResult
}

// ListAlerts lists all alerts
func (client *Client) ListAlerts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        AlertsPath,
		QueryParams: req.QueryParams,
		Result:      &ListAlertsResult{},
	})
}

// GetAlert gets an existing alert
func (client *Client) GetAlert(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", AlertsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetAlertResult{},
	})
}

// CreateAlert creates a new alert
func (client *Client) CreateAlert(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        AlertsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateAlertResult{},
	})
}

// UpdateAlert updates an existing alert
func (client *Client) UpdateAlert(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", AlertsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateAlertResult{},
	})
}

// DeleteAlert deletes an existing alert
func (client *Client) DeleteAlert(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", AlertsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteAlertResult{},
	})
}

// FindAlertByName gets an existing alert by name
func (client *Client) FindAlertByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListAlerts(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListAlertsResult)
	alertCount := len(*listResult.Alerts)
	if alertCount != 1 {
		return resp, fmt.Errorf("found %d Alerts for %v", alertCount, name)
	}
	firstRecord := (*listResult.Alerts)[0]
	alertID := firstRecord.ID
	return client.GetAlert(alertID, &Request{})
}
