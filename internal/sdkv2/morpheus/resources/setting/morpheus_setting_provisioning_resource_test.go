// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func RenderMorpheusSettingProvisioningConfig(t *testing.T, overrides map[string]string) string {
	t.Helper()

	defaults := map[string]string{
		"allow_zone_selection":         "false",
		"allow_host_selection":         "false",
		"require_environments":         "false",
		"show_pricing":                 "true",
		"hide_datastore_stats":         "true",
		"cross_tenant_naming_policies": "false",
		"cloudinit_username":           "cloudinit",
		"cloudinit_password":           "Pa55w0rd!",
		"windows_password":             "Pa55w0rd!",
		"pxe_root_password":            "Pa55w0rd!",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_setting_provisioning" "tf_example_provisioning_setting" {
  allow_zone_selection         = %s
  allow_host_selection         = %s
  require_environments         = %s
  show_pricing                 = %s
  hide_datastore_stats         = %s
  cross_tenant_naming_policies = %s
  cloudinit_username           = "%s"
  cloudinit_password           = "%s"
  windows_password             = "%s"
  pxe_root_password            = "%s"
}
`,
		defaults["allow_zone_selection"],
		defaults["allow_host_selection"],
		defaults["require_environments"],
		defaults["show_pricing"],
		defaults["hide_datastore_stats"],
		defaults["cross_tenant_naming_policies"],
		defaults["cloudinit_username"],
		defaults["cloudinit_password"],
		defaults["windows_password"],
		defaults["pxe_root_password"],
	)
}

func TestAccMorpheusSettingProvisioningExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := RenderMorpheusSettingProvisioningConfig(t, nil)

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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
