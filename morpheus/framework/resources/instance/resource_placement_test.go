// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Two provisioned hosts in the HVM cluster, supplied by the environment, are
// required to exercise the live-migration path (update preferred host).
const (
	envHVMHostID  = "TF_ACC_HVM_HOST_ID"
	envHVMHostID2 = "TF_ACC_HVM_HOST_ID_2"
)

// placementHVMConfig renders an HVM instance with a config_hvm placement block.
// placementLines is the placement_strategy / preferred_host_id lines to inject
// into config_hvm, and extraTimeouts is any additional timeouts entries (e.g.
// an update timeout for the migration step).
func placementHVMConfig(t *testing.T, label, name, placementLines, extraTimeouts string) string {
	t.Helper()

	cloudName := testhelpers.HVMCloudName(t)
	clusterName := testhelpers.HVMClusterName(t)
	networkName := testhelpers.HVMNetworkName(t)
	planName := testhelpers.HVMPlanName(t)
	layoutName, layoutVersion := testhelpers.HVMLayout(t)

	imageData, imageLine := "", ""
	if image := os.Getenv(testhelpers.EnvHVMImageName); image != "" {
		imageData = `
data "hpe_morpheus_image" "hvm" {
  name = "` + image + `"
}
`
		imageLine = "\n    image_id          = data.hpe_morpheus_image.hvm.id"
	}

	// A live migration between hosts needs the VM's volumes on shared storage;
	// automatic selection may pick a local datastore, so name the datastore
	// when the lab has a shared one. Falls back to automatic selection.
	datastoreData, datastoreLine := "", `      datastore_auto_selection = "auto"`
	if ds := os.Getenv(testhelpers.EnvHVMDatastoreName); ds != "" {
		datastoreData = `
data "hpe_morpheus_datastore" "hvm" {
  name       = "` + ds + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
`
		datastoreLine = `      datastore_id             = data.hpe_morpheus_datastore.hvm.id`
	}

	return `
data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "` + cloudName + `"
}

data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "` + layoutName + `"
  version = "` + layoutVersion + `"
}

