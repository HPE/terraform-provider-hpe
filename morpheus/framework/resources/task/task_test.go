package task_test

//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_conditional_workflow.tf example_conditional_workflow.tf.tmpl Name "Example Conditional Workflow Task" IfOperationalWorkflowId "4090" IfOperationalWorkflowName "Example If Workflow" ElseOperationalWorkflowId "4091" ElseOperationalWorkflowName "Example Else Workflow"
//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_generic_config.tf example_generic_config.tf.tmpl Name "Example Generic Task" OperationalWorkflowId "4090" OperationalWorkflowName "Example Workflow"

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusTaskExampleConditionalOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	ifOperationalWorkflowId := "140"
	elseOperationalWorkflowId := "102"
	ifOperationalWorkflowName := "A Failure"
	elseOperationalWorkflowName := "A Test"
	resourceConfig, err := testhelpers.RenderExample(t, "example_conditional_workflow.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", ifOperationalWorkflowId,
		"ElseOperationalWorkflowId", elseOperationalWorkflowId,
		"IfOperationalWorkflowName", ifOperationalWorkflowName,
		"ElseOperationalWorkflowName", elseOperationalWorkflowName,
	)
	if err != nil {
		t.Fatal(err)
	}

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
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			ifOperationalWorkflowId,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			elseOperationalWorkflowId,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

func TestAccMorpheusTaskConditionalWorkflowUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	ifOperationalWorkflowId := "140"
	elseOperationalWorkflowId := "102"
	ifOperationalWorkflowName := "A Failure"
	elseOperationalWorkflowName := "A Test"
	resourceConfig, err := testhelpers.RenderExample(t, "example_conditional_workflow.tf.tmpl",
		"Name", name,
		"IfOperationalWorkflowId", ifOperationalWorkflowId,
		"ElseOperationalWorkflowId", elseOperationalWorkflowId,
		"IfOperationalWorkflowName", ifOperationalWorkflowName,
		"ElseOperationalWorkflowName", elseOperationalWorkflowName,
	)
	if err != nil {
		t.Fatal(err)
	}

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
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.if_operational_workflow_id",
			ifOperationalWorkflowId,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			elseOperationalWorkflowId,
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
			},
			// update the task config
			{
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            checkFn,
				Config: providerConfig + `
					resource "hpe_morpheus_task" "example_task" {
						name = "` + name + `2"
						task_type_code = "conditionalWorkflow"
						config_conditional_workflow = {
							if_operational_workflow_id   = ` + ifOperationalWorkflowId + `
							if_operational_workflow_name = "` + ifOperationalWorkflowName + `"

							else_operational_workflow_id   = ` + elseOperationalWorkflowId + `
							else_operational_workflow_name = "` + elseOperationalWorkflowName + `"
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
				Config: providerConfig + `
					resource "hpe_morpheus_task" "example_task" {
						name = "` + name + `2"
						task_type_code = "nestedWorkflow"

						config = {
							operationalWorkflowId   = "` + ifOperationalWorkflowId + `"
							operationalWorkflowName = "` + ifOperationalWorkflowName + `"
						}

						execute_target = "local"
						retryable = false
					}`,
			},
		},
	})
}

func TestAccMorpheusTaskExampleGenericNestedOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	operationalWorkflowId := "102"
	operationalWorkflowName := "A Test"
	resourceConfig, err := testhelpers.RenderExample(t, "example_generic_config.tf.tmpl",
		"Name", name,
		"OperationalWorkflowId", operationalWorkflowId,
		"OperationalWorkflowName", operationalWorkflowName,
	)
	if err != nil {
		t.Fatal(err)
	}

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
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config.operationalWorkflowId",
			operationalWorkflowId,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config.operationalWorkflowName",
			operationalWorkflowName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
