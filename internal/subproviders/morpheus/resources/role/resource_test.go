// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "ExampleRole" Multitenant "false" Description "An example role" RoleType "user"
//go:generate go run ../../../../../cmd/render example-using-legacy-provider.tf.tmpl TaskDataSourceName "example_legacy_task" TaskName "example_task" ResourceName "example_with_legacy_provider" Name "ExampleRoleWithLegacyProvider" Description "An example role using legacy provider" RoleType "user" Task0Access "full"

//go:build experimental

package role_test

import (
	"os"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

// Some notes about what we expect to happen with Permissions in acceptance test import testing:

// On import, if the permissions have been computed at create,
// then the import step will pass happily.
// If the permissions have been set by the user at create,
// then the import verification step will fail,
// because the API permissions being imported do not match the
// existing resource's permissions in state.

// Therefore, for any tests using user-set permissions,
// we skip the permissions import verification check.

// Check that we can create a role with only required attributes specified
func TestAccMorpheusRoleRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "example_all" {
  name = "` + name + `"
  description = "test"
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
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.default_catalog_item_type_access",
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
	defer testhelpers.RecordResult(t)
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
				//nolint:lll
				ImportStateVerifyIgnore: []string{"permissions"}, // ignore verification on computed permissions (import)
				ResourceName:            "hpe_morpheus_role.example",
				Check:                   checkFn,
			},
		},
	})
}

func TestAccMorpheusRolePermissionsDefaultAccessPermissionsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
				ImportStateVerifyIgnore: []string{"permissions.feature_permissions"},
				ResourceName:            "hpe_morpheus_role.default_access_permissions_ok",
				Check:                   checkFn,
			},
		},
	})
}

// Tests that our mixed usage for legacy provider example
// file template used for docs is a valid config
func TestAccMorpheusRoleExampleLegacyProviderOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfigLegacy := testhelpers.ProviderBlockLegacy()
	providerConfigMixed := testhelpers.ProviderBlockMixed()

	// for setting up all of the required legacy resources to be tested
	resourceConfigLegacy := `
resource "morpheus_groovy_script_task" "testacc_role_example_legacy_provider_task" {
  name                = "` + name + `"
  source_type         = "local"
}
`

	resourceConfig, err := testhelpers.RenderExample(t, "example-using-legacy-provider.tf.tmpl",
		"TaskDataSourceName", "legacy_task_data_source",
		"TaskName", name,
		"ResourceName", "testacc_example_role_legacy_provider",
		"Name", name,
		"Description", "An example role using legacy provider",
		"RoleType", "user",
		"Task0Access", "full",
	)
	if err != nil {
		t.Fatal(err)
	}

	// perform these checks on the resource created with the old provider
	checksLegacy := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"morpheus_groovy_script_task.testacc_role_example_legacy_provider_task",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"morpheus_groovy_script_task.testacc_role_example_legacy_provider_task",
			"source_type",
			"local",
		),
	}

	// perform these checks on the resource created with the new provider
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"description",
			"An example role using legacy provider",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"role_type",
			"user",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"permissions.task_permissions.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"permissions.task_permissions.0.id",
			"morpheus_groovy_script_task.testacc_role_example_legacy_provider_task",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_example_role_legacy_provider",
			"permissions.task_permissions.0.access",
			"full",
		),
	}

	checkFnLegacy := resource.ComposeAggregateTestCheckFunc(checksLegacy...)
	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfigLegacy + resourceConfigLegacy,
				ExpectNonEmptyPlan: false,
				Check:              checkFnLegacy,
				PlanOnly:           false,
			},
			{
				Config:             providerConfigMixed + resourceConfigLegacy + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				// check only task permissions for import
				ImportStateVerifyIgnore: []string{
					"permissions.feature_permissions",
					"permissions.cloud_permissions",
					"permissions.catalog_item_type_permissions",
					"permissions.group_permissions",
					"permissions.blueprint_permissions",
					"permissions.instance_type_permissions",
					"permissions.persona_permissions",
					"permissions.report_type_permissions",
					"permissions.workflow_permissions",
					"permissions.vdi_pool_permissions",
					"permissions.default_group_access",
					"permissions.default_catalog_item_type_access",
					"permissions.default_instance_type_access",
					"permissions.default_persona_access",
					"permissions.default_report_type_access",
					"permissions.default_task_access",
					"permissions.default_workflow_access",
					"permissions.default_vdi_pool_access",
					"permissions.default_blueprint_access",
				},
				ResourceName: "hpe_morpheus_role.testacc_example_role_legacy_provider",
				Check:        checkFn,
			},
		},
	})
}

