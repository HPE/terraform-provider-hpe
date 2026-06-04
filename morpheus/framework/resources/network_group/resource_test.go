package network_group_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkGroupResourceExampleOk(t *testing.T) {
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

	resourceConfig, err := network_group.RenderNetworkGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_network_group.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_network_group.example", "description", "Network group for production workloads"),
		resource.TestCheckResourceAttr("hpe_morpheus_network_group.example", "visibility", "private"),
		resource.TestCheckResourceAttr("hpe_morpheus_network_group.example", "active", "true"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_network_group.example", "id"),
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
				ResourceName:      "hpe_morpheus_network_group.example",
			},
		},
	})
}

func TestAccMorpheusNetworkGroupResourceUpdateOk(t *testing.T) {
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

	createConfig, err := network_group.RenderNetworkGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_group" "example" {
  name        = "` + name + `"
  description = "Network group for updated production workloads"
  visibility  = "public"
  active      = false
}
`

	resourceName := "hpe_morpheus_network_group.example"

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Network group for production workloads"),
		resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Network group for updated production workloads"),
		resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
		resource.TestCheckResourceAttr(resourceName, "active", "false"),
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
