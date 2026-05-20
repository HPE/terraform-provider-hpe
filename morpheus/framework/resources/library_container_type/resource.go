package library_container_type

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
	_ resource.Resource                = &libraryContainerTypeResource{}
	_ resource.ResourceWithConfigure   = &libraryContainerTypeResource{}
	_ resource.ResourceWithImportState = &libraryContainerTypeResource{}
)

type libraryContainerTypeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &libraryContainerTypeResource{}
}

func (r *libraryContainerTypeResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_container_type"
}

func (r *libraryContainerTypeResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LibraryContainerTypeSchema(ctx)
}

func (r *libraryContainerTypeResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan libraryContainerTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddNodeTypeRequestContainerType{
		Name:              plan.Name.ValueString(),
		ShortName:         plan.ShortName.ValueString(),
		ContainerVersion:  plan.ContainerVersion.ValueString(),
		ProvisionTypeCode: plan.ProvisionTypeCode.ValueString(),
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Scripts.IsNull() && !plan.Scripts.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.Scripts.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Scripts = ids
	}
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.Templates.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Templates = ids
	}
	if !plan.VirtualImageID.IsNull() && !plan.VirtualImageID.IsUnknown() {
		body.VirtualImageId = plan.VirtualImageID.ValueInt64Pointer()
	}
	if !plan.StatTypeCode.IsNull() && !plan.StatTypeCode.IsUnknown() {
		body.StatTypeCode = plan.StatTypeCode.ValueStringPointer()
	}
	if !plan.LogTypeCode.IsNull() && !plan.LogTypeCode.IsUnknown() {
		body.LogTypeCode = plan.LogTypeCode.ValueStringPointer()
	}
	if !plan.ServerType.IsNull() && !plan.ServerType.IsUnknown() {
		body.ServerType = plan.ServerType.ValueStringPointer()
	}

	result, httpResp, err := client.LibraryAPI.AddNodeType(ctx).AddNodeTypeRequest(sdk.AddNodeTypeRequest{
		ContainerType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_container_type", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Extract ID from response
	if ct := result.GetContainerType(); ct.Id != nil {
		plan.ID = types.Int64Value(int64(*ct.Id))
	} else if ctData, ok := result.AdditionalProperties["containerType"]; ok {
		if ctMap, ok := ctData.(map[string]interface{}); ok {
			if idVal, ok := ctMap["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					plan.ID = types.Int64Value(int64(v))
				case int64:
					plan.ID = types.Int64Value(v)
				}
			}
		}
	}

	// Read back
	readResult, httpResp, err := client.LibraryAPI.GetNodeType(ctx, plan.ID.ValueInt64()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_type", "", err, httpResp)

		return
	}

	ct := readResult.GetContainerType()
	mapGetNodeTypeToModel(ctx, &plan, &ct, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryContainerTypeResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state libraryContainerTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetNodeType(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_type", "", err, httpResp)

		return
	}

	ct := result.GetContainerType()
	mapGetNodeTypeToModel(ctx, &state, &ct, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *libraryContainerTypeResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan libraryContainerTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateNodeTypeRequestContainerType{
		Name:              plan.Name.ValueStringPointer(),
		ShortName:         plan.ShortName.ValueStringPointer(),
		ContainerVersion:  plan.ContainerVersion.ValueStringPointer(),
		ProvisionTypeCode: plan.ProvisionTypeCode.ValueStringPointer(),
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Scripts.IsNull() && !plan.Scripts.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.Scripts.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Scripts = ids
	}
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.Templates.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Templates = ids
	}
	if !plan.VirtualImageID.IsNull() && !plan.VirtualImageID.IsUnknown() {
		body.VirtualImageId = plan.VirtualImageID.ValueInt64Pointer()
	}
	if !plan.StatTypeCode.IsNull() && !plan.StatTypeCode.IsUnknown() {
		body.StatTypeCode = plan.StatTypeCode.ValueStringPointer()
	}
	if !plan.LogTypeCode.IsNull() && !plan.LogTypeCode.IsUnknown() {
		body.LogTypeCode = plan.LogTypeCode.ValueStringPointer()
	}
	if !plan.ServerType.IsNull() && !plan.ServerType.IsUnknown() {
		body.ServerType = plan.ServerType.ValueStringPointer()
	}

	_, httpResp, err := client.LibraryAPI.UpdateNodeType(ctx, id).UpdateNodeTypeRequest(sdk.UpdateNodeTypeRequest{
		ContainerType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_container_type", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Read back
	readResult, httpResp, err := client.LibraryAPI.GetNodeType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_container_type", "", err, httpResp)

		return
	}

	ct := readResult.GetContainerType()
	mapGetNodeTypeToModel(ctx, &plan, &ct, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryContainerTypeResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state libraryContainerTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteNodeType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_container_type", "", err, httpResp)

		return
	}
}

func (r *libraryContainerTypeResource) ImportState(
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

func mapGetNodeTypeToModel(
	ctx context.Context,
	model *libraryContainerTypeModel,
	ct *sdk.GetNodeType200ResponseContainerType,
	diags *diag.Diagnostics,
) {
	if ct.Id != nil {
		model.ID = types.Int64Value(int64(*ct.Id))
	}
	if ct.Name != nil {
		model.Name = types.StringValue(*ct.Name)
	}
	if ct.ShortName != nil {
		model.ShortName = types.StringValue(*ct.ShortName)
	}
	if ct.ContainerVersion != nil {
		model.ContainerVersion = types.StringValue(*ct.ContainerVersion)
	}
	if ct.ProvisionType != nil && ct.ProvisionType.Code != nil {
		model.ProvisionTypeCode = types.StringValue(*ct.ProvisionType.Code)
	}

	// Description - API does not return this field; preserve existing state value
	if desc, ok := ct.AdditionalProperties["description"]; ok && desc != nil {
		if s, ok := desc.(string); ok {
			model.Description = types.StringValue(s)
		}
		// else: keep existing model.Description unchanged
	}
	// else: keep existing model.Description unchanged (API doesn't return it)

	// Labels
	if ct.Labels != nil {
		labelValues, d := types.ListValueFrom(ctx, types.StringType, ct.Labels)
		diags.Append(d...)
		model.Labels = labelValues
	} else {
		model.Labels = types.ListNull(types.StringType)
	}

	// Scripts - extract IDs from containerScripts objects
	if ct.ContainerScripts != nil {
		var ids []int64
		for _, s := range ct.ContainerScripts {
			if idVal, ok := s["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					ids = append(ids, int64(v))
				case int64:
					ids = append(ids, v)
				}
			}
		}
		idValues, d := types.ListValueFrom(ctx, types.Int64Type, ids)
		diags.Append(d...)
		model.Scripts = idValues
	} else {
		model.Scripts = types.ListNull(types.Int64Type)
	}

	// Templates - extract IDs from containerTemplates objects
	if ct.ContainerTemplates != nil {
		var ids []int64
		for _, t := range ct.ContainerTemplates {
			if idVal, ok := t["id"]; ok {
				switch v := idVal.(type) {
				case float64:
					ids = append(ids, int64(v))
				case int64:
					ids = append(ids, v)
				}
			}
		}
		idValues, d := types.ListValueFrom(ctx, types.Int64Type, ids)
		diags.Append(d...)
		model.Templates = idValues
	} else {
		model.Templates = types.ListNull(types.Int64Type)
	}

	// VirtualImage
	if ct.VirtualImage != nil && ct.VirtualImage.Id != nil {
		model.VirtualImageID = types.Int64Value(int64(*ct.VirtualImage.Id))
	} else {
		model.VirtualImageID = types.Int64Null()
	}

	// stat/log/server type from AdditionalProperties
	if v, ok := ct.AdditionalProperties["statTypeCode"]; ok && v != nil {
		if s, ok := v.(string); ok {
			model.StatTypeCode = types.StringValue(s)
		} else {
			model.StatTypeCode = types.StringNull()
		}
	} else {
		model.StatTypeCode = types.StringNull()
	}
	if v, ok := ct.AdditionalProperties["logTypeCode"]; ok && v != nil {
		if s, ok := v.(string); ok {
			model.LogTypeCode = types.StringValue(s)
		} else {
			model.LogTypeCode = types.StringNull()
		}
	} else {
		model.LogTypeCode = types.StringNull()
	}
	if v, ok := ct.AdditionalProperties["serverType"]; ok && v != nil {
		if s, ok := v.(string); ok {
			model.ServerType = types.StringValue(s)
		} else {
			model.ServerType = types.StringNull()
		}
	} else {
		model.ServerType = types.StringNull()
	}
}
