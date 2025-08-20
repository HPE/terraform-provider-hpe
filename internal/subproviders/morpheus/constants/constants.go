// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package constants

import "time"

const SubProviderName = "morpheus"

// TODO: properly implement resource timeouts similar to
// terraform-plugin-framework-timeouts
const NetworkDeleteTimeout = 5 * time.Minute
