resource "hpe_morpheus_spec_template_helm" "example" {
  name          = "tf-helm-spec-example-git"
  source_type   = "repository"
  repository_id = "2"
  version_ref   = "main"
  spec_path     = "./spec.yaml"
}
