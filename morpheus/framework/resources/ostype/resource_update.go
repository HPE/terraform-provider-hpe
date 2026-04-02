package ostype

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState OsTypeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	osType := sdk.NewUpdateOsTypeRequestOsTypeWithDefaults()

	// bit_count
	if !plan.BitCount.IsNull() && !plan.BitCount.IsUnknown() {
		osType.SetBitCount(plan.BitCount.ValueInt64())
	}

	// category
	if !plan.Category.IsNull() && !plan.Category.IsUnknown() {
		osType.SetCategory(plan.Category.ValueString())
	}

	// cloud_init_version
	if !plan.CloudInitVersion.IsNull() && !plan.CloudInitVersion.IsUnknown() {
		osType.SetCloudInitVersion(plan.CloudInitVersion.ValueString())
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		osType.SetDescription(plan.Description.ValueString())
	}

	// install_agent
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		osType.SetInstallAgent(plan.InstallAgent.ValueBool())
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		osType.SetName(plan.Name.ValueString())
	}

	// os_codename
	if !plan.OsCodename.IsNull() && !plan.OsCodename.IsUnknown() {
		osType.SetOsCodename(plan.OsCodename.ValueString())
	}

	// os_family
	if !plan.OsFamily.IsNull() && !plan.OsFamily.IsUnknown() {
		osType.SetOsFamily(plan.OsFamily.ValueString())
	}

	// os_name
	if !plan.OsName.IsNull() && !plan.OsName.IsUnknown() {
		osType.SetOsName(plan.OsName.ValueString())
	}

	// os_version
	if !plan.OsVersion.IsNull() && !plan.OsVersion.IsUnknown() {
		osType.SetOsVersion(plan.OsVersion.ValueString())
	}

	// platform
	if !plan.Platform.IsNull() && !plan.Platform.IsUnknown() {
		osType.SetPlatform(plan.Platform.ValueString())
	}

	// vendor
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		osType.SetVendor(plan.Vendor.ValueString())
	}

	updateReq := sdk.NewUpdateOsTypeRequestWithDefaults()
	updateReq.SetOsType(*osType)

	id := currentState.Id.ValueInt64()

	_, httpResp, err := client.LibraryAPI.UpdateOsType(ctx, id).
		UpdateOsTypeRequest(*updateReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating os type",
			"os type "+plan.Name.ValueString()+" PUT failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, diag := getOsTypeAsState(ctx, id, client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
