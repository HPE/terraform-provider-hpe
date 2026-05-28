package whitelabel_settings_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/whitelabel_settings"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusWhitelabelSettingsResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	applianceName := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := whitelabel_settings.RenderWhitelabelSettingsConfig(t, map[string]string{
		"ApplianceName": applianceName,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_whitelabel_settings.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.example", "enabled", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.example", "appliance_name", applianceName),
		resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.example", "primary_color", "#1a73e8"),
		resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.example", "secondary_color", "#ffffff"),
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

func TestAccMorpheusWhitelabelSettingsResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	applianceName := acctest.RandomWithPrefix(t.Name())
	updatedApplianceName := applianceName + "-updated"

	createConfig, err := whitelabel_settings.RenderWhitelabelSettingsConfig(t, map[string]string{
		"ApplianceName": applianceName,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := whitelabel_settings.RenderWhitelabelSettingsConfig(t, map[string]string{
		"ApplianceName":  updatedApplianceName,
		"PrimaryColor":   "#0f62fe",
		"SecondaryColor": "#161616",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_whitelabel_settings.example"
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
