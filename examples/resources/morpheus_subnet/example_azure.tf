resource "hpe_morpheus_subnet" "azure_example" {
  type_id    = 12
  visibility = "private"

  config = {
    subnetName = "my-azure-subnet"
    subnetCidr = "10.0.1.0/24"
  }
  config_version = 1
}
