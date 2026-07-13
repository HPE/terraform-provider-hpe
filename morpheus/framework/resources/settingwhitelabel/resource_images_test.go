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

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusSettingWhitelabelResourceImagesOk is the regression test for
// MORPH-12625. header_logo, footer_logo, login_logo and favicon are local image
// file paths that the provider uploads via the multipart images endpoint.
// Previously they were sent (and silently dropped) via the JSON settings
// endpoint and read back as null, producing a "provider produced an unexpected
// new value ... but now null" error on apply.
//
// Step 1 applies a config that sets all four images (this alone reproduced the
// original failure). Step 2 is a plan-only step asserting the follow-up plan is
// empty (the configured paths round-trip consistently). Steps 3-4 remove the
// images and assert they are reset.
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
	footerLogo := writeTestPNG(t, dir, "footer.png")
	loginLogo := writeTestPNG(t, dir, "login.png")
	// The favicon field also accepts image/png per the Morpheus domain
	// contentType constraint, so a PNG fixture is valid here too.
	favicon := writeTestPNG(t, dir, "favicon.ico")

	resourceName := "hpe_morpheus_setting_whitelabel.images"

	withImages := providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_setting_whitelabel" "images" {
  enabled        = true
  appliance_name = %q
  header_logo    = %q
  footer_logo    = %q
  login_logo     = %q
  favicon        = %q
}
`, applianceName, headerLogo, footerLogo, loginLogo, favicon)

	withoutImages := providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_setting_whitelabel" "images" {
  enabled        = true
  appliance_name = %q
}
`, applianceName)

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "appliance_name", applianceName),
		// The configured local paths are preserved verbatim in state (not
		// overwritten with the server-generated storage URL).
		resource.TestCheckResourceAttr(resourceName, "header_logo", headerLogo),
		resource.TestCheckResourceAttr(resourceName, "footer_logo", footerLogo),
		resource.TestCheckResourceAttr(resourceName, "login_logo", loginLogo),
		resource.TestCheckResourceAttr(resourceName, "favicon", favicon),
	)

	resetChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "appliance_name", applianceName),
		resource.TestCheckNoResourceAttr(resourceName, "header_logo"),
		resource.TestCheckNoResourceAttr(resourceName, "footer_logo"),
		resource.TestCheckNoResourceAttr(resourceName, "login_logo"),
		resource.TestCheckNoResourceAttr(resourceName, "favicon"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// Apply with all four images set. Prior to the fix this apply failed
			// with an inconsistent-result error.
			{Config: withImages, Check: createChecks},
			// The configured paths must round-trip: a re-plan is empty.
			{Config: withImages, ExpectNonEmptyPlan: false, PlanOnly: true},
			// Removing the images resets them to null.
			{Config: withoutImages, Check: resetChecks},
			{Config: withoutImages, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
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
