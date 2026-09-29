// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusStorageVolumeResourceWriteOnlyOk: the write-only config map
// is taken by Morpheus only at creation, so a change plans a replacement
// (MORPH-18005); the deprecated config_wo_version marker still replaces when
// bumped and does nothing when removed; both import forms leave a clean plan.
// The digest of a map is sha256 of its JSON form. Needs an Alletra MP storage
// server (id 1 in the example), as the other volume tests do.
func TestAccMorpheusStorageVolumeResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix("tfwo-vol")

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_storage_volume.wo",
		Attribute: "config",
		Marker:    "config_wo_version",
		APIPath:   "/api/storage-volumes",
		APIKey:    "storageVolume.config.hpe_storage_datastore",
		Config: func(value, extra string, unrelated bool) string {
			size := "30"
			if unrelated {
				size = "40"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_volume" "wo" {
  name              = "` + name + `"
  type_code         = "alletra-mp-volume"
  storage_server_id = 1
  max_storage       = ` + size + `

  config = {
    hpe_storage_datastore = ` + value + `
  }` + extra + `
}
`
		},
		ImportDigest: func(value string) string { return testhelpers.Sha256Hex(`{"hpe_storage_datastore":` + value + `}`) },
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyReplaceSteps(c, rec, "5", "6")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "6")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
