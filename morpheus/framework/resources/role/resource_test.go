// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package role_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Some notes about what we expect to happen with Permissions in acceptance test import testing:

// On import, if the permissions have been computed at create,
// then the import step will pass happily.
// If the permissions have been set by the user at create,
// then the import verification step will fail,
// because the API permissions being imported do not match the
// existing resource's permissions in state.

// Therefore, for any tests using user-set permissions,
// we skip the permissions import verification check.

// Check that we can create a user role with only required attributes specified
func TestAccMorpheusRoleUserRequiredAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "example_required" {
  name = "` + name + `"
}
`
	checks := []resource.TestCheckFunc{
		// required
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"name",
			name,
		),
		// checks for optional
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"description",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"default_persona_code",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"landing_url",
		),
		// checks for computed
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"multitenant_locked",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"role_type",
			"user",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true, // Check state post import
				ImportStateVerifyIgnore: []string{"permissions"},
				ResourceName:            "hpe_morpheus_role.example_required",
				Check:                   checkFn,
			},
		},
	})
}

// Check that we can create a tenant role with only required attributes specified
func TestAccMorpheusRoleTenantRequiredAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "example_required" {
  name = "` + name + `"
  role_type = "tenant"
}
`
	checks := []resource.TestCheckFunc{
		// required
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"name",
			name,
		),
		// checks for optional
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"description",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"default_persona_code",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"landing_url",
		),
		// check the role type
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"role_type",
			"tenant",
		),
		// checks for fields not applicable to tenant roles
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"multitenant_locked",
			"false",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true, // Check state post import
				ImportStateVerifyIgnore: []string{"permissions"},
				ResourceName:            "hpe_morpheus_role.example_required",
				Check:                   checkFn,
			},
		},
	})
}

// Check that we can create a role with all attributes specified
func TestAccMorpheusRoleAllAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.Ansible, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "example_all" {
  name = "` + name + `"
  description = "test"
  default_persona_code = "standard"
  landing_url = "https://test.com"
  multitenant = true
  multitenant_locked = true
  role_type = "user"
  permissions = {
	feature_permissions = [
	  {
		code   = "integrations-ansible"
		access = "full"
	  }
	]
	default_group_access = "full"
  }
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"description",
			"test",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"default_persona_code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"landing_url",
			"https://test.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"multitenant",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"multitenant_locked",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"role_type",
			"user",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"permissions.feature_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"permissions.feature_permissions.0.code",
			"integrations-ansible",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"permissions.feature_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"permissions.default_group_access",
			"full",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.default_catalog_item_type_access",
					"permissions.default_cloud_access",
					"permissions.default_instance_type_access",
					"permissions.default_persona_access",
					"permissions.default_report_type_access",
					"permissions.default_task_access",
					"permissions.default_workflow_access",
					"permissions.default_vdi_pool_access",
					"permissions.default_blueprint_access",
				},
				ResourceName: "hpe_morpheus_role.example_all",
				Check:        checkFn,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusRoleExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"Multitenant", "false",
		"Description", "a test of the example HCL config",
		"RoleType", "user")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"description",
			"a test of the example HCL config",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"role_type",
			"user",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
				//nolint:lll
				ImportStateVerifyIgnore: []string{"permissions"}, // ignore verification on computed permissions (import)
				ResourceName:            "hpe_morpheus_role.example",
				Check:                   checkFn,
			},
		},
	})
}

func TestAccMorpheusRolePermissionsDefaultAccessPermissionsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "default_access_permissions_ok" {
	name = "` + name + `"
	permissions = {
		default_group_access               = "full"
		default_instance_type_access      = "full"
		default_blueprint_access          = "full"
		default_catalog_item_type_access  = "full"
		default_persona_access            = "full"
		default_vdi_pool_access           = "full"
		default_report_type_access        = "full"
		default_task_access               = "full"
		default_workflow_access           = "full"
	}
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_group_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_catalog_item_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_persona_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_vdi_pool_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_report_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions.default_workflow_access",
			"full",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.default_cloud_access",
				},
				ResourceName: "hpe_morpheus_role.default_access_permissions_ok",
				Check:        checkFn,
			},
		},
	})
}

