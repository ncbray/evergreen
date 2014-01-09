package tree

import (
	"evergreen/dub/runtime"
)

type TextMatch interface {
	isTextMatch()
}
type RuneFilter struct {
	Min rune
	Max rune
}
type RuneRangeMatch struct {
	Invert  bool
	Filters []*RuneFilter
}

func (node *RuneRangeMatch) isTextMatch() {
}

type StringLiteralMatch struct {
	Value string
}

func (node *StringLiteralMatch) isTextMatch() {
}

type MatchSequence struct {
	Matches []TextMatch
}

func (node *MatchSequence) isTextMatch() {
}

type MatchChoice struct {
	Matches []TextMatch
}

func (node *MatchChoice) isTextMatch() {
}

type MatchRepeat struct {
	Match TextMatch
	Min   int
}

func (node *MatchRepeat) isTextMatch() {
}

type MatchLookahead struct {
	Invert bool
	Match  TextMatch
}

func (node *MatchLookahead) isTextMatch() {
}

type Id struct {
	Pos  int
	Text string
}
type ASTExpr interface {
	isASTExpr()
}
type RuneLiteral struct {
	Text  string
	Value rune
}

func (node *RuneLiteral) isASTExpr() {
}

type StringLiteral struct {
	Text  string
	Value string
}

func (node *StringLiteral) isASTExpr() {
}

type IntLiteral struct {
	Text  string
	Value int
}

func (node *IntLiteral) isASTExpr() {
}

type BoolLiteral struct {
	Text  string
	Value bool
}

func (node *BoolLiteral) isASTExpr() {
}

type StringMatch struct {
	Match TextMatch
}

func (node *StringMatch) isASTExpr() {
}

type RuneMatch struct {
	Match *RuneRangeMatch
}

func (node *RuneMatch) isASTExpr() {
}

type ASTDecl interface {
	isASTDecl()
}
type ASTType interface {
	isASTType()
}
type ASTTypeRef interface {
	isASTTypeRef()
}
type TypeRef struct {
	Name *Id
	T    ASTType
}

func (node *TypeRef) isASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    ASTType
}

func (node *ListTypeRef) isASTTypeRef() {
}

type Destructure interface {
	isDestructure()
}
type DestructureValue struct {
	Expr ASTExpr
}

func (node *DestructureValue) isDestructure() {
}

type DestructureField struct {
	Name        *Id
	Destructure Destructure
}
type DestructureStruct struct {
	Type *TypeRef
	Args []*DestructureField
}

func (node *DestructureStruct) isDestructure() {
}

type DestructureList struct {
	Type *ListTypeRef
	Args []Destructure
}

func (node *DestructureList) isDestructure() {
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

type Choice struct {
	Blocks [][]ASTExpr
}

func (node *Choice) isASTExpr() {
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

type Assign struct {
	Expr   ASTExpr
	Name   *Id
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) isASTExpr() {
}

type GetName struct {
	Name *Id
	Info int
}

func (node *GetName) isASTExpr() {
}

type NamedExpr struct {
	Name *Id
	Expr ASTExpr
}
type Construct struct {
	Type *TypeRef
	Args []*NamedExpr
}

func (node *Construct) isASTExpr() {
}

type ConstructList struct {
	Type *ListTypeRef
	Args []ASTExpr
}

func (node *ConstructList) isASTExpr() {
}

type Coerce struct {
	Type ASTTypeRef
	Expr ASTExpr
}

func (node *Coerce) isASTExpr() {
}

type Call struct {
	Name *Id
	T    ASTType
}

func (node *Call) isASTExpr() {
}

type Position struct {
}

func (node *Position) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
}

type Append struct {
	List ASTExpr
	Expr ASTExpr
	T    ASTType
}

func (node *Append) isASTExpr() {
}

type Return struct {
	Exprs []ASTExpr
}

func (node *Return) isASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     ASTType
}

func (node *BinaryOp) isASTExpr() {
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) isASTDecl() {
}
func (node *BuiltinType) isASTType() {
}

type ListType struct {
	Type ASTType
}

func (node *ListType) isASTDecl() {
}
func (node *ListType) isASTType() {
}

type FieldDecl struct {
	Name *Id
	Type ASTTypeRef
}
type StructDecl struct {
	Name       *Id
	Implements ASTTypeRef
	Fields     []*FieldDecl
}

func (node *StructDecl) isASTDecl() {
}
func (node *StructDecl) isASTType() {
}

type ASTFunc interface {
	isASTFunc()
}
type LocalInfo struct {
	Name string
	T    ASTType
}
type FuncDecl struct {
	Name        *Id
	ReturnTypes []ASTTypeRef
	Block       []ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) isASTDecl() {
}
func (node *FuncDecl) isASTFunc() {
}

type Test struct {
	Rule        *Id
	Name        *Id
	Type        ASTType
	Input       string
	Destructure Destructure
}
type File struct {
	Decls []ASTDecl
	Tests []*Test
}

func S(frame *runtime.State) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 bool
	goto block1
block1:
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = ' '
		r3 = r1 == r2
		if r3 {
			goto block2
		} else {
			r4 = '\t'
			r5 = r1 == r4
			if r5 {
				goto block2
			} else {
				r6 = '\r'
				r7 = r1 == r6
				if r7 {
					goto block2
				} else {
					r8 = '\n'
					r9 = r1 == r8
					if r9 {
						goto block2
					} else {
						frame.Fail()
						goto block3
					}
				}
			}
		}
	} else {
		goto block3
	}
block2:
	frame.Consume()
	goto block1
block3:
	frame.Recover(r0)
	return
}
func EndKeyword(frame *runtime.State) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 bool
	var r14 rune
	var r15 bool
	r0 = frame.LookaheadBegin()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = 'a'
		r3 = r1 >= r2
		if r3 {
			r4 = 'z'
			r5 = r1 <= r4
			if r5 {
				goto block3
			} else {
				goto block1
			}
		} else {
			goto block1
		}
	} else {
		goto block5
	}
block1:
	r6 = 'A'
	r7 = r1 >= r6
	if r7 {
		r8 = 'Z'
		r9 = r1 <= r8
		if r9 {
			goto block3
		} else {
			goto block2
		}
	} else {
		goto block2
	}
block2:
	r10 = '_'
	r11 = r1 == r10
	if r11 {
		goto block3
	} else {
		r12 = '0'
		r13 = r1 >= r12
		if r13 {
			r14 = '9'
			r15 = r1 <= r14
			if r15 {
				goto block3
			} else {
				goto block4
			}
		} else {
			goto block4
		}
	}
block3:
	frame.Consume()
	frame.LookaheadFail(r0)
	goto block6
block4:
	frame.Fail()
	goto block5
block5:
	frame.LookaheadNormal(r0)
	S(frame)
	if frame.Flow == 0 {
		return
	} else {
		goto block6
	}
block6:
	return
}
func Ident(frame *runtime.State) (ret0 *Id) {
	var r0 int
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 bool
	var r13 int
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 rune
	var r18 bool
	var r19 rune
	var r20 bool
	var r21 rune
	var r22 bool
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 bool
	var r27 rune
	var r28 bool
	var r29 string
	var r30 *Id
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		r3 = 'a'
		r4 = r2 >= r3
		if r4 {
			r5 = 'z'
			r6 = r2 <= r5
			if r6 {
				goto block3
			} else {
				goto block1
			}
		} else {
			goto block1
		}
	} else {
		goto block10
	}
block1:
	r7 = 'A'
	r8 = r2 >= r7
	if r8 {
		r9 = 'Z'
		r10 = r2 <= r9
		if r10 {
			goto block3
		} else {
			goto block2
		}
	} else {
		goto block2
	}
block2:
	r11 = '_'
	r12 = r2 == r11
	if r12 {
		goto block3
	} else {
		frame.Fail()
		goto block10
	}
block3:
	frame.Consume()
	goto block4
block4:
	r13 = frame.Checkpoint()
	r14 = frame.Peek()
	if frame.Flow == 0 {
		r15 = 'a'
		r16 = r14 >= r15
		if r16 {
			r17 = 'z'
			r18 = r14 <= r17
			if r18 {
				goto block7
			} else {
				goto block5
			}
		} else {
			goto block5
		}
	} else {
		goto block9
	}
block5:
	r19 = 'A'
	r20 = r14 >= r19
	if r20 {
		r21 = 'Z'
		r22 = r14 <= r21
		if r22 {
			goto block7
		} else {
			goto block6
		}
	} else {
		goto block6
	}
block6:
	r23 = '_'
	r24 = r14 == r23
	if r24 {
		goto block7
	} else {
		r25 = '0'
		r26 = r14 >= r25
		if r26 {
			r27 = '9'
			r28 = r14 <= r27
			if r28 {
				goto block7
			} else {
				goto block8
			}
		} else {
			goto block8
		}
	}
block7:
	frame.Consume()
	goto block4
block8:
	frame.Fail()
	goto block9
block9:
	frame.Recover(r13)
	r29 = frame.Slice(r1)
	S(frame)
	if frame.Flow == 0 {
		r30 = &Id{Pos: r0, Text: r29}
		ret0 = r30
		return
	} else {
		goto block10
	}
block10:
	return
}
func DecodeInt(frame *runtime.State) (ret0 int) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 bool
	var r6 int
	var r7 rune
	var r8 int
	var r9 int
	var r10 int
	var r11 int
	var r12 int
	var r13 int
	var r14 int
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 rune
	var r19 bool
	var r20 int
	var r21 rune
	var r22 int
	var r23 int
	var r24 int
	var r25 int
	var r26 int
	r0 = 0
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = '0'
		r3 = r1 >= r2
		if r3 {
			r4 = '9'
			r5 = r1 <= r4
			if r5 {
				frame.Consume()
				r6 = int(r1)
				r7 = '0'
				r8 = int(r7)
				r9 = r6 - r8
				r10 = 10
				r11 = r0 * r10
				r12 = r11 + r9
				r13 = r12
				goto block1
			} else {
				goto block4
			}
		} else {
			goto block4
		}
	} else {
		goto block5
	}
block1:
	r14 = frame.Checkpoint()
	r15 = frame.Peek()
	if frame.Flow == 0 {
		r16 = '0'
		r17 = r15 >= r16
		if r17 {
			r18 = '9'
			r19 = r15 <= r18
			if r19 {
				frame.Consume()
				r20 = int(r15)
				r21 = '0'
				r22 = int(r21)
				r23 = r20 - r22
				r24 = 10
				r25 = r13 * r24
				r26 = r25 + r23
				r13 = r26
				goto block1
			} else {
				goto block2
			}
		} else {
			goto block2
		}
	} else {
		goto block3
	}
block2:
	frame.Fail()
	goto block3
block3:
	frame.Recover(r14)
	ret0 = r13
	return
block4:
	frame.Fail()
	goto block5
