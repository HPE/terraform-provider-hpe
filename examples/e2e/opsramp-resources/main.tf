
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.6.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

# Create individual resources
resource "hpe_opsramp_resource" "resource1" {
  alias_name    = "MyFirstResource2"
  resource_name = "TestResource1"
  resource_type = "Other"
}

resource "hpe_opsramp_resource" "resource2" {
  alias_name    = "MySecondResource2"
  hostname      = "testresource.local"
  resource_type = "Other"
}