// test that we can create a user role with all possible permissions set using
// strongly-typed permissions
// we test all possible permissions EXCEPT VDI Pool.
// For now, the VDI pool section of the OpenAPI spec looks to be incorrect
// and needs to be updated so that we can create one using the generated SDK.
func TestAccMorpheusRoleAllPermissionsUserRoleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig := `
resource "hpe_morpheus_group" "testacc_group" {
  name = "` + name + `"
}

resource "hpe_morpheus_app_blueprint_terraform" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
}

resource "hpe_morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "hpe_morpheus_task_groovy_script" "testacc_task" {
  name = "` + name + `"
  source_type         = "local"
}

resource "hpe_morpheus_workflow_operational" "testacc_workflow" {
  name = "` + name + `"
}
`

	resourceConfig := `
data "hpe_morpheus_group" "testacc_group" {
  name = hpe_morpheus_group.testacc_group.name
}

data "hpe_morpheus_blueprint" "testacc_blueprint" {
  name = hpe_morpheus_app_blueprint_terraform.testacc_blueprint.name
}

data "hpe_morpheus_instance_type" "testacc_instance_type" {
  name = hpe_morpheus_instance_type.testacc_instance_type.name
}

data "hpe_morpheus_task" "testacc_task" {
  name = hpe_morpheus_task_groovy_script.testacc_task.name
}

data "hpe_morpheus_workflow" "testacc_workflow" {
  name = hpe_morpheus_workflow_operational.testacc_workflow.name
}

resource "hpe_morpheus_role" "testacc_role_all_permissions_user_role_ok" {
  name      = "` + name + `"
  role_type = "user"

  permissions = {
    feature_permissions = [
      {
        code   = "activity"
        access = "read"
      },
      {
        code   = "admin-accounts"
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
        id     = data.hpe_morpheus_blueprint.testacc_blueprint.id
        access = "full"
      }
    ]
    instance_type_permissions = [
      {
        id     = data.hpe_morpheus_instance_type.testacc_instance_type.id
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
        id     = data.hpe_morpheus_task.testacc_task.id
        access = "full"
      }
    ]
    workflow_permissions = [
      {
        id     = data.hpe_morpheus_workflow.testacc_workflow.id
        access = "full"
      }
    ]
    default_group_access             = "full"
    default_blueprint_access         = "full"
    default_catalog_item_type_access = "full"
    default_instance_type_access     = "full"
    default_persona_access           = "full"
    default_report_type_access       = "full"
    default_task_access              = "full"
    default_workflow_access          = "full"
    default_vdi_pool_access          = "full"
  }
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"role_type",
			"user",
		),
		// check the default permission access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_group_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_workflow_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.default_vdi_pool_access",
			"full",
		),
		// check the permissions for resources already existing in morpheus
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "activity",
				"access": "read",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "admin-accounts",
				"access": "full",
			},
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.persona_permissions.0.access",
			"full",
		),
		// check the permissions for the resources created with the legacy provider
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.group_permissions.0.id",
			"data.hpe_morpheus_group.testacc_group",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.group_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.blueprint_permissions.0.id",
			"data.hpe_morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.blueprint_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.instance_type_permissions.0.id",
			"data.hpe_morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.instance_type_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.task_permissions.0.id",
			"data.hpe_morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.workflow_permissions.0.id",
			"data.hpe_morpheus_workflow.testacc_workflow",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyResourceConfig,
				// one of the blueprints values will be computed
				// so this has to be set to `true`
				ExpectNonEmptyPlan: true,
				PlanOnly:           false,
			},
			{
				Config: providerConfig + dependencyResourceConfig + resourceConfig,
				// one of the blueprints values will be computed
				// so this has to be set to `true`
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.default_cloud_access",
				},
				ResourceName: "hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
				Check:        checkFn,
			},
		},
	})
}

// the difference between user and tenant role is that user roles can be assigned
// group permissions while tenant roles can be assigned cloud permissions
func TestAccMorpheusRoleTenantAllPermissionsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig := `
resource "hpe_morpheus_cloud" "testacc_cloud" {
  name = "` + name + `"
  tenant_id = 1
  visibility = "public" # cloud must be visible to the client for the zone permissions to be set
  config_hvm = {
    "enable_network_type_selection" = true
  }
}

