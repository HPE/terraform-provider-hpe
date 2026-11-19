# End-to-end: software-defined networking on an HPE VM Essentials (HVM) cluster
# with the HPE Morpheus SDN network plugin (a Juniper CN2 controller).
#
# The flow registers the CN2 controller as the cloud's network server, creates
# an SDN virtual network on the cluster, a logical router joining it, a BGP
# peering to the data-centre fabric, a load balancer with an HTTP health check,
# and finally allocates a floating IP from the router's pool. (Attaching the
# floating IP to an instance is a separate resource, pending platform fixes.)
#
# Prerequisites that exist outside Terraform: an HVM cloud with a cluster
# whose hosts can reach the CN2 controller.
#
# The resources this example uses were verified against a Morpheus 9.1.0.2 lab
# appliance with the SDN plugin; the example itself has not been applied as a
# whole, so validate the ids and addresses for your environment.

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.3.0"
    }
  }
}

provider "hpe" {
  morpheus {
    url      = var.morpheus_url
    username = var.morpheus_username
    password = var.morpheus_password
  }
}

variable "morpheus_url" { type = string }
variable "morpheus_username" { type = string }
variable "morpheus_password" {
  type      = string
  sensitive = true
}

variable "hvm_cloud_id" {
  type        = number
  description = "The HVM cloud whose cluster hosts the SDN controller."
}

variable "cluster_pool_id" {
  type        = number
  description = "The resource pool of the HVM cluster (the pool SDN networks and routers are placed on)."
}

variable "cn2_controller_address" {
  type        = string
  description = "The address of the Juniper CN2 controller."
}

variable "floating_ip_pool_id" {
  type        = number
  description = "A floating IP pool defined on the SDN network server, to be served by the logical router."
}

# 1. Register the CN2 controller as the cloud's network server. The controller
#    is reached from the cluster's hosts, so no service credentials are taken.
resource "hpe_morpheus_network_server" "sdn" {
  name      = "HPE Morpheus SDN"
  type_code = "morpheus-junipersdn-network"
  cloud_id  = var.hvm_cloud_id

  config_juniper_sdn = {
    controller_address = var.cn2_controller_address
    enable_sub_tenants = false
  }
}

# 2. An SDN virtual network on the cluster. The plugin edits nothing in place,
#    so any change replaces the network.
resource "hpe_morpheus_network" "web" {
  name              = "sdn-web"
  type_code         = "juniperSdn"
  cloud_id          = var.hvm_cloud_id
  group_id          = 1
  network_server_id = hpe_morpheus_network_server.sdn.id
  zone_pool_id      = var.cluster_pool_id
  cidr              = "10.90.0.0/24"
  gateway           = "10.90.0.1"
}

resource "hpe_morpheus_network" "app" {
  name              = "sdn-app"
  type_code         = "juniperSdn"
  cloud_id          = var.hvm_cloud_id
  group_id          = 1
  network_server_id = hpe_morpheus_network_server.sdn.id
  zone_pool_id      = var.cluster_pool_id
  cidr              = "10.91.0.0/24"
  gateway           = "10.91.0.1"
}

# 3. A logical router joining the two virtual networks and serving the floating
#    IP pool, so that the pool is available to allocate from in step 6.
resource "hpe_morpheus_network_router" "mesh" {
  name                   = "sdn-mesh"
  group_id               = 1
  network_integration_id = hpe_morpheus_network_server.sdn.id

  config_juniper_sdn = {
    resource_pool_id     = var.cluster_pool_id
    virtual_network_uids = [hpe_morpheus_network.web.unique_id, hpe_morpheus_network.app.unique_id]
    floating_ip_pool_ids = [var.floating_ip_pool_id]
  }
}

# 4. BGP peering between the SDN controller and the fabric.
resource "hpe_morpheus_network_router" "bgp" {
  name                   = "sdn-bgp"
  group_id               = 1
  network_integration_id = hpe_morpheus_network_server.sdn.id

  config_juniper_bgp = {
    vendor           = "juniper"
    asn              = 65000
    peer_address_1   = "10.0.0.1"
    address_families = ["inet"]
    hold_time        = 90
  }
}

# 5. A load balancer with a virtual IP on the web network, forwarding port 80
#    to 8080, with an HTTP health check.
#
#    CAUTION: on some appliances, deleting a load balancer has removed the pool
#    and virtual-server records of other load balancers on the same appliance
#    (MORPH-18174, observed with NSX-T; the suspected cause is not specific to a
#    load balancer type). Do not apply or destroy this example on a shared
#    appliance that hosts load balancers you care about until that is resolved.
resource "hpe_morpheus_load_balancer" "web" {
  name              = "sdn-web-lb"
  type_code         = "juniper-sdn-lb"
  cloud_id          = var.hvm_cloud_id
  network_server_id = hpe_morpheus_network_server.sdn.id

  config_juniper_sdn = {
    resource_pool_id    = var.cluster_pool_id
    virtual_network_uid = hpe_morpheus_network.web.unique_id
    virtual_ip          = "10.90.0.100"
    source_port         = 80
    destination_port    = 8080
    protocol            = "TCP"
  }
}

resource "hpe_morpheus_load_balancer_monitor" "web_health" {
  load_balancer_id = hpe_morpheus_load_balancer.web.id
  name             = "web-health"
  monitor_type     = "HTTP"
  monitor_interval = 5
  monitor_timeout  = 15
  max_retry        = 3

  config = {
    http_method      = "GET"
    http_uri         = "/health"
    http_port        = 8080
    http_return_code = "200"
  }
}

# 6. A floating IP from the pool the router serves. Pools are defined on the SDN
#    network server and discovered through the option source once a router of
#    the resource pool serves them; the lookup therefore waits for the router,
#    and the allocation checks the pool is listed before it runs.
data "hpe_morpheus_option_source_values" "floating_ip_pools" {
  source     = "juniperFloatingIpPools"
  parameters = { juniperResourcePool = tostring(var.cluster_pool_id) }

  depends_on = [hpe_morpheus_network_router.mesh]
}

resource "hpe_morpheus_network_floating_ip" "web" {
  network_server_id   = hpe_morpheus_network_server.sdn.id
  floating_ip_pool_id = var.floating_ip_pool_id

  lifecycle {
    precondition {
      condition     = contains([for v in data.hpe_morpheus_option_source_values.floating_ip_pools.values : tonumber(v.value)], var.floating_ip_pool_id)
      error_message = "Floating IP pool ${var.floating_ip_pool_id} is not served by a router on resource pool ${var.cluster_pool_id}."
    }
  }
}

output "web_floating_ip" {
  value = hpe_morpheus_network_floating_ip.web.ip_address
}
