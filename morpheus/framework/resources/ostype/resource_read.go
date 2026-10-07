package ostype

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data OsTypeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, found, diag := getOsTypeAsState(ctx, data.Id.ValueInt64(), client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	// The OS type was deleted out of band (Morpheus returned 404); remove it
	// from state so the next plan recreates it instead of erroring.
	if !found {
		tflog.Info(ctx, fmt.Sprintf("os type %d not found, removing from state", data.Id.ValueInt64()))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getOsTypeAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (OsTypeModel, bool, diag.Diagnostics) {
	var state OsTypeModel
	var diags diag.Diagnostics

	osTypeResp, httpResp, err := client.LibraryAPI.GetOsType(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return state, false, diags
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading os type",
			fmt.Sprintf("os type %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	osType := osTypeResp.OsType
	if osType == nil {
		diags.AddError(
			"error reading os type",
			fmt.Sprintf("os type %d GET returned no osType", id),
		)

		return state, false, diags
	}

	state.Id = convert.Int64ToType(osType.Id)
	state.Code = convert.StrToType(osType.Code)
	state.BitCount = convert.Int64ToType(osType.BitCount)
	state.Category = convert.StrToType(osType.Category.Get())
	state.CloudInitVersion = convert.StrToType(osType.CloudInitVersion)
	state.Description = convert.StrToType(osType.Description.Get())
	state.InstallAgent = convert.BoolToType(osType.InstallAgent.Get())
	state.Name = convert.StrToType(osType.Name)
	state.OsCodename = convert.StrToType(osType.OsCodename.Get())
	state.OsFamily = convert.StrToType(osType.OsFamily.Get())
	state.OsName = convert.StrToType(osType.OsName.Get())
	state.OsVersion = convert.StrToType(osType.OsVersion.Get())
	state.Platform = convert.StrToType(osType.Platform)
	state.Vendor = convert.StrToType(osType.Vendor.Get())

	// owner is a oneOf (string | object with name); extract as string
	if osType.Owner != nil {
		if osType.Owner.String != nil {
			state.Owner = types.StringValue(*osType.Owner.String)
		} else if osType.Owner.GetOsType200ResponseOsTypeOwnerOneOf != nil {
			state.Owner = convert.StrToType(osType.Owner.GetOsType200ResponseOsTypeOwnerOneOf.Name)
		} else {
			state.Owner = types.StringNull()
		}
	} else {
		state.Owner = types.StringNull()
	}

	return state, true, diags
}
