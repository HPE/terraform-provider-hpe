// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouterfirewallrule implements a data source for network_router_firewall_rule
package networkrouterfirewallrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	summary                                 = "read network router firewall rule data source"
	ErrorNoValidSearchTerms                 = `no valid search terms - an id or name is required`
	ErrorNoNetworkRouterFirewallRuleFound   = `no network router firewall rule found`
	ErrorMultipleNetworkRouterFirewallRules = `multiple network router firewall rules were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_firewall_rule"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterFirewallRuleDataSourceSchema(ctx)
}
