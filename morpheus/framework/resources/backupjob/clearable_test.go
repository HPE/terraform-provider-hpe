// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupjob_test

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

// Removing an optional attribute clears it (MORPH-17990): the attributes here
// are Optional without Computed, so Terraform plans the null and the apply
// used to fail with an inconsistent result because Update never sent it.
func TestAccMorpheusBackupJobRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_backup_job.test"

	config := func(with ...string) string {
		set := map[string]string{
			"retention_count": `retention_count = 3`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_backup_job" "test" {
  name = %q
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"retention_count"}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "retention_count"),
				),
			},
			{
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("retention_count"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "retention_count"),
				),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
