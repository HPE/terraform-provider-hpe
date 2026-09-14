// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package constants

import "time"

// Put this here, we get an import cycle in `configure` package if we store the
// `ProviderName` const in the `morpheus` package.
const ProviderName = "morpheus"

// TODO: properly implement resource timeouts similar to
// terraform-plugin-framework-timeouts
const NetworkDeleteTimeout = 5 * time.Minute

// AffinityGroupMinVersion is the minimum Morpheus appliance version required by
// the affinity group resources and data sources. It is declared here, rather
// than in any one of the six affinity group packages, so that all of them gate
// on exactly the same value.
//
// The affinity group API arrived in 8.0.8, but its semantics moved twice after
// that: 8.0.9 made `pool` required, and 8.0.10 dropped the constraint that a
// name be unique per refType+refId+pool. 8.0.10 is therefore the first release
// whose behaviour the provider can rely on, so that — not 8.0.8 — is the gate.
// (MORPH-15506.)
const AffinityGroupMinVersion = ">= 8.0.10"

// TenantParentMinVersion is the minimum Morpheus appliance version on which
// the tenant resource's parent_id attribute is honoured.
//
// Nominating a parent tenant on create (parentAccount.id) and the parent
// object in the tenant response both arrived with the tenant-hierarchy feature
// ("Tenants of Tenants") in 8.1.0 and were never backported to 8.0.x. Before
// 8.1.0 the API silently ignores parentAccount and returns no parent, so a
// configured parent_id can never be reflected in state: the apply fails with
// an inconsistent-result error and every subsequent plan forces a replacement.
// The gate turns that loop into a single, named diagnostic.
const TenantParentMinVersion = ">= 8.1.0"
