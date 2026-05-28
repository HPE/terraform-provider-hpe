package option_list

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type optionListModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	SourceURL   types.String `tfsdk:"source_url"`
	Visibility  types.String `tfsdk:"visibility"`
	ApiType     types.String `tfsdk:"api_type"`
	RealTime    types.Bool   `tfsdk:"real_time"`
}

func OptionListSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Option Type List resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the option type list.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the option type list.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the option type list.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the option list (rest, manual, ldap, api).",
			},
			"source_url": schema.StringAttribute{
				Optional:    true,
				Description: "The source URL for the option list.",
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("private"),
				Description: "The visibility of the option type list.",
			},
			"api_type": schema.StringAttribute{
				Optional:    true,
				Description: "The API type of the option list.",
			},
			"real_time": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the option list is fetched in real time.",
			},
		},
	}
}
