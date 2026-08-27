# Filter by name. Values are Go regular expressions; a security group matches
# the block if its name matches ANY value.
data "hpe_morpheus_security_groups" "by_name" {
  filter {
    name   = "name"
    values = ["^web-", "-prod$"]
  }
}

# Filter by visibility ("public" or "private").
data "hpe_morpheus_security_groups" "public" {
  filter {
    name   = "visibility"
    values = ["public"]
  }
}

# Filter by cloud (zone) id.
data "hpe_morpheus_security_groups" "by_cloud" {
  filter {
    name   = "cloud_id"
    values = ["^3$"]
  }
}

# Filter by active state.
data "hpe_morpheus_security_groups" "active_only" {
  filter {
    name   = "active"
    values = ["true"]
  }
}

# Multiple filter blocks are ANDed together. Sort the results by id descending.
data "hpe_morpheus_security_groups" "combined" {
  filter {
    name   = "name"
    values = ["web"]
  }

  filter {
    name   = "visibility"
    values = ["public"]
  }

  sort_ascending = false
}

# No filter blocks returns all security groups (up to 250).
data "hpe_morpheus_security_groups" "all" {}

# Consume the full objects returned in `security_groups`.
output "web_security_group_ids" {
  value = data.hpe_morpheus_security_groups.by_name.security_groups[*].id
}
