package dub

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
)

// Begin AST construction wrappers

func id(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func singleName(name string) []*ast.Ident {
	return []*ast.Ident{id(name)}
}

func strLiteral(name string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(name),
	}
}

func intLiteral(value int) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(int64(value), 10),
	}
}

func addr(expr ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  expr,
	}
}

func ptr(expr ast.Expr) ast.Expr {
	return &ast.StarExpr{X: expr}
}

func attr(expr ast.Expr, name string) ast.Expr {
	return &ast.SelectorExpr{X: expr, Sel: id(name)}
}

// End AST construction wrappers

func genFunc(decl *FuncDecl) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: id(decl.Name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: &ast.FieldList{List: []*ast.Field{&ast.Field{Type: ptr(id("bogus_type2"))}}},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{Lhs: []ast.Expr{id("block")}, Tok: token.DEFINE, Rhs: []ast.Expr{intLiteral(0)}},
				&ast.SwitchStmt{
					Tag: id("block"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.CaseClause{
								List: []ast.Expr{intLiteral(0)},
								Body: []ast.Stmt{
									&ast.ReturnStmt{Results: []ast.Expr{}},
								},
							},
						},
					},
				},
			},
		}}
}

func GenerateGo() string {
	data := []*FuncDecl{
		&FuncDecl{Name: "S"},
		&FuncDecl{Name: "integer"},
	}

	decls := []ast.Decl{}

	decls = append(decls, genFunc(data[0]))
	decls = append(decls, genFunc(data[1]))

	f := &ast.File{Name: id("dub"), Decls: decls}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)

	return buf.String()
}
