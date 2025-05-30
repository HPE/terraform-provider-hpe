package role_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

// Check that we can create a role with only
// required attributes specified
func TestAccMorpheusRoleRequiredAttrsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}
variable "testacc_morpheus_insecure" {}

provider "hpe" {
	morpheus {
		url = var.testacc_morpheus_url
		username = var.testacc_morpheus_username
		password = var.testacc_morpheus_password
		insecure = var.testacc_morpheus_insecure
	}
}

resource "hpe_morpheus_role" "foo" {
	name = "testacc-TestAccMorpheusRoleRequiredAttrsOk"
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"name",
			"testacc-TestAccMorpheusRoleRequiredAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"role_type",
			"user",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.foo",
			"description",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_role.foo",
				Check:             checkFn,
			},
		},
	})
}

// Check that we can create a role with all attributes specified
func TestAccMorpheusRoleAllAttrsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}
variable "testacc_morpheus_insecure" {}

provider "hpe" {
	morpheus {
		url = var.testacc_morpheus_url
		username = var.testacc_morpheus_username
		password = var.testacc_morpheus_password
		insecure = var.testacc_morpheus_insecure
	}
}

resource "hpe_morpheus_role" "foo" {
	name = "testacc-TestAccMorpheusRoleAllAttrsOk"
	description = "test"
	multitenant = true
	role_type = "user"
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"name",
			"testacc-TestAccMorpheusRoleAllAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"description",
			"test",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"multitenant",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.foo",
			"role_type",
			"user",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_role.foo",
				Check:             checkFn,
			},
		},
	})
}

// TODO: Add more acceptance tests
