// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing an optional attribute clears it on the image (MORPH-17990): the
// labels, the minimum disk, the SSH username, the user data and the OS type,
// all Optional+Computed and cleared once an apply has recorded them. min_ram
// stays: Morpheus rejects a null minimum memory with a server error.
func TestAccMorpheusImageRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_image.test"

	config := func(with ...string) string {
		set := map[string]string{
			"labels":       `labels       = ["terraform", "acctest"]`,
			"min_disk":     `min_disk     = 20`,
			"ssh_username": `ssh_username = "alpine"`,
			"user_data":    `user_data    = "#!/bin/sh\necho hello"`,
			"os_type_id":   `os_type_id   = data.hpe_morpheus_os_type.linux.id`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_os_type" "linux" {
  name = "linux"
}

data "hpe_morpheus_storage_bucket" "local" {
  name = %q
}

resource "hpe_morpheus_image" "test" {
  name                = %q
  image_type          = "qcow2"
  url                 = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.2-x86_64-bios-cloudinit-r0.qcow2"
  storage_provider_id = data.hpe_morpheus_storage_bucket.local.id
  visibility          = "private"
%s
}
`, testhelpers.ImageStorageBucket("Local Storage"), name, strings.Join(lines, "\n"))
	}

	all := []string{"labels", "min_disk", "ssh_username", "user_data", "os_type_id"}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	var noAttrs []resource.TestCheckFunc
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		key := a
		if a == "labels" {
			key = "labels.#"
		}
		noAttrs = append(noAttrs, resource.TestCheckNoResourceAttr(rn, key))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "labels.#", "2"),
					resource.TestCheckResourceAttr(rn, "min_disk", "20"),
					resource.TestCheckResourceAttr(rn, "ssh_username", "alpine"),
					resource.TestCheckResourceAttrSet(rn, "user_data"),
					resource.TestCheckResourceAttrSet(rn, "os_type_id"),
				),
			},
			{
				Config:           config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: nullChecks},
				Check:            resource.ComposeAggregateTestCheckFunc(noAttrs...),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
