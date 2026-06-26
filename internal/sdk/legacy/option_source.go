package morpheus

import (
	"fmt"
)

var (
	// OptionSourcePath is the API endpoint for option sources
	OptionSourcePath = "/api/options"
)

// GetOptionSourceResult generic type of result returned by this endpoint
type GetOptionSourceResult struct {
	Data *[]OptionSourceOption `json:"data"`
}

type OptionSourceOption struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // ugh, this can be a number or a string
	// sometimes ID and Code are also returned as a convenience
	ID         int64  `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"externalId"`
}

// types for /api/options/layoutsForCloud

type GetOptionSourceLayoutsResult struct {
	Data *[]LayoutOption `json:"data"`
}

type LayoutOption struct {
	OptionSourceOption
	Version string `json:"version"`
}

// types for /api/options/zoneNetworkOptions

type GetOptionSourceZoneNetworkOptionsResult struct {
	Data *ZoneNetworkOptionsData `json:"data"`
}

type ZoneNetworkOptionsData struct {
	Networks      *[]NetworkOption      `json:"networks"`
	NetworkGroups *[]NetworkGroupOption `json:"networkGroups"`
	NetworkTypes  *[]NetworkTypeOption  `json:"networkTypes"`
}

type NetworkOption struct {
	ID         string `json:"id"` // like network-45
	Name       string `json:"name"`
	DhcpServer bool   `json:"dchpServer"`
}

type NetworkGroupOption struct {
	ID   string `json:"id"` // like networkGroup-55
	Name string `json:"name"`
}

type NetworkTypeOption struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	ExternalId  string `json:"externalId"`
	Enabled     bool   `json:"enabled"`
	DefaultType bool   `json:"defaultType"`
}

// API endpoints

func (client *Client) GetOptionSource(optionSource string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%s", OptionSourcePath, optionSource),
		QueryParams: req.QueryParams,
		Result:      &GetOptionSourceResult{},
	})
}

func (client *Client) GetOptionSourceLayouts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/layoutsForCloud", OptionSourcePath),
		QueryParams: req.QueryParams,
		Result:      &GetOptionSourceLayoutsResult{},
	})
}

func (client *Client) GetOptionSourceZoneNetworkOptions(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/zoneNetworkOptions", OptionSourcePath),
		QueryParams: req.QueryParams,
		Result:      &GetOptionSourceZoneNetworkOptionsResult{},
	})
}
