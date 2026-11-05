// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupinstance"
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

func TestAccMorpheusBackupInstanceResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	instanceName := testhelpers.BackupInstanceName(t)
	bucketName := testhelpers.BackupBucketName(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	fixtureConfig := testhelpers.BackupFixtureBlock(instanceName, "", bucketName, "")
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        testhelpers.BackupInstanceIDRef,
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": testhelpers.BackupBucketIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_backup_instance.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttrPair(resourceName, "instance_id", testhelpers.BackupInstanceDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		// MORPH-16232: storage_provider_id must survive apply even though the
		// API may return a null storage provider.
		resource.TestCheckResourceAttrPair(resourceName, "storage_provider_id", testhelpers.BackupBucketDataSource, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttrSet(resourceName, "backup_type_code"),
		resource.TestCheckResourceAttrSet(resourceName, "container_id"),
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
				Config:            providerConfig + fixtureConfig + dependencyConfig + resourceConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
			// Do this to destroy the backup_instance (backup_job destroyed as side-effect of backup_instance destroy)...
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

func TestAccMorpheusBackupInstanceResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	instanceName := testhelpers.BackupInstanceName(t)
	bucketName, otherBucketName := testhelpers.BackupBucketNames(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	fixtureConfig := testhelpers.BackupFixtureBlock(instanceName, "", bucketName, otherBucketName)
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        testhelpers.BackupInstanceIDRef,
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": testhelpers.BackupBucketIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "updated"
	updateConfig := `
resource "hpe_morpheus_backup_instance" "example" {
  name                = "` + updatedName + `"
  instance_id         = ` + testhelpers.BackupInstanceIDRef + `
  job_id              = hpe_morpheus_backup_job.example.id
  storage_provider_id = ` + testhelpers.BackupOtherBucketIDRef + `
  enabled             = false
}
`

	resourceName := "hpe_morpheus_backup_instance.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttrPair(resourceName, "instance_id", testhelpers.BackupInstanceDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		// MORPH-16232: storage_provider_id must survive apply/update.
		resource.TestCheckResourceAttrPair(resourceName, "storage_provider_id", testhelpers.BackupBucketDataSource, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttrPair(resourceName, "instance_id", testhelpers.BackupInstanceDataSource, "id"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
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
			// Do this to destroy the backup_instance (backup_job destroyed as side-effect of backup_instance destroy)...
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
