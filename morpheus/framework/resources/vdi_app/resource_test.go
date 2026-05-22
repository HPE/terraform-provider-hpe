// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package vdi_app_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusVdiAppBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig := `
resource "hpe_morpheus_vdi_app" "test" {
  name          = "` + name + `"
  launch_prefix = "||/usr/bin/app"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_vdi_app.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_app.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_app.test", "launch_prefix", "||/usr/bin/app"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_vdi_app.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusVdiAppUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig := `
resource "hpe_morpheus_vdi_app" "test" {
  name          = "` + name + `"
  launch_prefix = "||/usr/bin/app"
}
`

	updateConfig := `
resource "hpe_morpheus_vdi_app" "test" {
  name          = "` + name + `"
  launch_prefix = "||/usr/bin/app"
  description   = "Updated description"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_app.test", "name", name),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_app.test", "description", "Updated description"),
				),
			},
		},
	})
}
