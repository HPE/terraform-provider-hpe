// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	datasourceprofile "github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancerprofile"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerProfileDataSourceByIdOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	// Not parallel: NSX-T load balancer profiles are children of a load balancer
	// and can be affected by concurrent LB lifecycle operations on the shared
	// integration. Run these profile tests sequentially to avoid that.
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

	profileName := acctest.RandomWithPrefix(t.Name())

	profileConfig := renderHTTPProfileResource(profileName)

	// Read the profile created above by its id.
	dataSourceConfig, err := datasourceprofile.RenderLoadBalancerProfileDataSourceByIDConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Id":             "hpe_morpheus_load_balancer_profile.http.id",
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + profileConfig + dataSourceConfig

	dsName := "data.hpe_morpheus_load_balancer_profile.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrPair(
						dsName, "id",
						"hpe_morpheus_load_balancer_profile.http", "id",
					),
					resource.TestCheckResourceAttr(dsName, "name", profileName),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileDataSourceByNameOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	// Not parallel — see TestAccMorpheusLoadBalancerProfileDataSourceByIdOk.
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

	profileName := acctest.RandomWithPrefix(t.Name())

	profileConfig := renderHTTPProfileResource(profileName)

	// Read the profile created above by name.
	dataSourceConfig, err := datasourceprofile.RenderLoadBalancerProfileDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer_profile.http.load_balancer_id",
		"Name":           profileName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + profileConfig + dataSourceConfig

	dsName := "data.hpe_morpheus_load_balancer_profile.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrPair(
						dsName, "id",
						"hpe_morpheus_load_balancer_profile.http", "id",
					),
					resource.TestCheckResourceAttr(dsName, "name", profileName),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileDataSourceNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

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

	dataSourceConfig, err := datasourceprofile.RenderLoadBalancerProfileDataSourceByNameConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.lb.id",
		"Name":           "nonexistent-profile-" + acctest.RandString(8),
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + lbConfig + dataSourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

// renderHTTPProfileResource creates a minimal HTTP profile resource config.
// This is a test helper that builds inline HCL until the resource builder
// provides a render helper.
func renderHTTPProfileResource(name string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "http" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBHttpProfile"
  config_http {
    https_redirect = true
  }
}
`, name)
}
