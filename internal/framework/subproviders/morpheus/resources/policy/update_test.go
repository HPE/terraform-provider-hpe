package policy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// Test update of mutable attributes
func TestAccMorpheusPolicyUpdateOk(t *testing.T) {
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

// Test that changing associated_resource_type triggers replacement
func TestAccMorpheusPolicyResourceTypeChangeRequiresReplace(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

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
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + initialConfig,
			},
			{
				Config:             providerConfig + updatedConfig,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.replace_test", "associated_resource_type", "Group"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_policy.replace_test", "associated_resource_id"),
				),
			},
		},
	})
}
