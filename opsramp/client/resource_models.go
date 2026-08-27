// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

type ResourceType struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	DisplayName string `json:"displayName"`
}

type CreateResource struct {
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
	HostName     string `json:"hostName"`
	AliasName    string `json:"aliasName"`
}

type UpdateResource struct {
	ResourceType string `json:"resourceType"`
	AliasName    string `json:"aliasName"`
}

type ResourceCreated struct {
	ResourceUuid string `json:"resourceUUID"`
	TenantId     string `json:"tenantID"`
}

type GeneralInfo struct {
	HostName     string `json:"hostName"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
	AliasName    string `json:"aliasName"`
}

type GetResource struct {
	Uuid           string      `json:"id"`
	AgentInstalled bool        `json:"agentInstalled"`
	GeneralInfo    GeneralInfo `json:"generalInfo"`
}

type SearchResources struct {
	Results []SearchResource `json:"results"`
}

type SearchResource struct {
	Uuid         string `json:"id"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
}
