// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const updateOperation = "update cloud resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
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

	// This won't do anything with the current API.
	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		updateCloud.SetGroupId(plan.GroupId.ValueInt64())
	}

	if !plan.AutoRecoverPowerState.IsNull() && !plan.AutoRecoverPowerState.IsUnknown() {
		updateCloud.SetAutoRecoverPowerState(plan.AutoRecoverPowerState.ValueBool())
	}

	if !plan.Code.IsNull() && !plan.Code.IsUnknown() {
		updateCloud.SetCode(plan.Code.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		updateCloud.SetEnabled(plan.Enabled.ValueBool())
	}

	if plan.Labels.IsNull() || plan.Labels.IsUnknown() {
		updateCloud.SetLabels([]string{})
	} else {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				updateOperation,
				"cloud "+name+": failed to parse labels: "+err.Error(),
			)

			return
		}

		updateCloud.SetLabels(labels)
	}

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		updateCloud.SetLocation(plan.Location.ValueString())
	}

	if !plan.SecurityMode.IsNull() && !plan.SecurityMode.IsUnknown() {
		updateCloud.SetSecurityMode(plan.SecurityMode.ValueString())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		updateCloud.SetVisibility(plan.Visibility.ValueString())
	}

	// TODO: Update spec to generate setters
	if !plan.AgentInstallMode.IsNull() && !plan.AgentInstallMode.IsUnknown() {
		updateCloud.AdditionalProperties["agentMode"] = plan.AgentInstallMode.ValueString()
	}

	if !plan.CostingMode.IsNull() && !plan.CostingMode.IsUnknown() {
		updateCloud.AdditionalProperties["costingMode"] = plan.CostingMode.ValueString()
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		updateCloud.AdditionalProperties["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.GuidanceMode.IsNull() && !plan.GuidanceMode.IsUnknown() {
		updateCloud.AdditionalProperties["guidanceMode"] = plan.GuidanceMode.ValueString()
	}

	updateCloud.Config = make(map[string]any)

	if !plan.ApplianceUrl.IsNull() && !plan.ApplianceUrl.IsUnknown() {
		updateCloud.Config["applianceUrl"] = plan.ApplianceUrl.ValueString()
	}

	if !plan.DataCenterName.IsNull() && !plan.DataCenterName.IsUnknown() {
		updateCloud.Config["datacenterName"] = plan.DataCenterName.ValueString()
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		updateCloud.Config["externalId"] = plan.ExternalId.ValueString()
	}

	if !plan.ImportExistingVms.IsNull() && !plan.ImportExistingVms.IsUnknown() {
		updateCloud.Config["inventoryLevel"] = plan.ImportExistingVms.ValueString()
	}

	if !plan.KeyboardLayout.IsNull() && !plan.KeyboardLayout.IsUnknown() {
		updateCloud.Config["consoleKeymap"] = plan.KeyboardLayout.ValueString()
	}

	switch {
	case !plan.ConfigAws.IsNull() && !plan.ConfigAws.IsUnknown():
		if !plan.ConfigAws.AccessKey.IsNull() && !plan.ConfigAws.AccessKey.IsUnknown() {
			updateCloud.Config["accessKey"] = plan.ConfigAws.AccessKey.ValueString()
		}

		if !plan.ConfigAws.ApiProxy.IsNull() && !plan.ConfigAws.ApiProxy.IsUnknown() {
			updateCloud.Config["apiProxy"] = plan.ConfigAws.ApiProxy.ValueString()
		}

		if !plan.ConfigAws.BackupProvider.IsNull() && !plan.ConfigAws.BackupProvider.IsUnknown() {
			updateCloud.Config["backupMode"] = plan.ConfigAws.BackupProvider.ValueString()
		}

		if !plan.ConfigAws.BypassProxy.IsNull() && !plan.ConfigAws.BypassProxy.IsUnknown() {
			updateCloud.Config["bypassProxyForCloud"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigAws.BypassProxy)
		}

		if !plan.ConfigAws.ChangeManagementConfig.IsNull() && !plan.ConfigAws.ChangeManagementConfig.IsUnknown() {
			updateCloud.Config["changeManagementConfig"] = plan.ConfigAws.ChangeManagementConfig.ValueString()
		}

		if !plan.ConfigAws.CmdbConfig.IsNull() && !plan.ConfigAws.CmdbConfig.IsUnknown() {
			updateCloud.Config["cmdbConfig"] = plan.ConfigAws.CmdbConfig.ValueString()
		}

		if !plan.ConfigAws.CmdbDiscovery.IsNull() && !plan.ConfigAws.CmdbDiscovery.IsUnknown() {
			updateCloud.Config["configCmdbDiscovery"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigAws.CmdbDiscovery)
		}

		if !plan.ConfigAws.ConfigManagementId.IsNull() && !plan.ConfigAws.ConfigManagementId.IsUnknown() {
			updateCloud.Config["configManagementId"] = plan.ConfigAws.ConfigManagementId.ValueString()
		}

		if !plan.ConfigAws.Costing.IsNull() && !plan.ConfigAws.Costing.IsUnknown() {
			updateCloud.Config["costing"] = plan.ConfigAws.Costing.ValueString()
		}

		if !plan.ConfigAws.CostingBucket.IsNull() && !plan.ConfigAws.CostingBucket.IsUnknown() {
			updateCloud.Config["costingBucket"] = plan.ConfigAws.CostingBucket.ValueString()
		}

		if !plan.ConfigAws.CostingFolder.IsNull() && !plan.ConfigAws.CostingFolder.IsUnknown() {
			updateCloud.Config["costingFolder"] = plan.ConfigAws.CostingFolder.ValueString()
		}

		if !plan.ConfigAws.CostingKey.IsNull() && !plan.ConfigAws.CostingKey.IsUnknown() {
			updateCloud.Config["costingKey"] = plan.ConfigAws.CostingKey.ValueString()
		}

		if !plan.ConfigAws.CostingReportName.IsNull() && !plan.ConfigAws.CostingReportName.IsUnknown() {
			updateCloud.Config["costingReportName"] = plan.ConfigAws.CostingReportName.ValueString()
		}

		if !plan.ConfigAws.CostingSecret.IsNull() && !plan.ConfigAws.CostingSecret.IsUnknown() {
			updateCloud.Config["costingSecret"] = plan.ConfigAws.CostingSecret.ValueString()
		}

		if !plan.ConfigAws.Credentials.IsNull() && !plan.ConfigAws.Credentials.IsUnknown() {
			updateCloud.Config["useHostCredentials"] = plan.ConfigAws.Credentials.ValueString()
		}

		if !plan.ConfigAws.DarkModeLogo.IsNull() && !plan.ConfigAws.DarkModeLogo.IsUnknown() {
			updateCloud.Config["darkModeLogo"] = plan.ConfigAws.DarkModeLogo.ValueString()
		}

		if !plan.ConfigAws.Domain.IsNull() && !plan.ConfigAws.Domain.IsUnknown() {
			updateCloud.Config["domain"] = plan.ConfigAws.Domain.ValueString()
		}

		if !plan.ConfigAws.EbsEncryption.IsNull() && !plan.ConfigAws.EbsEncryption.IsUnknown() {
			updateCloud.Config["ebsEncryption"] = plan.ConfigAws.EbsEncryption.ValueString()
		}

		if !plan.ConfigAws.Endpoint.IsNull() && !plan.ConfigAws.Endpoint.IsUnknown() {
			updateCloud.Config["endpoint"] = plan.ConfigAws.Endpoint.ValueString()
		}

		if !plan.ConfigAws.Guidance.IsNull() && !plan.ConfigAws.Guidance.IsUnknown() {
			updateCloud.Config["guidance"] = plan.ConfigAws.Guidance.ValueString()
		}

		if !plan.ConfigAws.Logo.IsNull() && !plan.ConfigAws.Logo.IsUnknown() {
			updateCloud.Config["logo"] = plan.ConfigAws.Logo.ValueString()
		}

		if !plan.ConfigAws.NetworkMode.IsNull() && !plan.ConfigAws.NetworkMode.IsUnknown() {
			updateCloud.Config["networkMode"] = plan.ConfigAws.NetworkMode.ValueString()
		}

		if !plan.ConfigAws.NoProxy.IsNull() && !plan.ConfigAws.NoProxy.IsUnknown() {
			updateCloud.Config["noProxy"] = plan.ConfigAws.NoProxy.ValueString()
		}

		if !plan.ConfigAws.Proxy.IsNull() && !plan.ConfigAws.Proxy.IsUnknown() {
			updateCloud.Config["proxy"] = plan.ConfigAws.Proxy.ValueString()
		}

		if !plan.ConfigAws.Region.IsNull() && !plan.ConfigAws.Region.IsUnknown() {
			updateCloud.Config["region"] = plan.ConfigAws.Region.ValueString()
		}

		if !plan.ConfigAws.ReplicationProvider.IsNull() && !plan.ConfigAws.ReplicationProvider.IsUnknown() {
			updateCloud.Config["replicationMode"] = plan.ConfigAws.ReplicationProvider.ValueString()
		}

		if !plan.ConfigAws.RoleArn.IsNull() && !plan.ConfigAws.RoleArn.IsUnknown() {
			updateCloud.Config["stsAssumeRole"] = plan.ConfigAws.RoleArn.ValueString()
		}

		if !plan.ConfigAws.SecretKey.IsNull() && !plan.ConfigAws.SecretKey.IsUnknown() {
			updateCloud.Config["secretKey"] = plan.ConfigAws.SecretKey.ValueString()
		}

		if !plan.ConfigAws.Timezone.IsNull() && !plan.ConfigAws.Timezone.IsUnknown() {
			updateCloud.Config["timezone"] = plan.ConfigAws.Timezone.ValueString()
		}

		if !plan.ConfigAws.UserData.IsNull() && !plan.ConfigAws.UserData.IsUnknown() {
			updateCloud.Config["userDataLinux"] = plan.ConfigAws.UserData.ValueString()
		}

		if !plan.ConfigAws.VdiGateway.IsNull() && !plan.ConfigAws.VdiGateway.IsUnknown() {
			updateCloud.Config["vdiGateway"] = plan.ConfigAws.VdiGateway.ValueString()
		}

		if !plan.ConfigAws.Vpc.IsNull() && !plan.ConfigAws.Vpc.IsUnknown() {
			updateCloud.Config["vpc"] = plan.ConfigAws.Vpc.ValueString()
		}

	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		if !plan.ConfigHvm.CertificateProvider.IsNull() && !plan.ConfigHvm.CertificateProvider.IsUnknown() {
			updateCloud.Config["certificateProvider"] = plan.ConfigHvm.CertificateProvider.ValueString()
		}

		if !plan.ConfigHvm.EnableNetworkTypeSelection.IsNull() &&
			!plan.ConfigHvm.EnableNetworkTypeSelection.IsUnknown() {
			updateCloud.Config["enableNetworkTypeSelection"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigHvm.EnableNetworkTypeSelection)
		}
	case !plan.ConfigVmware.IsNull() && !plan.ConfigVmware.IsUnknown():
		if !plan.ConfigVmware.ApiUrl.IsNull() && !plan.ConfigVmware.ApiUrl.IsUnknown() {
			updateCloud.Config["apiUrl"] = plan.ConfigVmware.ApiUrl.ValueString()
		}

		if !plan.ConfigVmware.ApiVersion.IsNull() && !plan.ConfigVmware.ApiVersion.IsUnknown() {
			updateCloud.Config["apiVersion"] = plan.ConfigVmware.ApiVersion.ValueString()
		}

		if !plan.ConfigVmware.CertificateProvider.IsNull() && !plan.ConfigVmware.CertificateProvider.IsUnknown() {
			updateCloud.Config["certificateProvider"] = plan.ConfigVmware.CertificateProvider.ValueString()
		}

		if !plan.ConfigVmware.Cluster.IsNull() && !plan.ConfigVmware.Cluster.IsUnknown() {
			updateCloud.Config["cluster"] = plan.ConfigVmware.Cluster.ValueString()
		}

		if !plan.ConfigVmware.ConfigManagementId.IsNull() && !plan.ConfigVmware.ConfigManagementId.IsUnknown() {
			updateCloud.Config["configManagementId"] = plan.ConfigVmware.ConfigManagementId.ValueString()
		}

		if !plan.ConfigVmware.Datacenter.IsNull() && !plan.ConfigVmware.Datacenter.IsUnknown() {
			updateCloud.Config["datacenter"] = plan.ConfigVmware.Datacenter.ValueString()
		}

		if !plan.ConfigVmware.Password.IsNull() && !plan.ConfigVmware.Password.IsUnknown() {
			updateCloud.Config["password"] = plan.ConfigVmware.Password.ValueString()
		}

		if !plan.ConfigVmware.ResourcePool.IsNull() && !plan.ConfigVmware.ResourcePool.IsUnknown() {
			updateCloud.Config["resourcePool"] = plan.ConfigVmware.ResourcePool.ValueString()
		}

		if !plan.ConfigVmware.RpcMode.IsNull() && !plan.ConfigVmware.RpcMode.IsUnknown() {
			updateCloud.Config["rpcMode"] = plan.ConfigVmware.RpcMode.ValueString()
		}

		if !plan.ConfigVmware.Username.IsNull() && !plan.ConfigVmware.Username.IsUnknown() {
			updateCloud.Config["username"] = plan.ConfigVmware.Username.ValueString()
		}

		if !plan.ConfigVmware.EnableDiskTypeSelection.IsNull() && !plan.ConfigVmware.EnableDiskTypeSelection.IsUnknown() {
			updateCloud.Config["enableDiskTypeSelection"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableDiskTypeSelection)
		}

		if !plan.ConfigVmware.EnableNetworkTypeSelection.IsNull() &&
			!plan.ConfigVmware.EnableNetworkTypeSelection.IsUnknown() {
			updateCloud.Config["enableNetworkTypeSelection"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableNetworkTypeSelection)
		}

		if !plan.ConfigVmware.EnableStorageTypeSelection.IsNull() &&
			!plan.ConfigVmware.EnableStorageTypeSelection.IsUnknown() {
			updateCloud.Config["enableStorageTypeSelection"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableStorageTypeSelection)
		}

		if !plan.ConfigVmware.EnableVnc.IsNull() && !plan.ConfigVmware.EnableVnc.IsUnknown() {
			updateCloud.Config["enableVnc"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigVmware.EnableVnc)
		}

		if !plan.ConfigVmware.HideHostSelection.IsNull() && !plan.ConfigVmware.HideHostSelection.IsUnknown() {
			updateCloud.Config["hideHostSelection"] = convert.
				BoolTypeToStringPointerOnOff(plan.ConfigVmware.HideHostSelection)
		}
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		if plan.CloudTypeCode.IsNull() || plan.CloudTypeCode.IsUnknown() {
			resp.Diagnostics.AddError(
				updateOperation,
				"cloud "+name+": cloud_type_code is required for generic configurations",
			)

			return
		}
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			updateOperation,
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
			updateOperation,
			"cloud "+name+" PUT failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if cloud.GetZone().Id == nil {
		resp.Diagnostics.AddError(
			updateOperation,
			"cloud "+name+": id is nil",
		)

		return
	}

	newid := *cloud.GetZone().Id
	if newid != id {
		resp.Diagnostics.AddError(
			updateOperation,
			"cloud "+name+": id mismatch "+fmt.Sprintf("%d != %d", id, newid),
		)

		return
	}

	state, pdiags := getCloudAsState(ctx, newid, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("cloud %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
