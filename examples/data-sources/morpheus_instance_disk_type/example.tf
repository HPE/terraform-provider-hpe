data "hpe_morpheus_cloud" "example" {
  name = "hvm"
}

data "hpe_morpheus_instance_disk_type" "example" {
  name      = "Standard"
  cloud_id  = data.hpe_morpheus_cloud.example.id
  layout_id = 77
  group_id  = 1
}
