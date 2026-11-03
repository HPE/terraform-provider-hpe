// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the bucket: the retention days, Optional-only,
// so the plan carries the null and Update sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "retention_days", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"access_key":            "write-only",
	"access_key_version":    "version companion of a write-only value",
	"bucket_name":           "Morpheus gives it a default for the provider type; the current name is kept when it is removed",
	"default_backup_target": "boolean",
	"endpoint":              "part of the provider configuration; the current endpoint is kept when it is removed",
	"secret_key":            "write-only",
	"secret_key_version":    "version companion of a write-only value",
}
