package task_test

//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_conditional_workflow.tf example_conditional_workflow.tf.tmpl Name "Example Conditional Workflow Task"
//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_generic_config.tf example_generic_config.tf.tmpl Name "Example Generic Task"

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
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

var zodiacTestParameters = systemoverride.SystemTestParameters{
	Name: "zodiac",
	Params: map[string]string{
		"IfOperationalWorkflowId":     "91",
		"IfOperationalWorkflowName":   "Hello World",
		"ElseOperationalWorkflowId":   "92",
		"ElseOperationalWorkflowName": "Hello World 2",
	},
}

var featureTestParameters = systemoverride.SystemTestParameters{
	Name: "feature",
	Params: map[string]string{
		"IfOperationalWorkflowId":     "4131",
		"IfOperationalWorkflowName":   "Hello World",
		"ElseOperationalWorkflowId":   "4432",
		"ElseOperationalWorkflowName": "Hello World 2",
	},
}

func TestAccMorpheusTaskExampleConditionalOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()
	testSystem := systemoverride.GetPreferred(t, "zodiac")
	params := systemoverride.GetParameters(testSystem, zodiacTestParameters, featureTestParameters)
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())
	resourceConfig, err := testhelpers.RenderExample(t, "example_conditional_workflow.tf.tmpl",
		append(params.ToSlice(), "Name", name)...,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(providerConfig + resourceConfig)
	fmt.Println(params)

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
			params["ifOperationalWorkflowId"],
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			params["elseOperationalWorkflowId"],
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
	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlock(testSystem)
	params := systemoverride.GetParameters(testSystem, zodiacTestParameters, featureTestParameters)

	name := acctest.RandomWithPrefix(t.Name())
	resourceConfig, err := testhelpers.RenderExample(t, "example_conditional_workflow.tf.tmpl",
		append(params.ToSlice(), "Name", name)...,
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
			params["ifOperationalWorkflowId"],
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config_conditional_workflow.else_operational_workflow_id",
			params["elseOperationalWorkflowId"],
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
							if_operational_workflow_id   =` + params["ifOperationalWorkflowId"] + `
							if_operational_workflow_name = "Test 1"

							else_operational_workflow_id   = ` + params["elseOperationalWorkflowId"] + `
							else_operational_workflow_name = "Test 2"
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

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlock(testSystem)
	params := systemoverride.GetParameters(testSystem,
		systemoverride.SystemTestParameters{
			Name: "zodiac",
			Params: map[string]string{
				"operationalWorkflowId":   "90",
				"operationalWorkflowName": "Hello World",
			},
		}, systemoverride.SystemTestParameters{
			Name: "feature",
			Params: map[string]string{
				"operationalWorkflowId":   "3143",
				"operationalWorkflowName": "Hello World 2",
			},
		},
	)

	name := acctest.RandomWithPrefix(t.Name())
	resourceConfig, err := testhelpers.RenderExample(t, "example_generic_config.tf.tmpl",
		append(params.ToSlice(), "Name", name)...,
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
			params["operationalWorkflowId"],
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task.example_task",
			"config.operationalWorkflowName",
			params["operationalWorkflowName"],
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
