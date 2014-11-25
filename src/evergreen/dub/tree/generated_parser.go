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

type NilLiteral struct {
}

func (node *NilLiteral) isASTExpr() {
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

type ASTTypeRef interface {
	isASTTypeRef()
}

type TypeRef struct {
	Name *Id
	T    DubType
}

func (node *TypeRef) isASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    DubType
}

func (node *ListTypeRef) isASTTypeRef() {
}

type QualifiedTypeRef struct {
	Package *Id
	Name    *Id
	T       DubType
}

func (node *QualifiedTypeRef) isASTTypeRef() {
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
	Type ASTTypeRef
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
	Expr    ASTExpr
	Targets []ASTExpr
	Type    ASTTypeRef
	Define  bool
}

func (node *Assign) isASTExpr() {
}

type NameRef struct {
	Name  *Id
	Local LocalInfo_Ref
}

func (node *NameRef) isASTExpr() {
}

type NamedExpr struct {
	Name *Id
	Expr ASTExpr
}

type Construct struct {
	Type ASTTypeRef
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
	Name   *Id
	Args   []ASTExpr
	Target ASTCallable
	T      []DubType
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
	T    DubType
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
	T     DubType
}

func (node *BinaryOp) isASTExpr() {
}

type FieldDecl struct {
	Name *Id
	Type ASTTypeRef
}

type StructDecl struct {
	Name       *Id
	Implements ASTTypeRef
	Fields     []*FieldDecl
	Scoped     bool
	Contains   []ASTTypeRef
	T          *StructType
}

func (node *StructDecl) isASTDecl() {
}

type ASTCallable interface {
	isASTCallable()
}

type LocalInfo_Ref uint32

const NoLocalInfo = ^LocalInfo_Ref(0)

type LocalInfo_Scope struct {
	objects []*LocalInfo
}

type LocalInfo struct {
	Name string
	T    DubType
}

type Param struct {
	Name *NameRef
	Type ASTTypeRef
}

type FuncDecl struct {
	Name            *Id
	Params          []*Param
	ReturnTypes     []ASTTypeRef
	Block           []ASTExpr
	LocalInfo_Scope *LocalInfo_Scope
}

func (node *FuncDecl) isASTDecl() {
}

func (node *FuncDecl) isASTCallable() {
}

type Test struct {
	Name        *Id
	Rule        ASTExpr
	Type        DubType
	Input       string
	Flow        string
	Destructure Destructure
}

type ImportDecl struct {
	Path *StringLiteral
}

type File struct {
	Imports []*ImportDecl
	Decls   []ASTDecl
	Tests   []*Test
}

type Package struct {
	Path  []string
	Files []*File
}

type BuiltinTypeIndex struct {
	String *BuiltinType
	Rune   *BuiltinType
	Int    *BuiltinType
	Int64  *BuiltinType
	Bool   *BuiltinType
	Graph  *BuiltinType
	Nil    *NilType
}

type Program struct {
	Builtins *BuiltinTypeIndex
	Packages []*Package
}

type DubType interface {
	isDubType()
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) isDubType() {
}

type NilType struct {
}

func (node *NilType) isDubType() {
}

type ListType struct {
	Type DubType
}

func (node *ListType) isDubType() {
}

type FieldType struct {
	Name *Id
	Type DubType
}

type StructType struct {
	Name       *Id
	Implements *StructType
	Fields     []*FieldType
	Scoped     bool
	Contains   []*StructType
	IsParent   bool
}

func (node *StructType) isDubType() {
}

