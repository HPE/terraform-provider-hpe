package parser

import (
	"fmt"
	"strings"
)

// TerraformFile represents a parsed Terraform configuration file
type TerraformFile struct {
	Resources []*Resource
}

// Resource represents a Terraform resource block
type Resource struct {
	Type       string // e.g., "hpe_morpheus_instance"
	Name       string // e.g., "example"
	Attributes map[string]*Attribute
}

// Attribute represents an attribute in a resource
type Attribute struct {
	Path  string      // Dot-notation path, e.g., "config.noAgent.Bool"
	Value interface{} // The actual value (can be primitive, map, or slice)
}

// GetAttribute retrieves an attribute by its dot-notation path
func (r *Resource) GetAttribute(path string) (*Attribute, bool) {
	attr, exists := r.Attributes[path]
	return attr, exists
}

// SetAttribute sets an attribute at the given dot-notation path
func (r *Resource) SetAttribute(path string, value interface{}) {
	if r.Attributes == nil {
		r.Attributes = make(map[string]*Attribute)
	}
	r.Attributes[path] = &Attribute{
		Path:  path,
		Value: value,
	}
}

// RemoveAttribute removes an attribute at the given dot-notation path
func (r *Resource) RemoveAttribute(path string) {
	delete(r.Attributes, path)
}

// GetAttributeValue retrieves the value at a nested path
func (r *Resource) GetAttributeValue(path string) (interface{}, bool) {
	parts := strings.Split(path, ".")

	// Start with the root attribute
	rootAttr, exists := r.Attributes[parts[0]]
	if !exists {
		return nil, false
	}

	// If this is the full path, return it
	if len(parts) == 1 {
		return rootAttr.Value, true
	}

	// Navigate through nested maps
	current := rootAttr.Value
	for i := 1; i < len(parts); i++ {
		switch v := current.(type) {
		case map[string]interface{}:
			next, ok := v[parts[i]]
			if !ok {
				return nil, false
			}
			current = next
		default:
			return nil, false
		}
	}

	return current, true
}

// SetAttributeValue sets the value at a nested path
func (r *Resource) SetAttributeValue(path string, value interface{}) error {
	parts := strings.Split(path, ".")

	if len(parts) == 1 {
		r.SetAttribute(path, value)
		return nil
	}

	// Ensure root exists
	rootAttr, exists := r.Attributes[parts[0]]
	if !exists {
		rootAttr = &Attribute{
			Path:  parts[0],
			Value: make(map[string]interface{}),
		}
		r.Attributes[parts[0]] = rootAttr
	}

	// Navigate and create intermediate maps
	current := rootAttr.Value
	for i := 1; i < len(parts)-1; i++ {
		m, ok := current.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot navigate through non-map at %s", strings.Join(parts[:i+1], "."))
		}

		next, exists := m[parts[i]]
		if !exists {
			next = make(map[string]interface{})
			m[parts[i]] = next
		}
		current = next
	}

	// Set the final value
	m, ok := current.(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot set value on non-map at %s", strings.Join(parts[:len(parts)-1], "."))
	}
	m[parts[len(parts)-1]] = value

	return nil
}

// RemoveAttributeValue removes the value at a nested path
func (r *Resource) RemoveAttributeValue(path string) error {
	parts := strings.Split(path, ".")

	if len(parts) == 1 {
		r.RemoveAttribute(path)
		return nil
	}

	// Navigate to parent
	rootAttr, exists := r.Attributes[parts[0]]
	if !exists {
		return nil // Nothing to remove
	}

	current := rootAttr.Value
	for i := 1; i < len(parts)-1; i++ {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil // Path doesn't exist
		}

		next, exists := m[parts[i]]
		if !exists {
			return nil // Path doesn't exist
		}
		current = next
	}

	// Remove the final key
	m, ok := current.(map[string]interface{})
	if !ok {
		return nil
	}
	delete(m, parts[len(parts)-1])

	return nil
}

// ListMatchingPaths returns all attribute paths matching a wildcard pattern
func (r *Resource) ListMatchingPaths(pattern string, except []string) []string {
	var matches []string

	// Build exception map for fast lookup
	exceptMap := make(map[string]bool)
	for _, e := range except {
		exceptMap[e] = true
	}

	// Check if pattern has wildcard
	if !strings.Contains(pattern, "*") {
		// No wildcard, just check if exact match exists
		if _, exists := r.Attributes[pattern]; exists && !exceptMap[pattern] {
			matches = append(matches, pattern)
		}
		return matches
	}

	// Handle wildcard patterns
	prefix := strings.TrimSuffix(pattern, "*")

	for path := range r.Attributes {
		// Check if it matches the prefix and is not in except list
		if strings.HasPrefix(path, prefix) && !exceptMap[path] {
			matches = append(matches, path)
		}
	}

	return matches
}
