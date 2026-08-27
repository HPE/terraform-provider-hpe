package morpheus

import (
	"fmt"
)

var (
	// PriceSetsPath is the API endpoint for priceSets
	PriceSetsPath = "/api/price-sets"
)

// PriceSets structures for use in request and response payloads
type PriceSet struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Code          string                 `json:"code"`
	Active        bool                   `json:"active"`
	PriceUnit     string                 `json:"priceUnit"`
	Type          string                 `json:"type"`
	RegionCode    string                 `json:"regionCode"`
	SystemCreated bool                   `json:"systemCreated"`
	Zone          map[string]interface{} `json:"zone"`
	ZonePool      map[string]interface{} `json:"zonePool"`
	Prices        []Price                `json:"prices"`
	RestartUsage  bool                   `json:"restartUsage"`
}

// ListPriceSetsResult structure parses the list priceSets response payload
type ListPriceSetsResult struct {
	PriceSets *[]PriceSet `json:"priceSets"`
	Meta      *MetaResult `json:"meta"`
}

// GetPriceSetResult structure parses the get priceSet response payload
type GetPriceSetResult struct {
	PriceSet *PriceSet `json:"priceSet"`
}

// CreatePriceSetResult structure parses the create priceSet response payload
type CreatePriceSetResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	ID      int64             `json:"id"`
}

// UpdatePriceSetResult structure parses the update priceSet response payload
type UpdatePriceSetResult struct {
	CreatePriceSetResult
}

// DeletePriceSetResult structure parses the delete priceSet response payload
type DeletePriceSetResult struct {
	DeleteResult
}

// ListPriceSetSets lists all priceSets
// https://apidocs.morpheusdata.com/#get-all-priceSets
func (client *Client) ListPriceSets(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PriceSetsPath,
		QueryParams: req.QueryParams,
		Result:      &ListPriceSetsResult{},
	})
}

// GetPriceSetSet gets an existing priceSet
// https://apidocs.morpheusdata.com/#get-a-specific-priceSet
func (client *Client) GetPriceSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PriceSetsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPriceSetResult{},
	})
}

// CreatePriceSetSet creates a new priceSetSet
// https://apidocs.morpheusdata.com/#create-a-priceSet
func (client *Client) CreatePriceSet(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        PriceSetsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreatePriceSetResult{},
	})
}

// UpdatePriceSet updates an existing priceSet
// https://apidocs.morpheusdata.com/#update-a-priceSet
func (client *Client) UpdatePriceSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PriceSetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePriceSetResult{},
	})
}

// DeletePriceSet deletes an existing priceSet
// https://apidocs.morpheusdata.com/#deactivate-a-priceSet
func (client *Client) DeletePriceSet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/deactivate", PriceSetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePriceSetResult{},
	})
}

// FindPriceSetByName gets an existing priceSet by name
func (client *Client) FindPriceSetByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPriceSets(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPriceSetsResult)
	priceSetCount := len(*listResult.PriceSets)
	if priceSetCount != 1 {
		return resp, fmt.Errorf("found %d PriceSets for %v", priceSetCount, name)
	}
	firstRecord := (*listResult.PriceSets)[0]
	priceSetID := firstRecord.ID
	return client.GetPriceSet(priceSetID, &Request{})
}
