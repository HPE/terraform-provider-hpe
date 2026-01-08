package instance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	errfmt "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HPE/terraform-provider-hpe/internal/framework/utils"
)

var (
	// All known statuses are:
	// pending, denied, cancelled, provisioning, finishing, failed, resizing,
	// running, warning, stopped, suspended, removing, restarting, cloning,
	// restoring, stopping, starting, suspending, pendingRemoval,
	// pendingDeleteApproval, pendingReconfigureApproval, unknown
	CreateTargetStatuses = []string{
		"running",
	}

	CreateErrorStatuses = []string{
		"denied",
		"cancelled",
		"failed",
		"stopped",
		"suspended",
		"removing",
		"pendingRemoval",
	}
)

// Create implements resource.Resource.
func (g *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan InstanceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	createTimeout, diags := plan.Timeouts.Create(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	// cloud_id
	reqInstance := sdk.NewAddInstanceRequestWithDefaults()
	if !plan.CloudId.IsNull() {
		reqInstance.SetZoneId(plan.CloudId.ValueInt64())
	}

	// config
	configMap := make(map[string]any)
	if !plan.Config.IsNull() {
		configValue := plan.Config.UnderlyingValue()
		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create instance resource",
				"instance: failed to convert config: "+
					err.Error(),
			)

			return
		}
		configDataMap, ok := configAny.(map[string]any)
		if ok {
			configMap = configDataMap
		} else {
			resp.Diagnostics.AddError(
				"error creating instance",
				"could not parse config value",
			)
		}
	}
	reqInstance.Config = sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigConfig{
		MapmapOfStringAny: &configMap,
	}

	// evars
	evars, diags := convert.FromSetType(ctx, plan.Evars, evarMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert evars")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetEvars(evars)

	// group_id
	reqInstance.Instance.SetSite(
		*sdk.NewAddInstanceRequestInstanceSite(plan.GroupId.ValueInt64()),
	)

	// instance_context
	if !plan.InstanceContext.IsNull() {
		reqInstance.Instance.SetInstanceContext(plan.InstanceContext.ValueString())
	}

	// instance_type_id
	if !plan.InstanceTypeId.IsNull() {
		// instance creation API does not support specifying an instance type
		// ID directly instead it expects an instance type code. For consistency
		// we prefer to use ID in Terraform. This means we need to make an extra
		// API call to get the instance type code value.
		code, diags := getInstanceTypeCode(ctx, client, plan.InstanceTypeId.ValueInt64())
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}

		reqInstance.Instance.SetInstanceType(
			*sdk.NewAddInstanceRequestInstanceInstanceType(
				code,
			),
		)
	}

	// layout_id
	if !plan.LayoutId.IsNull() {
		reqInstance.Instance.SetLayout(
			*sdk.NewAddInstanceRequestInstanceLayout(
				plan.LayoutId.ValueInt64(),
			),
		)
	}

	// name
	if !plan.Name.IsNull() {
		reqInstance.Instance.SetName(plan.Name.ValueString())
	}

	// network_interfaces
	networkInterfaces, diags := convert.FromListType(
		ctx,
		plan.NetworkInterfaces,
		networkInterfaceMapper(ctx),
	)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetNetworkInterfaces(networkInterfaces)

	// plan_id
	if !plan.PlanId.IsNull() {
		reqInstance.Instance.SetPlan(
			sdk.AddInstanceRequestInstancePlan{Id: plan.PlanId.ValueInt64()},
		)
	}

	// ports
	ports, diags := convert.FromSetType(
		ctx,
		plan.Ports,
		func(in PortsValue) sdk.AddInstanceRequestPortsInner {
			return sdk.AddInstanceRequestPortsInner{
				Port: in.Port.ValueInt64(),
				Name: in.Name.ValueStringPointer(),
				Lb:   *sdk.NewNullableString(in.LoadBalancerProtocol.ValueStringPointer()),
			}
		},
	)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert ports")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetPorts(ports)

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, tagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetTags(tags)

	// task_set_id
	if !plan.TaskSetId.IsNull() {
		reqInstance.SetTaskSetId(plan.TaskSetId.ValueInt64())
	}

	// volumes
	volumes, diags := convert.FromListType(ctx, plan.Volumes, volumeMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetVolumes(volumes)

	instance, httpResp, err := client.InstancesAPI.AddInstance(ctx).
		AddInstanceRequest(*reqInstance).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating instance", errfmt.ErrMsg(err, httpResp))

		return
	}

	// Store ID locally but not in state yet
	instanceId := instance.Instance.GetId()

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		utils.TaintResourceState(ctx, utils.TaintResourceStateConfig{
			ResourceType: "instance",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Wait for the instance to be ready
	waitForReady := func() (string, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceId).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			return "", backoff.Permanent(err)
		}

		status := resp.Instance.GetStatus()

		return status, checkStatusDone(
			status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if status, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	); err != nil {
		if status == "" {
			errUnwrapped := errors.Unwrap(err)
			if errUnwrapped != nil {
				resp.Diagnostics.AddError(
					"create instance resource",
					fmt.Sprintf(
						"instance %d: provisioning failed: %v",
						instanceId,
						errUnwrapped,
					),
				)
			} else {
				resp.Diagnostics.AddError(
					"create instance resource",
					fmt.Sprintf(
						"instance %d: provisioning failed",
						instanceId,
					),
				)
			}
		} else {
			resp.Diagnostics.AddError(
				"create instance resource",
				fmt.Sprintf(
					"instance %d: provisioning failed current status is: %s",
					instanceId,
					status,
				),
			)
		}
		taintResourceState(instanceId)

		return
	}

	state, diag := getInstanceAsState(ctx, instanceId, client, plan)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		resp.Diagnostics.AddError(
			"failed to read instance state",
			fmt.Sprintf("Instance %d was created but could not be read", instanceId),
		)
		taintResourceState(instanceId)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set instance state",
			fmt.Sprintf("Instance %d was created but state could not be saved", instanceId),
		)
		taintResourceState(instanceId)

		return
	}
}

func getInstanceTypeCode(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	_, resp, _ := client.InstancesAPI.
		GetInstanceTypeProvisioning(ctx, id).
		Execute()

	if resp == nil || resp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate instance resource",
			// TODO: better error message
			fmt.Sprintf("instance type %d GET failed", id),
		)

		return "", diags
	}

	instanceType := struct {
		InstanceType struct {
			Code string `json:"code"`
		} `json:"instanceType"`
	}{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&instanceType); err != nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance type %d decode failed: %v", id, err),
		)

		return "", diags
	}

	return instanceType.InstanceType.Code, diags
}
