// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancedisktype

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                        = "read instance disk type data source"
	ErrorNoInstanceDiskTypeFound   = `no instance disk type found`
	ErrorMultipleInstanceDiskTypes = `multiple instance disk types were returned`
)

// diskType is the disk type element returned by the service plans options
// endpoint, aliased to keep the generated SDK name out of the signatures below.
type diskType = sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner

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
	resp.TypeName = req.ProviderTypeName + "_" + "instance_disk_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceDiskTypeDataSourceSchema(ctx)
	resp.Schema.Description = "Looks up a Morpheus instance disk type (storage volume type) by " +
		"name within a cloud, layout and group, returning its id for use as storage_type_id on " +
		"hpe_morpheus_instance volumes."
	resp.Schema.MarkdownDescription = "Looks up a Morpheus instance disk type (storage volume type) by " +
		"name within a cloud, layout and group, returning its `id` for use as `storage_type_id` on " +
		"`hpe_morpheus_instance` volumes."
}

// normalizeName trims surrounding whitespace and lowercases, matching hpegl
// (strings.TrimSpace(strings.ToLower(...))) so ported configurations resolve
// identically.
func normalizeName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// collectDiskTypes flattens the storage types from every service plan and
// dedupes them by id. Morpheus repeats the same storage types in every plan, so
// without deduping a lookup would always see multiple matches. The first
// occurrence of each id wins; entries without an id are skipped as they cannot
// be used as a storage_type_id. It is a pure function so it is unit testable
// without an appliance.
func collectDiskTypes(plans []sdk.ListInstanceServicePlans200ResponsePlansInner) []diskType {
	seen := make(map[int32]struct{})

	var out []diskType

	for _, p := range plans {
		for _, st := range p.StorageTypes {
			if st.Id == nil {
				continue
			}

			if _, ok := seen[*st.Id]; ok {
				continue
			}

			seen[*st.Id] = struct{}{}
			out = append(out, st)
		}
	}

	return out
}

// matchDiskType finds the single disk type whose name matches the requested name
// (case-insensitive, whitespace-trimmed). It errors on zero or more than one
// match, so the data source fails clearly rather than silently picking an
// arbitrary one. It is a pure function so the match/error logic is unit testable
// without an appliance.
func matchDiskType(diskTypes []diskType, name string) (diskType, error) {
	target := normalizeName(name)

	var matches []diskType

	for _, st := range diskTypes {
		if st.Name == nil || normalizeName(*st.Name) != target {
			continue
		}

		matches = append(matches, st)
	}

	switch len(matches) {
	case 0:
		return diskType{}, errors.New(ErrorNoInstanceDiskTypeFound)
	case 1:
		return matches[0], nil
	default:
		return diskType{}, errors.New(ErrorMultipleInstanceDiskTypes)
	}
}

// diskTypeAsState maps a matched disk type onto the Terraform model. Every
// scalar on the model is exposed except:
//
//   - optionTypes, an array of free-form maps, which carries no useful value in
//     a lookup (the network_type data source omits its equivalent too).
//   - createDatastore, where the API returns a JSON boolean but the schema types
//     the field as a nullable string. The SDK therefore surfaces it as the
//     coerced string "0"/"1", so exposing it would be actively misleading. It is
//     omitted until the schema types it as a boolean.
func diskTypeAsState(st diskType, config InstanceDiskTypeModel) InstanceDiskTypeModel {
	return InstanceDiskTypeModel{
		// Arguments, echoed back unchanged.
		Name:     config.Name,
		CloudId:  config.CloudId,
		LayoutId: config.LayoutId,
		GroupId:  config.GroupId,

		// Identity and presentation.
		Id:           convert.Int32ToType(st.Id),
		Code:         convert.StrToType(st.Code),
		DisplayName:  convert.StrToType(st.DisplayName),
		DisplayOrder: convert.Int32ToType(st.DisplayOrder),
		Description:  convert.StrToType(st.Description.Get()),

		// Classification.
		VolumeType:     convert.StrToType(st.VolumeType),
		VolumeCategory: convert.StrToType(st.VolumeCategory),
		StorageType:    convert.StrToType(st.StorageType.Get()),
		ExternalId:     convert.StrToType(st.ExternalId.Get()),

		// Capabilities.
		DefaultType:          convert.BoolToType(st.DefaultType),
		Enabled:              convert.BoolToType(st.Enabled),
		Editable:             convert.BoolToType(st.Editable),
		Deletable:            convert.BoolToType(st.Deletable),
		Resizable:            convert.BoolToType(st.Resizable),
		PlanResizable:        convert.BoolToType(st.PlanResizable),
		NameEditable:         convert.BoolToType(st.NameEditable),
		CustomLabel:          convert.BoolToType(st.CustomLabel),
		CustomSize:           convert.BoolToType(st.CustomSize),
		AutoDelete:           convert.BoolToType(st.AutoDelete),
		AllowSearch:          convert.BoolToType(st.AllowSearch),
		HasDatastore:         convert.BoolToType(st.HasDatastore),
		HasIso:               convert.BoolToType(st.HasISO),
		NoStorage:            convert.BoolToType(st.NoStorage),
		MultiAttachSupported: convert.BoolToType(st.MultiAttachSupported),
		HasActiveReplica:     convert.BoolToType(st.HasActiveReplica.Get()),

		// Sizing and IOPS. The API returns these as strings, and null when no
		// limit applies.
		ConfigurableIops:   convert.BoolToType(st.ConfigurableIOPS),
		MinIops:            convert.StrToType(st.MinIOPS.Get()),
		MaxIops:            convert.StrToType(st.MaxIOPS.Get()),
		MinStorage:         convert.StrToType(st.MinStorage.Get()),
		MaxStorage:         convert.StrToType(st.MaxStorage.Get()),
		VolumeOptionSource: convert.StrToType(st.VolumeOptionSource.Get()),
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config InstanceDiskTypeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// The instance disk types available depend on the provisioning context, so
	// the options endpoint is scoped by cloud, layout and group (Morpheus
	// zone/layout/site). This mirrors hpegl's GetStorageVolTypeID call.
	rs, hresp, err := apiClient.InstancesAPI.ListInstanceServicePlans(ctx).
		ZoneId(config.CloudId.ValueInt64()).
		LayoutId(config.LayoutId.ValueInt64()).
		SiteId(config.GroupId.ValueInt64()).
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for instance service plans: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	if len(rs.Plans) == 0 {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"no service plans found for cloud id %d, layout id %d, group id %d",
			config.CloudId.ValueInt64(), config.LayoutId.ValueInt64(), config.GroupId.ValueInt64()))

		return
	}

	match, err := matchDiskType(collectDiskTypes(rs.Plans), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := diskTypeAsState(match, config)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
