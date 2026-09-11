// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Command valuestatecheck fails the build when a generated framework *Value type
// is constructed without setting its unexported `state attr.ValueState` field.
//
// Such a value defaults to attr.ValueStateNull and lowers to a null tftypes
// value, silently discarding every attribute set on it. In a set or list that
// makes multiple elements byte-identical, which Terraform rejects with
// "Duplicate Set Element"; the fault is invisible at the attr.Value layer and
// only appears once the value is lowered to tftypes, after it has left the
// provider. See MORPH-16289 (the audit) and MORPH-16245 (the motivating defect).
//
// It reports two defect shapes and FAILS (exit 1) on any found in non-test code:
//   - danger-literal: T{attrs...} with no `state:` key
//   - danger-mutate:  x := T{} (empty) whose fields are assigned later but whose
//     state is never set anywhere in the enclosing function
//
// It deliberately ignores the benign `T{}.AttributeTypes(ctx)` method-receiver
// idiom and empty `T{}` values that are never mutated (genuine nulls). Findings
// in _test.go are reported as warnings but do not fail the build.
//
// Detection is keyed on the actual value-type names harvested from the
// generated schema_gen.go files (a struct with a `state attr.ValueState`
// field), which makes it far more precise than a textual grep and lets it catch
// construction routed through helper/factory functions and inline slice/map
// elements. Pure stdlib (go/parser), so it adds no dependency to the module.
//
// Usage: valuestatecheck [dir]   (dir defaults to ./morpheus/framework)
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func identName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return x.Sel.Name
	}

	return ""
}

func containerElemName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.ArrayType:
		return identName(x.Elt)
	case *ast.MapType:
		return identName(x.Value)
	}

	return ""
}

func exprString(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		base := exprString(x.X)
		if base == "" {
			return ""
		}

		return base + "." + x.Sel.Name
	case *ast.StarExpr:
		return exprString(x.X)
	case *ast.ParenExpr:
		return exprString(x.X)
	}

	return ""
}

// isAttrValueState reports whether the field type is syntactically
// `attr.ValueState`.
func isAttrValueState(e ast.Expr) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := sel.X.(*ast.Ident)

	return ok && x.Name == "attr" && sel.Sel.Name == "ValueState"
}

// classify reports whether a composite literal sets `state` and whether it sets
// any non-state attribute field.
func classify(lit *ast.CompositeLit) (setsState, setsAttrs bool) {
	if len(lit.Elts) == 0 {
		return false, false
	}
	keyed := false
	for _, e := range lit.Elts {
		if _, ok := e.(*ast.KeyValueExpr); ok {
			keyed = true

			break
		}
	}
	if !keyed {
		// A positional literal must list every field, so state is included.
		return true, true
	}
	for _, e := range lit.Elts {
		if kv, ok := e.(*ast.KeyValueExpr); ok {
			if id, ok := kv.Key.(*ast.Ident); ok {
				if id.Name == "state" {
					setsState = true
				} else {
					setsAttrs = true
				}
			}
		}
	}

	return
}

func fileClass(path string) string {
	switch {
	case strings.HasSuffix(path, "schema_gen.go"):
		return "generated"
	case strings.HasSuffix(path, "_test.go"):
		return "test"
	default:
		return "handwritten"
	}
}

func rel(path string) string {
	if i := strings.Index(path, "morpheus/framework/"); i >= 0 {
		return path[i:]
	}

	return path
}

type finding struct {
	path, typ, cat, note string
	line                 int
}

