package morpheus

import (
	"fmt"
)

var (
	// MonitoringAppsPath is the API endpoint for monitoring apps
	MonitoringAppsPath = "/api/monitoring/apps"
)

// MonitorApp structures for use in request and response payloads
type MonitoringApp struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Active      bool    `json:"active"`
	MinHappy    string  `json:"minHappy"`
	Severity    string  `json:"severity"`
	InUptime    bool    `json:"inUptime"`
	Checks      []int64 `json:"checks"`
	CheckGroups []int64 `json:"checkGroups"`
}

type ListMonitoringAppsResult struct {
	MonitoringApps *[]MonitoringApp `json:"monitorApps"`
	Meta           *MetaResult      `json:"meta"`
}

type GetMonitoringAppResult struct {
	MonitoringApp *MonitoringApp `json:"monitorApp"`
}

type CreateMonitoringAppResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	MonitoringApp *MonitoringApp    `json:"monitorApp"`
}

type UpdateMonitoringAppResult struct {
	CreateMonitoringAppResult
}

type DeleteMonitoringAppResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListMonitoringApps(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        MonitoringAppsPath,
		QueryParams: req.QueryParams,
		Result:      &ListMonitoringAppsResult{},
	})
}

func (client *Client) GetMonitoringApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", MonitoringAppsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetMonitoringAppResult{},
	})
}

func (client *Client) CreateMonitoringApp(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        MonitoringAppsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateMonitoringAppResult{},
	})
}

func (client *Client) UpdateMonitoringApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", MonitoringAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateMonitoringAppResult{},
	})
}

func (client *Client) DeleteMonitoringApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", MonitoringAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteMonitoringAppResult{},
	})
}

// helper functions
func (client *Client) FindMonitoringAppByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListMonitoringApps(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListMonitoringAppsResult)
	monitoringAppCount := len(*listResult.MonitoringApps)
	if monitoringAppCount != 1 {
		return resp, fmt.Errorf("found %d Monitoring Apps for %v", monitoringAppCount, name)
	}
	firstRecord := (*listResult.MonitoringApps)[0]
	monitoringAppID := firstRecord.ID
	return client.GetMonitoringApp(monitoringAppID, &Request{})
}
