// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitches_test

// Real appliance acceptance tests for the virtual switches data source.
//
// These require a live Morpheus >= 9.1.0 appliance with a virtual-switch-capable
// HVM cluster and the capabilities.VirtualSwitch capability. They skip cleanly
// when the environment is not provisioned and are NOT expected to run until such
// an appliance is arranged (MORPH-17745).
//
// Required environment:
//
//	TF_ACC_CAPABILITIES   must include "virtual_switch"
//	TF_ACC_HVM_CLUSTER_ID the numeric id of a virtual-switch-capable HVM cluster

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusVirtualSwitchesDataSource(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VirtualSwitch)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	clusterID := os.Getenv("TF_ACC_HVM_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("TF_ACC_HVM_CLUSTER_ID must be set")
	}

	if _, err := strconv.ParseInt(clusterID, 10, 64); err != nil {
		t.Fatalf("TF_ACC_HVM_CLUSTER_ID must be numeric: %s", err)
	}

	ds := "data.hpe_morpheus_virtual_switches.example"

	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_virtual_switches" "example" {
  cluster_id = ` + clusterID + `
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(ds, "virtual_switches.#"),
				),
			},
		},
	})
}
