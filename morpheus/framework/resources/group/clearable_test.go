// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package group_test

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

// Removing an optional attribute clears it on the group (MORPH-17990): code
// and location (Optional, not computed — the apply used to fail with an
// inconsistent result), labels and cloud_ids (Optional+Computed — the removal
// used to be silently ignored).
func TestAccMorpheusGroupRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)
	const rn = "hpe_morpheus_group.test"

	config := func(with ...string) string {
		set := map[string]string{
			"code":      `code      = "` + code + `"`,
			"location":  `location  = "here"`,
			"labels":    `labels    = ["Label1", "Label2"]`,
			"cloud_ids": `cloud_ids = [1]`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_group" "test" {
  name = %q
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"code", "location", "labels", "cloud_ids"}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "code", code),
					resource.TestCheckResourceAttr(rn, "location", "here"),
					resource.TestCheckResourceAttr(rn, "labels.#", "2"),
					resource.TestCheckResourceAttr(rn, "cloud_ids.#", "1"),
				),
			},
			{
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("code"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("location"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("labels"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("cloud_ids"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "code"),
					resource.TestCheckNoResourceAttr(rn, "location"),
					resource.TestCheckNoResourceAttr(rn, "labels.#"),
					resource.TestCheckNoResourceAttr(rn, "cloud_ids.#"),
				),
			},
			{
				// Nothing recorded as configured any more: an empty plan.
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
