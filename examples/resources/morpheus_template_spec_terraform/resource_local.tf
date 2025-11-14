resource "hpe_morpheus_template_spec_terraform" "tfexample_terraform_spec_terraform_local" {
  name         = "tf-terraform-spec-example-local"
  source_type  = "local"
  spec_content = <<EOF
resource "aws_instance" "instance_1" {
  ami           = "ami-0b91a410940e82c54"
  instance_type = "t2.micro"
}
EOF
}