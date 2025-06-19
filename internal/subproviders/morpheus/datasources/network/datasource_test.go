// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/h2non/gock"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
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

func TestNetworkDataSourceExample(t *testing.T) {
	defer testhelpers.RecordResult(t)
	defer gock.Off()

	providerConfig := `
provider "hpe" {
	morpheus {
		url = "http://net1.test"
		access_token = "abc123"
		insecure = true
	}
}
`

	path := "../../../../../examples/data-sources/hpe_morpheus_network/data-source.tf"
	exampleConfig, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading example config: %v", err)
	}

	prefix := fmt.Sprintf("testacc-%s", uuid.New().String())

	networkResponseJSON := `{
    "network": {
        "id": 123,
        "name": "` + prefix + "-example" + `",
        "displayName": "` + prefix + "-example" + `",
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

	networksListJSON := `{
    "networks": [{
        "id": 123,
        "name": "` + prefix + "-example" + `",
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

	gock.New("http://net1.test").
		Get("/api/networks($)").
		MatchParam("name", prefix+"-example").
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + string(exampleConfig),
				ConfigVariables: config.Variables{
					"prefix": config.StringVariable(prefix),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"name",
						prefix+"-example",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"id",
						"123",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"display_name",
						prefix+"-example",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"description",
						"A test network for basic acceptance testing",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"cidr",
						"10.0.0.0/24",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"active",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"labels.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"labels.0",
						"test-label-1",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"labels.1",
						"test-label-2",
					),
				),
			},
		},
	})
	//if !gock.IsDone() {
	//	t.Errorf("Not all gock mocks were called")
	//}
}
