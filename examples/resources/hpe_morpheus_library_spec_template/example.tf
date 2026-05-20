resource "hpe_morpheus_library_spec_template" "example" {
  name    = "Kubernetes Deployment"
  type    = "kubernetes"
  source  = "local"
  content = "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: example-app"
}
