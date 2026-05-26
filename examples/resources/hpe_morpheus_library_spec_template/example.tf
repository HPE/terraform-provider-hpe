resource "hpe_morpheus_library_spec_template" "example" {
  name    = "Kubernetes Deployment"
  type    = "kubernetes"
  source  = "local"
  content = file("/path-to-file")
}
