// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &deviceManagementPolicyDataSource{}

func NewDeviceManagementPolicyDataSource() datasource.DataSource {
	return &deviceManagementPolicyDataSource{}
}

type deviceManagementPolicyDataSource struct {
	BaseData
}

type deviceManagementPolicyDataSourceModel struct {
	Client          types.String `tfsdk:"client"`
	Name            types.String `tfsdk:"name"`
	Id              types.Int64  `tfsdk:"id"`
	UID             types.String `tfsdk:"uid"`
	ResourceType    types.String `tfsdk:"resource_type"`
	MatchType       types.String `tfsdk:"match_type"`
	ExecutionOrder  types.Int64  `tfsdk:"execution_order"`
	Customized      types.Bool   `tfsdk:"customized"`
	GlobalPolicy    types.Bool   `tfsdk:"global_policy"`
	ParentPolicy    types.Int64  `tfsdk:"parent_policy"`
	ParentPolicyUID types.String `tfsdk:"parent_policy_uid"`
}

func (d *deviceManagementPolicyDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_device_management_policy"
}

func (d *deviceManagementPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an OpsRamp Device Management Policy by exact name and returns its identifiers and summary attributes.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:    true,
				Description: "Optional client (tenant) UUID to query against. Defaults to the provider tenant.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The device management policy name to look up.",
			},
			"parent_policy": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric `id` of the parent policy to clone this policy from.",
			},
			"parent_policy_uid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The `uid` of the parent policy, when this policy was cloned.",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The numeric identifier of the policy.",
			},
			"uid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (e.g. `POLICY-xxxxxxxx`) of the policy, usable with `hpe_opsramp_device_management_policy` import.",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of resource the policy applies to.",
			},
			"match_type": schema.StringAttribute{
				Computed:    true,
				Description: "How the policy's filter rules are combined (`ALL` or `ANY`).",
			},
			"execution_order": schema.Int64Attribute{
				Computed:    true,
				Description: "The execution order of the policy.",
			},
			"customized": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this policy has been customized (cloned) from a parent policy.",
			},
			"global_policy": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is a global policy rather than a tenant-specific policy.",
			},
		},
	}
}

func (d *deviceManagementPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.apiClient == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "Expected an authenticated API client from provider.Configure()")

		return
	}

	var data deviceManagementPolicyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := resolveLookupTenantID(d.apiClient, data.Client)

	policy, err := d.apiClient.FindDeviceManagementPolicyByName(tenantID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Device management policy lookup failed",
			fmt.Sprintf("Could not find device management policy %q: %s", data.Name.ValueString(), err.Error()),
		)

		return
	}

	data.Id = types.Int64Value(int64(policy.Id))
	data.UID = types.StringValue(policy.UID)
	data.ResourceType = types.StringValue(policy.ResourceType)
	data.MatchType = types.StringValue(policy.MatchType)
	data.ExecutionOrder = types.Int64Value(int64(policy.ExecutionOrder))
	data.Customized = types.BoolValue(policy.Customized)
	data.GlobalPolicy = types.BoolValue(policy.GlobalPolicy)
	data.ParentPolicy = types.Int64Value(int64(policy.ParentPolicy))
	data.ParentPolicyUID = types.StringValue(policy.ParentPolicyUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
