package tree

import (
	"evergreen/assert"
	"evergreen/base"
	"testing"
)

func binaryOpExample(swap bool) []Stmt {
	first := "a"
	second := "b"
	if swap {
		first, second = second, first
	}
	return []Stmt{
		&Assign{
			Sources: []Expr{
				&IntLiteral{
					Value: 2,
				},
			},
			Op: ":=",
			Targets: []Target{
				&SetName{Text: first},
			},
		},
		&Assign{
			Sources: []Expr{
				&IntLiteral{
					Value: 3,
				},
			},
			Op: ":=",
			Targets: []Target{
				&SetName{Text: second},
			},
		},
		&Assign{
			Sources: []Expr{
				&BinaryExpr{
					Left:  &GetName{Text: "a"},
					Op:    "+",
					Right: &GetName{Text: "b"},
				},
			},
			Op: ":=",
			Targets: []Target{
				&SetName{Text: "ret0"},
			},
		},
	}
}

func functionExample() *FuncDecl {
	return &FuncDecl{
		Name: "foo",
		Type: &FuncType{},
		Body: binaryOpExample(false),
	}
}

func TestBinaryExprDefUse(t *testing.T) {
	block := binaryOpExample(false)
	du := makeApproxDefUse(&FuncDecl{})
	approxDefUseBlock(block, du)

	info := du.GetInfo("a")
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 1)

	info = du.GetInfo("b")
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 1)

	info = du.GetInfo("ret0")
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 0)
}

func TestBinaryExprRetree(t *testing.T) {
	block := binaryOpExample(false)

	du := makeApproxDefUse(&FuncDecl{})

	approxDefUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	gen := &textGenerator{}
	generateBlock(gen, block, w)
	checkCode(b.String(), "ret0 := 2 + 3\n", t)
}

func TestBinaryExprRetreeSwap(t *testing.T) {
	block := binaryOpExample(true)

	du := makeApproxDefUse(&FuncDecl{})

	approxDefUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	gen := &textGenerator{}
	generateBlock(gen, block, w)
	checkCode(b.String(), "b := 2\nret0 := 3 + b\n", t)
}

func TestFuncRetree(t *testing.T) {
	decl := functionExample()
	retreeDecl(decl)

	b, w := base.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	GenerateFunc(gen, decl, w)
	checkCode(b.String(), "func foo() {\n\tret0 := 2 + 3\n}\n", t)
}
