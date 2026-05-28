package cluster_affinity_group_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_affinity_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusClusterAffinityGroupResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	clusterID := os.Getenv("TF_ACC_MORPHEUS_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("TF_ACC_MORPHEUS_CLUSTER_ID not set")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := cluster_affinity_group.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_cluster_affinity_group.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.example", "cluster_id", clusterID),
		resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.example", "name", name),
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

func TestAccMorpheusClusterAffinityGroupResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	clusterID := os.Getenv("TF_ACC_MORPHEUS_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("TF_ACC_MORPHEUS_CLUSTER_ID not set")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := cluster_affinity_group.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := cluster_affinity_group.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cluster_affinity_group.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
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