func LineTerminator(frame *runtime.State) {
	var r0 int
	var r1 rune
	var r4 rune
	var r7 rune
	var r10 rune
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 == '\n' {
			frame.Consume()
			goto block3
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r0)
	r4 = frame.Peek()
	if frame.Flow == 0 {
		if r4 == '\r' {
			frame.Consume()
			r7 = frame.Peek()
			if frame.Flow == 0 {
				if r7 == '\n' {
					frame.Consume()
					goto block3
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(r0)
	r10 = frame.Peek()
	if frame.Flow == 0 {
		if r10 == '\r' {
			frame.Consume()
			goto block3
		}
		frame.Fail()
		goto block4
	}
	goto block4
block3:
	return
block4:
	return
}

func S(frame *runtime.State) {
	var r0 int
	var r1 int
	var r2 rune
	var r7 rune
	var r10 rune
	var r13 int
	var r14 rune
	goto block1
block1:
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == ' ' {
			goto block2
		}
		if r2 == '\t' {
			goto block2
		}
		frame.Fail()
		goto block3
	}
	goto block3
block2:
	frame.Consume()
	goto block1
block3:
	frame.Recover(r1)
	LineTerminator(frame)
	if frame.Flow == 0 {
		goto block1
	}
	frame.Recover(r1)
	r7 = frame.Peek()
	if frame.Flow == 0 {
		if r7 == '/' {
			frame.Consume()
			r10 = frame.Peek()
			if frame.Flow == 0 {
				if r10 == '/' {
					frame.Consume()
					goto block4
				}
				frame.Fail()
				goto block7
			}
			goto block7
		}
		frame.Fail()
		goto block7
	}
	goto block7
block4:
	r13 = frame.Checkpoint()
	r14 = frame.Peek()
	if frame.Flow == 0 {
		if r14 == '\n' {
			goto block5
		}
		if r14 == '\r' {
			goto block5
		}
		frame.Consume()
		goto block4
	}
	goto block6
block5:
	frame.Fail()
	goto block6
block6:
	frame.Recover(r13)
	goto block1
block7:
	frame.Recover(r0)
	return
}

func EndKeyword(frame *runtime.State) {
	var r0 int
	var r1 rune
	r0 = frame.LookaheadBegin()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 >= 'a' {
			if r1 <= 'z' {
				goto block3
			}
			goto block1
		}
		goto block1
	}
	goto block5
block1:
	if r1 >= 'A' {
		if r1 <= 'Z' {
			goto block3
		}
		goto block2
	}
	goto block2
block2:
	if r1 == '_' {
		goto block3
	}
	if r1 >= '0' {
		if r1 <= '9' {
			goto block3
		}
		goto block4
	}
	goto block4
block3:
	frame.Consume()
	frame.LookaheadFail(r0)
	return
block4:
	frame.Fail()
	goto block5
block5:
	frame.LookaheadNormal(r0)
	return
}

func Ident(frame *runtime.State) (ret0 *Id) {
	var r0 int
	var r1 int
	var r2 int
	var r3 rune
	var r6 rune
	var r9 rune
	var r12 rune
	var r15 rune
	var r18 rune
	var r21 rune
	var r24 rune
	var r27 rune
	var r30 rune
	var r33 rune
	var r36 rune
	var r39 rune
	var r42 rune
	var r45 rune
	var r48 rune
	var r51 rune
	var r54 rune
	var r57 rune
	var r60 rune
	var r63 rune
	var r66 rune
	var r69 rune
	var r72 rune
	var r75 rune
	var r78 rune
	var r81 rune
	var r84 rune
	var r87 rune
	var r90 rune
	var r93 rune
	var r96 rune
	var r99 rune
	var r102 rune
	var r105 rune
	var r108 rune
	var r111 rune
	var r114 rune
	var r117 rune
	var r120 rune
	var r123 rune
	var r126 rune
	var r129 rune
	var r132 rune
	var r135 rune
	var r138 rune
	var r141 rune
	var r144 rune
	var r147 rune
	var r150 rune
	var r153 rune
	var r156 rune
	var r159 rune
	var r162 rune
	var r165 rune
	var r168 rune
	var r171 int
	var r172 rune
	var r187 int
	var r188 rune
	var r199 int
	var r200 rune
	r0 = frame.Checkpoint()
	r1 = frame.LookaheadBegin()
	r2 = frame.Checkpoint()
	r3 = frame.Peek()
	if frame.Flow == 0 {
		if r3 == 'f' {
			frame.Consume()
			r6 = frame.Peek()
			if frame.Flow == 0 {
				if r6 == 'u' {
					frame.Consume()
					r9 = frame.Peek()
					if frame.Flow == 0 {
						if r9 == 'n' {
							frame.Consume()
							r12 = frame.Peek()
							if frame.Flow == 0 {
								if r12 == 'c' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r2)
	r15 = frame.Peek()
	if frame.Flow == 0 {
		if r15 == 't' {
			frame.Consume()
			r18 = frame.Peek()
			if frame.Flow == 0 {
				if r18 == 'e' {
					frame.Consume()
					r21 = frame.Peek()
					if frame.Flow == 0 {
						if r21 == 's' {
							frame.Consume()
							r24 = frame.Peek()
							if frame.Flow == 0 {
								if r24 == 't' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block2
							}
							goto block2
						}
						frame.Fail()
						goto block2
					}
					goto block2
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(r2)
	r27 = frame.Peek()
	if frame.Flow == 0 {
		if r27 == 's' {
			frame.Consume()
			r30 = frame.Peek()
			if frame.Flow == 0 {
				if r30 == 't' {
					frame.Consume()
					r33 = frame.Peek()
					if frame.Flow == 0 {
						if r33 == 'r' {
							frame.Consume()
							r36 = frame.Peek()
							if frame.Flow == 0 {
								if r36 == 'u' {
									frame.Consume()
									r39 = frame.Peek()
									if frame.Flow == 0 {
										if r39 == 'c' {
											frame.Consume()
											r42 = frame.Peek()
											if frame.Flow == 0 {
												if r42 == 't' {
													frame.Consume()
													goto block13
												}
												frame.Fail()
												goto block3
											}
											goto block3
										}
										frame.Fail()
										goto block3
									}
									goto block3
								}
								frame.Fail()
								goto block3
							}
							goto block3
						}
						frame.Fail()
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(r2)
	r45 = frame.Peek()
	if frame.Flow == 0 {
		if r45 == 's' {
			frame.Consume()
			r48 = frame.Peek()
			if frame.Flow == 0 {
				if r48 == 't' {
					frame.Consume()
					r51 = frame.Peek()
					if frame.Flow == 0 {
						if r51 == 'a' {
							frame.Consume()
							r54 = frame.Peek()
							if frame.Flow == 0 {
								if r54 == 'r' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block4
							}
							goto block4
						}
						frame.Fail()
						goto block4
					}
					goto block4
				}
				frame.Fail()
				goto block4
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(r2)
	r57 = frame.Peek()
	if frame.Flow == 0 {
		if r57 == 'p' {
			frame.Consume()
			r60 = frame.Peek()
			if frame.Flow == 0 {
				if r60 == 'l' {
					frame.Consume()
					r63 = frame.Peek()
					if frame.Flow == 0 {
						if r63 == 'u' {
							frame.Consume()
							r66 = frame.Peek()
							if frame.Flow == 0 {
								if r66 == 's' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block5
							}
							goto block5
						}
						frame.Fail()
						goto block5
					}
					goto block5
				}
				frame.Fail()
				goto block5
			}
			goto block5
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(r2)
	r69 = frame.Peek()
	if frame.Flow == 0 {
		if r69 == 'c' {
			frame.Consume()
			r72 = frame.Peek()
			if frame.Flow == 0 {
				if r72 == 'h' {
					frame.Consume()
					r75 = frame.Peek()
					if frame.Flow == 0 {
						if r75 == 'o' {
							frame.Consume()
							r78 = frame.Peek()
							if frame.Flow == 0 {
								if r78 == 'o' {
									frame.Consume()
									r81 = frame.Peek()
									if frame.Flow == 0 {
										if r81 == 's' {
											frame.Consume()
											r84 = frame.Peek()
											if frame.Flow == 0 {
												if r84 == 'e' {
													frame.Consume()
													goto block13
												}
												frame.Fail()
												goto block6
											}
											goto block6
										}
										frame.Fail()
										goto block6
									}
									goto block6
								}
								frame.Fail()
								goto block6
							}
							goto block6
						}
						frame.Fail()
						goto block6
					}
					goto block6
				}
				frame.Fail()
				goto block6
			}
			goto block6
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(r2)
	r87 = frame.Peek()
	if frame.Flow == 0 {
		if r87 == 'o' {
			frame.Consume()
			r90 = frame.Peek()
			if frame.Flow == 0 {
				if r90 == 'r' {
					frame.Consume()
					goto block13
				}
				frame.Fail()
				goto block7
			}
			goto block7
		}
		frame.Fail()
		goto block7
	}
	goto block7
block7:
	frame.Recover(r2)
	r93 = frame.Peek()
	if frame.Flow == 0 {
		if r93 == 'q' {
			frame.Consume()
			r96 = frame.Peek()
			if frame.Flow == 0 {
				if r96 == 'u' {
					frame.Consume()
					r99 = frame.Peek()
					if frame.Flow == 0 {
						if r99 == 'e' {
							frame.Consume()
							r102 = frame.Peek()
							if frame.Flow == 0 {
								if r102 == 's' {
									frame.Consume()
									r105 = frame.Peek()
									if frame.Flow == 0 {
										if r105 == 't' {
											frame.Consume()
											r108 = frame.Peek()
											if frame.Flow == 0 {
												if r108 == 'i' {
													frame.Consume()
													r111 = frame.Peek()
													if frame.Flow == 0 {
														if r111 == 'o' {
															frame.Consume()
															r114 = frame.Peek()
															if frame.Flow == 0 {
																if r114 == 'n' {
																	frame.Consume()
																	goto block13
																}
																frame.Fail()
																goto block8
															}
															goto block8
														}
														frame.Fail()
														goto block8
													}
													goto block8
												}
												frame.Fail()
												goto block8
											}
											goto block8
										}
										frame.Fail()
										goto block8
									}
									goto block8
								}
								frame.Fail()
								goto block8
							}
							goto block8
						}
						frame.Fail()
						goto block8
					}
					goto block8
				}
				frame.Fail()
				goto block8
			}
			goto block8
		}
		frame.Fail()
		goto block8
	}
	goto block8
block8:
	frame.Recover(r2)
	r117 = frame.Peek()
	if frame.Flow == 0 {
		if r117 == 'i' {
			frame.Consume()
			r120 = frame.Peek()
			if frame.Flow == 0 {
				if r120 == 'f' {
					frame.Consume()
					goto block13
				}
				frame.Fail()
				goto block9
			}
			goto block9
		}
		frame.Fail()
		goto block9
	}
	goto block9
block9:
	frame.Recover(r2)
	r123 = frame.Peek()
	if frame.Flow == 0 {
		if r123 == 'e' {
			frame.Consume()
			r126 = frame.Peek()
			if frame.Flow == 0 {
				if r126 == 'l' {
					frame.Consume()
					r129 = frame.Peek()
					if frame.Flow == 0 {
						if r129 == 's' {
							frame.Consume()
							r132 = frame.Peek()
							if frame.Flow == 0 {
								if r132 == 'e' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block10
							}
							goto block10
						}
						frame.Fail()
						goto block10
					}
					goto block10
				}
				frame.Fail()
				goto block10
			}
			goto block10
		}
		frame.Fail()
		goto block10
	}
	goto block10
block10:
	frame.Recover(r2)
	r135 = frame.Peek()
	if frame.Flow == 0 {
		if r135 == 't' {
			frame.Consume()
			r138 = frame.Peek()
			if frame.Flow == 0 {
				if r138 == 'r' {
					frame.Consume()
					r141 = frame.Peek()
					if frame.Flow == 0 {
						if r141 == 'u' {
							frame.Consume()
							r144 = frame.Peek()
							if frame.Flow == 0 {
								if r144 == 'e' {
									frame.Consume()
									goto block13
								}
								frame.Fail()
								goto block11
							}
							goto block11
						}
						frame.Fail()
						goto block11
					}
					goto block11
				}
				frame.Fail()
				goto block11
			}
			goto block11
		}
		frame.Fail()
		goto block11
	}
	goto block11
block11:
	frame.Recover(r2)
	r147 = frame.Peek()
	if frame.Flow == 0 {
		if r147 == 'f' {
			frame.Consume()
			r150 = frame.Peek()
			if frame.Flow == 0 {
				if r150 == 'a' {
					frame.Consume()
					r153 = frame.Peek()
					if frame.Flow == 0 {
						if r153 == 'l' {
							frame.Consume()
							r156 = frame.Peek()
							if frame.Flow == 0 {
								if r156 == 's' {
									frame.Consume()
									r159 = frame.Peek()
									if frame.Flow == 0 {
										if r159 == 'e' {
											frame.Consume()
											goto block13
										}
										frame.Fail()
										goto block12
									}
									goto block12
								}
								frame.Fail()
								goto block12
							}
							goto block12
						}
						frame.Fail()
						goto block12
					}
					goto block12
				}
				frame.Fail()
				goto block12
			}
			goto block12
		}
		frame.Fail()
		goto block12
	}
	goto block12
block12:
	frame.Recover(r2)
	r162 = frame.Peek()
	if frame.Flow == 0 {
		if r162 == 'n' {
			frame.Consume()
			r165 = frame.Peek()
			if frame.Flow == 0 {
				if r165 == 'i' {
					frame.Consume()
					r168 = frame.Peek()
					if frame.Flow == 0 {
						if r168 == 'l' {
							frame.Consume()
							goto block13
						}
						frame.Fail()
						goto block19
					}
					goto block19
				}
				frame.Fail()
				goto block19
			}
			goto block19
		}
		frame.Fail()
		goto block19
	}
	goto block19
block13:
	r171 = frame.LookaheadBegin()
	r172 = frame.Peek()
	if frame.Flow == 0 {
		if r172 >= 'a' {
			if r172 <= 'z' {
				goto block16
			}
			goto block14
		}
		goto block14
	}
	goto block18
block14:
	if r172 >= 'A' {
		if r172 <= 'Z' {
			goto block16
		}
		goto block15
	}
	goto block15
block15:
	if r172 == '_' {
		goto block16
	}
	if r172 >= '0' {
		if r172 <= '9' {
			goto block16
		}
		goto block17
	}
	goto block17
block16:
	frame.Consume()
	frame.LookaheadFail(r171)
	goto block19
block17:
	frame.Fail()
	goto block18
block18:
	frame.LookaheadNormal(r171)
	frame.LookaheadFail(r1)
	goto block29
block19:
	frame.LookaheadNormal(r1)
	r187 = frame.Checkpoint()
	r188 = frame.Peek()
	if frame.Flow == 0 {
		if r188 >= 'a' {
			if r188 <= 'z' {
				goto block22
			}
			goto block20
		}
		goto block20
	}
	goto block29
block20:
	if r188 >= 'A' {
		if r188 <= 'Z' {
			goto block22
		}
		goto block21
	}
	goto block21
block21:
	if r188 == '_' {
		goto block22
	}
	frame.Fail()
	goto block29
block22:
	frame.Consume()
	goto block23
block23:
	r199 = frame.Checkpoint()
	r200 = frame.Peek()
	if frame.Flow == 0 {
		if r200 >= 'a' {
			if r200 <= 'z' {
				goto block26
			}
			goto block24
		}
		goto block24
	}
	goto block28
block24:
	if r200 >= 'A' {
		if r200 <= 'Z' {
			goto block26
		}
		goto block25
	}
	goto block25
block25:
	if r200 == '_' {
		goto block26
	}
	if r200 >= '0' {
		if r200 <= '9' {
			goto block26
		}
		goto block27
	}
	goto block27
block26:
	frame.Consume()
	goto block23
block27:
	frame.Fail()
	goto block28
block28:
	frame.Recover(r199)
	ret0 = &Id{Pos: r0, Text: frame.Slice(r187)}
	return
block29:
	return
}

func DecodeInt(frame *runtime.State) (ret0 int, ret1 string) {
	var r0 int
	var r1 int
	var r2 rune
	var r10 int
	var r14 int
	var r15 int
	var r16 rune
	var r24 int
	var r28 string
	r0 = 0
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 >= '0' {
			if r2 <= '9' {
				frame.Consume()
				r10 = int(r2) - int('0')
				r14 = r0*10 + r10
				goto block1
			}
			goto block4
		}
		goto block4
	}
	goto block5
block1:
	r15 = frame.Checkpoint()
	r16 = frame.Peek()
	if frame.Flow == 0 {
		if r16 >= '0' {
			if r16 <= '9' {
				frame.Consume()
				r24 = int(r16) - int('0')
				r14 = r14*10 + r24
				goto block1
			}
			goto block2
		}
		goto block2
	}
	goto block3
block2:
	frame.Fail()
	goto block3
block3:
	frame.Recover(r15)
	r28 = frame.Slice(r1)
	ret0 = r14
	ret1 = r28
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
	var r5 rune
	var r9 rune
	var r13 rune
	var r17 rune
	var r21 rune
	var r25 rune
	var r29 rune
	var r33 rune
	var r37 rune
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 == 'a' {
			frame.Consume()
			ret0 = '\a'
			goto block10
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r0)
	r5 = frame.Peek()
	if frame.Flow == 0 {
		if r5 == 'b' {
			frame.Consume()
			ret0 = '\b'
			goto block10
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(r0)
	r9 = frame.Peek()
	if frame.Flow == 0 {
		if r9 == 'f' {
			frame.Consume()
			ret0 = '\f'
			goto block10
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(r0)
	r13 = frame.Peek()
	if frame.Flow == 0 {
		if r13 == 'n' {
			frame.Consume()
			ret0 = '\n'
			goto block10
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(r0)
	r17 = frame.Peek()
	if frame.Flow == 0 {
		if r17 == 'r' {
			frame.Consume()
			ret0 = '\r'
			goto block10
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(r0)
	r21 = frame.Peek()
	if frame.Flow == 0 {
		if r21 == 't' {
			frame.Consume()
			ret0 = '\t'
			goto block10
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(r0)
	r25 = frame.Peek()
	if frame.Flow == 0 {
		if r25 == 'v' {
			frame.Consume()
			ret0 = '\v'
			goto block10
		}
		frame.Fail()
		goto block7
	}
	goto block7
block7:
	frame.Recover(r0)
	r29 = frame.Peek()
	if frame.Flow == 0 {
		if r29 == '\\' {
			frame.Consume()
			ret0 = '\\'
			goto block10
		}
		frame.Fail()
		goto block8
	}
	goto block8
block8:
	frame.Recover(r0)
	r33 = frame.Peek()
	if frame.Flow == 0 {
		if r33 == '\'' {
			frame.Consume()
			ret0 = '\''
			goto block10
		}
		frame.Fail()
		goto block9
	}
	goto block9
block9:
	frame.Recover(r0)
	r37 = frame.Peek()
	if frame.Flow == 0 {
		if r37 == '"' {
			frame.Consume()
			ret0 = '"'
			goto block10
		}
		frame.Fail()
		goto block11
	}
	goto block11
block10:
	return
block11:
	return
}

func DecodeString(frame *runtime.State) (ret0 string) {
	var r0 rune
	var r4 []rune
	var r5 int
	var r6 int
	var r7 rune
	var r13 rune
	var r16 rune
	var r18 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '"' {
			frame.Consume()
			r4 = []rune{}
			goto block1
		}
		frame.Fail()
		goto block5
	}
	goto block5
block1:
	r5 = frame.Checkpoint()
	r6 = frame.Checkpoint()
	r7 = frame.Peek()
	if frame.Flow == 0 {
		if r7 == '"' {
			goto block2
		}
		if r7 == '\\' {
			goto block2
		}
		frame.Consume()
		r4 = append(r4, r7)
		goto block1
	}
	goto block3
block2:
	frame.Fail()
	goto block3
block3:
	frame.Recover(r6)
	r13 = frame.Peek()
	if frame.Flow == 0 {
		if r13 == '\\' {
			frame.Consume()
			r16 = EscapedChar(frame)
			if frame.Flow == 0 {
				r4 = append(r4, r16)
				goto block1
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(r5)
	r18 = frame.Peek()
	if frame.Flow == 0 {
		if r18 == '"' {
			frame.Consume()
			ret0 = string(r4)
			return
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	return
}

func DecodeRune(frame *runtime.State) (ret0 rune, ret1 string) {
	var r0 int
	var r1 rune
	var r4 int
	var r5 rune
	var r10 rune
	var r11 rune
	var r14 rune
	var r15 rune
	var r18 string
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 == '\'' {
			frame.Consume()
			r4 = frame.Checkpoint()
			r5 = frame.Peek()
			if frame.Flow == 0 {
				if r5 == '\\' {
					goto block1
				}
				if r5 == '\'' {
					goto block1
				}
				frame.Consume()
				r10 = r5
				goto block3
			}
			goto block2
		}
		frame.Fail()
		goto block4
	}
	goto block4
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(r4)
	r11 = frame.Peek()
	if frame.Flow == 0 {
		if r11 == '\\' {
			frame.Consume()
			r14 = EscapedChar(frame)
			if frame.Flow == 0 {
				r10 = r14
				goto block3
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block3:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		if r15 == '\'' {
			frame.Consume()
			r18 = frame.Slice(r0)
			ret0 = r10
			ret1 = r18
			return
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	return
}

func DecodeBool(frame *runtime.State) (ret0 bool, ret1 string) {
	var r0 int
	var r1 int
	var r2 rune
	var r5 rune
	var r8 rune
	var r11 rune
	var r15 bool
	var r16 rune
	var r19 rune
	var r22 rune
	var r25 rune
	var r28 rune
	var r32 string
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == 't' {
			frame.Consume()
			r5 = frame.Peek()
			if frame.Flow == 0 {
				if r5 == 'r' {
					frame.Consume()
					r8 = frame.Peek()
					if frame.Flow == 0 {
						if r8 == 'u' {
							frame.Consume()
							r11 = frame.Peek()
							if frame.Flow == 0 {
								if r11 == 'e' {
									frame.Consume()
									r15 = true
									goto block2
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r1)
	r16 = frame.Peek()
	if frame.Flow == 0 {
		if r16 == 'f' {
			frame.Consume()
			r19 = frame.Peek()
			if frame.Flow == 0 {
				if r19 == 'a' {
					frame.Consume()
					r22 = frame.Peek()
					if frame.Flow == 0 {
						if r22 == 'l' {
							frame.Consume()
							r25 = frame.Peek()
							if frame.Flow == 0 {
								if r25 == 's' {
									frame.Consume()
									r28 = frame.Peek()
									if frame.Flow == 0 {
										if r28 == 'e' {
											frame.Consume()
											r15 = false
											goto block2
										}
										frame.Fail()
										goto block3
									}
									goto block3
								}
								frame.Fail()
								goto block3
							}
							goto block3
						}
						frame.Fail()
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block2:
	EndKeyword(frame)
	if frame.Flow == 0 {
		r32 = frame.Slice(r0)
		ret0 = r15
		ret1 = r32
		return
	}
	goto block3
block3:
	return
}

func ParseStringLiteral(frame *runtime.State) (ret0 *StringLiteral) {
	var r0 int
	var r1 string
	r0 = frame.Checkpoint()
	r1 = DecodeString(frame)
	if frame.Flow == 0 {
		ret0 = &StringLiteral{Text: frame.Slice(r0), Value: r1}
		return
	}
	return
}

func Literal(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 rune
	var r2 string
	var r4 *StringLiteral
	var r5 int
	var r6 string
	var r8 bool
	var r9 string
	var r11 rune
	var r14 rune
	var r17 rune
	r0 = frame.Checkpoint()
	r1, r2 = DecodeRune(frame)
	if frame.Flow == 0 {
		ret0 = &RuneLiteral{Text: r2, Value: r1}
		goto block1
	}
	frame.Recover(r0)
	r4 = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		ret0 = r4
		goto block1
	}
	frame.Recover(r0)
	r5, r6 = DecodeInt(frame)
	if frame.Flow == 0 {
		ret0 = &IntLiteral{Text: r6, Value: r5}
		goto block1
	}
	frame.Recover(r0)
	r8, r9 = DecodeBool(frame)
	if frame.Flow == 0 {
		ret0 = &BoolLiteral{Text: r9, Value: r8}
		goto block1
	}
	frame.Recover(r0)
	r11 = frame.Peek()
	if frame.Flow == 0 {
		if r11 == 'n' {
			frame.Consume()
			r14 = frame.Peek()
			if frame.Flow == 0 {
				if r14 == 'i' {
					frame.Consume()
					r17 = frame.Peek()
					if frame.Flow == 0 {
						if r17 == 'l' {
							frame.Consume()
							ret0 = &NilLiteral{}
							goto block1
						}
						frame.Fail()
						goto block2
					}
					goto block2
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block1:
	return
block2:
	return
}

func BinaryOperator(frame *runtime.State) (ret0 string, ret1 int) {
	var r0 int
	var r1 int
	var r2 rune
	var r9 string
	var r10 int
	var r11 int
	var r12 rune
	var r17 string
	var r18 int
	var r19 int
	var r20 int
	var r21 rune
	var r26 int
	var r27 rune
	var r30 rune
	var r35 rune
	var r38 string
	var r39 int
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == '*' {
			goto block1
		}
		if r2 == '/' {
			goto block1
		}
		if r2 == '%' {
			goto block1
		}
		frame.Fail()
		goto block2
	}
	goto block2
block1:
	frame.Consume()
	r9 = frame.Slice(r1)
	r10 = 5
	ret0 = r9
	ret1 = r10
	goto block10
block2:
	frame.Recover(r0)
	r11 = frame.Checkpoint()
	r12 = frame.Peek()
	if frame.Flow == 0 {
		if r12 == '+' {
			goto block3
		}
		if r12 == '-' {
			goto block3
		}
		frame.Fail()
		goto block4
	}
	goto block4
block3:
	frame.Consume()
	r17 = frame.Slice(r11)
	r18 = 4
	ret0 = r17
	ret1 = r18
	goto block10
block4:
	frame.Recover(r0)
	r19 = frame.Checkpoint()
	r20 = frame.Checkpoint()
	r21 = frame.Peek()
	if frame.Flow == 0 {
		if r21 == '<' {
			goto block5
		}
		if r21 == '>' {
			goto block5
		}
		frame.Fail()
		goto block7
	}
	goto block7
block5:
	frame.Consume()
	r26 = frame.Checkpoint()
	r27 = frame.Peek()
	if frame.Flow == 0 {
		if r27 == '=' {
			frame.Consume()
			goto block9
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(r26)
	goto block9
block7:
	frame.Recover(r20)
	r30 = frame.Peek()
	if frame.Flow == 0 {
		if r30 == '!' {
			goto block8
		}
		if r30 == '=' {
			goto block8
		}
		frame.Fail()
		goto block11
	}
	goto block11
block8:
	frame.Consume()
	r35 = frame.Peek()
	if frame.Flow == 0 {
		if r35 == '=' {
			frame.Consume()
			goto block9
		}
		frame.Fail()
		goto block11
	}
	goto block11
block9:
	r38 = frame.Slice(r19)
	r39 = 3
	ret0 = r38
	ret1 = r39
	goto block10
block10:
	return
block11:
	return
}

func StringMatchExpr(frame *runtime.State) (ret0 *StringMatch) {
	var r0 rune
	var r3 TextMatch
	var r4 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '/' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = ParseMatchChoice(frame)
				if frame.Flow == 0 {
					S(frame)
					if frame.Flow == 0 {
						r4 = frame.Peek()
						if frame.Flow == 0 {
							if r4 == '/' {
								frame.Consume()
								ret0 = &StringMatch{Match: r3}
								return
							}
							frame.Fail()
							goto block1
						}
						goto block1
					}
					goto block1
				}
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	return
}

func RuneMatchExpr(frame *runtime.State) (ret0 *RuneMatch) {
	var r0 rune
	var r3 *RuneRangeMatch
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '$' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = MatchRune(frame)
				if frame.Flow == 0 {
					ret0 = &RuneMatch{Match: r3}
					return
				}
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	return
}

func ParseStructTypeRef(frame *runtime.State) (ret0 ASTTypeRef) {
	var r0 int
	var r1 *Id
	var r2 rune
	var r5 *Id
	var r7 *Id
	r0 = frame.Checkpoint()
	r1 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r2 = frame.Peek()
			if frame.Flow == 0 {
				if r2 == '.' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r5 = Ident(frame)
						if frame.Flow == 0 {
							ret0 = &QualifiedTypeRef{Package: r1, Name: r5}
							goto block2
						}
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		goto block1
	}
	goto block1
block1:
	frame.Recover(r0)
	r7 = Ident(frame)
	if frame.Flow == 0 {
		ret0 = &TypeRef{Name: r7}
		goto block2
	}
	return
block2:
	return
}

func ParseListTypeRef(frame *runtime.State) (ret0 *ListTypeRef) {
	var r0 rune
	var r3 rune
	var r6 ASTTypeRef
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '[' {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				if r3 == ']' {
					frame.Consume()
					r6 = ParseTypeRef(frame)
					if frame.Flow == 0 {
						ret0 = &ListTypeRef{Type: r6}
						return
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	return
}

func ParseTypeRef(frame *runtime.State) (ret0 ASTTypeRef) {
	var r0 int
	var r1 ASTTypeRef
	var r2 *ListTypeRef
	r0 = frame.Checkpoint()
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block1
	}
	frame.Recover(r0)
	r2 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		ret0 = r2
		goto block1
	}
	return
block1:
	return
}

func ParseDestructure(frame *runtime.State) (ret0 Destructure) {
	var r0 int
	var r1 ASTTypeRef
	var r2 rune
	var r6 []*DestructureField
	var r7 int
	var r8 *Id
	var r9 rune
	var r12 Destructure
	var r15 rune
	var r19 *ListTypeRef
	var r20 rune
	var r24 []Destructure
	var r25 int
	var r26 Destructure
	var r27 []Destructure
	var r28 []Destructure
	var r29 rune
	var r33 ASTExpr
	r0 = frame.Checkpoint()
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r2 = frame.Peek()
			if frame.Flow == 0 {
				if r2 == '{' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r6 = []*DestructureField{}
						goto block1
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		goto block3
	}
	goto block3
block1:
	r7 = frame.Checkpoint()
	r8 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r9 = frame.Peek()
			if frame.Flow == 0 {
				if r9 == ':' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r12 = ParseDestructure(frame)
						if frame.Flow == 0 {
							S(frame)
							if frame.Flow == 0 {
								r6 = append(r6, &DestructureField{Name: r8, Destructure: r12})
								goto block1
							}
							goto block2
						}
						goto block2
					}
					goto block2
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r7)
	r15 = frame.Peek()
	if frame.Flow == 0 {
		if r15 == '}' {
			frame.Consume()
			ret0 = &DestructureStruct{Type: r1, Args: r6}
			goto block7
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(r0)
	r19 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r20 = frame.Peek()
			if frame.Flow == 0 {
				if r20 == '{' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r24 = []Destructure{}
						goto block4
					}
					goto block6
				}
				frame.Fail()
				goto block6
			}
			goto block6
		}
		goto block6
	}
	goto block6
block4:
	r25 = frame.Checkpoint()
	r26 = ParseDestructure(frame)
	if frame.Flow == 0 {
		r27 = append(r24, r26)
		S(frame)
		if frame.Flow == 0 {
			r24 = r27
			goto block4
		}
		r28 = r27
		goto block5
	}
	r28 = r24
	goto block5
block5:
	frame.Recover(r25)
	r29 = frame.Peek()
	if frame.Flow == 0 {
		if r29 == '}' {
			frame.Consume()
			ret0 = &DestructureList{Type: r19, Args: r28}
			goto block7
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(r0)
	r33 = Literal(frame)
	if frame.Flow == 0 {
		ret0 = &DestructureValue{Expr: r33}
		goto block7
	}
	return
block7:
	return
}

func ParseRuneFilterRune(frame *runtime.State) (ret0 rune) {
	var r0 int
	var r1 rune
	var r8 rune
	var r11 int
	var r12 rune
	var r13 rune
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 == ']' {
			goto block1
		}
		if r1 == '-' {
			goto block1
		}
		if r1 == '\\' {
			goto block1
		}
		frame.Consume()
		ret0 = r1
		goto block3
	}
	goto block2
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(r0)
	r8 = frame.Peek()
	if frame.Flow == 0 {
		if r8 == '\\' {
			frame.Consume()
			r11 = frame.Checkpoint()
			r12 = EscapedChar(frame)
			if frame.Flow == 0 {
				ret0 = r12
				goto block3
			}
			frame.Recover(r11)
			r13 = frame.Peek()
			if frame.Flow == 0 {
				frame.Consume()
				ret0 = r13
				goto block3
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block3:
	return
block4:
	return
}

func ParseRuneFilter(frame *runtime.State) (ret0 *RuneFilter) {
	var r0 rune
	var r1 int
	var r2 rune
	var r5 rune
	var r6 rune
	r0 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = frame.Peek()
		if frame.Flow == 0 {
			if r2 == '-' {
				frame.Consume()
				r5 = ParseRuneFilterRune(frame)
				if frame.Flow == 0 {
					r6 = r5
					goto block2
				}
				goto block1
			}
			frame.Fail()
			goto block1
		}
		goto block1
	}
	return
block1:
	frame.Recover(r1)
	r6 = r0
	goto block2
block2:
	ret0 = &RuneFilter{Min: r0, Max: r6}
	return
}

func MatchRune(frame *runtime.State) (ret0 *RuneRangeMatch) {
	var r0 rune
	var r3 bool
	var r4 []*RuneFilter
	var r5 int
	var r6 rune
	var r10 bool
	var r11 []*RuneFilter
	var r12 int
	var r13 *RuneFilter
	var r15 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '[' {
			frame.Consume()
			r3 = false
			r4 = []*RuneFilter{}
			r5 = frame.Checkpoint()
			r6 = frame.Peek()
			if frame.Flow == 0 {
				if r6 == '^' {
					frame.Consume()
					r10, r11 = true, r4
					goto block2
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block3
	}
	goto block3
block1:
	frame.Recover(r5)
	r10, r11 = r3, r4
	goto block2
block2:
	r12 = frame.Checkpoint()
	r13 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		r11 = append(r11, r13)
		goto block2
	}
	frame.Recover(r12)
	r15 = frame.Peek()
	if frame.Flow == 0 {
		if r15 == ']' {
			frame.Consume()
			ret0 = &RuneRangeMatch{Invert: r10, Filters: r11}
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	return
}

func Atom(frame *runtime.State) (ret0 TextMatch) {
	var r0 int
	var r1 *RuneRangeMatch
	var r2 string
	var r4 rune
	var r7 TextMatch
	var r8 rune
	r0 = frame.Checkpoint()
	r1 = MatchRune(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block1
	}
	frame.Recover(r0)
	r2 = DecodeString(frame)
	if frame.Flow == 0 {
		ret0 = &StringLiteralMatch{Value: r2}
		goto block1
	}
	frame.Recover(r0)
	r4 = frame.Peek()
	if frame.Flow == 0 {
		if r4 == '(' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r7 = ParseMatchChoice(frame)
				if frame.Flow == 0 {
					S(frame)
					if frame.Flow == 0 {
						r8 = frame.Peek()
						if frame.Flow == 0 {
							if r8 == ')' {
								frame.Consume()
								ret0 = r7
								goto block1
							}
							frame.Fail()
							goto block2
						}
						goto block2
					}
					goto block2
				}
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block1:
	return
block2:
	return
}

func MatchPostfix(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 rune
	var r7 rune
	var r12 rune
	r0 = Atom(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		S(frame)
		if frame.Flow == 0 {
			r2 = frame.Peek()
			if frame.Flow == 0 {
				if r2 == '*' {
					frame.Consume()
					ret0 = &MatchRepeat{Match: r0, Min: 0}
					goto block4
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		goto block1
	}
	return
block1:
	frame.Recover(r1)
	S(frame)
	if frame.Flow == 0 {
		r7 = frame.Peek()
		if frame.Flow == 0 {
			if r7 == '+' {
				frame.Consume()
				ret0 = &MatchRepeat{Match: r0, Min: 1}
				goto block4
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r1)
	S(frame)
	if frame.Flow == 0 {
		r12 = frame.Peek()
		if frame.Flow == 0 {
			if r12 == '?' {
				frame.Consume()
				ret0 = &MatchChoice{Matches: []TextMatch{r0, &MatchSequence{Matches: []TextMatch{}}}}
				goto block4
			}
			frame.Fail()
			goto block3
		}
		goto block3
	}
	goto block3
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
	var r7 bool
	var r8 rune
	var r11 TextMatch
	var r13 TextMatch
	r0 = frame.Checkpoint()
	r1 = false
	r2 = frame.Checkpoint()
	r3 = frame.Peek()
	if frame.Flow == 0 {
		if r3 == '!' {
			frame.Consume()
			r7 = true
			goto block2
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r2)
	r8 = frame.Peek()
	if frame.Flow == 0 {
		if r8 == '&' {
			frame.Consume()
			r7 = r1
			goto block2
		}
		frame.Fail()
		goto block3
	}
	goto block3
block2:
	S(frame)
	if frame.Flow == 0 {
		r11 = MatchPostfix(frame)
		if frame.Flow == 0 {
			ret0 = &MatchLookahead{Invert: r7, Match: r11}
			goto block4
		}
		goto block3
	}
	goto block3
block3:
	frame.Recover(r0)
	r13 = MatchPostfix(frame)
	if frame.Flow == 0 {
		ret0 = r13
		goto block4
	}
	return
block4:
	return
}

func Sequence(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 []TextMatch
	var r3 TextMatch
	var r5 []TextMatch
	var r6 int
	var r7 TextMatch
	r0 = MatchPrefix(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = []TextMatch{r0}
		S(frame)
		if frame.Flow == 0 {
			r3 = MatchPrefix(frame)
			if frame.Flow == 0 {
				r5 = append(r2, r3)
				goto block1
			}
			goto block3
		}
		goto block3
	}
	return
block1:
	r6 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r7 = MatchPrefix(frame)
		if frame.Flow == 0 {
			r5 = append(r5, r7)
			goto block1
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r6)
	ret0 = &MatchSequence{Matches: r5}
	goto block4
block3:
	frame.Recover(r1)
	ret0 = r0
	goto block4
block4:
	return
}

func ParseMatchChoice(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 []TextMatch
	var r3 rune
	var r6 TextMatch
	var r8 []TextMatch
	var r9 int
	var r10 rune
	var r13 TextMatch
	r0 = Sequence(frame)
	if frame.Flow == 0 {
		r1 = frame.Checkpoint()
		r2 = []TextMatch{r0}
		S(frame)
		if frame.Flow == 0 {
			r3 = frame.Peek()
			if frame.Flow == 0 {
				if r3 == '|' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r6 = Sequence(frame)
						if frame.Flow == 0 {
							r8 = append(r2, r6)
							goto block1
						}
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		goto block3
	}
	return
block1:
	r9 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r10 = frame.Peek()
		if frame.Flow == 0 {
			if r10 == '|' {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r13 = Sequence(frame)
					if frame.Flow == 0 {
						r8 = append(r8, r13)
						goto block1
					}
					goto block2
				}
				goto block2
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r9)
	ret0 = &MatchChoice{Matches: r8}
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
	var r4 []ASTExpr
	var r5 int
	var r6 rune
	var r9 ASTExpr
	var r11 []ASTExpr
	r0 = []ASTExpr{}
	r1 = frame.Checkpoint()
	r2 = ParseExpr(frame)
	if frame.Flow == 0 {
		r4 = append(r0, r2)
		goto block1
	}
	frame.Recover(r1)
	r11 = r0
	goto block3
block1:
	r5 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r6 = frame.Peek()
		if frame.Flow == 0 {
			if r6 == ',' {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r9 = ParseExpr(frame)
					if frame.Flow == 0 {
						r4 = append(r4, r9)
						goto block1
					}
					goto block2
				}
				goto block2
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r5)
	r11 = r4
	goto block3
block3:
	ret0 = r11
	return
}

func ParseTargetList(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 *NameRef
	var r2 []ASTExpr
	var r3 int
	var r4 rune
	var r7 *NameRef
	r0 = ParseNameRef(frame)
	if frame.Flow == 0 {
		r2 = []ASTExpr{r0}
		goto block1
	}
	return
block1:
	r3 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r4 = frame.Peek()
		if frame.Flow == 0 {
			if r4 == ',' {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r7 = ParseNameRef(frame)
					if frame.Flow == 0 {
						r2 = append(r2, r7)
						goto block1
					}
					goto block2
				}
				goto block2
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r3)
	ret0 = r2
	return
}

func ParseNamedExpr(frame *runtime.State) (ret0 *NamedExpr) {
	var r0 *Id
	var r1 rune
	var r4 ASTExpr
	r0 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r1 = frame.Peek()
			if frame.Flow == 0 {
				if r1 == ':' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r4 = ParseExpr(frame)
						if frame.Flow == 0 {
							ret0 = &NamedExpr{Name: r0, Expr: r4}
							return
						}
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		goto block1
	}
	goto block1
block1:
	return
}

func ParseNamedExprList(frame *runtime.State) (ret0 []*NamedExpr) {
	var r0 []*NamedExpr
	var r1 int
	var r2 *NamedExpr
	var r4 []*NamedExpr
	var r5 int
	var r6 rune
	var r9 *NamedExpr
	var r11 []*NamedExpr
	r0 = []*NamedExpr{}
	r1 = frame.Checkpoint()
	r2 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		r4 = append(r0, r2)
		goto block1
	}
	frame.Recover(r1)
	r11 = r0
	goto block3
block1:
	r5 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r6 = frame.Peek()
		if frame.Flow == 0 {
			if r6 == ',' {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r9 = ParseNamedExpr(frame)
					if frame.Flow == 0 {
						r4 = append(r4, r9)
						goto block1
					}
					goto block2
				}
				goto block2
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r5)
	r11 = r4
	goto block3
block3:
	ret0 = r11
	return
}

func ParseReturnTypeList(frame *runtime.State) (ret0 []ASTTypeRef) {
	var r0 int
	var r1 []ASTTypeRef
	var r2 ASTTypeRef
	r0 = frame.Checkpoint()
	r1 = ParseParenthTypeList(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block1
	}
	frame.Recover(r0)
	r2 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		ret0 = []ASTTypeRef{r2}
		goto block1
	}
	frame.Recover(r0)
	ret0 = []ASTTypeRef{}
	goto block1
block1:
	return
}

func PrimaryExpr(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 ASTExpr
	var r2 rune
	var r5 rune
	var r8 rune
	var r11 rune
	var r14 rune
	var r17 []ASTExpr
	var r19 rune
	var r22 rune
	var r25 rune
	var r28 rune
	var r31 rune
	var r34 rune
	var r37 rune
	var r40 rune
	var r44 rune
	var r47 rune
	var r50 rune
	var r53 rune
	var r56 rune
	var r59 rune
	var r62 rune
	var r65 ASTTypeRef
	var r66 rune
	var r69 ASTExpr
	var r70 rune
	var r74 rune
	var r77 rune
	var r80 rune
	var r83 rune
	var r86 rune
	var r89 rune
	var r92 rune
	var r95 *Id
	var r96 rune
	var r99 ASTExpr
	var r100 rune
	var r108 *Id
	var r109 rune
	var r112 []ASTExpr
	var r113 rune
	var r117 ASTTypeRef
	var r118 rune
	var r121 []*NamedExpr
	var r122 rune
	var r126 *ListTypeRef
	var r127 rune
	var r130 []ASTExpr
	var r131 rune
	var r135 *StringMatch
	var r136 *RuneMatch
	var r137 rune
	var r140 ASTExpr
	var r141 rune
	var r144 *NameRef
	r0 = frame.Checkpoint()
	r1 = Literal(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block9
	}
	frame.Recover(r0)
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == 's' {
			frame.Consume()
			r5 = frame.Peek()
			if frame.Flow == 0 {
				if r5 == 'l' {
					frame.Consume()
					r8 = frame.Peek()
					if frame.Flow == 0 {
						if r8 == 'i' {
							frame.Consume()
							r11 = frame.Peek()
							if frame.Flow == 0 {
								if r11 == 'c' {
									frame.Consume()
									r14 = frame.Peek()
									if frame.Flow == 0 {
										if r14 == 'e' {
											frame.Consume()
											EndKeyword(frame)
											if frame.Flow == 0 {
												S(frame)
												if frame.Flow == 0 {
													r17 = ParseCodeBlock(frame)
													if frame.Flow == 0 {
														ret0 = &Slice{Block: r17}
														goto block9
													}
													goto block1
												}
												goto block1
											}
											goto block1
										}
										frame.Fail()
										goto block1
									}
									goto block1
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r0)
	r19 = frame.Peek()
	if frame.Flow == 0 {
		if r19 == 'p' {
			frame.Consume()
			r22 = frame.Peek()
			if frame.Flow == 0 {
				if r22 == 'o' {
					frame.Consume()
					r25 = frame.Peek()
					if frame.Flow == 0 {
						if r25 == 's' {
							frame.Consume()
							r28 = frame.Peek()
							if frame.Flow == 0 {
								if r28 == 'i' {
									frame.Consume()
									r31 = frame.Peek()
									if frame.Flow == 0 {
										if r31 == 't' {
											frame.Consume()
											r34 = frame.Peek()
											if frame.Flow == 0 {
												if r34 == 'i' {
													frame.Consume()
													r37 = frame.Peek()
													if frame.Flow == 0 {
														if r37 == 'o' {
															frame.Consume()
															r40 = frame.Peek()
															if frame.Flow == 0 {
																if r40 == 'n' {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		ret0 = &Position{}
																		goto block9
																	}
																	goto block2
																}
																frame.Fail()
																goto block2
															}
															goto block2
														}
														frame.Fail()
														goto block2
													}
													goto block2
												}
												frame.Fail()
												goto block2
											}
											goto block2
										}
										frame.Fail()
										goto block2
									}
									goto block2
								}
								frame.Fail()
								goto block2
							}
							goto block2
						}
						frame.Fail()
						goto block2
					}
					goto block2
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(r0)
	r44 = frame.Peek()
	if frame.Flow == 0 {
		if r44 == 'c' {
			frame.Consume()
			r47 = frame.Peek()
			if frame.Flow == 0 {
				if r47 == 'o' {
					frame.Consume()
					r50 = frame.Peek()
					if frame.Flow == 0 {
						if r50 == 'e' {
							frame.Consume()
							r53 = frame.Peek()
							if frame.Flow == 0 {
								if r53 == 'r' {
									frame.Consume()
									r56 = frame.Peek()
									if frame.Flow == 0 {
										if r56 == 'c' {
											frame.Consume()
											r59 = frame.Peek()
											if frame.Flow == 0 {
												if r59 == 'e' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r62 = frame.Peek()
															if frame.Flow == 0 {
																if r62 == '(' {
																	frame.Consume()
																	S(frame)
																	if frame.Flow == 0 {
																		r65 = ParseTypeRef(frame)
																		if frame.Flow == 0 {
																			S(frame)
																			if frame.Flow == 0 {
																				r66 = frame.Peek()
																				if frame.Flow == 0 {
																					if r66 == ',' {
																						frame.Consume()
																						S(frame)
																						if frame.Flow == 0 {
																							r69 = ParseExpr(frame)
																							if frame.Flow == 0 {
																								S(frame)
																								if frame.Flow == 0 {
																									r70 = frame.Peek()
																									if frame.Flow == 0 {
																										if r70 == ')' {
																											frame.Consume()
																											ret0 = &Coerce{Type: r65, Expr: r69}
																											goto block9
																										}
																										frame.Fail()
																										goto block3
																									}
																									goto block3
																								}
																								goto block3
																							}
																							goto block3
																						}
																						goto block3
																					}
																					frame.Fail()
																					goto block3
																				}
																				goto block3
																			}
																			goto block3
																		}
																		goto block3
																	}
																	goto block3
																}
																frame.Fail()
																goto block3
															}
															goto block3
														}
														goto block3
													}
													goto block3
												}
												frame.Fail()
												goto block3
											}
											goto block3
										}
										frame.Fail()
										goto block3
									}
									goto block3
								}
								frame.Fail()
								goto block3
							}
							goto block3
						}
						frame.Fail()
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(r0)
	r74 = frame.Peek()
	if frame.Flow == 0 {
		if r74 == 'a' {
			frame.Consume()
			r77 = frame.Peek()
			if frame.Flow == 0 {
				if r77 == 'p' {
					frame.Consume()
					r80 = frame.Peek()
					if frame.Flow == 0 {
						if r80 == 'p' {
							frame.Consume()
							r83 = frame.Peek()
							if frame.Flow == 0 {
								if r83 == 'e' {
									frame.Consume()
									r86 = frame.Peek()
									if frame.Flow == 0 {
										if r86 == 'n' {
											frame.Consume()
											r89 = frame.Peek()
											if frame.Flow == 0 {
												if r89 == 'd' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r92 = frame.Peek()
															if frame.Flow == 0 {
																if r92 == '(' {
																	frame.Consume()
																	S(frame)
																	if frame.Flow == 0 {
																		r95 = Ident(frame)
																		if frame.Flow == 0 {
																			S(frame)
																			if frame.Flow == 0 {
																				r96 = frame.Peek()
																				if frame.Flow == 0 {
																					if r96 == ',' {
																						frame.Consume()
																						S(frame)
																						if frame.Flow == 0 {
																							r99 = ParseExpr(frame)
																							if frame.Flow == 0 {
																								S(frame)
																								if frame.Flow == 0 {
																									r100 = frame.Peek()
																									if frame.Flow == 0 {
																										if r100 == ')' {
																											frame.Consume()
																											ret0 = &Assign{Expr: &Append{List: &NameRef{Name: r95}, Expr: r99}, Targets: []ASTExpr{&NameRef{Name: r95}}}
																											goto block9
																										}
																										frame.Fail()
																										goto block4
																									}
																									goto block4
																								}
																								goto block4
																							}
																							goto block4
																						}
																						goto block4
																					}
																					frame.Fail()
																					goto block4
																				}
																				goto block4
																			}
																			goto block4
																		}
																		goto block4
																	}
																	goto block4
																}
																frame.Fail()
																goto block4
															}
															goto block4
														}
														goto block4
													}
													goto block4
												}
												frame.Fail()
												goto block4
											}
											goto block4
										}
										frame.Fail()
										goto block4
									}
									goto block4
								}
								frame.Fail()
								goto block4
							}
							goto block4
						}
						frame.Fail()
						goto block4
					}
					goto block4
				}
				frame.Fail()
				goto block4
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(r0)
	r108 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r109 = frame.Peek()
			if frame.Flow == 0 {
				if r109 == '(' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r112 = ParseExprList(frame)
						if frame.Flow == 0 {
							S(frame)
							if frame.Flow == 0 {
								r113 = frame.Peek()
								if frame.Flow == 0 {
									if r113 == ')' {
										frame.Consume()
										ret0 = &Call{Name: r108, Args: r112}
										goto block9
									}
									frame.Fail()
									goto block5
								}
								goto block5
							}
							goto block5
						}
						goto block5
					}
					goto block5
				}
				frame.Fail()
				goto block5
			}
			goto block5
		}
		goto block5
	}
	goto block5
block5:
	frame.Recover(r0)
	r117 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r118 = frame.Peek()
			if frame.Flow == 0 {
				if r118 == '{' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r121 = ParseNamedExprList(frame)
						if frame.Flow == 0 {
							S(frame)
							if frame.Flow == 0 {
								r122 = frame.Peek()
								if frame.Flow == 0 {
									if r122 == '}' {
										frame.Consume()
										ret0 = &Construct{Type: r117, Args: r121}
										goto block9
									}
									frame.Fail()
									goto block6
								}
								goto block6
							}
							goto block6
						}
						goto block6
					}
					goto block6
				}
				frame.Fail()
				goto block6
			}
			goto block6
		}
		goto block6
	}
	goto block6
block6:
	frame.Recover(r0)
	r126 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r127 = frame.Peek()
			if frame.Flow == 0 {
				if r127 == '{' {
					frame.Consume()
					S(frame)
					if frame.Flow == 0 {
						r130 = ParseExprList(frame)
						if frame.Flow == 0 {
							S(frame)
							if frame.Flow == 0 {
								r131 = frame.Peek()
								if frame.Flow == 0 {
									if r131 == '}' {
										frame.Consume()
										ret0 = &ConstructList{Type: r126, Args: r130}
										goto block9
									}
									frame.Fail()
									goto block7
								}
								goto block7
							}
							goto block7
						}
						goto block7
					}
					goto block7
				}
				frame.Fail()
				goto block7
			}
			goto block7
		}
		goto block7
	}
	goto block7
block7:
	frame.Recover(r0)
	r135 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		ret0 = r135
		goto block9
	}
	frame.Recover(r0)
	r136 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		ret0 = r136
		goto block9
	}
	frame.Recover(r0)
	r137 = frame.Peek()
	if frame.Flow == 0 {
		if r137 == '(' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r140 = ParseExpr(frame)
				if frame.Flow == 0 {
					S(frame)
					if frame.Flow == 0 {
						r141 = frame.Peek()
						if frame.Flow == 0 {
							if r141 == ')' {
								frame.Consume()
								ret0 = r140
								goto block9
							}
							frame.Fail()
							goto block8
						}
						goto block8
					}
					goto block8
				}
				goto block8
			}
			goto block8
		}
		frame.Fail()
		goto block8
	}
	goto block8
block8:
	frame.Recover(r0)
	r144 = ParseNameRef(frame)
	if frame.Flow == 0 {
		ret0 = r144
		goto block9
	}
	return
block9:
	return
}

func ParseNameRef(frame *runtime.State) (ret0 *NameRef) {
	var r0 *Id
	r0 = Ident(frame)
	if frame.Flow == 0 {
		ret0 = &NameRef{Name: r0}
		return
	}
	return
}

func ParseBinaryOp(frame *runtime.State, r0 int) (ret0 ASTExpr) {
	var r1 ASTExpr
	var r2 ASTExpr
	var r3 int
	var r4 string
	var r5 int
	var r9 ASTExpr
	r1 = PrimaryExpr(frame)
	if frame.Flow == 0 {
		r2 = r1
		goto block1
	}
	return
block1:
	r3 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r4, r5 = BinaryOperator(frame)
		if frame.Flow == 0 {
			if r5 < r0 {
				frame.Fail()
				goto block2
			}
			S(frame)
			if frame.Flow == 0 {
				r9 = ParseBinaryOp(frame, r5+1)
				if frame.Flow == 0 {
					r2 = &BinaryOp{Left: r2, Op: r4, Right: r9}
					goto block1
				}
				goto block2
			}
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r3)
	ret0 = r2
	return
}

func ParseExpr(frame *runtime.State) (ret0 ASTExpr) {
	var r1 ASTExpr
	r1 = ParseBinaryOp(frame, 1)
	if frame.Flow == 0 {
		ret0 = r1
		return
	}
	return
}

func ParseCompoundStatement(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 rune
	var r4 rune
	var r7 rune
	var r10 rune
	var r13 []ASTExpr
	var r16 rune
	var r19 rune
	var r22 rune
	var r25 rune
	var r28 []ASTExpr
	var r31 rune
	var r34 rune
	var r37 rune
	var r40 rune
	var r43 rune
	var r46 rune
	var r49 []ASTExpr
	var r50 [][]ASTExpr
	var r51 rune
	var r54 rune
	var r57 []ASTExpr
	var r59 [][]ASTExpr
	var r60 int
	var r61 rune
	var r64 rune
	var r67 []ASTExpr
	var r70 rune
	var r73 rune
	var r76 rune
	var r79 rune
	var r82 rune
	var r85 rune
	var r88 rune
	var r91 rune
	var r94 []ASTExpr
	var r96 rune
	var r99 rune
	var r102 ASTExpr
	var r103 []ASTExpr
	r0 = frame.Checkpoint()
	r1 = frame.Peek()
	if frame.Flow == 0 {
		if r1 == 's' {
			frame.Consume()
			r4 = frame.Peek()
			if frame.Flow == 0 {
				if r4 == 't' {
					frame.Consume()
					r7 = frame.Peek()
					if frame.Flow == 0 {
						if r7 == 'a' {
							frame.Consume()
							r10 = frame.Peek()
							if frame.Flow == 0 {
								if r10 == 'r' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										if frame.Flow == 0 {
											r13 = ParseCodeBlock(frame)
											if frame.Flow == 0 {
												ret0 = &Repeat{Block: r13, Min: 0}
												goto block7
											}
											goto block1
										}
										goto block1
									}
									goto block1
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r0)
	r16 = frame.Peek()
	if frame.Flow == 0 {
		if r16 == 'p' {
			frame.Consume()
			r19 = frame.Peek()
			if frame.Flow == 0 {
				if r19 == 'l' {
					frame.Consume()
					r22 = frame.Peek()
					if frame.Flow == 0 {
						if r22 == 'u' {
							frame.Consume()
							r25 = frame.Peek()
							if frame.Flow == 0 {
								if r25 == 's' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										if frame.Flow == 0 {
											r28 = ParseCodeBlock(frame)
											if frame.Flow == 0 {
												ret0 = &Repeat{Block: r28, Min: 1}
												goto block7
											}
											goto block2
										}
										goto block2
									}
									goto block2
								}
								frame.Fail()
								goto block2
							}
							goto block2
						}
						frame.Fail()
						goto block2
					}
					goto block2
				}
				frame.Fail()
				goto block2
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(r0)
	r31 = frame.Peek()
	if frame.Flow == 0 {
		if r31 == 'c' {
			frame.Consume()
			r34 = frame.Peek()
			if frame.Flow == 0 {
				if r34 == 'h' {
					frame.Consume()
					r37 = frame.Peek()
					if frame.Flow == 0 {
						if r37 == 'o' {
							frame.Consume()
							r40 = frame.Peek()
							if frame.Flow == 0 {
								if r40 == 'o' {
									frame.Consume()
									r43 = frame.Peek()
									if frame.Flow == 0 {
										if r43 == 's' {
											frame.Consume()
											r46 = frame.Peek()
											if frame.Flow == 0 {
												if r46 == 'e' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r49 = ParseCodeBlock(frame)
															if frame.Flow == 0 {
																r50 = [][]ASTExpr{r49}
																S(frame)
																if frame.Flow == 0 {
																	r51 = frame.Peek()
																	if frame.Flow == 0 {
																		if r51 == 'o' {
																			frame.Consume()
																			r54 = frame.Peek()
																			if frame.Flow == 0 {
																				if r54 == 'r' {
																					frame.Consume()
																					EndKeyword(frame)
																					if frame.Flow == 0 {
																						S(frame)
																						if frame.Flow == 0 {
																							r57 = ParseCodeBlock(frame)
																							if frame.Flow == 0 {
																								r59 = append(r50, r57)
																								goto block3
																							}
																							goto block5
																						}
																						goto block5
																					}
																					goto block5
																				}
																				frame.Fail()
																				goto block5
																			}
																			goto block5
																		}
																		frame.Fail()
																		goto block5
																	}
																	goto block5
																}
																goto block5
															}
															goto block5
														}
														goto block5
													}
													goto block5
												}
												frame.Fail()
												goto block5
											}
											goto block5
										}
										frame.Fail()
										goto block5
									}
									goto block5
								}
								frame.Fail()
								goto block5
							}
							goto block5
						}
						frame.Fail()
						goto block5
					}
					goto block5
				}
				frame.Fail()
				goto block5
			}
			goto block5
		}
		frame.Fail()
		goto block5
	}
	goto block5
block3:
	r60 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r61 = frame.Peek()
		if frame.Flow == 0 {
			if r61 == 'o' {
				frame.Consume()
				r64 = frame.Peek()
				if frame.Flow == 0 {
					if r64 == 'r' {
						frame.Consume()
						EndKeyword(frame)
						if frame.Flow == 0 {
							S(frame)
							if frame.Flow == 0 {
								r67 = ParseCodeBlock(frame)
								if frame.Flow == 0 {
									r59 = append(r59, r67)
									goto block3
								}
								goto block4
							}
							goto block4
						}
						goto block4
					}
					frame.Fail()
					goto block4
				}
				goto block4
			}
			frame.Fail()
			goto block4
		}
		goto block4
	}
	goto block4
block4:
	frame.Recover(r60)
	ret0 = &Choice{Blocks: r59}
	goto block7
block5:
	frame.Recover(r0)
	r70 = frame.Peek()
	if frame.Flow == 0 {
		if r70 == 'q' {
			frame.Consume()
			r73 = frame.Peek()
			if frame.Flow == 0 {
				if r73 == 'u' {
					frame.Consume()
					r76 = frame.Peek()
					if frame.Flow == 0 {
						if r76 == 'e' {
							frame.Consume()
							r79 = frame.Peek()
							if frame.Flow == 0 {
								if r79 == 's' {
									frame.Consume()
									r82 = frame.Peek()
									if frame.Flow == 0 {
										if r82 == 't' {
											frame.Consume()
											r85 = frame.Peek()
											if frame.Flow == 0 {
												if r85 == 'i' {
													frame.Consume()
													r88 = frame.Peek()
													if frame.Flow == 0 {
														if r88 == 'o' {
															frame.Consume()
															r91 = frame.Peek()
															if frame.Flow == 0 {
																if r91 == 'n' {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		S(frame)
																		if frame.Flow == 0 {
																			r94 = ParseCodeBlock(frame)
																			if frame.Flow == 0 {
																				ret0 = &Optional{Block: r94}
																				goto block7
																			}
																			goto block6
																		}
																		goto block6
																	}
																	goto block6
																}
																frame.Fail()
																goto block6
															}
															goto block6
														}
														frame.Fail()
														goto block6
													}
													goto block6
												}
												frame.Fail()
												goto block6
											}
											goto block6
										}
										frame.Fail()
										goto block6
									}
									goto block6
								}
								frame.Fail()
								goto block6
							}
							goto block6
						}
						frame.Fail()
						goto block6
					}
					goto block6
				}
				frame.Fail()
				goto block6
			}
			goto block6
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(r0)
	r96 = frame.Peek()
	if frame.Flow == 0 {
		if r96 == 'i' {
			frame.Consume()
			r99 = frame.Peek()
			if frame.Flow == 0 {
				if r99 == 'f' {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						S(frame)
						if frame.Flow == 0 {
							r102 = ParseExpr(frame)
							if frame.Flow == 0 {
								S(frame)
								if frame.Flow == 0 {
									r103 = ParseCodeBlock(frame)
									if frame.Flow == 0 {
										ret0 = &If{Expr: r102, Block: r103}
										goto block7
									}
									goto block8
								}
								goto block8
							}
							goto block8
						}
						goto block8
					}
					goto block8
				}
				frame.Fail()
				goto block8
			}
			goto block8
		}
		frame.Fail()
		goto block8
	}
	goto block8
block7:
	return
block8:
	return
}

func EOS(frame *runtime.State) {
	var r0 rune
	S(frame)
	if frame.Flow == 0 {
		r0 = frame.Peek()
		if frame.Flow == 0 {
			if r0 == ';' {
				frame.Consume()
				return
			}
			frame.Fail()
			goto block1
		}
		goto block1
	}
	goto block1
block1:
	return
}

func ParseStatement(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 ASTExpr
	var r2 rune
	var r5 rune
	var r8 rune
	var r11 *NameRef
	var r12 ASTTypeRef
	var r13 ASTExpr
	var r14 int
	var r15 rune
	var r18 ASTExpr
	var r19 ASTExpr
	var r23 rune
	var r26 rune
	var r29 rune
	var r32 rune
	var r36 rune
	var r39 rune
	var r42 rune
	var r45 rune
	var r48 rune
	var r51 rune
	var r54 []ASTExpr
	var r56 []ASTExpr
	var r57 bool
	var r58 int
	var r59 rune
	var r62 rune
	var r66 bool
	var r67 rune
	var r70 ASTExpr
	var r72 ASTExpr
	r0 = frame.Checkpoint()
	r1 = ParseCompoundStatement(frame)
	if frame.Flow == 0 {
		ret0 = r1
		goto block9
	}
	frame.Recover(r0)
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == 'v' {
			frame.Consume()
			r5 = frame.Peek()
			if frame.Flow == 0 {
				if r5 == 'a' {
					frame.Consume()
					r8 = frame.Peek()
					if frame.Flow == 0 {
						if r8 == 'r' {
							frame.Consume()
							EndKeyword(frame)
							if frame.Flow == 0 {
								S(frame)
								if frame.Flow == 0 {
									r11 = ParseNameRef(frame)
									if frame.Flow == 0 {
										S(frame)
										if frame.Flow == 0 {
											r12 = ParseTypeRef(frame)
											if frame.Flow == 0 {
												r13 = nil
												r14 = frame.Checkpoint()
												S(frame)
												if frame.Flow == 0 {
													r15 = frame.Peek()
													if frame.Flow == 0 {
														if r15 == '=' {
															frame.Consume()
															S(frame)
															if frame.Flow == 0 {
																r18 = ParseExpr(frame)
																if frame.Flow == 0 {
																	r19 = r18
																	goto block2
																}
																goto block1
															}
															goto block1
														}
														frame.Fail()
														goto block1
													}
													goto block1
												}
												goto block1
											}
											goto block3
										}
										goto block3
									}
									goto block3
								}
								goto block3
							}
							goto block3
						}
						frame.Fail()
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block1:
	frame.Recover(r14)
	r19 = r13
	goto block2
block2:
	EOS(frame)
	if frame.Flow == 0 {
		ret0 = &Assign{Expr: r19, Targets: []ASTExpr{r11}, Type: r12, Define: true}
		goto block9
	}
	goto block3
block3:
	frame.Recover(r0)
	r23 = frame.Peek()
	if frame.Flow == 0 {
		if r23 == 'f' {
			frame.Consume()
			r26 = frame.Peek()
			if frame.Flow == 0 {
				if r26 == 'a' {
					frame.Consume()
					r29 = frame.Peek()
					if frame.Flow == 0 {
						if r29 == 'i' {
							frame.Consume()
							r32 = frame.Peek()
							if frame.Flow == 0 {
								if r32 == 'l' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										EOS(frame)
										if frame.Flow == 0 {
											ret0 = &Fail{}
											goto block9
										}
										goto block4
									}
									goto block4
								}
								frame.Fail()
								goto block4
							}
							goto block4
						}
						frame.Fail()
						goto block4
					}
					goto block4
				}
				frame.Fail()
				goto block4
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(r0)
	r36 = frame.Peek()
	if frame.Flow == 0 {
		if r36 == 'r' {
			frame.Consume()
			r39 = frame.Peek()
			if frame.Flow == 0 {
				if r39 == 'e' {
					frame.Consume()
					r42 = frame.Peek()
					if frame.Flow == 0 {
						if r42 == 't' {
							frame.Consume()
							r45 = frame.Peek()
							if frame.Flow == 0 {
								if r45 == 'u' {
									frame.Consume()
									r48 = frame.Peek()
									if frame.Flow == 0 {
										if r48 == 'r' {
											frame.Consume()
											r51 = frame.Peek()
											if frame.Flow == 0 {
												if r51 == 'n' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r54 = ParseExprList(frame)
															if frame.Flow == 0 {
																EOS(frame)
																if frame.Flow == 0 {
																	ret0 = &Return{Exprs: r54}
																	goto block9
																}
																goto block5
															}
															goto block5
														}
														goto block5
													}
													goto block5
												}
												frame.Fail()
												goto block5
											}
											goto block5
										}
										frame.Fail()
										goto block5
									}
									goto block5
								}
								frame.Fail()
								goto block5
							}
							goto block5
						}
						frame.Fail()
						goto block5
					}
					goto block5
				}
				frame.Fail()
				goto block5
			}
			goto block5
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(r0)
	r56 = ParseTargetList(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r57 = false
			r58 = frame.Checkpoint()
			r59 = frame.Peek()
			if frame.Flow == 0 {
				if r59 == ':' {
					frame.Consume()
					r62 = frame.Peek()
					if frame.Flow == 0 {
						if r62 == '=' {
							frame.Consume()
							r66 = true
							goto block7
						}
						frame.Fail()
						goto block6
					}
					goto block6
				}
				frame.Fail()
				goto block6
			}
			goto block6
		}
		goto block8
	}
	goto block8
block6:
	frame.Recover(r58)
	r67 = frame.Peek()
	if frame.Flow == 0 {
		if r67 == '=' {
			frame.Consume()
			r66 = r57
			goto block7
		}
		frame.Fail()
		goto block8
	}
	goto block8
block7:
	S(frame)
	if frame.Flow == 0 {
		r70 = ParseExpr(frame)
		if frame.Flow == 0 {
			EOS(frame)
			if frame.Flow == 0 {
				ret0 = &Assign{Expr: r70, Targets: r56, Define: r66}
				goto block9
			}
			goto block8
		}
		goto block8
	}
	goto block8
block8:
	frame.Recover(r0)
	r72 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret0 = r72
			goto block9
		}
		goto block10
	}
	goto block10
block9:
	return
block10:
	return
}

func ParseCodeBlock(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 rune
	var r4 []ASTExpr
	var r5 int
	var r6 ASTExpr
	var r7 []ASTExpr
	var r8 []ASTExpr
	var r9 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '{' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r4 = []ASTExpr{}
				goto block1
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block1:
	r5 = frame.Checkpoint()
	r6 = ParseStatement(frame)
	if frame.Flow == 0 {
		r7 = append(r4, r6)
		S(frame)
		if frame.Flow == 0 {
			r4 = r7
			goto block1
		}
		r8 = r7
		goto block2
	}
	r8 = r4
	goto block2
block2:
	frame.Recover(r5)
	r9 = frame.Peek()
	if frame.Flow == 0 {
		if r9 == '}' {
			frame.Consume()
			ret0 = r8
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	return
}

func ParseParenthTypeList(frame *runtime.State) (ret0 []ASTTypeRef) {
	var r0 rune
	var r3 []ASTTypeRef
	var r4 int
	var r5 ASTTypeRef
	var r6 []ASTTypeRef
	var r7 []ASTTypeRef
	var r8 int
	var r9 rune
	var r12 ASTTypeRef
	var r13 []ASTTypeRef
	var r14 []ASTTypeRef
	var r15 []ASTTypeRef
	var r16 []ASTTypeRef
	var r17 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == '(' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r3 = []ASTTypeRef{}
				r4 = frame.Checkpoint()
				r5 = ParseTypeRef(frame)
				if frame.Flow == 0 {
					r6 = append(r3, r5)
					S(frame)
					if frame.Flow == 0 {
						r7 = r6
						goto block1
					}
					r16 = r6
					goto block3
				}
				r16 = r3
				goto block3
			}
			goto block5
		}
		frame.Fail()
		goto block5
	}
	goto block5
block1:
	r8 = frame.Checkpoint()
	r9 = frame.Peek()
	if frame.Flow == 0 {
		if r9 == ',' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r12 = ParseTypeRef(frame)
				if frame.Flow == 0 {
					r13 = append(r7, r12)
					S(frame)
					if frame.Flow == 0 {
						r7 = r13
						goto block1
					}
					r14 = r13
					goto block2
				}
				r14 = r7
				goto block2
			}
			r14 = r7
			goto block2
		}
		frame.Fail()
		r14 = r7
		goto block2
	}
	r14 = r7
	goto block2
block2:
	frame.Recover(r8)
	r15 = r14
	goto block4
block3:
	frame.Recover(r4)
	r15 = r16
	goto block4
block4:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		if r17 == ')' {
			frame.Consume()
			ret0 = r15
			return
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	return
}

func ParseStructDecl(frame *runtime.State) (ret0 *StructDecl) {
	var r0 rune
	var r3 rune
	var r6 rune
	var r9 rune
	var r12 rune
	var r15 rune
	var r18 *Id
	var r19 bool
	var r20 int
	var r21 rune
	var r24 rune
	var r27 rune
	var r30 rune
	var r33 rune
	var r36 rune
	var r40 bool
	var r41 []ASTTypeRef
	var r42 int
	var r43 rune
	var r46 rune
	var r49 rune
	var r52 rune
	var r55 rune
	var r58 rune
	var r61 rune
	var r64 rune
	var r67 []ASTTypeRef
	var r68 []ASTTypeRef
	var r69 []ASTTypeRef
	var r70 ASTTypeRef
	var r71 int
	var r72 rune
	var r75 rune
	var r78 rune
	var r81 rune
	var r84 rune
	var r87 rune
	var r90 rune
	var r93 rune
	var r96 rune
	var r99 rune
	var r102 ASTTypeRef
	var r103 ASTTypeRef
	var r104 ASTTypeRef
	var r105 rune
	var r109 []*FieldDecl
	var r110 int
	var r111 *Id
	var r112 ASTTypeRef
	var r115 rune
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == 's' {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				if r3 == 't' {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						if r6 == 'r' {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								if r9 == 'u' {
									frame.Consume()
									r12 = frame.Peek()
									if frame.Flow == 0 {
										if r12 == 'c' {
											frame.Consume()
											r15 = frame.Peek()
											if frame.Flow == 0 {
												if r15 == 't' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r18 = Ident(frame)
															if frame.Flow == 0 {
																S(frame)
																if frame.Flow == 0 {
																	r19 = false
																	r20 = frame.Checkpoint()
																	r21 = frame.Peek()
																	if frame.Flow == 0 {
																		if r21 == 's' {
																			frame.Consume()
																			r24 = frame.Peek()
																			if frame.Flow == 0 {
																				if r24 == 'c' {
																					frame.Consume()
																					r27 = frame.Peek()
																					if frame.Flow == 0 {
																						if r27 == 'o' {
																							frame.Consume()
																							r30 = frame.Peek()
																							if frame.Flow == 0 {
																								if r30 == 'p' {
																									frame.Consume()
																									r33 = frame.Peek()
																									if frame.Flow == 0 {
																										if r33 == 'e' {
																											frame.Consume()
																											r36 = frame.Peek()
																											if frame.Flow == 0 {
																												if r36 == 'd' {
																													frame.Consume()
																													EndKeyword(frame)
																													if frame.Flow == 0 {
																														S(frame)
																														if frame.Flow == 0 {
																															r40 = true
																															goto block2
																														}
																														goto block1
																													}
																													goto block1
																												}
																												frame.Fail()
																												goto block1
																											}
																											goto block1
																										}
																										frame.Fail()
																										goto block1
																									}
																									goto block1
																								}
																								frame.Fail()
																								goto block1
																							}
																							goto block1
																						}
																						frame.Fail()
																						goto block1
																					}
																					goto block1
																				}
																				frame.Fail()
																				goto block1
																			}
																			goto block1
																		}
																		frame.Fail()
																		goto block1
																	}
																	goto block1
																}
																goto block9
															}
															goto block9
														}
														goto block9
													}
													goto block9
												}
												frame.Fail()
												goto block9
											}
											goto block9
										}
										frame.Fail()
										goto block9
									}
									goto block9
								}
								frame.Fail()
								goto block9
							}
							goto block9
						}
						frame.Fail()
						goto block9
					}
					goto block9
				}
				frame.Fail()
				goto block9
			}
			goto block9
		}
		frame.Fail()
		goto block9
	}
	goto block9
block1:
	frame.Recover(r20)
	r40 = r19
	goto block2
block2:
	r41 = []ASTTypeRef{}
	r42 = frame.Checkpoint()
	r43 = frame.Peek()
	if frame.Flow == 0 {
		if r43 == 'c' {
			frame.Consume()
			r46 = frame.Peek()
			if frame.Flow == 0 {
				if r46 == 'o' {
					frame.Consume()
					r49 = frame.Peek()
					if frame.Flow == 0 {
						if r49 == 'n' {
							frame.Consume()
							r52 = frame.Peek()
							if frame.Flow == 0 {
								if r52 == 't' {
									frame.Consume()
									r55 = frame.Peek()
									if frame.Flow == 0 {
										if r55 == 'a' {
											frame.Consume()
											r58 = frame.Peek()
											if frame.Flow == 0 {
												if r58 == 'i' {
													frame.Consume()
													r61 = frame.Peek()
													if frame.Flow == 0 {
														if r61 == 'n' {
															frame.Consume()
															r64 = frame.Peek()
															if frame.Flow == 0 {
																if r64 == 's' {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		S(frame)
																		if frame.Flow == 0 {
																			r67 = ParseParenthTypeList(frame)
																			if frame.Flow == 0 {
																				S(frame)
																				if frame.Flow == 0 {
																					r68 = r67
																					goto block4
																				}
																				r69 = r67
																				goto block3
																			}
																			r69 = r41
																			goto block3
																		}
																		r69 = r41
																		goto block3
																	}
																	r69 = r41
																	goto block3
																}
																frame.Fail()
																r69 = r41
																goto block3
															}
															r69 = r41
															goto block3
														}
														frame.Fail()
														r69 = r41
														goto block3
													}
													r69 = r41
													goto block3
												}
												frame.Fail()
												r69 = r41
												goto block3
											}
											r69 = r41
											goto block3
										}
										frame.Fail()
										r69 = r41
										goto block3
									}
									r69 = r41
									goto block3
								}
								frame.Fail()
								r69 = r41
								goto block3
							}
							r69 = r41
							goto block3
						}
						frame.Fail()
						r69 = r41
						goto block3
					}
					r69 = r41
					goto block3
				}
				frame.Fail()
				r69 = r41
				goto block3
			}
			r69 = r41
			goto block3
		}
		frame.Fail()
		r69 = r41
		goto block3
	}
	r69 = r41
	goto block3
block3:
	frame.Recover(r42)
	r68 = r69
	goto block4
block4:
	r70 = nil
	r71 = frame.Checkpoint()
	r72 = frame.Peek()
	if frame.Flow == 0 {
		if r72 == 'i' {
			frame.Consume()
			r75 = frame.Peek()
			if frame.Flow == 0 {
				if r75 == 'm' {
					frame.Consume()
					r78 = frame.Peek()
					if frame.Flow == 0 {
						if r78 == 'p' {
							frame.Consume()
							r81 = frame.Peek()
							if frame.Flow == 0 {
								if r81 == 'l' {
									frame.Consume()
									r84 = frame.Peek()
									if frame.Flow == 0 {
										if r84 == 'e' {
											frame.Consume()
											r87 = frame.Peek()
											if frame.Flow == 0 {
												if r87 == 'm' {
													frame.Consume()
													r90 = frame.Peek()
													if frame.Flow == 0 {
														if r90 == 'e' {
															frame.Consume()
															r93 = frame.Peek()
															if frame.Flow == 0 {
																if r93 == 'n' {
																	frame.Consume()
																	r96 = frame.Peek()
																	if frame.Flow == 0 {
																		if r96 == 't' {
																			frame.Consume()
																			r99 = frame.Peek()
																			if frame.Flow == 0 {
																				if r99 == 's' {
																					frame.Consume()
																					EndKeyword(frame)
																					if frame.Flow == 0 {
																						S(frame)
																						if frame.Flow == 0 {
																							r102 = ParseTypeRef(frame)
																							if frame.Flow == 0 {
																								S(frame)
																								if frame.Flow == 0 {
																									r103 = r102
																									goto block6
																								}
																								r104 = r102
																								goto block5
																							}
																							r104 = r70
																							goto block5
																						}
																						r104 = r70
																						goto block5
																					}
																					r104 = r70
																					goto block5
																				}
																				frame.Fail()
																				r104 = r70
																				goto block5
																			}
																			r104 = r70
																			goto block5
																		}
																		frame.Fail()
																		r104 = r70
																		goto block5
																	}
																	r104 = r70
																	goto block5
																}
																frame.Fail()
																r104 = r70
																goto block5
															}
															r104 = r70
															goto block5
														}
														frame.Fail()
														r104 = r70
														goto block5
													}
													r104 = r70
													goto block5
												}
												frame.Fail()
												r104 = r70
												goto block5
											}
											r104 = r70
											goto block5
										}
										frame.Fail()
										r104 = r70
										goto block5
									}
									r104 = r70
									goto block5
								}
								frame.Fail()
								r104 = r70
								goto block5
							}
							r104 = r70
							goto block5
						}
						frame.Fail()
						r104 = r70
						goto block5
					}
					r104 = r70
					goto block5
				}
				frame.Fail()
				r104 = r70
				goto block5
			}
			r104 = r70
			goto block5
		}
		frame.Fail()
		r104 = r70
		goto block5
	}
	r104 = r70
	goto block5
block5:
	frame.Recover(r71)
	r103 = r104
	goto block6
block6:
	r105 = frame.Peek()
	if frame.Flow == 0 {
		if r105 == '{' {
			frame.Consume()
			S(frame)
			if frame.Flow == 0 {
				r109 = []*FieldDecl{}
				goto block7
			}
			goto block9
		}
		frame.Fail()
		goto block9
	}
	goto block9
block7:
	r110 = frame.Checkpoint()
	r111 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r112 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				S(frame)
				if frame.Flow == 0 {
					r109 = append(r109, &FieldDecl{Name: r111, Type: r112})
					goto block7
				}
				goto block8
			}
			goto block8
		}
		goto block8
	}
	goto block8
block8:
	frame.Recover(r110)
	r115 = frame.Peek()
	if frame.Flow == 0 {
		if r115 == '}' {
			frame.Consume()
			ret0 = &StructDecl{Name: r18, Implements: r103, Fields: r109, Scoped: r40, Contains: r68}
			return
		}
		frame.Fail()
		goto block9
	}
	goto block9
block9:
	return
}

func ParseParam(frame *runtime.State) (ret0 *Param) {
	var r0 *NameRef
	var r1 ASTTypeRef
	r0 = ParseNameRef(frame)
	if frame.Flow == 0 {
		S(frame)
		if frame.Flow == 0 {
			r1 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				ret0 = &Param{Name: r0, Type: r1}
				return
			}
			goto block1
		}
		goto block1
	}
	goto block1
block1:
	return
}

func ParseParamList(frame *runtime.State) (ret0 []*Param) {
	var r0 []*Param
	var r1 int
	var r2 *Param
	var r4 []*Param
	var r5 int
	var r6 rune
	var r9 *Param
	var r11 []*Param
	r0 = []*Param{}
	r1 = frame.Checkpoint()
	r2 = ParseParam(frame)
	if frame.Flow == 0 {
		r4 = append(r0, r2)
		goto block1
	}
	frame.Recover(r1)
	r11 = r0
	goto block3
block1:
	r5 = frame.Checkpoint()
	S(frame)
	if frame.Flow == 0 {
		r6 = frame.Peek()
		if frame.Flow == 0 {
			if r6 == ',' {
				frame.Consume()
				S(frame)
				if frame.Flow == 0 {
					r9 = ParseParam(frame)
					if frame.Flow == 0 {
						r4 = append(r4, r9)
						goto block1
					}
					goto block2
				}
				goto block2
			}
			frame.Fail()
			goto block2
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(r5)
	r11 = r4
	goto block3
block3:
	ret0 = r11
	return
}

func ParseFuncDecl(frame *runtime.State) (ret0 *FuncDecl) {
	var r0 rune
	var r3 rune
	var r6 rune
	var r9 rune
	var r12 *Id
	var r13 rune
	var r16 []*Param
	var r17 rune
	var r20 []ASTTypeRef
	var r21 []ASTExpr
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == 'f' {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				if r3 == 'u' {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						if r6 == 'n' {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								if r9 == 'c' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										if frame.Flow == 0 {
											r12 = Ident(frame)
											if frame.Flow == 0 {
												S(frame)
												if frame.Flow == 0 {
													r13 = frame.Peek()
													if frame.Flow == 0 {
														if r13 == '(' {
															frame.Consume()
															S(frame)
															if frame.Flow == 0 {
																r16 = ParseParamList(frame)
																if frame.Flow == 0 {
																	S(frame)
																	if frame.Flow == 0 {
																		r17 = frame.Peek()
																		if frame.Flow == 0 {
																			if r17 == ')' {
																				frame.Consume()
																				S(frame)
																				if frame.Flow == 0 {
																					r20 = ParseReturnTypeList(frame)
																					if frame.Flow == 0 {
																						S(frame)
																						if frame.Flow == 0 {
																							r21 = ParseCodeBlock(frame)
																							if frame.Flow == 0 {
																								ret0 = &FuncDecl{Name: r12, Params: r16, ReturnTypes: r20, Block: r21, LocalInfo_Scope: &LocalInfo_Scope{}}
																								return
																							}
																							goto block1
																						}
																						goto block1
																					}
																					goto block1
																				}
																				goto block1
																			}
																			frame.Fail()
																			goto block1
																		}
																		goto block1
																	}
																	goto block1
																}
																goto block1
															}
															goto block1
														}
														frame.Fail()
														goto block1
													}
													goto block1
												}
												goto block1
											}
											goto block1
										}
										goto block1
									}
									goto block1
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	return
}

func ParseMatchState(frame *runtime.State) (ret0 string) {
	var r0 int
	var r1 int
	var r2 int
	var r3 rune
	var r6 rune
	var r9 rune
	var r12 rune
	var r15 rune
	var r18 rune
	var r21 rune
	var r24 rune
	var r27 rune
	var r30 rune
	var r33 string
	r0 = frame.Checkpoint()
	r1 = frame.Checkpoint()
	r2 = frame.Checkpoint()
	r3 = frame.Peek()
	if frame.Flow == 0 {
		if r3 == 'N' {
			frame.Consume()
			r6 = frame.Peek()
			if frame.Flow == 0 {
				if r6 == 'O' {
					frame.Consume()
					r9 = frame.Peek()
					if frame.Flow == 0 {
						if r9 == 'R' {
							frame.Consume()
							r12 = frame.Peek()
							if frame.Flow == 0 {
								if r12 == 'M' {
									frame.Consume()
									r15 = frame.Peek()
									if frame.Flow == 0 {
										if r15 == 'A' {
											frame.Consume()
											r18 = frame.Peek()
											if frame.Flow == 0 {
												if r18 == 'L' {
													frame.Consume()
													goto block2
												}
												frame.Fail()
												goto block1
											}
											goto block1
										}
										frame.Fail()
										goto block1
									}
									goto block1
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(r2)
	r21 = frame.Peek()
	if frame.Flow == 0 {
		if r21 == 'F' {
			frame.Consume()
			r24 = frame.Peek()
			if frame.Flow == 0 {
				if r24 == 'A' {
					frame.Consume()
					r27 = frame.Peek()
					if frame.Flow == 0 {
						if r27 == 'I' {
							frame.Consume()
							r30 = frame.Peek()
							if frame.Flow == 0 {
								if r30 == 'L' {
									frame.Consume()
									goto block2
								}
								frame.Fail()
								goto block3
							}
							goto block3
						}
						frame.Fail()
						goto block3
					}
					goto block3
				}
				frame.Fail()
				goto block3
			}
			goto block3
		}
		frame.Fail()
		goto block3
	}
	goto block3
block2:
	r33 = frame.Slice(r1)
	EndKeyword(frame)
	if frame.Flow == 0 {
		ret0 = r33
		goto block4
	}
	goto block3
block3:
	frame.Recover(r0)
	ret0 = "NORMAL"
	goto block4
block4:
	return
}

func ParseTest(frame *runtime.State) (ret0 *Test) {
	var r0 rune
	var r3 rune
	var r6 rune
	var r9 rune
	var r12 *Id
	var r13 ASTExpr
	var r14 string
	var r15 string
	var r16 Destructure
	r0 = frame.Peek()
	if frame.Flow == 0 {
		if r0 == 't' {
			frame.Consume()
			r3 = frame.Peek()
			if frame.Flow == 0 {
				if r3 == 'e' {
					frame.Consume()
					r6 = frame.Peek()
					if frame.Flow == 0 {
						if r6 == 's' {
							frame.Consume()
							r9 = frame.Peek()
							if frame.Flow == 0 {
								if r9 == 't' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										if frame.Flow == 0 {
											r12 = Ident(frame)
											if frame.Flow == 0 {
												S(frame)
												if frame.Flow == 0 {
													r13 = ParseExpr(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r14 = DecodeString(frame)
															if frame.Flow == 0 {
																S(frame)
																if frame.Flow == 0 {
																	r15 = ParseMatchState(frame)
																	if frame.Flow == 0 {
																		S(frame)
																		if frame.Flow == 0 {
																			r16 = ParseDestructure(frame)
																			if frame.Flow == 0 {
																				ret0 = &Test{Name: r12, Rule: r13, Input: r14, Flow: r15, Destructure: r16}
																				return
																			}
																			goto block1
																		}
																		goto block1
																	}
																	goto block1
																}
																goto block1
															}
															goto block1
														}
														goto block1
													}
													goto block1
												}
												goto block1
											}
											goto block1
										}
										goto block1
									}
									goto block1
								}
								frame.Fail()
								goto block1
							}
							goto block1
						}
						frame.Fail()
						goto block1
					}
					goto block1
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	return
}

func ParseImports(frame *runtime.State) (ret0 []*ImportDecl) {
	var r0 []*ImportDecl
	var r1 int
	var r2 rune
	var r5 rune
	var r8 rune
	var r11 rune
	var r14 rune
	var r17 rune
	var r20 rune
	var r23 []*ImportDecl
	var r24 int
	var r25 *StringLiteral
	var r27 []*ImportDecl
	var r28 []*ImportDecl
	var r29 rune
	var r32 []*ImportDecl
	var r33 []*ImportDecl
	r0 = []*ImportDecl{}
	r1 = frame.Checkpoint()
	r2 = frame.Peek()
	if frame.Flow == 0 {
		if r2 == 'i' {
			frame.Consume()
			r5 = frame.Peek()
			if frame.Flow == 0 {
				if r5 == 'm' {
					frame.Consume()
					r8 = frame.Peek()
					if frame.Flow == 0 {
						if r8 == 'p' {
							frame.Consume()
							r11 = frame.Peek()
							if frame.Flow == 0 {
								if r11 == 'o' {
									frame.Consume()
									r14 = frame.Peek()
									if frame.Flow == 0 {
										if r14 == 'r' {
											frame.Consume()
											r17 = frame.Peek()
											if frame.Flow == 0 {
												if r17 == 't' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														if frame.Flow == 0 {
															r20 = frame.Peek()
															if frame.Flow == 0 {
																if r20 == '(' {
																	frame.Consume()
																	S(frame)
																	if frame.Flow == 0 {
																		r23 = r0
																		goto block1
																	}
																	r33 = r0
																	goto block3
																}
																frame.Fail()
																r33 = r0
																goto block3
															}
															r33 = r0
															goto block3
														}
														r33 = r0
														goto block3
													}
													r33 = r0
													goto block3
												}
												frame.Fail()
												r33 = r0
												goto block3
											}
											r33 = r0
											goto block3
										}
										frame.Fail()
										r33 = r0
										goto block3
									}
									r33 = r0
									goto block3
								}
								frame.Fail()
								r33 = r0
								goto block3
							}
							r33 = r0
							goto block3
						}
						frame.Fail()
						r33 = r0
						goto block3
					}
					r33 = r0
					goto block3
				}
				frame.Fail()
				r33 = r0
				goto block3
			}
			r33 = r0
			goto block3
		}
		frame.Fail()
		r33 = r0
		goto block3
	}
	r33 = r0
	goto block3
block1:
	r24 = frame.Checkpoint()
	r25 = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		r27 = append(r23, &ImportDecl{Path: r25})
		S(frame)
		if frame.Flow == 0 {
			r23 = r27
			goto block1
		}
		r28 = r27
		goto block2
	}
	r28 = r23
	goto block2
block2:
	frame.Recover(r24)
	r29 = frame.Peek()
	if frame.Flow == 0 {
		if r29 == ')' {
			frame.Consume()
			r32 = r28
			goto block4
		}
		frame.Fail()
		r33 = r28
		goto block3
	}
	r33 = r28
	goto block3
block3:
	frame.Recover(r1)
	r32 = r33
	goto block4
block4:
	ret0 = r32
	return
}

func ParseFile(frame *runtime.State) (ret0 *File) {
	var r0 []ASTDecl
	var r1 []*Test
	var r2 []*ImportDecl
	var r3 []ASTDecl
	var r4 []*Test
	var r5 int
	var r6 int
	var r7 *FuncDecl
	var r9 []ASTDecl
	var r10 []*Test
	var r11 *StructDecl
	var r13 *Test
	var r15 []ASTDecl
	var r16 []*Test
	var r17 int
	r0 = []ASTDecl{}
	r1 = []*Test{}
	S(frame)
	if frame.Flow == 0 {
		r2 = ParseImports(frame)
		if frame.Flow == 0 {
			S(frame)
			if frame.Flow == 0 {
				r3, r4 = r0, r1
				goto block1
			}
			goto block4
		}
		goto block4
	}
	goto block4
block1:
	r5 = frame.Checkpoint()
	r6 = frame.Checkpoint()
	r7 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		r9, r10 = append(r3, r7), r4
		goto block2
	}
	frame.Recover(r6)
	r11 = ParseStructDecl(frame)
	if frame.Flow == 0 {
		r9, r10 = append(r3, r11), r4
		goto block2
	}
	frame.Recover(r6)
	r13 = ParseTest(frame)
	if frame.Flow == 0 {
		r9, r10 = r3, append(r4, r13)
		goto block2
	}
	r15, r16 = r3, r4
	goto block3
block2:
	S(frame)
	if frame.Flow == 0 {
		r3, r4 = r9, r10
		goto block1
	}
	r15, r16 = r9, r10
	goto block3
block3:
	frame.Recover(r5)
	r17 = frame.LookaheadBegin()
	frame.Peek()
	if frame.Flow == 0 {
		frame.Consume()
		frame.LookaheadFail(r17)
		goto block4
	}
	frame.LookaheadNormal(r17)
	ret0 = &File{Imports: r2, Decls: r15, Tests: r16}
	return
block4:
	return
}
