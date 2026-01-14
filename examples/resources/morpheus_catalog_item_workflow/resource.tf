resource "hpe_morpheus_catalog_item_workflow" "example" {
  name                 = "tfexample_workflow_catalog_item"
  description          = "Example Terraform workflow catalog item"
  logo_image_path      = "wordpress.png"
  logo_image_name      = "wordpress.png"
  dark_logo_image_path = "wordpressbak.png"
  dark_logo_image_name = "wordpressbak.png"
  enabled              = true
  featured             = true
  labels               = ["terraform", "demo"]
  workflow_id          = 1
  context_type         = "appliance"
  content              = "Example catalog content"
  visibility           = "public"
}