block5:
	return
}
func EscapedChar(frame *runtime.State) (ret0 rune) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 rune
	var r15 bool
	var r16 rune
	var r17 rune
	var r18 rune
	var r19 bool
	var r20 rune
	var r21 rune
	var r22 rune
	var r23 bool
	var r24 rune
	var r25 rune
	var r26 rune
	var r27 bool
	var r28 rune
	var r29 rune
	var r30 rune
	var r31 bool
	var r32 rune
	var r33 rune
	var r34 rune
	var r35 bool
	var r36 rune
	var r37 rune
	var r38 rune
	var r39 bool
	var r40 rune
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = 'a'
		r3 = r1 == r2
		if r3 {
			frame.Consume()
			r4 = '\a'
			ret0 = r4
			goto block10
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	frame.Recover(r0)
	r5 = frame.Peek()
	if frame.Flow == 0 {
		r6 = 'b'
		r7 = r5 == r6
		if r7 {
			frame.Consume()
			r8 = '\b'
			ret0 = r8
			goto block10
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r0)
	r9 = frame.Peek()
	if frame.Flow == 0 {
		r10 = 'f'
		r11 = r9 == r10
		if r11 {
			frame.Consume()
			r12 = '\f'
			ret0 = r12
			goto block10
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r13 = frame.Peek()
	if frame.Flow == 0 {
		r14 = 'n'
		r15 = r13 == r14
		if r15 {
			frame.Consume()
			r16 = '\n'
			ret0 = r16
			goto block10
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r0)
	r17 = frame.Peek()
	if frame.Flow == 0 {
		r18 = 'r'
		r19 = r17 == r18
		if r19 {
			frame.Consume()
			r20 = '\r'
			ret0 = r20
			goto block10
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block5:
	frame.Recover(r0)
	r21 = frame.Peek()
	if frame.Flow == 0 {
		r22 = 't'
		r23 = r21 == r22
		if r23 {
			frame.Consume()
			r24 = '\t'
			ret0 = r24
			goto block10
		} else {
			frame.Fail()
			goto block6
		}
	} else {
		goto block6
	}
block6:
	frame.Recover(r0)
	r25 = frame.Peek()
	if frame.Flow == 0 {
		r26 = 'v'
		r27 = r25 == r26
		if r27 {
			frame.Consume()
			r28 = '\v'
			ret0 = r28
			goto block10
		} else {
			frame.Fail()
			goto block7
		}
	} else {
		goto block7
	}
block7:
	frame.Recover(r0)
	r29 = frame.Peek()
	if frame.Flow == 0 {
		r30 = '\\'
		r31 = r29 == r30
		if r31 {
			frame.Consume()
			r32 = '\\'
			ret0 = r32
			goto block10
		} else {
			frame.Fail()
			goto block8
		}
	} else {
		goto block8
	}
block8:
	frame.Recover(r0)
	r33 = frame.Peek()
	if frame.Flow == 0 {
		r34 = '\''
		r35 = r33 == r34
		if r35 {
			frame.Consume()
			r36 = '\''
			ret0 = r36
			goto block10
		} else {
			frame.Fail()
			goto block9
		}
	} else {
		goto block9
	}
block9:
	frame.Recover(r0)
	r37 = frame.Peek()
	if frame.Flow == 0 {
		r38 = '"'
		r39 = r37 == r38
		if r39 {
			frame.Consume()
			r40 = '"'
			ret0 = r40
			goto block10
		} else {
			frame.Fail()
			goto block11
		}
	} else {
		goto block11
	}
block10:
	return
block11:
	return
}
func DecodeString(frame *runtime.State) (ret0 string) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 []rune
	var r4 []rune
	var r5 int
	var r6 int
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 bool
	var r12 []rune
	var r13 rune
	var r14 rune
	var r15 bool
	var r16 rune
	var r17 []rune
	var r18 rune
	var r19 rune
	var r20 bool
	var r21 string
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '"'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = []rune{}
			r4 = r3
			goto block1
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block1:
	r5 = frame.Checkpoint()
	r6 = frame.Checkpoint()
	r7 = frame.Peek()
	if frame.Flow == 0 {
		r8 = '"'
		r9 = r7 == r8
		if r9 {
			goto block2
		} else {
			r10 = '\\'
			r11 = r7 == r10
			if r11 {
				goto block2
			} else {
				frame.Consume()
				r12 = append(r4, r7)
				r4 = r12
				goto block1
			}
		}
	} else {
		goto block3
	}
block2:
	frame.Fail()
	goto block3
block3:
	frame.Recover(r6)
	r13 = frame.Peek()
	if frame.Flow == 0 {
		r14 = '\\'
		r15 = r13 == r14
		if r15 {
			frame.Consume()
			r16 = EscapedChar(frame)
			if frame.Flow == 0 {
				r17 = append(r4, r16)
				r4 = r17
				goto block1
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r5)
	r18 = frame.Peek()
	if frame.Flow == 0 {
		r19 = '"'
		r20 = r18 == r19
		if r20 {
			frame.Consume()
			r21 = string(r4)
			ret0 = r21
			return
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block5:
	return
}
func DecodeRune(frame *runtime.State) (ret0 rune) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 int
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 rune
	var r14 rune
	var r15 rune
	var r16 bool
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '\''
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = frame.Checkpoint()
			r4 = frame.Peek()
			if frame.Flow == 0 {
				r5 = '\\'
				r6 = r4 == r5
				if r6 {
					goto block1
				} else {
					r7 = '\''
					r8 = r4 == r7
					if r8 {
						goto block1
					} else {
						frame.Consume()
						r9 = r4
						goto block3
					}
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(r3)
	r10 = frame.Peek()
	if frame.Flow == 0 {
		r11 = '\\'
		r12 = r10 == r11
		if r12 {
			frame.Consume()
			r13 = EscapedChar(frame)
			if frame.Flow == 0 {
				r9 = r13
				goto block3
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block3:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		r15 = '\''
		r16 = r14 == r15
		if r16 {
			frame.Consume()
			ret0 = r9
			return
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	return
}
func DecodeBool(frame *runtime.State) (ret0 bool) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 bool
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 rune
	var r18 rune
	var r19 bool
	var r20 rune
	var r21 rune
	var r22 bool
	var r23 rune
	var r24 rune
	var r25 bool
	var r26 rune
	var r27 rune
	var r28 bool
	var r29 bool
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = 't'
		r3 = r1 == r2
		if r3 {
			frame.Consume()
			r4 = frame.Peek()
			if frame.Flow == 0 {
				r5 = 'r'
				r6 = r4 == r5
				if r6 {
					frame.Consume()
					r7 = frame.Peek()
					if frame.Flow == 0 {
						r8 = 'u'
						r9 = r7 == r8
						if r9 {
							frame.Consume()
							r10 = frame.Peek()
							if frame.Flow == 0 {
								r11 = 'e'
								r12 = r10 == r11
								if r12 {
									frame.Consume()
									r13 = true
									ret0 = r13
									goto block2
								} else {
									frame.Fail()
									goto block1
								}
							} else {
								goto block1
							}
						} else {
							frame.Fail()
							goto block1
						}
					} else {
						goto block1
					}
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	frame.Recover(r0)
	r14 = frame.Peek()
	if frame.Flow == 0 {
		r15 = 'f'
		r16 = r14 == r15
		if r16 {
			frame.Consume()
			r17 = frame.Peek()
			if frame.Flow == 0 {
				r18 = 'a'
				r19 = r17 == r18
				if r19 {
					frame.Consume()
					r20 = frame.Peek()
					if frame.Flow == 0 {
						r21 = 'l'
						r22 = r20 == r21
						if r22 {
							frame.Consume()
							r23 = frame.Peek()
							if frame.Flow == 0 {
								r24 = 's'
								r25 = r23 == r24
								if r25 {
									frame.Consume()
									r26 = frame.Peek()
									if frame.Flow == 0 {
										r27 = 'e'
										r28 = r26 == r27
										if r28 {
											frame.Consume()
											r29 = false
											ret0 = r29
											goto block2
										} else {
											frame.Fail()
											goto block3
										}
									} else {
										goto block3
									}
								} else {
									frame.Fail()
									goto block3
								}
							} else {
								goto block3
							}
						} else {
							frame.Fail()
							goto block3
						}
					} else {
						goto block3
					}
				} else {
					frame.Fail()
					goto block3
				}
			} else {
				goto block3
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block2:
	return
block3:
	return
}
func Literal(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 int
	var r2 rune
	var r3 string
	var r4 *RuneLiteral
	var r5 int
	var r6 string
	var r7 string
	var r8 *StringLiteral
	var r9 int
	var r10 int
	var r11 string
	var r12 *IntLiteral
	var r13 int
	var r14 bool
	var r15 string
	var r16 *BoolLiteral
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = DecodeRune(frame)
	if frame.Flow == 0 {
		r3 = frame.Slice(r1)
		S(frame)
		if frame.Flow == 0 {
			r4 = &RuneLiteral{Text: r3, Value: r2}
			ret0 = r4
			goto block4
		} else {
			goto block1
		}
	} else {
		goto block1
	}
block1:
	frame.Recover(r0)
	r5 = frame.Checkpoint()
	r6 = DecodeString(frame)
	if frame.Flow == 0 {
		r7 = frame.Slice(r5)
		S(frame)
		if frame.Flow == 0 {
			r8 = &StringLiteral{Text: r7, Value: r6}
			ret0 = r8
			goto block4
		} else {
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r0)
	r9 = frame.Checkpoint()
	r10 = DecodeInt(frame)
	if frame.Flow == 0 {
		r11 = frame.Slice(r9)
		S(frame)
		if frame.Flow == 0 {
			r12 = &IntLiteral{Text: r11, Value: r10}
			ret0 = r12
			goto block4
		} else {
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r13 = frame.Checkpoint()
	r14 = DecodeBool(frame)
	if frame.Flow == 0 {
		r15 = frame.Slice(r13)
		EndKeyword(frame)
		if frame.Flow == 0 {
			r16 = &BoolLiteral{Text: r15, Value: r14}
			ret0 = r16
			goto block4
		} else {
			goto block5
		}
	} else {
		goto block5
	}
block4:
	return
block5:
	return
}
func BinaryOperator(frame *runtime.State) (ret0 string) {
	var r0 int
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 rune
	var r15 bool
	var r16 int
	var r17 rune
	var r18 rune
	var r19 bool
	var r20 rune
	var r21 rune
	var r22 bool
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 rune
	var r27 bool
	var r28 int
	var r29 rune
	var r30 rune
	var r31 bool
	var r32 rune
	var r33 bool
	var r34 rune
	var r35 bool
	var r36 rune
	var r37 bool
	var r38 rune
	var r39 bool
	var r40 rune
	var r41 bool
	var r42 rune
	var r43 bool
	var r44 rune
	var r45 bool
	var r46 string
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		r3 = '+'
		r4 = r2 == r3
		if r4 {
			goto block1
		} else {
			r5 = '-'
			r6 = r2 == r5
			if r6 {
				goto block1
			} else {
				r7 = '*'
				r8 = r2 == r7
				if r8 {
					goto block1
				} else {
					r9 = '/'
					r10 = r2 == r9
					if r10 {
						goto block1
					} else {
						frame.Fail()
						goto block2
					}
				}
			}
		}
	} else {
		goto block2
	}
block1:
	frame.Consume()
	goto block9
block2:
	frame.Recover(r1)
	r11 = frame.Peek()
	if frame.Flow == 0 {
		r12 = '<'
		r13 = r11 == r12
		if r13 {
			goto block3
		} else {
			r14 = '>'
			r15 = r11 == r14
			if r15 {
				goto block3
			} else {
				frame.Fail()
				goto block5
			}
		}
	} else {
		goto block5
	}
block3:
	frame.Consume()
	r16 = frame.Checkpoint()
	r17 = frame.Peek()
	if frame.Flow == 0 {
		r18 = '='
		r19 = r17 == r18
		if r19 {
			frame.Consume()
			goto block9
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r16)
	goto block9
block5:
	frame.Recover(r1)
	r20 = frame.Peek()
	if frame.Flow == 0 {
		r21 = '!'
		r22 = r20 == r21
		if r22 {
			goto block6
		} else {
			r23 = '='
			r24 = r20 == r23
			if r24 {
				goto block6
			} else {
				frame.Fail()
				goto block10
			}
		}
	} else {
		goto block10
	}
block6:
	frame.Consume()
	r25 = frame.Peek()
	if frame.Flow == 0 {
		r26 = '='
		r27 = r25 == r26
		if r27 {
			frame.Consume()
			r28 = frame.LookaheadBegin()
			r29 = frame.Peek()
			if frame.Flow == 0 {
				r30 = '+'
				r31 = r29 == r30
				if r31 {
					goto block7
				} else {
					r32 = '-'
					r33 = r29 == r32
					if r33 {
						goto block7
					} else {
						r34 = '*'
						r35 = r29 == r34
						if r35 {
							goto block7
						} else {
							r36 = '/'
							r37 = r29 == r36
							if r37 {
								goto block7
							} else {
								r38 = '<'
								r39 = r29 == r38
								if r39 {
									goto block7
								} else {
									r40 = '>'
									r41 = r29 == r40
									if r41 {
										goto block7
									} else {
										r42 = '!'
										r43 = r29 == r42
										if r43 {
											goto block7
										} else {
											r44 = '='
											r45 = r29 == r44
											if r45 {
												goto block7
											} else {
												frame.Fail()
												goto block8
											}
										}
									}
								}
							}
						}
					}
				}
			} else {
				goto block8
			}
		} else {
			frame.Fail()
			goto block10
		}
	} else {
		goto block10
	}
block7:
	frame.Consume()
	frame.LookaheadFail(r28)
	goto block10
block8:
	frame.LookaheadNormal(r28)
	goto block9
block9:
	r46 = frame.Slice(r0)
	S(frame)
	if frame.Flow == 0 {
		ret0 = r46
		return
	} else {
		goto block10
	}
block10:
	return
}
func StringMatchExpr(frame *runtime.State) (ret0 *StringMatch) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 TextMatch
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 *StringMatch
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '/'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = ParseMatchChoice(frame)
				if frame.Flow == 0 {
					r4 = frame.Peek()
					if frame.Flow == 0 {
						r5 = '/'
						r6 = r4 == r5
						if r6 {
							frame.Consume()
							S(frame)
							if frame.Flow == 0 {
								r7 = &StringMatch{Match: r3}
								ret0 = r7
								return
							} else {
								goto block1
							}
						} else {
							frame.Fail()
							goto block1
						}
					} else {
						goto block1
					}
				} else {
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func RuneMatchExpr(frame *runtime.State) (ret0 *RuneMatch) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 *RuneRangeMatch
	var r4 *RuneMatch
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '$'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = MatchRune(frame)
				if frame.Flow == 0 {
					r4 = &RuneMatch{Match: r3}
					ret0 = r4
					return
				} else {
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseStructTypeRef(frame *runtime.State) (ret0 *TypeRef) {
	var r0 *Id
	var r1 *TypeRef
	r0 = Ident(frame)
	if frame.Flow == 0 {
		r1 = &TypeRef{Name: r0}
		ret0 = r1
		return
	} else {
		return
	}
}
func ParseListTypeRef(frame *runtime.State) (ret0 *ListTypeRef) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 ASTTypeRef
	var r7 *ListTypeRef
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '['
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				r4 = ']'
				r5 = r3 == r4
				if r5 {
					frame.Consume()
					r6 = ParseTypeRef(frame)
					if frame.Flow == 0 {
						r7 = &ListTypeRef{Type: r6}
						ret0 = r7
						return
					} else {
						goto block1
					}
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseTypeRef(frame *runtime.State) (ret0 ASTTypeRef) {
	var r0 int
	var r1 *TypeRef
	var r2 *ListTypeRef
	r0 = frame.Checkpoint()
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block1
	} else {
		frame.Recover(r0)
		r2 = ParseListTypeRef(frame)
		if frame.Flow == 0 {
			ret0 = r2
			goto block1
		} else {
			return
		}
	}
block1:
	return
}
func ParseDestructure(frame *runtime.State) (ret0 Destructure) {
	var r0 int
	var r1 *TypeRef
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 []*DestructureField
	var r6 []*DestructureField
	var r7 int
	var r8 *Id
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 Destructure
	var r13 *DestructureField
	var r14 []*DestructureField
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 *DestructureStruct
	var r19 *ListTypeRef
	var r20 rune
	var r21 rune
	var r22 bool
	var r23 []Destructure
	var r24 []Destructure
	var r25 int
	var r26 Destructure
	var r27 []Destructure
	var r28 rune
	var r29 rune
	var r30 bool
	var r31 *DestructureList
	var r32 ASTExpr
	var r33 *DestructureValue
	r0 = frame.Checkpoint()
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		r2 = frame.Peek()
		if frame.Flow == 0 {
			r3 = '{'
			r4 = r2 == r3
			if r4 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r5 = []*DestructureField{}
					r6 = r5
					goto block1
				} else {
					goto block3
				}
			} else {
				frame.Fail()
				goto block3
			}
		} else {
			goto block3
		}
	} else {
		goto block3
	}
block1:
	r7 = frame.Checkpoint()
	r8 = Ident(frame)
	if frame.Flow == 0 {
		r9 = frame.Peek()
		if frame.Flow == 0 {
			r10 = ':'
			r11 = r9 == r10
			if r11 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r12 = ParseDestructure(frame)
					if frame.Flow == 0 {
						r13 = &DestructureField{Name: r8, Destructure: r12}
						r14 = append(r6, r13)
						r6 = r14
						goto block1
					} else {
						goto block2
					}
				} else {
					goto block2
				}
			} else {
				frame.Fail()
				goto block2
			}
		} else {
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r7)
	r15 = frame.Peek()
	if frame.Flow == 0 {
		r16 = '}'
		r17 = r15 == r16
		if r17 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r18 = &DestructureStruct{Type: r1, Args: r6}
				ret0 = r18
				goto block6
			} else {
				goto block3
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r19 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		r20 = frame.Peek()
		if frame.Flow == 0 {
			r21 = '{'
			r22 = r20 == r21
			if r22 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r23 = []Destructure{}
					r24 = r23
					goto block4
				} else {
					goto block5
				}
			} else {
				frame.Fail()
				goto block5
			}
		} else {
			goto block5
		}
	} else {
		goto block5
	}
block4:
	r25 = frame.Checkpoint()
	r26 = ParseDestructure(frame)
	if frame.Flow == 0 {
		r27 = append(r24, r26)
		r24 = r27
		goto block4
	} else {
		frame.Recover(r25)
		r28 = frame.Peek()
		if frame.Flow == 0 {
			r29 = '}'
			r30 = r28 == r29
			if r30 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r31 = &DestructureList{Type: r19, Args: r24}
					ret0 = r31
					goto block6
				} else {
					goto block5
				}
			} else {
				frame.Fail()
				goto block5
			}
		} else {
			goto block5
		}
	}
block5:
	frame.Recover(r0)
	r32 = Literal(frame)
	if frame.Flow == 0 {
		r33 = &DestructureValue{Expr: r32}
		ret0 = r33
		goto block6
	} else {
		return
	}
block6:
	return
}
func ParseRuneFilterRune(frame *runtime.State) (ret0 rune) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 int
	var r12 rune
	var r13 rune
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = ']'
		r3 = r1 == r2
		if r3 {
			goto block1
		} else {
			r4 = '-'
			r5 = r1 == r4
			if r5 {
				goto block1
			} else {
				r6 = '\\'
				r7 = r1 == r6
				if r7 {
					goto block1
				} else {
					frame.Consume()
					ret0 = r1
					goto block3
				}
			}
		}
	} else {
		goto block2
	}
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(r0)
	r8 = frame.Peek()
	if frame.Flow == 0 {
		r9 = '\\'
		r10 = r8 == r9
		if r10 {
			frame.Consume()
			r11 = frame.Checkpoint()
			r12 = EscapedChar(frame)
			if frame.Flow == 0 {
				ret0 = r12
				goto block3
			} else {
				frame.Recover(r11)
				r13 = frame.Peek()
				if frame.Flow == 0 {
					frame.Consume()
					ret0 = r13
					goto block3
				} else {
					goto block4
				}
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block3:
	return
block4:
	return
}
func ParseRuneFilter(frame *runtime.State) (ret0 *RuneFilter) {
	var r0 rune
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 rune
	var r7 *RuneFilter
	r0 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = frame.Peek()
		if frame.Flow == 0 {
			r3 = '-'
			r4 = r2 == r3
			if r4 {
				frame.Consume()
				r5 = ParseRuneFilterRune(frame)
				if frame.Flow == 0 {
					r6 = r5
					goto block2
				} else {
					goto block1
				}
			} else {
				frame.Fail()
				goto block1
			}
		} else {
			goto block1
		}
	} else {
		return
	}
block1:
	frame.Recover(r1)
	r6 = r0
	goto block2
block2:
	r7 = &RuneFilter{Min: r0, Max: r6}
	ret0 = r7
	return
}
func MatchRune(frame *runtime.State) (ret0 *RuneRangeMatch) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 bool
	var r4 []*RuneFilter
	var r5 int
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 bool
	var r10 bool
	var r11 []*RuneFilter
	var r12 int
	var r13 *RuneFilter
	var r14 []*RuneFilter
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 *RuneRangeMatch
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '['
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = false
			r4 = []*RuneFilter{}
			r5 = frame.Checkpoint()
			r6 = frame.Peek()
			if frame.Flow == 0 {
				r7 = '^'
				r8 = r6 == r7
				if r8 {
					frame.Consume()
					r9 = true
					r10, r11 = r9, r4
					goto block2
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block1:
	frame.Recover(r5)
	r10, r11 = r3, r4
	goto block2
block2:
	r12 = frame.Checkpoint()
	r13 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		r14 = append(r11, r13)
		r10, r11 = r10, r14
		goto block2
	} else {
		frame.Recover(r12)
		r15 = frame.Peek()
		if frame.Flow == 0 {
			r16 = ']'
			r17 = r15 == r16
			if r17 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r18 = &RuneRangeMatch{Invert: r10, Filters: r11}
					ret0 = r18
					return
				} else {
					goto block3
				}
			} else {
				frame.Fail()
				goto block3
			}
		} else {
			goto block3
		}
	}
block3:
	return
}
func Atom(frame *runtime.State) (ret0 TextMatch) {
	var r0 int
	var r1 *RuneRangeMatch
	var r2 string
	var r3 *StringLiteralMatch
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 TextMatch
	var r8 rune
	var r9 rune
	var r10 bool
	r0 = frame.Checkpoint()
	r1 = MatchRune(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block2
	} else {
		frame.Recover(r0)
		r2 = DecodeString(frame)
		if frame.Flow == 0 {
			S(frame)
			if frame.Flow == 0 {
				r3 = &StringLiteralMatch{Value: r2}
				ret0 = r3
				goto block2
			} else {
				goto block1
			}
		} else {
			goto block1
		}
	}
block1:
	frame.Recover(r0)
	r4 = frame.Peek()
	if frame.Flow == 0 {
		r5 = '('
		r6 = r4 == r5
		if r6 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r7 = ParseMatchChoice(frame)
				if frame.Flow == 0 {
					r8 = frame.Peek()
					if frame.Flow == 0 {
						r9 = ')'
						r10 = r8 == r9
						if r10 {
							frame.Consume()
							S(frame)
							if frame.Flow == 0 {
								ret0 = r7
								goto block2
							} else {
								goto block3
							}
						} else {
							frame.Fail()
							goto block3
						}
					} else {
						goto block3
					}
				} else {
					goto block3
				}
			} else {
				goto block3
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block2:
	return
block3:
	return
}
func MatchPostfix(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 int
	var r6 *MatchRepeat
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 int
	var r11 *MatchRepeat
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 []TextMatch
	var r16 *MatchSequence
	var r17 []TextMatch
	var r18 *MatchChoice
	r0 = Atom(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = frame.Peek()
		if frame.Flow == 0 {
			r3 = '*'
			r4 = r2 == r3
			if r4 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r5 = 0
					r6 = &MatchRepeat{Match: r0, Min: r5}
					ret0 = r6
					goto block4
				} else {
					goto block1
				}
			} else {
				frame.Fail()
				goto block1
			}
		} else {
			goto block1
		}
	} else {
		return
	}
block1:
	frame.Recover(r1)
	r7 = frame.Peek()
	if frame.Flow == 0 {
		r8 = '+'
		r9 = r7 == r8
		if r9 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r10 = 1
				r11 = &MatchRepeat{Match: r0, Min: r10}
				ret0 = r11
				goto block4
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r1)
	r12 = frame.Peek()
	if frame.Flow == 0 {
		r13 = '?'
		r14 = r12 == r13
		if r14 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r15 = []TextMatch{}
				r16 = &MatchSequence{Matches: r15}
				r17 = []TextMatch{r0, r16}
				r18 = &MatchChoice{Matches: r17}
				ret0 = r18
				goto block4
			} else {
				goto block3
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r1)
	ret0 = r0
	goto block4
block4:
	return
}
func MatchPrefix(frame *runtime.State) (ret0 TextMatch) {
	var r0 int
	var r1 bool
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 bool
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 TextMatch
	var r12 *MatchLookahead
	var r13 TextMatch
	r0 = frame.Checkpoint()
	r1 = false
	r2 = frame.Checkpoint()
	r3 = frame.Peek()
	if frame.Flow == 0 {
		r4 = '!'
		r5 = r3 == r4
		if r5 {
			frame.Consume()
			r6 = true
			r7 = r6
			goto block2
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	frame.Recover(r2)
	r8 = frame.Peek()
	if frame.Flow == 0 {
		r9 = '&'
		r10 = r8 == r9
		if r10 {
			frame.Consume()
			r7 = r1
			goto block2
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block2:
	S(frame)
	if frame.Flow == 0 {
		r11 = MatchPostfix(frame)
		if frame.Flow == 0 {
			r12 = &MatchLookahead{Invert: r7, Match: r11}
			ret0 = r12
			goto block4
		} else {
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r13 = MatchPostfix(frame)
	if frame.Flow == 0 {
		ret0 = r13
		goto block4
	} else {
		return
	}
block4:
	return
}
func Sequence(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 []TextMatch
	var r3 TextMatch
	var r4 []TextMatch
	var r5 []TextMatch
	var r6 int
	var r7 TextMatch
	var r8 []TextMatch
	var r9 *MatchSequence
	r0 = MatchPrefix(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = []TextMatch{r0}
		r3 = MatchPrefix(frame)
		if frame.Flow == 0 {
			r4 = append(r2, r3)
			r5 = r4
			goto block1
		} else {
			frame.Recover(r1)
			ret0 = r0
			goto block2
		}
	} else {
		return
	}
block1:
	r6 = frame.Checkpoint()
	r7 = MatchPrefix(frame)
	if frame.Flow == 0 {
		r8 = append(r5, r7)
		r5 = r8
		goto block1
	} else {
		frame.Recover(r6)
		r9 = &MatchSequence{Matches: r5}
		ret0 = r9
		goto block2
	}
block2:
	return
}
func ParseMatchChoice(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 []TextMatch
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 TextMatch
	var r7 []TextMatch
	var r8 []TextMatch
	var r9 int
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 TextMatch
	var r14 []TextMatch
	var r15 *MatchChoice
	r0 = Sequence(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = []TextMatch{r0}
		r3 = frame.Peek()
		if frame.Flow == 0 {
			r4 = '|'
			r5 = r3 == r4
			if r5 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r6 = Sequence(frame)
					if frame.Flow == 0 {
						r7 = append(r2, r6)
						r8 = r7
						goto block1
					} else {
						goto block3
					}
				} else {
					goto block3
				}
			} else {
				frame.Fail()
				goto block3
			}
		} else {
			goto block3
		}
	} else {
		return
	}
block1:
	r9 = frame.Checkpoint()
	r10 = frame.Peek()
	if frame.Flow == 0 {
		r11 = '|'
		r12 = r10 == r11
		if r12 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r13 = Sequence(frame)
				if frame.Flow == 0 {
					r14 = append(r8, r13)
					r8 = r14
					goto block1
				} else {
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r9)
	r15 = &MatchChoice{Matches: r8}
	ret0 = r15
	goto block4
block3:
	frame.Recover(r1)
	ret0 = r0
	goto block4
block4:
	return
}
func ParseExprList(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 []ASTExpr
	var r1 int
	var r2 ASTExpr
	var r3 []ASTExpr
	var r4 []ASTExpr
	var r5 int
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 ASTExpr
	var r10 []ASTExpr
	var r11 []ASTExpr
	r0 = []ASTExpr{}
	r1 = frame.Checkpoint()
	r2 = ParseExpr(frame)
	if frame.Flow == 0 {
		r3 = append(r0, r2)
		r4 = r3
		goto block1
	} else {
		frame.Recover(r1)
		r11 = r0
		goto block3
	}
block1:
	r5 = frame.Checkpoint()
	r6 = frame.Peek()
	if frame.Flow == 0 {
		r7 = ','
		r8 = r6 == r7
		if r8 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r9 = ParseExpr(frame)
				if frame.Flow == 0 {
					r10 = append(r4, r9)
					r4 = r10
					goto block1
				} else {
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r5)
	r11 = r4
	goto block3
block3:
	ret0 = r11
	return
}
func ParseNamedExpr(frame *runtime.State) (ret0 *NamedExpr) {
	var r0 *Id
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 ASTExpr
	var r5 *NamedExpr
	r0 = Ident(frame)
	if frame.Flow == 0 {
		r1 = frame.Peek()
		if frame.Flow == 0 {
			r2 = ':'
			r3 = r1 == r2
			if r3 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r4 = ParseExpr(frame)
					if frame.Flow == 0 {
						r5 = &NamedExpr{Name: r0, Expr: r4}
						ret0 = r5
						return
					} else {
						goto block1
					}
				} else {
					goto block1
				}
			} else {
				frame.Fail()
				goto block1
			}
		} else {
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseNamedExprList(frame *runtime.State) (ret0 []*NamedExpr) {
	var r0 []*NamedExpr
	var r1 int
	var r2 *NamedExpr
	var r3 []*NamedExpr
	var r4 []*NamedExpr
	var r5 int
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 *NamedExpr
	var r10 []*NamedExpr
	var r11 []*NamedExpr
	r0 = []*NamedExpr{}
	r1 = frame.Checkpoint()
	r2 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		r3 = append(r0, r2)
		r4 = r3
		goto block1
	} else {
		frame.Recover(r1)
		r11 = r0
		goto block3
	}
block1:
	r5 = frame.Checkpoint()
	r6 = frame.Peek()
	if frame.Flow == 0 {
		r7 = ','
		r8 = r6 == r7
		if r8 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r9 = ParseNamedExpr(frame)
				if frame.Flow == 0 {
					r10 = append(r4, r9)
					r4 = r10
					goto block1
				} else {
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r5)
	r11 = r4
	goto block3
block3:
	ret0 = r11
	return
}
func ParseTypeList(frame *runtime.State) (ret0 []ASTTypeRef) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 []ASTTypeRef
	var r4 int
	var r5 ASTTypeRef
	var r6 []ASTTypeRef
	var r7 []ASTTypeRef
	var r8 int
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 ASTTypeRef
	var r13 []ASTTypeRef
	var r14 []ASTTypeRef
	var r15 rune
	var r16 rune
	var r17 bool
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '('
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = []ASTTypeRef{}
				r4 = frame.Checkpoint()
				r5 = ParseTypeRef(frame)
				if frame.Flow == 0 {
					r6 = append(r3, r5)
					r7 = r6
					goto block1
				} else {
					frame.Recover(r4)
					r14 = r3
					goto block3
				}
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block1:
	r8 = frame.Checkpoint()
	r9 = frame.Peek()
	if frame.Flow == 0 {
		r10 = ','
		r11 = r9 == r10
		if r11 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r12 = ParseTypeRef(frame)
				if frame.Flow == 0 {
					r13 = append(r7, r12)
					r7 = r13
					goto block1
				} else {
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r8)
	r14 = r7
	goto block3
block3:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		r16 = ')'
		r17 = r15 == r16
		if r17 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				ret0 = r14
				return
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	return
}
func ParseExpr(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 ASTExpr
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 []ASTExpr
	var r18 *Slice
	var r19 rune
	var r20 rune
	var r21 bool
	var r22 rune
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 rune
	var r27 bool
	var r28 rune
	var r29 rune
	var r30 bool
	var r31 rune
	var r32 rune
	var r33 bool
	var r34 rune
	var r35 rune
	var r36 bool
	var r37 rune
	var r38 rune
	var r39 bool
	var r40 rune
	var r41 rune
	var r42 bool
	var r43 *Position
	var r44 rune
	var r45 rune
	var r46 bool
	var r47 rune
	var r48 rune
	var r49 bool
	var r50 rune
	var r51 rune
	var r52 bool
	var r53 rune
	var r54 rune
	var r55 bool
	var r56 rune
	var r57 rune
	var r58 bool
	var r59 rune
	var r60 rune
	var r61 bool
	var r62 ASTTypeRef
	var r63 ASTExpr
	var r64 *Coerce
	var r65 rune
	var r66 rune
	var r67 bool
	var r68 rune
	var r69 rune
	var r70 bool
	var r71 rune
	var r72 rune
	var r73 bool
	var r74 rune
	var r75 rune
	var r76 bool
	var r77 rune
	var r78 rune
	var r79 bool
	var r80 rune
	var r81 rune
	var r82 bool
	var r83 *Id
	var r84 ASTExpr
	var r85 *GetName
	var r86 *Append
	var r87 *Assign
	var r88 *Id
	var r89 rune
	var r90 rune
	var r91 bool
	var r92 rune
	var r93 rune
	var r94 bool
	var r95 *Call
	var r96 string
	var r97 ASTExpr
	var r98 ASTExpr
	var r99 *BinaryOp
	var r100 *TypeRef
	var r101 rune
	var r102 rune
	var r103 bool
	var r104 []*NamedExpr
	var r105 rune
	var r106 rune
	var r107 bool
	var r108 *Construct
	var r109 *ListTypeRef
	var r110 rune
	var r111 rune
	var r112 bool
	var r113 []ASTExpr
	var r114 rune
	var r115 rune
	var r116 bool
	var r117 *ConstructList
	var r118 *StringMatch
	var r119 *RuneMatch
	var r120 *Id
	var r121 *GetName
	r0 = frame.Checkpoint()
	r1 = Literal(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block9
	} else {
		frame.Recover(r0)
		r2 = frame.Peek()
		if frame.Flow == 0 {
			r3 = 's'
			r4 = r2 == r3
			if r4 {
				frame.Consume()
				r5 = frame.Peek()
				if frame.Flow == 0 {
					r6 = 'l'
					r7 = r5 == r6
					if r7 {
						frame.Consume()
						r8 = frame.Peek()
						if frame.Flow == 0 {
							r9 = 'i'
							r10 = r8 == r9
							if r10 {
								frame.Consume()
								r11 = frame.Peek()
								if frame.Flow == 0 {
									r12 = 'c'
									r13 = r11 == r12
									if r13 {
										frame.Consume()
										r14 = frame.Peek()
										if frame.Flow == 0 {
											r15 = 'e'
											r16 = r14 == r15
											if r16 {
												frame.Consume()
												EndKeyword(frame)
												if frame.Flow == 0 {
													r17 = ParseCodeBlock(frame)
													if frame.Flow == 0 {
														r18 = &Slice{Block: r17}
														ret0 = r18
														goto block9
													} else {
														goto block1
													}
												} else {
													goto block1
												}
											} else {
												frame.Fail()
												goto block1
											}
										} else {
											goto block1
										}
									} else {
										frame.Fail()
										goto block1
									}
								} else {
									goto block1
								}
							} else {
								frame.Fail()
								goto block1
							}
						} else {
							goto block1
						}
					} else {
						frame.Fail()
						goto block1
					}
				} else {
					goto block1
				}
			} else {
				frame.Fail()
				goto block1
			}
		} else {
			goto block1
		}
	}
block1:
	frame.Recover(r0)
	r19 = frame.Peek()
	if frame.Flow == 0 {
		r20 = 'p'
		r21 = r19 == r20
		if r21 {
			frame.Consume()
			r22 = frame.Peek()
			if frame.Flow == 0 {
				r23 = 'o'
				r24 = r22 == r23
				if r24 {
					frame.Consume()
					r25 = frame.Peek()
					if frame.Flow == 0 {
						r26 = 's'
						r27 = r25 == r26
						if r27 {
							frame.Consume()
							r28 = frame.Peek()
							if frame.Flow == 0 {
								r29 = 'i'
								r30 = r28 == r29
								if r30 {
									frame.Consume()
									r31 = frame.Peek()
									if frame.Flow == 0 {
										r32 = 't'
										r33 = r31 == r32
										if r33 {
											frame.Consume()
											r34 = frame.Peek()
											if frame.Flow == 0 {
												r35 = 'i'
												r36 = r34 == r35
												if r36 {
													frame.Consume()
													r37 = frame.Peek()
													if frame.Flow == 0 {
														r38 = 'o'
														r39 = r37 == r38
														if r39 {
															frame.Consume()
															r40 = frame.Peek()
															if frame.Flow == 0 {
																r41 = 'n'
																r42 = r40 == r41
																if r42 {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		r43 = &Position{}
																		ret0 = r43
																		goto block9
																	} else {
																		goto block2
																	}
																} else {
																	frame.Fail()
																	goto block2
																}
															} else {
																goto block2
															}
														} else {
															frame.Fail()
															goto block2
														}
													} else {
														goto block2
													}
												} else {
													frame.Fail()
													goto block2
												}
											} else {
												goto block2
											}
										} else {
											frame.Fail()
											goto block2
										}
									} else {
										goto block2
									}
								} else {
									frame.Fail()
									goto block2
								}
							} else {
								goto block2
							}
						} else {
							frame.Fail()
							goto block2
						}
					} else {
						goto block2
					}
				} else {
					frame.Fail()
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r0)
	r44 = frame.Peek()
	if frame.Flow == 0 {
		r45 = 'c'
		r46 = r44 == r45
		if r46 {
			frame.Consume()
			r47 = frame.Peek()
			if frame.Flow == 0 {
				r48 = 'o'
				r49 = r47 == r48
				if r49 {
					frame.Consume()
					r50 = frame.Peek()
					if frame.Flow == 0 {
						r51 = 'e'
						r52 = r50 == r51
						if r52 {
							frame.Consume()
							r53 = frame.Peek()
							if frame.Flow == 0 {
								r54 = 'r'
								r55 = r53 == r54
								if r55 {
									frame.Consume()
									r56 = frame.Peek()
									if frame.Flow == 0 {
										r57 = 'c'
										r58 = r56 == r57
										if r58 {
											frame.Consume()
											r59 = frame.Peek()
											if frame.Flow == 0 {
												r60 = 'e'
												r61 = r59 == r60
												if r61 {
													frame.Consume()
													S(frame)
													if frame.Flow == 0 {
														r62 = ParseTypeRef(frame)
														if frame.Flow == 0 {
															r63 = ParseExpr(frame)
															if frame.Flow == 0 {
																r64 = &Coerce{Type: r62, Expr: r63}
																ret0 = r64
																goto block9
															} else {
																goto block3
															}
														} else {
															goto block3
														}
													} else {
														goto block3
													}
												} else {
													frame.Fail()
													goto block3
												}
											} else {
												goto block3
											}
										} else {
											frame.Fail()
											goto block3
										}
									} else {
										goto block3
									}
								} else {
									frame.Fail()
									goto block3
								}
							} else {
								goto block3
							}
						} else {
							frame.Fail()
							goto block3
						}
					} else {
						goto block3
					}
				} else {
					frame.Fail()
					goto block3
				}
			} else {
				goto block3
			}
		} else {
			frame.Fail()
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r65 = frame.Peek()
	if frame.Flow == 0 {
		r66 = 'a'
		r67 = r65 == r66
		if r67 {
			frame.Consume()
			r68 = frame.Peek()
			if frame.Flow == 0 {
				r69 = 'p'
				r70 = r68 == r69
				if r70 {
					frame.Consume()
					r71 = frame.Peek()
					if frame.Flow == 0 {
						r72 = 'p'
						r73 = r71 == r72
						if r73 {
							frame.Consume()
							r74 = frame.Peek()
							if frame.Flow == 0 {
								r75 = 'e'
								r76 = r74 == r75
								if r76 {
									frame.Consume()
									r77 = frame.Peek()
									if frame.Flow == 0 {
										r78 = 'n'
										r79 = r77 == r78
										if r79 {
											frame.Consume()
											r80 = frame.Peek()
											if frame.Flow == 0 {
												r81 = 'd'
												r82 = r80 == r81
												if r82 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r83 = Ident(frame)
														if frame.Flow == 0 {
															r84 = ParseExpr(frame)
															if frame.Flow == 0 {
																r85 = &GetName{Name: r83}
																r86 = &Append{List: r85, Expr: r84}
																r87 = &Assign{Expr: r86, Name: r83}
																ret0 = r87
																goto block9
															} else {
																goto block4
															}
														} else {
															goto block4
														}
													} else {
														goto block4
													}
												} else {
													frame.Fail()
													goto block4
												}
											} else {
												goto block4
											}
										} else {
											frame.Fail()
											goto block4
										}
									} else {
										goto block4
									}
								} else {
									frame.Fail()
									goto block4
								}
							} else {
								goto block4
							}
						} else {
							frame.Fail()
							goto block4
						}
					} else {
						goto block4
					}
				} else {
					frame.Fail()
					goto block4
				}
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r0)
	r88 = Ident(frame)
	if frame.Flow == 0 {
		r89 = frame.Peek()
		if frame.Flow == 0 {
			r90 = '('
			r91 = r89 == r90
			if r91 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r92 = frame.Peek()
					if frame.Flow == 0 {
						r93 = ')'
						r94 = r92 == r93
						if r94 {
							frame.Consume()
							S(frame)
							if frame.Flow == 0 {
								r95 = &Call{Name: r88}
								ret0 = r95
								goto block9
							} else {
								goto block5
							}
						} else {
							frame.Fail()
							goto block5
						}
					} else {
						goto block5
					}
				} else {
					goto block5
				}
			} else {
				frame.Fail()
				goto block5
			}
		} else {
			goto block5
		}
	} else {
		goto block5
	}
block5:
	frame.Recover(r0)
	r96 = BinaryOperator(frame)
	if frame.Flow == 0 {
		r97 = ParseExpr(frame)
		if frame.Flow == 0 {
			r98 = ParseExpr(frame)
			if frame.Flow == 0 {
				r99 = &BinaryOp{Left: r97, Op: r96, Right: r98}
				ret0 = r99
				goto block9
			} else {
				goto block6
			}
		} else {
			goto block6
		}
	} else {
		goto block6
	}
block6:
	frame.Recover(r0)
	r100 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		r101 = frame.Peek()
		if frame.Flow == 0 {
			r102 = '{'
			r103 = r101 == r102
			if r103 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r104 = ParseNamedExprList(frame)
					if frame.Flow == 0 {
						r105 = frame.Peek()
						if frame.Flow == 0 {
							r106 = '}'
							r107 = r105 == r106
							if r107 {
								frame.Consume()
								S(frame)
								if frame.Flow == 0 {
									r108 = &Construct{Type: r100, Args: r104}
									ret0 = r108
									goto block9
								} else {
									goto block7
								}
							} else {
								frame.Fail()
								goto block7
							}
						} else {
							goto block7
						}
					} else {
						goto block7
					}
				} else {
					goto block7
				}
			} else {
				frame.Fail()
				goto block7
			}
		} else {
			goto block7
		}
	} else {
		goto block7
	}
block7:
	frame.Recover(r0)
	r109 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		r110 = frame.Peek()
		if frame.Flow == 0 {
			r111 = '{'
			r112 = r110 == r111
			if r112 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r113 = ParseExprList(frame)
					if frame.Flow == 0 {
						r114 = frame.Peek()
						if frame.Flow == 0 {
							r115 = '}'
							r116 = r114 == r115
							if r116 {
								frame.Consume()
								S(frame)
								if frame.Flow == 0 {
									r117 = &ConstructList{Type: r109, Args: r113}
									ret0 = r117
									goto block9
								} else {
									goto block8
								}
							} else {
								frame.Fail()
								goto block8
							}
						} else {
							goto block8
						}
					} else {
						goto block8
					}
				} else {
					goto block8
				}
			} else {
				frame.Fail()
				goto block8
			}
		} else {
			goto block8
		}
	} else {
		goto block8
	}
block8:
	frame.Recover(r0)
	r118 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		ret0 = r118
		goto block9
	} else {
		frame.Recover(r0)
		r119 = RuneMatchExpr(frame)
		if frame.Flow == 0 {
			ret0 = r119
			goto block9
		} else {
			frame.Recover(r0)
			r120 = Ident(frame)
			if frame.Flow == 0 {
				r121 = &GetName{Name: r120}
				ret0 = r121
				goto block9
			} else {
				return
			}
		}
	}
block9:
	return
}
func ParseCompoundStatement(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 []ASTExpr
	var r14 int
	var r15 *Repeat
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 rune
	var r20 rune
	var r21 bool
	var r22 rune
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 rune
	var r27 bool
	var r28 []ASTExpr
	var r29 int
	var r30 *Repeat
	var r31 rune
	var r32 rune
	var r33 bool
	var r34 rune
	var r35 rune
	var r36 bool
	var r37 rune
	var r38 rune
	var r39 bool
	var r40 rune
	var r41 rune
	var r42 bool
	var r43 rune
	var r44 rune
	var r45 bool
	var r46 rune
	var r47 rune
	var r48 bool
	var r49 []ASTExpr
	var r50 [][]ASTExpr
	var r51 [][]ASTExpr
	var r52 int
	var r53 rune
	var r54 rune
	var r55 bool
	var r56 rune
	var r57 rune
	var r58 bool
	var r59 []ASTExpr
	var r60 [][]ASTExpr
	var r61 *Choice
	var r62 rune
	var r63 rune
	var r64 bool
	var r65 rune
	var r66 rune
	var r67 bool
	var r68 rune
	var r69 rune
	var r70 bool
	var r71 rune
	var r72 rune
	var r73 bool
	var r74 rune
	var r75 rune
	var r76 bool
	var r77 rune
	var r78 rune
	var r79 bool
	var r80 rune
	var r81 rune
	var r82 bool
	var r83 rune
	var r84 rune
	var r85 bool
	var r86 []ASTExpr
	var r87 *Optional
	var r88 rune
	var r89 rune
	var r90 bool
	var r91 rune
	var r92 rune
	var r93 bool
	var r94 ASTExpr
	var r95 []ASTExpr
	var r96 *If
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		r2 = 's'
		r3 = r1 == r2
		if r3 {
			frame.Consume()
			r4 = frame.Peek()
			if frame.Flow == 0 {
				r5 = 't'
				r6 = r4 == r5
				if r6 {
					frame.Consume()
					r7 = frame.Peek()
					if frame.Flow == 0 {
						r8 = 'a'
						r9 = r7 == r8
						if r9 {
							frame.Consume()
							r10 = frame.Peek()
							if frame.Flow == 0 {
								r11 = 'r'
								r12 = r10 == r11
								if r12 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										r13 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											r14 = 0
											r15 = &Repeat{Block: r13, Min: r14}
											ret0 = r15
											goto block7
										} else {
											goto block1
										}
									} else {
										goto block1
									}
								} else {
									frame.Fail()
									goto block1
								}
							} else {
								goto block1
							}
						} else {
							frame.Fail()
							goto block1
						}
					} else {
						goto block1
					}
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	frame.Recover(r0)
	r16 = frame.Peek()
	if frame.Flow == 0 {
		r17 = 'p'
		r18 = r16 == r17
		if r18 {
			frame.Consume()
			r19 = frame.Peek()
			if frame.Flow == 0 {
				r20 = 'l'
				r21 = r19 == r20
				if r21 {
					frame.Consume()
					r22 = frame.Peek()
					if frame.Flow == 0 {
						r23 = 'u'
						r24 = r22 == r23
						if r24 {
							frame.Consume()
							r25 = frame.Peek()
							if frame.Flow == 0 {
								r26 = 's'
								r27 = r25 == r26
								if r27 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										r28 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											r29 = 1
											r30 = &Repeat{Block: r28, Min: r29}
											ret0 = r30
											goto block7
										} else {
											goto block2
										}
									} else {
										goto block2
									}
								} else {
									frame.Fail()
									goto block2
								}
							} else {
								goto block2
							}
						} else {
							frame.Fail()
							goto block2
						}
					} else {
						goto block2
					}
				} else {
					frame.Fail()
					goto block2
				}
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block2:
	frame.Recover(r0)
	r31 = frame.Peek()
	if frame.Flow == 0 {
		r32 = 'c'
		r33 = r31 == r32
		if r33 {
			frame.Consume()
			r34 = frame.Peek()
			if frame.Flow == 0 {
				r35 = 'h'
				r36 = r34 == r35
				if r36 {
					frame.Consume()
					r37 = frame.Peek()
					if frame.Flow == 0 {
						r38 = 'o'
						r39 = r37 == r38
						if r39 {
							frame.Consume()
							r40 = frame.Peek()
							if frame.Flow == 0 {
								r41 = 'o'
								r42 = r40 == r41
								if r42 {
									frame.Consume()
									r43 = frame.Peek()
									if frame.Flow == 0 {
										r44 = 's'
										r45 = r43 == r44
										if r45 {
											frame.Consume()
											r46 = frame.Peek()
											if frame.Flow == 0 {
												r47 = 'e'
												r48 = r46 == r47
												if r48 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r49 = ParseCodeBlock(frame)
														if frame.Flow == 0 {
															r50 = [][]ASTExpr{r49}
															r51 = r50
															goto block3
														} else {
															goto block5
														}
													} else {
														goto block5
													}
												} else {
													frame.Fail()
													goto block5
												}
											} else {
												goto block5
											}
										} else {
											frame.Fail()
											goto block5
										}
									} else {
										goto block5
									}
								} else {
									frame.Fail()
									goto block5
								}
							} else {
								goto block5
							}
						} else {
							frame.Fail()
							goto block5
						}
					} else {
						goto block5
					}
				} else {
					frame.Fail()
					goto block5
				}
			} else {
				goto block5
			}
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block3:
	r52 = frame.Checkpoint()
	r53 = frame.Peek()
	if frame.Flow == 0 {
		r54 = 'o'
		r55 = r53 == r54
		if r55 {
			frame.Consume()
			r56 = frame.Peek()
			if frame.Flow == 0 {
				r57 = 'r'
				r58 = r56 == r57
				if r58 {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						r59 = ParseCodeBlock(frame)
						if frame.Flow == 0 {
							r60 = append(r51, r59)
							r51 = r60
							goto block3
						} else {
							goto block4
						}
					} else {
						goto block4
					}
				} else {
					frame.Fail()
					goto block4
				}
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r52)
	r61 = &Choice{Blocks: r51}
	ret0 = r61
	goto block7
block5:
	frame.Recover(r0)
	r62 = frame.Peek()
	if frame.Flow == 0 {
		r63 = 'q'
		r64 = r62 == r63
		if r64 {
			frame.Consume()
			r65 = frame.Peek()
			if frame.Flow == 0 {
				r66 = 'u'
				r67 = r65 == r66
				if r67 {
					frame.Consume()
					r68 = frame.Peek()
					if frame.Flow == 0 {
						r69 = 'e'
						r70 = r68 == r69
						if r70 {
							frame.Consume()
							r71 = frame.Peek()
							if frame.Flow == 0 {
								r72 = 's'
								r73 = r71 == r72
								if r73 {
									frame.Consume()
									r74 = frame.Peek()
									if frame.Flow == 0 {
										r75 = 't'
										r76 = r74 == r75
										if r76 {
											frame.Consume()
											r77 = frame.Peek()
											if frame.Flow == 0 {
												r78 = 'i'
												r79 = r77 == r78
												if r79 {
													frame.Consume()
													r80 = frame.Peek()
													if frame.Flow == 0 {
														r81 = 'o'
														r82 = r80 == r81
														if r82 {
															frame.Consume()
															r83 = frame.Peek()
															if frame.Flow == 0 {
																r84 = 'n'
																r85 = r83 == r84
																if r85 {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		r86 = ParseCodeBlock(frame)
																		if frame.Flow == 0 {
																			r87 = &Optional{Block: r86}
																			ret0 = r87
																			goto block7
																		} else {
																			goto block6
																		}
																	} else {
																		goto block6
																	}
																} else {
																	frame.Fail()
																	goto block6
																}
															} else {
																goto block6
															}
														} else {
															frame.Fail()
															goto block6
														}
													} else {
														goto block6
													}
												} else {
													frame.Fail()
													goto block6
												}
											} else {
												goto block6
											}
										} else {
											frame.Fail()
											goto block6
										}
									} else {
										goto block6
									}
								} else {
									frame.Fail()
									goto block6
								}
							} else {
								goto block6
							}
						} else {
							frame.Fail()
							goto block6
						}
					} else {
						goto block6
					}
				} else {
					frame.Fail()
					goto block6
				}
			} else {
				goto block6
			}
		} else {
			frame.Fail()
			goto block6
		}
	} else {
		goto block6
	}
block6:
	frame.Recover(r0)
	r88 = frame.Peek()
	if frame.Flow == 0 {
		r89 = 'i'
		r90 = r88 == r89
		if r90 {
			frame.Consume()
			r91 = frame.Peek()
			if frame.Flow == 0 {
				r92 = 'f'
				r93 = r91 == r92
				if r93 {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						r94 = ParseExpr(frame)
						if frame.Flow == 0 {
							r95 = ParseCodeBlock(frame)
							if frame.Flow == 0 {
								r96 = &If{Expr: r94, Block: r95}
								ret0 = r96
								goto block7
							} else {
								goto block8
							}
						} else {
							goto block8
						}
					} else {
						goto block8
					}
				} else {
					frame.Fail()
					goto block8
				}
			} else {
				goto block8
			}
		} else {
			frame.Fail()
			goto block8
		}
	} else {
		goto block8
	}
block7:
	return
block8:
	return
}
func EOS(frame *runtime.State) {
	var r0 rune
	var r1 rune
	var r2 bool
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = ';'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				return
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseStatement(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 ASTExpr
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 *Id
	var r12 ASTTypeRef
	var r13 ASTExpr
	var r14 int
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 ASTExpr
	var r19 ASTExpr
	var r20 bool
	var r21 *Assign
	var r22 rune
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 rune
	var r27 bool
	var r28 rune
	var r29 rune
	var r30 bool
	var r31 rune
	var r32 rune
	var r33 bool
	var r34 *Fail
	var r35 rune
	var r36 rune
	var r37 bool
	var r38 rune
	var r39 rune
	var r40 bool
	var r41 rune
	var r42 rune
	var r43 bool
	var r44 rune
	var r45 rune
	var r46 bool
	var r47 rune
	var r48 rune
	var r49 bool
	var r50 rune
	var r51 rune
	var r52 bool
	var r53 int
	var r54 rune
	var r55 rune
	var r56 bool
	var r57 []ASTExpr
	var r58 rune
	var r59 rune
	var r60 bool
	var r61 *Return
	var r62 ASTExpr
	var r63 []ASTExpr
	var r64 *Return
	var r65 []ASTExpr
	var r66 *Return
	var r67 *Id
	var r68 bool
	var r69 int
	var r70 rune
	var r71 rune
	var r72 bool
	var r73 rune
	var r74 rune
	var r75 bool
	var r76 bool
	var r77 bool
	var r78 rune
	var r79 rune
	var r80 bool
	var r81 ASTExpr
	var r82 *Assign
	var r83 ASTExpr
	r0 = frame.Checkpoint()
	r1 = ParseCompoundStatement(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block11
	} else {
		frame.Recover(r0)
		r2 = frame.Peek()
		if frame.Flow == 0 {
			r3 = 'v'
			r4 = r2 == r3
			if r4 {
				frame.Consume()
				r5 = frame.Peek()
				if frame.Flow == 0 {
					r6 = 'a'
					r7 = r5 == r6
					if r7 {
						frame.Consume()
						r8 = frame.Peek()
						if frame.Flow == 0 {
							r9 = 'r'
							r10 = r8 == r9
							if r10 {
								frame.Consume()
								EndKeyword(frame)
								if frame.Flow == 0 {
									r11 = Ident(frame)
									if frame.Flow == 0 {
										r12 = ParseTypeRef(frame)
										if frame.Flow == 0 {
											r13 = nil
											r14 = frame.Checkpoint()
											r15 = frame.Peek()
											if frame.Flow == 0 {
												r16 = '='
												r17 = r15 == r16
												if r17 {
													frame.Consume()
													S(frame)
													if frame.Flow == 0 {
														r18 = ParseExpr(frame)
														if frame.Flow == 0 {
															r19 = r18
															goto block2
														} else {
															goto block1
														}
													} else {
														goto block1
													}
												} else {
													frame.Fail()
													goto block1
												}
											} else {
												goto block1
											}
										} else {
											goto block3
										}
									} else {
										goto block3
									}
								} else {
									goto block3
								}
							} else {
								frame.Fail()
								goto block3
							}
						} else {
							goto block3
						}
					} else {
						frame.Fail()
						goto block3
					}
				} else {
					goto block3
				}
			} else {
				frame.Fail()
				goto block3
			}
		} else {
			goto block3
		}
	}
block1:
	frame.Recover(r14)
	r19 = r13
	goto block2
block2:
	EOS(frame)
	if frame.Flow == 0 {
		r20 = true
		r21 = &Assign{Expr: r19, Name: r11, Type: r12, Define: r20}
		ret0 = r21
		goto block11
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r22 = frame.Peek()
	if frame.Flow == 0 {
		r23 = 'f'
		r24 = r22 == r23
		if r24 {
			frame.Consume()
			r25 = frame.Peek()
			if frame.Flow == 0 {
				r26 = 'a'
				r27 = r25 == r26
				if r27 {
					frame.Consume()
					r28 = frame.Peek()
					if frame.Flow == 0 {
						r29 = 'i'
						r30 = r28 == r29
						if r30 {
							frame.Consume()
							r31 = frame.Peek()
							if frame.Flow == 0 {
								r32 = 'l'
								r33 = r31 == r32
								if r33 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										EOS(frame)
										if frame.Flow == 0 {
											r34 = &Fail{}
											ret0 = r34
											goto block11
										} else {
											goto block4
										}
									} else {
										goto block4
									}
								} else {
									frame.Fail()
									goto block4
								}
							} else {
								goto block4
							}
						} else {
							frame.Fail()
							goto block4
						}
					} else {
						goto block4
					}
				} else {
					frame.Fail()
					goto block4
				}
			} else {
				goto block4
			}
		} else {
			frame.Fail()
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r0)
	r35 = frame.Peek()
	if frame.Flow == 0 {
		r36 = 'r'
		r37 = r35 == r36
		if r37 {
			frame.Consume()
			r38 = frame.Peek()
			if frame.Flow == 0 {
				r39 = 'e'
				r40 = r38 == r39
				if r40 {
					frame.Consume()
					r41 = frame.Peek()
					if frame.Flow == 0 {
						r42 = 't'
						r43 = r41 == r42
						if r43 {
							frame.Consume()
							r44 = frame.Peek()
							if frame.Flow == 0 {
								r45 = 'u'
								r46 = r44 == r45
								if r46 {
									frame.Consume()
									r47 = frame.Peek()
									if frame.Flow == 0 {
										r48 = 'r'
										r49 = r47 == r48
										if r49 {
											frame.Consume()
											r50 = frame.Peek()
											if frame.Flow == 0 {
												r51 = 'n'
												r52 = r50 == r51
												if r52 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r53 = frame.Checkpoint()
														r54 = frame.Peek()
														if frame.Flow == 0 {
															r55 = '('
															r56 = r54 == r55
															if r56 {
																frame.Consume()
																S(frame)
																if frame.Flow == 0 {
																	r57 = ParseExprList(frame)
																	if frame.Flow == 0 {
																		r58 = frame.Peek()
																		if frame.Flow == 0 {
																			r59 = ')'
																			r60 = r58 == r59
																			if r60 {
																				frame.Consume()
																				S(frame)
																				if frame.Flow == 0 {
																					EOS(frame)
																					if frame.Flow == 0 {
																						r61 = &Return{Exprs: r57}
																						ret0 = r61
																						goto block11
																					} else {
																						goto block5
																					}
																				} else {
																					goto block5
																				}
																			} else {
																				frame.Fail()
																				goto block5
																			}
																		} else {
																			goto block5
																		}
																	} else {
																		goto block5
																	}
																} else {
																	goto block5
																}
															} else {
																frame.Fail()
																goto block5
															}
														} else {
															goto block5
														}
													} else {
														goto block7
													}
												} else {
													frame.Fail()
													goto block7
												}
											} else {
												goto block7
											}
										} else {
											frame.Fail()
											goto block7
										}
									} else {
										goto block7
									}
								} else {
									frame.Fail()
									goto block7
								}
							} else {
								goto block7
							}
						} else {
							frame.Fail()
							goto block7
						}
					} else {
						goto block7
					}
				} else {
					frame.Fail()
					goto block7
				}
			} else {
				goto block7
			}
		} else {
			frame.Fail()
			goto block7
		}
	} else {
		goto block7
	}
block5:
	frame.Recover(r53)
	r62 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			r63 = []ASTExpr{r62}
			r64 = &Return{Exprs: r63}
			ret0 = r64
			goto block11
		} else {
			goto block6
		}
	} else {
		goto block6
	}
block6:
	frame.Recover(r53)
	EOS(frame)
	if frame.Flow == 0 {
		r65 = []ASTExpr{}
		r66 = &Return{Exprs: r65}
		ret0 = r66
		goto block11
	} else {
		goto block7
	}
block7:
	frame.Recover(r0)
	r67 = Ident(frame)
	if frame.Flow == 0 {
		r68 = false
		r69 = frame.Checkpoint()
		r70 = frame.Peek()
		if frame.Flow == 0 {
			r71 = ':'
			r72 = r70 == r71
			if r72 {
				frame.Consume()
				r73 = frame.Peek()
				if frame.Flow == 0 {
					r74 = '='
					r75 = r73 == r74
					if r75 {
						frame.Consume()
						r76 = true
						r77 = r76
						goto block9
					} else {
						frame.Fail()
						goto block8
					}
				} else {
					goto block8
				}
			} else {
				frame.Fail()
				goto block8
			}
		} else {
			goto block8
		}
	} else {
		goto block10
	}
block8:
	frame.Recover(r69)
	r78 = frame.Peek()
	if frame.Flow == 0 {
		r79 = '='
		r80 = r78 == r79
		if r80 {
			frame.Consume()
			r77 = r68
			goto block9
		} else {
			frame.Fail()
			goto block10
		}
	} else {
		goto block10
	}
block9:
	S(frame)
	if frame.Flow == 0 {
		r81 = ParseExpr(frame)
		if frame.Flow == 0 {
			EOS(frame)
			if frame.Flow == 0 {
				r82 = &Assign{Expr: r81, Name: r67, Define: r77}
				ret0 = r82
				goto block11
			} else {
				goto block10
			}
		} else {
			goto block10
		}
	} else {
		goto block10
	}
block10:
	frame.Recover(r0)
	r83 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret0 = r83
			goto block11
		} else {
			goto block12
		}
	} else {
		goto block12
	}
block11:
	return
block12:
	return
}
func ParseCodeBlock(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 []ASTExpr
	var r4 []ASTExpr
	var r5 int
	var r6 ASTExpr
	var r7 []ASTExpr
	var r8 rune
	var r9 rune
	var r10 bool
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = '{'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = []ASTExpr{}
				r4 = r3
				goto block1
			} else {
				goto block2
			}
		} else {
			frame.Fail()
			goto block2
		}
	} else {
		goto block2
	}
block1:
	r5 = frame.Checkpoint()
	r6 = ParseStatement(frame)
	if frame.Flow == 0 {
		r7 = append(r4, r6)
		r4 = r7
		goto block1
	} else {
		frame.Recover(r5)
		r8 = frame.Peek()
		if frame.Flow == 0 {
			r9 = '}'
			r10 = r8 == r9
			if r10 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					ret0 = r4
					return
				} else {
					goto block2
				}
			} else {
				frame.Fail()
				goto block2
			}
		} else {
			goto block2
		}
	}
block2:
	return
}
func ParseStructDecl(frame *runtime.State) (ret0 *StructDecl) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 *Id
	var r19 ASTTypeRef
	var r20 int
	var r21 rune
	var r22 rune
	var r23 bool
	var r24 rune
	var r25 rune
	var r26 bool
	var r27 rune
	var r28 rune
	var r29 bool
	var r30 rune
	var r31 rune
	var r32 bool
	var r33 rune
	var r34 rune
	var r35 bool
	var r36 rune
	var r37 rune
	var r38 bool
	var r39 rune
	var r40 rune
	var r41 bool
	var r42 rune
	var r43 rune
	var r44 bool
	var r45 rune
	var r46 rune
	var r47 bool
	var r48 rune
	var r49 rune
	var r50 bool
	var r51 ASTTypeRef
	var r52 ASTTypeRef
	var r53 rune
	var r54 rune
	var r55 bool
	var r56 []*FieldDecl
	var r57 []*FieldDecl
	var r58 int
	var r59 *Id
	var r60 ASTTypeRef
	var r61 *FieldDecl
	var r62 []*FieldDecl
	var r63 rune
	var r64 rune
	var r65 bool
	var r66 *StructDecl
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = 's'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				r4 = 't'
				r5 = r3 == r4
				if r5 {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						r7 = 'r'
						r8 = r6 == r7
						if r8 {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								r10 = 'u'
								r11 = r9 == r10
								if r11 {
									frame.Consume()
									r12 = frame.Peek()
									if frame.Flow == 0 {
										r13 = 'c'
										r14 = r12 == r13
										if r14 {
											frame.Consume()
											r15 = frame.Peek()
											if frame.Flow == 0 {
												r16 = 't'
												r17 = r15 == r16
												if r17 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r18 = Ident(frame)
														if frame.Flow == 0 {
															r19 = nil
															r20 = frame.Checkpoint()
															r21 = frame.Peek()
															if frame.Flow == 0 {
																r22 = 'i'
																r23 = r21 == r22
																if r23 {
																	frame.Consume()
																	r24 = frame.Peek()
																	if frame.Flow == 0 {
																		r25 = 'm'
																		r26 = r24 == r25
																		if r26 {
																			frame.Consume()
																			r27 = frame.Peek()
																			if frame.Flow == 0 {
																				r28 = 'p'
																				r29 = r27 == r28
																				if r29 {
																					frame.Consume()
																					r30 = frame.Peek()
																					if frame.Flow == 0 {
																						r31 = 'l'
																						r32 = r30 == r31
																						if r32 {
																							frame.Consume()
																							r33 = frame.Peek()
																							if frame.Flow == 0 {
																								r34 = 'e'
																								r35 = r33 == r34
																								if r35 {
																									frame.Consume()
																									r36 = frame.Peek()
																									if frame.Flow == 0 {
																										r37 = 'm'
																										r38 = r36 == r37
																										if r38 {
																											frame.Consume()
																											r39 = frame.Peek()
																											if frame.Flow == 0 {
																												r40 = 'e'
																												r41 = r39 == r40
																												if r41 {
																													frame.Consume()
																													r42 = frame.Peek()
																													if frame.Flow == 0 {
																														r43 = 'n'
																														r44 = r42 == r43
																														if r44 {
																															frame.Consume()
																															r45 = frame.Peek()
																															if frame.Flow == 0 {
																																r46 = 't'
																																r47 = r45 == r46
																																if r47 {
																																	frame.Consume()
																																	r48 = frame.Peek()
																																	if frame.Flow == 0 {
																																		r49 = 's'
																																		r50 = r48 == r49
																																		if r50 {
																																			frame.Consume()
																																			EndKeyword(frame)
																																			if frame.Flow == 0 {
																																				r51 = ParseTypeRef(frame)
																																				if frame.Flow == 0 {
																																					r52 = r51
																																					goto block2
																																				} else {
																																					goto block1
																																				}
																																			} else {
																																				goto block1
																																			}
																																		} else {
																																			frame.Fail()
																																			goto block1
																																		}
																																	} else {
																																		goto block1
																																	}
																																} else {
																																	frame.Fail()
																																	goto block1
																																}
																															} else {
																																goto block1
																															}
																														} else {
																															frame.Fail()
																															goto block1
																														}
																													} else {
																														goto block1
																													}
																												} else {
																													frame.Fail()
																													goto block1
																												}
																											} else {
																												goto block1
																											}
																										} else {
																											frame.Fail()
																											goto block1
																										}
																									} else {
																										goto block1
																									}
																								} else {
																									frame.Fail()
																									goto block1
																								}
																							} else {
																								goto block1
																							}
																						} else {
																							frame.Fail()
																							goto block1
																						}
																					} else {
																						goto block1
																					}
																				} else {
																					frame.Fail()
																					goto block1
																				}
																			} else {
																				goto block1
																			}
																		} else {
																			frame.Fail()
																			goto block1
																		}
																	} else {
																		goto block1
																	}
																} else {
																	frame.Fail()
																	goto block1
																}
															} else {
																goto block1
															}
														} else {
															goto block5
														}
													} else {
														goto block5
													}
												} else {
													frame.Fail()
													goto block5
												}
											} else {
												goto block5
											}
										} else {
											frame.Fail()
											goto block5
										}
									} else {
										goto block5
									}
								} else {
									frame.Fail()
									goto block5
								}
							} else {
								goto block5
							}
						} else {
							frame.Fail()
							goto block5
						}
					} else {
						goto block5
					}
				} else {
					frame.Fail()
					goto block5
				}
			} else {
				goto block5
			}
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block1:
	frame.Recover(r20)
	r52 = r19
	goto block2
