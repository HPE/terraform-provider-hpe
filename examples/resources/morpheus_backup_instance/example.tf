resource "hpe_morpheus_backup_instance" "example" {
  name                = "Example Backup"
  instance_id         = 1
  job_id              = 1
  storage_provider_id = 1
  enabled             = true
}
