data "hpe_opsramp_device_management_policy" "test_policy" {
  name = "mysql-servers"
}

output "device_management_policy_id" {
  value = data.hpe_opsramp_device_management_policy.test_policy.id
}
