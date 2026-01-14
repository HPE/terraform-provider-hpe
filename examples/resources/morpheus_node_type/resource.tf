resource "hpe_morpheus_node_type" "example" {
  name             = "tf_example_node_type"
  short_name       = "tfexamplenodetype"
  labels           = ["demo", "nodeType", "terraform"]
  technology       = "vmware"
  version          = "2.0"
  category         = "tfexample"
  virtual_image_id = 10

  service_port {
    name     = "web"
    port     = "8080"
    protocol = "HTTP"
  }

  service_port {
    name     = "secureweb"
    port     = "8443"
    protocol = "HTTPS"
  }
}

