package user_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userGroupModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SudoUser    types.Bool   `tfsdk:"sudo_user"`
	ServerGroup types.String `tfsdk:"server_group"`
}

func UserGroupSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus User Group resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the user group.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the user group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the user group.",
			},
			"sudo_user": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the user group has sudo access.",
			},
			"server_group": schema.StringAttribute{
				Optional:    true,
				Description: "The server group associated with the user group.",
			},
		},
	}
}
