// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/policy/consts"
	internalErrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

const summary = "read policy data source"

// apiTypeToResourceType converts API types back to user-facing resource types
func apiTypeToResourceType(apiType string) string {
	switch apiType {
	case "ComputeZone":
		return "Cloud"
	case "ComputeSite":
		return "Group"
	default:
		// For other types (User, Role, Network, Plan), pass through as-is
		return apiType
	}
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_policy"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = PolicyDataSourceSchema(ctx)
}

func getPolicyByName(
	ctx context.Context,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	name := data.Name.ValueString()
	ps, hresp, err := apiClient.PoliciesAPI.ListPolicies(ctx).Name(name).Execute()
	if ps == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for policy with name %s: %s",
			name, internalErrors.ErrMsg(err, hresp)))

		return diags
	}

	policies := ps.GetPolicies()

	// Additional filtering to ensure exact name match (API might return partial matches)
	var filteredPolicies []sdk.ListPolicies200ResponseAllOfPoliciesInner
	for _, p := range policies {
		if p.GetName() == data.Name.ValueString() {
			filteredPolicies = append(filteredPolicies, p)
		}
	}
	policies = filteredPolicies

	if len(policies) > 1 {
		diags.AddError(summary, consts.ErrorMultiplePolicies)

		return diags
	} else if len(policies) == 0 {
		diags.AddError(summary, consts.ErrorNoPolicyFound)

		return diags
	}

	policy := policies[0]

	return getPolicyByID(ctx, policy.GetId(), data, apiClient)
}

func getPolicyByID(
	ctx context.Context,
	id int64,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	p, hresp, err := apiClient.PoliciesAPI.GetPolicies(ctx, id).Execute()
	if p == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for policy with id %d: %s",
			id, internalErrors.ErrMsg(err, hresp)))

		return diags
	}
	policy, ok := p.GetPolicyOk()
	if !ok {
		diags.AddError(summary, consts.ErrorNoPolicyFound)

		return diags
	}

	// Map basic policy fields
	data.Id = convert.Int64ToType(policy.Id)
	data.Name = convert.StrToType(policy.Name)
	data.Description = convert.StrToType(policy.Description.Get())
	data.Enabled = convert.BoolToType(policy.Enabled)
	data.EachUser = convert.BoolToType(policy.EachUser.Get())

	// Handle AssociatedResourceId and AssociatedResourceType
	if policy.RefId.IsSet() && policy.RefId.Get() != nil {
		data.AssociatedResourceId = types.Int64Value(policy.GetRefId())
	} else {
		data.AssociatedResourceId = types.Int64Null()
	}

	if policy.RefType.IsSet() && policy.RefType.Get() != nil {
		apiType := *policy.RefType.Get()
		data.AssociatedResourceType = types.StringValue(apiTypeToResourceType(apiType))
	} else {
		data.AssociatedResourceType = types.StringValue("Global")
	}

	// Handle PolicyType
	if policy.PolicyType != nil {
		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.PolicyType.Id),
				"code": convert.StrToType(policy.PolicyType.Code),
				"name": convert.StrToType(policy.PolicyType.Name),
			},
		)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)
			return diags
		}
		data.PolicyType = policyTypeValue
	} else {
		data.PolicyType = NewPolicyTypeValueNull()
	}

	// Handle Config - for now set to null since it's complex nested structure
	// This would require a detailed mapping of all the different config types
	data.Config = NewConfigValueNull()

	// Handle Cloud (Zone)
	if policy.Zone != nil {
		cloudValue, cloudDiags := NewCloudValue(
			CloudValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.Zone.Id),
				"name": convert.StrToType(policy.Zone.Name),
			},
		)
		if cloudDiags.HasError() {
			diags.Append(cloudDiags...)
			return diags
		}
		data.Cloud = cloudValue
	} else {
		data.Cloud = NewCloudValueNull()
	}

	// Handle Group (Site)
	if policy.Site != nil {
		groupValue, groupDiags := NewGroupValue(
			GroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.Site.Id),
				"name": convert.StrToType(policy.Site.Name),
			},
		)
		if groupDiags.HasError() {
			diags.Append(groupDiags...)
			return diags
		}
		data.Group = groupValue
	} else {
		data.Group = NewGroupValueNull()
	}

	// Handle Owner
	if policy.Owner.IsSet() && policy.Owner.Get() != nil {
		owner := policy.Owner.Get()
		ownerValue, ownerDiags := NewOwnerValue(
			OwnerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(owner.Id),
				"name": convert.StrToType(owner.Name),
			},
		)
		if ownerDiags.HasError() {
			diags.Append(ownerDiags...)
			return diags
		}
		data.Owner = ownerValue
	} else {
		data.Owner = NewOwnerValueNull()
	}

	// Handle Role
	if policy.Role != nil {
		roleValue, roleDiags := NewRoleValue(
			RoleValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":        convert.Int64ToType(policy.Role.Id),
				"authority": convert.StrToType(policy.Role.Authority),
			},
		)
		if roleDiags.HasError() {
			diags.Append(roleDiags...)
			return diags
		}
		data.Role = roleValue
	} else {
		data.Role = NewRoleValueNull()
	}

	// Handle User
	if policy.User != nil {
		userValue, userDiags := NewUserValue(
			UserValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":       convert.Int64ToType(policy.User.Id),
				"username": convert.StrToType(policy.User.Username),
			},
		)
		if userDiags.HasError() {
			diags.Append(userDiags...)
			return diags
		}
		data.User = userValue
	} else {
		data.User = NewUserValueNull()
	}

	// Handle Tenants (Accounts)
	if len(policy.Accounts) > 0 {
		tenantValues := []attr.Value{}
		for _, account := range policy.Accounts {
			tenantValue, tenantDiags := types.ObjectValueFrom(ctx,
				map[string]attr.Type{
					"id":   types.Int64Type,
					"name": types.StringType,
				},
				map[string]attr.Value{
					"id":   convert.Int64ToType(account.Id),
					"name": convert.StrToType(account.Name),
				},
			)
			if tenantDiags.HasError() {
				diags.Append(tenantDiags...)
				return diags
			}
			tenantValues = append(tenantValues, tenantValue)
		}

		tenantsSet, setDiags := types.SetValueFrom(ctx,
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   types.Int64Type,
					"name": types.StringType,
				},
			},
			tenantValues,
		)
		if setDiags.HasError() {
			diags.Append(setDiags...)
			return diags
		}
		data.Tenants = tenantsSet
	} else {
		data.Tenants = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			},
		})
	}

	return diags
}

func getPolicy(
	ctx context.Context,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	if !data.Id.IsNull() {
		return getPolicyByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Name.IsNull() {
		return getPolicyByName(ctx, data, apiClient)
	}

	diags := diag.Diagnostics{}
	diags.AddError(summary, consts.ErrorNoValidPolicyTerms)

	return diags
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data PolicyModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	diags = getPolicy(ctx, &data, apiClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
