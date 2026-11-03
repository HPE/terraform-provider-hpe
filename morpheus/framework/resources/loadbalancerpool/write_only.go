// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only partition: see wodigest. A change marks
// write_only_revision unknown, which plans an in-place update in which Update
// sends the configured value. The resource has no marker attribute.
var wo = wodigest.Tracker{
	Resource:   "load_balancer_pool",
	Attributes: map[string]string{"partition": ""},
}
