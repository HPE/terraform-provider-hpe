// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	datasourcemonitor "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	resourcemonitor "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerMonitorDataSourceByIdOk(t *testing.T) {
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

	monitorName := acctest.RandomWithPrefix(t.Name())

	monitorConfig, err := resourcemonitor.RenderLoadBalancerMonitorNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render monitor config: %s", err)
	}

	dataSourceConfig, err := datasourcemonitor.RenderLoadBalancerMonitorDataSourceByIDConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Id":             "hpe_morpheus_load_balancer_monitor.nsxt.id",
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + monitorConfig + dataSourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_monitor.example",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_load_balancer_monitor.example",
						"name",
						monitorName,
					),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerMonitorDataSourceByNameOk(t *testing.T) {
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

	monitorName := acctest.RandomWithPrefix(t.Name())

	monitorConfig, err := resourcemonitor.RenderLoadBalancerMonitorNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render monitor config: %s", err)
	}

	dataSourceConfig, err := datasourcemonitor.RenderLoadBalancerMonitorDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + monitorConfig + dataSourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_monitor.example",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_load_balancer_monitor.example",
						"name",
						monitorName,
					),
				),
			},
		},
	})
}
