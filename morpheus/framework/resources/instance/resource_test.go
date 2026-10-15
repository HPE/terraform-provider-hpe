// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/containerip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusInstanceResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		// The instance type is the layout's, and the resource pool the
		// cluster's: both resolved by the configuration's data sources.
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_instance.example", "instance_type_id",
			"data.hpe_morpheus_instance_type_layout.vme", "instance_type.id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_instance.example", "layout_id",
			"data.hpe_morpheus_instance_type_layout.vme", "id",
		),
		resource.TestCheckResourceAttrWith(
			"hpe_morpheus_instance.example", "config_hvm.resource_pool_id",
			func(v string) error {
				if !strings.HasPrefix(v, "pool-") {
					return fmt.Errorf("expected the cluster's pool as pool-<id>, got %q", v)
				}

				return nil
			},
		),
		// status is a computed attribute, refreshed on read so an out-of-band
		// deletion of the underlying VM surfaces as a change on the next plan.
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_instance.example",
			"status",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}

// TestAccMorpheusInstanceResourceUserGroupStorageProfile provisions an HVM/KVM
// instance with a user_group and a volume storage_profile and asserts both
// round-trip into state. user_group is a provision-time, RequiresReplace input
// read back from the instance config; storage_profile is a write-mostly volume
// input preserved from the plan on read.
func TestAccMorpheusInstanceResourceUserGroupStorageProfile(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// The user group is created for the run and referenced, so no id is
	// assumed. kvm-cache-none is a standard KVM/HVM storage profile code.
	storageProfile := "kvm-cache-none"

	userGroupConfig, err := usergroup.RenderUserGroupConfig(t, map[string]string{
		"Name":        name,
		"Description": "user group for " + name,
		"UserIds":     "[]", // no members: the example's default references a whoami data source this config does not declare
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":           name,
		"UserGroup":      "hpe_morpheus_user_group.example.id",
		"StorageProfile": storageProfile,
	})
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = userGroupConfig + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_instance.example", "user_group",
			"hpe_morpheus_user_group.example", "id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"volumes.0.storage_profile",
			storageProfile,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusInstanceResourceAzureExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Azure)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "34"
	resourcePool := "pool-1"

	resourceConfig, err := instance.RenderInstanceAzureConfig(t, map[string]string{
		"Name":         name,
		"InstanceType": instanceTypeID,
		"ResourcePool": resourcePool,
		"AzureRegion":  "eastus",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"instance_type_id",
			instanceTypeID,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.resource_pool_id",
			resourcePool,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.create_user",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.azure_region",
			"eastus",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusInstanceResourceAzureSubnet(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Azure)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "34"
	resourcePool := "pool-1"
	subnetID := "1"

	resourceConfig, err := instance.RenderInstanceAzureSubnetConfig(t, map[string]string{
		"Name":         name,
		"InstanceType": instanceTypeID,
		"ResourcePool": resourcePool,
		"AzureRegion":  "eastus",
		"SubnetId":     subnetID,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"instance_type_id",
			instanceTypeID,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.resource_pool_id",
			resourcePool,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.create_user",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config_azure.azure_region",
			"eastus",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"network_interfaces.0.subnet_id",
			subnetID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}

func TestAccMorpheusInstanceResourceUpdateName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedResourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"name",
						name,
					),
				),
			},
			{
				Config: providerConfig + updatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"name",
						updatedName,
					),
				),
			},
		},
	})
}

func TestAccMorpheusInstanceResourceUpdateInstanceContext(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":            name,
		"InstanceContext": "dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedResourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":            name,
		"InstanceContext": "production",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"instance_context",
						"dev",
					),
				),
			},
			{
				Config: providerConfig + updatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"instance_context",
						"production",
					),
				),
			},
		},
	})
}

func TestAccMorpheusInstanceResourceUpdateTags(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedResourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":         name,
		"MultipleTags": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"tags.#",
						"1",
					),
				),
			},
			{
				Config: providerConfig + updatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"tags.#",
						"5",
					),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourceWaitForIPAddress provisions an HVM instance
// with wait_for_ip_address = true and verifies that connection_info contains a
// real address after create. A layout with the agent disabled does report an
// address (the address comes from the platform, not from anything inside the
// guest), so the standard layout is used.
func TestAccMorpheusInstanceResourceWaitForIPAddress(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// A second instance beside the rendered one, with wait_for_ip_address set,
	// built from the same data sources the rendered configuration declares
	// (cloud, cluster, layout, network, plan, and the image when the fixture
	// names one -- a layout without an image of its own provisions a VM with
	// nothing to boot, which never reports an address) and its per-run group
	// (MORPH-16458: not the master group, which a sub-tenant cannot use).
	imageLine := ""
	if os.Getenv(testhelpers.EnvHVMImageName) != "" {
		imageLine = "\n    image_id         = data.hpe_morpheus_image.vme.id"
	}

	resourceConfig += `
resource "hpe_morpheus_instance" "wait_ip" {
  name             = "` + name + `-waitip"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.vme.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.vme.instance_type.id
  group_id         = data.hpe_morpheus_group.example.id
  plan_id          = data.hpe_morpheus_service_plan.vme_plan.id

  instance_context = "dev"

  wait_for_ip_address = true

  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.vme.id
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    },
  ]

  # The agent is what reports the guest's address, so this instance always
  # installs it, whatever the fixture says for the others.
  config_hvm = {
    resource_pool_id = "pool-${data.hpe_morpheus_cluster.vme.permissions.resource_pool.id}"
    no_agent         = false
    create_user      = false` + imageLine + `
  }
}
`

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.wait_ip",
			"wait_for_ip_address",
			"true",
		),
		// The address must be a real one, not a placeholder. Checking only that
		// the attribute is set would pass even if the wait did nothing: a
		// container that has not reported yet returns the sentinel 0.0.0.0
		// rather than an absent value, so a presence check is satisfied
		// immediately and proves nothing.
		resource.TestCheckResourceAttrWith(
			"hpe_morpheus_instance.wait_ip",
			"connection_info.0",
			func(value string) error {
				if !containerip.Ready(value) {
					return fmt.Errorf(
						"connection_info.0 is %q, which is a placeholder or empty; "+
							"wait_for_ip_address should have waited for a real address",
						value,
					)
				}

				return nil
			},
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checks,
			},
		},
	})
}
