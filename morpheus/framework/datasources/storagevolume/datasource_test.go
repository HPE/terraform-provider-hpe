// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/storagevolume"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url      = ""
    username = ""
    password = ""
  }
}
`

// TestAccMorpheusStorageVolumeNoSearchTerms verifies that the data source
// errors when neither id nor name is supplied. The error is raised before any
// API call, so the offline provider configuration is sufficient.
func TestAccMorpheusStorageVolumeNoSearchTerms(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_storage_volume" "test" {
}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(storagevolume.ErrorNoValidSearchTerms),
			},
		},
	})
}
