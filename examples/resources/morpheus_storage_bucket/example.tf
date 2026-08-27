resource "hpe_morpheus_storage_bucket" "example" {
  name          = "Example Storage Bucket"
  provider_type = "s3"
  bucket_name   = "example-bucket"
}
