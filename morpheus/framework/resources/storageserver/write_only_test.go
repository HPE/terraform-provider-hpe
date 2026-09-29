// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusStorageServerResourceWriteOnlyOk: changes to
// service_password_wo are detected by digest (MORPH-18005); the deprecated
// marker still works; both import forms leave a clean plan. Needs an appliance
// with the Alletra MP storage server type.
func TestAccMorpheusStorageServerResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_storage_server.wo",
		Attribute: "service_password_wo",
		Marker:    "service_password_wo_version",
		APIPath:   "/api/storage-servers",
		APIKey:    "storageServer.servicePassword",
		Config: func(value, extra string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_server" "wo" {
  name                = "` + name + `"
  type                = "alletra-mp"
  description         = "` + desc + `"
  service_url         = "https://alletra.example.com"
  service_username    = "admin"
  service_password_wo = "` + value + `"` + extra + `
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
