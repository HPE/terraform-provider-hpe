package policy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// TestAccMorpheusPolicyResourceUpdateAllAttrsOk tests updating all mutable attributes and detecting changes
func TestAccMorpheusPolicyResourceUpdateAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	uniqueName := acctest.RandomWithPrefix(t.Name())

	baseConfigText := providerConfig + `
variable "name" { type = string }
variable "description" {
  type    = string
  default = "Test policy description"
}
variable "enabled" {
  type    = bool
  default = true
}
variable "each_user" {
  type    = bool
  default = false
}
variable "max_memory" {
  type    = number
  default = 4294967296
}
variable "max_containers" {
  type    = number
  default = 10
}
variable "max_storage" {
  type    = number
  default = 10737418240
}
variable "max_cores" {
  type    = number
  default = 4
}
variable "budget_max_price" {
  type    = number
  default = 1000.0
}
variable "budget_currency" {
  type    = string
  default = "USD"
}
variable "budget_unit" {
  type    = string
  default = "month"
}

resource "hpe_morpheus_policy" "test" {
	name        = var.name
	description = var.description
	enabled     = var.enabled
	each_user   = var.each_user

	# Required empty config field
	config = {}

	config_max_memory = {
		max_memory = {
			anyof1 = var.max_memory
		}
	}

	config_max_containers = {
		max_containers = var.max_containers
	}

	config_max_storage = {
		max_storage = var.max_storage
	}

	config_max_cores = {
		max_cores = var.max_cores
	}

	config_budget = {
		max_price          = var.budget_max_price
		max_price_currency = var.budget_currency
		max_price_unit     = var.budget_unit
	}
}
`

	baseConfigVars := config.Variables{
		"name": config.StringVariable(uniqueName),
	}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Test policy description"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "each_user", "false"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_containers.max_containers", "10"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_storage.max_storage", "10737418240"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_cores.max_cores", "4"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price", "1000"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price_currency", "USD"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price_unit", "month"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	// Test updates to various fields
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial creation
			{
				ConfigVariables:    baseConfigVars,
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			// Update description
			{
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"description": config.StringVariable("Updated policy description"),
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "true"),
				),
				PlanOnly: false,
			},
			// Update enabled status
			{
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"description": config.StringVariable("Updated policy description"),
					"enabled":     config.BoolVariable(false),
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "false"),
				),
				PlanOnly: false,
			},
			// Update each_user flag
			{
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"description": config.StringVariable("Updated policy description"),
					"enabled":     config.BoolVariable(false),
					"each_user":   config.BoolVariable(true),
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "each_user", "true"),
				),
				PlanOnly: false,
			},
			// Update max memory
			{
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(uniqueName),
					"description": config.StringVariable("Updated policy description"),
					"enabled":     config.BoolVariable(false),
					"each_user":   config.BoolVariable(true),
					"max_memory":  config.IntegerVariable(8589934592), // 8GB
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "each_user", "true"),
				),
				PlanOnly: false,
			},
			// Update max containers
			{
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(uniqueName),
					"description":    config.StringVariable("Updated policy description"),
					"enabled":        config.BoolVariable(false),
					"each_user":      config.BoolVariable(true),
					"max_memory":     config.IntegerVariable(8589934592), // 8GB
					"max_containers": config.IntegerVariable(20),
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "each_user", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_containers.max_containers", "20"),
				),
				PlanOnly: false,
			},
			// Update budget settings
			{
				ConfigVariables: config.Variables{
					"name":             config.StringVariable(uniqueName),
					"description":      config.StringVariable("Updated policy description"),
					"enabled":          config.BoolVariable(false),
					"each_user":        config.BoolVariable(true),
					"max_memory":       config.IntegerVariable(8589934592), // 8GB
					"max_containers":   config.IntegerVariable(20),
					"budget_max_price": config.FloatVariable(2000.0),
					"budget_currency":  config.StringVariable("EUR"),
					"budget_unit":      config.StringVariable("week"),
				},
				Config:             baseConfigText,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Updated policy description"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "each_user", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_containers.max_containers", "20"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price", "2000"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price_currency", "EUR"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_budget.max_price_unit", "week"),
				),
				PlanOnly: false,
			},
			// Final plan check
			{
				ConfigVariables: config.Variables{
					"name":             config.StringVariable(uniqueName),
					"description":      config.StringVariable("Updated policy description"),
					"enabled":          config.BoolVariable(false),
					"each_user":        config.BoolVariable(true),
					"max_memory":       config.IntegerVariable(8589934592), // 8GB
					"max_containers":   config.IntegerVariable(20),
					"budget_max_price": config.FloatVariable(2000.0),
					"budget_currency":  config.StringVariable("EUR"),
					"budget_unit":      config.StringVariable("week"),
				},
				Config:             baseConfigText,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // Ensures the plan is empty, validating no drift
			},
		},
	})
}

// TestAccMorpheusPolicyConfigUpdateOk tests updating specific config types
func TestAccMorpheusPolicyConfigUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Test updating between different config types
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with backup creation policy
			{
				Config: providerConfig + `
resource "hpe_morpheus_policy" "test_config" {
  name        = "` + uniqueName + `"
  description = "Test config update"
  enabled     = true

  # Required empty config field
  config = {}

  config_backup_creation = {
    create_backup      = true
    create_backup_type = "snapshot"
  }
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_backup_creation.create_backup", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_backup_creation.create_backup_type", "snapshot"),
				),
			},
			// Update to backup creation with different settings
			{
				Config: providerConfig + `
resource "hpe_morpheus_policy" "test_config" {
  name        = "` + uniqueName + `"
  description = "Test config update - modified"
  enabled     = true

  # Required empty config field
  config = {}

  config_backup_creation = {
    create_backup      = false
    create_backup_type = "backup"
  }
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "description", "Test config update - modified"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_backup_creation.create_backup", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_backup_creation.create_backup_type", "backup"),
				),
			},
			// Update tags configuration
			{
				Config: providerConfig + `
resource "hpe_morpheus_policy" "test_config" {
  name        = "` + uniqueName + `"
  description = "Test config update - tags"
  enabled     = true

  # Required empty config field
  config = {}

  config_tags = {
    key    = "Project"
    value  = "Alpha,Beta,Gamma"
    strict = true
  }
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "name", uniqueName),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "description", "Test config update - tags"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_tags.key", "Project"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_tags.value", "Alpha,Beta,Gamma"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test_config", "config_tags.strict", "true"),
				),
			},
			// Plan check
			{
				Config: providerConfig + `
resource "hpe_morpheus_policy" "test_config" {
  name        = "` + uniqueName + `"
  description = "Test config update - tags"
  enabled     = true

  # Required empty config field
  config = {}

  config_tags = {
    key    = "Project"
    value  = "Alpha,Beta,Gamma"
    strict = true
  }
}
				`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
