package storage_server_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

const resourceName = "hpe_morpheus_storage_server.test"

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageServerResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping: requires external storage infrastructure and stored credential")

		return
	}

	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create with local credentials
			{
				Config: providerConfig + testAccStorageServerConfigLocalCreds(rName, "nfs", "testuser", "testpass"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "nfs"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(resourceName, "service_username", "testuser"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_password_wo_version",
				},
			},
			// Update description and visibility
			{
				Config: providerConfig + testAccStorageServerConfigUpdated(
					rName, "nfs", "testuser", "testpass", "updated description", "public"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
		},
	})
}

func TestAccStorageServerResource_credential(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping: requires external storage infrastructure and stored credential")

		return
	}

	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create with stored credential
			{
				Config: providerConfig + testAccStorageServerConfigCredential(rName, "nfs", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "credential_id", "1"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_password_wo_version",
				},
			},
		},
	})
}

func TestAccStorageServerResource_tenants(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping: requires external storage infrastructure and multiple tenants")

		return
	}

	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create with tenants
			{
				Config: providerConfig + testAccStorageServerConfigTenants(rName, "nfs", "testuser", "testpass", []int{1, 2}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenants.#", "2"),
				),
			},
			// Update tenants
			{
				Config: providerConfig + testAccStorageServerConfigTenants(rName, "nfs", "testuser", "testpass", []int{1}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tenants.#", "1"),
				),
			},
		},
	})
}

// TestAccStorageServerResource_planOnly validates the schema and config without
// requiring a real Morpheus backend. This catches schema issues, conflictsWith
// validation, and default value problems.
func TestAccStorageServerResource_planOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Validate local creds config plans successfully
			{
				Config:   providerConfig + testAccStorageServerConfigLocalCreds(rName, "nfs", "testuser", "testpass"),
				PlanOnly: true,
			},
			// Validate credential config plans successfully
			{
				Config:   providerConfig + testAccStorageServerConfigCredential(rName, "nfs", 1),
				PlanOnly: true,
			},
			// Validate tenants config plans successfully
			{
				Config:   providerConfig + testAccStorageServerConfigTenants(rName, "nfs", "testuser", "testpass", []int{1, 2}),
				PlanOnly: true,
			},
		},
	})
}

// TestAccStorageServerResource_conflictsValidation verifies that credential_id
// and service_username/service_password_wo cannot be set together.
func TestAccStorageServerResource_conflictsValidation(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name            = %q
  type            = "nfs"
  credential_id   = 1
  service_username = "user"
}
`, rName),
				ExpectError: regexp.MustCompile(`(?i)conflict`),
			},
		},
	})
}

// testAccStorageServerConfigLocalCreds returns a config with local username/password auth.
func testAccStorageServerConfigLocalCreds(name, serverType, username, password string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name                       = %q
  type                       = %q
  service_username           = %q
  service_password_wo        = %q
  service_password_wo_version = 1
}
`, name, serverType, username, password)
}

// testAccStorageServerConfigUpdated returns a config with updated description and visibility.
func testAccStorageServerConfigUpdated(name, serverType, username, password, description, visibility string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name                       = %q
  type                       = %q
  service_username           = %q
  service_password_wo        = %q
  service_password_wo_version = 1
  description                = %q
  visibility                 = %q
}
`, name, serverType, username, password, description, visibility)
}

// testAccStorageServerConfigCredential returns a config using a stored credential.
func testAccStorageServerConfigCredential(name, serverType string, credentialID int) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name          = %q
  type          = %q
  credential_id = %d
}
`, name, serverType, credentialID)
}

// testAccStorageServerConfigTenants returns a config with a tenants list.
func testAccStorageServerConfigTenants(name, serverType, username, password string, tenants []int) string {
	tenantStr := ""
	for i, id := range tenants {
		if i > 0 {
			tenantStr += ", "
		}
		tenantStr += fmt.Sprintf("%d", id)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name                       = %q
  type                       = %q
  service_username           = %q
  service_password_wo        = %q
  service_password_wo_version = 1
  tenants                    = [%s]
}
`, name, serverType, username, password, tenantStr)
}
