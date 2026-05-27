resource "hpe_morpheus_backup" "example" {
  name            = "DB Backup"
  instance_id     = 1
  backup_type     = "lvmSnapshot"
  retention_count = 7
  enabled         = true
}
