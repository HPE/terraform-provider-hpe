package cluster_datastore_test

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

func TestAccClusterDatastoreResource_basic(t *testing.T) {
	clusterID := os.Getenv("TF_ACC_MORPHEUS_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("TF_ACC_MORPHEUS_CLUSTER_ID not set, skipping")
	}

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccClusterDatastoreConfig(clusterID, rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_cluster_datastore.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_datastore.test", "cluster_id", clusterID),
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_datastore.test", "active", "true"),
				),
			},
			// ImportState with composite ID
			{
				ResourceName:      "hpe_morpheus_cluster_datastore.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources["hpe_morpheus_cluster_datastore.test"]

					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["id"]), nil
				},
			},
			// Update active
			{
				Config: testAccClusterDatastoreConfig(clusterID, rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_cluster_datastore.test", "active", "false"),
				),
			},
		},
	})
}

func testAccClusterDatastoreConfig(clusterID, name string, active bool) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_cluster_datastore" "test" {
  cluster_id = %q
  name       = %q
  active     = %t
}
`, clusterID, name, active)
}
