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

	addCloud := sdk.NewAddCloudsRequestZoneWithDefaults()
	addCloud.SetName(name)

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		addCloud.SetGroupId(plan.GroupId.ValueInt64())
	}

	addCloudConfig := addCloud.GetConfig()

	var cloudTypeCode string

	switch {
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		cloudTypeCode = standardCloud

		config := sdkfuncs.NewHvmCloudConfig()

		if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
			config.SetApplianceUrl(plan.ApplianceUrl.ValueString())
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.SetDatacenterName(plan.DataCenterName.ValueString())
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.SetExternalId(plan.ExternalId.ValueString())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.SetInventoryLevel(plan.ImportExistingVms.ValueString())
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.SetConsoleKeymap(plan.KeyboardLayout.ValueString())
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
			config.SetApplianceUrl(plan.ApplianceUrl.ValueString())
		}

		if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
			config.SetDatacenterName(plan.DataCenterName.ValueString())
		}

		if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
			config.SetExternalId(plan.ExternalId.ValueString())
		}

		if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
			config.SetInventoryLevel(plan.ImportExistingVms.ValueString())
		}

		if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
			config.SetConsoleKeymap(plan.KeyboardLayout.ValueString())
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

	addCloud.SetConfig(addCloudConfig)

	cloudTypeWithCode := sdk.NewAddCloudsRequestZoneZoneTypeAnyOf1()
	cloudTypeWithCode.SetCode(cloudTypeCode)

	cloudType := addCloud.GetZoneType()
	cloudType.AddCloudsRequestZoneZoneTypeAnyOf1 = cloudTypeWithCode

	addCloud.SetZoneType(cloudType)

	addCloud.AdditionalProperties = make(map[string]any)

	addCloud.SetAccountId(tenantID)

	if !plan.AgentInstallMode.IsNull() && !plan.AgentInstallMode.IsUnknown() {
		addCloud.SetAgentMode(plan.AgentInstallMode.ValueString())
	}

	if !plan.AutoRecoverPowerState.IsNull() && !plan.AutoRecoverPowerState.IsUnknown() {
		addCloud.SetAutoRecoverPowerState(plan.AutoRecoverPowerState.ValueBool())
	}

	if !plan.Code.IsNull() && !plan.Code.IsUnknown() {
		addCloud.SetCode(plan.Code.ValueString())
	}

	if !plan.CostingMode.IsNull() && !plan.CostingMode.IsUnknown() {
		addCloud.AdditionalProperties["costingMode"] = plan.CostingMode.ValueString()
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		addCloud.SetEnabled(plan.Enabled.ValueBool())
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

		addCloud.SetLabels(labels)
	}

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		addCloud.SetLocation(plan.Location.ValueString())
	}

	if !plan.SecurityMode.IsNull() && !plan.SecurityMode.IsUnknown() {
		addCloud.SetSecurityMode(plan.SecurityMode.ValueString())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		addCloud.SetVisibility(plan.Visibility.ValueString())
	}

	createRequest := sdk.NewAddCloudsRequest(*addCloud)

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

	id := *cloud.GetZone().Id
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
