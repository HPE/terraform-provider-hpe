// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package role_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Test for updating all attributes of a user role, including permissions.
// Creates all permission dependencies (i.e. permissions for resources using ID)
// as part of the test.
// TODO: Add VDI pool permissions checks once they are fixed in the OpenAPI spec and SDK.

func TestAccMorpheusRoleUserUpdateAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlockMixed()

	name := acctest.RandomWithPrefix(t.Name())
	nameUpdated := name + "Updated"

	dependencyResourceConfig := `
resource "hpe_morpheus_group" "testacc_group" {
  name = "` + name + `"
}

resource "morpheus_terraform_app_blueprint" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
  spec_template_ids = [] # not a required field, but if we don't include it this will be computed
}

resource "morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "morpheus_groovy_script_task" "testacc_task" {
  name = "` + name + `"
  source_type = "local"
}

resource "morpheus_operational_workflow" "testacc_workflow" {
  name = "` + name + `"
}
`
	// Initial configuration
	initialResourceConfig := `
data "hpe_morpheus_group" "testacc_group" {
  name = hpe_morpheus_group.testacc_group.name
}

data "morpheus_blueprint" "testacc_blueprint" {
  name = morpheus_terraform_app_blueprint.testacc_blueprint.name
}

data "morpheus_instance_type" "testacc_instance_type" {
  name = morpheus_instance_type.testacc_instance_type.name
}

data "morpheus_task" "testacc_task" {
  name = morpheus_groovy_script_task.testacc_task.name
}

data "morpheus_workflow" "testacc_workflow" {
  name = morpheus_operational_workflow.testacc_workflow.name
}

resource "hpe_morpheus_role" "update_test" {
  name = "` + name + `"
  description = "Initial role description"
  landing_url = "https://initial.example.com"
  multitenant = false
  multitenant_locked = false
  role_type = "user"
  permissions = {
	default_group_access = "none"
	default_instance_type_access = "none"
	default_blueprint_access = "none"
	default_catalog_item_type_access = "none"
	default_persona_access = "none"
	default_report_type_access = "none"
	default_task_access = "none"
	default_vdi_pool_access = "none"
	default_workflow_access = "none"
	feature_permissions = [
	  {
		code   = "integrations-ansible"
		access = "none"
	  }
	]
	# test specific permission types
	group_permissions = [
	  {
		id     = data.hpe_morpheus_group.testacc_group.id
		access = "none"
	  }
	]
	blueprint_permissions = [
	  {
		id     = data.morpheus_blueprint.testacc_blueprint.id
		access = "none"
	  }
	]
	instance_type_permissions = [
	  {
		id     = data.morpheus_instance_type.testacc_instance_type.id
		access = "none"
	  }
	]
	persona_permissions = [
	  {
		code   = "standard"
		access = "none"
	  }
	]
	report_type_permissions = [
	  {
		code   = "appCost"
		access = "none"
	  }
	]
	task_permissions = [
	  {
		id     = data.morpheus_task.testacc_task.id
		access = "none"
	  }
	]
	workflow_permissions = [
	  {
		id     = data.morpheus_workflow.testacc_workflow.id
		access = "none"
	  }
	]
  }
}
`

	// Updated configuration
	updatedResourceConfig := `
data "hpe_morpheus_group" "testacc_group" {
  name = hpe_morpheus_group.testacc_group.name
}

data "morpheus_blueprint" "testacc_blueprint" {
  name = morpheus_terraform_app_blueprint.testacc_blueprint.name
}

data "morpheus_instance_type" "testacc_instance_type" {
  name = morpheus_instance_type.testacc_instance_type.name
}

data "morpheus_task" "testacc_task" {
  name = morpheus_groovy_script_task.testacc_task.name
}

data "morpheus_workflow" "testacc_workflow" {
  name = morpheus_operational_workflow.testacc_workflow.name
}

resource "hpe_morpheus_role" "update_test" {
  name = "` + nameUpdated + `"
  description = "Updated role description"
  landing_url = "https://updated.example.com"
  multitenant = true
  multitenant_locked = true
  role_type = "user"
  permissions = {
	default_group_access = "full"
	default_instance_type_access = "full"
	default_blueprint_access = "full"
	default_catalog_item_type_access = "full"
	default_persona_access = "full"
	default_report_type_access = "full"
	default_task_access = "full"
	default_vdi_pool_access = "full"
	default_workflow_access = "full"
	feature_permissions = [
	  {
		code   = "integrations-ansible"
		access = "full"
	  }
	]
	group_permissions = [
	  {
		id     = data.hpe_morpheus_group.testacc_group.id
		access = "full"
	  }
	]
	blueprint_permissions = [
	  {
		id     = data.morpheus_blueprint.testacc_blueprint.id
		access = "full"
	  }
	]
	instance_type_permissions = [
	  {
		id     = data.morpheus_instance_type.testacc_instance_type.id
		access = "full"
	  }
	]
	persona_permissions = [
	  {
		code   = "standard"
		access = "full"
	  }
	]
	report_type_permissions = [
	  {
		code   = "appCost"
		access = "full"
	  }
	]
	task_permissions = [
	  {
		id     = data.morpheus_task.testacc_task.id
		access = "full"
	  }
	]
	workflow_permissions = [
	  {
		id     = data.morpheus_workflow.testacc_workflow.id
		access = "full"
	  }
	]
  }
}
`

	removedPermissionsResourceConfig := `
resource "hpe_morpheus_role" "update_test" {
  name = "` + nameUpdated + `"
  description = "Updated role description"
  landing_url = "https://updated.example.com"
  multitenant = true
  multitenant_locked = true
  role_type = "user"
}
`
	// Initial checks
	initialChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Initial role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://initial.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"user",
		),
		// Default access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_group_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_instance_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_blueprint_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_catalog_item_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_persona_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_report_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_task_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_vdi_pool_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_workflow_access",
			"none",
		),
		// Feature permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.code",
			"integrations-ansible",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.access",
			"none",
		),
		// Group permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.0.id",
			"data.hpe_morpheus_group.testacc_group",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.0.access",
			"none",
		),
		// Blueprint permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.id",
			"data.morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.access",
			"none",
		),
		// Instance type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.id",
			"data.morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.access",
			"none",
		),
		// Persona permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.access",
			"none",
		),
		// Report type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.code",
			"appCost",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.access",
			"none",
		),
		// Task permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.id",
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.access",
			"none",
		),
		// Workflow permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.access",
			"none",
		),
	}

	// Updated checks
	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			nameUpdated,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Updated role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://updated.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"user",
		),
		// Default access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_group_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_catalog_item_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_persona_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_report_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_vdi_pool_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_workflow_access",
			"full",
		),
		// Feature permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.code",
			"integrations-ansible",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.access",
			"full",
		),
		// Group permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.0.id",
			"data.hpe_morpheus_group.testacc_group",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.group_permissions.0.access",
			"full",
		),
		// Blueprint permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.id",
			"data.morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.access",
			"full",
		),
		// Instance type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.id",
			"data.morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.access",
			"full",
		),
		// Persona permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.access",
			"full",
		),
		// Report type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.code",
			"appCost",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.access",
			"full",
		),
		// Task permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.id",
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.access",
			"full",
		),
		// Workflow permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.access",
			"full",
		),
	}

	// Removed permissions checks
	removedPermissionsChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			nameUpdated,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Updated role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://updated.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"user",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions",
		),
	}

	initialCheckFn := resource.ComposeAggregateTestCheckFunc(initialChecks...)
	updatedCheckFn := resource.ComposeAggregateTestCheckFunc(updatedChecks...)
	removedPermissionsCheckFn := resource.ComposeAggregateTestCheckFunc(removedPermissionsChecks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create dependencies
			{
				Config:             providerConfig + dependencyResourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
			},
			// Create initial role
			{
				Config: providerConfig + dependencyResourceConfig + initialResourceConfig,
				// blueprint resource in the old provider has a computed field
				// which doesn't use state for unknown, so plan will not be empty
				ExpectNonEmptyPlan: false,
				Check:              initialCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-apply plan detects no changes for dependencyResourceConfig + initial config
			{
				Config:             providerConfig + dependencyResourceConfig + initialResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              initialCheckFn,
				PlanOnly:           true,
			},
			// Update the role
			{
				Config:             providerConfig + dependencyResourceConfig + updatedResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              updatedCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-update plan detects no changes
			{
				Config:             providerConfig + dependencyResourceConfig + updatedResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              updatedCheckFn,
				PlanOnly:           true,
			},
			// Test import state verification with updated values
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Feature permissions on import will be the full set of computed permissions,
					// not just those we set.
					"permissions.feature_permissions",
					"permissions.default_cloud_access",
				},
				ResourceName: "hpe_morpheus_role.update_test",
				Check:        updatedCheckFn,
			},
			// Finally, test that nothing breaks if we remove the permissions block
			{
				Config:             providerConfig + removedPermissionsResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              removedPermissionsCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-update plan detects no changes
			{
				Config:             providerConfig + removedPermissionsResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              removedPermissionsCheckFn,
				PlanOnly:           true,
			},
		},
	})
}

