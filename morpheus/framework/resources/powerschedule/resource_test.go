// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package powerschedule_test

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

func TestAccMorpheusPowerScheduleBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_power_schedule" "test" {
  name           = "` + name + `"
  monday_on_time = "07:00"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_power_schedule.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_power_schedule.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_power_schedule.test", "monday_on_time", "07:00"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_power_schedule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusPowerScheduleUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_power_schedule" "test" {
  name           = "` + name + `"
  monday_on_time = "07:00"
}
`

	updateConfig := `
resource "hpe_morpheus_power_schedule" "test" {
  name           = "` + name + `"
  monday_on_time = "09:00"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_power_schedule.test", "monday_on_time", "07:00"),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_power_schedule.test", "monday_on_time", "09:00"),
				),
			},
		},
	})
}
