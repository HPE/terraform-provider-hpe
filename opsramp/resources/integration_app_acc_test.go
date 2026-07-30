// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

// ---------------------------------------------------------------------------
// hpe_opsramp_integration_app – Kubernetes-2.0 SDK APP
// ---------------------------------------------------------------------------

func TestAccIntegrationAppResourceKubernetes(t *testing.T) {
	configName := acctest.RandomName("k8s-cfg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationAppDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAppKubernetesConfig(configName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationAppExists(t, "hpe_opsramp_integration_app.test_with_cfg"),
					testAccEnsureIntegrationConfigExists(t, "hpe_opsramp_integration_config.test_app_cfg"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration_app.test_with_cfg", "id"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration_config.test_app_cfg", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test_app_cfg", "name", configName),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

func testAccIntegrationAppKubernetesConfig(configName string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration_app" "test_with_cfg" {
  application                    = "Kubernetes-2.0"
  version                        = "2.3.0"
  bypass_resource_reconciliation = true
}

resource "hpe_opsramp_integration_config" "test_app_cfg" {
  integration_id = hpe_opsramp_integration_app.test_with_cfg.id
  name           = %q
  config         = jsonencode({
    Etcd                       = true
    coreDNS                    = true
    KubeProxy                  = true
    enableLog                  = true
	enableTrace                = false
	ebpfFeatureFlag            = false
	enableEBPF                 = false
    kubeEvents                 = true
    KubeletStats               = true
    KubeAPIServer              = true
    KubeScheduler              = true
    KubeClusterReceiver        = true
    KubeControllerManager      = true
    replicaCount               = 1
    DistributionType           = "K8S"
	clientLevelLogPermission = false
	clientLevelTracePermission = false
  })

  all_resources = false
}
`, acctest.ProviderConfigHCL(), configName)
}

// ---------------------------------------------------------------------------
// Shared check helpers
// ---------------------------------------------------------------------------

func testAccEnsureIntegrationAppExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.ID)
		if id == "" {
			return fmt.Errorf("resource id is empty in state for %s", resourceName)
		}

		tenantID, _ := acctest.LookupProviderEnv("tenant")

		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetIntegrationV3(tenantID, id)
		if err != nil {
			return fmt.Errorf("integration app %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckIntegrationAppDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_integration_app" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			integration, err := apiClient.GetIntegrationV3(tenantID, rs.Primary.ID)

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "no installed integration found with") {
					return fmt.Errorf("unexpected error checking deleted integration app %s: %w", rs.Primary.ID, err)
				}
			} else if integration != nil && integration.Status != "disabled" {
				return fmt.Errorf("integration app still exists with status '%s': %s", integration.Status, rs.Primary.ID)
			}

		}

		return nil
	}
}
