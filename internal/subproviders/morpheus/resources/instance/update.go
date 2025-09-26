package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	errfmt "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Update implements resource.Resource.
func (g *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {

	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	var plan InstanceModel
	var state InstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Instance update state: %v", state.Volumes.Elements()))
	tflog.Info(ctx, fmt.Sprintf("Instance update plan: %v", plan.Volumes.Elements()))

	updateInstance := client.InstancesAPI.UpdateInstance(ctx, plan.Id.ValueInt64())
	updateRequest := sdk.NewUpdateInstanceRequest()
	instanceUpdateRequest := sdk.NewUpdateInstanceRequestInstance()

	// name
	if !plan.Name.IsNull() {
		instanceUpdateRequest.Name = plan.Name.ValueStringPointer()
		instanceUpdateRequest.DisplayName = plan.Name.ValueStringPointer()
	}

	// TODO: DESCRIPTION IS MISSING FROM THE SCHEMA??
	// if !plan.Description.IsNull() {
	// 	instanceUpdateRequest.Description = plan.Description.ValueStringPointer()
	// }

	// instance_context
	if !plan.InstanceContext.IsNull() {
		instanceUpdateRequest.InstanceContext = plan.InstanceContext.ValueStringPointer()
	}

	// group_id
	if !plan.GroupId.IsNull() {
		site := sdk.NewUpdateInstanceRequestInstanceSite()
		site.Id = plan.GroupId.ValueInt64Pointer()
		instanceUpdateRequest.Site = site
	}

	// tags
	if !plan.Tags.IsNull() {
		tags, diags := convert.FromSetType(ctx, plan.Tags, tagMapper)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert tags")
			resp.Diagnostics.Append(diags...)

			return
		}
		instanceUpdateRequest.SetTags(tags)
	}

	updateRequest.SetInstance(*instanceUpdateRequest)
	updateResp, httpResp, err := updateInstance.UpdateInstanceRequest(*updateRequest).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating instance", errfmt.ErrMsg(err, httpResp))

		return
	}

	_ = updateResp

	needResize := false

	if !state.PlanId.Equal(plan.PlanId) {
		needResize = true
	}

	// compare state and plan volumes so we only resize if required

	// compare state and plan network_interfaces so we only resize if required

	if needResize {
		resp.Diagnostics.AddWarning("instance requires a resize", "attributes changed in the plan which will require resizing the instance")
	}

	resizeInstance := client.InstancesAPI.ResizeInstance(ctx, plan.Id.ValueInt64())
	resizeRequest := sdk.NewResizeInstanceRequest()

	// plan_id
	if !plan.PlanId.IsNull() || !plan.PlanId.IsUnknown() {
		resizeRequest.Instance = sdk.NewResizeInstanceRequestInstance()
		resizeRequest.Instance.Plan = sdk.NewResizeInstanceRequestInstancePlan()
		resizeRequest.Instance.Plan.SetId(plan.PlanId.ValueInt64())
	}

	// volumes
	volumes, diags := convert.FromListType(ctx, plan.Volumes, volumeMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}

	resizeRequest.SetDeleteOriginalVolumes(true)
	resizeRequest.SetVolumes(volumes)

	// network_interfaces
	networkInterfaces, diags := convert.FromSetType(
		ctx,
		plan.NetworkInterfaces,
		networkInterfaceMapper,
	)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")
		resp.Diagnostics.Append(diags...)

		return
	}
	resizeRequest.SetNetworkInterfaces(networkInterfaces)

	resizeResp, httpResp, err := resizeInstance.ResizeInstanceRequest(*resizeRequest).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error resizing instance", errfmt.ErrMsg(err, httpResp))

		return
	}

	tflog.Info(ctx, fmt.Sprintln(resizeResp))

	waitForReady := func() (*sdk.GetInstance200Response, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, plan.Id.ValueInt64()).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			return nil, backoff.Permanent(err)
		}

		status := resp.Instance.GetStatus()

		return resp, checkStatusDone(
			status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if r, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		var status string

		if r.GetInstance().Status != nil {
			status = *r.GetInstance().Status
		}

		resp.Diagnostics.AddError(
			"resize instance resource",
			fmt.Sprintf(
				"instance %d: resizing failed current status is: %v",
				plan.Id.ValueInt64(),
				status,
			),
		)
	}

	state2, diag := getInstanceAsState(ctx, state.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state2)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Println(resp.State)
}
