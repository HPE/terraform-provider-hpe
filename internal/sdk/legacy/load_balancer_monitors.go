package morpheus

import (
	"fmt"
)

var (
	// LoadBalancerMonitorsPath is the API endpoint for load balancer monitors
	LoadBalancerMonitorsPath = "/api/load-balancers"
)

// LoadBalancer structures for use in request and response payloads
type LoadBalancerMonitor struct {
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
	Code                string `json:"code"`
	Category            string `json:"category"`
	Visibility          string `json:"visibility"`
	Description         string `json:"description"`
	MonitorType         string `json:"monitorType"`
	MonitorInterval     int64  `json:"monitorInterval"`
	MonitorTimeout      int64  `json:"monitorTimeout"`
	SendData            string `json:"sendData"`
	SendVersion         string `json:"sendVersion"`
	SendType            string `json:"sendType"`
	ReceiveData         string `json:"receiveData"`
	ReceiveCode         string `json:"receiveCode"`
	DisabledData        string `json:"disabledData"`
	MonitorUsername     string `json:"monitorUsername"`
	MonitorPassword     string `json:"monitorPassword"`
	MonitorPasswordHash string `json:"monitorPasswordHash"`
	MonitorDestination  string `json:"monitorDestination"`
	MonitorReverse      bool   `json:"monitorReverse"`
	MonitorTransparent  bool   `json:"monitorTransparent"`
	MonitorAdaptive     bool   `json:"monitorAdaptive"`
	AliasAddress        string `json:"aliasAddress"`
	AliasPort           int64  `json:"aliasPort"`
	InternalId          string `json:"internalId"`
	ExternalId          string `json:"externalId"`
	MonitorSource       string `json:"monitorSource"`
	Status              string `json:"status"`
	StatusMessage       string `json:"statusMessage"`
	StatusDate          string `json:"statusDate"`
	Enabled             bool   `json:"enabled"`
	MaxRetry            int64  `json:"maxRetry"`
	FallCount           int64  `json:"fallCount"`
	RiseCount           int64  `json:"riseCount"`
	DataLength          string `json:"dataLength"`
	Config              struct {
		Kind         string `json:"kind"`
		Name         string `json:"name"`
		Partition    string `json:"partition"`
		FullPath     string `json:"fullPath"`
		Generation   int64  `json:"generation"`
		SelfLink     string `json:"selfLink"`
		Count        string `json:"count"`
		Debug        string `json:"debug"`
		Destination  string `json:"destination"`
		Interval     int64  `json:"interval"`
		ManualResume string `json:"manualResume"`
		TimeUntilUp  int64  `json:"timeUntilUp"`
		Timeout      int64  `json:"timeout"`
		UpInterval   int64  `json:"upInterval"`
		ExternalId   string `json:"externalId"`
		ServiceType  string `json:"serviceType"`
	} `json:"config"`
	CreatedBy   interface{} `json:"createdBy"`
	DateCreated string      `json:"dateCreated"`
	LastUpdated string      `json:"lastUpdated"`
}

// ListLoadBalancerMonitorsResult structure parses the list load balancers response payload
type ListLoadBalancerMonitorsResult struct {
	LoadBalancerMonitors *[]LoadBalancerMonitor `json:"loadBalancerMonitors"`
	Meta                 *MetaResult            `json:"meta"`
}

type GetLoadBalancerMonitorResult struct {
	LoadBalancerMonitor *LoadBalancerMonitor `json:"loadBalancerMonitor"`
}

type CreateLoadBalancerMonitorResult struct {
	Success             bool                 `json:"success"`
	Message             string               `json:"msg"`
	Errors              map[string]string    `json:"errors"`
	LoadBalancerMonitor *LoadBalancerMonitor `json:"loadBalancerMonitor"`
}

type UpdateLoadBalancerMonitorResult struct {
	CreateLoadBalancerMonitorResult
}

type DeleteLoadBalancerMonitorResult struct {
	DeleteResult
}

// ListLoadBalancerMonitors lists all load balancer monitors
func (client *Client) ListLoadBalancerMonitors(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/monitors", LoadBalancerMonitorsPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancerMonitorsResult{},
	})
}

// GetLoadBalancerMonitor gets an existing load balancer monitor
func (client *Client) GetLoadBalancerMonitor(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/monitors/%d", LoadBalancerMonitorsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerMonitorResult{},
	})
}

// CreateLoadBalancerMonitor creates a new load balancer monitor
func (client *Client) CreateLoadBalancerMonitor(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/monitors", LoadBalancerMonitorsPath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateLoadBalancerMonitorResult{},
	})
}

// UpdateLoadBalancerMonitor updates an existing load balancer monitor
func (client *Client) UpdateLoadBalancerMonitor(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/monitors/%d", LoadBalancerMonitorsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerMonitorResult{},
	})
}

// DeleteLoadBalancerMonitor deletes an existing load balancer monitor
func (client *Client) DeleteLoadBalancerMonitor(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/monitors/%d", LoadBalancerMonitorsPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteLoadBalancerMonitorResult{},
	})
}

// FindLoadBalancerMonitorByName gets an existing load balancer monitor by name
func (client *Client) FindLoadBalancerMonitorByName(id int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancerMonitors(id, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancersResult)
	loadBalancerMonitorsCount := len(*listResult.LoadBalancers)
	if loadBalancerMonitorsCount != 1 {
		return resp, fmt.Errorf("found %d Load Balancers for %v", loadBalancerMonitorsCount, name)
	}
	firstRecord := (*listResult.LoadBalancers)[0]
	loadBalancerID := firstRecord.ID
	return client.GetLoadBalancer(loadBalancerID, &Request{})
}
