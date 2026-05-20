package libraryfiletemplate

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
	_ resource.Resource                = &libraryFileTemplateResource{}
	_ resource.ResourceWithConfigure   = &libraryFileTemplateResource{}
	_ resource.ResourceWithImportState = &libraryFileTemplateResource{}
)

type libraryFileTemplateResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &libraryFileTemplateResource{}
}

func (r *libraryFileTemplateResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_file_template"
}

func (r *libraryFileTemplateResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LibraryFileTemplateSchema(ctx)
}

func (r *libraryFileTemplateResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan libraryFileTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddFileTemplateRequestContainerTemplate{
		Name:     plan.Name.ValueString(),
		FileName: plan.FileName.ValueString(),
	}

	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.FilePath.IsNull() {
		body.FilePath = plan.FilePath.ValueStringPointer()
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.TemplatePhase.IsNull() {
		body.TemplatePhase = plan.TemplatePhase.ValueStringPointer()
	}
	if !plan.Template.IsNull() {
		body.Template = plan.Template.ValueStringPointer()
	}
	if !plan.SettingName.IsNull() {
		body.SettingName = plan.SettingName.ValueStringPointer()
	}
	if !plan.SettingCategory.IsNull() {
		body.SettingCategory = plan.SettingCategory.ValueStringPointer()
	}

	result, httpResp, err := client.LibraryAPI.AddFileTemplate(ctx).AddFileTemplateRequest(sdk.AddFileTemplateRequest{
		ContainerTemplate: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_file_template", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Extract ID from create response. The SDK's GetId() may return 0 if the API
	// returns the ID nested inside containerTemplate rather than at the top level.
	if id := result.GetId(); id != 0 {
		plan.ID = types.Int64Value(id)
	} else if ctData, ok := result.AdditionalProperties["containerTemplate"]; ok {
		if ctMap, ok := ctData.(map[string]interface{}); ok {
			if id, ok := ctMap["id"].(float64); ok {
				plan.ID = types.Int64Value(int64(id))
			}
		}
	}

	if plan.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Create Failed", "Could not extract ID from API response")

		return
	}

	// Read back the full resource
	if err := readFileTemplateIntoModel(ctx, client, plan.ID.ValueInt64(), &plan, &resp.Diagnostics); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_file_template", plan.Name.ValueString(), err, nil)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryFileTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state libraryFileTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	if err := readFileTemplateIntoModel(ctx, client, id, &state, &resp.Diagnostics); err != nil {
		if err.Error() == "not found" {
			resp.State.RemoveResource(ctx)

			return
		}
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_file_template", "", err, nil)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *libraryFileTemplateResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan libraryFileTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateFileTemplateRequestContainerTemplate{}

	body.Name = plan.Name.ValueStringPointer()
	body.FileName = plan.FileName.ValueStringPointer()

	var labels []string
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.FilePath.IsNull() {
		body.FilePath = plan.FilePath.ValueStringPointer()
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.TemplatePhase.IsNull() {
		body.TemplatePhase = plan.TemplatePhase.ValueStringPointer()
	}
	if !plan.Template.IsNull() {
		body.Template = plan.Template.ValueStringPointer()
	}
	if !plan.SettingName.IsNull() {
		body.SettingName = plan.SettingName.ValueStringPointer()
	}
	if !plan.SettingCategory.IsNull() {
		body.SettingCategory = plan.SettingCategory.ValueStringPointer()
	}

	_, httpResp, err := client.LibraryAPI.UpdateFileTemplate(ctx, id).
		UpdateFileTemplateRequest(sdk.UpdateFileTemplateRequest{
			ContainerTemplate: &body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_file_template", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Read back the full resource
	if err := readFileTemplateIntoModel(ctx, client, id, &plan, &resp.Diagnostics); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_file_template", plan.Name.ValueString(), err, nil)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryFileTemplateResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state libraryFileTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteFileTemplate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_file_template", "", err, httpResp)

		return
	}
}

func (r *libraryFileTemplateResource) ImportState(
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

func mapGetResponseToModel(
	ctx context.Context,
	model *libraryFileTemplateModel,
	tmpl *sdk.GetFileTemplate200ResponseContainerTemplate,
	diags *diag.Diagnostics,
) {
	if tmpl.Id != nil {
		model.ID = types.Int64Value(*tmpl.Id)
	}
	if tmpl.Name != nil {
		model.Name = types.StringValue(*tmpl.Name)
	}
	if tmpl.Labels != nil {
		labelsList, d := types.ListValueFrom(ctx, types.StringType, tmpl.Labels)
		diags.Append(d...)
		model.Labels = labelsList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
	if tmpl.FileName != nil {
		model.FileName = types.StringValue(*tmpl.FileName)
	}
	if tmpl.FilePath != nil {
		model.FilePath = types.StringValue(*tmpl.FilePath)
	} else {
		model.FilePath = types.StringNull()
	}
	if v, ok := tmpl.GetCategoryOk(); ok && v != nil {
		model.Category = types.StringValue(*v)
	} else {
		model.Category = types.StringNull()
	}
	if tmpl.TemplatePhase != nil {
		model.TemplatePhase = types.StringValue(*tmpl.TemplatePhase)
	}
	if tmpl.Template != nil {
		model.Template = types.StringValue(*tmpl.Template)
	} else {
		model.Template = types.StringNull()
	}
	if v, ok := tmpl.GetSettingNameOk(); ok && v != nil {
		model.SettingName = types.StringValue(*v)
	} else {
		model.SettingName = types.StringNull()
	}
	if v, ok := tmpl.GetSettingCategoryOk(); ok && v != nil {
		model.SettingCategory = types.StringValue(*v)
	} else {
		model.SettingCategory = types.StringNull()
	}
}

// readFileTemplateIntoModel calls GetFileTemplate and handles potential SDK
// deserialization issues by falling back to raw JSON parsing.
func readFileTemplateIntoModel(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	model *libraryFileTemplateModel,
	diags *diag.Diagnostics,
) error {
	result, httpResp, err := client.LibraryAPI.GetFileTemplate(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return fmt.Errorf("not found")
	}

	// If SDK deserialization succeeded, use typed response
	if err == nil && result != nil {
		tmpl := result.GetContainerTemplate()
		mapGetResponseToModel(ctx, model, &tmpl, diags)

		return nil
	}

	// SDK deserialization failed. Fall back to raw JSON from the HTTP response body.
	if httpResp == nil || httpResp.StatusCode >= 300 {
		return fmt.Errorf("API error: %w", err)
	}

	body, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		return fmt.Errorf("reading response body: %w", readErr)
	}

	var raw struct {
		ContainerTemplate map[string]interface{} `json:"containerTemplate"`
	}
	if jsonErr := json.Unmarshal(body, &raw); jsonErr != nil {
		return fmt.Errorf("parsing response JSON: %w", jsonErr)
	}

	mapGenericTemplateToModel(ctx, model, raw.ContainerTemplate, diags)

	return nil
}

func mapGenericTemplateToModel(
	ctx context.Context,
	model *libraryFileTemplateModel,
	m map[string]interface{},
	diags *diag.Diagnostics,
) {
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
	if v, ok := m["fileName"].(string); ok {
		model.FileName = types.StringValue(v)
	}
	if v, ok := m["filePath"].(string); ok && v != "" {
		model.FilePath = types.StringValue(v)
	} else {
		model.FilePath = types.StringNull()
	}
	if v, ok := m["category"].(string); ok && v != "" {
		model.Category = types.StringValue(v)
	} else {
		model.Category = types.StringNull()
	}
	if v, ok := m["templatePhase"].(string); ok {
		model.TemplatePhase = types.StringValue(v)
	}
	if v, ok := m["template"].(string); ok {
		model.Template = types.StringValue(v)
	} else {
		model.Template = types.StringNull()
	}
	if v, ok := m["settingName"].(string); ok && v != "" {
		model.SettingName = types.StringValue(v)
	} else {
		model.SettingName = types.StringNull()
	}
	if v, ok := m["settingCategory"].(string); ok && v != "" {
		model.SettingCategory = types.StringValue(v)
	} else {
		model.SettingCategory = types.StringNull()
	}
}
