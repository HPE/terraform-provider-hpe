package option_list

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
	_ resource.Resource                = &optionListResource{}
	_ resource.ResourceWithConfigure   = &optionListResource{}
	_ resource.ResourceWithImportState = &optionListResource{}
)

type optionListResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &optionListResource{}
}

func (r *optionListResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_option_list"
}

func (r *optionListResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OptionListSchema(ctx)
}

func (r *optionListResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan optionListModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddOptionListRequestOptionTypeList{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.SourceURL.IsNull() && !plan.SourceURL.IsUnknown() {
		body.SourceUrl = plan.SourceURL.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.ApiType.IsNull() && !plan.ApiType.IsUnknown() {
		body.ApiType = *sdk.NewNullableString(plan.ApiType.ValueStringPointer())
	}
	if !plan.RealTime.IsNull() && !plan.RealTime.IsUnknown() {
		body.RealTime = plan.RealTime.ValueBoolPointer()
	}

	_, httpResp, err := client.LibraryAPI.AddOptionList(ctx).AddOptionListRequest(sdk.AddOptionListRequest{
		OptionTypeList: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "option_list",
			plan.Name.ValueString(), err, httpResp)

		return
	}

	// Find the created resource by listing
	listResult, httpResp, err := client.LibraryAPI.ListOptionLists(ctx).Name(plan.Name.ValueString()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"Read Error After Create",
			"Option type list was created successfully but could not be read back. "+
				"The resource may exist in Morpheus. Check the Morpheus UI and import manually if needed: "+
				"'terraform import <resource_type>.<name> <id>'",
		)

		return
	}

	// SDK field mismatch: API returns "optionTypeLists" but SDK expects "optionTypes"
	optionLists := listResult.GetOptionTypes()
	if len(optionLists) == 0 {
		if rawLists, ok := listResult.AdditionalProperties["optionTypeLists"]; ok {
			if listsSlice, ok := rawLists.([]interface{}); ok && len(listsSlice) > 0 {
				if firstMap, ok := listsSlice[0].(map[string]interface{}); ok {
					mapOptionTypeListFromRaw(&plan, firstMap)
					resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

					return
				}
			}
		}
		resp.Diagnostics.AddError(
			"Not Found After Create",
			"Option type list was created successfully but could not be found by name. "+
				"The resource may exist in Morpheus. Check the Morpheus UI and import manually if needed: "+
				"'terraform import <resource_type>.<name> <id>'",
		)

		return
	}

	ol := optionLists[0]
	mapListOptionListToModel(&plan, &ol)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *optionListResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state optionListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetOptionList(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "option_list", "", err, httpResp)

		return
	}

	optionTypes := result.GetOptionTypes()
	if len(optionTypes) == 0 {
		// SDK field mismatch: API returns "optionTypeList" but SDK expects "optionTypes"
		if rawOL, ok := result.AdditionalProperties["optionTypeList"]; ok {
			if olMap, ok := rawOL.(map[string]interface{}); ok {
				mapOptionTypeListFromRaw(&state, olMap)
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

				return
			}
		}
		resp.State.RemoveResource(ctx)

		return
	}

	ol := optionTypes[0]
	mapGetOptionListToModel(&state, &ol)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *optionListResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan optionListModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateOptionListRequestOptionTypeList{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.SourceURL.IsNull() && !plan.SourceURL.IsUnknown() {
		body.SourceUrl = plan.SourceURL.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.ApiType.IsNull() && !plan.ApiType.IsUnknown() {
		body.ApiType = *sdk.NewNullableString(plan.ApiType.ValueStringPointer())
	}
	if !plan.RealTime.IsNull() && !plan.RealTime.IsUnknown() {
		body.RealTime = plan.RealTime.ValueBoolPointer()
	}

	_, httpResp, err := client.LibraryAPI.UpdateOptionList(ctx, id).UpdateOptionListRequest(sdk.UpdateOptionListRequest{
		OptionTypeList: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "option_list",
			plan.Name.ValueString(), err, httpResp)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *optionListResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state optionListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteOptionList(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "option_list", "", err, httpResp)

		return
	}
}

func (r *optionListResource) ImportState(
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

func mapListOptionListToModel(
	model *optionListModel,
	ol *sdk.ListOptionLists200ResponseAllOfOptionTypesInner,
) {
	if ol.Id != nil {
		model.ID = types.Int64Value(*ol.Id)
	}
	if ol.Name != nil {
		model.Name = types.StringValue(*ol.Name)
	}
	if v := ol.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ol.Type != nil {
		model.Type = types.StringValue(*ol.Type)
	} else {
		model.Type = types.StringNull()
	}
	if ol.SourceUrl != nil {
		model.SourceURL = types.StringValue(*ol.SourceUrl)
	} else {
		model.SourceURL = types.StringNull()
	}
	if ol.Visibility != nil {
		model.Visibility = types.StringValue(*ol.Visibility)
	}
	if v := ol.ApiType.Get(); v != nil {
		model.ApiType = types.StringValue(*v)
	} else {
		model.ApiType = types.StringNull()
	}
	if ol.RealTime != nil {
		model.RealTime = types.BoolValue(*ol.RealTime)
	}
}

func mapGetOptionListToModel(model *optionListModel, ol *sdk.GetOptionList200ResponseOptionTypesInner) {
	if ol.Id != nil {
		model.ID = types.Int64Value(*ol.Id)
	}
	if ol.Name != nil {
		model.Name = types.StringValue(*ol.Name)
	}
	if v := ol.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ol.Type != nil {
		model.Type = types.StringValue(*ol.Type)
	} else {
		model.Type = types.StringNull()
	}
	if ol.SourceUrl != nil {
		model.SourceURL = types.StringValue(*ol.SourceUrl)
	} else {
		model.SourceURL = types.StringNull()
	}
	if ol.Visibility != nil {
		model.Visibility = types.StringValue(*ol.Visibility)
	}
	if v := ol.ApiType.Get(); v != nil {
		model.ApiType = types.StringValue(*v)
	} else {
		model.ApiType = types.StringNull()
	}
	if ol.RealTime != nil {
		model.RealTime = types.BoolValue(*ol.RealTime)
	}
}

// mapOptionTypeListFromRaw maps a raw JSON map to the model (fallback for SDK field mismatch).
func mapOptionTypeListFromRaw(model *optionListModel, m map[string]interface{}) {
	if v, ok := m["id"].(float64); ok {
		model.ID = types.Int64Value(int64(v))
	}
	if v, ok := m["name"].(string); ok {
		model.Name = types.StringValue(v)
	}
	if v, ok := m["description"]; ok && v != nil {
		if s, ok := v.(string); ok {
			model.Description = types.StringValue(s)
		} else {
			model.Description = types.StringNull()
		}
	} else {
		model.Description = types.StringNull()
	}
	if v, ok := m["type"].(string); ok {
		model.Type = types.StringValue(v)
	} else {
		model.Type = types.StringNull()
	}
	if v, ok := m["sourceUrl"].(string); ok {
		model.SourceURL = types.StringValue(v)
	} else {
		model.SourceURL = types.StringNull()
	}
	if v, ok := m["visibility"].(string); ok {
		model.Visibility = types.StringValue(v)
	}
	if v, ok := m["apiType"]; ok && v != nil {
		if s, ok := v.(string); ok {
			model.ApiType = types.StringValue(s)
		} else {
			model.ApiType = types.StringNull()
		}
	} else {
		model.ApiType = types.StringNull()
	}
	if v, ok := m["realTime"].(bool); ok {
		model.RealTime = types.BoolValue(v)
	}
}
