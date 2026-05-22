package cluster_affinity_group_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(
			provider.New("test", morpheus.New())(),
		)()
	},
}

func TestAccMorpheusClusterAffinityGroupResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	clusterID := os.Getenv("TF_ACC_MORPHEUS_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("TF_ACC_MORPHEUS_CLUSTER_ID not set, skipping")
	}

	rName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccClusterAffinityGroupConfig(clusterID, rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_cluster_affinity_group.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.test", "cluster_id", clusterID),
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.test", "name", rName),
				),
			},
			// ImportState with composite ID
			{
				ResourceName:      "hpe_morpheus_cluster_affinity_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources["hpe_morpheus_cluster_affinity_group.test"]

					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["id"]), nil
				},
			},
			// Update description
			{
				Config: testAccClusterAffinityGroupConfig(clusterID, rName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccClusterAffinityGroupConfig(clusterID, name, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`  description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_cluster_affinity_group" "test" {
  cluster_id = %q
  name       = %q
%s
}
`, clusterID, name, desc)
}
