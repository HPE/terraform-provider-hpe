package morpheus

import (
	"fmt"
)

var (
	// LoadBalancerTypesPath is the API endpoint for load balancers types
	LoadBalancerTypesPath = "/api/load-balancers-types"
)

// LoadBalancerType structures for use in request and response payloads
type LoadBalancerType struct {
	ID                          int64       `json:"id"`
	Name                        string      `json:"name"`
	Code                        string      `json:"code"`
	Enabled                     bool        `json:"enabled"`
	Internal                    bool        `json:"internal"`
	Creatable                   bool        `json:"creatable"`
	SupportsCerts               bool        `json:"supportsCerts"`
	SupportsHostname            bool        `json:"supportsHostname"`
	SupportsVip                 bool        `json:"supportsVip"`
	SupportsSticky              bool        `json:"supportsSticky"`
	SupportsBalancing           bool        `json:"supportsBalancing"`
	SupportsScheme              bool        `json:"supportsScheme"`
	SupportsFloatingIp          bool        `json:"supportsFloatingIp"`
	SupportsMonitor             bool        `json:"supportsMonitor"`
	SupportsPoolDetail          bool        `json:"supportsPoolDetail"`
	Editable                    bool        `json:"editable"`
	Removable                   bool        `json:"removable"`
	SharedVipMode               string      `json:"sharedVipMode"`
	CreateType                  string      `json:"createType"`
	Format                      string      `json:"format"`
	ZoneType                    interface{} `json:"zoneType"`
	CertSize                    interface{} `json:"certSize"`
	HasVirtualServers           bool        `json:"hasVirtualServers"`
	HasVirtualServerPolicies    interface{} `json:"hasVirtualServerPolicies"`
	HasMonitors                 bool        `json:"hasMonitors"`
	HasNodes                    bool        `json:"hasNodes"`
	HasNodeMonitors             bool        `json:"hasNodeMonitors"`
	HasNodeWeight               bool        `json:"hasNodeWeight"`
	HasPolicies                 bool        `json:"hasPolicies"`
	HasProfiles                 bool        `json:"hasProfiles"`
	HasRules                    bool        `json:"hasRules"`
	HasScripts                  bool        `json:"hasScripts"`
	HasServices                 bool        `json:"hasServices"`
	HasPools                    bool        `json:"hasPools"`
	HasPrivateVip               bool        `json:"hasPrivateVip"`
	CreateVirtualServers        bool        `json:"createVirtualServers"`
	CreateVirtualServerPolicies interface{} `json:"createVirtualServerPolicies"`
	CreateMonitors              bool        `json:"createMonitors"`
	CreateNodes                 bool        `json:"createNodes"`
	CreatePolicies              bool        `json:"createPolicies"`
	CreateProfiles              bool        `json:"createProfiles"`
	CreateRules                 bool        `json:"createRules"`
	CreateScripts               bool        `json:"createScripts"`
	CreateServices              bool        `json:"createServices"`
	CreatePools                 bool        `json:"createPools"`
	NameEditable                bool        `json:"nameEditable"`
	PoolMemberType              interface{} `json:"poolMemberType"`
	NodeResourceType            interface{} `json:"nodeResourceType"`
	ImageCode                   string      `json:"imageCode"`
	PoolSupportsStatus          interface{} `json:"poolSupportsStatus"`
	NodeSupportsStatus          bool        `json:"nodeSupportsStatus"`
	InstanceSupportsStatus      interface{} `json:"instanceSupportsStatus"`
	ProfileSupportsProxy        bool        `json:"profileSupportsProxy"`
	CreatePricePlans            bool        `json:"createPricePlans"`
	ProfileSupportsPersistence  interface{} `json:"profileSupportsPersistence"`
	ProfilesEditable            interface{} `json:"profilesEditable"`
	OptionTypes                 interface{} `json:"optionTypes"`
	VipOptionTypes              interface{} `json:"vipOptionTypes"`
}

// ListLoadBalancerTypesResult structure parses the list load balancer types response payload
type ListLoadBalancerTypesResult struct {
	LoadBalancerTypes *[]LoadBalancerType `json:"loadBalancerTypes"`
	Meta              *MetaResult         `json:"meta"`
}

type GetLoadBalancerTypeResult struct {
	LoadBalancerType *LoadBalancerType `json:"loadBalancerType"`
}

// ListLoadBalancerTypes lists all load balancer types
func (client *Client) ListLoadBalancerTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        LoadBalancerTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancerTypesResult{},
	})
}

// GetLoadBalancerType gets an existing load balancer type
func (client *Client) GetLoadBalancerType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", LoadBalancerTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerTypeResult{},
	})
}

// FindLoadBalancerByName gets an existing load balancer by name
func (client *Client) FindLoadBalancerTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancerTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancerTypesResult)
	loadBalancerTypesCount := len(*listResult.LoadBalancerTypes)
	if loadBalancerTypesCount != 1 {
		return resp, fmt.Errorf("found %d Load Balancers Types for %v", loadBalancerTypesCount, name)
	}
	firstRecord := (*listResult.LoadBalancerTypes)[0]
	loadBalancerTypeID := firstRecord.ID
	return client.GetLoadBalancerType(loadBalancerTypeID, &Request{})
}
