package whitelabel_settings_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusWhitelabelSettingsBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping: whitelabel settings requires an appropriate license (whitelabel not approved by license)")
	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping: whitelabel settings requires an appropriate license (whitelabel not approved by license)")
	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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
