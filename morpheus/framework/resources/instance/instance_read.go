// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	awsCode    = "amazon"
	azureCode  = "azure"
	bmaasCode  = "hpe-baremetal-plugin.provision"
	hvmCode    = "mvm-cluster"
	kvmCode    = "kvm"
	vmwareCode = "vmware"
)

type apiConfigType = sdk.GetInstance200ResponseInstanceConfig

func instanceIDValue(instance sdk.GetInstance200ResponseInstance) int64 {
	if instance.Id == nil {
		return 0
	}

	return *instance.Id
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

	// Get timeout from HCL if set, the default is 45 minutes
	createTimeout, diags := data.Timeouts.Read(ctx, 45*time.Minute)
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

	// servicePlanOptions is not returned by the API
	servicePlanOptions := data.ServicePlanOptions

	state, found, diag := getInstanceAsState(ctx, data.Id.ValueInt64(), client, data, true)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	// The instance was deleted out of band (Morpheus returned 404); remove it
	// from state so the next apply recreates it instead of erroring.
	if !found {
		tflog.Warn(
			ctx,
			fmt.Sprintf("instance %d not found, removing from state", data.Id.ValueInt64()),
		)
		resp.State.RemoveResource(ctx)

		return
	}

	state.ServicePlanOptions = servicePlanOptions

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getInstanceAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan InstanceModel,
	// refresh is true for a plain Read (drift detection): API-backed fields like
	// volume storage_profile are taken from the API. It is false for create/update
	// post-apply reads, where the configured value is preferred so the final state
	// matches the plan.
	refresh bool,
) (InstanceModel, bool, diag.Diagnostics) {
	var state InstanceModel
	var diags diag.Diagnostics

	resp, hresp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		// The instance no longer exists in Morpheus (e.g. deleted out of band).
		// Signal not-found (without an error) so the caller can remove it from state.
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return state, false, diags
		}

		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, false, diags
	}

	if resp.Instance == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET returned empty instance", id),
		)

		return state, false, diags
	}

	instance := *resp.Instance

	// cloud_id
	state.CloudId = convert.Int64ToType(instance.Cloud.Id)

	// config
	state.Config = types.DynamicNull()
	state.ConfigAzure = NewConfigAzureValueNull()
	state.ConfigHvm = NewConfigHvmValueNull()
	state.ConfigVmware = NewConfigVmwareValueNull()
	state.ConfigBmaas = NewConfigBmaasValueNull()

	switch {
	case plan.Name.IsNull() || plan.Name.IsUnknown():
		// on import, always read the config from the API
		code, apiConfig, gdiags := getCodeAndConfig(id, instance)
		diags.Append(gdiags...)
		if diags.HasError() {
			return state, false, diags
		}

		switch *code {
		case awsCode:
			configAws, cdiags := getInstanceAWSConfig(ctx, id, apiConfig)
			diags = append(diags, cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.ConfigAws = configAws

		case azureCode:
			configAzure, cdiags := getInstanceAzureConfig(ctx, id, apiConfig)
			diags.Append(cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.ConfigAzure = configAzure

		case hvmCode:
			configHvm, cdiags := getInstanceHVMConfig(ctx, id, apiConfig)
			diags.Append(cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.ConfigHvm = configHvm

		case vmwareCode:
			configVMware, cdiags := getInstanceVMwareConfig(ctx, id, apiConfig)
			diags.Append(cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.ConfigVmware = configVMware

		case bmaasCode:
			configBmaas, cdiags := getInstanceBmaasConfig(ctx, id, apiConfig)
			diags.Append(cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.ConfigBmaas = configBmaas

		default:
			config, cdiags := getInstanceConfigGeneric(ctx, id, apiConfig)
			diags.Append(cdiags...)
			if diags.HasError() {
				return state, false, diags
			}
			state.Config = config

		}

	// For normal reads, use the config from the plan if it's set.
	case !plan.ConfigAws.IsNull() && !plan.ConfigAws.IsUnknown():
		state.ConfigAws = plan.ConfigAws

	case !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown():
		state.ConfigAzure = plan.ConfigAzure

	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		state.ConfigHvm = plan.ConfigHvm

	case !plan.ConfigVmware.IsNull() && !plan.ConfigVmware.IsUnknown():
		state.ConfigVmware = plan.ConfigVmware

	case !plan.ConfigBmaas.IsNull() && !plan.ConfigBmaas.IsUnknown():
		state.ConfigBmaas = plan.ConfigBmaas

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		state.Config = plan.Config
	}

	// connection_info
	cInfo, dc := getConnectionInfo(instance)
	diags.Append(dc...)
	if diags.HasError() {
		return state, false, diags
	}

	state.ConnectionInfo = cInfo

	// description
	if instance.Description.IsSet() {
		state.Description = convert.StrToType(instance.Description.Get())
	}

	// evars
	// API may respond with more evars than what the user set so we need to
	// check the /instance/{id}/envs endpoint which gives us the user specified
	// evars separately.
	envVars, diags := getInstanceEnvVars(ctx, id, client)
	if diags.HasError() {
		return state, false, diags
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
	if instance.Group.IsSet() && instance.Group.Get() != nil {
		state.GroupId = convert.Int64ToType(instance.Group.Get().Id)
	}

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

	// layout_size - from Config
	if instance.Config != nil {
		state.LayoutSize = convert.Int64ToType(instance.Config.LayoutSize)
	} else if !plan.LayoutSize.IsNull() && !plan.LayoutSize.IsUnknown() {
		// fallback to instance.layoutSize
		state.LayoutSize = plan.LayoutSize
	}

	// name
	state.Name = convert.StrToType(instance.Name)

	// host_name - Computed: populated from the API (Morpheus derives it from the
	// name when not explicitly set). Shared by Create/Read/Update via this mapper.
	state.HostName = convert.StrToType(instance.HostName)

	// labels - Computed set of organization labels; round-trips from the GET
	// response (empty -> null).
	state.Labels = convert.StrSliceToSet(instance.Labels)

	// server_uuids - RequiresReplace, create-only input. Preserve the incoming
	// value when the user set it (the API assigns exactly those UUIDs to the
	// servers, so preserving avoids any read-back mismatch). Otherwise read the
	// auto-generated UUIDs back from containerDetails[].server.uuid so the
	// Computed value is known after apply. It is an unordered set because Morpheus
	// does not guarantee containerDetails ordering matches the supplied order.
	if !plan.ServerUuids.IsNull() && !plan.ServerUuids.IsUnknown() {
		state.ServerUuids = plan.ServerUuids
	} else {
		state.ServerUuids = serverUUIDsFromContainerDetails(instance.ContainerDetails)
	}

	// status
	// Refreshed on every read so an out-of-band deletion of the underlying VM,
	// which Morpheus reports as "unknown" while retaining the instance record,
	// surfaces as a change on the next plan.
	state.Status = convert.StrToType(instance.Status)

	// interfaces
	ifaces, ifDiags := getStateInterfaces(ctx, instance, plan)
	diags = append(diags, ifDiags...)
	if diags.HasError() {
		return state, false, diags
	}

	networkInterfacesList, d := types.ListValueFrom(ctx, NetworkInterfacesValue{}.Type(ctx), ifaces)
	diags.Append(d...)

	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")

		return state, false, diags
	}

	state.NetworkInterfaces = networkInterfacesList

	// network_domain_id
	networkDomainId, ndiags := getInstanceNetworkDomainId(instance, plan)
	diags = append(diags, ndiags...)
	if diags.HasError() {
		return state, false, diags
	}
	state.NetworkDomainId = networkDomainId

	// user_group
	userGroupId, ugDiags := getInstanceUserGroupId(instance, plan)
	diags = append(diags, ugDiags...)
	if diags.HasError() {
		return state, false, diags
	}
	state.UserGroup = userGroupId

	// plan_id
	state.PlanId = convert.Int64ToType(instance.Plan.Id)

	// ports
	// assume ports always match the plan. Ports are a bit complicated in the
	// API. They are a part of container_details, which is an array which may
	// be of arbitrary size. The ports are contained within every element of the
	// array. The array of ports also may contain non-user controlled ports.
	// This should probably be replaced by a plan modifier.
	state.Ports = plan.Ports

	// timeouts
	state.Timeouts = plan.Timeouts

	// tags
	// we store the Name from the plan in the state, and compare the name returned from the API against
	// the state name while allowing for capitalisation changes
	planTagNames := make(map[string]string) // lowercase name -> plan name
	for _, planTag := range plan.Tags.Elements() {
		pt := planTag.(TagsValue)
		planTagNames[strings.ToLower(pt.Name.ValueString())] = pt.Name.ValueString()
	}

	tags, d := convert.ToSetType(
		ctx,
		instance.Tags,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner,
		) TagsValue {
			name := convert.StrToType(in.Name)
			if in.Name != nil {
				if planName, ok := planTagNames[strings.ToLower(*in.Name)]; ok {
					name = types.StringValue(planName)
				}
			}

			return TagsValue{
				Name:  name,
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
	volumes, d := getVolumes(ctx, instance, plan, refresh)
	diags.Append(d...)
	state.Volumes = volumes

	return state, true, diags
}

func getInstanceNetworkDomainId(
	instance sdk.GetInstance200ResponseInstance,
	plan InstanceModel,
) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	// if this isn't an import, return the plan value
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		return plan.NetworkDomainId, diags
	}

	// on import, get the network domain id from the config in the API response
	apiConfig := instance.Config
	if apiConfig == nil {
		diags.AddError(
			"get network_domain_id",
			fmt.Sprintf("instance %d config GET failed", instanceIDValue(instance)),
		)

		return types.Int64Null(), diags
	}

	if apiConfig.NetworkDomain == nil {
		return types.Int64Null(), diags
	}

	return convert.Int64ToType(apiConfig.NetworkDomain.Id), diags
}

// getInstanceUserGroupId returns the user_group id. On a normal read the plan
// value is preserved (user_group is provision-time only, so it never changes
// after create); on import it is read from the API config.
func getInstanceUserGroupId(
	instance sdk.GetInstance200ResponseInstance,
	plan InstanceModel,
) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	// if this isn't an import, return the plan value
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		return plan.UserGroup, diags
	}

	// on import, get the user group id from the config in the API response
	apiConfig := instance.Config
	if apiConfig == nil || apiConfig.UserGroup == nil {
		return types.Int64Null(), diags
	}

	return convert.Int64ToType(apiConfig.UserGroup.Id), diags
}

func getInstanceAzureConfig(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (ConfigAzureValue, diag.Diagnostics) {
	configAzure := ConfigAzureValue{}

	createUser, cdiags := getCreateUser(id, apiConfig)
	if cdiags.HasError() {
		return configAzure, cdiags
	}

	resourcePoolId, rdiags := getResourcePoolId(id, apiConfig)
	if rdiags.HasError() {
		return configAzure, rdiags
	}

	getAdditionalPropertyString := func(key string) types.String {
		if v, ok := apiConfig.AdditionalProperties[key]; ok && v != nil {
			return types.StringValue(fmt.Sprint(v))
		}

		return types.StringNull()
	}

	configAzure.AvailabilityOptions = getAdditionalPropertyString("availabilityOptions")
	configAzure.AvailabilitySet = getAdditionalPropertyString("availabilitySet")
	configAzure.AvailabilityZone = getAdditionalPropertyString("availabilityZone")
	configAzure.AzureRegion = getAdditionalPropertyString("azureRegion")
	configAzure.AzurefloatingIp = getAdditionalPropertyString("azurefloatingIp")
	configAzure.AzuresecurityGroupId = getAdditionalPropertyString("azuresecurityGroupId")
	configAzure.BootDiagnostics = getAdditionalPropertyString("bootDiagnostics")
	configAzure.CreateUser = convert.BoolToType(createUser)
	configAzure.DiagnosticsStorageAccount = getAdditionalPropertyString("diagnosticsStorageAccount")
	configAzure.OsGuestDiagnostics = getAdditionalPropertyString("osGuestDiagnostics")
	configAzure.ResourcePoolId = convert.StrToType(resourcePoolId)
	configAzure.state = attr.ValueStateKnown

	return configAzure, diag.Diagnostics{}
}

// getInstanceVMwareConfig builds the config_vmware block from the API response for vmware instances
func getInstanceVMwareConfig(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (ConfigVmwareValue, diag.Diagnostics) {
	configVmware := ConfigVmwareValue{}

	// CreateUser
	createUser, cdiags := getCreateUser(id, apiConfig)
	if cdiags.HasError() {
		return configVmware, cdiags
	}

	// NoAgent
	noAgent, ndiags := getNoAgent(id, apiConfig)
	if ndiags.HasError() {
		return configVmware, ndiags
	}

	// NestedVirtualization
	nestedVirtualization, ndiags := getNestedVirtualization(id, apiConfig)
	if ndiags.HasError() {
		return configVmware, ndiags
	}

	// ResourcePoolId
	resourcePoolId, rdiags := getResourcePoolId(id, apiConfig)
	if rdiags.HasError() {
		return configVmware, rdiags
	}

	// VMwareFolderId
	folderId := apiConfig.VmwareFolderId

	configVmware.CreateUser = convert.BoolToType(createUser)
	configVmware.NoAgent = convert.BoolToType(noAgent)
	configVmware.NestedVirtualization = convert.StrToType(nestedVirtualization)
	configVmware.ResourcePoolId = convert.StrToType(resourcePoolId)
	configVmware.VmwareFolderId = convert.StrToType(folderId)
	configVmware.state = attr.ValueStateKnown

	return configVmware, diag.Diagnostics{}
}

// getInstanceHVMConfig builds the config_hvm block from the API response for hvm instances
func getInstanceHVMConfig(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (ConfigHvmValue, diag.Diagnostics) {
	configHvm := ConfigHvmValue{}

	// CreateUser
	createUser, cdiags := getCreateUser(id, apiConfig)
	if cdiags.HasError() {
		return configHvm, cdiags
	}

	// NoAgent
	noAgent, ndiags := getNoAgent(id, apiConfig)
	if ndiags.HasError() {
		return configHvm, ndiags
	}

	// NestedVirtualization
	nestedVirtualization, ndiags := getNestedVirtualization(id, apiConfig)
	if ndiags.HasError() {
		return configHvm, ndiags
	}

	// ResourcePoolId
	resourcePoolId, rdiags := getResourcePoolId(id, apiConfig)
	if rdiags.HasError() {
		return configHvm, rdiags
	}

	// KvmHostId
	var kvmHostId *int64
	if apiConfig.KvmHostId.IsSet() {
		kvmHostId = apiConfig.KvmHostId.Get()
	}

	configHvm.CreateUser = convert.BoolToType(createUser)
	configHvm.NoAgent = convert.BoolToType(noAgent)
	configHvm.NestedVirtualization = convert.StrToType(nestedVirtualization)
	configHvm.ResourcePoolId = convert.StrToType(resourcePoolId)
	configHvm.KvmHostId = convert.Int64ToType(kvmHostId)
	configHvm.state = attr.ValueStateKnown

	return configHvm, diag.Diagnostics{}
}

// getInstanceAWSConfig builds the config_asw block from the API response for aws instances
func getInstanceAWSConfig(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (ConfigAwsValue, diag.Diagnostics) {
	configAws := ConfigAwsValue{}

	// CreateUser
	createUser, cdiags := getCreateUser(id, apiConfig)
	if cdiags.HasError() {
		return configAws, cdiags
	}

	// NoAgent
	noAgent, ndiags := getNoAgent(id, apiConfig)
	if ndiags.HasError() {
		return configAws, ndiags
	}

	// ResourcePoolId
	resourcePoolId, rdiags := getResourcePoolId(id, apiConfig)
	if rdiags.HasError() {
		return configAws, rdiags
	}

	// isEC2
	isEC2 := apiConfig.IsEC2

	// PublicIpType
	var publicIpType *string
	if apiConfig.PublicIpType.IsSet() {
		publicIpType = apiConfig.PublicIpType.Get()
	}

	// InstanceProfile
	var instanceProfile *string
	if apiConfig.InstanceProfile.IsSet() {
		instanceProfile = apiConfig.InstanceProfile.Get()
	}

	// KmsKeyId
	var kmsKeyId *string
	if apiConfig.KmsKeyId.IsSet() {
		kmsKeyId = apiConfig.KmsKeyId.Get()
	}

	// AvailabilityId
	var availabilityId *string
	if apiConfig.AvailabilityId.IsSet() {
		availabilityId = apiConfig.AvailabilityId.Get()
	}

	// SecGroups
	secGroupsList := types.ListNull(SecurityGroupsValue{}.Type(ctx))
	var sd diag.Diagnostics
	if apiConfig.SecurityGroups != nil {
		secGroupsList, sd = convert.ToListType(
			ctx,
			apiConfig.SecurityGroups,
			func(
				in sdk.AddInstance200ResponseAllOfOneOfInstanceConfigSecurityGroupsInner,
			) SecurityGroupsValue {
				v := SecurityGroupsValue{}
				v.Id = convert.StrToType(in.Id)
				v.state = attr.ValueStateKnown

				return v
			},
		)

		if sd.HasError() {
			return configAws, sd
		}
	}

	configAws.CreateUser = convert.BoolToType(createUser)
	configAws.NoAgent = convert.BoolToType(noAgent)
	configAws.ResourcePoolId = convert.StrToType(resourcePoolId)
	configAws.SecurityGroups = secGroupsList
	configAws.IsEc2 = convert.StringToBool(ctx, *isEC2)
	configAws.PublicIpType = convert.StrToType(publicIpType)
	// These next three can have a value of "" which corresponds to null
	configAws.KmsKeyId = basetypes.NewStringNull()
	if kmsKeyId != nil && *kmsKeyId != "" {
		configAws.KmsKeyId = convert.StrToType(kmsKeyId)
	}
	configAws.AvailabilityZoneId = basetypes.NewStringNull()
	if availabilityId != nil && *availabilityId != "" {
		configAws.AvailabilityZoneId = convert.StrToType(availabilityId)
	}
	configAws.InstanceProfile = basetypes.NewStringNull()
	if instanceProfile != nil && *instanceProfile != "" {
		configAws.InstanceProfile = convert.StrToType(instanceProfile)
	}
	configAws.state = attr.ValueStateKnown

	return configAws, diag.Diagnostics{}
}

// getInstanceBmaasConfig builds the config_bmaas block from the API response for
// HPE bare metal (BMaaS) instances. The baremetal plugin's option types are stored
// in the generic instance config map, so the plugin-specific fields (enforce RAID
// boot volume and selected hosts) are read from the config's untyped additional
// properties, while no_agent is a typed field on the instance config.
func getInstanceBmaasConfig(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (ConfigBmaasValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	configBmaas := ConfigBmaasValue{}

	noAgent, ndiags := getNoAgent(id, apiConfig)
	if ndiags.HasError() {
		return configBmaas, ndiags
	}

	createUser, cdiags := getCreateUser(id, apiConfig)
	if cdiags.HasError() {
		return configBmaas, cdiags
	}

	resourcePoolId, rdiags := getResourcePoolId(id, apiConfig)
	if rdiags.HasError() {
		return configBmaas, rdiags
	}

	selectedHosts, shdiags := selectedHostsFromConfig(
		ctx,
		apiConfig.AdditionalProperties["selectedHosts"],
	)
	diags.Append(shdiags...)
	if diags.HasError() {
		return configBmaas, diags
	}

	imageID := basetypes.NewInt64Null()
	if id, ok := numberToInt64(apiConfig.AdditionalProperties["imageId"]); ok {
		imageID = basetypes.NewInt64Value(id)
	}

	configBmaas.CreateUser = convert.BoolToType(createUser)
	configBmaas.ImageId = imageID
	configBmaas.NoAgent = convert.BoolToType(noAgent)
	configBmaas.ResourcePoolId = convert.StrToType(resourcePoolId)
	// enforce_raid_boot_volume defaults to true in the schema; mirror that default
	// when the value is absent from the config so an imported instance is stable.
	configBmaas.EnforceRaidBootVolume = boolFromConfig(
		apiConfig.AdditionalProperties["enforceRaidBootVolume"],
		true,
	)
	configBmaas.SelectedHosts = selectedHosts
	configBmaas.state = attr.ValueStateKnown

	return configBmaas, diags
}

// boolFromConfig coerces a value from the untyped instance config map into a bool
// value, tolerating the encodings Morpheus may use (native bool or strings such as
// "on"/"off"/"true"/"false"). The supplied default is used when the key is absent
// or cannot be interpreted.
func boolFromConfig(v interface{}, def bool) basetypes.BoolValue {
	switch val := v.(type) {
	case bool:
		return basetypes.NewBoolValue(val)
	case string:
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "on", "true", "1", "yes":
			return basetypes.NewBoolValue(true)
		case "off", "false", "0", "no":
			return basetypes.NewBoolValue(false)
		}
	}

	return basetypes.NewBoolValue(def)
}

// selectedHostsFromConfig parses the baremetal plugin's selectedHosts config value
// into a list of host ids. The plugin stores each entry as an object with a "value"
// field, but this tolerates bare ids and the various JSON number/string encodings.
// An absent or empty value yields a null list.
func selectedHostsFromConfig(
	ctx context.Context,
	v interface{},
) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	nullList := basetypes.NewListNull(types.Int64Type)

	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 {
		return nullList, diags
	}

	hostIDs := make([]int64, 0, len(raw))
	for _, elem := range raw {
		if hostID, ok := hostIDFromElement(elem); ok {
			hostIDs = append(hostIDs, hostID)
		}
	}

	if len(hostIDs) == 0 {
		return nullList, diags
	}

	list, d := types.ListValueFrom(ctx, types.Int64Type, hostIDs)
	diags.Append(d...)

	return list, diags
}

// hostIDFromElement extracts a host id from a selectedHosts element, which the
// baremetal plugin stores as an object with a "value" field, falling back to
// treating the element itself as the id.
func hostIDFromElement(elem interface{}) (int64, bool) {
	if m, ok := elem.(map[string]interface{}); ok {
		return numberToInt64(m["value"])
	}

	return numberToInt64(elem)
}

// numberToInt64 coerces the JSON-decoded representations of a number (float64 from
// encoding/json, native ints, or a numeric string) into an int64.
func numberToInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int64:
		return n, true
	case int:
		return int64(n), true
	case string:
		i, err := strconv.ParseInt(strings.TrimSpace(n), 10, 64)
		if err != nil {
			return 0, false
		}

		return i, true
	}

	return 0, false
}

func getCreateUser(
	id int64,
	apiConfig *apiConfigType,
) (*bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	createUser := apiConfig.CreateUser
	if createUser == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get config createUser", id),
		)

		return nil, diags
	}

	return createUser, nil
}

func getNoAgent(
	id int64,
	apiConfig *apiConfigType,
) (*bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	noAgent := apiConfig.NoAgent
	if noAgent == nil || noAgent.Bool == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get config noAgent", id),
		)

		return nil, diags

	}

	return noAgent.Bool, nil
}

