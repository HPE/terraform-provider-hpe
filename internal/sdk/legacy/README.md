# Legacy go-morpheus-sdk

This module contains a modified version of the legacy Morpheus Go SDK.

At a high level, changes will include using the Go standard library HTTP client as a backend rather than Resty, as well as general bugfixes for the SDK.

### WARNING: FOR INTERNAL USE WITH MORPHEUS TERRAFORM PROVIDER ONLY

Do NOT use as general purpose clients for interacting with Morpheus.

These clients are intended to be used in the [HPE Morpheus Terraform Provider](https://github.com/HPE/terraform-provider-hpe).

Not all endpoints have been tested as development of the Terraform provider is ongoing.
