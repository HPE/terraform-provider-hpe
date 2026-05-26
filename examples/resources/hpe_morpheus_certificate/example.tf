resource "hpe_morpheus_certificate" "example" {
  name      = "Example Certificate"
  cert_file = file("/path-to-file")
  key_file  = file("/path-to-file")
}
