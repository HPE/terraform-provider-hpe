// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// licenseModelV0 is the state shape produced by the retired
// terraform-plugin-sdk/v2 license resource (schema version 0). Its defining
// difference from the current model is that "key" was a plain Sensitive string
// stored IN STATE; it is now write-only and never stored. The other attributes
// carry over unchanged.
type licenseModelV0 struct {
	Id           types.String `tfsdk:"id"`
	Key          types.String `tfsdk:"key"`
	ProductTier  types.String `tfsdk:"product_tier"`
	AccountName  types.String `tfsdk:"account_name"`
	StartDate    types.String `tfsdk:"start_date"`
	EndDate      types.String `tfsdk:"end_date"`
	MaxInstances types.Int64  `tfsdk:"max_instances"`
	MaxMemory    types.Int64  `tfsdk:"max_memory"`
	MaxStorage   types.Int64  `tfsdk:"max_storage"`
	HardLimit    types.Bool   `tfsdk:"hard_limit"`
	FreeTrial    types.Bool   `tfsdk:"free_trial"`
	MultiTenant  types.Bool   `tfsdk:"multi_tenant"`
	SupportLevel types.String `tfsdk:"support_level"`
}

// licenseSchemaV0 describes the SDKv2 state shape so the framework can decode
// prior state during an upgrade. Only the attribute types matter for decoding;
// the flags mirror the SDKv2 declaration for clarity. Note "key" is a plain
// Sensitive string here, not write-only.
func licenseSchemaV0() schema.Schema {
	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"key":           schema.StringAttribute{Required: true, Sensitive: true},
			"product_tier":  schema.StringAttribute{Computed: true},
			"account_name":  schema.StringAttribute{Computed: true},
			"start_date":    schema.StringAttribute{Computed: true},
			"end_date":      schema.StringAttribute{Computed: true},
			"max_instances": schema.Int64Attribute{Computed: true},
			"max_memory":    schema.Int64Attribute{Computed: true},
			"max_storage":   schema.Int64Attribute{Computed: true},
			"hard_limit":    schema.BoolAttribute{Computed: true},
			"free_trial":    schema.BoolAttribute{Computed: true},
			"multi_tenant":  schema.BoolAttribute{Computed: true},
			"support_level": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *Resource) UpgradeState(
	_ context.Context,
) map[int64]resource.StateUpgrader {
	priorSchema := licenseSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeLicenseStateV0toV1,
		},
	}
}

// upgradeLicenseStateV0toV1 converts SDKv2-produced state (plaintext key) to
// the current framework model. The framework's same-version passthrough does
// NOT nullify write-only attributes, so this upgrader is what drops the
// plaintext key from state (Key is set to null). Every other attribute carries
// over as-is -- the SDKv2 resource wrote ""/0/false for a licence that was not
// installed; the first refresh replaces them or removes the resource. The
// write-only digest is not written here: the first Read seeds it from the
// API's hash.
//
// Legacy state without an id is refused rather than carried forward: Read would
// otherwise match nothing meaningfully. SDKv2 removed a resource from state
// when it cleared its id, so such state is corrupt and the practitioner must
// remove it and re-import.
func upgradeLicenseStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old licenseModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if old.Id.IsNull() || old.Id.IsUnknown() || old.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"upgrade license resource state",
			"the legacy license state has no id, so it cannot be upgraded; remove it with "+
				"'terraform state rm' and import the licence again",
		)

		return
	}

	upgraded := LicenseModel{
		Id:           old.Id,
		Key:          types.StringNull(),
		ProductTier:  old.ProductTier,
		AccountName:  old.AccountName,
		StartDate:    old.StartDate,
		EndDate:      old.EndDate,
		MaxInstances: old.MaxInstances,
		MaxMemory:    old.MaxMemory,
		MaxStorage:   old.MaxStorage,
		HardLimit:    old.HardLimit,
		FreeTrial:    old.FreeTrial,
		MultiTenant:  old.MultiTenant,
		SupportLevel: old.SupportLevel,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}
