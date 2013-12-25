package main

import (
	"bytes"
	"evergreen/base"
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"text/scanner"
)

type ASTExpr interface {
	isASTExpr()
}

type ASTType interface {
	DubType() dub.DubType
	isASTType()
}

type TypeRef struct {
	Name string
	T    dub.DubType
}

func (node *TypeRef) DubType() dub.DubType {
	return node.T
}

func (node *TypeRef) isASTType() {
}

type If struct {
	Expr  ASTExpr
	Block []ASTExpr
}

func (node *If) isASTExpr() {
}

type Repeat struct {
	Block []ASTExpr
	Min   int
}

func (node *Repeat) isASTExpr() {
}

type Slice struct {
	Block []ASTExpr
}

func (node *Slice) isASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     dub.DubType
}

func (node *BinaryOp) isASTExpr() {
}

type Call struct {
	Name string
	T    dub.DubType
}

func (node *Call) isASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) isASTExpr() {
}

type RuneLiteral struct {
	Value rune
}

func (node *RuneLiteral) isASTExpr() {
}

type SetName struct {
	Expr ASTExpr
	Name string
	Info int
}

func (node *SetName) isASTExpr() {
}

type Read struct {
}

func (node *Read) isASTExpr() {
}

type Return struct {
	Exprs []ASTExpr
}

func (node *Return) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
}

type LocalInfo struct {
	Name string
	T    dub.DubType
}

type Decl interface {
	isASTDecl()
}

type FuncDecl struct {
	Name        string
	ReturnTypes []ASTType
	Block       []ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) isASTDecl() {
}

type FieldDecl struct {
	Name string
	Type ASTType
}

type StructDecl struct {
	Name   string
	Fields []*FieldDecl
}

func (node *StructDecl) isASTDecl() {
}

type DASMTokenType int

const (
	Ident DASMTokenType = iota
	Int
	Char
	Punc
	EOF
)

type DASMToken struct {
	Type DASMTokenType
	Text string
	Pos  scanner.Position
}

type DASMScanner struct {
	scanner *scanner.Scanner
	Current DASMToken
	Next    DASMToken
}

func (s *DASMScanner) Scan() {
	s.Current = s.Next
	tok := s.scanner.Scan()
	s.Next.Text = s.scanner.TokenText()
	s.Next.Pos = s.scanner.Pos()
	switch tok {
	case scanner.Ident:
		s.Next.Type = Ident
	case scanner.Int:
		s.Next.Type = Int
	case scanner.Char:
		s.Next.Type = Char
	case scanner.EOF:
		s.Next.Type = EOF
	default:
		if tok > 0 {
			s.Next.Type = Punc
		} else {
			panic(tok)
		}
	}
}

func (s *DASMScanner) AssertType(t DASMTokenType) {
	if s.Current.Type != t {
		panic(s.Current.Type)
	}
}

func CreateScanner(data []byte) *DASMScanner {
	s := &DASMScanner{scanner: &scanner.Scanner{}}
	s.scanner.Init(bytes.NewReader(data))
	s.Scan()
	s.Scan()
	return s
}

func getName(s *DASMScanner) (string, bool) {
	if s.Current.Type == Ident {
		text := s.Current.Text
		s.Scan()
		return text, true
	}
	return "", false
}

func getPunc(s *DASMScanner, text string) bool {
	if s.Current.Type == Punc && s.Current.Text == text {
		s.Scan()
		return true
	}
	return false
}

func getInt(s *DASMScanner) (int, bool) {
	if s.Current.Type == Int {
		count, _ := strconv.Atoi(s.Current.Text)
		s.Scan()
		return count, true
	}
	return 0, false
}

func parseExprList(s *DASMScanner) ([]ASTExpr, bool) {
	ok := getPunc(s, "(")
	if !ok {
		return nil, false
	}
	exprs := []ASTExpr{}
	for {
		if getPunc(s, ")") {
			return exprs, true
		}
		e, ok := parseExpr(s)
		if !ok {
			return nil, false
		}
		exprs = append(exprs, e)
	}
}

func parseTypeList(s *DASMScanner) ([]ASTType, bool) {
	ok := getPunc(s, "(")
	if !ok {
		return nil, false
	}
	types := []ASTType{}
	for {
		if getPunc(s, ")") {
			return types, true
		}
		t, ok := parseType(s)
		if !ok {
			return nil, false
		}
		types = append(types, t)
	}
}

