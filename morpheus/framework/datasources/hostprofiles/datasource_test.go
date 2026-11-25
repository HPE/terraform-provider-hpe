// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofiles_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusHostProfilesDataSourceOk provisions a host profile with
// Terraform and lists it back through the data source, filtering on the
// fixture's unique name so the created profile is the single match. Creating
// host profiles is restricted to the master tenant, so the test is gated
// accordingly. The list is phrase-filtered server-side (a database query, not
// the eventually-consistent search index), so the fixture is visible
// immediately.
func TestAccMorpheusHostProfilesDataSourceOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile, capabilities.MasterTenant)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// phrase is a literal string, not a reference to the fixture, so an explicit
	// depends_on is required to order the create before the list read.
	config := providerConfig + `
resource "hpe_morpheus_host_profile" "fixture" {
  name        = "` + name + `"
  description = "host profiles data source acceptance fixture"
}

data "hpe_morpheus_host_profiles" "example" {
  phrase     = "` + name + `"
  depends_on = [hpe_morpheus_host_profile.fixture]
}
`

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_host_profiles.example", "host_profiles.#", "1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_host_profiles.example",
			"host_profiles.0.name", name,
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_host_profiles.example", "host_profiles.0.id",
			"hpe_morpheus_host_profile.fixture", "id",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}
