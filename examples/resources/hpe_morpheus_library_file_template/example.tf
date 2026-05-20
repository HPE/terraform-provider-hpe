resource "hpe_morpheus_library_file_template" "example" {
  name           = "Nginx Config"
  file_name      = "nginx.conf"
  file_path      = "/etc/nginx"
  template_phase = "provision"
  template       = "server {\n  listen 80;\n  server_name <%= instance.hostname %>;\n}"
}
