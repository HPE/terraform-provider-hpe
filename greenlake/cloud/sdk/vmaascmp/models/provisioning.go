// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package models

type GetAllProvisioningTypes struct {
	ProvisionTypes []ProvisionTypes `json:"provisionTypes"`
}

type ProvisionTypes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
