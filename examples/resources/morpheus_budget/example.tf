resource "hpe_morpheus_budget" "example" {
  name        = "Q1 Cloud Budget"
  description = "First quarter cloud spending budget"
  year        = 2025
  interval    = "year"
  scope       = "account"
  enabled     = true
}
