package environment_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name.tf.tmpl Name "Example name"

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

const providerConfig = `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    username     = var.testacc_morpheus_username
    password     = var.testacc_morpheus_password
  }
}
`

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindEnvironmentById(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	environment, err := testhelpers.CreateEnvironment(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteEnvironment(t, environment.GetId())
	})

	environmentID := fmt.Sprintf("%d", environment.GetId())
	environmentName := environment.GetName()

	config, err := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", environmentID)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_environment.test",
			"name",
			environmentName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_environment.test",
			"id",
			environmentID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindIdbyName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	environment, err := testhelpers.CreateEnvironment(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteEnvironment(t, environment.GetId())
	})

	environmentID := fmt.Sprintf("%d", environment.GetId())
	environmentName := environment.GetName()

	config, err := testhelpers.RenderExample(t, "example-name.tf.tmpl", "Name", environmentName)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_environment.test",
			"name",
			environmentName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_environment.test",
			"id",
			environmentID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindEnvironmentNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := providerConfig + `
      data "hpe_morpheus_environment" "test" {
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_environment.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorNoEnvironmentFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindEnvironmentNoSearchAttrs(t *testing.T) {
	config := providerConfigOffline + `
      data "hpe_morpheus_environment" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_environment.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindEnvironmentBothSearchAttrs(t *testing.T) {
	config := providerConfigOffline + `
      data "hpe_morpheus_environment" "test" {
        id = 1
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_environment.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
