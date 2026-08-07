# broker

Go SDK for the **PCE broker**, the service that hands out Morpheus connection
details for a PCE deployment.

The broker trades a GreenLake IAM token for a Morpheus URL and access token. It
is a credential and discovery broker, not a data-path proxy: once the URL and
token are returned, the provider talks to Morpheus directly.

The same client serves both PCE deployment types. Connected PCE uses the
HPE-hosted broker with a GLCS IAM token scoped by `?space=`; Disconnected PCE
uses an operator-supplied broker URL with a GLP IAM token scoped by `?tenantID=`
and an `X-Tenant-ID` header.

Adapted from `github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk`.

## Layout

- `client/` — API client and the broker exchange (`broker.go`).
- `models/` — request/response types for the broker API.
- `common/` — shared constants (API base path, the broker exchange path).

Identifiers such as `GetCMPDetails` and `CMPDetails` keep their upstream names,
and the exchange path is the one the broker serves; `CMP` was the term
previously used for Morpheus.
