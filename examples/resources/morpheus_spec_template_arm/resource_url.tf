resource "hpe_morpheus_spec_template_arm" "tfexample_arm_spec_template_url" {
  name        = "tf-arm-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/spec.json"
}
