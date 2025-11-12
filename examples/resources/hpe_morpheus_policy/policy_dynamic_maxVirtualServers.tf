# Max Virtual Servers Policy - Limits virtual server count
resource "hpe_morpheus_policy" "max_virtual_servers" {
  name                     = "Max Virtual Servers Policy"
  description              = "Limit maximum virtual server count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxVirtualServers"
  }

  config = {
    maxVirtualServers = "10" # Maximum number of virtual servers
  }
}
