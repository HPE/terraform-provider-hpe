// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func newProviderWithErrorGit() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactoriesGit = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithErrorGit,
}

func TestAccMorpheusTaskEmailGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "task_email_resource_git.tf.tmpl",
		"Name", name,
		"Code", "tfexample_email_git",
		"Labels", `["demo","terraform"]`,
		"EmailAddress", "<%=instance.createdByEmail%>",
		"Subject", "<%=instance.hostname%> provisioning complete",
		"Source", "repository",
		"ContentPath", "example.txt",
		"RepositoryId", "1",
		"VersionRef", "main",
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
			"hpe_morpheus_task_email.tfexample_email_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"code",
			"tfexample_email_git",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"email_address",
			"<%=instance.createdByEmail%>",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"subject",
			"<%=instance.hostname%> provisioning complete",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"source",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"content_path",
			"example.txt",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"repository_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"version_ref",
			"main",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"skip_wrapped_email_template",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_email.tfexample_email_git",
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesGit,
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
