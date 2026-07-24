// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package models

type ServicePlans struct {
	ServicePlansResponse []ServicePlanResponse `json:"servicePlans"`
}

type GetSpecificServicePlan struct {
	ServicePlansResponse ServicePlanResponse `json:"servicePlan"`
}

type ServicePlanResponse struct {
	ID                   int         `json:"id"`
	Name                 string      `json:"name"`
	Code                 string      `json:"code"`
	Active               bool        `json:"active"`
	SortOrder            int         `json:"sortOrder"`
	Description          interface{} `json:"description"`
	MaxStorage           int64       `json:"maxStorage"`
	MaxMemory            int64       `json:"maxMemory"`
	MaxCPU               interface{} `json:"maxCpu"`
	MaxCores             int         `json:"maxCores"`
	MaxDisks             int         `json:"maxDisks"`
	CoresPerSocket       int         `json:"coresPerSocket"`
	CustomCPU            bool        `json:"customCpu"`
	CustomCores          bool        `json:"customCores"`
	CustomMaxStorage     bool        `json:"customMaxStorage"`
	CustomMaxDataStorage bool        `json:"customMaxDataStorage"`
	CustomMaxMemory      bool        `json:"customMaxMemory"`
	AddVolumes           bool        `json:"addVolumes"`
	MemoryOptionSource   interface{} `json:"memoryOptionSource"`
	CPUOptionSource      interface{} `json:"cpuOptionSource"`
	DateCreated          string      `json:"dateCreated"`
	LastUpdated          string      `json:"lastUpdated"`
	RegionCode           interface{} `json:"regionCode"`
	Visibility           string      `json:"visibility"`
	Editable             bool        `json:"editable"`
	ProvisionType        struct {
		ID                        int    `json:"id"`
		Name                      string `json:"name"`
		Code                      string `json:"code"`
		RootDiskCustomizable      bool   `json:"rootDiskCustomizable"`
		AddVolumes                bool   `json:"addVolumes"`
		CustomizeVolume           bool   `json:"customizeVolume"`
		HasConfigurableCPUSockets bool   `json:"hasConfigurableCpuSockets"`
	} `json:"provisionType"`
	Tenants   string `json:"tenants"`
	PriceSets []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		PriceUnit string `json:"priceUnit"`
	} `json:"priceSets"`
	Config struct{} `json:"config"`
}
