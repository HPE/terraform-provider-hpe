package morpheus

import (
	"fmt"
)

var (
	// LoadBalancerPoolsPath is the API endpoint for load balancer pools
	LoadBalancerPoolsPath = "/api/load-balancers"
)

// LoadBalancerPool structures for use in request and response payloads
type LoadBalancerPool struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LoadBalancer struct {
		ID   int64 `json:"id"`
		Type struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
		Name string `json:"name"`
		IP   string `json:"ip"`
	} `json:"loadBalancer"`
	Category         string      `json:"category"`
	Visibility       string      `json:"visibility"`
	Description      string      `json:"description"`
	InternalId       string      `json:"internalId"`
	ExternalId       string      `json:"externalId"`
	Enabled          bool        `json:"enabled"`
	VipSticky        interface{} `json:"vipSticky"`
	VipBalance       string      `json:"vipBalance"`
	AllowNat         bool        `json:"allowNat"`
	AllowSnat        bool        `json:"allowSnat"`
	VipClientIpMode  string      `json:"vipClientIpMode"`
	VipServerIpMode  string      `json:"vipServerIpMode"`
	MinActive        int64       `json:"minActive"`
	MinInService     int64       `json:"minInService"`
	MinUpMonitor     string      `json:"minUpMonitor"`
	MinUpAction      string      `json:"minUpAction"`
	MaxQueueDepth    int64       `json:"maxQueueDepth"`
	MaxQueueTime     int64       `json:"maxQueueTime"`
	NumberActive     int64       `json:"numberActive"`
	NumberInService  int64       `json:"numberInService"`
	HealthScore      float64     `json:"healthScore"`
	PerformanceScore float64     `json:"performanceScore"`
	HealthPenalty    float64     `json:"healthPenalty"`
	SecurityPenalty  float64     `json:"securityPenalty"`
	ErrorPenalty     float64     `json:"errorPenalty"`
	DownAction       string      `json:"downAction"`
	RampTime         int64       `json:"rampTime"`
	Port             int64       `json:"port"`
	PortType         string      `json:"portType"`
	Status           string      `json:"status"`
	Monitors         []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"monitors"`
	Members []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"members"`
	Config struct {
		Kind                   string `json:"kind"`
		Name                   string `json:"name"`
		Partition              string `json:"partition"`
		FullPath               string `json:"fullPath"`
		Generation             int64  `json:"generation"`
		SelfLink               string `json:"selfLink"`
		AllowNat               bool   `json:"allowNat"`
		AllowSnat              bool   `json:"allowSnat"`
		Description            string `json:"description"`
		IgnorePersistedWeight  string `json:"ignorePersistedWeight"`
		IpTosToClient          string `json:"ipTosToClient"`
		IpTosToServer          string `json:"ipTosToServer"`
		LinkQosToClient        string `json:"linkQosToClient"`
		LinkQosToServer        string `json:"linkQosToServer"`
		LoadBalancingMode      string `json:"loadBalancingMode"`
		MinActiveMembers       int64  `json:"minActiveMembers"`
		MinUpMembers           int64  `json:"minUpMembers"`
		MinUpMembersAction     int64  `json:"minUpMembersAction"`
		MinUpMembersChecking   int64  `json:"minUpMembersChecking"`
		Monitor                string `json:"monitor"`
		QueueDepthLimit        int64  `json:"queueDepthLimit"`
		QueueOnConnectionLimit string `json:"queueOnConnectionLimit"`
		QueueTimeLimit         int64  `json:"queueTimeLimit"`
		ReselectTries          int64  `json:"reselectTries"`
		ServiceDownAction      string `json:"serviceDownAction"`
		SlowRampTime           int64  `json:"slowRampTime"`
		SnatTranslationType    string `json:"snatTranslationType"`
		TcpMultiplexing        bool   `json:"tcpMultiplexing"`
		TcpMultiplexingNumber  int64  `json:"tcpMultiplexingNumber"`
		ActiveMonitorPaths     int64  `json:"activeMonitorPaths"`
		PassiveMonitorPath     int64  `json:"passiveMonitorPath"`
		MemberGroup            struct {
			Name             string `json:"name"`
			Path             string `json:"path"`
			Port             int64  `json:"port"`
			MaxIpListSize    int64  `json:"maxIpListSize"`
			IpRevisionFilter string `json:"ipRevisionFilter"`
		} `json:"memberGroup"`
		MembersReference struct {
			Link            string `json:"link"`
			IsSubcollection bool   `json:"isSubcollection"`
		} `json:"membersReference"`
		ExternalId string `json:"externalId"`
	} `json:"config"`
	CreatedBy   interface{} `json:"createdBy"`
	DateCreated string      `json:"dateCreated"`
	LastUpdated string      `json:"lastUpdated"`
}

// ListLoadBalancerPoolsResult structure parses the list load balancers response payload
type ListLoadBalancerPoolsResult struct {
	LoadBalancerPools *[]LoadBalancerPool `json:"loadBalancerPools"`
	Meta              *MetaResult         `json:"meta"`
}

type GetLoadBalancerPoolResult struct {
	LoadBalancerPool *LoadBalancerPool `json:"loadBalancerPool"`
}

type CreateLoadBalancerPoolResult struct {
	Success          bool              `json:"success"`
	Message          string            `json:"msg"`
	Errors           map[string]string `json:"errors"`
	LoadBalancerPool *LoadBalancerPool `json:"loadBalancerPool"`
}

type UpdateLoadBalancerPoolResult struct {
	CreateLoadBalancerPoolResult
}

type DeleteLoadBalancerPoolResult struct {
	DeleteResult
}

// ListLoadBalancerPools lists all load balancer pools
func (client *Client) ListLoadBalancerPools(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/pools", LoadBalancerPoolsPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancerPoolsResult{},
	})
}

// GetLoadBalancerPool gets an existing load balancer pool
func (client *Client) GetLoadBalancerPool(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/pools/%d", LoadBalancerPoolsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerPoolResult{},
	})
}

// CreateLoadBalancerPool creates a new load balancer pool
func (client *Client) CreateLoadBalancerPool(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/pools", LoadBalancerPoolsPath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateLoadBalancerPoolResult{},
	})
}

// UpdateLoadBalancerPool updates an existing load balancer pool
func (client *Client) UpdateLoadBalancerPool(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/pools/%d", LoadBalancerPoolsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerPoolResult{},
	})
}

// DeleteLoadBalancerPool deletes an existing load balancer pool
func (client *Client) DeleteLoadBalancerPool(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/pools/%d", LoadBalancerPoolsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteLoadBalancerPoolResult{},
	})
}

// FindLoadBalancerPoolByName gets an existing load balancer pool by name
func (client *Client) FindLoadBalancerPoolByName(loadBalancerId int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancerPools(loadBalancerId, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancerPoolsResult)
	loadBalancerPoolsCount := len(*listResult.LoadBalancerPools)
	if loadBalancerPoolsCount != 1 {
		return resp, fmt.Errorf("found %d load balancer pools for %v", loadBalancerPoolsCount, name)
	}
	firstRecord := (*listResult.LoadBalancerPools)[0]
	loadBalancerPoolID := firstRecord.ID
	return client.GetLoadBalancerPool(loadBalancerId, loadBalancerPoolID, &Request{})
}
