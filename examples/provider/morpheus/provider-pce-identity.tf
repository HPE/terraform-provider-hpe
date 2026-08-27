# Copyright 2025-2026 Hewlett Packard Enterprise Development LP

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.0.0"
    }
  }
}

provider "hpe" {
  # Provide morpheus block if you want to create morpheus resources
  morpheus {
    pce_identity {
      client_id     = "client_id"
      client_secret = "client_secret"
      issuer_url    = "https://issuer.example.com"
      location      = "location"
      space         = "space"
    }
  }
}
