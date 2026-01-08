resource "hpe_morpheus_instance_type_layout" "example" {
  instance_type_id = hpe_morpheus_instance_type.tf_example_instance_type.id
  name             = "todo_app_frontend"
  labels           = ["demo", "layout", "terraform"]
  version          = "1.0"
  technology       = "vmware"
}
