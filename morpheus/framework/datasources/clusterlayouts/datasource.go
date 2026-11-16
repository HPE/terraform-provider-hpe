// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusterlayouts implements a plural data source for cluster layouts
package clusterlayouts

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

const summary = "read cluster layouts data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_layouts"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterLayoutsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClusterLayoutsModel

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

	// Every page is fetched before the filter blocks are applied. The endpoint
	// pages at 25 by default, and a regular expression evaluated over a
	// truncated fetch would report "no such layout" for one that exists -- on
	// an appliance with the stock library the HVM cluster layouts sit well
	// beyond the first page -- so completeness is a correctness requirement.
	layouts, err := paging.Collect(
		ctx,
		func(ctx context.Context, offset, maxItems int64) (
			[]sdk.ListClusterLayouts200ResponseAllOfLayoutsInner, int64, error,
		) {
			rs, hresp, err := apiClient.ClusterLayoutsAPI.ListClusterLayouts(ctx).
				Max(maxItems).
				Offset(offset).
				Execute()
			if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
				return nil, 0, fmt.Errorf("LIST failed for cluster layouts: %s",
					providererrors.ErrMsg(err, hresp))
			}

			var total int64
			if rs.Meta != nil && rs.Meta.Total != nil {
				total = *rs.Meta.Total
			}

			return rs.Layouts, total, nil
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	objs := make([]attr.Value, 0, len(layouts))

	for i := range layouts {
		l := &layouts[i]

		// Client-side: generic filter blocks.
		if !dsfilter.Matches(l, filters, fieldValue) {
			continue
		}

		v, diags := layoutToValue(ctx, l)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(ClusterLayoutsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ClusterLayouts = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fieldValue resolves the value of a filterable field from the raw SDK item.
//
// provision_type is not a materialized element attribute (schema-generator
// deviation): it is filtered against the raw item's ProvisionType.Code, falling
// back to ProvisionType.Name.
func fieldValue(
	l *sdk.ListClusterLayouts200ResponseAllOfLayoutsInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if l.Name != nil {
			return *l.Name, true
		}
	case "code":
		if l.Code != nil {
			return *l.Code, true
		}
	case "cluster_version":
		if l.ClusterVersion != nil {
			return *l.ClusterVersion, true
		}
	case "compute_version":
		if l.ComputeVersion != nil {
			return *l.ComputeVersion, true
		}
	case "provision_type":
		if pt := l.ProvisionType.Get(); pt != nil {
			if pt.Code != nil {
				return *pt.Code, true
			}

			if pt.Name != nil {
				return *pt.Name, true
			}
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

// layoutToValue maps an API cluster layout into the generated custom object value.
func layoutToValue(
	ctx context.Context,
	l *sdk.ListClusterLayouts200ResponseAllOfLayoutsInner,
) (ClusterLayoutsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Labels
	labelsVal := types.SetNull(types.StringType)

	if len(l.Labels) > 0 {
		vals := make([]attr.Value, 0, len(l.Labels))
		for _, label := range l.Labels {
			vals = append(vals, types.StringValue(label))
		}

		v, d := types.SetValue(types.StringType, vals)
		diags.Append(d...)
		labelsVal = v
	}

	v, d := NewClusterLayoutsValue(
		ClusterLayoutsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"cluster_version":           convert.StrToType(l.ClusterVersion),
			"code":                      convert.StrToType(l.Code),
			"compute_version":           convert.StrToType(l.ComputeVersion),
			"creatable":                 convert.BoolToType(l.Creatable),
			"date_created":              convert.TimeToType(l.DateCreated),
			"description":               convert.StrToType(l.Description),
			"enabled":                   convert.BoolToType(l.Enabled),
			"has_auto_scale":            convert.BoolToType(l.HasAutoScale),
			"has_config":                convert.BoolToType(l.HasConfig),
			"has_settings":              convert.BoolToType(l.HasSettings),
			"id":                        convert.Int64ToType(l.Id),
			"install_container_runtime": convert.BoolToType(l.InstallContainerRuntime),
			"internal_id":               convert.StrToType(l.InternalId),
			"labels":                    labelsVal,
			"last_updated":              convert.TimeToType(l.LastUpdated),
			"memory_requirement":        convert.Int64ToType(l.MemoryRequirement),
			"name":                      convert.StrToType(l.Name),
			"server_count":              convert.Int64ToType(l.ServerCount),
			"sort_order":                convert.Int64ToType(l.SortOrder),
		},
	)
	diags.Append(d...)

	return v, diags
}
