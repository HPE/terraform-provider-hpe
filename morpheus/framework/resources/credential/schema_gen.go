package credential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type credentialModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	IntegrationID types.Int64  `tfsdk:"integration_id"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
}

func CredentialSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Credential resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the credential.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the credential.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The credential type code.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the credential.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the credential is enabled.",
			},
			"integration_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The integration ID to associate with the credential.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The username or access key for the credential.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The password or secret key for the credential.",
			},
		},
	}
}
