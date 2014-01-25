package dub

import (
	"bytes"
	"evergreen/base"
	"evergreen/dub/tree"
	dst "evergreen/go/tree"
	"fmt"
)

func id(name string) dst.Expr {
	return &dst.NameRef{Text: name}
}

func attr(expr dst.Expr, name string) dst.Expr {
	return &dst.Selector{Expr: expr, Text: name}
}

func strLiteral(value string) dst.Expr {
	return &dst.StringLiteral{Value: value}
}

func runeLiteral(value rune) dst.Expr {
	return &dst.RuneLiteral{Value: value}
}

func intLiteral(value int) dst.Expr {
	return &dst.IntLiteral{Value: value}
}

func boolLiteral(value bool) dst.Expr {
	return &dst.BoolLiteral{Value: value}
}

func makeFatalTest(cond dst.Expr, f string, args ...dst.Expr) dst.Stmt {
	wrapped := []dst.Expr{strLiteral(f)}
	wrapped = append(wrapped, args...)
	return &dst.If{
		Cond: cond,
		Body: []dst.Stmt{
			&dst.Call{
				Expr: attr(&dst.NameRef{Text: "t"}, "Fatalf"),
				Args: wrapped,
			},
		},
	}
}

func makeLen(expr dst.Expr) dst.Expr {
	return &dst.Call{
		Expr: id("len"),
		Args: []dst.Expr{
			expr,
		},
	}
}

func checkEQ(x dst.Expr, y dst.Expr) dst.Expr {
	return &dst.BinaryExpr{
		Left:  x,
		Op:    "==",
		Right: y,
	}
}

func checkNE(x dst.Expr, y dst.Expr) dst.Expr {
	return &dst.BinaryExpr{
		Left:  x,
		Op:    "!=",
		Right: y,
	}
}

func generateDestructure(name string, path string, d tree.Destructure, general tree.ASTType, gbuilder *GlobalDubBuilder, stmts []dst.Stmt) []dst.Stmt {
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
			stmts = append(stmts, &dst.Assign{
				Targets: []dst.Expr{
					id(actual_name),
					id("ok"),
				},
				Op: ":=",
				Sources: []dst.Expr{
					// TODO typecast tree.
					&dst.TypeAssert{
						Expr: id(name),
						Type: &dst.PointerType{Element: &dst.TypeRef{Name: d.Type.Name.Text}},
					},
				},
			})
			stmts = append(stmts, makeFatalTest(
				&dst.UnaryExpr{Op: "!", Expr: id("ok")},
				fmt.Sprintf("%s: expected a *%s but got a %%#v", path, dt.Name.Text),
				id(name),
			))
		}
		cond := checkEQ(id(actual_name), id("nil"))
		stmts = append(stmts, makeFatalTest(cond, fmt.Sprintf("%s: nil", path)))

		for _, arg := range d.Args {
			fn := arg.Name.Text
			childstmts := []dst.Stmt{}
			child_name := fmt.Sprintf("%s_%s", name, fn)
			child_path := fmt.Sprintf("%s.%s", path, fn)
			childstmts = append(childstmts, &dst.Assign{
				Targets: []dst.Expr{
					id(child_name),
				},
				Op: ":=",
				Sources: []dst.Expr{
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
			stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
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
			childstmts := []dst.Stmt{}
			child_name := fmt.Sprintf("%s_%d", name, i)
			child_path := fmt.Sprintf("%s[%d]", path, i)
			childstmts = append(childstmts, &dst.Assign{
				Targets: []dst.Expr{
					id(child_name),
				},
				Op: ":=",
				Sources: []dst.Expr{
					&dst.Index{
						Expr:  id(name),
						Index: intLiteral(i),
					},
				},
			})
			childstmts = generateDestructure(child_name, child_path, arg, dt.Type, gbuilder, childstmts)
			stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
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

func generateExpr(state string, expr tree.ASTExpr) dst.Expr {
	switch expr := expr.(type) {
	case *tree.Call:
		args := []dst.Expr{
			id(state),
		}
		for _, arg := range expr.Args {
			args = append(args, generateExpr(state, arg))
		}
		return &dst.Call{
			Expr: id(expr.Name.Text),
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

func generateGoTest(tst *tree.Test, gbuilder *GlobalDubBuilder) *dst.FuncDecl {
	stmts := []dst.Stmt{}

	state := "state"
	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Expr{
			id(state),
		},
		Op: ":=",
		Sources: []dst.Expr{
			&dst.Call{
				Expr: attr(id("runtime"), "MakeState"),
				Args: []dst.Expr{
					strLiteral(tst.Input),
				},
			},
		},
	})

	root := "o"

	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Expr{
			id(root),
		},
		Op: ":=",
		Sources: []dst.Expr{
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

	return &dst.FuncDecl{
		Name: fmt.Sprintf("Test_%s", tst.Name.Text),
		Params: []*dst.Param{
			&dst.Param{
				Name: "t",
				Type: &dst.PointerType{Element: &dst.TypeRef{Name: "testing.T"}},
			},
		},
		Returns: []*dst.Param{},
		Body:    stmts,
	}
}

func GenerateTests(module string, tests []*tree.Test, gbuilder *GlobalDubBuilder) string {
	imports := []*dst.Import{
		&dst.Import{Path: "evergreen/dub/runtime"},
		&dst.Import{Path: "testing"},
	}

	decls := []dst.Decl{}

	for _, tst := range tests {
		decls = append(decls, generateGoTest(tst, gbuilder))
	}

	file := &dst.File{
		Package: "tree",
		Imports: imports,
		Decls:   decls,
	}

	b := &bytes.Buffer{}
	w := &base.CodeWriter{Out: b}
	dst.GenerateFile(file, w)
	return b.String()

}
