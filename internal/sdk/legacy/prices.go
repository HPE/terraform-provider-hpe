package morpheus

import (
	"fmt"
)

var (
	// PricesPath is the API endpoint for prices
	PricesPath = "/api/prices"
)

// Prices structures for use in request and response payloads
type Price struct {
	ID                  int64   `json:"id"`
	Name                string  `json:"name"`
	Code                string  `json:"code"`
	Active              bool    `json:"active"`
	PriceType           string  `json:"priceType"`
	PriceUnit           string  `json:"priceUnit"`
	AdditionalPriceUnit string  `json:"additionalPriceUnit"`
	Price               float64 `json:"price"`
	CustomPrice         float64 `json:"customPrice"`
	MarkupType          string  `json:"markupType"`
	Markup              float64 `json:"markup"`
	MarkupPercent       float64 `json:"markupPercent"`
	Cost                float64 `json:"cost"`
	Currency            string  `json:"currency"`
	IncurCharges        string  `json:"incurCharges"`
	Platform            string  `json:"platform"`
	Software            string  `json:"software"`
	RestartUsage        bool    `json:"restartUsage"`
	Volumetype          struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"volumeType"`
	Datastore struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"datastore"`
	CrossCloudApply bool `json:"crossCloudApply"`
	Zone            struct {
		ID int64 `json:"id"`
	} `json:"zone"`
	Zonepool struct {
		ID int64 `json:"id"`
	} `json:"zonePool"`
	Account interface{} `json:"account"`
}

// ListPricesResult structure parses the list prices response payload
type ListPricesResult struct {
	Prices *[]Price    `json:"prices"`
	Meta   *MetaResult `json:"meta"`
}

// GetPriceResult structure parses the get price response payload
type GetPriceResult struct {
	Price *Price `json:"price"`
}

// CreatePriceResult structure parses the create price response payload
type CreatePriceResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	ID      int64             `json:"id"`
}

// UpdatePriceResult structure parses the update price response payload
type UpdatePriceResult struct {
	CreatePriceResult
}

// DeletePriceResult structure parses the delete price response payload
type DeletePriceResult struct {
	DeleteResult
}

// ListPrices lists all prices
// https://apidocs.morpheusdata.com/#get-all-prices
func (client *Client) ListPrices(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PricesPath,
		QueryParams: req.QueryParams,
		Result:      &ListPricesResult{},
	})
}

// GetPrice gets an existing price
// https://apidocs.morpheusdata.com/#get-a-specific-price
func (client *Client) GetPrice(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PricesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPriceResult{},
	})
}

// CreatePrice creates a new price
// https://apidocs.morpheusdata.com/#create-a-price
func (client *Client) CreatePrice(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        PricesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreatePriceResult{},
	})
}

// UpdatePrice updates an existing price
// https://apidocs.morpheusdata.com/#update-a-price
func (client *Client) UpdatePrice(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PricesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePriceResult{},
	})
}

// DeletePrice deletes an existing price
// https://apidocs.morpheusdata.com/#deactivate-a-price
func (client *Client) DeletePrice(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/deactivate", PricesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePriceResult{},
	})
}

// FindPriceByName gets an existing price by name
func (client *Client) FindPriceByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPrices(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPricesResult)
	priceCount := len(*listResult.Prices)
	if priceCount != 1 {
		return resp, fmt.Errorf("found %d Prices for %v", priceCount, name)
	}
	firstRecord := (*listResult.Prices)[0]
	priceID := firstRecord.ID
	return client.GetPrice(priceID, &Request{})
}
