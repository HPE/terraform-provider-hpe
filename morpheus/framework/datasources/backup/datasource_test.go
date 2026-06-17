// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backup_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	databackup "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_job"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupinstance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
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
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func backupDataSourceChecks(name string) resource.TestCheckFunc {
	const dataSourceName = "data.hpe_morpheus_backup.example"

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			dataSourceName, "id",
			"hpe_morpheus_backup_instance.example", "id",
		),
		resource.TestCheckResourceAttr(dataSourceName, "name", name),
		resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceName, "instance.id", "26"),
		resource.TestCheckResourceAttr(dataSourceName, "storage_provider.id", "1"),
		resource.TestCheckResourceAttrPair(
			dataSourceName, "job.id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttrSet(dataSourceName, "backup_type.code"),
		resource.TestCheckResourceAttrSet(dataSourceName, "location_type"),
	)
}

// TestAccMorpheusFindBackupById creates a backup (a backup job and an instance
// backup) and reads it back through the data source by ID.
func TestAccMorpheusFindBackupById(t *testing.T) {
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

	backupJobConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	backupInstanceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackup.RenderBackupDataSourceByIDConfig(t, map[string]string{
		"Id": "hpe_morpheus_backup_instance.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create the backup job + instance backup and read it back by ID.
			{
				Config: providerConfig + backupJobConfig + backupInstanceConfig + dataSourceConfig,
				Check:  backupDataSourceChecks(name),
			},
			// Tear down the backup first (backup_job destroyed as a side-effect
			// of the backup_instance destroy)...
			{
				Config:             providerConfig + backupJobConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only.
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccMorpheusFindBackupByName creates a backup (a backup job and an
// instance backup) and reads it back through the data source by name.
func TestAccMorpheusFindBackupByName(t *testing.T) {
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

	backupJobConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	backupInstanceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackup.RenderBackupDataSourceByNameConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create the backup job + instance backup first so it can be found
			// by name.
			{
				Config: providerConfig + backupJobConfig + backupInstanceConfig,
			},
			// Read the backup back by name.
			{
				Config: providerConfig + backupJobConfig + backupInstanceConfig + dataSourceConfig,
				Check:  backupDataSourceChecks(name),
			},
			// Tear down the backup first (backup_job destroyed as a side-effect
			// of the backup_instance destroy)...
			{
				Config:             providerConfig + backupJobConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only.
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccMorpheusFindBackupNotFound asserts that looking up a non-existent
// backup by name surfaces the no-backup-found error.
func TestAccMorpheusFindBackupNotFound(t *testing.T) {
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

	config := providerConfig + `
data "hpe_morpheus_backup" "example" {
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackup.ErrorNoBackupFound),
			},
		},
	})
}

// TestAccMorpheusFindBackupNoSearchAttrs asserts that omitting both id and name
// surfaces the no-valid-search-terms error.
func TestAccMorpheusFindBackupNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_backup" "example" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackup.ErrorNoValidSearchTerms),
			},
		},
	})
}

// TestAccMorpheusFindBackupBothSearchAttrs asserts that providing both id and
// name fails the ConflictsWith validation during plan.
func TestAccMorpheusFindBackupBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_backup" "example" {
  id   = 1
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackup.ErrorRunningPreApply),
			},
		},
	})
}
