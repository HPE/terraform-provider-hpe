// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// TestAccMorpheusCloudResourceWriteOnlyOk: changes to the nested write-only
// config_vmware.password are detected by digest (MORPH-18005); the deprecated
// nested marker still works; both import forms (the nested attribute as a
// dotted path in the import ID) leave a clean plan. Needs a reachable vCenter:
// Morpheus validates the connection on create even for a disabled cloud.
func TestAccMorpheusCloudResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	visibility := testhelpers.TenantVisibility(t)

	apiURL, ok := testhelpers.LookupProviderEnv("", "vcenter_url")
	if !ok {
		t.Skip("TF_VAR_testacc_morpheus_vcenter_url (and _vcenter_username/_vcenter_datacenter) not set")
	}

	vcUser, _ := testhelpers.LookupProviderEnv("", "vcenter_username")
	vcDC, _ := testhelpers.LookupProviderEnv("", "vcenter_datacenter")

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:   "hpe_morpheus_cloud.wo",
		Attribute: "config_vmware.password",
		Marker:    "config_vmware.password_version",
		APIPath:   "/api/zones",
		APIKey:    "zone.config.password",
		Config: func(value, extra string, unrelated bool) string {
			location := "somewhere"
			if unrelated {
				location = "elsewhere"
			}

			return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + `
resource "hpe_morpheus_cloud" "wo" {
  name       = "` + name + `"
  code       = "` + name + `"
  tenant_id  = ` + testhelpers.WhoamiTenantIDRef + `
  group_id   = 1
  enabled    = false
  location   = "` + location + `"
  visibility = "` + visibility + `"

  config_vmware = {
    api_url    = "` + apiURL + `"
    username   = "` + vcUser + `"
    password   = "` + value + `"` + extra + `
    datacenter = "` + vcDC + `"
  }
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
