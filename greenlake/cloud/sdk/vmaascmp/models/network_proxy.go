//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

type GetAllNetworkProxies struct {
	GetNetworkProxies []GetNetworkProxy `json:"networkProxies"`
	NetworkProxyCount int               `json:"networkProxyCount"`
	Meta              Meta              `json:"meta"`
}

type GetSpecificNetworkProxy struct {
	GetNetworkProxies GetNetworkProxy `json:"networkProxy"`
}

type GetNetworkProxy struct {
	ID         int     `json:"id" tf:"id,computed"`
	Name       string  `json:"name" tf:"name"`
	ProxyHost  string  `json:"proxyHost"`
	ProxyPort  int     `json:"proxyPort"`
	Visibility string  `json:"visibility"`
	Account    IDModel `json:"account"`
	Owner      IDModel `json:"owner"`
}