func getNestedVirtualization(
	id int64,
	apiConfig *apiConfigType,
) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !apiConfig.NestedVirtualization.IsSet() {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get config nestedVirtualization", id),
		)

		return nil, diags
	}

	return apiConfig.NestedVirtualization.Get(), nil
}

func getResourcePoolId(
	id int64,
	apiConfig *apiConfigType,
) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourcePoolId := apiConfig.ResourcePoolId
	if resourcePoolId == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get config resourcePoolId", id),
		)

		return nil, diags
	}

	return resourcePoolId.String, nil
}

// getCodeAndConfig returns the "code" for the instance and the config struct from the API response
// The code is the layout provision type code for non-HVM clusters, and the cluster type code for HVM clusters.
func getCodeAndConfig(
	id int64,
	instance sdk.GetInstance200ResponseInstance,
) (*string, *apiConfigType, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get layout first
	layout := instance.Layout
	if layout == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get layout", id),
		)

		return nil, nil, diags
	}

	provisionTypeCode := layout.ProvisionTypeCode
	if provisionTypeCode == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get layout provision type code", id),
		)

		return nil, nil, diags
	}

	var code string
	switch *provisionTypeCode {
	case kvmCode:
		// If it's a kvm cluster, we need to check the cluster type code to see if it's actually an hvm cluster.
		// This is because the API returns "kvm" as the provision type code for both kvm and hvm clusters,
		// and we need to check the cluster type code to differentiate between them.
		cluster := instance.Cluster
		if cluster == nil {
			code = kvmCode

			break
		}

		clusterCode, cdiags := getClusterCode(id, cluster)
		if cdiags.HasError() {
			diags.Append(cdiags...)

			return nil, nil, diags
		}

		code = *clusterCode

	default:
		code = *provisionTypeCode
	}

	apiConfig := instance.Config
	if apiConfig == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get instance config", id),
		)

		return nil, nil, diags
	}

	return &code, apiConfig, diags
}

