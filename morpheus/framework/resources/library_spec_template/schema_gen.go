package library_spec_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type librarySpecTemplateModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Source     types.String `tfsdk:"source"`
	Content    types.String `tfsdk:"content"`
	ExternalID types.String `tfsdk:"external_id"`
}

func LibrarySpecTemplateSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Spec Template resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the spec template.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the spec template.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the spec template (terraform, kubernetes, helm, cloudformation, arm).",
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "The source type (local, repository, url).",
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the spec template.",
			},
			"external_id": schema.StringAttribute{
				Optional:    true,
				Description: "The external ID of the spec template.",
			},
		},
	}
}
