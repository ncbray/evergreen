package flow

import (
	"bytes"
	"evergreen/base"
	"fmt"
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

type blockInfo struct {
	label string
}

func blockName(i int) string {
	return fmt.Sprintf("block%d", i)
}

func emitOp(name string, args ...ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: attr(id("frame"), name), Args: args}}
}

func constInt(v int64) ast.Expr {
	// Type converstion hack
	return intLiteral(int(v))
}

func reg(r DubRegister) ast.Expr {
	return id(RegisterName(r))
}

func returnVarName(i int) string {
	return fmt.Sprintf("ret%d", i)
}

func opAssign(expr ast.Expr, dst DubRegister) ast.Stmt {
	if dst != NoRegister {
		return &ast.AssignStmt{
			Lhs: []ast.Expr{reg(dst)},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{expr},
		}
	} else {
		return &ast.ExprStmt{X: expr}
	}
}

var opToTok = map[string]token.Token{
	"+":  token.ADD,
	"-":  token.SUB,
	"*":  token.MUL,
	"/":  token.QUO,
	"==": token.EQL,
	"!=": token.NEQ,
	"<":  token.LSS,
	">":  token.GTR,
	"<=": token.LEQ,
	">=": token.GEQ,
}

func goTypeName(t DubType) ast.Expr {
	switch t := t.(type) {
	case *BoolType:
		return id("bool")
	case *IntType:
		return id("int")
	case *RuneType:
		return id("rune")
	case *StringType:
		return id("string")
	case *ListType:
		return &ast.ArrayType{Elt: goTypeName(t.Type)}
	case *LLStruct:
		if t.Abstract {
			return id(t.Name)
		} else {
			return ptr(id(t.Name))
		}
	default:
		panic(t)
	}
}

func GenerateOp(f *LLFunc, op DubOp, block []ast.Stmt) []ast.Stmt {
	if IsNop(op) {
		return block
	}

	switch op := op.(type) {
	case *BinaryOp:
		tok, ok := opToTok[op.Op]
		if !ok {
			panic(op.Op)
		}
		block = append(block, opAssign(
			&ast.BinaryExpr{
				X:  reg(op.Left),
				Op: tok,
				Y:  reg(op.Right),
			},
			op.Dst,
		))
	case *CallOp:
		block = append(block, opAssign(
			&ast.CallExpr{
				Fun: id(op.Name),
				Args: []ast.Expr{
					id("frame"),
				},
			},
			op.Dst,
		))
	case *ConstructOp:
		elts := make([]ast.Expr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = &ast.KeyValueExpr{Key: id(arg.Key), Value: reg(arg.Value)}
		}
		block = append(block, opAssign(
			addr(&ast.CompositeLit{
				Type: id(op.Type.Name),
				Elts: elts,
			}),
			op.Dst,
		))
	case *ConstructListOp:
		elts := make([]ast.Expr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = reg(arg)
		}
		block = append(block, opAssign(
			&ast.CompositeLit{
				Type: goTypeName(op.Type),
				Elts: elts,
			},
			op.Dst,
		))
	case *CoerceOp:
		block = append(block, opAssign(
			&ast.CallExpr{
				Fun: goTypeName(op.T),
				Args: []ast.Expr{
					reg(op.Src),
				},
			},
			op.Dst,
		))
	case *ConstantNilOp:
		block = append(block, opAssign(
			id("nil"),
			op.Dst,
		))
	case *ConstantBoolOp:
		block = append(block, opAssign(
			id(fmt.Sprintf("%v", op.Value)),
			op.Dst,
		))
	case *ConstantIntOp:
		block = append(block, opAssign(
			constInt(op.Value),
			op.Dst,
		))
	case *ConstantRuneOp:
		block = append(block, opAssign(
			&ast.BasicLit{
				Kind:  token.CHAR,
				Value: strconv.QuoteRune(op.Value),
			},
			op.Dst,
		))
	case *ConstantStringOp:
		block = append(block, opAssign(
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(op.Value),
			},
			op.Dst,
		))
	case *Peek:
		block = append(block, opAssign(
			emitOp("Peek").X,
			op.Dst,
		))
	case *Consume:
		block = append(block,
			emitOp("Consume"),
		)
	case *AppendOp:
		block = append(block, opAssign(
			&ast.CallExpr{
				Fun: id("append"),
				Args: []ast.Expr{
					reg(op.List),
					reg(op.Value),
				},
			},
			op.Dst,
		))
	case *ReturnOp:
		if len(op.Exprs) != len(f.ReturnTypes) {
			panic(fmt.Sprintf("Wrong number of return values.  Expected %d, got %d.", len(f.ReturnTypes), len(op.Exprs)))
		}
		for i, e := range op.Exprs {
			block = append(block, &ast.AssignStmt{
				Lhs: []ast.Expr{id(returnVarName(i))},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{reg(e)},
			})
		}
	case *Fail:
		block = append(block, emitOp("Fail"))
	case *Checkpoint:
		block = append(block, opAssign(
			emitOp("Checkpoint").X,
			op.Dst,
		))
	case *Recover:
		block = append(block, emitOp("Recover", reg(op.Src)))
	case *LookaheadBegin:
		block = append(block, opAssign(
			emitOp("LookaheadBegin").X,
			op.Dst,
		))
	case *LookaheadEnd:
		if op.Failed {
			block = append(block, emitOp("LookaheadFail", reg(op.Src)))
		} else {
			block = append(block, emitOp("LookaheadNormal", reg(op.Src)))
		}
	case *Slice:
		block = append(block, opAssign(
			emitOp("Slice", reg(op.Src)).X,
			op.Dst,
		))
	case *CopyOp:
		block = append(block, opAssign(
			reg(op.Src),
			op.Dst,
		))
	default:
		panic(op)
	}
	return block
}