block2:
	r53 = frame.Peek()
	if frame.Flow == 0 {
		r54 = '{'
		r55 = r53 == r54
		if r55 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r56 = []*FieldDecl{}
				r57 = r56
				goto block3
			} else {
				goto block5
			}
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block3:
	r58 = frame.Checkpoint()
	r59 = Ident(frame)
	if frame.Flow == 0 {
		r60 = ParseTypeRef(frame)
		if frame.Flow == 0 {
			r61 = &FieldDecl{Name: r59, Type: r60}
			r62 = append(r57, r61)
			r57 = r62
			goto block3
		} else {
			goto block4
		}
	} else {
		goto block4
	}
block4:
	frame.Recover(r58)
	r63 = frame.Peek()
	if frame.Flow == 0 {
		r64 = '}'
		r65 = r63 == r64
		if r65 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r66 = &StructDecl{Name: r18, Implements: r52, Fields: r57}
				ret0 = r66
				return
			} else {
				goto block5
			}
		} else {
			frame.Fail()
			goto block5
		}
	} else {
		goto block5
	}
block5:
	return
}
func ParseFuncDecl(frame *runtime.State) (ret0 *FuncDecl) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 *Id
	var r13 []ASTTypeRef
	var r14 []ASTExpr
	var r15 *FuncDecl
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = 'f'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				r4 = 'u'
				r5 = r3 == r4
				if r5 {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						r7 = 'n'
						r8 = r6 == r7
						if r8 {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								r10 = 'c'
								r11 = r9 == r10
								if r11 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										r12 = Ident(frame)
										if frame.Flow == 0 {
											r13 = ParseTypeList(frame)
											if frame.Flow == 0 {
												r14 = ParseCodeBlock(frame)
												if frame.Flow == 0 {
													r15 = &FuncDecl{Name: r12, ReturnTypes: r13, Block: r14}
													ret0 = r15
													return
												} else {
													goto block1
												}
											} else {
												goto block1
											}
										} else {
											goto block1
										}
									} else {
										goto block1
									}
								} else {
									frame.Fail()
									goto block1
								}
							} else {
								goto block1
							}
						} else {
							frame.Fail()
							goto block1
						}
					} else {
						goto block1
					}
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseTest(frame *runtime.State) (ret0 *Test) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 *Id
	var r13 *Id
	var r14 string
	var r15 Destructure
	var r16 *Test
	r0 = frame.Peek()
	if frame.Flow == 0 {
		r1 = 't'
		r2 = r0 == r1
		if r2 {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				r4 = 'e'
				r5 = r3 == r4
				if r5 {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						r7 = 's'
						r8 = r6 == r7
						if r8 {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								r10 = 't'
								r11 = r9 == r10
								if r11 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										r12 = Ident(frame)
										if frame.Flow == 0 {
											r13 = Ident(frame)
											if frame.Flow == 0 {
												r14 = DecodeString(frame)
												if frame.Flow == 0 {
													S(frame)
													if frame.Flow == 0 {
														r15 = ParseDestructure(frame)
														if frame.Flow == 0 {
															r16 = &Test{Rule: r12, Name: r13, Input: r14, Destructure: r15}
															ret0 = r16
															return
														} else {
															goto block1
														}
													} else {
														goto block1
													}
												} else {
													goto block1
												}
											} else {
												goto block1
											}
										} else {
											goto block1
										}
									} else {
										goto block1
									}
								} else {
									frame.Fail()
									goto block1
								}
							} else {
								goto block1
							}
						} else {
							frame.Fail()
							goto block1
						}
					} else {
						goto block1
					}
				} else {
					frame.Fail()
					goto block1
				}
			} else {
				goto block1
			}
		} else {
			frame.Fail()
			goto block1
		}
	} else {
		goto block1
	}
