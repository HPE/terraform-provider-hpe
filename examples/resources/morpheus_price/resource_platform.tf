resource "hpe_morpheus_price" "example" {
  name          = "terraform-test"
  code          = "terraform-test"
  tenant_id     = 1
  price_type    = "platform"
  platform      = "linux"
  price_unit    = "minute"
  incur_charges = "always"
  currency      = "USD"
  cost          = 38.00
}
