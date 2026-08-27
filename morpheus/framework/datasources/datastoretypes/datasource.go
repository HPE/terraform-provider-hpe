// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastoretypes

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

const summary = "read datastore types data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "datastore_types"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DatastoreTypesDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves a list of Morpheus datastore types, optionally " +
		"filtered using one or more filter blocks."
	resp.Schema.MarkdownDescription = "Retrieves a list of Morpheus datastore types, optionally " +
		"filtered using one or more filter blocks."
}

// filterModel decodes a single filter block.
type filterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
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
	var config DatastoreTypesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var blocks []filterModel
	resp.Diagnostics.Append(config.Filter.ElementsAs(ctx, &blocks, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := compileFilters(ctx, blocks, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	rs, hresp, err := apiClient.DatastoresAPI.ListDatastoreTypes(ctx).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for datastore types: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.DatastoreTypes))

	for i := range rs.DatastoreTypes {
		dt := rs.DatastoreTypes[i]
		if !datastoreTypeMatchesFilters(&dt, filters) {
			continue
		}

		v, diags := datastoreTypeInnerToValue(ctx, &dt)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(DatastoreTypesValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.DatastoreTypes = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// compileFilters converts the filter blocks from configuration into compiled
// regular expressions. Invalid patterns are reported as diagnostics.
func compileFilters(
	ctx context.Context,
	blocks []filterModel,
	diags *diag.Diagnostics,
) []compiledFilter {
	filters := make([]compiledFilter, 0, len(blocks))

	for _, b := range blocks {
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

		filters = append(filters, compiledFilter{field: field, res: res})
	}

	return filters
}

// datastoreTypeMatchesFilters reports whether dt satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func datastoreTypeMatchesFilters(
	dt *sdk.ListDatastoreTypes200ResponseDatastoreTypesInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := datastoreTypeFieldValue(dt, f.field)
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

// datastoreTypeFieldValue returns the string representation of the named filter
// field for regex matching, and whether the field is present.
func datastoreTypeFieldValue(
	dt *sdk.ListDatastoreTypes200ResponseDatastoreTypesInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if dt.Name != nil {
			return *dt.Name, true
		}
	case "code":
		if dt.Code != nil {
			return *dt.Code, true
		}
	case "external_type":
		if dt.ExternalType != nil {
			return *dt.ExternalType, true
		}
	case "disk_type":
		if dt.DiskType != nil {
			return *dt.DiskType, true
		}
	case "creatable":
		if dt.Creatable != nil {
			return strconv.FormatBool(*dt.Creatable), true
		}
	case "is_plugin":
		if dt.IsPlugin != nil {
			return strconv.FormatBool(*dt.IsPlugin), true
		}
	}

	return "", false
}

// datastoreTypeInnerToValue maps an API datastore type into the generated
// custom object value used as the element of the datastore_types set.
func datastoreTypeInnerToValue(
	ctx context.Context,
	dt *sdk.ListDatastoreTypes200ResponseDatastoreTypesInner,
) (DatastoreTypesValue, diag.Diagnostics) {
	var allDiags diag.Diagnostics

	// Build storage_server_type nested object.
	var storageServerTypeVal types.Object
	if dt.StorageServerType != nil {
		objVal, diags := types.ObjectValue(
			StorageServerTypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(dt.StorageServerType.Id),
			},
		)
		allDiags.Append(diags...)
		storageServerTypeVal = objVal
	} else {
		storageServerTypeVal = types.ObjectNull(StorageServerTypeValue{}.AttributeTypes(ctx))
	}

	if allDiags.HasError() {
		return DatastoreTypesValue{}, allDiags
	}

	// Build option_types set.
	optionTypeElems := make([]attr.Value, 0, len(dt.OptionTypes))
	for i := range dt.OptionTypes {
		ot := dt.OptionTypes[i]
		otVal, diags := optionTypeInnerToValue(ctx, &ot)
		allDiags.Append(diags...)
		if allDiags.HasError() {
			return DatastoreTypesValue{}, allDiags
		}

		otObj, diags := otVal.ToObjectValue(ctx)
		allDiags.Append(diags...)
		if allDiags.HasError() {
			return DatastoreTypesValue{}, allDiags
		}

		optionTypeElems = append(optionTypeElems, otObj)
	}

	optionTypesSet, diags := types.SetValue(
		OptionTypesType{
			ObjectType: types.ObjectType{
				AttrTypes: OptionTypesValue{}.AttributeTypes(ctx),
			},
		},
		optionTypeElems,
	)
	allDiags.Append(diags...)
	if allDiags.HasError() {
		return DatastoreTypesValue{}, allDiags
	}

	attrs := map[string]attr.Value{
		"code":                convert.StrToType(dt.Code),
		"creatable":           convert.BoolToType(dt.Creatable),
		"disk_type":           convert.StrToType(dt.DiskType),
		"editable":            convert.BoolToType(dt.Editable),
		"external_sub_type":   convert.StrToType(dt.ExternalSubType),
		"external_type":       convert.StrToType(dt.ExternalType),
		"id":                  convert.Int64ToType(dt.Id),
		"is_plugin":           convert.BoolToType(dt.IsPlugin),
		"local_storage":       convert.BoolToType(dt.LocalStorage),
		"name":                convert.StrToType(dt.Name),
		"option_types":        optionTypesSet,
		"removable":           convert.BoolToType(dt.Removable),
		"storage_server_type": storageServerTypeVal,
	}

	val, valDiags := NewDatastoreTypesValue(DatastoreTypesValue{}.AttributeTypes(ctx), attrs)
	allDiags.Append(valDiags...)

	return val, allDiags
}

// optionTypeInnerToValue maps an API option type into the generated
// OptionTypesValue.
func optionTypeInnerToValue(
	ctx context.Context,
	ot *sdk.ListDatastoreTypes200ResponseDatastoreTypesInnerOptionTypesInner,
) (OptionTypesValue, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"code":          convert.StrToType(ot.Code),
		"default_value": convert.StrToType(ot.DefaultValue.Get()),
		"display_order": convert.Int64ToType(ot.DisplayOrder),
		"field_context": convert.StrToType(ot.FieldContext),
		"field_label":   convert.StrToType(ot.FieldLabel),
		"field_name":    convert.StrToType(ot.FieldName),
		"id":            convert.Int64ToType(ot.Id),
		"name":          convert.StrToType(ot.Name),
		"option_source": convert.StrToType(ot.OptionSource.Get()),
		"required":      convert.BoolToType(ot.Required),
		"type":          convert.StrToType(ot.Type),
	}

	return NewOptionTypesValue(OptionTypesValue{}.AttributeTypes(ctx), attrs)
}
