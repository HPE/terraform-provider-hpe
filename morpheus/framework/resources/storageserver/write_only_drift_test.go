// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusStorageServerResourceWriteOnlyDriftDetected: the service
// password is rotated out of band (a direct API PUT, as the Morpheus UI would
// do). The next plan sees the appliance's stored-secret hash no longer matches
// the baseline recorded at the last send and plans an update that re-sends the
// configured value; the plan is then clean. Needs an appliance with the Alletra
// MP storage server type.
func TestAccMorpheusStorageServerResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)
	testhelpers.SkipUnlessAcc(t)

	// Morpheus 9.1.2 validates the storage server credentials on create and
	// update (a reachable host with a wrong password is refused: "Error
	// connecting to the storage server"), so an out-of-band rotation cannot be
	// staged against an ordinary live integration. Set
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
	addr := "hpe_morpheus_storage_server.wo"

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_server" "wo" {
  name                = "` + name + `"
  type                = "` + f.TypeCode + `"
  description         = "write-only drift"
  service_url         = "` + f.URL + `"
  service_username    = "` + f.Username + `"
  service_password_wo = "` + f.Password + `"
}
`

	var id int64

	capture := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not in state", addr)
		}

		var err error
		id, err = strconv.ParseInt(rs.Primary.ID, 10, 64)

		return err
	}

	rotate := func() {
		if id == 0 {
			t.Fatal("capture did not run: storage server id is 0")
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("client: %s", err)
		}

		rotated := "Rotated456!"
		body := sdk.UpdateStorageServersRequest{
			StorageServer: sdk.UpdateStorageServersRequestStorageServer{
				ServicePassword: &rotated,
			},
		}

		// Compare the hash around the PUT: a rotation that is accepted but not
		// stored would otherwise surface as a confusing plan-check failure.
		hash := func() string {
			resp, _, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
			if err != nil || resp == nil || resp.StorageServer == nil {
				return ""
			}
			if h := resp.StorageServer.ServicePasswordHash.Get(); h != nil {
				return *h
			}

			return ""
		}
		before := hash()

		// With TF_ACC_STORAGE_SERVER_ROTATABLE set the target storage server is
		// expected to accept the rotation; a rejected PUT is now a real failure
		// (the test can no longer silently pass by skipping here).
		if _, httpResp, err := client.StorageAPI.
			UpdateStorageServers(ctx, id).
			UpdateStorageServersRequest(body).Execute(); err != nil {
			t.Fatalf("out-of-band rotate rejected by appliance: %s (status %v)", err, httpResp)
		}

		if after := hash(); after == before {
			t.Fatalf("rotate PUT did not change servicePasswordHash (%q); drift cannot be detected", before)
		}
	}

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: config, Check: capture},
			{
				PreConfig:        rotate,
				Config:           config,
				ConfigPlanChecks: update,
			},
			{Config: config, PlanOnly: true},
		},
	})
}
