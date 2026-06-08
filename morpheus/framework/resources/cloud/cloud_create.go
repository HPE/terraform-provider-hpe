// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/sdkfuncs"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const createOperation = "create cloud resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan CloudModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	tenantID := plan.TenantId.ValueInt64()

	var config CloudModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addCloud := &sdk.AddCloudsRequestZone{}
	addCloud.Name = name

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		addCloud.GroupId = plan.GroupId.ValueInt64()
	}

	addCloudConfig := addCloud.Config

	var cloudTypeCode string

	switch {
	case !plan.ConfigAws.IsNull() && !plan.ConfigAws.IsUnknown():
		cloudTypeCode = awsCloud

		config := sdkfuncs.NewAwsCloudConfig(plan.ConfigAws.Endpoint.ValueString())

		// Shove common config fields into the config object if they are set in the plan
		// even if they aren't part of the specific config struct in the SDK
		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			config.ApplianceUrl = plan.ApplianceUrl.ValueStringPointer()
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.DatacenterName = plan.DataCenterName.ValueStringPointer()
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.ExternalId.Set(plan.ExternalId.ValueStringPointer())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.InventoryLevel = plan.ImportExistingVms.ValueStringPointer()
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.ConsoleKeymap = plan.KeyboardLayout.ValueStringPointer()
		}

		if !plan.ConfigAws.AccessKey.IsNull() &&
			!plan.ConfigAws.AccessKey.IsUnknown() {
			config.AccessKey = plan.ConfigAws.AccessKey.ValueStringPointer()
		}

		if !plan.ConfigAws.ApiProxy.IsNull() &&
			!plan.ConfigAws.ApiProxy.IsUnknown() {
			config.ApiProxy.Set(plan.ConfigAws.ApiProxy.ValueStringPointer())
		}

		if !plan.ConfigAws.BypassProxy.IsNull() &&
			!plan.ConfigAws.BypassProxy.IsUnknown() {
			config.BypassProxyForCloud = convert.BoolTypeToStringPointerOnOff(plan.ConfigAws.BypassProxy)
		}

		if !plan.ConfigAws.ChangeManagementConfig.IsNull() &&
			!plan.ConfigAws.ChangeManagementConfig.IsUnknown() {
			config.ChangeManagementConfig = plan.ConfigAws.ChangeManagementConfig.ValueStringPointer()
		}

		if !plan.ConfigAws.CmdbConfig.IsNull() &&
			!plan.ConfigAws.CmdbConfig.IsUnknown() {
			config.CmdbConfig = plan.ConfigAws.CmdbConfig.ValueStringPointer()
		}

		if !plan.ConfigAws.CmdbDiscovery.IsNull() &&
			!plan.ConfigAws.CmdbDiscovery.IsUnknown() {
			config.ConfigCmdbDiscovery = convert.BoolTypeToStringPointerOnOff(plan.ConfigAws.CmdbDiscovery)
		}

		if !plan.ConfigAws.ConfigManagementId.IsNull() &&
			!plan.ConfigAws.ConfigManagementId.IsUnknown() {
			config.ConfigManagementId = plan.ConfigAws.ConfigManagementId.ValueStringPointer()
		}

		if !plan.ConfigAws.Costing.IsNull() &&
			!plan.ConfigAws.Costing.IsUnknown() {
			config.Costing = plan.ConfigAws.Costing.ValueStringPointer()
		}

		if !plan.ConfigAws.CostingBucket.IsNull() &&
			!plan.ConfigAws.CostingBucket.IsUnknown() {
			config.CostingBucket = plan.ConfigAws.CostingBucket.ValueStringPointer()
		}

		if !plan.ConfigAws.CostingFolder.IsNull() &&
			!plan.ConfigAws.CostingFolder.IsUnknown() {
			config.CostingFolder.Set(plan.ConfigAws.CostingFolder.ValueStringPointer())
		}

		if !plan.ConfigAws.CostingKey.IsNull() &&
			!plan.ConfigAws.CostingKey.IsUnknown() {
			config.CostingKey.Set(plan.ConfigAws.CostingKey.ValueStringPointer())
		}

		if !plan.ConfigAws.CostingReportName.IsNull() &&
			!plan.ConfigAws.CostingReportName.IsUnknown() {
			config.CostingReportName.Set(plan.ConfigAws.CostingReportName.ValueStringPointer())
		}

		if !plan.ConfigAws.CostingSecret.IsNull() &&
			!plan.ConfigAws.CostingSecret.IsUnknown() {
			config.CostingSecret.Set(plan.ConfigAws.CostingSecret.ValueStringPointer())
		}

		if !plan.ConfigAws.Credentials.IsNull() &&
			!plan.ConfigAws.Credentials.IsUnknown() {
			config.Credentials = plan.ConfigAws.Credentials.ValueStringPointer()
		}

		if !plan.ConfigAws.DarkModeLogo.IsNull() &&
			!plan.ConfigAws.DarkModeLogo.IsUnknown() {
			config.DarkModeLogo.Set(plan.ConfigAws.DarkModeLogo.ValueStringPointer())
		}

		if !plan.ConfigAws.Domain.IsNull() &&
			!plan.ConfigAws.Domain.IsUnknown() {
			config.Domain = plan.ConfigAws.Domain.ValueStringPointer()
		}

		if !plan.ConfigAws.EbsEncryption.IsNull() &&
			!plan.ConfigAws.EbsEncryption.IsUnknown() {
			config.EbsEncryption = plan.ConfigAws.EbsEncryption.ValueStringPointer()
		}

		if !plan.ConfigAws.Guidance.IsNull() &&
			!plan.ConfigAws.Guidance.IsUnknown() {
			config.Guidance = plan.ConfigAws.Guidance.ValueStringPointer()
		}

		if !plan.ConfigAws.Logo.IsNull() &&
			!plan.ConfigAws.Logo.IsUnknown() {
			config.Logo.Set(plan.ConfigAws.Logo.ValueStringPointer())
		}

		if !plan.ConfigAws.NetworkMode.IsNull() &&
			!plan.ConfigAws.NetworkMode.IsUnknown() {
			config.NetworkMode.Set(plan.ConfigAws.NetworkMode.ValueStringPointer())
		}

		if !plan.ConfigAws.NoProxy.IsNull() &&
			!plan.ConfigAws.NoProxy.IsUnknown() {
			config.NoProxy = plan.ConfigAws.NoProxy.ValueStringPointer()
		}

		if !plan.ConfigAws.Proxy.IsNull() &&
			!plan.ConfigAws.Proxy.IsUnknown() {
			config.Proxy = plan.ConfigAws.Proxy.ValueStringPointer()
		}

		if !plan.ConfigAws.Region.IsNull() &&
			!plan.ConfigAws.Region.IsUnknown() {
			config.Region = plan.ConfigAws.Region.ValueStringPointer()
		}

		if !plan.ConfigAws.RoleArn.IsNull() &&
			!plan.ConfigAws.RoleArn.IsUnknown() {
			config.StsAssumeRole = plan.ConfigAws.RoleArn.ValueStringPointer()
		}

		if !plan.ConfigAws.SecretKey.IsNull() &&
			!plan.ConfigAws.SecretKey.IsUnknown() {
			config.SecretKey = plan.ConfigAws.SecretKey.ValueStringPointer()
		}

		if !plan.ConfigAws.Timezone.IsNull() &&
			!plan.ConfigAws.Timezone.IsUnknown() {
			config.Timezone = plan.ConfigAws.Timezone.ValueStringPointer()
		}

		if !plan.ConfigAws.UserData.IsNull() &&
			!plan.ConfigAws.UserData.IsUnknown() {
			config.UserDataLinux = plan.ConfigAws.UserData.ValueStringPointer()
		}

		if !plan.ConfigAws.VdiGateway.IsNull() &&
			!plan.ConfigAws.VdiGateway.IsUnknown() {
			config.VdiGateway = plan.ConfigAws.VdiGateway.ValueStringPointer()
		}

		if !plan.ConfigAws.Vpc.IsNull() &&
			!plan.ConfigAws.Vpc.IsUnknown() {
			config.Vpc = plan.ConfigAws.Vpc.ValueStringPointer()
		}

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf = config

	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		cloudTypeCode = standardCloud

		config := sdkfuncs.NewHvmCloudConfig()

		// Shove common config fields into the config object if they are set in the plan
		// even if they aren't part of the specific config struct in the SDK
		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			config.ApplianceUrl = plan.ApplianceUrl.ValueStringPointer()
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.DatacenterName = plan.DataCenterName.ValueStringPointer()
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.ExternalId.Set(plan.ExternalId.ValueStringPointer())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.InventoryLevel = plan.ImportExistingVms.ValueStringPointer()
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.ConsoleKeymap = plan.KeyboardLayout.ValueStringPointer()
		}

		if !plan.ConfigHvm.CertificateProvider.IsNull() &&
			!plan.ConfigHvm.CertificateProvider.IsUnknown() {
			config.CertificateProvider = plan.ConfigHvm.CertificateProvider.ValueStringPointer()
		}

		if !plan.ConfigHvm.EnableNetworkTypeSelection.IsNull() &&
			!plan.ConfigHvm.EnableNetworkTypeSelection.IsUnknown() {
			config.EnableNetworkTypeSelection.Set(
				convert.BoolTypeToStringPointerOnOff(plan.ConfigHvm.EnableNetworkTypeSelection),
			)
		}

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf2 = config

	case !plan.ConfigVmware.IsNull() && !plan.ConfigVmware.IsUnknown():
		cloudTypeCode = vmwareCloud

		config := sdkfuncs.NewVmwareCloudConfig(
			plan.ConfigVmware.ApiUrl.ValueString(),
			plan.ConfigVmware.ApiVersion.ValueString(),
			plan.ConfigVmware.Datacenter.ValueString(),
		)

		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			config.ApplianceUrl = plan.ApplianceUrl.ValueStringPointer()
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.DatacenterName = plan.DataCenterName.ValueStringPointer()
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.ExternalId.Set(plan.ExternalId.ValueStringPointer())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.InventoryLevel = plan.ImportExistingVms.ValueStringPointer()
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.ConsoleKeymap = plan.KeyboardLayout.ValueStringPointer()
		}

		if !plan.ConfigVmware.CertificateProvider.IsNull() &&
			!plan.ConfigVmware.CertificateProvider.IsUnknown() {
			config.CertificateProvider = plan.ConfigVmware.CertificateProvider.ValueStringPointer()
		}

		if !plan.ConfigVmware.Cluster.IsNull() &&
			!plan.ConfigVmware.Cluster.IsUnknown() {
			config.Cluster = plan.ConfigVmware.Cluster.ValueStringPointer()
		}

		if !plan.ConfigVmware.ConfigManagementId.IsNull() &&
			!plan.ConfigVmware.ConfigManagementId.IsUnknown() {
			config.ConfigManagementId = plan.ConfigVmware.ConfigManagementId.ValueStringPointer()
		}

		if !plan.ConfigVmware.EnableDiskTypeSelection.IsNull() &&
			!plan.ConfigVmware.EnableDiskTypeSelection.IsUnknown() {
			config.EnableDiskTypeSelection.Set(
				convert.BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableDiskTypeSelection),
			)
		}

		if !plan.ConfigVmware.EnableNetworkTypeSelection.IsNull() &&
			!plan.ConfigVmware.EnableNetworkTypeSelection.IsUnknown() {
			config.EnableNetworkTypeSelection.Set(
				convert.BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableNetworkTypeSelection),
			)
		}

		if !plan.ConfigVmware.EnableStorageTypeSelection.IsNull() &&
			!plan.ConfigVmware.EnableStorageTypeSelection.IsUnknown() {
			config.EnableStorageTypeSelection.Set(
				convert.BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableStorageTypeSelection),
			)
		}

		if !plan.ConfigVmware.EnableVnc.IsNull() &&
			!plan.ConfigVmware.EnableVnc.IsUnknown() {
			config.EnableVnc.Set(convert.BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableVnc))
		}

		if !plan.ConfigVmware.HideHostSelection.IsNull() &&
			!plan.ConfigVmware.HideHostSelection.IsUnknown() {
			config.HideHostSelection.Set(convert.BoolTypeToStringPointerOnOff(plan.ConfigVmware.HideHostSelection))
		}

		if !plan.ConfigVmware.Password.IsNull() &&
			!plan.ConfigVmware.Password.IsUnknown() {
			config.Password = plan.ConfigVmware.Password.ValueStringPointer()
		}

		if !plan.ConfigVmware.ResourcePool.IsNull() &&
			!plan.ConfigVmware.ResourcePool.IsUnknown() {
			config.ResourcePool = plan.ConfigVmware.ResourcePool.ValueStringPointer()
		}

		if !plan.ConfigVmware.RpcMode.IsNull() &&
			!plan.ConfigVmware.RpcMode.IsUnknown() {
			config.RpcMode.Set(plan.ConfigVmware.RpcMode.ValueStringPointer())
		}

		if !plan.ConfigVmware.Username.IsNull() &&
			!plan.ConfigVmware.Username.IsUnknown() {
			config.Username = plan.ConfigVmware.Username.ValueStringPointer()
		}

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf3 = config

	case !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown():
		cloudTypeCode = azureCloud

		config := sdkfuncs.NewAzureCloudConfig()
		config.AdditionalProperties = make(map[string]interface{})

		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			config.ApplianceUrl = plan.ApplianceUrl.ValueStringPointer()
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.DatacenterName = plan.DataCenterName.ValueStringPointer()
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.ExternalId.Set(plan.ExternalId.ValueStringPointer())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.InventoryLevel = plan.ImportExistingVms.ValueStringPointer()
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.ConsoleKeymap = plan.KeyboardLayout.ValueStringPointer()
		}

		if !plan.ConfigAzure.AzureRegion.IsNull() && !plan.ConfigAzure.AzureRegion.IsUnknown() {
			config.AdditionalProperties["azureRegion"] = plan.ConfigAzure.AzureRegion.ValueString()
		}

		if !plan.ConfigAzure.CmdbDiscovery.IsNull() && !plan.ConfigAzure.CmdbDiscovery.IsUnknown() {
			config.AdditionalProperties["configCmdbDiscovery"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigAzure.CmdbDiscovery)
		}

		if !plan.ConfigAzure.SubscriberId.IsNull() && !plan.ConfigAzure.SubscriberId.IsUnknown() {
			config.SubscriberId = plan.ConfigAzure.SubscriberId.ValueStringPointer()
		}

		if !plan.ConfigAzure.TenantId.IsNull() && !plan.ConfigAzure.TenantId.IsUnknown() {
			config.TenantId = plan.ConfigAzure.TenantId.ValueStringPointer()
		}

		if !plan.ConfigAzure.ClientId.IsNull() && !plan.ConfigAzure.ClientId.IsUnknown() {
			config.ClientId = plan.ConfigAzure.ClientId.ValueStringPointer()
		}

		if !plan.ConfigAzure.ClientSecret.IsNull() && !plan.ConfigAzure.ClientSecret.IsUnknown() {
			config.ClientSecret = plan.ConfigAzure.ClientSecret.ValueStringPointer()
		}

		if !plan.ConfigAzure.ResourceGroup.IsNull() && !plan.ConfigAzure.ResourceGroup.IsUnknown() {
			config.ResourceGroup = plan.ConfigAzure.ResourceGroup.ValueStringPointer()
		}

		if !plan.ConfigAzure.CloudType.IsNull() && !plan.ConfigAzure.CloudType.IsUnknown() {
			config.CloudType = plan.ConfigAzure.CloudType.ValueStringPointer()
		}

		if !plan.ConfigAzure.ImportExisting.IsNull() && !plan.ConfigAzure.ImportExisting.IsUnknown() {
			config.ImportExisting = plan.ConfigAzure.ImportExisting.ValueStringPointer()
		}

		if !plan.ConfigAzure.StorageAccount.IsNull() && !plan.ConfigAzure.StorageAccount.IsUnknown() {
			config.StorageAccount = plan.ConfigAzure.StorageAccount.ValueStringPointer()
		}

		if !plan.ConfigAzure.RpcMode.IsNull() && !plan.ConfigAzure.RpcMode.IsUnknown() {
			config.RpcMode.Set(plan.ConfigAzure.RpcMode.ValueStringPointer())
		}

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf1 = config

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		if plan.CloudTypeCode.IsNull() || plan.CloudTypeCode.IsUnknown() {
			resp.Diagnostics.AddError(
				createOperation,
				"cloud "+name+": cloud_type_code is required for generic configurations",
			)

			return
		}

		cloudTypeCode = plan.CloudTypeCode.ValueString()

		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"cloud "+name+": failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				createOperation,
				"cloud "+name+": config must be a valid object/map",
			)

			return
		}

		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			configDataMap["applianceUrl"] = plan.ApplianceUrl.ValueString()
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			configDataMap["datacenterName"] = plan.DataCenterName.ValueString()
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			configDataMap["externalId"] = plan.ExternalId.ValueString()
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			configDataMap["inventoryLevel"] = plan.ImportExistingVms.ValueString()
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			configDataMap["consoleKeymap"] = plan.KeyboardLayout.ValueString()
		}

		addCloudConfig.MapmapOfStringAny = &configDataMap
	}

	addCloud.Config = addCloudConfig

	cloudTypeWithCode := &sdk.AddCloudsRequestZoneZoneTypeAnyOf1{Code: &cloudTypeCode}

	cloudType := addCloud.ZoneType
	cloudType.AddCloudsRequestZoneZoneTypeAnyOf1 = cloudTypeWithCode

	addCloud.ZoneType = cloudType

	addCloud.AdditionalProperties = make(map[string]any)

	addCloud.AccountId = &tenantID

	if !plan.AgentInstallMode.IsNull() && !plan.AgentInstallMode.IsUnknown() {
		addCloud.AgentMode = plan.AgentInstallMode.ValueStringPointer()
	}

	if !plan.AutoRecoverPowerState.IsNull() && !plan.AutoRecoverPowerState.IsUnknown() {
		addCloud.AutoRecoverPowerState = plan.AutoRecoverPowerState.ValueBoolPointer()
	}

	if !plan.DefaultDatastoreSyncActive.IsNull() && !plan.DefaultDatastoreSyncActive.IsUnknown() {
		addCloud.DefaultDatastoreSyncActive = plan.DefaultDatastoreSyncActive.ValueBoolPointer()
	}
	if !plan.DefaultFolderSyncActive.IsNull() && !plan.DefaultFolderSyncActive.IsUnknown() {
		addCloud.DefaultFolderSyncActive = plan.DefaultFolderSyncActive.ValueBoolPointer()
	}
	if !plan.DefaultNetworkSyncActive.IsNull() && !plan.DefaultNetworkSyncActive.IsUnknown() {
		addCloud.DefaultNetworkSyncActive = plan.DefaultNetworkSyncActive.ValueBoolPointer()
	}
	if !plan.DefaultPlanSyncActive.IsNull() && !plan.DefaultPlanSyncActive.IsUnknown() {
		addCloud.DefaultPlanSyncActive = plan.DefaultPlanSyncActive.ValueBoolPointer()
	}
	if !plan.DefaultPoolSyncActive.IsNull() && !plan.DefaultPoolSyncActive.IsUnknown() {
		addCloud.DefaultPoolSyncActive = plan.DefaultPoolSyncActive.ValueBoolPointer()
	}
	if !plan.DefaultSecurityGroupSyncActive.IsNull() && !plan.DefaultSecurityGroupSyncActive.IsUnknown() {
		addCloud.DefaultSecurityGroupSyncActive = plan.DefaultSecurityGroupSyncActive.ValueBoolPointer()
	}

	if !plan.Code.IsNull() && !plan.Code.IsUnknown() {
		addCloud.Code = plan.Code.ValueStringPointer()
	}

	if !plan.CostingMode.IsNull() && !plan.CostingMode.IsUnknown() {
		addCloud.AdditionalProperties["costingMode"] = plan.CostingMode.ValueString()
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		addCloud.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		addCloud.AdditionalProperties["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.GuidanceMode.IsNull() && !plan.GuidanceMode.IsUnknown() {
		addCloud.AdditionalProperties["guidanceMode"] = plan.GuidanceMode.ValueString()
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				createOperation,
				"cloud "+name+": failed to parse label: "+err.Error(),
			)

			return
		}

		addCloud.Labels = labels
	}

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		addCloud.Location.Set(plan.Location.ValueStringPointer())
	}

	if !plan.SecurityMode.IsNull() && !plan.SecurityMode.IsUnknown() {
		addCloud.SecurityMode = plan.SecurityMode.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		addCloud.Visibility = plan.Visibility.ValueStringPointer()
	}

	createRequest := &sdk.AddCloudsRequest{Zone: *addCloud}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			"cloud "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	// Call API
	cloud, hresp, err := client.CloudsAPI.AddClouds(ctx).AddCloudsRequest(*createRequest).Execute()
	if cloud == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			"cloud "+name+" POST failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	id := *cloud.Zone.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cloud",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getCloudAsState(ctx, id, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("Cloud %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("Cloud %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}
