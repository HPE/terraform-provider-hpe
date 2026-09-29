// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusProvisioningLicenseResourceWriteOnlyOk: changes to the
// required write-only license_key_wo are detected by digest (MORPH-18005); the
// deprecated marker still works; both import forms leave a clean plan.
func TestAccMorpheusProvisioningLicenseResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_provisioning_license.wo",
		Attribute: "license_key_wo",
		Marker:    "license_key_wo_version",
		APIPath:   "/api/provisioning-licenses",
		APIKey:    "license.licenseKey",
		Config: func(value, extra string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_provisioning_license" "wo" {
  name           = "` + name + `"
  license_type   = "win"
  description    = "` + desc + `"
  license_key_wo = "` + value + `"` + extra + `
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec,
		"AAAAA-BBBBB-CCCCC-DDDDD-EEEE1", "AAAAA-BBBBB-CCCCC-DDDDD-EEEE2", "AAAAA-BBBBB-CCCCC-DDDDD-EEEE3")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "AAAAA-BBBBB-CCCCC-DDDDD-EEEE3")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
