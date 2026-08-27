package morpheus

import (
	"fmt"
	"time"
)

var (
	// ChecksPath is the API endpoint for check groups
	ChecksPath = "/api/monitoring/checks"
)

// Check structures for use in request and response payloads
type Check struct {
	ID      int64 `json:"id"`
	Account struct {
		ID int64 `json:"id"`
	} `json:"account"`
	Active        bool    `json:"active"`
	APIKey        string  `json:"apiKey"`
	Availability  float64 `json:"availability"`
	CheckAgent    string  `json:"checkAgent"`
	CheckInterval int64   `json:"checkInterval"`
	CheckSpec     string  `json:"checkSpec"`
	CheckType     struct {
		ID         int64  `json:"id"`
		Code       string `json:"code"`
		Name       string `json:"name"`
		MetricName string `json:"metricName"`
	} `json:"checkType"`
	Container struct {
		ID int64 `json:"id"`
	} `json:"container"`
	Config         interface{} `json:"config"`
	CreateIncident bool        `json:"createIncident"`
	Muted          bool        `json:"muted"`
	CreatedBy      struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"createdBy"`
	DateCreated     time.Time `json:"dateCreated"`
	Description     string    `json:"description"`
	EndDate         time.Time `json:"endDate"`
	Health          int64     `json:"health"`
	InUptime        bool      `json:"inUptime"`
	LastBoxStats    string    `json:"lastBoxStats"`
	LastCheckStatus string    `json:"lastCheckStatus"`
	LastError       string    `json:"lastError"`
	LastErrorDate   time.Time `json:"lastErrorDate"`
	LastMessage     string    `json:"lastMessage"`
	LastMetric      string    `json:"lastMetric"`
	LastRunDate     time.Time `json:"lastRunDate"`
	LastStats       string    `json:"lastStats"`
	LastSuccessDate time.Time `json:"lastSuccessDate"`
	LastTimer       int64     `json:"lastTimer"`
	LastUpdated     time.Time `json:"lastUpdated"`
	LastWarningDate time.Time `json:"lastWarningDate"`
	Name            string    `json:"name"`
	NextRunDate     time.Time `json:"nextRunDate"`
	OutageTime      int64     `json:"outageTime"`
	Severity        string    `json:"severity"`
	StartDate       time.Time `json:"startDate"`
}

// ListChecksResult structure parses the list check response payload
type ListChecksResult struct {
	Checks *[]Check    `json:"checks"`
	Meta   *MetaResult `json:"meta"`
}

type GetCheckResult struct {
	Check *Check `json:"check"`
}

type CreateCheckResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Check   *Check            `json:"check"`
}

type UpdateCheckResult struct {
	CreateCheckResult
}

type DeleteCheckResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListChecks(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ChecksPath,
		QueryParams: req.QueryParams,
		Result:      &ListChecksResult{},
	})
}

func (client *Client) GetCheck(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ChecksPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCheckResult{},
	})
}

// CreateCheck creates a new check
func (client *Client) CreateCheck(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ChecksPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCheckResult{},
	})
}

// UpdateCheck updates an existing check
func (client *Client) UpdateCheck(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ChecksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCheckResult{},
	})
}

// DeleteCheck deletes an existing check
func (client *Client) DeleteCheck(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ChecksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCheckResult{},
	})
}

func (client *Client) FindCheckByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListChecks(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListChecksResult)
	checkCount := len(*listResult.Checks)
	if checkCount != 1 {
		return resp, fmt.Errorf("found %d Checks for %v", checkCount, name)
	}
	firstRecord := (*listResult.Checks)[0]
	checkID := firstRecord.ID
	return client.GetCheck(checkID, &Request{})
}
