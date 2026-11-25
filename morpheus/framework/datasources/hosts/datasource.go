// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package hosts implements the hpe_morpheus_hosts data source.
package hosts

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
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/hosttypes"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const summary = "read hosts data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "hosts"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = HostsDataSourceSchema(ctx)
}

// matchesClientFilters applies the host-relevant filter that is not expressed
// as a server-side query parameter: compute_server_type (code). bare_metal is
// a query parameter of the list (bareMetalHost).
func matchesClientFilters(
	srv *sdk.ListHosts200ResponseAllOfServersInner,
	wantType string,
	filterByType bool,
) bool {
	if filterByType {
		code := ""
		if srv.ComputeServerType != nil && srv.ComputeServerType.Code != nil {
			code = *srv.ComputeServerType.Code
		}
		if code != wantType {
			return false
		}
	}

	return true
}

// hostToValue maps an API host summary into the generated custom object value
// used as the element of the hosts set.
func hostToValue(
	ctx context.Context,
	srv *sdk.ListHosts200ResponseAllOfServersInner,
	bareMetalCodes map[string]bool,
) (HostsValue, diag.Diagnostics) {
	// /api/servers does not say whether a type is bare metal; the server types
	// list does, and bareMetalCodes holds the codes that are.
	bareMetal := hosttypes.BareMetalValue(bareMetalCodes, nil)
	cstCode := types.StringNull()

	if srv.ComputeServerType != nil {
		cstCode = convert.StrToType(srv.ComputeServerType.Code)
		bareMetal = hosttypes.BareMetalValue(bareMetalCodes, srv.ComputeServerType.Code)
	}

	parentHostID := types.Int64Null()
	if p := srv.ParentServer.Get(); p != nil {
		parentHostID = convert.Int64ToType(p.Id)
	}

	return NewHostsValue(
		HostsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"bare_metal":          bareMetal,
			"cloud_id":            convert.Int64ToType(srv.ZoneId),
			"compute_server_type": cstCode,
			"external_ip":         convert.StrToType(srv.ExternalIp.Get()),
			"id":                  convert.Int64ToType(srv.Id),
			"internal_ip":         convert.StrToType(srv.InternalIp.Get()),
			"name":                convert.StrToType(srv.Name),
			"parent_host_id":      parentHostID,
			"power_state":         convert.StrToType(srv.PowerState),
			"status":              convert.StrToType(srv.Status),
		},
	)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config HostsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Server-side filters.
	var (
		cloudID      int64
		filterCloud  bool
		groupID      int64
		filterGroup  bool
		phrase       string
		filterPhrase bool
	)

	if !config.CloudId.IsNull() && !config.CloudId.IsUnknown() {
		cloudID = config.CloudId.ValueInt64()
		filterCloud = true
	}
	if !config.GroupId.IsNull() && !config.GroupId.IsUnknown() {
		groupID = config.GroupId.ValueInt64()
		filterGroup = true
	}
	if !config.Phrase.IsNull() && !config.Phrase.IsUnknown() {
		phrase = config.Phrase.ValueString()
		filterPhrase = true
	}

	// Client-side filters.
	wantType := ""
	filterByType := false
	if !config.ComputeServerType.IsNull() && !config.ComputeServerType.IsUnknown() {
		wantType = config.ComputeServerType.ValueString()
		filterByType = true
	}

	filterByBareMetal := !config.BareMetal.IsNull() && !config.BareMetal.IsUnknown()

	servers, err := paging.Collect(ctx, func(
		ctx context.Context,
		offset, max int64,
	) ([]sdk.ListHosts200ResponseAllOfServersInner, int64, error) {
		apiReq := apiClient.HostsAPI.ListHosts(ctx).Max(max).Offset(offset)
		if filterByBareMetal {
			apiReq = apiReq.BareMetalHost(config.BareMetal.ValueBool())
		}
		if filterCloud {
			apiReq = apiReq.ZoneId(cloudID)
		}
		if filterGroup {
			apiReq = apiReq.SiteId(groupID)
		}
		if filterPhrase {
			apiReq = apiReq.Phrase(phrase)
		}

		rs, hresp, err := apiReq.Execute()
		if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
			return nil, 0, fmt.Errorf(
				"LIST failed for hosts: %s", providererrors.ErrMsg(err, hresp))
		}

		total := int64(0)
		if rs.Meta != nil && rs.Meta.Total != nil {
			total = *rs.Meta.Total
		}

		return rs.Servers, total, nil
	})
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	// A list that cannot be read leaves bare_metal null ("not known") rather
	// than failing the read, as the resource does.
	bareMetalCodes, err := hosttypes.BareMetalCodes(ctx, apiClient)
	if err != nil {
		resp.Diagnostics.AddWarning("bare_metal could not be determined", err.Error())
	}

	objs := make([]attr.Value, 0, len(servers))

	for i := range servers {
		srv := &servers[i]

		if !matchesClientFilters(srv, wantType, filterByType) {
			continue
		}

		v, diags := hostToValue(ctx, srv, bareMetalCodes)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(HostsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Hosts = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
