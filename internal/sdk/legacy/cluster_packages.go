package morpheus

import (
	"fmt"
)

var (
	// ClusterPackagesPath is the API endpoint for cluster packages
	ClusterPackagesPath = "/api/library/cluster-packages"
)

// ClusterPackage structures for use in request and response payloads
type ClusterPackage struct {
	ID             int64       `json:"id"`
	Code           string      `json:"code"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Enabled        bool        `json:"enabled"`
	PackageVersion string      `json:"packageVersion"`
	PackageType    string      `json:"packageType"`
	Type           string      `json:"type"`
	Account        interface{} `json:"account"`
	RepeatInstall  bool        `json:"repeatInstall"`
	SortOrder      int64       `json:"sortOrder"`
	SpecTemplates  []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"specTemplates"`
}

// ListClusterPackagesResult structure parses the list cluster packages response payload
type ListClusterPackagesResult struct {
	ClusterPackages *[]ClusterPackage `json:"clusterPackages"`
	Meta            *MetaResult       `json:"meta"`
}

type GetClusterPackageResult struct {
	ClusterPackage *ClusterPackage `json:"clusterPackage"`
}

type CreateClusterPackageResult struct {
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
	ClusterPackage *ClusterPackage   `json:"clusterPackage"`
}

type UpdateClusterPackageResult struct {
	CreateClusterPackageResult
}

type DeleteClusterPackageResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListClusterPackages(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ClusterPackagesPath,
		QueryParams: req.QueryParams,
		Result:      &ListClusterPackagesResult{},
	})
}

// GetClusterPackage gets a cluster package
func (client *Client) GetClusterPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ClusterPackagesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetClusterPackageResult{},
	})
}

// CreateClusterPackage creates a new cluster package
func (client *Client) CreateClusterPackage(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ClusterPackagesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateClusterPackageResult{},
	})
}

// UpdateClusterPackage updates an existing cluster package
func (client *Client) UpdateClusterPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ClusterPackagesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateClusterPackageResult{},
	})
}

// DeleteClusterPackage deletes an existing cluster package
func (client *Client) DeleteClusterPackage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ClusterPackagesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteClusterPackageResult{},
	})
}

// FindClusterPackageByName gets an existing cluster package by name
func (client *Client) FindClusterPackageByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListClusterPackages(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListClusterPackagesResult)
	clusterPackageCount := len(*listResult.ClusterPackages)
	if clusterPackageCount != 1 {
		return resp, fmt.Errorf("found %d Cluster Packages for %v", clusterPackageCount, name)
	}
	firstRecord := (*listResult.ClusterPackages)[0]
	clusterPackageID := firstRecord.ID
	return client.GetClusterPackage(clusterPackageID, &Request{})
}
