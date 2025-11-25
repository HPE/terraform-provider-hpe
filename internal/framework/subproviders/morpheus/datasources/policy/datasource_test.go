// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy_test

//go:generate go run ../../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../../cmd/render example-name.tf.tmpl Name "Example name"

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/policy/consts"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
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

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

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

  # Config block testing all 36 policy type configurations
  config = {
    # 1. Lifecycle fields (lifecycle policy)
    lifecycleType                     = "fixed"
    lifecycleAge                      = "30"
    lifecycleRenewal                  = "7"
    lifecycleNotify                   = "3"
    lifecycleMessage                  = "Instance will expire soon"
    lifecycleAutoRenew                = "on"
    lifecycleAllowExtend              = "on"
    lifecycleExtensionsBeforeApproval = "2"
    lifecycleHideFixed                = false
    
    # 2. Naming fields (naming policy)
    namingType     = "fixed"
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}"
    namingConflict = true
    
    # 3. Hostname fields (hostNaming policy)
    hostNamingType    = "fixed"
    hostNamingPattern = "host-$${groupCode}-$${type}-$${sequence}"
    
    # 4. Server naming fields (serverNaming policy)
    serverNamingType     = "fixed"
    serverNamingPattern  = "server-$${groupCode}-$${type}-$${sequence}"
    serverNamingConflict = true
    
    # 5. Max Memory fields (maxMemory policy)
    maxMemory         = "16"
    excludeContainers = "on"
    
    # 6. Max Storage fields (maxStorage policy)
    maxStorage = "500"
    
    # 7. Max Cores fields (maxCores policy)
    maxCores = "16"
    
    # 8. Budget fields (maxPrice policy)
    maxPrice         = "1000"
    maxPriceCurrency = "USD"
    maxPriceUnit     = "month"
    
    # 9. Shutdown fields (shutdown policy)
    shutdownType                     = "fixed"
    shutdownAge                      = "30"
    shutdownRenewal                  = "7"
    shutdownNotify                   = "3"
    shutdownMessage                  = "Instance will shutdown soon"
    shutdownAutoRenew                = "on"
    shutdownAllowExtend              = "on"
    shutdownExtensionsBeforeApproval = "2"
    shutdownHideFixed                = false
    
    # 10. Backup Creation fields (createBackup policy)
    createBackupType = "fixed"
    createBackup     = true
    
    # 11. User Creation fields (createUser policy)
    createUserType = "fixed"
    createUser     = true
    
    # 12. Tags fields (tags policy)
    strict      = true
    key         = "environment"
    value       = "production"
    valueListId = ""
    
    # 13. Delayed Removal fields (delayedRemoval policy)
    removalAge = "30"
    
    # 14. Max VMs fields (maxVms policy)
    maxVms = "20"
    
    # 15. Max Containers fields (maxContainers policy)
    maxContainers = "50"
    
    # 16. Max Hosts fields (maxHosts policy)
    maxHosts = "10"
    
    # 17. Max Networks fields (maxNetworks policy)
    maxNetworks = "15"
    
    # 18. Max Routers fields (maxRouters policy)
    maxRouters = "5"
    
    # 19. Max Pools fields (maxPools policy)
    maxPools = "8"
    
    # 20. Max Pool Members fields (maxPoolMembers policy)
    maxPoolMembers = "12"
    
    # 21. Max Virtual Servers fields (maxVirtualServers policy)
    maxVirtualServers = "25"
    
    # 22. Max Snapshots fields (maxSnapshots policy)
    maxSnapshots = "5"
    
    # 23. Cypher Access fields (cypher policy)
    keyPattern = "secret/*"
    read       = true
    write      = true
    update     = true
    delete     = false
    list       = true
    
    # 24. Power Schedule fields (powerSchedule policy)
    powerScheduleType      = "fixed"
    powerSchedule          = "1"
    powerScheduleHideFixed = false
    
    # 25. MOTD fields (motd policy)
    "motd.title"     = "Welcome"
    "motd.message"   = "Welcome to the platform"
    "motd.type"      = "info"
    "motd._fullPage" = "off"
    
    # 26. User Group Creation fields (createUserGroup policy)
    userGroup = "1"
    
    # 27-30. Approval fields (deleteApproval, provisionApproval, reconfigureApproval, workflowApproval policies)
    accountIntegrationId = "1"
    
    # 31. Workflow fields (workflow policy)
    workflowId = "1"
    
    # 32. Storage Server Quota fields (storageServerQuota policy)
    storageServerId = "1"
    
    # 33. Storage Bucket Quota fields (storageBucketQuota policy)
    # Uses maxStorage (already defined above)
    
    # 34. Storage Share Quota fields (storageShareQuota policy)
    # Uses maxStorage (already defined above)
    
    # 35. Backup Storage fields (backupStorage policy)
    backupStorageIds = [1, 2]
    
    # 36. Required Network fields (requiredNetwork policy)
    requiredNetworks = [100, 200]
  }
}
`

	dataSourceConfig := `
