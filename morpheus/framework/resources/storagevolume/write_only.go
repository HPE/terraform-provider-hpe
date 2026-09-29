// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"

// Change detection for the write-only attributes config, config_alletramp_bmaas: see
// wodigest. Morpheus takes the configuration only at creation, so a change plans a replacement --
// the only way the new value can take effect -- rather than an update that would silently do
// nothing. The deprecated *_version markers keep working for the
// whole 2.2 line.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "storage_volume",
	Attributes: map[string]string{
		"config":                 "config_wo_version",
		"config_alletramp_bmaas": "config_alletramp_bmaas_wo_version",
	},
	Immutable: []string{"config", "config_alletramp_bmaas"},
}
