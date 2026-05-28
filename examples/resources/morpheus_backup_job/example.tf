resource "hpe_morpheus_backup_job" "example" {
  name            = "Nightly Backup Job"
  code            = "nightly-backup"
  retention_count = 14
  enabled         = true
}