var nameToOp = map[string]string{
	"eq": "==",
	"ne": "!=",
	"gt": ">",
	"lt": "<",
}

func parseType(s *DASMScanner) (ASTType, bool) {
	switch s.Current.Type {
	case Ident:
		result := &TypeRef{Name: s.Current.Text}
		s.Scan()
		return result, true
	default:
		panic(s.Current.Type)
	}
}

func parseExpr(s *DASMScanner) (ASTExpr, bool) {
	switch s.Current.Type {
	case Ident:
		switch s.Current.Text {
		case "star":
			s.Scan()
			block, ok := parseCodeBlock(s)
			if !ok {
				return nil, false
			}
			return &Repeat{Block: block, Min: 0}, true
		case "plus":
			s.Scan()
			block, ok := parseCodeBlock(s)
			if !ok {
				return nil, false
			}
			return &Repeat{Block: block, Min: 1}, true
		case "slice":
			s.Scan()
			block, ok := parseCodeBlock(s)
			if !ok {
				return nil, false
			}
			return &Slice{Block: block}, true
		case "if":
			s.Scan()
			expr, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			block, ok := parseCodeBlock(s)
			if !ok {
				return nil, false
			}
			return &If{Expr: expr, Block: block}, true
		case "define":
			s.Scan()
			expr, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			dst, ok := getName(s)
			if !ok {
				return nil, false
			}
			return &SetName{Expr: expr, Name: dst}, true
		case "read":
			s.Scan()
			return &Read{}, true
		case "fail":
			s.Scan()
			return &Fail{}, true
		case "eq", "ne", "gt", "lt":
			op := nameToOp[s.Current.Text]
			s.Scan()
			l, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			r, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			return &BinaryOp{Left: l, Op: op, Right: r}, true
		case "call":
			s.Scan()
			name, ok := getName(s)
			if !ok {
				return nil, false
			}
			return &Call{Name: name}, true
		case "return":
			s.Scan()
			exprs, ok := parseExprList(s)
			if !ok {
				return nil, false
			}
			return &Return{Exprs: exprs}, true
		default:
			text := s.Current.Text
			s.Scan()
			return &GetName{Name: text}, true
		}
	case Char:
		v, _ := strconv.Unquote(s.Current.Text)
		s.Scan()
		return &RuneLiteral{Value: []rune(v)[0]}, true
	default:
		return nil, false
	}
}

func parseCodeBlock(s *DASMScanner) ([]ASTExpr, bool) {
	ok := getPunc(s, "{")
	if !ok {
		return nil, false
	}
	result := []ASTExpr{}
	for {
		if s.Current.Type == Punc && s.Current.Text == "}" {
			s.Scan()
			return result, true
		}

		expr, ok := parseExpr(s)
		if !ok {
			return nil, false
		}
		result = append(result, expr)
		for s.Current.Type == Punc && s.Current.Text == ";" {
			s.Scan()
		}
	}
}

func parseFunction(s *DASMScanner) (*FuncDecl, bool) {
	name, ok := getName(s)
	if !ok {
		return nil, false
	}
	returnTypes, ok := parseTypeList(s)
	if !ok {
		return nil, false
	}
	block, ok := parseCodeBlock(s)
	if !ok {
		return nil, false
	}
	return &FuncDecl{Name: name, ReturnTypes: returnTypes, Block: block}, true
}

func parseStructure(s *DASMScanner) (*StructDecl, bool) {
	name, ok := getName(s)
	if !ok {
		return nil, false
	}

	ok = getPunc(s, "{")
	if !ok {
		return nil, false
	}

	fields := []*FieldDecl{}
	for {
		if getPunc(s, "}") {
			return &StructDecl{Name: name, Fields: fields}, true
		}

		name, ok := getName(s)
		if !ok {
			return nil, false
		}

		t, ok := parseType(s)
		if !ok {
			return nil, false
		}
		fields = append(fields, &FieldDecl{Name: name, Type: t})
	}
}

func parseFile(s *DASMScanner) ([]Decl, bool) {
	decls := []Decl{}
	for {
		switch s.Current.Type {
		case Ident:
			switch s.Current.Text {
			case "func":
				s.Scan()
				f, ok := parseFunction(s)
				if !ok {
					return nil, false
				}
				decls = append(decls, f)
			case "struct":
				s.Scan()
				f, ok := parseStructure(s)
				if !ok {
					return nil, false
				}
				decls = append(decls, f)
			default:
				panic(s.Current.Text)
			}
		case EOF:
			return decls, true
		default:
			return nil, false
		}
	}
}

