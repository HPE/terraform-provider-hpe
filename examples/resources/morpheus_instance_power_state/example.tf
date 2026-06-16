resource "hpe_morpheus_instance_power_state" "example" {
  instance_id   = 1
  desired_state = "running"
}
