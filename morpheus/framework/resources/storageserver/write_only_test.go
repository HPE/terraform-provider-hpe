// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver_test

import (
	"os"
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

	// Morpheus 9.1.2 validates the storage server credentials on create and
	// update (a reachable host with a wrong password is refused: "Error
	// connecting to the storage server"), so neither a change-detection matrix
	// that sends three different passwords nor an out-of-band rotation to a
	// wrong value can be staged against an ordinary live integration. The
	// password handling is covered by the unit tests and by the same wodigest
	// mechanism exercised live on resources whose secret the appliance stores
	// without validating (hpe_morpheus_image). Set
	// TF_ACC_STORAGE_SERVER_ROTATABLE=1 to re-enable this test against a storage
	// server whose credentials can be rotated. Skip first, before the fixture
	// helper, so the body is not executed.
	if os.Getenv("TF_ACC_STORAGE_SERVER_ROTATABLE") == "" {
		t.Skip("set TF_ACC_STORAGE_SERVER_ROTATABLE=1 against a storage server whose credentials can be rotated; MORPH-19005 follow-up")
	}

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable storage
	// service is required; without the fixture there is no appliance to create
	// against.
	f := testhelpers.RequireStorageServerFixture(t)

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
  type                = "` + f.TypeCode + `"
  description         = "` + desc + `"
  service_url         = "` + f.URL + `"
  service_username    = "` + f.Username + `"
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
