// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package serviceplan_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name-provision.tf.tmpl Name "Example name" ProvisionTypeCode "arm"

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/serviceplan"
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

func TestAccMorpheusFindServicePlanById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	servicePlanName := acctest.RandomWithPrefix(t.Name())
	provisionTypeCode := "arm"

	// TODO: switch to using new provider for this test
	servicePlanResourceConfig := `
resource "morpheus_service_plan" "test" {
  name = "` + servicePlanName + `"
  code = "standard"
  price_set_ids  = []
  provision_type = "` + provisionTypeCode + `"
}
`

	providerConfig := testhelpers.ProviderBlockMixed()

	dataSourceConfig, err := testhelpers.RenderExample(
		t, "example-id.tf.tmpl", "Id", "morpheus_service_plan.test.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"name",
			servicePlanName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.3",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePlanResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindServicePlanByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	provisionTypeCode := "arm"

	// TODO: switch to using new provider for this test
	servicePlanResourceConfig := `
resource "morpheus_service_plan" "test" {
  name = "` + name + `"
  code = "standard"
  price_set_ids  = []
  provision_type = "` + provisionTypeCode + `"
}
`
	providerConfig := testhelpers.ProviderBlockMixed()

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name-provision.tf.tmpl",
		"Name", name,
		"ProvisionTypeCode", provisionTypeCode)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"provision_type_code",
			provisionTypeCode,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.3",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create dependency
			{
				Config: providerConfig + servicePlanResourceConfig,
			},
			{
				Config: providerConfig + servicePlanResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindServicePlanNoPlanFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
		data "hpe_morpheus_service_plan" "test" {
			name = "____"
			provision_type_code = "arm"
		}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorNoServicePlanFound

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

// TODO: Investigate why the logs for this case show `Error running pre-apply plan: exit status 1`
// The cloud test in comparison just exits with `exit status 1`
// Is it possibly to do with the additional validators on this data source?
func TestAccMorpheusFindServicePlanNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
	data "hpe_morpheus_service_plan" "test" {
	}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	// ExpectError only matches on ErrorRunningPreApply, not on serviceplan.ErrorNoValidSearchTerms
	expected := serviceplan.ErrorRunningPreApply

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

func TestAccMorpheusFindServicePlanBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_service_plan" "test" {
				id = "1"
				name = "_____"
				provision_type_code = "______"
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorRunningPreApply

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

func TestAccMorpheusFindServicePlanByProvisionOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_service_plan" "test" {
				provision_type_code = "arm"
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorRunningPreApply

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

// test to verify that all of the attributes from a created service plan can be read
func TestAccMorpheusFindServicePlanVerifyAttributes(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_service_plan" "test_all" {
  name                   = "` + name + `"
  code                   = "` + code + `"
  max_memory             = 4294967296
  max_storage            = 0
  add_volumes            = true
  cores_per_socket       = 1
  custom_cores           = true
  custom_cpu             = false
  custom_max_memory      = true
  custom_max_storage     = true
  description            = "test serviceplan"
  max_cores              = 1
  max_cpu                = 0
  max_disks              = 2
  memory_size_type       = "mb"
  price_set_ids           = [1]
  provision_type_code    = "arm"
  sort_order             = 0
  storage_size_type      = "gb"

  config_ranges =  {
    min_storage           = 1
    max_storage           = 2
    min_memory            = 1048576
    max_memory            = 2097152
    min_cores             = 1
    max_cores             = 2
    min_sockets           = 1
    max_sockets           = 10
    min_cores_per_socket  = 1
    max_cores_per_socket  = 10
    min_per_disk_size     = 1
    max_per_disk_size     = 2
  }
}
`
	dataSourceConfig := `
data "hpe_morpheus_service_plan" "test_all" {
  name = "` + name + `"
  provision_type_code = "arm"
}

`

	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"max_memory",
			"4294967296",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"max_storage",
			"0",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"add_volumes",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"cores_per_socket",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"custom_cores",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"custom_cpu",
			"false",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"custom_max_memory",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"custom_max_storage",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"description",
			"test serviceplan",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"max_cores",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"max_cpu",
			"0",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"max_disks",
			"2",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"memory_size_type",
			"mb",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"provision_type_code",
			"arm",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"sort_order",
			"0",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"storage_size_type",
			"gb",
		),

		// Set (price_set_ids)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"price_set_ids.#",
			"1",
		),
		resource.TestCheckTypeSetElemAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"price_set_ids.*",
			"1",
		),

		// Nested block (config_ranges)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_storage",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_storage",
			"2",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_memory",
			"1048576",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_memory",
			"2097152",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_cores",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_cores",
			"2",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_sockets",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_sockets",
			"10",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_cores_per_socket",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_cores_per_socket",
			"10",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.min_per_disk_size",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test_all",
			"config_ranges.max_per_disk_size",
			"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}
