package ostype

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan OsTypeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	osType := &sdk.AddOsTypesRequestOsType{
		Name:     plan.Name.ValueString(),
		Platform: plan.Platform.ValueString(),
		Code:     plan.Code.ValueString(),
		BitCount: plan.BitCount.ValueInt64(),
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

	// vendor
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		osType.Vendor.Set(plan.Vendor.ValueStringPointer())
	}

	addReq := sdk.NewAddOsTypesRequestWithDefaults()
	addReq.OsType = osType

	createResp, httpResp, err := client.LibraryAPI.AddOsTypes(ctx).
		AddOsTypesRequest(*addReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating os type",
			"os type "+plan.Name.ValueString()+" POST failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	createdID := int64(0)
	if createResp.Id.Get() != nil {
		createdID = *createResp.Id.Get()
	}

	state, diag := getOsTypeAsState(ctx, createdID, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
