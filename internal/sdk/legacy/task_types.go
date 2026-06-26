package morpheus

import (
	"fmt"
)

var (
	// TaskTypesPath is the API endpoint for Task types
	TaskTypesPath = "/api/task-types"
)

// TaskType structures for use in request and response payloads
type TaskType struct {
	ID                   int64        `json:"id"`
	Code                 string       `json:"code"`
	Name                 string       `json:"name"`
	Category             string       `json:"category"`
	Description          string       `json:"description"`
	Scriptable           bool         `json:"scriptable"`
	Enabled              bool         `json:"enabled"`
	HasResults           bool         `json:"hasResults"`
	AllowExecuteLocal    bool         `json:"allowExecuteLocal"`
	AllowExecuteRemote   bool         `json:"allowExecuteRemote"`
	AllowExecuteResource bool         `json:"allowExecuteResource"`
	AllowLocalRepo       bool         `json:"allowLocalRepo"`
	AllowRemoteKeyAuth   bool         `json:"allowRemoteKeyAuth"`
	OptionTypes          []OptionType `json:"optionTypes"`
}

type ListTaskTypesResult struct {
	TaskTypes *[]TaskType `json:"taskTypes"`
	Meta      *MetaResult `json:"meta"`
}

type GetTaskTypeResult struct {
	TaskType *TaskType `json:"taskType"`
}

// Client request methods
func (client *Client) ListTaskTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        TaskTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListTaskTypesResult{},
	})
}

func (client *Client) GetTaskType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", TaskTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetTaskTypeResult{},
	})
}

func (client *Client) FindTaskTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListTaskTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListTaskTypesResult)
	TaskTypeCount := len(*listResult.TaskTypes)
	if TaskTypeCount != 1 {
		return resp, fmt.Errorf("found %d task types named %v", TaskTypeCount, name)
	}
	firstRecord := (*listResult.TaskTypes)[0]
	TaskTypeID := firstRecord.ID
	return client.GetTaskType(TaskTypeID, &Request{})
}
