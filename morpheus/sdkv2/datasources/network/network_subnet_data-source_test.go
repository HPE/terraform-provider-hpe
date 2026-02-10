// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsnetwork "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusDataSourceNetworkSubnetExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := "VM Network"

	var dependenciesConfig string

	datasourceConfig, err := dsnetwork.RenderNetworkSubnetConfig(t, map[string]string{
		"Name":      "\"" + name + "\"",
		"NetworkId": "2",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_subnet.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_subnet.example",
			"network_id",
			"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + "\n" + dependenciesConfig + "\n" + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
