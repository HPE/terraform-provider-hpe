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
	awsCloud      = "amazon"
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
	case cloudType == awsCloud && (!plan.ConfigAws.IsNull() || importing):
		cfg := cloud.GetConfig().GetClouds200ResponseZoneConfigAnyOf2

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl)
		state.DataCenterName = convert.StrToType(cfg.DatacenterName)

		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		if cfg.AccessKey != nil {
			attrTypes["access_key"] = types.StringType
			attrValues["access_key"] = convert.StrToType(cfg.AccessKey)
		}

		if cfg.BackupMode != nil {
			attrTypes["backup_provider"] = types.StringType
			attrValues["backup_provider"] = convert.StrToType(cfg.BackupMode)
		}

		if cfg.ConfigCmdbDiscovery != nil {
			attrTypes["cmdb_discovery"] = types.BoolType
			attrValues["cmdb_discovery"] = convert.BoolToType(cfg.ConfigCmdbDiscovery)
		}

		if cfg.ConfigManagementId != nil {
			attrTypes["config_management_id"] = types.StringType
			attrValues["config_management_id"] = convert.StrToType(cfg.ConfigManagementId)
		}

		if cfg.CostingBucket != nil {
			attrTypes["costing_bucket"] = types.StringType
			attrValues["costing_bucket"] = convert.StrToType(cfg.CostingBucket)
		}

		if cfg.CostingFolder != nil {
			attrTypes["costing_folder"] = types.StringType
			attrValues["costing_folder"] = convert.StrToType(cfg.CostingFolder)
		}

		if cfg.CostingAccessKey != nil {
			attrTypes["costing_key"] = types.StringType
			attrValues["costing_key"] = convert.StrToType(cfg.CostingAccessKey)
		}

		if cfg.CostingReportName != nil {
			attrTypes["costing_report_name"] = types.StringType
			attrValues["costing_report_name"] = convert.StrToType(cfg.CostingReportName)
		}

		if cfg.CostingSecretKey.Get() != nil {
			attrTypes["costing_secret"] = types.StringType
			attrValues["costing_secret"] = convert.StrToType(cfg.CostingSecretKey.Get())
		}

		if cfg.UseHostCredentials != nil {
			attrTypes["credentials"] = types.StringType
			attrValues["credentials"] = convert.StrToType(cfg.UseHostCredentials)
		}

		if cfg.EbsEncryption != nil {
			attrTypes["ebs_encryption"] = types.StringType
			attrValues["ebs_encryption"] = convert.StrToType(cfg.EbsEncryption)
		}

		if cfg.Endpoint != nil {
			attrTypes["endpoint"] = types.StringType
			attrValues["endpoint"] = convert.StrToType(cfg.Endpoint)
		}

		if cfg.ReplicationMode != nil {
			attrTypes["replication_provider"] = types.StringType
			attrValues["replication_provider"] = convert.StrToType(cfg.ReplicationMode)
		}

		if cfg.StsAssumeRole != nil {
			attrTypes["role_arn"] = types.StringType
			attrValues["role_arn"] = convert.StrToType(cfg.StsAssumeRole)
		}

		if cfg.SecretKey != nil {
			attrTypes["secret_key"] = types.StringType
			attrValues["secret_key"] = convert.StrToType(cfg.SecretKey)
		}

		if cfg.Vpc != nil {
			attrTypes["vpc"] = types.StringType
			attrValues["vpc"] = convert.StrToType(cfg.Vpc)
		}

		// These fields are not returned by the API yet.
		// We want to avoid overwriting any values set in the config.
		attrTypes["api_proxy"] = types.StringType
		if plan.ConfigAws.ApiProxy.IsUnknown() {
			attrValues["api_proxy"] = types.StringUnknown()
		} else {
			attrValues["api_proxy"] = plan.ConfigAws.ApiProxy
		}

		attrTypes["bypass_proxy"] = types.BoolType
		if plan.ConfigAws.BypassProxy.IsUnknown() {
			attrValues["bypass_proxy"] = types.BoolUnknown()
		} else {
			attrValues["bypass_proxy"] = plan.ConfigAws.BypassProxy
		}

		attrTypes["change_management_config"] = types.StringType
		if plan.ConfigAws.ChangeManagementConfig.IsUnknown() {
			attrValues["change_management_config"] = types.StringUnknown()
		} else {
			attrValues["change_management_config"] = plan.ConfigAws.ChangeManagementConfig
		}

		attrTypes["cmdb_config"] = types.StringType
		if plan.ConfigAws.CmdbConfig.IsUnknown() {
			attrValues["cmdb_config"] = types.StringUnknown()
		} else {
			attrValues["cmdb_config"] = plan.ConfigAws.CmdbConfig
		}

		attrTypes["costing"] = types.StringType
		if plan.ConfigAws.Costing.IsUnknown() {
			attrValues["costing"] = types.StringUnknown()
		} else {
			attrValues["costing"] = plan.ConfigAws.Costing
		}

		attrTypes["dark_mode_logo"] = types.StringType
		if plan.ConfigAws.DarkModeLogo.IsUnknown() {
			attrValues["dark_mode_logo"] = types.StringUnknown()
		} else {
			attrValues["dark_mode_logo"] = plan.ConfigAws.DarkModeLogo
		}

		attrTypes["domain"] = types.StringType
		if plan.ConfigAws.Domain.IsUnknown() {
			attrValues["domain"] = types.StringUnknown()
		} else {
			attrValues["domain"] = plan.ConfigAws.Domain
		}

		attrTypes["guidance"] = types.StringType
		if plan.ConfigAws.Guidance.IsUnknown() {
			attrValues["guidance"] = types.StringUnknown()
		} else {
			attrValues["guidance"] = plan.ConfigAws.Guidance
		}

		attrTypes["logo"] = types.StringType
		if plan.ConfigAws.Logo.IsUnknown() {
			attrValues["logo"] = types.StringUnknown()
		} else {
			attrValues["logo"] = plan.ConfigAws.Logo
		}

		attrTypes["network_mode"] = types.StringType
		if plan.ConfigAws.NetworkMode.IsUnknown() {
			attrValues["network_mode"] = types.StringUnknown()
		} else {
			attrValues["network_mode"] = plan.ConfigAws.NetworkMode
		}

		attrTypes["no_proxy"] = types.StringType
		if plan.ConfigAws.NoProxy.IsUnknown() {
			attrValues["no_proxy"] = types.StringUnknown()
		} else {
			attrValues["no_proxy"] = plan.ConfigAws.NoProxy
		}

		attrTypes["proxy"] = types.StringType
		if plan.ConfigAws.Proxy.IsUnknown() {
			attrValues["proxy"] = types.StringUnknown()
		} else {
			attrValues["proxy"] = plan.ConfigAws.Proxy
		}

		attrTypes["region"] = types.StringType
		if plan.ConfigAws.Region.IsUnknown() {
			attrValues["region"] = types.StringUnknown()
		} else {
			attrValues["region"] = plan.ConfigAws.Region
		}

		attrTypes["timezone"] = types.StringType
		if plan.ConfigAws.Timezone.IsUnknown() {
			attrValues["timezone"] = types.StringUnknown()
		} else {
			attrValues["timezone"] = plan.ConfigAws.Timezone
		}

		attrTypes["user_data"] = types.StringType
		if plan.ConfigAws.UserData.IsUnknown() {
			attrValues["user_data"] = types.StringUnknown()
		} else {
			attrValues["user_data"] = plan.ConfigAws.UserData
		}

		attrTypes["vdi_gateway"] = types.StringType
		if plan.ConfigAws.VdiGateway.IsUnknown() {
			attrValues["vdi_gateway"] = types.StringUnknown()
		} else {
			attrValues["vdi_gateway"] = plan.ConfigAws.VdiGateway
		}

		if len(attrValues) > 0 {
			configAws, diagsAws := NewConfigAwsValue(attrTypes, attrValues)
			if diagsAws.HasError() {
				diags.Append(diagsAws...)
				diags.AddError(
					readOperation,
					fmt.Sprintf("cloud %d: failed to decode AWS configuration", id),
				)

				return state, diags
			}

			state.ConfigAws = configAws
		}

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
