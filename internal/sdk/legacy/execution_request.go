package morpheus

import (
	"fmt"
	"strconv"
)

var (
	// ExecutionRequestPath is the API endpoint for execution requests
	ExecutionRequestPath = "/api/execution-request"
)

// ExeuctionRequestResult structure parses the response from the server when running or retrieving an exeuction request
type ExecutionRequest struct {
	ID            int64    `json:"id"`
	UniqueID      string   `json:"uniqueId"`
	ContainerID   int64    `json:"containerId"`
	ServerID      int64    `json:"serverId"`
	InstanceID    int64    `json:"instanceId"`
	ResourceID    int64    `json:"resourceId"`
	AppID         int64    `json:"appId"`
	StdOut        string   `json:"stdOut"`
	StdErr        string   `json:"stdErr"`
	ExitCode      int64    `json:"exitCode"`
	Status        string   `json:"status"`
	ExpiresAt     string   `json:"expiresAt"`
	CreatedById   int64    `json:"createdById"`
	StatusMessage string   `json:"statusMessage"`
	ErrorMessage  string   `json:"errorMessage"`
	Config        struct{} `json:"config"`
	RawData       string   `json:"rawData"`
}

type ExecutionRequestResult struct {
	ExecutionRequest ExecutionRequest `json:"executionRequest"`
}

// Execution Request Methods

func (client *Client) GetExecutionRequest(uniqueID string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%s", ExecutionRequestPath, uniqueID),
		QueryParams: req.QueryParams,
		Result:      &ExecutionRequestResult{},
	})
}

func (client *Client) CreateExecutionRequest(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%s", ExecutionRequestPath, "execute"),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ExecutionRequestResult{},
	})
}

// helper functions

func (client *Client) ExecuteScriptOnInstance(instance Instance, script string) (stdOut string, stdErr string, err error) {

	resp, err := client.CreateExecutionRequest(&Request{
		QueryParams: map[string]string{
			"instanceId": strconv.Itoa(int(instance.ID)),
		},
		Body: map[string]interface{}{
			"script": script,
		},
	})

	if err != nil {
		fmt.Println("unable to run script on instance ", instance.ID, " due to ", err)
	}

	stdOut = resp.Result.(*ExecutionRequestResult).ExecutionRequest.StdOut
	stdErr = resp.Result.(*ExecutionRequestResult).ExecutionRequest.StdErr

	return
}
