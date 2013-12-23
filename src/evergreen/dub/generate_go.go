package dub

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

func genFunc(decl *FuncDecl) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: id(decl.Name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: &ast.FieldList{List: []*ast.Field{&ast.Field{Type: ptr(id("bogus_type2"))}}},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.LabeledStmt{Label: id("foo"), Stmt: &ast.AssignStmt{Lhs: []ast.Expr{id("block")}, Tok: token.DEFINE, Rhs: []ast.Expr{intLiteral(0)}}},
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
		},
	}
}

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
	"==": token.EQL,
	"!=": token.NEQ,
	"<":  token.LSS,
	">":  token.GTR,
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

	gotoNode := func(next *base.Node) ast.Stmt {
		return &ast.BranchStmt{Tok: token.GOTO, Label: id(blockName(m[next]))}
	}

	emitSwitch := func(cond ast.Expr, t *base.Node, f *base.Node) ast.Stmt {
		if t != nil {
			if f != nil {
				return &ast.IfStmt{
					Cond: cond,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							gotoNode(t),
						},
					},
					Else: gotoNode(f),
				}
			} else {
				return gotoNode(t)
			}
		} else {
			return gotoNode(f)
		}
	}

	// Declare the variables.
	// It is easier to do this up front than calculate where they need to be defined.
	for i, info := range f.Registers {
		if info.T == nil {
			continue
		}
		typeName := "unknown"
		switch t := info.T.(type) {
		case *BoolType:
			typeName = "bool"
		case *IntType:
			typeName = "int"
		case *RuneType:
			typeName = "rune"
		case *StringType:
			typeName = "string"
		default:
			panic(t)
		}
		stmts = append(stmts, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: singleName(RegisterName(DubRegister(i))),
						Type:  id(typeName),
					},
				},
			},
		})
	}

	// HACK otherwise first label is unused.
	stmts = append(stmts, gotoNode(nodes[0]))

	// Generate Go code from flow blocks
	for i, node := range nodes {
		block := []ast.Stmt{}
		switch data := node.Data.(type) {
		case *DubEntry:
			block = append(block, gotoNode(node.GetNext(0)))
		case *DubExit:
			block = append(block, &ast.ReturnStmt{})
		case *DubBlock:
			// Generate statements.
			for _, op := range data.Ops {
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
				case *Read:
					block = append(block, opAssign(
						emitOp("Read").X,
						op.Dst,
					))
				case *Fail:
					block = append(block, emitOp("Fail"))
				case *Checkpoint:
					block = append(block, opAssign(
						emitOp("Checkpoint").X,
						op.Dst,
					))
				case *Recover:
					block = append(block, emitOp("Recover", reg(op.Src)))
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
			}
			cond := &ast.BinaryExpr{
				X:  attr(id("frame"), "Flow"),
				Op: token.EQL,
				Y:  constInt(0),
			}

			next := emitSwitch(cond, node.GetNext(0), node.GetNext(1))
			block = append(block, next)
		case *DubSwitch:
			next := emitSwitch(reg(data.Cond), node.GetNext(0), node.GetNext(1))
			block = append(block, next)
		default:
			panic(data)
		}
		// Label the first statement
		block[0] = &ast.LabeledStmt{Label: id(blockName(i)), Stmt: block[0]}
		// Extend the statement list
		stmts = append(stmts, block...)
	}

	funcDecl := &ast.FuncDecl{
		Name: id(f.Name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: singleName("frame"),
						Type:  ptr(attr(id("dub"), "DubState")),
					},
				},
			},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}

	return funcDecl
}

func GenerateGo(module string, funcs []*LLFunc) string {
	decls := []ast.Decl{}

	decls = append([]ast.Decl{&ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: 1,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: strLiteral("evergreen/dub")},
		},
	}}, decls...)

	for _, f := range funcs {
		decls = append(decls, GenerateGoFunc(f))
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
