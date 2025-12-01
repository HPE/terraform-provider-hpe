// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// populate policy resource model with current API values
func getPolicyAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan *PolicyModel,
) (PolicyModel, diag.Diagnostics) {
	var state PolicyModel
	var diags diag.Diagnostics

	policy, hresp, err := client.PoliciesAPI.GetPolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK || policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d GET failed", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	if policy.Policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d is nil", id),
		)

		return state, diags
	}

	p := policy.Policy

	// Set basic fields
	state.Id = convert.Int64ToType(p.Id)
	state.Name = convert.StrToType(p.Name)

	// Handle nullable fields properly
	if p.Description.IsSet() {
		state.Description = convert.StrToType(p.Description.Get())
	}

	state.Enabled = convert.BoolToType(p.Enabled)

	if p.EachUser.IsSet() {
		state.EachUser = convert.BoolToType(p.EachUser.Get())
	}

	// Handle RefId - convert string to int64
	if p.RefId.IsSet() && p.RefId.Get() != nil {
		state.AssociatedResourceId = types.Int64Value(p.GetRefId())
	}

	// Handle RefType
	// If RefType is null or not set, it means it's a Global policy
	if p.RefType.IsSet() && p.RefType.Get() != nil {
		apiType := *p.RefType.Get()
		// Convert API type to user-facing resource type
		state.AssociatedResourceType = types.StringValue(apiTypeToResourceType(apiType))
	} else {
		state.AssociatedResourceType = types.StringValue(AssociatedResourceTypeGlobal)
	}

	// Set Tenant IDs
	if len(p.Accounts) > 0 {
		tenantIDs := make([]int64, 0, len(p.Accounts))
		for _, acc := range p.Accounts {
			if acc.Id != nil {
				tenantIDs = append(tenantIDs, *acc.Id)
			}
		}
		tenantsSet, setDiags := types.SetValueFrom(ctx, types.Int64Type, tenantIDs)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return state, diags
		}
		state.Tenants = tenantsSet
	} else {
		// Set to null if no accounts
		state.Tenants = types.SetNull(types.Int64Type)
	}

	// Set PolicyType
	if p.PolicyType != nil {
		policyTypeAttrs := map[string]attr.Value{}
		if p.PolicyType.Code != nil {
			policyTypeAttrs["code"] = types.StringValue(*p.PolicyType.Code)
		} else {
			policyTypeAttrs["code"] = types.StringNull()
		}
		if p.PolicyType.Id != nil {
			policyTypeAttrs["id"] = types.Int64Value(*p.PolicyType.Id)
		} else {
			policyTypeAttrs["id"] = types.Int64Null()
		}
		if p.PolicyType.Name != nil {
			policyTypeAttrs["name"] = types.StringValue(*p.PolicyType.Name)
		} else {
			policyTypeAttrs["name"] = types.StringNull()
		}

		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx), policyTypeAttrs)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)

			return state, diags
		}
		state.PolicyType = policyTypeValue
	}

	// Handle Config - preserve from plan or convert from API
	state.Config = types.DynamicNull()

	if plan != nil && !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	} else if p.Config != nil {
		// Convert API config to dynamic type
		var err error
		state.Config, err = convert.StructToDynamic(ctx, p.Config)
		if err != nil {
			diags.AddError(
				"populate policy resource",
				fmt.Sprintf("policy %d: failed to convert config: %s", id, err.Error()),
			)

			return state, diags
		}
	}

	// Computed types
	// Set Cloud if present
	if p.Zone != nil {
		cloudAttrs := map[string]attr.Value{}
		if p.Zone.Id != nil {
			cloudAttrs["id"] = types.Int64Value(*p.Zone.Id)
		} else {
			cloudAttrs["id"] = types.Int64Null()
		}
		if p.Zone.Name != nil {
			cloudAttrs["name"] = types.StringValue(*p.Zone.Name)
		} else {
			cloudAttrs["name"] = types.StringNull()
		}

		cloudValue, cloudDiags := NewCloudValue(CloudValue{}.AttributeTypes(ctx), cloudAttrs)
		if cloudDiags.HasError() {
			diags.Append(cloudDiags...)

			return state, diags
		}
		state.Cloud = cloudValue
	} else {
		state.Cloud = NewCloudValueNull()
	}

	// Set Group if present
	if p.Site != nil {
		groupAttrs := map[string]attr.Value{}
		if p.Site.Id != nil {
			groupAttrs["id"] = types.Int64Value(*p.Site.Id)
		} else {
			groupAttrs["id"] = types.Int64Null()
		}
		if p.Site.Name != nil {
			groupAttrs["name"] = types.StringValue(*p.Site.Name)
		} else {
			groupAttrs["name"] = types.StringNull()
		}

		groupValue, groupDiags := NewGroupValue(GroupValue{}.AttributeTypes(ctx), groupAttrs)
		if groupDiags.HasError() {
			diags.Append(groupDiags...)

			return state, diags
		}
		state.Group = groupValue
	} else {
		state.Group = NewGroupValueNull()
	}

	// Set Owner if present
	state.Owner = NewOwnerValueNull()
	if p.Owner.IsSet() && p.Owner.Get() != nil {
		owner := p.Owner.Get()
		ownerAttrs := map[string]attr.Value{}
		if owner.Id != nil {
			ownerAttrs["id"] = types.Int64Value(*owner.Id)
		} else {
			ownerAttrs["id"] = types.Int64Null()
		}
		if owner.Name != nil {
			ownerAttrs["name"] = types.StringValue(*owner.Name)
		} else {
			ownerAttrs["name"] = types.StringNull()
		}

		ownerValue, ownerDiags := NewOwnerValue(OwnerValue{}.AttributeTypes(ctx), ownerAttrs)
		if ownerDiags.HasError() {
			diags.Append(ownerDiags...)

			return state, diags
		}
		state.Owner = ownerValue
	}

	// Set Role if present
	if p.Role != nil {
		roleAttrs := map[string]attr.Value{}
		if p.Role.Id != nil {
			roleAttrs["id"] = types.Int64Value(*p.Role.Id)
		} else {
			roleAttrs["id"] = types.Int64Null()
		}
		if p.Role.Authority != nil {
			roleAttrs["authority"] = types.StringValue(*p.Role.Authority)
		} else {
			roleAttrs["authority"] = types.StringNull()
		}

		roleValue, roleDiags := NewRoleValue(RoleValue{}.AttributeTypes(ctx), roleAttrs)
		if roleDiags.HasError() {
			diags.Append(roleDiags...)

			return state, diags
		}
		state.Role = roleValue
	} else {
		state.Role = NewRoleValueNull()
	}

	// Set User if present
	if p.User != nil {
		userAttrs := map[string]attr.Value{}
		if p.User.Id != nil {
			userAttrs["id"] = types.Int64Value(*p.User.Id)
		} else {
			userAttrs["id"] = types.Int64Null()
		}
		if p.User.Username != nil {
			userAttrs["username"] = types.StringValue(*p.User.Username)
		} else {
			userAttrs["username"] = types.StringNull()
		}

		userValue, userDiags := NewUserValue(UserValue{}.AttributeTypes(ctx), userAttrs)
		if userDiags.HasError() {
			diags.Append(userDiags...)

			return state, diags
		}
		state.User = userValue
	} else {
		state.User = NewUserValueNull()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan PolicyModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read policy resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
