// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute domain_password_wo: see
// wodigest. A change (or a change to the deprecated marker) marks
// write_only_revision unknown, which plans an in-place update in which Update
// re-sends the configured value. The deprecated *_version marker keeps working
// for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "network_domain",
	Attributes: map[string]string{
		"domain_password_wo": "domain_password_wo_version",
	},
}
