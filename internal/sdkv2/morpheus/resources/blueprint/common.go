// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
//
//nolint:unused
func matchTemplatesWithSchema(templates []int64, declaredTemplates []any) []int64 {
	result := make([]int64, len(declaredTemplates))

	rMap := make(map[int64]int64, len(templates))
	for _, template := range templates {
		rMap[template] = template
	}

	for i, definedTemplate := range declaredTemplates {
		// skip if type assertion failed
		if definedTemplateInt64, ok := definedTemplate.(int64); ok {
			if v, ok := rMap[definedTemplateInt64]; ok {
				// matched node type declared by ID
				result[i] = v
				delete(rMap, v)
			}
		}
	}
	// append unmatched node type to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}

	return result
}
