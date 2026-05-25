package whitelabel_settings_test

import (
	"os"
	"testing"

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

func TestAccMorpheusWhitelabelSettingsBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping: whitelabel settings requires an appropriate license (whitelabel not approved by license)")
	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccWhitelabelSettingsConfig("TF Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.test", "appliance_name", "TF Test"),
				),
			},
		},
	})
}

func TestAccMorpheusWhitelabelSettingsUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping: whitelabel settings requires an appropriate license (whitelabel not approved by license)")
	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccWhitelabelSettingsConfig("TF Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.test", "appliance_name", "TF Test"),
				),
			},
			{
				Config: providerConfig + testAccWhitelabelSettingsConfig("TF Test Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_whitelabel_settings.test", "appliance_name", "TF Test Updated"),
				),
			},
		},
	})
}

func testAccWhitelabelSettingsConfig(applianceName string) string {
	return `
resource "hpe_morpheus_whitelabel_settings" "test" {
  enabled        = true
  appliance_name = "` + applianceName + `"
}
`
}
