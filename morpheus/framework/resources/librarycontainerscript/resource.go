package librarycontainerscript

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	_ resource.Resource                = &libraryContainerScriptResource{}
	_ resource.ResourceWithConfigure   = &libraryContainerScriptResource{}
	_ resource.ResourceWithImportState = &libraryContainerScriptResource{}
)

type libraryContainerScriptResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &libraryContainerScriptResource{}
}

func (r *libraryContainerScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_container_script"
}

func (r *libraryContainerScriptResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LibraryContainerScriptSchema(ctx)
}

func (r *libraryContainerScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryContainerScriptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddScriptRequestContainerScript{
		Name: plan.Name.ValueString(),
	}

	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)

		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.ScriptVersion.IsNull() {
		body.ScriptVersion = plan.ScriptVersion.ValueStringPointer()
	}
	if !plan.ScriptPhase.IsNull() {
		body.ScriptPhase = plan.ScriptPhase.ValueStringPointer()
	}
	if !plan.ScriptType.IsNull() {
		body.ScriptType = plan.ScriptType.ValueStringPointer()
	}
	if !plan.Script.IsNull() {
		body.Script = plan.Script.ValueStringPointer()
	}
	if !plan.RunAsUser.IsNull() {
		body.RunAsUser = plan.RunAsUser.ValueStringPointer()
	}
	if !plan.SudoUser.IsNull() {
		body.SudoUser = plan.SudoUser.ValueBoolPointer()
	}

	result, httpResp, err := client.LibraryAPI.AddScript(ctx).AddScriptRequest(sdk.AddScriptRequest{
		ContainerScript: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_container_script", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Extract ID from response — SDK may not deserialize correctly due to account field mismatch
	var scriptID int64
	if result != nil {
		script := result.GetContainerScript()
		if script.Id != nil {
			scriptID = *script.Id
		}
	}
	if scriptID == 0 {
		// Try AdditionalProperties
		if result != nil {
			if csData, ok := result.AdditionalProperties["containerScript"]; ok {
				if csMap, ok := csData.(map[string]interface{}); ok {
					if id, ok := csMap["id"].(float64); ok {
						scriptID = int64(id)
					}
				}
			}
		}
	}
	if scriptID == 0 {
		resp.Diagnostics.AddError("Error creating library_container_script", "Could not extract ID from response")
		return
	}

	// Read back the full resource
	if err := readScriptIntoModel(ctx, client, scriptID, &plan, &resp.Diagnostics); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_script", plan.Name.ValueString(), err, nil)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryContainerScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryContainerScriptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	if err := readScriptIntoModel(ctx, client, id, &state, &resp.Diagnostics); err != nil {
		if err.Error() == "not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_script", "", err, nil)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *libraryContainerScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryContainerScriptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateScriptRequestContainerScript{
		Name: plan.Name.ValueStringPointer(),
	}

	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.ScriptVersion.IsNull() {
		body.ScriptVersion = plan.ScriptVersion.ValueStringPointer()
	}
	if !plan.ScriptPhase.IsNull() {
		body.ScriptPhase = plan.ScriptPhase.ValueStringPointer()
	}
	if !plan.ScriptType.IsNull() {
		body.ScriptType = plan.ScriptType.ValueStringPointer()
	}
	if !plan.Script.IsNull() {
		body.Script = plan.Script.ValueStringPointer()
	}
	if !plan.RunAsUser.IsNull() {
		body.RunAsUser = plan.RunAsUser.ValueStringPointer()
	}
	if !plan.SudoUser.IsNull() {
		body.SudoUser = plan.SudoUser.ValueBoolPointer()
	}

	_, httpResp, err := client.LibraryAPI.UpdateScript(ctx, id).UpdateScriptRequest(sdk.UpdateScriptRequest{
		ContainerScript: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_container_script", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back
	if err := readScriptIntoModel(ctx, client, id, &plan, &resp.Diagnostics); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_script", plan.Name.ValueString(), err, nil)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryContainerScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryContainerScriptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteScript(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_container_script", "", err, httpResp)
		return
	}
}

