// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	backuphost "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_host"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_job"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusBackupHostResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Backup) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := backuphost.RenderBackupHostConfig(t, map[string]string{
		"Name":              name,
		"HostId":            "139",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_backup_host.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "host_id", "139"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hostname"),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:            providerConfig + dependencyConfig + resourceConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
			// Do this to destroy the backup_host (backup_job destroyed as side-effect of backup_host destroy)...
			{
				Config:             providerConfig + dependencyConfig,
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
	if capabilities.Missing(t, capabilities.Backup) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := backuphost.RenderBackupHostConfig(t, map[string]string{
		"Name":              name,
		"HostId":            "139",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "updated"
	updateConfig := `
resource "hpe_morpheus_backup_host" "example" {
  name                = "` + updatedName + `"
  host_id             = 139
  job_id              = hpe_morpheus_backup_job.example.id
  backup_type_code    = "fileBackup"
  path                = "/etc/hosts"
  storage_provider_id = 2
  enabled             = false
}
`

	resourceName := "hpe_morpheus_backup_host.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "host_id", "139"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hostname"),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "host_id", "139"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "backup_type_code", "fileBackup"),
		resource.TestCheckResourceAttr(resourceName, "path", "/etc/hosts"),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "2"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
	)
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + dependencyConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + dependencyConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// Do this to destroy the backup_host (backup_job destroyed as side-effect of backup_host destroy)...
			{
				Config:             providerConfig + dependencyConfig,
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
