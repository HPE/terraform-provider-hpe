// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_cluster_hks_vsphere/resource.tf cluster_hks_vsphere_resource.tf.tmpl Name "tfvsphere" ResourcePrefix "vmpre" HostnamePrefix "ospre" Description "Terraform HKS cluster example" CloudId "data.hpe_morpheus_cloud.morpheus_vsphere.id" GroupId "data.hpe_morpheus_group.morpheus_lab.id" ClusterLayoutId "1070" PodCidr "172.20.0.0/16" ServiceCidr "172.30.0.0/16" ClusterRepoAccountId "1" MasterNodePoolPlanId "data.hpe_morpheus_service_plan.master_nodes.id" MasterNodePoolResourcePoolId "data.hpe_morpheus_resource_pool.vsphere_resource_pool.id" MasterNodePoolNetworkInterfaceNetworkId "data.hpe_morpheus_network.vm_network.id" MasterNodePoolStorageVolumeRoot "true" MasterNodePoolStorageVolumeSize "20" MasterNodePoolStorageVolumeName "root" MasterNodePoolStorageVolumeStorageType "1" MasterNodePoolStorageVolumeDatastoreId "9" MasterNodePoolTagsApp "hksmaster" WorkerNodePoolCount "3" WorkerNodePoolPlanId "data.hpe_morpheus_service_plan.worker_nodes.id" WorkerNodePoolResourcePoolId "data.hpe_morpheus_resource_pool.vsphere_resource_pool.id" WorkerNodePoolNetworkInterfaceNetworkId "data.hpe_morpheus_network.vm_network.id" WorkerNodePoolStorageVolume0Root "true" WorkerNodePoolStorageVolume0Size "20" WorkerNodePoolStorageVolume0Name "root" WorkerNodePoolStorageVolume0StorageType "1" WorkerNodePoolStorageVolume0DatastoreId "2" WorkerNodePoolStorageVolume1Root "false" WorkerNodePoolStorageVolume1Size "20" WorkerNodePoolStorageVolume1Name "data" WorkerNodePoolStorageVolume1StorageType "1" WorkerNodePoolStorageVolume1DatastoreId "2" WorkerNodePoolTagsApp "hksworker"

// RenderClusterHKSVsphereConfig generates a Terraform configuration for the cluster_hks_vsphere resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderClusterHKSVsphereConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                         "tfvsphere",
		"ResourcePrefix":               "vmpre",
		"HostnamePrefix":               "ospre",
		"Description":                  "Terraform HKS cluster example",
		"CloudId":                      "data.hpe_morpheus_cloud.morpheus_vsphere.id",
		"GroupId":                      "data.hpe_morpheus_group.morpheus_lab.id",
		"ClusterLayoutId":              "1070",
		"PodCidr":                      "172.20.0.0/16",
		"ServiceCidr":                  "172.30.0.0/16",
		"WorkflowId":                   "100",
		"ClusterRepoAccountId":         "1",
		"MasterNodePoolPlanId":         "data.hpe_morpheus_service_plan.master_nodes.id",
		"MasterNodePoolResourcePoolId": "data.hpe_morpheus_resource_pool.vsphere_resource_pool.id",
		"MasterNodePoolNetworkInterfaceNetworkId": "data.hpe_morpheus_network.vm_network.id",
		"MasterNodePoolStorageVolumeRoot":         "true",
		"MasterNodePoolStorageVolumeSize":         "20",
		"MasterNodePoolStorageVolumeName":         "root",
		"MasterNodePoolStorageVolumeStorageType":  "1",
		"MasterNodePoolStorageVolumeDatastoreId":  "9",
		"MasterNodePoolTagsApp":                   "hksmaster",
		"WorkerNodePoolCount":                     "3",
		"WorkerNodePoolPlanId":                    "data.hpe_morpheus_service_plan.worker_nodes.id",
		"WorkerNodePoolResourcePoolId":            "data.hpe_morpheus_resource_pool.vsphere_resource_pool.id",
		"WorkerNodePoolNetworkInterfaceNetworkId": "data.hpe_morpheus_network.vm_network.id",
		"WorkerNodePoolStorageVolume0Root":        "true",
		"WorkerNodePoolStorageVolume0Size":        "20",
		"WorkerNodePoolStorageVolume0Name":        "root",
		"WorkerNodePoolStorageVolume0StorageType": "1",
		"WorkerNodePoolStorageVolume0DatastoreId": "2",
		"WorkerNodePoolStorageVolume1Root":        "false",
		"WorkerNodePoolStorageVolume1Size":        "20",
		"WorkerNodePoolStorageVolume1Name":        "data",
		"WorkerNodePoolStorageVolume1StorageType": "1",
		"WorkerNodePoolStorageVolume1DatastoreId": "2",
		"WorkerNodePoolTagsApp":                   "hksworker",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "cluster_hks_vsphere_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
