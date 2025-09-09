package instance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	errfmt "github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
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

	DeleteErrorStatuses = []string{
		"denied",
		"cancelled",
		"failed",
		"stopped",
		"suspended",
		"restoring",
	}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

// Metadata implements resource.Resource.
func (g *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance"
	resp.TypeName = strings.Join(
		[]string{req.ProviderTypeName, constants.SubProviderName, "instance"},
		"_",
	)
}

// Schema implements resource.Resource.
func (g *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = InstanceResourceSchema(ctx)
}

// ImportState implements resource.ResourceWithImportState.
func (g *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import instance resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(diags...)
}

func getInstanceEnvVars(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) ([]sdk.GetEnvVariables200ResponseEnvsInner, diag.Diagnostics) {
	var diags diag.Diagnostics
	resp, _, _ := client.InstancesAPI.GetEnvVariables(ctx, id).Execute()
	// ignoring errors for now, the sdk can't parse some of the unused fields
	// due to polymorphic values

	// if err != nil || hresp.StatusCode != http.StatusOK {
	// diags.AddError(
	// 	"populate instance resource",
	// 	fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
	// )

	// return nil, diags
	// }

	return resp.GetEnvs(), diags
}

func getInstanceAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan InstanceModel,
) (InstanceModel, diag.Diagnostics) {
	var state InstanceModel
	var diags diag.Diagnostics

	resp, hresp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	instance := resp.GetInstance()

	// cloud_id
	state.CloudId = convert.Int64ToType(instance.Cloud.Id)

	// config
	state.Config = types.DynamicNull()

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	// evars
	// API may respond with more evars than what the user set so we need to
	// check the /instance/{id}/envs endpoint which gives us the user specified
	// evars separately.
	envVars, diags := getInstanceEnvVars(ctx, id, client)
	if diags.HasError() {
		return state, diags
	}

	evars, d := convert.ToSetType(
		ctx,
		envVars,
		func(
			in sdk.GetEnvVariables200ResponseEnvsInner,
		) EvarsValue {
			return EvarsValue{
				Name:  types.StringValue(in.Name),
				Value: types.StringValue(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	state.Evars = evars

	// group_id
	state.GroupId = convert.Int64ToType(instance.Group.Id)

	// id
	state.Id = convert.Int64ToType(instance.Id)

	// instance_context
	state.InstanceContext = convert.StrToType(instance.InstanceContext.Get())

	// instance_type_id
	state.InstanceTypeId = convert.Int64ToType(instance.InstanceType.Id)

	// instance_type_id
	state.InstanceTypeId = convert.Int64ToType(instance.InstanceType.Id)

	// layout_id
	state.LayoutId = convert.Int64ToType(instance.Layout.Id)

	// layout_size
	state.LayoutSize = convert.Int64ToType(instance.Config.LayoutSize)

	// name
	state.Name = convert.StrToType(instance.Name)

	// network_interfaces
	interfaces, d := convert.ToSetType(
		ctx,
		resp.GetInstance().Interfaces,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
		) NetworkInterfacesValue {
			v := NetworkInterfacesValue{}
			v.IpAddress = convert.StrToType(in.IpAddress)
			v.IpMode = convert.StrToType(in.IpMode)

			groupID := int64(in.Network.GetGroup())
			v.NetworkGroupId = types.Int64Value(groupID)

			v.NetworkId = convert.Int64ToType(in.Network.Id)

			v.state = attr.ValueStateKnown

			return v
		},
	)
	diags.Append(d...)
	state.NetworkInterfaces = interfaces

	// plan_id
	state.PlanId = convert.Int64ToType(instance.Plan.Id)

	// ports
	// assume ports always match the plan. Ports are a bit complicated in the
	// API. They are a part of container_details, which is an array which may
	// be of arbitrary size. The ports are contained within every element of the
	// array. The array of ports also may contain non-user controlled ports.
	// This should probably be replaced by a plan modifier.
	state.Ports = plan.Ports

	// tags
	tags, d := convert.ToSetType(
		ctx,
		resp.GetInstance().Tags,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner,
		) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	state.Tags = tags

	// task_set_id
	// task_set_id is not included in the API response, it is a write only value
	// we might need modify the schema to add some plan modifiers to it
	state.TaskSetId = plan.TaskSetId

	// volumes
	apiVolumes := slices.DeleteFunc(
		resp.GetInstance().Volumes,
		func(v sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner) bool {
			if v.Name == nil {
				return false
			}

			if strings.HasPrefix(*v.Name, "CD ROM") {
				return true
			}

			return false
		},
	)

	volumes, d := convert.ToSetType(
		ctx,
		apiVolumes,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
		) VolumesValue {
			v := VolumesValue{}
			v.Id = convert.Int64ToType(in.Id)
			v.RootVolume = convert.BoolToType(in.RootVolume)
			v.Name = convert.StrToType(in.Name)
			v.Size = convert.Int64ToType(in.Size)
			v.StorageTypeId = convert.Int64ToType(in.StorageType)

			if in.DatastoreId != nil {
				datastore, err := strconv.ParseInt(in.GetDatastoreId(), 10, 64)
				if err != nil {
					v.DatastoreId = basetypes.NewInt64Unknown()
				}

				v.DatastoreId = types.Int64Value(datastore)
			}

			v.ControllerMountPoint = convert.StrToType(in.ControllerMountPoint)
			v.state = attr.ValueStateKnown

			return v
		},
	)
	diags.Append(d...)
	state.Volumes = volumes

	return state, diags
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

func checkStatusDone(status string, targetStatuses []string, errorStatuses []string) error {
	switch {
	case slices.Contains(errorStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
	}
}

// Create implements resource.Resource.
func (g *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data InstanceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	// cloud_id
	reqInstance := sdk.NewAddInstanceRequestWithDefaults()
	if !data.CloudId.IsNull() {
		reqInstance.SetZoneId(data.CloudId.ValueInt64())
	}

	// config
	configMap := make(map[string]any)
	if !data.Config.IsNull() {
		configValue := data.Config.UnderlyingValue()
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
	evars, diags := convert.FromSetType(ctx, data.Evars, evarMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert evars")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetEvars(evars)

	// group_id
	reqInstance.Instance.SetSite(
		*sdk.NewAddInstanceRequestInstanceSite(data.GroupId.ValueInt64()),
	)

	// instance_context
	if !data.InstanceContext.IsNull() {
		reqInstance.Instance.SetInstanceContext(data.InstanceContext.ValueString())
	}

	// instance_type_id
	if !data.InstanceTypeId.IsNull() {
		// instance creation API does not support specifying an instance type
		// ID directly instead it expects an instance type code. For consistency
		// we prefer to use ID in Terraform. This means we need to make an extra
		// API call to get the instance type code value.
		code, diags := getInstanceTypeCode(ctx, client, data.InstanceTypeId.ValueInt64())
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
	if !data.LayoutId.IsNull() {
		reqInstance.Instance.SetLayout(
			*sdk.NewAddInstanceRequestInstanceLayout(
				data.LayoutId.ValueInt64(),
			),
		)
	}

	// layout_size
	if !data.LayoutSize.IsNull() {
		reqInstance.SetLayoutSize(data.LayoutSize.ValueInt64())
	}

	// name
	if !data.Name.IsNull() {
		reqInstance.Instance.SetName(data.Name.ValueString())
	}

	// network_interfaces
	networkInterfaces, diags := convert.FromSetType(
		ctx,
		data.NetworkInterfaces,
		networkInterfaceMapper,
	)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetNetworkInterfaces(networkInterfaces)

	// plan_id
	if !data.PlanId.IsNull() {
		reqInstance.Instance.SetPlan(
			sdk.AddInstanceRequestInstancePlan{Id: data.PlanId.ValueInt64()},
		)
	}

	// ports
	ports, diags := convert.FromSetType(
		ctx,
		data.Ports,
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
	tags, diags := convert.FromSetType(ctx, data.Tags, tagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqInstance.SetTags(tags)

	// task_set_id
	if !data.TaskSetId.IsNull() {
		reqInstance.SetTaskSetId(data.TaskSetId.ValueInt64())
	}

	// volumes
	volumes, diags := convert.FromSetType(ctx, data.Volumes, volumeMapper)
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

	data.Id = convert.Int64ToType(instance.Instance.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	waitForReady := func() (*sdk.GetInstance200Response, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, data.Id.ValueInt64()).Execute()
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
			"create instance resource",
			fmt.Sprintf(
				"instance %d: provisioning failed current status is: %v",
				data.Id.ValueInt64(),
				status,
			),
		)
	}

	state, diag := getInstanceAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (g *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data InstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id
	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	deleteReq := client.InstancesAPI.DeleteInstance(ctx, id.ValueInt64())
	_, hresp, err := deleteReq.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	waitForDeleted := func() (*sdk.GetInstance200Response, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, data.Id.ValueInt64()).Execute()
		if err != nil && hresp.StatusCode != http.StatusNotFound {
			return nil, backoff.Permanent(err)
		}

		// 404 status code counts as a successful delete
		if hresp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		status := resp.Instance.GetStatus()

		return resp, checkStatusDone(
			status,
			nil,
			DeleteErrorStatuses,
		)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"delete instance resource",
			fmt.Sprintf("instance %d: DELETE failed ", id)+err.Error(),
		)
	}
}

// Read implements resource.Resource.
func (g *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data InstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diag := getInstanceAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements resource.Resource.
func (g *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update instance resource",
		"update of 'instance' resources has not been implemented",
	)
}
