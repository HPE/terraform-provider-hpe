// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package models

type CmpVersionModel struct {
	Appliance Appliance `json:"appliance"`
}

type Appliance struct {
	BuildVersion string `json:"buildVersion"`
}
