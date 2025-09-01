// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config CloudModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	updateCloud := sdk.NewUpdateCloudsRequestZoneWithDefaults()
	updateCloud.AdditionalProperties = make(map[string]any)
	updateCloud.SetName(name)

	if !plan.AgentInstallMode.IsNull() {
		updateCloud.AdditionalProperties["agentMode"] = plan.AgentInstallMode.ValueString()
	}

	if !plan.AutoRecoverPowerState.IsNull() {
		updateCloud.SetAutoRecoverPowerState(plan.AutoRecoverPowerState.ValueBool())
	}

	if !plan.Code.IsNull() {
		updateCloud.SetCode(plan.Code.ValueString())
	}

	if !plan.CostingMode.IsNull() {
		updateCloud.AdditionalProperties["costingMode"] = plan.CostingMode.ValueString()
	}

	if !plan.Enabled.IsNull() {
		updateCloud.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.ExternalId.IsNull() {
		updateCloud.AdditionalProperties["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.GuidanceMode.IsNull() {
		updateCloud.AdditionalProperties["guidanceMode"] = plan.GuidanceMode.ValueString()
	}

	if plan.Labels.IsNull() {
		updateCloud.SetLabels([]string{})
	} else {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"update cloud resource",
				"cloud "+name+": failed to parse labels: "+err.Error(),
			)

			return
		}

		updateCloud.SetLabels(labels)
	}

	if !plan.Location.IsNull() {
		updateCloud.SetLocation(plan.Location.ValueString())
	}

	if !plan.SecurityMode.IsNull() {
		updateCloud.SetSecurityMode(plan.SecurityMode.ValueString())
	}

	if !plan.Visibility.IsNull() {
		updateCloud.SetVisibility(plan.Visibility.ValueString())
	}

	// TODO: We have to revisit the spec to fix this
	updateCloud.Config = make(map[string]any)

	if !plan.ApplianceUrl.IsNull() {
		updateCloud.Config["applianceUrl"] = plan.ApplianceUrl.ValueString()
	}

	if !plan.DataCenterName.IsNull() {
		updateCloud.Config["datacenterName"] = plan.DataCenterName.ValueString()
	}

	if !plan.ExternalId.IsNull() {
		updateCloud.Config["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.ImportExistingVms.IsNull() {
		updateCloud.Config["inventoryLevel"] = plan.ImportExistingVms.ValueString()
	}

	if !plan.KeyboardLayout.IsNull() {
		updateCloud.Config["consoleKeymap"] = plan.KeyboardLayout.ValueString()
	}

	if !plan.ConfigHvm.CertificateProvider.IsNull() {
		updateCloud.Config["certificateProvider"] = plan.ConfigHvm.CertificateProvider.ValueString()
	}

	if !plan.ConfigHvm.EnableNetworkTypeSelection.IsNull() {
		updateCloud.Config["enableNetworkTypeSelection"] = plan.ConfigHvm.EnableNetworkTypeSelection.ValueBool()
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update cloud resource",
			"cloud "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	updateCloudReq := sdk.NewUpdateCloudsRequest(*updateCloud)

	cloud, hresp, err := client.CloudsAPI.UpdateClouds(ctx, id).
		UpdateCloudsRequest(*updateCloudReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update cloud resource",
			"cloud "+name+" PUT failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if cloud.GetZone().Id == nil {
		resp.Diagnostics.AddError(
			"update cloud resource",
			"cloud "+name+": id is nil",
		)

		return
	}

	newid := *cloud.GetZone().Id
	if newid != id {
		resp.Diagnostics.AddError(
			"update cloud resource",
			"cloud "+name+": id mismatch "+fmt.Sprintf("%d != %d", id, newid),
		)

		return
	}

	state, pdiags := getCloudAsState(ctx, newid, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"update cloud resource",
			fmt.Sprintf("cloud %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
