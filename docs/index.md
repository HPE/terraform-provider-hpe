---
layout: ""
page_title: "hpe Provider"
subcategory: ""
description: |-
  This is the hpe terraform provider
---

# hpe Provider

This is the hpe terraform provider which is still in development.  It will eventually replace the
[HPE GreenLake Terraform Provider](https://registry.terraform.io/providers/HPE/hpegl/latest) and the
[Morpheus Terraform Provider](https://registry.terraform.io/providers/gomorpheus/morpheus/latest).

Initially this provider will support Morpheus, but will in time expand to cover other HPE offerings.

This provider requires 64-bit versions of the terraform binary to work properly.

->In some circumstances users may need to use this provider and the [Morpheus Terraform Provider](https://registry.terraform.io/providers/gomorpheus/morpheus/latest)
in the same configuration.  To ensure that Morpheus API SSL cert checking is consistent across both providers
the `MORPHEUS_API_SECURE` environment variable (set to true) can be used to enable SSL cert checking in the Morpheus provider.
This environment variable is supported by versions 0.14.0 and later of the Morpheus provider.  Note that
Morpheus API SSL cert checking is enabled by default in this provider.

## Morpheus

This provider can be used to manage Morpheus resources.  Support will grow over time.  See below for
release notes for the current version (v0.3.0).

### Authentication

There are two ways to authenticate with Morpheus:
1. Using a username and password
2. Using an access_token

With either method the URL of the Morpheus instance must be provided as `url`.

By default the provider will check the Morpheus server certificate and will fail if it is not valid.  This can be
be toggled off by setting `insecure` to `true` in the provider block.

### Example Usage

#### Using a username and password

```terraform
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
```

#### Using an access token

```terraform
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
    access_token = "access_token"
    url          = "https://morpheus.example.com"
  }
}
```

#### Using an access token with insecure

```terraform
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
    access_token = "access_token"
    insecure     = true
    url          = "https://morpheus.example.com"
  }
}
```

## Release Notes

->The `hpe_morpheus_user` resource uses a `WriteOnly` password field.  `WriteOnly` attributes are supported
by Terraform versions 1.11 and later.

->The `hpe_morpheus_instance`, `hpe_morpheus_network`, `hpe_morpheus_datastore` and `hpe_morpheus_cloud` resources use
a `Dynamic` attribute for `config`.  This means that the `config` block can contain arbitrary nested attributes which
will be evaluated at run-time.  Examples of these are shown in the documentation.

### New functionality

In this release (v0.3.0) we have added the following resource functionality:
- hpe_morpheus_image resource has been added (Create, Delete, Read - no Update)
- hpe_morpheus_policy resource has been added (Create, Delete, Read, Update)
- hpe_morpheus_instance Update functionality has been added
- hpe_morpheus_service_plan `cores_per_socket` is now required
- hpe_morpheus_datastore import will now populate `resource_permissions` (`groups` only) and `tenants`

In this release (v0.3.0) we have added the following data-source functionality:
- hpe_morpheus_datastore data-source has been added

### New known issues

- hpe_morpheus_datastore data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.  This is an API bug which is being investigated.

### Known issues from previous releases

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- hpe_morpheus_datastore when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- hpe_morpheus_datastore delete is not guaranteed to succeed.  AlletraMP HVM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- hpe_morpheus_instance only supports 1 network
- hpe_morpheus_instance in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `morpheus` (Block List, Max: 1) (see [below for nested schema](#nestedblock--morpheus))

<a id="nestedblock--morpheus"></a>
### Nested Schema for `morpheus`

Required:

- `url` (String) Morpheus instance URL

Optional:

- `access_token` (String, Sensitive) Morpheus access token for authentication
- `insecure` (Boolean) Explicitly allow the provider to perform "insecure" SSL requests. If omitted, default value is `false`
- `password` (String, Sensitive) Morpheus password for authentication, required if username is set
- `username` (String) Morpheus username for authentication, required if password is set
