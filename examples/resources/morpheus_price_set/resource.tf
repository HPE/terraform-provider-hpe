resource "hpe_morpheus_price_set" "example" {
  name        = "terraform-test"
  code        = "terraform-test"
  region_code = "us-west-2"
  price_unit  = "minute"
  type        = "fixed"
  price_ids   = [1]
}
