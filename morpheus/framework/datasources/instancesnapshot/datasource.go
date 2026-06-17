// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancesnapshot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read instance snapshot data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance_snapshot"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = InstanceSnapshotDataSourceSchema(ctx)
}

// snapshotFilesAsSet maps the API snapshotFiles array into the generated
// SnapshotFilesValue set type.
func snapshotFilesAsSet(
	ctx context.Context,
	files []sdk.GetHostSnpshots200ResponseSnapshotsInnerSnapshotFilesInner,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	elemType := SnapshotFilesValue{}.Type(ctx)
	if len(files) == 0 {
		return types.SetNull(elemType), diags
	}

	volAttrTypes := VolumeValue{}.AttributeTypes(ctx)
	values := make([]attr.Value, 0, len(files))

	for _, f := range files {
		volume := types.ObjectNull(volAttrTypes)

		if f.Volume != nil {
			obj, d := types.ObjectValue(volAttrTypes, map[string]attr.Value{
				"id": convert.Int64ToType(f.Volume.Id),
			})
			diags.Append(d...)
			volume = obj
		}

		diskIndex := types.Int64Null()
		if f.DiskIndex != nil {
			diskIndex = types.Int64Value(int64(*f.DiskIndex))
		}

		values = append(values, SnapshotFilesValue{
			DiskIndex:         diskIndex,
			ExportPath:        convert.StrToType(f.ExportPath),
			ExternalId:        convert.StrToType(f.ExternalId),
			Id:                convert.Int64ToType(f.Id),
			Name:              convert.StrToType(f.Name),
			Path:              convert.StrToType(f.Path),
			SnapshotFilesType: convert.StrToType(f.Type),
			Volume:            volume,
			state:             attr.ValueStateKnown,
		})
	}

	set, d := types.SetValue(elemType, values)
	diags.Append(d...)

	return set, diags
}

// snapshotAsState maps the API snapshot response into the data source model.
func snapshotAsState(
	ctx context.Context,
	snap *sdk.GetSnapshotInstance200ResponseSnapshot,
) (InstanceSnapshotModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := InstanceSnapshotModel{
		Id:              convert.Int64ToType(snap.Id),
		Name:            convert.StrToType(snap.Name),
		Description:     convert.StrToType(snap.Description.Get()),
		ExternalId:      convert.StrToType(snap.ExternalId.Get()),
		Status:          convert.StrToType(snap.Status),
		State:           convert.StrToType(snap.State.Get()),
		SnapshotType:    convert.StrToType(snap.SnapshotType),
		Datastore:       convert.StrToType(snap.Datastore.Get()),
		ParentSnapshot:  convert.StrToType(snap.ParentSnapshot.Get()),
		CurrentlyActive: convert.BoolToType(snap.CurrentlyActive),
		MemorySnapshot:  convert.BoolToType(snap.MemorySnapshot),
		ForExport:       convert.BoolToType(snap.ForExport),
		ForBackup:       convert.BoolToType(snap.ForBackup),
	}

	if snap.Zone != nil {
		model.CloudId = convert.Int64ToType(snap.Zone.Id)
	} else {
		model.CloudId = types.Int64Null()
	}

	if snap.SnapshotCreated.IsSet() && snap.SnapshotCreated.Get() != nil {
		model.SnapshotCreated = types.StringValue(snap.SnapshotCreated.Get().Format(time.RFC3339))
	} else {
		model.SnapshotCreated = types.StringNull()
	}

	if snap.DateCreated != nil {
		model.DateCreated = types.StringValue(snap.DateCreated.Format(time.RFC3339))
	} else {
		model.DateCreated = types.StringNull()
	}

	files, d := snapshotFilesAsSet(ctx, snap.SnapshotFiles)
	diags.Append(d...)
	model.SnapshotFiles = files

	return model, diags
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config InstanceSnapshotModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client: "+err.Error(),
		)

		return
	}

	snapshotID := config.Id.ValueInt64()

	snapResp, hresp, err := client.InstancesAPI.GetSnapshotInstance(ctx, snapshotID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("snapshot %d GET failed: %s", snapshotID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if snapResp == nil || snapResp.Snapshot == nil {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("snapshot %d: response missing snapshot", snapshotID),
		)

		return
	}

	if hresp != nil && hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("snapshot %d GET returned status %d", snapshotID, hresp.StatusCode),
		)

		return
	}

	state, d2 := snapshotAsState(ctx, snapResp.Snapshot)
	resp.Diagnostics.Append(d2...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
