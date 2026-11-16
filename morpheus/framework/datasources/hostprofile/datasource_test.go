// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile_test

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

// TestAccMorpheusHostProfileDataSourceOk provisions a host profile with
// Terraform and reads it back through the data source by id, asserting the
// data source echoes the fixture's name. Creating host profiles is restricted
// to the master tenant, so the test is gated accordingly; this removes the
// dependence on a pre-existing profile of an assumed id.
func TestAccMorpheusHostProfileDataSourceOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile, capabilities.MasterTenant)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// The data source reads by the fixture's id, so the reference creates an
	// implicit dependency: Terraform creates the profile before reading it.
	config := providerConfig + `
resource "hpe_morpheus_host_profile" "fixture" {
  name        = "` + name + `"
  description = "host profile data source acceptance fixture"
}

data "hpe_morpheus_host_profile" "example" {
  id = hpe_morpheus_host_profile.fixture.id
}
`

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_host_profile.example", "id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_host_profile.example", "name", name,
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_host_profile.example", "id",
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
