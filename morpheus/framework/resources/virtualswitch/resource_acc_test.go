// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch_test

// Real appliance acceptance tests for the virtual switch resource.
//
// These require a live Morpheus >= 9.1.0 appliance with an HVM (mvm-cluster)
// cluster on a virtual-switch-capable (2.0+) layout, hosts with spare NICs, and
// the capabilities.VirtualSwitch capability. They are gated so they skip
// cleanly when the environment is not provisioned; they are NOT expected to run
// until such an appliance is arranged (MORPH-17745).
//
// Required environment:
//
//	TF_ACC_CAPABILITIES   must include "virtual_switch"
//	TF_ACC_HVM_CLUSTER_ID the numeric id of a virtual-switch-capable HVM cluster
//	TF_ACC_VSW_NIC1       the primary uplink NIC name (e.g. "eth2")
//
// Optional:
//
//	TF_ACC_VSW_NIC2       a secondary uplink NIC name for bonded uplinks

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/virtualswitch"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusVirtualSwitchResource(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VirtualSwitch)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	clusterID := os.Getenv("TF_ACC_HVM_CLUSTER_ID")
	nic1 := os.Getenv("TF_ACC_VSW_NIC1")
	if clusterID == "" || nic1 == "" {
		t.Skip("TF_ACC_HVM_CLUSTER_ID and TF_ACC_VSW_NIC1 must be set")
	}

	if _, err := strconv.ParseInt(clusterID, 10, 64); err != nil {
		t.Fatalf("TF_ACC_HVM_CLUSTER_ID must be numeric: %s", err)
	}

	nic2 := os.Getenv("TF_ACC_VSW_NIC2")

	// Switch names are unique per cluster and limited to 16 characters, so
	// use a short prefix plus a random suffix (t.Name() alone would truncate
	// to the same fixed string on every run and collide with any leftover).
	// The prefix is what the virtual switch sweeper matches, since the
	// suite-wide "TestAccMorpheus" prefix cannot fit.
	// The create uses a 16-character name, the API's maximum, so the
	// schema's limit is exercised; the rename is shorter.
	name := virtualswitch.TestNamePrefix + acctest.RandString(11) // "tfacc" + 11 = 16
	renamed := virtualswitch.TestNamePrefix + acctest.RandString(5) + "-2"

	// A NIC bond needs two NICs; with the optional second NIC unset the
	// switch is created with single-NIC uplinks and no bond_mode.
	bondMode := "active-backup"
	if nic2 == "" {
		bondMode = ""
	}

	resourceName := "hpe_morpheus_virtual_switch.example"

	createCfg, err := virtualswitch.RenderVirtualSwitchExample(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
		"BondMode":  bondMode,
		"Mtu":       "9000",
		"Nic1":      nic1,
		"Nic2":      nic2,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	updateCfg, err := virtualswitch.RenderVirtualSwitchExample(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      renamed,
		"BondMode":  bondMode,
		"Mtu":       "9000",
		"Nic1":      nic1,
		"Nic2":      nic2,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
				),
			},
			{
				Config: providerConfig + updateCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", renamed),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// segments is ignored for the segment map KEYS only: import
				// synthesises keys (e.g. "data-nfs_100") from each segment's
				// identity, which differ from the user's config keys (e.g. "nfs").
				// The keys are cosmetic -- re-keying after import is a no-op (see
				// computeUpdatePlan identity matching) -- but ImportStateVerify
				// compares map keys literally, so the ignore must stay.
				ImportStateVerifyIgnore: []string{"timeouts", "segments"},
			},
		},
	})
}
