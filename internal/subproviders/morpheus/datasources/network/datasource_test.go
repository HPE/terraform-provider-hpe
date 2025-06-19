// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const networkResponseJSON = `{
    "network": {
        "id": 123,
        "name": "testacc-TestAccNetworkDataSourceBasic",
        "displayName": "testacc-TestAccNetworkDataSourceBasic",
        "description": "A test network for basic acceptance testing",
        "labels": ["test-label-1", "test-label-2"],
        "tags": [],
        "group": null,
        "zone": null,
        "type": {
            "id": 52,
            "name": "ACI Endpoint Group",
            "code": "aciVxlan"
        },
        "owner": {
            "id": 1,
            "name": "Morpheus QA"
        },
        "ipv4Enabled": true,
        "ipv6Enabled": false,
        "category": "aci.epg.44",
        "cidr": "10.0.0.0/24",
        "visibility": "private",
        "active": true,
        "defaultNetwork": false,
        "subnets": [],
        "tenants": []
    }
}`

const networksListJSON = `{
    "networks": [{
        "id": 123,
        "name": "testacc-TestAccNetworkDataSourceBasic",
        "displayName": "testacc-TestAccNetworkDataSourceBasic",
        "description": "A test network for basic acceptance testing",
        "labels": ["test-label-1", "test-label-2"],
        "tags": [],
        "group": null,
        "zone": null,
        "type": {
            "id": 52,
            "name": "ACI Endpoint Group",
            "code": "aciVxlan"
        },
        "owner": {
            "id": 1,
            "name": "Morpheus QA"
        },
        "ipv4Enabled": true,
        "ipv6Enabled": false,
        "category": "aci.epg.44",
        "cidr": "10.0.0.0/24",
        "visibility": "private",
        "active": true,
        "defaultNetwork": false,
        "subnets": [],
        "tenants": []
    }]
}`

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	httpClient := &http.Client{}
	gock.InterceptClient(httpClient)

	clientFactoryFunc := func(m model.SubModel) *clientfactory.ClientFactory {
		return clientfactory.New(
			m,
			clientfactory.WithFactoryHTTPClient(httpClient),
		)
	}

	providerInstance := provider.New(
		"test",
		morpheus.New(morpheus.WithClientFactory(clientFactoryFunc)),
	)()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer,
	error,
){
	"hpe": newProviderWithError,
}

func TestNetworkDataSourceBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)
	defer gock.Off()

	gock.New("http://net1.test").
		Get("/api/networks($)").
		MatchParam("name", "testacc-TestAccNetworkDataSourceBasic").
		Persist().
		Reply(200).
		SetHeader("Content-Type", "application/json").
		JSON(networksListJSON)

	gock.New("http://net1.test").
		Get("/api/networks/123").
		Persist().
		Reply(200).
		SetHeader("Content-Type", "application/json").
		JSON(networkResponseJSON)

	providerConfig := `
provider "hpe" {
	morpheus {
		url = "http://net1.test"
		access_token = "abc123"
		insecure = true
	}
}
`

	resourceConfig := `
data "hpe_morpheus_network" "test" {
  name = "testacc-TestAccNetworkDataSourceBasic"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"name",
						"testacc-TestAccNetworkDataSourceBasic",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"id",
						"123",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"display_name",
						"testacc-TestAccNetworkDataSourceBasic",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"description",
						"A test network for basic acceptance testing",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"cidr",
						"10.0.0.0/24",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"active",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"labels.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"labels.0",
						"test-label-1",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.test",
						"labels.1",
						"test-label-2",
					),
				),
			},
		},
	})
}
