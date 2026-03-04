// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	standardCloud = "standard"
	vmwareCloud   = "vmware"
	readOperation = "read cloud resource"
)

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
			readOperation,
			fmt.Sprintf("cloud %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	cloud := c.GetZone()

	if !cloud.IsSetConfig() {
		diags.AddError(
			readOperation,
			fmt.Sprintf("cloud %d missing config", id),
		)

		return state, diags
	}

	if len(cloud.Groups) == 0 {
		diags.AddError(
			readOperation,
			fmt.Sprintf("cloud %d no associated groups", id),
		)

		return state, diags
	}

	importing := plan.Name.IsNull()
	cloudType := *cloud.ZoneType.Code

	state.GroupId = convert.Int64ToType(cloud.Groups[0].Id)
	state.AgentInstallMode = convert.StrToType(cloud.AgentMode)
	state.AutoRecoverPowerState = convert.BoolToType(cloud.AutoRecoverPowerState)
	if cloud.GetCode() != "standard" { // workaround an API bug
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

	switch {
	case cloudType == standardCloud && (!plan.ConfigHvm.IsNull() || importing):
		cfg := cloud.GetConfig().GetClouds200ResponseZoneConfigAnyOf

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
		state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
		state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
		state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
		state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())

		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		if cfg.CertificateProvider.Get() != nil {
			attrTypes["certificate_provider"] = types.StringType
			attrValues["certificate_provider"] = convert.StrToType(cfg.CertificateProvider.Get())
		}

		if cfg.EnableNetworkTypeSelection.Get() != nil {
			attrTypes["enable_network_type_selection"] = types.BoolType
			attrValues["enable_network_type_selection"] = convert.StringToBool(ctx, *cfg.EnableNetworkTypeSelection.Get())
		}

		if len(attrValues) > 0 {
			configHvm, diagsHvm := NewConfigHvmValue(attrTypes, attrValues)
			if diagsHvm.HasError() {
				diags.Append(diagsHvm...)
				diags.AddError(
					readOperation,
					fmt.Sprintf("cloud %d: failed to decode HVM configuration", id),
				)

				return state, diags
			}

			state.ConfigHvm = configHvm
		}

	case cloudType == vmwareCloud && (!plan.ConfigVmware.IsNull() || importing):
		cfg := cloud.GetConfig().GetClouds200ResponseZoneConfigAnyOf1

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
		state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
		state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
		state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
		state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())

		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		if cfg.ApiUrl != nil {
			attrTypes["api_url"] = types.StringType
			attrValues["api_url"] = convert.StrToType(cfg.ApiUrl)
		}

		if cfg.ApiVersion.Get() != nil {
			attrTypes["api_version"] = types.StringType
			attrValues["api_version"] = convert.StrToType(cfg.ApiVersion.Get())
		}

		if cfg.CertificateProvider.Get() != nil {
			attrTypes["certificate_provider"] = types.StringType
			attrValues["certificate_provider"] = convert.StrToType(cfg.CertificateProvider.Get())
		}

		if cfg.Cluster != nil {
			attrTypes["cluster"] = types.StringType
			attrValues["cluster"] = convert.StrToType(cfg.Cluster)
		}

		if cfg.ConfigManagementId.Get() != nil {
			attrTypes["config_management_id"] = types.StringType
			attrValues["config_management_id"] = convert.StrToType(cfg.ConfigManagementId.Get())
		}

		if cfg.Datacenter != nil {
			attrTypes["datacenter"] = types.StringType
			attrValues["datacenter"] = convert.StrToType(cfg.Datacenter)
		}

		if cfg.EnableDiskTypeSelection.Get() != nil {
			attrTypes["enable_disk_type_selection"] = types.BoolType
			attrValues["enable_disk_type_selection"] = convert.StringToBool(ctx, *cfg.EnableDiskTypeSelection.Get())
		}

		if cfg.EnableNetworkTypeSelection.Get() != nil {
			attrTypes["enable_network_type_selection"] = types.BoolType
			attrValues["enable_network_type_selection"] = convert.StringToBool(ctx, *cfg.EnableNetworkTypeSelection.Get())
		}

		if cfg.EnableStorageTypeSelection.Get() != nil {
			attrTypes["enable_storage_type_selection"] = types.BoolType
			attrValues["enable_storage_type_selection"] = convert.StringToBool(ctx, *cfg.EnableStorageTypeSelection.Get())
		}

		if cfg.EnableVnc.Get() != nil {
			attrTypes["enable_vnc"] = types.BoolType
			attrValues["enable_vnc"] = convert.StringToBool(ctx, *cfg.EnableVnc.Get())
		}

		if cfg.HideHostSelection.Get() != nil {
			attrTypes["hide_host_selection"] = types.BoolType
			attrValues["hide_host_selection"] = convert.StringToBool(ctx, *cfg.HideHostSelection.Get())
		} else {
			attrTypes["hide_host_selection"] = types.BoolType
			attrValues["hide_host_selection"] = types.BoolNull()
		}

		if cfg.Password.Get() != nil {
			attrTypes["password"] = types.StringType
			attrValues["password"] = convert.StrToType(cfg.Password.Get())
		} else {
			attrTypes["password"] = types.StringType
			attrValues["password"] = types.StringNull()
		}

		if cfg.ResourcePool != nil {
			attrTypes["resource_pool"] = types.StringType
			attrValues["resource_pool"] = convert.StrToType(cfg.ResourcePool)
		}

		if cfg.RpcMode != nil {
			attrTypes["rpc_mode"] = types.StringType
			attrValues["rpc_mode"] = convert.StrToType(cfg.RpcMode)
		}

		if cfg.Username != nil {
			attrTypes["username"] = types.StringType
			attrValues["username"] = convert.StrToType(cfg.Username)
		}

		if len(attrValues) > 0 {
			configVmware, diagsHvm := NewConfigVmwareValue(attrTypes, attrValues)
			if diagsHvm.HasError() {
				diags.Append(diagsHvm...)
				diags.AddError(
					readOperation,
					fmt.Sprintf("cloud %d: failed to decode VMware configuration", id),
				)

				return state, diags
			}

			state.ConfigVmware = configVmware
		}
	case !plan.Config.IsNull() || importing:
		state.CloudTypeCode = convert.StrToType(cloud.ZoneType.Code)

		state.Config = types.DynamicNull()

		cfg := cloud.GetConfig().MapmapOfStringAny
		if cfg == nil {
			diags.AddError(
				readOperation,
				"cloud: generic config missing from API response",
			)

			return state, diags
		}

		cfgValue := *cfg

		// Move these common fields up
		if v, ok := cfgValue["applianceUrl"].(string); ok {
			state.ApplianceUrl = convert.StrToType(&v)
			delete(cfgValue, "applianceUrl")
		}
		if v, ok := cfgValue["datacenterName"].(string); ok {
			state.DataCenterName = convert.StrToType(&v)
			delete(cfgValue, "datacenterName")
		}
		if v, ok := cfgValue["externalId"].(string); ok {
			state.ExternalId = convert.StrToType(&v)
			delete(cfgValue, "externalId")
		}
		if v, ok := cfgValue["inventoryLevel"].(string); ok {
			state.ImportExistingVms = convert.StrToType(&v)
			delete(cfgValue, "inventoryLevel")
		}
		if v, ok := cfgValue["consoleKeymap"].(string); ok {
			state.KeyboardLayout = convert.StrToType(&v)
			delete(cfgValue, "consoleKeymap")
		}

		if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
			state.Config = plan.Config
		} else {
			cfg := cloud.GetConfig().MapmapOfStringAny

			cfgValue := *cfg

			// Move these common fields up
			if v, ok := cfgValue["applianceUrl"].(string); ok {
				state.ApplianceUrl = convert.StrToType(&v)
				delete(cfgValue, "applianceUrl")
			}
			if v, ok := cfgValue["datacenterName"].(string); ok {
				state.DataCenterName = convert.StrToType(&v)
				delete(cfgValue, "datacenterName")
			}
			if v, ok := cfgValue["externalId"].(string); ok {
				state.ExternalId = convert.StrToType(&v)
				delete(cfgValue, "externalId")
			}
			if v, ok := cfgValue["inventoryLevel"].(string); ok {
				state.ImportExistingVms = convert.StrToType(&v)
				delete(cfgValue, "inventoryLevel")
			}
			if v, ok := cfgValue["consoleKeymap"].(string); ok {
				state.KeyboardLayout = convert.StrToType(&v)
				delete(cfgValue, "consoleKeymap")
			}

			state.Config, err = convert.MapToDynamic(ctx, cfgValue)
			if err != nil {
				diags.AddError(
					readOperation,
					"cloud: failed to convert generic config: "+err.Error(),
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
			readOperation,
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, pdiags := getCloudAsState(ctx, id, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			readOperation,
			fmt.Sprintf("cloud %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
