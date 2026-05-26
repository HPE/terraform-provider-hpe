resource "hpe_morpheus_cypher" "example" {
  id    = "secret/my-api-key"
  value = "sk-abc123def456"
  ttl   = 0
}