// getClusterCode returns the cluster type code for an instance if the cluster information is present
// in the API response
func getClusterCode(
	id int64,
	cluster *sdk.GetInstance200ResponseInstanceCluster,
) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	clusterType := cluster.Type
	if clusterType == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get cluster type", id),
		)

		return nil, diags
	}

	clusterCode := clusterType.Code
	if clusterCode == nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed to get cluster type code", id),
		)

		return nil, diags

	}

	return clusterCode, diags
}

// getInstanceConfigGeneric converts the instance config from the API response to a generic dynamic config
// for unknown instance types
func getInstanceConfigGeneric(
	ctx context.Context,
	id int64,
	apiConfig *apiConfigType,
) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert apiConfig to map[string]any
	apiConfigForConfig, err := apiConfig.ToMap()
	if err != nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d: failed to convert apiConfig to map: %s", id, err.Error()),
		)

		return types.DynamicNull(), diags
	}

	// Dereference pointers in apiConfigForConfig
	for k, v := range apiConfigForConfig {
		if v != nil {
			vType := reflect.TypeOf(v)
			if vType != nil && vType.Kind() == reflect.Pointer {
				vValue := reflect.ValueOf(v)
				if !vValue.IsNil() {
					if vValue.Elem().Kind() == reflect.Struct {
						// Convert struct to map
						apiConfigForConfig[k] = convertStructToMap(vValue.Elem().Interface())
					} else {
						apiConfigForConfig[k] = vValue.Elem().Interface()
					}
				}
			}
		}
	}

	configDynamic, err := convert.MapToDynamic(ctx, apiConfigForConfig)
	if err != nil {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d: failed to convert config to dynamic config: %s", id, err.Error()),
		)

		return types.DynamicNull(), diags
	}

	return configDynamic, diags
}

