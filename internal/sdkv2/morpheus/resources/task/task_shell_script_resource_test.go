// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func RenderTaskShellScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "local",
		"ScriptContent":     "  echo \"testing\"",
		"Sudo":              "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(t,
		"morpheus_task_shell_script_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"Sudo", defaults["Sudo"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

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

func TestAccMorpheusTaskShellScriptResourceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderTaskShellScriptConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"source_type",
			"local",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"sudo",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_local",
			"allow_custom_config",
			"true",
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
