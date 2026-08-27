package morpheus

import (
	"fmt"
	"time"
)

var (
	// ScriptTemplatesPath is the API endpoint for container scripts
	ScriptTemplatesPath = "/api/library/container-scripts"
)

// ScriptTemplate structures for use in request and response payloads
type ScriptTemplate struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Labels  []string `json:"labels"`
	Account struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Code          string      `json:"code"`
	Category      interface{} `json:"category"`
	SortOrder     int         `json:"sortOrder"`
	ScriptVersion string      `json:"scriptVersion"`
	Script        string      `json:"script"`
	Scriptservice interface{} `json:"scriptService"`
	Scriptmethod  interface{} `json:"scriptMethod"`
	ScriptType    string      `json:"scriptType"`
	ScriptPhase   string      `json:"scriptPhase"`
	RunAsUser     string      `json:"runAsUser"`
	Runaspassword interface{} `json:"runAsPassword"`
	SudoUser      bool        `json:"sudoUser"`
	FailOnError   bool        `json:"failOnError"`
	Datecreated   time.Time   `json:"dateCreated"`
	Lastupdated   time.Time   `json:"lastUpdated"`
}

type ListScriptTemplatesResult struct {
	ScriptTemplates *[]ScriptTemplate `json:"containerScripts"`
	Meta            *MetaResult       `json:"meta"`
}

type GetScriptTemplateResult struct {
	ScriptTemplate *ScriptTemplate `json:"containerScript"`
}

type CreateScriptTemplateResult struct {
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
	ScriptTemplate *ScriptTemplate   `json:"containerScript"`
}

type UpdateScriptTemplateResult struct {
	CreateScriptTemplateResult
}

type DeleteScriptTemplateResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListScriptTemplates(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ScriptTemplatesPath,
		QueryParams: req.QueryParams,
		Result:      &ListScriptTemplatesResult{},
	})
}

func (client *Client) GetScriptTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ScriptTemplatesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetScriptTemplateResult{},
	})
}

func (client *Client) CreateScriptTemplate(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ScriptTemplatesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateScriptTemplateResult{},
	})
}

func (client *Client) UpdateScriptTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ScriptTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateScriptTemplateResult{},
	})
}

func (client *Client) DeleteScriptTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ScriptTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteScriptTemplateResult{},
	})
}

// FindScriptTemplateByName gets an existing spec template by name
func (client *Client) FindScriptTemplateByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListScriptTemplates(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListScriptTemplatesResult)
	scriptTemplateCount := len(*listResult.ScriptTemplates)
	if scriptTemplateCount != 1 {
		return resp, fmt.Errorf("found %d Script Templates for %v", scriptTemplateCount, name)
	}
	firstRecord := (*listResult.ScriptTemplates)[0]
	scriptTemplateID := firstRecord.ID
	return client.GetScriptTemplate(scriptTemplateID, &Request{})
}
