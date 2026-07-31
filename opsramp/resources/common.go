// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import "github.com/hashicorp/terraform-plugin-framework/types"

// setToStringSlice converts a types.Set of strings to a []string.
func setToStringSlice(s types.Set) []string {
	elements := s.Elements()
	result := make([]string, 0, len(elements))
	for _, e := range elements {
		if sv, ok := e.(types.String); ok {
			result = append(result, sv.ValueString())
		}
	}

	return result
}

// stringSetDiff returns elements in a that are not in b.
func stringSetDiff(a, b []string) []string {
	bSet := make(map[string]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}
	var diff []string
	for _, v := range a {
		if _, found := bSet[v]; !found {
			diff = append(diff, v)
		}
	}

	return diff
}
