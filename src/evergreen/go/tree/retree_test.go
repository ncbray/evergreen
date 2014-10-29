package tree

import (
	"evergreen/assert"
	"evergreen/base"
	"testing"
)

func binaryOpExample(swap bool) (*FuncDecl, []Stmt) {
	decl := &FuncDecl{
		Name:            "foo",
		Type:            &FuncType{},
		LocalInfo_Scope: &LocalInfo_Scope{},
	}

	intType := &ExternalType{Name: "int"}
	a := decl.CreateLocalInfo("a", &TypeRef{Impl: intType})
	b := decl.CreateLocalInfo("b", &TypeRef{Impl: intType})
	ret := decl.CreateLocalInfo("ret0", &TypeRef{Impl: intType})

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
	decl.Body = block
	return decl, block
}

func functionExample() *FuncDecl {
	decl, _ := binaryOpExample(false)
	return decl
}

func TestBinaryExprDefUse(t *testing.T) {
	decl, block := binaryOpExample(false)
	du := makeApproxDefUse(decl)
	defUseBlock(block, du)

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

func TestBinaryExprRetree(t *testing.T) {
	decl, block := binaryOpExample(false)

	du := makeApproxDefUse(decl)

	defUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	generateBlock(gen, block, w)
	checkCode(b.String(), "ret0 := 2 + 3\n", t)
}

func TestBinaryExprRetreeSwap(t *testing.T) {
	decl, block := binaryOpExample(true)

	du := makeApproxDefUse(decl)

	defUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
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
