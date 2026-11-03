// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceCloneResourceWriteOnlyOk runs the write-only matrix on
// source_instance_id: a change of source replaces the clone, the deprecated
// marker still replaces when bumped and does nothing when removed; both import
// forms leave a clean plan. Needs two source instances; the test creates both
// from the HVM fixture (see cloneSourceConfig) rather than cloning instances
// named from outside, so nothing pre-existing is cloned or raced.
func TestAccMorpheusInstanceCloneResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix("tfwo-clone")

	// Two sources: the rendered hpe_morpheus_instance.example and a second
	// instance beside it, built from the same data sources the rendered
	// configuration declares.
	sources := cloneSourceConfig(t, name+"-src1") + secondSourceConfig(name+"-src2")

	datastoreLine := ""
	if os.Getenv(testhelpers.EnvHVMDatastoreName) != "" {
		datastoreLine = ", datastore_id = " + instance.ExampleDatastoreData + ".id"
	}

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_instance_clone.wo",
		Attribute: "source_instance_id",
		Marker:    "source_instance_id_version",
		// The source instance is in the URL (/api/instances/{id}/clone), not the
		// body: any write to the clone path counts as "sent".
		APIPath: "/clone",
		APIKey:  "",
		// The clone updates only its plan and volume sizes in place, so the
		// matrix's unrelated-change step is left out (unrelated is ignored).
		Config: func(value, extra string, _ bool) string {
			return testhelpers.ProviderBlock() + sources + `
resource "hpe_morpheus_instance_clone" "wo" {
  source_instance_id = ` + value + extra + `
  name               = "` + name + `"

  volumes = [{ name = "root", size = 20, root_volume = true` + datastoreLine + ` }]

  network_interfaces = [{ network_id = ` + instance.ExampleNetworkData + `.id, ip_mode = "" }]
}
`
		},
		// The values are references to the two source resources, so the
		// import digest has to be taken of the id the second source got, read
		// from state, not of the reference text.
		ValueFromState: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources["hpe_morpheus_instance.src2"]
			if !ok {
				return "", fmt.Errorf("hpe_morpheus_instance.src2 not in state")
			}

			return rs.Primary.ID, nil
		},
	}

	rec := testhelpers.RecordRequests(t)

	src1, src2 := "hpe_morpheus_instance.example.id", "hpe_morpheus_instance.src2.id"
	steps := testhelpers.WriteOnlyReplaceSteps(c, rec, src1, src2)
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, src2, "volumes", "network_interfaces")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

// secondSourceConfig is a second HVM instance beside the rendered one, on the
// same data sources (cloud, cluster, layout, network, plan, image when the
// fixture names one) and the same per-run group.
func secondSourceConfig(name string) string {
	imageLine := ""
	if os.Getenv(testhelpers.EnvHVMImageName) != "" {
		imageLine = "\n    image_id              = " + instance.ExampleImageData + ".id"
	}

	volume := `
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"`
	if os.Getenv(testhelpers.EnvHVMDatastoreName) != "" {
		volume = `
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = ` + instance.ExampleDatastoreData + `.id`
	}

	return `
resource "hpe_morpheus_instance" "src2" {
  name             = "` + name + `"
  cloud_id         = ` + instance.ExampleCloudData + `.id
  layout_id        = ` + instance.ExampleLayoutData + `.id
  instance_type_id = ` + instance.ExampleLayoutData + `.instance_type.id
  group_id         = ` + instance.ExampleGroupData + `.id
  plan_id          = ` + instance.ExamplePlanData + `.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = ` + instance.ExampleNetworkData + `.id
    }
  ]

  volumes = [
    {` + volume + `
    },
  ]

  tags = [
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_hvm = {
    resource_pool_id      = "pool-${` + instance.ExampleClusterData + `.permissions.resource_pool.id}"
    nested_virtualization = "off"
    no_agent              = ` + fmt.Sprintf("%t", testhelpers.HVMNoAgent()) + `
    create_user           = false` + imageLine + `
  }
}
`
}
