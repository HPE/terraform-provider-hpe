// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	datasourcevs "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	resourcevs "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerVirtualServerDataSourceByIdOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkLoadBalancer) {
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

	vipName := acctest.RandomWithPrefix(t.Name())

	// Use port 443 with protocol "http" intentionally: this verifies the API
	// stores the configured port as-is rather than normalizing based on protocol.
	vsConfig, err := resourcevs.RenderLoadBalancerVirtualServerNsxtMinimalConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"VipName":            vipName,
		"Description":        "datasource test vs",
		"VipAddress":         "10.0.0.201",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"ApplicationProfile": "13",
	})
	if err != nil {
		t.Fatalf("failed to render vs config: %s", err)
	}

	dsConfig, err := datasourcevs.RenderLoadBalancerVirtualServerDataSourceByIDConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Id":             "hpe_morpheus_load_balancer_virtual_server.nsxt_minimal.id",
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + vsConfig + dsConfig
	dsName := "data.hpe_morpheus_load_balancer_virtual_server.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "vip_name", vipName),
					resource.TestCheckResourceAttr(dsName, "description", "datasource test vs"),
					resource.TestCheckResourceAttr(dsName, "vip_address", "10.0.0.201"),
					resource.TestCheckResourceAttr(dsName, "vip_port", "443"),
					resource.TestCheckResourceAttr(dsName, "vip_protocol", "http"),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerDataSourceByNameOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkLoadBalancer) {
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

	vipName := acctest.RandomWithPrefix(t.Name())

	// Use port 443 with protocol "http" intentionally: this verifies the API
	// stores the configured port as-is rather than normalizing based on protocol.
	vsConfig, err := resourcevs.RenderLoadBalancerVirtualServerNsxtMinimalConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"VipName":            vipName,
		"Description":        "datasource test vs",
		"VipAddress":         "10.0.0.202",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"ApplicationProfile": "13",
	})
	if err != nil {
		t.Fatalf("failed to render vs config: %s", err)
	}

	dsConfig, err := datasourcevs.RenderLoadBalancerVirtualServerDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"VipName":        vipName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + vsConfig + dsConfig
	dsName := "data.hpe_morpheus_load_balancer_virtual_server.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "vip_name", vipName),
					resource.TestCheckResourceAttr(dsName, "description", "datasource test vs"),
					resource.TestCheckResourceAttr(dsName, "vip_address", "10.0.0.202"),
					resource.TestCheckResourceAttr(dsName, "vip_port", "443"),
					resource.TestCheckResourceAttr(dsName, "vip_protocol", "http"),
				),
			},
		},
	})
}
