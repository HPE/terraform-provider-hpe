package monitoring_check_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoring_check"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusMonitoringCheckResourceExampleOk(t *testing.T) {
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

	resourceConfig, err := monitoring_check.RenderMonitoringCheckConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_monitoring_check.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_monitoring_check.example", "check_type_id", "1"),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_monitoring_check.example",
			"description",
			"HTTP health check for production website",
		),
		resource.TestCheckResourceAttr("hpe_morpheus_monitoring_check.example", "check_interval", "60"),
		resource.TestCheckResourceAttr("hpe_morpheus_monitoring_check.example", "active", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_monitoring_check.example", "severity", "critical"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_monitoring_check.example", "id"),
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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_monitoring_check.example",
			},
		},
	})
}

func TestAccMorpheusMonitoringCheckResourceUpdateOk(t *testing.T) {
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

	createConfig, err := monitoring_check.RenderMonitoringCheckConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_monitoring_check" "example" {
  name           = "` + name + `"
  check_type_id  = 1
  description    = "HTTPS health check for staging website"
  check_interval = 120
  active         = true
  severity       = "warning"
}
`

	resourceName := "hpe_morpheus_monitoring_check.example"

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "check_type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "description", "HTTP health check for production website"),
		resource.TestCheckResourceAttr(resourceName, "check_interval", "60"),
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
		resource.TestCheckResourceAttr(resourceName, "severity", "critical"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "check_type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "description", "HTTPS health check for staging website"),
		resource.TestCheckResourceAttr(resourceName, "check_interval", "120"),
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
		resource.TestCheckResourceAttr(resourceName, "severity", "warning"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
