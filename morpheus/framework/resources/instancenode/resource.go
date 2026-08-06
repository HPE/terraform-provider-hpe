// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// Interface compliance assertions.
var (
	_ resource.Resource                     = &Resource{}
	_ resource.ResourceWithModifyPlan       = &Resource{}
	_ resource.ResourceWithConfigValidators = &Resource{}
)

// instanceNodeModel is the Terraform state model.
type instanceNodeModel struct {
	InstanceID       types.Int64    `tfsdk:"instance_id"`
	ResourcePoolID   types.Int64    `tfsdk:"resource_pool_id"`
	PreProvisioned   types.Bool     `tfsdk:"pre_provisioned"`
	SelectedServerID types.Int64    `tfsdk:"selected_server_id"`
	SshHost          types.String   `tfsdk:"ssh_host"`
	SshUsername      types.String   `tfsdk:"ssh_username"`
	SshPassword      types.String   `tfsdk:"ssh_password"`
	SshKeyPairID     types.Int64    `tfsdk:"ssh_key_pair_id"`
	WaitForIPAddress types.Bool     `tfsdk:"wait_for_ip_address"`
	ContainerID      types.Int64    `tfsdk:"container_id"`
	ServerID         types.Int64    `tfsdk:"server_id"`
	IPAddress        types.String   `tfsdk:"ip_address"`
	Hostname         types.String   `tfsdk:"hostname"`
	InternalIP       types.String   `tfsdk:"internal_ip"`
	ExternalFQDN     types.String   `tfsdk:"external_fqdn"`
	MacAddress       types.String   `tfsdk:"mac_address"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

// Resource implements the hpe_morpheus_instance_node resource.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata implements resource.Resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_instance_node"
}

// Schema implements resource.Resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Adds a single node to an existing instance. Volumes, networks, " +
			"and service plan are inherited from the instance and must not be set here. " +
			"Optionally, for HPE bare-metal instances, resource_pool_id can place the " +
			"node in a specific resource pool. For bare-metal instances, destroying " +
			"this resource returns the server to its pool rather than destroying it.",
		MarkdownDescription: "Adds a single node to an existing instance. Volumes, networks, " +
			"and service plan are inherited from the instance and must not be set here. " +
			"Optionally, for HPE bare-metal instances, `resource_pool_id` can place the " +
			"node in a specific resource pool. For bare-metal instances, destroying " +
			"this resource returns the server to its pool rather than destroying it.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The ID of the instance to add the node to.",
				MarkdownDescription: "The ID of the instance to add the node to.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"resource_pool_id": schema.Int64Attribute{
				Optional: true,
				Description: "The resource pool where the node will be placed. " +
					"Only honoured for HPE bare-metal instances; omit for " +
					"all other instance types.",
				MarkdownDescription: "The resource pool where the node will be placed. " +
					"Only honoured for HPE bare-metal instances; omit for " +
					"all other instance types.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"pre_provisioned": schema.BoolAttribute{
				Optional: true,
				Description: "Attach an existing server instead of " +
					"provisioning a new one.",
				MarkdownDescription: "Attach an existing server instead of " +
					"provisioning a new one.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"selected_server_id": schema.Int64Attribute{
				Optional: true,
				Description: "The ID of the pre-provisioned server. " +
					"Required when pre_provisioned is set.",
				MarkdownDescription: "The ID of the pre-provisioned server. " +
					"Required when `pre_provisioned` is set.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"ssh_host": schema.StringAttribute{
				Optional:            true,
				Description:         "SSH host for the pre-provisioned server.",
				MarkdownDescription: "SSH host for the pre-provisioned server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_username": schema.StringAttribute{
				Optional:            true,
				Description:         "SSH username for the pre-provisioned server.",
				MarkdownDescription: "SSH username for the pre-provisioned server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				WriteOnly: true,
				Description: "SSH password for the pre-provisioned server. " +
					"Write-only.",
				MarkdownDescription: "SSH password for the pre-provisioned server. " +
					"Write-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_key_pair_id": schema.Int64Attribute{
				Optional: true,
				Description: "The ID of the SSH key pair for the " +
					"pre-provisioned server.",
				MarkdownDescription: "The ID of the SSH key pair for the " +
					"pre-provisioned server.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"wait_for_ip_address": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Wait for the node to obtain an IP address " +
					"during create.",
				MarkdownDescription: "Wait for the node to obtain an IP address " +
					"during create.",
			},
			"container_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The container ID of the added node.",
				MarkdownDescription: "The container ID of the added node.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The compute server ID of the added node.",
				MarkdownDescription: "The compute server ID of the added node.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ip_address": schema.StringAttribute{
				Computed:            true,
				Description:         "The IP address of the node.",
				MarkdownDescription: "The IP address of the node.",
			},
			"hostname": schema.StringAttribute{
				Computed:            true,
				Description:         "The hostname of the node container.",
				MarkdownDescription: "The hostname of the node container.",
			},
			"internal_ip": schema.StringAttribute{
				Computed:            true,
				Description:         "The internal IP address of the node.",
				MarkdownDescription: "The internal IP address of the node.",
			},
			"external_fqdn": schema.StringAttribute{
				Computed:            true,
				Description:         "The external fully-qualified domain name of the node.",
				MarkdownDescription: "The external fully-qualified domain name of the node.",
			},
			"mac_address": schema.StringAttribute{
				Computed: true,
				Description: "The MAC address of the node server's primary network " +
					"interface. Only the primary interface address is surfaced; nodes " +
					"with bonded or multiple interfaces expose only this one.",
				MarkdownDescription: "The MAC address of the node server's primary network " +
					"interface. Only the primary interface address is surfaced; nodes " +
					"with bonded or multiple interfaces expose only this one.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (r *Resource) ConfigValidators(
	_ context.Context,
) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		preProvisionedRequiresServerID{},
	}
}

// preProvisionedRequiresServerID is a config validator that requires
// selected_server_id only when pre_provisioned is explicitly true.
// Unlike RequiredTogether, it does not fire when pre_provisioned is
// false or absent.
type preProvisionedRequiresServerID struct{}

func (v preProvisionedRequiresServerID) Description(_ context.Context) string {
	return "selected_server_id is required when pre_provisioned is true"
}

func (v preProvisionedRequiresServerID) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v preProvisionedRequiresServerID) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var preProvisioned types.Bool

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("pre_provisioned"), &preProvisioned)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If pre_provisioned is null, unknown, or false, no constraint applies.
	if preProvisioned.IsNull() || preProvisioned.IsUnknown() || !preProvisioned.ValueBool() {
		return
	}

	// pre_provisioned is true — selected_server_id must be set.
	var selectedServerID types.Int64

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("selected_server_id"), &selectedServerID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if selectedServerID.IsNull() || selectedServerID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("selected_server_id"),
			"selected_server_id is required when pre_provisioned is true",
			"When pre_provisioned is set to true, selected_server_id must "+
				"also be configured to identify the server to attach.",
		)
	}
}
