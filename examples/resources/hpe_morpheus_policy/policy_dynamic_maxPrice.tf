# Budget Policy - Limits instance costs
resource "hpe_morpheus_policy" "budget" {
  name                     = "Budget Policy"
  description              = "Limit maximum instance costs"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPrice"
  }

  config = {
    maxPrice         = "1000"  # Maximum price limit
    maxPriceCurrency = "USD"   # Currency code
    maxPriceUnit     = "month" # Options: "hour", "day", "month", "year"
  }
}
