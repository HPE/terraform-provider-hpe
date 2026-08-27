data "hpe_morpheus_resource_pool" "example" {
  name     = "morpheuspool"
  cloud_id = data.hpe_morpheus_cloud.vspherecloud.id
}
