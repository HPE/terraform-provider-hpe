// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// Flatten converts a nested map into a flat map with dot-notation keys.
// e.g., {"zone": {"name": "foo", "config": {"apiUrl": "bar"}}}
// becomes {"zone.name": "foo", "zone.config.apiUrl": "bar"}
func Flatten(nested map[string]any) map[string]any {
	flat := make(map[string]any)
	flattenRecursive("", nested, flat)

	return flat
}

func flattenRecursive(prefix string, m map[string]any, flat map[string]any) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]any:
			flattenRecursive(key, val, flat)
		default:
			flat[key] = val
		}
	}
}

// Unflatten converts a flat dot-notation map back into a nested map.
// e.g., {"zone.name": "foo", "zone.config.apiUrl": "bar"}
// becomes {"zone": {"name": "foo", "config": {"apiUrl": "bar"}}}
//
// Special handling: segments matching "anyofN" or "oneofN" are transparent
// (skipped during nesting) since they are OpenAPI codegen artifacts, not real JSON keys.
func Unflatten(flat map[string]any) map[string]any {
	nested := make(map[string]any)

	for key, val := range flat {
		parts := strings.Split(key, ".")
		// Filter out anyofN/oneofN segments
		filteredParts := make([]string, 0, len(parts))
		for _, part := range parts {
			if isCodegenArtifact(part) {
				continue
			}

			filteredParts = append(filteredParts, part)
		}

		if len(filteredParts) == 0 {
			continue
		}

		current := nested

		for i, part := range filteredParts {
			if i == len(filteredParts)-1 {
				current[part] = val
			} else {
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]any)
				}

				next, ok := current[part].(map[string]any)
				if !ok {
					// Conflict: a leaf value exists where we need a branch.
					// Create a new map and overwrite.
					next = make(map[string]any)
					current[part] = next
				}

				current = next
			}
		}
	}

	return nested
}

// isCodegenArtifact returns true for path segments that are OpenAPI codegen
// artifacts (anyof0, anyof1, oneof0, etc.) and should not appear in actual JSON.
func isCodegenArtifact(s string) bool {
	if len(s) < 6 {
		return false
	}

	if strings.HasPrefix(s, "anyof") || strings.HasPrefix(s, "oneof") {
		// Check that the rest is digits
		suffix := s[5:]
		for _, r := range suffix {
			if r < '0' || r > '9' {
				return false
			}
		}

		return true
	}

	return false
}

// TransformForRead applies forward moves to an API response flat map,
// converting API field names to TF field names.
// This is the same direction as go-split applies at code-gen time.
func (cc *CompiledConfig) TransformForRead(apiFlat map[string]any) map[string]any {
	result := make(map[string]any, len(apiFlat))
	for k, v := range apiFlat {
		result[k] = v
	}

	// Apply forward moves in order
	for _, m := range cc.forwardMoves {
		applyMove(m.from, m.to, result)
	}

	// Apply removes
	for path := range cc.removes {
		removeByPrefix(path, result)
	}

	return result
}

// TransformForWrite applies inverse moves to a TF flat map,
// converting TF field names back to API field names.
func (cc *CompiledConfig) TransformForWrite(tfFlat map[string]any) map[string]any {
	result := make(map[string]any, len(tfFlat))
	for k, v := range tfFlat {
		result[k] = v
	}

	// Apply inverse moves in order (reversed from forward)
	for _, m := range cc.inverseMoves {
		applyMove(m.from, m.to, result)
	}

	return result
}

// applyMove renames keys in a flat map. Handles both exact matches and prefix matches.
// If `to` is "", it unnests (promotes children to root by stripping the from prefix).
// If `from` is "", it nests (puts all root keys under the `to` prefix).
func applyMove(from, to string, flat map[string]any) {
	if from == "" && to != "" {
		// Nest: wrap all current keys under `to`
		keys := make([]string, 0, len(flat))
		for k := range flat {
			keys = append(keys, k)
		}
		for _, key := range keys {
			val := flat[key]
			delete(flat, key)
			flat[to+"."+key] = val
		}

		return
	}

	// Move exact key
	if val, exists := flat[from]; exists {
		if to == "" {
			// Unnesting: just delete the parent, children are handled below
			delete(flat, from)
		} else {
			delete(flat, from)
			flat[to] = val
		}
	}

	// Move all children with the prefix
	prefix := from + "."
	for key, val := range flat {
		if strings.HasPrefix(key, prefix) {
			suffix := strings.TrimPrefix(key, prefix)
			delete(flat, key)

			if to == "" {
				// Unnest: promote to root
				flat[suffix] = val
			} else {
				flat[to+"."+suffix] = val
			}
		}
	}
}

// removeByPrefix removes all keys that match the path exactly or are children.
func removeByPrefix(path string, flat map[string]any) {
	delete(flat, path)

	prefix := path + "."
	for key := range flat {
		if strings.HasPrefix(key, prefix) {
			delete(flat, key)
		}
	}
}

// ApplyTypeConversionsForWrite converts TF types to API types based on templates.
// For example, template-bool fields convert true→"on", false→"off".
func (cc *CompiledConfig) ApplyTypeConversionsForWrite(flat map[string]any) {
	for path, ttype := range cc.templateTypes {
		val, exists := flat[path]
		if !exists {
			continue
		}

		switch ttype {
		case templateBool:
			if b, ok := val.(bool); ok {
				if b {
					flat[path] = "on"
				} else {
					flat[path] = "off"
				}
			}
		}
	}
}

// ApplyTypeConversionsForRead converts API types to TF types based on templates.
// For example, template-bool fields convert "on"→true, "off"→false.
func (cc *CompiledConfig) ApplyTypeConversionsForRead(flat map[string]any) {
	for path, ttype := range cc.templateTypes {
		val, exists := flat[path]
		if !exists {
			continue
		}

		switch ttype {
		case templateBool:
			if s, ok := val.(string); ok {
				switch strings.ToLower(s) {
				case "on", "true", "yes", "1":
					flat[path] = true
				case "off", "false", "no", "0":
					flat[path] = false
				}
			}
		}
	}
}

// SnakeToCamelKeys converts all snake_case keys in a flat map to camelCase.
// This is applied after inverse transforms when marshaling for the API.
// Each segment of a dot-notation key is converted independently.
func SnakeToCamelKeys(flat map[string]any) map[string]any {
	result := make(map[string]any, len(flat))

	for key, val := range flat {
		parts := strings.Split(key, ".")
		for i, part := range parts {
			// Only convert if the segment contains underscores (is snake_case)
			if strings.Contains(part, "_") {
				parts[i] = strcase.ToLowerCamel(part)
			}
		}

		result[strings.Join(parts, ".")] = val
	}

	return result
}

// CamelToSnakeKeys converts all camelCase keys in a flat map to snake_case.
// This is applied before forward transforms when unmarshaling from the API.
// Each segment of a dot-notation key is converted independently.
func CamelToSnakeKeys(flat map[string]any) map[string]any {
	result := make(map[string]any, len(flat))

	for key, val := range flat {
		parts := strings.Split(key, ".")
		for i, part := range parts {
			// Only convert if the segment has uppercase letters (is camelCase)
			if hasUpperCase(part) {
				parts[i] = strcase.ToSnake(part)
			}
		}

		result[strings.Join(parts, ".")] = val
	}

	return result
}

// hasUpperCase returns true if the string contains any uppercase letter.
func hasUpperCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}

	return false
}
