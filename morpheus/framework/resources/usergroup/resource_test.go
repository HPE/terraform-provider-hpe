// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package usergroup_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusUserGroupBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_user_group" "test" {
  name      = "` + name + `"
  sudo_user = false
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_user_group.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_user_group.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_user_group.test", "sudo_user", "false"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_user_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusUserGroupUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_user_group" "test" {
  name      = "` + name + `"
  sudo_user = false
}
`

	updateConfig := `
resource "hpe_morpheus_user_group" "test" {
  name      = "` + name + `"
  sudo_user = true
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_user_group.test", "sudo_user", "false"),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_user_group.test", "sudo_user", "true"),
				),
			},
		},
	})
}
