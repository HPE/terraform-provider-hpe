// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageservers_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusStorageServersEmptyResult verifies that a filter matching no
// storage servers yields an empty result. This does not require creating a
// storage server (which needs real backend connectivity).
func TestAccMorpheusStorageServersEmptyResult(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_storage_servers" "test" {
        filter {
          name   = "name"
          values = ["this-name-should-not-exist-______"]
        }
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_servers.test", "storage_servers.#", "0",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
