package dub

import (
	"bytes"
	"evergreen/dub/tree"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
)

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

func generateDestructure(name string, path string, d tree.Destructure, general tree.ASTType, gbuilder *GlobalDubBuilder, stmts []ast.Stmt) []ast.Stmt {
	switch d := d.(type) {
	case *tree.DestructureStruct:
		actual_name := name

		t := tree.ResolveType(d.Type)
		dt, ok := t.(*tree.StructDecl)
		if !ok {
			panic(t)
		}

		at := gbuilder.TranslateType(t)
		gt := gbuilder.TranslateType(general)

		if gt != at {
			actual_name = fmt.Sprintf("typed_%s", name)
			stmts = append(stmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					id(actual_name),
					id("ok"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					attr(id(name), fmt.Sprintf("(*%s)", d.Type.Name.Text)),
				},
			})
			stmts = append(stmts, makeFatalTest(
				&ast.UnaryExpr{Op: token.NOT, X: id("ok")},
				fmt.Sprintf("%s: expected a *%s but got a %%#v", path, dt.Name.Text),
				id(name),
			))
		}
		cond := checkEQ(id(actual_name), id("nil"))
		stmts = append(stmts, makeFatalTest(cond, fmt.Sprintf("%s: nil", path)))

		for _, arg := range d.Args {
			fn := arg.Name.Text
			childstmts := []ast.Stmt{}
			child_name := fmt.Sprintf("%s_%s", name, fn)
			child_path := fmt.Sprintf("%s.%s", path, fn)
			childstmts = append(childstmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					id(child_name),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					attr(id(actual_name), fn),
				},
			})
			f := tree.GetField(dt, fn)
			childstmts = generateDestructure(
				child_name,
				child_path,
				arg.Destructure,
				tree.ResolveType(f.Type),
				gbuilder,
				childstmts,
			)
			stmts = append(stmts, &ast.BlockStmt{List: childstmts})
		}
	case *tree.DestructureList:
		stmts = append(stmts, makeFatalTest(
			checkNE(makeLen(id(name)), intLiteral(len(d.Args))),
			fmt.Sprintf("%s: expected length %d but got %%d", path, len(d.Args)),
			makeLen(id(name)),
		))
		t := tree.ResolveType(d.Type)
		dt, ok := t.(*tree.ListType)
		if !ok {
			panic(t)
		}
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
			childstmts = generateDestructure(child_name, child_path, arg, dt.Type, gbuilder, childstmts)
			stmts = append(stmts, &ast.BlockStmt{List: childstmts})
		}
	case *tree.DestructureValue:
		switch expr := d.Expr.(type) {
		case *tree.StringLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), strLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), strLiteral(expr.Value), id(name)))
		case *tree.RuneLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), runeLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#U but got %%#U", path), runeLiteral(expr.Value), id(name)))
		case *tree.IntLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), intLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), intLiteral(expr.Value), id(name)))
		case *tree.BoolLiteral:
			stmts = append(stmts, makeFatalTest(checkNE(id(name), boolLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), boolLiteral(expr.Value), id(name)))
		default:
			panic(expr)
		}
	default:
		panic(d)
	}

	return stmts
}

func generateExpr(state string, expr tree.ASTExpr) ast.Expr {
	switch expr := expr.(type) {
	case *tree.Call:
		args := []ast.Expr{
			id(state),
		}
		for _, arg := range expr.Args {
			args = append(args, generateExpr(state, arg))
		}
		return &ast.CallExpr{
			Fun:  id(expr.Name.Text),
			Args: args,
		}
	case *tree.StringLiteral:
		return strLiteral(expr.Value)
	case *tree.RuneLiteral:
		return runeLiteral(expr.Value)
	case *tree.IntLiteral:
		return intLiteral(expr.Value)
	case *tree.BoolLiteral:
		return boolLiteral(expr.Value)
	default:
		panic(expr)
	}

}

func generateGoTest(tst *tree.Test, gbuilder *GlobalDubBuilder) *ast.FuncDecl {
	stmts := []ast.Stmt{}

	state := "state"
	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			id(state),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: attr(id("runtime"), "MakeState"),
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
			generateExpr(state, tst.Rule),
		},
	})

	stmts = append(stmts, makeFatalTest(
		checkNE(attr(id(state), "Index"), intLiteral(len(tst.Input))),
		fmt.Sprintf("Only consumed %%d/%d (deepest %%d) runes", len(tst.Input)),
		attr(id(state), "Index"),
		attr(id(state), "Deepest"),
	))

	stmts = append(stmts, makeFatalTest(
		checkNE(attr(id(state), "Flow"), intLiteral(0)),
		"Expected flow to be 0, but got %d",
		attr(id(state), "Flow"),
	))

	stmts = generateDestructure(root, root, tst.Destructure, tst.Type, gbuilder, stmts)

	return &ast.FuncDecl{
		Name: id(fmt.Sprintf("Test_%s", tst.Name.Text)),
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

func GenerateTests(module string, tests []*tree.Test, gbuilder *GlobalDubBuilder) string {
	decls := []ast.Decl{}
	decls = append([]ast.Decl{&ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: strLiteral("evergreen/dub/runtime")},
			&ast.ImportSpec{Path: strLiteral("testing")},
		},
	}}, decls...)

	// TODO
	for _, tst := range tests {
		decls = append(decls, generateGoTest(tst, gbuilder))
	}

	file := &ast.File{
		Name:  id("tree"),
		Decls: decls,
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, file)

	return buf.String()

}
