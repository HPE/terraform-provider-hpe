// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleIAM exercises the simple-iam e2e scenario:
// permission sets, roles, users, and user groups working together.
func TestAccE2ESimpleIAM(t *testing.T) {
	t.Run("full IAM stack", func(t *testing.T) {
		adminPerms := acctest.RandomName("e2e-admin-perms")
		viewPerms := acctest.RandomName("e2e-view-perms")
		adminRole := acctest.RandomName("e2e-admin-role")
		viewRole := acctest.RandomName("e2e-view-role")
		userName := acctest.RandomName("e2e-user")
		groupName := acctest.RandomName("e2e-group")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleIAMConfig(adminPerms, viewPerms, adminRole, viewRole, userName, groupName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsurePermissionSetExists(t, "opsramp_permission_set.admin_perms"),
						testAccEnsurePermissionSetExists(t, "opsramp_permission_set.view_perms"),
						testAccEnsureRoleExists(t, "opsramp_role.admin_role"),
						testAccEnsureRoleExists(t, "opsramp_role.view_role"),
						testAccEnsureUserExists(t, "opsramp_user.admin"),
						testAccEnsureUserGroupExists(t, "opsramp_user_group.admin_group"),
						resource.TestCheckResourceAttr("opsramp_permission_set.admin_perms", "name", adminPerms),
						resource.TestCheckResourceAttr("opsramp_permission_set.view_perms", "name", viewPerms),
						resource.TestCheckResourceAttr("opsramp_role.admin_role", "name", adminRole),
						resource.TestCheckResourceAttr("opsramp_role.view_role", "name", viewRole),
						resource.TestCheckResourceAttr("opsramp_user.admin", "login_name", userName),
						resource.TestCheckResourceAttr("opsramp_user_group.admin_group", "name", groupName),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleIAMConfig(adminPerms, viewPerms, adminRole, viewRole, userName, groupName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_permission_set" "admin_perms" {
	name        = "%s"
	description = "Full administrative access"

	permissions = [
		{
			name = "Alerts"
			type = "Manage Alerts "
		},
		{
			name = "Devices"
			type = "Manage Device "
		}
	]
}

resource "opsramp_permission_set" "view_perms" {
	name        = "%s"
	description = "View-only access"

	permissions = [
		{
			name = "Alerts"
			type = "View Alerts "
		},
		{
			name = "Devices"
			type = "View Device "
		}
	]
}

resource "opsramp_role" "admin_role" {
	name        = "%s"
	description = "Administrative role"

	permissions = [
		opsramp_permission_set.admin_perms.unique_id
	]
}

resource "opsramp_role" "view_role" {
	name        = "%s"
	description = "View-only role"

	permissions = [
		opsramp_permission_set.view_perms.unique_id
	]
}

resource "opsramp_user" "admin" {
	login_name = "%s"
	password   = "E2ETestP@ss1234!"
	first_name = "E2E"
	last_name  = "Admin"
	email      = "%s@example.com"
	time_zone  = "Europe/Paris"
	country    = "Spain"

	user_notifications = [
		{
			notify_type             = "Account Information"
			notify_method           = "Email"
			notify_input_type       = "Primary Email"
			notify_recurring_report = false
		},
		{
			notify_type             = "Alert Notification"
			notify_method           = "No Notify"
			notify_recurring_report = false
		},
		{
			notify_type             = "Report Notification"
			notify_method           = "No Notify"
			notify_recurring_report = false
		}
	]

	change_password = false
}

resource "opsramp_user_group" "admin_group" {
	name        = "%s"
	description = "E2E test admin group"

	roles = [
		opsramp_role.admin_role.id
	]

	users = [
		opsramp_user.admin.id
	]
}
`, acctest.ProviderConfigHCL(), adminPerms, viewPerms, adminRole, viewRole, userName, userName, groupName)
}