func convertStructToMap(s any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// Validate it's a struct
	if typ.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		fieldValue := val.Field(i).Interface()
		if reflect.TypeOf(fieldValue).Kind() == reflect.Pointer {
			if !reflect.ValueOf(fieldValue).IsNil() {
				fieldValue = reflect.ValueOf(fieldValue).Elem().Interface()
			} else {
				fieldValue = nil
			}
		}
		result[field.Name] = fieldValue
	}

	return result
}

func getInstanceEnvVars(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) ([]sdk.GetEnvVariables200ResponseEnvsInner, diag.Diagnostics) {
	var diags diag.Diagnostics
	resp, hresp, err := client.InstancesAPI.GetEnvVariables(ctx, id).Execute()

	// GET /api/instances/{id}/envs is gated by the "provisioning-environment"
	// permission, which requires the "environmentVariables" license feature.
	// Reduced editions that do not license it (e.g. HPE VM Essentials) — or a role
	// lacking the permission — disable the endpoint, so it returns a non-200
	// (MORPH-13817 saw a 404 on VME). On a nil response or any non-200 there are no
	// env vars to set: surface a warning for visibility, but never error (that
	// would fail the whole instance Read) and never dereference the nil resp (the
	// panic this guards against). The SDK's decode error is unreliable on valid
	// 200s (polymorphic fields it cannot model), so the HTTP status is the primary
	// signal and err is only used for context.
	if hresp == nil || hresp.StatusCode != http.StatusOK || resp == nil {
		diags.AddWarning(
			"Could not read instance environment variables",
			fmt.Sprintf("The environment variables endpoint returned no usable result "+
				"for instance %d, so no environment variables were read into state: %s",
				id, errfmt.ErrMsg(err, hresp)),
		)

		return nil, diags
	}

	return resp.Envs, diags
}

// serverUUIDsFromContainerDetails builds the server_uuids set from
// instance.containerDetails[].server.uuid, skipping containers with no server or
// no uuid. Returns a null set when no UUIDs are present. server_uuids is an
// unordered set because Morpheus does not guarantee containerDetails ordering.
func serverUUIDsFromContainerDetails(containers []sdk.InstanceContainer2) types.Set {
	uuids := make([]string, 0, len(containers))
	for _, cont := range containers {
		if cont.Server != nil && cont.Server.Uuid != nil {
			uuids = append(uuids, *cont.Server.Uuid)
		}
	}

	return convert.StrSliceToSet(uuids)
}

