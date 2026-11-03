// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusImageResourceWriteOnlyOk: changes to ssh_password_wo are
// detected by digest (MORPH-18005); the deprecated marker still works; both
// import forms leave a clean plan.
func TestAccMorpheusImageResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	visibility := testhelpers.TenantVisibility(t)

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:   "hpe_morpheus_image.wo",
		Attribute: "ssh_password_wo",
		Marker:    "ssh_password_wo_version",
		APIPath:   "/api/virtual-images",
		APIKey:    "virtualImage.sshPassword",
		Config: func(value, extra string, unrelated bool) string {
			// description requires replacement on this resource; min_ram does not.
			minRAM := "1"
			if unrelated {
				minRAM = "2"
			}

			return testhelpers.ProviderBlock() + `
data "hpe_morpheus_os_type" "test" {
  name = "linux"
}

data "hpe_morpheus_storage_bucket" "test" {
  name = "Local Storage"
}

resource "hpe_morpheus_image" "wo" {
  name                = "` + name + `"
  description         = "write-only matrix"
  min_ram             = ` + minRAM + `
  image_type          = "qcow2"
  url                 = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.2-x86_64-bios-cloudinit-r0.qcow2"
  storage_provider_id = data.hpe_morpheus_storage_bucket.test.id
  os_type_id          = data.hpe_morpheus_os_type.test.id
  ssh_username        = "alpine"
  ssh_password_wo     = "` + value + `"` + extra + `
  visibility          = "` + visibility + `"
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	// url is not returned by the API (as in the existing import test).
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!", "url")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps:                    steps,
	})
}
