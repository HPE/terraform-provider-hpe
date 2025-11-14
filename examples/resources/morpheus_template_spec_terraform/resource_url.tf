resource "hpe_morpheus_template_spec_terraform" "tfexample_terraform_spec_terraform_url" {
  name        = "tf-terraform-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/spec.tf"
}