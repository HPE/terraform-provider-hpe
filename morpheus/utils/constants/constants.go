// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package constants

import "time"

// Put this here, we get an import cycle in `configure` package if we store the
// `ProviderName` const in the `morpheus` package.
const ProviderName = "morpheus"

// TODO: properly implement resource timeouts similar to
// terraform-plugin-framework-timeouts
const NetworkDeleteTimeout = 5 * time.Minute
