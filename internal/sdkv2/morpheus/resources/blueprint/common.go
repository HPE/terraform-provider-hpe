package blueprint

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
//
//nolint:unused
func matchTemplatesWithSchema(templates []int64, declaredTemplates []interface{}) []int64 {
	result := make([]int64, len(declaredTemplates))

	rMap := make(map[int64]int64, len(templates))
	for _, template := range templates {
		rMap[template] = template
	}

	for i, definedTemplate := range declaredTemplates {
		definedTemplate := int64(definedTemplate.(int))

		if v, ok := rMap[definedTemplate]; ok {
			// matched node type declared by ID
			result[i] = v
			delete(rMap, v)
		}
	}
	// append unmatched node type to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}

	return result
}
