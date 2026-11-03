// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const (
	woAPIPath = "/api/tasks"
	woAPIKey  = "task.taskOptions.password"
)

// woConfig renders a remote-execution PowerShell task so the remote target password
// is meaningful. extra is appended inside the block (the deprecated marker
// line, or ""); unrelated toggles the code to force an unrelated change.
func woConfig(name, value, extra string, unrelated bool) string {
	code := name
	if unrelated {
		code = name + "-changed"
	}

	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_task_powershell_script" "wo" {
  name                      = "` + name + `"
  code                      = "` + code + `"
  source_type               = "local"
  script_content            = "echo hi"
  execute_target            = "remote"
  remote_target_host        = "10.0.0.1"
  remote_target_port        = "5985"
  remote_target_username    = "Administrator"
  remote_target_password_wo = "` + value + `"` + extra + `
}
`
}

// TestAccMorpheusTaskPowershellScriptResourceWriteOnlyOk runs the standard write-only
// change-detection matrix: changes to remote_target_password_wo are detected by
// digest, the deprecated marker still works, and an unrelated update does not
// send the key.
func TestAccMorpheusTaskPowershellScriptResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	addr := "hpe_morpheus_task_powershell_script.wo"

	const (
		v1 = "Secret123!"
		v2 = "Secret456!"
		v3 = "Secret789!"
	)

	rec := testhelpers.RecordRequests(t)

	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "remote_target_password_wo",
		Marker:    "remote_target_password_wo_version",
		APIPath:   woAPIPath,
		APIKey:    woAPIKey,
		Config: func(value, extra string, unrelated bool) string {
			return woConfig(name, value, extra, unrelated)
		},
	}

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, v1, v2, v3)
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, v3)...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
