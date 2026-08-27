package morpheus

import (
	"fmt"
	"time"
)

var (
	// CheckAppsPath is the API endpoint for check apps
	CheckAppsPath = "/api/monitoring/apps"
)

// CheckApp structures for use in request and response payloads
type CheckApp struct {
	ID      int64 `json:"id"`
	Account struct {
		ID int64 `json:"id"`
	} `json:"account"`
	Active bool `json:"active"`
	App    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"app"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	InUptime        bool      `json:"inUptime"`
	LastCheckStatus string    `json:"lastCheckStatus"`
	LastWarningDate time.Time `json:"lastWarningDate"`
	LastErrorDate   time.Time `json:"lastErrorDate"`
	LastSuccessDate time.Time `json:"lastSuccessDate"`
	LastRunDate     time.Time `json:"lastRunDate"`
	LastError       string    `json:"lastError"`
	LastTimer       int64     `json:"lastTimer"`
	Health          int64     `json:"health"`
	History         string    `json:"history"`
	Severity        string    `json:"severity"`
	CreateIncident  bool      `json:"createIncident"`
	Muted           bool      `json:"muted"`
	CreatedBy       struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"createdBy"`
	DateCreated  time.Time `json:"dateCreated"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Availability float64   `json:"availability"`
	Checks       []int64   `json:"checks"`
	CheckGroups  []int64   `json:"checkGroups"`
}

// ListCheckAppsResult structure parses the list check apps response payload
type ListCheckAppsResult struct {
	CheckApps *[]CheckApp `json:"monitorApps"`
	Meta      *MetaResult `json:"meta"`
}

type GetCheckAppResult struct {
	CheckApp *CheckApp `json:"monitorApp"`
}

type CreateCheckAppResult struct {
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
	CheckApp *CheckApp         `json:"monitorApp"`
}

type UpdateCheckAppResult struct {
	CreateCheckAppResult
}

type DeleteCheckAppResult struct {
	DeleteResult
}

// Client request methods

// ListCheckApps list all check apps
func (client *Client) ListCheckApps(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CheckAppsPath,
		QueryParams: req.QueryParams,
		Result:      &ListCheckAppsResult{},
	})
}

// GetCheckApp gets a check app
func (client *Client) GetCheckApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CheckAppsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCheckAppResult{},
	})
}

// CreateCheckApp creates a new check app
func (client *Client) CreateCheckApp(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        CheckAppsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCheckGroupResult{},
	})
}

// UpdateCheckApp updates an existing check app
func (client *Client) UpdateCheckApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", CheckAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCheckAppResult{},
	})
}

// DeleteCheckApp deletes an existing check app
func (client *Client) DeleteCheckApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", CheckAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCheckAppResult{},
	})
}

// FindCheckAppByName gets an existing check app by name
func (client *Client) FindCheckAppByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListCheckApps(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCheckAppsResult)
	checkAppCount := len(*listResult.CheckApps)
	if checkAppCount != 1 {
		return resp, fmt.Errorf("found %d Check Apps for %v", checkAppCount, name)
	}
	firstRecord := (*listResult.CheckApps)[0]
	checkGroupID := firstRecord.ID
	return client.GetCheckApp(checkGroupID, &Request{})
}
