package security_group_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type securityGroupRuleModel struct {
	ID              types.Int64  `tfsdk:"id"`
	SecurityGroupID types.Int64  `tfsdk:"security_group_id"`
	Name            types.String `tfsdk:"name"`
	Direction       types.String `tfsdk:"direction"`
	Protocol        types.String `tfsdk:"protocol"`
	RuleType        types.String `tfsdk:"rule_type"`
	Source          types.String `tfsdk:"source"`
	SourceType      types.String `tfsdk:"source_type"`
	Destination     types.String `tfsdk:"destination"`
	DestinationType types.String `tfsdk:"destination_type"`
	PortRange       types.String `tfsdk:"port_range"`
	Policy          types.String `tfsdk:"policy"`
}

func SecurityGroupRuleSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Security Group Rule resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the security group rule.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"security_group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the parent security group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the rule.",
			},
			"direction": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ingress"),
				Description: "The direction of the rule (ingress or egress).",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "The protocol (tcp, udp, icmp, etc.).",
			},
			"rule_type": schema.StringAttribute{
				Required:    true,
				Description: "The rule type.",
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "The source for the rule.",
			},
			"source_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The source type for the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination": schema.StringAttribute{
				Optional:    true,
				Description: "The destination for the rule.",
			},
			"destination_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The destination type for the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port_range": schema.StringAttribute{
				Optional:    true,
				Description: "The port range for the rule.",
			},
			"policy": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("accept"),
				Description: "The policy (accept or deny).",
			},
		},
	}
}
