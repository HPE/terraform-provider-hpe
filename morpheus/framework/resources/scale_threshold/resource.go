package scale_threshold

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &scaleThresholdResource{}
	_ resource.ResourceWithConfigure   = &scaleThresholdResource{}
	_ resource.ResourceWithImportState = &scaleThresholdResource{}
)

type scaleThresholdResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &scaleThresholdResource{}
}

func (r *scaleThresholdResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_scale_threshold"
}

func (r *scaleThresholdResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ScaleThresholdSchema(ctx)
}

func (r *scaleThresholdResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan scaleThresholdModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddScaleThresholdsRequest{
		ScaleThreshold: sdk.AddScaleThresholdsRequestScaleThreshold{
			Name: plan.Name.ValueString(),
		},
	}

	if !plan.AutoUpscale.IsNull() {
		body.ScaleThreshold.AutoUp = plan.AutoUpscale.ValueBoolPointer()
	}
	if !plan.AutoDownscale.IsNull() {
		body.ScaleThreshold.AutoDown = plan.AutoDownscale.ValueBoolPointer()
	}
	if !plan.MinCount.IsNull() {
		v := int32(plan.MinCount.ValueInt64()) //nolint:gosec // value range is safe
		body.ScaleThreshold.MinCount = &v
	}
	if !plan.MaxCount.IsNull() {
		v := int32(plan.MaxCount.ValueInt64()) //nolint:gosec // value range is safe
		body.ScaleThreshold.MaxCount = &v
	}
	if !plan.CPUEnabled.IsNull() {
		body.ScaleThreshold.CpuEnabled = plan.CPUEnabled.ValueBoolPointer()
	}
	if !plan.MinCPU.IsNull() {
		v := plan.MinCPU.ValueFloat64()
		body.ScaleThreshold.MinCpu = &v
	}
	if !plan.MaxCPU.IsNull() {
		v := plan.MaxCPU.ValueFloat64()
		body.ScaleThreshold.MaxCpu = &v
	}
	if !plan.MemoryEnabled.IsNull() {
		body.ScaleThreshold.MemoryEnabled = plan.MemoryEnabled.ValueBoolPointer()
	}
	if !plan.MinMemory.IsNull() {
		v := plan.MinMemory.ValueFloat64()
		body.ScaleThreshold.MinMemory = &v
	}
	if !plan.MaxMemory.IsNull() {
		v := plan.MaxMemory.ValueFloat64()
		body.ScaleThreshold.MaxMemory = &v
	}

	result, httpResp, err := client.AutomationAPI.AddScaleThresholds(ctx).AddScaleThresholdsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "scale_threshold", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.ScaleThreshold != nil && result.ScaleThreshold.Id != nil {
		plan.ID = types.Int64Value(*result.ScaleThreshold.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scaleThresholdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state scaleThresholdModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.AutomationAPI.GetScaleThresholds(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "scale_threshold", "", err, httpResp)

		return
	}

	st := result.ScaleThreshold
	if st != nil {
		state.ID = types.Int64Value(*st.Id)
		if st.Name != nil {
			state.Name = types.StringValue(*st.Name)
		}
		if st.AutoUp != nil {
			state.AutoUpscale = types.BoolValue(*st.AutoUp)
		}
		if st.AutoDown != nil {
			state.AutoDownscale = types.BoolValue(*st.AutoDown)
		}
		if st.MinCount != nil {
			state.MinCount = types.Int64Value(*st.MinCount)
		}
		if st.MaxCount != nil {
			state.MaxCount = types.Int64Value(*st.MaxCount)
		}
		if st.CpuEnabled != nil {
			state.CPUEnabled = types.BoolValue(*st.CpuEnabled)
		}
		if st.MinCpu != nil {
			state.MinCPU = types.Float64Value(float64(*st.MinCpu))
		}
		if st.MaxCpu != nil {
			state.MaxCPU = types.Float64Value(float64(*st.MaxCpu))
		}
		if st.MemoryEnabled != nil {
			state.MemoryEnabled = types.BoolValue(*st.MemoryEnabled)
		}
		if st.MinMemory != nil {
			state.MinMemory = types.Float64Value(float64(*st.MinMemory))
		}
		if st.MaxMemory != nil {
			state.MaxMemory = types.Float64Value(float64(*st.MaxMemory))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scaleThresholdResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan scaleThresholdModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	name := plan.Name.ValueString()
	body := sdk.UpdateScaleThresholdsRequest{
		ScaleThreshold: sdk.UpdateScaleThresholdsRequestScaleThreshold{
			Name: &name,
		},
	}

	if !plan.AutoUpscale.IsNull() {
		body.ScaleThreshold.AutoUp = plan.AutoUpscale.ValueBoolPointer()
	}
	if !plan.AutoDownscale.IsNull() {
		body.ScaleThreshold.AutoDown = plan.AutoDownscale.ValueBoolPointer()
	}
	if !plan.MinCount.IsNull() {
		v := int32(plan.MinCount.ValueInt64()) //nolint:gosec // value range is safe
		body.ScaleThreshold.MinCount = &v
	}
	if !plan.MaxCount.IsNull() {
		v := int32(plan.MaxCount.ValueInt64()) //nolint:gosec // value range is safe
		body.ScaleThreshold.MaxCount = &v
	}
	if !plan.CPUEnabled.IsNull() {
		body.ScaleThreshold.CpuEnabled = plan.CPUEnabled.ValueBoolPointer()
	}
	if !plan.MinCPU.IsNull() {
		v := float32(plan.MinCPU.ValueFloat64())
		body.ScaleThreshold.MinCpu = &v
	}
	if !plan.MaxCPU.IsNull() {
		v := float32(plan.MaxCPU.ValueFloat64())
		body.ScaleThreshold.MaxCpu = &v
	}
	if !plan.MemoryEnabled.IsNull() {
		body.ScaleThreshold.MemoryEnabled = plan.MemoryEnabled.ValueBoolPointer()
	}
	if !plan.MinMemory.IsNull() {
		v := float32(plan.MinMemory.ValueFloat64())
		body.ScaleThreshold.MinMemory = &v
	}
	if !plan.MaxMemory.IsNull() {
		v := float32(plan.MaxMemory.ValueFloat64())
		body.ScaleThreshold.MaxMemory = &v
	}

	_, httpResp, err := client.AutomationAPI.UpdateScaleThresholds(ctx, id).UpdateScaleThresholdsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "scale_threshold", plan.Name.ValueString(), err, httpResp)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scaleThresholdResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state scaleThresholdModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.AutomationAPI.RemoveScaleThresholds(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "scale_threshold", "", err, httpResp)

		return
	}
}

func (r *scaleThresholdResource) ImportState(
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

// Ensure unused imports are satisfied.
var _ *http.Response
