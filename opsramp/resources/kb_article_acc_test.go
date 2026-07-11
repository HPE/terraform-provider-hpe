// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccKBArticleResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("kb-art-cat")
		articleSubject := acctest.RandomName("kb-article")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckKBArticleDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccKBArticleConfig(catName, articleSubject),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureKBArticleExists(t, "opsramp_kb_article.test_article"),
						resource.TestCheckResourceAttrSet("opsramp_kb_article.test_article", "id"),
						resource.TestCheckResourceAttr("opsramp_kb_article.test_article", "subject", articleSubject),
					),
				},
			},
		})
	})
}

func testAccKBArticleConfig(catName string, subject string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_kb_category" "test_art_category" {
	name        = "%s"
	description = "Category for article test"
}

resource "opsramp_kb_article" "test_article" {
	subject     = "%s"
	content     = "Acceptance test article content"
	category_id = opsramp_kb_category.test_art_category.id
}
`, acctest.ProviderConfigHCL(), catName, subject)
}

func testAccEnsureKBArticleExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetKBArticle(tenantID, id)
		if err != nil {
			return fmt.Errorf("KB article %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckKBArticleDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_kb_article" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			article, err := apiClient.GetKBArticle(tenantID, rs.Primary.ID)

			if article != nil && article.State != "TRASH" {
				return fmt.Errorf("KB article still exists: %s (%s), object %+v", rs.Primary.ID, article.Subject, article)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "no article found with id") {
					return fmt.Errorf("unexpected error checking deleted KB article %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
