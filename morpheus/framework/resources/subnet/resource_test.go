package subnet_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/subnet"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSubnetResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Subnet) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	// NOTE: intentionally NOT t.Parallel(). Both subnet resource tests provision
	// on the same Azure VNet (network 88), and Azure serializes write operations
	// per-VNet regardless of CIDR. Running them in parallel makes one create fail
	// with "Another operation on this or dependent resource is in progress".
	// Keeping them serial (no t.Parallel) avoids the collision.

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := subnet.RenderSubnetConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_subnet.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_subnet.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_subnet.example", "type_id", "8"),
		resource.TestCheckResourceAttr("hpe_morpheus_subnet.example", "visibility", "private"),
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
				ResourceName:      "hpe_morpheus_subnet.example",
			},
		},
	})
}

func TestAccMorpheusSubnetResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Subnet) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	// NOTE: intentionally NOT t.Parallel() — see TestAccMorpheusSubnetResourceExampleOk.
	// Both tests share Azure VNet 88, which serializes per-VNet operations.

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := subnet.RenderSubnetConfig(t, map[string]string{
		"Name": name,
		// Distinct from ExampleOk's 10.0.250.0/24. The tests run serially now,
		// but distinct CIDRs remain as defensive hygiene against a delete→create
		// CIDR-reuse race on the shared VNet.
		"SubnetCidr": "10.0.251.0/24",
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := subnet.RenderSubnetConfig(t, map[string]string{
		"Name":       name,
		"Visibility": "public",
		"SubnetCidr": "10.0.251.0/24",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_subnet.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "8"),
		resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "8"),
		resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
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
