// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud_test

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
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud/consts"
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

func TestAccMorpheusFindCloudById(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	group, err := testhelpers.CreateGroup(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteGroup(t, group.GetId())
	})

	cloud, err := testhelpers.CreateCloud(t, group.GetId())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteCloud(t, cloud.GetId())
	})

	cloudID := fmt.Sprintf("%d", cloud.GetId())
	cloudName := cloud.GetName()

	config, err := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", cloudID)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"name",
			cloudName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
			cloudID,
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

func TestAccMorpheusFindCloudByName(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	group, err := testhelpers.CreateGroup(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteGroup(t, group.GetId())
	})

	cloud, err := testhelpers.CreateCloud(t, group.GetId())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteCloud(t, cloud.GetId())
	})

	cloudID := fmt.Sprintf("%d", cloud.GetId())
	cloudName := cloud.GetName()

	config, err := testhelpers.RenderExample(t, "example-name.tf.tmpl", "Name", cloudName)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"name",
			cloudName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
			cloudID,
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

func TestAccMorpheusFindCloudNotFound(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := providerConfig + `
      data "hpe_morpheus_cloud" "test" {
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoCloudFound

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

func TestAccMorpheusFindCloudNoSearchAttrs(t *testing.T) {
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cloud" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindCloudBothSearchAttrs(t *testing.T) {
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cloud" "test" {
        id = 1
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorRunningPreApply

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
