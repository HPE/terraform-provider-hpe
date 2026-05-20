resource "hpe_morpheus_price" "example" {
  name       = "Standard Compute"
  code       = "price.compute.standard"
  price_type = "compute"
  price_unit = "hour"
  cost       = 0.05
  currency   = "USD"
  markup     = 20.0
}
