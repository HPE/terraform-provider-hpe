// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkfloatingips implements a plural data source that lists Morpheus
// network floating IPs, optionally narrowed by server-side query filters and by
// client-side filter blocks.
package networkfloatingips

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/dsfilter"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const summary = "read network floating ips data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "network_floating_ips"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkFloatingIpsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkFloatingIpsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := compileFilters(ctx, config.Filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Server-side query filters. cloud_id maps to the ZoneId parameter.
	items, err := listFloatingIps(ctx, apiClient, &config)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	objs := make([]attr.Value, 0, len(items))

	for i := range items {
		item := &items[i]

		// Client-side generic filter blocks, applied after the server-side query.
		if !dsfilter.Matches(item, filters, fieldValue) {
			continue
		}

		v, diags := floatingIpToValue(ctx, item)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(NetworkFloatingIpsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.NetworkFloatingIps = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fieldValue reads a named filter field off a floating IP list item. The set of
// valid names is constrained by the schema; an item that lacks the field does
// not match, as documented by dsfilter.
func fieldValue(
	item *sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner,
	field string,
) (string, bool) {
	switch field {
	case "ip_address":
		if item.IpAddress != nil {
			return *item.IpAddress, true
		}
	case "ip_status":
		if item.IpStatus != nil {
			return *item.IpStatus, true
		}
	case "external_id":
		if v := item.ExternalId.Get(); v != nil {
			return *v, true
		}
	case "cloud_name":
		if item.Cloud != nil && item.Cloud.Name != nil {
			return *item.Cloud.Name, true
		}
	case "server_name":
		if item.Server != nil && item.Server.Name != nil {
			return *item.Server.Name, true
		}
	}

	return "", false
}

// compileFilters converts the configured filter blocks into compiled filters.
//
// Only the extraction is data-source specific: FilterValue is generated per
// data source. Compilation and matching live in dsfilter, shared with the other
// plural data sources.
func compileFilters(
	ctx context.Context,
	filterSet types.Set,
	diags *diag.Diagnostics,
) []dsfilter.Compiled {
	if filterSet.IsNull() || filterSet.IsUnknown() {
		return nil
	}

	var filterBlocks []FilterValue

	diags.Append(filterSet.ElementsAs(ctx, &filterBlocks, false)...)

	if diags.HasError() {
		return nil
	}

	blocks := make([]dsfilter.Block, 0, len(filterBlocks))

	for _, b := range filterBlocks {
		var values []string

		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)

		if diags.HasError() {
			return nil
		}

		blocks = append(blocks, dsfilter.Block{Name: b.Name.ValueString(), Values: values})
	}

	return dsfilter.Compile(blocks, summary, diags)
}

// floatingIpToValue maps an API floating IP list item into the generated custom
// object value.
func floatingIpToValue(
	ctx context.Context,
	item *sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner,
) (NetworkFloatingIpsValue, diag.Diagnostics) {
	cloudId := types.Int64Null()
	cloudName := types.StringNull()

	if item.Cloud != nil {
		cloudId = convert.Int64ToType(item.Cloud.Id)
		cloudName = convert.StrToType(item.Cloud.Name)
	}

	serverId := types.Int64Null()
	serverName := types.StringNull()

	if item.Server != nil {
		serverId = convert.Int64ToType(item.Server.Id)
		serverName = convert.StrToType(item.Server.Name)
	}

	ipRangeStart := types.StringNull()
	ipRangeEnd := types.StringNull()

	if item.IpRange != nil {
		ipRangeStart = convert.StrToType(item.IpRange.StartAddress)
		ipRangeEnd = convert.StrToType(item.IpRange.EndAddress)
	}

	networkDomainId := types.Int64Null()
	networkDomainName := types.StringNull()

	if item.NetworkDomain != nil {
		networkDomainId = convert.Int64ToType(item.NetworkDomain.Id)
		networkDomainName = convert.StrToType(item.NetworkDomain.Name)
	}

	return NewNetworkFloatingIpsValue(
		NetworkFloatingIpsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"cloud_id":               cloudId,
			"cloud_name":             cloudName,
			"date_created":           convert.TimeToType(item.DateCreated),
			"external_id":            convert.StrToType(item.ExternalId.Get()),
			"id":                     convert.Int64ToType(item.Id),
			"ip_address":             convert.StrToType(item.IpAddress),
			"ip_range_end_address":   ipRangeEnd,
			"ip_range_start_address": ipRangeStart,
			"ip_status":              convert.StrToType(item.IpStatus),
			"last_updated":           convert.TimeToType(item.LastUpdated),
			"network_domain_id":      networkDomainId,
			"network_domain_name":    networkDomainName,
			"ptr_id":                 convert.StrToType(item.PtrId.Get()),
			"server_id":              serverId,
			"server_name":            serverName,
		},
	)
}

// listFloatingIps collects every page of the floating IP listing that matches
// the configured server-side filters, so the client-side filter blocks see the
// whole set rather than the first page.
func listFloatingIps(
	ctx context.Context, apiClient *sdk.APIClient, config *NetworkFloatingIpsModel,
) ([]sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner, error) {
	return paging.Collect(ctx,
		func(ctx context.Context, offset, max int64) (
			[]sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner, int64, error,
		) {
			apiReq := apiClient.NetworksAPI.GetAllNetworkFloatingIps(ctx).Max(max).Offset(offset)

			if !config.CloudId.IsNull() && !config.CloudId.IsUnknown() {
				apiReq = apiReq.ZoneId(config.CloudId.ValueInt64())
			}

			if !config.ServerId.IsNull() && !config.ServerId.IsUnknown() {
				apiReq = apiReq.ServerId(config.ServerId.ValueInt64())
			}

			if !config.IpStatus.IsNull() && !config.IpStatus.IsUnknown() {
				apiReq = apiReq.IpStatus(config.IpStatus.ValueString())
			}

			if !config.Phrase.IsNull() && !config.Phrase.IsUnknown() {
				apiReq = apiReq.Phrase(config.Phrase.ValueString())
			}

			rs, hresp, err := apiReq.Execute()
			if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
				return nil, 0, fmt.Errorf("LIST failed for network floating ips: %s",
					providererrors.ErrMsg(err, hresp))
			}

			var total int64
			if rs.Meta != nil && rs.Meta.Total != nil {
				total = *rs.Meta.Total
			}

			return rs.NetworkFloatingIps, total, nil
		})
}
