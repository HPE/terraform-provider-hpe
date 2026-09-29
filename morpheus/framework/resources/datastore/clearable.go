// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastore

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the datastore. There are none.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"active":                 "boolean",
	"config_alletramp_bmaas": "set when the datastore is created; a change replaces it, and the current value is kept when it is removed",
	"config_alletramp_hvm":   "set when the datastore is created; a change replaces it, and the current value is kept when it is removed",
	"config_nfs":             "set when the datastore is created; a change replaces it, and the current value is kept when it is removed",
	"datastores":             "a nested configuration block, sent as configured and not cleared as a whole",
	"default_store":          "boolean",
	"resource_permissions":   "a nested configuration block, sent as configured and not cleared as a whole",
	"resource_pool":          "a nested configuration block, sent as configured and not cleared as a whole",
	"storage_server":         "a nested configuration block, sent as configured and not cleared as a whole",
	"tenants":                "a nested configuration block, sent as configured and not cleared as a whole",
	"config":                 "set when the datastore is created; a change replaces it, and the current value is kept when it is removed",
	"visibility":             "Morpheus requires a visibility; the current one is kept when it is removed",
}
