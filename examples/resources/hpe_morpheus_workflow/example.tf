resource "hpe_morpheus_workflow" "example" {
  name        = "App Deployment"
  description = "Standard application deployment workflow"
  type        = "provision"
  platform    = "all"
  visibility  = "public"
}
