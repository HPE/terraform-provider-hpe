// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerMonitorResourceNsxtExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

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

	resourceName := "hpe_morpheus_load_balancer_monitor.nsxt"

	monitorConfig, err := loadbalancermonitor.RenderLoadBalancerMonitorNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	config := providerConfig + lbConfig + monitorConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "name", monitorName),
		resource.TestCheckResourceAttr(resourceName, "description", "An NSX-T HTTP health check monitor"),
		resource.TestCheckResourceAttr(resourceName, "monitor_type", "http"),
		resource.TestCheckResourceAttr(resourceName, "monitor_interval", "5"),
		resource.TestCheckResourceAttr(resourceName, "monitor_timeout", "15"),
		resource.TestCheckResourceAttr(resourceName, "monitor_destination", "/"),
		resource.TestCheckResourceAttr(resourceName, "fall_count", "3"),
		resource.TestCheckResourceAttr(resourceName, "rise_count", "3"),
		resource.TestCheckResourceAttr(resourceName, "alias_port", "8080"),
		resource.TestCheckResourceAttr(resourceName, "send_type", "GET"),
		resource.TestCheckResourceAttr(resourceName, "send_version", "HTTP_VERSION_1_1"),
		resource.TestCheckResourceAttr(resourceName, "receive_code", "200"),
		resource.TestCheckResourceAttr(resourceName, "data_length", "0"),
		resource.TestCheckResourceAttr(resourceName, "max_retry", "3"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "load_balancer_id"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ResourceName:      resourceName,
				ImportStateVerifyIgnore: []string{
					"monitor_password_wo_version",
					"config",
					// The API accepts these on create/update but never
					// returns them, so an imported resource cannot know
					// them. See preserveUnreturnedFields in resource_read.go.
					"receive_data",
					"data_length",

					"write_only_revision",
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerMonitorResourceNsxtUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

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
	updatedName := monitorName + "-updated"

	resourceName := "hpe_morpheus_load_balancer_monitor.nsxt"

	createConfig, err := loadbalancermonitor.RenderLoadBalancerMonitorNsxtConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render create config: %s", err)
	}

	createConfig = providerConfig + lbConfig + createConfig

	updateConfig, err := loadbalancermonitor.RenderLoadBalancerMonitorNsxtConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"Name":               updatedName,
		"Description":        "Updated NSX-T monitor",
		"MonitorTimeout":     "30",
		"FallCount":          "5",
		"RiseCount":          "5",
		"MonitorDestination": "/status",
		"ReceiveCode":        "201",
	})
	if err != nil {
		t.Fatalf("failed to render update config: %s", err)
	}

	updateConfig = providerConfig + lbConfig + updateConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", monitorName),
					resource.TestCheckResourceAttr(resourceName, "monitor_timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "fall_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "rise_count", "3"),
				),
				PlanOnly: false,
			},
			{
				Config:             createConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated NSX-T monitor"),
					resource.TestCheckResourceAttr(resourceName, "monitor_timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "monitor_destination", "/status"),
					resource.TestCheckResourceAttr(resourceName, "fall_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "rise_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "receive_code", "201"),
				),
				PlanOnly: false,
			},
			{
				Config:             updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

// importStateIDFunc returns an ImportStateIdFunc that produces the composite
// import ID "loadBalancerId.monitorId" from the resource state.
func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		lbID := rs.Primary.Attributes["load_balancer_id"]
		id := rs.Primary.Attributes["id"]

		return lbID + "." + id, nil
	}
}

// TestAccMorpheusLoadBalancerMonitorResourceJuniperSdnHttpExampleOk creates an
// HTTP health-check monitor on an SDN (juniper-sdn-lb) load balancer. The LB
// and the juniperSdn network it sits on are created per-test as prerequisites.
func TestAccMorpheusLoadBalancerMonitorResourceJuniperSdnHttpExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)

	// The monitor needs a load balancer of its own, and any load balancer
	// delete on Morpheus 9.1.x damages the other load balancers' pool and
	// virtual-server records (MORPH-18174), so this test is held with the
	// load balancer test until that fix lands.
	if os.Getenv(testhelpers.EnvSdnLoadBalancer) == "" {
		t.Skip("SDN load balancer monitor test held for MORPH-18174; set " + testhelpers.EnvSdnLoadBalancer + "=1 to run")
	}

	f := testhelpers.RequireSdnFixture(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	lbName := name[0:16] + name[len(name)-16:]

	network := `
resource "hpe_morpheus_network" "sdn" {
  name              = "` + lbName + `-net"
  cloud_id          = ` + f.CloudID + `
  group_id          = ` + f.GroupID + `
  type_code         = "juniperSdn"
  network_server_id = ` + f.NetworkServerID + `
  zone_pool_id      = ` + f.ClusterPoolID + `
  cidr              = "10.93.0.0/24"
  gateway           = "10.93.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
`

	lbConfig, err := loadbalancer.RenderLoadBalancerJuniperSdnConfig(t, map[string]string{
		"Name":              lbName,
		"CloudId":           f.CloudID,
		"NetworkServerId":   f.NetworkServerID,
		"ResourcePoolId":    f.ClusterPoolID,
		"VirtualNetworkUid": "hpe_morpheus_network.sdn.unique_id",
		"VirtualIp":         "10.93.0.100",
		"SourcePort":        "80",
		"DestinationPort":   "8080",
		"Protocol":          "TCP",
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	monitorName := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_load_balancer_monitor.sdn"

	monitorConfig, err := loadbalancermonitor.RenderLoadBalancerMonitorJuniperSdnConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.sdn.id",
		"Name":           monitorName,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	config := testhelpers.ProviderBlock() + network + lbConfig + monitorConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", monitorName),
					resource.TestCheckResourceAttr(resourceName, "monitor_type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "config.http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "config.http_uri", "/health"),
					resource.TestCheckResourceAttr(resourceName, "config.http_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "config.http_return_code", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "load_balancer_id"),
				),
			},
			{
				// The create step's config.http_* checks are satisfied by the
				// planned values echoed into state (overlayPlanConfig); the
				// import reads them back from the appliance and so proves the
				// read-back mapping.
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ResourceName:      resourceName,
				ImportStateVerifyIgnore: []string{
					// Accepted on create but never returned by the API; see
					// preserveUnreturnedFields in resource_read.go.
					"receive_data",
					"data_length",
					"write_only_revision",
				},
			},
		},
	})
}
