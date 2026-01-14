resource "hpe_morpheus_network_domain" "example" {
  name        = "tfexampledomain"
  description = "Terraform example network domain"
  public_zone = true
  visibility  = "private"
  tenant_id   = 1
  active      = true
}
