package tree

import (
	"evergreen/text"
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

type textGenerator struct {
	decl *FuncDecl
}

func GeneratePrecExpr(gen *textGenerator, expr Expr) (string, int) {
	if expr == nil {
		panic("expr is nil")
	}

	switch expr := expr.(type) {
	case *IntLiteral:
		return strconv.Itoa(expr.Value), operandPrec
	case *Float32Literal:
		return strconv.FormatFloat(float64(expr.Value), 'g', -1, 32), operandPrec
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
	case *GetGlobal:
		return expr.Text, operandPrec
	case *GetLocal:
		info := expr.Info
		return info.Name, operandPrec
	case *GetName:
		return expr.Text, operandPrec
	case *UnaryExpr:
		return fmt.Sprintf("%s%s", expr.Op, GenerateSafeExpr(gen, expr.Expr, prefixPrec)), prefixPrec
	case *BinaryExpr:
		prec, ok := binaryOpToPrec[expr.Op]
		if !ok {
			panic(expr.Op)
		}
		return fmt.Sprintf("%s %s %s", GenerateSafeExpr(gen, expr.Left, prec), expr.Op, GenerateSafeExpr(gen, expr.Right, prec+1)), prec
	case *Selector:
		base := GenerateSafeExpr(gen, expr.Expr, postfixPrec)
		return fmt.Sprintf("%s.%s", base, expr.Text), postfixPrec
	case *Index:
		base := GenerateSafeExpr(gen, expr.Expr, postfixPrec)
		index := GenerateSafeExpr(gen, expr.Index, anyPrec)
		return fmt.Sprintf("%s[%s]", base, index), postfixPrec
	case *Call:
		base := GenerateSafeExpr(gen, expr.Expr, postfixPrec)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = GenerateSafeExpr(gen, arg, anyPrec)
		}
		return fmt.Sprintf("%s(%s)", base, strings.Join(args, ", ")), postfixPrec
	case *StructLiteral:
		t := GenerateType(expr.Type)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = fmt.Sprintf("%s: %s", arg.Name, GenerateSafeExpr(gen, arg.Expr, anyPrec))
		}
		return fmt.Sprintf("%s{%s}", t, strings.Join(args, ", ")), postfixPrec
	case *ListLiteral:
		t := GenerateType(expr.Type)
		args := make([]string, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = GenerateSafeExpr(gen, arg, anyPrec)
		}
		return fmt.Sprintf("%s{%s}", t, strings.Join(args, ", ")), postfixPrec
	case *TypeAssert:
		base := GenerateSafeExpr(gen, expr.Expr, postfixPrec)
		t := GenerateType(expr.Type)
		return fmt.Sprintf("%s.(%s)", base, t), postfixPrec
	case *TypeCoerce:
		t := GenerateType(expr.Type)
		e := GenerateSafeExpr(gen, expr.Expr, anyPrec)
		return fmt.Sprintf("%s(%s)", t, e), postfixPrec
	default:
		panic(expr)
	}
}

func GenerateSafeExpr(gen *textGenerator, expr Expr, requiredPrec int) string {
	result, actualPrec := GeneratePrecExpr(gen, expr)
	if requiredPrec > actualPrec {
		result = fmt.Sprintf("(%s)", result)
	}
	return result
}

func GenerateExpr(gen *textGenerator, expr Expr) string {
	result, _ := GeneratePrecExpr(gen, expr)
	return result
}

func GenerateExprList(gen *textGenerator, exprs []Expr) string {
	parts := make([]string, len(exprs))
	for i, e := range exprs {
		parts[i] = GenerateExpr(gen, e)
	}
	return strings.Join(parts, ", ")
}

func GenerateTarget(gen *textGenerator, expr Target) string {
	switch expr := expr.(type) {
	case *SetLocal:
		info := expr.Info
		return info.Name
	case *SetName:
		return expr.Text
	default:
		panic(expr)
	}
}

func GenerateTargetList(gen *textGenerator, exprs []Target) string {
	parts := make([]string, len(exprs))
	for i, e := range exprs {
		parts[i] = GenerateTarget(gen, e)
	}
	return strings.Join(parts, ", ")
}

func Dedent(w *text.CodeWriter) {
	margin := w.GetMargin()
	w.SetMargin(margin[:len(margin)-1])
}

