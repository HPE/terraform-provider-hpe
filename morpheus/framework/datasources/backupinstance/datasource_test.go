// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	databackupinstance "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backupinstance"
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

func backupInstanceDataSourceChecks(name string) resource.TestCheckFunc {
	const dataSourceName = "data.hpe_morpheus_backup_instance.example"

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			dataSourceName, "id",
			"hpe_morpheus_backup_instance.example", "id",
		),
		resource.TestCheckResourceAttr(dataSourceName, "name", name),
		resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceName, "instance_id", "26"),
		resource.TestCheckResourceAttr(dataSourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttrPair(
			dataSourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttrSet(dataSourceName, "container_id"),
	)
}

// TestAccMorpheusFindBackupInstanceById creates an instance backup and reads
// it back through the data source by ID.
func TestAccMorpheusFindBackupInstanceById(t *testing.T) {
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

	resourceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackupinstance.RenderBackupInstanceDataSourceByIDConfig(t, map[string]string{
		"Id": "hpe_morpheus_backup_instance.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + resourceConfig + dataSourceConfig,
				Check:  backupInstanceDataSourceChecks(name),
			},
			// Tear down the backup_instance (backup_job destroyed as side-effect)...
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

// TestAccMorpheusFindBackupInstanceByName creates an instance backup and reads
// it back through the data source by name.
func TestAccMorpheusFindBackupInstanceByName(t *testing.T) {
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

	resourceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := databackupinstance.RenderBackupInstanceDataSourceByNameConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Create the backup job + instance backup first.
			{
				Config: providerConfig + dependencyConfig + resourceConfig,
			},
			// Read the backup back by name.
			{
				Config: providerConfig + dependencyConfig + resourceConfig + dataSourceConfig,
				Check:  backupInstanceDataSourceChecks(name),
			},
			// Tear down the backup_instance (backup_job destroyed as side-effect)...
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

// TestAccMorpheusFindBackupInstanceNotFound asserts that looking up a
// non-existent backup instance by name surfaces the no-backup-found error.
func TestAccMorpheusFindBackupInstanceNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
data "hpe_morpheus_backup_instance" "example" {
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackupinstance.ErrorNoBackupFound),
			},
		},
	})
}

// TestAccMorpheusFindBackupInstanceNoSearchAttrs asserts that omitting both id
// and name surfaces the no-valid-search-terms error.
func TestAccMorpheusFindBackupInstanceNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_backup_instance" "example" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackupinstance.ErrorNoValidSearchTerms),
			},
		},
	})
}
