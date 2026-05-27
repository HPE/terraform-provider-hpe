resource "hpe_morpheus_library_instance_type" "example" {
  name        = "Custom App Server"
  code        = "custom-app-server"
  description = "Custom application server instance type"
  category    = "web"
  visibility  = "public"
}
