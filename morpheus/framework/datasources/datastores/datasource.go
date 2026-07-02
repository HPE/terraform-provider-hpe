// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastores

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

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
	summary = "read datastores data source"
	// listMax bounds the number of records fetched from the API in one call.
	listMax = 250
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
	resp.TypeName = req.ProviderTypeName + "_" + "datastores"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DatastoresDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves a list of Morpheus datastores, optionally " +
		"filtered using one or more filter blocks."
	resp.Schema.MarkdownDescription = "Retrieves a list of Morpheus datastores, optionally " +
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
	var config DatastoresModel

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

	rs, hresp, err := apiClient.DatastoresAPI.ListDatastores(ctx).
		Max(listMax).
		Sort("id").
		Direction("asc").
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for datastores: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.Datastores))

	for i := range rs.Datastores {
		ds := rs.Datastores[i]
		if !datastoreMatchesFilters(&ds, filters) {
			continue
		}

		v, diags := datastoreInnerToValue(ctx, &ds)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(DatastoresValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Datastores = setVal

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

// datastoreMatchesFilters reports whether ds satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func datastoreMatchesFilters(
	ds *sdk.ListDatastores200ResponseAllOfDatastoresInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := datastoreFieldValue(ds, f.field)
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

// datastoreFieldValue returns the string representation of the named filter
// field for regex matching, and whether the field is present.
func datastoreFieldValue(
	ds *sdk.ListDatastores200ResponseAllOfDatastoresInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		return ds.Name, true
	case "code":
		if ds.Code.IsSet() {
			if v := ds.Code.Get(); v != nil {
				return *v, true
			}
		}
	case "type":
		return ds.Type, true
	case "status":
		return ds.Status, true
	case "visibility":
		if ds.Visibility != nil {
			return *ds.Visibility, true
		}
	case "associated_resource_type":
		return associatedResourceTypeString(ds.RefType)
	}

	return "", false
}

// datastoreInnerToValue maps an API datastore into the generated custom object
// value used as the element of the datastores set.
func datastoreInnerToValue(
	ctx context.Context,
	ds *sdk.ListDatastores200ResponseAllOfDatastoresInner,
) (DatastoresValue, diag.Diagnostics) {
	// datastore_type is required in the SDK, map directly.
	dtID := types.Int64Value(ds.DatastoreType.Id)
	dtCode := types.StringValue(ds.DatastoreType.Code)
	dtName := types.StringValue(ds.DatastoreType.Name)

	datastoreTypeObj, dtDiags := types.ObjectValue(
		DatastoreTypeValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":   dtID,
			"code": dtCode,
			"name": dtName,
		},
	)
	if dtDiags.HasError() {
		return DatastoresValue{}, dtDiags
	}

	// storage_server nested object (optional pointer).
	var storageServerObj types.Object
	if ds.StorageServer != nil {
		storageServerObj, _ = types.ObjectValue(
			StorageServerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ds.StorageServer.Id),
			},
		)
	} else {
		storageServerObj = types.ObjectNull(StorageServerValue{}.AttributeTypes(ctx))
	}

	// zone nested object (optional pointer).
	var zoneObj types.Object
	if ds.Zone != nil {
		zoneObj, _ = types.ObjectValue(
			ZoneValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ds.Zone.Id),
			},
		)
	} else {
		zoneObj = types.ObjectNull(ZoneValue{}.AttributeTypes(ctx))
	}

	// zone_pool nested object (optional pointer).
	var zonePoolObj types.Object
	if ds.ZonePool != nil {
		zonePoolObj, _ = types.ObjectValue(
			ZonePoolValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ds.ZonePool.Id),
			},
		)
	} else {
		zonePoolObj = types.ObjectNull(ZonePoolValue{}.AttributeTypes(ctx))
	}

	// owner nested object (optional pointer).
	var ownerObj types.Object
	if ds.Owner != nil {
		ownerObj, _ = types.ObjectValue(
			OwnerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ds.Owner.Id),
			},
		)
	} else {
		ownerObj = types.ObjectNull(OwnerValue{}.AttributeTypes(ctx))
	}

	// config: empty object (the schema defines no inner attributes).
	configObj, _ := types.ObjectValue(
		ConfigValue{}.AttributeTypes(ctx),
		map[string]attr.Value{},
	)

	// Build the attribute map for the DatastoresValue.
	attrs := map[string]attr.Value{
		"active":                   convert.BoolToType(ds.Active),
		"allow_provision":          convert.BoolToType(ds.AllowProvision),
		"allow_read":               convert.BoolToType(ds.AllowRead),
		"allow_write":              convert.BoolToType(ds.AllowWrite),
		"code":                     convert.StrToType(ds.Code.Get()),
		"config":                   configObj,
		"datastore_type":           datastoreTypeObj,
		"default_store":            convert.BoolToType(ds.DefaultStore),
		"drs_enabled":              convert.BoolToType(ds.DrsEnabled),
		"external_id":              convert.StrToType(ds.ExternalId),
		"external_path":            convert.StrToType(ds.ExternalPath),
		"external_type":            convert.StrToType(ds.ExternalType),
		"free_space":               convert.Int64ToType(ds.FreeSpace.Get()),
		"heart_beat_target":        convert.BoolToType(ds.HeartBeatTarget),
		"id":                       types.Int64Value(ds.Id),
		"name":                     types.StringValue(ds.Name),
		"online":                   convert.BoolToType(ds.Online),
		"owner":                    ownerObj,
		"associated_resource_id":   convert.Int64ToType(ds.RefId),
		"associated_resource_type": associatedResourceType(ds.RefType),
		"status":                   types.StringValue(ds.Status),
		"status_message":           convert.StrToType(ds.StatusMessage),
		"storage_server":           storageServerObj,
		"storage_size":             convert.Int64ToType(ds.StorageSize.Get()),
		"type":                     types.StringValue(ds.Type),
		"visibility":               convert.StrToType(ds.Visibility),
		"zone":                     zoneObj,
		"zone_pool":                zonePoolObj,
	}

	return NewDatastoresValue(DatastoresValue{}.AttributeTypes(ctx), attrs)
}

const (
	refTypeComputeZone        = "ComputeZone"
	refTypeComputeServerGroup = "ComputeServerGroup"
	associatedResourceCloud   = "Cloud"
	associatedResourceCluster = "Cluster"
)

// associatedResourceTypeString maps a datastore's raw refType (ComputeZone /
// ComputeServerGroup) to the friendly value used by the hpe_morpheus_datastore
// resource (Cloud / Cluster). The second return is false when refType is nil;
// any unrecognised refType is returned unchanged.
func associatedResourceTypeString(refType *string) (string, bool) {
	if refType == nil {
		return "", false
	}

	switch *refType {
	case refTypeComputeZone:
		return associatedResourceCloud, true
	case refTypeComputeServerGroup:
		return associatedResourceCluster, true
	default:
		return *refType, true
	}
}

// associatedResourceType is the types.String form of
// associatedResourceTypeString (null when refType is nil).
func associatedResourceType(refType *string) types.String {
	if s, ok := associatedResourceTypeString(refType); ok {
		return types.StringValue(s)
	}

	return types.StringNull()
}
