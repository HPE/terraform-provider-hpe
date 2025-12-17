// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/task"
)

func TestAccMorpheusTaskEmailUrlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskEmailUrlConfig(t, name, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"code",
			"tfexample_email_url",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"email_address",
			"<%=instance.createdByEmail%>",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"subject",
			"<%=instance.hostname%> provisioning complete",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"source",
			"url",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"content_url",
			"https://example.com/example.txt",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"skip_wrapped_email_template",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_url",
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
