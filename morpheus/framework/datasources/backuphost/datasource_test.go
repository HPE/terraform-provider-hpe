// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	databackuphost "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backuphost"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backuphost"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupjob"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
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

func backupHostDataSourceChecks(name string) resource.TestCheckFunc {
	const dataSourceName = "data.hpe_morpheus_backup_host.example"

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			dataSourceName, "id",
			"hpe_morpheus_backup_host.example", "id",
		),
		resource.TestCheckResourceAttr(dataSourceName, "name", name),
		resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceName, "host.id", "139"),
		resource.TestCheckResourceAttr(dataSourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttrPair(
			dataSourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
	)
}

// TestAccMorpheusFindBackupHostById creates a host backup and reads it back
// through the data source by ID.
func TestAccMorpheusFindBackupHostById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
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
		"HostId":            "139",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackuphost.RenderBackupHostDataSourceByIDConfig(t, map[string]string{
		"Id": "hpe_morpheus_backup_host.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + resourceConfig + dataSourceConfig,
				Check:  backupHostDataSourceChecks(name),
			},
			// Tear down the backup_host (backup_job destroyed as side-effect)...
			{
				Config:             providerConfig + dependencyConfig,
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

// TestAccMorpheusFindBackupHostByName creates a host backup and reads it back
// through the data source by name.
func TestAccMorpheusFindBackupHostByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
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
		"HostId":            "139",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackuphost.RenderBackupHostDataSourceByNameConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Create the backup job + host backup first.
			{
				Config: providerConfig + dependencyConfig + resourceConfig,
			},
			// Read the backup back by name.
			{
				Config: providerConfig + dependencyConfig + resourceConfig + dataSourceConfig,
				Check:  backupHostDataSourceChecks(name),
			},
			// Tear down the backup_host (backup_job destroyed as side-effect)...
			{
				Config:             providerConfig + dependencyConfig,
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

// TestAccMorpheusFindBackupHostNotFound asserts that looking up a non-existent
// backup host by name surfaces the no-backup-found error.
func TestAccMorpheusFindBackupHostNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
data "hpe_morpheus_backup_host" "example" {
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackuphost.ErrorNoBackupFound),
			},
		},
	})
}

// TestAccMorpheusFindBackupHostNoSearchAttrs asserts that omitting both id and
// name surfaces the no-valid-search-terms error.
func TestAccMorpheusFindBackupHostNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_backup_host" "example" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackuphost.ErrorNoValidSearchTerms),
			},
		},
	})
}
