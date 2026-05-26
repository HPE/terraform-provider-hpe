package storage_bucket_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageBucketResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Skip("Skipping: requires external storage provider credentials (S3, etc.)")

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccStorageBucketConfig(rName, "s3", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_storage_bucket.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_storage_bucket.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_storage_bucket.test", "provider_type", "s3"),
				),
			},
			// ImportState (ignore sensitive fields)
			{
				ResourceName:            "hpe_morpheus_storage_bucket.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_key", "secret_key"},
			},
			// Update description
			{
				Config: providerConfig + testAccStorageBucketConfig(rName, "s3", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_storage_bucket.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccStorageBucketConfig(name, providerType, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`  description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_storage_bucket" "test" {
  name          = %q
  provider_type = %q
%s
}
`, name, providerType, desc)
}
