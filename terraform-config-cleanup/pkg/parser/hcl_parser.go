package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// ParseFile parses a Terraform configuration file
func ParseFile(filename string) (*TerraformFile, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(src, filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	tfFile := &TerraformFile{
		Resources: []*Resource{},
	}

	// Parse the body
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("unexpected body type")
	}

	// Extract resources
	for _, block := range body.Blocks {
		if block.Type == "resource" {
			resource, err := parseResource(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse resource: %w", err)
			}
			tfFile.Resources = append(tfFile.Resources, resource)
		}
	}

	return tfFile, nil
}

// parseResource parses a resource block into our AST
func parseResource(block *hclsyntax.Block) (*Resource, error) {
	if len(block.Labels) < 2 {
		return nil, fmt.Errorf("resource block must have type and name labels")
	}

	resource := &Resource{
		Type:       block.Labels[0],
		Name:       block.Labels[1],
		Attributes: make(map[string]*Attribute),
	}

	// Parse attributes
	for name, attr := range block.Body.Attributes {
		value, err := exprToValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse attribute %s: %w", name, err)
		}

		resource.Attributes[name] = &Attribute{
			Path:  name,
			Value: value,
		}
	}

	// Parse nested blocks (like config, lifecycle, etc.)
	for _, nestedBlock := range block.Body.Blocks {
		blockValue, err := parseBlock(nestedBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block %s: %w", nestedBlock.Type, err)
		}

		resource.Attributes[nestedBlock.Type] = &Attribute{
			Path:  nestedBlock.Type,
			Value: blockValue,
		}
	}

	return resource, nil
}

// parseBlock recursively parses a nested block
func parseBlock(block *hclsyntax.Block) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Parse attributes in this block
	for name, attr := range block.Body.Attributes {
		value, err := exprToValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse attribute %s: %w", name, err)
		}
		result[name] = value
	}

	// Parse nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		blockValue, err := parseBlock(nestedBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to parse nested block %s: %w", nestedBlock.Type, err)
		}

		// If there are multiple blocks with the same type, make it a list
		if existing, exists := result[nestedBlock.Type]; exists {
			// Convert to slice if not already
			switch v := existing.(type) {
			case []interface{}:
				result[nestedBlock.Type] = append(v, blockValue)
			default:
				result[nestedBlock.Type] = []interface{}{v, blockValue}
			}
		} else {
			result[nestedBlock.Type] = blockValue
		}
	}

	return result, nil
}

// exprToValue converts an HCL expression to a Go value
func exprToValue(expr hclsyntax.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return ctyToGo(e.Val), nil

	case *hclsyntax.TemplateExpr:
		// Simple string template
		if len(e.Parts) == 1 {
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return ctyToGo(lit.Val), nil
			}
		}
		// For complex templates, try to evaluate
		val, diags := e.Value(nil)
		if !diags.HasErrors() {
			return ctyToGo(val), nil
		}
		return fmt.Sprintf("${...}"), nil // Fallback for complex expressions

	case *hclsyntax.ObjectConsExpr:
		result := make(map[string]interface{})
		for _, item := range e.Items {
			keyVal, diags := item.KeyExpr.Value(nil)
			if diags.HasErrors() {
				continue
			}
			key := ctyToGo(keyVal).(string)

			val, err := exprToValue(item.ValueExpr)
			if err != nil {
				return nil, err
			}
			result[key] = val
		}
		return result, nil

	case *hclsyntax.TupleConsExpr:
		var result []interface{}
		for _, elem := range e.Exprs {
			val, err := exprToValue(elem)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		}
		return result, nil

	case *hclsyntax.ScopeTraversalExpr:
		// Variable reference - format as reference expression
		parts := make([]string, len(e.Traversal))
		for i, traverser := range e.Traversal {
			switch t := traverser.(type) {
			case hcl.TraverseRoot:
				parts[i] = t.Name
			case hcl.TraverseAttr:
				parts[i] = "." + t.Name
			case hcl.TraverseIndex:
				parts[i] = fmt.Sprintf("[%v]", t.Key)
			}
		}
		return fmt.Sprintf("${%s}", strings.Join(parts, "")), nil

	case *hclsyntax.FunctionCallExpr:
		// Function call - preserve as string for now
		return fmt.Sprintf("${%s(...)}", e.Name), nil

	default:
		// Try to evaluate the expression
		val, diags := expr.Value(nil)
		if !diags.HasErrors() {
			return ctyToGo(val), nil
		}
		// Fallback: return a placeholder
		return fmt.Sprintf("${...}"), nil
	}
}

// ctyToGo converts a cty.Value to a Go value
func ctyToGo(val cty.Value) interface{} {
	if val.IsNull() {
		return nil
	}

	ty := val.Type()

	switch {
	case ty == cty.String:
		return val.AsString()
	case ty == cty.Number:
		bf := val.AsBigFloat()
		if bf.IsInt() {
			i, _ := bf.Int64()
			return i
		}
		f, _ := bf.Float64()
		return f
	case ty == cty.Bool:
		return val.True()
	case ty.IsListType() || ty.IsTupleType() || ty.IsSetType():
		var result []interface{}
		it := val.ElementIterator()
		for it.Next() {
			_, v := it.Element()
			result = append(result, ctyToGo(v))
		}
		return result
	case ty.IsMapType() || ty.IsObjectType():
		result := make(map[string]interface{})
		it := val.ElementIterator()
		for it.Next() {
			k, v := it.Element()
			result[k.AsString()] = ctyToGo(v)
		}
		return result
	default:
		return nil
	}
}
