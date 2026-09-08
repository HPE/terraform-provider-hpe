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

// MORPH-16949. Morpheus may provision a larger disk than was requested, and the
// provider previously recorded the request and discarded what was created. The
// real size was therefore invisible to plan and refresh, and an import disagreed
// with the state an apply had produced.
//
// `size` now records the request and `actual_size` records what exists.
//
// `size` is excluded from import verification throughout. It is not an
// observable property of a provisioned disk — import has no configuration to
// recover the request from and reads back what was actually created.
// `actual_size` is the attribute that must survive the round trip.
//
// When the platform grows a disk: a request below the image's minDisk is
// rejected outright ("A volume is smaller than the image allows"), not grown.
// Growth happens in the narrow band where the request meets minDisk but falls
// short of the image's rawSize, which is then rounded up to the next whole GB.

// TestAccMorpheusInstanceResourceVolumeActualSizeVMware exercises the common
// case, where no growth occurs and actual_size simply agrees with the request.
// It proves the attribute is populated from a live API response and survives
// import, and that a repeat apply is a no-op.
func TestAccMorpheusInstanceResourceVolumeActualSizeVMware(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	name := acctest.RandomWithPrefix("tfacc-volsize")

	cfg, err := instance.RenderInstanceVMwareConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.0.size", "10"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.1.size", "10"),

					// No image constraint applies here, so what was created
					// matches what was asked for.
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.0.actual_size", "10"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.1.actual_size", "10"),
				),
			},
			{
				Config:   testhelpers.ProviderBlock() + cfg,
				PlanOnly: true,
			},
			{
				ResourceName:      "hpe_morpheus_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Pre-existing import gaps, unrelated to volume sizing.
					"network_interfaces",
					"volumes.0.datastore_auto_selection",
					"volumes.1.datastore_auto_selection",
					// A request, not an observable — see the note above.
					"volumes.0.size",
					"volumes.1.size",
				},
			},
		},
	})
}

// TestAccMorpheusInstanceResourceVolumeActualSizeHVM exercises the case the
// ticket was raised for, where the platform grows the disk and size and
// actual_size genuinely differ.
//
// The image declares minDisk 10GB with a rawSize fractionally above 10GB, so a
// 10GB request is accepted and then rounded up to 11GB. This combination was
// confirmed against a live appliance while diagnosing MORPH-16949.
//
// This test is not currently runnable. The growth band contains no vmdk image,
// so it cannot be expressed on the VMware path, and the HVM example this renders
// from targets an appliance whose images have no backing files and therefore
// cannot provision anything. It is written so the behaviour is covered the
// moment an HVM appliance with usable images is available.
func TestAccMorpheusInstanceResourceVolumeActualSizeHVM(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	name := acctest.RandomWithPrefix("tfacc-volsize-hvm")

	cfg, err := instance.RenderInstanceHVMConfig(t, map[string]string{
		"Name": name,
		// minDisk 10GB, rawSize 10.002GB — accepted, then rounded up to 11GB.
		"ImageName": "dand-ubuntu-20.04-from-iso",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					// The request is preserved verbatim...
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.0.size", "10"),
					// ...while actual_size reports the disk that exists.
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example", "volumes.0.actual_size", "11"),
				),
			},
			// The regression this guards: without the plan modifier the grown
			// volume plans a shrink back to the requested 10GB on every run.
			{
				Config:   testhelpers.ProviderBlock() + cfg,
				PlanOnly: true,
			},
			{
				ResourceName:      "hpe_morpheus_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"network_interfaces",
					"volumes.0.datastore_auto_selection",
					"volumes.0.size",
				},
			},
		},
	})
}