func main() {
	dir := "morpheus/framework"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fset := token.NewFileSet()
	parsed := map[string]*ast.File{}
	var parseErrs []string
	//nolint:gosec // dir is a source tree given on the command line, not attacker input
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// testdata holds intentionally-unbuildable fixtures; skip them so a
			// fixture that will not parse cannot fail the gate.
			if name := d.Name(); name == "testdata" || name == "vendor" {
				return filepath.SkipDir
			}

			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		af, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			// A file we cannot parse is a file we cannot verify: surface it and
			// fail, rather than silently dropping it and still reporting OK — an
			// unchecked file could hide the very construction this gate exists to
			// catch.
			parseErrs = append(parseErrs, fmt.Sprintf("%s: %v", rel(path), perr))

			return nil
		}
		parsed[path] = af

		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "valuestatecheck: walking %s: %v\n", dir, err)
		os.Exit(2)
	}
	if len(parseErrs) > 0 {
		fmt.Fprintf(os.Stderr, "valuestatecheck: cannot verify %d unparseable Go file(s):\n", len(parseErrs))
		for _, e := range parseErrs {
			fmt.Fprintf(os.Stderr, "  %s\n", e)
		}
		os.Exit(2)
	}
	if len(parsed) == 0 {
		fmt.Fprintf(os.Stderr, "valuestatecheck: no Go files found under %s\n", dir)
		os.Exit(2)
	}

	// Harvest value-type names: a struct with a `state attr.ValueState` field.
	valueTypes := map[string]bool{}
	for _, af := range parsed {
		for _, decl := range af.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if st, ok := ts.Type.(*ast.StructType); ok {
					for _, f := range st.Fields.List {
						if len(f.Names) == 1 && f.Names[0].Name == "state" && isAttrValueState(f.Type) {
							valueTypes[ts.Name.Name] = true
						}
					}
				}
			}
		}
	}

	var findings []finding
	seen := map[string]bool{}
	add := func(lit *ast.CompositeLit, name, cat, note string) {
		pos := fset.Position(lit.Pos())
		key := fmt.Sprintf("%s:%d:%d:%s", pos.Filename, pos.Line, pos.Column, cat)
		if seen[key] {
			return
		}
		seen[key] = true
		findings = append(findings, finding{pos.Filename, name, cat, note, pos.Line})
	}

	for _, af := range parsed {
		// Literals used as a method receiver (T{}.Method) are not construction.
		receiver := map[*ast.CompositeLit]bool{}
		ast.Inspect(af, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if cl, ok := sel.X.(*ast.CompositeLit); ok {
					receiver[cl] = true
				}
			}

			return true
		})

		checkLit := func(lit *ast.CompositeLit, name string) {
			if receiver[lit] {
				return
			}
			if s, a := classify(lit); !s && a {
				add(lit, name, "danger-literal", "attributes set in the literal, state omitted")
			}
		}
		ast.Inspect(af, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if name := identName(lit.Type); name != "" && valueTypes[name] {
				checkLit(lit, name)
			}
			if elem := containerElemName(lit.Type); elem != "" && valueTypes[elem] {
				for _, e := range lit.Elts {
					el := e
					if kv, ok := e.(*ast.KeyValueExpr); ok {
						el = kv.Value
					}
					if inner, ok := el.(*ast.CompositeLit); ok && inner.Type == nil {
						checkLit(inner, elem)
					}
				}
			}

			return true
		})

		// Zero-value-then-mutate, analysed within each function.
		for _, decl := range af.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			stateSet, fieldSet := map[string]bool{}, map[string]bool{}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if asn, ok := n.(*ast.AssignStmt); ok {
					for _, lhs := range asn.Lhs {
						if sel, ok := lhs.(*ast.SelectorExpr); ok {
							if base := exprString(sel.X); base != "" {
								if sel.Sel.Name == "state" {
									stateSet[base] = true
								} else {
									fieldSet[base] = true
								}
							}
						}
					}
				}

				return true
			})
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				asn, ok := n.(*ast.AssignStmt)
				if !ok {
					return true
				}
				for i, rhs := range asn.Rhs {
					lit, ok := rhs.(*ast.CompositeLit)
					if !ok {
						continue
					}
					name := identName(lit.Type)
					if name == "" || !valueTypes[name] || receiver[lit] {
						continue
					}
					if s, a := classify(lit); s || a {
						continue
					}
					if i >= len(asn.Lhs) {
						continue
					}
					if target := exprString(asn.Lhs[i]); target != "" && fieldSet[target] && !stateSet[target] {
						add(lit, name, "danger-mutate",
							"empty "+name+"{} -> "+target+": fields assigned later but state never set in function")
					}
				}

				return true
			})
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].path != findings[j].path {
			return findings[i].path < findings[j].path
		}

		return findings[i].line < findings[j].line
	})

	var prod, tests []finding
	for _, f := range findings {
		if fileClass(f.path) == "test" {
			tests = append(tests, f)
		} else {
			prod = append(prod, f)
		}
	}

	for _, f := range tests {
		fmt.Fprintf(os.Stderr, "warning: %s:%d: %s constructed without state [%s] (test)\n",
			rel(f.path), f.line, f.typ, f.cat)
	}
	if len(prod) == 0 {
		fmt.Printf("valuestatecheck: OK — %d value types, no production *Value constructed without state\n", len(valueTypes))

		return
	}
	for _, f := range prod {
		fmt.Fprintf(os.Stderr, "%s:%d: %s constructed without state: %s [%s]\n",
			rel(f.path), f.line, f.typ, f.note, f.cat)
	}
	fmt.Fprintf(os.Stderr, "\nvaluestatecheck: FAIL — %d production *Value construction(s) omit state (MORPH-16289).\n"+
		"Set `state: attr.ValueStateKnown`, use the New*Value(...) factory, or New*ValueNull().\n", len(prod))
	os.Exit(1)
}
