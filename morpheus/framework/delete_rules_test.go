// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package framework_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Morpheus renders a refused delete as HTTP 200 with success:false, so a
// Delete that judges the DELETE by status code alone removes the resource from
// state while it remains on the appliance (MORPH-17987). errfmt.CheckDeleted
// reads the body; this rule keeps every Delete on it.
//
// The rule is on the source, not the schema: for every framework resource
// package, every SDK call in a Delete method (or in a function it calls) that
// looks like a delete -- a chain ending in .Execute() whose method name starts
// with Delete or Remove, or an ExecuteContainerAction -- must not discard its
// first result, and that result must reach errfmt.CheckDeleted.
//
// Resources that delete through some other means (a membership rewrite that is
// verified by re-reading, a resource that removes nothing on the appliance)
// have no such call and pass trivially.

// deleteCallPattern is what a delete-shaped SDK call's method name looks like.
func isDeleteMethod(name string) bool {
	return strings.HasPrefix(name, "Delete") || strings.HasPrefix(name, "Remove") ||
		name == "ExecuteContainerAction"
}

// executeChainMethod returns the SDK method name of an X.Y(...).Z(...).Execute()
// chain, or "" when call is not such a chain.
func executeChainMethod(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Execute" {
		return ""
	}

	// Walk down the chain to the outermost method with arguments (the API
	// method); option setters like .Force(true) sit in between.
	cur := sel.X

	for {
		inner, ok := cur.(*ast.CallExpr)
		if !ok {
			return ""
		}

		isel, ok := inner.Fun.(*ast.SelectorExpr)
		if !ok {
			return ""
		}

		if isDeleteMethod(isel.Sel.Name) {
			return isel.Sel.Name
		}

		cur = isel.X
	}
}

// deleteCallsInFunc finds delete-shaped Execute() chains assigned in fn and
// reports, per call, the name the result is bound to ("_" when discarded).
func deleteCallsInFunc(fn *ast.FuncDecl) map[*ast.CallExpr]string {
	out := map[*ast.CallExpr]string{}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || len(as.Rhs) != 1 {
			return true
		}

		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok || executeChainMethod(call) == "" {
			return true
		}

		if len(as.Lhs) == 0 {
			return true
		}

		if id, ok := as.Lhs[0].(*ast.Ident); ok {
			out[call] = id.Name
		}

		return true
	})

	return out
}

// resultReachesCheckDeleted reports whether fn contains a call
// errfmt.CheckDeleted(<name>, ...).
func resultReachesCheckDeleted(fn *ast.FuncDecl, name string) bool {
	found := false

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "CheckDeleted" {
			return true
		}

		if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "errfmt" {
			return true
		}

		if len(call.Args) > 0 {
			if id, ok := call.Args[0].(*ast.Ident); ok && id.Name == name {
				found = true
			}
		}

		return true
	})

	return found
}

func TestUnitDeleteImplementationsCheckTheResponseBody(t *testing.T) {
	t.Parallel()

	root := filepath.Join("resources")

	dirs, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}

	var violations []string

	checked := 0

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		fset := token.NewFileSet()

		entries, err := os.ReadDir(filepath.Join(root, d.Name()))
		if err != nil {
			t.Fatalf("%s: %v", d.Name(), err)
		}

		var files []*ast.File

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}

			file, err := parser.ParseFile(fset, filepath.Join(root, d.Name(), e.Name()), nil, 0)
			if err != nil {
				t.Fatalf("%s/%s: %v", d.Name(), e.Name(), err)
			}

			files = append(files, file)
		}

		{
			for _, file := range files {
				for _, decl := range file.Decls {
					fn, ok := decl.(*ast.FuncDecl)
					if !ok || fn.Body == nil {
						continue
					}

					for call, bound := range deleteCallsInFunc(fn) {
						checked++

						method := executeChainMethod(call)
						where := d.Name() + "." + fn.Name.Name + " (" + method + ")"

						switch {
						case bound == "_":
							violations = append(violations, where+": the delete's result is discarded; bind it and pass it to errfmt.CheckDeleted")
						case !resultReachesCheckDeleted(fn, bound):
							violations = append(violations, where+": result "+bound+" is not passed to errfmt.CheckDeleted")
						}
					}
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("found no delete calls at all; the rule's pattern no longer matches the code")
	}

	sort.Strings(violations)

	for _, v := range violations {
		t.Error(v)
	}

	t.Logf("checked %d delete-shaped SDK calls", checked)
}
