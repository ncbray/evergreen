package dub

import (
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	core "evergreen/dub/tree"
	dst "evergreen/go/tree"
	"fmt"
)

type TestingContext struct {
	link     flow.DubToGoLinker
	t        *dst.StructDecl
	stateT   *dst.StructDecl
	state    int
	funcDecl *dst.FuncDecl
	tInfo    int
	okInfo   int
	index    *flow.BuiltinIndex
}

func (ctx *TestingContext) GetState() dst.Expr {
	return &dst.GetLocal{Info: ctx.state}
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

func translateType(ctx *TestingContext, at core.DubType) dst.Type {
	switch cat := at.(type) {
	case *core.StructType:
		ref := ctx.link.TypeRef(cat, flow.STRUCT)
		if cat.IsParent {
			return ref
		} else {
			return &dst.PointerType{Element: ref}
		}
	case *core.ListType:
		return &dst.SliceType{Element: translateType(ctx, cat.Type)}
	case *core.BuiltinType:
		switch cat.Name {
		case "string":
			return &dst.TypeRef{Impl: ctx.index.StringType}
		case "rune":
			return &dst.TypeRef{Impl: ctx.index.RuneType}
		case "int":
			return &dst.TypeRef{Impl: ctx.index.IntType}
		case "bool":
			return &dst.TypeRef{Impl: ctx.index.BoolType}
		default:
			panic(cat.Name)
		}
	default:
		panic(at)
	}
}

func generateDestructure(value int, nameX string, path string, d tree.Destructure, generalType core.DubType, ctx *TestingContext, stmts []dst.Stmt) []dst.Stmt {
	switch d := d.(type) {
	case *tree.DestructureStruct:
		actual_value := value

		actualType := tree.ResolveType(d.Type)
		structType, ok := actualType.(*core.StructType)
		if !ok {
			panic(actualType)
		}

		if generalType != actualType {
			actual_name := fmt.Sprintf("typed_%s", nameX)

			lref := translateType(ctx, actualType)
			actual_value = ctx.funcDecl.CreateLocalInfo(actual_name, lref)

			ref := translateType(ctx, actualType)
			stmts = append(stmts, &dst.Assign{
				Targets: []dst.Target{
					// HACK
					ctx.funcDecl.MakeSetLocal(actual_value),
					ctx.funcDecl.MakeSetLocal(ctx.okInfo),
				},
				Op: "=",
				Sources: []dst.Expr{
					// TODO typecast tree.
					&dst.TypeAssert{
						Expr: ctx.funcDecl.MakeGetLocal(value),
						Type: ref,
					},
				},
			})
			stmts = append(stmts, ctx.makeFatalTest(
				&dst.UnaryExpr{Op: "!", Expr: ctx.funcDecl.MakeGetLocal(ctx.okInfo)},
				fmt.Sprintf("%s: expected a *%s but got a %%#v", path, structType.Name.Text),
				ctx.funcDecl.MakeGetLocal(value),
			))
		}
		cond := checkEQ(ctx.funcDecl.MakeGetLocal(actual_value), &dst.NilLiteral{})
		stmts = append(stmts, ctx.makeFatalTest(cond, fmt.Sprintf("%s: nil", path)))

		for _, arg := range d.Args {
			fn := arg.Name.Text
			childstmts := []dst.Stmt{}
			child_name := fmt.Sprintf("%s_%s", nameX, fn)
			child_path := fmt.Sprintf("%s.%s", path, fn)

			f := tree.GetField(structType, fn)
			t := f.Type

			child_value := ctx.funcDecl.CreateLocalInfo(child_name, translateType(ctx, t))
			childstmts = append(childstmts, &dst.Assign{
				Targets: []dst.Target{
					&dst.SetLocal{Info: child_value},
				},
				Op: "=",
				Sources: []dst.Expr{
					attr(ctx.funcDecl.MakeGetLocal(actual_value), fn),
				},
			})
			childstmts = generateDestructure(
				child_value,
				child_name,
				child_path,
				arg.Destructure,
				t,
				ctx,
				childstmts,
			)
			//stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
			stmts = append(stmts, childstmts...)
		}
	case *tree.DestructureList:
		stmts = append(stmts, ctx.makeFatalTest(
			checkNE(makeLen(ctx.funcDecl.MakeGetLocal(value)), intLiteral(len(d.Args))),
			fmt.Sprintf("%s: expected length %d but got %%d", path, len(d.Args)),
			makeLen(ctx.funcDecl.MakeGetLocal(value)),
		))
		t := tree.ResolveType(d.Type)
		dt, ok := t.(*tree.ListType)
		if !ok {
			panic(t)
		}
		for i, arg := range d.Args {
			childstmts := []dst.Stmt{}
			child_name := fmt.Sprintf("%s_%d", nameX, i)
			child_path := fmt.Sprintf("%s[%d]", path, i)
			child_value := ctx.funcDecl.CreateLocalInfo(child_name, translateType(ctx, dt.Type))
			childstmts = append(childstmts, &dst.Assign{
				Targets: []dst.Target{
					ctx.funcDecl.MakeSetLocal(child_value),
				},
				Op: "=",
				Sources: []dst.Expr{
					&dst.Index{
						Expr:  ctx.funcDecl.MakeGetLocal(value),
						Index: intLiteral(i),
					},
				},
			})
			childstmts = generateDestructure(child_value, child_name, child_path, arg, dt.Type, ctx, childstmts)
			//stmts = append(stmts, &dst.BlockStmt{Body: childstmts})
			stmts = append(stmts, childstmts...)
		}
	case *tree.DestructureValue:
		switch expr := d.Expr.(type) {
		case *tree.StringLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(ctx.funcDecl.MakeGetLocal(value), strLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), strLiteral(expr.Value), ctx.funcDecl.MakeGetLocal(value)))
		case *tree.RuneLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(ctx.funcDecl.MakeGetLocal(value), runeLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#U but got %%#U", path), runeLiteral(expr.Value), ctx.funcDecl.MakeGetLocal(value)))
		case *tree.IntLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(ctx.funcDecl.MakeGetLocal(value), intLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), intLiteral(expr.Value), ctx.funcDecl.MakeGetLocal(value)))
		case *tree.BoolLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(ctx.funcDecl.MakeGetLocal(value), boolLiteral(expr.Value)), fmt.Sprintf("%s: expected %%#v but got %%#v", path), boolLiteral(expr.Value), ctx.funcDecl.MakeGetLocal(value)))
		case *tree.NilLiteral:
			stmts = append(stmts, ctx.makeFatalTest(checkNE(ctx.funcDecl.MakeGetLocal(value), nilLiteral()), fmt.Sprintf("%s: expected nil but got %%#v", path), ctx.funcDecl.MakeGetLocal(value)))
		default:
			panic(expr)
		}
	default:
		panic(d)
	}

	return stmts
}

