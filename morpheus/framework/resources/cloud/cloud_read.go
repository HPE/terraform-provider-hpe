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
	azureCloud    = "azure"
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

	if c.Zone == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("cloud %d missing zone", id),
		)

		return state, diags
	}

	cloud := c.Zone

	if cloud.Config == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("cloud %d missing config", id),
		)

		return state, diags
	}

	importing := plan.Name.IsNull()
	cloudType := *cloud.ZoneType.Code

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		state.GroupId = plan.GroupId
	} else {
		// On import, use the first API value
		state.GroupId = types.Int64Null()
		if len(cloud.Groups) > 0 {
			state.GroupId = convert.Int64ToType(cloud.Groups[0].Id)
		}
	}

	state.AgentInstallMode = convert.StrToType(cloud.AgentMode)
	state.AutoRecoverPowerState = convert.BoolToType(cloud.AutoRecoverPowerState)
	if cloud.Code == nil || *cloud.Code != "standard" { // workaround an API bug
		state.Code = convert.StrToType(cloud.Code)
	}
	state.CostingMode = convert.StrToType(cloud.CostingMode.Get())
	state.DefaultDatastoreSyncActive = convert.BoolToType(cloud.DefaultDatastoreSyncActive)
	state.DefaultFolderSyncActive = convert.BoolToType(cloud.DefaultFolderSyncActive)
	state.DefaultNetworkSyncActive = convert.BoolToType(cloud.DefaultNetworkSyncActive)
	state.DefaultPlanSyncActive = convert.BoolToType(cloud.DefaultPlanSyncActive)
	state.DefaultPoolSyncActive = convert.BoolToType(cloud.DefaultPoolSyncActive)
	state.DefaultSecurityGroupSyncActive = convert.BoolToType(cloud.DefaultSecurityGroupSyncActive)
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
		cfg := cloud.Config.GetClouds200ResponseZoneConfigAnyOf2

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl)
		state.DataCenterName = convert.StrToType(cfg.DatacenterName)

		attrTypes := ConfigAwsValue{}.AttributeTypes(ctx)
		attrValues := make(map[string]attr.Value)

		if cfg.AccessKey != nil {
			attrValues["access_key"] = convert.StrToType(cfg.AccessKey)
		} else {
			attrValues["access_key"] = types.StringNull()
		}

		if cfg.BackupMode != nil {
			attrValues["backup_provider"] = convert.StrToType(cfg.BackupMode)
		} else {
			attrValues["backup_provider"] = types.StringNull()
		}

		if cfg.ConfigCmdbDiscovery != nil {
			attrValues["cmdb_discovery"] = convert.BoolToType(cfg.ConfigCmdbDiscovery)
		} else {
			attrValues["cmdb_discovery"] = types.BoolNull()
		}

		if cfg.ConfigManagementId != nil {
			attrValues["config_management_id"] = convert.StrToType(cfg.ConfigManagementId)
		} else {
			attrValues["config_management_id"] = types.StringNull()
		}

		if cfg.CostingBucket != nil {
			attrValues["costing_bucket"] = convert.StrToType(cfg.CostingBucket)
		} else {
			attrValues["costing_bucket"] = types.StringNull()
		}

		if cfg.CostingFolder != nil {
			attrValues["costing_folder"] = convert.StrToType(cfg.CostingFolder)
		} else {
			attrValues["costing_folder"] = types.StringNull()
		}

		if cfg.CostingAccessKey != nil {
			attrValues["costing_key"] = convert.StrToType(cfg.CostingAccessKey)
		} else {
			attrValues["costing_key"] = types.StringNull()
		}

		if cfg.CostingReportName != nil {
			attrValues["costing_report_name"] = convert.StrToType(cfg.CostingReportName)
		} else {
			attrValues["costing_report_name"] = types.StringNull()
		}

		if cfg.CostingSecretKey.Get() != nil {
			attrValues["costing_secret"] = convert.StrToType(cfg.CostingSecretKey.Get())
		} else {
			attrValues["costing_secret"] = types.StringNull()
		}

		if cfg.UseHostCredentials != nil {
			attrValues["credentials"] = convert.StrToType(cfg.UseHostCredentials)
		} else {
			attrValues["credentials"] = types.StringNull()
		}

		if cfg.EbsEncryption != nil {
			attrValues["ebs_encryption"] = convert.StrToType(cfg.EbsEncryption)
		} else {
			attrValues["ebs_encryption"] = types.StringNull()
		}

		if cfg.Endpoint != nil {
			attrValues["endpoint"] = convert.StrToType(cfg.Endpoint)
		} else {
			attrValues["endpoint"] = types.StringNull()
		}

		if cfg.ReplicationMode != nil {
			attrValues["replication_provider"] = convert.StrToType(cfg.ReplicationMode)
		} else {
			attrValues["replication_provider"] = types.StringNull()
		}

		if cfg.StsAssumeRole != nil {
			attrValues["role_arn"] = convert.StrToType(cfg.StsAssumeRole)
		} else {
			attrValues["role_arn"] = types.StringNull()
		}

		if cfg.SecretKey != nil {
			attrValues["secret_key"] = convert.StrToType(cfg.SecretKey)
		} else {
			attrValues["secret_key"] = types.StringNull()
		}

		if cfg.Vpc != nil {
			attrValues["vpc"] = convert.StrToType(cfg.Vpc)
		} else {
			attrValues["vpc"] = types.StringNull()
		}

		// These fields are not returned by the API yet.
		// We want to avoid overwriting any values set in the config.
		if plan.ConfigAws.ApiProxy.IsUnknown() {
			attrValues["api_proxy"] = types.StringUnknown()
		} else {
			attrValues["api_proxy"] = plan.ConfigAws.ApiProxy
		}

		if plan.ConfigAws.BypassProxy.IsUnknown() {
			attrValues["bypass_proxy"] = types.BoolUnknown()
		} else {
			attrValues["bypass_proxy"] = plan.ConfigAws.BypassProxy
		}

		if plan.ConfigAws.ChangeManagementConfig.IsUnknown() {
			attrValues["change_management_config"] = types.StringUnknown()
		} else {
			attrValues["change_management_config"] = plan.ConfigAws.ChangeManagementConfig
		}

		if plan.ConfigAws.CmdbConfig.IsUnknown() {
			attrValues["cmdb_config"] = types.StringUnknown()
		} else {
			attrValues["cmdb_config"] = plan.ConfigAws.CmdbConfig
		}

		if plan.ConfigAws.Costing.IsUnknown() {
			attrValues["costing"] = types.StringUnknown()
		} else {
			attrValues["costing"] = plan.ConfigAws.Costing
		}

		if plan.ConfigAws.DarkModeLogo.IsUnknown() {
			attrValues["dark_mode_logo"] = types.StringUnknown()
		} else {
			attrValues["dark_mode_logo"] = plan.ConfigAws.DarkModeLogo
		}

		if plan.ConfigAws.Domain.IsUnknown() {
			attrValues["domain"] = types.StringUnknown()
		} else {
			attrValues["domain"] = plan.ConfigAws.Domain
		}

		if plan.ConfigAws.Guidance.IsUnknown() {
			attrValues["guidance"] = types.StringUnknown()
		} else {
			attrValues["guidance"] = plan.ConfigAws.Guidance
		}

		if plan.ConfigAws.Logo.IsUnknown() {
			attrValues["logo"] = types.StringUnknown()
		} else {
			attrValues["logo"] = plan.ConfigAws.Logo
		}

		if plan.ConfigAws.NetworkMode.IsUnknown() {
			attrValues["network_mode"] = types.StringUnknown()
		} else {
			attrValues["network_mode"] = plan.ConfigAws.NetworkMode
		}

		if plan.ConfigAws.NoProxy.IsUnknown() {
			attrValues["no_proxy"] = types.StringUnknown()
		} else {
			attrValues["no_proxy"] = plan.ConfigAws.NoProxy
		}

		if plan.ConfigAws.Proxy.IsUnknown() {
			attrValues["proxy"] = types.StringUnknown()
		} else {
			attrValues["proxy"] = plan.ConfigAws.Proxy
		}

		if plan.ConfigAws.Region.IsUnknown() {
			attrValues["region"] = types.StringUnknown()
		} else {
			attrValues["region"] = plan.ConfigAws.Region
		}

		if plan.ConfigAws.Timezone.IsUnknown() {
			attrValues["timezone"] = types.StringUnknown()
		} else {
			attrValues["timezone"] = plan.ConfigAws.Timezone
		}

		if plan.ConfigAws.UserData.IsUnknown() {
			attrValues["user_data"] = types.StringUnknown()
		} else {
			attrValues["user_data"] = plan.ConfigAws.UserData
		}

		if plan.ConfigAws.VdiGateway.IsUnknown() {
			attrValues["vdi_gateway"] = types.StringUnknown()
		} else {
			attrValues["vdi_gateway"] = plan.ConfigAws.VdiGateway
		}

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

	case cloudType == standardCloud && (!plan.ConfigHvm.IsNull() || importing):
		cfg := cloud.Config.GetClouds200ResponseZoneConfigAnyOf

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
		state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
		state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
		state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
		state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())

		attrTypes := ConfigHvmValue{}.AttributeTypes(ctx)
		attrValues := make(map[string]attr.Value)

		if cfg.CertificateProvider.Get() != nil {
			attrValues["certificate_provider"] = convert.StrToType(cfg.CertificateProvider.Get())
		} else {
			attrValues["certificate_provider"] = types.StringNull()
		}

		if cfg.EnableNetworkTypeSelection.Get() != nil {
			attrValues["enable_network_type_selection"] = convert.StringToBool(ctx, *cfg.EnableNetworkTypeSelection.Get())
		} else {
			attrValues["enable_network_type_selection"] = types.BoolNull()
		}

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

	case cloudType == vmwareCloud && (!plan.ConfigVmware.IsNull() || importing):
		cfg := cloud.Config.GetClouds200ResponseZoneConfigAnyOf1

		// Move these common fields up
		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl.Get())
		state.DataCenterName = convert.StrToType(cfg.DatacenterName.Get())
		state.ExternalId = convert.StrToType(cfg.ExternalId.Get())
		state.ImportExistingVms = convert.StrToType(cfg.InventoryLevel.Get())
		state.KeyboardLayout = convert.StrToType(cfg.ConsoleKeymap.Get())

		attrTypes := ConfigVmwareValue{}.AttributeTypes(ctx)
		attrValues := make(map[string]attr.Value)

		if cfg.ApiUrl != nil {
			attrValues["api_url"] = convert.StrToType(cfg.ApiUrl)
		} else {
			attrValues["api_url"] = types.StringNull()
		}

		if cfg.ApiVersion.Get() != nil {
			attrValues["api_version"] = convert.StrToType(cfg.ApiVersion.Get())
		} else {
			attrValues["api_version"] = types.StringNull()
		}

		if cfg.CertificateProvider.Get() != nil {
			attrValues["certificate_provider"] = convert.StrToType(cfg.CertificateProvider.Get())
		} else {
			attrValues["certificate_provider"] = types.StringNull()
		}

		if cfg.Cluster != nil {
			attrValues["cluster"] = convert.StrToType(cfg.Cluster)
		} else {
			attrValues["cluster"] = types.StringNull()
		}

		if cfg.ConfigManagementId.Get() != nil {
			attrValues["config_management_id"] = convert.StrToType(cfg.ConfigManagementId.Get())
		} else {
			attrValues["config_management_id"] = types.StringNull()
		}

		if cfg.Datacenter != nil {
			attrValues["datacenter"] = convert.StrToType(cfg.Datacenter)
		} else {
			attrValues["datacenter"] = types.StringNull()
		}

		if cfg.EnableDiskTypeSelection.Get() != nil {
			attrValues["enable_disk_type_selection"] = convert.StringToBool(ctx, *cfg.EnableDiskTypeSelection.Get())
		} else {
			attrValues["enable_disk_type_selection"] = types.BoolNull()
		}

		if cfg.EnableNetworkTypeSelection.Get() != nil {
			attrValues["enable_network_type_selection"] = convert.StringToBool(ctx, *cfg.EnableNetworkTypeSelection.Get())
		} else {
			attrValues["enable_network_type_selection"] = types.BoolNull()
		}

		if cfg.EnableStorageTypeSelection.Get() != nil {
			attrValues["enable_storage_type_selection"] = convert.StringToBool(ctx, *cfg.EnableStorageTypeSelection.Get())
		} else {
			attrValues["enable_storage_type_selection"] = types.BoolNull()
		}

		if cfg.EnableVnc.Get() != nil {
			attrValues["enable_vnc"] = convert.StringToBool(ctx, *cfg.EnableVnc.Get())
		} else {
			attrValues["enable_vnc"] = types.BoolNull()
		}

		if cfg.HideHostSelection.Get() != nil {
			attrValues["hide_host_selection"] = convert.StringToBool(ctx, *cfg.HideHostSelection.Get())
		} else {
			attrValues["hide_host_selection"] = types.BoolNull()
		}

		if cfg.Password.Get() != nil {
			attrValues["password"] = convert.StrToType(cfg.Password.Get())
		} else {
			attrValues["password"] = types.StringNull()
		}

		if cfg.ResourcePool != nil {
			attrValues["resource_pool"] = convert.StrToType(cfg.ResourcePool)
		} else {
			attrValues["resource_pool"] = types.StringNull()
		}

		if cfg.RpcMode != nil {
			attrValues["rpc_mode"] = convert.StrToType(cfg.RpcMode)
		} else {
			attrValues["rpc_mode"] = types.StringNull()
		}

		if cfg.Username != nil {
			attrValues["username"] = convert.StrToType(cfg.Username)
		} else {
			attrValues["username"] = types.StringNull()
		}

		configVmware, diagsVmware := NewConfigVmwareValue(attrTypes, attrValues)
		if diagsVmware.HasError() {
			diags.Append(diagsVmware...)
			diags.AddError(
				readOperation,
				fmt.Sprintf("cloud %d: failed to decode VMware configuration", id),
			)

			return state, diags
		}

		state.ConfigVmware = configVmware

	case cloudType == azureCloud && (!plan.ConfigAzure.IsNull() || importing):
		// Azure read uses GetClouds200ResponseZoneConfigAnyOf3 (the GET response model),
		// while create uses AddCloudsRequestZoneConfigAnyOf1 (the POST request model).
		// The AnyOf index differs between request/response specs but both are Azure.
		cfg := cloud.Config.GetClouds200ResponseZoneConfigAnyOf3

		state.ApplianceUrl = convert.StrToType(cfg.ApplianceUrl)
		state.DataCenterName = convert.StrToType(cfg.DatacenterName)

		attrTypes := ConfigAzureValue{}.AttributeTypes(ctx)
		attrValues := make(map[string]attr.Value)

		if azureRegion, ok := cfg.AdditionalProperties["azureRegion"].(string); ok {
			attrValues["azure_region"] = types.StringValue(azureRegion)
		} else {
			attrValues["azure_region"] = types.StringNull()
		}

		if cfg.ClientId != nil {
			attrValues["client_id"] = convert.StrToType(cfg.ClientId)
		} else {
			attrValues["client_id"] = types.StringNull()
		}

		attrValues["client_secret"] = types.StringNull()

		if cfg.CloudType != nil {
			attrValues["cloud_type"] = convert.StrToType(cfg.CloudType)
		} else {
			attrValues["cloud_type"] = types.StringNull()
		}

		if cfg.ConfigCmdbDiscovery != nil {
			attrValues["cmdb_discovery"] = convert.BoolToType(cfg.ConfigCmdbDiscovery)
		} else {
			attrValues["cmdb_discovery"] = types.BoolNull()
		}

		if cfg.ImportExisting != nil {
			attrValues["import_existing"] = convert.StrToType(cfg.ImportExisting)
		} else {
			attrValues["import_existing"] = types.StringNull()
		}

		if cfg.ResourceGroup != nil {
			attrValues["resource_group"] = convert.StrToType(cfg.ResourceGroup)
		} else {
			attrValues["resource_group"] = types.StringNull()
		}

		if cfg.RpcMode != nil {
			attrValues["rpc_mode"] = convert.StrToType(cfg.RpcMode)
		} else {
			attrValues["rpc_mode"] = types.StringNull()
		}

		if cfg.StorageAccount != nil {
			attrValues["storage_account"] = convert.StrToType(cfg.StorageAccount)
		} else {
			attrValues["storage_account"] = types.StringNull()
		}

		if cfg.SubscriberId != nil {
			attrValues["subscriber_id"] = convert.StrToType(cfg.SubscriberId)
		} else {
			attrValues["subscriber_id"] = types.StringNull()
		}

		if cfg.TenantId != nil {
			attrValues["tenant_id"] = convert.StrToType(cfg.TenantId)
		} else {
			attrValues["tenant_id"] = types.StringNull()
		}

		configAzure, diagsAzure := NewConfigAzureValue(attrTypes, attrValues)
		if diagsAzure.HasError() {
			diags.Append(diagsAzure...)
			diags.AddError(
				readOperation,
				fmt.Sprintf("cloud %d: failed to decode Azure configuration", id),
			)

			return state, diags
		}

		state.ConfigAzure = configAzure
	case !plan.Config.IsNull() || importing:
		state.CloudTypeCode = convert.StrToType(cloud.ZoneType.Code)

		state.Config = types.DynamicNull()

		cfg := cloud.Config.MapmapOfStringAny
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
			cfg := cloud.Config.MapmapOfStringAny

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
