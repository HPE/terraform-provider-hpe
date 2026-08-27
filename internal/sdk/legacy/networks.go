package morpheus

import (
	"fmt"
)

var (
	// NetworksPath is the API endpoint for networks
	NetworksPath = "/api/networks"
)

// Network structures for use in request and response payloads
type Network struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Labels      []string `json:"labels"`
	Zone        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Type struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"type"`
	Owner struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	Code                    string      `json:"code"`
	Category                string      `json:"category"`
	InterfaceName           string      `json:"interfaceName"`
	BridgeName              string      `json:"bridgeName"`
	BridgeInterface         string      `json:"bridgeInterface"`
	Description             string      `json:"description"`
	ExternalId              string      `json:"externalId"`
	InternalId              string      `json:"internalId"`
	UniqueId                string      `json:"uniqueId"`
	ExternalType            string      `json:"externalType"`
	RefUrl                  string      `json:"refUrl"`
	RefType                 string      `json:"refType"`
	RefId                   int64       `json:"refId"`
	VlanId                  int64       `json:"vlanId"`
	VswitchName             string      `json:"vswitchName"`
	DhcpServer              bool        `json:"dhcpServer"`
	DhcpIp                  string      `json:"dhcpIp"`
	Gateway                 string      `json:"gateway"`
	Netmask                 string      `json:"netmask"`
	Broadcast               string      `json:"broadcast"`
	SubnetAddress           string      `json:"subnetAddress"`
	DnsPrimary              string      `json:"dnsPrimary"`
	DnsSecondary            string      `json:"dnsSecondary"`
	Cidr                    string      `json:"cidr"`
	TftpServer              string      `json:"tftpServer"`
	BootFile                string      `json:"bootFile"`
	SwitchId                int64       `json:"switchId"`
	FabricId                int64       `json:"fabricId"`
	NetworkRole             string      `json:"networkRole"`
	Status                  string      `json:"status"`
	AvailabilityZone        string      `json:"availabilityZone"`
	Pool                    string      `json:"pool"`
	NetworkProxy            string      `json:"networkProxy"`
	NetworkDomain           string      `json:"networkDomain"`
	SearchDomains           interface{} `json:"searchDomains"`
	PrefixLength            string      `json:"prefixLength"`
	Visibility              string      `json:"visibility"`
	EnableAdmin             bool        `json:"enableAdmin"`
	ScanNetwork             bool        `json:"scanNetwork"`
	Active                  bool        `json:"active"`
	DefaultNetwork          bool        `json:"defaultNetwork"`
	AssignPublicIp          bool        `json:"assignPublicIp"`
	ApplianceUrlProxyBypass bool        `json:"applianceUrlProxyBypass"`
	ZonePool                struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zonePool"`
	AllowStaticOverride bool `json:"allowStaticOverride"`
	Tenants             []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
	Subnets []struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		Cidr       string `json:"cidr"`
		DhcpServer bool   `json:"dhcpServer"`
		Visibility string `json:"visibility"`
		Active     bool   `json:"active"`
		Pool       string `json:"pool"`
	} `json:"subnets"`
	ResourcePermission struct {
		All      bool `json:"all"`
		AllPlans bool `json:"allPlans"`
	}
	Config struct {
		VlanIDs                 string `json:"vlanIDs"`
		ConnectedGateway        bool   `json:"connectedGateway"`
		SubnetIpManagementType  string `json:"subnetIpManagementType"`
		SubnetIpServerId        int64  `json:"subnetIpServerId"`
		DhcpRange               string `json:"dhcpRange"`
		SubnetDhcpServerAddress string `json:"subnetDhcpServerAddress"`
		SubnetDhcpLeaseTime     string `json:"subnetDhcpLeaseTime"`
	} `json:"config"`
}

type ListNetworksResult struct {
	Networks *[]Network  `json:"networks"`
	Meta     *MetaResult `json:"meta"`
}

type GetNetworkResult struct {
	Network *Network `json:"network"`
}

type CreateNetworkResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Network *Network          `json:"network"`
}

type UpdateNetworkResult struct {
	CreateNetworkResult
}

type DeleteNetworkResult struct {
	DeleteResult
}

type NetworkPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Visibility  string `json:"visibility"`
}

type CreateNetworkPayload struct {
	NetworkPayload *NetworkPayload `json:"network"`
}

type UpdateNetworkPayload struct {
	CreateNetworkPayload
}

// ListNetworks lists all networks
func (client *Client) ListNetworks(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworksPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworksResult{},
	})
}

// GetNetwork gets an existing network
func (client *Client) GetNetwork(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworksPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkResult{},
	})
}

// CreateNetwork creates a new network
func (client *Client) CreateNetwork(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworksPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkResult{},
	})
}

// UpdateNetwork updates an existing network
func (client *Client) UpdateNetwork(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkResult{},
	})
}

// DeleteNetwork deletes an existing network
func (client *Client) DeleteNetwork(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkResult{},
	})
}

// FindNetworkByName finds an existing network by the network name
func (client *Client) FindNetworkByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworks(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworksResult)
	networkCount := len(*listResult.Networks)
	if networkCount != 1 {
		return resp, fmt.Errorf("found %d networks for %v", networkCount, name)
	}
	firstRecord := (*listResult.Networks)[0]
	networkID := firstRecord.ID
	return client.GetNetwork(networkID, &Request{})
}
