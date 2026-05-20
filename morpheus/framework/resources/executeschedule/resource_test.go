// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package executeschedule_test

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

func TestAccMorpheusExecuteScheduleBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_execute_schedule" "test" {
  name = "` + name + `"
  cron = "0 0 * * *"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_execute_schedule.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_execute_schedule.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_execute_schedule.test", "cron", "0 0 * * *"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_execute_schedule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusExecuteScheduleUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_execute_schedule" "test" {
  name = "` + name + `"
  cron = "0 0 * * *"
}
`

	updateConfig := `
resource "hpe_morpheus_execute_schedule" "test" {
  name = "` + name + `"
  cron = "0 6 * * 1-5"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_execute_schedule.test", "cron", "0 0 * * *"),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_execute_schedule.test", "cron", "0 6 * * 1-5"),
				),
			},
		},
	})
}
