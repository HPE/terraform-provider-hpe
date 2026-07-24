// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

type GetAllEnvironment struct {
	Environments []GetEnvironment `json:"environments"`
}

type GetEnvironment struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetSpecificEnvironment struct {
	Environments GetEnvironment `json:"environment"`
}
