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
	case *NilLiteral:
		return "nil", operandPrec
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
	case *StructLiteral:
		t := GenerateType(expr.Type)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = fmt.Sprintf("%s: %s", arg.Name, GenerateSafeExpr(arg.Expr, anyPrec))
		}
		return fmt.Sprintf("%s{%s}", t, strings.Join(args, ", ")), postfixPrec
	case *ListLiteral:
		t := GenerateType(expr.Type)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = GenerateSafeExpr(arg, anyPrec)
		}
		return fmt.Sprintf("%s{%s}", t, strings.Join(args, ", ")), postfixPrec
	case *TypeAssert:
		base := GenerateSafeExpr(expr.Expr, postfixPrec)
		t := GenerateType(expr.Type)
		return fmt.Sprintf("%s.(%s)", base, t), postfixPrec
	case *TypeCoerce:
		t := GenerateType(expr.Type)
		e := GenerateSafeExpr(expr.Expr, anyPrec)
		return fmt.Sprintf("%s(%s)", t, e), postfixPrec
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

func GenerateExprList(exprs []Expr) string {
	gen := make([]string, len(exprs))
	for i, e := range exprs {
		gen[i] = GenerateExpr(e)
	}
	return strings.Join(gen, ", ")
}

func Dedent(w *base.CodeWriter) {
	margin := w.GetMargin()
	w.SetMargin(margin[:len(margin)-1])
}

func GenerateStmt(stmt Stmt, w *base.CodeWriter) {
	expr, ok := stmt.(Expr)
	if ok {
		w.Line(GenerateExpr(expr))
		return
	}
	switch stmt := stmt.(type) {
	case *BlockStmt:
		w.Line("{")
		GenerateBody(stmt.Body, w)
		w.Line("}")
	case *If:
		w.Linef("if %s {", GenerateExpr(stmt.Cond))
		GenerateBody(stmt.Body, w)
		next := stmt.Else
		for next != nil {
			switch stmt := next.(type) {
			case *If:
				w.Linef("} else if %s {", GenerateExpr(stmt.Cond))
				GenerateBody(stmt.Body, w)
				next = stmt.Else
			case *BlockStmt:
				w.Line("} else {")
				GenerateBody(stmt.Body, w)
				next = nil
			default:
				panic(next)
			}
		}
		w.Line("}")
	case *Assign:
		sources := GenerateExprList(stmt.Sources)
		targets := GenerateExprList(stmt.Targets)
		w.Linef("%s %s %s", targets, stmt.Op, sources)
	case *Var:
		t := GenerateType(stmt.Type)
		if stmt.Expr != nil {
			w.Linef("var %s %s = %s", stmt.Name, t, GenerateExpr(stmt.Expr))
		} else {
			w.Linef("var %s %s", stmt.Name, t)
		}
	case *Goto:
		w.Linef("goto %s", stmt.Text)
	case *Label:
		Dedent(w)
		w.Linef("%s:", stmt.Text)
		w.RestoreMargin()
	case *Return:
		if len(stmt.Args) > 0 {
			w.Linef("return %s", GenerateExprList(stmt.Args))
		} else {
			w.Line("return")
		}
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
	case *FuncType:
		return GenerateFuncType(t)
	default:
		panic(t)
	}
}

func GenerateBody(stmts []Stmt, w *base.CodeWriter) {
	w.AppendMargin(indent)
	for _, stmt := range stmts {
		GenerateStmt(stmt, w)
	}
	w.RestoreMargin()
}

func GenerateParam(p *Param) string {
	t := GenerateType(p.Type)
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
		return " " + GenerateType(returns[0].Type)
	} else {
		params := make([]string, len(returns))
		for i, p := range returns {
			params[i] = GenerateParam(p)
		}
		return fmt.Sprintf(" (%s)", strings.Join(params, ", "))
	}
}

func GenerateFuncType(t *FuncType) string {
	params := make([]string, len(t.Params))
	for i, p := range t.Params {
		params[i] = GenerateParam(p)
	}
	returns := GenerateReturns(t.Results)
	return fmt.Sprintf("(%s)%s", strings.Join(params, ", "), returns)
}

func GenerateFunc(decl *FuncDecl, w *base.CodeWriter) {
	recv := ""
	if decl.Recv != nil {
		recv = fmt.Sprintf("(%s %s) ", decl.Recv.Name, GenerateType(decl.Recv.Type))
	}
	t := GenerateFuncType(decl.Type)
	w.Linef("func %s%s%s {", recv, decl.Name, t)
	GenerateBody(decl.Body, w)
	w.Line("}")
}

func GenerateDecl(decl Decl, w *base.CodeWriter) {
	switch decl := decl.(type) {
	case *StructDecl:
		w.Linef("type %s struct {", decl.Name)
		w.AppendMargin(indent)
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
			w.Linef("%s%s %s", field.Name, padding, GenerateType(field.Type))
		}
		w.RestoreMargin()
		w.Line("}")
	case *InterfaceDecl:
		w.Linef("type %s interface {", decl.Name)
		w.AppendMargin(indent)
		for _, field := range decl.Fields {
			w.Linef("%s%s", field.Name, GenerateType(field.Type))
		}
		w.RestoreMargin()
		w.Line("}")
	case *FuncDecl:
		GenerateFunc(decl, w)
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
		w.AppendMargin(indent)

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
		w.RestoreMargin()
		w.Line(")")
		w.EmptyLines(1)
	}
	for _, decl := range file.Decls {
		GenerateDecl(decl, w)
		w.EmptyLines(1)
	}
}
