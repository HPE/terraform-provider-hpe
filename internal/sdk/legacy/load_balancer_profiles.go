package morpheus

import (
	"fmt"
)

var (
	// LoadBalancerProfilesPath is the API endpoint for load balancer profiles
	LoadBalancerProfilesPath = "/api/load-balancers"
)

// LoadBalancer structures for use in request and response payloads
type LoadBalancerProfile struct {
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
	Code                  string      `json:"code"`
	Category              string      `json:"category"`
	ServiceType           string      `json:"serviceType"`
	ServiceTypeDisplay    string      `json:"serviceTypeDisplay"`
	Visibility            string      `json:"visibility"`
	Description           string      `json:"description"`
	InternalId            string      `json:"internalId"`
	ExternalId            string      `json:"externalId"`
	ProxyType             string      `json:"ProxyType"`
	RedirectRewrite       interface{} `json:"RedirectRewrite"`
	PersistenceType       string      `json:"PersistenceType"`
	SslEnabled            bool        `json:"SslEnabled"`
	SslCert               string      `json:"SslCert"`
	SslCertHash           string      `json:"SslCertHash"`
	AccountCertificate    string      `json:"AccountCertificate"`
	Enabled               bool        `json:"enabled"`
	RedirectUrl           string      `json:"RedirectUrl"`
	InsertXforwardedFor   bool        `json:"insertXforwardedFor"`
	PersistenceCookieName string      `json:"persistenceCookieName"`
	PersistenceExpiresIn  int64       `json:"persistenceExpiresIn"`
	Editable              bool        `json:"editable"`
	Config                struct {
		Kind                    string      `json:"kind"`
		Name                    string      `json:"name"`
		Partition               string      `json:"partition"`
		FullPath                string      `json:"fullPath"`
		Generation              int         `json:"generation"`
		SelfLink                string      `json:"selfLink"`
		AppService              string      `json:"appService"`
		ConnectionTimeout       int64       `json:"connectionTimeout"`
		EntryVirtualServer      string      `json:"entryVirtualServer"`
		ServiceDownAction       string      `json:"serviceDownAction"`
		ConnectionCloseTimeout  int64       `json:"connectionCloseTimeout"`
		FastTcpIdleTimeout      int64       `json:"fastTcpIdleTimeout"`
		FastUdpIdleTimeout      int64       `json:"fastUdpIdleTimeout"`
		HaFlowMirroring         bool        `json:"haFlowMirroring"`
		ProfileType             string      `json:"profileType"`
		XForwardedFor           string      `json:"xForwardedFor"`
		SharePersistence        bool        `json:"sharePersistence"`
		HaPersistenceMirroring  bool        `json:"haPersistenceMirroring"`
		PersistenceEntryTimeout int64       `json:"persistenceEntryTimeout"`
		ResponseTimeout         int64       `json:"responseTimeout"`
		ResponseHeaderSize      int64       `json:"responseHeaderSize"`
		RequestHeaderSize       int64       `json:"requestHeaderSize"`
		NtlmAuthentication      bool        `json:"ntlmAuthentication"`
		HttpIdleTimeout         int64       `json:"httpIdleTimeout"`
		HttpsRedirect           interface{} `json:"httpsRedirect"`
		SslSuite                string      `json:"sslSuite"`
		SupportedSslCiphers     []string    `json:"supportedSslCiphers"`
		SupportedSslProtocols   []string    `json:"supportedSslProtocols"`
		SessionCache            bool        `json:"sessionCache"`
		SessionCacheTimeout     int64       `json:"sessionCacheTimeout"`
		PreferServerCipher      bool        `json:"preferServerCipher"`
		CookieFallback          bool        `json:"cookieFallback"`
		CookieGarbling          bool        `json:"cookieGarbling"`
		CookieMode              string      `json:"cookieMode"`
		CookieName              string      `json:"cookieName"`
		Purge                   bool        `json:"purge"`
		ResourceType            string      `json:"resource_type"`
	} `json:"config"`
	CreatedBy   interface{} `json:"createdBy"`
	DateCreated string      `json:"dateCreated"`
	LastUpdated string      `json:"lastUpdated"`
}

// ListLoadBalancerProfilesResult structure parses the list load balancers response payload
type ListLoadBalancerProfilesResult struct {
	LoadBalancerProfiles *[]LoadBalancerProfile `json:"loadBalancerProfiles"`
	Meta                 *MetaResult            `json:"meta"`
}

type GetLoadBalancerProfileResult struct {
	LoadBalancerProfile *LoadBalancerProfile `json:"loadBalancerProfile"`
}

type CreateLoadBalancerProfileResult struct {
	Success             bool                 `json:"success"`
	Message             string               `json:"msg"`
	Errors              map[string]string    `json:"errors"`
	LoadBalancerProfile *LoadBalancerProfile `json:"loadBalancerProfile"`
}

type UpdateLoadBalancerProfileResult struct {
	CreateLoadBalancerProfileResult
}

type DeleteLoadBalancerProfileResult struct {
	DeleteResult
}

// ListLoadBalancerProfiles lists all load balancer profiles
func (client *Client) ListLoadBalancerProfiles(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/profiles", LoadBalancerProfilesPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancerProfilesResult{},
	})
}

// GetLoadBalancerProfile gets an existing load balancer profile
func (client *Client) GetLoadBalancerProfile(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/profiles/%d", LoadBalancerProfilesPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerProfileResult{},
	})
}

// CreateLoadBalancerProfile creates a new load balancer profile
func (client *Client) CreateLoadBalancerProfile(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/profiles", LoadBalancerProfilesPath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateLoadBalancerProfileResult{},
	})
}

// UpdateLoadBalancerProfile updates an existing load balancer profile
func (client *Client) UpdateLoadBalancerProfile(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/profiles/%d", LoadBalancerProfilesPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerProfileResult{},
	})
}

// DeleteLoadBalancerProfile deletes an existing load balancer profile
func (client *Client) DeleteLoadBalancerProfile(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/profiles/%d", LoadBalancerProfilesPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteLoadBalancerProfileResult{},
	})
}

// FindLoadBalancerProfileByName gets an existing load balancer profile by name
func (client *Client) FindLoadBalancerProfileByName(loadBalancerId int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancerProfiles(loadBalancerId, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancerProfilesResult)
	loadBalancerProfilesCount := len(*listResult.LoadBalancerProfiles)
	if loadBalancerProfilesCount != 1 {
		return resp, fmt.Errorf("found %d Load Balancer Profiles for %v", loadBalancerProfilesCount, name)
	}
	firstRecord := (*listResult.LoadBalancerProfiles)[0]
	loadBalancerProfileID := firstRecord.ID
	return client.GetLoadBalancerProfile(loadBalancerId, loadBalancerProfileID, &Request{})
}
