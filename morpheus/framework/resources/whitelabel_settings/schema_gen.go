package whitelabel_settings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type whitelabelSettingsModel struct {
	ID               types.String `tfsdk:"id"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	ApplianceName    types.String `tfsdk:"appliance_name"`
	HeaderLogo       types.String `tfsdk:"header_logo"`
	FooterLogo       types.String `tfsdk:"footer_logo"`
	LoginLogo        types.String `tfsdk:"login_logo"`
	Favicon          types.String `tfsdk:"favicon"`
	PrimaryColor     types.String `tfsdk:"primary_color"`
	SecondaryColor   types.String `tfsdk:"secondary_color"`
	SupportMenuLinks types.String `tfsdk:"support_menu_links"`
}

func WhitelabelSettingsSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages Morpheus Whitelabel Settings. This is a singleton resource — only one instance should exist.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The fixed identifier for the whitelabel settings singleton.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the whitelabel feature is enabled.",
			},
			"appliance_name": schema.StringAttribute{
				Optional:    true,
				Description: "The appliance name. Master account only.",
			},
			"header_logo": schema.StringAttribute{
				Optional:    true,
				Description: "The header logo URL or path.",
			},
			"footer_logo": schema.StringAttribute{
				Optional:    true,
				Description: "The footer logo URL or path.",
			},
			"login_logo": schema.StringAttribute{
				Optional:    true,
				Description: "The login logo URL or path.",
			},
			"favicon": schema.StringAttribute{
				Optional:    true,
				Description: "The favicon URL or path.",
			},
			"primary_color": schema.StringAttribute{
				Optional:    true,
				Description: "The primary button background color.",
			},
			"secondary_color": schema.StringAttribute{
				Optional:    true,
				Description: "The header background color.",
			},
			"support_menu_links": schema.StringAttribute{
				Optional:    true,
				Description: "Support menu links as a JSON string.",
			},
		},
	}
}
