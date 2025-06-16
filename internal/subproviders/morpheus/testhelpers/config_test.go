package testhelpers_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/subprovider"

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

// we explicitly unset env vars here, so do not run in parallel
func TestAccProviderBlockWithAccessToken(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_access_token", "abcdefg")
	t.Setenv("TF_VAR_testacc_morpheus_insecure", "false")

	username, setA := os.LookupEnv("TF_VAR_testacc_morpheus_username")
	password, setB := os.LookupEnv("TF_VAR_testacc_morpheus_password")

	os.Unsetenv("TF_VAR_testacc_morpheus_username")
	os.Unsetenv("TF_VAR_testacc_morpheus_password")
	t.Cleanup(func() {
		if setA {
			os.Setenv("TF_VAR_testacc_morpheus_username", username)
		}
		if setB {
			os.Setenv("TF_VAR_testacc_morpheus_password", password)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

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
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// we explicitly unset env vars here, so do not run in parallel
func TestAccProviderBlockWithCredentials(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_username", "foo@test.com")
	t.Setenv("TF_VAR_testacc_morpheus_password", "testpass")
	t.Setenv("TF_VAR_testacc_morpheus_insecure", "false")

	accessToken, set := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")

	os.Unsetenv("TF_VAR_testacc_morpheus_access_token")
	t.Cleanup(func() {
		if set {
			os.Setenv("TF_VAR_testacc_morpheus_access_token", accessToken)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

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
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// if all access token and creds are provided, then it'll prefer access token
func TestAccProviderBlockAllAuth(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_username", "foo@test.com")
	t.Setenv("TF_VAR_testacc_morpheus_password", "testpass")
	t.Setenv("TF_VAR_testacc_morpheus_access_token", "abcdefg")
	t.Setenv("TF_VAR_testacc_morpheus_insecure", "false")

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

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
	}
	checkFn := testresource.ComposeAggregateTestCheckFunc(checks...)
	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccProviderBlockMissingAuth(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")

	username, setA := os.LookupEnv("TF_VAR_testacc_morpheus_username")
	password, setB := os.LookupEnv("TF_VAR_testacc_morpheus_password")
	accessToken, setC := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")

	os.Unsetenv("TF_VAR_testacc_morpheus_username")
	os.Unsetenv("TF_VAR_testacc_morpheus_password")
	os.Unsetenv("TF_VAR_testacc_morpheus_access_token")
	t.Cleanup(func() {
		if setA {
			os.Setenv("TF_VAR_testacc_morpheus_username", username)
		}
		if setB {
			os.Setenv("TF_VAR_testacc_morpheus_password", password)
		}
		if setC {
			os.Setenv("TF_VAR_testacc_morpheus_access_token", accessToken)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Attribute "morpheus\[0\].(username|access_token)" must be specified`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccProviderBlockMissingPassword(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_username", "foo@test.com")

	password, setA := os.LookupEnv("TF_VAR_testacc_morpheus_password")
	accessToken, setB := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")

	os.Unsetenv("TF_VAR_testacc_morpheus_password")
	os.Unsetenv("TF_VAR_testacc_morpheus_access_token")
	t.Cleanup(func() {
		if setA {
			os.Setenv("TF_VAR_testacc_morpheus_password", password)
		}
		if setB {
			os.Setenv("TF_VAR_testacc_morpheus_access_token", accessToken)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Attribute "morpheus\[0\].password" must be specified when\n` +
		`"morpheus\[0\].username" is specified`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccProviderBlockMissingUsername(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_password", "testpass")

	username, setA := os.LookupEnv("TF_VAR_testacc_morpheus_username")
	accessToken, setB := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")

	os.Unsetenv("TF_VAR_testacc_morpheus_username")
	os.Unsetenv("TF_VAR_testacc_morpheus_access_token")
	t.Cleanup(func() {
		if setA {
			os.Setenv("TF_VAR_testacc_morpheus_username", username)
		}
		if setB {
			os.Setenv("TF_VAR_testacc_morpheus_access_token", accessToken)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `(Attribute "morpheus\[0\].(username|access_token)" must be specified` +
		`|Attribute "morpheus\[0\].password" must be specified when\n` +
		`"morpheus\[0\].username" is specified)`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccProviderBlockNoneSet(t *testing.T) {
	// can't unset these for t without a cleanup func,
	// but setting to "" is good enough

	url, setA := os.LookupEnv("TF_VAR_testacc_morpheus_url")
	username, setB := os.LookupEnv("TF_VAR_testacc_morpheus_username")
	password, setC := os.LookupEnv("TF_VAR_testacc_morpheus_password")
	accessToken, setD := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")

	os.Unsetenv("TF_VAR_testacc_morpheus_url")
	os.Unsetenv("TF_VAR_testacc_morpheus_username")
	os.Unsetenv("TF_VAR_testacc_morpheus_password")
	os.Unsetenv("TF_VAR_testacc_morpheus_access_token")
	t.Cleanup(func() {
		if setA {
			os.Setenv("TF_VAR_testacc_morpheus_url", url)
		}
		if setB {
			os.Setenv("TF_VAR_testacc_morpheus_username", username)
		}
		if setC {
			os.Setenv("TF_VAR_testacc_morpheus_password", password)
		}
		if setD {
			os.Setenv("TF_VAR_testacc_morpheus_access_token", accessToken)
		}
	})

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := testhelpers.FakeResourceConfig()

	expected := `Must set a configuration value for the morpheus\[0\].url attribute as the\n` +
		`provider has marked it as required.`

	testresource.Test(t, testresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testresource.TestStep{
			{
				ExpectError:        regexp.MustCompile(expected),
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
