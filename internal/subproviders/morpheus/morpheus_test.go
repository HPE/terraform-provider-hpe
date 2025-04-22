package morpheus_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func fakeResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

type FakeModel struct {
	Name types.String `tfsdk:"name"`
}

type SubProviderTest struct {
	morpheus.SubProvider
}

func (t SubProviderTest) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		NewResource,
	}

	return resources
}

func New() *SubProviderTest {
	m := morpheus.SubProvider{}
	t := SubProviderTest{SubProvider: m}

	return &t
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusSubProviderMissingURL(t *testing.T) {
	providerConfig := `
provider "hpe" {
	morpheus {
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `The argument "url" is required, but no definition was found`
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusSubProviderStrayResource(t *testing.T) {
	providerConfig := `
provider "hpe" {
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `missing morpheus provider block`
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig,
				ExpectError:        regexp.MustCompile(expected),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusSubProviderTooManyBlocks(t *testing.T) {
	providerConfig := `
provider "hpe" {
	morpheus {url = "https://example1.com"}
	morpheus {url = "https://example2.com"}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `Attribute morpheus list must contain` +
		` at least 0 elements and at most 1`
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig,
				ExpectError:        regexp.MustCompile(expected),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusSubProviderEmptyBlock checks that
// the absence of a block does not raise an error
func TestAccMorpheusSubProviderEmptyBlock(t *testing.T) {
	providerConfig := `
provider "hpe" {
}
`
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
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
	var plan FakeModel
	req.Plan.Get(ctx, &plan)

	_, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			"Unable to create client: "+err.Error(),
		)

		return
	}
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
