package image_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusImageLocalFile(t *testing.T) {
	defer testhelpers.RecordResult(t)
	capabilities.MustHaveOrSkip(t, capabilities.All)
	if testing.Short() {
		t.Skip("live image upload")
	}
	name := acctest.RandomWithPrefix(t.Name())
	first := filepath.Join(t.TempDir(), "disk.qcow2")
	second := filepath.Join(t.TempDir(), "disk.qcow2")
	for filename, contents := range map[string]string{first: "first image", second: "other image"} {
		if err := os.WriteFile(filename, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	provider := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_os_type" "local" { name = "linux" }
data "hpe_morpheus_storage_bucket" "local" { name = %q }
`, testhelpers.ImageStorageBucket("Local Storage"))
	render := func(filename string) string {
		return provider + fmt.Sprintf(`
resource "hpe_morpheus_image" "local" {
  name = %q
  image_type = "qcow2"
  storage_provider_id = data.hpe_morpheus_storage_bucket.local.id
  os_type_id = data.hpe_morpheus_os_type.local.id
  file_path = %q
}
`, name, filename)
	}
	const address = "hpe_morpheus_image.local"
	const importedAddress = "hpe_morpheus_image.imported"
	imported := func(filename string) string {
		return strings.Replace(render(filename), `resource "hpe_morpheus_image" "local"`, `resource "hpe_morpheus_image" "imported"`, 1)
	}
	var identity string
	checkIdentity := func(state *terraform.State) error {
		object := state.RootModule().Resources[importedAddress]
		if object == nil {
			object = state.RootModule().Resources[address]
		}
		actual := object.Primary.ID
		if identity == "" {
			identity = actual
		} else if identity != actual {
			return fmt.Errorf("image ID changed from %s to %s", identity, actual)
		}

		return nil
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_7_0)},
		Steps: []resource.TestStep{
			{
				Config:   strings.Replace(render(first), "  file_path =", "  url = \"https://example.invalid/disk.qcow2\"\n  file_path =", 1),
				PlanOnly: true, ExpectError: regexp.MustCompile("cannot be specified"),
			},
			{Config: render(first), Check: resource.ComposeAggregateTestCheckFunc(checkIdentity,
				resource.TestCheckResourceAttr(address, "file_name", "disk.qcow2"),
				resource.TestCheckResourceAttr(address, "file_sha256", fmt.Sprintf("%x", sha256.Sum256([]byte("first image")))),
			)},
			{Config: render(first), PlanOnly: true},
			{
				ResourceName: address, ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"file_path", "file_sha256", "file_name", "file_size", "write_only_revision"},
			},
			{
				ResourceName: importedAddress, ImportState: true, ImportStatePersist: true, Config: imported(first),
				ImportStateIdFunc: func(_ *terraform.State) (string, error) { return identity, nil },
			},
			{Config: imported(first) + `
removed {
  from = hpe_morpheus_image.local
  lifecycle { destroy = false }
}
`, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(importedAddress, plancheck.ResourceActionUpdate),
			}}, Check: checkIdentity},
			{Config: imported(second), ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(importedAddress, plancheck.ResourceActionUpdate),
			}}, Check: resource.ComposeAggregateTestCheckFunc(checkIdentity,
				resource.TestCheckResourceAttr(importedAddress, "file_sha256", fmt.Sprintf("%x", sha256.Sum256([]byte("other image")))),
			)},
			{Config: imported(second), PlanOnly: true},
		},
	})
}
