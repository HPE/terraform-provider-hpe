resource "hpe_morpheus_app_blueprint_kubernetes" "tfexample_kubernetes_app_blueprint_yaml" {
  name              = "tf-kubernetes-app-blueprint-example-yaml"
  description       = "tf example kubernetes app blueprint"
  category          = "k8s"
  source_type       = "yaml"
  blueprint_content = <<TFEOF
...
TFEOF
}
