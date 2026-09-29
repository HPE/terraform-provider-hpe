// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// tenantModelV0 is the state shape produced by the retired
// terraform-plugin-sdk/v2 tenant resource (schema version 0). Its defining
// difference from the current model is that "id" is a string; SDKv2 always
// stored the resource ID as a string, whereas every framework schema uses
// int64. The other attributes are a subset of the current model and carry over
// unchanged.
type tenantModelV0 struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Subdomain      types.String `tfsdk:"subdomain"`
	BaseRoleId     types.Int64  `tfsdk:"base_role_id"`
	Currency       types.String `tfsdk:"currency"`
	AccountNumber  types.String `tfsdk:"account_number"`
	AccountName    types.String `tfsdk:"account_name"`
	CustomerNumber types.String `tfsdk:"customer_number"`
}

// tenantSchemaV0 describes the SDKv2 state shape so the framework can decode
// prior state during an upgrade. Only the attribute types matter for decoding;
// the optional/computed flags mirror the SDKv2 declaration for clarity.
func tenantSchemaV0() schema.Schema {
	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true},
			"name":            schema.StringAttribute{Required: true},
			"description":     schema.StringAttribute{Optional: true, Computed: true},
			"enabled":         schema.BoolAttribute{Optional: true},
			"subdomain":       schema.StringAttribute{Optional: true, Computed: true},
			"base_role_id":    schema.Int64Attribute{Required: true},
			"currency":        schema.StringAttribute{Optional: true},
			"account_number":  schema.StringAttribute{Optional: true, Computed: true},
			"account_name":    schema.StringAttribute{Optional: true, Computed: true},
			"customer_number": schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *Resource) UpgradeState(
	_ context.Context,
) map[int64]resource.StateUpgrader {
	priorSchema := tenantSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeTenantStateV0toV1,
		},
	}
}

// upgradeTenantStateV0toV1 converts SDKv2-produced state (string id) to the
// current framework model (int64 id). Attributes introduced by the framework
// port are set to null; the refresh that follows the upgrade repopulates the
// computed ones from the API.
//
// A legacy state without an id is refused rather than carried forward: Read
// would otherwise fetch tenant 0. SDKv2 removed a resource from state when it
// cleared its id, so such state is corrupt and the practitioner must remove it
// and re-import.
func upgradeTenantStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old tenantModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if old.Id.IsNull() || old.Id.IsUnknown() || old.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"upgrade tenant resource state",
			"the legacy tenant state has no id, so it cannot be upgraded; remove it with "+
				"'terraform state rm' and import the tenant again",
		)

		return
	}

	parsed, err := strconv.ParseInt(old.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"upgrade tenant resource state",
			"could not convert legacy string id "+old.Id.ValueString()+
				" to an integer: "+err.Error(),
		)

		return
	}

	upgraded := TenantModel{
		Id:             types.Int64Value(parsed),
		Name:           old.Name,
		Description:    old.Description,
		Enabled:        old.Enabled,
		Subdomain:      old.Subdomain,
		BaseRoleId:     old.BaseRoleId,
		Currency:       old.Currency,
		AccountNumber:  old.AccountNumber,
		AccountName:    old.AccountName,
		CustomerNumber: old.CustomerNumber,
		// Attributes introduced by the framework port; refreshed on next Read.
		BaseRoleName:    types.StringNull(),
		ParentName:      types.StringNull(),
		ParentSubdomain: types.StringNull(),
		ExternalId:      types.StringNull(),
		Master:          types.BoolNull(),
		InstanceCount:   types.Int64Null(),
		UserCount:       types.Int64Null(),
		DateCreated:     types.StringNull(),
		LastUpdated:     types.StringNull(),
		ParentId:        types.Int64Null(),
		RemoveResources: types.BoolNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}
