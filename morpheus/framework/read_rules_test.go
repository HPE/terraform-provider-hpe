// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package framework

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// A resource deleted outside Terraform must leave state on the next refresh,
// not fail it. The framework convention is that Read answers a 404 from the
// API with resp.State.RemoveResource(ctx), so that the next plan recreates the
// resource; a Read that instead reports the 404 as an error leaves the
// practitioner with `terraform state rm` as the only way forward (MORPH-18584,
// first seen on hpe_morpheus_cluster).
//
// The rule is on the source: for every framework resource package, the Read
// method — or a function in the same package that Read reaches — must call
// RemoveResource. Reachability is by name within the package, which is how
// every resource here structures its read (Read → getXAsState / populateX).
//
// readRemoveExemptions lists the packages whose resource cannot be deleted out
// of band, with the reason.
var readRemoveExemptions = map[string]string{
	// A singleton appliance setting: there is nothing to delete, GET never 404s.
	"settingwhitelabel": "singleton setting",
}

// funcsByName indexes a package's function and method declarations by name.
// Methods and functions share the namespace, which is good enough for a
// reachability walk: a resource's read helpers have distinctive names.
func funcsByName(files []*ast.File) map[string][]*ast.FuncDecl {
	out := map[string][]*ast.FuncDecl{}
	for _, file := range files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
				out[fn.Name.Name] = append(out[fn.Name.Name], fn)
			}
		}
	}

	return out
}

// isReadMethod reports whether fn is a resource Read: a method named Read
// taking a resource.ReadRequest.
func isReadMethod(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || fn.Name.Name != "Read" || fn.Type.Params == nil {
		return false
	}
	for _, p := range fn.Type.Params.List {
		if sel, ok := p.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "ReadRequest" {
			return true
		}
	}

	return false
}

// callsRemoveResource reports whether fn, or any same-package function it
// calls (transitively), calls something named RemoveResource.
func callsRemoveResource(fn *ast.FuncDecl, index map[string][]*ast.FuncDecl, seen map[*ast.FuncDecl]bool) bool {
	if seen[fn] {
		return false
	}
	seen[fn] = true

	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		var name string
		switch f := call.Fun.(type) {
		case *ast.SelectorExpr:
			name = f.Sel.Name
		case *ast.Ident:
			name = f.Name
		default:
			return true
		}
		if name == "RemoveResource" {
			found = true

			return false
		}
		for _, callee := range index[name] {
			if callsRemoveResource(callee, index, seen) {
				found = true

				return false
			}
		}

		return true
	})

	return found
}

