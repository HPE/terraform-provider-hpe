// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancedisktype

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

// matchedDiskType is the resolved instance disk type. Fields are held as
// pointers so a value the API omits stays null in state (mapped via the convert
// helpers) rather than being flattened to a zero value.
type matchedDiskType struct {
	id             *int32
	code           *string
	displayName    *string
	displayOrder   *int32
	volumeCategory *string
	defaultType    *bool
	enabled        *bool
}

// collectDiskTypes flattens the storage types from every service plan and
// dedupes them by id. Morpheus repeats the same storage types in every plan, so
// without deduping a lookup would always see multiple matches. The first
// occurrence of each id wins; entries without an id are skipped as they cannot
// be used as a storage_type_id. It is a pure function so it is unit testable
// without an appliance.
func collectDiskTypes(
	plans []sdk.ListInstanceServicePlans200ResponsePlansInner,
) []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner {
	seen := make(map[int32]struct{})

	var out []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner

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
func matchDiskType(
	diskTypes []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner,
	name string,
) (matchedDiskType, error) {
	target := normalizeName(name)

	var matches []matchedDiskType

	for _, st := range diskTypes {
		if st.Name == nil || normalizeName(*st.Name) != target {
			continue
		}

		matches = append(matches, matchedDiskType{
			id:             st.Id,
			code:           st.Code,
			displayName:    st.DisplayName,
			displayOrder:   st.DisplayOrder,
			volumeCategory: st.VolumeCategory,
			defaultType:    st.DefaultType,
			enabled:        st.Enabled,
		})
	}

	switch len(matches) {
	case 0:
		return matchedDiskType{}, errors.New(ErrorNoInstanceDiskTypeFound)
	case 1:
		return matches[0], nil
	default:
		return matchedDiskType{}, errors.New(ErrorMultipleInstanceDiskTypes)
	}
}

// int32PtrToType converts an optional int32 (as several service-plan storage
// type fields are typed) to a Terraform Int64 value, preserving null. The
// convert package only provides an Int64 helper.
func int32PtrToType(p *int32) types.Int64 {
	if p == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*p))
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

	// The disk type id is the value consumers need (storage_type_id); a match
	// without one cannot produce a usable result.
	if match.id == nil {
		resp.Diagnostics.AddError(summary, "matched instance disk type has no id")

		return
	}

	state := InstanceDiskTypeModel{
		Name:           config.Name,
		CloudId:        config.CloudId,
		LayoutId:       config.LayoutId,
		GroupId:        config.GroupId,
		Id:             int32PtrToType(match.id),
		Code:           convert.StrToType(match.code),
		DisplayName:    convert.StrToType(match.displayName),
		DisplayOrder:   int32PtrToType(match.displayOrder),
		VolumeCategory: convert.StrToType(match.volumeCategory),
		DefaultType:    convert.BoolToType(match.defaultType),
		Enabled:        convert.BoolToType(match.enabled),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