func GenerateGoFunc(f *LLFunc) ast.Decl {
	nodes := base.ReversePostorder(f.Region)
	//info := make([]blockInfo, len(nodes))
	m := map[*base.Node]int{}
	stmts := []ast.Stmt{}

	// Map node -> ID.
	for i, node := range nodes {
		m[node] = i
	}

	gotoNode := func(e *base.Edge, n *base.Node, block []ast.Stmt) []ast.Stmt {
		if e != nil {
			transfers, ok := e.Data.([]Transfer)
			if ok {
				clobbered := map[DubRegister]bool{}
				for _, t := range transfers {
					_, is_clobbered := clobbered[t.Src]
					if is_clobbered {
						// TODO handle conflicitng transfers.
						panic(t.Src)
					}
					if t.Src != t.Dst {
						block = append(block, opAssign(
							reg(t.Src),
							t.Dst,
						))
					}
					clobbered[t.Dst] = true
				}
			}
		}
		return append(block, &ast.BranchStmt{Tok: token.GOTO, Label: id(blockName(m[n]))})
	}

	emitSwitch := func(cond ast.Expr, te *base.Edge, t *base.Node, fe *base.Edge, f *base.Node, block []ast.Stmt) []ast.Stmt {
		if t != nil {
			if f != nil {
				block = append(block, &ast.IfStmt{
					Cond: cond,
					Body: &ast.BlockStmt{
						List: gotoNode(te, t, nil),
					},
					Else: &ast.BlockStmt{
						List: gotoNode(fe, f, nil),
					},
				})
				return block
			} else {
				return gotoNode(te, t, block)
			}
		} else {
			return gotoNode(fe, f, block)
		}
	}

	GenerateFlowSwitch := func(node *base.Node, block []ast.Stmt) []ast.Stmt {
		cond := &ast.BinaryExpr{
			X:  attr(id("frame"), "Flow"),
			Op: token.EQL,
			Y:  constInt(0),
		}
		return emitSwitch(cond, node.GetExit(0), node.GetNext(0), node.GetExit(1), node.GetNext(1), block)
	}

	// Declare the variables.
	// It is easier to do this up front than calculate where they need to be defined.
	for i, info := range f.Registers {
		stmts = append(stmts, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: singleName(RegisterName(DubRegister(i))),
						Type:  goTypeName(info.T),
					},
				},
			},
		})
	}

	// HACK otherwise first label is unused.
	stmts = gotoNode(nil, nodes[0], stmts)

	// Generate Go code from flow blocks
	for i, node := range nodes {
		block := []ast.Stmt{}
		switch data := node.Data.(type) {
		case *DubEntry:
			block = gotoNode(node.GetExit(0), node.GetNext(0), block)
		case *DubExit:
			block = append(block, &ast.ReturnStmt{})
		case *DubSwitch:
			block = emitSwitch(reg(data.Cond), node.GetExit(0), node.GetNext(0), node.GetExit(1), node.GetNext(1), block)
		case DubOp:
			block = GenerateOp(f, data, block)
			block = GenerateFlowSwitch(node, block)
		default:
			panic(data)
		}
		// Label the first statement
		block[0] = &ast.LabeledStmt{Label: id(blockName(i)), Stmt: block[0]}
		// Extend the statement list
		stmts = append(stmts, block...)
	}

	results := []*ast.Field{}
	for i, t := range f.ReturnTypes {
		results = append(results, &ast.Field{
			Names: singleName(returnVarName(i)),
			Type:  goTypeName(t),
		})
	}

	funcDecl := &ast.FuncDecl{
		Name: id(f.Name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: singleName("frame"),
						Type:  ptr(attr(id("runtime"), "State")),
					},
				},
			},
			Results: &ast.FieldList{List: results},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}

	return funcDecl
}

func tagName(s *LLStruct) string {
	return fmt.Sprintf("is%s", s.Name)
}

func addTags(base *LLStruct, parent *LLStruct, decls []ast.Decl) []ast.Decl {
	if parent != nil {
		decls = addTags(base, parent.Implements, decls)
		decls = append(decls, &ast.FuncDecl{
			Name: id(tagName(parent)),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: singleName("node"),
						Type:  goTypeName(base),
					},
				},
			},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{},
			},
			Body: &ast.BlockStmt{},
		})
	}
	return decls
}

func GenerateGoStruct(s *LLStruct, decls []ast.Decl) []ast.Decl {
	var t ast.Expr
	if s.Abstract {
		if len(s.Fields) != 0 {
			panic(s.Name)
		}
		fields := []*ast.Field{
			&ast.Field{
				Names: singleName(tagName(s)),
				Type: &ast.FuncType{
					Params:  &ast.FieldList{},
					Results: &ast.FieldList{},
				},
			},
		}

		t = &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: fields,
			},
		}
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: id(s.Name),
					Type: t,
				},
			},
		})
	} else {
		fields := []*ast.Field{}
		for _, f := range s.Fields {
			fields = append(fields, &ast.Field{
				Names: singleName(f.Name),
				Type:  goTypeName(f.T),
			})
		}
		t = &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		}
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: id(s.Name),
					Type: t,
				},
			},
		})

		decls = addTags(s, s.Implements, decls)
	}
	return decls
}

func GenerateGo(module string, structs []*LLStruct, funcs []*LLFunc) string {
	decls := []ast.Decl{}

	decls = append([]ast.Decl{&ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: strLiteral("evergreen/dub/runtime")},
		},
	}}, decls...)

	for _, f := range structs {
		decls = GenerateGoStruct(f, decls)
	}

	for _, f := range funcs {
		decls = append(decls, GenerateGoFunc(f))
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
