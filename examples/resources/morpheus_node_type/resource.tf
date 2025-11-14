resource "hpe_morpheus_node_type" "tf_example_node" {
  name             = "tf_example_node_type"
  short_name       = "tfexamplenodetype"
  labels           = ["demo", "nodeType", "terraform"]
  technology       = "vmware"
  version          = "2.0"
  category         = "tfexample"
  virtual_image_id = 10

  file_template_ids = [
    data.hpe_morpheus_file_template.tfexample.id,
    113
  ]

  script_template_ids = [
    data.hpe_morpheus_script_template.tfscript1.id,
    data.hpe_morpheus_script_template.tfscript2.id
  ]

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

