// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

import "encoding/json"

type IDModel struct {
	ID int `json:"id,omitempty" tf:"id"`
}

type IDStringModel struct {
	ID string `json:"id,omitempty" tf:"id"`
}

type UserNameModel struct {
	UserName string `json:"username" tf:"username"`
}

type NameModel struct {
	Name string `json:"name" tf:"name"`
}

type Meta struct {
	Max    json.Number `json:"max"`
	Offset json.Number `json:"offset"`
	Size   json.Number `json:"size"`
	Total  json.Number `json:"total"`
}
