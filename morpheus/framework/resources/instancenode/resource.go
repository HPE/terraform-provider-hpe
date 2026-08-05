// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
		resourcevalidator.RequiredTogether(
			path.MatchRoot("pre_provisioned"),
			path.MatchRoot("selected_server_id"),
		),
	}
}
