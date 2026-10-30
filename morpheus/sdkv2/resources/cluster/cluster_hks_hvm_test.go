// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusClusterHKSHVMExampleOk provisions a whole HKS cluster -- one master
// and three workers on 8 GB plans -- from the example's QA-appliance ids (cloud 1,
// layout 132, plan 23, pool 1, network 1, datastores 2 and 9). That is more than an
// appliance that merely has a usable Kubernetes cluster (the kubernetes_cluster
// capability, which the namespace tests need) can be assumed to afford, and the ids
// exist nowhere else, so the test is opt-in like the other cluster-building tests:
// set RUN_SKIPPED_BY_DEFAULT to run it.
func TestAccMorpheusClusterHKSHVMExampleOk(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.KubernetesCluster)

	if skip.SkipByDefault(t) {
		t.Log("Skipping test not opted in via RUN_SKIPPED_BY_DEFAULT")

		return
	}

	defer testhelpers.RecordResult(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	prefix := strings.ToLower(name)

	resourceConfig, err := cluster.RenderClusterHKSHVMConfig(t, map[string]string{
		"Name":                               name,
		"WorkflowId":                         "null",
		"ServerStorageDataVolumeDatastoreId": "1",
		"ServerStorageRootVolumeDatastoreId": "1",
		"ResourcePrefix":                     prefix,
		"HostnamePrefix":                     prefix,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"resource_prefix",
			prefix,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"hostname_prefix",
			prefix,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"description",
			"Terraform HKS cluster example",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_hvm.example",
			"cloud_id",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_hvm.example",
			"group_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"cluster_layout_id",
			"132",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"pod_cidr",
			"172.20.0.0/16",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"service_cidr",
			"172.30.0.0/16",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.plan_id",
			"23",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.resource_pool_id",
			"1",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.network_interface.0.network_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.storage_volume.0.root",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.storage_volume.0.size",
			"20",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.storage_volume.0.name",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.storage_volume.0.storage_type",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.storage_volume.0.datastore_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"server.0.tags.app",
			"hksmaster",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_hks_hvm.example",
			"workers",
			"3",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
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
