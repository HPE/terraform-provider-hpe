// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Client (multi-tenancy) models

// CreateClient represents the request to create a new client (sub-tenant)
type CreateClient struct {
	Name        string   `json:"name"`
	UniqueId    string   `json:"uniqueId,omitempty"`
	Address     string   `json:"address,omitempty"`
	City        string   `json:"city,omitempty"`
	State       string   `json:"state,omitempty"`
	Country     string   `json:"country,omitempty"`
	Zip         string   `json:"zip,omitempty"`
	TimeZone    string   `json:"timeZone,omitempty"`
	PhoneNumber string   `json:"phoneNumber,omitempty"`
	MspId       string   `json:"mspId,omitempty"`
	ParentId    string   `json:"parentId,omitempty"`
	Addons      []string `json:"addOns,omitempty"`
	Packages    []string `json:"packages,omitempty"`
}

// ClientResponse represents the API response for a client
type ClientResponse struct {
	Id                uint     `json:"id"`
	UniqueId          string   `json:"uniqueId"`
	Name              string   `json:"name"`
	Activated         bool     `json:"activated"`
	Address           string   `json:"address"`
	City              string   `json:"city"`
	State             string   `json:"state"`
	Country           string   `json:"country"`
	Zip               string   `json:"zip"`
	TimeZone          string   `json:"timeZone"`
	PhoneNumber       string   `json:"phoneNumber"`
	MspId             string   `json:"mspId"`
	Mspid             int      `json:"mspid"`
	ParentId          string   `json:"parentId"`
	ShowCopyClipBoard bool     `json:"showCopyClipboard"`
	CreatedTime       string   `json:"createdTime"`
	UpdatedTime       string   `json:"updatedTime"`
	Addons            []string `json:"addOns"`
	Packages          []string `json:"packages"`
}

// ClientMinimal represents minimal client info for listing
type ClientMinimal struct {
	Id       uint   `json:"id"`
	UniqueId string `json:"uniqueId"`
	Name     string `json:"name"`
}

// ClientListResponse represents the response for listing clients
type ClientListResponse struct {
	Results    []ClientMinimal `json:"results"`
	TotalCount int             `json:"totalCount"`
	PageNo     int             `json:"pageNo"`
	PageSize   int             `json:"pageSize"`
}