// getVolumes builds the volumes list from instance.containerDetails.server.volumes
func getVolumes(
	ctx context.Context,
	instance sdk.GetInstance200ResponseInstance,
	plan InstanceModel,
	refresh bool,
) (basetypes.ListValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Get volumes from instance.containerDetails.server.volumes
	contDetails := instance.ContainerDetails
	if len(contDetails) == 0 {
		diags.AddError(
			"cannot get instance containerDetails",
			fmt.Sprintf("instance %d GET containerDetails failed", instanceIDValue(instance)),
		)

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	server := contDetails[0].Server
	if server == nil {
		diags.AddError(
			"cannot get instance containerDetails server",
			fmt.Sprintf("instance %d GET containerDetails.server failed", instanceIDValue(instance)),
		)

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	serverVolumes := server.Volumes
	if serverVolumes == nil {
		diags.AddError(
			"cannot get instance containerDetails server volumes",
			fmt.Sprintf("instance %d GET containerDetails.server.volumes failed", instanceIDValue(instance)),
		)

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	// Remove any CD ROM volumes from the list
	apiVolumes := slices.DeleteFunc(
		serverVolumes,
		func(v sdk.InstanceContainerServerVolume1) bool {
			if v.Name == nil {
				return false
			}

			if strings.HasPrefix(*v.Name, "CD ROM") {
				return true
			}

			return false
		},
	)

	// Remove externally-attached storage-server (SAN) volumes — e.g. Alletra MP
	// BMaaS LUNs created by hpe_morpheus_storage_volume and exported to this
	// instance's host — which are not part of the instance's provisioned disks.
	apiVolumes = removeExternalStorageVolumes(apiVolumes)

	// Import
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		nonRaidVolumes := removeRaidDisks(apiVolumes)
		bootVolumeInFirst := bootVolumeFirst(nonRaidVolumes)

		return convertAPIVolumesToStateVolumes(ctx, bootVolumeInFirst)
	}

	// If the number of volumes is the same as the plan, we can do a direct conversion
	if len(apiVolumes) == len(plan.Volumes.Elements()) {
		autoselectVolumes := setDatastoreAutoSelectionAndSize(apiVolumes, plan, refresh)

		return convertAPIVolumesToStateVolumes(ctx, autoselectVolumes)
	}

	// The number of volumes is different to the plan
	nonRaidVolumes := removeRaidDisks(apiVolumes)
	reorderedVolumes := reorderVolumes(nonRaidVolumes, plan)
	filledVolumes := fillVolumeFieldsFromPlan(reorderedVolumes, plan)
	autoselectVolumes := setDatastoreAutoSelectionAndSize(filledVolumes, plan, refresh)

	return convertAPIVolumesToStateVolumes(ctx, autoselectVolumes)
}

// setDatastoreAutoSelectionAndSize sets the AdditionalProperties field for volumes with
// DatastoreAutoSelection set to that in the plan
// Also sets the MaxStorage field from the plan size
// Handles cases where the number of API volumes differs from plan volumes
func setDatastoreAutoSelectionAndSize(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
	refresh bool,
) []sdk.InstanceContainerServerVolume1 {
	autoSelection := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	planVolumes := plan.Volumes.Elements()

	// Determine how many volumes we can safely access from the plan
	maxIndex := len(apiVolumes)
	if len(planVolumes) < maxIndex {
		maxIndex = len(planVolumes)
	}

	for i, apiVol := range apiVolumes {
		// Only set fields from plan if we have a corresponding plan volume
		if i < maxIndex {
			planVol := planVolumes[i].(VolumesValue)

			// Initialize AdditionalProperties map if it doesn't exist
			if apiVol.AdditionalProperties == nil {
				apiVol.AdditionalProperties = make(map[string]interface{})
			}

			if !planVol.DatastoreAutoSelection.IsNull() && !planVol.DatastoreAutoSelection.IsUnknown() {
				// Set AdditionalProperties to indicate auto-selection
				apiVol.AdditionalProperties["DatastoreAutoSelection"] = planVol.DatastoreAutoSelection.ValueString()
			}

			apiVol.MaxStorage = planVol.Size.ValueInt64Pointer()
			// We set this flag to indicate that Terraform set the MaxStorage value
			apiVol.AdditionalProperties["TerraformSetMaxStorage"] = true

			// storage_profile: on a post-apply read (create/update) prefer the
			// configured value so the final state matches the plan, and an
			// API-side default for an unset volume is absorbed by the computed
			// attribute. On a plain refresh, leave the API value so a server-side
			// change is detected as drift.
			if !refresh && !planVol.StorageProfile.IsNull() && !planVol.StorageProfile.IsUnknown() {
				apiVol.StorageProfile = planVol.StorageProfile.ValueStringPointer()
			}
		}
		// If i >= maxIndex, just append the apiVol as-is (unmatched volumes)

		autoSelection = append(autoSelection, apiVol)
	}

	return autoSelection
}

// fillVolumeFieldsFromPlan fills in some missing fields in the API volumes from the plan
// This is needed because some fields are not returned by the API
// for certain volume types (e.g. Metal RAID volumes)
// We assume that the volumes are in the same order as the plan, but handle cases where
// the number of volumes differs
func fillVolumeFieldsFromPlan(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
) []sdk.InstanceContainerServerVolume1 {
	// Now fill in any missing fields from the plan
	filledVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	planVolumes := plan.Volumes.Elements()

	// Determine how many volumes we can safely fill from the plan
	maxIndex := len(apiVolumes)
	if len(planVolumes) < maxIndex {
		maxIndex = len(planVolumes)
	}

	for i, apiVol := range apiVolumes {
		// Only fill from plan if we have a corresponding plan volume
		if i < maxIndex {
			planVol := planVolumes[i].(VolumesValue)
			if apiVol.DatastoreId == nil {
				apiVol.DatastoreId = planVol.DatastoreId.ValueInt64Pointer()
			}
			// We set TypeId for all volume types since for some types (e.g. Metal RAID) the API TypeId is different
			apiVol.TypeId = planVol.StorageTypeId.ValueInt64Pointer()
		}
		// If i >= maxIndex, just append the apiVol as-is (unmatched volumes)

		filledVolumes = append(filledVolumes, apiVol)
	}

	return filledVolumes
}

// bootVolumeFirst puts the boot volume first in the list of volumes
// We hope to remove this function in future when the API names RAID volumes correctly
func bootVolumeFirst(
	apiVolumes []sdk.InstanceContainerServerVolume1,
) []sdk.InstanceContainerServerVolume1 {
	// If there are no volumes, return the original list
	if len(apiVolumes) == 0 {
		return apiVolumes
	}

	// Put boot volume first
	bootVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	otherVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	for _, v := range apiVolumes {
		if v.RootVolume != nil && *v.RootVolume {
			bootVolumes = append(bootVolumes, v)
		} else {
			otherVolumes = append(otherVolumes, v)
		}
	}

	// Combine boot volumes first, then other volumes
	result := append(bootVolumes, otherVolumes...)

	return result
}

// reorderVolumes re-orders the list of volumes to match the plan
func reorderVolumes(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
) []sdk.InstanceContainerServerVolume1 {
	// If there are no volumes, return empty list
	if len(apiVolumes) == 0 {
		return apiVolumes
	}

	// Track which API volumes have been matched to avoid duplicates
	matchedVolumes := make(map[int]bool)

	// Now re-order to match plan
	orderedVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))

	// Match remaining volumes with plan volumes by name
	for _, planVol := range plan.Volumes.Elements() {
		planVolTyped := planVol.(VolumesValue)
		for i, apiVol := range apiVolumes {
			if matchedVolumes[i] {
				continue // Skip already matched volumes
			}

			if apiVol.Name != nil && planVolTyped.Name.ValueString() == *apiVol.Name {
				orderedVolumes = append(orderedVolumes, apiVol)
				matchedVolumes[i] = true

				break
			}
		}
	}

	// Append any unmatched volumes at the end to avoid data loss
	for i, apiVol := range apiVolumes {
		if !matchedVolumes[i] {
			orderedVolumes = append(orderedVolumes, apiVol)
		}
	}

	return orderedVolumes
}

