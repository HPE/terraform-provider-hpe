// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusTaskEmailExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskEmailConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"code",
			"tfexample_email",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"email_address",
			"<%=instance.createdByEmail%>",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"subject",
			"<%=instance.hostname%> provisioning complete",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"source",
			"local",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"content",
			"Your instance <%=instance.hostname%> was provisioned.",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"skip_wrapped_email_template",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email",
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
