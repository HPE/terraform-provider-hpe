package setting_whitelabel_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/setting_whitelabel"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSettingWhitelabelResourceExampleOk(t *testing.T) {
	// We can't run this test in parallel as it's a singleton resource in Morpheus.
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	applianceName := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := setting_whitelabel.RenderSettingWhitelabelConfig(t, map[string]string{
		"ApplianceName": applianceName,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_setting_whitelabel.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_setting_whitelabel.example", "enabled", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_setting_whitelabel.example", "appliance_name", applianceName),
		resource.TestCheckResourceAttr("hpe_morpheus_setting_whitelabel.example", "primary_color", "#1a73e8"),
		resource.TestCheckResourceAttr("hpe_morpheus_setting_whitelabel.example", "secondary_color", "#ffffff"),
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

func TestAccMorpheusSettingWhitelabelResourceUpdateOk(t *testing.T) {
	// We can't run this test in parallel as it's a singleton resource in Morpheus.
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	applianceName := acctest.RandomWithPrefix(t.Name())
	updatedApplianceName := applianceName + "-updated"

	createConfig, err := setting_whitelabel.RenderSettingWhitelabelConfig(t, map[string]string{
		"ApplianceName": applianceName,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := setting_whitelabel.RenderSettingWhitelabelConfig(t, map[string]string{
		"ApplianceName":  updatedApplianceName,
		"PrimaryColor":   "#0f62fe",
		"SecondaryColor": "#161616",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_setting_whitelabel.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "appliance_name", applianceName),
		resource.TestCheckResourceAttr(resourceName, "primary_color", "#1a73e8"),
		resource.TestCheckResourceAttr(resourceName, "secondary_color", "#ffffff"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "appliance_name", updatedApplianceName),
		resource.TestCheckResourceAttr(resourceName, "primary_color", "#0f62fe"),
		resource.TestCheckResourceAttr(resourceName, "secondary_color", "#161616"),
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
