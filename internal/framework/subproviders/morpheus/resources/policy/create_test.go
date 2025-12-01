// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// Test creating a policy with required attributes only
func TestAccMorpheusPolicyRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "required" {
  name = "` + name + `"
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

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"associated_resource_type",
			"Group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"policy_type.code",
			"maxMemory",
		),
		// Check computed defaults
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"enabled",
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
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
				ResourceName:            "hpe_morpheus_policy.required",
				Check:                   checkFn,
			},
		},
	})
}

// Test creating policies with different policy types which apply to Bare Metal
func TestAccMorpheusPolicyAllBareMetalPolicyTypesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlockMixed()
	namePrefix := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")
	cloudName := acctest.RandomWithPrefix(t.Name() + "-cloud")

	resourceConfig := `
variable "policy_name" {
  type = string
}

variable "policy_description" {
  type = string
}

variable "policy_type_code" {
  type = string
}

variable "policy_config" {
  type = map(any)
}

resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
  tenant_id = 1
  group_id = hpe_morpheus_group.test.id
  code = "` + cloudName + `"
  cloud_type_code = "standard"
  
  config = {
    certificateProvider = "internal"
    enableNetworkTypeSelection = false
  }
}

resource "morpheus_operational_workflow" "test" {
  name = "` + namePrefix + `-workflow"
}

data "morpheus_workflow" "test" {
  name = morpheus_operational_workflow.test.name
}

resource "hpe_morpheus_policy" "test" {
  name = var.policy_name
  description = var.policy_description
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = var.policy_type_code
  }
  
  config = var.policy_config
}
`

	resourceConfigCloud := `
variable "policy_name" {
  type = string
}

variable "policy_description" {
  type = string
}

variable "policy_type_code" {
  type = string
}

variable "policy_config" {
  type = map(any)
}

resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
  tenant_id = 1
  group_id = hpe_morpheus_group.test.id
  code = "` + cloudName + `"
  cloud_type_code = "standard"
  
  config = {
    certificateProvider = "internal"
    enableNetworkTypeSelection = false
  }
}

resource "morpheus_operational_workflow" "test" {
  name = "` + namePrefix + `-workflow"
}

data "morpheus_workflow" "test" {
  name = morpheus_operational_workflow.test.name
}

resource "hpe_morpheus_policy" "test" {
  name = var.policy_name
  description = var.policy_description
  associated_resource_type = "Cloud"
  associated_resource_id = hpe_morpheus_cloud.test.id
  enabled = true
  
  policy_type = {
    code = var.policy_type_code
  }
  
  config = var.policy_config
}
`

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Approve Delete
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-deleteApproval"),
					"policy_description": config.StringVariable("Delete approval policy"),
					"policy_type_code":   config.StringVariable("deleteApproval"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"accountIntegrationId": config.StringVariable("1"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-deleteApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "deleteApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Delete approval policy"),
				),
			},
			// Step 2: Approve Provision
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-provisionApproval"),
					"policy_description": config.StringVariable("Provision approval policy"),
					"policy_type_code":   config.StringVariable("provisionApproval"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"accountIntegrationId": config.StringVariable("1"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-provisionApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "provisionApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Provision approval policy"),
				),
			},
			// Step 3: Approve Reconfigure
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-reconfigureApproval"),
					"policy_description": config.StringVariable("Reconfigure approval policy"),
					"policy_type_code":   config.StringVariable("reconfigureApproval"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"accountIntegrationId": config.StringVariable("1"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-reconfigureApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "reconfigureApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Reconfigure approval policy"),
				),
			},
			// Step 4: Approve Workflow Execute
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-workflowApproval"),
					"policy_description": config.StringVariable("Workflow approval policy"),
					"policy_type_code":   config.StringVariable("workflowApproval"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"accountIntegrationId": config.StringVariable("1"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-workflowApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "workflowApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Workflow approval policy"),
				),
			},
			// Step 5: Delayed Delete
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-delayedRemoval"),
					"policy_description": config.StringVariable("Delayed removal policy"),
					"policy_type_code":   config.StringVariable("delayedRemoval"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"removalAge": config.StringVariable("7"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-delayedRemoval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "delayedRemoval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Delayed removal policy"),
				),
			},
			// Step 6: Instance Naming
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-naming"),
					"policy_description": config.StringVariable("Naming policy"),
					"policy_type_code":   config.StringVariable("naming"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"namingType":     config.StringVariable("fixed"),
						"namingPattern":  config.StringVariable("instance-${sequence}"),
						"namingConflict": config.BoolVariable(true),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-naming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "naming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Naming policy"),
				),
			},
			// Step 7: Instance Networks
			{
				Config: providerConfig + resourceConfigCloud,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-requiredNetwork"),
					"policy_description": config.StringVariable("Required network policy"),
					"policy_type_code":   config.StringVariable("requiredNetwork"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"requiredNetworks": config.ListVariable(
							config.IntegerVariable(1),
						),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-requiredNetwork"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "requiredNetwork"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Required network policy"),
				),
			},
			// Step 8: Max Memory
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxMemory"),
					"policy_description": config.StringVariable("Max memory policy"),
					"policy_type_code":   config.StringVariable("maxMemory"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxMemory": config.IntegerVariable(1073741824),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxMemory"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxMemory"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max memory policy"),
				),
			},
			// Step 9: Max Pool Members
			{
				Config: providerConfig + resourceConfigCloud,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxPoolMembers"),
					"policy_description": config.StringVariable("Max pool members policy"),
					"policy_type_code":   config.StringVariable("maxPoolMembers"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxPoolMembers": config.StringVariable("10"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPoolMembers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPoolMembers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max pool members policy"),
				),
			},
			// Step 10: Max Snapshots
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxSnapshots"),
					"policy_description": config.StringVariable("Max snapshots policy"),
					"policy_type_code":   config.StringVariable("maxSnapshots"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxSnapshots": config.StringVariable("5"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxSnapshots"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxSnapshots"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max snapshots policy"),
				),
			},
			// Step 11: Max Storage
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxStorage"),
					"policy_description": config.StringVariable("Max storage policy"),
					"policy_type_code":   config.StringVariable("maxStorage"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxStorage":        config.StringVariable("10737418240"),
						"excludeContainers": config.BoolVariable(false),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxStorage"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxStorage"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max storage policy"),
				),
			},
			// Step 12: Max VMs
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxVms"),
					"policy_description": config.StringVariable("Max VMs policy"),
					"policy_type_code":   config.StringVariable("maxVms"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxVms": config.StringVariable("20"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxVms"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxVms"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max VMs policy"),
				),
			},
			// Step 13: Max Networks
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxNetworks"),
					"policy_description": config.StringVariable("Max networks policy"),
					"policy_type_code":   config.StringVariable("maxNetworks"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxNetworks": config.StringVariable("5"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxNetworks"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxNetworks"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max networks policy"),
				),
			},
			/* // Step 14: Storage Server Quota - Commented out, requires Global scope which impacts other users
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-storageServerQuota"),
					"policy_description": config.StringVariable("Storage server quota policy"),
					"policy_type_code":   config.StringVariable("storageServerQuota"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"storageServerId": config.StringVariable("1"),
						"maxStorage":      config.StringVariable("10737418240"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-storageServerQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "storageServerQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Storage server quota policy"),
				),
			},
			*/ // Step 15: Tags
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-tags"),
					"policy_description": config.StringVariable("Tags policy"),
					"policy_type_code":   config.StringVariable("tags"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"strict": config.BoolVariable(true),
						"key":    config.StringVariable("environment"),
						"value":  config.StringVariable("production"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-tags"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "tags"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Tags policy"),
				),
			},
			// Step 16: User Creation
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-createUser"),
					"policy_description": config.StringVariable("Create user policy"),
					"policy_type_code":   config.StringVariable("createUser"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"createUserType": config.StringVariable("fixed"),
						"createUser":     config.StringVariable("on"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createUser"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createUser"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Create user policy"),
				),
			},
			// Step 17: User Group Creation
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-createUserGroup"),
					"policy_description": config.StringVariable("Create user group policy"),
					"policy_type_code":   config.StringVariable("createUserGroup"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"userGroup": config.StringVariable("1"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createUserGroup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createUserGroup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Create user group policy"),
				),
			},
			// Step 18: Workflow (uses Morpheus provider to create workflow resource)
			{
				Config: providerConfig + `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "morpheus_operational_workflow" "test" {
  name = "` + namePrefix + `-workflow"
}

data "morpheus_workflow" "test" {
  name = morpheus_operational_workflow.test.name
}

resource "hpe_morpheus_policy" "test" {
  name = "` + namePrefix + `-workflow"
  description = "Workflow policy"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "workflow"
  }
  
  config = {
    workflowId = data.morpheus_workflow.test.id
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-workflow"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "workflow"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Workflow policy"),
				),
			},
		},
	})
}

// Test creating policies scoped to different resource types (Group, Cloud, User, Role)
func TestAccMorpheusPolicyResourceTypesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Create dependency resources
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")
	cloudName := acctest.RandomWithPrefix(t.Name() + "-cloud")
	roleName := acctest.RandomWithPrefix(t.Name() + "-role")
	userName := acctest.RandomWithPrefix(t.Name() + "-user")
	planCode := acctest.RandomWithPrefix(t.Name() + "-plan")
	networkName := acctest.RandomWithPrefix(t.Name() + "-network")

	dependencyConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
  tenant_id = 1
  group_id = hpe_morpheus_group.test.id
  code = "` + cloudName + `"
  cloud_type_code = "standard"
  
  config = {
    certificateProvider = "internal"
    enableNetworkTypeSelection = false
  }
}

resource "hpe_morpheus_role" "test" {
  name = "` + roleName + `"
  role_type = "user"
}

resource "hpe_morpheus_user" "test" {
  username = "` + userName + `"
  email = "` + userName + `@test.com"
  role_ids = [hpe_morpheus_role.test.id]
  password_wo = "TestPassword123!"
}

resource "hpe_morpheus_service_plan" "test" {
  name = "` + planCode + `"
  code = "` + planCode + `"
  sort_order = 10000
  max_memory = 4294967296
  max_storage = 536870912
  provision_type_code = "arm"
  custom_max_storage = true
  cores_per_socket = 1
  config_ranges = {
    min_storage = 268435456
    max_storage = 536870912
  }
}

resource "hpe_morpheus_network" "test" {
  name = "` + networkName + `"
  description = "Test network"
  cloud_id = hpe_morpheus_cloud.test.id
  pool_id = 1
  group_id = hpe_morpheus_group.test.id
  type_id = 1
  config = {}
  active = true
  dhcp_server = false
  appliance_url_proxy_bypass = true
  tenant_ids = [1]
  visibility = "private"
  cidr = "10.0.0.0/24"
  labels = ["terraform", "test"]
}
`

	policyNameGroup := acctest.RandomWithPrefix(t.Name() + "-group-policy")
	policyNameCloud := acctest.RandomWithPrefix(t.Name() + "-cloud-policy")
	policyNameRole := acctest.RandomWithPrefix(t.Name() + "-role-policy")
	policyNameUser := acctest.RandomWithPrefix(t.Name() + "-user-policy")
	policyNamePlan := acctest.RandomWithPrefix(t.Name() + "-plan-policy")
	policyNameNetwork := acctest.RandomWithPrefix(t.Name() + "-network-policy")

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "group_policy",
		"Name", policyNameGroup,
		"Description", "Example group-scoped policy",
		"AssociatedResourceType", "Group",
		"AssociatedResourceID", "hpe_morpheus_group.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	cloudResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "cloud_policy",
		"Name", policyNameCloud,
		"Description", "Example cloud-scoped policy",
		"AssociatedResourceType", "Cloud",
		"AssociatedResourceID", "hpe_morpheus_cloud.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	roleResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "role_policy",
		"Name", policyNameRole,
		"Description", "Example role-scoped policy",
		"AssociatedResourceType", "Role",
		"AssociatedResourceID", "hpe_morpheus_role.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	userResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "user_policy",
		"Name", policyNameUser,
		"Description", "Example user-scoped policy",
		"AssociatedResourceType", "User",
		"AssociatedResourceID", "hpe_morpheus_user.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	planResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "plan_policy",
		"Name", policyNamePlan,
		"Description", "Example plan-scoped policy",
		"AssociatedResourceType", "Plan",
		"AssociatedResourceID", "hpe_morpheus_service_plan.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	networkResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "network_policy",
		"Name", policyNameNetwork,
		"Description", "Example network-scoped policy",
		"AssociatedResourceType", "Network",
		"AssociatedResourceID", "hpe_morpheus_network.test.id",
		"PolicyTypeCode", "maxVms",
		"ConfigKey", "maxVms",
		"ConfigValue", "20",
	)
	if err != nil {
		t.Fatal(err)
	}

	checksGroup := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.group_policy", "name", policyNameGroup),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.group_policy", "associated_resource_type", "Group"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.group_policy", "associated_resource_id"),
	}

	checksCloud := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.cloud_policy", "name", policyNameCloud),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.cloud_policy", "associated_resource_type", "Cloud"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.cloud_policy", "associated_resource_id"),
	}

	checksRole := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.role_policy", "name", policyNameRole),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.role_policy", "associated_resource_type", "Role"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.role_policy", "associated_resource_id"),
	}

	checksUser := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.user_policy", "name", policyNameUser),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.user_policy", "associated_resource_type", "User"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.user_policy", "associated_resource_id"),
	}

	checksPlan := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.plan_policy", "name", policyNamePlan),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.plan_policy", "associated_resource_type", "Plan"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.plan_policy", "associated_resource_id"),
	}

	checksNetwork := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.network_policy", "name", policyNameNetwork),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.network_policy", "associated_resource_type", "Network"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.network_policy", "associated_resource_id"),
	}

	allChecks := append(checksGroup, checksCloud...)
	allChecks = append(allChecks, checksRole...)
	allChecks = append(allChecks, checksUser...)
	allChecks = append(allChecks, checksPlan...)
	allChecks = append(allChecks, checksNetwork...)
	checkFn := resource.ComposeAggregateTestCheckFunc(
		allChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + resourceConfig +
					cloudResourceConfig + roleResourceConfig +
					userResourceConfig + planResourceConfig +
					networkResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
		},
	})
}

