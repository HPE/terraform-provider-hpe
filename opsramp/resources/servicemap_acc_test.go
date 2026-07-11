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

func TestAccServicemapResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		smName := acctest.RandomName("servicemap")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicemapDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapConfig(smName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicemapExists(t, "opsramp_servicemap.test_sm"),
						resource.TestCheckResourceAttrSet("opsramp_servicemap.test_sm", "id"),
						resource.TestCheckResourceAttr("opsramp_servicemap.test_sm", "name", smName),
						resource.TestCheckResourceAttr("opsramp_servicemap.test_sm", "type", "Service"),
					),
				},
			},
		})
	})
}

func TestAccServicemapWithChildResource(t *testing.T) {
	t.Run("parent and child", func(t *testing.T) {
		rootName := acctest.RandomName("sm-root")
		childName := acctest.RandomName("sm-child")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicemapDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapWithChildConfig(rootName, childName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicemapExists(t, "opsramp_servicemap.test_sm_root"),
						testAccEnsureServicemapExists(t, "opsramp_servicemap.test_sm_child"),
						resource.TestCheckResourceAttr("opsramp_servicemap.test_sm_root", "name", rootName),
						resource.TestCheckResourceAttr("opsramp_servicemap.test_sm_child", "name", childName),
					),
				},
			},
		})
	})
}

func testAccServicemapConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_servicemap" "test_sm" {
	name = "%s"
	type = "Service"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccServicemapWithChildConfig(rootName string, childName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_servicemap" "test_sm_root" {
	name = "%s"
	type = "Service"
}

resource "opsramp_servicemap" "test_sm_child" {
	name   = "%s"
	type   = "Service"
	parent = opsramp_servicemap.test_sm_root.id
}
`, acctest.ProviderConfigHCL(), rootName, childName)
}

func testAccEnsureServicemapExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetServicemap(tenantID, id)
		if err != nil {
			return fmt.Errorf("servicemap %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckServicemapDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_servicemap" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetServicemap(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("servicemap still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no service group exists") {
				return fmt.Errorf("unexpected error checking deleted servicemap %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
