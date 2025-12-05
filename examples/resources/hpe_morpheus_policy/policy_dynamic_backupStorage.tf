# # Backup Targets Policy - Restricts available backup storage targets
# Allowed associated_resource_types: Group, Cloud, User, Global, Role
# Tenant specification: NOT allowed (cannot specify tenants array)
# resource "hpe_morpheus_policy" "backup_targets" {
#   name                     = "Backup Targets Policy"
#   description              = "Restrict available backup targets"
#   associated_resource_type = "User"
#   associated_resource_id   = 9969
#   enabled                  = true
#
#   policy_type = {
#     code = "backupStorage"
#   }
#
#   config = {
#     backupStorageIds = [5, 6] # Array of backup storage IDs to restrict available backup targets
#   }
# }
