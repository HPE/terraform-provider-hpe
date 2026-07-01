// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	datasourcepool "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancerpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	resourcepool "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancerpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerPoolDataSourceByIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	// Not parallel: a standalone NSX-T load balancer pool stays "offline"
	// (unrealized) and can be pruned by a Morpheus/NSX-T inventory sync that
	// fires when a concurrent test tears down its load balancer on the shared
	// integration. Run these pool tests sequentially to avoid that.

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

	poolName := acctest.RandomWithPrefix(t.Name())

	poolConfig, err := resourcepool.RenderLoadBalancerPoolNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           poolName,
	})
	if err != nil {
		t.Fatalf("failed to render pool config: %s", err)
	}

	// Read the pool created above by its id. Referencing the resource's id
	// makes the data source depend on (and resolve after) the pool.
	dataSourceConfig, err := datasourcepool.RenderLoadBalancerPoolDataSourceByIDConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Id":             "hpe_morpheus_load_balancer_pool.nsxt.id",
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + poolConfig + dataSourceConfig

	dsName := "data.hpe_morpheus_load_balancer_pool.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrPair(
						dsName, "id",
						"hpe_morpheus_load_balancer_pool.nsxt", "id",
					),
					resource.TestCheckResourceAttr(dsName, "name", poolName),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerPoolDataSourceByNameOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	// Not parallel: a standalone NSX-T load balancer pool stays "offline"
	// (unrealized) and can be pruned by a Morpheus/NSX-T inventory sync that
	// fires when a concurrent test tears down its load balancer on the shared
	// integration. Run these pool tests sequentially to avoid that.

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

	poolName := acctest.RandomWithPrefix(t.Name())

	poolConfig, err := resourcepool.RenderLoadBalancerPoolNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           poolName,
	})
	if err != nil {
		t.Fatalf("failed to render pool config: %s", err)
	}

	// Read the pool created above by name. Sourcing load_balancer_id from the
	// pool resource makes the data source depend on (and resolve after) it.
	dataSourceConfig, err := datasourcepool.RenderLoadBalancerPoolDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer_pool.nsxt.load_balancer_id",
		"Name":           poolName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + poolConfig + dataSourceConfig

	dsName := "data.hpe_morpheus_load_balancer_pool.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrPair(
						dsName, "id",
						"hpe_morpheus_load_balancer_pool.nsxt", "id",
					),
					resource.TestCheckResourceAttr(dsName, "name", poolName),
				),
			},
		},
	})
}
