resource "hpe_morpheus_price_set" "example" {
  name        = "Standard Price Set"
  code        = "priceset.standard"
  price_unit  = "hour"
  type        = "compute_plus_storage"
  region_code = "us-east-1"
}
