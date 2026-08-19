---
page_title: "Authentication"
subcategory: "Morpheus"
---

# Morpheus authentication

There are four ways to authenticate with Morpheus:

1. Using a username and password
2. Using an access_token
3. Using PCE Identity, for a Connected PCE (Private Cloud Enterprise) deployment
4. Using PCE Disconnected Identity, for a Disconnected PCE deployment

For the first two the URL of the Morpheus instance must be provided as `url`.

For the last two the URL and access token are obtained from GreenLake, so `url` must not be set.
Provide a `pce_identity` or `pce_disconnected_identity` block instead, using either GreenLake API
client credentials (`client_id`, `client_secret` and `issuer_url`) or a pre-generated `iam_token`.
The two blocks are mutually exclusive, and neither can be combined with `url`, `username`,
`password`, `access_token` or `tenant_subdomain`.

Note that if authenticating with username and password a tenant can be specified by including the
`SUBDOMAIN` value for the tenant as `tenant_subdomain` in the provider block.

By default the provider will check the Morpheus server certificate and will fail if it is not
valid.  This can be toggled off by setting `insecure` to `true` in the provider block.

## Using a username and password

```terraform
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
    username = "username"
    password = "password"
    url      = "https://morpheus.example.com"
  }
}
```

## Using a username and password with tenant_subdomain

```terraform
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
    username         = "username"
    password         = "password"
    tenant_subdomain = "tenant"
    url              = "https://morpheus.example.com"
  }
}
```

## Using an access token

```terraform
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
    access_token = "access_token"
    url          = "https://morpheus.example.com"
  }
}
```

## Using an access token with insecure

```terraform
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
    access_token = "access_token"
    insecure     = true
    url          = "https://morpheus.example.com"
  }
}
```

## Using PCE Identity

```terraform
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
```

## Using PCE Identity with an IAM token

```terraform
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
      iam_token = "iam-token"
      location  = "location"
      space     = "space"
    }
  }
}
```

## Using PCE Disconnected Identity

```terraform
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
      client_id     = "client_id"
      client_secret = "client_secret"
      issuer_url    = "https://issuer.example.com"
      location      = "location"
      workspace_id  = "workspace_id"
      broker_url    = "https://broker.example.com"
    }
  }
}
```

## Using PCE Disconnected Identity with an IAM token

```terraform
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
      iam_token    = "iam-token"
      location     = "location"
      workspace_id = "workspace_id"
      broker_url   = "https://broker.example.com"
    }
  }
}
```
