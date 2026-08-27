package morpheus

import (
	_ "fmt"
)

var (
	SetupPath = "/api/setup"
)

type SetupCheckResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"msg"`
	BuildVersion string `json:"buildVersion"`
	SetupNeeded  bool   `json:"setupNeeded"`
}

type SetupInitResult struct {
	StandardResult
}

type SetupInitPayload struct {
	// Request
	AccountName string `json:"accountName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (client *Client) SetupCheck(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        (SetupPath + "/check"),
		QueryParams: req.QueryParams,
		Result:      &SetupCheckResult{},
	})
}

func (client *Client) SetupInit(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        (SetupPath + "/init"),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		// Body: payload.(map[string]interface {}),
		Result: &SetupInitResult{},
	})
}
