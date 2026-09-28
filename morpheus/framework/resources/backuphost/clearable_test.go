// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupjob"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing path or ssh_username clears it on the backup (MORPH-17990); both
// are Optional without Computed, so Terraform plans the null and Update sends
// it. Needs an appliance with a backup-capable host (the fixture the other
// backup host tests use); skipped without the Backup capability.
func TestAccMorpheusBackupHostRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_backup_host.test"

	jobConfig, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	config := func(with ...string) string {
		set := map[string]string{
			"path":         `path         = "/etc/hostname"`,
			"ssh_username": `ssh_username = "backup"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + jobConfig + fmt.Sprintf(`
resource "hpe_morpheus_backup_host" "test" {
  name                = %q
  host_id             = 139
  job_id              = hpe_morpheus_backup_job.example.id
  backup_type_code    = "fileBackup"
  storage_provider_id = 1
%s
}
`, name, strings.Join(lines, "\n"))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config("path", "ssh_username"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "path", "/etc/hostname"),
					resource.TestCheckResourceAttr(rn, "ssh_username", "backup"),
				),
			},
			{
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("path"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("ssh_username"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "path"),
					resource.TestCheckNoResourceAttr(rn, "ssh_username"),
				),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
