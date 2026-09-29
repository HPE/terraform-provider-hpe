// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusSettingProvisioningExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := setting.RenderSettingProvisioningConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"allow_zone_selection",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"allow_host_selection",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"require_environments",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"show_pricing",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"hide_datastore_stats",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"cross_tenant_naming_policies",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"cloudinit_username",
			"cloudinit",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

const provisioningResourceAddr = "hpe_morpheus_setting_provisioning.tf_example_provisioning_setting"

// TestAccMorpheusProvisioningSetting_updateEnvelope is the regression for
// MORPH-14742: the update (PUT) path used to dereference a provisioningSettings
// envelope that the PUT never returns, so create succeeded but any subsequent
// update failed with "Not found in response: ProvisioningSettings". This
// two-step create-then-update test asserts the second apply completes cleanly.
func TestAccMorpheusProvisioningSetting_updateEnvelope(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	createConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"ShowPricing": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Flip an unrelated boolean to force a genuine update on step 2.
	updateConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"ShowPricing": "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_pricing", "true",
				),
			},
			// The update path is the regression target: it must apply cleanly.
			{
				Config: providerConfig + updateConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_pricing", "false",
				),
			},
		},
	})
}

// TestAccMorpheusProvisioningSetting_showConsoleKeyboardSettings is the
// regression for MORPH-14736: show_console_keyboard_settings was present in the
// schema and read back, but was never included in the create/update payloads,
// so it never persisted. This asserts the value round-trips as "true" on the
// create path and again after an update.
func TestAccMorpheusProvisioningSetting_showConsoleKeyboardSettings(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	enabledConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"ShowConsoleKeyboardSettings": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Toggle an unrelated attribute to force an update while the keyboard
	// setting stays true, proving it persists on the update path too.
	enabledUpdatedConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"ShowConsoleKeyboardSettings": "true",
		"HideDatastoreStats":          "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Create path.
			{
				Config: providerConfig + enabledConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_console_keyboard_settings", "true",
				),
			},
			// Update path.
			{
				Config: providerConfig + enabledUpdatedConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_console_keyboard_settings", "true",
				),
			},
		},
	})
}

// TestAccMorpheusProvisioningSetting_validationEmptyCloudInitUsername is the
// regression for MORPH-16150(a): cloudinit_username had no ValidateFunc, so an
// empty string was silently accepted. It must now be rejected at plan time.
func TestAccMorpheusProvisioningSetting_validationEmptyCloudInitUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	emptyUsernameConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"CloudinitUsername": "",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + emptyUsernameConfig,
				ExpectError: regexp.MustCompile(`(?i)expected .*cloudinit_username.* to not be an empty string|empty string`),
				PlanOnly:    true,
			},
		},
	})
}

// TestAccMorpheusProvisioningSetting_passwordOmittedWhenUnset is the regression
// for MORPH-16150(b): unset password attributes were sent as "" which is
// destructive server-side (it clears a stored credential). A previously-set
// password must survive an unrelated update in which the password attribute is
// omitted from the config.
func TestAccMorpheusProvisioningSetting_passwordOmittedWhenUnset(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	// Step 1: set a windows_password explicitly.
	withPasswordConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"WindowsPassword": "InitialW!ndows1",
		"ShowPricing":     "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Step 2: omit the password entirely from config (rendered as an empty
	// string in the template) while changing an unrelated attribute. The unset
	// password must be omitted from the request, so the plan must be clean
	// (the sha256 DiffSuppressFunc treats "" as no change) and the apply must
	// not clear the stored credential.
	withoutPasswordConfig, err := setting.RenderSettingProvisioningConfig(t, map[string]string{
		"WindowsPassword": "",
		"ShowPricing":     "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + withPasswordConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_pricing", "true",
				),
			},
			{
				Config: providerConfig + withoutPasswordConfig,
				Check: resource.TestCheckResourceAttr(
					provisioningResourceAddr, "show_pricing", "false",
				),
			},
			// A follow-up plan must be empty: omitting the password must not
			// produce a perpetual diff or attempt to re-clear the credential.
			{
				Config:             providerConfig + withoutPasswordConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
