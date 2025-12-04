// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusTaskEmailResourceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "task_email_resource.tf.tmpl",
		"Name", name,
		"Code", "tfexample_email",
		"Labels", "[\"demo\", \"terraform\"]",
		"EmailAddress", "<%=instance.createdByEmail%>",
		"Subject", "<%=instance.hostname%> provisioning complete",
		"Source", "local",
		"Content", "Your instance <%=instance.hostname%> was provisioned.",
		"SkipWrappedEmailTemplate", "false",
		"Retryable", "true",
		"RetryCount", "1",
		"RetryDelaySeconds", "10",
		"AllowCustomConfig", "true",
	)
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
