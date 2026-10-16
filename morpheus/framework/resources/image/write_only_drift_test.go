// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusImageResourceWriteOnlyDriftDetected: the SSH password is
// rotated out of band (a direct API PUT, as the Morpheus UI would do). The next
// plan sees the appliance's stored-secret hash no longer matches the baseline
// recorded at the last send and plans an update that re-sends the configured
// value; the plan is then clean.
func TestAccMorpheusImageResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	visibility := testhelpers.TenantVisibility(t)
	addr := "hpe_morpheus_image.wo"

	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_os_type" "test" {
  name = "linux"
}

data "hpe_morpheus_storage_bucket" "test" {
  name = "` + testhelpers.ImageStorageBucket("Local Storage") + `"
}

resource "hpe_morpheus_image" "wo" {
  name                = "` + name + `"
  description         = "write-only drift"
  min_ram             = 1
  image_type          = "qcow2"
  url                 = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.2-x86_64-bios-cloudinit-r0.qcow2"
  storage_provider_id = data.hpe_morpheus_storage_bucket.test.id
  os_type_id          = data.hpe_morpheus_os_type.test.id
  ssh_username        = "alpine"
  ssh_password_wo     = "Secret123!"
  visibility          = "` + visibility + `"
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

	// hashOf reads the appliance's stored-secret hash for the image.
	hashOf := func(ctx context.Context, client *sdk.APIClient) string {
		resp, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
		if err != nil {
			t.Fatalf("read image back: %s (status %v)", err, httpResp)
		}
		if resp.VirtualImage == nil {
			t.Fatalf("read image back: nil VirtualImage")
		}

		if v := resp.VirtualImage.SshPasswordHash.Get(); v != nil {
			return *v
		}

		return ""
	}

	rotate := func() {
		if id == 0 {
			t.Fatal("capture did not run: image id is 0")
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("client: %s", err)
		}

		before := hashOf(ctx, client)

		rotated := "Rotated456!"
		body := sdk.UpdateVirtualImageRequest{
			VirtualImage: sdk.UpdateVirtualImageRequestVirtualImage{
				SshPassword: *sdk.NewNullableString(&rotated),
			},
		}

		if _, httpResp, err := client.LibraryAPI.
			UpdateVirtualImage(ctx, id).
			UpdateVirtualImageRequest(body).Execute(); err != nil {
			t.Fatalf("out-of-band rotate failed: %s (status %v)", err, httpResp)
		}

		// Assert the appliance actually re-stored the secret: if the stored
		// hash is unchanged there is no out-of-band drift to detect and the
		// test would pass even with drift detection removed -- a false pass.
		if after := hashOf(ctx, client); after == before {
			t.Fatalf("appliance did not re-store the SSH password: hash unchanged after rotate (%q); nothing to detect", after)
		}
	}

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
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
