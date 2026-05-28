package vdi_pool_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdi_pool"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusVdiPoolResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
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

	resourceConfig, err := vdi_pool.RenderVdiPoolConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_vdi_pool.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "description", "VDI pool for development team"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "max_pool_size", "10"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "min_idle", "2"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "initial_pool_size", "3"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "enabled", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "persistent_user", "true"),
		resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.example", "idle_timeout", "30"),
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
				ResourceName:      "hpe_morpheus_vdi_pool.example",
			},
		},
	})
}

func TestAccMorpheusVdiPoolResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
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

	createConfig, err := vdi_pool.RenderVdiPoolConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := vdi_pool.RenderVdiPoolConfig(t, map[string]string{
		"Name":        name,
		"Description": "Updated VDI pool for development team",
		"MaxPoolSize": "12",
		"MinIdle":     "1",
		"Enabled":     "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_vdi_pool.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "VDI pool for development team"),
		resource.TestCheckResourceAttr(resourceName, "max_pool_size", "10"),
		resource.TestCheckResourceAttr(resourceName, "min_idle", "2"),
		resource.TestCheckResourceAttr(resourceName, "initial_pool_size", "3"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "persistent_user", "true"),
		resource.TestCheckResourceAttr(resourceName, "idle_timeout", "30"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated VDI pool for development team"),
		resource.TestCheckResourceAttr(resourceName, "max_pool_size", "12"),
		resource.TestCheckResourceAttr(resourceName, "min_idle", "1"),
		resource.TestCheckResourceAttr(resourceName, "initial_pool_size", "3"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "persistent_user", "true"),
		resource.TestCheckResourceAttr(resourceName, "idle_timeout", "30"),
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
