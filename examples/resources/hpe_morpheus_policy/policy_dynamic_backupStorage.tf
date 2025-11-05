# NOTE: This policy type is currently not supported due to a bug in the Morpheus API
# Until this is resolved in the Morpheus API, this example is commented out

# # Backup Targets Policy - Restricts available backup storage targets
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
#     backupStorageIds = [5, 6] # Array of backup storage IDs
#   }
# }
