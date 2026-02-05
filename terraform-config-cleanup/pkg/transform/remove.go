package transform

import (
	"fmt"

	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// ApplyRemoves applies remove operations to a resource
// Removes can be:
// - Simple strings: "config.backup"
// - Maps with path, except, and when conditions
func ApplyRemoves(resource *parser.Resource, removes []interface{}) error {
	for _, removeOp := range removes {
		if err := applyRemove(resource, removeOp); err != nil {
			return err
		}
	}
	return nil
}

// applyRemove applies a single remove operation
func applyRemove(resource *parser.Resource, removeOp interface{}) error {
	switch r := removeOp.(type) {
	case string:
		// Simple remove: just a path string
		return applySimpleRemove(resource, r)

	case map[string]interface{}:
		// Complex remove with path, except, and/or when
		return applyConditionalRemove(resource, r)

	default:
		return fmt.Errorf("invalid remove operation format: %T", removeOp)
	}
}

// applySimpleRemove removes a simple path
func applySimpleRemove(resource *parser.Resource, path string) error {
	// Check if it's a wildcard pattern
	if containsWildcard(path) {
		matches := resource.ListMatchingPaths(path, nil)
		for _, match := range matches {
			if err := resource.RemoveAttributeValue(match); err != nil {
				return fmt.Errorf("failed to remove %s: %w", match, err)
			}
		}
		return nil
	}

	// Simple removal
	return resource.RemoveAttributeValue(path)
}

// applyConditionalRemove applies a remove with conditions
func applyConditionalRemove(resource *parser.Resource, removeMap map[string]interface{}) error {
	// Extract path
	pathValue, ok := removeMap["path"]
	if !ok {
		return fmt.Errorf("conditional remove must have 'path' field")
	}
	path, ok := pathValue.(string)
	if !ok {
		return fmt.Errorf("'path' must be a string")
	}

	// Check when condition if present
	if whenValue, ok := removeMap["when"]; ok {
		whenMap, ok := whenValue.(map[string]interface{})
		if !ok {
			return fmt.Errorf("'when' must be a map")
		}

		if !evaluateWhenCondition(resource, whenMap) {
			// Condition not met, skip this removal
			return nil
		}
	}

	// Extract except list if present
	var except []string
	if exceptValue, ok := removeMap["except"]; ok {
		exceptList, ok := exceptValue.([]interface{})
		if !ok {
			return fmt.Errorf("'except' must be a list")
		}
		for _, e := range exceptList {
			if s, ok := e.(string); ok {
				except = append(except, s)
			}
		}
	}

	// Apply the removal with wildcard support
	if containsWildcard(path) {
		matches := resource.ListMatchingPaths(path, except)
		for _, match := range matches {
			if err := resource.RemoveAttributeValue(match); err != nil {
				return fmt.Errorf("failed to remove %s: %w", match, err)
			}
		}
		return nil
	}

	// Check if path is in except list
	for _, e := range except {
		if path == e {
			return nil
		}
	}

	return resource.RemoveAttributeValue(path)
}

// evaluateWhenCondition checks if a when condition is satisfied
func evaluateWhenCondition(resource *parser.Resource, whenMap map[string]interface{}) bool {
	// Check for "and" operator
	if andValue, ok := whenMap["and"]; ok {
		andList, ok := andValue.([]interface{})
		if !ok {
			return false
		}

		for _, condItem := range andList {
			condMap, ok := condItem.(map[string]interface{})
			if !ok {
				return false
			}
			if !evaluateSingleCondition(resource, condMap) {
				return false
			}
		}
		return true
	}

	// Check for "or" operator
	if orValue, ok := whenMap["or"]; ok {
		orList, ok := orValue.([]interface{})
		if !ok {
			return false
		}

		for _, condItem := range orList {
			condMap, ok := condItem.(map[string]interface{})
			if !ok {
				continue
			}
			if evaluateSingleCondition(resource, condMap) {
				return true
			}
		}
		return false
	}

	// No and/or, evaluate as single condition
	return evaluateSingleCondition(resource, whenMap)
}

// evaluateSingleCondition evaluates a single condition
func evaluateSingleCondition(resource *parser.Resource, condMap map[string]interface{}) bool {
	// Get the attribute to check
	attrPath, ok := condMap["attribute"].(string)
	if !ok {
		return false
	}

	attrValue, exists := resource.GetAttributeValue(attrPath)

	// Check "is_null" condition
	if isNullValue, ok := condMap["is_null"]; ok {
		isNull, _ := isNullValue.(bool)
		if isNull {
			return !exists || attrValue == nil
		}
		return exists && attrValue != nil
	}

	if !exists {
		return false
	}

	// Check "in" condition
	if inValue, ok := condMap["in"]; ok {
		inList, ok := inValue.([]interface{})
		if !ok {
			return false
		}

		attrStr := fmt.Sprintf("%v", attrValue)
		for _, item := range inList {
			if fmt.Sprintf("%v", item) == attrStr {
				return true
			}
		}
		return false
	}

	// Check "equals" condition
	if equalsValue, ok := condMap["equals"]; ok {
		return fmt.Sprintf("%v", attrValue) == fmt.Sprintf("%v", equalsValue)
	}

	// Check "not_equals" condition
	if notEqualsValue, ok := condMap["not_equals"]; ok {
		return fmt.Sprintf("%v", attrValue) != fmt.Sprintf("%v", notEqualsValue)
	}

	return false
}

// containsWildcard checks if a path contains a wildcard character
func containsWildcard(path string) bool {
	return len(path) > 0 && path[len(path)-1] == '*'
}
