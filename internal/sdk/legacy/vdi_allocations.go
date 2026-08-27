package morpheus

import (
	"fmt"
	"time"
)

var (
	// VDIAllocationsPath is the API endpoint for vdi allocations
	VDIAllocationsPath = "/api/vdi-allocations"
)

// VDIAllocation structures for use in request and response payloads
type VDIAllocation struct {
	ID   int64 `json:"id"`
	Pool struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"pool"`
	Instance struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"instance"`
	User struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"user"`
	Localusercreated bool        `json:"localUserCreated"`
	Persistent       bool        `json:"persistent"`
	Recyclable       bool        `json:"recyclable"`
	Status           string      `json:"status"`
	Datecreated      time.Time   `json:"dateCreated"`
	Lastupdated      time.Time   `json:"lastUpdated"`
	Lastreserved     interface{} `json:"lastReserved"`
	Releasedate      time.Time   `json:"releaseDate"`
}

type ListVDIAllocationsResult struct {
	VDIAllocations *[]VDIAllocation `json:"vdiAllocations"`
	Meta           *MetaResult      `json:"meta"`
}

type GetVDIAllocationResult struct {
	VDIAllocation *VDIAllocation `json:"vdiAllocation"`
}

// Client request methods
func (client *Client) ListVDIAllocations(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        VDIAllocationsPath,
		QueryParams: req.QueryParams,
		Result:      &ListVDIAllocationsResult{},
	})
}

func (client *Client) GetVDIAllocation(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", VDIAllocationsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetVDIAllocationResult{},
	})
}
