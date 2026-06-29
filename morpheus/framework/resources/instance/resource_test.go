// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusInstanceResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Skip("Skipping test due to missing capabilities")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "34"
	resourcePool := "pool-1"

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":         name,
		"InstanceType": instanceTypeID,
		"ResourcePool": resourcePool,
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
			"config_hvm.resource_pool_id",
			resourcePool,
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.All) {
		t.Skip("Skipping test due to missing capabilities")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// user_group and storage_profile use hard-coded values that exist in the
	// reference test environment, consistent with the other instance tests
	// (resource pool, layout, network, datastore). kvm-cache-none is a standard
	// KVM/HVM storage profile code.
	userGroup := "1"
	storageProfile := "kvm-cache-none"

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":           name,
		"UserGroup":      userGroup,
		"StorageProfile": storageProfile,
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
			"user_group",
			userGroup,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"volumes.0.storage_profile",
			storageProfile,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.Azure) {
		t.Skip("Skipping test due to missing capabilities")
	}

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.Azure) {
		t.Skip("Skipping test due to missing capabilities")
	}

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.All) {
		t.Skip("Skipping test due to missing capabilities")
	}

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.All) {
		t.Skip("Skipping test due to missing capabilities")
	}

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	if capabilities.Missing(t, capabilities.All) {
		t.Skip("Skipping test due to missing capabilities")
	}

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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
