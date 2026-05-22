// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// Test update of mutable attributes
func TestAccMorpheusPolicyUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	nameUpdated := name + "-updated"
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	initialConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "update_test" {
  name = "` + name + `"
  description = "Initial description"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	updatedConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "update_test" {
  name = "` + nameUpdated + `"
  description = "Updated description"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = false
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 16
  }
  
  tenants = [1]
}
`

	initialChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "description", "Initial description"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "enabled", "true"),
	}

	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "name", nameUpdated),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "description", "Updated description"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.update_test", "enabled", "false"),
	}

	initialCheckFn := resource.ComposeAggregateTestCheckFunc(initialChecks...)
	updatedCheckFn := resource.ComposeAggregateTestCheckFunc(updatedChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + initialConfig,
				ExpectNonEmptyPlan: false,
				Check:              initialCheckFn,
			},
			{
				Config:             providerConfig + updatedConfig,
				ExpectNonEmptyPlan: false,
				Check:              updatedCheckFn,
			},
		},
	})
}

// Test that changing associated_resource_id triggers replacement
func TestAccMorpheusPolicyAssociatedResourceIdChangeRequiresReplace(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	const resourceName = "hpe_morpheus_policy.replace_test"

	providerConfig := testhelpers.ProviderBlock()
	groupName := acctest.RandomWithPrefix(t.Name())
	policyName := acctest.RandomWithPrefix(t.Name())

	initialConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	updatedConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_group" "test2" {
  name = "` + groupName + `-2"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test2.id
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	var initialResourceID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with initial associated_resource_id
				Config: providerConfig + initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_resource_id"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxMemory"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Store the initial ID for comparison
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}
						initialResourceID = rs.Primary.ID

						return nil
					},
				),
			},
			{
				// Step 2: Change associated_resource_id and verify resource is replaced (new ID)
				Config: providerConfig + updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_resource_id"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxMemory"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Verify the ID changed (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_policy.replace_test")
						}
						currentResourceID := rs.Primary.ID
						if currentResourceID == initialResourceID {
							return fmt.Errorf("Expected resource ID to change due to associated_resource_id change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceID)
						}

						return nil
					},
				),
			},
		},
	})
}

// Test that changing associated_resource_type triggers replacement
func TestAccMorpheusPolicyAssociatedResourceTypeChangeRequiresReplace(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	const resourceName = "hpe_morpheus_policy.replace_test"

	providerConfig := testhelpers.ProviderBlock()
	groupName := acctest.RandomWithPrefix(t.Name())
	policyName := acctest.RandomWithPrefix(t.Name())

	initialConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	updatedConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Global"
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	var initialResourceID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with Group associated_resource_type
				Config: providerConfig + initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_resource_id"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxMemory"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Store the initial ID for comparison
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}
						initialResourceID = rs.Primary.ID

						return nil
					},
				),
			},
			{
				// Step 2: Change associated_resource_type to Global and verify resource is replaced
				Config: providerConfig + updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Global"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxMemory"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Verify the ID changed (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_policy.replace_test")
						}
						currentResourceID := rs.Primary.ID
						if currentResourceID == initialResourceID {
							return fmt.Errorf("Expected resource ID to change due to associated_resource_type change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceID)
						}

						return nil
					},
				),
			},
		},
	})
}

// Test that changing policy_type.code triggers replacement
func TestAccMorpheusPolicyTypeCodeChangeRequiresReplace(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	const resourceName = "hpe_morpheus_policy.replace_test"

	providerConfig := testhelpers.ProviderBlock()
	groupName := acctest.RandomWithPrefix(t.Name())
	policyName := acctest.RandomWithPrefix(t.Name())

	initialConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
  
  tenants = [1]
}
`

	updatedConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "replace_test" {
  name = "` + policyName + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "maxStorage"
  }
  
  config = {
    maxStorage = 10
  }
  
  tenants = [1]
}
`

	var initialResourceID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with maxMemory policy type
				Config: providerConfig + initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_resource_id"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxMemory"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Store the initial ID for comparison
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}
						initialResourceID = rs.Primary.ID

						return nil
					},
				),
			},
			{
				// Step 2: Change policy_type.code to maxStorage and verify resource is replaced
				Config: providerConfig + updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_resource_id"),
					resource.TestCheckResourceAttr(resourceName, "policy_type.code", "maxStorage"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Verify the ID changed (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_policy.replace_test")
						}
						currentResourceID := rs.Primary.ID
						if currentResourceID == initialResourceID {
							return fmt.Errorf("Expected resource ID to change due to policy_type.code change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceID)
						}

						return nil
					},
				),
			},
		},
	})
}
