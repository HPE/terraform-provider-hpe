package morpheus

import (
	"fmt"
)

var (
	// SecurityPackagesPath is the API endpoint for security packages
	SecurityPackagesPath = "/api/security-packages"
)

// SecurityPackage structures for use in request and response payloads
type SecurityPackage struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	Type   struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"account"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
	Url         string      `json:"url"`
	UUID        string      `json:"uuid"`
	Config      interface{} `json:"config"`
	DateCreated string      `json:"dateCreated"`
	LastUpdated string      `json:"lastUpdated"`
}

type ListSecurityPackagesResult struct {
	SecurityPackages *[]SecurityPackage `json:"securityPackages"`
	Meta             *MetaResult        `json:"meta"`
}

type GetSecurityPackageResult struct {
	SecurityPackage *SecurityPackage `json:"securityPackage"`
}

type CreateSecurityPackageResult struct {
	Success         bool              `json:"success"`
	Message         string            `json:"msg"`
	Errors          map[string]string `json:"errors"`
	SecurityPackage *SecurityPackage  `json:"securityPackage"`
}

type UpdateSecurityPackageResult struct {
	CreateSecurityPackageResult
}

type DeleteSecurityPackageResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListSecurityPackages(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        SecurityPackagesPath,
		QueryParams: req.QueryParams,
		Result:      &ListSecurityPackagesResult{},
	})
}

func (client *Client) GetSecurityPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", SecurityPackagesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetSecurityPackageResult{},
	})
}

func (client *Client) CreateSecurityPackage(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        SecurityPackagesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateSecurityPackageResult{},
	})
}

func (client *Client) UpdateSecurityPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", SecurityPackagesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateSecurityPackageResult{},
	})
}

func (client *Client) DeleteSecurityPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", SecurityPackagesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteSecurityPackageResult{},
	})
}

// FindSecurityPackageByName gets an existing security package by name
func (client *Client) FindSecurityPackageByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListSecurityPackages(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListSecurityPackagesResult)
	securityPackageCount := len(*listResult.SecurityPackages)
	if securityPackageCount != 1 {
		return resp, fmt.Errorf("found %d security packages for %v", securityPackageCount, name)
	}
	firstRecord := (*listResult.SecurityPackages)[0]
	securityPackageID := firstRecord.ID
	return client.GetSecurityPackage(securityPackageID, &Request{})
}
