package main

import (
	"bytes"
	"evergreen/base"
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"text/scanner"
)

type ASTExpr interface {
	isASTExpr()
}

type If struct {
	Expr  ASTExpr
	Block []ASTExpr
}

func (node *If) isASTExpr() {
}

type Repeat struct {
	Block []ASTExpr
	Min   int
}

func (node *Repeat) isASTExpr() {
}

type Slice struct {
	Block []ASTExpr
}

func (node *Slice) isASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     dub.DubType
}

func (node *BinaryOp) isASTExpr() {
}

type Call struct {
	Name string
}

func (node *Call) isASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) isASTExpr() {
}

type RuneLiteral struct {
	Value rune
}

func (node *RuneLiteral) isASTExpr() {
}

type SetName struct {
	Expr ASTExpr
	Name string
	Info int
}

func (node *SetName) isASTExpr() {
}

type Read struct {
}

func (node *Read) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
}

type LocalInfo struct {
	Name string
	T    dub.DubType
}

type FuncDecl struct {
	Name   string
	Block  []ASTExpr
	Locals []*LocalInfo
}

func getName(s *scanner.Scanner) string {
	tok := s.Scan()
	if tok != scanner.Ident {
		panic(tok)
	}
	return s.TokenText()
}

func getPunc(s *scanner.Scanner, punc rune) {
	tok := s.Scan()
	if tok != punc {
		panic(tok)
	}
}

var nameToOp = map[string]string{
	"eq": "==",
	"ne": "!=",
	"gt": ">",
	"lt": "<",
}

func parseExpr(s *scanner.Scanner) (result ASTExpr) {
	tok := s.Scan()
	text := s.TokenText()

	switch tok {
	case scanner.Ident:
		switch text {
		case "star":
			block := parseCodeBlock(s)
			result = &Repeat{Block: block, Min: 0}
		case "plus":
			block := parseCodeBlock(s)
			result = &Repeat{Block: block, Min: 1}
		case "slice":
			block := parseCodeBlock(s)
			result = &Slice{Block: block}
		case "if":
			expr := parseExpr(s)
			block := parseCodeBlock(s)
			result = &If{Expr: expr, Block: block}
		case "define":
			expr := parseExpr(s)
			dst := getName(s)
			result = &SetName{Expr: expr, Name: dst}
		case "read":
			result = &Read{}
		case "fail":
			result = &Fail{}
		case "eq", "ne", "gt", "lt":
			l := parseExpr(s)
			r := parseExpr(s)
			op := nameToOp[text]
			result = &BinaryOp{Left: l, Op: op, Right: r}
		case "call":
			name := getName(s)
			result = &Call{Name: name}
		default:
			result = &GetName{Name: text}
		}
	case scanner.Char:
		v, _ := strconv.Unquote(s.TokenText())
		result = &RuneLiteral{Value: []rune(v)[0]}
	default:
		panic(tok)
	}
	return
}

func parseCodeBlock(s *scanner.Scanner) (result []ASTExpr) {
	getPunc(s, '{')
	result = []ASTExpr{}
	defer func() {
		if r := recover(); r != nil {
			if s.TokenText() == "}" {
				// End of a block
			} else {
				panic(r)
			}
		}
	}()
	for {
		result = append(result, parseExpr(s))
		if s.TokenText() != "}" {
			tok := s.Scan()
			switch tok {
			case ';':
			case '}':
				return
			default:
				panic(tok)
			}
		}
	}
}

func parseFunction(s *scanner.Scanner) *FuncDecl {
	name := getName(s)
	block := parseCodeBlock(s)
	return &FuncDecl{Name: name, Block: block}
}

func parseFile(s *scanner.Scanner) (decls []*FuncDecl) {
	decls = []*FuncDecl{}

	tok := s.Scan()
	for {
		switch tok {
		case scanner.Ident:
			text := s.TokenText()
			switch text {
			case "func":
				decls = append(decls, parseFunction(s))
			default:
				panic(tok)
			}
		case scanner.EOF:
			return
		default:
			panic(tok)
		}
		// do something with tok
		tok = s.Scan()
	}
}

func parseDASM(filename string) []*FuncDecl {
	data, _ := ioutil.ReadFile(filename)
	s := &scanner.Scanner{}
	s.Init(bytes.NewReader(data))

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unexpected %s @ %s\n", s.TokenText(), s.Pos())
		}
	}()
	return parseFile(s)
}

type semanticScope struct {
	parent *semanticScope
	locals map[string]int
}

func (scope *semanticScope) localInfo(name string) (int, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return -1, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{parent: scope, locals: map[string]int{}}
}

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope) dub.DubType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope)
		return &dub.VoidType{}
	case *If:
		semanticExprPass(decl, expr.Expr, scope)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope))
		return &dub.VoidType{}
	case *BinaryOp:
		semanticExprPass(decl, expr.Left, scope)
		semanticExprPass(decl, expr.Right, scope)
		// HACK assume compare
		t := &dub.BoolType{}
		expr.T = t
		return t
	case *GetName:
		info, found := scope.localInfo(expr.Name)
		if !found {
			panic(expr.Name)
		}
		expr.Info = info
		return decl.Locals[info].T
	case *SetName:
		t := semanticExprPass(decl, expr.Expr, scope)
		info := len(decl.Locals)
		decl.Locals = append(decl.Locals, &LocalInfo{Name: expr.Name, T: t})
		scope.locals[expr.Name] = info
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope)
		return &dub.StringType{}
	case *Read:
		return &dub.RuneType{}
	case *RuneLiteral:
		return &dub.RuneType{}
	case *Fail:
		return &dub.VoidType{}
	default:
		panic(expr)
	}
}

