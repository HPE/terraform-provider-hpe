package task

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

func (g *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "task"
}

func (g *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = TaskResourceSchema(ctx)
}
