// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	readOperation     = "read cluster resource"
	populateOperation = "populate cluster resource"
)

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	readTimeout, diags := data.Timeouts.Read(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diag := getClusterAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getClusterAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan ClusterModel,
) (ClusterModel, diag.Diagnostics) {
	var state ClusterModel
	var diags diag.Diagnostics

	importing := plan.Name.IsNull()

	clusterResp, httpResp, err := client.ClustersAPI.GetCluster(ctx, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("cluster %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	cluster := clusterResp.GetCluster()

	// top-level fields
	state.Id = convert.Int64ToType(cluster.Id)

	clusterTypeCodeFromName := clusterTypeCodeForName(cluster.Type.GetName())
	state.ClusterTypeCode = convert.StrToType(&clusterTypeCodeFromName)

	state.CloudId = convert.Int64ToType(cluster.Zone.Id)
	state.Uuid = convert.StrToType(cluster.Uuid)
	state.Description = convert.StrToType(cluster.Description.Get())
	state.GroupId = convert.Int64ToType(cluster.Site.Id)
	state.LayoutId = convert.Int64ToType(cluster.Layout.Id)
	state.Name = convert.StrToType(cluster.Name)
	state.ServiceUrl = convert.StrToType(cluster.ServiceUrl.Get())

	// requires replace, and not present in GET
	state.WorkflowId = plan.WorkflowId

	// timeout
	state.Timeouts = plan.Timeouts

	// handle different types of cluster configs...
	switch {
	case clusterTypeCodeFromName == clusterTypeCodeMVM &&
		(!plan.ConfigHvm.IsNull() || importing):

		// doesn't require replace - cpu_arch, cpu_model, dynamic_placement
		var cpuArchVal string
		if v, ok := cluster.Config["cpuArch"].(string); ok {
			cpuArchVal = v
		}

		var cpuModelVal string
		if v, ok := cluster.Config["cpuModel"].(string); ok {
			cpuModelVal = v
		}

		var dynamicPlacementModeVal string
		if v, ok := cluster.Config["dynamicPlacementMode"].(string); ok {
			dynamicPlacementModeVal = v
		}

		var powerPolicyVal string
		if v, ok := cluster.Config["powerPolicy"].(string); ok {
			powerPolicyVal = v
		}

		var configHvmVal ConfigHvmValue
		if importing {

			var computeInterfaceNameVal string
			if v, ok := cluster.Config["computeInterfaceName"].(string); ok {
				computeInterfaceNameVal = v
			}

			var storageInterfaceNameVal string
			if v, ok := cluster.Config["storageInterfaceName"].(string); ok {
				storageInterfaceNameVal = v
			}

			var overlayInterfaceNameVal string
			if v, ok := cluster.Config["overlayInterfaceName"].(string); ok {
				overlayInterfaceNameVal = v
			}

			var computeVlansVal string
			if v, ok := cluster.Config["computeVlans"].(string); ok {
				computeVlansVal = v
			}

			// read as many values as possible from API
			configHvmVal, diags = NewConfigHvmValue(ConfigHvmValue{}.AttributeTypes(ctx), map[string]attr.Value{
				// can't read createUser from config GET
				"compute_interface_name": convert.StrToType(&computeInterfaceNameVal),
				"storage_interface_name": convert.StrToType(&storageInterfaceNameVal),
				"overlay_interface_name": convert.StrToType(&overlayInterfaceNameVal),
				"compute_vlans":          convert.StrToType(&computeVlansVal),
				// can't read vcpuPlacementMode from config GET
				"dynamic_placement": convert.StringToBool(ctx, dynamicPlacementModeVal),
				"cpu_arch":          convert.StrToType(&cpuArchVal),
				"cpu_model":         convert.StrToType(&cpuModelVal),
				"power_policy":      convert.StrToType(&powerPolicyVal),
			})

		} else {
			configHvmVal, diags = NewConfigHvmValue(ConfigHvmValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"create_user":            plan.ConfigHvm.CreateUser,
				"compute_interface_name": plan.ConfigHvm.ComputeInterfaceName,
				"storage_interface_name": plan.ConfigHvm.StorageInterfaceName,
				"overlay_interface_name": plan.ConfigHvm.OverlayInterfaceName,
				"compute_vlans":          plan.ConfigHvm.ComputeVlans,
				"vcpu_placement_mode":    plan.ConfigHvm.VcpuPlacementMode,
				"dynamic_placement":      convert.StringToBool(ctx, dynamicPlacementModeVal),
				"cpu_arch":               convert.StrToType(&cpuArchVal),
				"cpu_model":              convert.StrToType(&cpuModelVal),
				"power_policy":           convert.StrToType(&powerPolicyVal),
			})
		}

		if diags.HasError() {
			return state, diags
		}

		state.ConfigHvm = configHvmVal
	// for importing to static blocks, handle generic config last
	case !plan.Config.IsNull() || importing:
		state.Config = basetypes.NewDynamicNull()

		cfg := cluster.GetConfig()
		if cfg == nil {
			diags.AddError(
				readOperation,
				"cluster: generic config missing from API response",
			)

			return state, diags
		}

		// Set plan config to state if it's not null.
		// This means that on Read, we don't get clashes with
		// additional or missing properties from the GET config block.
		// i.e. only manage in state what we've set in the config.
		if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
			state.Config = plan.Config
		} else {
			// importing
			o, err := convert.MapToDynamic(ctx, cfg)
			if err != nil {
				diags.AddError(populateOperation, err.Error())
			}

			state.Config = o
		}

	}

	server, diags := buildReadClusterServer(ctx, plan)
	if diags != nil {
		diags.Append(diags...)

		return state, diags
	}

	state.Server = server

	// labels
	respLabels := cluster.GetLabels()
	labels, err := convert.SetToStrSlice(plan.Labels)
	if err != nil {
		diags.AddError(
			populateOperation,
			"could not parse a slice of labels",
		)

		return state, diags
	}

	// Morpheus API may change the casing of the labels, to avoid Terraform
	// throwing a gasket we convert the casing of labels to be as specified
	// by the user.
	for _, label := range labels {
		for i, respLabel := range respLabels {
			if strings.EqualFold(label, respLabel) {
				if label != respLabel {
					respLabels[i] = label
				}
			}
		}
	}

	state.Labels = convert.StrSliceToSet(respLabels)

	return state, diags
}

