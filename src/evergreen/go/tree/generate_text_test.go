package tree

import (
	"evergreen/assert"
	"evergreen/text"
	"go/format"
	"testing"
)

func checkFormat(code string, t *testing.T) {
	out, err := format.Source([]byte(code))
	if err != nil {
		t.Fatal(err)
	}
	outs := string(out)
	assert.StringEquals(t, code, outs)
}

func checkCode(actual string, expected string, t *testing.T) {
	assert.StringEquals(t, actual, expected)
	checkFormat(actual, t)
}

func genExpr(expr Expr) string {
	gen := &textGenerator{}
	return GenerateExpr(gen, expr)
}

func TestIntLiteral(t *testing.T) {
	expr := &IntLiteral{Value: 17}
	result := genExpr(expr)
	checkCode(result, "17", t)
}

func TestBoolLiteral(t *testing.T) {
	expr := &BoolLiteral{Value: true}
	result := genExpr(expr)
	checkCode(result, "true", t)
}

func TestStringLiteral(t *testing.T) {
	expr := &StringLiteral{Value: "foo\nbar"}
	result := genExpr(expr)
	checkCode(result, "\"foo\\nbar\"", t)
}

func TestRuneLiteral(t *testing.T) {
	expr := &RuneLiteral{Value: 'x'}
	result := genExpr(expr)
	checkCode(result, "'x'", t)
}

func TestNilLiteral(t *testing.T) {
	expr := &NilLiteral{}
	result := genExpr(expr)
	checkCode(result, "nil", t)
}

func TestSimpleUnaryExpr(t *testing.T) {
	expr := &UnaryExpr{
		Op:   "-",
		Expr: &IntLiteral{Value: 17},
	}
	result := genExpr(expr)
	checkCode(result, "-17", t)
}

func TestSimpleBinaryExpr(t *testing.T) {
	expr := &BinaryExpr{
		Left:  &IntLiteral{Value: 12},
		Op:    "+",
		Right: &IntLiteral{Value: 34},
	}
	result := genExpr(expr)
	checkCode(result, "12 + 34", t)
}

func TestSimpleundExpr(t *testing.T) {
	expr := &BinaryExpr{
		Left: &BinaryExpr{
			Left:  &IntLiteral{Value: 12},
			Op:    "+",
			Right: &IntLiteral{Value: 34},
		},
		Op:    "*",
		Right: &IntLiteral{Value: 56},
	}
	result := genExpr(expr)
	checkCode(result, "(12 + 34) * 56", t)
}

func TestUseAssosiativity(t *testing.T) {
	expr := &BinaryExpr{
		Left: &BinaryExpr{
			Left:  &IntLiteral{Value: 12},
			Op:    "*",
			Right: &IntLiteral{Value: 34},
		},
		Op:    "*",
		Right: &IntLiteral{Value: 56},
	}
	result := genExpr(expr)
	checkCode(result, "12 * 34 * 56", t)
}

func TestPreserveAssosiativity(t *testing.T) {
	expr := &BinaryExpr{
		Left: &IntLiteral{Value: 12},
		Op:   "*",
		Right: &BinaryExpr{
			Left:  &IntLiteral{Value: 34},
			Op:    "*",
			Right: &IntLiteral{Value: 56},
		},
	}
	result := genExpr(expr)
	checkCode(result, "12 * (34 * 56)", t)
}

func TestMethodCall(t *testing.T) {
	expr := &Call{
		Expr: &Selector{
			Expr: &GetName{Text: "foo"},
			Text: "bar",
		},
		Args: []Expr{
			&IntLiteral{Value: 12},
			&IntLiteral{Value: 34},
		},
	}
	result := genExpr(expr)
	checkCode(result, "foo.bar(12, 34)", t)
}

func TestIndex(t *testing.T) {
	expr := &Index{
		Expr:  &GetName{Text: "foo"},
		Index: &GetName{Text: "bar"},
	}
	result := genExpr(expr)
	checkCode(result, "foo[bar]", t)
}

func TestTypeCoerce(t *testing.T) {
	expr := &TypeCoerce{
		Type: &SliceRef{Element: &NameRef{Name: "rune"}},
		Expr: &GetName{Text: "s"},
	}
	result := genExpr(expr)
	checkCode(result, "[]rune(s)", t)
}

func TestStructLiteral(t *testing.T) {
	expr := &StructLiteral{
		Type: &NameRef{Name: "Foo"},
		Args: []*KeywordExpr{
			&KeywordExpr{
				Name: "Bar",
				Expr: &IntLiteral{Value: 12},
			},
			&KeywordExpr{
				Name: "Baz",
				Expr: &IntLiteral{Value: 34},
			},
		},
	}
	result := genExpr(expr)
	checkCode(result, "Foo{Bar: 12, Baz: 34}", t)
}