block1:
	return
}
func ParseFile(frame *runtime.State) (ret0 *File) {
	var r0 []ASTDecl
	var r1 []*Test
	var r2 []ASTDecl
	var r3 []*Test
	var r4 int
	var r5 int
	var r6 *FuncDecl
	var r7 []ASTDecl
	var r8 *StructDecl
	var r9 []ASTDecl
	var r10 *Test
	var r11 []*Test
	var r12 int
	var r13 *File
	r0 = []ASTDecl{}
	r1 = []*Test{}
	r2, r3 = r0, r1
	goto block1
block1:
	r4 = frame.Checkpoint()
	r5 = frame.Checkpoint()
	r6 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		r7 = append(r2, r6)
		r2, r3 = r7, r3
		goto block1
	} else {
		frame.Recover(r5)
		r8 = ParseStructDecl(frame)
		if frame.Flow == 0 {
			r9 = append(r2, r8)
			r2, r3 = r9, r3
			goto block1
		} else {
			frame.Recover(r5)
			r10 = ParseTest(frame)
			if frame.Flow == 0 {
				r11 = append(r3, r10)
				r2, r3 = r2, r11
				goto block1
			} else {
				frame.Recover(r4)
				r12 = frame.LookaheadBegin()
				frame.Peek()
				if frame.Flow == 0 {
					frame.Consume()
					frame.LookaheadFail(r12)
					return
				} else {
					frame.LookaheadNormal(r12)
					r13 = &File{Decls: r2, Tests: r3}
					ret0 = r13
					return
				}
			}
		}
	}
}
