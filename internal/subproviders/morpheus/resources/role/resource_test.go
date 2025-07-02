// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "ExampleRole" Multitenant "false" Description "An example role" RoleType "user"
//go:generate go run ../../../../../cmd/render example-using-legacy-provider.tf.tmpl TaskDataSourceName "example_legacy_task" TaskName "example_task" ResourceName "example_with_legacy_provider" Name "ExampleRoleWithLegacyProvider" Description "An example role using legacy provider" RoleType "user" Task0Access "full"

package role_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

	resourceConfig := `
resource "hpe_morpheus_role" "example_required" {
  name = "TestAccMorpheusRoleRequiredAttrsOk"
}
`
	checks := []resource.TestCheckFunc{
		// required
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"name",
			"TestAccMorpheusRoleRequiredAttrsOk",
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
		composeCheckFnStatePermissionsEqAPIPermissions(
			t,
			"hpe_morpheus_role.example_required",
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
				ResourceName:      "hpe_morpheus_role.example_required",
				Check:             checkFn,
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

	resourceConfig := `
resource "hpe_morpheus_role" "example_all" {
  name = "TestAccMorpheusRoleAllAttrsOk"
  description = "test"
  landing_url = "https://test.com"
  multitenant = true
  multitenant_locked = true
  role_type = "user"
  permissions = jsonencode({
    "featurePermissions": [
      {
        "code" = "integrations-ansible"
        "access" = "full"
      }
    ],
    "globalSiteAccess" = "full"
  })
}
`
	// jsonencode() will have sorted the keys in objects
	//nolint:lll
	expectedPermissionsJSON := `{"featurePermissions":[{"access":"full","code":"integrations-ansible"}],"globalSiteAccess":"full"}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"name",
			"TestAccMorpheusRoleAllAttrsOk",
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
			"permissions",
			expectedPermissionsJSON,
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
				ResourceName:            "hpe_morpheus_role.example_all",
				Check:                   checkFn,
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

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", "TestAccMorpheusRoleExampleOk",
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
			"TestAccMorpheusRoleExampleOk",
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

// default == global
func TestAccMorpheusRolePermissionsDefaultAccessPermissionsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_role" "default_access_permissions_ok" {
	name = "TestAccMorpheusRolePermissionsDefaultAccessPermissionsOk"
	permissions = jsonencode({
  "globalSiteAccess" = "full"
  "globalZoneAccess" = "full"
  "globalInstanceTypeAccess" = "full"
  "globalAppTemplateAccess" = "full"
  "globalCatalogItemTypeAccess" = "full"
  "globalPersonaAccess" = "full"
  "globalVdiPoolAccess" = "full"
  "globalReportTypeAccess" = "full"
  "globalTaskAccess" = "full"
  "globalTaskSetAccess" = "full"
})
}
`
	//nolint:lll
	// the input will have been sorted by the jsonencode() function
	expectedDefaultPermissionsJSON := `{"globalAppTemplateAccess":"full","globalCatalogItemTypeAccess":"full","globalInstanceTypeAccess":"full","globalPersonaAccess":"full","globalReportTypeAccess":"full","globalSiteAccess":"full","globalTaskAccess":"full","globalTaskSetAccess":"full","globalVdiPoolAccess":"full","globalZoneAccess":"full"}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"name",
			"TestAccMorpheusRolePermissionsDefaultAccessPermissionsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.default_access_permissions_ok",
			"permissions",
			expectedDefaultPermissionsJSON,
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
				ResourceName:      "hpe_morpheus_role.default_access_permissions_ok",
				//nolint:lll
				ImportStateVerifyIgnore: []string{"permissions"}, // ignore verification on computed permissions (import)
				Check:                   checkFn,
			},
		},
	})
}

// check that we correctly store the API-computed permissions in the statefile when
// the user has not set any permissions
func TestAccMorpheusRolePermissionsComputedPermissionsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_role" "computed_permissions_ok" {
	name = "TestAccMorpheusRolePermissionsComputedPermissionsOk"
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.computed_permissions_ok",
			"name",
			"TestAccMorpheusRolePermissionsComputedPermissionsOk",
		),
		composeCheckFnStatePermissionsEqAPIPermissions(
			t,
			"hpe_morpheus_role.computed_permissions_ok",
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
				ResourceName:      "hpe_morpheus_role.computed_permissions_ok",
				Check:             checkFn,
			},
		},
	})
}

// test that providing feature permissions with a JSON string literal works
func TestAccMorpheusRolePermissionsFeaturePermissionsJSONStringOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	featurePermissionsJSON := `
{
  "featurePermissions": [
    {
      "code": "integrations-ansible",
      "access": "full"
    },
    {
      "code": "admin-appliance",
      "access": "none"
    },
    {
      "code": "app-templates",
      "access": "none"
    }
  ]
}
`
	resourceConfig := fmt.Sprintf(`resource "hpe_morpheus_role" "json_string_ok" {
	name = "TestAccMorpheusRolePermissionsFeaturePermissionsOk"
	permissions = <<-EOT
%sEOT
}
`, featurePermissionsJSON)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.json_string_ok",
			"name",
			"TestAccMorpheusRolePermissionsFeaturePermissionsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.json_string_ok",
			"permissions",
			featurePermissionsJSON,
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
				ResourceName:            "hpe_morpheus_role.json_string_ok",
				Check:                   checkFn,
			},
		},
	})
}

// test that we can set permissions using jsonencode()
func TestAccMorpheusRolePermissionsFeaturePermissionsJSONEncodeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `resource "hpe_morpheus_role" "json_encode_ok" {
name = "TestAccMorpheusRolePermissionsFeaturePermissionsOk"
permissions = jsonencode({
  "featurePermissions": [
    {
      "code" = "integrations-ansible"
      "access" = "full"
    },
    {
      "code" = "admin-appliance"
      "access" = "none"
    }
  ]
})
}
`

	// remember, jsonencode() sorts the keys of an object
	//nolint:lll
	expectedFeaturePermissionsJSON := `{"featurePermissions":[{"access":"full","code":"integrations-ansible"},{"access":"none","code":"admin-appliance"}]}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.json_encode_ok",
			"name",
			"TestAccMorpheusRolePermissionsFeaturePermissionsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.json_encode_ok",
			"permissions",
			expectedFeaturePermissionsJSON,
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
				ResourceName:            "hpe_morpheus_role.json_encode_ok",
				Check:                   checkFn,
			},
		},
	})
}

