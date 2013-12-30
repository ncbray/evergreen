package dasm

import (
	"bytes"
	"evergreen/dub"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
)

// Utility functions for generated tests
func MakeState(input string) *dub.DubState {
	return &dub.DubState{Stream: []rune(input)}
}

func id(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func singleName(name string) []*ast.Ident {
	return []*ast.Ident{id(name)}
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

func strLiteral(name string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(name),
	}
}

func runeLiteral(value rune) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.CHAR,
		Value: strconv.QuoteRune(value),
	}
}

func intLiteral(value int) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(int64(value), 10),
	}
}

func boolLiteral(value bool) ast.Expr {
	return id(fmt.Sprintf("%v", value))
}

func makeFatalTest(cond ast.Expr, f string, args ...ast.Expr) ast.Stmt {
	wrapped := []ast.Expr{strLiteral(f)}
	wrapped = append(wrapped, args...)
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun:  attr(id("t"), "Fatalf"),
						Args: wrapped,
					},
				},
			},
		},
	}
}

func makeLen(expr ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: id("len"),
		Args: []ast.Expr{
			expr,
		},
	}
}

func checkEQ(x ast.Expr, y ast.Expr) ast.Expr {
	return &ast.BinaryExpr{
		X:  x,
		Op: token.EQL,
		Y:  y,
	}
}

func checkNE(x ast.Expr, y ast.Expr) ast.Expr {
	return &ast.BinaryExpr{
		X:  x,
		Op: token.NEQ,
		Y:  y,
	}
}

func generateDestructure(name string, path string, d Destructure, stmts []ast.Stmt) []ast.Stmt {
	switch d := d.(type) {
	case *DestructureStruct:
		actual_name := name
		if d.GT != d.AT {
			actual_name = fmt.Sprintf("typed_%s", name)
			stmts = append(stmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					id(actual_name),
					id("ok"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					attr(id(name), fmt.Sprintf("(*%s)", d.Type.Name)),
				},
			})
			stmts = append(stmts, makeFatalTest(
				&ast.UnaryExpr{Op: token.NOT, X: id("ok")},
				fmt.Sprintf("%s: expected a *%s but got a %%#v", path, d.AT.Name),
				id(name),
			))
		}
		cond := checkEQ(id(actual_name), id("nil"))
		stmts = append(stmts, makeFatalTest(cond, fmt.Sprintf("%s: nil", path)))

		for _, arg := range d.Args {
			childstmts := []ast.Stmt{}
			child_name := fmt.Sprintf("%s_%s", name, arg.Name)
			child_path := fmt.Sprintf("%s.%s", path, arg.Name)
			childstmts = append(childstmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					id(child_name),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					attr(id(actual_name), arg.Name),
				},
			})
			childstmts = generateDestructure(child_name, child_path, arg.Destructure, childstmts)
			stmts = append(stmts, &ast.BlockStmt{List: childstmts})
		}
	case *DestructureList:
		stmts = append(stmts, makeFatalTest(
			checkNE(makeLen(id(name)), intLiteral(len(d.Args))),
			fmt.Sprintf("%s: expected length %d but got %%d", path, len(d.Args)),
			makeLen(id(name)),
		))
		for i, arg := range d.Args {
			childstmts := []ast.Stmt{}
			child_name := fmt.Sprintf("%s_%d", name, i)
			child_path := fmt.Sprintf("%s[%d]", path, i)
			childstmts = append(childstmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					id(child_name),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.IndexExpr{
						X:     id(name),
						Index: intLiteral(i),
					},
				},
			})
			childstmts = generateDestructure(child_name, child_path, arg, childstmts)
			stmts = append(stmts, &ast.BlockStmt{List: childstmts})
		}
	case *DestructureValue:
		switch expr := d.Expr.(type) {
		case *StringLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), strLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), strLiteral(expr.Value), id(name)))
		case *RuneLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), runeLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#U but got %%#U", path), runeLiteral(expr.Value), id(name)))
		case *IntLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), intLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), intLiteral(expr.Value), id(name)))
		case *BoolLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), boolLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), boolLiteral(expr.Value), id(name)))
		default:
			panic(expr)
		}
	default:
		panic(d)
	}

	return stmts
}

func generateGoTest(tst *Test) *ast.FuncDecl {
	stmts := []ast.Stmt{}

	state := "state"
	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			id(state),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: attr(id("dasm"), "MakeState"),
				Args: []ast.Expr{
					strLiteral(tst.Input),
				},
			},
		},
	})

	root := "o"

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			id(root),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: id(tst.Rule),
				Args: []ast.Expr{
					id(state),
				},
			},
		},
	})

	stmts = append(stmts, makeFatalTest(
		checkNE(attr(id(state), "Index"), intLiteral(len(tst.Input))),
		fmt.Sprintf("Only consumed %%d/%d runes", len(tst.Input)),
		attr(id(state), "Index"),
	))

	stmts = append(stmts, makeFatalTest(
		checkNE(attr(id(state), "Flow"), intLiteral(0)),
		"Expected flow to be 0, but got %d",
		attr(id(state), "Flow"),
	))

	stmts = generateDestructure(root, root, tst.Destructure, stmts)

	return &ast.FuncDecl{
		Name: id(fmt.Sprintf("Test_%s_%s", tst.Rule, tst.Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: singleName("t"),
						Type:  ptr(attr(id("testing"), "T")),
					},
				},
			},
			Results: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{List: stmts},
	}
}

func GenerateTests(module string, tests []*Test) string {
	decls := []ast.Decl{}
	decls = append([]ast.Decl{&ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: strLiteral("evergreen/dasm")},
			&ast.ImportSpec{Path: strLiteral("testing")},
		},
	}}, decls...)

	// TODO
	for _, tst := range tests {
		decls = append(decls, generateGoTest(tst))
	}

	file := &ast.File{
		Name:  id(module),
		Decls: decls,
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, file)

	return buf.String()

}
