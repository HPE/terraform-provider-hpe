// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskshellscript_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func legacyConfig(name, password, code string) string {
	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_task_shell_script" "legacy" {
  name                   = "` + name + `"
  code                   = "` + code + `"
  source_type            = "local"
  script_content         = "echo hi"
  execute_target         = "remote"
  remote_target_host     = "10.0.0.1"
  remote_target_port     = "22"
  remote_target_username = "root"
  remote_target_password = "` + password + `"
}
`
}

// TestAccMorpheusTaskShellScriptResourceDeprecatedPasswordOk proves a
// configuration written for the SDKv2 resource still works: the deprecated
// plaintext remote_target_password is sent on create and when it changes, and
// not on an unrelated update (omitting it leaves the task's password alone).
func TestAccMorpheusTaskShellScriptResourceDeprecatedPasswordOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	rec := testhelpers.RecordRequests(t)

	const apiPath, apiKey = "/api/tasks", "task.taskOptions.password"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy123!", name),
				Check: resource.ComposeAggregateTestCheckFunc(
					rec.ExpectSent(apiPath, apiKey, true),
					resource.TestCheckResourceAttr("hpe_morpheus_task_shell_script.legacy", "remote_target_password", "Legacy123!"),
				),
			},
			{
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy123!", name+"-changed"),
				Check:     rec.ExpectSent(apiPath, apiKey, false),
			},
			{
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy456!", name+"-changed"),
				Check:     rec.ExpectSent(apiPath, apiKey, true),
			},
		},
	})
}
