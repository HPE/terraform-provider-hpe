// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	dsnetwork "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/network"
)

func TestAccMorpheusDataSourceNetworkSubnetExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

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
