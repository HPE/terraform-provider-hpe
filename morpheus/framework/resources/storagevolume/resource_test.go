// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storagevolume"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageVolumeResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := storagevolume.RenderStorageVolumeConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_storage_volume.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_storage_volume.example", "type_code", "hpealletraMPLUN"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_storage_volume.example",
			},
		},
	})
}

func TestAccMorpheusStorageVolumeResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := storagevolume.RenderStorageVolumeConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := storagevolume.RenderStorageVolumeConfig(t, map[string]string{
		"Name": updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_storage_volume.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceCompleteOk exercises the max_storage
// (size) path, which the example-based tests do not cover. max_storage is
// expressed in GiB; before the units fix the Morpheus API rejected the request
// with "The size must be between 1 and 65536 GiB" because the value was treated
// as bytes (MORPH-13021). The PlanOnly step additionally confirms that the
// bytes->GiB read-back round-trips without drift, and the import step confirms
// that the computed_optional type_id/type_code pair produces a clean import.
func TestAccMorpheusStorageVolumeResourceCompleteOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_storage_volume.example"

	resourceConfig, err := storagevolume.RenderStorageVolumeCompleteConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "hpealletraMPLUN"),
					// max_storage is in GiB.
					resource.TestCheckResourceAttr(resourceName, "max_storage", "10"),
					// type_id is computed_optional and is populated from the API.
					resource.TestCheckResourceAttrSet(resourceName, "type_id"),
				),
			},
			{
				// A no-op re-plan confirms the bytes->GiB read-back round-trips
				// without drift.
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceWriteOnlyConfigOk exercises the typed
// write-only config_alletramp_bmaas block (buildCreateConfig) end to end, and
// confirms that incrementing config_alletramp_bmaas_wo_version forces a
// replacement (the write-only block cannot be diffed in state, so the version
// trigger drives the change).
func TestAccMorpheusStorageVolumeResourceWriteOnlyConfigOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_storage_volume.alletramp_bmaas"

	createConfig, err := storagevolume.RenderStorageVolumeAlletraMPBMaaSConfig(t, map[string]string{
		"Name":      name,
		"WoVersion": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	bumpConfig, err := storagevolume.RenderStorageVolumeAlletraMPBMaaSConfig(t, map[string]string{
		"Name":      name,
		"WoVersion": "2",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectReplace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "hpealletraMPLUN"),
					// config_alletramp_bmaas is write-only and not persisted in state.
					resource.TestCheckNoResourceAttr(resourceName, "config_alletramp_bmaas.datastore_id"),
					resource.TestCheckResourceAttr(resourceName, "config_alletramp_bmaas_wo_version", "1"),
				),
			},
			{
				// Bumping the write-only version trigger forces a replacement.
				Config:           providerConfig + bumpConfig,
				ConfigPlanChecks: expectReplace,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "config_alletramp_bmaas_wo_version", "2"),
				),
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceRequiresServerOrGroup verifies the schema
// rejects a volume with neither storage_server_id nor storage_group_id at plan
// time (MORPH-12939), instead of failing apply with a generic API
// "error saving volume". The error is raised during config validation, so no
// storage backend is required.
func TestAccMorpheusStorageVolumeResourceRequiresServerOrGroup(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_volume" "test" {
  name      = "tf-acc-no-server"
  type_code = "hpealletraMPLUN"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)[Aa]t least one attribute out of.*storage_server_id`),
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceConfigExportConflict verifies the
// config_alletramp_bmaas block rejects setting both compute_server_id and
// instance_ids (mutually exclusive export targets) at plan time.
func TestAccMorpheusStorageVolumeResourceConfigExportConflict(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_volume" "test" {
  name              = "tf-acc-export-conflict"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  config_alletramp_bmaas = {
    datastore_id      = 5
    compute_server_id = 10
    instance_ids      = [7, 8]
  }
  config_alletramp_bmaas_wo_version = 1
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination|cannot be specified when`),
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceInstanceIDsRequiresShared verifies that
// instance_ids (a multi-attach export target) requires shared = true, and is
// rejected at plan time when shared is false.
func TestAccMorpheusStorageVolumeResourceInstanceIDsRequiresShared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_volume" "test" {
  name              = "tf-acc-instance-ids-shared"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  config_alletramp_bmaas = {
    datastore_id = 5
    instance_ids = [7, 8]
    shared       = false
  }
  config_alletramp_bmaas_wo_version = 1
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)shared is required with instance_ids`),
			},
		},
	})
}

// TestAccMorpheusStorageVolumeResourceComputeServerForbidsShared verifies that
// compute_server_id (a single-attach export target) is incompatible with
// shared = true, and is rejected at plan time.
func TestAccMorpheusStorageVolumeResourceComputeServerForbidsShared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_volume" "test" {
  name              = "tf-acc-compute-server-shared"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  config_alletramp_bmaas = {
    datastore_id      = 5
    compute_server_id = 10
    shared            = true
  }
  config_alletramp_bmaas_wo_version = 1
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)shared conflicts with compute_server_id`),
			},
		},
	})
}
