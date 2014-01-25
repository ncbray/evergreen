package tree

import (
	"evergreen/base"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
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

func GenerateType(t Type) string {
	switch t := t.(type) {
	case *TypeRef:
		return t.Name
	case *PointerType:
		return fmt.Sprintf("*%s", GenerateType(t.Element))
	case *SliceType:
		return fmt.Sprintf("[]%s", GenerateType(t.Element))
	default:
		panic(t)
	}
}

func GenerateBody(stmts []Stmt, w *base.CodeWriter) {
	w.PushMargin(indent)
	for _, stmt := range stmts {
		GenerateStmt(stmt, w)
	}
	w.PopMargin()
}

func GenerateParam(p *Param) string {
	t := GenerateType(p.T)
	if p.Name != "" {
		return fmt.Sprintf("%s %s", p.Name, t)
	} else {
		return t
	}
}

func GenerateReturns(returns []*Param) string {
	if len(returns) == 0 {
		return ""
	} else if len(returns) == 1 && returns[0].Name == "" {
		return " " + GenerateType(returns[0].T)
	} else {
		params := make([]string, len(returns))
		for i, p := range returns {
			params[i] = GenerateParam(p)
		}
		return fmt.Sprintf(" (%s)", strings.Join(params, ", "))
	}
}

func GenerateFunc(decl *FuncDecl, w *base.CodeWriter) {
	params := make([]string, len(decl.Params))
	for i, p := range decl.Params {
		params[i] = GenerateParam(p)
	}
	returns := GenerateReturns(decl.Returns)
	w.Linef("func %s(%s)%s {", decl.Name, strings.Join(params, ", "), returns)
	GenerateBody(decl.Body, w)
	w.Line("}")
}

func GenerateDecl(decl Decl, w *base.CodeWriter) {
	switch decl := decl.(type) {
	case *StructDecl:
		w.Linef("type %s struct {", decl.Name)
		w.PushMargin(indent)
		biggestName := 0
		for _, field := range decl.Fields {
			size := utf8.RuneCountInString(field.Name)
			if size > biggestName {
				biggestName = size
			}
		}
		for _, field := range decl.Fields {
			// Align the types
			padding := strings.Repeat(" ", biggestName-utf8.RuneCountInString(field.Name))
			w.Linef("%s%s %s", field.Name, padding, GenerateType(field.T))
		}
		w.PopMargin()
		w.Line("}")
	default:
		panic(decl)
	}
}

type ImportOrder []*Import

func (imports ImportOrder) Len() int {
	return len(imports)
}

func (imports ImportOrder) Swap(i, j int) {
	imports[i], imports[j] = imports[j], imports[i]
}

func (imports ImportOrder) Less(i, j int) bool {
	return imports[i].Path < imports[j].Path
}

func GenerateFile(file *File, w *base.CodeWriter) {
	w.Linef("package %s", file.Package)
	w.EmptyLines(1)
	if len(file.Imports) > 0 {
		w.Line("import (")
		w.PushMargin(indent)

		// Sort imports
		imports := make([]*Import, len(file.Imports))
		copy(imports, file.Imports)
		sort.Sort(ImportOrder(imports))

		for _, imp := range imports {
			path := strconv.Quote(imp.Path)
			if imp.Name != "" {
				w.Linef("%s %s", imp.Name, path)
			} else {
				w.Line(path)
			}
		}
		w.PopMargin()
		w.Line(")")
		w.EmptyLines(1)
	}
	for _, decl := range file.Decls {
		GenerateDecl(decl, w)
		w.EmptyLines(1)
	}
}
