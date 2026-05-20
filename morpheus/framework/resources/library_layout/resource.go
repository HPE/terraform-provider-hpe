package library_layout

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
	_ resource.Resource                = &libraryLayoutResource{}
	_ resource.ResourceWithConfigure   = &libraryLayoutResource{}
	_ resource.ResourceWithImportState = &libraryLayoutResource{}
)

type libraryLayoutResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &libraryLayoutResource{}
}

func (r *libraryLayoutResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_layout"
}

func (r *libraryLayoutResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LibraryLayoutSchema(ctx)
}

func (r *libraryLayoutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryLayoutModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddLayoutRequestInstanceTypeLayout{
		Name:              plan.Name.ValueString(),
		InstanceVersion:   plan.InstanceVersion.ValueString(),
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
	if !plan.SortOrder.IsNull() && !plan.SortOrder.IsUnknown() {
		body.SortOrder = plan.SortOrder.ValueInt64Pointer()
	}
	if !plan.Creatable.IsNull() && !plan.Creatable.IsUnknown() {
		body.Creatable = plan.Creatable.ValueBoolPointer()
	}
	if !plan.MemoryRequirement.IsNull() && !plan.MemoryRequirement.IsUnknown() {
		body.MemoryRequirement = plan.MemoryRequirement.ValueStringPointer()
	}
	if !plan.HasAutoScale.IsNull() && !plan.HasAutoScale.IsUnknown() {
		body.HasAutoScale = plan.HasAutoScale.ValueBoolPointer()
	}
	if !plan.ContainerTypes.IsNull() && !plan.ContainerTypes.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.ContainerTypes.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.ContainerTypes = ids
	}
	if !plan.OptionTypes.IsNull() && !plan.OptionTypes.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.OptionTypes.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.OptionTypes = ids
	}
	if !plan.SpecTemplates.IsNull() && !plan.SpecTemplates.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.SpecTemplates.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.SpecTemplates = ids
	}

	instanceTypeID := plan.InstanceTypeID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.AddLayout(ctx, instanceTypeID).AddLayoutRequest(sdk.AddLayoutRequest{
		InstanceTypeLayout: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_layout", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Extract ID from response
	if layout := result.GetInstanceTypeLayout(); layout.Id != nil {
		plan.ID = types.Int64Value(*layout.Id)
	} else if ltData, ok := result.AdditionalProperties["instanceTypeLayout"]; ok {
		if ltMap, ok := ltData.(map[string]interface{}); ok {
			if idVal, ok := ltMap["id"]; ok {
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
	readResult, httpResp, err := client.LibraryAPI.GetLayout(ctx, plan.ID.ValueInt64()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_layout", "", err, httpResp)
		return
	}

	lt := readResult.GetInstanceTypeLayout()
	mapGetLayoutToModel(ctx, &plan, &lt, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryLayoutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryLayoutModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetLayout(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_layout", "", err, httpResp)
		return
	}

	lt := result.GetInstanceTypeLayout()
	mapGetLayoutToModel(ctx, &state, &lt, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *libraryLayoutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryLayoutModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateLayoutRequestInstanceTypeLayout{
		Name:              plan.Name.ValueStringPointer(),
		InstanceVersion:   plan.InstanceVersion.ValueStringPointer(),
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
	if !plan.SortOrder.IsNull() && !plan.SortOrder.IsUnknown() {
		body.SortOrder = plan.SortOrder.ValueInt64Pointer()
	}
	if !plan.Creatable.IsNull() && !plan.Creatable.IsUnknown() {
		body.Creatable = plan.Creatable.ValueBoolPointer()
	}
	if !plan.MemoryRequirement.IsNull() && !plan.MemoryRequirement.IsUnknown() {
		body.MemoryRequirement = plan.MemoryRequirement.ValueStringPointer()
	}
	if !plan.HasAutoScale.IsNull() && !plan.HasAutoScale.IsUnknown() {
		body.HasAutoScale = plan.HasAutoScale.ValueBoolPointer()
	}
	if !plan.ContainerTypes.IsNull() && !plan.ContainerTypes.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.ContainerTypes.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.ContainerTypes = ids
	}
	if !plan.OptionTypes.IsNull() && !plan.OptionTypes.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.OptionTypes.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.OptionTypes = ids
	}
	if !plan.SpecTemplates.IsNull() && !plan.SpecTemplates.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.SpecTemplates.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.SpecTemplates = ids
	}

	_, httpResp, err := client.LibraryAPI.UpdateLayout(ctx, id).UpdateLayoutRequest(sdk.UpdateLayoutRequest{
		InstanceTypeLayout: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_layout", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back
	readResult, httpResp, err := client.LibraryAPI.GetLayout(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_layout", "", err, httpResp)
		return
	}

	lt := readResult.GetInstanceTypeLayout()
	mapGetLayoutToModel(ctx, &plan, &lt, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryLayoutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryLayoutModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteLayout(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_layout", "", err, httpResp)
		return
	}
}

func (r *libraryLayoutResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetLayoutToModel(ctx context.Context, model *libraryLayoutModel, lt *sdk.GetLayout200ResponseInstanceTypeLayout, diags *diag.Diagnostics) {
	if lt.Id != nil {
		model.ID = types.Int64Value(*lt.Id)
	}
	if lt.Name != nil {
		model.Name = types.StringValue(*lt.Name)
	}
	if lt.InstanceVersion != nil {
		model.InstanceVersion = types.StringValue(*lt.InstanceVersion)
	}
	if v := lt.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if lt.SortOrder != nil {
		model.SortOrder = types.Int64Value(*lt.SortOrder)
	}
	if lt.Creatable != nil {
		model.Creatable = types.BoolValue(*lt.Creatable)
	}
	if lt.ProvisionType != nil && lt.ProvisionType.Code != nil {
		model.ProvisionTypeCode = types.StringValue(*lt.ProvisionType.Code)
	}
	if v := lt.MemoryRequirement.Get(); v != nil {
		model.MemoryRequirement = types.StringValue(fmt.Sprintf("%d", *v))
	} else {
		model.MemoryRequirement = types.StringNull()
	}

	// Labels
	if lt.Labels != nil {
		labelValues, d := types.ListValueFrom(ctx, types.StringType, lt.Labels)
		diags.Append(d...)
		model.Labels = labelValues
	} else {
		model.Labels = types.ListNull(types.StringType)
	}

	// ContainerTypes - extract IDs from objects
	if lt.ContainerTypes != nil {
		var ids []int64
		for _, ct := range lt.ContainerTypes {
			if ct.Id != nil {
				ids = append(ids, int64(*ct.Id))
			}
		}
		idValues, d := types.ListValueFrom(ctx, types.Int64Type, ids)
		diags.Append(d...)
		model.ContainerTypes = idValues
	} else {
		model.ContainerTypes = types.ListNull(types.Int64Type)
	}

	// OptionTypes - extract IDs from objects
	if lt.OptionTypes != nil {
		var ids []int64
		for _, ot := range lt.OptionTypes {
			if idVal, ok := ot["id"]; ok {
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
		model.OptionTypes = idValues
	} else {
		model.OptionTypes = types.ListNull(types.Int64Type)
	}

	// SpecTemplates - extract IDs from objects
	if lt.SpecTemplates != nil {
		var ids []int64
		for _, st := range lt.SpecTemplates {
			if idVal, ok := st["id"]; ok {
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
		model.SpecTemplates = idValues
	} else {
		model.SpecTemplates = types.ListNull(types.Int64Type)
	}

	// HasAutoScale - not in get response, preserve state
	// instance_type_id is preserved from state/plan (not returned in layout GET)
}
