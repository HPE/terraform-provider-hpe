package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	errfmt "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
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
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state InstanceModel
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
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		instanceUpdateRequest.Name = plan.Name.ValueStringPointer()
		instanceUpdateRequest.DisplayName = plan.Name.ValueStringPointer()
	}

	// TODO: DESCRIPTION IS MISSING FROM THE SCHEMA??
	// if !plan.Description.IsNull() {
	// 	instanceUpdateRequest.Description = plan.Description.ValueStringPointer()
	// }

	// instance_context
	if !plan.InstanceContext.IsNull() && !plan.InstanceContext.IsUnknown() {
		instanceUpdateRequest.InstanceContext = plan.InstanceContext.ValueStringPointer()
	}

	// group_id
	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		site := sdk.NewUpdateInstanceRequestInstanceSite()
		site.Id = plan.GroupId.ValueInt64Pointer()
		instanceUpdateRequest.Site = site
	}

	// tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, diags := convert.FromSetType(ctx, plan.Tags, tagMapper)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert tags")
			resp.Diagnostics.Append(diags...)

			return
		}
		instanceUpdateRequest.SetTags(tags)
	}

	updateRequest.SetInstance(*instanceUpdateRequest)
	_, httpResp, err := updateInstance.UpdateInstanceRequest(*updateRequest).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating instance", errfmt.ErrMsg(err, httpResp))

		return
	}

	resizeInstance := client.InstancesAPI.ResizeInstance(ctx, plan.Id.ValueInt64())
	resizeRequest := sdk.NewResizeInstanceRequest()

	resizing := !state.PlanId.Equal(plan.PlanId)

	// compare state and plan volumes so we only resize if required
	// TODO: make this compare each volume rather than just length
	if len(state.Volumes.Elements()) != len(plan.Volumes.Elements()) {
		resizing = true

		volumes, diags := convert.FromListType(ctx, plan.Volumes, volumeMapper)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert volumes")
			resp.Diagnostics.Append(diags...)

			return
		}

		resizeRequest.SetDeleteOriginalVolumes(false)
		resizeRequest.SetVolumes(volumes)
	}

	// compare state and plan network_interfaces so we only resize if required
	// TODO: make this compare each network interface rather than just length
	if len(state.NetworkInterfaces.Elements()) != len(plan.NetworkInterfaces.Elements()) {
		resizing = true

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
	}

	if resizing {
		// plan_id
		if !plan.PlanId.IsNull() || !plan.PlanId.IsUnknown() {
			resizeRequest.Instance = sdk.NewResizeInstanceRequestInstance()
			resizeRequest.Instance.Plan = sdk.NewResizeInstanceRequestInstancePlan()
			resizeRequest.Instance.Plan.SetId(plan.PlanId.ValueInt64())
		}

		resizeResp, httpResp, err := resizeInstance.ResizeInstanceRequest(*resizeRequest).Execute()
		if err != nil || httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError("error resizing instance", errfmt.ErrMsg(err, httpResp))

			return
		}

		tflog.Info(ctx, fmt.Sprintln(resizeResp))

		waitForReady := func() (string, error) {
			resp, hresp, err := client.InstancesAPI.GetInstance(ctx, plan.Id.ValueInt64()).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}

			status := resp.Instance.GetStatus()

			return status, checkStatusDone(
				status,
				UpdateTargetStatuses,
				UpdateErrorStatuses,
			)
		}

		if status, err := backoff.Retry(
			ctx,
			waitForReady,
			backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
			backoff.WithMaxElapsedTime(45*time.Minute),
		); err != nil {
			resp.Diagnostics.AddError(
				"resize instance resource",
				fmt.Sprintf(
					"instance %d: resizing failed current status is: %s",
					plan.Id.ValueInt64(),
					status,
				),
			)
		}
	}

	newState, diag := getInstanceAsState(ctx, state.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