resource "hpe_morpheus_app_blueprint_terraform" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
}

resource "hpe_morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "hpe_morpheus_task_groovy_script" "testacc_task" {
  name = "` + name + `"
  source_type         = "local"
}

resource "hpe_morpheus_workflow_operational" "testacc_workflow" {
  name = "` + name + `"
}
`

	resourceConfig := `
data "hpe_morpheus_cloud" "testacc_cloud" {
  name = hpe_morpheus_cloud.testacc_cloud.name
}

data "hpe_morpheus_blueprint" "testacc_blueprint" {
  name = hpe_morpheus_app_blueprint_terraform.testacc_blueprint.name
}

data "hpe_morpheus_instance_type" "testacc_instance_type" {
  name = hpe_morpheus_instance_type.testacc_instance_type.name
}

data "hpe_morpheus_task" "testacc_task" {
  name = hpe_morpheus_task_groovy_script.testacc_task.name
}

data "hpe_morpheus_workflow" "testacc_workflow" {
  name = hpe_morpheus_workflow_operational.testacc_workflow.name
}

resource "hpe_morpheus_role" "testacc_role_all_permissions_tenant_role_ok" {
  name      = "` + name + `"
  role_type = "tenant"

  permissions = {
	feature_permissions = [
	  {
		code   = "activity"
		access = "read"
	  },
	  {
		code   = "admin-accounts"
		access = "full"
	  }
	]
	cloud_permissions = [
	  {
		id     = data.hpe_morpheus_cloud.testacc_cloud.id
		access = "full"
	  }
	]
	blueprint_permissions = [
	  {
		id     = data.hpe_morpheus_blueprint.testacc_blueprint.id
		access = "full"
	  }
	]
	instance_type_permissions = [
	  {
		id     = data.hpe_morpheus_instance_type.testacc_instance_type.id
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
		id     = data.hpe_morpheus_task.testacc_task.id
		access = "full"
	  }
	]
	workflow_permissions = [
	  {
		id     = data.hpe_morpheus_workflow.testacc_workflow.id
		access = "full"
	  }
	]
	default_cloud_access             = "full"
	default_blueprint_access         = "full"
	default_catalog_item_type_access = "full"
	default_instance_type_access     = "full"
	default_persona_access           = "full"
	default_report_type_access       = "full"
	default_task_access              = "full"
	default_workflow_access          = "full"
	default_vdi_pool_access          = "full"
  }
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"role_type",
			"tenant",
		),
		// check fields not available for tenant roles
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"multitenant_locked",
			"false",
		),
		// check fields not available for tenant roles
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_group_access",
		),
		// check the default permission access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_cloud_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_workflow_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.default_vdi_pool_access",
			"full",
		),
		// check the permissions for resources already existing in morpheus
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "activity",
				"access": "read",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "admin-accounts",
				"access": "full",
			},
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.persona_permissions.0.access",
			"full",
		),
		// check the permissions for the resources created with the legacy provider
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.cloud_permissions.0.id",
			"data.hpe_morpheus_cloud.testacc_cloud",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.cloud_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.blueprint_permissions.0.id",
			"data.hpe_morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.blueprint_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.instance_type_permissions.0.id",
			"data.hpe_morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.instance_type_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.task_permissions.0.id",
			"data.hpe_morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
			"permissions.workflow_permissions.0.id",
			"data.hpe_morpheus_workflow.testacc_workflow",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyResourceConfig,
				// one of the blueprints values will be computed
				// so this has to be set to `true`
				ExpectNonEmptyPlan: true,
				PlanOnly:           false,
			},
			{
				Config: providerConfig + dependencyResourceConfig + resourceConfig,
				// one of the blueprints values will be computed
				// so this has to be set to `true`
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.default_group_access",
				},
				ResourceName: "hpe_morpheus_role.testacc_role_all_permissions_tenant_role_ok",
				Check:        checkFn,
			},
		},
	})
}
