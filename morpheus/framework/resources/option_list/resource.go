package option_list

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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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
	resp.Schema = OptionListResourceSchema(ctx)
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

	var plan OptionListModel
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
	if !plan.SourceUrl.IsNull() && !plan.SourceUrl.IsUnknown() {
		body.SourceUrl = plan.SourceUrl.ValueStringPointer()
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

	// Find the created resource by listing to extract the ID
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
	id := extractOptionListID(listResult, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// GET by ID with SDK mismatch workaround (same as Read)
	readResult, httpResp, err := client.LibraryAPI.GetOptionList(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "option_list", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "option_list",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if !applyGetOptionListResponse(readResult, &plan) {
		resp.Diagnostics.AddError(
			"Not Found After Create",
			"Option type list was created but could not be read by ID. "+
				"The resource may exist in Morpheus. Import manually if needed: "+
				"'terraform import <resource_type>.<name> <id>'",
		)

		return
	}

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

	var state OptionListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetOptionList(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "option_list", "", err, httpResp)

		return
	}

	if !applyGetOptionListResponse(result, &state) {
		resp.State.RemoveResource(ctx)

		return
	}

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

	var plan OptionListModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateOptionListRequestOptionTypeList{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.SourceUrl.IsNull() && !plan.SourceUrl.IsUnknown() {
		body.SourceUrl = plan.SourceUrl.ValueStringPointer()
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

	// GET by ID with SDK mismatch workaround (same as Read)
	readResult, httpResp, err := client.LibraryAPI.GetOptionList(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "option_list", plan.Name.ValueString(), err, httpResp)

		return
	}

	if !applyGetOptionListResponse(readResult, &plan) {
		resp.Diagnostics.AddError(
			"Read Error After Update",
			"Option type list was updated successfully but could not be read back by ID.",
		)

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

	var state OptionListModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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

func mapGetOptionListToModel(model *OptionListModel, ol *sdk.GetOptionList200ResponseOptionTypesInner) {
	if ol.Id != nil {
		model.Id = types.Int64Value(*ol.Id)
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
		model.SourceUrl = types.StringValue(*ol.SourceUrl)
	} else {
		model.SourceUrl = types.StringNull()
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
func mapOptionTypeListFromRaw(model *OptionListModel, m map[string]interface{}) {
	if v, ok := m["id"].(float64); ok {
		model.Id = types.Int64Value(int64(v))
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
		model.SourceUrl = types.StringValue(v)
	} else {
		model.SourceUrl = types.StringNull()
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

// extractOptionListID extracts the option list ID from a list-by-name response,
// handling the SDK field mismatch where the API may return "optionTypeLists" instead of "optionTypes".
// Returns 0 and appends to diags on failure.
func extractOptionListID(result *sdk.ListOptionLists200Response, diags *diag.Diagnostics) int64 {
	if len(result.OptionTypes) > 0 {
		ol := result.OptionTypes[0]
		if ol.Id == nil {
			diags.AddError("API returned nil ID", "OptionList ID is nil in the list response")

			return 0
		}

		return *ol.Id
	}

	// Fallback: SDK field mismatch - API returns "optionTypeLists", SDK expects "optionTypes"
	if rawLists, ok := result.AdditionalProperties["optionTypeLists"]; ok {
		if listsSlice, ok := rawLists.([]interface{}); ok && len(listsSlice) > 0 {
			if firstMap, ok := listsSlice[0].(map[string]interface{}); ok {
				if v, ok := firstMap["id"].(float64); ok {
					return int64(v)
				}
			}
		}
	}

	diags.AddError(
		"Not Found After Create",
		"Option type list was created successfully but could not be found by name. "+
			"The resource may exist in Morpheus. Check the Morpheus UI and import manually if needed: "+
			"'terraform import <resource_type>.<name> <id>'",
	)

	return 0
}

// applyGetOptionListResponse populates model from a GetOptionList response,
// handling the SDK field mismatch where the API returns "optionTypeList" instead of "optionTypes".
// Returns true if the model was populated (caller should set state).
// Returns false if not found (caller handles via AddError or RemoveResource).
func applyGetOptionListResponse(result *sdk.GetOptionList200Response, model *OptionListModel) bool {
	if len(result.OptionTypes) > 0 {
		mapGetOptionListToModel(model, &result.OptionTypes[0])

		return true
	}

	// Fallback: SDK field mismatch - API returns "optionTypeList", SDK expects "optionTypes"
	if rawOL, ok := result.AdditionalProperties["optionTypeList"]; ok {
		if olMap, ok := rawOL.(map[string]interface{}); ok {
			mapOptionTypeListFromRaw(model, olMap)

			return true
		}
	}

	return false
}
