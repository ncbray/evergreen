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
	isASTType()
}

type ASTTypeRef interface {
	Resolve() ASTType
	isASTTypeRef()
}

type TypeRef struct {
	Name string
	T    ASTType
}

func (node *TypeRef) Resolve() ASTType {
	return node.T
}

func (node *TypeRef) isASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    ASTType
}

func (node *ListTypeRef) Resolve() ASTType {
	return node.T
}

func (node *ListTypeRef) isASTTypeRef() {
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

type Optional struct {
	Block []ASTExpr
}

func (node *Optional) isASTExpr() {
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
	T     ASTType
}

func (node *BinaryOp) isASTExpr() {
}

type Call struct {
	Name string
	T    ASTType
}

func (node *Call) isASTExpr() {
}

type KeyValue struct {
	Key   string
	Value ASTExpr
}

type Construct struct {
	Type ASTTypeRef
	Args []*KeyValue
}

func (node *Construct) isASTExpr() {
}

type ConstructList struct {
	Type ASTTypeRef
	Args []ASTExpr
}

func (node *ConstructList) isASTExpr() {
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

type StringLiteral struct {
	Value string
}

func (node *StringLiteral) isASTExpr() {
}

type Assign struct {
	Expr   ASTExpr
	Name   string
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) isASTExpr() {
}

type Read struct {
}

func (node *Read) isASTExpr() {
}

type Append struct {
	List  ASTExpr
	Value ASTExpr
	T     ASTType
}

func (node *Append) isASTExpr() {
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
	T    ASTType
}

type Decl interface {
	AsType() (ASTType, bool)
	AsFunc() (ASTFunc, bool)
	isASTDecl()
}

type ASTFunc interface {
	ReturnType() ASTType
}

type FuncDecl struct {
	Name        string
	ReturnTypes []ASTTypeRef
	Block       []ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) AsType() (ASTType, bool) {
	return nil, false
}

func (node *FuncDecl) AsFunc() (ASTFunc, bool) {
	return node, true
}

func (node *FuncDecl) ReturnType() ASTType {
	// HACK assume single return value
	if len(node.ReturnTypes) == 0 {
		return nil
	}
	if len(node.ReturnTypes) != 1 {
		panic(node.Name)
	}
	return node.ReturnTypes[0].Resolve()
}

func (node *FuncDecl) isASTDecl() {
}

type FieldDecl struct {
	Name string
	Type ASTTypeRef
}

type StructDecl struct {
	Name       string
	Implements ASTTypeRef
	Fields     []*FieldDecl
}

func (node *StructDecl) AsType() (ASTType, bool) {
	return node, true
}

func (node *StructDecl) AsFunc() (ASTFunc, bool) {
	return nil, false
}

func (node *StructDecl) isASTDecl() {
}

func (node *StructDecl) isASTType() {
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) AsType() (ASTType, bool) {
	return node, true
}

func (node *BuiltinType) AsFunc() (ASTFunc, bool) {
	return nil, false
}

func (node *BuiltinType) isASTDecl() {
}

func (node *BuiltinType) isASTType() {
}

type ListType struct {
	Type ASTType
}

func (node *ListType) isASTType() {
}

type DASMTokenType int

const (
	Ident DASMTokenType = iota
	Int
	Char
	String
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
	case scanner.String:
		s.Next.Type = String
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

func getKeyword(s *DASMScanner, text string) bool {
	if s.Current.Type == Ident && s.Current.Text == text {
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

func parseKeyValueList(s *DASMScanner) ([]*KeyValue, bool) {
	ok := getPunc(s, "(")
	if !ok {
		return nil, false
	}
	args := []*KeyValue{}
	for {
		if getPunc(s, ")") {
			return args, true
		}
		key, ok := getName(s)
		if !ok {
			return nil, false
		}
		if !getPunc(s, ":") {
			return nil, false
		}
		e, ok := parseExpr(s)
		if !ok {
			return nil, false
		}
		args = append(args, &KeyValue{Key: key, Value: e})
	}
}

func parseTypeList(s *DASMScanner) ([]ASTTypeRef, bool) {
	ok := getPunc(s, "(")
	if !ok {
		return nil, false
	}
	types := []ASTTypeRef{}
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

func parseType(s *DASMScanner) (ASTTypeRef, bool) {
	switch s.Current.Type {
	case Ident:
		result := &TypeRef{Name: s.Current.Text}
		s.Scan()
		return result, true
	case Punc:
		if s.Current.Text == "[" && s.Next.Text == "]" {
			s.Scan()
			s.Scan()
			child, ok := parseType(s)
			if !ok {
				return nil, false
			}
			return &ListTypeRef{Type: child}, true
		}
		return nil, false
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
		case "question":
			s.Scan()
			block, ok := parseCodeBlock(s)
			if !ok {
				return nil, false
			}
			return &Optional{Block: block}, true
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
		case "var":
			s.Scan()
			name, ok := getName(s)
			if !ok {
				return nil, false
			}

			t, ok := parseType(s)
			if !ok {
				return nil, false
			}

			var expr ASTExpr
			if getPunc(s, "=") {
				expr, ok = parseExpr(s)
				if !ok {
					return nil, false
				}
			}
			return &Assign{Expr: expr, Name: name, Type: t, Define: true}, true
		case "define":
			s.Scan()
			name, ok := getName(s)
			if !ok {
				return nil, false
			}
			expr, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			return &Assign{Expr: expr, Name: name, Define: true}, true
		case "assign":
			s.Scan()
			name, ok := getName(s)
			if !ok {
				return nil, false
			}
			expr, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			return &Assign{Expr: expr, Name: name, Define: false}, true
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
		case "cons":
			s.Scan()
			t, ok := parseType(s)
			if !ok {
				return nil, false
			}
			args, ok := parseKeyValueList(s)
			if !ok {
				return nil, false
			}
			return &Construct{Type: t, Args: args}, true
		case "conl":
			s.Scan()
			t, ok := parseType(s)
			if !ok {
				return nil, false
			}
			args, ok := parseExprList(s)
			if !ok {
				return nil, false
			}
			return &ConstructList{Type: t, Args: args}, true
		case "append":
			s.Scan()
			name, ok := getName(s)
			if !ok {
				return nil, false
			}
			expr, ok := parseExpr(s)
			if !ok {
				return nil, false
			}
			return &Assign{
				Expr: &Append{
					List: &GetName{
						Name: name,
					},
					Value: expr,
				},
				Name: name,
			}, true
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
	case String:
		v, _ := strconv.Unquote(s.Current.Text)
		s.Scan()
		return &StringLiteral{Value: v}, true
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

	var implements ASTTypeRef
	ok = getKeyword(s, "implements")
	if ok {
		implements, ok = parseType(s)
		if !ok {
			return nil, false
		}
	}

	ok = getPunc(s, "{")
	if !ok {
		return nil, false
	}

	fields := []*FieldDecl{}
	for {
		if getPunc(s, "}") {
			return &StructDecl{
				Name:       name,
				Implements: implements,
				Fields:     fields,
			}, true
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

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope) ASTType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.Void
	case *Optional:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.Void
	case *If:
		semanticExprPass(decl, expr.Expr, scope, glbls)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope), glbls)
		return glbls.Void
	case *BinaryOp:
		semanticExprPass(decl, expr.Left, scope, glbls)
		semanticExprPass(decl, expr.Right, scope, glbls)
		// HACK assume compare
		t := glbls.Bool
		expr.T = t
		return t
	case *GetName:
		info, found := scope.localInfo(expr.Name)
		if !found {
			panic(fmt.Sprintf("Could not resolve name %#v", expr.Name))
		}
		expr.Info = info
		return decl.Locals[info].T
	case *Assign:
		var t ASTType
		if expr.Expr != nil {
			t = semanticExprPass(decl, expr.Expr, scope, glbls)
		}
		if expr.Type != nil {
			t = semanticTypePass(expr.Type, glbls)
		}
		if t == nil {
			panic(fmt.Sprintf("Cannot infer the type of %#v", expr.Name))
		}
		var info int
		var exists bool
		if expr.Define {
			_, exists = scope.locals[expr.Name]
			if exists {
				panic(fmt.Sprintf("Tried to redefine %#v", expr.Name))
			}

			info = len(decl.Locals)
			decl.Locals = append(decl.Locals, &LocalInfo{Name: expr.Name, T: t})
			scope.locals[expr.Name] = info
		} else {
			info, exists = scope.locals[expr.Name]
			if !exists {
				panic(fmt.Sprintf("Tried to assign to unknown variable %#v", expr.Name))
			}
		}
		expr.Info = info
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.String
	case *Read:
		return glbls.Rune
	case *RuneLiteral:
		return glbls.Rune
	case *StringLiteral:
		return glbls.String
	case *Return:
		for _, e := range expr.Exprs {
			semanticExprPass(decl, e, scope, glbls)
		}
		return glbls.Void
	case *Fail:
		return glbls.Void
	case *Call:
		// HACK resolve other scopes?
		decl, ok := glbls.Module[expr.Name]
		if !ok {
			panic(expr.Name)
		}
		f, ok := decl.AsFunc()
		if !ok {
			panic(expr.Name)
		}
		t := f.ReturnType()
		expr.T = t
		return t
	case *Append:
		t := semanticExprPass(decl, expr.List, scope, glbls)
		semanticExprPass(decl, expr.Value, scope, glbls)
		expr.T = t
		return t
	case *Construct:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg.Value, scope, glbls)
		}
		return t
	case *ConstructList:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg, scope, glbls)
		}
		return t
	default:
		panic(expr)
	}
}

func semanticTypePass(node ASTTypeRef, glbls *ModuleScope) ASTType {
	switch node := node.(type) {
	case *TypeRef:
		d, ok := glbls.Module[node.Name]
		if !ok {
			d, ok = glbls.Builtin[node.Name]
		}
		if !ok {
			panic(node.Name)
		}
		t, ok := d.AsType()
		if !ok {
			panic(node.Name)
		}
		node.T = t
		return t
	case *ListTypeRef:
		t := semanticTypePass(node.Type, glbls)
		// TODO memoize list types
		node.T = &ListType{Type: t}
		return node.T
	default:
		panic(node)
	}
}

func semanticBlockPass(decl *FuncDecl, block []ASTExpr, scope *semanticScope, glbls *ModuleScope) {
	for _, expr := range block {
		semanticExprPass(decl, expr, scope, glbls)
	}
}

func semanticFuncPass(decl *FuncDecl, glbls *ModuleScope) {
	for _, t := range decl.ReturnTypes {
		semanticTypePass(t, glbls)
	}
	semanticBlockPass(decl, decl.Block, childScope(nil), glbls)
}

func semanticStructPass(decl *StructDecl, glbls *ModuleScope) {
	if decl.Implements != nil {
		semanticTypePass(decl.Implements, glbls)
	}
	for _, f := range decl.Fields {
		semanticTypePass(f.Type, glbls)
	}
}

type ModuleScope struct {
	Builtin map[string]Decl
	Module  map[string]Decl

	String *BuiltinType
	Rune   *BuiltinType
	Int    *BuiltinType
	Bool   *BuiltinType
	Void   *BuiltinType
}

func semanticPass(decls []Decl) *ModuleScope {
	glbls := &ModuleScope{
		Builtin: map[string]Decl{},
		Module:  map[string]Decl{},
	}
	glbls.String = &BuiltinType{"string"}
	glbls.Builtin["string"] = glbls.String

	glbls.Rune = &BuiltinType{"rune"}
	glbls.Builtin["rune"] = glbls.Rune

	glbls.Int = &BuiltinType{"int"}
	glbls.Builtin["int"] = glbls.Int

	glbls.Bool = &BuiltinType{"bool"}
	glbls.Builtin["bool"] = glbls.Bool

	glbls.Void = &BuiltinType{"void"}
	glbls.Builtin["void"] = glbls.Void

	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			glbls.Module[decl.Name] = decl
		case *StructDecl:
			glbls.Module[decl.Name] = decl
		default:
			panic(decl)
		}
	}
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			semanticFuncPass(decl, glbls)
		case *StructDecl:
			semanticStructPass(decl, glbls)
		default:
			panic(decl)
		}
	}
	return glbls
}

