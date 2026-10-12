// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// EnvLicenseKey names a real licence key for the test appliance.
const EnvLicenseKey = "TF_VAR_testacc_morpheus_license_key"

const licenseResourceName = "hpe_morpheus_license.test"

// TestAccMorpheusLicenseInstalledKeyOk: the framework resource installs the
// configured key, reads its details back, reports no drift while installed,
// recreates it after an out-of-band uninstall, rejects a bogus key, and
// imports by key id (MORPH-18310, MORPH-18584, MORPH-18805).
//
// Not exercised: replacing the key with a second valid one (the update
// action, the new key id and digest, and the empty plan that follows). The
// test environment has a single licence key, and a licence key cannot be
// fabricated. The bogus-key step covers the change-detection half of that
// path -- the plan is an update and the new key is posted -- up to the
// appliance's refusal.
func TestAccMorpheusLicenseInstalledKeyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.License)

	testhelpers.SkipUnlessAcc(t)

	key := strings.TrimSpace(os.Getenv(EnvLicenseKey))
	if key == "" {
		t.Skip(EnvLicenseKey + " not set; skipping test requiring a real licence key")
	}
	if len(key) < 8 {
		t.Fatalf("%s must hold a licence key (at least 8 characters); got %d", EnvLicenseKey, len(key))
	}

	keyID := key[:8]
	providerConfig := testhelpers.ProviderBlock()
	factories := testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil)

	// The key is passed through a variable so it never appears in rendered
	// configuration or test output.
	baseConfig := providerConfig + `
variable "license_key" {
  type      = string
  sensitive = true
  default   = "` + key + `"
}

resource "hpe_morpheus_license" "test" {
  key = var.license_key
}
`

	// A whitespace variant of the same key: trimming makes it the same licence.
	whitespaceConfig := providerConfig + `
variable "license_key" {
  type      = string
  sensitive = true
  default   = "` + key + `"
}

resource "hpe_morpheus_license" "test" {
  key = "  ${var.license_key}\n"
}
`

	bogusConfig := providerConfig + `
resource "hpe_morpheus_license" "test" {
  key = "NOTAREALKEY-0123456789"
}
`

	installedChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(licenseResourceName, "id", keyID),
		resource.TestCheckResourceAttrSet(licenseResourceName, "product_tier"),
		resource.TestCheckResourceAttrSet(licenseResourceName, "account_name"),
		resource.TestCheckResourceAttrSet(licenseResourceName, "start_date"),
		resource.TestCheckNoResourceAttr(licenseResourceName, "key"),
	)

	rec := testhelpers.RecordRequests(t)

	uninstallOutOfBand := func() {
		rec.Reset()

		ctx := context.Background()
		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("building out-of-band client: %v", err)
		}
		if _, _, err := client.LicenseAPI.UninstallLicense(ctx).KeyId(keyID).Execute(); err != nil {
			t.Fatalf("out-of-band uninstall failed: %v", err)
		}
	}

	// A step whose configuration drops the resource from state WITHOUT
	// uninstalling (the import step then re-adds it by key id).
	forget := providerConfig + `
removed {
  from = hpe_morpheus_license.test

  lifecycle {
    destroy = false
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			// 1. Create: the key is POSTed.
			{
				PreConfig: func() { rec.Reset() },
				Config:    baseConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					installedChecks,
					rec.ExpectSent("/api/license", "license", true),
				),
			},
			// 2. No drift: the digest matches the API hash.
			{
				Config:             baseConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 3. Whitespace variant trims to the same key: empty plan.
			{
				Config:             whitespaceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 4. Uninstalled out of band: refresh removes it, plan is a CREATE.
			{
				PreConfig: uninstallOutOfBand,
				Config:    baseConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(licenseResourceName, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					installedChecks,
					rec.ExpectSent("/api/license", "license", true),
				),
			},
			// 5. No drift again.
			{
				Config:             baseConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 6. A bogus key: the plan is an update (the digest differs), the
			// key is posted, and the appliance refuses it with HTTP 400 and
			// its "Error decoding license" message, which the update reports.
			{
				PreConfig: func() { rec.Reset() },
				Config:    bogusConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(licenseResourceName, plancheck.ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`(?s)update license.*Error decoding license`),
			},
			// The failed apply must not have changed state: the real key still
			// plans empty.
			{
				Config:             baseConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 7. Forget without uninstalling, then import by key id. The forget
			// step removes the resource from the test state, so there is no
			// prior state for ImportStateVerify to compare with; the imported
			// attributes are checked directly instead.
			{
				Config: forget,
			},
			{
				ResourceName:       licenseResourceName,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateId:      keyID,
				Config:             baseConfig,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("imported %d instances, want 1", len(states))
					}

					attrs := states[0].Attributes
					if got := attrs["id"]; got != keyID {
						return fmt.Errorf("imported id = %q, want %q", got, keyID)
					}

					for _, k := range []string{"product_tier", "account_name", "start_date"} {
						if attrs[k] == "" {
							return fmt.Errorf("imported %s is empty, want the installed licence's value", k)
						}
					}

					if v, ok := attrs["key"]; ok && v != "" {
						return fmt.Errorf("imported key = %q, want absent: the key is write-only", v)
					}

					return nil
				},
			},
			// The digest was recorded from the API hash on import: empty plan.
			{
				Config:             baseConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
