resource "hpe_morpheus_library_option_type_list" "example" {
  name        = "Region List"
  description = "List of available regions"
  type        = "rest"
  source_url  = "https://api.example.com/regions"
  visibility  = "public"
  real_time   = false
}
