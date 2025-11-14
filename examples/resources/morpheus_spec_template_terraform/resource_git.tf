resource "hpe_morpheus_spec_template_terraform" "tfexample_terraform_spec_terraform_git" {
  name          = "tf-terraform-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "Instance Types/Terraform/CloudResource/aws/vpc.tf"
}
