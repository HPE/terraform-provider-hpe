package morpheus

import (
	"fmt"
	"time"
)

var (
	// JobsPath is the API endpoint for jobs
	JobsPath = "/api/jobs"
)

// Job structures for use in request and response payloads
type Job struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Labels     []string `json:"labels"`
	Enabled    bool     `json:"enabled"`
	TargetType string   `json:"targetType"`
	Task       struct {
		ID int64 `json:"id"`
	} `json:"task"`
	Workflow struct {
		ID int64 `json:"id"`
	} `json:"workflow"`
	Category  string `json:"category"`
	CreatedBy struct {
		DisplayName string `json:"displayName"`
		ID          int64  `json:"id"`
		Username    string `json:"username"`
	} `json:"createdBy"`
	CustomConfig  string      `json:"customConfig"`
	CustomOptions interface{} `json:"customOptions"`
	DateCreated   time.Time   `json:"dateCreated"`
	DateTime      interface{} `json:"dateTime"`
	Description   string      `json:"description"`
	JobSummary    string      `json:"jobSummary"`
	LastResult    string      `json:"lastResult"`
	LastRun       time.Time   `json:"lastRun"`
	LastUpdated   time.Time   `json:"lastUpdated"`
	Namespace     interface{} `json:"namespace"`
	ScheduleMode  string      `json:"scheduleMode"`
	Status        string      `json:"status"`
	Targets       []struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		RefId      int64  `json:"refId"`
		TargetType string `json:"targetType"`
	} `json:"targets"`
	Type struct {
		Code string `json:"code"`
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	SecurityProfile string `json:"securityProfile"`
	ScanPath        string `json:"scanPath"`
}

// ListJobsResult structure parses the list jobs response payload
type ListJobsResult struct {
	Jobs *[]Job      `json:"jobs"`
	Meta *MetaResult `json:"meta"`
}

type GetJobResult struct {
	Job *Job `json:"job"`
}

type CreateJobResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Job     *Job              `json:"job"`
}

type UpdateJobResult struct {
	CreateJobResult
}

type DeleteJobResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListJobs(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        JobsPath,
		QueryParams: req.QueryParams,
		Result:      &ListJobsResult{},
	})
}

func (client *Client) GetJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", JobsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetJobResult{},
	})
}

func (client *Client) CreateJob(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        JobsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateJobResult{},
	})
}

func (client *Client) UpdateJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", JobsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateJobResult{},
	})
}

func (client *Client) DeleteJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", JobsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteJobResult{},
	})
}

// helper functions

func (client *Client) FindJobByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListJobs(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListJobsResult)
	jobCount := len(*listResult.Jobs)
	if jobCount != 1 {
		return resp, fmt.Errorf("found %d Jobs named %v", jobCount, name)
	}
	firstRecord := (*listResult.Jobs)[0]
	jobID := firstRecord.ID
	return client.GetJob(jobID, &Request{})
}
