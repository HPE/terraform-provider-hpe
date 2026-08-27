// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusInstanceNodeVirtual tests adding a node to a virtual
// (non-metal) instance without resource_pool_id. This is the primary
// use case: scale-out on any instance type.
//
// Not parallel: this test provisions a real instance and scales it.
// Running it concurrently with TestAccMorpheusInstanceNodePoolOnVirtualError
// has been observed to exceed appliance capacity, causing intermittent
// failures where the add-node action is silently refused.
func TestAccMorpheusInstanceNodeVirtual(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	instanceID, err := testhelpers.CreateInstance(t)
	if err != nil {
		t.Fatalf("failed to create parent instance: %v", err)
	}
	t.Cleanup(func() {
		if err := testhelpers.DeleteInstance(t, instanceID); err != nil {
			t.Logf("WARNING: failed to delete instance %d: %v", instanceID, err)
		}
	})

	providerConfig := testhelpers.ProviderBlock()

	config := fmt.Sprintf(`
resource "hpe_morpheus_instance_node" "test_virtual" {
  instance_id = %d
}
`, instanceID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_instance_node.test_virtual",
						"container_id",
					),
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_instance_node.test_virtual",
						"server_id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceNodeMetalPool tests adding a node to a bare-metal
// instance in a specific resource pool.
func TestAccMorpheusInstanceNodeMetalPool(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Metal)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	instanceID := os.Getenv("TF_VAR_metal_instance_id")
	poolID := os.Getenv("TF_VAR_metal_pool_id")

	if instanceID == "" || poolID == "" {
		t.Skip("TF_VAR_metal_instance_id and TF_VAR_metal_pool_id not set")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := fmt.Sprintf(`
resource "hpe_morpheus_instance_node" "test_metal" {
  instance_id      = %s
  resource_pool_id = %s
}
`, instanceID, poolID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_instance_node.test_metal",
						"container_id",
					),
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_instance_node.test_metal",
						"server_id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceNodePoolOnVirtualError tests that setting
// resource_pool_id on a non-metal instance produces an error.
//
// Not parallel: this test provisions a real instance and scales it.
// Running it concurrently with TestAccMorpheusInstanceNodeVirtual
// has been observed to exceed appliance capacity, causing intermittent
// failures where the add-node action is silently refused.
func TestAccMorpheusInstanceNodePoolOnVirtualError(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	instanceID, err := testhelpers.CreateInstance(t)
	if err != nil {
		t.Fatalf("failed to create parent instance: %v", err)
	}
	t.Cleanup(func() {
		if err := testhelpers.DeleteInstance(t, instanceID); err != nil {
			t.Logf("WARNING: failed to delete instance %d: %v", instanceID, err)
		}
	})

	providerConfig := testhelpers.ProviderBlock()

	// Setting resource_pool_id on a virtual instance must fail.
	config := fmt.Sprintf(`
resource "hpe_morpheus_instance_node" "test_guard" {
  instance_id      = %d
  resource_pool_id = 999
}
`, instanceID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + config,
				ExpectError: regexp.MustCompile(`resource_pool_id is only valid for bare-metal`),
			},
		},
	})
}

// TestAccMorpheusInstanceNodeWaitForIP tests wait_for_ip_address=true.
func TestAccMorpheusInstanceNodeWaitForIP(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Metal)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	instanceID := os.Getenv("TF_VAR_metal_instance_id")
	poolID := os.Getenv("TF_VAR_metal_pool_id")

	if instanceID == "" || poolID == "" {
		t.Skip("TF_VAR_metal_instance_id and TF_VAR_metal_pool_id not set")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := fmt.Sprintf(`
resource "hpe_morpheus_instance_node" "test_ip" {
  instance_id         = %s
  resource_pool_id    = %s
  wait_for_ip_address = true
}
`, instanceID, poolID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_instance_node.test_ip",
						"ip_address",
					),
				),
			},
		},
	})
}
