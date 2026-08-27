package morpheus

import (
	"fmt"
	"time"
)

var (
	// PluginsPath is the API endpoint for plugins
	PluginsPath = "/api/plugins"
)

// Plugin structures for use in request and response payloads
type Plugin struct {
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Version               string `json:"version"`
	RefType               string `json:"refType"`
	Enabled               bool   `json:"enabled"`
	Author                string `json:"author"`
	WebsiteUrl            string `json:"websiteUrl"`
	SourceCodeLocationUrl string `json:"sourceCodeLocationUrl"`
	IssueTrackerUrl       string `json:"issueTrackerUrl"`
	Valid                 bool   `json:"valid"`
	Status                string `json:"status"`
	StatusMessage         string `json:"statusMessage"`
	Providers             []struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"providers"`
	Config      interface{}   `json:"config"`
	OptionTypes []interface{} `json:"optionTypes"`
	DateCreated time.Time     `json:"dateCreated"`
	LastUpdated time.Time     `json:"lastUpdated"`
}

// ListPluginsResult structure parses the list plugins response payload
type ListPluginsResult struct {
	Plugins *[]Plugin   `json:"plugins"`
	Meta    *MetaResult `json:"meta"`
}

type GetPluginResult struct {
	Plugin *Plugin `json:"plugin"`
}

type CreatePluginResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Plugin  *Plugin           `json:"plugin"`
}

type UpdatePluginResult struct {
	CreatePluginResult
}

type DeletePluginResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListPlugins(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PluginsPath,
		QueryParams: req.QueryParams,
		Result:      &ListPluginsResult{},
	})
}

func (client *Client) GetPlugin(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PluginsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPluginResult{},
	})
}

// UploadPlugin uploads a new plugin
func (client *Client) UploadPlugin(filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("%s/upload", PluginsPath),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &CreateCatalogItemResult{},
	})
}

// UpdatePlugin updates an existing plugin
func (client *Client) UpdatePlugin(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PluginsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePluginResult{},
	})
}

// DeletePlugin deletes an existing plugin
func (client *Client) DeletePlugin(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", PluginsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePluginResult{},
	})
}

func (client *Client) FindPluginByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPlugins(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPluginsResult)
	pluginItemCount := len(*listResult.Plugins)
	if pluginItemCount != 1 {
		return resp, fmt.Errorf("found %d Plugins for %v", pluginItemCount, name)
	}
	firstRecord := (*listResult.Plugins)[0]
	pluginID := firstRecord.ID
	return client.GetPlugin(pluginID, &Request{})
}