// test that there's no change in plan after running an apply
func TestAccMorpheusRolePermissionsPlanAfterApply(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfigGood := `resource "hpe_morpheus_role" "plan_after_apply_good" {
name = "TestAccMorpheusRolePermissionsPlanAfterApplyGoodPermissions"
permissions = jsonencode({
  "featurePermissions": [
    {
      "code" = "integrations-ansible"
      "access" = "full"
    }
  ],
  "globalSiteAccess" = "full"
})
}
`
	resourceConfigBad := `resource "hpe_morpheus_role" "plan_after_apply_bad" {
name = "TestAccMorpheusRolePermissionsPlanAfterApplyBadPermissions"
permissions = jsonencode({
  "globalSiteAccessFoo" = "full"
})
}
`

	// remember, jsonencode() sorts the keys of an object
	//nolint:lll
	expectedGoodPermissionsJSON := `{"featurePermissions":[{"access":"full","code":"integrations-ansible"}],"globalSiteAccess":"full"}`

	expectedBadPermissionsJSON := `{"globalSiteAccessFoo":"full"}`

	checksGood := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.plan_after_apply_good",
			"name",
			"TestAccMorpheusRolePermissionsPlanAfterApplyGoodPermissions",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.plan_after_apply_good",
			"permissions",
			expectedGoodPermissionsJSON,
		),
	}

	checksBad := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.plan_after_apply_bad",
			"name",
			"TestAccMorpheusRolePermissionsPlanAfterApplyBadPermissions",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.plan_after_apply_bad",
			"permissions",
			expectedBadPermissionsJSON,
		),
	}

	checkFnGood := resource.ComposeAggregateTestCheckFunc(checksGood...)
	checkFnBad := resource.ComposeAggregateTestCheckFunc(checksBad...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfigGood,
				ExpectNonEmptyPlan: false, // works on refresh plan after apply, too
				Check:              checkFnGood,
				ResourceName:       "hpe_morpheus_role.plan_after_apply_good",
				PlanOnly:           false,
			},
			{
				Config:             providerConfig + resourceConfigBad,
				ExpectNonEmptyPlan: true, // works on refresh plan after apply, too
				Check:              checkFnBad,
				ResourceName:       "hpe_morpheus_role.plan_after_apply_bad",
				PlanOnly:           false,
			},
		},
	})
}

