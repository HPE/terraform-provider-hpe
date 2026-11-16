resource "hpe_morpheus_host_profile" "example" {
  name        = "Example Host Profile"
  description = "An example host profile"

  config = {
    storage = {
      iscsi = {
        transport    = "tcp"
        chap_enabled = false
      }
      multipath = {
        find_multipaths     = true
        user_friendly_names = true
      }
    }
    system = {
      # config sub-attributes are Optional: one you omit stays null and
      # unmanaged, so Terraform neither sets nor reconciles it. Declare it
      # explicitly (e.g. timezone = "America/Denver") to manage it, and remove
      # it later to clear it.
      cpu = {
        power_management = "performance"
      }
      logging = {
        level = "info"
      }
    }
    baseline = {
      # Extra packages to install, one per line in package=version format.
      additional_packages = "openvswitch=2.17.0"
    }
  }
}
