resource "hpe_morpheus_scale_threshold" "example" {
  name           = "Web Tier Scaling"
  auto_upscale   = true
  auto_downscale = true
  min_count      = 1
  max_count      = 5
  cpu_enabled    = true
  min_cpu        = 20.0
  max_cpu        = 80.0
}