// Test for updating all attributes of a tenant role, including permissions.
// Creates all permission dependencies (i.e. permissions for resources using ID)
// as part of the test.
// TODO: Add VDI pool permissions checks once they are fixed in the OpenAPI spec and SDK.

func TestAccMorpheusRoleTenantUpdateAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlockMixed()

	name := acctest.RandomWithPrefix(t.Name())
	nameUpdated := name + "Updated"

	dependencyResourceConfig := `
resource "morpheus_standard_cloud" "testacc_cloud" {
  name = "` + name + `"
  code = "standard"
  tenant_id = 1
  visibility = "public"
}

resource "morpheus_terraform_app_blueprint" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
  spec_template_ids = [] # not a required field, but if we don't include it this will be computed
}

resource "morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "morpheus_groovy_script_task" "testacc_task" {
  name = "` + name + `"
  source_type = "local"
}

resource "morpheus_operational_workflow" "testacc_workflow" {
  name = "` + name + `"
}
`
	// Initial configuration
	initialResourceConfig := `
data "morpheus_cloud" "testacc_cloud" {
  name = morpheus_standard_cloud.testacc_cloud.name
}

data "morpheus_blueprint" "testacc_blueprint" {
  name = morpheus_terraform_app_blueprint.testacc_blueprint.name
}

data "morpheus_instance_type" "testacc_instance_type" {
  name = morpheus_instance_type.testacc_instance_type.name
}

data "morpheus_task" "testacc_task" {
  name = morpheus_groovy_script_task.testacc_task.name
}

data "morpheus_workflow" "testacc_workflow" {
  name = morpheus_operational_workflow.testacc_workflow.name
}

resource "hpe_morpheus_role" "update_test" {
  name = "` + name + `"
  description = "Initial role description"
  landing_url = "https://initial.example.com"
  role_type = "tenant"
  permissions = {
	default_cloud_access = "none"
	default_instance_type_access = "none"
	default_blueprint_access = "none"
	default_catalog_item_type_access = "none"
	default_persona_access = "none"
	default_report_type_access = "none"
	default_task_access = "none"
	default_vdi_pool_access = "none"
	default_workflow_access = "none"
	feature_permissions = [
	  {
		code   = "integrations-ansible"
		access = "none"
	  }
	]
	# test specific permission types
	cloud_permissions = [
	  {
		id     = data.morpheus_cloud.testacc_cloud.id
		access = "none"
	  }
	]
	blueprint_permissions = [
	  {
		id     = data.morpheus_blueprint.testacc_blueprint.id
		access = "none"
	  }
	]
	instance_type_permissions = [
	  {
		id     = data.morpheus_instance_type.testacc_instance_type.id
		access = "none"
	  }
	]
	persona_permissions = [
	  {
		code   = "standard"
		access = "none"
	  }
	]
	report_type_permissions = [
	  {
		code   = "appCost"
		access = "none"
	  }
	]
	task_permissions = [
	  {
		id     = data.morpheus_task.testacc_task.id
		access = "none"
	  }
	]
	workflow_permissions = [
	  {
		id     = data.morpheus_workflow.testacc_workflow.id
		access = "none"
	  }
	]
  }
}
`

	// Updated configuration
	updatedResourceConfig := `
data "morpheus_cloud" "testacc_cloud" {
  name = morpheus_standard_cloud.testacc_cloud.name
}

data "morpheus_blueprint" "testacc_blueprint" {
  name = morpheus_terraform_app_blueprint.testacc_blueprint.name
}

data "morpheus_instance_type" "testacc_instance_type" {
  name = morpheus_instance_type.testacc_instance_type.name
}

data "morpheus_task" "testacc_task" {
  name = morpheus_groovy_script_task.testacc_task.name
}

data "morpheus_workflow" "testacc_workflow" {
  name = morpheus_operational_workflow.testacc_workflow.name
}

resource "hpe_morpheus_role" "update_test" {
  name = "` + nameUpdated + `"
  description = "Updated role description"
  landing_url = "https://updated.example.com"
  role_type = "tenant"
  permissions = {
	default_cloud_access = "full"
	default_instance_type_access = "full"
	default_blueprint_access = "full"
	default_catalog_item_type_access = "full"
	default_persona_access = "full"
	default_report_type_access = "full"
	default_task_access = "full"
	default_vdi_pool_access = "full"
	default_workflow_access = "full"
	feature_permissions = [
	  {
		code   = "integrations-ansible"
		access = "full"
	  }
	]
	cloud_permissions = [
	  {
		id     = data.morpheus_cloud.testacc_cloud.id
		access = "full"
	  }
	]
	blueprint_permissions = [
	  {
		id     = data.morpheus_blueprint.testacc_blueprint.id
		access = "full"
	  }
	]
	instance_type_permissions = [
	  {
		id     = data.morpheus_instance_type.testacc_instance_type.id
		access = "full"
	  }
	]
	persona_permissions = [
	  {
		code   = "standard"
		access = "full"
	  }
	]
	report_type_permissions = [
	  {
		code   = "appCost"
		access = "full"
	  }
	]
	task_permissions = [
	  {
		id     = data.morpheus_task.testacc_task.id
		access = "full"
	  }
	]
	workflow_permissions = [
	  {
		id     = data.morpheus_workflow.testacc_workflow.id
		access = "full"
	  }
	]
  }
}
`

	removedPermissionsResourceConfig := `
resource "hpe_morpheus_role" "update_test" {
  name = "` + nameUpdated + `"
  description = "Updated role description"
  landing_url = "https://updated.example.com"
  role_type = "tenant"
}
`
	// Initial checks
	initialChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Initial role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://initial.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"tenant",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"false",
		),
		// checks for fields not applicable to tenant roles
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_group_access",
		),
		// Default access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_cloud_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_instance_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_blueprint_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_catalog_item_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_persona_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_report_type_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_task_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_vdi_pool_access",
			"none",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_workflow_access",
			"none",
		),
		// Feature permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.code",
			"integrations-ansible",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.access",
			"none",
		),
		// Cloud permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.0.id",
			"data.morpheus_cloud.testacc_cloud",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.0.access",
			"none",
		),
		// Blueprint permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.id",
			"data.morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.access",
			"none",
		),
		// Instance type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.id",
			"data.morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.access",
			"none",
		),
		// Persona permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.access",
			"none",
		),
		// Report type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.code",
			"appCost",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.access",
			"none",
		),
		// Task permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.id",
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.access",
			"none",
		),
		// Workflow permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.access",
			"none",
		),
	}

	// Updated checks
	updatedChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			nameUpdated,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Updated role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://updated.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"tenant",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"false",
		),
		// checks for fields not applicable to tenant roles
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_group_access",
		),
		// Default access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_cloud_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_catalog_item_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_persona_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_report_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_vdi_pool_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.default_workflow_access",
			"full",
		),
		// Feature permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.code",
			"integrations-ansible",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.feature_permissions.0.access",
			"full",
		),
		// Cloud permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.0.id",
			"data.morpheus_cloud.testacc_cloud",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.cloud_permissions.0.access",
			"full",
		),
		// Blueprint permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.id",
			"data.morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.blueprint_permissions.0.access",
			"full",
		),
		// Instance type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.id",
			"data.morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.instance_type_permissions.0.access",
			"full",
		),
		// Persona permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.persona_permissions.0.access",
			"full",
		),
		// Report type permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.code",
			"appCost",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.report_type_permissions.0.access",
			"full",
		),
		// Task permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.id",
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.task_permissions.0.access",
			"full",
		),
		// Workflow permissions
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions.workflow_permissions.0.access",
			"full",
		),
	}

	// Removed permissions checks
	removedPermissionsChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"name",
			nameUpdated,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"description",
			"Updated role description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"landing_url",
			"https://updated.example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"role_type",
			"tenant",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.update_test",
			"multitenant_locked",
			"false",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.update_test",
			"permissions",
		),
	}

	initialCheckFn := resource.ComposeAggregateTestCheckFunc(initialChecks...)
	updatedCheckFn := resource.ComposeAggregateTestCheckFunc(updatedChecks...)
	removedPermissionsCheckFn := resource.ComposeAggregateTestCheckFunc(removedPermissionsChecks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create dependencies
			{
				Config:             providerConfig + dependencyResourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
			},
			// Create initial role
			{
				Config: providerConfig + dependencyResourceConfig + initialResourceConfig,
				// blueprint resource in the old provider has a computed field
				// which doesn't use state for unknown, so plan will not be empty
				ExpectNonEmptyPlan: false,
				Check:              initialCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-apply plan detects no changes for dependencyResourceConfig + initial config
			{
				Config:             providerConfig + dependencyResourceConfig + initialResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              initialCheckFn,
				PlanOnly:           true,
			},
			// Update the role
			{
				Config:             providerConfig + dependencyResourceConfig + updatedResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              updatedCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-update plan detects no changes
			{
				Config:             providerConfig + dependencyResourceConfig + updatedResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              updatedCheckFn,
				PlanOnly:           true,
			},
			// Test import state verification with updated values
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Feature permissions on import will be the full set of computed permissions,
					// not just those we set.
					"permissions.feature_permissions",
					"permissions.default_group_access",
				},
				ResourceName: "hpe_morpheus_role.update_test",
				Check:        updatedCheckFn,
			},
			// Finally, test that nothing breaks if we remove the permissions block
			{
				Config:             providerConfig + removedPermissionsResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              removedPermissionsCheckFn,
				PlanOnly:           false,
			},
			// Check that a post-update plan detects no changes
			{
				Config:             providerConfig + removedPermissionsResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              removedPermissionsCheckFn,
				PlanOnly:           true,
			},
		},
	})
}
