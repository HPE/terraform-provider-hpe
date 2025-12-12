// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderHpeMorpheusFileTemplateConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            name,
		"Labels":          `["demo", "template", "terraform"]`,
		"FileName":        "tfcustom.cnf",
		"FilePath":        "/etc/my.cnf.d",
		"Phase":           "preProvision",
		"FileContent":     `"# Test MySQL Configuration\n[mysqld]\ninnodb_buffer_pool_size = 128M"`,
		"FileOwner":       "root",
		"SettingName":     "myCnf",
		"SettingCategory": "master",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_file_template_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Labels", defaults["Labels"],
		"FileName", defaults["FileName"],
		"FilePath", defaults["FilePath"],
		"Phase", defaults["Phase"],
		"FileContent", defaults["FileContent"],
		"FileOwner", defaults["FileOwner"],
		"SettingName", defaults["SettingName"],
		"SettingCategory", defaults["SettingCategory"],
	)
}

func TestAccMorpheusFileTemplateExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderHpeMorpheusFileTemplateConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"file_name",
			"tfcustom.cnf",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"file_path",
			"/etc/my.cnf.d",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"phase",
			"preProvision",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"file_owner",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"setting_name",
			"myCnf",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_file_template.tfexample_file_template",
			"setting_category",
			"master",
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
