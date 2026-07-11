// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		loginName := acctest.RandomName("user")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckUserDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccUserConfig(loginName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureUserExists(t, "opsramp_user.test_user"),
						resource.TestCheckResourceAttrSet("opsramp_user.test_user", "id"),
						resource.TestCheckResourceAttr("opsramp_user.test_user", "login_name", loginName),
						resource.TestCheckResourceAttr("opsramp_user.test_user", "first_name", "AccTest"),
						resource.TestCheckResourceAttr("opsramp_user.test_user", "last_name", "User"),
						resource.TestCheckResourceAttr("opsramp_user.test_user", "country", "Spain"),
					),
				},
			},
		})
	})
}

func testAccUserConfig(loginName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_user" "test_user" {
	login_name = "%s"
	password   = "AccTestP@ss1234!"
	first_name = "AccTest"
	last_name  = "User"
	email      = "%s@example.com"
	time_zone  = "Europe/Paris"
	country    = "Spain"

	user_notifications = [
		{
			notify_type             = "Account Information"
			notify_method           = "Email"
			notify_input_type       = "Primary Email"
			notify_recurring_report = false
		},
		{
			notify_type             = "Alert Notification"
			notify_method           = "No Notify"
			notify_recurring_report = false
		},
		{
			notify_type             = "Report Notification"
			notify_method           = "No Notify"
			notify_recurring_report = false
		}
	]

	change_password = false
}
`, acctest.ProviderConfigHCL(), loginName, loginName)
}

func testAccEnsureUserExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		tenantID := os.Getenv("OPSRAMP_TENANT")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetUser(tenantID, id)
		if err != nil {
			return fmt.Errorf("user %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckUserDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_user" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			user, err := apiClient.GetUser(tenantID, rs.Primary.ID)

			if user != nil && user.Status != "terminate" {
				return fmt.Errorf("user still exists: %s, object: %+v", rs.Primary.ID, user)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "404") && !strings.Contains(errText, "not found") {
					return fmt.Errorf("unexpected error checking deleted user %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
