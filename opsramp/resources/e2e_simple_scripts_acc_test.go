// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleScripts exercises the simple-scripts e2e scenario:
// script categories with hierarchy and a script with attachment and parameters.
func TestAccE2ESimpleScripts(t *testing.T) {
	t.Run("category hierarchy and script", func(t *testing.T) {
		parentCat := acctest.RandomName("e2e-script-parent")
		childCat := acctest.RandomName("e2e-script-child")
		scriptName := acctest.RandomName("e2e-script")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScriptCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleScriptsConfig(parentCat, childCat, scriptName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScriptCategoryExists(t, "hpe_opsramp_script_category.automation"),
						testAccEnsureScriptCategoryExists(t, "hpe_opsramp_script_category.linux"),
						testAccEnsureScriptExists(t, "hpe_opsramp_script.restart_service"),
						resource.TestCheckResourceAttr("hpe_opsramp_script_category.automation", "name", parentCat),
						resource.TestCheckResourceAttr("hpe_opsramp_script_category.linux", "name", childCat),
						resource.TestCheckResourceAttr("hpe_opsramp_script.restart_service", "name", scriptName),
						resource.TestCheckResourceAttr("hpe_opsramp_script.restart_service", "execution_type", "SHELL"),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleScriptsConfig(parentCat, childCat, scriptName string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_script_category" "automation" {
	name = "%s"
}

resource "hpe_opsramp_script_category" "linux" {
	name      = "%s"
	parent_id = hpe_opsramp_script_category.automation.uuid
}

resource "hpe_opsramp_script" "restart_service" {
	category_id     = hpe_opsramp_script_category.linux.uuid
	name            = "%s"
	description     = "Restart a service on a Linux machine."
	platforms       = ["LINUX"]
	execution_type  = "SHELL"
	install_timeout = 120

	attachment = {
		name = "restart_service_linux.sh"
		file = "#!/bin/bash\nsystemctl restart $1"
	}

	parameters = [
		{
			name          = "service_name"
			description   = ""
			default_value = ""
			type          = "REQUIRED"
			data_type     = "STRING"
		}
	]
}
`, acctest.ProviderConfigHCL(), parentCat, childCat, scriptName)
}