// Test creating policies using static schema fields (config_* attributes)
func TestAccMorpheusPolicyAllStaticSchemaOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlockMixed()
	namePrefix := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	roleName := acctest.RandomWithPrefix(t.Name() + "-role")
	userName := acctest.RandomWithPrefix(t.Name() + "-user")
	cloudName := acctest.RandomWithPrefix(t.Name() + "-cloud")

	dependencyConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_role" "test" {
  name = "` + roleName + `"
  role_type = "user"
}

resource "hpe_morpheus_user" "test" {
  username = "` + userName + `"
  email = "` + userName + `@test.com"
  role_ids = [hpe_morpheus_role.test.id]
  password_wo = "TestPassword123!"
}

resource "morpheus_operational_workflow" "test" {
  name = "` + namePrefix + `-workflow"
}

resource "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
  tenant_id = 1
  group_id = hpe_morpheus_group.test.id
  code = "` + cloudName + `"
  cloud_type_code = "standard"
  
  config = {
    certificateProvider = "internal"
    enableNetworkTypeSelection = false
  }
}

data "morpheus_workflow" "test" {
  name = morpheus_operational_workflow.test.name
}
`

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Approve Delete (using config_approval)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-deleteApproval"
			  description = "Delete approval policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "deleteApproval"
			  }

			  config_approval = {
			    account_integration_id = "1"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-deleteApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "deleteApproval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_approval.account_integration_id", "1"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Step 2: Backup Creation (using config_create_backup)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-createBackup"
			  description = "Create backup policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "createBackup"
			  }

			  config_create_backup = {
			    create_backup = true
			    create_backup_type = "user"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createBackup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createBackup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_create_backup.create_backup", "true"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Step 3: Backup Storage (using config_backup_storage) - uses Group scope - skipped, API bug
			// {
			// 	Config: providerConfig + dependencyConfig + `
			// resource "hpe_morpheus_policy" "test" {
			//   name = "` + namePrefix + `-backupStorage"
			//   description = "Backup storage policy using static schema"
			//   associated_resource_type = "Group"
			//   associated_resource_id = hpe_morpheus_group.test.id
			//   enabled = true

			//   policy_type = {
			//     code = "backupStorage"
			//   }

			//   config_backup_storage = {
			//     backup_storage_ids = [1]
			//   }
			// }`,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-backupStorage"),
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "backupStorage"),
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_backup_storage.backup_storage_ids.#", "1"),
			// 	),
			// 	ExpectNonEmptyPlan: false,
			// },
			// Step 4: Cypher Access (using config_cypher) - uses User scope
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-cypher"
			  description = "Cypher access policy using static schema"
			  associated_resource_type = "User"
			  associated_resource_id = hpe_morpheus_user.test.id
			  enabled = true

			  policy_type = {
			    code = "cypher"
			  }

			  config_cypher = {
			    read = true
			    write = true
			    update = false
			    delete = false
			    list = true
			    key_pattern = "secret/*"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-cypher"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "cypher"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_cypher.read", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_cypher.write", "true"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Step 5: Delayed Delete (using config_delayed_removal)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-delayedRemoval"
			  description = "Delayed removal policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "delayedRemoval"
			  }

			  config_delayed_removal = {
			    removal_age = "7"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-delayedRemoval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "delayedRemoval"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_delayed_removal.removal_age", "7"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Step 6: Hostname (using config_host_naming)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-hostNaming"
			  description = "Host naming policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "hostNaming"
			  }

			  config_host_naming = {
			    host_naming_type = "user"
			    host_naming_pattern = "host-$$$${sequence}"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-hostNaming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "hostNaming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_host_naming.host_naming_type", "user"),
				),
			},
			// Step 7: Expiration (using config_lifecycle) - TEMPORARILY SKIPPED
			{
				Config: providerConfig + dependencyConfig + `
							resource "hpe_morpheus_policy" "test" {
							  name = "` + namePrefix + `-lifecycle"
							  description = "Lifecycle policy using static schema"
							  associated_resource_type = "Group"
							  associated_resource_id = hpe_morpheus_group.test.id
							  enabled = true

							  policy_type = {
							    code = "lifecycle"
							  }

							  config_lifecycle = {
							  	account_integration_id = "1"
							    lifecycle_type = "fixed"
							    lifecycle_age = "30"
							  }
							}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-lifecycle"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "lifecycle"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_lifecycle.lifecycle_age", "30"),
				),
			},
			// Step 7: Max Containers (using config_max_containers)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxContainers"
			  description = "Max containers policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxContainers"
			  }

			  config_max_containers = {
			    max_containers = "10"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxContainers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxContainers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_containers.max_containers", "10"),
				),
			},
			// Step 9: Max Cores (using config_max_cores)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxCores"
			  description = "Max cores policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxCores"
			  }

			  config_max_cores = {
			    max_cores = "16"
			    exclude_containers = false
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxCores"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxCores"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_cores.max_cores", "16"),
				),
			},
			// Step 10: Max Hosts (using config_max_hosts)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxHosts"
			  description = "Max hosts policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxHosts"
			  }

			  config_max_hosts = {
			    max_hosts = "50"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxHosts"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxHosts"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_hosts.max_hosts", "50"),
				),
			},
			// Step 11: Max Memory (using config_max_memory)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxMemory"
			  description = "Max memory policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxMemory"
			  }

			  config_max_memory = {
			    max_memory = "1073741824"
			    exclude_containers = true
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxMemory"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxMemory"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_memory.max_memory", "1073741824"),
				),
			},
			// Step 12: Network Quota (using config_max_networks)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxNetworks"
			  description = "Max networks policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxNetworks"
			  }

			  config_max_networks = {
			    max_networks = "5"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxNetworks"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxNetworks"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_networks.max_networks", "5"),
				),
			},
			// Step 13: Max Pool Members (using config_max_pool_members)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxPoolMembers"
			  description = "Max pool members policy using static schema"
			  associated_resource_type = "Cloud"
			  associated_resource_id = hpe_morpheus_cloud.test.id
			  enabled = true

			  policy_type = {
			    code = "maxPoolMembers"
			  }

			  config_max_pool_members = {
			    max_pool_members = "10"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPoolMembers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPoolMembers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_pool_members.max_pool_members", "10"),
				),
			},
			// Step 14: Max Load Balancer Pools (using config_max_pools)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxPools"
			  description = "Max pools policy using static schema"
			  associated_resource_type = "Cloud"
			  associated_resource_id = hpe_morpheus_cloud.test.id
			  enabled = true

			  policy_type = {
			    code = "maxPools"
			  }

			  config_max_pools = {
			    max_pools = "8"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPools"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPools"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_pools.max_pools", "8"),
				),
			},
			// Step 15: Budget (using config_max_price)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxPrice"
			  description = "Max price policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxPrice"
			  }

			  config_max_price = {
			    max_price = 1000.50
			    max_price_currency = "USD"
			    max_price_unit = "month"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPrice"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPrice"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_price.max_price_currency", "USD"),
				),
			},
			// Step 16: Router Quota (using config_max_routers)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxRouters"
			  description = "Max routers policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxRouters"
			  }

			  config_max_routers = {
			    max_routers = "3"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxRouters"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxRouters"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_routers.max_routers", "3"),
				),
			},
			// Step 17: Max Snapshots (using config_max_snapshots)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxSnapshots"
			  description = "Max snapshots policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxSnapshots"
			  }

			  config_max_snapshots = {
			    max_snapshots = "5"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxSnapshots"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxSnapshots"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_snapshots.max_snapshots", "5"),
				),
			},
			// Step 18: Max Storage (using config_max_storage)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxStorage"
			  description = "Max storage policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxStorage"
			  }

			  config_max_storage = {
			    max_storage = "10737418240"
			    exclude_containers = false
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxStorage"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxStorage"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_storage.max_storage", "10737418240"),
				),
			},
			// Step 19: Max Virtual Servers (using config_max_virtual_servers)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxVirtualServers"
			  description = "Max virtual servers policy using static schema"
			  associated_resource_type = "Cloud"
			  associated_resource_id = hpe_morpheus_cloud.test.id
			  enabled = true

			  policy_type = {
			    code = "maxVirtualServers"
			  }

			  config_max_virtual_servers = {
			    max_virtual_servers = "15"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxVirtualServers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxVirtualServers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_virtual_servers.max_virtual_servers", "15"),
				),
			},
			// Step 20: Max VMs (using config_max_vms)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-maxVms"
			  description = "Max VMs policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "maxVms"
			  }

			  config_max_vms = {
			    max_vms = "20"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxVms"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxVms"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_max_vms.max_vms", "20"),
				),
			},
			// Step 21: Message of the Day (using config_motd) - skipped, uses global
			/*
							{
								Config: providerConfig + dependencyConfig + `
				resource "hpe_morpheus_policy" "test" {
				  name = "` + namePrefix + `-motd"
				  description = "MOTD policy using static schema"
				  associated_resource_type = "Group"
				  associated_resource_id = hpe_morpheus_group.test.id
				  enabled = true

				  policy_type = {
				    code = "motd"
				  }

				  config_motd = {
				    motdtitle = "Welcome"
				    motdmessage = "Welcome to the system"
				    motdtype = "text"
				    motddate = "2024-01-01"
				  }
				}`,
								Check: resource.ComposeAggregateTestCheckFunc(
									resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-motd"),
									resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "motd"),
									resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_motd.motdtitle", "Welcome"),
								),
							},
			*/
			// Step 22: Instance Name (using config_naming)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-naming"
			  description = "Naming policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "naming"
			  }

			  config_naming = {
			    naming_type = "user"
			    naming_pattern = "instance-$$$${sequence}"
			    naming_conflict = true
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-naming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "naming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_naming.naming_type", "user"),
				),
			},
			// Step 23: Power Scheduling (using config_power_schedule)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-powerSchedule"
			  description = "Power schedule policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "powerSchedule"
			  }

			  config_power_schedule = {
			    power_schedule = "1"
			    power_schedule_type = "user"
			    power_schedule_hide_fixed = false
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-powerSchedule"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "powerSchedule"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_power_schedule.power_schedule", "1"),
				),
			},
			// Step 24: Instance Networks (using config_required_network)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-requiredNetwork"
			  description = "Required network policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "requiredNetwork"
			  }

			  config_required_network = {
			    required_networks = [1, 2]
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-requiredNetwork"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "requiredNetwork"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_required_network.required_networks.#", "2"),
				),
			},
			// Step 25: Cluster Resource Name (using config_server_naming)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-serverNaming"
			  description = "Server naming policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "serverNaming"
			  }

			  config_server_naming = {
			    server_naming_type = "user"
			    server_naming_pattern = "server-$$$${sequence}"
			    server_naming_conflict = true
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-serverNaming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "serverNaming"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_server_naming.server_naming_type", "user"),
				),
			},
			// Step 26: Shutdown (using config_shutdown)
			{
				Config: providerConfig + dependencyConfig + `
			resource "hpe_morpheus_policy" "test" {
			  name = "` + namePrefix + `-shutdown"
			  description = "Shutdown policy using static schema"
			  associated_resource_type = "Group"
			  associated_resource_id = hpe_morpheus_group.test.id
			  enabled = true

			  policy_type = {
			    code = "shutdown"
			  }

			  config_shutdown = {
			    shutdown_type = "fixed"
			    shutdown_age = "30"
			  }
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-shutdown"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "shutdown"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_shutdown.shutdown_age", "30"),
				),
			},
			// Step 27: Storage Server Storage Quota (using config_storage_server_quota) - skipped, uses Global
			// {
			// 	Config: providerConfig + dependencyConfig + `
			// resource "hpe_morpheus_policy" "test" {
			//   name = "` + namePrefix + `-storageServerQuota"
			//   description = "Storage server quota policy using static schema"
			//   associated_resource_type = "Group"
			//   associated_resource_id = hpe_morpheus_group.test.id
			//   enabled = true

			//   policy_type = {
			//     code = "storageServerQuota"
			//   }

			//   config_storage_server_quota = {
			//     storage_server_id = "1"
			//     max_storage = "10737418240"
			//   }
			// }`,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-storageServerQuota"),
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "storageServerQuota"),
			// 		resource.TestCheckResourceAttr("hpe_morpheus_policy.test",
			// 			"config_storage_server_quota.storage_server_id", "1"),
			// 	),
			// },
			// Step 28: Tags (using config_tags)
			{
				Config: providerConfig + dependencyConfig + `
resource "hpe_morpheus_policy" "test" {
  name = "` + namePrefix + `-tags"
  description = "Tags policy using static schema"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "tags"
  }
  
  config_tags = {
    strict = true
    key = "environment"
    value = "production"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-tags"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "tags"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_tags.key", "environment"),
				),
			},
			// Step 29: User Creation (using config_create_user)
			{
				Config: providerConfig + dependencyConfig + `
resource "hpe_morpheus_policy" "test" {
  name = "` + namePrefix + `-createUser"
  description = "Create user policy using static schema"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "createUser"
  }
  
  config_create_user = {
    create_user_type = "fixed"
    create_user = true
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createUser"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createUser"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_create_user.create_user", "true"),
				),
			},
			// Step 30: User Group Creation (using config_create_user_group)
			{
				Config: providerConfig + dependencyConfig + `
resource "hpe_morpheus_policy" "test" {
  name = "` + namePrefix + `-createUserGroup"
  description = "Create user group policy using static schema"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "createUserGroup"
  }
  
  config_create_user_group = {
    user_group = "1"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createUserGroup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createUserGroup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "config_create_user_group.user_group", "1"),
				),
			},
			// Step 31: Workflow (using config_workflow)
			{
				Config: providerConfig + dependencyConfig + `
resource "hpe_morpheus_policy" "test" {
  name = "` + namePrefix + `-workflow"
  description = "Workflow policy using static schema"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "workflow"
  }
  
  config_workflow = {
    workflow_id = data.morpheus_workflow.test.id
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-workflow"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "workflow"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_policy.test", "config_workflow.workflow_id"),
				),
			},
		},
	})
}
