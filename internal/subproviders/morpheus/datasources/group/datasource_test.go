package group_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const providerConfig = `
terraform {
  required_providers {
    hpe = {
      source  = "hpe/hpe"
      version = "0.0.1"
    }
  }
}

variable "morpheus_url" {}
variable "morpheus_access_token" {}
variable "morpheus_insecure" {}

provider "hpe" {
  morpheus {
    url          = var.morpheus_url
    access_token = var.morpheus_access_token
    insecure     = var.morpheus_insecure
  }
}
`

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusGroupAttrsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := providerConfig + `
      data "hpe_morpheus_group" "test" {
        id = 1
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"name",
			"Jasonx",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_group.test",
			"id",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   config,
				Check:    checkFn,
				PlanOnly: true,
			},
		},
	})
}
