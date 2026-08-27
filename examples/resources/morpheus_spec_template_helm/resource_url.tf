resource "hpe_morpheus_spec_template_helm" "example" {
  name        = "tf-helm-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/chart.yaml"
}
