package testhelpers_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/subprovider"

	"github.com/hashicorp/terraform-plugin-testing/config"
	testresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type SubProviderTest struct {
	subprovider.SubProvider
}

func (t SubProviderTest) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		testhelpers.NewResource,
	}

	return resources
}

func New() *SubProviderTest {
	m := morpheus.New()
	t := SubProviderTest{SubProvider: m}

	return &t
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

func TestAccProviderBlockWithAccessToken(t *testing.T) {
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

func TestAccProviderBlockWithCredentials(t *testing.T) {
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
func TestAccProviderBlockAllAuth(t *testing.T) {
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

func TestAccProviderBlockMissingURL(t *testing.T) {
	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Must set a configuration value for the morpheus\[0\].url attribute as the\n` +
		`provider has marked it as required.`

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

func TestAccProviderBlockMissingAuth(t *testing.T) {
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

func TestAccProviderBlockMissingUsername(t *testing.T) {
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

func TestAccProviderBlockMissingPassword(t *testing.T) {
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

func TestAccProviderBlockNoneSet(t *testing.T) {
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Must set a configuration value for the morpheus\[0\].url attribute as the\n` +
		`provider has marked it as required.`

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
