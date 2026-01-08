// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package job_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/job"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/task"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccJobTaskDateAndTimeExampleOk(t *testing.T) {
	// this test has a dependency on Instance and will not be considered for now.
	return

	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependencyResourceConfig, currentDependency string
	var err error

	currentDependency, err = task.RenderTaskShellScriptConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dependencyResourceConfig += currentDependency

	resourceConfig, err := job.RenderJobTaskDateAndTimeConfig(t, map[string]string{
		"Name":   name,
		"TaskId": "hpe_morpheus_task_shell_script.tfexample_shell_local.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_task.example",
			"task_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"schedule_mode",
			"date_and_time",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"scheduled_date_and_time",
			"2022-12-30T06:00:00Z",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"context_type",
			"instance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_ids.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_ids.1",
			"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccJobTaskScheduleExampleOk(t *testing.T) {
	// this test has a dependency on Instance and will not be considered for now.
	return
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependencyResourceConfig, currentDependency string
	var err error

	currentDependency, err = task.RenderTaskShellScriptConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dependencyResourceConfig += currentDependency

	resourceConfig, err := job.RenderJobTaskScheduleConfig(t, map[string]string{
		"Name":   name,
		"TaskId": "hpe_morpheus_task_shell_script.tfexample_shell_local.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_task.example",
			"task_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"schedule_mode",
			"scheduled",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_task.example",
			"execution_schedule_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"context_type",
			"instance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_ids.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_ids.0",
			"91",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"custom_config",
			"{\"test\":\"new\"}",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccJobTaskManualExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependencyResourceConfig, currentDependency string
	var err error

	currentDependency, err = task.RenderTaskShellScriptConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dependencyResourceConfig += currentDependency

	resourceConfig, err := job.RenderJobTaskManualConfig(t, map[string]string{
		"Name":   name,
		"TaskId": "hpe_morpheus_task_shell_script.tfexample_shell_local.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_job_task.example",
			"task_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"schedule_mode",
			"manual",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"context_type",
			"instance-label",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_job_task.example",
			"instance_label",
			"demo",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
