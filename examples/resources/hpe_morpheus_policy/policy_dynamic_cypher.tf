# Cypher Access Policy - Controls access to Cypher secrets
resource "hpe_morpheus_policy" "cypher_access" {
  name                     = "Cypher Access Policy"
  description              = "Control Cypher key access permissions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "cypher"
  }

  config = {
    keyPattern = "secret/*" # Pattern to match Cypher keys (e.g., "secret/*", "password/*")
    read       = true       # Allow read access
    write      = true       # Allow write access
    update     = true       # Allow update access
    delete     = false      # Deny delete access
    list       = true       # Allow list access
  }
}
