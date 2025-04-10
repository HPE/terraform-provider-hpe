// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package subprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type SubProvider interface {
	Configure(f func(any)) error
	GetName() string
	GetSchema() map[string]schema.Attribute
	GetDataSources() []func() datasource.DataSource
	GetResources() []func() resource.Resource
}
