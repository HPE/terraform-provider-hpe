package user_source

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userSourceModel struct {
	ID                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	AccountId             types.Int64  `tfsdk:"account_id"`
	Description           types.String `tfsdk:"description"`
	DefaultAccountRoleId  types.Int64  `tfsdk:"default_account_role_id"`
}

func UserSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Identity Source (User Source) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the identity source.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the identity source.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type code of the identity source (e.g. ldap, saml, activeDirectory).",
			},
			"account_id": schema.Int64Attribute{
				Required:    true,
				Description: "The account (tenant) ID to create the identity source under.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the identity source.",
			},
			"default_account_role_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the default account role.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
