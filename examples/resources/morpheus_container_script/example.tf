resource "hpe_morpheus_container_script" "example" {
  name          = "Install Dependencies"
  script_phase  = "provision"
  script_type   = "bash"
  script        = file("/path-to-file")
  sudo_user     = true
  fail_on_error = true
}
