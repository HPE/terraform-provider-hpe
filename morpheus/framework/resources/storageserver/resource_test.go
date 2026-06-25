package storageserver_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storageserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageServerResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")
	}
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := storageserver.RenderStorageServerConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_storage_server.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_storage_server.example", "type", "nfs"),
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
				ResourceName:      "hpe_morpheus_storage_server.example",
			},
		},
	})
}

func TestAccMorpheusStorageServerResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")
	}
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := storageserver.RenderStorageServerConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := storageserver.RenderStorageServerConfig(t, map[string]string{
		"Name": updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_storage_server.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type", "nfs"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "type", "nfs"),
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