type GlobalDubBuilder struct {
	types  map[ASTType]dub.DubType
	String dub.DubType
	Rune   dub.DubType
	Int    dub.DubType
	Bool   dub.DubType
}

func (builder *GlobalDubBuilder) TranslateType(t ASTType) dub.DubType {
	switch t := t.(type) {
	case *StructDecl, *BuiltinType:
		dt, ok := builder.types[t]
		if !ok {
			panic(t)
		}
		return dt
	case *ListType:
		parent := builder.TranslateType(t.Type)
		// TODO memoize
		return &dub.ListType{Type: parent}
	default:
		panic(t)
	}
}

type DubBuilder struct {
	decl      *FuncDecl
	registers []dub.RegisterInfo
	localMap  []dub.DubRegister
	glbl      *GlobalDubBuilder
}

func (builder *DubBuilder) CreateRegister(t ASTType) dub.DubRegister {
	return builder.CreateLLRegister(builder.glbl.TranslateType(t))
}

func (builder *DubBuilder) CreateLLRegister(t dub.DubType) dub.DubRegister {
	builder.registers = append(builder.registers, dub.RegisterInfo{T: t})
	return dub.DubRegister(len(builder.registers) - 1)
}

func (builder *DubBuilder) ZeroRegister(dst dub.DubRegister) dub.DubOp {
	info := builder.registers[dst]
	switch info.T.(type) {
	case *dub.LLStruct:
		return &dub.ConstantNilOp{Dst: dst}
	default:
		panic(info.T)
	}
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

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})

		// Splice in the checkpoint as the first operation.
		// Note: block may be empty, but this code is carefully designed to work in that case.
		// Sure it's an infinite loop, but it would stange if that loop vanished.
		oldHead := block.Head()
		oldHead.TransferEntries(head)
		head.SetExit(dub.NORMAL, oldHead)

		// Normal flow iterates
		block.GetExit(dub.NORMAL).TransferEntries(head)

		// Stop iterating on failure and recover
		block.GetExit(dub.FAIL).TransferEntries(block.GetExit(dub.NORMAL))
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			block.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, block.GetExit(dub.NORMAL))
		}

		r.Splice(dub.NORMAL, block)

		return dub.NoRegister

	case *Optional:
		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})
		r.Connect(dub.NORMAL, head)
		head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))

		block := lowerBlock(expr.Block, builder)

		restore := dub.CreateBlock([]dub.DubOp{
			&dub.Recover{Src: checkpoint},
		})

		block.Connect(dub.FAIL, restore)
		restore.SetExit(dub.NORMAL, block.GetExit(dub.NORMAL))

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

	case *Assign:
		dst := builder.localMap[expr.Info]
		var op dub.DubOp
		if expr.Expr != nil {
			src := lowerExpr(expr.Expr, r, builder, true)
			op = &dub.CopyOp{Src: src, Dst: dst}
		} else {
			op = builder.ZeroRegister(dst)
		}
		body := dub.CreateBlock([]dub.DubOp{op})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *RuneLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Rune)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *StringLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.String)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantStringOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Read:
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(builder.glbl.Rune)
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
	case *Append:
		l := lowerExpr(expr.List, r, builder, true)
		v := lowerExpr(expr.Value, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}

		body := dub.CreateBlock([]dub.DubOp{
			&dub.AppendOp{
				List:  l,
				Value: v,
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
	case *Construct:
		args := make([]*dub.KeyValue, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = &dub.KeyValue{
				Key:   arg.Key,
				Value: lowerExpr(arg.Value, r, builder, true),
			}
		}
		t := builder.glbl.TranslateType(expr.Type.Resolve())
		s, ok := t.(*dub.LLStruct)
		if !ok {
			panic(t)
		}
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstructOp{
				Type: s,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *ConstructList:
		args := make([]dub.DubRegister, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, r, builder, true)
		}
		t := builder.glbl.TranslateType(expr.Type.Resolve())
		l, ok := t.(*dub.ListType)
		if !ok {
			panic(t)
		}
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstructListOp{
				Type: l,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Slice:
		start := builder.CreateLLRegister(builder.glbl.Int)
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
			dst = builder.CreateLLRegister(builder.glbl.String)
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

func lowerAST(decl *FuncDecl, glbl *GlobalDubBuilder) *dub.LLFunc {
	builder := &DubBuilder{decl: decl, glbl: glbl}

	f := &dub.LLFunc{Name: decl.Name}
	types := make([]dub.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = builder.glbl.TranslateType(node.Resolve())
	}
	f.ReturnTypes = types
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

func lowerStruct(decl *StructDecl, s *dub.LLStruct, gbuilder *GlobalDubBuilder) *dub.LLStruct {
	fields := []*dub.LLField{}
	var implements *dub.LLStruct
	if decl.Implements != nil {
		t := gbuilder.TranslateType(decl.Implements.Resolve())
		var ok bool
		implements, ok = t.(*dub.LLStruct)
		if !ok {
			panic(t)
		}
	}
	for _, field := range decl.Fields {
		fields = append(fields, &dub.LLField{
			Name: field.Name,
			T:    gbuilder.TranslateType(field.Type.Resolve()),
		})
	}
	*s = dub.LLStruct{
		Name:       decl.Name,
		Implements: implements,
		Fields:     fields,
	}
	return s
}

func processDASM(name string) {
	decls := parseDASM(fmt.Sprintf("dasm/%s.dasm", name))
	glbls := semanticPass(decls)
	gbuilder := &GlobalDubBuilder{types: map[ASTType]dub.DubType{}}

	gbuilder.String = &dub.StringType{}
	gbuilder.types[glbls.String] = gbuilder.String

	gbuilder.Rune = &dub.RuneType{}
	gbuilder.types[glbls.Rune] = gbuilder.Rune

	gbuilder.Int = &dub.IntType{}
	gbuilder.types[glbls.Int] = gbuilder.Int

	gbuilder.Bool = &dub.BoolType{}
	gbuilder.types[glbls.Bool] = gbuilder.Bool

	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
		case *StructDecl:
			gbuilder.types[decl] = &dub.LLStruct{}
		default:
			panic(decl)
		}
	}

	structs := []*dub.LLStruct{}
	funcs := []*dub.LLFunc{}
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			f := lowerAST(decl, gbuilder)
			funcs = append(funcs, f)

			// Dump flowgraph
			dot := base.RegionToDot(f.Region)
			outfile := filepath.Join("output", name, fmt.Sprintf("%s.svg", f.Name))
			io.WriteDot(dot, outfile)
		case *StructDecl:
			t, _ := gbuilder.types[decl]
			s, _ := t.(*dub.LLStruct)
			structs = append(structs, lowerStruct(decl, s, gbuilder))
		default:
			panic(decl)
		}
	}
	for _, s := range structs {
		if s.Implements != nil {
			s.Implements.Abstract = true
		}
	}

	code := dub.GenerateGo(name, structs, funcs)
	fmt.Println(code)
	io.WriteFile(fmt.Sprintf("src/generated/%s/parser.go", name), []byte(code))
}

func main() {
	processDASM("math")
	processDASM("dubx")
}
