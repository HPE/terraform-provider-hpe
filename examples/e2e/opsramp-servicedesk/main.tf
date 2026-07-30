
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

# Create multiple categories
resource "hpe_opsramp_servicedesk_category" "category1" {
  name        = "Category1"
  description = "Category1 Description"
  ticket_type = "serviceRequests"
}

resource "hpe_opsramp_servicedesk_category" "category2" {
  name        = "Category2"
  description = "Category2 Description"
  ticket_type = "incidents"
}

resource "hpe_opsramp_servicedesk_category" "category3" {
  name        = "Category3"
  description = "Category3 Description"
  ticket_type = "problems"
}


# Create  multiple business impacts
resource "hpe_opsramp_servicedesk_business_impact" "business_impact1" {
  name        = "Business Impact1"
  description = "Business Impact1 Description"
}

resource "hpe_opsramp_servicedesk_business_impact" "business_impact2" {
  name        = "Business Impact2"
  description = "Business Impact2 Description"
}

resource "hpe_opsramp_servicedesk_business_impact" "business_impact3" {
  name        = "Business Impact3"
  description = "Business Impact3 Description"
}

# # # Create multiple urgencies
resource "hpe_opsramp_servicedesk_urgency" "urgency1" {
  name        = "Urgency1"
  description = "Urgency1 Description"
}

resource "hpe_opsramp_servicedesk_urgency" "urgency2" {
  name        = "Urgency2"
  description = "Urgency2 Description"
}

resource "hpe_opsramp_servicedesk_urgency" "urgency3" {
  name        = "Urgency3"
  description = "Urgency3 Description"
}