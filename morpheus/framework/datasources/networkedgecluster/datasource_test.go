// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkedgecluster_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkedgecluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkEdgeClusterByID(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dsConfig, err := networkedgecluster.RenderNetworkEdgeClusterDataSourceByIDConfig(t, map[string]string{
		"NetworkServerId": os.Getenv("TF_ACC_NETWORK_SERVER_ID"),
		"Id":              os.Getenv("TF_ACC_EDGE_CLUSTER_ID"),
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_edge_cluster.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrSet(dsName, "name"),
					resource.TestCheckResourceAttrSet(dsName, "network_server_id"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkEdgeClusterByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dsConfig, err := networkedgecluster.RenderNetworkEdgeClusterDataSourceByNameConfig(t, map[string]string{
		"NetworkServerId": os.Getenv("TF_ACC_NETWORK_SERVER_ID"),
		"Name":            os.Getenv("TF_ACC_EDGE_CLUSTER_NAME"),
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_edge_cluster.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttrSet(dsName, "name"),
					resource.TestCheckResourceAttrSet(dsName, "network_server_id"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkEdgeClusterNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_edge_cluster" "test" {
        network_server_id = ` + os.Getenv("TF_ACC_NETWORK_SERVER_ID") + `
        name              = "______nonexistent______"
      }`

	expected := regexp.MustCompile(networkedgecluster.ErrorNoEdgeClusterFound)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusNetworkEdgeClusterNoSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_edge_cluster" "test" {
        network_server_id = 1
      }`

	expected := regexp.MustCompile(networkedgecluster.ErrorNoValidSearchTerms)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}
