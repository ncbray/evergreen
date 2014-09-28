package tree

import (
	"evergreen/base"
	"fmt"
	"go/format"
	"testing"
)

func checkString(actual string, expected string, t *testing.T) {
	if actual != expected {
		fmt.Println(actual)
		fmt.Println(expected)
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func checkFormat(code string, t *testing.T) {
	out, err := format.Source([]byte(code))
	if err != nil {
		t.Fatal(err)
	}
	outs := string(out)
	checkString(code, outs, t)
}

func checkCode(actual string, expected string, t *testing.T) {
	checkString(actual, expected, t)
	checkFormat(actual, t)
}

func genExpr(expr Expr) string {
	return GenerateExpr(expr)
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
			Expr: &NameRef{Text: "foo"},
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
		Expr:  &NameRef{Text: "foo"},
		Index: &NameRef{Text: "bar"},
	}
	result := genExpr(expr)
	checkCode(result, "foo[bar]", t)
}

func TestTypeCoerce(t *testing.T) {
	expr := &TypeCoerce{
		Type: &SliceType{Element: &TypeRef{Name: "rune"}},
		Expr: &NameRef{Text: "s"},
	}
	result := genExpr(expr)
	checkCode(result, "[]rune(s)", t)
}

func TestStructLiteral(t *testing.T) {
	expr := &StructLiteral{
		Type: &TypeRef{Name: "Foo"},
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
		Type: &SliceType{Element: &TypeRef{Name: "int"}},
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
		Recv: &Param{Name: "o", Type: &PointerType{Element: &TypeRef{Name: "Obj"}}},
		Type: &FuncType{
			Params: []*Param{
				&Param{Name: "cond", Type: &TypeRef{Name: "bool"}},
				&Param{Name: "names", Type: &SliceType{Element: &TypeRef{Name: "string"}}},
			},
			Results: []*Param{
				&Param{Name: "biz", Type: &TypeRef{Name: "int"}},
				&Param{Name: "baz", Type: &PointerType{Element: &TypeRef{Name: "int"}}},
			},
		},
		Body: []Stmt{
			&If{
				Cond: &NameRef{Text: "cond"},
				Body: []Stmt{
					&StringLiteral{Value: "hello"},
				},
				Else: &If{
					Cond: &UnaryExpr{Op: "!", Expr: &NameRef{Text: "cond"}},
					Body: []Stmt{
						&StringLiteral{Value: "goodbye"},
					},
					Else: &BlockStmt{
						Body: []Stmt{
							&StringLiteral{Value: "impossible"},
						},
					},
				},
			},
			&Assign{
				Sources: []Expr{
					&Call{
						Expr: &NameRef{Text: "bar"},
						Args: []Expr{&NameRef{Text: "names"}},
					},
					&IntLiteral{Value: 7},
				},
				Op: ":=",
				Targets: []Expr{
					&NameRef{Text: "biz"},
					&NameRef{Text: "baz"},
				},
			},
		},
	}
	b, w := base.BufferedCodeWriter()
	GenerateFunc(decl, w)
	checkCode(b.String(), "func (o *Obj) foo(cond bool, names []string) (biz int, baz *int) {\n\tif cond {\n\t\t\"hello\"\n\t} else if !cond {\n\t\t\"goodbye\"\n\t} else {\n\t\t\"impossible\"\n\t}\n\tbiz, baz := bar(names), 7\n}\n", t)
}

func TestFile(t *testing.T) {
	file := &File{
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
				Fields: []*Field{
					&Field{
						Name: "Baz",
						Type: &TypeRef{
							Name: "other.Biz",
						},
					},
					&Field{
						Name: "BazXYZ",
						Type: &TypeRef{
							Name: "more.Biz",
						},
					},
				},
			},
			&FuncDecl{
				Name: "F",
				Type: &FuncType{},
				Body: []Stmt{
					&BlockStmt{
						Body: []Stmt{
							&Var{Name: "foo", Type: &TypeRef{Name: "int"}, Expr: &IntLiteral{Value: 7}},
							&Goto{Text: "block"},
							&Label{Text: "block"},
							&Return{},
						},
					},
				},
			},
			&InterfaceDecl{
				Name: "I",
				Fields: []*Field{
					&Field{
						Name: "Touch",
						Type: &FuncType{},
					},
					&Field{
						Name: "Process",
						Type: &FuncType{
							Params: []*Param{
								&Param{Name: "inp", Type: &TypeRef{Name: "int"}},
							},
							Results: []*Param{
								&Param{Name: "outp", Type: &TypeRef{Name: "string"}},
							},
						},
					},
				},
			},
		},
	}

	b, w := base.BufferedCodeWriter()
	GenerateFile(file, w)
	checkCode(b.String(), "package foo\n\nimport (\n\tmore \"more/other\"\n\t\"some/other\"\n\tx \"x/other\"\n)\n\ntype Bar struct {\n\tBaz    other.Biz\n\tBazXYZ more.Biz\n}\n\nfunc F() {\n\t{\n\t\tvar foo int = 7\n\t\tgoto block\n\tblock:\n\t\treturn\n\t}\n}\n\ntype I interface {\n\tTouch()\n\tProcess(inp int) (outp string)\n}\n", t)

}
