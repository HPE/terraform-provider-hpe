resource "hpe_morpheus_environment" "example" {
  name        = "staging"
  code        = "staging"
  description = "Staging environment"
  visibility  = "public"
  active      = true
}
