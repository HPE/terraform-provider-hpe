package option_type_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/option_type"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusOptionTypeResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	fieldName := acctest.RandomWithPrefix("environment")

	resourceConfig, err := option_type.RenderOptionTypeConfig(t, map[string]string{
		"Name":      name,
		"FieldName": fieldName,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "field_name", fieldName),
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "type", "select"),
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "field_label", "Environment"),
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "default_value", "development"),
		resource.TestCheckResourceAttr("hpe_morpheus_option_type.example", "required", "true"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_option_type.example", "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_option_type.example",
			},
		},
	})
}

func TestAccMorpheusOptionTypeResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	fieldName := acctest.RandomWithPrefix("environment")

	createConfig, err := option_type.RenderOptionTypeConfig(t, map[string]string{
		"Name":      name,
		"FieldName": fieldName,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_option_type" "example" {
  name          = "` + name + `"
  field_name    = "` + fieldName + `"
  type          = "select"
  field_label   = "Environment Updated"
  default_value = "staging"
  required      = true
}
`

	resourceName := "hpe_morpheus_option_type.example"

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "field_name", fieldName),
		resource.TestCheckResourceAttr(resourceName, "type", "select"),
		resource.TestCheckResourceAttr(resourceName, "field_label", "Environment"),
		resource.TestCheckResourceAttr(resourceName, "default_value", "development"),
		resource.TestCheckResourceAttr(resourceName, "required", "true"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "field_name", fieldName),
		resource.TestCheckResourceAttr(resourceName, "type", "select"),
		resource.TestCheckResourceAttr(resourceName, "field_label", "Environment Updated"),
		resource.TestCheckResourceAttr(resourceName, "default_value", "staging"),
		resource.TestCheckResourceAttr(resourceName, "required", "true"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
