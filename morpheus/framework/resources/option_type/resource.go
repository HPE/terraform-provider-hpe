package option_type

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &optionTypeResource{}
	_ resource.ResourceWithConfigure   = &optionTypeResource{}
	_ resource.ResourceWithImportState = &optionTypeResource{}
)

type optionTypeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &optionTypeResource{}
}

func (r *optionTypeResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_option_type"
}

func (r *optionTypeResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OptionTypeSchema(ctx)
}

func (r *optionTypeResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan optionTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddOptionTypeRequestOptionType{
		Name: plan.Name.ValueString(),
	}
	if !plan.FieldName.IsNull() {
		body.FieldName = plan.FieldName.ValueStringPointer()
	}
	if !plan.Type.IsNull() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.FieldLabel.IsNull() {
		body.FieldLabel = plan.FieldLabel.ValueStringPointer()
	}
	if !plan.Placeholder.IsNull() {
		body.PlaceHolder = plan.Placeholder.ValueStringPointer()
	}
	if !plan.DefaultValue.IsNull() {
		body.DefaultValue = plan.DefaultValue.ValueStringPointer()
	}
	if !plan.Required.IsNull() {
		body.Required = plan.Required.ValueBoolPointer()
	}
	if !plan.ExportMeta.IsNull() {
		body.ExportMeta = plan.ExportMeta.ValueBoolPointer()
	}

	result, httpResp, err := client.LibraryAPI.AddOptionType(ctx).AddOptionTypeRequest(sdk.AddOptionTypeRequest{
		OptionType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "option_type", plan.Name.ValueString(), err, httpResp)

		return
	}

	if id := result.Id.Get(); id != nil {
		plan.ID = types.Int64Value(*id)
	} else if otData, ok := result.AdditionalProperties["optionType"]; ok {
		if otMap, ok := otData.(map[string]interface{}); ok {
			if idVal, ok := otMap["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					plan.ID = types.Int64Value(int64(v))
				case int64:
					plan.ID = types.Int64Value(v)
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *optionTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// The SDK does not expose a GetOptionType endpoint.
	// State is preserved from Create/Update. Import relies on the ID being set and
	// a subsequent refresh via Update. This is a known SDK limitation.
	var state optionTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *optionTypeResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan optionTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateOptionTypeRequestOptionType{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.FieldName.IsNull() {
		body.FieldName = plan.FieldName.ValueStringPointer()
	}
	if !plan.Type.IsNull() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.FieldLabel.IsNull() {
		body.FieldLabel = plan.FieldLabel.ValueStringPointer()
	}
	if !plan.Placeholder.IsNull() {
		body.PlaceHolder = plan.Placeholder.ValueStringPointer()
	}
	if !plan.DefaultValue.IsNull() {
		body.DefaultValue = plan.DefaultValue.ValueStringPointer()
	}
	if !plan.Required.IsNull() {
		body.Required = plan.Required.ValueBoolPointer()
	}
	if !plan.ExportMeta.IsNull() {
		body.ExportMeta = plan.ExportMeta.ValueBoolPointer()
	}

	_, httpResp, err := client.LibraryAPI.UpdateOptionType(ctx, id).UpdateOptionTypeRequest(sdk.UpdateOptionTypeRequest{
		OptionType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "option_type", plan.Name.ValueString(), err, httpResp)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *optionTypeResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state optionTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteOptionType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "option_type", "", err, httpResp)

		return
	}
}

func (r *optionTypeResource) ImportState(
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
