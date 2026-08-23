package workspace

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
)

// FieldShape is one struct field's name and rendered type expression.
type FieldShape struct {
	Name string
	Type string
}

// TypeShape is one top-level type declaration: its name, the kind of
// declaration ("struct", "interface" or the rendered underlying type for
// anything else) and its fields when it is a struct.
type TypeShape struct {
	Name   string
	Kind   string
	Fields []FieldShape
}

// GoDeclarations lists the top-level identifiers AST-visible in one Go
// source file, resolved by parsing rather than regex (Decision 1;
// requirement R6): structural criteria can be checked against what the
// learner actually declared without executing their code.
type GoDeclarations struct {
	Funcs []string
	Types []TypeShape
}

// HasFunc reports whether name is declared as a top-level function
// (methods are not top-level; PROJECT.md structural criteria target
// declarations, not implementations).
func (d GoDeclarations) HasFunc(name string) bool {
	return slices.Contains(d.Funcs, name)
}

// Type returns the declaration named name, if any.
func (d GoDeclarations) Type(name string) (TypeShape, bool) {
	for _, t := range d.Types {
		if t.Name == name {
			return t, true
		}
	}
	return TypeShape{}, false
}

// HasField reports whether typeName declares a field named fieldName.
func (d GoDeclarations) HasField(typeName, fieldName string) bool {
	t, ok := d.Type(typeName)
	if !ok {
		return false
	}
	for _, f := range t.Fields {
		if f.Name == fieldName {
			return true
		}
	}
	return false
}

// ParseGoFile parses the Go source at path and lists its top-level
// declarations. It never type-checks and never resolves imports, so it
// cannot execute or load learner code (constraint: read-only).
func ParseGoFile(path string) (GoDeclarations, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return GoDeclarations{}, err
	}
	var decl GoDeclarations
	for _, d := range f.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			if v.Recv == nil {
				decl.Funcs = append(decl.Funcs, v.Name.Name)
			}
		case *ast.GenDecl:
			if v.Tok != token.TYPE {
				continue
			}
			for _, spec := range v.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				decl.Types = append(decl.Types, typeShapeOf(ts))
			}
		}
	}
	return decl, nil
}

func typeShapeOf(ts *ast.TypeSpec) TypeShape {
	shape := TypeShape{Name: ts.Name.Name}
	switch t := ts.Type.(type) {
	case *ast.StructType:
		shape.Kind = "struct"
		shape.Fields = fieldsOf(t)
	case *ast.InterfaceType:
		shape.Kind = "interface"
	default:
		shape.Kind = renderExpr(ts.Type)
	}
	return shape
}

func fieldsOf(t *ast.StructType) []FieldShape {
	var fields []FieldShape
	if t.Fields == nil {
		return fields
	}
	for _, f := range t.Fields.List {
		typ := renderExpr(f.Type)
		if len(f.Names) == 0 {
			// Embedded field: the type name is also the field name.
			fields = append(fields, FieldShape{Name: typ, Type: typ})
			continue
		}
		for _, n := range f.Names {
			fields = append(fields, FieldShape{Name: n.Name, Type: typ})
		}
	}
	return fields
}

// renderExpr renders a type expression back to source-like text without
// importing go/printer, covering the shapes learner code realistically
// uses at this level (identifiers, selectors, pointers, slices, maps).
func renderExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + renderExpr(e.X)
	case *ast.SelectorExpr:
		return renderExpr(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + renderExpr(e.Elt)
		}
		return "[...]" + renderExpr(e.Elt)
	case *ast.MapType:
		return "map[" + renderExpr(e.Key) + "]" + renderExpr(e.Value)
	default:
		return "?"
	}
}
