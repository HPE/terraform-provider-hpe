resource "hpe_morpheus_library_file_template" "example" {
  name           = "Nginx Config"
  file_name      = "nginx.conf"
  file_path      = "/etc/nginx"
  template_phase = "provision"
  template       = file("/path-to-file")
}
