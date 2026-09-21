// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// MORPH-16489: the nestedWorkflow / conditionalWorkflow task APIs validate the
// OWNERSHIP (not merely the visibility) of the referenced operational
// workflow. Hardcoded master-owned workflows (ids 797/796, names
// "full-operational-workflow-*") therefore fail when the acceptance tests run
// as a sub-tenant. Instead of pinning those, each test now declares its own
// operational workflow(s) per run and wires both the id and the name from the
// created resource — the API requires a non-blank name, then overwrites it from
// the resolved workflow. Workflows are created with visibility = "private":
// sub-tenants cannot own public workflows, and public workflows cause perpetual
// plan churn (MORPH-16419).

// operationalWorkflowHCL renders a private operational workflow resource with a
// per-run unique name. The Terraform resource label and the workflow name are
// both derived from label so multiple workflows in one config do not collide.
func operationalWorkflowHCL(label, name string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_workflow_operational" %[1]q {
  name                = %[2]q
  description         = "MORPH-16489 per-run operational workflow"
  labels              = ["terraform", "acctest"]
  platform            = "all"
  visibility          = "private"
  allow_custom_config = true
}
`, label, name)
}

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// taskAccFactories muxes the framework provider (hpe_morpheus_task) with the
// SDKv2 provider so tests can also create hpe_morpheus_workflow_operational
// resources (an SDKv2 resource) as per-run fixtures. See MORPH-16489.
func taskAccFactories(t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()

	return testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider())
}

func TestAccMorpheusTaskResourceExampleConditionalOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// MORPH-16489: create the two referenced operational workflows per run and
	// wire both id and name from the created resources, instead of hardcoding
	// master-owned workflows.
	ifWF := operationalWorkflowHCL("if_wf", name+"-if")
	elseWF := operationalWorkflowHCL("else_wf", name+"-else")
	resourceConfig, err := testhelpers.RenderExample(
		t, "example_conditional_workflow.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", "hpe_morpheus_workflow_operational.if_wf.id",
		"ElseOperationalWorkflowId", "hpe_morpheus_workflow_operational.else_wf.id",
		"IfOperationalWorkflowName", "hpe_morpheus_workflow_operational.if_wf.name",
		"ElseOperationalWorkflowName", "hpe_morpheus_workflow_operational.else_wf.name",
	)
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = ifWF + elseWF + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"task_type_code",
			"conditionalWorkflow",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			"hpe_morpheus_workflow_operational.if_wf", "id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			"hpe_morpheus_workflow_operational.else_wf", "id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: taskAccFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_task.example_task",
				ImportStateVerifyIgnore: []string{
					// trailing whitespace causes a mismatch during import
					// StringSemanticEquals method of TrimmedString type is not
					// used on import
					"config_conditional_workflow.conditional_script",
				},
				Check: checkFn,
			},
		},
	})
}

func TestAccMorpheusTaskResourceConditionalWorkflowUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// MORPH-16489: create the referenced operational workflows per run and wire
	// both id and name from them, rather than hardcoding master-owned workflows.
	ifWF := operationalWorkflowHCL("if_wf", name+"-if")
	elseWF := operationalWorkflowHCL("else_wf", name+"-else")
	workflowConfig := ifWF + elseWF
	resourceConfig, err := testhelpers.RenderExample(
		t, "example_conditional_workflow.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", "hpe_morpheus_workflow_operational.if_wf.id",
		"ElseOperationalWorkflowId", "hpe_morpheus_workflow_operational.else_wf.id",
		"IfOperationalWorkflowName", "hpe_morpheus_workflow_operational.if_wf.name",
		"ElseOperationalWorkflowName", "hpe_morpheus_workflow_operational.else_wf.name",
	)
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = workflowConfig + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name+"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"task_type_code",
			"conditionalWorkflow",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			"hpe_morpheus_workflow_operational.if_wf", "id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			"hpe_morpheus_workflow_operational.else_wf", "id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"allow_custom_config",
			"false",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	// this check verifies that the resource is going to be updated in place
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_task.example_task",
				plancheck.ResourceActionUpdate,
			),
		},
	}

	// verify that the resource will be destroyed before creation
	checkReplaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_task.example_task",
				plancheck.ResourceActionDestroyBeforeCreate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: taskAccFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
			// update the task config
			{
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            checkFn,
				Config: providerConfig + workflowConfig + `
					resource "hpe_morpheus_task" "example_task" {
						name = "` + name + `2"
						task_type_code = "conditionalWorkflow"
						config_conditional_workflow = {
							if_operational_workflow_id   = hpe_morpheus_workflow_operational.if_wf.id
							if_operational_workflow_name = hpe_morpheus_workflow_operational.if_wf.name

							else_operational_workflow_id   = hpe_morpheus_workflow_operational.else_wf.id
							else_operational_workflow_name = hpe_morpheus_workflow_operational.else_wf.name
						}

						execute_target = "local"
						retryable = false
						allow_custom_config = false
					}`,
			},
			// check if changing task_type_code requires replace
			{
				ConfigPlanChecks: checkReplaceUpdate,
				PlanOnly:         false,
				Config: providerConfig + workflowConfig + `
					resource "hpe_morpheus_task" "example_task" {
						name = "` + name + `2"
						task_type_code = "nestedWorkflow"

						config = {
							operationalWorkflowId   = hpe_morpheus_workflow_operational.if_wf.id
							operationalWorkflowName = hpe_morpheus_workflow_operational.if_wf.name
						}

						execute_target = "local"
						retryable = false
					}`,
			},
		},
	})
}

func TestAccMorpheusTaskResourceExampleGenericNestedOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// MORPH-16489: create the referenced operational workflow per run and wire
	// both id and name from it, rather than hardcoding a master-owned workflow.
	wf := operationalWorkflowHCL("wf", name+"-wf")
	resourceConfig, err := testhelpers.RenderExample(
		t, "example_generic_config.tf.tmpl",
		"Name", name,
		// config is a Dynamic (string) map, so the numeric id is coerced to a
		// string to match how the value round-trips.
		"OperationalWorkflowId", "tostring(hpe_morpheus_workflow_operational.wf.id)",
		"OperationalWorkflowName", "hpe_morpheus_workflow_operational.wf.name",
	)
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = wf + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"task_type_code",
			"nestedWorkflow",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config.operationalWorkflowId",
			"hpe_morpheus_workflow_operational.wf", "id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config.operationalWorkflowName",
			"hpe_morpheus_workflow_operational.wf", "name",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: taskAccFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_task.example_task",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusTaskResourceExampleConditionalNullElseOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// MORPH-16489: create the referenced "if" operational workflow per run.
	ifWF := operationalWorkflowHCL("if_wf", name+"-if")
	resourceConfig, err := testhelpers.RenderExample(
		t, "example_conditional_workflow_null_else.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", "hpe_morpheus_workflow_operational.if_wf.id",
		"IfOperationalWorkflowName", "hpe_morpheus_workflow_operational.if_wf.name",
	)
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = ifWF + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"task_type_code",
			"conditionalWorkflow",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			"hpe_morpheus_workflow_operational.if_wf", "id",
		),
		// Else id should be null in state
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: taskAccFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_task.example_task",
				ImportStateVerifyIgnore: []string{
					"config_conditional_workflow.conditional_script",
				},
				Check: checkFn,
			},
		},
	})
}

func TestAccMorpheusTaskResourceExampleConditionalNullElseUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	// MORPH-16489: create the referenced operational workflows per run. The
	// "else" workflow is only added in the update step, matching the transition
	// from a null else to a populated one.
	ifWF := operationalWorkflowHCL("if_wf", name+"-if")
	elseWF := operationalWorkflowHCL("else_wf", name+"-else")
	resourceConfig, err := testhelpers.RenderExample(
		t, "example_conditional_workflow_null_else.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", "hpe_morpheus_workflow_operational.if_wf.id",
		"IfOperationalWorkflowName", "hpe_morpheus_workflow_operational.if_wf.name",
	)
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = ifWF + resourceConfig

	nullElseChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name,
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
		),
	)

	nonNullElseChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			"hpe_morpheus_workflow_operational.if_wf", "id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			"hpe_morpheus_workflow_operational.else_wf", "id",
		),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				"hpe_morpheus_task.example_task",
				plancheck.ResourceActionUpdate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: taskAccFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              nullElseChecks,
			},
			// update to add a non-null else workflow
			{
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            nonNullElseChecks,
				Config: providerConfig + ifWF + elseWF + `
					resource "hpe_morpheus_task" "example_task" {
						name = "` + name + `"
						task_type_code = "conditionalWorkflow"
						config_conditional_workflow = {
							if_operational_workflow_id   = hpe_morpheus_workflow_operational.if_wf.id
							if_operational_workflow_name = hpe_morpheus_workflow_operational.if_wf.name

							else_operational_workflow_id   = hpe_morpheus_workflow_operational.else_wf.id
							else_operational_workflow_name = hpe_morpheus_workflow_operational.else_wf.name
						}

						execute_target = "local"
						retryable = false
						allow_custom_config = true
					}`,
			},
		},
	})
}
