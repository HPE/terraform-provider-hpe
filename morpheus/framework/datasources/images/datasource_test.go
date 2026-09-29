// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// The filtering, conversion and parameter mapping are covered by the internal
// unit tests, which run in under a second and can force cases a live appliance
// will not reliably produce.
//
// This is the one live check, and its job is narrow: prove the SDK call, the
// paging envelope and the conversion work end to end against a real API.
//
// It is deliberately scoped to a small result. The data source fetches every
// image its filters admit, and estates run large — an appliance used while
// developing this held 7040 images, of which 1134 were qcow2 or raw, so an
// unfiltered read took over a minute per plan. filter_type = "User" with a
// single image_type keeps this to a couple of dozen records.
const minimalConfig = `
data "hpe_morpheus_images" "minimal" {
  filter_type = "User"
  image_type  = ["qcow2"]
}
`

// TestMain lets the test harness aggregate results for CI, as the other
// acceptance test files in this repository do.
func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + minimalConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_images.minimal", "images.#"),

					// Elements carry the fields the singular data source
					// reports, which is what makes the two interchangeable.
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_images.minimal", "images.0.id"),
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_images.minimal", "images.0.name"),

					// The server-side filters were applied rather than ignored:
					// an ignored image_type returns the whole estate, and an
					// ignored filter_type returns synced and system images too.
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_images.minimal", "images.0.image_type", "qcow2"),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_images.minimal", "images.0.user_uploaded", "true"),
				),
			},
		},
	})
}
