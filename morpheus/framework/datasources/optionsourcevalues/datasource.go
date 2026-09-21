// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package optionsourcevalues implements hpe_morpheus_option_source_values: the
// full list of values an option source returns, optionally narrowed by filter
// blocks. For discovery, iteration and for_each. The singular
// hpe_morpheus_option_source_value is the plan-time check for exactly one.
package optionsourcevalues

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/optionsourcecommon"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/dsfilter"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/optionsource"
)

const summary = "read option source values data source"

// filterFieldPattern is the allow-list for a filter block's name: the two
// fixed fields, or an attributes key.
var filterFieldPattern = regexp.MustCompile(`^(name|value|attributes\.[A-Za-z0-9_.-]+)$`)

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

// Model is the data source's state.
type Model struct {
	optionsourcecommon.InputsModel

	Filter []filterModel `tfsdk:"filter"`
	Values types.List    `tfsdk:"values"`
}

type filterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "option_source_values"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	attrs := optionsourcecommon.InputAttributes()

	attrs["values"] = schema.ListNestedAttribute{
		Computed: true,
		Description: "The values the source returned, in the source's order, after the filter blocks. " +
			"Header rows the UI uses as group labels are not included.",
		MarkdownDescription: "The values the source returned, in the source's order, after the filter blocks. " +
			"Header rows the UI uses as group labels are not included.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: optionsourcecommon.ValueAttributes(),
		},
	}

	description := "Lists the values of a Morpheus option source, optionally narrowed by filter blocks. " +
		optionsourcecommon.Description

	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"filter": schema.SetNestedBlock{
				Description: "Filter block, applied client-side after the source has been read. Repeat to apply " +
					"several (all must match). Filter values are Go regular expressions and are case-sensitive.",
				MarkdownDescription: "Filter block, applied client-side after the source has been read. Repeat to apply " +
					"several (all must match). Filter values are Go regular expressions and are case-sensitive.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							Description: "The field to filter on: `name`, `value`, or `attributes.<key>` for a " +
								"source-specific attribute, e.g. `attributes.providerType`.",
							MarkdownDescription: "The field to filter on: `name`, `value`, or `attributes.<key>` for a " +
								"source-specific attribute, e.g. `attributes.providerType`.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(filterFieldPattern,
									"must be name, value, or attributes.<key>"),
							},
						},
						"values": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The patterns. A value matches the block if the field matches ANY of them. " +
								"A value that lacks the field does not match.",
							MarkdownDescription: "The patterns. A value matches the block if the field matches ANY of them. " +
								"A value that lacks the field does not match.",
							Validators: []validator.Set{setvalidator.SizeAtLeast(1)},
						},
					},
				},
			},
		},
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

	filters := compileFilters(ctx, config.Filter, &resp.Diagnostics)
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

	if len(result.Values) == 0 {
		// An empty list is a legitimate answer for most sources, so it is not
		// an error; but for the sources known to answer with an empty list
		// when a context parameter is missing, say so.
		if hint := optionsource.MissingContextHint(request.Source, request.Parameters); hint != "" {
			resp.Diagnostics.AddAttributeWarning(path.Root("source"),
				fmt.Sprintf("option source %q returned no values", request.Source), hint)
		}
	}

	elemType := types.ObjectType{AttrTypes: optionsourcecommon.ValueAttrTypes()}
	objs := make([]attr.Value, 0, len(result.Values))

	for _, v := range result.Values {
		if !dsfilter.Matches(v, filters, fieldValue) {
			continue
		}

		obj, diags := optionsourcecommon.ValueObject(v)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, obj)
	}

	list, diags := types.ListValue(elemType, objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Values = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fieldValue reads a filterable field off a value: name, value, or an
// attributes key. A value without that attribute has no field, and so does
// not match a filter on it.
func fieldValue(v optionsource.Value, field string) (string, bool) {
	switch {
	case field == "name":
		return v.Name, true
	case field == "value":
		return v.Value, true
	case strings.HasPrefix(field, "attributes."):
		s, ok := v.Attributes[strings.TrimPrefix(field, "attributes.")]

		return s, ok
	default:
		return "", false
	}
}

// compileFilters converts the configured blocks into compiled filters.
func compileFilters(ctx context.Context, blocks []filterModel, diags *diag.Diagnostics) []dsfilter.Compiled {
	if len(blocks) == 0 {
		return nil
	}

	out := make([]dsfilter.Block, 0, len(blocks))

	for _, b := range blocks {
		var values []string

		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)
		if diags.HasError() {
			return nil
		}

		out = append(out, dsfilter.Block{Name: b.Name.ValueString(), Values: values})
	}

	return dsfilter.Compile(out, summary, diags)
}
