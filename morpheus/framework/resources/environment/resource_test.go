// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package environment_test

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

func TestAccMorpheusEnvironmentBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")
	code := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_environment" "test" {
  name        = "` + name + `"
  code        = "` + code + `"
  description = "Initial description"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_environment.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_environment.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_environment.test", "code", code),
					resource.TestCheckResourceAttr("hpe_morpheus_environment.test", "description", "Initial description"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusEnvironmentUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")
	code := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_environment" "test" {
  name        = "` + name + `"
  code        = "` + code + `"
  description = "Initial description"
}
`

	updateConfig := `
resource "hpe_morpheus_environment" "test" {
  name        = "` + name + `"
  code        = "` + code + `"
  description = "Updated description"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_environment.test", "description", "Initial description"),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_environment.test", "description", "Updated description"),
				),
			},
		},
	})
}
