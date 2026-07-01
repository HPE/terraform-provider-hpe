// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume_test

import (
	"os"
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
