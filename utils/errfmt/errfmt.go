// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package errfmt

import "fmt"

// ChildProviderSchemaErr returns a diagnostic-ready error message when a child
// provider does not expose the expected schema block keyed by its type name.
func ChildProviderSchemaErr(providerName string) (summary, detail string) {
	summary = "Invalid child provider schema"
	detail = fmt.Sprintf(
		"Child provider %q did not return a schema block keyed by its type name. "+
			"Ensure the child provider is wrapped with adapter.NewAdaptedProvider() "+
			"before passing it to the parent provider.",
		providerName,
	)

	return summary, detail
}