func (r *libraryContainerScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetResponseToModel(ctx context.Context, model *libraryContainerScriptModel, script *sdk.GetScript200ResponseContainerScript, diags *diag.Diagnostics) {
	if script.Id != nil {
		model.ID = types.Int64Value(*script.Id)
	}
	if script.Name != nil {
		model.Name = types.StringValue(*script.Name)
	}
	if script.Labels != nil {
		labelsList, d := types.ListValueFrom(ctx, types.StringType, script.Labels)
		diags.Append(d...)
		model.Labels = labelsList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
	if v, ok := script.GetCategoryOk(); ok && v != nil {
		model.Category = types.StringValue(*v)
	} else {
		model.Category = types.StringNull()
	}
	if script.ScriptVersion != nil {
		model.ScriptVersion = types.StringValue(*script.ScriptVersion)
	}
	if script.ScriptPhase != nil {
		model.ScriptPhase = types.StringValue(*script.ScriptPhase)
	}
	if script.ScriptType != nil {
		model.ScriptType = types.StringValue(*script.ScriptType)
	}
	if script.Script != nil {
		model.Script = types.StringValue(*script.Script)
	} else {
		model.Script = types.StringNull()
	}
	if v, ok := script.GetRunAsUserOk(); ok && v != nil {
		model.RunAsUser = types.StringValue(*v)
	} else {
		model.RunAsUser = types.StringNull()
	}
	if script.SudoUser != nil {
		model.SudoUser = types.BoolValue(*script.SudoUser)
	}
	if script.FailOnError != nil {
		model.FailOnError = types.BoolValue(*script.FailOnError)
	}
}

// readScriptIntoModel calls GetScript and handles the SDK deserialization error
// caused by the 'account' field being an object in the API response but typed as
// a string in the SDK. When the SDK fails, it falls back to raw JSON parsing.
func readScriptIntoModel(ctx context.Context, client *sdk.APIClient, id int64, model *libraryContainerScriptModel, diags *diag.Diagnostics) error {
	result, httpResp, err := client.LibraryAPI.GetScript(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return fmt.Errorf("not found")
	}

	// If SDK deserialization succeeded, use typed response
	if err == nil && result != nil {
		script := result.GetContainerScript()
		mapGetResponseToModel(ctx, model, &script, diags)
		return nil
	}

	// SDK deserialization failed (likely due to 'account' field type mismatch).
	// Fall back to raw JSON from the HTTP response body.
	if httpResp == nil || httpResp.StatusCode >= 300 {
		return fmt.Errorf("API error: %w", err)
	}

	body, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		return fmt.Errorf("reading response body: %w", readErr)
	}

	var raw struct {
		ContainerScript map[string]interface{} `json:"containerScript"`
	}
	if jsonErr := json.Unmarshal(body, &raw); jsonErr != nil {
		return fmt.Errorf("parsing response JSON: %w", jsonErr)
	}

	mapGenericScriptToModel(ctx, model, raw.ContainerScript, diags)
	return nil
}

func mapGenericScriptToModel(ctx context.Context, model *libraryContainerScriptModel, m map[string]interface{}, diags *diag.Diagnostics) {
	if id, ok := m["id"].(float64); ok {
		model.ID = types.Int64Value(int64(id))
	}
	if name, ok := m["name"].(string); ok {
		model.Name = types.StringValue(name)
	}
	if labels, ok := m["labels"].([]interface{}); ok {
		strs := make([]string, 0, len(labels))
		for _, l := range labels {
			if s, ok := l.(string); ok {
				strs = append(strs, s)
			}
		}
		labelsList, d := types.ListValueFrom(ctx, types.StringType, strs)
		diags.Append(d...)
		model.Labels = labelsList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
	if v, ok := m["category"].(string); ok && v != "" {
		model.Category = types.StringValue(v)
	} else {
		model.Category = types.StringNull()
	}
	if v, ok := m["scriptVersion"].(string); ok {
		model.ScriptVersion = types.StringValue(v)
	}
	if v, ok := m["scriptPhase"].(string); ok {
		model.ScriptPhase = types.StringValue(v)
	}
	if v, ok := m["scriptType"].(string); ok {
		model.ScriptType = types.StringValue(v)
	}
	if v, ok := m["script"].(string); ok {
		model.Script = types.StringValue(v)
	} else {
		model.Script = types.StringNull()
	}
	if v, ok := m["runAsUser"].(string); ok && v != "" {
		model.RunAsUser = types.StringValue(v)
	} else {
		model.RunAsUser = types.StringNull()
	}
	if v, ok := m["sudoUser"].(bool); ok {
		model.SudoUser = types.BoolValue(v)
	}
	if v, ok := m["failOnError"].(bool); ok {
		model.FailOnError = types.BoolValue(v)
	}
}
