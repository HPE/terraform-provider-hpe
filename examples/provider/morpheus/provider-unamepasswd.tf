# Copyright 2025 Hewlett Packard Enterprise Development LP

terraform {
  required_providers {
    hpegl = {
      source  = "HPE/hpe"
      version = "= 0.0.1"
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
