package testhelpers

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FakeResourceConfig() string {
	return `
resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
}

func fakeResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"testattr": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type fakeModel struct {
	Name     types.String `tfsdk:"name"`
	TestAttr types.String `tfsdk:"testattr"`
}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "hpe" + "_" + "morpheus" + "_" + "fake"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = fakeResourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data fakeModel
	req.Plan.Get(ctx, &data)

	c, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			"Unable to create client: "+err.Error(),
		)

		return
	}

	data.TestAttr = types.StringValue(c.GetConfig().Servers[0].URL)
	resp.State.Set(ctx, &data)
}

func (r *Resource) Read(
	_ context.Context,
	_ resource.ReadRequest,
	_ *resource.ReadResponse,
) {
}

func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
}

func (r *Resource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}