// removeExternalStorageVolumes drops storage-server (SAN) volumes from the list.
// A volume that belongs to a storage server — e.g. an Alletra MP BMaaS LUN
// created by hpe_morpheus_storage_volume and exported to this instance's host —
// is added to the compute server's volume collection as a side effect of that
// export (see updateVolumeAttachment in the Alletra MP plugin). Such volumes are
// not part of the instance's provisioned disks, so including them here would
// cause spurious drift on this resource's volumes (and the resize-on-count-change
// update path could try to detach them). Provisioned VM/metal disks are managed
// via a datastore/zone and carry no storageServer, so this only excludes
// externally-managed array LUNs.
func removeExternalStorageVolumes(
	apiVolumes []sdk.InstanceContainerServerVolume1,
) []sdk.InstanceContainerServerVolume1 {
	return slices.DeleteFunc(
		apiVolumes,
		func(v sdk.InstanceContainerServerVolume1) bool {
			return v.StorageServer != nil
		},
	)
}

// removeRaidDisks removes any RAID disks from the list of volumes
func removeRaidDisks(
	apiVolumes []sdk.InstanceContainerServerVolume1,
) []sdk.InstanceContainerServerVolume1 {
	// build a map of device-counts for the API volumes
	deviceCount := make(map[string]int)
	for _, volume := range apiVolumes {
		if volume.DeviceName != nil {
			deviceCount[*volume.DeviceName]++
		}
	}

	// remove the RAID disks from the apiVolumes list
	nonRaidDiskVolumes := slices.DeleteFunc(
		apiVolumes,
		func(v sdk.InstanceContainerServerVolume1) bool {
			// Skip volumes without a device name
			if v.DeviceName == nil {
				return false
			}

			if deviceCount[*v.DeviceName] > 1 {
				// We're going to remove volumes which have a diskType
				if v.DiskType != nil && v.DiskMode == nil {
					return true
				}
			}

			return false
		},
	)

	return nonRaidDiskVolumes
}

// convertBytesPtrToGBBytes converts a pointer to int64 bytes to a pointer to int64 GB bytes
func convertBytesPtrToGBBytes(b *int64) *int64 {
	if b == nil {
		return nil
	}
	gb := *b / (1 << 30)

	return &gb
}

func convertAPIVolumesToStateVolumes(
	ctx context.Context,
	apiVolumes []sdk.InstanceContainerServerVolume1,
) (basetypes.ListValue, diag.Diagnostics) {
	volumes, d := convert.ToListType(
		ctx,
		apiVolumes,
		func(
			in sdk.InstanceContainerServerVolume1,
		) VolumesValue {
			v := VolumesValue{}
			v.Id = convert.Int64ToType(in.Id)
			v.RootVolume = convert.BoolToType(in.RootVolume)
			v.Name = convert.StrToType(in.Name)
			v.StorageTypeId = convert.Int64ToType(in.TypeId)
			v.DatastoreId = convert.Int64ToType(in.DatastoreId)
			v.ControllerMountPoint = convert.StrToType(in.ControllerMountPoint)
			v.StorageProfile = convert.StrToType(in.StorageProfile)

			// Handle DatastoreAutoSelection and TerraformSetMaxStorage from AdditionalProperties
			// TerraformSetMaxStorage flag indicates that MaxStorage was set from plan (already in GB)
			// and should not be converted from bytes
			terraformSetMaxStorage := false
			if in.AdditionalProperties != nil {
				if dsAutoSel, ok := in.AdditionalProperties["DatastoreAutoSelection"]; ok {
					if dsAutoSelStr, ok := dsAutoSel.(string); ok {
						v.DatastoreAutoSelection = convert.StrToType(&dsAutoSelStr)
					}
				}

				if tsms, ok := in.AdditionalProperties["TerraformSetMaxStorage"]; ok {
					if tsmsBool, ok := tsms.(bool); ok {
						terraformSetMaxStorage = tsmsBool
					}
				}
			}

			// Set Size: if TerraformSetMaxStorage is true, MaxStorage is already in GB from plan
			// Otherwise, convert from bytes to GB
			if terraformSetMaxStorage {
				v.Size = convert.Int64ToType(in.MaxStorage)
			} else {
				v.Size = convert.Int64ToType(convertBytesPtrToGBBytes(in.MaxStorage))
			}

			v.state = attr.ValueStateKnown

			return v
		},
	)

	return volumes, d
}

// getConnectionInfo builds the connection_info list
func getConnectionInfo(
	instance sdk.GetInstance200ResponseInstance,
) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	cInfo := instance.ConnectionInfo
	if cInfo == nil {
		diags.AddError(
			"cannot get instance connectionInfo",
			fmt.Sprintf("instance %d GET connectionInfo failed", instanceIDValue(instance)),
		)

		return types.ListNull(types.StringType), diags
	}

	if len(cInfo) == 0 {
		return types.ListNull(types.StringType), diags
	}

	var vals []attr.Value
	for _, c := range cInfo {
		if c.Ip != nil {
			vals = append(vals, types.StringValue(*c.Ip))
		}
	}

	cList, dl := types.ListValue(types.StringType, vals)
	diags = append(diags, dl...)

	return cList, diags
}

// getStateInterfaces get the interfaces to be returned as state entries
func getStateInterfaces(
	ctx context.Context,
	instance sdk.GetInstance200ResponseInstance,
	plan InstanceModel,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Generate []NetworkInterfacesValue from instance.interfaces
	intfsFromInstance, id := getStateInterfacesFromInstance(ctx, instance)
	diags = append(diags, id...)
	if diags.HasError() {
		return nil, diags
	}

	// If this is an import, then return intfsFromInstance
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		return intfsFromInstance, diags
	}

	// Generate []NetworkInterfacesValue from the instance.containerDetails.server.interfaces
	intfsFromServer, id := getStateInterfacesFromInstanceServer(ctx, instance)
	diags = append(diags, id...)
	if diags.HasError() {
		return nil, diags
	}

	// Get []NetworkInterfacesValue from the plan
	var intfsFromPlan []NetworkInterfacesValue
	pd := plan.NetworkInterfaces.ElementsAs(ctx, &intfsFromPlan, false)
	if pd.HasError() {
		return nil, pd
	}

	// Compare intfsFromServer against intfsFromPlan, to see if the "shapes" are the same.
	// subnet_id is read back from the server interface itself (see
	// getStateInterfacesFromInstanceServer / getChildNetworks), so no plan-preservation
	// is required.
	if compareServerPlanIntfs(intfsFromServer, intfsFromPlan) {
		return intfsFromServer, diags
	}

	// "Shape" isn't the same, return intfsFromPlan
	return intfsFromPlan, diags
}