data "hpe_morpheus_network" "hvm" {
  name       = "` + networkName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_service_plan" "hvm" {
  name                = "` + planName + `"
  provision_type_code = "kvm"
}
` + imageData + datastoreData + `
resource "hpe_morpheus_group" "test" {
  name      = "` + name + `-grp"
  cloud_ids = [data.hpe_morpheus_cloud.hvm_cloud.id]
}

resource "hpe_morpheus_instance" "placement" {
  name             = "` + name + `"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.hvm.instance_type.id
  group_id         = hpe_morpheus_group.test.id
  plan_id          = data.hpe_morpheus_service_plan.hvm.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.hvm.id
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
` + datastoreLine + `
    }
  ]

  tags = [
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_hvm = {
    resource_pool_id  = "pool-${data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id}"
    no_agent          = ` + fmt.Sprintf("%t", testhelpers.HVMNoAgent()) + `
    create_user       = false` + imageLine + `
` + placementLines + `
  }

  timeouts = {
    create = "1h"
    delete = "20m"` + extraTimeouts + `
  }
}
`
}

// TestAccMorpheusInstanceHvmPlacementOk provisions an HVM instance pinned to a
// first host, then updates the preferred host to a second (exercising the live
// migration path with an update timeout), then switches the strategy to auto
// (releasing the pin), and finally imports the instance.
//
// It requires two provisioned hosts in the HVM cluster, supplied via
// TF_ACC_HVM_HOST_ID and TF_ACC_HVM_HOST_ID_2.
func TestAccMorpheusInstanceHvmPlacementOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	host1 := os.Getenv(envHVMHostID)
	host2 := os.Getenv(envHVMHostID2)
	if host1 == "" || host2 == "" {
		t.Skipf("skipping: %s and %s must both be set to two provisioned hosts",
			envHVMHostID, envHVMHostID2)
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_instance.placement"

	pinnedHost1 := fmt.Sprintf(
		"    kvm_host_id        = %s\n    placement_strategy = \"pinned\"\n    preferred_host_id  = %s", host1, host1)
	pinnedHost2 := fmt.Sprintf(
		"    kvm_host_id        = %s\n    placement_strategy = \"pinned\"\n    preferred_host_id  = %s", host1, host2)
	autoStrategy := fmt.Sprintf("    kvm_host_id        = %s\n    placement_strategy = \"auto\"", host1)

	// Moving a running VM between hosts needs shared storage: Morpheus refuses
	// a live move for a VM on local storage ("Server must be powered off to
	// move hosts because it is on local storage"). The re-pin to a second host
	// is therefore opt-in for labs whose datastore is shared.
	moveStep := []resource.TestStep{}
	if os.Getenv("TF_ACC_HVM_SHARED_STORAGE") != "" {
		moveStep = append(moveStep, resource.TestStep{
			// Update preferred host to host 2 (live migration). Allow a
			// longer update window for the migration job.
			Config: providerConfig + placementHVMConfig(
				t, "placement", name, pinnedHost2, "\n    update = \"30m\""),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					resourceName, "config_hvm.placement_strategy", "pinned"),
				resource.TestCheckResourceAttr(
					resourceName, "config_hvm.preferred_host_id", host2),
				// The stored preference alone would pass with a failed or
				// still-pending migration: assert the VM actually runs on
				// host 2 now.
				checkInstanceOnHost(resourceName, host2),
			),
		})
	}

	steps := []resource.TestStep{
		{
			// Create pinned to host 1.
			Config: providerConfig + placementHVMConfig(t, "placement", name, pinnedHost1, ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					resourceName, "config_hvm.placement_strategy", "pinned"),
				resource.TestCheckResourceAttr(
					resourceName, "config_hvm.preferred_host_id", host1),
			),
		},
	}
	steps = append(steps, moveStep...)
	steps = append(steps,
		resource.TestStep{
			// Switch to auto, releasing the host preference.
			Config: providerConfig + placementHVMConfig(t, "placement", name, autoStrategy, ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					resourceName, "config_hvm.placement_strategy", "auto"),
			),
		},
		resource.TestStep{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"compute_servers",
				"config",
				"connection_info",
				"labels",
				"network_interfaces",
				// On import tags is left unset: the resource cannot tell which
				// of Morpheus's tags the configuration manages (tags_all has
				// them all).
				"tags",
				"timeouts",
				"volumes",
			},
		},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

// checkInstanceOnHost asserts that every node server of the instance reports
// the given host as its current parent, by reading each server through the API.
func checkInstanceOnHost(resourceName, hostID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not in state", resourceName)
		}

		want, err := strconv.ParseInt(hostID, 10, 64)
		if err != nil {
			return err
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			return fmt.Errorf("client: %w", err)
		}

		// compute_servers is a set of the node server ids.
		n, _ := strconv.Atoi(rs.Primary.Attributes["compute_servers.#"])
		if n == 0 {
			return fmt.Errorf("%s has no compute_servers in state", resourceName)
		}

		for i := range n {
			sid, err := strconv.ParseInt(rs.Primary.Attributes[fmt.Sprintf("compute_servers.%d", i)], 10, 64)
			if err != nil {
				return err
			}

			resp, _, err := client.HostsAPI.GetHost(ctx, sdk.GetHostIdParameter{Int64: &sid}).Execute()
			if err != nil || resp == nil || resp.Server == nil {
				return fmt.Errorf("server %d: GET failed: %w", sid, err)
			}

			if resp.Server.ParentServer == nil || resp.Server.ParentServer.Id == nil {
				return fmt.Errorf("server %d reports no parent host", sid)
			}

			if *resp.Server.ParentServer.Id != want {
				return fmt.Errorf("server %d is on host %d, want %d", sid, *resp.Server.ParentServer.Id, want)
			}
		}

		return nil
	}
}
