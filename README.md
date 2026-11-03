<!-- Mirrors templates/index.md.tmpl above its "## Release Notes" section; keep them in step. -->

# HPE Provider

This is the HPE Terraform provider.  It replaces the
[HPE GreenLake Terraform Provider](https://registry.terraform.io/providers/HPE/hpegl/latest), for its
VMaaS resources and data sources, and the
[Morpheus Terraform Provider](https://registry.terraform.io/providers/gomorpheus/morpheus/latest).

This provider manages Morpheus and OpsRamp resources, configured with `morpheus` and `opsramp`
blocks, and will in time expand to cover other HPE offerings.  GreenLake VMaaS is covered by the
Morpheus resources, so there is no separate `hpegl` block.

This provider requires 64-bit versions of the Terraform binary to work properly.

## Support

**Support policy** — Bug fixes and new features are shipped in the next release of the provider. We
do not backport changes to previous releases; if you encounter a problem, please upgrade to the
latest version before opening an issue. Security vulnerabilities may be backported to the final
release of the previous major version at the maintainers' discretion, on a case-by-case basis. The
same policy applies to the migration tool (tfmigrator).

**Support note** — HPE Morpheus functionality in this provider is fully supported by HPE subject to
the support policy above. Support for HPE OpsRamp functionality is currently provided on a
best-efforts basis.

As of v2.0.0 this provider includes all functionality from the
[Morpheus Terraform Provider](https://registry.terraform.io/providers/gomorpheus/morpheus/latest),
and closes the remaining parity gaps with the `hpegl` VMaaS resources and data sources, so an `hpegl`
VMaaS configuration has a direct translation here.  Users of either provider are encouraged to
migrate; the Morpheus provider will be deprecated in the future.

We have developed tooling [tfmigrator](docs/guides/tfmigrator_migration.md) to assist with
migration, which supports `hpegl` configurations as well as Morpheus ones.  Please see the
[documentation](docs/guides/tfmigrator_migration.md) for more information on how to use this tool.

See the [Morpheus to HPE](docs/guides/morpheus_to_hpe_migration.md) guide for the mapping of
Morpheus provider resources and data sources to those in this provider.  The data sources that
replace an `hpegl` lookup each carry a "Migrating from hpegl" section showing the two configurations
side by side.

## Morpheus

This provider can be used to manage Morpheus resources.  Support will grow over time.  Release notes
for the current version are in [CHANGELOG.md](CHANGELOG.md).

### Authentication

There are four ways to authenticate with Morpheus: a username and password, an access token, PCE
Identity for a Connected PCE (Private Cloud Enterprise) deployment, and PCE Disconnected Identity
for a Disconnected deployment.  See the [Authentication](docs/guides/morpheus_authentication.md)
guide for the details of each, with a worked example of every supported provider configuration.

### Example Usage

```hcl
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.2.0"
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

## OpsRamp

This provider can also be used to manage OpsRamp resources.  See the
[OpsRamp to HPE](docs/guides/opsramp_to_hpe_migration.md) guide for moving an existing OpsRamp
provider configuration across.

### Authentication

OpsRamp resources are configured with an `opsramp` block, authenticating against the OpsRamp API
with OAuth client credentials scoped to a tenant.  Each value may instead be supplied through the
environment.  See the [Authentication](docs/guides/opsramp_authentication.md) guide for the details.

### Example Usage

```hcl
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.2.0"
    }
  }
}

provider "hpe" {
  # Provide opsramp block if you want to create opsramp resources
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}
```

## Documentation

- Provider documentation on the
  [Terraform Registry](https://registry.terraform.io/providers/HPE/hpe/latest/docs), generated from
  the `docs/` directory of this repository.
- Release notes for the current version in [CHANGELOG.md](CHANGELOG.md); earlier versions in
  [HISTORY.md](HISTORY.md).
- Manual installation, for environments without access to the registry, in [INSTALL.md](INSTALL.md).