func buildReadClusterServer(
	ctx context.Context,
	plan ClusterModel,
) (ServerValue, diag.Diagnostics) {
	importing := plan.Name.IsNull()

	if importing {
		// TODO: investigate which properties we can actually set from an import, if any.
		return NewServerValue(ServerValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	} else {
		return NewServerValue(ServerValue{}.AttributeTypes(ctx), map[string]attr.Value{
			// the requires replace properties...
			"data_device":              plan.Server.DataDevice,
			"hostname":                 plan.Server.Hostname,
			"lvm_enabled":              plan.Server.LvmEnabled,
			"name":                     plan.Server.Name,
			"management_net_interface": plan.Server.ManagementNetInterface,
			"network_domain":           plan.Server.NetworkDomain,
			"network_interfaces":       plan.Server.NetworkInterfaces,
			"security_groups":          plan.Server.SecurityGroups,
			"service_plan_id":          plan.Server.ServicePlanId,
			"ssh_port":                 plan.Server.SshPort,
			"ssh_password_wo":          plan.Server.SshPasswordWo, // write-only
			"ssh_password_wo_version":  plan.Server.SshPasswordWoVersion,
			"ssh_hosts":                plan.Server.SshHosts,
			"ssh_key_pair_id":          plan.Server.SshKeyPairId,
			"ssh_username":             plan.Server.SshUsername,
			"tags":                     plan.Server.Tags,
			"user_group_id":            plan.Server.UserGroupId,
			"visibility":               plan.Server.Visibility,
			"volumes":                  plan.Server.Volumes,
		})
	}
}

// Maintains a local mapping of cluster type names to their codes
// so that we don't have to perform additional API calls.
// The "type" field in cluster GET only provides name and id,
// so without this we'd need to get the code from an additional API call.
func clusterTypeCodeForName(name string) string {
	switch name {
	case "HVM":
		return clusterTypeCodeMVM
	case "Kubernetes Cluster":
		return clusterTypeCodeKubernetes
	case "AKS Cluster":
		return clusterTypeCodeAKS
	case "GKE Cluster":
		return clusterTypeCodeGKE
	case "External Kubernetes Cluster":
		return clusterTypeCodeExternalKubernetes
	case "EKS Cluster":
		return clusterTypeCodeEKS
	case "Docker Cluster":
		return clusterTypeCodeDocker
	default:
		return "Unknown cluster type name to code mapping."
	}
}