// test that we can create a user role with all possible permissions set using strongly-typed permissions
// we test all possible permissions EXCEPT VDI Pool.
// For now, the VDI pool section of the OpenAPI spec looks to be incorrect
// and needs to be updated so that we can create one using the generated SDK.
func TestAccMorpheusRoleAllPermissionsUserRoleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlockMixed()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig := `
resource "hpe_morpheus_group" "testacc_group" {
  name = "` + name + `"
}

resource "morpheus_terraform_app_blueprint" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
}

resource "morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "morpheus_groovy_script_task" "testacc_task" {
  name = "` + name + `"
  source_type         = "local"
}

resource "morpheus_operational_workflow" "testacc_workflow" {
  name = "` + name + `"
}
`

	resourceConfig := `
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
			"data.morpheus_blueprint.testacc_blueprint",
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
			"data.morpheus_instance_type.testacc_instance_type",
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
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
				ImportState:             true,
				ImportStateVerify:       true, // Check state post import
				ImportStateVerifyIgnore: []string{"permissions.feature_permissions"},
				ResourceName:            "hpe_morpheus_role.testacc_role_all_permissions_user_role_ok",
				Check:                   checkFn,
			},
		},
	})

}

// the difference between user and account role is that user roles can be assigned
// group permissions while account roles can be assigned cloud permissions
func TestAccMorpheusRoleAllPermissionsAccountRoleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlockMixed()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig := `
resource "morpheus_standard_cloud" "testacc_cloud" {
  name = "` + name + `"
  code = "standard"
  tenant_id = 1
  visibility = "public" # cloud must be visible to the client for the zone permissions to be set
}

resource "morpheus_terraform_app_blueprint" "testacc_blueprint" {
  name = "` + name + `"
  source_type = "hcl"
}

resource "morpheus_instance_type" "testacc_instance_type" {
  name = "` + name + `"
  code = "` + name + `"
  visibility = "public"
  category = "cloud"
}

resource "morpheus_groovy_script_task" "testacc_task" {
  name = "` + name + `"
  source_type         = "local"
}

resource "morpheus_operational_workflow" "testacc_workflow" {
  name = "` + name + `"
}
`

	resourceConfig := `
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

resource "hpe_morpheus_role" "testacc_role_all_permissions_account_role_ok" {
  name      = "` + name + `"
  role_type = "account"

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
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"role_type",
			"account",
		),
		// check the default permission access levels
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_cloud_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_workflow_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.default_vdi_pool_access",
			"full",
		),
		// check the permissions for resources already existing in morpheus
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "activity",
				"access": "read",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.feature_permissions.*",
			map[string]string{
				"code":   "admin-accounts",
				"access": "full",
			},
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.persona_permissions.0.code",
			"standard",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.persona_permissions.0.access",
			"full",
		),
		// check the permissions for the resources created with the legacy provider
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.cloud_permissions.0.id",
			"data.morpheus_cloud.testacc_cloud",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.cloud_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.blueprint_permissions.0.id",
			"data.morpheus_blueprint.testacc_blueprint",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.blueprint_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.instance_type_permissions.0.id",
			"data.morpheus_instance_type.testacc_instance_type",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.instance_type_permissions.0.access",
			"full",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.task_permissions.0.id",
			"data.morpheus_task.testacc_task",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
			"permissions.workflow_permissions.0.id",
			"data.morpheus_workflow.testacc_workflow",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
				ImportState:             true,
				ImportStateVerify:       true, // Check state post import
				ImportStateVerifyIgnore: []string{"permissions.feature_permissions"},
				ResourceName:            "hpe_morpheus_role.testacc_role_all_permissions_account_role_ok",
				Check:                   checkFn,
			},
		},
	})

}
