// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusTaskChefBootstrapResourceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "task_chef_bootstrap_resource.tf.tmpl",
		"Name", name,
		"Code", "terraform_example_chef",
		"Labels", "[\"demo\", \"terraform\"]",
		"ChefServerId", "9",
		"Environment", "dev",
		"RunList", "role[web]",
		"DataBagKey", "test123",
		"DataBagKeyPath", "/etc/chef/databag_secret",
		"NodeName", "demonode",
		"NodeAttributes", "{\n  \"test\":\"demo\"\n}",
		"Retryable", "true",
		"RetryCount", "1",
		"RetryDelaySeconds", "10",
		"AllowCustomConfig", "true",
		"Visibility", "public",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"code",
			"terraform_example_chef",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"chef_server_id",
			"9",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"environment",
			"dev",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"run_list",
			"role[web]",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"data_bag_key",
			"test123",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"data_bag_key_path",
			"/etc/chef/databag_secret",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"node_name",
			"demonode",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"node_attributes",
			"{\n  \"test\":\"demo\"\n}",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"allow_custom_config",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"visibility",
			"public",
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
