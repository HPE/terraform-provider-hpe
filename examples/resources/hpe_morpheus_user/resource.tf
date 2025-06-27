resource "hpe_morpheus_user" "example" {
  tenant_id                   = 1
  username                    = "testacc-example"
  email                       = "user@example.com"
  password_wo                 = "Secret123!"
  password_wo_version         = 1
  role_ids                    = [1]
  first_name                  = "Joe"
  last_name                   = "User"
  linux_username              = "linuser"
  linux_password_wo           = "Linux123!"
  linux_password_wo_version   = 1
  linux_key_pair_id           = 100
  receive_notifications       = false
  windows_username            = "winuser"
  windows_password_wo         = "Windows123!"
  windows_password_wo_version = 1
}
