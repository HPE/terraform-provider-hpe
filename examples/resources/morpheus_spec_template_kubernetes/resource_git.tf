resource "hpe_morpheus_spec_template_kubernetes" "tfexample_kubernetes_spec_template_git" {
  name          = "tf-kubernetes-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "./spec.yaml"
}
