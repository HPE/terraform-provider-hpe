// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package job_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/job"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/workflow"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusJobWorkflowDateAndTimeExampleOk(t *testing.T) {
	t.Skip("resource hpe_morpheus_job_workflow has external dependencies not yet implemented in test (instance)")
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := workflow.RenderWorkflowOperationalConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := job.RenderJobWorkflowDateAndTimeConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.0",
			"aws",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.1",
			"demo",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_workflow.example",
			"workflow_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"schedule_mode",
			"date_and_time",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"scheduled_date_and_time",
			"2022-12-30T06:00:00Z",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"context_type",
			"instance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_ids.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_ids.1",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"custom_options.demo",
			"testing",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, nil, sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusJobWorkflowScheduleExampleOk(t *testing.T) {
	t.Skip("resource hpe_morpheus_job_workflow has external dependencies not yet implemented in test (instance)")
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := workflow.RenderWorkflowOperationalConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := job.RenderJobWorkflowScheduleConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.0",
			"aws",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.1",
			"demo",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_workflow.example",
			"workflow_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"schedule_mode",
			"scheduled",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"execution_schedule_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"context_type",
			"instance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_ids.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_ids.0",
			"91",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"custom_options.demo",
			"testing",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, nil, sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusJobWorkflowManualExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := workflow.RenderWorkflowOperationalConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := job.RenderJobWorkflowManualConfig(t, map[string]string{
		"Name":       name,
		"WorkflowId": "hpe_morpheus_workflow_operational.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.0",
			"aws",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"labels.1",
			"demo",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_workflow.example",
			"workflow_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"schedule_mode",
			"manual",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"context_type",
			"instance-label",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"instance_label",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_workflow.example",
			"custom_options.demo",
			"testing",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, nil, sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
