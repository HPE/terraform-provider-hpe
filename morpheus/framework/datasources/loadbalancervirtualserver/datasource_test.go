// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestAccMorpheusLoadBalancerVirtualServerDataSourceByIdExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerDataSourceByIDConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"id",
					),
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"vip_name",
					),
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"load_balancer_id",
					),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerDataSourceByNameExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerDataSourceByNameConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"id",
					),
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"vip_name",
					),
					resource.TestCheckResourceAttrSet(
						"data.hpe_morpheus_load_balancer_virtual_server.example",
						"load_balancer_id",
					),
				),
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerDataSourceNotFound(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerDataSourceByNameConfig(t,
		map[string]string{
			"VipName": "______",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(loadbalancervirtualserver.ErrorNoLoadBalancerVirtualServer)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerDataSourceNoSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_load_balancer_virtual_server" "test" {
        load_balancer_id = 1
      }`

	expected := regexp.MustCompile(loadbalancervirtualserver.ErrorNoValidSearchTerms)

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

func TestAccMorpheusLoadBalancerVirtualServerDataSourceBothSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_load_balancer_virtual_server" "test" {
        id               = 1
        vip_name         = "______"
        load_balancer_id = 1
      }`

	expected := loadbalancervirtualserver.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
