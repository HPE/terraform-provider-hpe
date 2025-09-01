// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// populate cloud resource model with current API values
func getCloudAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (CloudModel, diag.Diagnostics) {
	var state CloudModel
	var diags diag.Diagnostics

	c, hresp, err := client.CloudsAPI.GetClouds(ctx, id).Execute()
	if err != nil || c == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate cloud resource",
			fmt.Sprintf("cloud %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	cloud := c.GetZone()

	if len(cloud.Groups) == 0 {
		diags.AddError(
			"populate cloud resource",
			fmt.Sprintf("cloud %d no associated groups", id),
		)

		return state, diags
	}

	state.GroupId = convert.Int64ToType(cloud.Groups[0].Id)

	cfg := cloud.GetConfig()

	switch cloud.ZoneType.GetCode() {
	case "standard":
		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		if cfg.CertificateProvider.Get() != nil {
			certificateProvider := convert.StrToType(cfg.CertificateProvider.Get())
			attrTypes["certificate_provider"] = types.StringType
			attrValues["certificate_provider"] = certificateProvider
		}

		if cfg.EnableNetworkTypeSelection != nil {
			enableNetworkTypeSelection := convert.BoolToType(cfg.EnableNetworkTypeSelection)
			attrTypes["enable_network_type_selection"] = types.BoolType
			attrValues["enable_network_type_selection"] = enableNetworkTypeSelection
		}

		if len(attrValues) > 0 {
			configHvm, diagsHvm := NewConfigHvmValue(attrTypes, attrValues)
			if diagsHvm.HasError() {
				diags.Append(diagsHvm...)
				diags.AddError(
					"populate cloud resource",
					fmt.Sprintf("cloud %d: failed to decode HVM configuration", id),
				)

				return CloudModel{}, diags
			}

			state.ConfigHvm = configHvm
		}
	default:
		diags.AddError(
			"populate cloud resource",
			fmt.Sprintf("cloud %d: unsupported cloud type configuration", id),
		)

		return CloudModel{}, diags
	}

	state.AgentInstallMode = convert.StrToType(cloud.AgentMode)
	state.AutoRecoverPowerState = convert.BoolToType(cloud.AutoRecoverPowerState)
	if *cloud.Code != "standard" {
		state.Code = convert.StrToType(cloud.Code)
	}
	state.CostingMode = convert.StrToType(cloud.CostingMode.Get())
	state.Enabled = convert.BoolToType(cloud.Enabled)
	state.GuidanceMode = convert.StrToType(cloud.GuidanceMode.Get())
	state.Id = convert.Int64ToType(cloud.Id)
	state.Labels = convert.StrSliceToSet(cloud.Labels)
	state.Location = convert.StrToType(cloud.Location.Get())
	state.Name = convert.StrToType(cloud.Name)
	state.SecurityMode = convert.StrToType(cloud.SecurityMode)
	state.TenantId = convert.Int64ToType(cloud.AccountId)
	state.Visibility = convert.StrToType(cloud.Visibility)

	state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
	state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
	state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
	state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
	state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())

	return state, diags
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan CloudModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read cloud resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, pdiags := getCloudAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read cloud resource",
			fmt.Sprintf("cloud %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
