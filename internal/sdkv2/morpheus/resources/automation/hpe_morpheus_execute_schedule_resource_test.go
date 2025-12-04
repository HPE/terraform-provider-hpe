// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation_test

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

func TestAccMorpheusHpeMorpheusExecuteScheduleExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "hpe_morpheus_execute_schedule_resource.tf.tmpl",
		"Name", name,
		"Description", "This schedule runs daily at 7 AM Mountain Time",
		"Enabled", "false",
		"TimeZone", "America/Denver",
		"Schedule", "7 0 * * *",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_execute_schedule.tf_example_execute_schedule",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_execute_schedule.tf_example_execute_schedule",
			"description",
			"This schedule runs daily at 7 AM Mountain Time",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_execute_schedule.tf_example_execute_schedule",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_execute_schedule.tf_example_execute_schedule",
			"time_zone",
			"America/Denver",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_execute_schedule.tf_example_execute_schedule",
			"schedule",
			"7 0 * * *",
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
