// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package computeservers implements a data source for compute_servers
package computeservers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary    = "read compute servers data source"
	maxResults = 10000
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
	resp.TypeName = req.ProviderTypeName + "_" + "compute_servers"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ComputeServersDataSourceSchema(ctx)
}

// compiledFilter is a filter block with its values pre-compiled as regular
// expressions.
type compiledFilter struct {
	field string
	res   []*regexp.Regexp
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ComputeServersModel

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

	// Build API request with server-side filters
	apiReq := apiClient.HostsAPI.ListHosts(ctx).Max(maxResults)

	// cloud_id → server-side ZoneId filter
	if !config.CloudId.IsNull() && !config.CloudId.IsUnknown() {
		apiReq = apiReq.ZoneId(config.CloudId.ValueInt64())
	}

	// bare_metal → server-side BareMetalHost filter
	if !config.BareMetal.IsNull() && !config.BareMetal.IsUnknown() {
		apiReq = apiReq.BareMetalHost(config.BareMetal.ValueBool())
	}

	rs, hresp, err := apiReq.Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"LIST failed for compute servers: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	// Client-side instance_id filter
	wantInstanceID := int64(0)
	filterByInstance := false

	if !config.InstanceId.IsNull() && !config.InstanceId.IsUnknown() {
		wantInstanceID = config.InstanceId.ValueInt64()
		filterByInstance = true
	}

	objs := make([]attr.Value, 0, len(rs.Servers))

	for i := range rs.Servers {
		srv := &rs.Servers[i]

		// Client-side: instance_id filter
		if filterByInstance {
			if srv.Instance == nil || srv.Instance.Id == nil || *srv.Instance.Id != wantInstanceID {
				continue
			}
		}

		// Client-side: generic filter blocks
		if !matchesFilters(srv, filters) {
			continue
		}

		v, diags := serverToValue(ctx, srv)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	// The generated schema declares the element type of servers as the custom
	// ServersType, so both the element type used here and every element must be
	// the generated type. A bare types.Object element fails the set's element
	// type check with "Invalid Set Element Type".
	setVal, diags := types.SetValue(ServersValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Servers = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func compileFilters(
	ctx context.Context,
	filterSet types.Set,
	diags *diag.Diagnostics,
) []compiledFilter {
	if filterSet.IsNull() || filterSet.IsUnknown() {
		return nil
	}

	var filterBlocks []FilterValue
	diags.Append(filterSet.ElementsAs(ctx, &filterBlocks, false)...)
	if diags.HasError() {
		return nil
	}

	compiled := make([]compiledFilter, 0, len(filterBlocks))

	for _, b := range filterBlocks {
		field := b.Name.ValueString()

		var values []string
		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)
		if diags.HasError() {
			return nil
		}

		res := make([]*regexp.Regexp, 0, len(values))
		for _, v := range values {
			re, err := regexp.Compile(v)
			if err != nil {
				diags.AddError(summary,
					fmt.Sprintf("invalid regular expression %q for filter %q: %s", v, field, err))

				return nil
			}
			res = append(res, re)
		}

		compiled = append(compiled, compiledFilter{field: field, res: res})
	}

	return compiled
}

func matchesFilters(
	srv *sdk.ListHosts200ResponseAllOfServersInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := fieldValue(srv, f.field)
		if !ok {
			return false
		}

		matched := false
		for _, re := range f.res {
			if re.MatchString(val) {
				matched = true

				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

func fieldValue(
	srv *sdk.ListHosts200ResponseAllOfServersInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if srv.Name != nil {
			return *srv.Name, true
		}
	case "status":
		if srv.Status != nil {
			return *srv.Status, true
		}
	case "power_state":
		if srv.PowerState != nil {
			return *srv.PowerState, true
		}
	case "platform":
		if v := srv.Platform.Get(); v != nil {
			return *v, true
		}
	case "hostname":
		if srv.Hostname != nil {
			return *srv.Hostname, true
		}
	case "uuid":
		if srv.Uuid != nil {
			return *srv.Uuid, true
		}
	case "external_id":
		if v := srv.ExternalId.Get(); v != nil {
			return *v, true
		}
	case "internal_id":
		if v := srv.InternalId.Get(); v != nil {
			return *v, true
		}
	case "cloud_id":
		if srv.ZoneId != nil {
			return strconv.FormatInt(*srv.ZoneId, 10), true
		}
	case "cloud_name":
		if z := srv.Zone.Get(); z != nil && z.Name != nil {
			return *z.Name, true
		}
	case "group_id":
		if srv.SiteId != nil {
			return strconv.FormatInt(*srv.SiteId, 10), true
		}
	case "compute_server_type_code":
		if srv.ComputeServerType != nil && srv.ComputeServerType.Code != nil {
			return *srv.ComputeServerType.Code, true
		}
	case "compute_server_type_name":
		if srv.ComputeServerType != nil && srv.ComputeServerType.Name != nil {
			return *srv.ComputeServerType.Name, true
		}
	case "plan_code":
		if srv.Plan != nil {
			if v := srv.Plan.Code.Get(); v != nil {
				return *v, true
			}
		}
	case "plan_name":
		if srv.Plan != nil {
			if v := srv.Plan.Name.Get(); v != nil {
				return *v, true
			}
		}
	case "visibility":
		if srv.Visibility != nil {
			return *srv.Visibility, true
		}
	}

	return "", false
}

// serverToValue maps an API host into the generated custom object value used as
// the element of the servers set. It must return ServersValue rather than a
// bare object: the schema declares the set's element type as ServersType, and a
// set rejects any element whose type is not exactly the declared element type.
func serverToValue(
	ctx context.Context,
	srv *sdk.ListHosts200ResponseAllOfServersInner,
) (ServersValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Labels
	labelsVal := types.SetNull(types.StringType)

	if len(srv.Labels) > 0 {
		vals := make([]attr.Value, 0, len(srv.Labels))
		for _, l := range srv.Labels {
			vals = append(vals, types.StringValue(l))
		}

		v, d := types.SetValue(types.StringType, vals)
		diags.Append(d...)
		labelsVal = v
	}

	// Cloud name from Zone (nullable)
	cloudName := types.StringNull()
	if z := srv.Zone.Get(); z != nil && z.Name != nil {
		cloudName = convert.StrToType(z.Name)
	}

	// ComputeServerType
	cstId := types.Int64Null()
	cstCode := types.StringNull()
	cstName := types.StringNull()
	cstManaged := types.BoolNull()

	if srv.ComputeServerType != nil {
		cstId = convert.Int64ToType(srv.ComputeServerType.Id)
		cstCode = convert.StrToType(srv.ComputeServerType.Code)
		cstName = convert.StrToType(srv.ComputeServerType.Name)
		cstManaged = convert.BoolToType(srv.ComputeServerType.Managed)
	}

	// Plan
	planId := types.Int64Null()
	planCode := types.StringNull()
	planName := types.StringNull()

	if srv.Plan != nil {
		planId = convert.Int64ToType(srv.Plan.Id.Get())
		planCode = convert.StrToType(srv.Plan.Code.Get())
		planName = convert.StrToType(srv.Plan.Name.Get())
	}

	// Instance
	instanceId := types.Int64Null()
	if srv.Instance != nil && srv.Instance.Id != nil {
		instanceId = convert.Int64ToType(srv.Instance.Id)
	}

	v, d := NewServersValue(
		ServersValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"agent_installed":             convert.BoolToType(srv.AgentInstalled),
			"cloud_id":                    convert.Int64ToType(srv.ZoneId),
			"cloud_name":                  cloudName,
			"compute_server_type_code":    cstCode,
			"compute_server_type_id":      cstId,
			"compute_server_type_managed": cstManaged,
			"compute_server_type_name":    cstName,
			"description":                 convert.StrToType(srv.Description.Get()),
			"external_id":                 convert.StrToType(srv.ExternalId.Get()),
			"external_ip":                 convert.StrToType(srv.ExternalIp.Get()),
			"group_id":                    convert.Int64ToType(srv.SiteId),
			"hostname":                    convert.StrToType(srv.Hostname),
			"id":                          convert.Int64ToType(srv.Id),
			"instance_id":                 instanceId,
			"internal_id":                 convert.StrToType(srv.InternalId.Get()),
			"internal_ip":                 convert.StrToType(srv.InternalIp.Get()),
			"labels":                      labelsVal,
			"max_memory":                  convert.Int64ToType(srv.MaxMemory),
			"max_storage":                 convert.Int64ToType(srv.MaxStorage),
			"name":                        convert.StrToType(srv.Name),
			"plan_code":                   planCode,
			"plan_id":                     planId,
			"plan_name":                   planName,
			"platform":                    convert.StrToType(srv.Platform.Get()),
			"power_state":                 convert.StrToType(srv.PowerState),
			"status":                      convert.StrToType(srv.Status),
			"uuid":                        convert.StrToType(srv.Uuid),
			"visibility":                  convert.StrToType(srv.Visibility),
		},
	)
	diags.Append(d...)

	return v, diags
}