func parseDASM(filename string) []Decl {
	data, _ := ioutil.ReadFile(filename)
	s := CreateScanner(data)
	f, ok := parseFile(s)
	if !ok {
		fmt.Printf("Unexpected %s @ %s\n", s.Current.Text, s.Current.Pos)
		panic(s.Current.Pos)
	}
	return f
}

type semanticScope struct {
	parent *semanticScope
	locals map[string]int
}

func (scope *semanticScope) localInfo(name string) (int, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return -1, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{parent: scope, locals: map[string]int{}}
}

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope) dub.DubType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope)
		return &dub.VoidType{}
	case *If:
		semanticExprPass(decl, expr.Expr, scope)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope))
		return &dub.VoidType{}
	case *BinaryOp:
		semanticExprPass(decl, expr.Left, scope)
		semanticExprPass(decl, expr.Right, scope)
		// HACK assume compare
		t := &dub.BoolType{}
		expr.T = t
		return t
	case *GetName:
		info, found := scope.localInfo(expr.Name)
		if !found {
			panic(expr.Name)
		}
		expr.Info = info
		return decl.Locals[info].T
	case *SetName:
		t := semanticExprPass(decl, expr.Expr, scope)
		info := len(decl.Locals)
		decl.Locals = append(decl.Locals, &LocalInfo{Name: expr.Name, T: t})
		scope.locals[expr.Name] = info
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope)
		return &dub.StringType{}
	case *Read:
		return &dub.RuneType{}
	case *RuneLiteral:
		return &dub.RuneType{}
	case *Return:
		return &dub.VoidType{}
	case *Fail:
		return &dub.VoidType{}
	case *Call:
		// HACK need to infer actual return type
		t := &dub.StringType{}
		expr.T = t
		return t
	default:
		panic(expr)
	}
}

func resolveTypeName(name string) dub.DubType {
	switch name {
	case "string":
		return &dub.StringType{}
	default:
		panic(name)
	}
}

func semanticTypePass(node ASTType) dub.DubType {
	switch node := node.(type) {
	case *TypeRef:
		t := resolveTypeName(node.Name)
		node.T = t
		return t
	default:
		panic(node)
	}
}

func semanticBlockPass(decl *FuncDecl, block []ASTExpr, scope *semanticScope) {
	for _, expr := range block {
		semanticExprPass(decl, expr, scope)
	}
}

func semanticFuncPass(decl *FuncDecl) {
	for _, t := range decl.ReturnTypes {
		semanticTypePass(t)
	}
	semanticBlockPass(decl, decl.Block, childScope(nil))
}

func semanticStructPass(decl *StructDecl) {
	for _, f := range decl.Fields {
		semanticTypePass(f.Type)
	}
}

func semanticPass(decls []Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			semanticFuncPass(decl)
		case *StructDecl:
			semanticStructPass(decl)
		default:
			panic(decl)
		}
	}
}

type DubBuilder struct {
	decl      *FuncDecl
	registers []dub.RegisterInfo
	localMap  []dub.DubRegister
}

func (builder *DubBuilder) CreateRegister(t dub.DubType) dub.DubRegister {
	builder.registers = append(builder.registers, dub.RegisterInfo{T: t})
	return dub.DubRegister(len(builder.registers) - 1)
}

