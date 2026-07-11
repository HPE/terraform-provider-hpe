# Create a resource using resource_name and resource_type
resource "hpe_opsramp_resource" "test_resource" {
  alias_name    = "TestResourceAlias"
  resource_name = "TestResource"
  resource_type = "Server"
}