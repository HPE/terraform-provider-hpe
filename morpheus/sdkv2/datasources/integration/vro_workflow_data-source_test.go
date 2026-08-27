// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsintegration "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusDataSourceVroWorkflowExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VRO)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var dependenciesConfig string

	datasourceConfig, err := dsintegration.RenderVroWorkflowConfig(t, map[string]string{
		"Name": strconv.Quote("Create an AD Computer Object"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_vro_workflow.example",
			"name",
			"Create an AD Computer Object",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}

// TestAccMorpheusDataSourceVroWorkflowNotFound verifies that looking up a
// non-existent vRO workflow returns a clear "not found" diagnostic instead of
// the previous nil value type-assertion error.
func TestAccMorpheusDataSourceVroWorkflowNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VRO)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	datasourceConfig, err := dsintegration.RenderVroWorkflowConfig(t, map[string]string{
		"Name": strconv.Quote("______nonexistent vRO workflow______"),
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + datasourceConfig,
				ExpectError: regexp.MustCompile(`no vRO workflow found named`),
			},
		},
	})
}
