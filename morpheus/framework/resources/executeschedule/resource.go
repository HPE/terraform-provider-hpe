package executeschedule

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
	_ resource.Resource                = &executeScheduleResource{}
	_ resource.ResourceWithConfigure   = &executeScheduleResource{}
	_ resource.ResourceWithImportState = &executeScheduleResource{}
)

type executeScheduleResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &executeScheduleResource{}
}

func (r *executeScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_execute_schedule"
}

func (r *executeScheduleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ExecuteScheduleSchema(ctx)
}

func (r *executeScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan executeScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Build request body from plan and call API
	// result, httpResp, err := client.JobsAPI.AddExecuteSchedule(ctx).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "execute_schedule", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = polling.ForCreate

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *executeScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state executeScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to read resource
	// result, httpResp, err := client.JobsAPI.GetExecuteSchedule(ctx, id).Execute()
	var httpResp *http.Response
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "execute_schedule", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	// TODO: Map response to state model
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *executeScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan executeScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	// TODO: Build request body from plan and call API
	// _, httpResp, err := client.JobsAPI.UpdateExecuteSchedule(ctx, id).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "execute_schedule", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *executeScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state executeScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to delete resource
	// httpResp, err := client.JobsAPI.DeleteExecuteSchedule(ctx, id).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "execute_schedule", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id
	_ = fmt.Sprintf
	_ = polling.ForDelete
}

func (r *executeScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
