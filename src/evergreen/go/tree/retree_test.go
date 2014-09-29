package tree

import (
	"evergreen/base"
	"testing"
)

func checkInt(actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

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
			Targets: []Expr{
				&NameRef{Text: first},
			},
		},
		&Assign{
			Sources: []Expr{
				&IntLiteral{
					Value: 3,
				},
			},
			Op: ":=",
			Targets: []Expr{
				&NameRef{Text: second},
			},
		},
		&Assign{
			Sources: []Expr{
				&BinaryExpr{
					Left:  &NameRef{Text: "a"},
					Op:    "+",
					Right: &NameRef{Text: "b"},
				},
			},
			Op: ":=",
			Targets: []Expr{
				&NameRef{Text: "ret0"},
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
	du := makeApproxDefUse()
	approxDefUseBlock(block, du)

	info := du.GetInfo("a")
	checkInt(info.Defs, 1, t)
	checkInt(info.Uses, 1, t)

	info = du.GetInfo("b")
	checkInt(info.Defs, 1, t)
	checkInt(info.Uses, 1, t)

	info = du.GetInfo("ret0")
	checkInt(info.Defs, 1, t)
	checkInt(info.Uses, 0, t)
}

func TestBinaryExprRetree(t *testing.T) {
	block := binaryOpExample(false)

	du := makeApproxDefUse()

	approxDefUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	generateBlock(block, w)
	checkCode(b.String(), "ret0 := 2 + 3\n", t)
}

func TestBinaryExprRetreeSwap(t *testing.T) {
	block := binaryOpExample(true)

	du := makeApproxDefUse()

	approxDefUseBlock(block, du)
	block = retreeBlock(block, du)

	b, w := base.BufferedCodeWriter()
	generateBlock(block, w)
	checkCode(b.String(), "b := 2\nret0 := 3 + b\n", t)
}

func TestFuncRetree(t *testing.T) {
	decl := functionExample()
	retreeDecl(decl)

	b, w := base.BufferedCodeWriter()
	GenerateFunc(decl, w)
	checkCode(b.String(), "func foo() {\n\tret0 := 2 + 3\n}\n", t)
}
