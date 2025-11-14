resource "hpe_morpheus_template_spec_arm" "tfexample_arm_spec_template_git" {
  name          = "tf-arm-spec-example-git"
  source_type   = "repository"
  repository_id = 2
  version_ref   = "main"
  spec_path     = "./test.json"
}