data "hpe_morpheus_policy" "test_all_attrs" {
  name = "` + policyName + `"
}
`

	// Validation checks for all 36 policy type configurations
	checks := []resource.TestCheckFunc{
		// Basic policy attributes
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "name", policyName),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "description", policyDescription),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "enabled", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "associated_resource_type", "Group"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "policy_type.code", "maxMemory"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_policy.test_all_attrs", "id"),

		// 1. Lifecycle config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleType", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleAge", "30"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleRenewal", "7"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleNotify", "3"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleMessage", "Instance will expire soon"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleAutoRenew", "on"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleAllowExtend", "on"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleExtensionsBeforeApproval", "2"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.lifecycleHideFixed", "false"),

		// 2. Naming config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.namingType", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.namingPattern", "vm-${groupCode}-${type}-${sequence}"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.namingConflict", "true"),

		// 3. Hostname config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.hostNamingType", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.hostNamingPattern", "host-${groupCode}-${type}-${sequence}"),

		// 4. Server naming config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.serverNamingType", "fixed"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.serverNamingPattern", "server-${groupCode}-${type}-${sequence}"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.serverNamingConflict", "true"),

		// 5. Max Memory config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxMemory", "16"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.excludeContainers", "on"),

		// 6. Max Storage config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxStorage", "500"),

		// 7. Max Cores config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxCores", "16"),

		// 8. Budget config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxPrice", "1000"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxPriceCurrency", "USD"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxPriceUnit", "month"),

		// 9. Shutdown config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownType", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownAge", "30"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownRenewal", "7"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownNotify", "3"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.shutdownMessage", "Instance will shutdown soon"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownAutoRenew", "on"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownAllowExtend", "on"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.shutdownExtensionsBeforeApproval", "2"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.shutdownHideFixed", "false"),

		// 10. Backup Creation config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.createBackupType", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.createBackup", "true"),

		// 11. User Creation config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.createUserType", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.createUser", "true"),

		// 12. Tags config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.strict", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.key", "environment"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.value", "production"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.valueListId", ""),

		// 13. Delayed Removal config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.removalAge", "30"),

		// 14. Max VMs config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxVms", "20"),

		// 15. Max Containers config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxContainers", "50"),

		// 16. Max Hosts config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxHosts", "10"),

		// 17. Max Networks config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxNetworks", "15"),

		// 18. Max Routers config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxRouters", "5"),

		// 19. Max Pools config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxPools", "8"),

		// 20. Max Pool Members config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxPoolMembers", "12"),

		// 21. Max Virtual Servers config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxVirtualServers", "25"),

		// 22. Max Snapshots config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.maxSnapshots", "5"),

		// 23. Cypher Access config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.keyPattern", "secret/*"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.read", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.write", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.update", "true"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.delete", "false"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.list", "true"),

		// 24. Power Schedule config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.powerScheduleType", "fixed"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.powerSchedule", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.powerScheduleHideFixed", "false"),

		// 25. MOTD config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.motd.title", "Welcome"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policy.test_all_attrs", "config.motd.message", "Welcome to the platform"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.motd.type", "info"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.motd._fullPage", "off"),

		// 26. User Group Creation config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.userGroup", "1"),

		// 27-30. Approval config fields (deleteApproval, provisionApproval, reconfigureApproval, workflowApproval)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.accountIntegrationId", "1"),

		// 31. Workflow config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.workflowId", "1"),

		// 32. Storage Server Quota config fields
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.storageServerId", "1"),

		// 35. Backup Storage config fields (array)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.backupStorageIds.#", "2"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.backupStorageIds.0", "1"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.backupStorageIds.1", "2"),

		// 36. Required Network config fields (array)
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.requiredNetworks.#", "2"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.requiredNetworks.0", "100"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_policy.test_all_attrs", "config.requiredNetworks.1", "200"),
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
