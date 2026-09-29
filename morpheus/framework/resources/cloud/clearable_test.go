// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// Removing an optional attribute clears it on the cloud (MORPH-17990): the
// Optional+Computed code and labels, once an apply has recorded them, and the
// Optional-only location and configuration keys, whose removal used to fail
// the apply. A standard (generic) cloud carries all of them.
func TestAccMorpheusCloudRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)
	const rn = "hpe_morpheus_cloud.test"

	config := func(with ...string) string {
		set := map[string]string{
			"code":                `code                = "` + code + `"`,
			"labels":              `labels              = ["terraform", "acctest"]`,
			"location":            `location            = "somewhere"`,
			"appliance_url":       `appliance_url       = "https://appliance.example.test"`,
			"data_center_name":    `data_center_name    = "aDatacenter"`,
			"external_id":         `external_id         = "` + code + `"`,
			"import_existing_vms": `import_existing_vms = "off"`,
			"keyboard_layout":     `keyboard_layout     = "us"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + fmt.Sprintf(`
resource "hpe_morpheus_cloud" "test" {
  name            = %q
  tenant_id       = data.hpe_morpheus_whoami.current.tenant_id
  group_id        = 1
  cloud_type_code = "standard"
  visibility      = "private"
%s

  config = {
    certificateProvider = "internal"
  }
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{
		"code", "labels", "location",
		"appliance_url", "data_center_name", "external_id", "import_existing_vms", "keyboard_layout",
	}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	noAttrs := []resource.TestCheckFunc{}
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		key := a
		if a == "labels" {
			key = "labels.#"
		}
		noAttrs = append(noAttrs, resource.TestCheckNoResourceAttr(rn, key))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "code", code),
					resource.TestCheckResourceAttr(rn, "labels.#", "2"),
					resource.TestCheckResourceAttr(rn, "location", "somewhere"),
					resource.TestCheckResourceAttr(rn, "appliance_url", "https://appliance.example.test"),
					resource.TestCheckResourceAttr(rn, "data_center_name", "aDatacenter"),
					resource.TestCheckResourceAttr(rn, "external_id", code),
					resource.TestCheckResourceAttr(rn, "import_existing_vms", "off"),
					resource.TestCheckResourceAttr(rn, "keyboard_layout", "us"),
				),
			},
			{
				Config:           config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: nullChecks},
				Check:            resource.ComposeAggregateTestCheckFunc(noAttrs...),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
