// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package securitygrouprule implements a data source for security_group_rule
package securitygrouprule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	summary                         = "read security group rule data source"
	ErrorNoValidSearchTerms         = `no valid search terms - an id or name is required`
	ErrorNoSecurityGroupRuleFound   = `no security group rule found`
	ErrorMultipleSecurityGroupRules = `multiple security group rules were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "security_group_rule"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SecurityGroupRuleDataSourceSchema(ctx)
}
