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
	return &ast.ExprStmt{X: &ast.CallExpr{Fun: attr(id("space"), name), Args: args}}
}

func constInt(v int64) ast.Expr {
	// Type converstion hack
	return intLiteral(int(v))
}

func reg(r DubRegister) ast.Expr {
	return intLiteral(int(r))
}

var opToName = map[string]string{
	"+": "Add",
	"<": "Lt",
}

func GenerateGo(r *base.Region) string {
	nodes := base.ReversePostorder(r)
	//info := make([]blockInfo, len(nodes))
	m := map[*base.Node]int{}
	stmts := []ast.Stmt{}

	for i, node := range nodes {
		m[node] = i
	}

	gotoNode := func(next *base.Node) ast.Stmt {
		return &ast.BranchStmt{Tok: token.GOTO, Label: id(blockName(m[next]))}
	}

	for i, node := range nodes {
		block := []ast.Stmt{}
		switch data := node.Data.(type) {
		case *DubEntry:
			block = append(block, gotoNode(node.GetNext(0)))
		case *DubExit:
			block = append(block, &ast.ReturnStmt{})
		case *DubBlock:
			for _, op := range data.Ops {
				switch op := op.(type) {
				case *GetLocalOp:
					block = append(block, emitOp("GetLocal", strLiteral(op.Name), reg(op.Dst)))
				case *SetLocalOp:
					block = append(block, emitOp("GetLocal", reg(op.Src), strLiteral(op.Name)))
				case *BinaryOp:
					name, ok := opToName[op.Op]
					if !ok {
						panic(op.Op)
					}
					block = append(block, emitOp(name, reg(op.Left), reg(op.Right), reg(op.Dst)))
				case *ConstantIntOp:
					block = append(block, emitOp("ConstantInt", constInt(op.Value), reg(op.Dst)))
				default:
					panic(op)
				}
			}
			block = append(block, gotoNode(node.GetNext(0)))
		case *DubSwitch:
			ifs := &ast.IfStmt{
				Cond: emitOp("Bool", reg(data.Cond)).X,
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						gotoNode(node.GetNext(0)),
					},
				},
				Else: gotoNode(node.GetNext(1)),
			}
			block = append(block, ifs)
		default:
			panic(data)
		}
		// Label the first statement
		block[0] = &ast.LabeledStmt{Label: id(blockName(i)), Stmt: block[0]}
		// Extend the statement list
		stmts = append(stmts, block...)
	}

	funcDecl := &ast.FuncDecl{
		Name: id("generated_function"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: singleName("space"),
						Type:  ptr(id("DubSpace")),
					},
					&ast.Field{
						Names: singleName("frame"),
						Type:  ptr(id("DubFrame")),
					},
				},
			},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}

	decls := []ast.Decl{funcDecl}
	f := &ast.File{
		Name:  id("dub"),
		Decls: decls,
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)

	return buf.String()
}
