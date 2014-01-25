package tree

import (
	"evergreen/base"
	"fmt"
	"strconv"
	"strings"
)

var binaryOpToPrec map[string]int = map[string]int{
	"*":  5,
	"/":  5,
	"%":  5,
	"<<": 5,
	">>": 5,
	"&":  5,
	"&^": 5,

	"+": 4,
	"-": 4,
	"|": 4,
	"^": 4,

	"==": 3,
	"!=": 3,
	"<":  3,
	"<=": 3,
	">":  3,
	">=": 3,

	"&&": 2,

	"||": 1,
}

const operandPrec = 8
const postfixPrec = 7
const prefixPrec = 6
const anyPrec = 0

const indent = "\t"

func GeneratePrecExpr(expr Expr) (string, int) {
	if expr == nil {
		panic("expr is nil")
	}

	switch expr := expr.(type) {
	case *IntLiteral:
		return strconv.Itoa(expr.Value), operandPrec
	case *BoolLiteral:
		if expr.Value {
			return "true", operandPrec
		} else {
			return "false", operandPrec
		}
	case *StringLiteral:
		return strconv.Quote(expr.Value), operandPrec
	case *RuneLiteral:
		return strconv.QuoteRune(expr.Value), operandPrec
	case *NameRef:
		return expr.Text, operandPrec
	case *UnaryExpr:
		return fmt.Sprintf("%s%s", expr.Op, GenerateSafeExpr(expr.Expr, prefixPrec)), prefixPrec
	case *BinaryExpr:
		prec, ok := binaryOpToPrec[expr.Op]
		if !ok {
			panic(expr.Op)
		}
		return fmt.Sprintf("%s %s %s", GenerateSafeExpr(expr.Left, prec), expr.Op, GenerateSafeExpr(expr.Right, prec+1)), prec
	case *Selector:
		base := GenerateSafeExpr(expr.Expr, postfixPrec)
		return fmt.Sprintf("%s.%s", base, expr.Text), postfixPrec
	case *Index:
		base := GenerateSafeExpr(expr.Expr, postfixPrec)
		index := GenerateSafeExpr(expr.Index, anyPrec)
		return fmt.Sprintf("%s[%s]", base, index), postfixPrec
	case *Call:
		base := GenerateSafeExpr(expr.Expr, postfixPrec)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = GenerateSafeExpr(arg, anyPrec)
		}
		return fmt.Sprintf("%s(%s)", base, strings.Join(args, ", ")), postfixPrec
	default:
		panic(expr)
	}
}

func GenerateSafeExpr(expr Expr, requiredPrec int) string {
	result, actualPrec := GeneratePrecExpr(expr)
	if requiredPrec > actualPrec {
		result = fmt.Sprintf("(%s)", result)
	}
	return result
}

func GenerateExpr(expr Expr) string {
	result, _ := GeneratePrecExpr(expr)
	return result
}

func GenerateStmt(stmt Stmt, w *base.CodeWriter) {
	expr, ok := stmt.(Expr)
	if ok {
		w.Line(GenerateExpr(expr))
		return
	}
	switch stmt := stmt.(type) {
	case *If:
		w.Linef("if %s {", GenerateExpr(stmt.Cond))
		GenerateBody(stmt.Body, w)
		w.Linef("}")
	case *Assign:
		sources := make([]string, len(stmt.Sources))
		for i, src := range stmt.Sources {
			sources[i] = GenerateExpr(src)
		}
		targets := make([]string, len(stmt.Targets))
		for i, tgt := range stmt.Targets {
			targets[i] = GenerateExpr(tgt)
		}
		w.Linef("%s %s %s", strings.Join(targets, ", "), stmt.Op, strings.Join(sources, ", "))
	default:
		panic(stmt)
	}
}

func GenerateBody(stmts []Stmt, w *base.CodeWriter) {
	w.PushMargin(indent)
	for _, stmt := range stmts {
		GenerateStmt(stmt, w)
	}
	w.PopMargin()
}

func GenerateFunc(decl *FuncDecl, w *base.CodeWriter) {
	w.Linef("func %s() {", decl.Name)
	GenerateBody(decl.Body, w)
	w.Line("}")
}
