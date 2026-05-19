package account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type accountModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Subdomain   types.String `tfsdk:"subdomain"`
	Active      types.Bool   `tfsdk:"active"`
	RoleID      types.Int64  `tfsdk:"role_id"`
	Currency    types.String `tfsdk:"currency"`
}

func AccountSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Account (Tenant) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the account.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the account.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the account.",
			},
			"subdomain": schema.StringAttribute{
				Optional:    true,
				Description: "The subdomain for the account. This will be part of the login URL and username for sub tenant users.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the account is active.",
			},
			"role_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The default base role ID for the account.",
			},
			"currency": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("USD"),
				Description: "The currency code (ISO 4217) for the account.",
			},
		},
	}
}
