resource "hpe_morpheus_credential" "example" {
  name        = "Service Account"
  type        = "username-password"
  description = "Service account credentials"
  username    = "svc-account"
  password    = "supersecret"
  enabled     = true
}
