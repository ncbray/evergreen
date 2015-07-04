package tree

import (
	"evergreen/assert"
	"evergreen/go/core"
	"evergreen/text"
	"testing"
)

func binaryOpExample(swap bool) *FuncDecl {
	decl := &FuncDecl{
		Name:            "foo",
		Type:            &FuncTypeRef{},
		LocalInfo_Scope: &LocalInfo_Scope{},
		Block:           &Block{},
	}

	intType := &core.ExternalType{Name: "int"}
	a := decl.CreateLocalInfo("a", &NameRef{T: intType})
	b := decl.CreateLocalInfo("b", &NameRef{T: intType})
	ret := decl.CreateLocalInfo("ret0", &NameRef{T: intType})

	first := a
	second := b
	if swap {
		first, second = second, first
	}

	block := []Stmt{
		&Assign{
			Sources: []Expr{
				&IntLiteral{
					Value: 2,
				},
			},
			Op: ":=",
			Targets: []Target{
				&SetLocal{Info: first},
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
				&SetLocal{Info: second},
			},
		},
		&Assign{
			Sources: []Expr{
				&BinaryExpr{
					Left:  &GetLocal{Info: a},
					Op:    "+",
					Right: &GetLocal{Info: b},
				},
			},
			Op: ":=",
			Targets: []Target{
				&SetLocal{Info: ret},
			},
		},
	}
	decl.Block.Body = block
	return decl
}

func functionExample() *FuncDecl {
	return binaryOpExample(false)
}

func TestBinaryExprDefUse(t *testing.T) {
	decl := binaryOpExample(false)
	du := makeApproxDefUse(decl)
	defUseBlock(decl.Block, du)

	info := du.GetLocalInfo(0)
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 1)

	info = du.GetLocalInfo(1)
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 1)

	info = du.GetLocalInfo(2)
	assert.IntEquals(t, info.Defs, 1)
	assert.IntEquals(t, info.Uses, 0)
}

func TestBinaryExprConsolidate(t *testing.T) {
	decl := binaryOpExample(false)

	du := makeApproxDefUse(decl)

	defUseBlock(decl.Block, du)
	decl.Block = consolidateBlock(decl.Block, du)

	b, w := text.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	generateBlock(gen, decl.Block, w)
	checkCode(b.String(), "ret0 := 2 + 3\n", t)
}

func TestBinaryExprConsolidateSwap(t *testing.T) {
	decl := binaryOpExample(true)

	du := makeApproxDefUse(decl)

	defUseBlock(decl.Block, du)
	decl.Block = consolidateBlock(decl.Block, du)

	b, w := text.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	generateBlock(gen, decl.Block, w)
	checkCode(b.String(), "b := 2\nret0 := 3 + b\n", t)
}

func TestFuncConsolidate(t *testing.T) {
	decl := functionExample()
	consolidateDecl(decl)

	b, w := text.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	GenerateFunc(gen, decl, w)
	checkCode(b.String(), "func foo() {\n\tret0 := 2 + 3\n}\n", t)
}
