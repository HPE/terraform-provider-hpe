data "hpe_morpheus_cloud" "example" {
  name = "hvm"
}

data "hpe_morpheus_group" "example" {
  name = "Zodiac"
}

resource "hpe_morpheus_load_balancer" "haproxy_generic" {
  name        = "example-terraform-haproxy-lb"
  description = "HAProxy load balancer via generic config"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  group_id    = data.hpe_morpheus_group.example.id
  visibility  = "public"
  type_code   = "haproxyContainer"

  config = {
    plan = {
      id = 8
    }
    pool = {
      id = "pool-574"
    }
  }
}
