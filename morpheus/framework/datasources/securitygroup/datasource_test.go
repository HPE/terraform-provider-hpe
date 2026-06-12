// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSecurityGroupDataSourceByNameExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := securitygroup.RenderSecurityGroupDataSourceByNameConfig(t, map[string]string{"Name":"MorpheusUbuntu-nsg"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_security_group.example", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_security_group.example", "name", "MorpheusUbuntu-nsg"),
				),
			},
		},
	})
}

func TestAccMorpheusSecurityGroupDataSourceByIdExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := securitygroup.RenderSecurityGroupDataSourceByIDConfig(t, map[string]string{"Id":"31"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_security_group.example", "name"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_security_group.example", "id", "31"),
				),
			},
		},
	})
}
