// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// settingBackupModelV0 is the state written by the retired SDKv2 resource
// (schema version 0). Its ids and retention were TypeInt (non-pointer in the
// legacy SDK), so an unset default storage bucket or schedule was stored as 0,
// and an unset retention as 0; the retention lived under retention_days only.
type settingBackupModelV0 struct {
	Id                           types.String `tfsdk:"id"`
	ScheduledBackups             types.Bool   `tfsdk:"scheduled_backups"`
	CreateBackups                types.Bool   `tfsdk:"create_backups"`
	BackupAppliance              types.Bool   `tfsdk:"backup_appliance"`
	DefaultBackupStorageBucketId types.Int64  `tfsdk:"default_backup_storage_bucket_id"`
	DefaultBackupScheduleId      types.Int64  `tfsdk:"default_backup_schedule_id"`
	RetentionDays                types.Int64  `tfsdk:"retention_days"`
}

// settingBackupSchemaV0 describes the SDKv2 state so the framework can decode
// it during an upgrade. Only the attribute types matter for decoding; the
// flags mirror the SDKv2 declaration. It is pinned by testdata/state_v0.json:
// the unit test fails if the two disagree.
func settingBackupSchemaV0() schema.Schema {
	optionalComputedInt := func() schema.Int64Attribute {
		return schema.Int64Attribute{Optional: true, Computed: true}
	}

	optionalComputedBool := func() schema.BoolAttribute {
		return schema.BoolAttribute{Optional: true, Computed: true}
	}

	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":                               schema.StringAttribute{Computed: true},
			"scheduled_backups":                optionalComputedBool(),
			"create_backups":                   optionalComputedBool(),
			"backup_appliance":                 optionalComputedBool(),
			"default_backup_storage_bucket_id": optionalComputedInt(),
			"default_backup_schedule_id":       optionalComputedInt(),
			"retention_days":                   optionalComputedInt(),
		},
	}
}

func (r *Resource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	priorSchema := settingBackupSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeStateV0toV1,
		},
	}
}

// upgradeStateV0toV1 converts the SDKv2 state to the current model. It keeps
// id ("1"), carries the three flags, maps an id of 0 to null (no entity has
// id 0; the legacy SDK stored an unset id as 0), and carries retention_days
// into both retention_days and retention_count, the two attributes over the
// one retentionCount setting. The refresh after the upgrade reads the
// appliance's values again.
func upgradeStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old settingBackupModelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := resourceModel{
		SettingBackupModel: SettingBackupModel{
			ScheduledBackups:             old.ScheduledBackups,
			CreateBackups:                old.CreateBackups,
			BackupAppliance:              old.BackupAppliance,
			DefaultBackupStorageBucketId: zeroToNull(old.DefaultBackupStorageBucketId),
			DefaultBackupScheduleId:      zeroToNull(old.DefaultBackupScheduleId),
			RetentionCount:               old.RetentionDays,
			RetentionDays:                old.RetentionDays,
		},
		Id: types.StringValue(singletonID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}

// zeroToNull maps the SDKv2 resource's 0 -- its unset value for an id, the
// legacy SDK's structs being non-pointer -- to null; no entity has id 0.
func zeroToNull(v types.Int64) types.Int64 {
	if v.IsNull() || v.IsUnknown() || v.ValueInt64() == 0 {
		return types.Int64Null()
	}

	return v
}
