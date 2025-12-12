resource "hpe_morpheus_spec_template_kubernetes" "tfexample_kubernetes_spec_template_url" {
  name        = "tf-kubernetes-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/spec.yaml"
}
