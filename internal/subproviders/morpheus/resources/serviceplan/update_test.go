package serviceplan_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

// TestAccMorpheusServicePlanResourceUpdateAllAttrsOk tests updating all mutable attributes and detecting changes
func TestAccMorpheusServicePlanResourceUpdateAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	uniqueName := acctest.RandomWithPrefix(t.Name())

	baseConfigText := providerConfig + `
variable "name" { type = string }
variable "code" { 
  type = string 
  default = "test-serviceplan-code"
}
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
variable "provision_type_code" {
  type    = string
  default = "kvm"
}
variable "max_memory" {
  type    = number
  default = 4294967296
}
variable "max_storage" {
  type    = number
  default = 0
}
variable "cores_per_socket" {
  type    = number
  default = 1
}
variable "custom_cores" {
  type    = bool
  default = true
}
variable "custom_cpu" {
  type    = bool
  default = false
}
variable "custom_max_memory" {
  type    = bool
  default = true
}
variable "max_cpu" {
  type    = number
  default = 0
}
variable "max_disks" {
  type    = number
  default = 2
}
variable "memory_size_type" {
  type    = string
  default = "mb"
}
variable "price_set_ids" {
  type    = list(number)
  default = [1]
}
variable "sort_order" {
  type    = number
  default = 10000
}
variable "storage_size_type" {
  type    = string
  default = "gb"
}
variable "config_ranges_min_memory" {
  type    = number
  default = 1048576
}
variable "config_ranges_max_memory" {
  type    = number
  default = 2097152
}
variable "config_ranges_min_cores" {
  type    = number
  default = 1
}
variable "config_ranges_max_cores" {
  type    = number
  default = 2
}
variable "config_ranges_min_sockets" {
  type    = number
  default = 1
}
variable "config_ranges_max_sockets" {
  type    = number
  default = 10
}
variable "config_ranges_min_cores_per_socket" {
  type    = number
  default = 1
}
variable "config_ranges_max_cores_per_socket" {
  type    = number
  default = 10
}

resource "hpe_morpheus_service_plan" "test" {
	name                   = var.name
	code                   = var.code
	max_memory             = var.max_memory
	max_storage            = var.max_storage
	add_volumes            = var.add_volumes
	cores_per_socket       = var.cores_per_socket
	custom_cores           = var.custom_cores
	custom_cpu             = var.custom_cpu
	custom_max_memory      = var.custom_max_memory
	custom_max_storage     = var.custom_max_storage
	description            = var.description
	max_cores              = var.max_cores
	max_cpu                = var.max_cpu
	max_disks              = var.max_disks
	memory_size_type       = var.memory_size_type
	price_set_ids          = var.price_set_ids
	provision_type_code    = var.provision_type_code
	sort_order             = var.sort_order
	storage_size_type      = var.storage_size_type

	config_ranges = {
		min_storage           = var.custom_max_storage ? var.config_ranges_min_storage : null
		max_storage           = var.custom_max_storage ? var.config_ranges_max_storage : null
		min_memory            = var.custom_max_memory ? var.config_ranges_min_memory : null
		max_memory            = var.custom_max_memory ? var.config_ranges_max_memory : null
		min_cores             = var.custom_cores ? var.config_ranges_min_cores : null
		max_cores             = var.custom_cores ? var.config_ranges_max_cores : null
		min_sockets           = var.config_ranges_min_sockets
		max_sockets           = var.config_ranges_max_sockets
		min_cores_per_socket  = var.config_ranges_min_cores_per_socket
		max_cores_per_socket  = var.config_ranges_max_cores_per_socket
	}
}
`

	baseConfigVars := config.Variables{
		"name": config.StringVariable(uniqueName),
	}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "name", uniqueName),
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "code", "test-serviceplan-code"),
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
		resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "provision_type_code", "kvm"),
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
			}, { // name change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName + "-updated"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // code change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					"code": config.StringVariable("updated-serviceplan-code"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // description change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"description": config.StringVariable("Updated description"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_memory change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(uniqueName),
					"max_memory": config.IntegerVariable(8589934592), // 8GB
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_storage change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"max_storage": config.IntegerVariable(10737418240), // 10GB
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // cores_per_socket change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":             config.StringVariable(uniqueName),
					"cores_per_socket": config.IntegerVariable(2),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // custom_cores toggle
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":         config.StringVariable(uniqueName),
					"custom_cores": config.BoolVariable(false),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // custom_cpu toggle
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(uniqueName),
					"custom_cpu": config.BoolVariable(true),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // custom_max_memory toggle
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":              config.StringVariable(uniqueName),
					"custom_max_memory": config.BoolVariable(false),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // custom_max_storage toggle
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":               config.StringVariable(uniqueName),
					"custom_max_storage": config.BoolVariable(false),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // add_volumes toggle
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"add_volumes": config.BoolVariable(false),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_cores change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":      config.StringVariable(uniqueName),
					"max_cores": config.IntegerVariable(4),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_cpu change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":    config.StringVariable(uniqueName),
					"max_cpu": config.IntegerVariable(100),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // max_disks change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":      config.StringVariable(uniqueName),
					"max_disks": config.IntegerVariable(5),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // memory_size_type change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":             config.StringVariable(uniqueName),
					"memory_size_type": config.StringVariable("gb"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // sort_order change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(uniqueName),
					"sort_order": config.IntegerVariable(20000),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // storage_size_type change
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":              config.StringVariable(uniqueName),
					"storage_size_type": config.StringVariable("mb"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // provision_type_code change from kvm to arm
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                config.StringVariable(uniqueName),
					"provision_type_code": config.StringVariable("arm"),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // config_ranges changes
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                               config.StringVariable(uniqueName),
					"config_ranges_min_storage":          config.IntegerVariable(2),
					"config_ranges_max_storage":          config.IntegerVariable(4),
					"config_ranges_min_memory":           config.IntegerVariable(2097152),
					"config_ranges_max_memory":           config.IntegerVariable(4194304),
					"config_ranges_min_cores":            config.IntegerVariable(2),
					"config_ranges_max_cores":            config.IntegerVariable(4),
					"config_ranges_min_sockets":          config.IntegerVariable(2),
					"config_ranges_max_sockets":          config.IntegerVariable(20),
					"config_ranges_min_cores_per_socket": config.IntegerVariable(2),
					"config_ranges_max_cores_per_socket": config.IntegerVariable(20),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}, { // comprehensive apply with multiple updates
				Config: baseConfigText,
				ConfigVariables: config.Variables{
					"name":                               config.StringVariable(uniqueName),
					"description":                        config.StringVariable("Comprehensive update test"),
					"max_memory":                         config.IntegerVariable(8589934592),
					"max_storage":                        config.IntegerVariable(10737418240),
					"cores_per_socket":                   config.IntegerVariable(2),
					"custom_cpu":                         config.BoolVariable(true),
					"add_volumes":                        config.BoolVariable(false),
					"max_cores":                          config.IntegerVariable(4),
					"max_cpu":                            config.IntegerVariable(100),
					"max_disks":                          config.IntegerVariable(5),
					"sort_order":                         config.IntegerVariable(20000),
					"provision_type_code":                config.StringVariable("arm"),
					"config_ranges_min_storage":          config.IntegerVariable(2),
					"config_ranges_max_storage":          config.IntegerVariable(4),
					"config_ranges_min_memory":           config.IntegerVariable(2097152),
					"config_ranges_max_memory":           config.IntegerVariable(4194304),
					"config_ranges_min_cores":            config.IntegerVariable(2),
					"config_ranges_max_cores":            config.IntegerVariable(4),
					"config_ranges_min_sockets":          config.IntegerVariable(2),
					"config_ranges_max_sockets":          config.IntegerVariable(20),
					"config_ranges_min_cores_per_socket": config.IntegerVariable(2),
					"config_ranges_max_cores_per_socket": config.IntegerVariable(20),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "code", "test-serviceplan-code"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "description", "Comprehensive update test"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_memory", "8589934592"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_storage", "10737418240"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "cores_per_socket", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "custom_cpu", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "add_volumes", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_cores", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_cpu", "100"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "max_disks", "5"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "sort_order", "20000"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "provision_type_code", "arm"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_storage", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_storage", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_memory", "2097152"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_memory", "4194304"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_cores", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_cores", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_sockets", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_sockets", "20"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.min_cores_per_socket", "2"),
					resource.TestCheckResourceAttr("hpe_morpheus_service_plan.test", "config_ranges.max_cores_per_socket", "20"),
				),
			},
		},
	})
}
