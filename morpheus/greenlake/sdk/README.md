# morpheus/greenlake/sdk

Common SDK packages shared across GreenLake deployment types
(e.g. `greenlake_connected`, and on-premise in future).

Code here is deployment-agnostic: it implements the low-level GreenLake platform
integrations — authentication, token exchange, and related clients — that more
than one deployment type may need. Deployment-specific wiring (schema handling,
`Configure`, token exchange) lives with the Morpheus provider, not here.

## Layout

- `token/` — GreenLake IAM token generation and automatic refresh. The
  client is IAM-version-agnostic (callers pass the version at call time), so it
  serves both GLCS and GLP token exchanges.

## Guidelines

- Keep packages here free of deployment-specific assumptions so any GreenLake
  deployment type can depend on them.
- Prefer adding shared, reusable integrations here rather than duplicating them
  inside individual deployment types.
