// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package power_schedule_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/power_schedule"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusPowerScheduleResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
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
	resourceName := "hpe_morpheus_power_schedule.example"

	resourceConfig, err := power_schedule.RenderPowerScheduleConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Power on during business hours"),
		resource.TestCheckResourceAttr(resourceName, "schedule_type", "power"),
		resource.TestCheckResourceAttr(resourceName, "schedule_timezone", "America/New_York"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "monday_on_time", "08:00"),
		resource.TestCheckResourceAttr(resourceName, "monday_off_time", "18:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_on_time", "08:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_off_time", "18:00"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusPowerScheduleResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
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
	resourceName := "hpe_morpheus_power_schedule.example"

	createConfig, err := power_schedule.RenderPowerScheduleConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_power_schedule" "example" {
  name              = "` + name + `"
  description       = "Updated business hours"
  schedule_type     = "power"
  schedule_timezone = "America/New_York"
  enabled           = false
  monday_on_time    = "09:00"
  monday_off_time   = "18:00"
  tuesday_on_time   = "08:00"
  tuesday_off_time  = "18:00"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Power on during business hours"),
		resource.TestCheckResourceAttr(resourceName, "schedule_type", "power"),
		resource.TestCheckResourceAttr(resourceName, "schedule_timezone", "America/New_York"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "monday_on_time", "08:00"),
		resource.TestCheckResourceAttr(resourceName, "monday_off_time", "18:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_on_time", "08:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_off_time", "18:00"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated business hours"),
		resource.TestCheckResourceAttr(resourceName, "schedule_type", "power"),
		resource.TestCheckResourceAttr(resourceName, "schedule_timezone", "America/New_York"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "monday_on_time", "09:00"),
		resource.TestCheckResourceAttr(resourceName, "monday_off_time", "18:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_on_time", "08:00"),
		resource.TestCheckResourceAttr(resourceName, "tuesday_off_time", "18:00"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
