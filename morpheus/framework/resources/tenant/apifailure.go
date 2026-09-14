// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"fmt"
	"sort"
	"strings"
)

// apiFailure reports whether a 2xx tenant response nonetheless carried
// success:false and, if so, assembles a human-readable detail from the
// envelope's msg and errors.
//
// Morpheus appliances before 8.1.0 render some tenant validation failures —
// "Invalid role id", "Missing required object: account" — with
// `render rtn as JSON` and no status code: HTTP 200, a body of
// {success:false, msg, errors} and no account object. 8.1.0 and later return
// HTTP 400 for the same cases, which errfmt.CheckResponse already turns into an
// error. Checking the envelope keeps the older appliances' real message rather
// than the misleading "Account ID is nil" a missing account would otherwise
// produce.
//
// success is the typed field where the SDK model has one (AddTenant); when it
// is nil the envelope's "success" key is read from props instead
// (UpdateTenant, whose model routes it to AdditionalProperties). msg and errors
// are never typed and always come from props.
func apiFailure(success *bool, props map[string]any) (string, bool) {
	failed := success != nil && !*success
	if success == nil {
		if v, ok := props["success"].(bool); ok && !v {
			failed = true
		}
	}

	if !failed {
		return "", false
	}

	var parts []string

	if msg, ok := props["msg"].(string); ok && msg != "" {
		parts = append(parts, msg)
	}

	if errs, ok := props["errors"].(map[string]any); ok && len(errs) > 0 {
		keys := make([]string, 0, len(errs))
		for k := range errs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s: %v", k, errs[k]))
		}
	}

	if len(parts) == 0 {
		return "the API reported success=false without a message", true
	}

	return strings.Join(parts, "; "), true
}
