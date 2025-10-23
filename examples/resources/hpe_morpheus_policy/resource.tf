# Configure the HPE provider
terraform {
  required_providers {
    hpe = {
      source = "HPE/hpe"
    }
  }
}

# Create a Morpheus policy with backup creation configuration
resource "hpe_morpheus_policy" "backup_policy" {
  name        = "backup-creation-policy"
  description = "Policy to enforce backup creation for instances"
  enabled     = true

  # Required empty config field
  config = {}

  config_backup_creation = {
    create_backup      = true
    create_backup_type = "snapshot"
  }
}

# Create a Morpheus policy with budget limits
resource "hpe_morpheus_policy" "budget_policy" {
  name        = "monthly-budget-policy"
  description = "Policy to enforce monthly budget limits"
  enabled     = true

  # Required empty config field
  config = {}

  config_budget = {
    max_price          = 1000.0
    max_price_currency = "USD"
    max_price_unit     = "month"
  }
}

# Create a Morpheus policy with resource limits
resource "hpe_morpheus_policy" "resource_limits_policy" {
  name        = "resource-limits-policy"
  description = "Policy to enforce resource limits"
  enabled     = true

  # Required empty config field
  config = {}

  config_max_memory = {
    max_memory = {
      anyof1 = 8589934592  # 8GB in bytes
    }
  }

  config_max_cores = {
    max_cores = 8
  }

  config_max_storage = {
    max_storage = 107374182400  # 100GB in bytes
  }
}

# Create a Morpheus policy for instance naming
resource "hpe_morpheus_policy" "naming_policy" {
  name        = "instance-naming-policy"
  description = "Policy to enforce instance naming patterns"
  enabled     = true

  # Required empty config field
  config = {}

  config_instance_name = {
    naming_type    = "pattern"
    naming_pattern = "vm-$${sequence}-$${cloudCode}"
  }
}

# Create a Morpheus policy for tags enforcement
resource "hpe_morpheus_policy" "tags_policy" {
  name        = "tags-enforcement-policy"
  description = "Policy to enforce required tags"
  enabled     = true

  # Required empty config field
  config = {}

  config_tags = {
    key    = "Environment"
    value  = "Production,Development,Testing"
    strict = true
  }
}