// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../../cmd/render example.tf.tmpl Name "ExamplePolicy" PolicyTypeCode "maxMemory" ConfigMaxMemoryMaxMemory "4294967296"

package policy_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
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

func TestAccMorpheusPolicyRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "example_required" {
  name    = "` + name + `"
  enabled = false
  ref_id  = 9969
  
  ref_type = {
    oneof0 = "User"
  }
  
  config = {
    max_memory_policy_type_configuration = {
      max_memory = {
        anyof1 = 4294967296  # 4GB in bytes
      }
    }
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_required",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_required",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_required",
			"ref_id",
			"9969",
		),
		// Check that config.max_memory_policy_type_configuration is set
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_policy.example_required",
			"config.max_memory_policy_type_configuration.max_memory.anyof1",
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
				ResourceName:      "hpe_morpheus_policy.example_required",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyBackupCreationOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "example_backup" {
  name        = "` + name + `"
  description = "Backup creation policy for testing"
  enabled     = false
  ref_id      = 9969
  
  ref_type = {
    oneof0 = "User"
  }
  
  config = {
    backup_creation_policy_type_configuration = {
      create_backup      = true
      create_backup_type = "snapshot"
    }
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"description",
			"Backup creation policy for testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"ref_id",
			"9969",
		),
		// Backup creation config
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"config.backup_creation_policy_type_configuration.create_backup",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_backup",
			"config.backup_creation_policy_type_configuration.create_backup_type",
			"snapshot",
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
				ResourceName:      "hpe_morpheus_policy.example_backup",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyBudgetOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "example_budget" {
  name        = "` + name + `"
  description = "Budget policy for testing"
  enabled     = false
  ref_id      = 9969
  
  ref_type = {
    oneof0 = "User"
  }
  
  config = {
    budget_policy_type_configuration = {
      max_price          = 1000.50
      max_price_currency = "USD"
      max_price_unit     = "month"
    }
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"description",
			"Budget policy for testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"ref_id",
			"9969",
		),
		// Budget config
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"config.budget_policy_type_configuration.max_price",
			"1000.5",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"config.budget_policy_type_configuration.max_price_currency",
			"USD",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_budget",
			"config.budget_policy_type_configuration.max_price_unit",
			"month",
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
				ResourceName:      "hpe_morpheus_policy.example_budget",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyInstanceNameOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "example_instance_name" {
  name        = "` + name + `"
  description = "Instance naming policy"
  enabled     = false
  ref_id      = 9969
  
  ref_type = {
    oneof0 = "User"
  }
  
  config = {
    instance_name_policy_type_configuration = {
      naming_type    = "pattern"
      naming_pattern = "vm-$${sequence}-$${cloudCode}"
    }
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"description",
			"Instance naming policy",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"ref_id",
			"9969",
		),
		// Instance name config
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"config.instance_name_policy_type_configuration.naming_type",
			"pattern",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_instance_name",
			"config.instance_name_policy_type_configuration.naming_pattern",
			"vm-${sequence}-${cloudCode}",
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
				ResourceName:      "hpe_morpheus_policy.example_instance_name",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyTagsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "example_tags" {
  name        = "` + name + `"
  description = "Tags policy for testing"
  enabled     = false
  ref_id      = 9969
  
  ref_type = {
    oneof0 = "User"
  }
  
  config = {
    tags_policy_type_configuration = {
      key    = "Environment"
      value  = "Production,Development,Testing"
      strict = true
    }
  }
}
	`
	checks := []resource.TestCheckFunc{
		// Top-level attributes
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"description",
			"Tags policy for testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"ref_id",
			"9969",
		),
		// Tags config
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"config.tags_policy_type_configuration.key",
			"Environment",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"config.tags_policy_type_configuration.value",
			"Production,Development,Testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_tags",
			"config.tags_policy_type_configuration.strict",
			"true",
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
				ResourceName:      "hpe_morpheus_policy.example_tags",
				Check:             checkFn,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusPolicyExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"PolicyTypeCode", "maxMemory",
		"ConfigMaxMemoryMaxMemory", "4294967296")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.example_policy",
			"name",
			name,
		),
		// Check that some config is set
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_policy.example_policy",
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
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_policy.example_policy",
				Check:             checkFn,
			},
		},
	})
}