func lowerExpr(expr ASTExpr, r *base.Region, builder *DubBuilder, used bool) dub.DubRegister {
	switch expr := expr.(type) {
	case *If:
		// TODO Min
		//l := dub.CreateRegion()

		cond := lowerExpr(expr.Expr, r, builder, true)
		block := lowerBlock(expr.Block, builder)

		// TODO conditional
		decide := dub.CreateSwitch(cond)

		r.Connect(dub.NORMAL, decide)
		decide.SetExit(0, r.GetExit(dub.NORMAL))
		r.Splice(dub.NORMAL, block)
		decide.SetExit(1, r.GetExit(dub.NORMAL))

		return dub.NoRegister

	case *Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			block := lowerBlock(expr.Block, builder)
			r.Splice(dub.NORMAL, block)
		}

		// Checkpoint
		checkpoint := builder.CreateRegister(&dub.IntType{})
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})

		r.Connect(dub.NORMAL, head)
		head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Normal flow iterates
		// NOTE actually connects nodes in two different regions.  Kinda hackish.
		block.GetExit(dub.NORMAL).TransferEntries(head)
		// Stop iterating on failure
		block.GetExit(dub.FAIL).TransferEntries(block.GetExit(dub.NORMAL))

		// Recover
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			block.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, block.GetExit(dub.NORMAL))
		}

		r.Splice(dub.NORMAL, block)

		return dub.NoRegister

	case *GetName:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateRegister(builder.decl.Locals[expr.Info].T)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CopyOp{Src: builder.localMap[expr.Info], Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *SetName:
		src := lowerExpr(expr.Expr, r, builder, true)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CopyOp{Src: src, Dst: builder.localMap[expr.Info]},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dub.NoRegister

	case *RuneLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateRegister(&dub.RuneType{})
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Read:
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(&dub.RuneType{})
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.Read{Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		r.AttachDefaultExits(body)
		return dst

	case *Return:
		exprs := make([]dub.DubRegister, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, r, builder, true)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ReturnOp{Exprs: exprs},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.RETURN))
		return dub.NoRegister

	case *Fail:
		body := dub.CreateBlock([]dub.DubOp{
			&dub.Fail{},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.FAIL, r.GetExit(dub.FAIL))

		return dub.NoRegister

	case *BinaryOp:
		left := lowerExpr(expr.Left, r, builder, true)
		right := lowerExpr(expr.Right, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.BinaryOp{
				Left:  left,
				Op:    expr.Op,
				Right: right,
				Dst:   dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst
	case *Call:
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CallOp{
				Name: expr.Name,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		r.AttachDefaultExits(body)
		return dst
	case *Slice:
		start := builder.CreateRegister(&dub.IntType{})
		// HACK assume checkpoint is just the index
		{
			head := dub.CreateBlock([]dub.DubOp{
				&dub.Checkpoint{Dst: start},
			})
			r.Connect(dub.NORMAL, head)
			head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		}
		block := lowerBlock(expr.Block, builder)
		r.Splice(dub.NORMAL, block)

		// Create a slice
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(&dub.StringType{})
		}
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Slice{Src: start, Dst: dst},
			})

			r.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		}
		return dst
	default:
		panic(expr)
	}

}

func lowerBlock(block []ASTExpr, builder *DubBuilder) *base.Region {
	r := dub.CreateRegion()
	for _, expr := range block {
		lowerExpr(expr, r, builder, false)
	}
	return r
}

func lowerAST(decl *FuncDecl) *dub.LLFunc {
	f := &dub.LLFunc{Name: decl.Name}
	types := make([]dub.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		switch node := node.(type) {
		case *TypeRef:
			types[i] = node.T
		default:
			panic(node)
		}
	}
	f.ReturnTypes = types
	builder := &DubBuilder{decl: decl}
	// Allocate register for locals
	builder.localMap = make([]dub.DubRegister, len(decl.Locals))
	for i, info := range decl.Locals {
		builder.localMap[i] = builder.CreateRegister(info.T)
	}
	f.Region = lowerBlock(decl.Block, builder)
	f.Region.GetExit(dub.RETURN).TransferEntries(f.Region.GetExit(dub.NORMAL))
	f.Registers = builder.registers
	return f
}

func lowerStruct(decl *StructDecl) *dub.LLStruct {
	fields := []*dub.LLField{}
	for _, field := range decl.Fields {
		fields = append(fields, &dub.LLField{
			Name: field.Name,
			T:    field.Type.DubType(),
		})
	}
	return &dub.LLStruct{Name: decl.Name, Fields: fields}
}

func main() {
	decls := parseDASM("dasm/math.dasm")
	semanticPass(decls)

	structs := []*dub.LLStruct{}
	funcs := []*dub.LLFunc{}
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			f := lowerAST(decl)
			funcs = append(funcs, f)

			// Dump flowgraph
			dot := base.RegionToDot(f.Region)
			outfile := filepath.Join("output", fmt.Sprintf("%s.svg", f.Name))
			io.WriteDot(dot, outfile)
		case *StructDecl:
			structs = append(structs, lowerStruct(decl))
		default:
			panic(decl)
		}
	}

	code := dub.GenerateGo("math", structs, funcs)
	fmt.Println(code)
	io.WriteFile("src/generated/math/parser.go", []byte(code))
}