// compareServerPlanIntfs compares the []NetworkInterfacesValues from instance.containerDetails.server.interfaces
// and plan see if they are the same shape
// Returns true if they are, false otherwise
func compareServerPlanIntfs(
	intfsFromServer, intfsFromPlan []NetworkInterfacesValue,
) bool {
	// Check length of lists first
	if len(intfsFromServer) != len(intfsFromPlan) {
		return false
	}

	// Get list of lengths of child interfaces for instance.containerDetails.server.interfaces list
	serverSubIntfs := make([]int, 0, len(intfsFromServer))
	for _, serverIntf := range intfsFromServer {
		serverSubIntfs = append(serverSubIntfs, len(serverIntf.ChildVirtualNetworks.Elements()))
	}

	// Get list of lengths of child interfaces for plan list
	planSubIntfs := make([]int, 0, len(intfsFromPlan))
	for _, planIntf := range intfsFromPlan {
		planSubIntfs = append(planSubIntfs, len(planIntf.ChildVirtualNetworks.Elements()))
	}

	// Compare lengths of child interfaces for "server" and "instance" lists
	for i := range serverSubIntfs {
		if serverSubIntfs[i] != planSubIntfs[i] {
			return false
		}
	}

	return true
}

// getStateInterfacesFromInstance build []NetworkInterfacesValue from interfaces, used on import
func getStateInterfacesFromInstance(
	ctx context.Context,
	instance sdk.GetInstance200ResponseInstance,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	instIntfs := instance.Interfaces
	if instIntfs == nil {
		diags.AddError(
			"instance GetInterfaces failed",
			fmt.Sprintf("instance %d GET interfaces failed", instanceIDValue(instance)),
		)

		return nil, diags
	}

	var ifaces []NetworkInterfacesValue
	for _, instIntf := range instIntfs {
		ifaceVal := NetworkInterfacesValue{}
		ifaceVal.Id = types.Int64Null()
		ifaceVal.IpAddress = convert.StrToType(instIntf.IpAddress)
		ifaceVal.IpMode = convert.StrToType(instIntf.IpMode)
		ifaceVal.PrimaryInterface = types.BoolNull()
		ifaceVal.Name = types.StringNull()
		ifaceVal.NetworkId = types.Int64Null()
		// subnet_id is not recoverable on import: the instance.interfaces list used
		// here does not include the subnet (only the server-interface read path does).
		ifaceVal.SubnetId = types.Int64Null()
		if net := instIntf.Network; net != nil {
			ifaceVal.NetworkId = convert.Int64ToType(net.Id)
			ifaceVal.IpPool = types.Int64Null()
			if pool := net.Pool; pool != nil {
				ifaceVal.IpPool = convert.Int64ToType(pool.Id)
			}
			ifaceVal.NetworkGroupId = convert.Int64ToType(net.Group)
		}
		ifaceVal.NetworkTypeId = networkTypeId(instIntf.NetworkInterfaceTypeId)
		ifaceVal.ChildVirtualNetworks = types.ListNull(ChildVirtualNetworksValue{}.Type(ctx))
		if len(instIntf.NetworkInterfaces) > 0 {
			childNetworks, cd := getInstanceInterfacesChildNetworks(ctx, instIntf.NetworkInterfaces)
			ifaceVal.ChildVirtualNetworks = childNetworks
			diags = append(diags, cd...)
		}
		ifaceVal.state = attr.ValueStateKnown

		ifaces = append(ifaces, ifaceVal)
	}

	return ifaces, diags
}

// networkTypeId helper function to handle NetworkInterfaceTypeId
// We only use this function on import, and it looks as if NetworkInterfaceTypeId will have a value of 0 instead
// of null.  In this case we return a null value to avoid inadvertent instance recreations when the HCL has been
// generated on import and a plan/apply following the initial import is performed
func networkTypeId(i *int64) basetypes.Int64Value {
	if i != nil && *i == 0 {
		return basetypes.NewInt64Null()
	}

	return convert.Int64ToType(i)
}

// getInstanceInterfacesChildNetworks returns child_networks from interfaces.networkInterfaces, used on import
func getInstanceInterfacesChildNetworks(
	ctx context.Context,
	nets []sdk.InstanceInterfacesNetworkInterfacesInner1,
) (basetypes.ListValue, diag.Diagnostics) {
	children := make([]ChildVirtualNetworksValue, 0)
	for _, instIntf := range nets {
		ifaceVal := ChildVirtualNetworksValue{}
		ifaceVal.Id = types.Int64Null()
		ifaceVal.IpAddress = convert.StrToType(instIntf.IpAddress)
		ifaceVal.IpMode = convert.StrToType(instIntf.IpMode)
		ifaceVal.PrimaryInterface = types.BoolNull()
		ifaceVal.Name = types.StringNull()
		ifaceVal.NetworkId = types.Int64Null()
		// subnet_id is not recoverable on import: the instance.interfaces list used
		// here does not include the subnet (only the server-interface read path does).
		ifaceVal.SubnetId = types.Int64Null()
		if net := instIntf.Network; net != nil {
			ifaceVal.NetworkId = convert.Int64ToType(net.Id)
			ifaceVal.IpPool = types.Int64Null()
			if pool := net.Pool; pool != nil {
				ifaceVal.IpPool = convert.Int64ToType(pool.Id)
			}
			ifaceVal.NetworkGroupId = convert.Int64ToType(net.Group)
		}
		ifaceVal.NetworkTypeId = networkTypeId(instIntf.NetworkInterfaceTypeId)

		ifaceVal.state = attr.ValueStateKnown
		children = append(children, ifaceVal)
	}

	return types.ListValueFrom(ctx, ChildVirtualNetworksValue{}.Type(ctx), children)
}

// getStateInterfacesFromInstanceServer get the []NetworkInterfacesValue from containerDetails.server.interfaces
func getStateInterfacesFromInstanceServer(
	ctx context.Context,
	instance sdk.GetInstance200ResponseInstance,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	// network_interfaces
	// We are going to read network interface information from containerDetails.server.interfaces
	// Note that, at present, all network IP addresses will not be available to us when we reach
	// this stage on instance creation, we will have enough in the state-file that a plan will
	// be a no-op, and that when all IP addresses are available (this can be seen in the UI) an
	// apply will update the state-file with the IP addresses etc.
	procIntfs := getAllServerInterfaces(instance)

	var ifaces []NetworkInterfacesValue
	var childInterfaces basetypes.ListValue
	var diags diag.Diagnostics

	for _, iface := range procIntfs.serverIntfsList {
		// Skip sub-interfaces
		if _, ok := procIntfs.isSubIntf[*iface.Id]; ok {
			continue
		}
		ifaceVal := NetworkInterfacesValue{}
		ifaceVal.Id = convert.Int64ToType(iface.Id)
		ifaceVal.IpAddress = convert.StrToType(iface.IpAddress)
		ifaceVal.IpMode = convert.StrToType(iface.IpMode)
		ifaceVal.NetworkGroupId = types.Int64Null()
		if iface.NetworkGroup != nil {
			ifaceVal.NetworkGroupId = convert.Int64ToType(iface.NetworkGroup.Id)
		}
		if iface.NetworkPool != nil {
			ifaceVal.IpPool = convert.Int64ToType(iface.NetworkPool.Id)
		}
		if iface.Network != nil {
			ifaceVal.NetworkId = convert.Int64ToType(iface.Network.Id)
		}
		// subnet_id is read back from the interface's subnet association. The API
		// resolves a subnet to its parent network (reported as network_id) but also
		// returns the subnet itself, so subnet_id round-trips on refresh.
		ifaceVal.SubnetId = types.Int64Null()
		if iface.Subnet != nil {
			ifaceVal.SubnetId = convert.Int64ToType(iface.Subnet.Id)
		}
		ifaceVal.Name = convert.StrToType(iface.Name)
		ifaceVal.PrimaryInterface = convert.BoolToType(iface.PrimaryInterface)

		childInterfaces, diags = getChildNetworks(ctx, iface.Id, procIntfs.subIntfsMap, procIntfs.serverIntfsMap)

		ifaceVal.ChildVirtualNetworks = childInterfaces

		ifaceVal.state = attr.ValueStateKnown

		ifaces = append(ifaces, ifaceVal)
	}

	return ifaces, diags
}

