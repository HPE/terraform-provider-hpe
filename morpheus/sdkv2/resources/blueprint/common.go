// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
//
//nolint:unused
func matchTemplatesWithSchema(templates []int64, declaredTemplates []any) []int64 {
	// templates are the IDs returned by the API; declaredTemplates are the IDs
	// from the practitioner's config/state. SDKv2 stores schema.TypeInt list
	// elements as Go int, so accept both int and int64. The previous int64-only
	// assertion silently failed for every element, which left zero-valued entries
	// (from make([]int64, len(declared))) and then appended the entire API list —
	// doubling the list and injecting phantom 0 IDs (a perpetual plan diff).
	present := make(map[int64]bool, len(templates))
	for _, t := range templates {
		present[t] = true
	}

	result := make([]int64, 0, len(templates))
	// Keep the declared order for templates the API still returns.
	for _, dt := range declaredTemplates {
		var id int64
		switch v := dt.(type) {
		case int:
			id = int64(v)
		case int64:
			id = v
		default:
			continue
		}
		if present[id] {
			result = append(result, id)
			delete(present, id)
		}
	}
	// Append any server-added templates that were not declared, in the API's
	// order (ranging the map directly would reorder the list on every read).
	for _, t := range templates {
		if present[t] {
			result = append(result, t)
			delete(present, t)
		}
	}

	return result
}
