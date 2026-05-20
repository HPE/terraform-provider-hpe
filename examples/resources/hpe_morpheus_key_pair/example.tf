resource "hpe_morpheus_key_pair" "example" {
  name       = "deploy-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ... user@host"
}
