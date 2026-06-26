// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus_test

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	providermod "github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
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

type FakeModel struct {
	Name     types.String `tfsdk:"name"`
	TestAttr types.String `tfsdk:"testattr"`
}

// For now, test using this method until we switch to testing the
// Morpheus Framework Provider with a flat provider config.
type SubProviderTest struct {
	adapter.ProviderAdapter
	underlyingProvider fwprovider.Provider
}

var _ fwprovider.Provider = &SubProviderTest{}

// Resources returns the Morpheus provider's resources plus the fake test resource.
func (t *SubProviderTest) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var adaptedResources []func() resource.Resource

	// Add Morpheus provider resources (adapted)
	for _, f := range t.underlyingProvider.Resources(ctx) {
		adaptedResources = append(
			adaptedResources,
			func() resource.Resource {
				return adapter.NewAdaptedResource(f(), &t.ProviderAdapter)
			},
		)
	}

	// Add fake test resource (also needs to be adapted for correct TypeName)
	adaptedResources = append(
		adaptedResources,
		func() resource.Resource {
			return adapter.NewAdaptedResource(NewResource(), &t.ProviderAdapter)
		},
	)

	return adaptedResources
}

func New() *SubProviderTest {
	morpheusProvider := morpheus.NewMorpheusProvider()
	return &SubProviderTest{
		ProviderAdapter:    *adapter.NewProviderAdapter(morpheusProvider),
		underlyingProvider: morpheusProvider,
	}
}

func NewWithCustomHTTPClient() *SubProviderTest {
	f := func(m model.MorpheusProviderModel) *clientfactory.ClientFactory {
		// example of passing in custom http client
		hc := &http.Client{}

		return clientfactory.New(
			m,
			clientfactory.WithFactoryHTTPClient(hc),
		)
	}

	morpheusProvider := morpheus.NewMorpheusProvider(morpheus.MorpheusWithClientFactory(f))
	return &SubProviderTest{
		ProviderAdapter:    *adapter.NewProviderAdapter(morpheusProvider),
		underlyingProvider: morpheusProvider,
	}
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := providermod.New("test", New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

// These tests are bad as the logic for parent provider config parsing by children
// is tightly coupled to the HPE provider implementation (the passing of []ListNestedBlock)
// to children as a flat Provider schema.
func TestAccMorpheusSubProviderMissingURL(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
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

func TestAccMorpheusSubProviderOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	providerConfig := `
provider "hpe" {
	morpheus {
		url = "https://127.0.0.1:0"
		username = "test-user"
		password = "test-password"
		insecure = true
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	testresource.TestCheckResourceAttr(
		"hpe_morpheus_fake.foo",
		"name",
		"bar",
	)
	checks := []testresource.TestCheckFunc{
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"name",
			"bar",
		),
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"testattr",
			"https://127.0.0.1:0",
		),
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// TestAccMorpheusSubProviderWithCustomHTTPClient is mainly
// an example of passing in a custom client
func TestAccMorpheusSubProviderWithCustomHTTPClient(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	newLocalProviderWithError := func() (tfprotov6.ProviderServer, error) {
		providerInstance := providermod.New("test", NewWithCustomHTTPClient())()

		return providerserver.NewProtocol6WithError(providerInstance)()
	}

	localTestAccProtoV6ProviderFactories := map[string]func() (
		tfprotov6.ProviderServer, error,
	){
		"hpe": newLocalProviderWithError,
	}

	providerConfig := `
provider "hpe" {
	morpheus {
		url = "https://127.0.0.1:0"
		username = "test-user"
		password = "test-password"
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	testresource.TestCheckResourceAttr(
		"hpe_morpheus_fake.foo",
		"name",
		"bar",
	)
	checks := []testresource.TestCheckFunc{
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"name",
			"bar",
		),
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"testattr",
			"https://127.0.0.1:0",
		),
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: localTestAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusSubProviderMissingAuth(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	providerConfig := `
provider "hpe" {
	morpheus {
		url = "http://127.0.0.1:0"
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `Attribute "morpheus\[0\].(username|access_token)" must be specified`
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

func TestAccMorpheusSubProviderMissingPassword(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	providerConfig := `
provider "hpe" {
	morpheus {
		url = "http://127.0.0.1:0"
		username = "test-user"
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `Attribute "morpheus\[0\]\.password" must be specified when\n` +
		`"morpheus\[0\]\.username" is specified`
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

func TestAccMorpheusSubProviderTooMuchAuth(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	providerConfig := `
provider "hpe" {
	morpheus {
		url = "http://example.com"
		username = "test-user"
		password = "test-password"
		tenant_subdomain = "foo"
		access_token = "this-is-not-a-token"
	}
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `Attribute "morpheus\[0\]\.(username|password|tenant_subdomain)" cannot be specified when\n` +
		`"morpheus\[0\]\.access_token" is specified`
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
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	providerConfig := `
provider "hpe" {
}

resource "hpe_morpheus_fake" "foo" {
	name = "bar"
}
`
	expected := `missing morpheus`
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
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
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
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
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
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "fake"
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
	var data FakeModel
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