// processedServerInterfaces struct that contains maps and a list produced from containerDetails.server.interfaces
type processedServerInterfaces struct {
	// subIntfsMap a map of interface-ids with a list of the ids of any sub-interfaces
	// note that it is possible for two (or maybe more) interfaces to have the same sub-interfaces (bonds)
	subIntfsMap map[int64][]int64
	// isSubIntf a map of interface-ids with a boolean saying if they are sub-interfaces
	isSubIntf map[int64]bool
	// serverIntfsMap a map of interface-ids with the corresponding interface information
	serverIntfsMap map[int64]sdk.InstanceContainerServerInterfacesInner1
	// serverIntfsList a list of the interfaces, which should (hopefully be in the same order as those specified
	// in network_interfaces
	serverIntfsList []sdk.InstanceContainerServerInterfacesInner1
}

// Process the set of interfaces in an instance
// This function takes an "instance" input and returns processedServerInterfaces
func getAllServerInterfaces(
	instance sdk.GetInstance200ResponseInstance,
) processedServerInterfaces {
	subIntfsMap := make(map[int64][]int64)
	isSubIntf := make(map[int64]bool)
	serverIntfsMap := make(map[int64]sdk.InstanceContainerServerInterfacesInner1)
	serverIntfsList := make([]sdk.InstanceContainerServerInterfacesInner1, 0)

	// First clean-up the server interface list
	// The list of interfaces is malleable, and changes after the instance has been created and returned
	// to us for reading.  In early stages the interface "name" (eth0, eth1 etc) will have repeated
	// entries.
	// Key here is the "UniqueId".  If an interface doesn't have a value for that or for Network then all
	// it has is an IP address that will be assigned to the interface with the same name (eth0, eth1, etc)
	for _, container := range instance.ContainerDetails {
		server := container.Server
		if server == nil {
			continue
		}
		serverIntfList := server.Interfaces
		serverIntfsNameMap := make(map[string][]sdk.InstanceContainerServerInterfacesInner1)
		serverIntfsNameListPosition := make([]string, 0)
		serverIntfsNameListMap := make(map[string]struct{})
		serverIntfsMergedNameMap := make(map[string]sdk.InstanceContainerServerInterfacesInner1)
		for _, serverIntf := range serverIntfList {
			// Skip this list entry if it doesn't have a name
			if serverIntf.Name == nil {
				continue
			}

			serverIntfsNameMap[*serverIntf.Name] = append(serverIntfsNameMap[*serverIntf.Name], serverIntf)
			// Keep a record of the order of the interface name ("eth0" etc) in the input list
			// We allow for duplicate name entries, so we're looking for the first entry
			if _, ok := serverIntfsNameListMap[*serverIntf.Name]; !ok {
				serverIntfsNameListPosition = append(serverIntfsNameListPosition, *serverIntf.Name)
				serverIntfsNameListMap[*serverIntf.Name] = struct{}{}
			}
		}

		for intfName, v := range serverIntfsNameMap {
			if len(v) == 1 {
				serverIntfsMergedNameMap[intfName] = v[0]

				continue
			}

			// Hopefully there will only be two entries in the other lists
			// What we are going to do here is use the entry that has Network information
			// as the base of the cumulative interface, and then hunt through the rest for
			// an ip-address
			var cumulativeIntf sdk.InstanceContainerServerInterfacesInner1
			var ipAddress *string
			for _, serverIntf := range v {
				if serverIntf.Network != nil {
					cumulativeIntf = serverIntf

					break
				}
			}
			for _, serverIntf := range v {
				if serverIntf.IpAddress != nil {
					ipAddress = serverIntf.IpAddress

					break
				}
			}

			cumulativeIntf.IpAddress = ipAddress
			serverIntfsMergedNameMap[intfName] = cumulativeIntf
		}

		// Order serverIntfsList by order that the entries appeared in the original list
		for _, intfName := range serverIntfsNameListPosition {
			serverIntfsList = append(serverIntfsList, serverIntfsMergedNameMap[intfName])
		}
	}

	// Build the maps that we're going to return
	for _, serverInterface := range serverIntfsList {
		if serverInterface.Id == nil {
			continue
		}

		serverIntfsMap[*serverInterface.Id] = serverInterface
		if serverInterface.Interfaces != nil {
			intfList := make([]int64, 0)
			for _, subIntf := range serverInterface.Interfaces {
				if subIntf.Id == nil {
					continue
				}
				intfList = append(intfList, *subIntf.Id)
				isSubIntf[*subIntf.Id] = true
			}
			if len(intfList) > 0 {
				subIntfsMap[*serverInterface.Id] = intfList
			}
		}
	}

	ret := processedServerInterfaces{}
	ret.subIntfsMap = subIntfsMap
	ret.isSubIntf = isSubIntf
	ret.serverIntfsMap = serverIntfsMap
	ret.serverIntfsList = serverIntfsList

	return ret
}

// Get the child virtual network interface values
func getChildNetworks(
	ctx context.Context,
	id *int64,
	subIntfMap map[int64][]int64,
	serverIntfsMap map[int64]sdk.InstanceContainerServerInterfacesInner1,
) (basetypes.ListValue, diag.Diagnostics) {
	if id == nil {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), nil
	}

	if len(subIntfMap[*id]) == 0 {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), nil
	}

	children := make([]ChildVirtualNetworksValue, 0)
	for _, subIntf := range subIntfMap[*id] {
		ifaceVal := ChildVirtualNetworksValue{}
		iface := serverIntfsMap[subIntf]
		ifaceVal.Id = convert.Int64ToType(iface.Id)
		ifaceVal.IpAddress = convert.StrToType(iface.IpAddress)
		ifaceVal.IpMode = convert.StrToType(iface.IpMode)
		ifaceVal.NetworkGroupId = types.Int64Null()
		if iface.NetworkGroup != nil {
			ifaceVal.NetworkGroupId = convert.Int64ToType(iface.NetworkGroup.Id)
		}
		if iface.NetworkPool != nil {
			ifaceVal.IpPool = convert.Int64ToType(iface.NetworkPool.Id)
		}
		ifaceVal.NetworkId = convert.Int64ToType(iface.Network.Id)
		// subnet_id round-trips from the interface's subnet association (see
		// getStateInterfacesFromInstanceServer).
		ifaceVal.SubnetId = types.Int64Null()
		if iface.Subnet != nil {
			ifaceVal.SubnetId = convert.Int64ToType(iface.Subnet.Id)
		}
		ifaceVal.Name = convert.StrToType(iface.Name)
		ifaceVal.PrimaryInterface = convert.BoolToType(iface.PrimaryInterface)
		ifaceVal.state = attr.ValueStateKnown
		children = append(children, ifaceVal)
	}

	return types.ListValueFrom(ctx, ChildVirtualNetworksValue{}.Type(ctx), children)
}
