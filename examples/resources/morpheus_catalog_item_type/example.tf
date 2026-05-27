resource "hpe_morpheus_catalog_item_type" "example" {
  name        = "Ubuntu VM"
  description = "Standard Ubuntu virtual machine"
  type        = "instance"
  enabled     = true
  visibility  = "public"
  featured    = true
}
