package morpheus

import (
	"fmt"
)

var (
	// LoadBalancersPath is the API endpoint for load balancers
	LoadBalancersPath = "/api/load-balancers"
)

// LoadBalancer structures for use in request and response payloads
type LoadBalancer struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AccountId int64  `json:"accountId"`
	Cloud     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"cloud"`
	Type struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"type"`
	Visibility   string `json:"visibility"`
	Description  string `json:"description"`
	Host         string `json:"host"`
	Port         int64  `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
	IP           string `json:"ip"`
	InternalIp   string `json:"internalIp"`
	ExternalIp   string `json:"externalIp"`
	ApiPort      int64  `json:"apiPort"`
	AdminPort    int64  `json:"adminPort"`
	SslEnabled   bool   `json:"sslEnabled"`
	SslCert      string `json:"sslCert"`
	Enabled      bool   `json:"enabled"`
	Config       struct {
		Scheme                    string   `json:"scheme"`
		Arn                       string   `json:"arn"`
		AmazonVpc                 string   `json:"amazonVpc"`
		SubnetIds                 []int64  `json:"subnetIds"`
		SecurityGroupIds          []string `json:"securityGroupIds"`
		CreatedDuringProvisioning bool     `json:"createdDuringProvisioning"`
		Loglevel                  string   `json:"loglevel"`
		Tier1                     string   `json:"tier1"`
		Size                      string   `json:"size"`
		AdminState                bool     `json:"adminState"`
		ServerVersion             string   `json:"serverVersion"`
		SystemVersion             string   `json:"systemVersion"`
		ResourceGroup             string   `json:"resourceGroup"`
	}
	DateCreated string `json:"dateCreated"`
	LastUpdated string `json:"lastUpdated"`
}

// ListLoadBalancersResult structure parses the list load balancers response payload
type ListLoadBalancersResult struct {
	LoadBalancers *[]LoadBalancer `json:"loadBalancers"`
	Meta          *MetaResult     `json:"meta"`
}

type GetLoadBalancerResult struct {
	LoadBalancer *LoadBalancer `json:"loadBalancer"`
}

type CreateLoadBalancerResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	LoadBalancer *LoadBalancer     `json:"loadBalancer"`
}

type UpdateLoadBalancerResult struct {
	CreateLoadBalancerResult
}

type DeleteLoadBalancerResult struct {
	DeleteResult
}

// ListLoadBalancers lists all load balancers
func (client *Client) ListLoadBalancers(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        LoadBalancersPath,
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancersResult{},
	})
}

// GetLoadBalancer gets an existing load balancer
func (client *Client) GetLoadBalancer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", LoadBalancersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerResult{},
	})
}

// CreateLoadBalancer creates a new load balancer
func (client *Client) CreateLoadBalancer(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        LoadBalancersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateLoadBalancerResult{},
	})
}

// UpdateLoadBalancer updates an existing load balancer
func (client *Client) UpdateLoadBalancer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", LoadBalancersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerResult{},
	})
}

// DeleteLoadBalancer deletes an existing load balancer
func (client *Client) DeleteLoadBalancer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", LoadBalancersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteLoadBalancerResult{},
	})
}

// RefreshLoadBalancer refreshes an existing load balancer
func (client *Client) RefreshLoadBalancer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/refresh", LoadBalancersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerResult{},
	})
}

// FindLoadBalancerByName gets an existing load balancer by name
func (client *Client) FindLoadBalancerByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancers(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancersResult)
	loadBalancersCount := len(*listResult.LoadBalancers)
	if loadBalancersCount != 1 {
		return resp, fmt.Errorf("found %d Load Balancers for %v", loadBalancersCount, name)
	}
	firstRecord := (*listResult.LoadBalancers)[0]
	loadBalancerID := firstRecord.ID
	return client.GetLoadBalancer(loadBalancerID, &Request{})
}
