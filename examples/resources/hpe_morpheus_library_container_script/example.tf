resource "hpe_morpheus_library_container_script" "example" {
  name          = "Install Dependencies"
  script_phase  = "provision"
  script_type   = "bash"
  script        = "#!/bin/bash\napt-get update && apt-get install -y curl wget"
  sudo_user     = true
  fail_on_error = true
}
