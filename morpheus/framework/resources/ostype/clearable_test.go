// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostype_test

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
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing an optional attribute clears it on the OS type (MORPH-17990): the
// description, family, version and cloud-init version are Optional+Computed
// and cleared once an apply has recorded them; the category, codename, name
// and vendor are Optional-only and Terraform plans the null itself.
func TestAccMorpheusOsTypeRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_os_type.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":        `description        = "to be cleared"`,
			"os_family":          `os_family          = "debian"`,
			"os_version":         `os_version         = "12"`,
			"cloud_init_version": `cloud_init_version = "23.4"`,
			"category":           `category           = "linux"`,
			"os_codename":        `os_codename        = "bookworm"`,
			"os_name":            `os_name            = "debian"`,
			"vendor":             `vendor             = "debian"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_os_type" "test" {
  name      = %q
  code      = "clr.test.os"
  platform  = "linux"
  bit_count = 64
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "os_family", "os_version", "cloud_init_version", "category", "os_codename", "os_name", "vendor"}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	var noAttrs []resource.TestCheckFunc
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		noAttrs = append(noAttrs, resource.TestCheckNoResourceAttr(rn, a))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "os_family", "debian"),
					resource.TestCheckResourceAttr(rn, "os_version", "12"),
					resource.TestCheckResourceAttr(rn, "cloud_init_version", "23.4"),
					resource.TestCheckResourceAttr(rn, "category", "linux"),
					resource.TestCheckResourceAttr(rn, "os_codename", "bookworm"),
					resource.TestCheckResourceAttr(rn, "os_name", "debian"),
					resource.TestCheckResourceAttr(rn, "vendor", "debian"),
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
