resource "hpe_morpheus_backup_host" "example" {
  name                = "Host Backup"
  host_id             = 1
  job_id              = 1
  backup_type_code    = "fileBackup"
  path                = "/etc/hostname"
  storage_provider_id = 1
  enabled             = true
}
