// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name.tf.tmpl Name "\"Example name\""

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud/consts"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindCloudById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlockMixed()

	cloudResourceConfig := `
# assume tenant_id 1 exists
resource "morpheus_standard_cloud" "test_cloud" {
  name = "` + name + `"
  code = "standard"
  tenant_id = 1
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "morpheus_standard_cloud.test_cloud.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cloudResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindCloudByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlockMixed()

	cloudResourceConfig := `
# assume tenant_id 1 exists
resource "morpheus_standard_cloud" "test_cloud" {
  name = "` + name + `"
  code = "standard"
  tenant_id = 1
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "morpheus_standard_cloud.test_cloud.name")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cloudResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindCloudNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

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
	defer testhelpers.RecordResult(t)
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
	defer testhelpers.RecordResult(t)
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
