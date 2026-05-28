data "hpe_morpheus_cloud" "azure" {
  name = "Azure Cloud"
}

data "hpe_morpheus_tenant" "example" {
  name = "Example Tenant"
}

resource "hpe_morpheus_security_group" "azure" {
  name                           = "azure-web-servers"
  description                    = "Security group for Azure web servers"
  cloud_id                       = data.hpe_morpheus_cloud.azure.id
  active                         = true
  visibility                     = "private"
  tenant_ids                     = [data.hpe_morpheus_tenant.example.id]
  resource_permission_groups_all = true
}

resource "hpe_morpheus_security_group_rule" "https" {
  security_group_id = hpe_morpheus_security_group.azure.id
  name              = "Allow HTTPS"
  protocol          = "tcp"
  rule_type         = "customRule"
  direction         = "ingress"
  port_range        = "443"
  source            = "0.0.0.0/0"
  policy            = "accept"
}
