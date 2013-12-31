package dasm

import (
	"evergreen/dub"
	"evergreen/dubx"
	"fmt"
	"io/ioutil"
)

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

func getKeyword(state *dub.DubState, expected string) {
	for _, expected := range []rune(expected) {
		actual := state.Peek()
		if state.Flow != 0 {
			return
		}
		if actual != expected {
			state.Fail()
			return
		}
		state.Consume()
	}
	dubx.S(state)
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

		checkpoint = state.Checkpoint()
		getPunc(state, ",")
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
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
		key := dubx.Ident(state)
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

		checkpoint = state.Checkpoint()
		getPunc(state, ",")
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
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
		t := dubx.ParseTypeRef(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		types = append(types, t)

		checkpoint = state.Checkpoint()
		getPunc(state, ",")
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
	}
	getPunc(state, ")")
	if state.Flow != 0 {
		return nil
	}
	return types
}

func parseExpr(state *dub.DubState) ASTExpr {
	checkpoint := state.Checkpoint()
	e := dubx.ParseExpr(state)
	if state.Flow == 0 {
		return e
	}
	state.Recover(checkpoint)
	text := dubx.Ident(state)
	if state.Flow == 0 {
		switch text {
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
		case "choose":
			blocks := [][]ASTExpr{parseCodeBlock(state)}
			if state.Flow == 0 {
				for {
					checkpoint := state.Checkpoint()
					getKeyword(state, "or")
					if state.Flow != 0 {
						state.Recover(checkpoint)
						break
					}
					block := parseCodeBlock(state)
					if state.Flow != 0 {
						state.Recover(checkpoint)
						break
					}
					blocks = append(blocks, block)
				}
				return &Choice{Blocks: blocks}
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
			name := dubx.Ident(state)
			if state.Flow == 0 {
				t := dubx.ParseTypeRef(state)
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
			name := dubx.Ident(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Assign{Expr: expr, Name: name, Define: true}
				}
			}
		case "assign":
			name := dubx.Ident(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Assign{Expr: expr, Name: name, Define: false}
				}
			}
		case "cons":
			t := dubx.ParseTypeRef(state)
			if state.Flow == 0 {
				args := parseKeyValueList(state)
				if state.Flow == 0 {
					return &Construct{Type: t, Args: args}
				}
			}
		case "conl":
			t := dubx.ParseTypeRef(state)
			if state.Flow == 0 {
				args := parseExprList(state)
				if state.Flow == 0 {
					return &ConstructList{Type: t, Args: args}
				}
			}
		case "coerce":
			t := dubx.ParseTypeRef(state)
			if state.Flow == 0 {
				expr := parseExpr(state)
				if state.Flow == 0 {
					return &Coerce{Type: t, Expr: expr}
				}
			}

		case "append":
			name := dubx.Ident(state)
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
			checkpoint := state.Checkpoint()
			exprs := parseExprList(state)
			if state.Flow == 0 {
				return &Return{Exprs: exprs}
			}
			state.Recover(checkpoint)
			expr := parseExpr(state)
			if state.Flow == 0 {
				return &Return{Exprs: []ASTExpr{expr}}
			}
			state.Recover(checkpoint)
			return &Return{Exprs: []ASTExpr{}}
		default:
			return &GetName{Name: text}
		}
	}
	state.Recover(checkpoint)
	{
		op := dubx.BinaryOperator(state)
		if state.Flow == 0 {
			l := parseExpr(state)
			if state.Flow == 0 {
				r := parseExpr(state)
				if state.Flow == 0 {
					return &BinaryOp{Left: l, Op: op, Right: r}
				}
			}
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
	name := dubx.Ident(state)
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
	return dubx.ParseTypeRef(state)
}

func parseStructure(state *dub.DubState) *StructDecl {
	name := dubx.Ident(state)
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
		name := dubx.Ident(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		t := dubx.ParseTypeRef(state)
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

func parseTest(state *dub.DubState) *Test {
	rule := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}
	name := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}
	input := dubx.DecodeString(state)
	if state.Flow != 0 {
		return nil
	}
	dubx.S(state)
	destructure := dubx.ParseDestructure(state)
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
		text := dubx.Ident(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			goto done
		}
		switch text {
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
		start = end
		end = s
		line = i + 1
		if deepest < end {
			break
		}
	}
	col = deepest - start
	text := string(stream[start:end])
	fmt.Printf("%s: Unexpected @ %d:%d\n%s", filename, line, col, text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < col; i++ {
		arrow = append(arrow, ' ')
	}
	arrow = append(arrow, '^')
	fmt.Println(string(arrow))
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
