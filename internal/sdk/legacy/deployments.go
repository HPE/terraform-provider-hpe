package morpheus

import (
	"fmt"
	"time"
)

var (
	// DeploymentsPath is the API endpoint for deployments
	DeploymentsPath = "/api/deployments"
)

// Deployment structures for use in request and response payloads
type Deployment struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	AccountID    int64       `json:"accountId"`
	ExternalID   interface{} `json:"externalId"`
	DateCreated  time.Time   `json:"dateCreated"`
	LastUpdated  time.Time   `json:"lastUpdated"`
	VersionCount int64       `json:"versionCount"`
}

// ListDeploymentsResult structure parses the list deployments response payload
type ListDeploymentsResult struct {
	Deployments *[]Deployment `json:"deployments"`
	Meta        *MetaResult   `json:"meta"`
}

// GetDeploymentResult structure parses the deployment response payload
type GetDeploymentResult struct {
	Deployment *Deployment `json:"deployment"`
}

// CreateDeploymentResult structure parses the create deployment response payload
type CreateDeploymentResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	Deployment *Deployment       `json:"deployment"`
}

type UpdateDeploymentResult struct {
	CreateDeploymentResult
}

type DeleteDeploymentResult struct {
	DeleteResult
}

// Client request methods

// ListDeployments get all existing deployments
func (client *Client) ListDeployments(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        DeploymentsPath,
		QueryParams: req.QueryParams,
		Result:      &ListDeploymentsResult{},
	})
}

// GetDeployment gets an existing deployment
func (client *Client) GetDeployment(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", DeploymentsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetDeploymentResult{},
	})
}

// CreateDeployment creates a new deployment
func (client *Client) CreateDeployment(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        DeploymentsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateDeploymentResult{},
	})
}

// UpdateDeployment updates an existing deployment
func (client *Client) UpdateDeployment(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", DeploymentsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateDeploymentResult{},
	})
}

// DeleteDeployment deletes an existing deployment
func (client *Client) DeleteDeployment(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", DeploymentsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteDeploymentResult{},
	})
}

func (client *Client) FindDeploymentByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListDeployments(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListDeploymentsResult)
	deploymentCount := len(*listResult.Deployments)
	if deploymentCount != 1 {
		return resp, fmt.Errorf("found %d Deployments for %v", deploymentCount, name)
	}
	firstRecord := (*listResult.Deployments)[0]
	checkID := firstRecord.ID
	return client.GetDeployment(checkID, &Request{})
}
