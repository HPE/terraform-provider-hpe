// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import (
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute ssh_password_wo: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value. The deprecated
// *_version marker keeps working for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "image",
	Attributes: map[string]string{
		"ssh_password_wo": "ssh_password_wo_version",
	},
}
