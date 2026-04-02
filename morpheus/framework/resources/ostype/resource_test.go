// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	ostype "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusOsTypeExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := acctest.RandomWithPrefix("os.type")

	resourceConfig, err := ostype.RenderOsTypeConfig(t,
		map[string]string{
			"Name": name,
			"Code": code,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"platform",
			"linux",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"bit_count",
			"64",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"description",
			"An example OS type",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"os_family",
			"debian",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"os_version",
			"22.04",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"install_agent",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example",
			"cloud_init_version",
			"2",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_os_type.example",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_os_type.example",
			},
		},
	})
}

func TestAccMorpheusOsTypeUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := acctest.RandomWithPrefix("os.type")

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"platform",
			"linux",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"bit_count",
			"64",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	updatedName := name + "-updated"

	updateChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"name",
			updatedName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"description",
			"Updated description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"bit_count",
			"32",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.test",
			"platform",
			"windows",
		),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(updateChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_os_type" "test" {
  name      = "` + name + `"
  code      = "` + code + `"
  platform  = "linux"
  bit_count = 64
}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + `
resource "hpe_morpheus_os_type" "test" {
  name        = "` + name + `"
  code        = "` + code + `"
  platform    = "linux"
  bit_count   = 64
}`,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
resource "hpe_morpheus_os_type" "test" {
  name        = "` + updatedName + `"
  code        = "` + code + `"
  platform    = "windows"
  bit_count   = 32
  description = "Updated description"
}`,
				Check:    checkUpdateFn,
				PlanOnly: false,
			},
		},
	})
}

func TestAccMorpheusOsTypeRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	code := acctest.RandomWithPrefix("os.type")

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_os_type" "example_required" {
  name      = "` + name + `"
  code      = "` + code + `"
  platform  = "linux"
  bit_count = 64
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example_required",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example_required",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example_required",
			"platform",
			"linux",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type.example_required",
			"bit_count",
			"64",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_os_type.example_required",
			},
		},
	})
}
