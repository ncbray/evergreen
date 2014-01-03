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
	Name string
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
	Name        string
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
	Name   string
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) isASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) isASTExpr() {
}

type NamedExpr struct {
	Name string
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
	Name string
	T    ASTType
}

func (node *Call) isASTExpr() {
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
	Name string
	Type ASTTypeRef
}
type StructDecl struct {
	Name       string
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
	Name        string
	ReturnTypes []ASTTypeRef
	Block       []ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) isASTDecl() {
}
func (node *FuncDecl) isASTFunc() {
}

type Test struct {
	Rule        string
	Name        string
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
func Ident(frame *runtime.State) (ret0 string) {
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
	var r12 int
	var r13 rune
	var r14 rune
	var r15 bool
	var r16 rune
	var r17 bool
	var r18 rune
	var r19 bool
	var r20 rune
	var r21 bool
	var r22 rune
	var r23 bool
	var r24 rune
	var r25 bool
	var r26 rune
	var r27 bool
	var r28 string
	r0 = frame.Checkpoint()
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
		goto block10
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
		frame.Fail()
		goto block10
	}
block3:
	frame.Consume()
	goto block4
block4:
	r12 = frame.Checkpoint()
	r13 = frame.Peek()
	if frame.Flow == 0 {
		r14 = 'a'
		r15 = r13 >= r14
		if r15 {
			r16 = 'z'
			r17 = r13 <= r16
			if r17 {
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
	r18 = 'A'
	r19 = r13 >= r18
	if r19 {
		r20 = 'Z'
		r21 = r13 <= r20
		if r21 {
			goto block7
		} else {
			goto block6
		}
	} else {
		goto block6
	}
block6:
	r22 = '_'
	r23 = r13 == r22
	if r23 {
		goto block7
	} else {
		r24 = '0'
		r25 = r13 >= r24
		if r25 {
			r26 = '9'
			r27 = r13 <= r26
			if r27 {
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
	frame.Recover(r12)
	r28 = frame.Slice(r0)
	S(frame)
	if frame.Flow == 0 {
		ret0 = r28
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
	var r0 string
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
	var r8 string
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
	var r0 string
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
	var r2 int
	var r3 int
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 rune
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
	var r24 rune
	var r25 bool
	var r26 rune
	var r27 rune
	var r28 bool
	var r29 int
	var r30 []ASTExpr
	var r31 *Repeat
	var r32 rune
	var r33 rune
	var r34 bool
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
	var r50 []ASTExpr
	var r51 [][]ASTExpr
	var r52 [][]ASTExpr
	var r53 int
	var r54 rune
	var r55 rune
	var r56 bool
	var r57 rune
	var r58 rune
	var r59 bool
	var r60 []ASTExpr
	var r61 [][]ASTExpr
	var r62 *Choice
	var r63 rune
	var r64 rune
	var r65 bool
	var r66 rune
	var r67 rune
	var r68 bool
	var r69 rune
	var r70 rune
	var r71 bool
	var r72 rune
	var r73 rune
	var r74 bool
	var r75 rune
	var r76 rune
	var r77 bool
	var r78 rune
	var r79 rune
	var r80 bool
	var r81 rune
	var r82 rune
	var r83 bool
	var r84 rune
	var r85 rune
	var r86 bool
	var r87 []ASTExpr
	var r88 *Optional
	var r89 rune
	var r90 rune
	var r91 bool
	var r92 rune
	var r93 rune
	var r94 bool
	var r95 rune
	var r96 rune
	var r97 bool
	var r98 rune
	var r99 rune
	var r100 bool
	var r101 rune
	var r102 rune
	var r103 bool
	var r104 []ASTExpr
	var r105 *Slice
	var r106 rune
	var r107 rune
	var r108 bool
	var r109 rune
	var r110 rune
	var r111 bool
	var r112 ASTExpr
	var r113 []ASTExpr
	var r114 *If
	var r115 rune
	var r116 rune
	var r117 bool
	var r118 rune
	var r119 rune
	var r120 bool
	var r121 rune
	var r122 rune
	var r123 bool
	var r124 string
	var r125 ASTTypeRef
	var r126 ASTExpr
	var r127 int
	var r128 rune
	var r129 rune
	var r130 bool
	var r131 ASTExpr
	var r132 ASTExpr
	var r133 bool
	var r134 *Assign
	var r135 rune
	var r136 rune
	var r137 bool
	var r138 rune
	var r139 rune
	var r140 bool
	var r141 rune
	var r142 rune
	var r143 bool
	var r144 rune
	var r145 rune
	var r146 bool
	var r147 *Fail
	var r148 rune
	var r149 rune
	var r150 bool
	var r151 rune
	var r152 rune
	var r153 bool
	var r154 rune
	var r155 rune
	var r156 bool
	var r157 rune
	var r158 rune
	var r159 bool
	var r160 rune
	var r161 rune
	var r162 bool
	var r163 rune
	var r164 rune
	var r165 bool
	var r166 ASTTypeRef
	var r167 ASTExpr
	var r168 *Coerce
	var r169 rune
	var r170 rune
	var r171 bool
	var r172 rune
	var r173 rune
	var r174 bool
	var r175 rune
	var r176 rune
	var r177 bool
	var r178 rune
	var r179 rune
	var r180 bool
	var r181 rune
	var r182 rune
	var r183 bool
	var r184 rune
	var r185 rune
	var r186 bool
	var r187 string
	var r188 ASTExpr
	var r189 *GetName
	var r190 *Append
	var r191 *Assign
	var r192 rune
	var r193 rune
	var r194 bool
	var r195 rune
	var r196 rune
	var r197 bool
	var r198 rune
	var r199 rune
	var r200 bool
	var r201 rune
	var r202 rune
	var r203 bool
	var r204 rune
	var r205 rune
	var r206 bool
	var r207 rune
	var r208 rune
	var r209 bool
	var r210 int
	var r211 rune
	var r212 rune
	var r213 bool
	var r214 []ASTExpr
	var r215 rune
	var r216 rune
	var r217 bool
	var r218 *Return
	var r219 ASTExpr
	var r220 []ASTExpr
	var r221 *Return
	var r222 []ASTExpr
	var r223 *Return
	var r224 string
	var r225 rune
	var r226 rune
	var r227 bool
	var r228 rune
	var r229 rune
	var r230 bool
	var r231 *Call
	var r232 string
	var r233 ASTExpr
	var r234 ASTExpr
	var r235 *BinaryOp
	var r236 *TypeRef
	var r237 rune
	var r238 rune
	var r239 bool
	var r240 []*NamedExpr
	var r241 rune
	var r242 rune
	var r243 bool
	var r244 *Construct
	var r245 *ListTypeRef
	var r246 rune
	var r247 rune
	var r248 bool
	var r249 []ASTExpr
	var r250 rune
	var r251 rune
	var r252 bool
	var r253 *ConstructList
	var r254 *StringMatch
	var r255 *RuneMatch
	var r256 string
	var r257 int
	var r258 bool
	var r259 int
	var r260 rune
	var r261 rune
	var r262 bool
	var r263 rune
	var r264 rune
	var r265 bool
	var r266 bool
	var r267 bool
	var r268 rune
	var r269 rune
	var r270 bool
	var r271 ASTExpr
	var r272 *Assign
	var r273 *GetName
	r0 = frame.Checkpoint()
	r1 = Literal(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block25
	} else {
		frame.Recover(r0)
		r2 = 0
		r3 = frame.Checkpoint()
		r4 = frame.Peek()
		if frame.Flow == 0 {
			r5 = 's'
			r6 = r4 == r5
			if r6 {
				frame.Consume()
				r7 = frame.Peek()
				if frame.Flow == 0 {
					r8 = 't'
					r9 = r7 == r8
					if r9 {
						frame.Consume()
						r10 = frame.Peek()
						if frame.Flow == 0 {
							r11 = 'a'
							r12 = r10 == r11
							if r12 {
								frame.Consume()
								r13 = frame.Peek()
								if frame.Flow == 0 {
									r14 = 'r'
									r15 = r13 == r14
									if r15 {
										frame.Consume()
										r16 = r2
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
	}
block1:
	frame.Recover(r3)
	r17 = frame.Peek()
	if frame.Flow == 0 {
		r18 = 'p'
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
						r24 = 'u'
						r25 = r23 == r24
						if r25 {
							frame.Consume()
							r26 = frame.Peek()
							if frame.Flow == 0 {
								r27 = 's'
								r28 = r26 == r27
								if r28 {
									frame.Consume()
									r29 = 1
									r16 = r29
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
block2:
	EndKeyword(frame)
	if frame.Flow == 0 {
		r30 = ParseCodeBlock(frame)
		if frame.Flow == 0 {
			r31 = &Repeat{Block: r30, Min: r16}
			ret0 = r31
			goto block25
		} else {
			goto block3
		}
	} else {
		goto block3
	}
block3:
	frame.Recover(r0)
	r32 = frame.Peek()
	if frame.Flow == 0 {
		r33 = 'c'
		r34 = r32 == r33
		if r34 {
			frame.Consume()
			r35 = frame.Peek()
			if frame.Flow == 0 {
				r36 = 'h'
				r37 = r35 == r36
				if r37 {
					frame.Consume()
					r38 = frame.Peek()
					if frame.Flow == 0 {
						r39 = 'o'
						r40 = r38 == r39
						if r40 {
							frame.Consume()
							r41 = frame.Peek()
							if frame.Flow == 0 {
								r42 = 'o'
								r43 = r41 == r42
								if r43 {
									frame.Consume()
									r44 = frame.Peek()
									if frame.Flow == 0 {
										r45 = 's'
										r46 = r44 == r45
										if r46 {
											frame.Consume()
											r47 = frame.Peek()
											if frame.Flow == 0 {
												r48 = 'e'
												r49 = r47 == r48
												if r49 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r50 = ParseCodeBlock(frame)
														if frame.Flow == 0 {
															r51 = [][]ASTExpr{r50}
															r52 = r51
															goto block4
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
block4:
	r53 = frame.Checkpoint()
	r54 = frame.Peek()
	if frame.Flow == 0 {
		r55 = 'o'
		r56 = r54 == r55
		if r56 {
			frame.Consume()
			r57 = frame.Peek()
			if frame.Flow == 0 {
				r58 = 'r'
				r59 = r57 == r58
				if r59 {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						r60 = ParseCodeBlock(frame)
						if frame.Flow == 0 {
							r61 = append(r52, r60)
							r52 = r61
							goto block4
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
block5:
	frame.Recover(r53)
	r62 = &Choice{Blocks: r52}
	ret0 = r62
	goto block25
block6:
	frame.Recover(r0)
	r63 = frame.Peek()
	if frame.Flow == 0 {
		r64 = 'q'
		r65 = r63 == r64
		if r65 {
			frame.Consume()
			r66 = frame.Peek()
			if frame.Flow == 0 {
				r67 = 'u'
				r68 = r66 == r67
				if r68 {
					frame.Consume()
					r69 = frame.Peek()
					if frame.Flow == 0 {
						r70 = 'e'
						r71 = r69 == r70
						if r71 {
							frame.Consume()
							r72 = frame.Peek()
							if frame.Flow == 0 {
								r73 = 's'
								r74 = r72 == r73
								if r74 {
									frame.Consume()
									r75 = frame.Peek()
									if frame.Flow == 0 {
										r76 = 't'
										r77 = r75 == r76
										if r77 {
											frame.Consume()
											r78 = frame.Peek()
											if frame.Flow == 0 {
												r79 = 'i'
												r80 = r78 == r79
												if r80 {
													frame.Consume()
													r81 = frame.Peek()
													if frame.Flow == 0 {
														r82 = 'o'
														r83 = r81 == r82
														if r83 {
															frame.Consume()
															r84 = frame.Peek()
															if frame.Flow == 0 {
																r85 = 'n'
																r86 = r84 == r85
																if r86 {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		r87 = ParseCodeBlock(frame)
																		if frame.Flow == 0 {
																			r88 = &Optional{Block: r87}
																			ret0 = r88
																			goto block25
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
		} else {
			frame.Fail()
			goto block7
		}
	} else {
		goto block7
	}
block7:
	frame.Recover(r0)
	r89 = frame.Peek()
	if frame.Flow == 0 {
		r90 = 's'
		r91 = r89 == r90
		if r91 {
			frame.Consume()
			r92 = frame.Peek()
			if frame.Flow == 0 {
				r93 = 'l'
				r94 = r92 == r93
				if r94 {
					frame.Consume()
					r95 = frame.Peek()
					if frame.Flow == 0 {
						r96 = 'i'
						r97 = r95 == r96
						if r97 {
							frame.Consume()
							r98 = frame.Peek()
							if frame.Flow == 0 {
								r99 = 'c'
								r100 = r98 == r99
								if r100 {
									frame.Consume()
									r101 = frame.Peek()
									if frame.Flow == 0 {
										r102 = 'e'
										r103 = r101 == r102
										if r103 {
											frame.Consume()
											EndKeyword(frame)
											if frame.Flow == 0 {
												r104 = ParseCodeBlock(frame)
												if frame.Flow == 0 {
													r105 = &Slice{Block: r104}
													ret0 = r105
													goto block25
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
			frame.Fail()
			goto block8
		}
	} else {
		goto block8
	}
block8:
	frame.Recover(r0)
	r106 = frame.Peek()
	if frame.Flow == 0 {
		r107 = 'i'
		r108 = r106 == r107
		if r108 {
			frame.Consume()
			r109 = frame.Peek()
			if frame.Flow == 0 {
				r110 = 'f'
				r111 = r109 == r110
				if r111 {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						r112 = ParseExpr(frame)
						if frame.Flow == 0 {
							r113 = ParseCodeBlock(frame)
							if frame.Flow == 0 {
								r114 = &If{Expr: r112, Block: r113}
								ret0 = r114
								goto block25
							} else {
								goto block9
							}
						} else {
							goto block9
						}
					} else {
						goto block9
					}
				} else {
					frame.Fail()
					goto block9
				}
			} else {
				goto block9
			}
		} else {
			frame.Fail()
			goto block9
		}
	} else {
		goto block9
	}
block9:
	frame.Recover(r0)
	r115 = frame.Peek()
	if frame.Flow == 0 {
		r116 = 'v'
		r117 = r115 == r116
		if r117 {
			frame.Consume()
			r118 = frame.Peek()
			if frame.Flow == 0 {
				r119 = 'a'
				r120 = r118 == r119
				if r120 {
					frame.Consume()
					r121 = frame.Peek()
					if frame.Flow == 0 {
						r122 = 'r'
						r123 = r121 == r122
						if r123 {
							frame.Consume()
							EndKeyword(frame)
							if frame.Flow == 0 {
								r124 = Ident(frame)
								if frame.Flow == 0 {
									r125 = ParseTypeRef(frame)
									if frame.Flow == 0 {
										r126 = nil
										r127 = frame.Checkpoint()
										r128 = frame.Peek()
										if frame.Flow == 0 {
											r129 = '='
											r130 = r128 == r129
											if r130 {
												frame.Consume()
												S(frame)
												if frame.Flow == 0 {
													r131 = ParseExpr(frame)
													if frame.Flow == 0 {
														r132 = r131
														goto block11
													} else {
														goto block10
													}
												} else {
													goto block10
												}
											} else {
												frame.Fail()
												goto block10
											}
										} else {
											goto block10
										}
									} else {
										goto block12
									}
								} else {
									goto block12
								}
							} else {
								goto block12
							}
						} else {
							frame.Fail()
							goto block12
						}
					} else {
						goto block12
					}
				} else {
					frame.Fail()
					goto block12
				}
			} else {
				goto block12
			}
		} else {
			frame.Fail()
			goto block12
		}
	} else {
		goto block12
	}
block10:
	frame.Recover(r127)
	r132 = r126
	goto block11
block11:
	r133 = true
	r134 = &Assign{Expr: r132, Name: r124, Type: r125, Define: r133}
	ret0 = r134
	goto block25
block12:
	frame.Recover(r0)
	r135 = frame.Peek()
	if frame.Flow == 0 {
		r136 = 'f'
		r137 = r135 == r136
		if r137 {
			frame.Consume()
			r138 = frame.Peek()
			if frame.Flow == 0 {
				r139 = 'a'
				r140 = r138 == r139
				if r140 {
					frame.Consume()
					r141 = frame.Peek()
					if frame.Flow == 0 {
						r142 = 'i'
						r143 = r141 == r142
						if r143 {
							frame.Consume()
							r144 = frame.Peek()
							if frame.Flow == 0 {
								r145 = 'l'
								r146 = r144 == r145
								if r146 {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										r147 = &Fail{}
										ret0 = r147
										goto block25
									} else {
										goto block13
									}
								} else {
									frame.Fail()
									goto block13
								}
							} else {
								goto block13
							}
						} else {
							frame.Fail()
							goto block13
						}
					} else {
						goto block13
					}
				} else {
					frame.Fail()
					goto block13
				}
			} else {
				goto block13
			}
		} else {
			frame.Fail()
			goto block13
		}
	} else {
		goto block13
	}
block13:
	frame.Recover(r0)
	r148 = frame.Peek()
	if frame.Flow == 0 {
		r149 = 'c'
		r150 = r148 == r149
		if r150 {
			frame.Consume()
			r151 = frame.Peek()
			if frame.Flow == 0 {
				r152 = 'o'
				r153 = r151 == r152
				if r153 {
					frame.Consume()
					r154 = frame.Peek()
					if frame.Flow == 0 {
						r155 = 'e'
						r156 = r154 == r155
						if r156 {
							frame.Consume()
							r157 = frame.Peek()
							if frame.Flow == 0 {
								r158 = 'r'
								r159 = r157 == r158
								if r159 {
									frame.Consume()
									r160 = frame.Peek()
									if frame.Flow == 0 {
										r161 = 'c'
										r162 = r160 == r161
										if r162 {
											frame.Consume()
											r163 = frame.Peek()
											if frame.Flow == 0 {
												r164 = 'e'
												r165 = r163 == r164
												if r165 {
													frame.Consume()
													S(frame)
													if frame.Flow == 0 {
														r166 = ParseTypeRef(frame)
														if frame.Flow == 0 {
															r167 = ParseExpr(frame)
															if frame.Flow == 0 {
																r168 = &Coerce{Type: r166, Expr: r167}
																ret0 = r168
																goto block25
															} else {
																goto block14
															}
														} else {
															goto block14
														}
													} else {
														goto block14
													}
												} else {
													frame.Fail()
													goto block14
												}
											} else {
												goto block14
											}
										} else {
											frame.Fail()
											goto block14
										}
									} else {
										goto block14
									}
								} else {
									frame.Fail()
									goto block14
								}
							} else {
								goto block14
							}
						} else {
							frame.Fail()
							goto block14
						}
					} else {
						goto block14
					}
				} else {
					frame.Fail()
					goto block14
				}
			} else {
				goto block14
			}
		} else {
			frame.Fail()
			goto block14
		}
	} else {
		goto block14
	}
block14:
	frame.Recover(r0)
	r169 = frame.Peek()
	if frame.Flow == 0 {
		r170 = 'a'
		r171 = r169 == r170
		if r171 {
			frame.Consume()
			r172 = frame.Peek()
			if frame.Flow == 0 {
				r173 = 'p'
				r174 = r172 == r173
				if r174 {
					frame.Consume()
					r175 = frame.Peek()
					if frame.Flow == 0 {
						r176 = 'p'
						r177 = r175 == r176
						if r177 {
							frame.Consume()
							r178 = frame.Peek()
							if frame.Flow == 0 {
								r179 = 'e'
								r180 = r178 == r179
								if r180 {
									frame.Consume()
									r181 = frame.Peek()
									if frame.Flow == 0 {
										r182 = 'n'
										r183 = r181 == r182
										if r183 {
											frame.Consume()
											r184 = frame.Peek()
											if frame.Flow == 0 {
												r185 = 'd'
												r186 = r184 == r185
												if r186 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r187 = Ident(frame)
														if frame.Flow == 0 {
															r188 = ParseExpr(frame)
															if frame.Flow == 0 {
																r189 = &GetName{Name: r187}
																r190 = &Append{List: r189, Expr: r188}
																r191 = &Assign{Expr: r190, Name: r187}
																ret0 = r191
																goto block25
															} else {
																goto block15
															}
														} else {
															goto block15
														}
													} else {
														goto block15
													}
												} else {
													frame.Fail()
													goto block15
												}
											} else {
												goto block15
											}
										} else {
											frame.Fail()
											goto block15
										}
									} else {
										goto block15
									}
								} else {
									frame.Fail()
									goto block15
								}
							} else {
								goto block15
							}
						} else {
							frame.Fail()
							goto block15
						}
					} else {
						goto block15
					}
				} else {
					frame.Fail()
					goto block15
				}
			} else {
				goto block15
			}
		} else {
			frame.Fail()
			goto block15
		}
	} else {
		goto block15
	}
block15:
	frame.Recover(r0)
	r192 = frame.Peek()
	if frame.Flow == 0 {
		r193 = 'r'
		r194 = r192 == r193
		if r194 {
			frame.Consume()
			r195 = frame.Peek()
			if frame.Flow == 0 {
				r196 = 'e'
				r197 = r195 == r196
				if r197 {
					frame.Consume()
					r198 = frame.Peek()
					if frame.Flow == 0 {
						r199 = 't'
						r200 = r198 == r199
						if r200 {
							frame.Consume()
							r201 = frame.Peek()
							if frame.Flow == 0 {
								r202 = 'u'
								r203 = r201 == r202
								if r203 {
									frame.Consume()
									r204 = frame.Peek()
									if frame.Flow == 0 {
										r205 = 'r'
										r206 = r204 == r205
										if r206 {
											frame.Consume()
											r207 = frame.Peek()
											if frame.Flow == 0 {
												r208 = 'n'
												r209 = r207 == r208
												if r209 {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														r210 = frame.Checkpoint()
														r211 = frame.Peek()
														if frame.Flow == 0 {
															r212 = '('
															r213 = r211 == r212
															if r213 {
																frame.Consume()
																S(frame)
																if frame.Flow == 0 {
																	r214 = ParseExprList(frame)
																	if frame.Flow == 0 {
																		r215 = frame.Peek()
																		if frame.Flow == 0 {
																			r216 = ')'
																			r217 = r215 == r216
																			if r217 {
																				frame.Consume()
																				S(frame)
																				if frame.Flow == 0 {
																					r218 = &Return{Exprs: r214}
																					ret0 = r218
																					goto block25
																				} else {
																					goto block16
																				}
																			} else {
																				frame.Fail()
																				goto block16
																			}
																		} else {
																			goto block16
																		}
																	} else {
																		goto block16
																	}
																} else {
																	goto block16
																}
															} else {
																frame.Fail()
																goto block16
															}
														} else {
															goto block16
														}
													} else {
														goto block17
													}
												} else {
													frame.Fail()
													goto block17
												}
											} else {
												goto block17
											}
										} else {
											frame.Fail()
											goto block17
										}
									} else {
										goto block17
									}
								} else {
									frame.Fail()
									goto block17
								}
							} else {
								goto block17
							}
						} else {
							frame.Fail()
							goto block17
						}
					} else {
						goto block17
					}
				} else {
					frame.Fail()
					goto block17
				}
			} else {
				goto block17
			}
		} else {
			frame.Fail()
			goto block17
		}
	} else {
		goto block17
	}
block16:
	frame.Recover(r210)
	r219 = ParseExpr(frame)
	if frame.Flow == 0 {
		r220 = []ASTExpr{r219}
		r221 = &Return{Exprs: r220}
		ret0 = r221
		goto block25
	} else {
		frame.Recover(r210)
		r222 = []ASTExpr{}
		r223 = &Return{Exprs: r222}
		ret0 = r223
		goto block25
	}
block17:
	frame.Recover(r0)
	r224 = Ident(frame)
	if frame.Flow == 0 {
		r225 = frame.Peek()
		if frame.Flow == 0 {
			r226 = '('
			r227 = r225 == r226
			if r227 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r228 = frame.Peek()
					if frame.Flow == 0 {
						r229 = ')'
						r230 = r228 == r229
						if r230 {
							frame.Consume()
							S(frame)
							if frame.Flow == 0 {
								r231 = &Call{Name: r224}
								ret0 = r231
								goto block25
							} else {
								goto block18
							}
						} else {
							frame.Fail()
							goto block18
						}
					} else {
						goto block18
					}
				} else {
					goto block18
				}
			} else {
				frame.Fail()
				goto block18
			}
		} else {
			goto block18
		}
	} else {
		goto block18
	}
block18:
	frame.Recover(r0)
	r232 = BinaryOperator(frame)
	if frame.Flow == 0 {
		r233 = ParseExpr(frame)
		if frame.Flow == 0 {
			r234 = ParseExpr(frame)
			if frame.Flow == 0 {
				r235 = &BinaryOp{Left: r233, Op: r232, Right: r234}
				ret0 = r235
				goto block25
			} else {
				goto block19
			}
		} else {
			goto block19
		}
	} else {
		goto block19
	}
block19:
	frame.Recover(r0)
	r236 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		r237 = frame.Peek()
		if frame.Flow == 0 {
			r238 = '{'
			r239 = r237 == r238
			if r239 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r240 = ParseNamedExprList(frame)
					if frame.Flow == 0 {
						r241 = frame.Peek()
						if frame.Flow == 0 {
							r242 = '}'
							r243 = r241 == r242
							if r243 {
								frame.Consume()
								S(frame)
								if frame.Flow == 0 {
									r244 = &Construct{Type: r236, Args: r240}
									ret0 = r244
									goto block25
								} else {
									goto block20
								}
							} else {
								frame.Fail()
								goto block20
							}
						} else {
							goto block20
						}
					} else {
						goto block20
					}
				} else {
					goto block20
				}
			} else {
				frame.Fail()
				goto block20
			}
		} else {
			goto block20
		}
	} else {
		goto block20
	}
block20:
	frame.Recover(r0)
	r245 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		r246 = frame.Peek()
		if frame.Flow == 0 {
			r247 = '{'
			r248 = r246 == r247
			if r248 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r249 = ParseExprList(frame)
					if frame.Flow == 0 {
						r250 = frame.Peek()
						if frame.Flow == 0 {
							r251 = '}'
							r252 = r250 == r251
							if r252 {
								frame.Consume()
								S(frame)
								if frame.Flow == 0 {
									r253 = &ConstructList{Type: r245, Args: r249}
									ret0 = r253
									goto block25
								} else {
									goto block21
								}
							} else {
								frame.Fail()
								goto block21
							}
						} else {
							goto block21
						}
					} else {
						goto block21
					}
				} else {
					goto block21
				}
			} else {
				frame.Fail()
				goto block21
			}
		} else {
			goto block21
		}
	} else {
		goto block21
	}
block21:
	frame.Recover(r0)
	r254 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		ret0 = r254
		goto block25
	} else {
		frame.Recover(r0)
		r255 = RuneMatchExpr(frame)
		if frame.Flow == 0 {
			ret0 = r255
			goto block25
		} else {
			frame.Recover(r0)
			r256 = Ident(frame)
			if frame.Flow == 0 {
				r257 = frame.Checkpoint()
				r258 = false
				r259 = frame.Checkpoint()
				r260 = frame.Peek()
				if frame.Flow == 0 {
					r261 = ':'
					r262 = r260 == r261
					if r262 {
						frame.Consume()
						r263 = frame.Peek()
						if frame.Flow == 0 {
							r264 = '='
							r265 = r263 == r264
							if r265 {
								frame.Consume()
								r266 = true
								r267 = r266
								goto block23
							} else {
								frame.Fail()
								goto block22
							}
						} else {
							goto block22
						}
					} else {
						frame.Fail()
						goto block22
					}
				} else {
					goto block22
				}
			} else {
				return
			}
		}
	}
block22:
	frame.Recover(r259)
	r268 = frame.Peek()
	if frame.Flow == 0 {
		r269 = '='
		r270 = r268 == r269
		if r270 {
			frame.Consume()
			r267 = r258
			goto block23
		} else {
			frame.Fail()
			goto block24
		}
	} else {
		goto block24
	}
block23:
	S(frame)
	if frame.Flow == 0 {
		r271 = ParseExpr(frame)
		if frame.Flow == 0 {
			r272 = &Assign{Expr: r271, Name: r256, Define: r267}
			ret0 = r272
			goto block25
		} else {
			goto block24
		}
	} else {
		goto block24
	}
block24:
	frame.Recover(r257)
	r273 = &GetName{Name: r256}
	ret0 = r273
	goto block25
block25:
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
	var r8 int
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 bool
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
	r5 = frame.Checkpoint()
	r6 = ParseExpr(frame)
	if frame.Flow == 0 {
		r7 = append(r4, r6)
		goto block2
	} else {
		frame.Recover(r5)
		r12 = frame.Peek()
		if frame.Flow == 0 {
			r13 = '}'
			r14 = r12 == r13
			if r14 {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					ret0 = r4
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
	}
block2:
	r8 = frame.Checkpoint()
	r9 = frame.Peek()
	if frame.Flow == 0 {
		r10 = ';'
		r11 = r9 == r10
		if r11 {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
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
block3:
	frame.Recover(r8)
	r4 = r7
	goto block1
block4:
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
	var r18 string
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
	var r59 string
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
	var r12 string
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
	var r12 string
	var r13 string
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
