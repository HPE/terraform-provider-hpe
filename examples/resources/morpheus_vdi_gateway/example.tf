resource "hpe_morpheus_vdi_gateway" "example" {
  name        = "Primary VDI Gateway"
  gateway_url = "https://vdi-gateway.example.com"
  description = "Main VDI gateway for remote access"
}