func generateExpr(ctx *TestingContext, expr tree.ASTExpr) dst.Expr {
	switch expr := expr.(type) {
	case *tree.Call:
		args := []dst.Expr{
			ctx.GetState(),
		}
		for _, arg := range expr.Args {
			args = append(args, generateExpr(ctx, arg))
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
		Name:            fmt.Sprintf("Test_%s", tst.Name.Text),
		LocalInfo_Scope: &dst.LocalInfo_Scope{},
	}

	// HACK
	ctx.funcDecl = decl
	ctx.tInfo = decl.CreateLocalInfo("t", &dst.PointerType{
		Element: &dst.TypeRef{
			Impl: ctx.t,
		},
	})
	ctx.state = decl.CreateLocalInfo("state", &dst.PointerType{
		Element: &dst.TypeRef{
			Impl: ctx.stateT,
		},
	})
	ctx.okInfo = decl.CreateLocalInfo("ok", &dst.TypeRef{
		// HACK should actual reference a type.
		Name: "bool",
		Impl: &dst.ExternalType{Name: "bool"},
	})

	stmts := []dst.Stmt{}
	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Target{
			&dst.SetLocal{
				Info: ctx.state,
			},
		},
		Op: "=",
		Sources: []dst.Expr{
			&dst.Call{
				Expr: attr(glbl("runtime"), "MakeState"),
				Args: []dst.Expr{
					strLiteral(tst.Input),
				},
			},
		},
	})

	root_name := "o"
	root_value := decl.CreateLocalInfo(root_name, translateType(ctx, tst.Type))

	stmts = append(stmts, &dst.Assign{
		Targets: []dst.Target{
			// HACK
			ctx.funcDecl.MakeSetLocal(root_value),
		},
		Op: "=",
		Sources: []dst.Expr{
			generateExpr(ctx, tst.Rule),
		},
	})

	flowName := tst.Flow

	// Runes consumed should only be checked if the call succeeds.
	if flowName == "NORMAL" {
		stmts = append(stmts, ctx.makeFatalTest(
			checkNE(attr(ctx.GetState(), "Index"), intLiteral(len(tst.Input))),
			fmt.Sprintf("Only consumed %%d/%d (deepest %%d) runes", len(tst.Input)),
			attr(ctx.GetState(), "Index"),
			attr(ctx.GetState(), "Deepest"),
		))
	}

	// Make sure the flow is what we expect.
	stmts = append(stmts, ctx.makeFatalTest(
		checkNE(attr(ctx.GetState(), "Flow"), attr(glbl("runtime"), flowName)),
		"Expected flow to be %d, but got %d",
		attr(glbl("runtime"), flowName), attr(ctx.GetState(), "Flow"),
	))

	stmts = generateDestructure(root_value, root_name, root_name, tst.Destructure, tst.Type, ctx, stmts)

	decl.Type = &dst.FuncType{
		Params: []*dst.Param{
			decl.MakeParam(ctx.tInfo),
		},
		Results: []*dst.Param{},
	}
	decl.Body = stmts
	return decl
}

func GenerateTests(leaf string, tests []*tree.Test, t *dst.StructDecl, stateT *dst.StructDecl, link flow.DubToGoLinker) *dst.File {
	_, index := flow.ExternBuiltinRuntime()
	ctx := &TestingContext{link: link, t: t, stateT: stateT, index: index}

	decls := []dst.Decl{}

	for _, tst := range tests {
		decls = append(decls, generateGoTest(tst, ctx))
	}

	file := &dst.File{
		Name:    "generated_dub_test.go",
		Package: leaf,
		Decls:   decls,
	}
	return file
}
