# Backup Creation Policy
resource "hpe_morpheus_policy" "backup_creation" {
  name                     = "Backup Creation Policy"
  description              = "Enforce backup creation for instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createBackup"
  }

  config = {
    createBackupType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)
    createBackup     = true   # Enforce backup creation
  }
}
