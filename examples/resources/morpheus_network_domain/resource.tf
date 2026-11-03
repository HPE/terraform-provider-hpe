resource "hpe_morpheus_network_domain" "example" {
  name        = "tfexampledomain"
  description = "Terraform example network domain"
  public_zone = false
  visibility  = "private"
  active      = true

  # Write-only: never stored in state. Changes are detected through a
  # private-state digest (see write_only_revision).
  domain_password_wo = "Secret123!"
  domain_username    = "joiner"

  # domain_password_wo_version is deprecated: changes to domain_password_wo are
  # detected automatically. It will be removed in the 2.6 release.
}
