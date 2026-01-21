// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/cluster"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccClusterHKSVsphereExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := cluster.RenderClusterHKSVsphereConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig += `
	data "hpe_morpheus_group" "morpheus_lab" {
		name = "QA VMware"
	}

	# we get multiple clouds on name search
	data "hpe_morpheus_cloud" "morpheus_vsphere" {
		id = 2
	}

	# we get multiple resource pools with name pool-2 on name search
	data "hpe_morpheus_resource_pool" "vsphere_resource_pool" {
		name     = "pool-2"
		cloud_id = data.hpe_morpheus_cloud.morpheus_vsphere.id
		id = 5
	}

	# we get multiple networks on name search
	data "hpe_morpheus_network" "vm_network" {
		id = 86657
	}

	data "hpe_morpheus_service_plan" "master_nodes" {
		name                = "2 CPU, 8GB Memory"
		provision_type_code = "vmware"
	}

	data "hpe_morpheus_service_plan" "worker_nodes" {
		name                = "2 CPU, 8GB Memory"
		provision_type_code = "vmware"
	}
	`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"resource_prefix",
			"vmpre",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"hostname_prefix",
			"ospre",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"description",
			"Terraform HKS cluster example",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"cloud_id",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"group_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"cluster_layout_id",
			"1070",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"pod_cidr",
			"172.20.0.0/16",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"service_cidr",
			"172.30.0.0/16",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"cluster_repo_account_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.plan_id",
			"244",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.resource_pool_id",
			"5",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.network_interface.0.network_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.storage_volume.0.root",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.storage_volume.0.size",
			"20",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.storage_volume.0.name",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.storage_volume.0.storage_type",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.storage_volume.0.datastore_id",
			"9",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"master_node_pool.0.tags.app",
			"hksmaster",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.count",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.plan_id",
			"244",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.resource_pool_id",
			"5",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.network_interface.0.network_id",
			"86657",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.0.root",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.0.size",
			"20",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.0.name",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.0.storage_type",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.0.datastore_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.1.root",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.1.size",
			"20",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.1.name",
			"data",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.1.storage_type",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.storage_volume.1.datastore_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_vsphere.example",
			"worker_node_pool.0.tags.app",
			"hksworker",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
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
