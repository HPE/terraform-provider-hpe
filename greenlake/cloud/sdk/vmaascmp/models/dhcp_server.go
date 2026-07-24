package models

import (
	"encoding/json"
	"time"
)

type CreateNetworkDhcpServerRequest struct {
	NetworkDhcpServer CreateNetworkDhcpServer `json:"networkDhcpServer"`
}

type CreateNetworkDhcpServer struct {
	ID              int         `json:"-" tf:"id,computed"`
	NetworkServerID int         `json:"networkServerId" tf:"network_server_id,computed"`
	Name            string      `json:"name" tf:"name"`
	LeaseTime       int         `json:"leaseTime" tf:"lease_time"`
	ServerIPAddress string      `json:"serverIpAddress" tf:"server_address"`
	Config          *DhcpConfig `json:"config" tf:"config,sub"`
}

type DhcpConfig struct {
	EdgeCluster string `json:"edgeCluster"  tf:"edge_cluster"`
}

type CreateNetworkDhcpServerResp struct {
	Success bool `json:"success"`
	ID      int  `json:"id" tf:"id,computed"`
}

type GetNetworkDhcpServers struct {
	GetNetworkDhcpServerResp []GetNetworkDhcpServerResp `json:"networkDhcpServers"`
}

type GetNetworkDhcpServerResp struct {
	ID              int           `json:"id" tf:"id,computed"`
	DateCreated     time.Time     `json:"dateCreated"`
	ProviderID      string        `json:"providerId" tf:"provider_id,computed"`
	LastUpdated     time.Time     `json:"lastUpdated"`
	LeaseTime       int           `json:"leaseTime"`
	Name            string        `json:"name"`
	ExternalID      string        `json:"externalId"`
	ServerIPAddress string        `json:"serverIpAddress"`
	NetworkServerID NetworkServer `json:"networkServer"`
	Meta            MetaData      `json:"meta"`
}

type NetworkServer struct {
	ID int `json:"id"`
}

type MetaData struct {
	Max    json.Number `json:"max"`
	Offset json.Number `json:"offset"`
	Size   json.Number `json:"size"`
	Total  json.Number `json:"total"`
}

type GetSpecificNetworkDhcpServer struct {
	GetSpecificNetworkDhcpServerResp GetSpecificNetworkDhcpServerResp `json:"networkDhcpServer"`
}

type GetSpecificNetworkDhcpServerResp struct {
	ID              int           `json:"id" tf:"id,computed"`
	DateCreated     time.Time     `json:"dateCreated"`
	ProviderID      string        `json:"providerId" tf:"provider_id,computed"`
	LastUpdated     time.Time     `json:"lastUpdated"`
	LeaseTime       int           `json:"leaseTime"`
	Name            string        `json:"name"`
	ExternalID      string        `json:"externalId"`
	ServerIPAddress string        `json:"serverIpAddress"`
	NetworkServerID NetworkServer `json:"networkServer"`
}
