// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy_test

//go:generate ../../../../../../bin/render example-id.tf.tmpl Id 99
//go:generate ../../../../../../bin/render example-name.tf.tmpl Name "Example name"

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/policy/consts"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusPolicyDataSourceFindByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	policyName := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name     = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "test" {
  name                     = "` + policyName + `"
  description              = "Test policy for datasource"
  associated_resource_type = "Group"
  associated_resource_id   = hpe_morpheus_group.test.id
  enabled                  = true

  policy_type = {
    code = "maxMemory"
  }

  config = {
    maxMemory = "8"
  }
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name.tf.tmpl",
		"Name", policyName)
	if err != nil {
		t.Fatalf("Failed to render example: %v", err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test",
			"name",
			policyName,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyDataSourceFindById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	policyName := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name     = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "test" {
  name                     = "` + policyName + `"
  description              = "Test policy for datasource"
  associated_resource_type = "Group"
  associated_resource_id   = hpe_morpheus_group.test.id
  enabled                  = true

  policy_type = {
    code = "maxMemory"
  }

  config = {
    maxMemory = "8"
  }
}
`

	dataSourceConfig, err := testhelpers.RenderExample(
		t, "example-id.tf.tmpl", "Id", "hpe_morpheus_policy.test.id")
	if err != nil {
		t.Fatalf("Failed to render example: %v", err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test",
			"name",
			policyName,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusPolicyDataSourceNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_policy" "test" {
        name = "nonexistent-policy"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_policy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoPolicyFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusPolicyDataSourceNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_policy" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_policy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidPolicyTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusPolicyDataSourceBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_policy" "test" {
        id = "1"
        name = "testpolicy"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_policy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusPolicyDataSourceVerifyAllAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	policyName := acctest.RandomWithPrefix(t.Name())
	policyDescription := "Comprehensive test policy with all config fields"
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	providerConfig := testhelpers.ProviderBlock()

	// Create a single policy resource with ALL config fields combined for testing
	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name     = "` + groupName + `"
  location = "test"
}

# Policy with all possible config fields (not a valid real-world policy, but tests datasource)
resource "hpe_morpheus_policy" "test_all_attrs" {
  name                     = "` + policyName + `"
  description              = "` + policyDescription + `"
  associated_resource_type = "Group"
  associated_resource_id   = hpe_morpheus_group.test.id
  enabled                  = true

  policy_type = {
    code = "maxMemory"
  }

  # Config block testing all 31 policy type configurations
  config = {
    # 1. Approval fields (deleteApproval, provisionApproval, reconfigureApproval, workflowApproval policies)
    accountIntegrationId = "1"
    flowId               = "1"
    workflowId           = "1"
    workflowType         = "flow"
    
    # 2. Backup Storage fields (backupStorage policy)
    backupStorageIds = [1, 2]
    
    # 3. Backup Creation fields (createBackup policy)
    createBackupType = "fixed"
    createBackup     = true
    
    # 4. User Creation fields (createUser policy)
    createUserType = "fixed"
    createUser     = true
    
    # 5. User Group Creation fields (createUserGroup policy)
    userGroup = "1"
    
    # 6. Cypher Access fields (cypher policy)
    keyPattern = "secret/*"
    read       = true
    write      = true
    update     = true
    delete     = false
    list       = true
    
    # 7. Budget fields (maxPrice policy)
    maxPrice         = 1000
    maxPriceCurrency = "USD"
    maxPriceUnit     = "month"
    
    # 8. Max Memory fields (maxMemory policy)
    maxMemory         = "16"
    excludeContainers = "on"
    
    # 9. Max Cores fields (maxCores policy)
    maxCores = "16"
    
    # 10. Delayed Removal fields (delayedRemoval policy)
    removalAge = "30"
    
    # 11. Lifecycle fields (lifecycle policy)
    lifecycleType                     = "fixed"
    lifecycleAge                      = "30"
    lifecycleRenewal                  = "7"
    lifecycleNotify                   = "3"
    lifecycleMessage                  = "Instance will expire soon"
    lifecycleAutoRenew                = "on"
    lifecycleAllowExtend              = "on"
    lifecycleExtensionsBeforeApproval = "2"
    lifecycleHideFixed                = false
    # flowId                            = "1"  # Set in approval section
    lifecycleWorkflowId               = "1"
    # workflowType                      = "flow"  # Set in approval section
    # accountIntegrationId              = "1"  # Set in approval section
    
    # 12. Hostname fields (hostNaming policy)
    hostNamingType    = "fixed"
    hostNamingPattern = "host-$${groupCode}-$${sequence}"
    
    # 13. Naming fields (naming policy)
    namingType     = "fixed"
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}"
    namingConflict = true
    
    # 14. Max Containers fields (maxContainers policy)
    maxContainers = "50"
    
    # 15. Max Hosts fields (maxHosts policy)
    maxHosts = "10"
    
    # 16. Max Networks fields (maxNetworks policy)
    maxNetworks = "15"
    
    # 17. Max Pool Members fields (maxPoolMembers policy)
    maxPoolMembers = "12"
    
    # 18. Max Pools fields (maxPools policy)
    maxPools = "8"
    
    # 19. Max Routers fields (maxRouters policy)
    maxRouters = "5"
    
    # 20. Max Snapshots fields (maxSnapshots policy)
    maxSnapshots = "5"
    
    # 21. Max Storage fields (maxStorage policy)
    maxStorage = "500"
    
    # 22. Max Virtual Servers fields (maxVirtualServers policy)
    maxVirtualServers = "25"
    
    # 23. Max VMs fields (maxVms policy)
    maxVms = "20"
    
    # 24. MOTD fields (motd policy)
    "motd.title"     = "Welcome"
    "motd.message"   = "Welcome to the platform"
    "motd.type"      = "info"
    "motd.fullPage" = "off"
    "motd.date"      = "2025-10-31 14:53:07"
    
    # 25. Power Schedule fields (powerSchedule policy)
    powerScheduleType      = "fixed"
    powerSchedule          = "1"
    powerScheduleHideFixed = false
    
    # 26. Required Network fields (requiredNetwork policy)
    requiredNetworks = [100, 200]
    
    # 27. Server naming fields (serverNaming policy)
    serverNamingType     = "fixed"
    serverNamingPattern  = "server-$${groupCode}-$${sequence}"
    serverNamingConflict = true
    
    # 28. Shutdown fields (shutdown policy)
    shutdownType                     = "fixed"
    shutdownAge                      = "30"
    shutdownRenewal                  = "7"
    shutdownNotify                   = "3"
    shutdownMessage                  = "Instance will shutdown soon"
    shutdownAutoRenew                = "on"
    shutdownAllowExtend              = "on"
    shutdownExtensionsBeforeApproval = "2"
    shutdownHideFixed                = false
    # accountIntegrationId             = "1"  # Set in approval section
    # flowId                           = "1"  # Set in approval section
    shutdownWorkflowId               = "1"
    # workflowType                     = "flow"  # Set in approval section
    
    # 29. Storage Server Quota fields (storageServerQuota policy)
    storageServerId = "1"
    # maxStorage      = "500"  # Set in Max Storage section
    
    # 30. Tags fields (tags policy)
    strict      = true
    key         = "environment"
    value       = "production"
    valueListId = ""
    
    # 31. Workflow fields (workflow policy)
    workflowId = "1"
  }
}
`

	dataSourceConfig := `
data "hpe_morpheus_policy" "test_all_attrs" {
  name = "` + policyName + `"
}
`

	// Validation checks for all 31 policy type configurations using static config fields
	checks := []resource.TestCheckFunc{
		// Basic policy attributes
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "name", policyName),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "description", policyDescription),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "associated_resource_type", "Group"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "policy_type.code", "maxMemory"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.test_all_attrs", "id"),

		// 1. Approval config fields (all attributes: account_integration_id, flow_id, workflow_id, workflow_type)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_approval.account_integration_id", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_approval.flow_id", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_approval.workflow_id", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_approval.workflow_type", "flow"),

		// 2. Backup Storage config fields (all attributes: backup_storage_ids)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_backup_storage.backup_storage_ids.#", "2"),

		// 3. Backup Creation config fields (all attributes: create_backup, create_backup_type)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_create_backup.create_backup_type", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_create_backup.create_backup", "true"),

		// 4. User Creation config fields (all attributes: create_user, create_user_type)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_create_user.create_user_type", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_create_user.create_user", "true"),

		// 5. User Group Creation config fields (all attributes: user_group)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_create_user_group.user_group", "1"),

		// 6. Cypher Access config fields (all attributes: key_pattern, read, write, update, delete, list)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.key_pattern", "secret/*"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.read", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.write", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.update", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.delete", "false"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_cypher.list", "true"),

		// 7. Budget config fields (all attributes: max_price, max_price_currency, max_price_unit)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_price.max_price", "1000"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_max_price.max_price_currency", "USD"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_price.max_price_unit", "month"),

		// 8. Max Memory config fields (all attributes: max_memory, exclude_containers)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_memory.max_memory", "16"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_max_memory.exclude_containers", "true"),

		// 9. Max Cores config fields (all attributes: max_cores, exclude_containers)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_cores.max_cores", "16"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.test_all_attrs", "config_max_cores.exclude_containers"),

		// 10. Delayed Removal config fields (all attributes: removal_age)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_delayed_removal.removal_age", "30"),

		// 11. Lifecycle config fields (all attributes)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_type", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_age", "30"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_renewal", "7"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_notify", "3"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_message", "Instance will expire soon"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_auto_renew", "true"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_allow_extend", "true"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_extensions_before_approval", "2"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_hide_fixed", "false"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.account_integration_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.flow_id", "1"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.lifecycle_workflow_id", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_lifecycle.workflow_type", "flow"),

		// 12. Hostname config fields (all attributes: host_naming_type, host_naming_pattern)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_host_naming.host_naming_type", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs",
			"config_host_naming.host_naming_pattern",
			"host-${groupCode}-${sequence}"),

		// 13. Naming config fields (all attributes: naming_type, naming_pattern, naming_conflict)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_naming.naming_type", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_naming.naming_pattern", "vm-${groupCode}-${type}-${sequence}"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_naming.naming_conflict", "true"),

		// 14. Max Containers config fields (all attributes: max_containers)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_max_containers.max_containers", "50"),

		// 15. Max Hosts config fields (all attributes: max_hosts)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_hosts.max_hosts", "10"),

		// 16. Max Networks config fields (all attributes: max_networks)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_networks.max_networks", "15"),

		// 17. Max Pool Members config fields (all attributes: max_pool_members)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_max_pool_members.max_pool_members", "12"),

		// 18. Max Pools config fields (all attributes: max_pools)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_pools.max_pools", "8"),

		// 19. Max Routers config fields (all attributes: max_routers)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_routers.max_routers", "5"),

		// 20. Max Snapshots config fields (all attributes: max_snapshots)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_snapshots.max_snapshots", "5"),

		// 21. Max Storage config fields (all attributes: max_storage, exclude_containers)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_storage.max_storage", "500"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.test_all_attrs", "config_max_storage.exclude_containers"),

		// 22. Max Virtual Servers config fields (all attributes: max_virtual_servers)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_max_virtual_servers.max_virtual_servers", "25"),

		// 23. Max VMs config fields (all attributes: max_vms)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_max_vms.max_vms", "20"),

		// 24. MOTD config fields (all attributes: motdtitle, motdmessage, motdtype, motd_fullpage, motddate)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_motd.motdtitle", "Welcome"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_motd.motdmessage", "Welcome to the platform"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_motd.motdtype", "info"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.test_all_attrs", "config_motd.motd_fullpage"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_motd.motddate", "2025-10-31 14:53:07"),

		// 25. Power Schedule config fields (all attributes: power_schedule_type, power_schedule, power_schedule_hide_fixed)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_power_schedule.power_schedule_type", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_power_schedule.power_schedule", "1"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_power_schedule.power_schedule_hide_fixed", "false"),

		// 26. Required Network config fields (all attributes: required_networks)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_required_network.required_networks.#", "2"),

		// 27. Server naming config fields (all attributes: server_naming_type, server_naming_pattern, server_naming_conflict)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_server_naming.server_naming_type", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs",
			"config_server_naming.server_naming_pattern",
			"server-${groupCode}-${sequence}"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_server_naming.server_naming_conflict", "true"),

		// 28. Shutdown config fields (all attributes)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_type", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_age", "30"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_renewal", "7"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_notify", "3"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_message", "Instance will shutdown soon"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_auto_renew", "true"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_allow_extend", "true"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_extensions_before_approval", "2"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_hide_fixed", "false"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.account_integration_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.flow_id", "1"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.shutdown_workflow_id", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_shutdown.workflow_type", "flow"),

		// 29. Storage Server Quota config fields (all attributes: storage_server_id, max_storage)
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config_storage_server_quota.storage_server_id", "1"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policy.test_all_attrs", "config_storage_server_quota.max_storage"),

		// 30. Tags config fields (all attributes: key, value, strict, value_list_id)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_tags.strict", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_tags.key", "environment"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_tags.value", "production"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_tags.value_list_id", ""),

		// 31. Workflow config fields (all attributes: workflow_id)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config_workflow.workflow_id", "1"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}