func TestListLiteral(t *testing.T) {
	expr := &ListLiteral{
		Type: &SliceRef{Element: &NameRef{Name: "int"}},
		Args: []Expr{
			&IntLiteral{Value: 12},
			&IntLiteral{Value: 34},
		},
	}
	result := genExpr(expr)
	checkCode(result, "[]int{12, 34}", t)
}

func TestFuncDecl(t *testing.T) {
	decl := &FuncDecl{
		Name: "foo",
		Recv: &Param{Name: "o", Type: &PointerRef{Element: &NameRef{Name: "Obj"}}},
		Type: &FuncTypeRef{
			Params: []*Param{
				&Param{Name: "cond", Type: &NameRef{Name: "bool"}},
				&Param{Name: "names", Type: &SliceRef{Element: &NameRef{Name: "string"}}},
			},
			Results: []*Param{
				&Param{Name: "biz", Type: &NameRef{Name: "int"}},
				&Param{Name: "baz", Type: &PointerRef{Element: &NameRef{Name: "int"}}},
			},
		},
		Block: &Block{Body: []Stmt{
			&If{
				Cond: &GetName{Text: "cond"},
				T: &Block{
					Body: []Stmt{
						&StringLiteral{Value: "hello"},
					},
				},
				F: &Block{
					Body: []Stmt{
						&If{
							Cond: &UnaryExpr{Op: "!", Expr: &GetName{Text: "cond"}},
							T: &Block{
								Body: []Stmt{
									&StringLiteral{Value: "goodbye"},
								},
							},
							F: &Block{
								Body: []Stmt{
									&StringLiteral{Value: "impossible"},
								},
							},
						},
					},
				},
			},
			&Assign{
				Sources: []Expr{
					&Call{
						Expr: &GetName{Text: "bar"},
						Args: []Expr{&GetName{Text: "names"}},
					},
					&IntLiteral{Value: 7},
				},
				Op: ":=",
				Targets: []Target{
					&SetName{Text: "biz"},
					&SetName{Text: "baz"},
				},
			},
		}},
		LocalInfo_Scope: &LocalInfo_Scope{},
	}
	b, w := text.BufferedCodeWriter()
	gen := &textGenerator{decl: decl}
	GenerateFunc(gen, decl, w)
	checkCode(b.String(), "func (o *Obj) foo(cond bool, names []string) (biz int, baz *int) {\n\tif cond {\n\t\t\"hello\"\n\t} else if !cond {\n\t\t\"goodbye\"\n\t} else {\n\t\t\"impossible\"\n\t}\n\tbiz, baz := bar(names), 7\n}\n", t)
}

func TestFile(t *testing.T) {
	file := &FileAST{
		Package: "foo",
		Imports: []*Import{
			&Import{
				Path: "some/other",
			},
			&Import{
				Name: "more",
				Path: "more/other",
			},
			&Import{
				Name: "x",
				Path: "x/other",
			},
		},
		Decls: []Decl{
			&StructDecl{
				Name: "Bar",
				Fields: []*FieldDecl{
					&FieldDecl{
						Name: "Baz",
						Type: &NameRef{
							Name: "other.Biz",
						},
					},
					&FieldDecl{
						Name: "BazXYZ",
						Type: &NameRef{
							Name: "more.Biz",
						},
					},
				},
			},
			&FuncDecl{
				Name: "F",
				Type: &FuncTypeRef{},
				Block: &Block{
					Body: []Stmt{
						&Var{Name: "foo", Type: &NameRef{Name: "int"}, Expr: &IntLiteral{Value: 7}},
						&Goto{Text: "block"},
						&Label{Text: "block"},
						&Return{},
					},
				},
				LocalInfo_Scope: &LocalInfo_Scope{},
			},
			&InterfaceDecl{
				Name: "I",
				Fields: []*FieldDecl{
					&FieldDecl{
						Name: "Touch",
						Type: &FuncTypeRef{},
					},
					&FieldDecl{
						Name: "Process",
						Type: &FuncTypeRef{
							Params: []*Param{
								&Param{Name: "inp", Type: &NameRef{Name: "int"}},
							},
							Results: []*Param{
								&Param{Name: "outp", Type: &NameRef{Name: "string"}},
							},
						},
					},
				},
			},
		},
	}

	b, w := text.BufferedCodeWriter()
	GenerateFile(file, w)
	expected := "package foo\n\nimport (\n\tmore \"more/other\"\n\t\"some/other\"\n\tx \"x/other\"\n)\n\ntype Bar struct {\n\tBaz    other.Biz\n\tBazXYZ more.Biz\n}\n\nfunc F() {\n\tvar foo int = 7\n\tgoto block\nblock:\n\treturn\n}\n\ntype I interface {\n\tTouch()\n\tProcess(inp int) (outp string)\n}\n"
	checkCode(b.String(), expected, t)

}
