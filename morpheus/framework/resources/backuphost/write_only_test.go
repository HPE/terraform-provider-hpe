// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupjob"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusBackupHostResourceWriteOnlyOk: changes to ssh_password_wo
// are detected by digest (MORPH-18005); the deprecated marker still works;
// both import forms leave a clean plan.
func TestAccMorpheusBackupHostResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Backup)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	job, err := backupjob.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:      "hpe_morpheus_backup_host.wo",
		Attribute:    "ssh_password_wo",
		Marker:       "ssh_password_wo_version",
		APIPath:      "/api/backups",
		APIKey:       "backup.targetPassword",
		CreateAPIKey: "backup.sshPassword",
		Config: func(value, extra string, unrelated bool) string {
			path := "/etc/hostname"
			if unrelated {
				path = "/etc/hosts"
			}

			return testhelpers.ProviderBlock() + job + `
resource "hpe_morpheus_backup_host" "wo" {
  name                = "` + name + `"
  host_id             = 139
  job_id              = hpe_morpheus_backup_job.example.id
  backup_type_code    = "fileBackup"
  path                = "` + path + `"
  storage_provider_id = 1
  enabled             = true
  ssh_username        = "root"
  ssh_password_wo     = "` + value + `"` + extra + `
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	// ssh_username is not returned by the API on a plain GET, so import cannot
	// verify it (as in the existing import test).
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!", "ssh_username")...)

	// Destroying the backup host also destroys its job on the appliance, so
	// remove the host first and then refresh, as the existing tests do.
	steps = append(steps,
		resource.TestStep{Config: testhelpers.ProviderBlock() + job, ExpectNonEmptyPlan: true},
		resource.TestStep{RefreshState: true, ExpectNonEmptyPlan: true},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
