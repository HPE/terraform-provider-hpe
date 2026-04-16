package ostype

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

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

	state, diag := getOsTypeAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getOsTypeAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	prior OsTypeModel,
) (OsTypeModel, diag.Diagnostics) {
	var state OsTypeModel
	var diags diag.Diagnostics

	osTypeResp, httpResp, err := client.LibraryAPI.GetOsType(ctx, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading os type",
			fmt.Sprintf("os type %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	osType := osTypeResp.GetOsType()

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

	// code may not be returned in the GET response; preserve from prior state if set
	if !prior.Code.IsNull() && !prior.Code.IsUnknown() {
		state.Code = prior.Code
	}

	return state, diags
}
