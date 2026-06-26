package morpheus

import (
	"fmt"
)

var (
	// ExecuteSchedulesPath is the API endpoint for execute schedules
	ExecuteSchedulesPath = "/api/execute-schedules"
)

// ExecuteSchedule structures for use in request and response payloads
type ExecuteSchedule struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Enabled          bool   `json:"enabled"`
	ScheduleType     string `json:"scheduleType"`
	ScheduleTimeZone string `json:"scheduleTimezone"`
	Desription       string `json:"description"`
	Cron             string `json:"cron"`
	DateCreated      string `json:"dateCreated"`
	LastUpdated      string `json:"lastUpdated"`
}

// ListExecuteSchedulesResult structure parses the list execute schedules response payload
type ListExecuteSchedulesResult struct {
	ExecuteSchedules *[]ExecuteSchedule `json:"schedules"`
	Meta             *MetaResult        `json:"meta"`
}

type GetExecuteScheduleResult struct {
	ExecuteSchedule *ExecuteSchedule `json:"schedule"`
}

type CreateExecuteScheduleResult struct {
	Success         bool              `json:"success"`
	Message         string            `json:"msg"`
	Errors          map[string]string `json:"errors"`
	ExecuteSchedule *ExecuteSchedule  `json:"schedule"`
}

type UpdateExecuteScheduleResult struct {
	CreateExecuteScheduleResult
}

type DeleteExecuteScheduleResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListExecuteSchedules(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ExecuteSchedulesPath,
		QueryParams: req.QueryParams,
		Result:      &ListExecuteSchedulesResult{},
	})
}

func (client *Client) GetExecuteSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ExecuteSchedulesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetExecuteScheduleResult{},
	})
}

func (client *Client) CreateExecuteSchedule(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ExecuteSchedulesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateExecuteScheduleResult{},
	})
}

// UpdateExecuteSchedule updates an existing execution schedule
func (client *Client) UpdateExecuteSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ExecuteSchedulesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateExecuteScheduleResult{},
	})
}

// DeleteExecuteSchedule deletes an existing execution schedule
func (client *Client) DeleteExecuteSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ExecuteSchedulesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteExecuteScheduleResult{},
	})
}

func (client *Client) FindExecuteScheduleByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListExecuteSchedules(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListExecuteSchedulesResult)
	executeScheduleCount := len(*listResult.ExecuteSchedules)
	if executeScheduleCount != 1 {
		return resp, fmt.Errorf("found %d execution schedules named %v", executeScheduleCount, name)
	}
	firstRecord := (*listResult.ExecuteSchedules)[0]
	executeScheduleID := firstRecord.ID
	return client.GetExecuteSchedule(executeScheduleID, &Request{})
}
