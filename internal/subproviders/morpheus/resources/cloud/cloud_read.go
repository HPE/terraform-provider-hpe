// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

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

const standardCloud = "standard"

// populate cloud resource model with current API values
func getCloudAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan CloudModel,
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

	if !cloud.IsSetConfig() {
		diags.AddError(
			"populate cloud resource",
			fmt.Sprintf("cloud %d missing config", id),
		)

		return state, diags
	}

	if len(cloud.Groups) == 0 {
		diags.AddError(
			"populate cloud resource",
			fmt.Sprintf("cloud %d no associated groups", id),
		)

		return state, diags
	}

	importing := false
	if plan.Name.IsNull() {
		importing = true
	}

	state.GroupId = convert.Int64ToType(cloud.Groups[0].Id)
	state.AgentInstallMode = convert.StrToType(cloud.AgentMode)
	state.AutoRecoverPowerState = convert.BoolToType(cloud.AutoRecoverPowerState)
	if *cloud.Code != "standard" { // workaround an API bug
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

	cfg := cloud.GetConfig()

	// Move these common fields up
	state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
	cfg.ApplianceUrl.Unset()
	state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
	cfg.DatacenterName.Unset()
	state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
	cfg.ExternalId.Unset()
	state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
	cfg.InventoryLevel.Unset()
	state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())
	cfg.ConsoleKeymap.Unset()

	// Remove possibly buggy API fields
	cfg.ConfigCmdbDiscovery = nil

	switch {
	case importing && *cloud.ZoneType.Code == standardCloud,
		!plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
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

				return state, diags
			}

			state.ConfigHvm = configHvm
		}
	default:
		state.CloudTypeCode = convert.StrToType(cloud.ZoneType.Code)

		state.Config = types.DynamicNull()

		if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
			state.Config = plan.Config
		} else {
			state.Config, err = convert.StructToDynamic(ctx, cfg)
			if err != nil {
				diags.AddError(
					"create cloud resource",
					"cloud: failed to convert config: "+err.Error(),
				)

				return state, diags
			}
		}
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
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

	state, pdiags := getCloudAsState(ctx, id, client, plan)
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
