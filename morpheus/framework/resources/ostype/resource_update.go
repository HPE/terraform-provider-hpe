package ostype

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
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

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	osType := &sdk.UpdateOsTypeRequestOsType{}

	// bit_count
	if !plan.BitCount.IsNull() && !plan.BitCount.IsUnknown() {
		osType.BitCount = plan.BitCount.ValueInt64Pointer()
	}

	// category
	clearable.NullableString(pending, "category", plan.Category, &osType.Category)

	// cloud_init_version
	clearable.NullableString(pending, "cloud_init_version", plan.CloudInitVersion, &osType.CloudInitVersion)

	// description
	clearable.NullableString(pending, "description", plan.Description, &osType.Description)

	// install_agent
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		osType.InstallAgent.Set(plan.InstallAgent.ValueBoolPointer())
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		osType.Name = plan.Name.ValueStringPointer()
	}

	// os_codename
	clearable.NullableString(pending, "os_codename", plan.OsCodename, &osType.OsCodename)

	// os_family
	clearable.NullableString(pending, "os_family", plan.OsFamily, &osType.OsFamily)

	// os_name
	clearable.NullableString(pending, "os_name", plan.OsName, &osType.OsName)

	// os_version
	clearable.NullableString(pending, "os_version", plan.OsVersion, &osType.OsVersion)

	// platform
	if !plan.Platform.IsNull() && !plan.Platform.IsUnknown() {
		osType.Platform = plan.Platform.ValueStringPointer()
	}

	// vendor
	clearable.NullableString(pending, "vendor", plan.Vendor, &osType.Vendor)

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

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