func semanticBlockPass(decl *FuncDecl, block []ASTExpr, scope *semanticScope) {
	for _, expr := range block {
		semanticExprPass(decl, expr, scope)
	}
}

func semanticFuncPass(decl *FuncDecl) {
	semanticBlockPass(decl, decl.Block, childScope(nil))
}

func semanticPass(decls []*FuncDecl) {
	for _, decl := range decls {
		semanticFuncPass(decl)
	}
}

type DubBuilder struct {
	decl      *FuncDecl
	registers []dub.RegisterInfo
	localMap  []dub.DubRegister
}

func (builder *DubBuilder) CreateRegister(t dub.DubType) dub.DubRegister {
	builder.registers = append(builder.registers, dub.RegisterInfo{T: t})
	return dub.DubRegister(len(builder.registers) - 1)
}

func lowerExpr(expr ASTExpr, r *base.Region, builder *DubBuilder, used bool) dub.DubRegister {
	switch expr := expr.(type) {
	case *If:
		// TODO Min
		//l := dub.CreateRegion()

		cond := lowerExpr(expr.Expr, r, builder, true)
		block := lowerBlock(expr.Block, builder)

		// TODO conditional
		decide := dub.CreateSwitch(cond)

		r.Connect(0, decide)
		decide.SetExit(0, r.GetExit(0))
		r.Splice(0, block)
		decide.SetExit(1, r.GetExit(0))

		return dub.NoRegister

	case *Repeat:

		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			block := lowerBlock(expr.Block, builder)
			r.Splice(0, block)
		}

		// Checkpoint
		checkpoint := builder.CreateRegister(&dub.IntType{})
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})

		r.Connect(0, head)
		head.SetExit(0, r.GetExit(0))

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Normal flow iterates
		// NOTE actually connects nodes in two different regions.  Kinda hackish.
		block.GetExit(0).TransferEntries(head)
		// Stop iterating on failure
		block.GetExit(1).TransferEntries(block.GetExit(0))

		// Recover
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			block.Connect(0, body)
			body.SetExit(0, block.GetExit(0))
		}

		r.Splice(0, block)

		return dub.NoRegister

	case *GetName:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateRegister(builder.decl.Locals[expr.Info].T)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CopyOp{Src: builder.localMap[expr.Info], Dst: dst},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dst

	case *SetName:
		src := lowerExpr(expr.Expr, r, builder, true)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CopyOp{Src: src, Dst: builder.localMap[expr.Info]},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dub.NoRegister

	case *RuneLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateRegister(&dub.RuneType{})
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dst

	case *Read:
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(&dub.RuneType{})
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.Read{Dst: dst},
		})
		r.Connect(0, body)
		r.AttachDefaultExits(body)
		return dst

	case *Fail:
		body := dub.CreateBlock([]dub.DubOp{
			&dub.Fail{},
		})
		r.Connect(0, body)
		body.SetExit(1, r.GetExit(1))
		return dub.NoRegister

	case *BinaryOp:
		left := lowerExpr(expr.Left, r, builder, true)
		right := lowerExpr(expr.Right, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.BinaryOp{
				Left:  left,
				Op:    expr.Op,
				Right: right,
				Dst:   dst,
			},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dst
	case *Slice:
		start := builder.CreateRegister(&dub.IntType{})
		// HACK assume checkpoint is just the index
		{
			head := dub.CreateBlock([]dub.DubOp{
				&dub.Checkpoint{Dst: start},
			})
			r.Connect(0, head)
			head.SetExit(0, r.GetExit(0))
		}
		block := lowerBlock(expr.Block, builder)
		r.Splice(0, block)

		// Create a slice
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(&dub.StringType{})
		}
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Slice{Src: start, Dst: dst},
			})

			r.Connect(0, body)
			body.SetExit(0, r.GetExit(0))
		}
		return dst
	default:
		panic(expr)
	}

}

func lowerBlock(block []ASTExpr, builder *DubBuilder) *base.Region {
	r := dub.CreateRegion()
	for _, expr := range block {
		lowerExpr(expr, r, builder, false)
	}
	return r
}

func lowerAST(decl *FuncDecl) *dub.LLFunc {
	f := &dub.LLFunc{Name: decl.Name}
	builder := &DubBuilder{decl: decl}
	// Allocate register for locals
	builder.localMap = make([]dub.DubRegister, len(decl.Locals))
	for i, info := range decl.Locals {
		builder.localMap[i] = builder.CreateRegister(info.T)
	}
	f.Region = lowerBlock(decl.Block, builder)
	f.Registers = builder.registers
	return f
}

func main() {
	decls := parseDASM("dasm/math.dasm")
	semanticPass(decls)

	funcs := []*dub.LLFunc{}
	for _, decl := range decls {
		f := lowerAST(decl)
		funcs = append(funcs, f)

		// Dump flowgraph
		dot := base.RegionToDot(f.Region)
		outfile := filepath.Join("output", fmt.Sprintf("%s.svg", f.Name))
		io.WriteDot(dot, outfile)
	}

	code := dub.GenerateGo("math", funcs)
	fmt.Println(code)
	io.WriteFile("src/generated/math/parser.go", []byte(code))

	/*
		i := "integer"
		b := "boolean"

		registers := []dub.RegisterInfo{
			dub.RegisterInfo{T: i},
			dub.RegisterInfo{T: i},
			dub.RegisterInfo{T: b},
			dub.RegisterInfo{T: i},
		}

		dot := base.RegionToDot(l)
		outfile := filepath.Join("output", "test.svg")

		result := make(chan error, 2)
		go func() {
			err := io.WriteDot(dot, outfile)
			result <- err
		}()

		fmt.Println(dub.GenerateGo(l, registers))

		err := <-result
		if err != nil {
			fmt.Println(err)
		}
	*/
}
