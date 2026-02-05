// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package workflow_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsworkflow "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/workflow"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/workflow"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceWorkflowExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := workflow.RenderWorkflowOperationalConfig(t, map[string]string{
		"Name":       name,
		"Visibility": "public",
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dsworkflow.RenderWorkflowConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_workflow_operational.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_workflow.example",
			"id",
			"hpe_morpheus_workflow_operational.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_workflow.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
