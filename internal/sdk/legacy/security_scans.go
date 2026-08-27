package morpheus

import (
	"fmt"
	"time"
)

var (
	// SecurityScansPath is the API endpoint for securityScans
	SecurityScansPath = "/api/security-scans"
)

// SecurityScan structures for use in request and response payloads
type SecurityScan struct {
	ID              int64 `json:"id"`
	SecurityPackage struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"securityPackage"`
	Server struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"server"`
	Status       string      `json:"status"`
	ScanDate     time.Time   `json:"scanDate"`
	ScanDuration int64       `json:"scanDuration"`
	TestCount    int64       `json:"testCount"`
	RunCount     int64       `json:"runCount"`
	PassCount    int64       `json:"passCount"`
	FailCount    int64       `json:"failCount"`
	OtherCount   int64       `json:"otherCount"`
	ScanScore    int64       `json:"scanScore"`
	ExternalId   interface{} `json:"externalId"`
	CreatedBy    string      `json:"createdBy"`
	UpdatedBy    string      `json:"updatedBy"`
	DateCreated  time.Time   `json:"dateCreated"`
	LastUpdated  time.Time   `json:"lastUpdated"`
}

// ListSecurityScansResult structure parses the list securityScans response payload
type ListSecurityScansResult struct {
	SecurityScans *[]SecurityScan `json:"securityScans"`
	Meta          *MetaResult     `json:"meta"`
}

type GetSecurityScanResult struct {
	SecurityScan *SecurityScan `json:"securityScan"`
}

// ListSecurityScans lists all security scans
func (client *Client) ListSecurityScans(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        SecurityScansPath,
		QueryParams: req.QueryParams,
		Result:      &ListSecurityScansResult{},
	})
}

// GetSecurityScan gets a security scan
func (client *Client) GetSecurityScan(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", SecurityScansPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetSecurityScanResult{},
	})
}
