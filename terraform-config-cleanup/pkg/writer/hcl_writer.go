package writer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// WriteFile writes a TerraformFile back to HCL format
func WriteFile(tfFile *parser.TerraformFile, filename string) error {
	var sb strings.Builder

	for i, resource := range tfFile.Resources {
		if i > 0 {
			sb.WriteString("\n")
		}
		if err := writeResource(&sb, resource); err != nil {
			return fmt.Errorf("failed to write resource %s.%s: %w", resource.Type, resource.Name, err)
		}
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// writeResource writes a single resource block
func writeResource(sb *strings.Builder, resource *parser.Resource) error {
	sb.WriteString(fmt.Sprintf("resource %q %q {\n", resource.Type, resource.Name))

	// Sort attributes for consistent output
	var keys []string
	for key := range resource.Attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		attr := resource.Attributes[key]
		if err := writeAttribute(sb, key, attr.Value, 1); err != nil {
			return err
		}
	}

	sb.WriteString("}\n")
	return nil
}

// writeAttribute writes an attribute with proper indentation
func writeAttribute(sb *strings.Builder, name string, value interface{}, indent int) error {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case nil:
		sb.WriteString(fmt.Sprintf("%s%s = null\n", indentStr, name))

	case bool:
		sb.WriteString(fmt.Sprintf("%s%s = %t\n", indentStr, name, v))

	case int:
		sb.WriteString(fmt.Sprintf("%s%s = %d\n", indentStr, name, v))

	case int64:
		sb.WriteString(fmt.Sprintf("%s%s = %d\n", indentStr, name, v))

	case float64:
		// Check if it's actually an integer
		if v == float64(int64(v)) {
			sb.WriteString(fmt.Sprintf("%s%s = %d\n", indentStr, name, int64(v)))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s = %g\n", indentStr, name, v))
		}

	case string:
		// Check if it's a reference or expression
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			// It's an expression, don't quote it
			expr := strings.TrimSuffix(strings.TrimPrefix(v, "${"), "}")
			sb.WriteString(fmt.Sprintf("%s%s = %s\n", indentStr, name, expr))
		} else {
			// Regular string, quote it
			sb.WriteString(fmt.Sprintf("%s%s = %q\n", indentStr, name, v))
		}

	case []interface{}:
		if len(v) == 0 {
			sb.WriteString(fmt.Sprintf("%s%s = []\n", indentStr, name))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s = [\n", indentStr, name))
			for _, item := range v {
				if err := writeListItem(sb, item, indent+1); err != nil {
					return err
				}
			}
			sb.WriteString(fmt.Sprintf("%s]\n", indentStr))
		}

	case map[string]interface{}:
		// Check if this should be a block or an object
		if shouldBeBlock(v) {
			// Write as nested block
			sb.WriteString(fmt.Sprintf("%s%s {\n", indentStr, name))

			// Sort keys for consistent output
			var keys []string
			for key := range v {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				if err := writeAttribute(sb, key, v[key], indent+1); err != nil {
					return err
				}
			}
			sb.WriteString(fmt.Sprintf("%s}\n", indentStr))
		} else {
			// Write as object expression
			if err := writeMapValue(sb, indentStr, name, v); err != nil {
				return err
			}
		}

	default:
		sb.WriteString(fmt.Sprintf("%s%s = %v\n", indentStr, name, v))
	}

	return nil
}

// writeListItem writes a single item in a list
func writeListItem(sb *strings.Builder, value interface{}, indent int) error {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case nil:
		sb.WriteString(fmt.Sprintf("%snull,\n", indentStr))

	case bool:
		sb.WriteString(fmt.Sprintf("%s%t,\n", indentStr, v))

	case int, int64:
		sb.WriteString(fmt.Sprintf("%s%d,\n", indentStr, v))

	case float64:
		if v == float64(int64(v)) {
			sb.WriteString(fmt.Sprintf("%s%d,\n", indentStr, int64(v)))
		} else {
			sb.WriteString(fmt.Sprintf("%s%g,\n", indentStr, v))
		}

	case string:
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			expr := strings.TrimSuffix(strings.TrimPrefix(v, "${"), "}")
			sb.WriteString(fmt.Sprintf("%s%s,\n", indentStr, expr))
		} else {
			sb.WriteString(fmt.Sprintf("%s%q,\n", indentStr, v))
		}

	case map[string]interface{}:
		sb.WriteString(fmt.Sprintf("%s{\n", indentStr))
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			itemIndent := strings.Repeat("  ", indent+1)
			sb.WriteString(fmt.Sprintf("%s%s = ", itemIndent, key))
			if err := writeInlineValue(sb, v[key]); err != nil {
				return err
			}
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%s},\n", indentStr))

	default:
		sb.WriteString(fmt.Sprintf("%s%v,\n", indentStr, v))
	}

	return nil
}

// writeInlineValue writes a value inline (for use in objects)
func writeInlineValue(sb *strings.Builder, value interface{}) error {
	switch v := value.(type) {
	case nil:
		sb.WriteString("null")
	case bool:
		sb.WriteString(fmt.Sprintf("%t", v))
	case int, int64:
		sb.WriteString(fmt.Sprintf("%d", v))
	case float64:
		if v == float64(int64(v)) {
			sb.WriteString(fmt.Sprintf("%d", int64(v)))
		} else {
			sb.WriteString(fmt.Sprintf("%g", v))
		}
	case string:
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			expr := strings.TrimSuffix(strings.TrimPrefix(v, "${"), "}")
			sb.WriteString(expr)
		} else {
			sb.WriteString(fmt.Sprintf("%q", v))
		}
	case []interface{}:
		sb.WriteString("[")
		for i, item := range v {
			if i > 0 {
				sb.WriteString(", ")
			}
			if err := writeInlineValue(sb, item); err != nil {
				return err
			}
		}
		sb.WriteString("]")
	case map[string]interface{}:
		sb.WriteString("{")
		i := 0
		for key, val := range v {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s = ", key))
			if err := writeInlineValue(sb, val); err != nil {
				return err
			}
			i++
		}
		sb.WriteString("}")
	default:
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	return nil
}

// writeMapValue writes a map as an object expression
func writeMapValue(sb *strings.Builder, indentStr, name string, m map[string]interface{}) error {
	sb.WriteString(fmt.Sprintf("%s%s = {\n", indentStr, name))

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		valueIndent := indentStr + "  "
		sb.WriteString(fmt.Sprintf("%s%s = ", valueIndent, key))
		if err := writeInlineValue(sb, m[key]); err != nil {
			return err
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("%s}\n", indentStr))
	return nil
}

// shouldBeBlock determines if a map should be written as a block vs object
// Heuristic: if all values are simple types or nested maps, it's likely a block
func shouldBeBlock(m map[string]interface{}) bool {
	// Common block names in Terraform configs
	blockNames := map[string]bool{
		"config":          true,
		"lifecycle":       true,
		"timeouts":        true,
		"connection":      true,
		"provisioner":     true,
		"config_approval": true,
		"policy_type":     true,
	}

	// If this map has keys that match common block names, it's likely a block
	for key := range m {
		if blockNames[key] {
			return true
		}
	}

	// Default to block style for complex nested structures
	return true
}
