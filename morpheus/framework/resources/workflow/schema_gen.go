package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type workflowModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Platform    types.String `tfsdk:"platform"`
	Visibility  types.String `tfsdk:"visibility"`
	Labels      types.List   `tfsdk:"labels"`
}

func WorkflowSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Workflow resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the workflow.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the workflow.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the workflow.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the workflow (provision or operation).",
			},
			"platform": schema.StringAttribute{
				Optional:    true,
				Description: "The platform of the workflow.",
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("private"),
				Description: "The visibility of the workflow (private or public).",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Labels for filtering.",
			},
		},
	}
}
