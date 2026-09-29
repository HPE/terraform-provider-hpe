// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the storage server: the description, service
// host, URL and username (Optional-only; Update sends null) and the tenants
// (Optional+Computed, cleared once an apply has recorded them; Update sends
// an empty list).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "service_host", Mode: clearable.Clear},
	{Name: "service_url", Mode: clearable.Clear},
	{Name: "service_username", Mode: clearable.Clear},
	{Name: "tenants", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{ //nolint:gosec // attribute names and reasons, not credentials
	"timeouts":                    "nested block; removing a whole block is not yet supported",
	"credential_id":               "a credential reference; the current credential is kept when it is removed",
	"enabled":                     "has a default; Terraform plans the default when it is removed",
	"service_password_wo":         "a write-only value; it is sent only when set",
	"service_password_wo_version": "version companion of a write-only value",
	"visibility":                  "has a default; Terraform plans the default when it is removed",
}
