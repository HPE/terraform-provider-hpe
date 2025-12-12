// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderTaskWriteAttributesConfig generates a Terraform configuration for testing
// the task_write_attributes resource. It accepts overrides to customize field values.
func RenderTaskWriteAttributesConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	code := strings.ToLower(name)

	defaults := map[string]string{
		"Label1":            "demo",
		"Label2":            "terraform",
		"Attributes":        `{"demo":"test"}`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"../task/morpheus_task_write_attributes_resource.tf.tmpl",
		"Name", name,
		"Code", code,
		"Label1", defaults["Label1"],
		"Label2", defaults["Label2"],
		"Attributes", defaults["Attributes"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func TestAccMorpheusTaskWriteAttributesExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderTaskWriteAttributesConfig(t, name, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"code",
			strings.ToLower(name),
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"labels.*",
			"demo",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"labels.*",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"attributes",
			`{"demo":"test"}`,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_write_attributes."+name,
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
