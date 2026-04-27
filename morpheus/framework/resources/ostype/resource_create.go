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

	osType := sdk.NewAddOsTypesRequestOsType(
		plan.Name.ValueString(),
		plan.Platform.ValueString(),
		plan.Code.ValueString(),
		plan.BitCount.ValueInt64(),
	)

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

	// vendor
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		osType.SetVendor(plan.Vendor.ValueString())
	}

	addReq := sdk.NewAddOsTypesRequestWithDefaults()
	addReq.SetOsType(*osType)

	createResp, httpResp, err := client.LibraryAPI.AddOsTypes(ctx).
		AddOsTypesRequest(*addReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating os type",
			"os type "+plan.Name.ValueString()+" POST failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	createdID := createResp.GetId()

	state, diag := getOsTypeAsState(ctx, createdID, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
