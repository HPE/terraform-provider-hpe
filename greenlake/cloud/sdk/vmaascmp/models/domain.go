// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

type GetAllDomains struct {
	NetworkDomains []GetDomain `json:"networkDomains"`

	NetworkDomainCount int  `json:"networkDomainCount"`
	Meta               Meta `json:"meta"`
	AccountID          int  `json:"accountId"`
}

type GetDomain struct {
	ID               int    `json:"id" tf:"id,computed"`
	Name             string `json:"name" tf:"name"`
	Active           bool   `json:"active" tf:"active,computed"`
	Visibility       string `json:"visibility"`
	DomainController bool   `json:"domainController"`
	PublicZone       bool   `json:"publicZone"`
}

type GetSpecificDomain struct {
	NetworkDomain GetDomain `json:"networkDomain"`
	SuccessOrErrorMessage
}
