package powerschedule

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/polling"
)

var (
	_ resource.Resource                = &powerScheduleResource{}
	_ resource.ResourceWithConfigure   = &powerScheduleResource{}
	_ resource.ResourceWithImportState = &powerScheduleResource{}
)

type powerScheduleResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &powerScheduleResource{}
}

func (r *powerScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_power_schedule"
}

func (r *powerScheduleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PowerScheduleSchema(ctx)
}

func (r *powerScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan powerScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Build request body from plan and call API
	// result, httpResp, err := client.AutomationAPI.AddPowerSchedule(ctx).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "power_schedule", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = polling.ForCreate

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *powerScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state powerScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to read resource
	// result, httpResp, err := client.AutomationAPI.GetPowerSchedule(ctx, id).Execute()
	var httpResp *http.Response
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "power_schedule", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	// TODO: Map response to state model
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *powerScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan powerScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	// TODO: Build request body from plan and call API
	// _, httpResp, err := client.AutomationAPI.UpdatePowerSchedule(ctx, id).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "power_schedule", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *powerScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state powerScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to delete resource
	// httpResp, err := client.AutomationAPI.DeletePowerSchedule(ctx, id).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "power_schedule", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id
	_ = fmt.Sprintf
	_ = polling.ForDelete
}

func (r *powerScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
