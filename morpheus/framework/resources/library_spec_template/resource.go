package library_spec_template

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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &librarySpecTemplateResource{}
	_ resource.ResourceWithConfigure   = &librarySpecTemplateResource{}
	_ resource.ResourceWithImportState = &librarySpecTemplateResource{}
)

type librarySpecTemplateResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &librarySpecTemplateResource{}
}

func (r *librarySpecTemplateResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_spec_template"
}

func (r *librarySpecTemplateResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LibrarySpecTemplateSchema(ctx)
}

func (r *librarySpecTemplateResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan librarySpecTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceType := "local"
	if !plan.Source.IsNull() {
		sourceType = plan.Source.ValueString()
	}

	file := sdk.AddSpecTemplateRequestSpecTemplateFile{
		SourceType: sourceType,
	}
	if !plan.Content.IsNull() {
		file.Content = plan.Content.ValueStringPointer()
	}

	body := sdk.AddSpecTemplateRequestSpecTemplate{
		Name: plan.Name.ValueString(),
		Type: sdk.AddSpecTemplateRequestSpecTemplateType{
			Code: plan.Type.ValueString(),
		},
		File: file,
	}

	result, httpResp, err := client.LibraryAPI.AddSpecTemplate(ctx).AddSpecTemplateRequest(sdk.AddSpecTemplateRequest{
		SpecTemplate: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_spec_template", plan.Name.ValueString(), err, httpResp)

		return
	}

	if id := result.Id.Get(); id != nil {
		plan.ID = types.Int64Value(*id)
	} else if stData, ok := result.AdditionalProperties["specTemplate"]; ok {
		if stMap, ok := stData.(map[string]interface{}); ok {
			if idVal, ok := stMap["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					plan.ID = types.Int64Value(int64(v))
				case int64:
					plan.ID = types.Int64Value(v)
				}
			}
		}
	}

	// Read back to populate all fields
	readResult, httpResp, err := client.LibraryAPI.GetSpecTemplate(ctx, plan.ID.ValueInt64()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_spec_template", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "library_spec_template",
			ResourceID:   plan.ID.ValueInt64(),
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	st := readResult.GetSpecTemplate()
	mapGetSpecTemplateToModel(&plan, &st)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *librarySpecTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state librarySpecTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetSpecTemplate(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_spec_template", "", err, httpResp)

		return
	}

	st := result.GetSpecTemplate()
	mapGetSpecTemplateToModel(&state, &st)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *librarySpecTemplateResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan librarySpecTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	sourceType := "local"
	if !plan.Source.IsNull() {
		sourceType = plan.Source.ValueString()
	}

	updateFile := sdk.UpdateSpecTemplateRequestSpecTemplateFile{
		SourceType: &sourceType,
	}
	if !plan.Content.IsNull() {
		updateFile.Content = plan.Content.ValueStringPointer()
	}

	typeCode := plan.Type.ValueString()
	updateBody := sdk.UpdateSpecTemplateRequestSpecTemplate{
		Name: plan.Name.ValueStringPointer(),
		Type: &sdk.UpdateSpecTemplateRequestSpecTemplateType{
			Code: &typeCode,
		},
		File: &updateFile,
	}

	_, httpResp, err := client.LibraryAPI.UpdateSpecTemplate(ctx, id).
		UpdateSpecTemplateRequest(sdk.UpdateSpecTemplateRequest{
			SpecTemplate: &updateBody,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_spec_template", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Read back
	result, httpResp, err := client.LibraryAPI.GetSpecTemplate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_spec_template", "", err, httpResp)

		return
	}

	st := result.GetSpecTemplate()
	mapGetSpecTemplateToModel(&plan, &st)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *librarySpecTemplateResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state librarySpecTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteSpecTemplate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_spec_template", "", err, httpResp)

		return
	}
}

func (r *librarySpecTemplateResource) ImportState(
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

func mapGetSpecTemplateToModel(model *librarySpecTemplateModel, st *sdk.GetSpecTemplate200ResponseSpecTemplate) {
	if st.Id != nil {
		model.ID = types.Int64Value(*st.Id)
	}
	if st.Name != nil {
		model.Name = types.StringValue(*st.Name)
	}
	if st.Type != nil && st.Type.Code != nil {
		model.Type = types.StringValue(*st.Type.Code)
	}
	if v := st.ExternalId.Get(); v != nil {
		model.ExternalID = types.StringValue(*v)
	} else {
		model.ExternalID = types.StringNull()
	}
	if st.File != nil && st.File.SourceType != nil {
		model.Source = types.StringValue(*st.File.SourceType)
	} else {
		model.Source = types.StringNull()
	}
	if st.File != nil && st.File.Content != nil {
		model.Content = types.StringValue(*st.File.Content)
	} else {
		model.Content = types.StringNull()
	}
}
