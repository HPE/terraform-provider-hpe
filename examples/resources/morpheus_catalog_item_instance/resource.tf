resource "hpe_morpheus_catalog_item_instance" "tf_example_instance_catalog_item" {
  name        = "tfexample_instance_catalog"
  description = "terraform example instance catalog item"
  image_path  = "tfexample.png"
  image_name  = "tfexample.png"
  enabled     = true
  featured    = true
  content     = <<TFEOF
  {"name":"test"}
  TFEOF
  config      = <<TFEOF
  {"name":"test"}
  TFEOF
  visibility  = "private"
}
