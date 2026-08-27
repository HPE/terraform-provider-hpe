package morpheus

import (
	"fmt"
	"time"
)

var (
	// CheckGroupsPath is the API endpoint for check groups
	CheckGroupsPath = "/api/monitoring/groups"
)

// CheckGroup structures for use in request and response payloads
type CheckGroup struct {
	ID      int64 `json:"id"`
	Account struct {
		ID int64 `json:"id"`
	} `json:"account"`
	Instance struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"instance"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	InUptime        bool      `json:"inUptime"`
	LastCheckStatus string    `json:"lastCheckStatus"`
	LastWarningDate time.Time `json:"lastWarningDate"`
	LastErrorDate   time.Time `json:"lastErrorDate"`
	LastSuccessDate time.Time `json:"lastSuccessDate"`
	LastRunDate     time.Time `json:"lastRunDate"`
	LastError       string    `json:"lastError"`
	OutageTime      int64     `json:"outageTime"`
	LastTimer       int64     `json:"lastTimer"`
	Health          int64     `json:"health"`
	History         string    `json:"history"`
	MinHappy        int64     `json:"minHappy"`
	LastMetric      string    `json:"lastMetric"`
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
	CheckType    struct {
		ID         int64  `json:"id"`
		Code       string `json:"code"`
		Name       string `json:"name"`
		MetricName string `json:"metricName"`
	} `json:"checkType"`
	Checks []int64 `json:"checks"`
}

// ListCheckGroupsResult structure parses the list check groups response payload
type ListCheckGroupsResult struct {
	CheckGroups *[]CheckGroup `json:"checkGroups"`
	Meta        *MetaResult   `json:"meta"`
}

type GetCheckGroupResult struct {
	CheckGroup *CheckGroup `json:"checkGroup"`
}

type CreateCheckGroupResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	CheckGroup *CheckGroup       `json:"checkGroup"`
}

type UpdateCheckGroupResult struct {
	CreateCheckGroupResult
}

type DeleteCheckGroupResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListCheckGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CheckGroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListCheckGroupsResult{},
	})
}

func (client *Client) GetCheckGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CheckGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCheckGroupResult{},
	})
}

// CreateCheckGroup creates a new check group
func (client *Client) CreateCheckGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        CheckGroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCheckGroupResult{},
	})
}

// UpdateCheckGroup updates an existing check group
func (client *Client) UpdateCheckGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", CheckGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCheckGroupResult{},
	})
}

// DeleteCheckGroup deletes an existing check group
func (client *Client) DeleteCheckGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", CheckGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCheckGroupResult{},
	})
}

func (client *Client) FindCheckGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListCheckGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCheckGroupsResult)
	checkGroupCount := len(*listResult.CheckGroups)
	if checkGroupCount != 1 {
		return resp, fmt.Errorf("found %d Check Groups for %v", checkGroupCount, name)
	}
	firstRecord := (*listResult.CheckGroups)[0]
	checkGroupID := firstRecord.ID
	return client.GetCheckGroup(checkGroupID, &Request{})
}
