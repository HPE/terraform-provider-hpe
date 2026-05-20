package key_pair

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type keyPairModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	PublicKey     types.String `tfsdk:"public_key"`
	PrivateKey    types.String `tfsdk:"private_key"`
	Passphrase    types.String `tfsdk:"passphrase"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
	Fingerprint   types.String `tfsdk:"fingerprint"`
}

func KeyPairSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Key Pair resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the key pair.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the key pair.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The public key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The private key. Write-only; not returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"passphrase": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The passphrase for the private key. Write-only; not returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"has_private_key": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the key pair has a private key.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "The fingerprint of the key pair.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
