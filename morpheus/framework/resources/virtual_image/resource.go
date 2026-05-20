package virtual_image

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &virtualImageResource{}
	_ resource.ResourceWithConfigure   = &virtualImageResource{}
	_ resource.ResourceWithImportState = &virtualImageResource{}
)

type virtualImageResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &virtualImageResource{}
}

func (r *virtualImageResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_virtual_image"
}

func (r *virtualImageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VirtualImageSchema(ctx)
}

func (r *virtualImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan virtualImageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddVirtualImageRequestVirtualImage{
		Name:      plan.Name.ValueStringPointer(),
		ImageType: plan.ImageType.ValueStringPointer(),
	}
	if !plan.OsTypeId.IsNull() && !plan.OsTypeId.IsUnknown() {
		body.OsType = *sdk.NewNullableInt64(plan.OsTypeId.ValueInt64Pointer())
	}
	if !plan.IsCloudInit.IsNull() && !plan.IsCloudInit.IsUnknown() {
		body.IsCloudInit = plan.IsCloudInit.ValueBoolPointer()
	}
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		body.InstallAgent = plan.InstallAgent.ValueBoolPointer()
	}
	if !plan.MinRAM.IsNull() && !plan.MinRAM.IsUnknown() {
		body.MinRam = *sdk.NewNullableInt64(plan.MinRAM.ValueInt64Pointer())
	}
	if !plan.MinDisk.IsNull() && !plan.MinDisk.IsUnknown() {
		body.MinDisk = *sdk.NewNullableInt64(plan.MinDisk.ValueInt64Pointer())
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}

	createResp, httpResp, err := client.LibraryAPI.AddVirtualImage(ctx).AddVirtualImageRequest(sdk.AddVirtualImageRequest{
		VirtualImage: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "virtual_image", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Extract ID from create response
	var id int64
	if vi := createResp.GetVirtualImage(); vi.Id != nil {
		id = *vi.Id
	} else if createResp.AdditionalProperties != nil {
		if idVal, ok := createResp.AdditionalProperties["id"]; ok {
			switch v := idVal.(type) {
			case float64:
				id = int64(v)
			case int64:
				id = v
			}
		}
	}
	if id == 0 {
		resp.Diagnostics.AddError("Create Error", "Could not determine ID from create response")

		return
	}

	// Read back
	result, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "virtual_image", "", err, httpResp)

		return
	}

	vi := result.GetVirtualImage()
	mapGetVirtualImageToModel(ctx, &plan, &vi, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *virtualImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state virtualImageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "virtual_image", "", err, httpResp)

		return
	}

	vi := result.GetVirtualImage()
	mapGetVirtualImageToModel(ctx, &state, &vi, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *virtualImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan virtualImageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateVirtualImageRequestVirtualImage{
		Name:      plan.Name.ValueStringPointer(),
		ImageType: plan.ImageType.ValueStringPointer(),
	}
	if !plan.OsTypeId.IsNull() && !plan.OsTypeId.IsUnknown() {
		body.OsType = *sdk.NewNullableInt64(plan.OsTypeId.ValueInt64Pointer())
	}
	if !plan.IsCloudInit.IsNull() && !plan.IsCloudInit.IsUnknown() {
		body.IsCloudInit = plan.IsCloudInit.ValueBoolPointer()
	}
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		body.InstallAgent = plan.InstallAgent.ValueBoolPointer()
	}
	if !plan.MinDisk.IsNull() && !plan.MinDisk.IsUnknown() {
		body.MinDisk = *sdk.NewNullableInt64(plan.MinDisk.ValueInt64Pointer())
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.MinRAM.IsNull() && !plan.MinRAM.IsUnknown() {
		if body.AdditionalProperties == nil {
			body.AdditionalProperties = make(map[string]interface{})
		}
		body.AdditionalProperties["minRam"] = plan.MinRAM.ValueInt64()
	}

	_, httpResp, err := client.LibraryAPI.UpdateVirtualImage(ctx, id).
		UpdateVirtualImageRequest(sdk.UpdateVirtualImageRequest{
			VirtualImage: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "virtual_image", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Read back
	result, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "virtual_image", "", err, httpResp)

		return
	}

	vi := result.GetVirtualImage()
	mapGetVirtualImageToModel(ctx, &plan, &vi, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *virtualImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state virtualImageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.RemoveVirtualImage(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "virtual_image", "", err, httpResp)

		return
	}
}

func (r *virtualImageResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetVirtualImageToModel(
	ctx context.Context,
	model *virtualImageModel,
	vi *sdk.GetVirtualImage200ResponseVirtualImage,
	diags *diag.Diagnostics,
) {
	if vi.Id != nil {
		model.ID = types.Int64Value(*vi.Id)
	}
	if vi.Name != nil {
		model.Name = types.StringValue(*vi.Name)
	}
	if vi.ImageType != nil {
		model.ImageType = types.StringValue(*vi.ImageType)
	}
	if vi.OsType != nil && vi.OsType.Id != nil {
		model.OsTypeId = types.Int64Value(*vi.OsType.Id)
	} else if model.OsTypeId.IsNull() || model.OsTypeId.IsUnknown() {
		model.OsTypeId = types.Int64Null()
	}
	if vi.IsCloudInit != nil {
		model.IsCloudInit = types.BoolValue(*vi.IsCloudInit)
	}
	if vi.InstallAgent != nil {
		model.InstallAgent = types.BoolValue(*vi.InstallAgent)
	}
	if vi.MinRam.IsSet() && vi.MinRam.Get() != nil {
		model.MinRAM = types.Int64Value(*vi.MinRam.Get())
	} else if model.MinRAM.IsNull() || model.MinRAM.IsUnknown() {
		model.MinRAM = types.Int64Null()
	}
	if vi.MinDisk.IsSet() && vi.MinDisk.Get() != nil {
		model.MinDisk = types.Int64Value(*vi.MinDisk.Get())
	} else if model.MinDisk.IsNull() || model.MinDisk.IsUnknown() {
		model.MinDisk = types.Int64Null()
	}
	if vi.Labels != nil {
		labelValues, d := types.ListValueFrom(ctx, types.StringType, vi.Labels)
		diags.Append(d...)
		model.Labels = labelValues
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
}
