# broker

Go SDK for the **VMaaS broker** — the GreenLake Private Cloud VMaaS credential
and discovery broker (historically "VMaaS-CMP", where `CMP` was the term
previously used for Morpheus).

The broker trades an IAM token for Morpheus connection details (`cmp_details` →
Morpheus URL + access token). It is a credential broker, not a data-path proxy:
once the URL and token are returned, the provider talks to Morpheus directly.

The same client serves both PCE deployment types. Connected PCE uses the
HPE-hosted broker with a GLCS IAM token scoped by `?space=`; Disconnected PCE
uses the same software running on-premise with a GLP IAM token scoped by
`?tenantID=` and an `X-Tenant-ID` header.

Adapted from `github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk`.

## Layout

- `client/` — API client and the broker `cmp_details` exchange (`broker.go`).
- `models/` — request/response types for the broker API.
- `common/` — shared constants (API base path, the broker `cmp_details` path).
