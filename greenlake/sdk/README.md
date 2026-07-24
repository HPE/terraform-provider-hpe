# greenlake/sdk

Common SDK packages shared across the GreenLake Terraform providers
(e.g. `greenlake_cloud`).

Code here is provider-agnostic: it implements the low-level GreenLake platform
integrations — authentication, token exchange, and related clients — that more
than one child provider may potentially need. Provider-specific wiring (schemas, `Configure`,
resources, data sources) lives with each provider, not here.

## Layout

- `token/` — GreenLake IAM token generation and automatic refresh. The
  client is IAM-version-agnostic (callers pass the version at call time), so it
  serves both GLCS and GLP token exchanges.

## Guidelines

- Keep packages here free of provider-specific assumptions so any GreenLake
  provider can depend on them.
- Prefer adding shared, reusable integrations here rather than duplicating them
  inside individual providers.
