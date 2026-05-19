package certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type certificateModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CertFile    types.String `tfsdk:"cert_file"`
	KeyFile     types.String `tfsdk:"key_file"`
	DomainName  types.String `tfsdk:"domain_name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func CertificateSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Certificate resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the certificate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the certificate.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the certificate.",
			},
			"cert_file": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The contents of the certificate file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_file": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The contents of the key file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The domain name this certificate is tied to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the certificate is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
