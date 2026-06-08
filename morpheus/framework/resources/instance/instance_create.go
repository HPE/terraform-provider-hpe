// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

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

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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
	reqInstance := &sdk.AddInstanceRequest{}
	if !plan.CloudId.IsNull() {
		reqInstance.ZoneId = plan.CloudId.ValueInt64Pointer()
	}

	// config
	switch {
	// AWS config
	case !plan.ConfigAws.IsNull() && !plan.ConfigAws.IsUnknown():
		noAgent := plan.ConfigAws.NoAgent.ValueBool()
		isEC2 := convert.BoolToStringTrueFalse(plan.ConfigAws.IsEc2.ValueBool()).ValueString()
		kmsKeyId := plan.ConfigAws.KmsKeyId.ValueString()
		instanceProfile := plan.ConfigAws.InstanceProfile.ValueString()
		publicIpType := plan.ConfigAws.PublicIpType.ValueString()
		availabilityId := plan.ConfigAws.AvailabilityZoneId.ValueString()
		resourcePoolId := plan.ConfigAws.ResourcePoolId.ValueString()

		configAWS := &sdk.AmazonInstanceConfiguration2{
			NoAgent:         *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:  &resourcePoolId,
			IsEC2:           &isEC2,
			KmsKeyId:        &kmsKeyId,
			InstanceProfile: &instanceProfile,
			PublicIpType:    &publicIpType,
			AvailabilityId:  &availabilityId,
		}

		// Security Groups
		if !plan.ConfigAws.SecurityGroups.IsNull() && !plan.ConfigAws.SecurityGroups.IsUnknown() {
			securityGroups, diags := convert.FromListType(
				ctx,
				plan.ConfigAws.SecurityGroups,
				func(in SecurityGroupsValue) sdk.AddInstanceRequestSecurityGroupsInner {
					id := in.Id.ValueString()

					return sdk.AddInstanceRequestSecurityGroupsInner{
						Id: &id,
					}
				},
			)
			if diags.HasError() {
				tflog.Error(ctx, "cannot convert AWS security groups")
				resp.Diagnostics.Append(diags...)

				return
			}
			reqInstance.SecurityGroups = securityGroups
		}

		reqInstance.Config = sdk.AddInstanceRequestConfig{
			AmazonInstanceConfiguration2: configAWS,
		}

	// HVM config
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		// The provisionTypeCode default is "mvm" which is the code for the HVM provisioning type.
		createUser := plan.ConfigHvm.CreateUser.ValueBool()
		nestedVirtualization := plan.ConfigHvm.NestedVirtualization.ValueString()
		noAgent := plan.ConfigHvm.NoAgent.ValueBool()
		resourcePoolId := plan.ConfigHvm.ResourcePoolId.ValueString()

		configHvm := &sdk.HVMInstanceConfiguration{
			CreateUser:           *sdk.NewNullableBool(&createUser),
			NestedVirtualization: &nestedVirtualization,
			NoAgent:              *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:       &resourcePoolId,
		}

		if !plan.ConfigHvm.KvmHostId.IsNull() {
			configHvm.KvmHostId = plan.ConfigHvm.KvmHostId.ValueInt64Pointer()
		}

		reqInstance.Config = sdk.AddInstanceRequestConfig{
			HVMInstanceConfiguration: configHvm,
		}

	// VMware config
	case !plan.ConfigVmware.IsNull() && !plan.ConfigVmware.IsUnknown():
		nestedVirtualization := plan.ConfigVmware.NestedVirtualization.ValueString()
		createUser := plan.ConfigVmware.CreateUser.ValueBool()
		noAgent := plan.ConfigVmware.NoAgent.ValueBool()
		resourcePoolId := plan.ConfigVmware.ResourcePoolId.ValueString()
		vmwareFolderId := plan.ConfigVmware.VmwareFolderId.ValueString()

		configVMware := &sdk.VMWareInstanceConfiguration2{
			NestedVirtualization: &nestedVirtualization,
			CreateUser:           *sdk.NewNullableBool(&createUser),
			NoAgent:              *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:       &resourcePoolId,
			VmwareFolderId:       &vmwareFolderId,
		}

		reqInstance.Config = sdk.AddInstanceRequestConfig{
			VMWareInstanceConfiguration2: configVMware,
		}

	// Azure config
	case !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown():
		createUser := plan.ConfigAzure.CreateUser.ValueBool()
		resourcePoolId := plan.ConfigAzure.ResourcePoolId.ValueString()

		configAzure := &sdk.AzureInstanceConfiguration2{
			CreateUser:     &createUser,
			ResourcePoolId: &resourcePoolId,
		}

		if !plan.ConfigAzure.AzureRegion.IsNull() && !plan.ConfigAzure.AzureRegion.IsUnknown() {
			configAzure.AzureRegion = plan.ConfigAzure.AzureRegion.ValueStringPointer()
		}

		if !plan.ConfigAzure.AzuresecurityGroupId.IsNull() && !plan.ConfigAzure.AzuresecurityGroupId.IsUnknown() {
			configAzure.AzuresecurityGroupId = plan.ConfigAzure.AzuresecurityGroupId.ValueStringPointer()
		}

		if !plan.ConfigAzure.AvailabilityOptions.IsNull() && !plan.ConfigAzure.AvailabilityOptions.IsUnknown() {
			configAzure.AvailabilityOptions = plan.ConfigAzure.AvailabilityOptions.ValueStringPointer()
		}

		if !plan.ConfigAzure.AvailabilitySet.IsNull() && !plan.ConfigAzure.AvailabilitySet.IsUnknown() {
			configAzure.AvailabilitySet = plan.ConfigAzure.AvailabilitySet.ValueStringPointer()
		}

		if !plan.ConfigAzure.AvailabilityZone.IsNull() && !plan.ConfigAzure.AvailabilityZone.IsUnknown() {
			if configAzure.AdditionalProperties == nil {
				configAzure.AdditionalProperties = make(map[string]interface{})
			}
			configAzure.AdditionalProperties["availabilityZone"] = plan.ConfigAzure.AvailabilityZone.ValueString()
		}

		if !plan.ConfigAzure.AzurefloatingIp.IsNull() && !plan.ConfigAzure.AzurefloatingIp.IsUnknown() {
			configAzure.AzurefloatingIp = plan.ConfigAzure.AzurefloatingIp.ValueStringPointer()
		}

		if !plan.ConfigAzure.BootDiagnostics.IsNull() && !plan.ConfigAzure.BootDiagnostics.IsUnknown() {
			configAzure.BootDiagnostics = plan.ConfigAzure.BootDiagnostics.ValueStringPointer()
		}

		if !plan.ConfigAzure.OsGuestDiagnostics.IsNull() && !plan.ConfigAzure.OsGuestDiagnostics.IsUnknown() {
			configAzure.OsGuestDiagnostics = plan.ConfigAzure.OsGuestDiagnostics.ValueStringPointer()
		}

		if !plan.ConfigAzure.DiagnosticsStorageAccount.IsNull() && !plan.ConfigAzure.DiagnosticsStorageAccount.IsUnknown() {
			configAzure.DiagnosticsStorageAccount = plan.ConfigAzure.DiagnosticsStorageAccount.ValueStringPointer()
		}

		reqInstance.Config = sdk.AddInstanceRequestConfig{
			AzureInstanceConfiguration2: configAzure,
		}

	// Generic config
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
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
		// use interface{} to satisfy SDK AdditionalProperties
		configMap := make(map[string]interface{})
		configDataMap, ok := configAny.(map[string]interface{})
		if ok {
			configMap = configDataMap
		} else {
			resp.Diagnostics.AddError(
				"error creating instance",
				"could not parse config value",
			)
		}

		reqInstance.Config = sdk.AddInstanceRequestConfig{
			GenericInstanceConfiguration2: &sdk.GenericInstanceConfiguration2{
				AdditionalProperties: configMap,
			},
		}
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqInstance.Instance.Description = plan.Description.ValueStringPointer()
	}

	// evars
	evars, diags := convert.FromSetType(ctx, plan.Evars, evarMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert evars")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.Evars = evars

	// group_id
	reqInstance.Instance.Site = sdk.AddInstanceRequestInstanceSite{
		Id: plan.GroupId.ValueInt64(),
	}

	// instance_context
	if !plan.InstanceContext.IsNull() {
		reqInstance.Instance.InstanceContext = plan.InstanceContext.ValueStringPointer()
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

		reqInstance.Instance.InstanceType = sdk.AddInstanceRequestInstanceInstanceType{
			Code: code,
		}
	}

	// layout_id
	if !plan.LayoutId.IsNull() {
		reqInstance.Instance.Layout = sdk.AddInstanceRequestInstanceLayout{
			Id: plan.LayoutId.ValueInt64(),
		}
	}

	// layout_size
	if !plan.LayoutSize.IsNull() {
		reqInstance.LayoutSize = plan.LayoutSize.ValueInt64Pointer()
	}

	// name
	if !plan.Name.IsNull() {
		reqInstance.Instance.Name = plan.Name.ValueString()
	}

	// network_domain_id
	if !plan.NetworkDomainId.IsNull() {
		netDomain := &sdk.AddInstanceRequestInstanceNetworkDomain{
			Id: plan.NetworkDomainId.ValueInt64(),
		}
		reqInstance.Instance.NetworkDomain = netDomain
	}

	// network_interfaces
	networkInterfaces, diags := convert.FromListType(
		ctx,
		plan.NetworkInterfaces,
		createNetworkInterfaceMapper(ctx),
	)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.NetworkInterfaces = networkInterfaces

	// plan_id
	if !plan.PlanId.IsNull() {
		reqInstance.Instance.Plan = sdk.AddInstanceRequestInstancePlan{
			Id: plan.PlanId.ValueInt64(),
		}
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
	reqInstance.Ports = ports

	// service_plan_options
	if !plan.ServicePlanOptions.IsNull() {
		memory := *plan.ServicePlanOptions.MaxMemory.ValueInt64Pointer() << 20
		servicePlanOptions := &sdk.AddInstanceRequestServicePlanOptions{
			MaxMemory:      &memory,
			MaxCores:       plan.ServicePlanOptions.MaxCores.ValueInt64Pointer(),
			CoresPerSocket: plan.ServicePlanOptions.CoresPerSocket.ValueInt64Pointer(),
		}
		reqInstance.ServicePlanOptions = servicePlanOptions
	}

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, createTagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.Tags = tags

	// task_set_id
	if !plan.TaskSetId.IsNull() {
		reqInstance.TaskSetId = plan.TaskSetId.ValueInt64Pointer()
	}

	// volumes
	volumes, diags := convert.FromListType(ctx, plan.Volumes, createVolumeMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.Volumes = volumes

	instance, httpResp, err := client.InstancesAPI.AddInstance(ctx).
		AddInstanceRequest(*reqInstance).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating instance", errfmt.ErrMsg(err, httpResp))

		return
	}

	if instance.Instance.Id == nil {
		resp.Diagnostics.AddError("error creating instance", "POST returned empty instance ID")

		return
	}

	// Store ID locally but not in state yet
	plan.Id = convert.Int64ToType(instance.Instance.Id)
	instanceId := *instance.Instance.Id

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "instance",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Wait for the instance to be ready
	waitForReady := func() (string, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceId).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}
		}

		inst := resp.Instance
		if inst == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty instance", instanceId))
		}

		if inst.Status == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", instanceId))
		}

		return *inst.Status, checkStatusDone(
			*inst.Status,
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
						"instance %d: provisioning failed - unknown error",
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

	// Read the instance state
	state, d := getInstanceAsState(ctx, instanceId, client, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		resp.Diagnostics.AddError(
			"failed to read instance state",
			fmt.Sprintf("Instance %d was created but could not be read", instanceId),
		)
		taintResourceState(instanceId)

		return
	}

	// Set ServicePlanOptions in state.
	state.ServicePlanOptions = plan.ServicePlanOptions

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
