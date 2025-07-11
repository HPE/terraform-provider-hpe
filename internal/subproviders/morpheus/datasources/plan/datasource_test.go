// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name-provision.tf.tmpl Name "Example name" ProvisionType "ARM"

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/plan"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
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

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindPlanById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	plan, err := testhelpers.CreatePlan(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeletePlan(t, plan.GetId())
	})

	planID := fmt.Sprintf("%d", plan.GetId())
	planName := plan.GetName()

	providerConfig := testhelpers.ProviderBlock()

	datasourceConfig, err := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", planID)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plan.test",
			"id",
			planID,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plan.test",
			"name",
			planName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datasourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindPlanByNameProvision(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	plan, err := testhelpers.CreatePlan(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testhelpers.DeletePlan(t, plan.GetId())
	})

	planID := fmt.Sprintf("%d", plan.GetId())
	planName := plan.GetName()
	planProvisonType := plan.GetProvisionType()
	planProvisionTypeName := planProvisonType.GetName()

	providerConfig := testhelpers.ProviderBlock()

	datasourceConfig, err := testhelpers.RenderExample(t, "example-name-provision.tf.tmpl",
		"Name", planName,
		"ProvisionType", planProvisionTypeName)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plan.test",
			"id",
			planID,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plan.test",
			"name",
			planName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plan.test",
			"provision_type",
			planProvisionTypeName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datasourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAcccMorpheusFindPlanNoPlanFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
		data "hpe_morpheus_plan" "test" {
			name = "____"
			provision_type = "ARM"
		}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := plan.ErrorNoPlanFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindPlanNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_plan" "test" {
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := plan.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindPlanBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_plan" "test" {
				id = "1"
				name = "_____"
				provision_type = "______"
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := plan.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
