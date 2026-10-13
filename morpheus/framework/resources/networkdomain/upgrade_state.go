// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// networkDomainModelV0 is the state shape produced by the retired
// terraform-plugin-sdk/v2 network domain resource (schema version 0). Its
// defining differences from the current model are that "id" is a string (SDKv2
// always stored resource IDs as strings). Everything else is carried unchanged,
// including the deprecated auto_join_domain, domain_password and the
// domain_password_wo_version marker (all still accepted until 2.6). The one
// attribute not carried is domain_password_wo: a write-only value is never kept
// in state (it is always null in a version-0 state anyway).
type networkDomainModelV0 struct {
	Active                  types.Bool   `tfsdk:"active"`
	AutoJoinDomain          types.Bool   `tfsdk:"auto_join_domain"`
	Description             types.String `tfsdk:"description"`
	DomainController        types.Bool   `tfsdk:"domain_controller"`
	DomainPassword          types.String `tfsdk:"domain_password"`
	DomainPasswordWo        types.String `tfsdk:"domain_password_wo"`
	DomainPasswordWoVersion types.Int64  `tfsdk:"domain_password_wo_version"`
	DomainUsername          types.String `tfsdk:"domain_username"`
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	PublicZone              types.Bool   `tfsdk:"public_zone"`
	TenantId                types.Int64  `tfsdk:"tenant_id"`
	Visibility              types.String `tfsdk:"visibility"`
}

// networkDomainSchemaV0 describes the SDKv2 state shape so the framework can
// decode prior state during an upgrade. Only the attribute types matter for
// decoding; the optional/computed flags mirror the SDKv2 declaration for
// clarity. It is pinned by testdata/state_v0.json (a real SDKv2 state): the
// unit test fails if this shape and that state disagree.
func networkDomainSchemaV0() schema.Schema {
	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"active":                     schema.BoolAttribute{Optional: true, Computed: true},
			"auto_join_domain":           schema.BoolAttribute{Optional: true},
			"description":                schema.StringAttribute{Optional: true},
			"domain_controller":          schema.BoolAttribute{Optional: true, Computed: true},
			"domain_password":            schema.StringAttribute{Optional: true, Sensitive: true},
			"domain_password_wo":         schema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
			"domain_password_wo_version": schema.Int64Attribute{Optional: true},
			"domain_username":            schema.StringAttribute{Optional: true},
			"id":                         schema.StringAttribute{Computed: true},
			"name":                       schema.StringAttribute{Required: true},
			"public_zone":                schema.BoolAttribute{Optional: true, Computed: true},
			"tenant_id":                  schema.Int64Attribute{Optional: true, Computed: true},
			"visibility":                 schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *Resource) UpgradeState(
	_ context.Context,
) map[int64]resource.StateUpgrader {
	priorSchema := networkDomainSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeNetworkDomainStateV0toV1,
		},
	}
}

// upgradeNetworkDomainStateV0toV1 converts SDKv2-produced state (string id) to
// the current framework model (int64 id). It parses the id, drops
// domain_password_wo (write-only, never in state), keeps the deprecated
// attributes and the marker, and carries the rest.
//
// write_only_revision is set null. The resource has no private state after an
// upgrade, so wodigest.Tracker.ModifyPlan takes its stored-digest-missing path
// (Load returns ok=false → DecideConfigured returns Seed), which marks
// write_only_revision unknown. The first plan after the upgrade therefore shows
// exactly one "~ write_only_revision" in-place update that re-sends the
// configured password and records its digest; every plan after that is clean.
func upgradeNetworkDomainStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old networkDomainModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if old.Id.IsNull() || old.Id.IsUnknown() || old.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"upgrade network_domain resource state",
			"the legacy network domain state has no id, so it cannot be upgraded; remove it "+
				"with 'terraform state rm' and import the network domain again",
		)

		return
	}

	parsed, err := strconv.ParseInt(old.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"upgrade network_domain resource state",
			"could not convert legacy string id "+old.Id.ValueString()+
				" to an integer: "+err.Error(),
		)

		return
	}

	upgraded := NetworkDomainModel{
		Id:                      types.Int64Value(parsed),
		Name:                    old.Name,
		Description:             old.Description,
		Active:                  old.Active,
		DomainController:        old.DomainController,
		DomainUsername:          old.DomainUsername,
		PublicZone:              old.PublicZone,
		TenantId:                old.TenantId,
		Visibility:              old.Visibility,
		DomainPasswordWoVersion: old.DomainPasswordWoVersion,
		// Deprecated but still accepted until 2.6; carried as they were.
		DomainPassword: old.DomainPassword,
		AutoJoinDomain: old.AutoJoinDomain,
		// Write-only: never in state.
		DomainPasswordWo: types.StringNull(),
		// Null so the first plan after upgrade seeds the digest; see above.
		WriteOnlyRevision: types.Int64Null(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}
