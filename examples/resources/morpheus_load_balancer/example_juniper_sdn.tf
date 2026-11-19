resource "hpe_morpheus_load_balancer" "sdn" {
  name              = "example-terraform-sdn-lb"
  type_code         = "juniper-sdn-lb"
  visibility        = "public"
  cloud_id          = 1
  network_server_id = 5

  config_juniper_sdn = {
    resource_pool_id    = 1
    virtual_network_uid = "0be6b79b-757c-4d92-affc-8084e328cb4c"
    virtual_ip          = "10.90.0.100"
    source_port         = 80
    destination_port    = 8080
    protocol            = "TCP"
  }
}
