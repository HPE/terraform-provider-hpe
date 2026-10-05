// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	sdklegacy "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// EnvLicenseKey names a real licence key for the test appliance. The test
// installs it, uninstalls it behind Terraform's back and expects the next
// apply to put it back; destroy leaves the appliance without it.
const EnvLicenseKey = "TF_VAR_testacc_morpheus_license_key"

// legacyTestClient is an API client for out-of-band changes to the appliance.
func legacyTestClient(t *testing.T) *sdklegacy.Client {
	t.Helper()

	url, ok := testhelpers.LookupProviderEnv("", "url")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_url must be set")
	}

	token, _ := testhelpers.LookupProviderEnv("", "access_token")
	username, _ := testhelpers.LookupProviderEnv("", "username")
	password, _ := testhelpers.LookupProviderEnv("", "password")
	tenantSubdomain, _ := testhelpers.LookupProviderEnv("", "tenant_subdomain")
	_, insecure := testhelpers.LookupProviderEnv("", "insecure")

	return client.NewLegacyClient(
		context.Background(),
		url,
		username,
		password,
		tenantSubdomain,
		token,
		sdklegacy.WithInsecure(insecure),
	)
}

// TestAccMorpheusLicenseInstalledKeyOk: the resource installs the configured
// key, reads the licence's details back, reports no drift while the key is
// installed, and re-installs it after it was uninstalled out of band
// (MORPH-18310).
func TestAccMorpheusLicenseInstalledKeyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.License)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	key := strings.TrimSpace(os.Getenv(EnvLicenseKey))
	if key == "" {
		t.Skip(EnvLicenseKey + " not set; skipping test requiring a real licence key")
	}

	// Morpheus identifies the licence by the key's first eight characters; a
	// real key is a few thousand characters long, so anything shorter is a
	// misconfigured variable, not a key.
	if len(key) < 8 {
		t.Fatalf("%s must hold a licence key (at least 8 characters); got %d", EnvLicenseKey, len(key))
	}

	keyID := key[:8]
	providerConfig := testhelpers.ProviderBlock()

	// The key is passed through a variable so it never appears in the
	// rendered configuration or in test output.
	resourceConfig := `
variable "license_key" {
  type      = string
  sensitive = true
  default   = ` + "\"" + key + "\"" + `
}

resource "hpe_morpheus_license" "test" {
  key = var.license_key
}
`

	const rn = "hpe_morpheus_license.test"

	installedChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(rn, "id", keyID),
		resource.TestCheckResourceAttrSet(rn, "product_tier"),
		resource.TestCheckResourceAttrSet(rn, "account_name"),
		resource.TestCheckResourceAttrSet(rn, "start_date"),
	)

	uninstallOutOfBand := func() {
		t.Log("uninstalling the licence behind Terraform's back")

		resp, err := legacyTestClient(t).UninstallLicense(keyID, &sdklegacy.Request{})
		if err != nil {
			t.Fatalf("out-of-band uninstall failed: %v (%s)", err, resp)
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Install and read back.
			{
				Config: providerConfig + resourceConfig,
				Check:  installedChecks,
			},
			// No drift while the key is installed.
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Uninstalled behind Terraform's back: the refresh notices the key is
			// gone, the plan is an in-place update, and the apply re-installs it.
			{
				PreConfig: uninstallOutOfBand,
				Config:    providerConfig + resourceConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				Check: installedChecks,
			},
			// And no drift again.
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Import by key id: the key is not readable, so the import leaves it
			// to be set by the next apply; everything else is verified.
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateId:           keyID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
		},
	})
}
