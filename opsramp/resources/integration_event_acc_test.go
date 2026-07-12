// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ---------------------------------------------------------------------------
// hpe_opsramp_integration_event – uses base notifier
// ---------------------------------------------------------------------------

func TestAccIntegrationEventResource_BaseNotifier(t *testing.T) {
	eventName := acctest.RandomName("event")
	eventNameUpdated := eventName + "-upd"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationEventDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccIntegrationEventBaseNotifierConfig(eventName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationEventExists(t, "hpe_opsramp_integration_event.test"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration_event.test", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "name", eventName),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "entity", "DEFAULT_RESOURCE"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "event_type", "CREATE"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "use_base_notifier", "true"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "active", "true"),
				),
			},
			// Update – rename and change active state
			{
				Config: testAccIntegrationEventBaseNotifierConfig(eventNameUpdated, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "name", eventNameUpdated),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test", "active", "false"),
				),
			},
			// ImportState testing — import ID is <integration_id>:<event_id>
			{
				ResourceName:      "hpe_opsramp_integration_event.test",
				ImportState:       true,
				ImportStateIdFunc: testAccIntegrationEventImportStateIdFunc("hpe_opsramp_integration_event.test"),
				ImportStateVerify: true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// hpe_opsramp_integration_event – overrides notifier (OAUTH2)
// ---------------------------------------------------------------------------

func TestAccIntegrationEventResource_OverrideNotifier(t *testing.T) {
	eventName := acctest.RandomName("event-override")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationEventDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationEventOverrideNotifierConfig(eventName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationEventExists(t, "hpe_opsramp_integration_event.test_override"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test_override", "use_base_notifier", "false"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test_override", "notifier.auth_type", "OAUTH2"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_event.test_override", "notifier.grant_type", "PASSWORD"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

// The integration_event tests share a parent CUSTOM integration (outbound required).
func testAccIntegrationEventParentConfig(pos int) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration" "event_parent" {
  display_name = "tf-acc-event-parent - %d"
  application  = "CUSTOM"
  category     = "Custom"

  inbound = {
    auth_type = "OAUTH2"
  }

  outbound = {
    base_uri  = "https://httpbin.org/post"
    auth_type = "NONE"
  }
}
`, acctest.ProviderConfigHCL(), pos)
}

func testAccIntegrationEventBaseNotifierConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration_event" "test" {
  integration_id         = hpe_opsramp_integration.event_parent.id
  name                   = %q
  entity                 = "DEFAULT_RESOURCE"
  event_type             = "CREATE"
  use_base_notifier      = true
  third_party_event_type = "POST"
  endpoint_uri           = "https://httpbin.org/post"
  event_payload          = "test payload"
  active                 = %t
  headers = {
    "Content-Type" = "application/json"
    "Accept" = "application/json"
  }
}
`, testAccIntegrationEventParentConfig(1), name, enabled)
}

func testAccIntegrationEventOverrideNotifierConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration_event" "test_override" {
  integration_id         = hpe_opsramp_integration.event_parent.id
  name                   = %q
  entity                 = "DEFAULT_RESOURCE"
  event_type             = "UPDATE"
  use_base_notifier      = false
  third_party_event_type = "POST"
  endpoint_uri           = "https://httpbin.org/post"
  event_payload          = "override payload"
  active                 = true
  headers = {
    "Content-Type" = "application/json"
    "Accept" = "application/json"
  }

  notifier = {
    type             = "REST_API"
    base_uri         = "https://httpbin.org/post"
    auth_type        = "OAUTH2"
    grant_type       = "PASSWORD"
    username         = "test-user"
	password		 = "test-password"
    api_key          = "test-client-id"
    api_secret          = "test-client-id"
    access_token_url = "https://httpbin.org/post"
    scope            = "read"
  }
}
`, testAccIntegrationEventParentConfig(3), name)
}

// ---------------------------------------------------------------------------
// Shared check helpers
// ---------------------------------------------------------------------------

func testAccEnsureIntegrationEventExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		integrationID := rs.Primary.Attributes["integration_id"]

		tenantID, _ := acctest.LookupProviderEnv("tenant")

		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetIntegrationEvent(tenantID, integrationID, id)
		if err != nil {
			return fmt.Errorf("integration event %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckIntegrationEventDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_integration_event" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			integrationID := rs.Primary.Attributes["integration_id"]

			_, err := apiClient.GetIntegrationEvent(tenantID, integrationID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("integration event still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no installed integration found with") {
				return fmt.Errorf("unexpected error checking deleted integration event %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

// testAccIntegrationEventImportStateIdFunc builds the composite import ID: <integration_id>:<event_id>
func testAccIntegrationEventImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		integrationID := rs.Primary.Attributes["integration_id"]
		eventID := rs.Primary.ID

		return integrationID + ":" + eventID, nil
	}
}
