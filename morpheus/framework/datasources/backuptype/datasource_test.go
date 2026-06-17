// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuptype_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	databackuptype "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backuptype"
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

// TestAccMorpheusFindBackupTypeById reads a platform backup type by ID.
func TestAccMorpheusFindBackupTypeById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := databackuptype.RenderBackupTypeDataSourceByIDConfig(t, map[string]string{
		"Id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	const dataSourceName = "data.hpe_morpheus_backup_type.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dataSourceName, "id", "1"),
		resource.TestCheckResourceAttrSet(dataSourceName, "name"),
		resource.TestCheckResourceAttrSet(dataSourceName, "code"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checks,
			},
		},
	})
}

// TestAccMorpheusFindBackupTypeByName reads a platform backup type by name.
func TestAccMorpheusFindBackupTypeByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	const backupTypeName = "File Backup"

	dataSourceConfig, err := databackuptype.RenderBackupTypeDataSourceByNameConfig(t, map[string]string{
		"Name": backupTypeName,
	})
	if err != nil {
		t.Fatal(err)
	}

	const dataSourceName = "data.hpe_morpheus_backup_type.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dataSourceName, "name", backupTypeName),
		resource.TestCheckResourceAttrSet(dataSourceName, "id"),
		resource.TestCheckResourceAttrSet(dataSourceName, "code"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checks,
			},
		},
	})
}

// TestAccMorpheusFindBackupTypeNotFound asserts that looking up a non-existent
// backup type by name surfaces the no-backup-type-found error.
func TestAccMorpheusFindBackupTypeNotFound(t *testing.T) {
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
data "hpe_morpheus_backup_type" "example" {
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackuptype.ErrorNoBackupTypeFound),
			},
		},
	})
}

// TestAccMorpheusFindBackupTypeNoSearchAttrs asserts that omitting both id and
// name surfaces the no-valid-search-terms error.
func TestAccMorpheusFindBackupTypeNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_backup_type" "example" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackuptype.ErrorNoValidSearchTerms),
			},
		},
	})
}

// TestAccMorpheusFindBackupTypeBothSearchAttrs asserts that providing both id
// and name fails the ConflictsWith validation during plan.
func TestAccMorpheusFindBackupTypeBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_backup_type" "example" {
  id   = 1
  name = "______"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(databackuptype.ErrorRunningPreApply),
			},
		},
	})
}
