// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"

// Change detection for the nested write-only attributes
// config_vmware.password and config_azure.client_secret: see wodigest. A change
// marks write_only_revision unknown, which plans an in-place update in which
// Update re-sends the configured value. The deprecated *_version markers keep
// working for the whole 2.2 line. ModifyPlan is in modify_plan.go.
var wo = wodigest.Tracker{
	Resource: "cloud",
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"config_vmware.password":     "config_vmware.password_version",
		"config_azure.client_secret": "config_azure.client_secret_version",
	},
}
