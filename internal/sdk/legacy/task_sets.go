package morpheus

import (
	"fmt"
)

var (
	// TaskSetsPath is the API endpoint for task sets
	TaskSetsPath = "/api/task-sets"
)

// TaskSet structures for use in request and response payloads
type TaskSet struct {
	ID                int64         `json:"id"`
	Name              string        `json:"name"`
	Labels            []string      `json:"labels"`
	Description       string        `json:"description"`
	Type              string        `json:"type"`
	AccountID         int64         `json:"accountId"`
	Visibility        string        `json:"visibility"`
	Platform          string        `json:"platform"`
	AllowCustomConfig bool          `json:"allowCustomConfig"`
	OptionTypes       []interface{} `json:"optionTypes"`
	Tasks             []int64       `json:"tasks"`
	TaskSetTasks      []TaskSetTask `json:"taskSetTasks"`
}

type TaskSetTask struct {
	ID        int64  `json:"id"`
	TaskPhase string `json:"taskPhase"`
	TaskOrder int64  `json:"taskOrder"`
	Task      Task   `json:"task"`
}

type ListTaskSetsResult struct {
	TaskSets *[]TaskSet  `json:"taskSets"`
	Meta     *MetaResult `json:"meta"`
}

type GetTaskSetResult struct {
	TaskSet *TaskSet `json:"taskSet"`
}

type CreateTaskSetResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	TaskSet *TaskSet          `json:"taskSet"`
}

type UpdateTaskSetResult struct {
	CreateTaskSetResult
}

type DeleteTaskSetResult struct {
	DeleteResult
}

// ListTaskSets lists all task sets
func (client *Client) ListTaskSets(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        TaskSetsPath,
		QueryParams: req.QueryParams,
		Result:      &ListTaskSetsResult{},
	})
}

// GetTaskSet gets an existing task set
func (client *Client) GetTaskSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", TaskSetsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetTaskSetResult{},
	})
}

// CreateTaskSet creates a new task set
func (client *Client) CreateTaskSet(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        TaskSetsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateTaskSetResult{},
	})
}

// UpdateTaskSet updates an existing task set
func (client *Client) UpdateTaskSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", TaskSetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateTaskSetResult{},
	})
}

// DeleteTaskSet deletes an existing task set
func (client *Client) DeleteTaskSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", TaskSetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteTaskSetResult{},
	})
}

// FindTaskSetByName gets an existing task set by name
func (client *Client) FindTaskSetByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListTaskSets(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListTaskSetsResult)
	taskSetCount := len(*listResult.TaskSets)
	if taskSetCount != 1 {
		return resp, fmt.Errorf("found %d TaskSets for %v", taskSetCount, name)
	}
	firstRecord := (*listResult.TaskSets)[0]
	taskSetID := firstRecord.ID
	return client.GetTaskSet(taskSetID, &Request{})
}
