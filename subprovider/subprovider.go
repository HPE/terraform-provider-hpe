// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type SubProvider interface {
	Configure(context.Context, func(any)) error
	GetName(context.Context) string
	GetSchema(context.Context) map[string]schema.Attribute
	GetDataSources(context.Context) []func() datasource.DataSource
	GetResources(context.Context) []func() resource.Resource
}
