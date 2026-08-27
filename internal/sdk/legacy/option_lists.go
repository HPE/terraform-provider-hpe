package morpheus

import (
	"fmt"
)

var (
	// OptionListsPath is the API endpoint for option lists
	OptionListsPath = "/api/library/option-type-lists"
)

// OptionLists structures for use in request and response payloads
type OptionList struct {
	ID                  int64    `json:"id"`
	Name                string   `json:"name"`
	Labels              []string `json:"labels"`
	Description         string   `json:"description"`
	Type                string   `json:"type"`
	SourceURL           string   `json:"sourceUrl"`
	Visibility          string   `json:"visibility"`
	SourceMethod        string   `json:"sourceMethod"`
	APIType             string   `json:"apiType,omitempty"`
	IgnoreSSLErrors     bool     `json:"ignoreSSLErrors"`
	RealTime            bool     `json:"realTime"`
	InitialDataset      string   `json:"initialDataset"`
	TranslationScript   string   `json:"translationScript"`
	RequestScript       string   `json:"requestScript"`
	ServiceUsername     string   `json:"serviceUsername"`
	ServicePassword     string   `json:"servicePassword"`
	ServicePasswordHash string   `json:"servicePasswordHash"`
	Config              struct {
		SourceHeaders []SourceHeader `json:"sourceHeaders"`
	} `json:"config"`
	Credential struct {
		Type string `json:"type"`
	} `json:"credential"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
}

type SourceHeader struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Masked bool   `json:"masked"`
}

type ListOptionListsResult struct {
	OptionLists *[]OptionList `json:"optionTypeLists"`
	Meta        *MetaResult   `json:"meta"`
}

type ListOptionListItemsResult struct {
	OptionListItems *[]struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"listItems"`
}

type GetOptionListResult struct {
	OptionList *OptionList `json:"optionTypeList"`
}

type CreateOptionListResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	OptionList *OptionList       `json:"optionTypeList"`
}

type UpdateOptionListResult struct {
	CreateOptionListResult
}

type DeleteOptionListResult struct {
	DeleteResult
}

// ListOptionListItems lists all option list items
func (client *Client) ListOptionListItems(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/items", OptionListsPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListOptionListsResult{},
	})
}

// ListOptionLists lists all option lists
func (client *Client) ListOptionLists(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        OptionListsPath,
		QueryParams: req.QueryParams,
		Result:      &ListOptionListsResult{},
	})
}

// GetOptionList gets an option list
func (client *Client) GetOptionList(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", OptionListsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetOptionListResult{},
	})
}

// CreateOptionList creates a new option list
func (client *Client) CreateOptionList(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        OptionListsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateOptionListResult{},
	})
}

// UpdateOptionList updates an existing option list
func (client *Client) UpdateOptionList(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", OptionListsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateOptionListResult{},
	})
}

// DeleteOptionList deletes an existing option list
func (client *Client) DeleteOptionList(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", OptionListsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteOptionListResult{},
	})
}

// FindOptionListByName gets an existing option list by name
func (client *Client) FindOptionListByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListOptionLists(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListOptionListsResult)
	optionListCount := len(*listResult.OptionLists)
	if optionListCount != 1 {
		return resp, fmt.Errorf("found %d OptionLists for %v", optionListCount, name)
	}
	firstRecord := (*listResult.OptionLists)[0]
	optionListID := firstRecord.ID
	return client.GetOptionList(optionListID, &Request{})
}
