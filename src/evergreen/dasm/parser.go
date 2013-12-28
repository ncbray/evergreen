package dasm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"text/scanner"
)

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

func getString(s *DASMScanner) (string, bool) {
	if s.Current.Type == String {
		value, _ := strconv.Unquote(s.Current.Text)
		s.Scan()
		return value, true
	}
	return "", false
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
		return nil, false
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

func parseLiteralDestructure(s *DASMScanner) (Destructure, bool) {
	switch s.Current.Type {
	case String:
		s, _ := getString(s)
		return &DestructureString{Value: s}, true
	default:
		return nil, false
	}
}

func parseDestructure(s *DASMScanner) (Destructure, bool) {
	t, ok := parseType(s)
	if !ok {
		return parseLiteralDestructure(s)
	}
	ok = getPunc(s, "{")
	if !ok {
		return nil, false
	}
	switch t := t.(type) {
	case *ListTypeRef:
		args := []Destructure{}
		for {
			if getPunc(s, "}") {
				return &DestructureList{
					Type: t,
					Args: args,
				}, true
			}

			arg, ok := parseDestructure(s)
			if !ok {
				return nil, false
			}
			getPunc(s, ",")
			args = append(args, arg)
		}
	case *TypeRef:
		args := []*DestructureField{}
		for {
			if getPunc(s, "}") {
				return &DestructureStruct{
					Type: t,
					Args: args,
				}, true
			}
			name, ok := getName(s)
			if !ok {
				return nil, false
			}
			ok = getPunc(s, ":")
			if !ok {
				return nil, false
			}
			arg, ok := parseDestructure(s)
			if !ok {
				return nil, false
			}
			getPunc(s, ",")
			args = append(args, &DestructureField{Name: name, Destructure: arg})
		}
	default:
		panic(t)
	}
}

func parseTest(s *DASMScanner) (*Test, bool) {
	rule, ok := getName(s)
	if !ok {
		return nil, false
	}
	name, ok := getName(s)
	if !ok {
		return nil, false
	}
	input, ok := getString(s)
	if !ok {
		return nil, false
	}
	destructure, ok := parseDestructure(s)
	if !ok {
		return nil, false
	}
	return &Test{Name: name, Rule: rule, Input: input, Destructure: destructure}, true
}

func parseFile(s *DASMScanner) (*File, bool) {
	decls := []Decl{}
	tests := []*Test{}
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
			case "test":
				s.Scan()
				t, ok := parseTest(s)
				if !ok {
					return nil, false
				}
				tests = append(tests, t)
			default:
				panic(s.Current.Text)
			}
		case EOF:
			return &File{
				Decls: decls,
				Tests: tests,
			}, true
		default:
			return nil, false
		}
	}
}

func ParseDASM(filename string) *File {
	data, _ := ioutil.ReadFile(filename)
	s := CreateScanner(data)
	f, ok := parseFile(s)
	if !ok {
		fmt.Printf("Unexpected %s @ %s\n", s.Current.Text, s.Current.Pos)
		panic(s.Current.Pos)
	}
	return f
}
