package morpheus

import (
	"fmt"
	"time"
)

var (
	// FileTemplatesPath is the API endpoint for container templates
	FileTemplatesPath = "/api/library/container-templates"
)

// FileTemplate structures for use in request and response payloads
type FileTemplate struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Name            string      `json:"name"`
	Labels          []string    `json:"labels"`
	FileName        string      `json:"fileName"`
	FilePath        string      `json:"filePath"`
	TemplateType    interface{} `json:"templateType"`
	TemplatePhase   string      `json:"templatePhase"`
	Template        string      `json:"template"`
	Category        interface{} `json:"category"`
	SettingCategory string      `json:"settingCategory"`
	SettingName     string      `json:"settingName"`
	AutoRun         bool        `json:"autoRun"`
	RunOnScale      bool        `json:"runOnScale"`
	RunOnDeploy     bool        `json:"runOnDeploy"`
	FileOwner       string      `json:"fileOwner"`
	FileGroup       interface{} `json:"fileGroup"`
	Permissions     interface{} `json:"permissions"`
	DateCreated     time.Time   `json:"dateCreated"`
	LastUpdated     time.Time   `json:"lastUpdated"`
}

type ListFileTemplatesResult struct {
	FileTemplates *[]FileTemplate `json:"containerTemplates"`
	Meta          *MetaResult     `json:"meta"`
}

type GetFileTemplateResult struct {
	FileTemplate *FileTemplate `json:"containerTemplate"`
}

type CreateFileTemplateResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	FileTemplate *FileTemplate     `json:"containerTemplate"`
}

type UpdateFileTemplateResult struct {
	CreateFileTemplateResult
}

type DeleteFileTemplateResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListFileTemplates(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        FileTemplatesPath,
		QueryParams: req.QueryParams,
		Result:      &ListFileTemplatesResult{},
	})
}

func (client *Client) GetFileTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", FileTemplatesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetFileTemplateResult{},
	})
}

func (client *Client) CreateFileTemplate(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        FileTemplatesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateFileTemplateResult{},
	})
}

func (client *Client) UpdateFileTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", FileTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateFileTemplateResult{},
	})
}

func (client *Client) DeleteFileTemplate(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", FileTemplatesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteFileTemplateResult{},
	})
}

// FindFileTemplateByName gets an existing spec template by name
func (client *Client) FindFileTemplateByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListFileTemplates(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListFileTemplatesResult)
	fileTemplateCount := len(*listResult.FileTemplates)
	if fileTemplateCount != 1 {
		return resp, fmt.Errorf("found %d File Templates for %v", fileTemplateCount, name)
	}
	firstRecord := (*listResult.FileTemplates)[0]
	fileTemplateID := firstRecord.ID
	return client.GetFileTemplate(fileTemplateID, &Request{})
}
