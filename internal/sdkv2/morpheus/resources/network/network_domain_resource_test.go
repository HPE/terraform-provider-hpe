// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/network"
)

func TestAccMorpheusNetworkDomainExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := network.RenderNetworkDomainConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"active",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"description",
			"Terraform example network domain",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"name",
			strings.ToLower(name),
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"public_zone",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"tenant_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_domain.example",
			"visibility",
			"private",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}
