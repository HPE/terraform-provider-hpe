// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktransportzone_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networktransportzone"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkTransportZoneByID(t *testing.T) {
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

	dsConfig, err := networktransportzone.RenderNetworkTransportZoneDataSourceByIDConfig(t, map[string]string{
		"NetworkServerId": os.Getenv("TF_ACC_NETWORK_SERVER_ID"),
		"Id":              os.Getenv("TF_ACC_TRANSPORT_ZONE_ID"),
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_transport_zone.example"

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

func TestAccMorpheusNetworkTransportZoneByName(t *testing.T) {
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

	dsConfig, err := networktransportzone.RenderNetworkTransportZoneDataSourceByNameConfig(t, map[string]string{
		"NetworkServerId": os.Getenv("TF_ACC_NETWORK_SERVER_ID"),
		"Name":            os.Getenv("TF_ACC_TRANSPORT_ZONE_NAME"),
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_transport_zone.example"

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

func TestAccMorpheusNetworkTransportZoneNotFound(t *testing.T) {
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
      data "hpe_morpheus_network_transport_zone" "test" {
        network_server_id = ` + os.Getenv("TF_ACC_NETWORK_SERVER_ID") + `
        name              = "______nonexistent______"
      }`

	expected := regexp.MustCompile(networktransportzone.ErrorNoTransportZone)

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

func TestAccMorpheusNetworkTransportZoneNoSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_transport_zone" "test" {
        network_server_id = 1
      }`

	expected := regexp.MustCompile(networktransportzone.ErrorNoValidSearchTerms)

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
