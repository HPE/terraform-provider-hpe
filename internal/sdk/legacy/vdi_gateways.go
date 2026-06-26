package morpheus

import (
	"fmt"
	"time"
)

var (
	// VDIGatewaysPath is the API endpoint for vdi gateways
	VDIGatewaysPath = "/api/vdi-gateways"
)

// VDIGateway structures for use in request and response payloads
type VDIGateway struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GatewayURL  string    `json:"gatewayUrl"`
	ApiKey      string    `json:"apiKey"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ListVDIGatewaysResult struct {
	VDIGateways *[]VDIGateway `json:"vdiGateways"`
	Meta        *MetaResult   `json:"meta"`
}

type GetVDIGatewayResult struct {
	VDIGateway *VDIGateway `json:"vdiGateway"`
}

type CreateVDIGatewayResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	VDIGateway *VDIGateway       `json:"vdiGateway"`
}

type UpdateVDIGatewayResult struct {
	CreateVDIGatewayResult
}

type DeleteVDIGatewayResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListVDIGateways(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        VDIGatewaysPath,
		QueryParams: req.QueryParams,
		Result:      &ListVDIGatewaysResult{},
	})
}

func (client *Client) GetVDIGateway(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", VDIGatewaysPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetVDIGatewayResult{},
	})
}

func (client *Client) CreateVDIGateway(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        VDIGatewaysPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateVDIGatewayResult{},
	})
}

func (client *Client) UpdateVDIGateway(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", VDIGatewaysPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateVDIGatewayResult{},
	})
}

func (client *Client) DeleteVDIGateway(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", VDIGatewaysPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteVDIGatewayResult{},
	})
}

// FindVDIGatewayByName gets an existing vdi gateway by name
func (client *Client) FindVDIGatewayByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListVDIGateways(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListVDIGatewaysResult)
	vdiGatewayCount := len(*listResult.VDIGateways)
	if vdiGatewayCount != 1 {
		return resp, fmt.Errorf("found %d VDI Gateways for %v", vdiGatewayCount, name)
	}
	firstRecord := (*listResult.VDIGateways)[0]
	vdiGatewayID := firstRecord.ID
	return client.GetVDIGateway(vdiGatewayID, &Request{})
}
