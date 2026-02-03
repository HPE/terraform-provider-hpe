// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_cluster_hks_hvm/resource.tf cluster_hks_hvm_resource.tf.tmpl Name "tfhvm" ResourcePrefix "vmpre" HostnamePrefix "ospre" Description "Terraform HKS cluster example" CloudId "data.hpe_morpheus_cloud.morpheus_hvm.id" GroupId "data.hpe_morpheus_group.morpheus_lab.id" ClusterLayoutId "1070" PodCidr "172.20.0.0/16" ServiceCidr "172.30.0.0/16" ServerPlanId "data.hpe_morpheus_service_plan.master_nodes.id" ServerResourcePoolId "data.hpe_morpheus_resource_pool.hvm.id" ServerNetworkInterfaceNetworkId "data.hpe_morpheus_network.vm_network.id" ServerStorageRootVolumeRoot "true" ServerStorageRootVolumeSize "20" ServerStorageRootVolumeName "root" ServerStorageRootVolumeStorageType "1" ServerStorageRootVolumeDatastoreId "9" ServerTagsApp "hksmaster" WorkerCount "3" ServerStorageDataVolumeRoot "false" ServerStorageDataVolumeSize "20" ServerStorageDataVolumeName "data" ServerStorageDataVolumeStorageType "1" ServerStorageDataVolumeDatastoreId "2"

// RenderClusterHKSHVMConfig generates a Terraform configuration for the cluster_hks_hvm resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderClusterHKSHVMConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                               "tfhvm",
		"ResourcePrefix":                     "vmpre",
		"HostnamePrefix":                     "ospre",
		"Description":                        "Terraform HKS cluster example",
		"CloudId":                            "1",
		"GroupId":                            "1",
		"ClusterLayoutId":                    "132",
		"PodCidr":                            "172.20.0.0/16",
		"ServiceCidr":                        "172.30.0.0/16",
		"WorkflowId":                         "100",
		"ServerPlanId":                       "23",
		"ServerResourcePoolId":               "1",
		"ServerNetworkInterfaceNetworkId":    "1",
		"ServerStorageRootVolumeRoot":        "true",
		"ServerStorageRootVolumeSize":        "20",
		"ServerStorageRootVolumeName":        "root",
		"ServerStorageRootVolumeStorageType": "1",
		"ServerStorageRootVolumeDatastoreId": "9",
		"ServerTagsApp":                      "hksmaster",
		"WorkerCount":                        "3",
		"ServerStorageDataVolumeRoot":        "false",
		"ServerStorageDataVolumeSize":        "20",
		"ServerStorageDataVolumeName":        "data",
		"ServerStorageDataVolumeStorageType": "1",
		"ServerStorageDataVolumeDatastoreId": "2",
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
	templatePath := filepath.Join(dir, "cluster_hks_hvm_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
