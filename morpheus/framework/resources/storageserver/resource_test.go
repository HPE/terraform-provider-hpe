package storageserver_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageServerResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable storage
	// service is required; without the fixture there is no appliance to create
	// against. The nfs example type does not reach a service, but every
	// create-based test in this package is gated on the fixture so the suite
	// is consistent on a 9.1.2 appliance.
	f := testhelpers.RequireStorageServerFixture(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_storage_server.example"

	resourceConfig := storageServerFixtureConfig(name, f)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type", f.TypeCode),
		resource.TestCheckResourceAttr(resourceName, "service_url", f.URL),
		resource.TestCheckResourceAttr(resourceName, "service_username", f.Username),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"write_only_revision", "service_username",
					"service_password_wo", "service_password_wo_version",
				},
				ResourceName: resourceName,
			},
		},
	})
}

// storageServerFixtureConfig renders a storage server against the reachable
// fixture service. The example render helper only emits name+type, so the
// create-based tests build HCL directly to include the service URL and
// credentials the 9.1.2 reachability check requires.
func storageServerFixtureConfig(name string, f testhelpers.StorageServerFixture) string {
	return `
resource "hpe_morpheus_storage_server" "example" {
  name             = "` + name + `"
  type             = "` + f.TypeCode + `"
  service_url      = "` + f.URL + `"
  service_username = "` + f.Username + `"
  service_password_wo = "` + f.Password + `"
}
`
}

func TestAccMorpheusStorageServerResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable storage
	// service is required; without the fixture there is no appliance to create
	// against.
	f := testhelpers.RequireStorageServerFixture(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig := storageServerFixtureConfig(name, f)
	updateConfig := storageServerFixtureConfig(updatedName, f)

	resourceName := "hpe_morpheus_storage_server.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type", f.TypeCode),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "type", f.TypeCode),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
