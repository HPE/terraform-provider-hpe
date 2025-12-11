// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
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

func RenderMorpheusInstanceTypeConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) string {
	t.Helper()

	defaults := map[string]string{
		"Name":                  name,
		"Code":                  strings.ToLower(name),
		"Description":           "Terraform Example Instance Type",
		"Labels":                `["demo", "instance", "terraform"]`,
		"Category":              "web",
		"Visibility":            "private",
		"ImagePath":             "tfexample.png",
		"ImageName":             "tfexample.png",
		"Featured":              "false",
		"EnableDeployments":     "true",
		"EnableScaling":         "true",
		"EnableSettings":        "true",
		"EnvironmentPrefix":     "TFEXAMPLE_DEMO",
		"OptionTypeIds":         "[1910, 1912]",
		"EvarFirstName":         "first",
		"EvarFirstValue":        "first",
		"EvarFirstExport":       "true",
		"EvarSecondName":        "second",
		"EvarSecondMaskedValue": "second",
		"EvarSecondExport":      "false",
	}

	for k, v := range overrides {
		defaults[k] = v
	}

	args := make([]string, 0, len(defaults)*2)
	for k, v := range defaults {
		args = append(args, k, v)
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_instance_type_resource.tf.tmpl",
		args...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return resourceConfig
}

func TestAccMorpheusInstanceTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := RenderMorpheusInstanceTypeConfig(t, name, map[string]string{"Name": name})

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"code",
			strings.ToLower(name),
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"description",
			"Terraform Example Instance Type",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"category",
			"web",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"image_path",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"image_name",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"featured",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_deployments",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_scaling",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_settings",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"environment_prefix",
			"TFEXAMPLE_DEMO",
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
