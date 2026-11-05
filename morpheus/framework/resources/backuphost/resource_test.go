// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backuphost"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupjob"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusBackupHostResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	serverName := testhelpers.BackupServerName(t)
	bucketName := testhelpers.BackupBucketName(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	fixtureConfig := testhelpers.BackupFixtureBlock("", serverName, bucketName, "")
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := backuphost.RenderBackupHostConfig(t, map[string]string{
		"Name":              name,
		"HostId":            testhelpers.BackupServerIDRef,
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": testhelpers.BackupBucketIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_backup_host.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttrPair(resourceName, "host_id", testhelpers.BackupServerDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hostname"),
		resource.TestCheckResourceAttrPair(resourceName, "storage_provider_id", testhelpers.BackupBucketDataSource, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fixtureConfig + dependencyConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + fixtureConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:                  providerConfig + fixtureConfig + dependencyConfig + resourceConfig,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"write_only_revision"},
				ResourceName:            resourceName,
			},
			// Do this to destroy the backup_host (backup_job destroyed as side-effect of backup_host destroy)...
			{
				Config:             providerConfig + fixtureConfig + dependencyConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMorpheusBackupHostResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	serverName := testhelpers.BackupServerName(t)
	bucketName, otherBucketName := testhelpers.BackupBucketNames(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	fixtureConfig := testhelpers.BackupFixtureBlock("", serverName, bucketName, otherBucketName)
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := backuphost.RenderBackupHostConfig(t, map[string]string{
		"Name":              name,
		"HostId":            testhelpers.BackupServerIDRef,
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": testhelpers.BackupBucketIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "updated"
	updateConfig := `
resource "hpe_morpheus_backup_host" "example" {
  name                = "` + updatedName + `"
  host_id             = ` + testhelpers.BackupServerIDRef + `
  job_id              = hpe_morpheus_backup_job.example.id
  backup_type_code    = "fileBackup"
  path                = "/etc/hosts"
  storage_provider_id = ` + testhelpers.BackupOtherBucketIDRef + `
  enabled             = false
}
`

	resourceName := "hpe_morpheus_backup_host.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttrPair(resourceName, "host_id", testhelpers.BackupServerDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hostname"),
		resource.TestCheckResourceAttrPair(resourceName, "storage_provider_id", testhelpers.BackupBucketDataSource, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttrPair(resourceName, "host_id", testhelpers.BackupServerDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hosts"),
		resource.TestCheckResourceAttrPair(resourceName, "storage_provider_id", testhelpers.BackupOtherBucketDataSource, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
	)
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fixtureConfig + dependencyConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + fixtureConfig + dependencyConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + fixtureConfig + dependencyConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// Do this to destroy the backup_host (backup_job destroyed as side-effect of backup_host destroy)...
			{
				Config:             providerConfig + fixtureConfig + dependencyConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
