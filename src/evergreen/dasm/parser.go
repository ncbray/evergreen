package dasm

import (
	"evergreen/dub"
	"evergreen/dubx"
	"fmt"
	"io/ioutil"
	"strconv"
)

type DASMTokenType int

type DASMToken struct {
	Tok dubx.Token
	Pos int
}

type DASMScanner struct {
	state   *dub.DubState
	Current DASMToken
	Next    DASMToken
}

func (s *DASMScanner) Scan() {
	s.Current = s.Next

	// HACK
	s.Next.Pos = s.state.Index

	var next dubx.Token
	if s.state.Flow == 0 {
		next = dubx.Tokenize(s.state)
	}
	// HACK
	if s.state.Flow != 0 {
		if s.state.Index != len(s.state.Stream) {
			panic(s.state.Index)
		}
		s.Next.Tok = nil
		return
	}
	s.Next.Tok = next
	switch next := next.(type) {
	case *dubx.IdTok:
	case *dubx.PuncTok:
	case *dubx.RuneTok:
		v, _ := strconv.Unquote(next.Text)
		next.Value = []rune(v)[0]
	case *dubx.StrTok:
		v, _ := strconv.Unquote(next.Text)
		next.Value = v
	case *dubx.IntTok:
		v, _ := strconv.Atoi(next.Text)
		next.Value = v
	default:
		panic(next)
	}
}

func CreateScanner(data []byte) *DASMScanner {
	state := &dub.DubState{Stream: []rune(string(data))}
	s := &DASMScanner{state: state}
	s.Scan()
	s.Scan()
	return s
}

func getName(s *DASMScanner) (string, bool) {
	switch tok := s.Current.Tok.(type) {
	case *dubx.IdTok:
		text := tok.Text
		s.Scan()
		return text, true
	default:
		return "", false
	}
}

func getPunc(s *DASMScanner, text string) bool {
	switch tok := s.Current.Tok.(type) {
	case *dubx.PuncTok:
		if tok.Text == text {
			s.Scan()
			return true
		}
	}
	return false
}

func getKeyword(s *DASMScanner, text string) bool {
	switch tok := s.Current.Tok.(type) {
	case *dubx.IdTok:
		if tok.Text == text {
			s.Scan()
			return true
		}
	}
	return false
}

func getString(s *DASMScanner) (string, bool) {
	switch tok := s.Current.Tok.(type) {
	case *dubx.StrTok:
		value := tok.Value
		s.Scan()
		return value, true
	default:
		return "", false
	}
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
	"ge": ">=",
	"le": "<=",
}

func parseType(s *DASMScanner) (ASTTypeRef, bool) {
	switch tok := s.Current.Tok.(type) {
	case *dubx.IdTok:
		result := &TypeRef{Name: tok.Text}
		s.Scan()
		return result, true
	case *dubx.PuncTok:
		if tok.Text != "[" {
			return nil, false
		}
		// HACKish lookahead
		n, ok := s.Next.Tok.(*dubx.PuncTok)
		if !ok || n.Text != "]" {
			return nil, false
		}
		s.Scan()
		s.Scan()
		child, ok := parseType(s)
		if !ok {
			return nil, false
		}
		return &ListTypeRef{Type: child}, true
	default:
		return nil, false
	}
}

func parseExpr(s *DASMScanner) (ASTExpr, bool) {
	switch tok := s.Current.Tok.(type) {
	case *dubx.IdTok:
		switch tok.Text {
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
		case "eq", "ne", "gt", "lt", "ge", "le":
			op := nameToOp[tok.Text]
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
			text := tok.Text
			s.Scan()
			return &GetName{Name: text}, true
		}
	case *dubx.RuneTok:
		v := tok.Value
		s.Scan()
		return &RuneLiteral{Value: v}, true
	case *dubx.StrTok:
		v := tok.Value
		s.Scan()
		return &StringLiteral{Value: v}, true
	case *dubx.IntTok:
		v := tok.Value
		s.Scan()
		return &IntLiteral{Value: v}, true
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
		if getPunc(s, "}") {
			return result, true
		}

		expr, ok := parseExpr(s)
		if !ok {
			return nil, false
		}
		result = append(result, expr)
		for getPunc(s, ";") {
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
	switch tok := s.Current.Tok.(type) {
	case *dubx.StrTok:
		v := tok.Value
		s.Scan()
		return &DestructureString{Value: v}, true
	case *dubx.RuneTok:
		v := tok.Value
		s.Scan()
		return &DestructureRune{Value: v}, true
	case *dubx.IntTok:
		v := tok.Value
		s.Scan()
		return &DestructureInt{Value: v}, true
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
		if s.Current.Tok == nil {
			return &File{
				Decls: decls,
				Tests: tests,
			}, true
		}
		switch tok := s.Current.Tok.(type) {
		case *dubx.IdTok:
			switch tok.Text {
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
				panic(tok.Text)
			}
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
		fmt.Printf("%s: Unexpected %v @ %v\n", filename, s.Current.Tok, s.Current.Pos)
		panic(s.Current.Pos)
	}
	return f
}
