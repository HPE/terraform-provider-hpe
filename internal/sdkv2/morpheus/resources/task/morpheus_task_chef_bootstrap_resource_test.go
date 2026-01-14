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
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/task"
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

func TestAccMorpheusTaskChefBootstrapExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderChefBootstrapConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
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
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"labels.#",
			"2",
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
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_task_chef_bootstrap.cheftask",
			"data_bag_key",
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
