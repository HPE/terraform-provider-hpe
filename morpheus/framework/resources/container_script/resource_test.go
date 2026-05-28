package container_script_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/container_script"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusContainerScriptResourceExampleOk(t *testing.T) {
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

	resourceConfig, err := container_script.RenderContainerScriptConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_container_script.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "script_phase", "provision"),
		resource.TestCheckResourceAttr(resourceName, "script_type", "bash"),
		resource.TestCheckResourceAttr(resourceName, "script", "apt-get update && apt-get install -y curl wget"),
		resource.TestCheckResourceAttr(resourceName, "sudo_user", "true"),
		resource.TestCheckResourceAttr(resourceName, "fail_on_error", "true"),
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
		},
	})
}

func TestAccMorpheusContainerScriptResourceUpdateOk(t *testing.T) {
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

	createConfig, err := container_script.RenderContainerScriptConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_container_script" "example" {
  name          = "` + name + `"
  script_phase  = "teardown"
  script_type   = "bash"
  script        = "apt-get update && apt-get install -y curl wget"
  sudo_user     = true
  fail_on_error = false
}
`

	resourceName := "hpe_morpheus_container_script.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "script_phase", "provision"),
		resource.TestCheckResourceAttr(resourceName, "script_type", "bash"),
		resource.TestCheckResourceAttr(resourceName, "script", "apt-get update && apt-get install -y curl wget"),
		resource.TestCheckResourceAttr(resourceName, "sudo_user", "true"),
		resource.TestCheckResourceAttr(resourceName, "fail_on_error", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "script_phase", "teardown"),
		resource.TestCheckResourceAttr(resourceName, "script_type", "bash"),
		resource.TestCheckResourceAttr(resourceName, "script", "apt-get update && apt-get install -y curl wget"),
		resource.TestCheckResourceAttr(resourceName, "sudo_user", "true"),
		resource.TestCheckResourceAttr(resourceName, "fail_on_error", "false"),
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
