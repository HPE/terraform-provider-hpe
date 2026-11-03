// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the Morpheus user. All five are Optional+Computed
// but carry the provider's Nullable*UpdateModifier, which plans the null
// itself whenever the configuration is null and the state is not — they are
// values only the practitioner sets, so there is no server-filled value to
// protect and no private-state record to keep. Update sends the null (the
// names and the key pair) or an empty string (the OS usernames, which the API
// rejects as null).
var clearableAttributes = clearable.Attributes{
	{Name: "first_name", Mode: clearable.Clear},
	{Name: "last_name", Mode: clearable.Clear},
	{Name: "linux_key_pair_id", Mode: clearable.Clear},
	{Name: "linux_username", Mode: clearable.Clear},
	{Name: "windows_username", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{ //nolint:gosec // attribute names, not credentials
	"receive_notifications":       "boolean",
	"tenant_id":                   "a change replaces the user; the current tenant is kept when it is removed",
	"password_wo_version":         "version companion of a write-only value",
	"linux_password_wo_version":   "version companion of a write-only value",
	"windows_password_wo_version": "version companion of a write-only value",
	"require_password_change":     "boolean; not read back, only acted on when the password is sent",
}
