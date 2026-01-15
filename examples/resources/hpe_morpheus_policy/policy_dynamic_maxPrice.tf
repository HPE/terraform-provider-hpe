# Budget Policy - Limits instance costs
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
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
    # Required
    maxPrice = 1000 # Maximum price limit

    # Optional
    maxPriceCurrency = "USD"   # Currency code
    maxPriceUnit     = "month" # Options: "hour", "month"
  }
}
