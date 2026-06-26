package morpheus

import (
	"fmt"
)

var (
	// OptionTypesPath is the API endpoint for option types
	OptionTypesPath = "/api/library/option-types"
)

// OptionType structures for use in request and response payloads
type OptionType struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Labels       []string `json:"labels"`
	Code         string   `json:"code"`
	Type         string   `json:"type"`
	FieldName    string   `json:"fieldName"`
	FieldLabel   string   `json:"fieldLabel"`
	PlaceHolder  string   `json:"placeHolder"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	ExportMeta   bool     `json:"exportMeta"`
	OptionSource string   `json:"optionSource"`
	OptionList   struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"optionList"`
	HelpBlock             string `json:"helpBlock"`
	Editable              bool   `json:"editable"`
	Creatable             bool   `json:"creatable"`
	DependsOnCode         string `json:"dependsOnCode"`
	VerifyPattern         string `json:"verifyPattern"`
	VisibleOnCode         string `json:"visibleOnCode"`
	RequireOnCode         string `json:"requireOnCode"`
	ContextualDefault     bool   `json:"contextualDefault"`
	DisplayValueOnDetails bool   `json:"displayValueOnDetails"`
	ShowOnCreate          bool   `json:"showOnCreate"`
	ShowOnEdit            bool   `json:"showOnEdit"`
	Config                struct {
		Rows                     string      `json:"rows"`
		BrowseEnabled            bool        `json:"browseEnabled"`
		ClearPlaceholderOnSelect bool        `json:"clearPlaceholderOnSelect"`
		MultiSelect              interface{} `json:"multiSelect"`
	} `json:"config"`
}

type ListOptionTypesResult struct {
	OptionTypes *[]OptionType `json:"optionTypes"`
	Meta        *MetaResult   `json:"meta"`
}

type GetOptionTypeResult struct {
	OptionType *OptionType `json:"optionType"`
}

type CreateOptionTypeResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	OptionType *OptionType       `json:"optionType"`
}

type UpdateOptionTypeResult struct {
	CreateOptionTypeResult
}

type DeleteOptionTypeResult struct {
	DeleteResult
}

// ListOptionTypes list all option types
func (client *Client) ListOptionTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        OptionTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListOptionTypesResult{},
	})
}

// GetOptionType gets an option type
func (client *Client) GetOptionType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", OptionTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetOptionTypeResult{},
	})
}

// CreateOptionType creates a new option type
func (client *Client) CreateOptionType(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        OptionTypesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateOptionTypeResult{},
	})
}

// UpdateOptionType updates an existing option type
func (client *Client) UpdateOptionType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", OptionTypesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateOptionTypeResult{},
	})
}

// DeleteOptionType deletes an existing option type
func (client *Client) DeleteOptionType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", OptionTypesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteOptionTypeResult{},
	})
}

// FindOptionTypeByName gets an existing option type by name
func (client *Client) FindOptionTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListOptionTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListOptionTypesResult)
	optionTypeCount := len(*listResult.OptionTypes)
	if optionTypeCount != 1 {
		return resp, fmt.Errorf("found %d OptionTypes for %v", optionTypeCount, name)
	}
	firstRecord := (*listResult.OptionTypes)[0]
	optionTypeID := firstRecord.ID
	return client.GetOptionType(optionTypeID, &Request{})
}
