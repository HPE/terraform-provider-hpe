resource "hpe_morpheus_certificate" "example" {
  name        = "wildcard-example-com"
  cert_file   = file("${path.module}/certs/wildcard.crt")
  key_file    = file("${path.module}/certs/wildcard.key")
  domain_name = "*.example.com"
  description = "Wildcard certificate for example.com"
}
