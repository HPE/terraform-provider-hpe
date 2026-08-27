resource "hpe_morpheus_instance_snapshot" "example" {
  instance_id      = 1
  name             = "baseline"
  description      = "Golden baseline snapshot"
  memory_snapshot  = false
  retain_on_delete = false
}