func TestUnitReadImplementationsRemoveTheResourceOn404(t *testing.T) {
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
		if reason, ok := readRemoveExemptions[d.Name()]; ok {
			t.Logf("%s: exempt (%s)", d.Name(), reason)

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

		index := funcsByName(files)

		for _, fns := range index {
			for _, fn := range fns {
				if !isReadMethod(fn) {
					continue
				}
				checked++
				if !callsRemoveResource(fn, index, map[*ast.FuncDecl]bool{}) {
					violations = append(violations, d.Name()+
						": Read never calls resp.State.RemoveResource; a 404 from the API must remove the resource from state, not fail the refresh")
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("found no Read methods at all; the rule's pattern no longer matches the code")
	}

	sort.Strings(violations)

	for _, v := range violations {
		t.Error(v)
	}

	t.Logf("checked %d Read methods", checked)
}

// The `found bool` convention has a second invariant the reachability rule
// above cannot see: a helper may answer "not found" only for a 404 or after
// recording an error — never on a success path. A success path returning false
// makes Read remove a healthy resource from state (the ostypeimage review
// finding on MORPH-18584). The rule: in any function reachable from Read that
// returns a bool, every `return` whose bool is the literal false must sit in a
// block that either tests errfmt.IsNotFound (or http.StatusNotFound) or has
// already called AddError / Append on a diagnostics value.
//
// Helpers that use the inverse convention — a `notFound` bool — return true for
// a 404; they are recognised by the result being named notFound (or similar)
// and are skipped here, since their false is the success value.

// boolResultIndex returns the index of the bool result of a found-signalling
// helper — a function whose results include both a bool and diag.Diagnostics
// — or -1 for anything else (an Equal method, a plain predicate, ...). It also
// reports whether that result is named as an inverse flag.
func boolResultIndex(fn *ast.FuncDecl) (idx int, inverse bool) {
	if fn.Type.Results == nil {
		return -1, false
	}
	hasDiags := false
	for _, f := range fn.Type.Results.List {
		if sel, ok := f.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "Diagnostics" {
			hasDiags = true
		}
	}
	if !hasDiags {
		return -1, false
	}
	i := 0
	for _, f := range fn.Type.Results.List {
		n := max(len(f.Names), 1)
		if id, ok := f.Type.(*ast.Ident); ok && id.Name == "bool" {
			inv := false
			for _, name := range f.Names {
				lower := strings.ToLower(name.Name)
				if strings.Contains(lower, "notfound") || strings.Contains(lower, "absent") || strings.Contains(lower, "gone") {
					inv = true
				}
			}

			return i, inv
		}
		i += n
	}

	return -1, false
}

// isFalse reports whether e is the literal false.
func isFalse(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)

	return ok && id.Name == "false"
}

// blockJustifiesNotFound reports whether the statements of block before stmt
// (inclusive of the enclosing if-condition) contain a 404 test or an error
// diagnostic call.
func blockJustifiesNotFound(cond ast.Expr, stmts []ast.Stmt, upto ast.Stmt) bool {
	found := false
	check := func(n ast.Node) bool {
		if found || n == nil {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectorExpr:
			switch x.Sel.Name {
			case "IsNotFound", "StatusNotFound", "AddError", "Append", "HasError":
				found = true
			}
		}

		return true
	}
	if cond != nil {
		ast.Inspect(cond, check)
	}
	for _, s := range stmts {
		ast.Inspect(s, check)
		if s == upto {
			break
		}
	}

	return found
}

// falseReturnsUnjustified lists the positions of `return …false…` statements in
// fn that no 404 test or error diagnostic justifies.
func falseReturnsUnjustified(fn *ast.FuncDecl, fset *token.FileSet, idx int) []string {
	var out []string
	// Walk blocks so that each return knows its enclosing block and the if
	// condition that guards it.
	var walk func(block *ast.BlockStmt, cond ast.Expr, inherited bool)
	walk = func(block *ast.BlockStmt, cond ast.Expr, inherited bool) {
		for _, st := range block.List {
			switch s := st.(type) {
			case *ast.ReturnStmt:
				if len(s.Results) > idx && isFalse(s.Results[idx]) {
					if !inherited && !blockJustifiesNotFound(cond, block.List, st) {
						out = append(out, fset.Position(s.Pos()).String())
					}
				}
			case *ast.IfStmt:
				just := inherited || blockJustifiesNotFound(s.Cond, nil, nil)
				walk(s.Body, s.Cond, just)
				if els, ok := s.Else.(*ast.BlockStmt); ok {
					walk(els, nil, inherited)
				} else if els, ok := s.Else.(*ast.IfStmt); ok {
					walk(&ast.BlockStmt{List: []ast.Stmt{els}}, nil, inherited)
				}
			case *ast.BlockStmt:
				walk(s, nil, inherited)
			case *ast.ForStmt:
				walk(s.Body, nil, inherited)
			case *ast.RangeStmt:
				walk(s.Body, nil, inherited)
			case *ast.SwitchStmt:
				for _, c := range s.Body.List {
					if cc, ok := c.(*ast.CaseClause); ok {
						walk(&ast.BlockStmt{List: cc.Body}, nil, inherited)
					}
				}
			}
		}
	}
	walk(fn.Body, nil, false)

	return out
}

// readReachable returns the functions reachable from fn by name in the package.
func readReachable(fn *ast.FuncDecl, index map[string][]*ast.FuncDecl, seen map[*ast.FuncDecl]bool) {
	if seen[fn] {
		return
	}
	seen[fn] = true
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		var name string
		switch f := call.Fun.(type) {
		case *ast.SelectorExpr:
			name = f.Sel.Name
		case *ast.Ident:
			name = f.Name
		}
		for _, callee := range index[name] {
			readReachable(callee, index, seen)
		}

		return true
	})
}

func TestUnitReadHelpersReturnNotFoundOnlyFor404OrError(t *testing.T) {
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
			t.Fatal(err)
		}
		var files []*ast.File
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") ||
				e.Name() == "schema_gen.go" {
				continue
			}
			file, err := parser.ParseFile(fset, filepath.Join(root, d.Name(), e.Name()), nil, 0)
			if err != nil {
				t.Fatal(err)
			}
			files = append(files, file)
		}
		index := funcsByName(files)

		reach := map[*ast.FuncDecl]bool{}
		for _, fns := range index {
			for _, fn := range fns {
				if isReadMethod(fn) {
					readReachable(fn, index, reach)
				}
			}
		}
		for fn := range reach {
			if isReadMethod(fn) {
				continue
			}
			idx, inverse := boolResultIndex(fn)
			if idx < 0 || inverse {
				continue
			}
			checked++
			for _, pos := range falseReturnsUnjustified(fn, fset, idx) {
				violations = append(violations, fmt.Sprintf("%s.%s: returns found=false on a path with no 404 test or error diagnostic (%s); a healthy resource would be removed from state", d.Name(), fn.Name.Name, pos))
			}
		}
	}

	if checked == 0 {
		t.Fatal("found no found-bool helpers reachable from a Read; the rule's pattern no longer matches the code")
	}
	sort.Strings(violations)
	for _, v := range violations {
		t.Error(v)
	}
	t.Logf("checked %d found-bool helpers reachable from Read", checked)
}
