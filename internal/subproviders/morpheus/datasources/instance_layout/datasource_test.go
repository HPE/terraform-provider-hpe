// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instance_layout_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name.tf.tmpl Name "Example name"
//go:generate go run ../../../../../cmd/render example-name-version.tf.tmpl Name "Example name" Version "1.2.3"

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/instance_layout"
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

func TestAccMorpheusFindInstanceLayoutById(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layout, err := testhelpers.CreateInstanceLayout(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteInstanceLayout(t, layout.GetId())
	})

	layoutID := fmt.Sprintf("%d", layout.GetId())
	layoutName := layout.GetName()

	config := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", layoutID)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
			layoutID,
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

func TestAccMorpheusFindInstanceLayoutByName(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layout, err := testhelpers.CreateInstanceLayout(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteInstanceLayout(t, layout.GetId())
	})

	layoutID := fmt.Sprintf("%d", layout.GetId())
	layoutName := layout.GetName()

	config := testhelpers.RenderExample(t, "example-name.tf.tmpl", "Name", layoutName)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
			layoutID,
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

func TestAccMorpheusFindInstanceLayoutByNameAndVersion(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layout, err := testhelpers.CreateInstanceLayout(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeleteInstanceLayout(t, layout.GetId())
	})

	layoutID := fmt.Sprintf("%d", layout.GetId())
	layoutName := layout.GetName()
	layoutVersion := layout.GetInstanceVersion()

	config := testhelpers.RenderExample(t, "example-name-version.tf.tmpl", "Name", layoutName, "Version", layoutVersion)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"version",
			layoutVersion,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
			layoutID,
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

func TestAccMorpheusFindInstanceLayoutNotFound(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := providerConfig + `
      data "hpe_morpheus_instance_layout" "test" {
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := instance_layout.ErrorNoInstanceLayoutFound

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

func TestAccMorpheusFindInstanceLayoutNoSearchAttrs(t *testing.T) {
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_instance_layout" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := instance_layout.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindInstanceLayoutWithIdAndName(t *testing.T) {
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_instance_layout" "test" {
        id = 1
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := instance_layout.ErrorRunningPreApply

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

func TestAccMorpheusFindInstanceLayoutWithIdAndVersion(t *testing.T) {
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_instance_layout" "test" {
        id = 1
        version = "123" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := instance_layout.ErrorRunningPreApply

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
