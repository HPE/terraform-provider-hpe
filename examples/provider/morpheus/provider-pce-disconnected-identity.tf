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
    pce_disconnected_identity {
      client_id        = "client_id"
      client_secret    = "client_secret"
      token_issuer_url = "https://issuer.example.com"
      location         = "location"
      workspace_id     = "workspace_id"
      broker_url       = "https://broker.example.com"
    }
  }
}
