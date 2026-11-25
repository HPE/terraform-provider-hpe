// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package host implements the hpe_morpheus_host data source.
package host

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/hosttypes"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const (
	summary = "read host data source"
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
	resp.TypeName = req.ProviderTypeName + "_" + "host"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = HostDataSourceSchema(ctx)
}

// findHostIDByName resolves a host id from a name (and optional cloud_id) using
// the list endpoint with phrase, then an exact-name match. It errors on zero or
// more than one match.
func findHostIDByName(
	ctx context.Context,
	data *HostModel,
	apiClient *sdk.APIClient,
) (int64, error) {
	name := data.Name.ValueString()

	// Every page of the phrase search, so an exact match (or a duplicate)
	// beyond the first page is not missed.
	servers, err := paging.Collect(ctx, func(
		ctx context.Context,
		offset, max int64,
	) ([]sdk.ListHosts200ResponseAllOfServersInner, int64, error) {
		req := apiClient.HostsAPI.ListHosts(ctx).Phrase(name).Max(max).Offset(offset)
		if !data.CloudId.IsNull() && !data.CloudId.IsUnknown() {
			req = req.ZoneId(data.CloudId.ValueInt64())
		}

		rs, hresp, err := req.Execute()
		if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
			return nil, 0, fmt.Errorf(
				"LIST failed for host %q: %s", name, providererrors.ErrMsg(err, hresp))
		}

		total := int64(0)
		if rs.Meta != nil && rs.Meta.Total != nil {
			total = *rs.Meta.Total
		}

		return rs.Servers, total, nil
	})
	if err != nil {
		return 0, err
	}

	var matches []int64

	for i := range servers {
		srv := &servers[i]
		if srv.Name == nil || *srv.Name != name || srv.Id == nil {
			continue
		}

		matches = append(matches, *srv.Id)
	}

	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("no host found with name %q", name)
	case 1:
		return matches[0], nil
	default:
		return 0, fmt.Errorf(
			"%d hosts found with name %q; refine with cloud_id or use id", len(matches), name)
	}
}

// getHostByID fetches a host by id and maps it onto the model.
func getHostByID(
	ctx context.Context,
	id int64,
	data *HostModel,
	apiClient *sdk.APIClient,
	warnings *[]string,
) error {
	h, hresp, err := apiClient.HostsAPI.GetHost(ctx, sdk.GetHostIdParameter{Int64: &id}).Execute()
	if h == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"GET failed for host %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	srv := h.Server
	if srv == nil {
		return fmt.Errorf("API response was missing the expected server object for host %d", id)
	}

	mapServer(srv, data)

	// Whether the type is bare metal is not in the server's rendering; the
	// server types list says.
	// A list that cannot be read leaves bare_metal null ("not known") rather
	// than failing the read, as the resource does.
	codes, err := hosttypes.BareMetalCodes(ctx, apiClient)
	if err != nil {
		*warnings = append(*warnings, err.Error()+"; it is null until the next read")
	}

	var code *string
	if srv.ComputeServerType != nil {
		code = srv.ComputeServerType.Code
	}
	data.BareMetal = hosttypes.BareMetalValue(codes, code)

	return nil
}

// mapServer maps an API host (server) onto the data source model.
func mapServer(srv *sdk.GetHost200ResponseServer, data *HostModel) {
	data.Id = convert.Int64ToType(srv.Id)
	data.Uuid = convert.StrToType(srv.Uuid)
	data.ExternalId = convert.StrToType(srv.ExternalId.Get())
	data.Name = convert.StrToType(srv.Name)
	data.Description = convert.StrToType(srv.Description.Get())
	data.Status = convert.StrToType(srv.Status)
	data.PowerState = convert.StrToType(srv.PowerState)
	data.InternalIp = convert.StrToType(srv.InternalIp.Get())
	data.ExternalIp = convert.StrToType(srv.ExternalIp.Get())
	data.Hostname = convert.StrToType(srv.Hostname)
	data.Managed = convert.BoolToType(srv.Managed)
	data.Enabled = convert.BoolToType(srv.Enabled)
	data.Visibility = convert.StrToType(srv.Visibility)
	data.AgentInstalled = convert.BoolToType(srv.AgentInstalled)
	data.AgentVersion = convert.StrToType(srv.AgentVersion.Get())
	data.DateCreated = convert.TimeToType(srv.DateCreated)
	data.LastUpdated = convert.TimeToType(srv.LastUpdated)
	data.CloudId = convert.Int64ToType(srv.ZoneId)
	data.GroupId = convert.Int64ToType(srv.SiteId)
	data.Labels = convert.StrSliceToSet(srv.Labels)

	// Compute server type.
	data.ComputeServerType = types.StringNull()
	data.ComputeServerTypeId = types.Int64Null()
	data.ComputeServerTypeName = types.StringNull()
	data.BareMetal = types.BoolNull()

	if cst := srv.ComputeServerType; cst != nil {
		data.ComputeServerType = convert.StrToType(cst.Code)
		data.ComputeServerTypeId = convert.Int64ToType(cst.Id)
		data.ComputeServerTypeName = convert.StrToType(cst.Name)
	}

	// Plan.
	data.PlanId = types.Int64Null()
	if srv.Plan != nil {
		data.PlanId = convert.Int64ToType(srv.Plan.Id.Get())
	}

	// Parent host (the hypervisor this host runs on, if any).
	data.ParentHostId = types.Int64Null()
	if p := srv.ParentServer; p != nil {
		data.ParentHostId = convert.Int64ToType(p.Id)
	}

	// Interfaces: the IP addresses of each network interface.
	ifaceIPs := make([]string, 0, len(srv.Interfaces))
	for i := range srv.Interfaces {
		if ip := srv.Interfaces[i].IpAddress; ip != nil && *ip != "" {
			ifaceIPs = append(ifaceIPs, *ip)
		}
	}
	data.Interfaces = strSliceToList(ifaceIPs)

	// Volumes: the name of each volume.
	volNames := make([]string, 0, len(srv.Volumes))
	for i := range srv.Volumes {
		if n := srv.Volumes[i].Name; n != nil && *n != "" {
			volNames = append(volNames, *n)
		}
	}
	data.Volumes = strSliceToList(volNames)
}

// strSliceToList converts a string slice to a types.List, null when empty.
func strSliceToList(items []string) types.List {
	if len(items) == 0 {
		return types.ListNull(types.StringType)
	}

	vals := make([]types.String, 0, len(items))
	for _, i := range items {
		vals = append(vals, types.StringValue(i))
	}

	list, diags := types.ListValueFrom(context.Background(), types.StringType, vals)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}

	return list
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data HostModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	var id int64

	switch {
	case !data.Id.IsNull() && !data.Id.IsUnknown():
		id = data.Id.ValueInt64()
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		id, err = findHostIDByName(ctx, &data, apiClient)
		if err != nil {
			resp.Diagnostics.AddError(summary, err.Error())

			return
		}
	default:
		resp.Diagnostics.AddError(
			summary, errors.New("one of id or name must be set").Error())

		return
	}

	var warnings []string
	if err := getHostByID(ctx, id, &data, apiClient, &warnings); err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	for _, w := range warnings {
		resp.Diagnostics.AddWarning("bare_metal could not be determined", w)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
