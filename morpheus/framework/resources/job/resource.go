package job

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
	_ resource.Resource                = &jobResource{}
	_ resource.ResourceWithConfigure   = &jobResource{}
	_ resource.ResourceWithImportState = &jobResource{}
)

type jobResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &jobResource{}
}

func (r *jobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_job"
}

func (r *jobResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = JobSchema(ctx)
}

func (r *jobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan jobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := sdk.WorkflowJobPayload{
		Name:       plan.Name.ValueString(),
		TargetType: plan.TargetType.ValueString(),
	}
	if !plan.ScheduleMode.IsNull() {
		schedMode := plan.ScheduleMode.ValueString()
		payload.ScheduleMode = sdk.StringAsWorkflowJobPayloadScheduleMode(&schedMode)
	} else {
		manual := "manual"
		payload.ScheduleMode = sdk.StringAsWorkflowJobPayloadScheduleMode(&manual)
	}
	if !plan.Enabled.IsNull() {
		payload.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.WorkflowID.IsNull() {
		wfID := plan.WorkflowID.ValueInt64()
		payload.Workflow = sdk.WorkflowJobPayloadWorkflow{Id: &wfID}
	}
	if !plan.CustomConfig.IsNull() {
		payload.CustomConfig = plan.CustomConfig.ValueStringPointer()
	}

	body := sdk.WorkflowJobPayloadAsAddJobsRequestJob(&payload)

	result, httpResp, err := client.JobsAPI.AddJobs(ctx).AddJobsRequest(sdk.AddJobsRequest{
		Job: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "job", plan.Name.ValueString(), err, httpResp)
		return
	}

	job := result.GetJob()
	mapAddResponseToModel(&plan, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state jobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.JobsAPI.GetJobs(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "job", "", err, httpResp)
		return
	}

	job := result.GetJob()
	mapGetResponseToModel(&state, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan jobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateJobsRequestJob{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.TargetType.IsNull() {
		body.TargetType = plan.TargetType.ValueStringPointer()
	}
	if !plan.CustomConfig.IsNull() {
		body.CustomConfig = plan.CustomConfig.ValueStringPointer()
	}

	_, httpResp, err := client.JobsAPI.UpdateJobs(ctx, id).UpdateJobsRequest(sdk.UpdateJobsRequest{
		Job: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "job", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read
	getResult, httpResp, err := client.JobsAPI.GetJobs(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "job", "", err, httpResp)
		return
	}

	job := getResult.GetJob()
	mapGetResponseToModel(&plan, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state jobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.JobsAPI.RemoveJobs(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "job", "", err, httpResp)
		return
	}
}

func (r *jobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *jobModel, job *sdk.AddJobs200ResponseAllOfJob) {
	if job.Id != nil {
		model.ID = types.Int64Value(*job.Id)
	}
	if job.Name != nil {
		model.Name = types.StringValue(*job.Name)
	}
	if job.Enabled != nil {
		model.Enabled = types.BoolValue(*job.Enabled)
	}
	if job.TargetType.IsSet() && job.TargetType.Get() != nil {
		model.TargetType = types.StringValue(*job.TargetType.Get())
	} else {
		model.TargetType = types.StringNull()
	}
	if job.CustomConfig.IsSet() && job.CustomConfig.Get() != nil {
		model.CustomConfig = types.StringValue(*job.CustomConfig.Get())
	} else {
		model.CustomConfig = types.StringNull()
	}
}

func mapGetResponseToModel(model *jobModel, job *sdk.GetJobs200ResponseAllOfJob) {
	// Use the AnyOf variant which has the common fields
	if job.GetJobs200ResponseAllOfJobAnyOf != nil {
		j := job.GetJobs200ResponseAllOfJobAnyOf
		if j.Id != nil {
			model.ID = types.Int64Value(*j.Id)
		}
		if j.Name != nil {
			model.Name = types.StringValue(*j.Name)
		}
		if j.Enabled != nil {
			model.Enabled = types.BoolValue(*j.Enabled)
		}
		if j.TargetType.IsSet() && j.TargetType.Get() != nil {
			model.TargetType = types.StringValue(*j.TargetType.Get())
		} else {
			model.TargetType = types.StringNull()
		}
		if j.CustomConfig.IsSet() && j.CustomConfig.Get() != nil {
			model.CustomConfig = types.StringValue(*j.CustomConfig.Get())
		} else {
			model.CustomConfig = types.StringNull()
		}
		if j.ScheduleMode != nil {
			sm := j.ScheduleMode
			if sm.String != nil {
				model.ScheduleMode = types.StringValue(*sm.String)
			} else if sm.Int64 != nil {
				model.ScheduleMode = types.StringValue(strconv.FormatInt(*sm.Int64, 10))
			} else {
				model.ScheduleMode = types.StringNull()
			}
		} else {
			model.ScheduleMode = types.StringNull()
		}
		if j.Workflow != nil && j.Workflow.Id != nil {
			model.WorkflowID = types.Int64Value(*j.Workflow.Id)
		}
	}
}
