data "hpe_morpheus_cloud" "nsxt" {
  name = "NSX-T Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

resource "hpe_morpheus_security_group" "nsxt" {
  name                           = "nsxt-app-tier"
  description                    = "Security group for NSX-T application tier"
  cloud_id                       = data.hpe_morpheus_cloud.nsxt.id
  active                         = true
  visibility                     = "public"
  resource_permission_groups_all = false
  resource_permission_group_ids  = [data.hpe_morpheus_group.example.id]
}

resource "hpe_morpheus_security_group_rule" "app_http" {
  security_group_id = hpe_morpheus_security_group.nsxt.id
  name              = "Allow HTTP"
  protocol          = "tcp"
  rule_type         = "customRule"
  direction         = "ingress"
  port_range        = "8080"
  source            = "10.0.0.0/8"
  policy            = "accept"
}