func TestAccMorpheusRolePermissionsTaskPermissionsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfigLegacy := testhelpers.ProviderBlockLegacy()
	providerConfigMixed := testhelpers.ProviderBlockMixed()

	legacyTaskResourceConfig := `
resource "morpheus_groovy_script_task" "testacc_groovy" {
  name                = "testacc_groovy"
  source_type         = "local"
}
`
	resourceConfigMixed := `
data "morpheus_task" "testacc_task" {
  name = "testacc_groovy"
}

resource "hpe_morpheus_role" "testacc_role_task_permissions" {
  name               = "TestAccMorpheusRolePermissionsTaskPermissionsOk"
  permissions = jsonencode({
    "taskPermissions" : [
      {
        "id" = data.morpheus_task.testacc_task.id
        "access" : "full"
      }
    ]
    }
  )
}
`
	// confirms that the task dependency was set up correctly
	checksTask := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"morpheus_groovy_script_task.testacc_groovy",
			"name",
			"testacc_groovy",
		),
		resource.TestCheckResourceAttr(
			"morpheus_groovy_script_task.testacc_groovy",
			"source_type",
			"local",
		),
	}
	checkFnTask := resource.ComposeAggregateTestCheckFunc(checksTask...)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.testacc_role_task_permissions",
			"name",
			"TestAccMorpheusRolePermissionsTaskPermissionsOk",
		),
		func(s *terraform.State) error {
			taskID, err := testhelpers.ExtractValue(
				s,
				"data.morpheus_task.testacc_task",
				"id",
			)
			if err != nil {
				return err
			}
			taskIDInt, err := strconv.Atoi(taskID)
			if err != nil {
				return err
			}

			rs := s.RootModule().Resources["hpe_morpheus_role.testacc_role_task_permissions"]
			if rs == nil {
				return errors.New("resource not found: hpe_morpheus_role.testacc_role_task_permissions")
			}

			statePermisions := rs.Primary.Attributes["permissions"]

			expectedPermissionsJSON := fmt.Sprintf(
				`{"taskPermissions":[{"access":"full","id":%d}]}`,
				taskIDInt)

			if statePermisions != expectedPermissionsJSON {
				return fmt.Errorf(
					"permissions in state do not match expected permissions:\nexpected: %s\ngot: %s",
					expectedPermissionsJSON, statePermisions)
			}

			return nil
		},
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	// Use the hpe/morpheus provider to create a Role using the ID from the legacy task resource

	// Set up the task resource using the legacy provider
	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Set up the task resource in a separate step so it's created for the following one
			{
				Config: providerConfigLegacy + legacyTaskResourceConfig,
				Check:  checkFnTask,
			},
			// Need to include the task resource config from the previous step
			// or else Terraform will assume it no longer exists
			{
				Config:             legacyTaskResourceConfig + providerConfigMixed + resourceConfigMixed,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
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

	providerConfigLegacy := testhelpers.ProviderBlockLegacy()
	providerConfigMixed := testhelpers.ProviderBlockMixed()

	// for setting up all of the required legacy resources to be tested
	resourceConfigLegacy := `
resource "morpheus_groovy_script_task" "testacc_role_example_legacy_provider_task" {
  name                = "testacc_role_example_legacy_provider_task"
  source_type         = "local"
}
`

	resourceConfig, err := testhelpers.RenderExample(t, "example-using-legacy-provider.tf.tmpl",
		"TaskDataSourceName", "testacc_role_example_legacy_provider_task_datasource",
		"TaskName", "testacc_role_example_legacy_provider_task",
		"ResourceName", "testacc_example_role_legacy_provider",
		"Name", "TestAccMorpheusRoleExampleLegacyProviderOk",
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
			"testacc_role_example_legacy_provider_task",
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
			"TestAccMorpheusRoleExampleLegacyProviderOk",
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
		// to test the permissions
		func(s *terraform.State) error {
			taskID, err := testhelpers.ExtractValue(
				s,
				"data.morpheus_task.testacc_role_example_legacy_provider_task_datasource",
				"id",
			)
			if err != nil {
				return err
			}
			taskIDInt, err := strconv.Atoi(taskID)
			if err != nil {
				return err
			}

			rs := s.RootModule().Resources["hpe_morpheus_role.testacc_example_role_legacy_provider"]
			if rs == nil {
				return errors.New("resource not found: hpe_morpheus_role.testacc_example_role_legacy_provider")
			}

			statePermisions := rs.Primary.Attributes["permissions"]

			expectedPermissionsJSON := fmt.Sprintf(
				`{"taskPermissions":[{"access":"full","id":%d}]}`,
				taskIDInt)

			if statePermisions != expectedPermissionsJSON {
				return fmt.Errorf(
					"permissions in state do not match expected permissions:\nexpected: %s\ngot: %s",
					expectedPermissionsJSON, statePermisions)
			}

			return nil
		},
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
				//nolint:lll
				ImportStateVerifyIgnore: []string{"permissions"}, // ignore verification on computed permissions (import)
				ResourceName:            "hpe_morpheus_role.testacc_example_role_legacy_provider",
				Check:                   checkFn,
			},
		},
	})
}

// Needed for when we want to verify entirely computed permissions in state.
// We can't compare against a string constant because the IDs of the featurePermissions can
// differ between Morpheus installs; presumably computed in parallel at Morpheus initialisation.
func composeCheckFnStatePermissionsEqAPIPermissions(
	t *testing.T,
	resource string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resource]
		if rs == nil {
			return fmt.Errorf("resource not found: %s", resource)
		}

		roleID := rs.Primary.Attributes["id"]
		roleIDInt, err := strconv.Atoi(roleID)
		if err != nil {
			return err
		}

		roleResp, err := testhelpers.GetRole(t, int64(roleIDInt))
		if err != nil {
			return err
		}

		// don't need it for marshaling to do comparison
		roleResp.Role = nil

		apiPermissions, err := json.Marshal(roleResp)
		if err != nil {
			return err
		}

		apiPermissionsStr := string(apiPermissions)

		statePermisions := rs.Primary.Attributes["permissions"]

		// the state Permissions should have already been sorted by a json.Marshal at create time
		if apiPermissionsStr != statePermisions {
			return fmt.Errorf("permissions in state do not match API permissions:\nexpected: %s\ngot: %s",
				statePermisions, apiPermissions)
		}

		return nil
	}
}
