// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerProfileResourceHttpExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.http"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "http" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  description      = "acceptance test http profile"
  service_type     = "LBHttpProfile"

  config_http = {
    http_idle_timeout    = 15
    request_header_size  = 2048
    response_header_size = 4096
    response_timeout     = 60
    https_redirect       = true
    x_forwarded_for      = "INSERT"
  }

  tags = [
    {
      name  = "env"
      value = "test"
    },
    {
      name  = "app"
      value = "web"
    },
  ]
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", profileName),
		resource.TestCheckResourceAttr(resourceName, "service_type", "LBHttpProfile"),
		resource.TestCheckResourceAttr(resourceName, "config_http.https_redirect", "true"),
		resource.TestCheckResourceAttr(resourceName, "config_http.x_forwarded_for", "INSERT"),
		resource.TestCheckResourceAttr(resourceName, "config_http.request_header_size", "2048"),
		resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				// The config block is reconstructed from the API on import (see
				// read_config.go), so the imported resource is complete and plans
				// cleanly. It is still excluded from ImportStateVerify because the
				// API echoes server-applied config defaults (e.g.
				// x_forwarded_for=INSERT) that the applied state does not carry
				// where the user omitted them, so the imported and applied states
				// legitimately differ on those fields.
				ImportStateVerifyIgnore: []string{"config_http"},
				ResourceName:            resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					lbID := rs.Primary.Attributes["load_balancer_id"]
					id := rs.Primary.Attributes["id"]

					return lbID + "/" + id, nil
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceHttpRedirectOff(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.http_off"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "http_off" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBHttpProfile"

  config_http = {
    https_redirect   = false
    redirect_address = "https://example.com"
  }
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "config_http.https_redirect", "false"),
		resource.TestCheckResourceAttr(resourceName, "config_http.redirect_address", "https://example.com"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceCookiePersistenceOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.cookie"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "cookie" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBCookiePersistenceProfile"

  config_cookie_persistence = {
    cookie_mode       = "INSERT"
    cookie_type       = "session"
    cookie_name       = "SERVERID"
    cookie_fallback   = true
    cookie_garbling   = true
    share_persistence = false
    cookie_domain     = ".example.com"
  }
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "service_type", "LBCookiePersistenceProfile"),
		resource.TestCheckResourceAttr(resourceName, "config_cookie_persistence.cookie_type", "session"),
		resource.TestCheckResourceAttr(resourceName, "config_cookie_persistence.cookie_mode", "INSERT"),
		resource.TestCheckResourceAttr(resourceName, "config_cookie_persistence.cookie_name", "SERVERID"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				// The config block is reconstructed from the API on import (see
				// read_config.go), so the imported resource is complete and plans
				// cleanly. It is still excluded from ImportStateVerify because the
				// API echoes server-applied config defaults that the applied state
				// does not carry where the user omitted them, so the imported and
				// applied states legitimately differ on those fields.
				ImportStateVerifyIgnore: []string{"config_cookie_persistence"},
				ResourceName:            resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					lbID := rs.Primary.Attributes["load_balancer_id"]
					id := rs.Primary.Attributes["id"]

					return lbID + "/" + id, nil
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceSourceIpPersistenceOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.source_ip"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "source_ip" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBSourceIpPersistenceProfile"

  config_source_ip_persistence = {
    purge_entries            = true
    ha_persistence_mirroring = false
    share_persistence        = true
    persistence_entry_timeout = 600
  }
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "service_type", "LBSourceIpPersistenceProfile"),
		resource.TestCheckResourceAttr(resourceName, "config_source_ip_persistence.purge_entries", "true"),
		resource.TestCheckResourceAttr(resourceName, "config_source_ip_persistence.persistence_entry_timeout", "600"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceClientSslCustomOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.client_ssl"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "client_ssl" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBClientSslProfile"

  config_client_ssl = {
    ssl_suite               = "CUSTOM"
    supported_ssl_ciphers   = ["TLS_RSA_WITH_AES_128_GCM_SHA256"]
    supported_ssl_protocols = ["TLS_V1_2"]
    session_cache           = true
    session_cache_timeout   = 300
    prefer_server_cipher    = true
  }
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "service_type", "LBClientSslProfile"),
		resource.TestCheckResourceAttr(resourceName, "config_client_ssl.ssl_suite", "CUSTOM"),
		resource.TestCheckResourceAttr(resourceName, "config_client_ssl.session_cache", "true"),
		resource.TestCheckResourceAttr(resourceName, "config_client_ssl.prefer_server_cipher", "true"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				// The config block is reconstructed from the API on import (see
				// read_config.go), so the imported resource is complete and plans
				// cleanly. It is still excluded from ImportStateVerify because the
				// API echoes server-applied config defaults that the applied state
				// does not carry where the user omitted them, so the imported and
				// applied states legitimately differ on those fields.
				ImportStateVerifyIgnore: []string{"config_client_ssl"},
				ResourceName:            resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					lbID := rs.Primary.Attributes["load_balancer_id"]
					id := rs.Primary.Attributes["id"]

					return lbID + "/" + id, nil
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceClientSslBalancedOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.client_ssl_balanced"

	config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "client_ssl_balanced" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBClientSslProfile"

  config_client_ssl = {
    ssl_suite = "BALANCED"
  }
}
`, profileName)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "config_client_ssl.ssl_suite", "BALANCED"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceHttpUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.http_update"

	createConfig := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "http_update" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBHttpProfile"

  config_http = {
    https_redirect     = true
    ntlm_authentication = false
  }
}
`, profileName)

	updateConfig := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "http_update" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBHttpProfile"

  config_http = {
    https_redirect     = false
    redirect_address   = "https://redirect.example.com"
    ntlm_authentication = true
  }
}
`, profileName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "config_http.https_redirect", "true"),
					resource.TestCheckResourceAttr(resourceName, "config_http.ntlm_authentication", "false"),
				),
			},
			{
				Config: updateConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "config_http.https_redirect", "false"),
					resource.TestCheckResourceAttr(resourceName, "config_http.redirect_address", "https://redirect.example.com"),
					resource.TestCheckResourceAttr(resourceName, "config_http.ntlm_authentication", "true"),
				),
			},
			// Idempotent plan-only step
			{
				Config:             updateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceServiceTypeRequiresReplace(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	profileName := acctest.RandomWithPrefix("TestAccMorpheus")
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_profile.replace"

	httpConfig := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "replace" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBHttpProfile"

  config_http = {
    http_idle_timeout = 30
  }
}
`, profileName)

	fastTCPConfig := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "replace" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBFastTcpProfile"

  config_fast_tcp = {
    fast_tcp_idle_timeout = 1800
  }
}
`, profileName)

	var initialID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: httpConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_type", "LBHttpProfile"),
					testExtractResourceAttr(resourceName, "id", &initialID),
				),
			},
			{
				Config: fastTCPConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_type", "LBFastTcpProfile"),
					testCheckResourceAttrNotEqual(resourceName, "id", &initialID),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerProfileResourceSmokeVariants(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	variants := []struct {
		serviceType string
		configBlock string
	}{
		{
			serviceType: "LBFastTcpProfile",
			configBlock: `config_fast_tcp = { fast_tcp_idle_timeout = 1800 }`,
		},
		{
			serviceType: "LBFastUdpProfile",
			configBlock: `config_fast_udp = { fast_udp_idle_timeout = 300 }`,
		},
		{
			serviceType: "LBGenericPersistenceProfile",
			configBlock: `config_generic_persistence = { persistence_entry_timeout = 300 }`,
		},
		{
			serviceType: "LBServerSslProfile",
			configBlock: `config_server_ssl = { ssl_suite = "BALANCED" }`,
		},
	}

	for _, v := range variants {
		v := v
		t.Run(v.serviceType, func(t *testing.T) {
			defer testhelpers.RecordResult(t)

			lbName := acctest.RandomWithPrefix(t.Name())
			lbName = lbName[0:16] + lbName[len(lbName)-16:]

			lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
				"Name": lbName,
			})
			if err != nil {
				t.Fatalf("failed to render lb config: %s", err)
			}

			profileName := acctest.RandomWithPrefix("TestAccMorpheus")
			providerConfig := testhelpers.ProviderBlock()
			resourceName := "hpe_morpheus_load_balancer_profile.smoke"

			config := providerConfig + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "smoke" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = %q
  %s
}
`, profileName, v.serviceType, v.configBlock)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, "service_type", v.serviceType),
						),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						// The config block is reconstructed from the API on import
						// (see read_config.go), so the imported resource is complete
						// and plans cleanly. It is still excluded from
						// ImportStateVerify because the API echoes server-applied
						// config defaults that the applied state does not carry where
						// the user omitted them, so the imported and applied states
						// legitimately differ on those fields.
						ImportStateVerifyIgnore: []string{
							"config_fast_tcp",
							"config_fast_udp",
							"config_generic_persistence",
							"config_server_ssl",
						},
						ResourceName: resourceName,
						ImportStateIdFunc: func(s *terraform.State) (string, error) {
							rs, ok := s.RootModule().Resources[resourceName]
							if !ok {
								return "", fmt.Errorf("resource not found: %s", resourceName)
							}

							lbID := rs.Primary.Attributes["load_balancer_id"]
							id := rs.Primary.Attributes["id"]

							return lbID + "/" + id, nil
						},
					},
				},
			})
		})
	}
}

// --- Test helpers ---

func testExtractResourceAttr(name, key string, dest *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found", name)
		}

		*dest = rs.Primary.Attributes[key]

		return nil
	}
}

func testCheckResourceAttrNotEqual(name, key string, expected *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found", name)
		}

		actual := rs.Primary.Attributes[key]
		if actual == *expected {
			return fmt.Errorf("expected %s.%s to differ from %q, but got %q", name, key, *expected, actual)
		}

		return nil
	}
}
