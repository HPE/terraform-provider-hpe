// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolumes

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
	summary = "read storage volumes data source"
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
	resp.TypeName = req.ProviderTypeName + "_" + "storage_volumes"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = StorageVolumesDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves a list of Morpheus storage volumes, optionally " +
		"filtered using one or more filter blocks."
	resp.Schema.MarkdownDescription = "Retrieves a list of Morpheus storage volumes, optionally " +
		"filtered using one or more filter blocks."
}

// storageVolumesFilterModel decodes a single filter block.
type storageVolumesFilterModel struct {
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
	var config StorageVolumesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var blocks []storageVolumesFilterModel
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

	rs, hresp, err := apiClient.StorageAPI.ListStorageVolumes(ctx).
		Max(listMax).
		Sort("id").
		Direction("asc").
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for storage volumes: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.StorageVolumes))

	for i := range rs.StorageVolumes {
		sv := rs.StorageVolumes[i]
		if !storageVolumeMatchesFilters(&sv, filters) {
			continue
		}

		v, diags := storageVolumeInnerToValue(ctx, &sv)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(StorageVolumesValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.StorageVolumes = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// compileFilters converts the filter blocks from configuration into compiled
// regular expressions. Invalid patterns are reported as diagnostics.
func compileFilters(
	ctx context.Context,
	blocks []storageVolumesFilterModel,
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

// storageVolumeMatchesFilters reports whether sv satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func storageVolumeMatchesFilters(
	sv *sdk.ListStorageVolumes200ResponseAllOfStorageVolumesInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := storageVolumeFieldValue(sv, f.field)
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

// storageVolumeFieldValue returns the string representation of the named filter
// field for regex matching, and whether the field is present.
func storageVolumeFieldValue(
	sv *sdk.ListStorageVolumes200ResponseAllOfStorageVolumesInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if sv.Name != nil {
			return *sv.Name, true
		}
	case "type_id":
		if sv.TypeId != nil {
			return strconv.FormatInt(*sv.TypeId, 10), true
		}
	case "type_code":
		if sv.Type != nil && sv.Type.Code != nil {
			return *sv.Type.Code, true
		}
	case "datastore_id":
		if id := sv.DatastoreId.Get(); id != nil {
			return strconv.FormatInt(*id, 10), true
		}
	case "datastore_name":
		if sv.Datastore != nil && sv.Datastore.Name != nil {
			return *sv.Datastore.Name, true
		}
	case "cloud_id":
		if sv.ZoneId != nil {
			return strconv.FormatInt(*sv.ZoneId, 10), true
		}
	case "cloud_name":
		if sv.Zone != nil && sv.Zone.Name != nil {
			return *sv.Zone.Name, true
		}
	case "storage_server_id":
		if sv.StorageServer != nil && sv.StorageServer.Id != nil {
			return strconv.FormatInt(*sv.StorageServer.Id, 10), true
		}
	case "storage_server_name":
		if sv.StorageServer != nil && sv.StorageServer.Name != nil {
			return *sv.StorageServer.Name, true
		}
	case "status":
		if sv.Status != nil {
			return *sv.Status, true
		}
	case "uuid":
		if sv.Uuid != nil {
			return *sv.Uuid, true
		}
	case "provision_type":
		if pt := sv.ProvisionType.Get(); pt != nil {
			return *pt, true
		}
	case "ref_type":
		if sv.RefType != nil {
			return *sv.RefType, true
		}
	case "ref_id":
		if sv.RefId != nil {
			return strconv.FormatInt(*sv.RefId, 10), true
		}
	case "pool_name":
		if sv.PoolName != nil {
			return *sv.PoolName, true
		}
	}

	return "", false
}

// storageVolumeInnerToValue maps an API storage volume into the generated
// custom object value used as the element of the storage_volumes set.
func storageVolumeInnerToValue(
	ctx context.Context,
	sv *sdk.ListStorageVolumes200ResponseAllOfStorageVolumesInner,
) (StorageVolumesValue, diag.Diagnostics) {
	typeID := types.Int64Null()
	typeCode := types.StringNull()
	typeName := types.StringNull()
	if sv.Type != nil {
		typeID = convert.Int64ToType(sv.Type.Id)
		typeCode = convert.StrToType(sv.Type.Code)
		typeName = convert.StrToType(sv.Type.Name)
	}

	cloudID := convert.Int64ToType(sv.ZoneId)
	cloudName := types.StringNull()
	if sv.Zone != nil {
		cloudName = convert.StrToType(sv.Zone.Name)
	}

	datastoreID := convert.Int64ToType(sv.DatastoreId.Get())
	datastoreName := types.StringNull()
	if sv.Datastore != nil {
		datastoreName = convert.StrToType(sv.Datastore.Name)
	}

	storageServerID := types.Int64Null()
	storageServerName := types.StringNull()
	if sv.StorageServer != nil {
		storageServerID = convert.Int64ToType(sv.StorageServer.Id)
		storageServerName = convert.StrToType(sv.StorageServer.Name)
	}

	storageGroupID := types.Int64Null()
	storageGroupName := types.StringNull()
	if sv.StorageGroup != nil {
		storageGroupID = convert.Int64ToType(sv.StorageGroup.Id)
		storageGroupName = convert.StrToType(sv.StorageGroup.Name)
	}

	attrs := map[string]attr.Value{
		"active":                  convert.BoolToType(sv.Active),
		"category":                convert.StrToType(sv.Category),
		"claim_name":              convert.StrToType(sv.ClaimName.Get()),
		"cloud_id":                cloudID,
		"cloud_name":              cloudName,
		"configurable_iops":       convert.BoolToType(sv.ConfigurableIOPS),
		"controller_id":           convert.Int64ToType(sv.ControllerId.Get()),
		"controller_mount_point":  convert.StrToType(sv.ControllerMountPoint.Get()),
		"copy_type":               convert.StrToType(sv.CopyType.Get()),
		"create_for_multi_attach": convert.BoolToType(sv.CreateForMultiAttach),
		"datastore_id":            datastoreID,
		"datastore_name":          datastoreName,
		"datastore_option":        convert.StrToType(sv.DatastoreOption),
		"description":             convert.StrToType(sv.Description.Get()),
		"device_display_name":     convert.StrToType(sv.DeviceDisplayName),
		"device_name":             convert.StrToType(sv.DeviceName),
		"disk_mode":               convert.StrToType(sv.DiskMode),
		"disk_type":               convert.StrToType(sv.DiskType),
		"display_order":           convert.Int64ToType(sv.DisplayOrder),
		"external_id":             convert.StrToType(sv.ExternalId),
		"fiber_wwn":               convert.StrToType(sv.FiberWwn.Get()),
		"file_name":               convert.StrToType(sv.FileName.Get()),
		"id":                      convert.Int64ToType(sv.Id),
		"image_type":              convert.StrToType(sv.ImageType),
		"internal_id":             convert.StrToType(sv.InternalId.Get()),
		"is_multi_attach":         convert.BoolToType(sv.IsMultiAttach),
		"max_iops":                convert.StrToType(sv.MaxIOPS.Get()),
		"max_storage":             convert.Int64ToType(sv.MaxStorage),
		"name":                    convert.StrToType(sv.Name),
		"namespace":               convert.StrToType(sv.Namespace.Get()),
		"online":                  convert.BoolToType(sv.Online),
		"pool_name":               convert.StrToType(sv.PoolName),
		"provision_type":          convert.StrToType(sv.ProvisionType.Get()),
		"read_only":               convert.BoolToType(sv.ReadOnly),
		"ref_id":                  convert.Int64ToType(sv.RefId),
		"ref_type":                convert.StrToType(sv.RefType),
		"removable":               convert.BoolToType(sv.Removable),
		"resizeable":              convert.BoolToType(sv.Resizeable.Get()),
		"root_volume":             convert.BoolToType(sv.RootVolume),
		"share_path":              convert.StrToType(sv.SharePath.Get()),
		"source":                  convert.StrToType(sv.Source),
		"source_id":               convert.StrToType(sv.SourceId),
		"status":                  convert.StrToType(sv.Status),
		"status_message":          convert.StrToType(sv.StatusMessage.Get()),
		"storage_group_id":        storageGroupID,
		"storage_group_name":      storageGroupName,
		"storage_profile":         convert.StrToType(sv.StorageProfile.Get()),
		"storage_server_id":       storageServerID,
		"storage_server_name":     storageServerName,
		"type":                    typeName,
		"type_code":               typeCode,
		"type_id":                 typeID,
		"type_name":               typeName,
		"unique_id":               convert.StrToType(sv.UniqueId.Get()),
		"unit_number":             convert.StrToType(sv.UnitNumber.Get()),
		"used_storage":            convert.Int64ToType(sv.UsedStorage),
		"uuid":                    convert.StrToType(sv.Uuid),
		"volume_name":             convert.StrToType(sv.VolumeName),
		"volume_path":             convert.StrToType(sv.VolumePath),
		"volume_type":             convert.StrToType(sv.VolumeType),
		"wwn":                     convert.StrToType(sv.Wwn.Get()),
	}

	return NewStorageVolumesValue(StorageVolumesValue{}.AttributeTypes(ctx), attrs)
}
