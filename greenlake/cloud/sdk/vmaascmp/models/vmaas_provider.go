// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package models

type TFVmaas struct {
	Location  string `tf:"location"`
	SpaceName string `tf:"space_name"`
}

type TFProvider struct {
	Vmaas TFVmaas `tf:"vmaas,sub"`
}
