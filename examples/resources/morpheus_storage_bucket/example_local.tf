# A bucket on the appliance's own disk: the first bucket a new appliance needs,
# since it ships with none and an image fetched by URL has to land somewhere.
# Local storage is available to the master tenant only; Morpheus names the
# bucket "." itself, so bucket_name is left out.
resource "hpe_morpheus_storage_bucket" "images" {
  name          = "images"
  provider_type = "local"
  base_path     = "/var/opt/morpheus/morpheus-ui/images"

  # Left out, the bucket keeps whatever it has; set it to make this bucket the
  # default image store, or to false to take the default away.
  default_virtual_image_target = true
}