// Test datasource can find policies scoped to different resource types (Group, Cloud, User, Role)
func TestAccMorpheusPolicyDataSourceResourceTypesOk(t *testing.T) {
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
`

	policyNameGroup := acctest.RandomWithPrefix(t.Name() + "-group-policy")
	policyNameCloud := acctest.RandomWithPrefix(t.Name() + "-cloud-policy")
	policyNameRole := acctest.RandomWithPrefix(t.Name() + "-role-policy")
	policyNameUser := acctest.RandomWithPrefix(t.Name() + "-user-policy")

	// Create policies for each resource type
	resourceConfig := `
resource "hpe_morpheus_policy" "group_policy" {
  name = "` + policyNameGroup + `"
  description = "Group-scoped policy"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = "1073741824"
  }
}

resource "hpe_morpheus_policy" "cloud_policy" {
  name = "` + policyNameCloud + `"
  description = "Cloud-scoped policy"
  associated_resource_type = "Cloud"
  associated_resource_id = hpe_morpheus_cloud.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = "10"
  }
}

resource "hpe_morpheus_policy" "role_policy" {
  name = "` + policyNameRole + `"
  description = "Role-scoped policy"
  associated_resource_type = "Role"
  associated_resource_id = hpe_morpheus_role.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = "10"
  }
}

