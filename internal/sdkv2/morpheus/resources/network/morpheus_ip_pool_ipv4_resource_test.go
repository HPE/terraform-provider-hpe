// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/network"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusIpPoolIpv4ResourceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := network.RenderIPPoolIPv4Config(t, name, map[string]string{
		"Name": "\"" + name + "\"",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_ip_pool_ipv4.tf_example_ipv4_pool",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_ip_pool_ipv4.tf_example_ipv4_pool",
			"ip_range.0.starting_address",
			"192.168.1.1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_ip_pool_ipv4.tf_example_ipv4_pool",
			"ip_range.0.ending_address",
			"192.168.1.10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_ip_pool_ipv4.tf_example_ipv4_pool",
			"ip_range.1.starting_address",
			"10.0.0.1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_ip_pool_ipv4.tf_example_ipv4_pool",
			"ip_range.1.ending_address",
			"10.0.0.10",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
