# Copyright 2025 Hewlett Packard Enterprise Development LP

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = "= 0.3.0"
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
