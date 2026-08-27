resource "hpe_morpheus_app_blueprint_helm" "example" {
  name           = "helmappblueprint"
  description    = "tf example helm app blueprint"
  category       = "helm"
  integration_id = 3
  repository_id  = 1
  version_ref    = "main"
  working_path   = "./test"
}
