package dub

import (
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	dst "evergreen/go/tree"
	"fmt"
)

type TestingContext struct {
	gbuilder *GlobalDubBuilder
	link     flow.DubToGoLinker
	t        *dst.StructDecl
	funcDecl *dst.FuncDecl
	tInfo    int
}

func lcl(name string) dst.Expr {
	return &dst.GetName{
		Text: name,
	}
}

func glbl(name string) dst.Expr {
	return &dst.GetGlobal{
		Text: name,
	}
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

func nilLiteral() dst.Expr {
	return &dst.NilLiteral{}
}

func (ctx *TestingContext) makeFatalTest(cond dst.Expr, f string, args ...dst.Expr) dst.Stmt {
	wrapped := []dst.Expr{strLiteral(f)}
	wrapped = append(wrapped, args...)
	return &dst.If{
		Cond: cond,
		Body: []dst.Stmt{
			&dst.Call{
				Expr: attr(ctx.funcDecl.MakeGetLocal(ctx.tInfo), "Fatalf"),
				Args: wrapped,
			},
		},
	}
}

func makeLen(expr dst.Expr) dst.Expr {
	return &dst.Call{
		Expr: glbl("len"),
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

func generateDestructure(name string, path string, d tree.Destructure, general tree.ASTType, ctx *TestingContext, stmts []dst.Stmt) []dst.Stmt {
	switch d := d.(type) {
	case *tree.DestructureStruct:
		actual_name := name

		t := tree.ResolveType(d.Type)
		dt, ok := t.(*tree.StructDecl)
		if !ok {
			panic(t)
		}

		at := ctx.gbuilder.TranslateType(t)
		gt := ctx.gbuilder.TranslateType(general)

		cat, ok := at.(*flow.LLStruct)
		if !ok {
			panic(at)
		}

		if gt != at {
			actual_name = fmt.Sprintf("typed_%s", name)
			ref := &dst.TypeRef{Name: cat.Name}
			ctx.link.TypeRef(ref, cat)
			stmts = append(stmts, &dst.Assign{
				Targets: []dst.Target{
					// HACK
					&dst.SetName{Text: actual_name},
					&dst.SetName{Text: "ok"},
				},
				Op: ":=",
				Sources: []dst.Expr{
					// TODO typecast tree.
					&dst.TypeAssert{
						Expr: lcl(name),
						Type: &dst.PointerType{Element: ref},
					},
				},
			})
			stmts = append(stmts, ctx.makeFatalTest(
				&dst.UnaryExpr{Op: "!", Expr: lcl("ok")},
				fmt.Sprintf("%s: expected a *%s but got a %%#v", path, dt.Name.Text),
				lcl(name),
			))
		}
		cond := checkEQ(lcl(actual_name), &dst.NilLiteral{})
		stmts = append(stmts, ctx.makeFatalTest(cond, fmt.Sprintf("%s: nil", path)))

		for _, arg := range d.Args {
			fn := arg.Name.Text
			childstmts := []dst.Stmt{}
			child_name := fmt.Sprintf("%s_%s", name, fn)
			child_path := fmt.Sprintf("%s.%s", path, fn)
			childstmts = append(childstmts, &dst.Assign{
				Targets: []dst.Target{
					&dst.SetName{Text: child_name},
				},
				Op: ":=",
				Sources: []dst.Expr{
					attr(lcl(actual_name), fn),
				},
			})
			f := tree.GetField(dt, fn)
			childstmts = generateDestructure(
				child_name,
				child_path,
				arg.Destructure,
				tree.ResolveType(f.Type),
				ctx,
				childstmts,
			)
			stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
		}
	case *tree.DestructureList:
		stmts = append(stmts, ctx.makeFatalTest(
			checkNE(makeLen(lcl(name)), intLiteral(len(d.Args))),
			fmt.Sprintf("%s: expected length %d but got %%d", path, len(d.Args)),
			makeLen(lcl(name)),
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
				Targets: []dst.Target{
					&dst.SetName{Text: child_name},
				},
				Op: ":=",
				Sources: []dst.Expr{
					&dst.Index{
						Expr:  lcl(name),
						Index: intLiteral(i),
					},
				},
			})
			childstmts = generateDestructure(child_name, child_path, arg, dt.Type, ctx, childstmts)
			stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
		}
	case *tree.DestructureValue:
		switch expr := d.Expr.(type) {
		case *tree.StringLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(lcl(name), strLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), strLiteral(expr.Value), lcl(name)))
		case *tree.RuneLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(lcl(name), runeLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#U but got %%#U", path), runeLiteral(expr.Value), lcl(name)))
		case *tree.IntLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(lcl(name), intLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), intLiteral(expr.Value), lcl(name)))
		case *tree.BoolLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(lcl(name), boolLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), boolLiteral(expr.Value), lcl(name)))
		case *tree.NilLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(lcl(name), nilLiteral()), fmt.Sprintf("%s: expected nil but got %%#v", path), lcl(name)))
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
			lcl(state),
		}
		for _, arg := range expr.Args {
			args = append(args, generateExpr(state, arg))
		}
		return &dst.Call{
			Expr: glbl(expr.Name.Text),
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

func generateGoTest(tst *tree.Test, ctx *TestingContext) *dst.FuncDecl {
	decl := &dst.FuncDecl{
		Name: fmt.Sprintf("Test_%s", tst.Name.Text),
	}
	tInfo := decl.CreateLocalInfo("t", &dst.PointerType{Element: &dst.TypeRef{Impl: ctx.t}})

	// HACK
	ctx.funcDecl = decl
	ctx.tInfo = tInfo

	stmts := []dst.Stmt{}

	state := "state"
	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Target{
			// HACK
			&dst.SetName{Text: state},
		},
		Op: ":=",
		Sources: []dst.Expr{
			&dst.Call{
				Expr: attr(glbl("runtime"), "MakeState"),
				Args: []dst.Expr{
					strLiteral(tst.Input),
				},
			},
		},
	})

	root := "o"

	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Target{
			// HACK
			&dst.SetName{Text: root},
		},
		Op: ":=",
		Sources: []dst.Expr{
			generateExpr(state, tst.Rule),
		},
	})

	flowName := tst.Flow

	// Runes consumed should only be checked if the call succeeds.
	if flowName == "NORMAL" {
		stmts = append(stmts, ctx.makeFatalTest(
			checkNE(attr(lcl(state), "Index"), intLiteral(len(tst.Input))),
			fmt.Sprintf("Only consumed %%d/%d (deepest %%d) runes", len(tst.Input)),
			attr(lcl(state), "Index"),
			attr(lcl(state), "Deepest"),
		))
	}

	// Make sure the flow is what we expect.
	stmts = append(stmts, ctx.makeFatalTest(
		checkNE(attr(lcl(state), "Flow"), attr(glbl("runtime"), flowName)),
		"Expected flow to be %d, but got %d",
		attr(glbl("runtime"), flowName), attr(lcl(state), "Flow"),
	))

	stmts = generateDestructure(root, root, tst.Destructure, tst.Type, ctx, stmts)

	decl.Type = &dst.FuncType{
		Params: []*dst.Param{
			decl.MakeParam(tInfo),
		},
		Results: []*dst.Param{},
	}
	decl.Body = stmts
	return decl
}

func ExternTestingRuntime() (*dst.Package, *dst.StructDecl) {
	t := &dst.StructDecl{
		Name: "T",
	}
	pkg := &dst.Package{
		Extern: true,
		Path:   []string{"testing"},
		Files: []*dst.File{
			&dst.File{
				Decls: []dst.Decl{
					t,
				},
			},
		},
	}
	return pkg, t
}

func GenerateTests(module string, tests []*tree.Test, gbuilder *GlobalDubBuilder, t *dst.StructDecl, link flow.DubToGoLinker) *dst.File {
	ctx := &TestingContext{gbuilder: gbuilder, link: link, t: t}

	imports := []*dst.Import{
		// HACK for runtime.MakeState
		&dst.Import{Path: "evergreen/dub/runtime"},
	}

	decls := []dst.Decl{}

	for _, tst := range tests {
		decls = append(decls, generateGoTest(tst, ctx))
	}

	file := &dst.File{
		Name:    "generated_parser_test.go",
		Package: "tree",
		Imports: imports,
		Decls:   decls,
	}
	return file
}
