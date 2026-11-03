// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package optionsourcevalue implements hpe_morpheus_option_source_value: exactly
// one value of an option source, or an error. This is the plan-time check —
// "there must be exactly one, and I want to know now" — for a value such as an
// SCVMM library share or the resource pool of an HVM cluster. The plural
// hpe_morpheus_option_source_values returns the whole list.
package optionsourcevalue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/optionsourcecommon"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/optionsource"
)

const (
	summary = "read option source value data source"
	// candidateLimit bounds how many candidates a multiple-match error lists.
	candidateLimit = 10
)

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

// Model is the data source's state. name and value are both inputs (exact
// filters) and outputs, following the convention of the singular lookup data
// sources.
type Model struct {
	optionsourcecommon.InputsModel

	Name       types.String `tfsdk:"name"`
	Value      types.String `tfsdk:"value"`
	Attributes types.Map    `tfsdk:"attributes"`
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "option_source_value"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	attrs := optionsourcecommon.InputAttributes()

	attrs["name"] = schema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "The display label of the value to select. When set, only values with exactly this name are " +
			"considered. Set on output.",
		MarkdownDescription: "The display label of the value to select. When set, only values with exactly this name are " +
			"considered. Set on output.",
	}
	attrs["value"] = schema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "The value to select, e.g. `pool-2`. When set, only values with exactly this value are " +
			"considered. On output, the value to use in configuration: the source's `value`, or its `id` when " +
			"it has no `value`.",
		MarkdownDescription: "The value to select, e.g. `pool-2`. When set, only values with exactly this value are " +
			"considered. On output, the value to use in configuration: the source's `value`, or its `id` when " +
			"it has no `value`.",
	}
	attrs["attributes"] = optionsourcecommon.ValueAttributes()["attributes"]

	description := "Reads exactly one value of a Morpheus option source and fails the plan when the source " +
		"returns none or several. " + optionsourcecommon.Description

	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes:          attrs,
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := config.Request(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client: "+err.Error())

		return
	}

	result := optionsourcecommon.Fetch(ctx, client, request, summary, &resp.Diagnostics)
	if result == nil {
		return
	}

	wantName, filterName := config.Name.ValueString(), !config.Name.IsNull() && !config.Name.IsUnknown()
	wantValue, filterValue := config.Value.ValueString(), !config.Value.IsNull() && !config.Value.IsUnknown()

	matches := make([]optionsource.Value, 0, 1)

	for _, v := range result.Values {
		if filterName && v.Name != wantName {
			continue
		}

		if filterValue && v.Value != wantValue {
			continue
		}

		matches = append(matches, v)
	}

	switch len(matches) {
	case 1:
		// fallthrough to state below
	case 0:
		resp.Diagnostics.AddError(
			fmt.Sprintf("no match for option source %q", request.Source),
			noMatchDetail(request, result, filterName, wantName, filterValue, wantValue),
		)

		return
	default:
		resp.Diagnostics.AddError(
			fmt.Sprintf("multiple matches for option source %q", request.Source),
			fmt.Sprintf("Expected exactly one value, found %d:\n    %s\n\nNarrow the result with name or value, "+
				"or use hpe_morpheus_option_source_values to select from the list.",
				len(matches), optionsourcecommon.DescribeCandidates(matches, candidateLimit)),
		)

		return
	}

	match := matches[0]

	attrsMap, diags := optionsourcecommon.AttributesMap(match)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Name = types.StringValue(match.Name)
	config.Value = types.StringValue(match.Value)
	config.Attributes = attrsMap

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// noMatchDetail explains a zero-match result. Three cases read differently:
// the source returned nothing at all (possibly for want of a context
// parameter); it returned values but none passed the name/value filter (list
// what it did return); or nothing was filtered and there was nothing to match.
func noMatchDetail(
	req optionsource.Request,
	res *optionsource.Result,
	filterName bool, wantName string,
	filterValue bool, wantValue string,
) string {
	if len(res.Values) == 0 {
		detail := "The source returned no values."

		if hint := optionsource.MissingContextHint(req.Source, req.Parameters); hint != "" {
			return detail + " " + hint
		}

		return detail + " Parameters sent: " + optionsource.DescribeParameters(req.Parameters) + "."
	}

	var criteria []string

	if filterName {
		criteria = append(criteria, fmt.Sprintf("name %q", wantName))
	}

	if filterValue {
		criteria = append(criteria, fmt.Sprintf("value %q", wantValue))
	}

	return fmt.Sprintf("None of the %d values the source returned has %s. The values are:\n    %s",
		len(res.Values), joinAnd(criteria), optionsourcecommon.DescribeCandidates(res.Values, candidateLimit))
}

func joinAnd(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	default:
		return parts[0] + " and " + parts[1]
	}
}
