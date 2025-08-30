// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	groupId := plan.GroupId.ValueInt64()
	tenantId := plan.TenantId.ValueInt64()

	var config CloudModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addCloudConfig := sdk.AddCloudsRequestZoneConfig{
		AddCloudsRequestZoneConfigAnyOf: &sdk.AddCloudsRequestZoneConfigAnyOf{},
	}

	cloudTypeCode := "standard"

	switch {
	case !plan.ConfigHvm.IsNull():
		cloudTypeCode = "standard"

		config := sdk.AddCloudsRequestZoneConfigAnyOfOneOf2{}

		if !plan.ConfigHvm.CertificateProvider.IsNull() {
			config.CertificateProvider = plan.ConfigHvm.CertificateProvider.ValueStringPointer()
		}
		if !plan.ConfigHvm.EnableNetworkTypeSelection.IsNull() {
			config.EnableNetworkTypeSelection = plan.ConfigHvm.EnableNetworkTypeSelection.ValueBoolPointer()
		}

		addCloudConfig.AddCloudsRequestZoneConfigAnyOf.AddCloudsRequestZoneConfigAnyOfOneOf2 = &config
	}

	// TODO: support other cloud types
	genConfig := addCloudConfig.AddCloudsRequestZoneConfigAnyOf.AddCloudsRequestZoneConfigAnyOfOneOf2
	if genConfig == nil {
		genConfig = &sdk.AddCloudsRequestZoneConfigAnyOfOneOf2{}
		addCloudConfig.AddCloudsRequestZoneConfigAnyOf.AddCloudsRequestZoneConfigAnyOfOneOf2 = genConfig
	}

	if genConfig.AdditionalProperties == nil {
		genConfig.AdditionalProperties = make(map[string]any)
	}

	if !plan.ApplianceUrl.IsNull() {
		genConfig.AdditionalProperties["applianceUrl"] = plan.ApplianceUrl.ValueStringPointer()
	}

	if !plan.DataCenterName.IsNull() {
		genConfig.AdditionalProperties["datacenterName"] = plan.DataCenterName.ValueStringPointer()
	}

	if !plan.ExternalId.IsNull() {
		genConfig.AdditionalProperties["externalId"] = plan.ExternalId.ValueStringPointer()
	}

	if !plan.ImportExistingVms.IsNull() {
		genConfig.AdditionalProperties["inventoryLevel"] = plan.ImportExistingVms.ValueStringPointer()
	}

	if !plan.KeyboardLayout.IsNull() {
		genConfig.AdditionalProperties["consoleKeymap"] = plan.KeyboardLayout.ValueStringPointer()
	}

	cloudType := sdk.AddCloudsRequestZoneZoneType{
		AddCloudsRequestZoneZoneTypeAnyOf1: &sdk.AddCloudsRequestZoneZoneTypeAnyOf1{
			Code: &cloudTypeCode,
		},
	}

	addCloud := sdk.NewAddCloudsRequestZone(name, groupId, cloudType, addCloudConfig)
	addCloud.AdditionalProperties = make(map[string]any)

	addCloud.SetAccountId(tenantId)

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
	tflog.Info(ctx, "==== createRequest ====", map[string]any{"request": spew.Sdump(createRequest)})

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
			"cloud "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	id := *cloud.GetZone().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getCloudAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create cloud resource",
			fmt.Sprintf("cloud %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
