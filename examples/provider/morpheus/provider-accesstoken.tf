# Copyright 2025-2026 Hewlett Packard Enterprise Development LP

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = "= 1.0.0"
    }
  }
}

provider "hpe" {
  # Provide morpheus block if you want to create morpheus resources
  morpheus {
    access_token = "access_token"
    url          = "https://morpheus.example.com"
  }
}
