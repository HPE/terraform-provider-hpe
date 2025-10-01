package serviceplan_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

// TestAccMorpheusServicePlanResourceUpdateOk tests updating mutable attributes and detecting changes
func TestAccMorpheusServicePlanResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	uniqueName := acctest.RandomWithPrefix(t.Name())
	uniqueCode := strings.ToLower(uniqueName)

	baseConfigText := providerConfig + `
variable "name" { type = string }
variable "code" { type = string }
variable "description" {
  type    = string
  default = "test serviceplan"
}
variable "add_volumes" {
  type    = bool
  default = true
}
variable "custom_max_storage" {
  type    = bool
  default = true
}
variable "max_cores" {
  type    = number
  default = 1
}
variable "config_ranges_min_storage" {
  type    = number
  default = 1
}
variable "config_ranges_max_storage" {
  type    = number
  default = 2
}

resource "hpe_morpheus_service_plan" "test" {
	name                   = var.name
	code                   = var.code
	max_memory             = 4294967296
	max_storage            = 0
	add_volumes            = var.add_volumes
	cores_per_socket       = 1
	custom_cores           = true
	custom_cpu             = false
	custom_max_memory      = true
	custom_max_storage     = var.custom_max_storage
	description            = var.description
	max_cores              = var.max_cores
	max_cpu                = 0
	max_disks              = 2
	memory_size_type       = "mb"
	price_set_ids          = [1]
	provision_type_code    = "arm"
	sort_order             = 10000
	storage_size_type      = "gb"

	config_ranges =  {
		min_storage           = var.config_ranges_min_storage
		max_storage           = var.config_ranges_max_storage
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

	baseConfigVars := config.Variables{
		"name":                      config.StringVariable(uniqueName),
		"code":                      config.StringVariable(uniqueCode),
		"description":               config.StringVariable("test serviceplan"),
		"add_volumes":               config.BoolVariable(true),
		"custom_max_storage":        config.BoolVariable(true),
		"max_cores":                 config.IntegerVariable(1),
		"config_ranges_min_storage": config.IntegerVariable(1),
		"config_ranges_max_storage": config.IntegerVariable(2),
	}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "name", uniqueName),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "code", uniqueCode),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "description", "test serviceplan"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_memory", "4294967296"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_storage", "0"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "add_volumes", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "cores_per_socket", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_cores", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_cpu", "false"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_max_memory", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_max_storage", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_cores", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_cpu", "0"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_disks", "2"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "memory_size_type", "mb"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "sort_order", "10000"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "storage_size_type", "gb"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "price_set_ids.#", "1"),
		resource.TestCheckTypeSetElemAttr("hpe_morpheus_service_plan.test", "price_set_ids.*", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_storage", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_storage", "2"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_memory", "1048576"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_memory", "2097152"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_cores", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_cores", "2"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_sockets", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_sockets", "10"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_cores_per_socket", "1"),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_cores_per_socket", "10"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // create
				Config:          baseConfigText,
				ConfigVariables: baseConfigVars,
				Check:           checkFn,
			}, { // no-change plan
				Config:             baseConfigText,
				ConfigVariables:    baseConfigVars,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			}, { // description change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("Changed description"),
					"add_volumes":               config.BoolVariable(true),
					"custom_max_storage":        config.BoolVariable(true),
					"max_cores":                 config.IntegerVariable(1),
					"config_ranges_min_storage": config.IntegerVariable(1),
					"config_ranges_max_storage": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // add_volumes toggle (set false -> expect diff vs base true)
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("test serviceplan"),
					"add_volumes":               config.BoolVariable(false),
					"custom_max_storage":        config.BoolVariable(true),
					"max_cores":                 config.IntegerVariable(1),
					"config_ranges_min_storage": config.IntegerVariable(1),
					"config_ranges_max_storage": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // custom_max_storage toggle (set false)
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("test serviceplan"),
					"add_volumes":               config.BoolVariable(true),
					"custom_max_storage":        config.BoolVariable(false),
					"max_cores":                 config.IntegerVariable(1),
					"config_ranges_min_storage": config.IntegerVariable(1),
					"config_ranges_max_storage": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_cores change to 2
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("test serviceplan"),
					"add_volumes":               config.BoolVariable(true),
					"custom_max_storage":        config.BoolVariable(true),
					"max_cores":                 config.IntegerVariable(2),
					"config_ranges_min_storage": config.IntegerVariable(1),
					"config_ranges_max_storage": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // config_ranges min_storage change to 2
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("test serviceplan"),
					"add_volumes":               config.BoolVariable(true),
					"custom_max_storage":        config.BoolVariable(true),
					"max_cores":                 config.IntegerVariable(1),
					"config_ranges_min_storage": config.IntegerVariable(2),
					"config_ranges_max_storage": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // comprehensive apply with multiple updates
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(uniqueName),
					"code":                      config.StringVariable(uniqueCode),
					"description":               config.StringVariable("Comprehensive update test"),
					"add_volumes":               config.BoolVariable(false),
					"custom_max_storage":        config.BoolVariable(true),
					"max_cores":                 config.IntegerVariable(2),
					"config_ranges_min_storage": config.IntegerVariable(2),
					"config_ranges_max_storage": config.IntegerVariable(3),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "code", uniqueCode),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "description", "Comprehensive update test"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "add_volumes", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_max_storage", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_cores", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "price_set_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("hpe_morpheus_service_plan.test", "price_set_ids.*", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_storage", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_storage", "3"),
				),
			},
		},
	})
}

// Name/code update test (should be in-place updates)
func TestAccMorpheusServicePlanResourceUpdateNameCode(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	initialName := acctest.RandomWithPrefix(t.Name() + "-initial")
	initialCode := strings.ToLower(initialName)
	updatedName := acctest.RandomWithPrefix(t.Name() + "-updated")
	updatedCode := strings.ToLower(updatedName)

	configText := providerConfig + `
variable "name" { type = string }
variable "code" { type = string }

resource "hpe_morpheus_service_plan" "name_code_test" {
  name                = var.name
  code                = var.code
  max_memory          = 4294967296
  max_storage         = 0
  provision_type_code = "arm"
  sort_order          = 10000
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // create
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(initialName),
					"code": config.StringVariable(initialCode),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.name_code_test", "name", initialName),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.name_code_test", "code", initialCode),
				),
			}, { // no-change plan
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(initialName),
					"code": config.StringVariable(initialCode),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			}, { // change plan
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(updatedName),
					"code": config.StringVariable(updatedCode),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // apply changes
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(updatedName),
					"code": config.StringVariable(updatedCode),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.name_code_test", "name", updatedName),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.name_code_test", "code", updatedCode),
				),
			},
		},
	})
}
