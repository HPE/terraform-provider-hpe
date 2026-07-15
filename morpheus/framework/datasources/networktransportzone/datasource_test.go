// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktransportzone_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networktransportzone"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkTransportZoneByID(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	networkServerID := os.Getenv("TF_ACC_NETWORK_SERVER_ID")
	transportZoneID := os.Getenv("TF_ACC_TRANSPORT_ZONE_ID")
	if networkServerID == "" || transportZoneID == "" {
		t.Skip("TF_ACC_NETWORK_SERVER_ID and TF_ACC_TRANSPORT_ZONE_ID must be set; skipping test requiring a known NSX-T transport zone")
	}

	providerConfig := testhelpers.ProviderBlock()

	dsConfig, err := networktransportzone.RenderNetworkTransportZoneDataSourceByIDConfig(t, map[string]string{
		"NetworkServerId": networkServerID,
		"Id":              transportZoneID,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_transport_zone.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
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
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	networkServerID := os.Getenv("TF_ACC_NETWORK_SERVER_ID")
	transportZoneName := os.Getenv("TF_ACC_TRANSPORT_ZONE_NAME")
	if networkServerID == "" || transportZoneName == "" {
		t.Skip("TF_ACC_NETWORK_SERVER_ID and TF_ACC_TRANSPORT_ZONE_NAME must be set; skipping test requiring a known NSX-T transport zone")
	}

	providerConfig := testhelpers.ProviderBlock()

	dsConfig, err := networktransportzone.RenderNetworkTransportZoneDataSourceByNameConfig(t, map[string]string{
		"NetworkServerId": networkServerID,
		"Name":            transportZoneName,
	})
	if err != nil {
		t.Fatalf("failed to render data source config: %s", err)
	}

	config := providerConfig + dsConfig
	dsName := "data.hpe_morpheus_network_transport_zone.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
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
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	networkServerID := os.Getenv("TF_ACC_NETWORK_SERVER_ID")
	if networkServerID == "" {
		t.Skip("TF_ACC_NETWORK_SERVER_ID not set; skipping test requiring a known NSX-T network server")
	}

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_transport_zone" "test" {
        network_server_id = ` + networkServerID + `
        name              = "______nonexistent______"
      }`

	expected := regexp.MustCompile(networktransportzone.ErrorNoTransportZone)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}
