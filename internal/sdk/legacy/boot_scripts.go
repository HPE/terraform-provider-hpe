package morpheus

import (
	"fmt"
)

var (
	// BootScriptsPath is the API endpoint for boot scripts
	BootScriptsPath = "/api/boot-scripts"
)

// BootScript structures for use in request and response payloads
type BootScript struct {
	ID      int64 `json:"id"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	FileName    string `json:"fileName"`
	Description string `json:"description"`
	Content     string `json:"content"`
	CreatedBy   struct {
		Username string `json:"username"`
	} `json:"createdBy"`
	Visibility string `json:"visibility"`
}

// ListBootScriptsResult structure parses the list bootScript response payload
type ListBootScriptsResult struct {
	BootScripts *[]BootScript `json:"bootScripts"`
	Meta        *MetaResult   `json:"meta"`
}

// GetBootScriptResult structure parses the get bootScript response payload
type GetBootScriptResult struct {
	BootScript *BootScript `json:"bootScript"`
}

// CreateBootScriptResult structure parses the create bootScript response payload
type CreateBootScriptResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	BootScript *BootScript       `json:"bootScript"`
}

// UpdateBootScriptResult structure parses the update bootScript response payload
type UpdateBootScriptResult struct {
	CreateBootScriptResult
}

// DeleteBootScriptResult structure parses the delete bootScript response payload
type DeleteBootScriptResult struct {
	DeleteResult
}

// ListBootScriptSets lists all bootScripts
func (client *Client) ListBootScripts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BootScriptsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBootScriptsResult{},
	})
}

// GetBootScriptSet gets an existing bootScript
func (client *Client) GetBootScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BootScriptsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBootScriptResult{},
	})
}

// CreateBootScriptSet creates a new bootScript
func (client *Client) CreateBootScript(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BootScriptsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBootScriptResult{},
	})
}

// UpdateBootScript updates an existing bootScript
func (client *Client) UpdateBootScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BootScriptsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBootScriptResult{},
	})
}

// DeleteBootScript deletes an existing bootScript
func (client *Client) DeleteBootScript(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BootScriptsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBootScriptResult{},
	})
}

// FindBootScriptByName gets an existing bootScript by name
func (client *Client) FindBootScriptByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListBootScripts(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListBootScriptsResult)
	bootScriptCount := len(*listResult.BootScripts)
	if bootScriptCount != 1 {
		return resp, fmt.Errorf("found %d Boot Scripts for %v", bootScriptCount, name)
	}
	firstRecord := (*listResult.BootScripts)[0]
	bootScriptID := firstRecord.ID
	return client.GetBootScript(bootScriptID, &Request{})
}
