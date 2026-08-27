package morpheus

import (
	"fmt"
)

var (
	// PreseedScriptsPath is the API endpoint for preseed scripts
	PreseedScriptsPath = "/api/preseed-scripts"
)

// PreseedScript structures for use in request and response payloads
type PreseedScript struct {
	ID      int64 `json:"id"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	FileName    string      `json:"fileName"`
	Description interface{} `json:"description"`
	Content     string      `json:"content"`
	CreatedBy   struct {
		Username string `json:"username"`
	} `json:"createdBy"`
}

// ListPreseedScriptsResult structure parses the list preseedScript response payload
type ListPreseedScriptsResult struct {
	PreseedScripts *[]PreseedScript `json:"preseedScripts"`
	Meta           *MetaResult      `json:"meta"`
}

// GetPreseedScriptResult structure parses the get preseedScript response payload
type GetPreseedScriptResult struct {
	PreseedScript *PreseedScript `json:"preseedScript"`
}

// CreatePreseedScriptResult structure parses the create preseedScript response payload
type CreatePreseedScriptResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	PreseedScript *PreseedScript    `json:"preseedScript"`
}

// UpdatePreseedScriptResult structure parses the update preseedScript response payload
type UpdatePreseedScriptResult struct {
	CreatePreseedScriptResult
}

// DeletePreseedScriptResult structure parses the delete preseedScript response payload
type DeletePreseedScriptResult struct {
	DeleteResult
}

// ListPreseedScriptSets lists all preseedScripts
func (client *Client) ListPreseedScripts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PreseedScriptsPath,
		QueryParams: req.QueryParams,
		Result:      &ListPreseedScriptsResult{},
	})
}

// GetPreseedScriptSet gets an existing preseedScript
func (client *Client) GetPreseedScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PreseedScriptsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPreseedScriptResult{},
	})
}

// CreatePreseedScriptSet creates a new preseedScript
func (client *Client) CreatePreseedScript(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        PreseedScriptsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreatePreseedScriptResult{},
	})
}

// UpdatePreseedScript updates an existing preseedScript
func (client *Client) UpdatePreseedScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PreseedScriptsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePreseedScriptResult{},
	})
}

// DeletePreseedScript deletes an existing preseedScript
func (client *Client) DeletePreseedScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", PreseedScriptsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePreseedScriptResult{},
	})
}

// FindPreseedScriptByName gets an existing preseedScript by name
func (client *Client) FindPreseedScriptByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPreseedScripts(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPreseedScriptsResult)
	preseedScriptCount := len(*listResult.PreseedScripts)
	if preseedScriptCount != 1 {
		return resp, fmt.Errorf("found %d Preseed Scripts for %v", preseedScriptCount, name)
	}
	firstRecord := (*listResult.PreseedScripts)[0]
	priceSetID := firstRecord.ID
	return client.GetPreseedScript(priceSetID, &Request{})
}
