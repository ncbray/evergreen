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

func opMultiAssign(expr ast.Expr, dsts []DubRegister) ast.Stmt {
	if len(dsts) != 0 {
		lhs := make([]ast.Expr, len(dsts))
		for i, dst := range dsts {
			if dst != NoRegister {
				lhs[i] = reg(dst)
			} else {
				lhs[i] = id("_")
			}
		}
		return &ast.AssignStmt{
			Lhs: lhs,
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
		args := []ast.Expr{
			id("frame"),
		}
		for _, arg := range op.Args {
			args = append(args, reg(arg))
		}
		block = append(block, opMultiAssign(
			&ast.CallExpr{
				Fun:  id(op.Name),
				Args: args,
			},
			op.Dsts,
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

	case *TransferOp:
		lhs := make([]ast.Expr, len(op.Dsts))
		for i, dst := range op.Dsts {
			lhs[i] = reg(dst)
		}
		rhs := make([]ast.Expr, len(op.Srcs))
		for i, src := range op.Srcs {
			rhs[i] = reg(src)
		}
		block = append(block, &ast.AssignStmt{
			Lhs: lhs,
			Tok: token.ASSIGN,
			Rhs: rhs,
		})
	default:
		panic(op)
	}
	return block
}

func generateNode(info *regionInfo, node base.NodeID, block []ast.Stmt) []ast.Stmt {
	g := info.decl.CFG
	op := info.decl.Ops[node]
	switch data := op.(type) {
	case *EntryOp:
		block = gotoNode(info, g.GetExit(node, 0), block)
	case *FlowExitOp:
		block = append(block, &ast.ReturnStmt{})
	case *ExitOp:
	case *SwitchOp:
		block = emitSwitch(info, reg(data.Cond), g.GetExit(node, 0), g.GetExit(node, 1), block)
	case DubOp:
		block = GenerateOp(info.decl, data, block)
		block = generateFlowSwitch(info, node, block)
	default:
		panic(data)
	}
	return block
}

type regionInfo struct {
	decl   *LLFunc
	labels map[base.NodeID]int
}

func gotoNode(info *regionInfo, n base.NodeID, block []ast.Stmt) []ast.Stmt {
	label, ok := info.labels[n]
	if ok {
		return append(block, &ast.BranchStmt{Tok: token.GOTO, Label: id(blockName(label))})
	} else {
		return generateNode(info, n, block)
	}
}

func emitSwitch(info *regionInfo, cond ast.Expr, t base.NodeID, f base.NodeID, block []ast.Stmt) []ast.Stmt {
	if t != base.NoNode {
		if f != base.NoNode {
			block = append(block, &ast.IfStmt{
				Cond: cond,
				Body: &ast.BlockStmt{
					List: gotoNode(info, t, nil),
				},
				Else: &ast.BlockStmt{
					List: gotoNode(info, f, nil),
				},
			})
			return block
		} else {
			return gotoNode(info, t, block)
		}
	} else {
		return gotoNode(info, f, block)
	}
}

func generateFlowSwitch(info *regionInfo, node base.NodeID, block []ast.Stmt) []ast.Stmt {
	g := info.decl.CFG
	numExits := g.NumExits(node)

	if numExits == 2 {
		cond := &ast.BinaryExpr{
			X:  attr(id("frame"), "Flow"),
			Op: token.EQL,
			Y:  constInt(0),
		}
		t := g.GetExit(node, 0)
		f := g.GetExit(node, 1)
		return emitSwitch(info, cond, t, f, block)
	} else if numExits == 1 {
		return gotoNode(info, g.GetExit(node, 0), block)
	} else {
		panic(info.decl.Ops[node])
	}
}

func IsParam(f *LLFunc, r DubRegister) bool {
	for _, p := range f.Params {
		if p == r {
			return true
		}
	}
	return false
}

func GenerateGoFunc(f *LLFunc) ast.Decl {
	g := f.CFG
	order, _ := base.ReversePostorder(g)

	heads := []base.NodeID{}
	labels := map[base.NodeID]int{}
	uid := 0

	nit := base.OrderedIterator(order)
	for nit.Next() {
		n := nit.Value()
		if (n == g.Entry() || g.NumEntries(n) >= 2) && n != g.Exit() {
			heads = append(heads, n)
			labels[n] = uid
			uid = uid + 1
		}
	}
	info := &regionInfo{decl: f, labels: labels}

	stmts := []ast.Stmt{}

	// Declare the variables.
	// It is easier to do this up front than calculate where they need to be defined.
	for i, info := range f.Registers {
		r := DubRegister(i)
		if IsParam(f, r) {
			continue
		}
		stmts = append(stmts, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: singleName(RegisterName(r)),
						Type:  goTypeName(info.T),
					},
				},
			},
		})
	}

	// Generate Go code from flow blocks
	for _, node := range heads {
		block := []ast.Stmt{}
		block = generateNode(info, node, block)
		// Label the first statement
		label, _ := info.labels[node]
		// HACK assume label 0 is always the entry node.
		if label != 0 {
			block[0] = &ast.LabeledStmt{Label: id(blockName(label)), Stmt: block[0]}
		}
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

	fields := []*ast.Field{
		&ast.Field{
			Names: singleName("frame"),
			Type:  ptr(attr(id("runtime"), "State")),
		},
	}

	for _, p := range f.Params {
		fields = append(fields, &ast.Field{
			Names: singleName(RegisterName(p)),
			Type:  goTypeName(f.Registers[p].T),
		})
	}

	funcDecl := &ast.FuncDecl{
		Name: id(f.Name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: fields,
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

	imports := []ast.Spec{}
	if len(funcs) > 0 {
		imports = append(imports, &ast.ImportSpec{Path: strLiteral("evergreen/dub/runtime")})
	}

	if len(imports) > 0 {
		decls = append([]ast.Decl{&ast.GenDecl{
			Tok:    token.IMPORT,
			Lparen: 1,
			Specs:  imports,
		}}, decls...)
	}

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
