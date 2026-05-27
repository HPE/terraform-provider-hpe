resource "hpe_morpheus_backup_job" "example" {
  name            = "Nightly Backup Job"
  code            = "nightly-backup"
  schedule_id     = 1
  retention_count = 14
  enabled         = true
}
