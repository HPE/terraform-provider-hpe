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

	osType := &sdk.UpdateOsTypeRequestOsType{}

	// bit_count
	if !plan.BitCount.IsNull() && !plan.BitCount.IsUnknown() {
		osType.BitCount = plan.BitCount.ValueInt64Pointer()
	}

	// category
	if !plan.Category.IsNull() && !plan.Category.IsUnknown() {
		osType.Category.Set(plan.Category.ValueStringPointer())
	}

	// cloud_init_version
	if !plan.CloudInitVersion.IsNull() && !plan.CloudInitVersion.IsUnknown() {
		osType.CloudInitVersion = plan.CloudInitVersion.ValueStringPointer()
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		osType.Description.Set(plan.Description.ValueStringPointer())
	}

	// install_agent
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		osType.InstallAgent.Set(plan.InstallAgent.ValueBoolPointer())
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		osType.Name = plan.Name.ValueStringPointer()
	}

	// os_codename
	if !plan.OsCodename.IsNull() && !plan.OsCodename.IsUnknown() {
		osType.OsCodename.Set(plan.OsCodename.ValueStringPointer())
	}

	// os_family
	if !plan.OsFamily.IsNull() && !plan.OsFamily.IsUnknown() {
		osType.OsFamily.Set(plan.OsFamily.ValueStringPointer())
	}

	// os_name
	if !plan.OsName.IsNull() && !plan.OsName.IsUnknown() {
		osType.OsName.Set(plan.OsName.ValueStringPointer())
	}

	// os_version
	if !plan.OsVersion.IsNull() && !plan.OsVersion.IsUnknown() {
		osType.OsVersion.Set(plan.OsVersion.ValueStringPointer())
	}

	// platform
	if !plan.Platform.IsNull() && !plan.Platform.IsUnknown() {
		osType.Platform = plan.Platform.ValueStringPointer()
	}

	// vendor
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		osType.Vendor.Set(plan.Vendor.ValueStringPointer())
	}

	updateReq := &sdk.UpdateOsTypeRequest{}
	updateReq.OsType = osType

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

	state, diag := getOsTypeAsState(ctx, id, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
