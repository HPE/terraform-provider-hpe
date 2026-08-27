resource "hpe_morpheus_option_list" "example" {
  name        = "Region List"
  description = "List of available regions"
  type        = "manual"
  visibility  = "public"
  real_time   = false
}
