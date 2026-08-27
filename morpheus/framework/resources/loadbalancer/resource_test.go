// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerResourceHAProxyExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer, capabilities.NetworkLoadBalancerHAProxy)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	// This resource only allows name to be 32 characters maximum.
	name = name[0:16] + name[len(name)-16:]

	config, err := loadbalancer.RenderLoadBalancerHAProxyConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.haproxy"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config",
					"config_haproxy",
					"cloud_id",
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerResourceHAProxyGenericExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer, capabilities.NetworkLoadBalancerHAProxy)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	// This resource only allows name to be 32 characters maximum.
	name = name[0:16] + name[len(name)-16:]

	config, err := loadbalancer.RenderLoadBalancerHAProxyGenericConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.haproxy_generic"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "haproxyContainer"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config",
					"config_haproxy",
					"cloud_id",
					"type_code",
					"permissions.groups",
				},
			},
		},
	})
}

// Test validation: permissions.all conflicts with permissions.groups
func TestAccMorpheusLoadBalancerResourceValidationPermissionsConflict(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	name = name[0:16] + name[len(name)-16:]

	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name = "` + name + `"

  config = {}

  permissions = {
    all    = true
    groups = [1, 2]
  }
}
`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile(`(?i)attribute "permissions.groups" cannot be specified when`),
			},
		},
	})
}

// TestMorpheusLoadBalancerWriteOnlyConfigParsing verifies that write-only attributes
// (group_id, network_server_id) are accepted in config and produce a valid plan.
// The framework nullifies write-only values in the plan, so Create must source them
// from req.Config — this test confirms the schema accepts the attributes correctly.
func TestAccMorpheusLoadBalancerWriteOnlyConfigParsing(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name              = "test-wo-attrs"
  group_id          = 1
  network_server_id = 42
  config_nsxt = {
    admin_state   = true
    log_level     = "INFO"
    size          = "SMALL"
    tier1_gateway = "/infra/tier-1s/test"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestMorpheusLoadBalancerWriteOnlyVersionCompanions verifies that the version
// companion attributes (group_id_version, network_server_id_version) are accepted
// in config alongside the write-only attributes they track.
func TestAccMorpheusLoadBalancerWriteOnlyVersionCompanions(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name                      = "test-wo-companions"
  group_id                  = 1
  group_id_version          = 1
  network_server_id         = 42
  network_server_id_version = 1
  config_nsxt = {
    admin_state   = true
    log_level     = "INFO"
    size          = "SMALL"
    tier1_gateway = "/infra/tier-1s/test"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
