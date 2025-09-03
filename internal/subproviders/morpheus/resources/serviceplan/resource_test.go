// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "ExampleServicePlan" Code "exampleserviceplan" MaxMemory "4294967296" MaxStorage "536870912"  ProvisionTypeCode "arm" CustomMaxStorage "true" ConfigRangesMinStorage "268435456" ConfigRangesMaxStorage "536870912"

package serviceplan_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
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

func TestAccMorpheusServicePlanRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	code := strings.ToLower(name)

	resourceConfig := `
resource "hpe_morpheus_service_plan" "example_required" {
  name                   = "` + name + `"
  code                   = "` + code + `"
  max_memory             = 4294967296
  max_storage            = 0
	provision_type_code    = "arm"
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_required",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_required",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_required",
			"max_memory",
			"4294967296",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_required",
			"max_storage",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_required",
			"provision_type_code",
			"arm",
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
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // Ensures the plan is empty, validating no drift
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_service_plan.example_required",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusServicePlanAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	resourceConfig := `
resource "hpe_morpheus_service_plan" "example_all" {
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
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"max_memory",
			"4294967296",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"max_storage",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"add_volumes",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"cores_per_socket",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"custom_cores",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"custom_cpu",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"custom_max_memory",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"custom_max_storage",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"description",
			"test serviceplan",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"max_cores",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"max_cpu",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"max_disks",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"memory_size_type",
			"mb",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"provision_type_code",
			"arm",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"sort_order",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"storage_size_type",
			"gb",
		),

		// Set (price_set_ids)
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"price_set_ids.#",
			"1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_service_plan.example_all",
			"price_set_ids.*",
			"1",
		),

		// Nested block (config_ranges)
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.min_storage",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.max_storage",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.min_memory",
			"1048576",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.max_memory",
			"2097152",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.min_cores",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.max_cores",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.min_sockets",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.max_sockets",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.min_cores_per_socket",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_all",
			"config_ranges.max_cores_per_socket",
			"10",
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
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // Ensures the plan is empty, validating no drift
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_service_plan.example_all",
				Check:             checkFn,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusServicePlanExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"Code", code,
		"MaxMemory", "4294967296",
		"MaxStorage", "536870912",
		"ProvisionTypeCode", "arm",
		"CustomMaxStorage", "true",
		"ConfigRangesMinStorage", "268435456",
		"ConfigRangesMaxStorage", "536870912")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"max_memory",
			"4294967296",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"max_storage",
			"536870912",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"provision_type_code",
			"arm",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"custom_max_storage",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"config_ranges.min_storage",
			"268435456",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_service_plan.example_service_plan",
			"config_ranges.max_storage",
			"536870912",
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
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_service_plan.example_service_plan",
				Check:             checkFn,
			},
		},
	})
}