func GenerateStmt(gen *textGenerator, stmt Stmt, w *text.CodeWriter) {
	expr, ok := stmt.(Expr)
	if ok {
		w.Line(GenerateExpr(gen, expr))
		return
	}
	switch stmt := stmt.(type) {
	case *BlockStmt:
		w.Line("{")
		GenerateBody(gen, stmt.Block, w)
		w.Line("}")
	case *If:
		w.Linef("if %s {", GenerateExpr(gen, stmt.Cond))
		GenerateBody(gen, stmt.T, w)
		next := stmt.F
		for next != nil && len(next.Body) > 0 {
			if len(next.Body) == 1 {
				switch stmt := next.Body[0].(type) {
				case *If:
					w.Linef("} else if %s {", GenerateExpr(gen, stmt.Cond))
					GenerateBody(gen, stmt.T, w)
					next = stmt.F
					continue
				}
			}
			w.Line("} else {")
			GenerateBody(gen, next, w)
			next = nil
		}
		w.Line("}")
	case *For:
		w.Linef("for {")
		GenerateBody(gen, stmt.Block, w)
		w.Line("}")
	case *Assign:
		sources := GenerateExprList(gen, stmt.Sources)
		targets := GenerateTargetList(gen, stmt.Targets)
		w.Linef("%s %s %s", targets, stmt.Op, sources)
	case *Var:
		t := GenerateType(stmt.Type)
		if stmt.Expr != nil {
			w.Linef("var %s %s = %s", stmt.Name, t, GenerateExpr(gen, stmt.Expr))
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
			w.Linef("return %s", GenerateExprList(gen, stmt.Args))
		} else {
			w.Line("return")
		}
	default:
		panic(stmt)
	}
}

func GenerateType(t TypeRef) string {
	switch t := t.(type) {
	case *NameRef:
		return t.Name
	case *PointerRef:
		return fmt.Sprintf("*%s", GenerateType(t.Element))
	case *SliceRef:
		return fmt.Sprintf("[]%s", GenerateType(t.Element))
	case *FuncTypeRef:
		return GenerateFuncType(t)
	default:
		panic(t)
	}
}

func generateBlock(gen *textGenerator, block *Block, w *text.CodeWriter) {
	for _, stmt := range block.Body {
		GenerateStmt(gen, stmt, w)
	}
}

func GenerateBody(gen *textGenerator, block *Block, w *text.CodeWriter) {
	w.AppendMargin(indent)
	generateBlock(gen, block, w)
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

func GenerateFuncType(t *FuncTypeRef) string {
	params := make([]string, len(t.Params))
	for i, p := range t.Params {
		params[i] = GenerateParam(p)
	}
	returns := GenerateReturns(t.Results)
	return fmt.Sprintf("(%s)%s", strings.Join(params, ", "), returns)
}

func GenerateFunc(gen *textGenerator, decl *FuncDecl, w *text.CodeWriter) {
	recv := ""
	if decl.Recv != nil {
		recv = fmt.Sprintf("(%s %s) ", decl.Recv.Name, GenerateType(decl.Recv.Type))
	}
	t := GenerateFuncType(decl.Type)
	w.Linef("func %s%s%s {", recv, decl.Name, t)
	GenerateBody(gen, decl.Block, w)
	w.Line("}")
}

func GenerateDecl(decl Decl, w *text.CodeWriter) {
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
	case *TypeDefDecl:
		w.Linef("type %s %s", decl.Name, GenerateType(decl.Type))
	case *FuncDecl:
		gen := &textGenerator{decl: decl}
		GenerateFunc(gen, decl, w)
	case *VarDecl:
		gen := &textGenerator{decl: nil} // HACK
		keyword := "var"
		if decl.Const {
			keyword = "const"
		}
		w.Linef("%s %s = %s", keyword, decl.Name, GenerateExpr(gen, decl.Expr))
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

func NeedsName(imp *Import) bool {
	if imp.Name != "" {
		parts := strings.Split(imp.Path, "/")
		name := parts[len(parts)-1]
		return name != imp.Name
	}
	return false
}

func GenerateFile(file *FileAST, w *text.CodeWriter) {
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
			if NeedsName(imp) {
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
