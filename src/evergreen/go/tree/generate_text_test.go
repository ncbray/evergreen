package tree

import (
	"bytes"
	"evergreen/base"
	"fmt"
	"testing"
)

func bufferedWriter() (*bytes.Buffer, *base.CodeWriter) {
	b := &bytes.Buffer{}
	w := &base.CodeWriter{Out: b}
	return b, w
}

func checkString(actual string, expected string, t *testing.T) {
	if actual != expected {
		fmt.Println(actual)
		fmt.Println(expected)
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func genExpr(expr Expr) string {
	return GenerateExpr(expr)
}

func TestIntLiteral(t *testing.T) {
	expr := &IntLiteral{Value: 17}
	result := genExpr(expr)
	checkString(result, "17", t)
}

func TestBoolLiteral(t *testing.T) {
	expr := &BoolLiteral{Value: true}
	result := genExpr(expr)
	checkString(result, "true", t)
}

func TestStringLiteral(t *testing.T) {
	expr := &StringLiteral{Value: "foo\nbar"}
	result := genExpr(expr)
	checkString(result, "\"foo\\nbar\"", t)
}

func TestRuneLiteral(t *testing.T) {
	expr := &RuneLiteral{Value: 'x'}
	result := genExpr(expr)
	checkString(result, "'x'", t)
}

func TestSimpleUnaryExpr(t *testing.T) {
	expr := &UnaryExpr{
		Op:   "-",
		Expr: &IntLiteral{Value: 17},
	}
	result := genExpr(expr)
	checkString(result, "-17", t)
}

func TestSimpleBinaryExpr(t *testing.T) {
	expr := &BinaryExpr{
		Left:  &IntLiteral{Value: 12},
		Op:    "+",
		Right: &IntLiteral{Value: 34},
	}
	result := genExpr(expr)
	checkString(result, "12 + 34", t)
}

func TestSimpleCompoundExpr(t *testing.T) {
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
	checkString(result, "(12 + 34) * 56", t)
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
	checkString(result, "12 * 34 * 56", t)
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
	checkString(result, "12 * (34 * 56)", t)
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
	checkString(result, "foo.bar(12, 34)", t)
}

func TestIndex(t *testing.T) {
	expr := &Index{
		Expr:  &NameRef{Text: "foo"},
		Index: &NameRef{Text: "bar"},
	}
	result := genExpr(expr)
	checkString(result, "foo[bar]", t)
}

func TestFuncDecl(t *testing.T) {
	b, w := bufferedWriter()
	decl := &FuncDecl{
		Name: "foo",
		Body: []Stmt{
			&If{
				Cond: &NameRef{Text: "cond"},
				Body: []Stmt{
					&StringLiteral{Value: "hello"},
				},
			},
			&Assign{
				Sources: []Expr{
					&Call{
						Expr: &NameRef{Text: "bar"},
						Args: []Expr{},
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
	GenerateFunc(decl, w)
	checkString(b.String(), "func foo() {\n\tif cond {\n\t\t\"hello\"\n\t}\n\tbiz, baz := bar(), 7\n}\n", t)
}
