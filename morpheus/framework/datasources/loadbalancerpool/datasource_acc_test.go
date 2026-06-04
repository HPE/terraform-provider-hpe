// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	datasourcepool "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancerpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerPoolDataSourceByIdOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:min(32, len(lbName))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	// NOTE: Pool resource creation is not yet available; this test will be
	// updated once the resource CRUD is implemented. For now the data source
	// code compiles and the test structure is correct.
	dataSourceConfig, err := datasourcepool.RenderLoadBalancerPoolDataSourceByIDConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Id":             "1",
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + dataSourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_pool.example",
						"id",
					),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerPoolDataSourceByNameOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:min(32, len(lbName))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	poolName := "test-pool"

	dataSourceConfig, err := datasourcepool.RenderLoadBalancerPoolDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           poolName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + dataSourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_pool.example",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_load_balancer_pool.example",
						"name",
						poolName,
					),
				),
			},
		},
	})
}
