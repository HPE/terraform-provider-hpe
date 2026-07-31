# User Group for Client Admins
resource "hpe_opsramp_user_group" "client_admin_group" {
  name        = "Admin User Group"
  description = "User group for client administrators"

  roles = [
    hpe_opsramp_role.client_admin_role.id
  ]

  users = [
    hpe_opsramp_user.admin.id
  ]
}