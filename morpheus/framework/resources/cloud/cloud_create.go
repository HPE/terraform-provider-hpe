// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/sdkfuncs"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const defaultCloudType = "standard"

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
	groupID := plan.GroupId.ValueInt64()
	tenantID := plan.TenantId.ValueInt64()

	var config CloudModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addCloud := sdk.NewAddCloudsRequestZoneWithDefaults()
	addCloud.SetName(name)
	addCloud.SetGroupId(groupID)

	addCloudConfig := addCloud.GetConfig()

	cloudTypeCode := defaultCloudType

	configAdditionalProperties := make(map[string]any)

	// TODO: update the openapi spec to get SDK fields for these values
	if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
		configAdditionalProperties["applianceUrl"] = plan.ApplianceUrl.ValueStringPointer()
	}

	if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
		configAdditionalProperties["datacenterName"] = plan.DataCenterName.ValueStringPointer()
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		configAdditionalProperties["externalId"] = plan.ExternalId.ValueStringPointer()
	}

	if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
		configAdditionalProperties["inventoryLevel"] = plan.ImportExistingVms.ValueStringPointer()
	}

	if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
		configAdditionalProperties["consoleKeymap"] = plan.KeyboardLayout.ValueStringPointer()
	}

	switch {
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		config := sdkfuncs.NewHvmCloudConfig()

		config.AdditionalProperties = configAdditionalProperties

		if !plan.ConfigHvm.CertificateProvider.IsNull() && !plan.ConfigHvm.CertificateProvider.IsUnknown() {
			config.CertificateProvider = plan.ConfigHvm.CertificateProvider.ValueStringPointer()
		}
		if !plan.ConfigHvm.EnableNetworkTypeSelection.IsNull() && !plan.ConfigHvm.EnableNetworkTypeSelection.IsUnknown() {
			config.EnableNetworkTypeSelection = plan.ConfigHvm.EnableNetworkTypeSelection.ValueBoolPointer()
		}

		addHvmConfig := sdkfuncs.AddCloudRequestHVMConfig(config)

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf = &addHvmConfig

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create cloud resource",
				"cloud "+name+": failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"create cloud resource",
				"cloud "+name+": config must be a valid object/map",
			)

			return
		}

		maps.Copy(configDataMap, configAdditionalProperties)

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

	if !plan.AgentInstallMode.IsNull() {
		addCloud.SetAgentMode(plan.AgentInstallMode.ValueString())
	}

	if !plan.AutoRecoverPowerState.IsNull() {
		addCloud.SetAutoRecoverPowerState(plan.AutoRecoverPowerState.ValueBool())
	}

	if !plan.Code.IsNull() {
		addCloud.SetCode(plan.Code.ValueString())
	}

	if !plan.CostingMode.IsNull() {
		addCloud.AdditionalProperties["costingMode"] = plan.CostingMode.ValueString()
	}

	if !plan.Enabled.IsNull() {
		addCloud.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.ExternalId.IsNull() {
		addCloud.AdditionalProperties["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.GuidanceMode.IsNull() {
		addCloud.AdditionalProperties["guidanceMode"] = plan.GuidanceMode.ValueString()
	}

	if !plan.Labels.IsNull() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"create cloud resource",
				"cloud "+name+": failed to parse label: "+err.Error(),
			)

			return
		}

		addCloud.SetLabels(labels)
	}

	if !plan.Location.IsNull() {
		addCloud.SetLocation(plan.Location.ValueString())
	}

	if !plan.SecurityMode.IsNull() {
		addCloud.SetSecurityMode(plan.SecurityMode.ValueString())
	}

	if !plan.Visibility.IsNull() {
		addCloud.SetVisibility(plan.Visibility.ValueString())
	}

	createRequest := sdk.NewAddCloudsRequest(*addCloud)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create cloud resource",
			"cloud "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	// Call API
	cloud, hresp, err := client.CloudsAPI.AddClouds(ctx).AddCloudsRequest(*createRequest).Execute()
	if cloud == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create cloud resource",
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
			"failed to read cloud state",
			fmt.Sprintf("Cloud %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set cloud state",
			fmt.Sprintf("Cloud %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}
