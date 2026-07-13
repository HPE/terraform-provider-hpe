// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingwhitelabel_test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusSettingWhitelabelResourceImagesOk is the regression test for
// MORPH-12625. The logo/favicon attributes are write-only local file paths
// (header_logo_wo, ...) uploaded via the multipart images endpoint, each paired
// with a *_wo_version trigger. Previously the logos were sent (and silently
// dropped) via the JSON settings endpoint and read back as null, producing a
// "provider produced an unexpected new value ... but now null" error on apply.
//
// Step 1 applies a config that sets all four images (this alone reproduced the
// original failure). Step 2 is a plan-only step asserting the follow-up plan is
// empty (the values round-trip consistently). Step 3 bumps a single version to
// re-upload one image in place. Because the *_wo values are write-only they are
// never present in state, so assertions use the *_wo_version attributes.
func TestAccMorpheusSettingWhitelabelResourceImagesOk(t *testing.T) {
	// Singleton resource in Morpheus: must not run in parallel.
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	applianceName := acctest.RandomWithPrefix(t.Name())

	dir := t.TempDir()
	headerLogo := writeTestPNG(t, dir, "header.png")
	headerLogoV2 := writeTestPNG(t, dir, "header-v2.png")
	footerLogo := writeTestPNG(t, dir, "footer.png")
	loginLogo := writeTestPNG(t, dir, "login.png")
	// The favicon field also accepts image/png per the Morpheus domain
	// contentType constraint, so a PNG fixture is valid here too.
	favicon := writeTestPNG(t, dir, "favicon.ico")

	resourceName := "hpe_morpheus_setting_whitelabel.images"

	create := imagesConfig(applianceName, headerLogo, 1, footerLogo, loginLogo, favicon)
	updated := imagesConfig(applianceName, headerLogoV2, 2, footerLogo, loginLogo, favicon)

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "appliance_name", applianceName),
		resource.TestCheckResourceAttr(resourceName, "header_logo_wo_version", "1"),
		resource.TestCheckResourceAttr(resourceName, "footer_logo_wo_version", "1"),
		resource.TestCheckResourceAttr(resourceName, "login_logo_wo_version", "1"),
		resource.TestCheckResourceAttr(resourceName, "favicon_wo_version", "1"),
		// Write-only values are never stored in state.
		resource.TestCheckNoResourceAttr(resourceName, "header_logo_wo"),
		resource.TestCheckNoResourceAttr(resourceName, "favicon_wo"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "appliance_name", applianceName),
		resource.TestCheckResourceAttr(resourceName, "header_logo_wo_version", "2"),
		resource.TestCheckResourceAttr(resourceName, "footer_logo_wo_version", "1"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// Apply with all four images set. Prior to the fix this apply failed
			// with an inconsistent-result error.
			{Config: providerConfig + create, Check: createChecks},
			// The values must round-trip: a re-plan is empty.
			{Config: providerConfig + create, ExpectNonEmptyPlan: false, PlanOnly: true},
			// Bumping header_logo_wo_version re-uploads that image in place.
			{Config: providerConfig + updated, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updated, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}

// imagesConfig renders a whitelabel resource that sets all four write-only image
// paths. headerVersion is parameterised so a test step can bump it to trigger a
// re-upload of the header logo.
func imagesConfig(applianceName, headerLogo string, headerVersion int, footerLogo, loginLogo, favicon string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_setting_whitelabel" "images" {
  enabled                = true
  appliance_name         = %q
  header_logo_wo         = %q
  header_logo_wo_version = %d
  footer_logo_wo         = %q
  footer_logo_wo_version = 1
  login_logo_wo          = %q
  login_logo_wo_version  = 1
  favicon_wo             = %q
  favicon_wo_version     = 1
}
`, applianceName, headerLogo, headerVersion, footerLogo, loginLogo, favicon)
}

// writeTestPNG writes a minimal valid 1x1 PNG image to dir/name and returns its
// absolute path. A real PNG is used (rather than a placeholder) so the Morpheus
// server accepts it against the whitelabel image contentType constraints.
func writeTestPNG(t *testing.T, dir, name string) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff})

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("creating test image %s: %v", path, err)
	}
	defer func() { _ = f.Close() }()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encoding test image %s: %v", path, err)
	}

	return path
}
