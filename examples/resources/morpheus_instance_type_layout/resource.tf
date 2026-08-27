resource "hpe_morpheus_instance_type_layout" "example" {
  instance_type_id = data.hpe_morpheus_instance_type.example.id
  name             = "todo_app_frontend"
  labels           = ["demo", "layout", "terraform"]
  version          = "1.0"
  technology       = "vmware"
}
