// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package serviceplan_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name-provision.tf.tmpl Name "Example name" ProvisionTypeCode "arm"

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/serviceplan"
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

func TestAccMorpheusFindServicePlanById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	servicePlanName := acctest.RandomWithPrefix(t.Name())
	provisionTypeCode := "arm"

	servicePlanResourceConfig := `
resource "morpheus_service_plan" "test" {
  name = "` + servicePlanName + `"
  code = "standard"
  price_set_ids  = []
  provision_type = "` + provisionTypeCode + `"
}
`

	providerConfig := testhelpers.ProviderBlockMixed()

	dataSourceConfig, err := testhelpers.RenderExample(
		t, "example-id.tf.tmpl", "Id", "morpheus_service_plan.test.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"name",
			servicePlanName,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.3",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePlanResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindServicePlanByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	servicePlanName := acctest.RandomWithPrefix(t.Name())
	provisionTypeCode := "arm"

	servicePlanResourceConfig := `
resource "morpheus_service_plan" "test" {
  name = "` + servicePlanName + `"
  code = "standard"
  price_set_ids  = []
  provision_type = "` + provisionTypeCode + `"
}
`
	providerConfig := testhelpers.ProviderBlockMixed()

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name-provision.tf.tmpl",
		"Name", "morpheus_service_plan.test.name",
		"ProvisionTypeCode", "morpheus_service_plan.test.provision_type")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"name",
			servicePlanName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"provision_type_code",
			provisionTypeCode,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.3",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePlanResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindServicePlanNoPlanFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
		data "hpe_morpheus_service_plan" "test" {
			name = "____"
			provision_type_code = "arm"
		}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorNoServicePlanFound

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

func TestAccMorpheusFindServicePlanNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_service_plan" "test" {
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindServicePlanBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_service_plan" "test" {
				id = "1"
				name = "_____"
				provision_type_code = "______"
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorRunningPreApply

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

func TestAccMorpheusFindServicePlanByProvisionOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
			data "hpe_morpheus_service_plan" "test" {
				provision_type_code = "arm"
			}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := serviceplan.ErrorRunningPreApply

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
