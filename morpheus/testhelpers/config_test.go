package testhelpers_test

import (
	"context"
	"regexp"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

type SubProviderTest struct {
	adapter.ProviderAdapter
	underlyingProvider fwprovider.Provider
}

var _ fwprovider.Provider = &SubProviderTest{}

func (t *SubProviderTest) Resources(
	ctx context.Context,
) []func() resource.Resource {
	var adaptedResources []func() resource.Resource

	for _, f := range t.underlyingProvider.Resources(ctx) {
		adaptedResources = append(
			adaptedResources,
			func() resource.Resource {
				return adapter.NewAdaptedResource(f(), &t.ProviderAdapter)
			},
		)
	}

	adaptedResources = append(
		adaptedResources,
		func() resource.Resource {
			return adapter.NewAdaptedResource(testhelpers.NewResource(), &t.ProviderAdapter)
		},
	)

	return adaptedResources
}

func New() *SubProviderTest {
	morpheusProvider := morpheus.New()

	return &SubProviderTest{
		ProviderAdapter:    *adapter.NewProviderAdapter(morpheusProvider),
		underlyingProvider: morpheusProvider,
	}
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

func TestAccMorpheusProviderBlockWithAccessToken(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	checks := []testresource.TestCheckFunc{
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"name",
			"bar",
		),
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": config.StringVariable("abcdefg"),
					"insecure":                      config.BoolVariable(false),
				},
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusProviderBlockWithCredentials(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	checks := []testresource.TestCheckFunc{
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"name",
			"bar",
		),
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     config.StringVariable("foo@test.com"),
					"testacc_morpheus_password":     config.StringVariable("testpass"),
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// if all access token and creds are provided, then it'll prefer access token
func TestAccMorpheusProviderBlockAllAuth(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	checks := []testresource.TestCheckFunc{
		testresource.TestCheckResourceAttr(
			"hpe_morpheus_fake.foo",
			"name",
			"bar",
		),
	}

	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     config.StringVariable("foo@test.com"),
					"testacc_morpheus_password":     config.StringVariable("testpass"),
					"testacc_morpheus_access_token": config.StringVariable("abcdefg"),
					"insecure":                      config.BoolVariable(false),
				},
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// missingConnectionDetails matches the summary line of the diagnostic the
// provider reports from Configure when neither "url" nor a usable identity
// block is set. "url" is Optional in the schema because an identity block may
// supply it, so this is no longer a schema-level required-attribute violation.
//
// Deliberately pinned to the summary line alone. The detail body carries
// example HCL and identity block attribute names, and asserting on that text
// is exactly what left these tests stale when the wording last moved.
const missingConnectionDetails = `Missing\s+Morpheus\s+connection\s+details`

func TestAccMorpheusProviderBlockMissingURL(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := missingConnectionDetails

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          nil,
					"testacc_morpheus_username":     config.StringVariable("foo@test.com"),
					"testacc_morpheus_password":     config.StringVariable("testpass"),
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          nil,
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": config.StringVariable("abcdefg"),
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusProviderBlockMissingAuth(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Attribute "morpheus\[0\].(username|access_token)" must be specified`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusProviderBlockMissingUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expectedA := `Attribute "morpheus\[0\].(username|access_token)" must be specified`
	expectedB := `(Attribute "morpheus\[0\].(username|access_token)" must be specified` +
		`|Attribute "morpheus\[0\].password" must be specified when\n` +
		`"morpheus\[0\].username" is specified)`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expectedA),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          config.StringVariable("https://test.morpheus.com"),
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     config.StringVariable("testpass"),
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expectedB),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusProviderBlockMissingPassword(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Attribute "morpheus\[0\].password" must be specified when\n` +
		`"morpheus\[0\].username" is specified`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          nil,
					"testacc_morpheus_username":     config.StringVariable("foo@test.com"),
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": nil,
					"insecure":                      config.BoolVariable(false),
				},
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusProviderBlockNoneSet(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := missingConnectionDetails

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ConfigVariables: config.Variables{
					"testacc_morpheus_url":          nil,
					"testacc_morpheus_username":     nil,
					"testacc_morpheus_password":     nil,
					"testacc_morpheus_access_token": nil,
					"insecure":                      nil,
				},
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
