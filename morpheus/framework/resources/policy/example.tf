resource "hpe_morpheus_policy" "role_policy" {
  name = "TestMaxMemoryRolePolicy"
  description = "Example role-scoped policy"
  associated_resource_type = "Role"
  associated_resource_id = 1
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
}
