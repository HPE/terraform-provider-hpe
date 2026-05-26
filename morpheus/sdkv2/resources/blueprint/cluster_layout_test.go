// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusClusterLayoutExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Kubernetes) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderClusterLayoutConfig(t, map[string]string{
		"Name":                      name,
		"ProvisionTypeId":           "data.hpe_morpheus_provision_type.example.id",
		"ClusterTypeId":             "data.hpe_morpheus_cluster_type.example.id",
		"WorkerNodePool1NodeTypeId": "data.hpe_morpheus_node_type.example.id",
		"WorkerNodePool2NodeTypeId": "data.hpe_morpheus_node_type.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	// use existing system resources to populate dependency fields
	// provision_type_id and cluster_type_id are required
	resourceConfig += `
	data "hpe_morpheus_provision_type" "example" {
		name = "KVM"
	}

	data "hpe_morpheus_cluster_type" "example" {
		name = "Kubernetes Cluster"
	}

	data "hpe_morpheus_node_type" "example" {
		name = "HVM Kubernetes Worker on Ubuntu 22.04"
	}
	`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"description",
			"Terraform example cluster layout",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"version",
			"1.0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"creatable",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"minimum_memory",
			"4294967296",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_layout.example",
			"cluster_type_id",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_layout.example",
			"provision_type_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"enable_scaling",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"evar.0.name",
			"application",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"evar.0.value",
			"first",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"evar.0.export",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"master_node_pool.0.count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"master_node_pool.0.node_type_id",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"master_node_pool.0.priority_order",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.0.count",
			"4",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.0.node_type_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.0.priority_order",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.1.count",
			"4",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.1.node_type_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_layout.example",
			"worker_node_pool.1.priority_order",
			"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
