// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// These tests cover the plan stability of hpe_morpheus_instance across in-place
// updates.
//
// The plugin framework marks every computed attribute that is null in the
// configuration as unknown whenever anything on the resource changes, so without
// a plan modifier pass a small edit shows unrelated attributes as "(known after
// apply)". The resource restores those values from prior state, but only when
// nothing that could affect them has changed — pinning a value the API
// legitimately changes during an update would raise "Provider produced
// inconsistent result after apply", which is a worse failure than a noisy plan.
//
// Each test therefore asserts one half of that contract:
//
//   - a metadata-only edit restores the computed values, so the plan stays
//     reviewable;
//   - an edit that reconfigures a collection leaves that collection's computed
//     values unknown, and the apply still succeeds.

const instanceResourceName = "hpe_morpheus_instance.example"

// TestAccMorpheusInstanceResourceTagUpdateDoesNotChurnPlan verifies that editing
// only tags does not show unrelated computed attributes as "(known after
// apply)".
//
// This is the regression test for the reported behaviour: changing one
// attribute churned connection_info, labels and every computed field of
// network_interfaces and volumes, producing a diff too large to review.
//
// Interface and volume identity are asserted because they are populated by the
// API on create and are stable across an unrelated edit, so a known value in the
// plan proves the restore happened.
func TestAccMorpheusInstanceResourceTagUpdateDoesNotChurnPlan(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	config, err := instance.RenderInstanceConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	taggedConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":         name,
		"MultipleTags": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
			},
			{
				Config: providerConfig + taggedConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							instanceResourceName, plancheck.ResourceActionUpdate,
						),
						// Restored from prior state rather than left unknown.
						plancheck.ExpectKnownValue(
							instanceResourceName,
							tfjsonpath.New("volumes").AtSliceIndex(0).AtMapKey("id"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							instanceResourceName,
							tfjsonpath.New("network_interfaces").AtSliceIndex(0).AtMapKey("id"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							instanceResourceName,
							tfjsonpath.New("network_interfaces").AtSliceIndex(0).AtMapKey("name"),
							knownvalue.NotNull(),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "tags.#", "5"),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourceResizeKeepsInterfaceAndVolumeIdentity verifies
// that resizing an instance does not churn interface and volume identity, and
// that the resize still applies cleanly.
//
// This is the scenario reported from the field: a memory change also showed
// interface and volume attributes as "(known after apply)". Addresses are
// deliberately not asserted as known — a resize that cannot be performed hot
// restarts the instance, so connection_info can legitimately change and is left
// unknown by design.
//
// The apply completing is itself part of the assertion: if identity were pinned
// while the API changed it, the step would fail with "Provider produced
// inconsistent result after apply".
func TestAccMorpheusInstanceResourceResizeKeepsInterfaceAndVolumeIdentity(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	config, err := instance.RenderInstanceConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	resizedConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":      name,
		"MaxMemory": "2048",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
			},
			{
				Config: providerConfig + resizedConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							instanceResourceName, plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectKnownValue(
							instanceResourceName,
							tfjsonpath.New("volumes").AtSliceIndex(0).AtMapKey("id"),
							knownvalue.NotNull(),
						),
						plancheck.ExpectKnownValue(
							instanceResourceName,
							tfjsonpath.New("network_interfaces").AtSliceIndex(0).AtMapKey("id"),
							knownvalue.NotNull(),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						instanceResourceName, "service_plan_options.max_memory", "2048",
					),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourceVolumeAddLeavesVolumeIdentityUnknown verifies
// the other half of the contract: when the volumes themselves are reconfigured,
// volume identity is *not* restored from prior state.
//
// Restoring it would be wrong twice over. The API may assign a new id, and
// because these are list elements a positional restore could pin the previous
// element's values onto a different volume. Either would surface as "Provider
// produced inconsistent result after apply", so the value is deliberately left
// unknown and the apply must still succeed.
func TestAccMorpheusInstanceResourceVolumeAddLeavesVolumeIdentityUnknown(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	singleVolumeConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":         name,
		"SingleVolume": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	twoVolumeConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + singleVolumeConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "volumes.#", "1"),
				),
			},
			{
				Config: providerConfig + twoVolumeConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							instanceResourceName, plancheck.ResourceActionUpdate,
						),
						// The volumes were reconfigured, so identity must not be
						// taken from prior state.
						plancheck.ExpectUnknownValue(
							instanceResourceName,
							tfjsonpath.New("volumes").AtSliceIndex(1).AtMapKey("id"),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "volumes.#", "2"),
				),
			},
		},
	})
}
