resource "hpe_morpheus_blueprint_app_kubernetes" "tfexample_kubernetes_app_blueprint_spec" {
  name              = "tf-kubernetes-app-blueprint-example-spec"
  description       = "tf example kubernetes app blueprint"
  category          = "k8s"
  source_type       = "spec"
  spec_template_ids = [2, 3]
}

