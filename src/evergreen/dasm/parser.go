package dasm

import (
	"evergreen/dub"
	"evergreen/dubx"
	"fmt"
	"io/ioutil"
	"strconv"
)

func getName(state *dub.DubState) string {
	tok := dubx.Ident(state)
	if state.Flow != 0 {
		return ""
	}
	return tok.Text
}

func getPunc(state *dub.DubState, value string) {
	for _, expected := range []rune(value) {
		c := state.Read()
		if state.Flow != 0 || c != expected {
			state.Fail()
			return
		}
	}
	dubx.S(state)
}

func getKeyword(state *dub.DubState, text string) {
	tok := dubx.Ident(state)
	if state.Flow != 0 || tok.Text != text {
		state.Fail()
	}
}

func getString(state *dub.DubState) string {
	tok := dubx.StrT(state)
	if state.Flow != 0 {
		return ""
	}
	v, _ := strconv.Unquote(tok.Text)
	return v
}

func parseExprList(state *dub.DubState) []ASTExpr {
	getPunc(state, "(")
	if state.Flow != 0 {
		return nil
	}
	exprs := []ASTExpr{}
	for {
		checkpoint := state.Checkpoint()
		e := parseExpr(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		exprs = append(exprs, e)
	}
	getPunc(state, ")")
	if state.Flow != 0 {
		return nil
	}
	return exprs
}

func parseKeyValueList(state *dub.DubState) []*KeyValue {
	getPunc(state, "(")
	if state.Flow != 0 {
		return nil
	}
	args := []*KeyValue{}
	for {
		checkpoint := state.Checkpoint()
		key := getName(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		getPunc(state, ":")
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		value := parseExpr(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		args = append(args, &KeyValue{Key: key, Value: value})
	}
	getPunc(state, ")")
	if state.Flow != 0 {
		return nil
	}
	return args
}

func parseTypeList(state *dub.DubState) []ASTTypeRef {
	getPunc(state, "(")
	if state.Flow != 0 {
		return nil
	}
	types := []ASTTypeRef{}
	for {
		checkpoint := state.Checkpoint()
		t := parseType(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		types = append(types, t)
	}
	getPunc(state, ")")
	if state.Flow != 0 {
		return nil
	}
	return types
}

var nameToOp = map[string]string{
	"eq": "==",
	"ne": "!=",
	"gt": ">",
	"lt": "<",
	"ge": ">=",
	"le": "<=",
}

func parseType(state *dub.DubState) ASTTypeRef {
	checkpoint := state.Checkpoint()
	tok := dubx.Ident(state)
	if state.Flow == 0 {
		return &TypeRef{Name: tok.Text}
	}
	state.Recover(checkpoint)

	getPunc(state, "[]")
	if state.Flow != 0 {
		return nil
	}
	t := parseType(state)
	if state.Flow != 0 {
		return nil
	}
	return &ListTypeRef{Type: t}
}

func parseExpr(state *dub.DubState) ASTExpr {
	checkpoint := state.Checkpoint()
	tok := dubx.Ident(state)
	if state.Flow == 0 {
		switch tok.Text {
		case "star":
			block := parseCodeBlock(state)
			if state.Flow == 0 {
				return &Repeat{Block: block, Min: 0}
			}
		case "plus":
			block := parseCodeBlock(state)
			if state.Flow == 0 {
				return &Repeat{Block: block, Min: 1}
			}
		case "question":
			block := parseCodeBlock(state)
			if state.Flow == 0 {
				return &Optional{Block: block}
			}
		case "slice":
			block := parseCodeBlock(state)
			if state.Flow == 0 {
				return &Slice{Block: block}
			}
		case "if":
			expr := parseExpr(state)
			if state.Flow == 0 {
				block := parseCodeBlock(state)
				if state.Flow == 0 {
					return &If{Expr: expr, Block: block}
				}
			}
		case "var":
			name := getName(state)
			if state.Flow == 0 {
				t := parseType(state)
				if state.Flow == 0 {
					checkpoint := state.Checkpoint()
					getPunc(state, "=")
					if state.Flow == 0 {
						expr := parseExpr(state)
						if state.Flow == 0 {
							return &Assign{Expr: expr, Name: name, Type: t, Define: true}
						}
					}
					state.Recover(checkpoint)
					return &Assign{Name: name, Type: t, Define: true}
				}
			}
		case "define":
			name := getName(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Assign{Expr: expr, Name: name, Define: true}
				}
			}
		case "assign":
			name := getName(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Assign{Expr: expr, Name: name, Define: false}
				}
			}
		case "read":
			return &Read{}
		case "fail":
			return &Fail{}
		case "eq", "ne", "gt", "lt", "ge", "le":
			op := nameToOp[tok.Text]
			l := parseExpr(state)
			if state.Flow == 0 {
				r := parseExpr(state)
				if state.Flow == 0 {
					return &BinaryOp{Left: l, Op: op, Right: r}
				}
			}
		case "call":
			name := getName(state)
			if state.Flow == 0 {
				return &Call{Name: name}
			}
		case "cons":
			t := parseType(state)
			if state.Flow == 0 {
				args := parseKeyValueList(state)
				if state.Flow == 0 {
					return &Construct{Type: t, Args: args}
				}
			}
		case "conl":
			t := parseType(state)
			if state.Flow == 0 {
				args := parseExprList(state)
				if state.Flow == 0 {
					return &ConstructList{Type: t, Args: args}
				}
			}
		case "append":
			name := getName(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Assign{
						Expr: &Append{
							List: &GetName{
								Name: name,
							},
							Value: expr,
						},
						Name: name,
					}
				}
			}
		case "return":
			exprs := parseExprList(state)
			if state.Flow == 0 {
				return &Return{Exprs: exprs}
			}
		default:
			text := tok.Text
			return &GetName{Name: text}
		}
	}

	state.Recover(checkpoint)
	{
		tok := dubx.Rune(state)
		if state.Flow == 0 {
			v, _ := strconv.Unquote(tok.Text)
			return &RuneLiteral{Value: []rune(v)[0]}
		}
	}
	state.Recover(checkpoint)
	{
		tok := dubx.StrT(state)
		if state.Flow == 0 {
			v, _ := strconv.Unquote(tok.Text)
			return &StringLiteral{Value: v}
		}
	}
	state.Recover(checkpoint)
	{
		tok := dubx.Int(state)
		if state.Flow == 0 {
			v, _ := strconv.Atoi(tok.Text)
			return &IntLiteral{Value: v}
		}
	}
	// Fail through
	return nil
}

func parseCodeBlock(state *dub.DubState) []ASTExpr {
	getPunc(state, "{")
	if state.Flow != 0 {
		return nil
	}
	result := []ASTExpr{}
	for {
		checkpoint := state.Checkpoint()
		expr := parseExpr(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		result = append(result, expr)
		for {
			checkpoint := state.Checkpoint()
			getPunc(state, ";")
			if state.Flow != 0 {
				state.Recover(checkpoint)
				break
			}

		}
	}
	getPunc(state, "}")
	if state.Flow != 0 {
		return nil
	}
	return result
}

func parseFunction(state *dub.DubState) *FuncDecl {
	name := getName(state)
	if state.Flow != 0 {
		return nil
	}
	returnTypes := parseTypeList(state)
	if state.Flow != 0 {
		return nil
	}
	block := parseCodeBlock(state)
	if state.Flow != 0 {
		return nil
	}
	return &FuncDecl{Name: name, ReturnTypes: returnTypes, Block: block}
}

func parseImplements(state *dub.DubState) ASTTypeRef {
	checkpoint := state.Checkpoint()
	getKeyword(state, "implements")
	if state.Flow != 0 {
		state.Recover(checkpoint)
		return nil
	}
	return parseType(state)
}

func parseStructure(state *dub.DubState) *StructDecl {
	name := getName(state)
	if state.Flow != 0 {
		return nil
	}

	implements := parseImplements(state)
	if state.Flow != 0 {
		return nil
	}

	getPunc(state, "{")
	if state.Flow != 0 {
		return nil
	}

	fields := []*FieldDecl{}
	for {
		checkpoint := state.Checkpoint()
		name := getName(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		t := parseType(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		fields = append(fields, &FieldDecl{Name: name, Type: t})
	}

	getPunc(state, "}")
	if state.Flow != 0 {
		return nil
	}

	return &StructDecl{
		Name:       name,
		Implements: implements,
		Fields:     fields,
	}
}

func parseLiteralDestructure(state *dub.DubState) Destructure {
	checkpoint := state.Checkpoint()

	{
		tok := dubx.Rune(state)
		if state.Flow == 0 {
			v, _ := strconv.Unquote(tok.Text)
			return &DestructureRune{Value: []rune(v)[0]}
		}
	}
	state.Recover(checkpoint)
	{
		tok := dubx.StrT(state)
		if state.Flow == 0 {
			v, _ := strconv.Unquote(tok.Text)
			return &DestructureString{Value: v}
		}
	}
	state.Recover(checkpoint)
	{
		tok := dubx.Int(state)
		if state.Flow == 0 {
			v, _ := strconv.Atoi(tok.Text)
			return &DestructureInt{Value: v}
		}
	}
	return nil
}

func parseDestructure(state *dub.DubState) Destructure {
	checkpoint := state.Checkpoint()
	t := parseType(state)
	if state.Flow != 0 {
		state.Recover(checkpoint)
		return parseLiteralDestructure(state)
	}
	getPunc(state, "{")
	if state.Flow != 0 {
		return nil
	}
	switch t := t.(type) {
	case *ListTypeRef:
		args := []Destructure{}
		for {
			checkpoint := state.Checkpoint()
			arg := parseDestructure(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				break
			}
			args = append(args, arg)

			/*
				checkpoint = state.Checkpoint()
				getPunc(state, ",")
				if state.Flow != 0 {
					state.Recover(checkpoint)
					break
				}
			*/
		}
		getPunc(state, "}")
		if state.Flow != 0 {
			return nil
		}
		return &DestructureList{
			Type: t,
			Args: args,
		}
	case *TypeRef:
		args := []*DestructureField{}
		for {
			checkpoint := state.Checkpoint()
			name := getName(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				break
			}
			getPunc(state, ":")
			if state.Flow != 0 {
				state.Recover(checkpoint)
				break
			}

			arg := parseDestructure(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				break
			}
			args = append(args, &DestructureField{Name: name, Destructure: arg})

			/*
				checkpoint = state.Checkpoint()
				getPunc(state, ",")
				if state.Flow != 0 {
					state.Recover(checkpoint)
					break
				}
			*/
		}
		getPunc(state, "}")
		if state.Flow != 0 {
			return nil
		}
		return &DestructureStruct{
			Type: t,
			Args: args,
		}
	default:
		panic(t)
	}
}

func parseTest(state *dub.DubState) *Test {
	rule := getName(state)
	if state.Flow != 0 {
		return nil
	}
	name := getName(state)
	if state.Flow != 0 {
		return nil
	}
	input := getString(state)
	if state.Flow != 0 {
		return nil
	}
	destructure := parseDestructure(state)
	if state.Flow != 0 {
		return nil
	}
	return &Test{Name: name, Rule: rule, Input: input, Destructure: destructure}
}

func parseFile(state *dub.DubState) *File {
	decls := []Decl{}
	tests := []*Test{}
	for {
		checkpoint := state.Checkpoint()
		tok := dubx.Ident(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			goto done
		}
		switch tok.Text {
		case "func":
			f := parseFunction(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			decls = append(decls, f)
		case "struct":
			f := parseStructure(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			decls = append(decls, f)
		case "test":
			t := parseTest(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			tests = append(tests, t)
		default:
			panic("foo")
			state.Recover(checkpoint)
			goto done
		}
	}
done:
	if state.Index != len(state.Stream) {
		state.Fail()
		return nil
	}
	return &File{
		Decls: decls,
		Tests: tests,
	}
}

func FindLines(stream []rune) []int {
	lines := []int{}
	for i, r := range stream {
		if r == '\n' {
			lines = append(lines, i+1)
		}
	}
	lines = append(lines, len(stream))
	return lines
}

func PrintError(filename string, deepest int, stream []rune, lines []int) {
	// Stupid linear search
	var line int
	var col int
	var start int
	var end int
	for i, s := range lines {
		end = s
		if deepest < s {
			break
		}
		start = s
		line = i
		col = deepest - start
	}
	text := string(stream[start:end])
	fmt.Printf("%s: Unexpected @ %d:%d\n%s", filename, line, col, text)
}

func ParseDASM(filename string) *File {
	data, _ := ioutil.ReadFile(filename)
	stream := []rune(string(data))
	state := &dub.DubState{Stream: stream}
	f := parseFile(state)
	if state.Flow != 0 {
		lines := FindLines(stream)
		PrintError(filename, state.Deepest, stream, lines)
		panic(state.Deepest)
	}
	return f
}
