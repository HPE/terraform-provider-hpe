// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerscript_test

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

// Removing the category, run-as user or labels clears it on the script
// (MORPH-17990); the labels are Optional+Computed and cleared once an apply
// has recorded them, the others Optional-only.
func TestAccMorpheusContainerScriptRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_container_script.test"

	config := func(with ...string) string {
		set := map[string]string{
			"category":    `category    = "clr"`,
			"run_as_user": `run_as_user = "root"`,
			"labels":      `labels      = ["terraform"]`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_container_script" "test" {
  name         = %q
  script_phase = "provision"
  script_type  = "bash"
  script       = "echo hello"
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"category", "run_as_user", "labels"}

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
					resource.TestCheckResourceAttr(rn, "category", "clr"),
					resource.TestCheckResourceAttr(rn, "run_as_user", "root"),
					resource.TestCheckResourceAttr(rn, "labels.#", "1"),
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
