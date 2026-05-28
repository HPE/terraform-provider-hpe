package network_pool_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_pool"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkPoolResourceExampleOk(t *testing.T) {
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
	resourceName := "hpe_morpheus_network_pool.example"

	resourceConfig, err := network_pool.RenderNetworkPoolConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.1.0"),
		resource.TestCheckResourceAttr(resourceName, "netmask", "255.255.255.0"),
		resource.TestCheckResourceAttr(resourceName, "gateway", "10.0.1.1"),
		resource.TestCheckResourceAttr(resourceName, "dns_domain", "example.com"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
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

func TestAccMorpheusNetworkPoolResourceUpdateOk(t *testing.T) {
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
	resourceName := "hpe_morpheus_network_pool.example"

	createConfig, err := network_pool.RenderNetworkPoolConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_pool" "example" {
  name           = "` + name + `"
  type_id        = 1
  subnet_address = "10.0.1.0"
  netmask        = "255.255.255.0"
  gateway        = "10.0.1.254"
  dns_domain     = "updated.example.com"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.1.0"),
		resource.TestCheckResourceAttr(resourceName, "netmask", "255.255.255.0"),
		resource.TestCheckResourceAttr(resourceName, "gateway", "10.0.1.1"),
		resource.TestCheckResourceAttr(resourceName, "dns_domain", "example.com"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "subnet_address", "10.0.1.0"),
		resource.TestCheckResourceAttr(resourceName, "netmask", "255.255.255.0"),
		resource.TestCheckResourceAttr(resourceName, "gateway", "10.0.1.254"),
		resource.TestCheckResourceAttr(resourceName, "dns_domain", "updated.example.com"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
