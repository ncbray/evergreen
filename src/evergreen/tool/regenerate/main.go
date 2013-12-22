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

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
}

func (node *BinaryOp) isASTExpr() {
}

type GetName struct {
	Name string
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

type FuncDecl struct {
	Name  string
	Block []ASTExpr
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

func parseExpr(s *scanner.Scanner) (result ASTExpr) {
	tok := s.Scan()
	text := s.TokenText()

	switch tok {
	case scanner.Ident:
		switch text {
		case "plus":
			block := parseCodeBlock(s)
			result = &Repeat{Block: block, Min: 1}
		case "if":
			expr := parseExpr(s)
			block := parseCodeBlock(s)
			result = &If{Expr: expr, Block: block}
		case "read":
			result = &Read{}
		case "fail":
			result = &Fail{}
		default:
			result = &GetName{Name: text}
		}
	case scanner.Char:
		v, _ := strconv.Unquote(s.TokenText())
		result = &RuneLiteral{Value: []rune(v)[0]}
	case '=':
		expr := parseExpr(s)
		dst := getName(s)
		result = &SetName{Expr: expr, Name: dst}
	case '<', '>':
		l := parseExpr(s)
		r := parseExpr(s)
		result = &BinaryOp{Left: l, Op: text, Right: r}
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

type DubBuilder struct {
	nextRegister dub.DubRegister
}

func (builder *DubBuilder) CreateRegister() dub.DubRegister {
	temp := builder.nextRegister
	builder.nextRegister += 1
	return temp
}

func lowerExpr(expr ASTExpr, r *base.Region, builder *DubBuilder) dub.DubRegister {
	switch expr := expr.(type) {
	case *If:
		// TODO Min
		//l := dub.CreateRegion()

		cond := lowerExpr(expr.Expr, r, builder)
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
		checkpoint := builder.CreateRegister()
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Checkpoint{Dst: checkpoint},
			})

			r.Connect(0, body)
			body.SetExit(0, r.GetExit(0))
		}

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Normal flow iterates
		block.GetExit(0).TransferEntries(block.Head())
		// Stop iterating on failure
		block.GetExit(1).TransferEntries(block.GetExit(0))

		r.Splice(0, block)

		// Recover
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			r.Connect(0, body)
			body.SetExit(0, r.GetExit(0))
		}

		return dub.NoRegister

	case *GetName:
		dst := builder.CreateRegister()
		body := dub.CreateBlock([]dub.DubOp{
			&dub.GetLocalOp{Name: expr.Name, Dst: dst},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dst

	case *SetName:
		src := lowerExpr(expr.Expr, r, builder)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.SetLocalOp{Src: src, Name: expr.Name},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dub.NoRegister

	case *RuneLiteral:
		dst := builder.CreateRegister()
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(0, body)
		body.SetExit(0, r.GetExit(0))
		return dst

	case *Read:
		dst := builder.CreateRegister()
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
		left := lowerExpr(expr.Left, r, builder)
		right := lowerExpr(expr.Right, r, builder)
		dst := builder.CreateRegister()
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
	default:
		panic(expr)
	}

}

func lowerBlock(block []ASTExpr, builder *DubBuilder) *base.Region {
	r := dub.CreateRegion()
	for _, expr := range block {
		lowerExpr(expr, r, builder)
	}
	return r
}

func lowerAST(decl *FuncDecl) *base.Region {
	builder := &DubBuilder{}
	return lowerBlock(decl.Block, builder)
}

func main() {
	decls := parseDASM("dasm/math.dasm")
	for _, decl := range decls {
		r := lowerAST(decl)
		dot := base.RegionToDot(r)
		outfile := filepath.Join("output", fmt.Sprintf("%s.svg", decl.Name))
		io.WriteDot(dot, outfile)

		registers := []dub.RegisterInfo{}
		fmt.Println(dub.GenerateGo(r, registers))
	}

	l := dub.CreateRegion()
	cond := dub.CreateBlock([]dub.DubOp{
		&dub.BinaryOp{
			Left:  0,
			Op:    "<",
			Right: 1,
			Dst:   2,
		},
	})
	decide := dub.CreateSwitch(2)
	body := dub.CreateBlock([]dub.DubOp{
		&dub.ConstantIntOp{Value: 1, Dst: 3},
		&dub.BinaryOp{
			Left:  0,
			Op:    "+",
			Right: 3,
			Dst:   0,
		},
	})

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.GetExit(0))

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

}