resource "hpe_morpheus_policy" "user_policy" {
  name = "` + policyNameUser + `"
  description = "User-scoped policy"
  associated_resource_type = "User"
  associated_resource_id = hpe_morpheus_user.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = "10"
  }
}
`

	// Data sources to find each policy
	dataSourceConfig := `
data "hpe_morpheus_policy" "group_policy" {
  name = hpe_morpheus_policy.group_policy.name
}

data "hpe_morpheus_policy" "cloud_policy" {
  name = hpe_morpheus_policy.cloud_policy.name
}

data "hpe_morpheus_policy" "role_policy" {
  name = hpe_morpheus_policy.role_policy.name
}

data "hpe_morpheus_policy" "user_policy" {
  name = hpe_morpheus_policy.user_policy.name
}
`

	checksGroup := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "name", policyNameGroup),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "description", "Group-scoped policy"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "associated_resource_type", "Group"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.group_policy", "associated_resource_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "policy_type.code", "maxMemory"),
		// Check group-specific attributes
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.group_policy", "group.id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.group_policy", "group.name", groupName),
		// Verify other resource type attributes are not set
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.group_policy", "cloud.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.group_policy", "user.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.group_policy", "role.id"),
	}

	checksCloud := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "name", policyNameCloud),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "description", "Cloud-scoped policy"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "associated_resource_type", "Cloud"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.cloud_policy", "associated_resource_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "policy_type.code", "maxMemory"),
		// Check cloud-specific attributes
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.cloud_policy", "cloud.id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.cloud_policy", "cloud.name", cloudName),
		// Verify other resource type attributes are not set
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.cloud_policy", "group.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.cloud_policy", "user.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.cloud_policy", "role.id"),
	}

	checksRole := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "name", policyNameRole),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "description", "Role-scoped policy"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "associated_resource_type", "Role"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.role_policy", "associated_resource_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "policy_type.code", "maxMemory"),
		// Check role-specific attributes
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.role_policy", "role.id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.role_policy", "role.authority", roleName),
		// Verify other resource type attributes are not set
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.role_policy", "group.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.role_policy", "cloud.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.role_policy", "user.id"),
	}

	checksUser := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "name", policyNameUser),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "description", "User-scoped policy"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "associated_resource_type", "User"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.user_policy", "associated_resource_id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "policy_type.code", "maxMemory"),
		// Check user-specific attributes
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.user_policy", "user.id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.user_policy", "user.username", userName),
		// Verify other resource type attributes are not set
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.user_policy", "group.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.user_policy", "cloud.id"),
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_policy.user_policy", "role.id"),
	}

	allChecks := append(checksGroup, checksCloud...)
	allChecks = append(allChecks, checksRole...)
	allChecks = append(allChecks, checksUser...)
	checkFn := resource.ComposeAggregateTestCheckFunc(
		allChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:             providerConfig + dependencyConfig + resourceConfig + dataSourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
