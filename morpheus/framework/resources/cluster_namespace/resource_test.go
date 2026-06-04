package cluster_namespace_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_namespace"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusClusterNamespaceResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Kubernetes) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	clusterID := "571"

	providerConfig := testhelpers.ProviderBlock()
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	// to get around the 63-character name limitation
	name = name[:63]

	resourceConfig, err := cluster_namespace.RenderClusterNamespaceConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("hpe_morpheus_cluster_namespace.example", "id"),
		resource.TestCheckResourceAttr("hpe_morpheus_cluster_namespace.example", "cluster_id", clusterID),
		resource.TestCheckResourceAttr("hpe_morpheus_cluster_namespace.example", "name", name),
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
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"active"},
				ResourceName:            "hpe_morpheus_cluster_namespace.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_cluster_namespace.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["cluster_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

func TestAccMorpheusClusterNamespaceResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Kubernetes) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	clusterID := "571"

	providerConfig := testhelpers.ProviderBlock()
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	// to get around the 63-character name limitation
	name = name[:62]

	// 63 character name limit
	updatedName := name + "u"

	createConfig, err := cluster_namespace.RenderClusterNamespaceConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := cluster_namespace.RenderClusterNamespaceConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cluster_namespace.example"
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
