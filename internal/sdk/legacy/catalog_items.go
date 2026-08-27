package morpheus

import (
	"fmt"
	"time"
)

var (
	// CatalogItemsPath is the API endpoint for catalog items
	CatalogItemsPath = "/api/catalog-item-types"
)

// CatalogItem structures for use in request and response payloads
type CatalogItem struct {
	ID             int64       `json:"id"`
	Name           string      `json:"name"`
	Code           string      `json:"code"`
	Description    string      `json:"description"`
	Category       string      `json:"category"`
	Labels         []string    `json:"labels"`
	Type           string      `json:"type"`
	Visibility     string      `json:"visibility"`
	LayoutCode     interface{} `json:"layoutCode"`
	FormType       string      `json:"formType"`
	RefType        string      `json:"refType"`
	RefID          interface{} `json:"refId"`
	Active         bool        `json:"active"`
	Enabled        bool        `json:"enabled"`
	Featured       bool        `json:"featured"`
	AllowQuantity  bool        `json:"allowQuantity"`
	IconPath       string      `json:"iconPath"`
	ImagePath      string      `json:"imagePath"`
	DarkImagePath  string      `json:"darkImagePath"`
	Context        string      `json:"context"`
	WorkflowConfig interface{} `json:"workflowConfig"`
	Content        string      `json:"content"`
	AppSpec        string      `json:"appSpec"`
	Blueprint      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	Workflow struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"workflow"`
	Config      interface{}   `json:"config"`
	FormConfig  interface{}   `json:"formConfig"`
	Form        Form          `json:"form"`
	OptionTypes []interface{} `json:"optionTypes"`
	CreatedBy   interface{}   `json:"createdBy"`
	Owner       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// ListCatalogItemsResult structure parses the list catalog items response payload
type ListCatalogItemsResult struct {
	CatalogItems *[]CatalogItem `json:"catalogItemTypes"`
	Meta         *MetaResult    `json:"meta"`
}

type GetCatalogItemResult struct {
	CatalogItem *CatalogItem `json:"catalogItemType"`
}

type CreateCatalogItemResult struct {
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
	CatalogItem *CatalogItem      `json:"catalogItemType"`
}

type UpdateCatalogItemResult struct {
	CreateCatalogItemResult
}

type DeleteCatalogItemResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListCatalogItems(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CatalogItemsPath,
		QueryParams: req.QueryParams,
		Result:      &ListCatalogItemsResult{},
	})
}

func (client *Client) GetCatalogItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CatalogItemsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCatalogItemResult{},
	})
}

// CreateCatalogItem creates a new catalog item
func (client *Client) CreateCatalogItem(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        CatalogItemsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCatalogItemResult{},
	})
}

// UpdateCatalogItem updates an existing catalog item
func (client *Client) UpdateCatalogItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", CatalogItemsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCatalogItemResult{},
	})
}

func (client *Client) UpdateCatalogItemLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("/api/catalog-item-types/%d/update-logo", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateInstanceTypeResult{},
	})
}

// DeleteCatalogItem deletes an existing catalog item
func (client *Client) DeleteCatalogItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", CatalogItemsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCatalogItemResult{},
	})
}

func (client *Client) FindCatalogItemByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListCatalogItems(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCatalogItemsResult)
	catalogItemCount := len(*listResult.CatalogItems)
	if catalogItemCount != 1 {
		return resp, fmt.Errorf("found %d Catalog Items for %v", catalogItemCount, name)
	}
	firstRecord := (*listResult.CatalogItems)[0]
	optionTypeID := firstRecord.ID
	return client.GetCatalogItem(optionTypeID, &Request{})
}
