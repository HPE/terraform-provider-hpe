# Tags Policy - Enforces instance tagging
resource "hpe_morpheus_policy" "tags" {
  name                     = "Tags Policy"
  description              = "Enforce instance tagging requirements"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "tags"
  }

  config = {
    strict      = true          # Strict enforcement
    key         = "environment" # Tag key to enforce
    value       = "production"  # Tag value (optional, can be left empty for any value)
    valueListId = ""            # ID of value from value list (optional)
  }
}
