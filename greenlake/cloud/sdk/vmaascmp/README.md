# vmaascmp

Go SDK for the **VMaaS-CMP API** — the GreenLake Private Cloud VMaaS "Cloud
Management Platform" (Morpheus) API surface.

Within the `greenlake_cloud` auth exchange this SDK provides the **VMaaS broker**
client, which trades a GLCS IAM token for Morpheus connection details
(`cmp_details` → Morpheus URL + access token). It also carries the broader
VMaaS-CMP resource client and models (instances, networks, load balancers, …).

`CMP` (Cloud Management Platform) is the term previously used for Morpheus.

Adapted from `github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk`.

## Layout

- `client/` — API client and per-resource service methods, including the broker
  `cmp_details` exchange (`broker.go`).
- `models/` — request/response types for the VMaaS-CMP API.
- `common/` — shared constants (API base paths, the broker `cmp_details` path).
