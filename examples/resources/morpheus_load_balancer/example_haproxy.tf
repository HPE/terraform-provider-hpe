data "hpe_morpheus_cloud" "example" {
  name = "hvm"
}

data "hpe_morpheus_group" "example" {
  name = "Zodiac"
}

resource "hpe_morpheus_load_balancer" "haproxy" {
  name        = "example-terraform-haproxy-lb"
  description = "HAProxy load balancer"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  group_id    = data.hpe_morpheus_group.example.id
  visibility  = "public"

  config_haproxy = {
    plan_id = 8
    pool    = "pool-574"
  }
}
