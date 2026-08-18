// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusters implements a plural data source for clusters
package clusters

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/dsfilter"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary    = "read clusters data source"
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
	resp.TypeName = req.ProviderTypeName + "_" + "clusters"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClustersDataSourceSchema(ctx)
}

// compiledFilter is a filter block with its values pre-compiled as regular
// expressions.
// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClustersModel

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

	// Cluster type codes, resolved once for the whole result set.
	//
	// The cluster response embeds only {id, name} for its type, so the stable
	// code has to come from the cluster types endpoint. It is worth the extra
	// call: type names are display strings that change, and the mvm-cluster
	// type is already displayed as "HVM".
	typeCodes, err := clusterTypeCodes(ctx, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not list cluster types: "+err.Error())

		return
	}

	// Build API request with server-side filters
	apiReq := apiClient.ClustersAPI.ListClusters(ctx).Max(maxResults)

	// cloud_id -> server-side ZoneId filter
	if !config.CloudId.IsNull() && !config.CloudId.IsUnknown() {
		apiReq = apiReq.ZoneId(config.CloudId.ValueInt64())
	}

	// type_id -> server-side TypeId filter
	if !config.TypeId.IsNull() && !config.TypeId.IsUnknown() {
		apiReq = apiReq.TypeId(config.TypeId.ValueInt64())
	}

	rs, hresp, err := apiReq.Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"LIST failed for clusters: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.Clusters))

	for i := range rs.Clusters {
		cl := &rs.Clusters[i]

		// Client-side: generic filter blocks
		if !dsfilter.Matches(cl, filters, func(c *sdk.ListClusters200ResponseAllOfClustersInner, field string) (string, bool) {
			return fieldValue(c, typeCodes, field)
		}) {
			continue
		}

		v, diags := clusterToValue(ctx, typeCodes, cl)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(ClustersValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Clusters = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func fieldValue(
	cl *sdk.ListClusters200ResponseAllOfClustersInner,
	typeCodes map[int64]string,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if cl.Name != nil {
			return *cl.Name, true
		}
	case "status":
		if cl.Status != nil {
			return *cl.Status, true
		}
	case "type_code":
		if t := cl.Type.Get(); t != nil && t.Id != nil {
			if code, ok := typeCodes[*t.Id]; ok {
				return code, true
			}
		}
	case "type_name":
		if t := cl.Type.Get(); t != nil && t.Name != nil {
			return *t.Name, true
		}
	case "cloud_id":
		if cl.Zone != nil && cl.Zone.Id != nil {
			return strconv.FormatInt(*cl.Zone.Id, 10), true
		}
	case "cloud_name":
		if cl.Zone != nil && cl.Zone.Name != nil {
			return *cl.Zone.Name, true
		}
	case "dynamic_placement_mode":
		if v := extractConfigString(cl.Config, "dynamicPlacementMode"); v != "" {
			return v, true
		}
	}

	return "", false
}

// extractConfigString safely extracts a string value from the untyped config map.
func extractConfigString(cfg map[string]interface{}, key string) string {
	if cfg == nil {
		return ""
	}

	v, ok := cfg[key]
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

// clusterToValue maps an API cluster into the generated custom object value.
// clusterTypeCodes returns cluster type id -> stable code.
//
// Codes are not present on the cluster response, which embeds only {id, name}
// for its type. They are worth resolving separately because codes are stable
// identifiers while display names are not: the type whose code is
// "mvm-cluster" is displayed as "HVM".
func clusterTypeCodes(
	ctx context.Context,
	apiClient *sdk.APIClient,
) (map[int64]string, error) {
	out := map[int64]string{}

	types, hresp, err := apiClient.ClustersAPI.ListClusterTypes(ctx).Execute()
	if err != nil {
		return nil, errors.New(providererrors.ErrMsg(err, hresp))
	}

	if types == nil {
		return out, nil
	}

	for _, t := range types.ClusterTypes {
		if t.Id != nil && t.Code != nil {
			out[*t.Id] = *t.Code
		}
	}

	return out, nil
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

func clusterToValue(
	ctx context.Context,
	typeCodes map[int64]string,
	cl *sdk.ListClusters200ResponseAllOfClustersInner,
) (ClustersValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Labels
	labelsVal := types.SetNull(types.StringType)

	if len(cl.Labels) > 0 {
		vals := make([]attr.Value, 0, len(cl.Labels))
		for _, l := range cl.Labels {
			vals = append(vals, types.StringValue(l))
		}

		v, d := types.SetValue(types.StringType, vals)
		diags.Append(d...)
		labelsVal = v
	}

	// Cloud (Zone)
	cloudId := types.Int64Null()
	cloudName := types.StringNull()

	if cl.Zone != nil {
		cloudId = convert.Int64ToType(cl.Zone.Id)
		cloudName = convert.StrToType(cl.Zone.Name)
	}

	// Type (nullable). The cluster response carries only {id, name}; the
	// stable code is resolved from the cluster types endpoint.
	typeId := types.Int64Null()
	typeCode := types.StringNull()
	typeName := types.StringNull()

	if t := cl.Type.Get(); t != nil {
		typeId = convert.Int64ToType(t.Id)
		typeName = convert.StrToType(t.Name)

		if t.Id != nil {
			if code, ok := typeCodes[*t.Id]; ok {
				typeCode = types.StringValue(code)
			}
		}
	}

	// Group (Site -- nullable)
	groupId := types.Int64Null()

	if s := cl.Site.Get(); s != nil {
		groupId = convert.Int64ToType(s.Id)
	}

	// Dynamic placement mode from config map
	dynamicPlacement := types.StringNull()
	if v := extractConfigString(cl.Config, "dynamicPlacementMode"); v != "" {
		dynamicPlacement = types.StringValue(v)
	}

	v, d := NewClustersValue(
		ClustersValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"cloud_id":               cloudId,
			"cloud_name":             cloudName,
			"description":            convert.StrToType(cl.Description.Get()),
			"dynamic_placement_mode": dynamicPlacement,
			"group_id":               groupId,
			"id":                     convert.Int64ToType(cl.Id),
			"labels":                 labelsVal,
			"name":                   convert.StrToType(cl.Name),
			"status":                 convert.StrToType(cl.Status),
			"type_code":              typeCode,
			"type_id":                typeId,
			"type_name":              typeName,
			"visibility":             convert.StrToType(cl.Visibility),
		},
	)
	diags.Append(d...)

	return v, diags
}
