resource "hpe_morpheus_key_pair" "example" {
  name        = "example-key-pair"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC..."
  private_key = <<EOF
  -----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----
  EOF
}
