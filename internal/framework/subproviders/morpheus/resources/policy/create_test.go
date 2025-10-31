package policy_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
    maxMemory = 1073741824
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

// Test creating policies with different policy types
func TestAccMorpheusPolicyAllPolicyTypesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlockMixed()
	namePrefix := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

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

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: maxMemory policy
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
			// Step 2: maxCores policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxCores"),
					"policy_description": config.StringVariable("Max cores policy"),
					"policy_type_code":   config.StringVariable("maxCores"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxCores":          config.StringVariable("4"),
						"excludeContainers": config.BoolVariable(false),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxCores"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxCores"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max cores policy"),
				),
			},
			// Step 3: maxStorage policy
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
			// Step 4: deleteApproval policy
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
			// Step 5: provisionApproval policy
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
			// Step 6: reconfigureApproval policy
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
			// Step 7: workflowApproval policy - Note: this is NOT in the API spec enum, may not be valid
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
			// Step 8: createBackup policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-createBackup"),
					"policy_description": config.StringVariable("Create backup policy"),
					"policy_type_code":   config.StringVariable("createBackup"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"createBackupType": config.StringVariable("fixed"),
						"createBackup":     config.StringVariable("on"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createBackup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createBackup"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Create backup policy"),
				),
			},
			// Step 9: backupStorage policy - SKIPPED due to test framework limitation with list variables in map(any)
			// The test framework wraps list variables in an extra array layer when passed through map(any)
			// This would need a dedicated test with inline config instead of variables
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-backupStorage"),
						"policy_description": config.StringVariable("Backup storage policy"),
						"policy_type_code":   config.StringVariable("backupStorage"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"backupStorageIds": config.TupleVariable(
								config.StringVariable("1"),
							),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-backupStorage"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "backupStorage"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Backup storage policy"),
					),
				},
			*/
			// Step 10: maxPrice policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxPrice"),
					"policy_description": config.StringVariable("Max price policy"),
					"policy_type_code":   config.StringVariable("maxPrice"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxPrice":         config.IntegerVariable(100),
						"maxPriceCurrency": config.StringVariable("USD"),
						"maxPriceUnit":     config.StringVariable("month"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPrice"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPrice"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max price policy"),
				),
			},
			// Step 11: serverNaming policy - SKIPPED due to API validation issue
			// API rejects with "server naming type must be set" even when serverNamingType is provided
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-serverNaming"),
						"policy_description": config.StringVariable("Server naming policy"),
						"policy_type_code":   config.StringVariable("serverNaming"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"serverNamingType":     config.StringVariable("name"),
							"serverNamingPattern":  config.StringVariable("server-${sequence}"),
							"serverNamingConflict": config.StringVariable("on"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-serverNaming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "serverNaming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Server naming policy"),
					),
				},
			*/
			// Step 12: cypher policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-cypher"),
					"policy_description": config.StringVariable("Cypher policy"),
					"policy_type_code":   config.StringVariable("cypher"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"keyPattern": config.StringVariable("secret/*"),
						"read":       config.BoolVariable(true),
						"write":      config.BoolVariable(false),
						"update":     config.BoolVariable(false),
						"delete":     config.BoolVariable(false),
						"list":       config.BoolVariable(true),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-cypher"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "cypher"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Cypher policy"),
				),
			},
			// Step 13: delayedRemoval policy
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
			// Step 14: lifecycle policy - SKIPPED due to complex approval integration dependencies
			// When lifecycleAutoRenew is off and lifecycleExtensionsBeforeApproval > 0, requires approval integration
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-lifecycle"),
						"policy_description": config.StringVariable("Lifecycle policy"),
						"policy_type_code":   config.StringVariable("lifecycle"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"lifecycleType":                     config.StringVariable("user"),
							"lifecycleAge":                      config.StringVariable("30"),
							"lifecycleRenewal":                  config.StringVariable("7"),
							"lifecycleNotify":                   config.StringVariable("3"),
							"lifecycleMessage":                  config.StringVariable("Instance will expire"),
							"lifecycleAutoRenew":                config.StringVariable("off"),
							"lifecycleAllowExtend":              config.StringVariable("on"),
							"lifecycleExtensionsBeforeApproval": config.StringVariable("2"),
							"lifecycleHideFixed":                config.StringVariable("off"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-lifecycle"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "lifecycle"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Lifecycle policy"),
					),
				},
			*/
			// Step 15: storageShareQuota policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-storageShareQuota"),
					"policy_description": config.StringVariable("Storage share quota policy"),
					"policy_type_code":   config.StringVariable("storageShareQuota"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxStorage": config.StringVariable("10737418240"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-storageShareQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "storageShareQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Storage share quota policy"),
				),
			},
			// Step 16: hostNaming policy - SKIPPED due to API validation issue (similar to serverNaming)
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-hostNaming"),
						"policy_description": config.StringVariable("Host naming policy"),
						"policy_type_code":   config.StringVariable("hostNaming"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"hostNamingType":    config.StringVariable("name"),
							"hostNamingPattern": config.StringVariable("host-${sequence}"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-hostNaming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "hostNaming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Host naming policy"),
					),
				},
			*/
			// Step 17: naming policy - SKIPPED due to API validation issue (similar to serverNaming)
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-naming"),
						"policy_description": config.StringVariable("Naming policy"),
						"policy_type_code":   config.StringVariable("naming"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"namingType":     config.StringVariable("name"),
							"namingPattern":  config.StringVariable("instance-${sequence}"),
							"namingConflict": config.StringVariable("on"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-naming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "naming"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Naming policy"),
					),
				},
			*/
			// Step 18: maxContainers policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxContainers"),
					"policy_description": config.StringVariable("Max containers policy"),
					"policy_type_code":   config.StringVariable("maxContainers"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxContainers": config.StringVariable("10"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxContainers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxContainers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max containers policy"),
				),
			},
			// Step 19: maxHosts policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxHosts"),
					"policy_description": config.StringVariable("Max hosts policy"),
					"policy_type_code":   config.StringVariable("maxHosts"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxHosts": config.StringVariable("5"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxHosts"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxHosts"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max hosts policy"),
				),
			},
			// Step 20: maxPools policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxPools"),
					"policy_description": config.StringVariable("Max pools policy"),
					"policy_type_code":   config.StringVariable("maxPools"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxPools": config.StringVariable("3"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxPools"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxPools"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max pools policy"),
				),
			},
			// Step 21: maxPoolMembers policy
			{
				Config: providerConfig + resourceConfig,
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
			// Step 22: maxSnapshots
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
			// Step 23: maxVirtualServers policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxVirtualServers"),
					"policy_description": config.StringVariable("Max virtual servers policy"),
					"policy_type_code":   config.StringVariable("maxVirtualServers"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxVirtualServers": config.StringVariable("10"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxVirtualServers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxVirtualServers"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max virtual servers policy"),
				),
			},
			// Step 24: maxVms policy
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
			// Step 25: motd policy - SKIPPED due to API date format issue
			// API returns motd.date in format "2025-10-31 11:43:10" which SDK cannot parse as RFC3339
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-motd"),
						"policy_description": config.StringVariable("MOTD policy"),
						"policy_type_code":   config.StringVariable("motd"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"motd.title":     config.StringVariable("Welcome"),
							"motd.message":   config.StringVariable("Welcome to the system"),
							"motd.type":      config.StringVariable("info"),
							"motd._fullPage": config.BoolVariable(false),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-motd"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "motd"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "MOTD policy"),
					),
				},
			*/
			// Step 26: maxNetworks policy
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
			// Step 27: storageBucketQuota policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-storageBucketQuota"),
					"policy_description": config.StringVariable("Storage bucket quota policy"),
					"policy_type_code":   config.StringVariable("storageBucketQuota"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxStorage":        config.StringVariable("10737418240"),
						"excludeContainers": config.BoolVariable(false),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-storageBucketQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "storageBucketQuota"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Storage bucket quota policy"),
				),
			},
			// Step 28: powerSchedule policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-powerSchedule"),
					"policy_description": config.StringVariable("Power schedule policy"),
					"policy_type_code":   config.StringVariable("powerSchedule"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"powerScheduleType":      config.StringVariable("fixed"),
						"powerSchedule":          config.StringVariable("1"),
						"powerScheduleHideFixed": config.BoolVariable(false),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-powerSchedule"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "powerSchedule"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Power schedule policy"),
				),
			},
			// Step 29: maxRouters policy
			{
				Config: providerConfig + resourceConfig,
				ConfigVariables: config.Variables{
					"policy_name":        config.StringVariable(namePrefix + "-maxRouters"),
					"policy_description": config.StringVariable("Max routers policy"),
					"policy_type_code":   config.StringVariable("maxRouters"),
					"policy_config": config.ObjectVariable(map[string]config.Variable{
						"maxRouters": config.StringVariable("3"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-maxRouters"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "maxRouters"),
					resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Max routers policy"),
				),
			},
			// Step 30: shutdown policy - SKIPPED due to complex approval integration dependencies (similar to lifecycle)
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-shutdown"),
						"policy_description": config.StringVariable("Shutdown policy"),
						"policy_type_code":   config.StringVariable("shutdown"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"shutdownType":                     config.StringVariable("user"),
							"shutdownAge":                      config.StringVariable("7"),
							"shutdownRenewal":                  config.StringVariable("3"),
							"shutdownNotify":                   config.StringVariable("1"),
							"shutdownMessage":                  config.StringVariable("Instance will shutdown"),
							"shutdownAutoRenew":                config.StringVariable("on"),
							"shutdownAllowExtend":              config.StringVariable("on"),
							"shutdownExtensionsBeforeApproval": config.StringVariable("0"),
							"shutdownHideFixed":                config.StringVariable("off"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-shutdown"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "shutdown"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Shutdown policy"),
					),
				},
			*/
			// Step 31: storageServerQuota policy
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
			// Step 32: tags policy
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
			// Step 33: createUser policy
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
			// Step 34: createUserGroup policy - SKIPPED due to API validation requiring a valid user group
			/*
				{
					Config: providerConfig + resourceConfig,
					ConfigVariables: config.Variables{
						"policy_name":        config.StringVariable(namePrefix + "-createUserGroup"),
						"policy_description": config.StringVariable("Create user group policy"),
						"policy_type_code":   config.StringVariable("createUserGroup"),
						"policy_config": config.ObjectVariable(map[string]config.Variable{
							"userGroup": config.StringVariable("Administrators"),
						}),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "name", namePrefix+"-createUserGroup"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "policy_type.code", "createUserGroup"),
						resource.TestCheckResourceAttr("hpe_morpheus_policy.test", "description", "Create user group policy"),
					),
				},
			*/
			// Step 35: workflow policy - uses workflow created in resourceConfig
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
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
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
