# Copyright 2025-2026 Hewlett Packard Enterprise Development LP

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.4.0"
    }
  }
}

provider "hpe" {
  # Provide morpheus block if you want to create morpheus resources
  morpheus {
    username = "username"
    password = "password"
    url      = "https://morpheus.example.com"
  }
}
