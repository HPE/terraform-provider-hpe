package morpheus

import (
	"fmt"
)

var (
	// JobsPath is the API endpoint for job executions
	JobExecutionsPath = "/api/job-executions"
)

// JobExecution structures for use in request and response payloads
type JobExecution struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Process struct {
		ID          int64  `json:"id"`
		AccountId   int64  `json:"accountId"`
		UniqueId    string `json:"uniqueId"`
		ProcessType struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"processType"`
		DisplayName   string  `json:"displayName"`
		Description   string  `json:"description"`
		SubType       string  `json:"subType"`
		SubId         string  `json:"subId"`
		ZoneId        int64   `json:"zoneId"`
		IntegrationId string  `json:"integrationId"`
		AppId         int64   `json:"appId"`
		InstanceId    int64   `json:"instanceId"`
		ContainerId   int64   `json:"containerId"`
		ServerId      int64   `json:"serverId"`
		ContainerName string  `json:"containerName"`
		Status        string  `json:"status"`
		Reason        string  `json:"reason"`
		Percent       float64 `json:"percent"`
		StatusEta     int64   `json:"statusEta"`
		Message       string  `json:"message"`
		Output        string  `json:"output"`
		Error         string  `json:"error"`
		StartDate     string  `json:"startDate"`
		EndDate       string  `json:"endDate"`
		Duration      int64   `json:"duration"`
		DateCreated   string  `json:"dateCreated"`
		LastUpdated   string  `json:"lastUpdated"`
		CreatedBy     struct {
			Username    string `json:"username"`
			DisplayName string `json:"displayName"`
		} `json:"createdBy"`
		UpdatedBy struct {
			Username    string `json:"username"`
			DisplayName string `json:"displayName"`
		} `json:"updatedBy"`
		Events []struct {
			ID          int64  `json:"id"`
			ProcessId   int64  `json:"processId"`
			AccountId   int64  `json:"accountId"`
			UniqueId    string `json:"uniqueId"`
			ProcessType struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"processType"`
			Description   string  `json:"description"`
			RefType       string  `json:"refType"`
			RefId         int64   `json:"refId"`
			SubType       string  `json:"subType"`
			SubId         string  `json:"subId"`
			ZoneId        int64   `json:"zoneId"`
			IntegrationId string  `json:"integrationId"`
			InstanceId    int64   `json:"instanceId"`
			ContainerId   int64   `json:"containerId"`
			ServerId      int64   `json:"serverId"`
			ContainerName string  `json:"containerName"`
			DisplayName   string  `json:"displayName"`
			Status        string  `json:"status"`
			Reason        string  `json:"reason"`
			Percent       float64 `json:"percent"`
			StatusEta     int64   `json:"statusEta"`
			Message       string  `json:"message"`
			Output        string  `json:"output"`
			Error         string  `json:"error"`
			StartDate     string  `json:"startDate"`
			EndDate       string  `json:"endDate"`
			Duration      int64   `json:"duration"`
			DateCreated   string  `json:"dateCreated"`
			LastUpdated   string  `json:"lastUpdated"`
			CreatedBy     struct {
				Username    string `json:"username"`
				DisplayName string `json:"displayName"`
			} `json:"createdBy"`
			UpdatedBy struct {
				Username    string `json:"username"`
				DisplayName string `json:"displayName"`
			} `json:"updatedBy"`
		} `json:"events"`
	} `json:"process"`
	Job struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"job"`
	Description   string `json:"description"`
	DateCreated   string `json:"dateCreated"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Duration      int64  `json:"duration"`
	ResultData    string `json:"resultData"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	CreatedBy     struct {
		Id          int64  `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	} `json:"createdBy"`
}

type JobExecutionEvent struct {
	ID          int64  `json:"id"`
	ProcessId   int64  `json:"processId"`
	AccountId   int64  `json:"accountId"`
	UniqueId    string `json:"uniqueId"`
	ProcessType struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"processType"`
	Description   string `json:"description"`
	RefType       string `json:"refType"`
	RefId         int64  `json:"refId"`
	SubType       string `json:"subType"`
	SubId         string `json:"subId"`
	ZoneId        int64  `json:"zoneId"`
	IntegrationId string `json:"integrationId"`
	InstanceId    int64  `json:"instanceId"`
	ContainerId   int64  `json:"containerId"`
	ServerId      int64  `json:"serverId"`
	ContainerName string `json:"containerName"`
	DisplayName   string `json:"displayName"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Percent       int64  `json:"percent"`
	StatusEta     int64  `json:"statusEta"`
	Message       string `json:"message"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Duration      int64  `json:"duration"`
	DateCreated   string `json:"dateCreated"`
	LastUpdated   string `json:"lastUpdated"`
	CreatedBy     struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	} `json:"createdBy"`
	UpdatedBy struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	} `json:"updatedBy"`
}

// ListJobExecutionsResult structure parses the list job executions response payload
type ListJobExecutionsResult struct {
	JobExecutions *[]JobExecution `json:"jobExecutions"`
	Meta          *MetaResult     `json:"meta"`
}

type GetJobExecutionResult struct {
	JobExecution *JobExecution `json:"jobExecution"`
}

type GetJobExecutionEventResult struct {
	ProcessEvent *JobExecutionEvent `json:"processEvent"`
}

// Client request methods
func (client *Client) ListJobExecutions(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        JobExecutionsPath,
		QueryParams: req.QueryParams,
		Result:      &ListJobExecutionsResult{},
	})
}

func (client *Client) GetJobExecution(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", JobExecutionsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetJobExecutionResult{},
	})
}

func (client *Client) GetJobExecutionEvent(id int64, eventId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/events/%d", JobExecutionsPath, id, eventId),
		QueryParams: req.QueryParams,
		Result:      &GetJobExecutionEventResult{},
	})
}
