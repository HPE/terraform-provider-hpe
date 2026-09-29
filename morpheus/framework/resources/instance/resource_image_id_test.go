// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// These cover MORPH-16314: image_id on config_vmware and config_hvm.
//
// The interesting assertion is not that image_id survives in state after
// apply. That would pass even if the value were only ever echoed back from the
// plan. It is the ImportState step: import discards prior state, so image_id
// has to be recovered from the API's config object. If it is not, it reads
// back null and, because the attribute forces replacement, the next plan
// proposes destroying and recreating the instance.
//
// That is the MORPH-16190 defect class, and it is what imageIDFromConfig
// exists to prevent. Both platforms were confirmed against live appliances to
// echo imageId back verbatim, with no rewriting to template.
//
// Both tests render the same templates the published examples are generated
// from, so the examples cannot drift away from what is actually tested.

func TestAccMorpheusInstanceResourceImageIDVMware(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())

	// ImageName activates the image_id block in the template. It is omitted
	// from the published example, so this is the only place it is exercised.
	//
	// The image has to be one the hpe_morpheus_image data source can see. That
	// data source lists without a filterType, so the API applies its default of
	// userUploaded = true; a synced or system image resolves through the raw
	// API but not through the data source.
	resourceConfig, err := instance.RenderInstanceVMwareConfig(t, map[string]string{
		"Name":      name,
		"ImageName": "auto-local-vmdk",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"hpe_morpheus_instance.example", "config_vmware.image_id",
						"data.hpe_morpheus_image.vmware", "id"),
				),
			},
			{
				// The guard. Import discards prior state, so image_id must come
				// back from the API or the following plan is non-empty.
				ResourceName:      "hpe_morpheus_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
				// These do not survive import, for reasons unrelated to
				// image_id. network_interfaces loses its computed identity and
				// address, and datastore_auto_selection is a placement input
				// the API never echoes back -- it returns the datastore it
				// actually chose. The affinity group test in this package
				// ignores network_interfaces for the same reason.
				//
				// config_vmware is deliberately NOT ignored: it carries
				// image_id, which is the whole point of this test.
				ImportStateVerifyIgnore: []string{
					"network_interfaces",
					"volumes.0.datastore_auto_selection",
					"volumes.1.datastore_auto_selection",
				},
			},
		},
	})
}

func TestAccMorpheusInstanceResourceImageIDHVM(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())

	// The HVM template sets image_id from its defaults, so only the name
	// needs overriding here.
	resourceConfig, err := instance.RenderInstanceHVMConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"hpe_morpheus_instance.example", "config_hvm.image_id",
						"data.hpe_morpheus_image.hvm", "id"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
				// See the VMware test for why these are ignored. config_hvm is
				// deliberately NOT ignored: it carries image_id.
				ImportStateVerifyIgnore: []string{
					"network_interfaces",
					"volumes.0.datastore_auto_selection",
				},
			},
		},
	})
}
