package tree

import (
	"evergreen/dub/core"
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
	Pos   int
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

type Float32Literal struct {
	Text  string
	Value float32
}

func (node *Float32Literal) isASTExpr() {
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
	T    core.DubType
}

func (node *TypeRef) isASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    core.DubType
}

func (node *ListTypeRef) isASTTypeRef() {
}

type QualifiedTypeRef struct {
	Package *Id
	Name    *Id
	T       core.DubType
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
	Else  []ASTExpr
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
	Local *LocalInfo
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
	Target core.Callable
	T      []core.DubType
}

func (node *Call) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
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
	T     core.DubType
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
	T          *core.StructType
}

func (node *StructDecl) isASTDecl() {
}

type LocalInfo_Ref uint32

type LocalInfo_Scope struct {
	objects []*LocalInfo
}

type LocalInfo struct {
	Name  string
	T     core.DubType
	Index LocalInfo_Ref
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
	F               *core.Function
	LocalInfo_Scope *LocalInfo_Scope
}

func (node *FuncDecl) isASTDecl() {
}

type Test struct {
	Name        *Id
	Rule        ASTExpr
	Type        core.DubType
	Input       string
	Flow        string
	Destructure Destructure
}

type ImportDecl struct {
	Path *StringLiteral
}

type File struct {
	Name    string
	Imports []*ImportDecl
	Decls   []ASTDecl
	Tests   []*Test
	F       *core.File
}

type Package struct {
	Path  []string
	Files []*File
	P     *core.Package
}

type Program struct {
	Builtins *core.BuiltinTypeIndex
	Packages []*Package
}

func LineTerminator(frame *runtime.State) {
	var checkpoint int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	checkpoint = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '\n' {
			frame.Consume()
			goto block3
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(checkpoint)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '\r' {
			frame.Consume()
			c2 = frame.Peek()
			if frame.Flow == 0 {
				if c2 == '\n' {
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
	frame.Recover(checkpoint)
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == '\r' {
			frame.Consume()
			goto block3
		}
		frame.Fail()
		return
	}
	return
block3:
	return
}

func S(frame *runtime.State) {
	var checkpoint0 int
	var checkpoint1 int
	var c0 rune
	var c1 rune
	var c2 rune
	var checkpoint2 int
	var c3 rune
	goto block1
block1:
	checkpoint0 = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == ' ' {
			goto block2
		}
		if c0 == '\t' {
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
	frame.Recover(checkpoint1)
	LineTerminator(frame)
	if frame.Flow == 0 {
		goto block1
	}
	frame.Recover(checkpoint1)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '/' {
			frame.Consume()
			c2 = frame.Peek()
			if frame.Flow == 0 {
				if c2 == '/' {
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
	checkpoint2 = frame.Checkpoint()
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == '\n' {
			goto block5
		}
		if c3 == '\r' {
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
	frame.Recover(checkpoint2)
	goto block1
block7:
	frame.Recover(checkpoint0)
	return
}

func EndKeyword(frame *runtime.State) {
	var checkpoint int
	var c rune
	checkpoint = frame.LookaheadBegin()
	c = frame.Peek()
	if frame.Flow == 0 {
		if c >= 'a' {
			if c <= 'z' {
				goto block3
			}
			goto block1
		}
		goto block1
	}
	goto block5
block1:
	if c >= 'A' {
		if c <= 'Z' {
			goto block3
		}
		goto block2
	}
	goto block2
block2:
	if c == '_' {
		goto block3
	}
	if c >= '0' {
		if c <= '9' {
			goto block3
		}
		goto block4
	}
	goto block4
block3:
	frame.Consume()
	frame.LookaheadFail(checkpoint)
	return
block4:
	frame.Fail()
	goto block5
block5:
	frame.LookaheadNormal(checkpoint)
	return
}

func Ident(frame *runtime.State) (ret *Id) {
	var r int
	var checkpoint0 int
	var checkpoint1 int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var c12 rune
	var c13 rune
	var c14 rune
	var c15 rune
	var c16 rune
	var c17 rune
	var c18 rune
	var c19 rune
	var c20 rune
	var c21 rune
	var c22 rune
	var c23 rune
	var c24 rune
	var c25 rune
	var c26 rune
	var c27 rune
	var c28 rune
	var c29 rune
	var c30 rune
	var c31 rune
	var c32 rune
	var c33 rune
	var c34 rune
	var c35 rune
	var c36 rune
	var c37 rune
	var c38 rune
	var c39 rune
	var c40 rune
	var c41 rune
	var c42 rune
	var c43 rune
	var c44 rune
	var c45 rune
	var c46 rune
	var c47 rune
	var c48 rune
	var c49 rune
	var c50 rune
	var c51 rune
	var c52 rune
	var c53 rune
	var c54 rune
	var c55 rune
	var c56 rune
	var c57 rune
	var c58 rune
	var c59 rune
	var c60 rune
	var c61 rune
	var c62 rune
	var c63 rune
	var c64 rune
	var c65 rune
	var c66 rune
	var c67 rune
	var c68 rune
	var c69 rune
	var c70 rune
	var c71 rune
	var c72 rune
	var c73 rune
	var c74 rune
	var checkpoint2 int
	var c75 rune
	var begin int
	var c76 rune
	var checkpoint3 int
	var c77 rune
	r = frame.Checkpoint()
	checkpoint0 = frame.LookaheadBegin()
	checkpoint1 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'f' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'u' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'n' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'c' {
									frame.Consume()
									goto block16
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
	frame.Recover(checkpoint1)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == 't' {
			frame.Consume()
			c5 = frame.Peek()
			if frame.Flow == 0 {
				if c5 == 'e' {
					frame.Consume()
					c6 = frame.Peek()
					if frame.Flow == 0 {
						if c6 == 's' {
							frame.Consume()
							c7 = frame.Peek()
							if frame.Flow == 0 {
								if c7 == 't' {
									frame.Consume()
									goto block16
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
	frame.Recover(checkpoint1)
	c8 = frame.Peek()
	if frame.Flow == 0 {
		if c8 == 's' {
			frame.Consume()
			c9 = frame.Peek()
			if frame.Flow == 0 {
				if c9 == 't' {
					frame.Consume()
					c10 = frame.Peek()
					if frame.Flow == 0 {
						if c10 == 'r' {
							frame.Consume()
							c11 = frame.Peek()
							if frame.Flow == 0 {
								if c11 == 'u' {
									frame.Consume()
									c12 = frame.Peek()
									if frame.Flow == 0 {
										if c12 == 'c' {
											frame.Consume()
											c13 = frame.Peek()
											if frame.Flow == 0 {
												if c13 == 't' {
													frame.Consume()
													goto block16
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
	frame.Recover(checkpoint1)
	c14 = frame.Peek()
	if frame.Flow == 0 {
		if c14 == 'i' {
			frame.Consume()
			c15 = frame.Peek()
			if frame.Flow == 0 {
				if c15 == 'm' {
					frame.Consume()
					c16 = frame.Peek()
					if frame.Flow == 0 {
						if c16 == 'p' {
							frame.Consume()
							c17 = frame.Peek()
							if frame.Flow == 0 {
								if c17 == 'l' {
									frame.Consume()
									c18 = frame.Peek()
									if frame.Flow == 0 {
										if c18 == 'e' {
											frame.Consume()
											c19 = frame.Peek()
											if frame.Flow == 0 {
												if c19 == 'm' {
													frame.Consume()
													c20 = frame.Peek()
													if frame.Flow == 0 {
														if c20 == 'e' {
															frame.Consume()
															c21 = frame.Peek()
															if frame.Flow == 0 {
																if c21 == 'n' {
																	frame.Consume()
																	c22 = frame.Peek()
																	if frame.Flow == 0 {
																		if c22 == 't' {
																			frame.Consume()
																			c23 = frame.Peek()
																			if frame.Flow == 0 {
																				if c23 == 's' {
																					frame.Consume()
																					goto block16
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
	frame.Recover(checkpoint1)
	c24 = frame.Peek()
	if frame.Flow == 0 {
		if c24 == 's' {
			frame.Consume()
			c25 = frame.Peek()
			if frame.Flow == 0 {
				if c25 == 't' {
					frame.Consume()
					c26 = frame.Peek()
					if frame.Flow == 0 {
						if c26 == 'a' {
							frame.Consume()
							c27 = frame.Peek()
							if frame.Flow == 0 {
								if c27 == 'r' {
									frame.Consume()
									goto block16
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
	frame.Recover(checkpoint1)
	c28 = frame.Peek()
	if frame.Flow == 0 {
		if c28 == 'p' {
			frame.Consume()
			c29 = frame.Peek()
			if frame.Flow == 0 {
				if c29 == 'l' {
					frame.Consume()
					c30 = frame.Peek()
					if frame.Flow == 0 {
						if c30 == 'u' {
							frame.Consume()
							c31 = frame.Peek()
							if frame.Flow == 0 {
								if c31 == 's' {
									frame.Consume()
									goto block16
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
	frame.Recover(checkpoint1)
	c32 = frame.Peek()
	if frame.Flow == 0 {
		if c32 == 'c' {
			frame.Consume()
			c33 = frame.Peek()
			if frame.Flow == 0 {
				if c33 == 'h' {
					frame.Consume()
					c34 = frame.Peek()
					if frame.Flow == 0 {
						if c34 == 'o' {
							frame.Consume()
							c35 = frame.Peek()
							if frame.Flow == 0 {
								if c35 == 'o' {
									frame.Consume()
									c36 = frame.Peek()
									if frame.Flow == 0 {
										if c36 == 's' {
											frame.Consume()
											c37 = frame.Peek()
											if frame.Flow == 0 {
												if c37 == 'e' {
													frame.Consume()
													goto block16
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
	frame.Recover(checkpoint1)
	c38 = frame.Peek()
	if frame.Flow == 0 {
		if c38 == 'o' {
			frame.Consume()
			c39 = frame.Peek()
			if frame.Flow == 0 {
				if c39 == 'r' {
					frame.Consume()
					goto block16
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
	frame.Recover(checkpoint1)
	c40 = frame.Peek()
	if frame.Flow == 0 {
		if c40 == 'q' {
			frame.Consume()
			c41 = frame.Peek()
			if frame.Flow == 0 {
				if c41 == 'u' {
					frame.Consume()
					c42 = frame.Peek()
					if frame.Flow == 0 {
						if c42 == 'e' {
							frame.Consume()
							c43 = frame.Peek()
							if frame.Flow == 0 {
								if c43 == 's' {
									frame.Consume()
									c44 = frame.Peek()
									if frame.Flow == 0 {
										if c44 == 't' {
											frame.Consume()
											c45 = frame.Peek()
											if frame.Flow == 0 {
												if c45 == 'i' {
													frame.Consume()
													c46 = frame.Peek()
													if frame.Flow == 0 {
														if c46 == 'o' {
															frame.Consume()
															c47 = frame.Peek()
															if frame.Flow == 0 {
																if c47 == 'n' {
																	frame.Consume()
																	goto block16
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
	frame.Recover(checkpoint1)
	c48 = frame.Peek()
	if frame.Flow == 0 {
		if c48 == 'i' {
			frame.Consume()
			c49 = frame.Peek()
			if frame.Flow == 0 {
				if c49 == 'f' {
					frame.Consume()
					goto block16
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
	frame.Recover(checkpoint1)
	c50 = frame.Peek()
	if frame.Flow == 0 {
		if c50 == 'e' {
			frame.Consume()
			c51 = frame.Peek()
			if frame.Flow == 0 {
				if c51 == 'l' {
					frame.Consume()
					c52 = frame.Peek()
					if frame.Flow == 0 {
						if c52 == 's' {
							frame.Consume()
							c53 = frame.Peek()
							if frame.Flow == 0 {
								if c53 == 'e' {
									frame.Consume()
									goto block16
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
	frame.Recover(checkpoint1)
	c54 = frame.Peek()
	if frame.Flow == 0 {
		if c54 == 'r' {
			frame.Consume()
			c55 = frame.Peek()
			if frame.Flow == 0 {
				if c55 == 'e' {
					frame.Consume()
					c56 = frame.Peek()
					if frame.Flow == 0 {
						if c56 == 't' {
							frame.Consume()
							c57 = frame.Peek()
							if frame.Flow == 0 {
								if c57 == 'u' {
									frame.Consume()
									c58 = frame.Peek()
									if frame.Flow == 0 {
										if c58 == 'r' {
											frame.Consume()
											c59 = frame.Peek()
											if frame.Flow == 0 {
												if c59 == 'n' {
													frame.Consume()
													goto block16
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
		}
		frame.Fail()
		goto block12
	}
	goto block12
block12:
	frame.Recover(checkpoint1)
	c60 = frame.Peek()
	if frame.Flow == 0 {
		if c60 == 'v' {
			frame.Consume()
			c61 = frame.Peek()
			if frame.Flow == 0 {
				if c61 == 'a' {
					frame.Consume()
					c62 = frame.Peek()
					if frame.Flow == 0 {
						if c62 == 'r' {
							frame.Consume()
							goto block16
						}
						frame.Fail()
						goto block13
					}
					goto block13
				}
				frame.Fail()
				goto block13
			}
			goto block13
		}
		frame.Fail()
		goto block13
	}
	goto block13
block13:
	frame.Recover(checkpoint1)
	c63 = frame.Peek()
	if frame.Flow == 0 {
		if c63 == 't' {
			frame.Consume()
			c64 = frame.Peek()
			if frame.Flow == 0 {
				if c64 == 'r' {
					frame.Consume()
					c65 = frame.Peek()
					if frame.Flow == 0 {
						if c65 == 'u' {
							frame.Consume()
							c66 = frame.Peek()
							if frame.Flow == 0 {
								if c66 == 'e' {
									frame.Consume()
									goto block16
								}
								frame.Fail()
								goto block14
							}
							goto block14
						}
						frame.Fail()
						goto block14
					}
					goto block14
				}
				frame.Fail()
				goto block14
			}
			goto block14
		}
		frame.Fail()
		goto block14
	}
	goto block14
block14:
	frame.Recover(checkpoint1)
	c67 = frame.Peek()
	if frame.Flow == 0 {
		if c67 == 'f' {
			frame.Consume()
			c68 = frame.Peek()
			if frame.Flow == 0 {
				if c68 == 'a' {
					frame.Consume()
					c69 = frame.Peek()
					if frame.Flow == 0 {
						if c69 == 'l' {
							frame.Consume()
							c70 = frame.Peek()
							if frame.Flow == 0 {
								if c70 == 's' {
									frame.Consume()
									c71 = frame.Peek()
									if frame.Flow == 0 {
										if c71 == 'e' {
											frame.Consume()
											goto block16
										}
										frame.Fail()
										goto block15
									}
									goto block15
								}
								frame.Fail()
								goto block15
							}
							goto block15
						}
						frame.Fail()
						goto block15
					}
					goto block15
				}
				frame.Fail()
				goto block15
			}
			goto block15
		}
		frame.Fail()
		goto block15
	}
	goto block15
block15:
	frame.Recover(checkpoint1)
	c72 = frame.Peek()
	if frame.Flow == 0 {
		if c72 == 'n' {
			frame.Consume()
			c73 = frame.Peek()
			if frame.Flow == 0 {
				if c73 == 'i' {
					frame.Consume()
					c74 = frame.Peek()
					if frame.Flow == 0 {
						if c74 == 'l' {
							frame.Consume()
							goto block16
						}
						frame.Fail()
						goto block22
					}
					goto block22
				}
				frame.Fail()
				goto block22
			}
			goto block22
		}
		frame.Fail()
		goto block22
	}
	goto block22
block16:
	checkpoint2 = frame.LookaheadBegin()
	c75 = frame.Peek()
	if frame.Flow == 0 {
		if c75 >= 'a' {
			if c75 <= 'z' {
				goto block19
			}
			goto block17
		}
		goto block17
	}
	goto block21
block17:
	if c75 >= 'A' {
		if c75 <= 'Z' {
			goto block19
		}
		goto block18
	}
	goto block18
block18:
	if c75 == '_' {
		goto block19
	}
	if c75 >= '0' {
		if c75 <= '9' {
			goto block19
		}
		goto block20
	}
	goto block20
block19:
	frame.Consume()
	frame.LookaheadFail(checkpoint2)
	goto block22
block20:
	frame.Fail()
	goto block21
block21:
	frame.LookaheadNormal(checkpoint2)
	frame.LookaheadFail(checkpoint0)
	return
block22:
	frame.LookaheadNormal(checkpoint0)
	begin = frame.Checkpoint()
	c76 = frame.Peek()
	if frame.Flow == 0 {
		if c76 >= 'a' {
			if c76 <= 'z' {
				goto block25
			}
			goto block23
		}
		goto block23
	}
	return
block23:
	if c76 >= 'A' {
		if c76 <= 'Z' {
			goto block25
		}
		goto block24
	}
	goto block24
block24:
	if c76 == '_' {
		goto block25
	}
	frame.Fail()
	return
block25:
	frame.Consume()
	goto block26
block26:
	checkpoint3 = frame.Checkpoint()
	c77 = frame.Peek()
	if frame.Flow == 0 {
		if c77 >= 'a' {
			if c77 <= 'z' {
				goto block29
			}
			goto block27
		}
		goto block27
	}
	goto block31
block27:
	if c77 >= 'A' {
		if c77 <= 'Z' {
			goto block29
		}
		goto block28
	}
	goto block28
block28:
	if c77 == '_' {
		goto block29
	}
	if c77 >= '0' {
		if c77 <= '9' {
			goto block29
		}
		goto block30
	}
	goto block30
block29:
	frame.Consume()
	goto block26
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(checkpoint3)
	ret = &Id{Pos: r, Text: frame.Slice(begin, frame.Checkpoint())}
	return
}

func ParseNumericLiteral(frame *runtime.State) (ret ASTExpr) {
	var value0 int
	var c_i int
	var r0 int
	var c0 rune
	var r1 int
	var value1 int
	var checkpoint0 int
	var c1 rune
	var r2 int
	var checkpoint1 int
	var c2 rune
	var c3 rune
	var r3 int
	var value2 int
	var divisor0 int
	var checkpoint2 int
	var c4 rune
	var r4 int
	var value3 int
	var divisor1 int
	var r5 string
	value0 = 0
	c_i = 1
	r0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 >= '0' {
			if c0 <= '9' {
				frame.Consume()
				r1 = int(c0) - int('0')
				value1 = value0*10 + r1
				goto block1
			}
			goto block10
		}
		goto block10
	}
	return
block1:
	checkpoint0 = frame.Checkpoint()
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 >= '0' {
			if c1 <= '9' {
				frame.Consume()
				r2 = int(c1) - int('0')
				value1 = value1*10 + r2
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
	frame.Recover(checkpoint0)
	checkpoint1 = frame.Checkpoint()
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '.' {
			frame.Consume()
			c3 = frame.Peek()
			if frame.Flow == 0 {
				if c3 >= '0' {
					if c3 <= '9' {
						frame.Consume()
						r3 = int(c3) - int('0')
						value2, divisor0 = value1*10+r3, c_i*10
						goto block4
					}
					goto block7
				}
				goto block7
			}
			goto block8
		}
		frame.Fail()
		goto block8
	}
	goto block8
block4:
	checkpoint2 = frame.Checkpoint()
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 >= '0' {
			if c4 <= '9' {
				frame.Consume()
				r4 = int(c4) - int('0')
				value2, divisor0 = value2*10+r4, divisor0*10
				goto block4
			}
			goto block5
		}
		goto block5
	}
	goto block6
block5:
	frame.Fail()
	goto block6
block6:
	frame.Recover(checkpoint2)
	value3, divisor1 = value2, divisor0
	goto block9
block7:
	frame.Fail()
	goto block8
block8:
	frame.Recover(checkpoint1)
	value3, divisor1 = value1, c_i
	goto block9
block9:
	r5 = frame.Slice(r0, frame.Checkpoint())
	if divisor1 > 1 {
		ret = &Float32Literal{Text: r5, Value: float32(value3) / float32(divisor1)}
		return
	}
	ret = &IntLiteral{Text: r5, Value: value3}
	return
block10:
	frame.Fail()
	return
}

func EscapedChar(frame *runtime.State) (ret rune) {
	var checkpoint int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	checkpoint = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'a' {
			frame.Consume()
			ret = '\a'
			return
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(checkpoint)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == 'b' {
			frame.Consume()
			ret = '\b'
			return
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == 'f' {
			frame.Consume()
			ret = '\f'
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(checkpoint)
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == 'n' {
			frame.Consume()
			ret = '\n'
			return
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(checkpoint)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == 'r' {
			frame.Consume()
			ret = '\r'
			return
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(checkpoint)
	c5 = frame.Peek()
	if frame.Flow == 0 {
		if c5 == 't' {
			frame.Consume()
			ret = '\t'
			return
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(checkpoint)
	c6 = frame.Peek()
	if frame.Flow == 0 {
		if c6 == 'v' {
			frame.Consume()
			ret = '\v'
			return
		}
		frame.Fail()
		goto block7
	}
	goto block7
block7:
	frame.Recover(checkpoint)
	c7 = frame.Peek()
	if frame.Flow == 0 {
		if c7 == '\\' {
			frame.Consume()
			ret = '\\'
			return
		}
		frame.Fail()
		goto block8
	}
	goto block8
block8:
	frame.Recover(checkpoint)
	c8 = frame.Peek()
	if frame.Flow == 0 {
		if c8 == '\'' {
			frame.Consume()
			ret = '\''
			return
		}
		frame.Fail()
		goto block9
	}
	goto block9
block9:
	frame.Recover(checkpoint)
	c9 = frame.Peek()
	if frame.Flow == 0 {
		if c9 == '"' {
			frame.Consume()
			ret = '"'
			return
		}
		frame.Fail()
		return
	}
	return
}

func DecodeString(frame *runtime.State) (ret string) {
	var c0 rune
	var contents []rune
	var checkpoint0 int
	var checkpoint1 int
	var c1 rune
	var c2 rune
	var r rune
	var c3 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '"' {
			frame.Consume()
			contents = []rune{}
			goto block1
		}
		frame.Fail()
		return
	}
	return
block1:
	checkpoint0 = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '"' {
			goto block2
		}
		if c1 == '\\' {
			goto block2
		}
		frame.Consume()
		contents = append(contents, c1)
		goto block1
	}
	goto block3
block2:
	frame.Fail()
	goto block3
block3:
	frame.Recover(checkpoint1)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '\\' {
			frame.Consume()
			r = EscapedChar(frame)
			if frame.Flow == 0 {
				contents = append(contents, r)
				goto block1
			}
			goto block4
		}
		frame.Fail()
		goto block4
	}
	goto block4
block4:
	frame.Recover(checkpoint0)
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == '"' {
			frame.Consume()
			ret = string(contents)
			return
		}
		frame.Fail()
		return
	}
	return
}

func DecodeRune(frame *runtime.State) (ret0 rune, ret1 string) {
	var r0 int
	var c0 rune
	var checkpoint int
	var c1 rune
	var value rune
	var c2 rune
	var r1 rune
	var c3 rune
	r0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '\'' {
			frame.Consume()
			checkpoint = frame.Checkpoint()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == '\\' {
					goto block1
				}
				if c1 == '\'' {
					goto block1
				}
				frame.Consume()
				value = c1
				goto block3
			}
			goto block2
		}
		frame.Fail()
		return
	}
	return
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(checkpoint)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '\\' {
			frame.Consume()
			r1 = EscapedChar(frame)
			if frame.Flow == 0 {
				value = r1
				goto block3
			}
			return
		}
		frame.Fail()
		return
	}
	return
block3:
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == '\'' {
			frame.Consume()
			ret0, ret1 = value, frame.Slice(r0, frame.Checkpoint())
			return
		}
		frame.Fail()
		return
	}
	return
}

func DecodeBool(frame *runtime.State) (ret0 bool, ret1 string) {
	var r int
	var checkpoint int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var value bool
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var c8 rune
	r = frame.Checkpoint()
	checkpoint = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 't' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'r' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'u' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'e' {
									frame.Consume()
									value = true
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
	frame.Recover(checkpoint)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == 'f' {
			frame.Consume()
			c5 = frame.Peek()
			if frame.Flow == 0 {
				if c5 == 'a' {
					frame.Consume()
					c6 = frame.Peek()
					if frame.Flow == 0 {
						if c6 == 'l' {
							frame.Consume()
							c7 = frame.Peek()
							if frame.Flow == 0 {
								if c7 == 's' {
									frame.Consume()
									c8 = frame.Peek()
									if frame.Flow == 0 {
										if c8 == 'e' {
											frame.Consume()
											value = false
											goto block2
										}
										frame.Fail()
										return
									}
									return
								}
								frame.Fail()
								return
							}
							return
						}
						frame.Fail()
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
block2:
	EndKeyword(frame)
	if frame.Flow == 0 {
		ret0, ret1 = value, frame.Slice(r, frame.Checkpoint())
		return
	}
	return
}

func ParseStringLiteral(frame *runtime.State) (ret *StringLiteral) {
	var r0 int
	var r1 string
	r0 = frame.Checkpoint()
	r1 = DecodeString(frame)
	if frame.Flow == 0 {
		ret = &StringLiteral{Pos: r0, Text: frame.Slice(r0, frame.Checkpoint()), Value: r1}
		return
	}
	return
}

func Literal(frame *runtime.State) (ret ASTExpr) {
	var checkpoint int
	var r0 rune
	var r1 string
	var r2 *StringLiteral
	var r3 ASTExpr
	var r4 bool
	var r5 string
	var c0 rune
	var c1 rune
	var c2 rune
	checkpoint = frame.Checkpoint()
	r0, r1 = DecodeRune(frame)
	if frame.Flow == 0 {
		ret = &RuneLiteral{Text: r1, Value: r0}
		return
	}
	frame.Recover(checkpoint)
	r2 = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		ret = r2
		return
	}
	frame.Recover(checkpoint)
	r3 = ParseNumericLiteral(frame)
	if frame.Flow == 0 {
		ret = r3
		return
	}
	frame.Recover(checkpoint)
	r4, r5 = DecodeBool(frame)
	if frame.Flow == 0 {
		ret = &BoolLiteral{Text: r5, Value: r4}
		return
	}
	frame.Recover(checkpoint)
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'n' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'i' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'l' {
							frame.Consume()
							ret = &NilLiteral{}
							return
						}
						frame.Fail()
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func BinaryOperator(frame *runtime.State) (ret0 string, ret1 int) {
	var checkpoint0 int
	var begin0 int
	var c0 rune
	var begin1 int
	var c1 rune
	var begin2 int
	var checkpoint1 int
	var c2 rune
	var checkpoint2 int
	var c3 rune
	var c4 rune
	var c5 rune
	checkpoint0 = frame.Checkpoint()
	begin0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '*' {
			goto block1
		}
		if c0 == '/' {
			goto block1
		}
		if c0 == '%' {
			goto block1
		}
		frame.Fail()
		goto block2
	}
	goto block2
block1:
	frame.Consume()
	ret0, ret1 = frame.Slice(begin0, frame.Checkpoint()), 5
	return
block2:
	frame.Recover(checkpoint0)
	begin1 = frame.Checkpoint()
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '+' {
			goto block3
		}
		if c1 == '-' {
			goto block3
		}
		frame.Fail()
		goto block4
	}
	goto block4
block3:
	frame.Consume()
	ret0, ret1 = frame.Slice(begin1, frame.Checkpoint()), 4
	return
block4:
	frame.Recover(checkpoint0)
	begin2 = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '<' {
			goto block5
		}
		if c2 == '>' {
			goto block5
		}
		frame.Fail()
		goto block7
	}
	goto block7
block5:
	frame.Consume()
	checkpoint2 = frame.Checkpoint()
	c3 = frame.Peek()
	if frame.Flow == 0 {
		if c3 == '=' {
			frame.Consume()
			goto block9
		}
		frame.Fail()
		goto block6
	}
	goto block6
block6:
	frame.Recover(checkpoint2)
	goto block9
block7:
	frame.Recover(checkpoint1)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == '!' {
			goto block8
		}
		if c4 == '=' {
			goto block8
		}
		frame.Fail()
		return
	}
	return
block8:
	frame.Consume()
	c5 = frame.Peek()
	if frame.Flow == 0 {
		if c5 == '=' {
			frame.Consume()
			goto block9
		}
		frame.Fail()
		return
	}
	return
block9:
	ret0, ret1 = frame.Slice(begin2, frame.Checkpoint()), 3
	return
}

func StringMatchExpr(frame *runtime.State) (ret *StringMatch) {
	var c0 rune
	var r TextMatch
	var c1 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '/' {
			frame.Consume()
			S(frame)
			r = ParseMatchChoice(frame)
			if frame.Flow == 0 {
				S(frame)
				c1 = frame.Peek()
				if frame.Flow == 0 {
					if c1 == '/' {
						frame.Consume()
						ret = &StringMatch{Match: r}
						return
					}
					frame.Fail()
					return
				}
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func RuneMatchExpr(frame *runtime.State) (ret *RuneMatch) {
	var c rune
	var r *RuneRangeMatch
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == '$' {
			frame.Consume()
			S(frame)
			r = MatchRune(frame)
			if frame.Flow == 0 {
				ret = &RuneMatch{Match: r}
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseStructTypeRef(frame *runtime.State) (ret ASTTypeRef) {
	var checkpoint int
	var r0 *Id
	var c rune
	var r1 *Id
	var r2 *Id
	checkpoint = frame.Checkpoint()
	r0 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == '.' {
				frame.Consume()
				S(frame)
				r1 = Ident(frame)
				if frame.Flow == 0 {
					ret = &QualifiedTypeRef{Package: r0, Name: r1}
					return
				}
				goto block1
			}
			frame.Fail()
			goto block1
		}
		goto block1
	}
	goto block1
block1:
	frame.Recover(checkpoint)
	r2 = Ident(frame)
	if frame.Flow == 0 {
		ret = &TypeRef{Name: r2}
		return
	}
	return
}

func ParseListTypeRef(frame *runtime.State) (ret *ListTypeRef) {
	var c0 rune
	var c1 rune
	var r ASTTypeRef
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '[' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == ']' {
					frame.Consume()
					r = ParseTypeRef(frame)
					if frame.Flow == 0 {
						ret = &ListTypeRef{Type: r}
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseTypeRef(frame *runtime.State) (ret ASTTypeRef) {
	var checkpoint int
	var r0 ASTTypeRef
	var r1 *ListTypeRef
	checkpoint = frame.Checkpoint()
	r0 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint)
	r1 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		ret = r1
		return
	}
	return
}

func ParseDestructure(frame *runtime.State) (ret Destructure) {
	var checkpoint0 int
	var r0 ASTTypeRef
	var c0 rune
	var fields0 []*DestructureField
	var checkpoint1 int
	var r1 *Id
	var c1 rune
	var r2 Destructure
	var c2 rune
	var r3 *ListTypeRef
	var c3 rune
	var fields1 []Destructure
	var checkpoint2 int
	var r4 Destructure
	var r5 []Destructure
	var fields2 []Destructure
	var c4 rune
	var r6 ASTExpr
	checkpoint0 = frame.Checkpoint()
	r0 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c0 = frame.Peek()
		if frame.Flow == 0 {
			if c0 == '{' {
				frame.Consume()
				S(frame)
				fields0 = []*DestructureField{}
				goto block1
			}
			frame.Fail()
			goto block3
		}
		goto block3
	}
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	r1 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c1 = frame.Peek()
		if frame.Flow == 0 {
			if c1 == ':' {
				frame.Consume()
				S(frame)
				r2 = ParseDestructure(frame)
				if frame.Flow == 0 {
					S(frame)
					fields0 = append(fields0, &DestructureField{Name: r1, Destructure: r2})
					goto block1
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
	frame.Recover(checkpoint1)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '}' {
			frame.Consume()
			ret = &DestructureStruct{Type: r0, Args: fields0}
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(checkpoint0)
	r3 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c3 = frame.Peek()
		if frame.Flow == 0 {
			if c3 == '{' {
				frame.Consume()
				S(frame)
				fields1 = []Destructure{}
				goto block4
			}
			frame.Fail()
			goto block5
		}
		goto block5
	}
	goto block5
block4:
	checkpoint2 = frame.Checkpoint()
	r4 = ParseDestructure(frame)
	if frame.Flow == 0 {
		r5 = append(fields1, r4)
		S(frame)
		fields1 = r5
		goto block4
	}
	fields2 = fields1
	frame.Recover(checkpoint2)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == '}' {
			frame.Consume()
			ret = &DestructureList{Type: r3, Args: fields2}
			return
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(checkpoint0)
	r6 = Literal(frame)
	if frame.Flow == 0 {
		ret = &DestructureValue{Expr: r6}
		return
	}
	return
}

func ParseRuneFilterRune(frame *runtime.State) (ret rune) {
	var checkpoint0 int
	var c0 rune
	var c1 rune
	var checkpoint1 int
	var r rune
	var c2 rune
	checkpoint0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == ']' {
			goto block1
		}
		if c0 == '-' {
			goto block1
		}
		if c0 == '\\' {
			goto block1
		}
		frame.Consume()
		ret = c0
		return
	}
	goto block2
block1:
	frame.Fail()
	goto block2
block2:
	frame.Recover(checkpoint0)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '\\' {
			frame.Consume()
			checkpoint1 = frame.Checkpoint()
			r = EscapedChar(frame)
			if frame.Flow == 0 {
				ret = r
				return
			}
			frame.Recover(checkpoint1)
			c2 = frame.Peek()
			if frame.Flow == 0 {
				frame.Consume()
				ret = c2
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseRuneFilter(frame *runtime.State) (ret *RuneFilter) {
	var r0 rune
	var checkpoint int
	var c rune
	var r1 rune
	var max rune
	r0 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		checkpoint = frame.Checkpoint()
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == '-' {
				frame.Consume()
				r1 = ParseRuneFilterRune(frame)
				if frame.Flow == 0 {
					max = r1
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
	frame.Recover(checkpoint)
	max = r0
	goto block2
block2:
	ret = &RuneFilter{Min: r0, Max: max}
	return
}

func MatchRune(frame *runtime.State) (ret *RuneRangeMatch) {
	var c0 rune
	var c_b bool
	var r0 []*RuneFilter
	var checkpoint0 int
	var c1 rune
	var invert bool
	var filters []*RuneFilter
	var checkpoint1 int
	var r1 *RuneFilter
	var c2 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '[' {
			frame.Consume()
			c_b = false
			r0 = []*RuneFilter{}
			checkpoint0 = frame.Checkpoint()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == '^' {
					frame.Consume()
					invert, filters = true, r0
					goto block2
				}
				frame.Fail()
				goto block1
			}
			goto block1
		}
		frame.Fail()
		return
	}
	return
block1:
	frame.Recover(checkpoint0)
	invert, filters = c_b, r0
	goto block2
block2:
	checkpoint1 = frame.Checkpoint()
	r1 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		filters = append(filters, r1)
		goto block2
	}
	frame.Recover(checkpoint1)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == ']' {
			frame.Consume()
			ret = &RuneRangeMatch{Invert: invert, Filters: filters}
			return
		}
		frame.Fail()
		return
	}
	return
}

func Atom(frame *runtime.State) (ret TextMatch) {
	var checkpoint int
	var r0 *RuneRangeMatch
	var r1 string
	var c0 rune
	var r2 TextMatch
	var c1 rune
	checkpoint = frame.Checkpoint()
	r0 = MatchRune(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint)
	r1 = DecodeString(frame)
	if frame.Flow == 0 {
		ret = &StringLiteralMatch{Value: r1}
		return
	}
	frame.Recover(checkpoint)
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '(' {
			frame.Consume()
			S(frame)
			r2 = ParseMatchChoice(frame)
			if frame.Flow == 0 {
				S(frame)
				c1 = frame.Peek()
				if frame.Flow == 0 {
					if c1 == ')' {
						frame.Consume()
						ret = r2
						return
					}
					frame.Fail()
					return
				}
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func MatchPostfix(frame *runtime.State) (ret TextMatch) {
	var r TextMatch
	var checkpoint int
	var c0 rune
	var c1 rune
	var c2 rune
	r = Atom(frame)
	if frame.Flow == 0 {
		checkpoint = frame.Checkpoint()
		S(frame)
		c0 = frame.Peek()
		if frame.Flow == 0 {
			if c0 == '*' {
				frame.Consume()
				ret = &MatchRepeat{Match: r, Min: 0}
				return
			}
			frame.Fail()
			goto block1
		}
		goto block1
	}
	return
block1:
	frame.Recover(checkpoint)
	S(frame)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '+' {
			frame.Consume()
			ret = &MatchRepeat{Match: r, Min: 1}
			return
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint)
	S(frame)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == '?' {
			frame.Consume()
			ret = &MatchChoice{Matches: []TextMatch{r, &MatchSequence{Matches: []TextMatch{}}}}
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(checkpoint)
	ret = r
	return
}

func MatchPrefix(frame *runtime.State) (ret TextMatch) {
	var checkpoint0 int
	var invert0 bool
	var checkpoint1 int
	var c0 rune
	var invert1 bool
	var c1 rune
	var r0 TextMatch
	var r1 TextMatch
	checkpoint0 = frame.Checkpoint()
	invert0 = false
	checkpoint1 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '!' {
			frame.Consume()
			invert1 = true
			goto block2
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(checkpoint1)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '&' {
			frame.Consume()
			invert1 = invert0
			goto block2
		}
		frame.Fail()
		goto block3
	}
	goto block3
block2:
	S(frame)
	r0 = MatchPostfix(frame)
	if frame.Flow == 0 {
		ret = &MatchLookahead{Invert: invert1, Match: r0}
		return
	}
	goto block3
block3:
	frame.Recover(checkpoint0)
	r1 = MatchPostfix(frame)
	if frame.Flow == 0 {
		ret = r1
		return
	}
	return
}

func Sequence(frame *runtime.State) (ret TextMatch) {
	var r0 TextMatch
	var checkpoint0 int
	var r1 []TextMatch
	var r2 TextMatch
	var l []TextMatch
	var checkpoint1 int
	var r3 TextMatch
	r0 = MatchPrefix(frame)
	if frame.Flow == 0 {
		checkpoint0 = frame.Checkpoint()
		r1 = []TextMatch{r0}
		S(frame)
		r2 = MatchPrefix(frame)
		if frame.Flow == 0 {
			l = append(r1, r2)
			goto block1
		}
		frame.Recover(checkpoint0)
		ret = r0
		return
	}
	return
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	r3 = MatchPrefix(frame)
	if frame.Flow == 0 {
		l = append(l, r3)
		goto block1
	}
	frame.Recover(checkpoint1)
	ret = &MatchSequence{Matches: l}
	return
}

func ParseMatchChoice(frame *runtime.State) (ret TextMatch) {
	var r0 TextMatch
	var checkpoint0 int
	var r1 []TextMatch
	var c0 rune
	var r2 TextMatch
	var l []TextMatch
	var checkpoint1 int
	var c1 rune
	var r3 TextMatch
	r0 = Sequence(frame)
	if frame.Flow == 0 {
		checkpoint0 = frame.Checkpoint()
		r1 = []TextMatch{r0}
		S(frame)
		c0 = frame.Peek()
		if frame.Flow == 0 {
			if c0 == '|' {
				frame.Consume()
				S(frame)
				r2 = Sequence(frame)
				if frame.Flow == 0 {
					l = append(r1, r2)
					goto block1
				}
				goto block3
			}
			frame.Fail()
			goto block3
		}
		goto block3
	}
	return
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '|' {
			frame.Consume()
			S(frame)
			r3 = Sequence(frame)
			if frame.Flow == 0 {
				l = append(l, r3)
				goto block1
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint1)
	ret = &MatchChoice{Matches: l}
	return
block3:
	frame.Recover(checkpoint0)
	ret = r0
	return
}

func ParseExprList(frame *runtime.State) (ret []ASTExpr) {
	var r0 []ASTExpr
	var checkpoint0 int
	var r1 ASTExpr
	var exprs0 []ASTExpr
	var checkpoint1 int
	var c rune
	var r2 ASTExpr
	var exprs1 []ASTExpr
	r0 = []ASTExpr{}
	checkpoint0 = frame.Checkpoint()
	r1 = ParseExpr(frame)
	if frame.Flow == 0 {
		exprs0 = append(r0, r1)
		goto block1
	}
	frame.Recover(checkpoint0)
	exprs1 = r0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r2 = ParseExpr(frame)
			if frame.Flow == 0 {
				exprs0 = append(exprs0, r2)
				goto block1
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint1)
	exprs1 = exprs0
	goto block3
block3:
	ret = exprs1
	return
}

func ParseTargetList(frame *runtime.State) (ret []ASTExpr) {
	var r0 *NameRef
	var exprs []ASTExpr
	var checkpoint int
	var c rune
	var r1 *NameRef
	r0 = ParseNameRef(frame)
	if frame.Flow == 0 {
		exprs = []ASTExpr{r0}
		goto block1
	}
	return
block1:
	checkpoint = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r1 = ParseNameRef(frame)
			if frame.Flow == 0 {
				exprs = append(exprs, r1)
				goto block1
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint)
	ret = exprs
	return
}

func ParseNamedExpr(frame *runtime.State) (ret *NamedExpr) {
	var r0 *Id
	var c rune
	var r1 ASTExpr
	r0 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == ':' {
				frame.Consume()
				S(frame)
				r1 = ParseExpr(frame)
				if frame.Flow == 0 {
					ret = &NamedExpr{Name: r0, Expr: r1}
					return
				}
				return
			}
			frame.Fail()
			return
		}
		return
	}
	return
}

func ParseNamedExprList(frame *runtime.State) (ret []*NamedExpr) {
	var r0 []*NamedExpr
	var checkpoint0 int
	var r1 *NamedExpr
	var exprs0 []*NamedExpr
	var checkpoint1 int
	var c rune
	var r2 *NamedExpr
	var exprs1 []*NamedExpr
	r0 = []*NamedExpr{}
	checkpoint0 = frame.Checkpoint()
	r1 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		exprs0 = append(r0, r1)
		goto block1
	}
	frame.Recover(checkpoint0)
	exprs1 = r0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r2 = ParseNamedExpr(frame)
			if frame.Flow == 0 {
				exprs0 = append(exprs0, r2)
				goto block1
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint1)
	exprs1 = exprs0
	goto block3
block3:
	ret = exprs1
	return
}

func ParseReturnTypeList(frame *runtime.State) (ret []ASTTypeRef) {
	var checkpoint int
	var r0 []ASTTypeRef
	var r1 ASTTypeRef
	checkpoint = frame.Checkpoint()
	r0 = ParseParenthTypeList(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint)
	r1 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		ret = []ASTTypeRef{r1}
		return
	}
	frame.Recover(checkpoint)
	ret = []ASTTypeRef{}
	return
}

func PrimaryExpr(frame *runtime.State) (ret ASTExpr) {
	var checkpoint int
	var r0 ASTExpr
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var r1 ASTTypeRef
	var c7 rune
	var r2 ASTExpr
	var c8 rune
	var r3 *Id
	var c9 rune
	var r4 []ASTExpr
	var c10 rune
	var r5 ASTTypeRef
	var c11 rune
	var r6 []*NamedExpr
	var c12 rune
	var r7 *ListTypeRef
	var c13 rune
	var r8 []ASTExpr
	var c14 rune
	var r9 *StringMatch
	var r10 *RuneMatch
	var c15 rune
	var r11 ASTExpr
	var c16 rune
	var r12 *NameRef
	checkpoint = frame.Checkpoint()
	r0 = Literal(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint)
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'c' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'o' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'e' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'r' {
									frame.Consume()
									c4 = frame.Peek()
									if frame.Flow == 0 {
										if c4 == 'c' {
											frame.Consume()
											c5 = frame.Peek()
											if frame.Flow == 0 {
												if c5 == 'e' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														c6 = frame.Peek()
														if frame.Flow == 0 {
															if c6 == '(' {
																frame.Consume()
																S(frame)
																r1 = ParseTypeRef(frame)
																if frame.Flow == 0 {
																	S(frame)
																	c7 = frame.Peek()
																	if frame.Flow == 0 {
																		if c7 == ',' {
																			frame.Consume()
																			S(frame)
																			r2 = ParseExpr(frame)
																			if frame.Flow == 0 {
																				S(frame)
																				c8 = frame.Peek()
																				if frame.Flow == 0 {
																					if c8 == ')' {
																						frame.Consume()
																						ret = &Coerce{Type: r1, Expr: r2}
																						return
																					}
																					frame.Fail()
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
															frame.Fail()
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
		}
		frame.Fail()
		goto block1
	}
	goto block1
block1:
	frame.Recover(checkpoint)
	r3 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c9 = frame.Peek()
		if frame.Flow == 0 {
			if c9 == '(' {
				frame.Consume()
				S(frame)
				r4 = ParseExprList(frame)
				S(frame)
				c10 = frame.Peek()
				if frame.Flow == 0 {
					if c10 == ')' {
						frame.Consume()
						ret = &Call{Name: r3, Args: r4}
						return
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
	goto block2
block2:
	frame.Recover(checkpoint)
	r5 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c11 = frame.Peek()
		if frame.Flow == 0 {
			if c11 == '{' {
				frame.Consume()
				S(frame)
				r6 = ParseNamedExprList(frame)
				S(frame)
				c12 = frame.Peek()
				if frame.Flow == 0 {
					if c12 == '}' {
						frame.Consume()
						ret = &Construct{Type: r5, Args: r6}
						return
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
	goto block3
block3:
	frame.Recover(checkpoint)
	r7 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c13 = frame.Peek()
		if frame.Flow == 0 {
			if c13 == '{' {
				frame.Consume()
				S(frame)
				r8 = ParseExprList(frame)
				S(frame)
				c14 = frame.Peek()
				if frame.Flow == 0 {
					if c14 == '}' {
						frame.Consume()
						ret = &ConstructList{Type: r7, Args: r8}
						return
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
	frame.Recover(checkpoint)
	r9 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		ret = r9
		return
	}
	frame.Recover(checkpoint)
	r10 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		ret = r10
		return
	}
	frame.Recover(checkpoint)
	c15 = frame.Peek()
	if frame.Flow == 0 {
		if c15 == '(' {
			frame.Consume()
			S(frame)
			r11 = ParseExpr(frame)
			if frame.Flow == 0 {
				S(frame)
				c16 = frame.Peek()
				if frame.Flow == 0 {
					if c16 == ')' {
						frame.Consume()
						ret = r11
						return
					}
					frame.Fail()
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
block5:
	frame.Recover(checkpoint)
	r12 = ParseNameRef(frame)
	if frame.Flow == 0 {
		ret = r12
		return
	}
	return
}

func ParseNameRef(frame *runtime.State) (ret *NameRef) {
	var r *Id
	r = Ident(frame)
	if frame.Flow == 0 {
		ret = &NameRef{Name: r}
		return
	}
	return
}

func ParseBinaryOp(frame *runtime.State, min_prec int) (ret ASTExpr) {
	var r0 ASTExpr
	var e ASTExpr
	var checkpoint int
	var r1 string
	var r2 int
	var r3 ASTExpr
	r0 = PrimaryExpr(frame)
	if frame.Flow == 0 {
		e = r0
		goto block1
	}
	return
block1:
	checkpoint = frame.Checkpoint()
	S(frame)
	r1, r2 = BinaryOperator(frame)
	if frame.Flow == 0 {
		if r2 < min_prec {
			frame.Fail()
			goto block2
		}
		S(frame)
		r3 = ParseBinaryOp(frame, r2+1)
		if frame.Flow == 0 {
			e = &BinaryOp{Left: e, Op: r1, Right: r3}
			goto block1
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint)
	ret = e
	return
}

func ParseExpr(frame *runtime.State) (ret ASTExpr) {
	var r ASTExpr
	r = ParseBinaryOp(frame, 1)
	if frame.Flow == 0 {
		ret = r
		return
	}
	return
}

func ParseCompoundStatement(frame *runtime.State) (ret ASTExpr) {
	var checkpoint0 int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var r0 []ASTExpr
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var r1 []ASTExpr
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var c12 rune
	var c13 rune
	var r2 []ASTExpr
	var r3 [][]ASTExpr
	var c14 rune
	var c15 rune
	var r4 []ASTExpr
	var blocks [][]ASTExpr
	var checkpoint1 int
	var c16 rune
	var c17 rune
	var r5 []ASTExpr
	var c18 rune
	var c19 rune
	var c20 rune
	var c21 rune
	var c22 rune
	var c23 rune
	var c24 rune
	var c25 rune
	var r6 []ASTExpr
	var c26 rune
	var c27 rune
	var r7 ASTExpr
	var r8 []ASTExpr
	var r9 []ASTExpr
	var checkpoint2 int
	var c28 rune
	var c29 rune
	var c30 rune
	var c31 rune
	var r10 []ASTExpr
	var else_ []ASTExpr
	checkpoint0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 's' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 't' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'a' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'r' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										r0 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											ret = &Repeat{Block: r0, Min: 0}
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
	frame.Recover(checkpoint0)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == 'p' {
			frame.Consume()
			c5 = frame.Peek()
			if frame.Flow == 0 {
				if c5 == 'l' {
					frame.Consume()
					c6 = frame.Peek()
					if frame.Flow == 0 {
						if c6 == 'u' {
							frame.Consume()
							c7 = frame.Peek()
							if frame.Flow == 0 {
								if c7 == 's' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										r1 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											ret = &Repeat{Block: r1, Min: 1}
											return
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
	frame.Recover(checkpoint0)
	c8 = frame.Peek()
	if frame.Flow == 0 {
		if c8 == 'c' {
			frame.Consume()
			c9 = frame.Peek()
			if frame.Flow == 0 {
				if c9 == 'h' {
					frame.Consume()
					c10 = frame.Peek()
					if frame.Flow == 0 {
						if c10 == 'o' {
							frame.Consume()
							c11 = frame.Peek()
							if frame.Flow == 0 {
								if c11 == 'o' {
									frame.Consume()
									c12 = frame.Peek()
									if frame.Flow == 0 {
										if c12 == 's' {
											frame.Consume()
											c13 = frame.Peek()
											if frame.Flow == 0 {
												if c13 == 'e' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														r2 = ParseCodeBlock(frame)
														if frame.Flow == 0 {
															r3 = [][]ASTExpr{r2}
															S(frame)
															c14 = frame.Peek()
															if frame.Flow == 0 {
																if c14 == 'o' {
																	frame.Consume()
																	c15 = frame.Peek()
																	if frame.Flow == 0 {
																		if c15 == 'r' {
																			frame.Consume()
																			EndKeyword(frame)
																			if frame.Flow == 0 {
																				S(frame)
																				r4 = ParseCodeBlock(frame)
																				if frame.Flow == 0 {
																					blocks = append(r3, r4)
																					goto block3
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
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c16 = frame.Peek()
	if frame.Flow == 0 {
		if c16 == 'o' {
			frame.Consume()
			c17 = frame.Peek()
			if frame.Flow == 0 {
				if c17 == 'r' {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						S(frame)
						r5 = ParseCodeBlock(frame)
						if frame.Flow == 0 {
							blocks = append(blocks, r5)
							goto block3
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
block4:
	frame.Recover(checkpoint1)
	ret = &Choice{Blocks: blocks}
	return
block5:
	frame.Recover(checkpoint0)
	c18 = frame.Peek()
	if frame.Flow == 0 {
		if c18 == 'q' {
			frame.Consume()
			c19 = frame.Peek()
			if frame.Flow == 0 {
				if c19 == 'u' {
					frame.Consume()
					c20 = frame.Peek()
					if frame.Flow == 0 {
						if c20 == 'e' {
							frame.Consume()
							c21 = frame.Peek()
							if frame.Flow == 0 {
								if c21 == 's' {
									frame.Consume()
									c22 = frame.Peek()
									if frame.Flow == 0 {
										if c22 == 't' {
											frame.Consume()
											c23 = frame.Peek()
											if frame.Flow == 0 {
												if c23 == 'i' {
													frame.Consume()
													c24 = frame.Peek()
													if frame.Flow == 0 {
														if c24 == 'o' {
															frame.Consume()
															c25 = frame.Peek()
															if frame.Flow == 0 {
																if c25 == 'n' {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		S(frame)
																		r6 = ParseCodeBlock(frame)
																		if frame.Flow == 0 {
																			ret = &Optional{Block: r6}
																			return
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
	frame.Recover(checkpoint0)
	c26 = frame.Peek()
	if frame.Flow == 0 {
		if c26 == 'i' {
			frame.Consume()
			c27 = frame.Peek()
			if frame.Flow == 0 {
				if c27 == 'f' {
					frame.Consume()
					EndKeyword(frame)
					if frame.Flow == 0 {
						S(frame)
						r7 = ParseExpr(frame)
						if frame.Flow == 0 {
							S(frame)
							r8 = ParseCodeBlock(frame)
							if frame.Flow == 0 {
								r9 = []ASTExpr{}
								checkpoint2 = frame.Checkpoint()
								S(frame)
								c28 = frame.Peek()
								if frame.Flow == 0 {
									if c28 == 'e' {
										frame.Consume()
										c29 = frame.Peek()
										if frame.Flow == 0 {
											if c29 == 'l' {
												frame.Consume()
												c30 = frame.Peek()
												if frame.Flow == 0 {
													if c30 == 's' {
														frame.Consume()
														c31 = frame.Peek()
														if frame.Flow == 0 {
															if c31 == 'e' {
																frame.Consume()
																EndKeyword(frame)
																if frame.Flow == 0 {
																	S(frame)
																	r10 = ParseCodeBlock(frame)
																	if frame.Flow == 0 {
																		else_ = r10
																		goto block8
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
													frame.Fail()
													goto block7
												}
												goto block7
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
							}
							return
						}
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
block7:
	frame.Recover(checkpoint2)
	else_ = r9
	goto block8
block8:
	ret = &If{Expr: r7, Block: r8, Else: else_}
	return
}

func EOS(frame *runtime.State) {
	var c rune
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ';' {
			frame.Consume()
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseStatement(frame *runtime.State) (ret ASTExpr) {
	var checkpoint0 int
	var r0 ASTExpr
	var c0 rune
	var c1 rune
	var c2 rune
	var r1 *NameRef
	var r2 ASTTypeRef
	var expr0 ASTExpr
	var checkpoint1 int
	var c3 rune
	var r3 ASTExpr
	var expr1 ASTExpr
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var c12 rune
	var c13 rune
	var r4 []ASTExpr
	var r5 []ASTExpr
	var defined0 bool
	var checkpoint2 int
	var c14 rune
	var c15 rune
	var defined1 bool
	var c16 rune
	var r6 ASTExpr
	var r7 ASTExpr
	checkpoint0 = frame.Checkpoint()
	r0 = ParseCompoundStatement(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint0)
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'v' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'a' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'r' {
							frame.Consume()
							EndKeyword(frame)
							if frame.Flow == 0 {
								S(frame)
								r1 = ParseNameRef(frame)
								if frame.Flow == 0 {
									S(frame)
									r2 = ParseTypeRef(frame)
									if frame.Flow == 0 {
										expr0 = nil
										checkpoint1 = frame.Checkpoint()
										S(frame)
										c3 = frame.Peek()
										if frame.Flow == 0 {
											if c3 == '=' {
												frame.Consume()
												S(frame)
												r3 = ParseExpr(frame)
												if frame.Flow == 0 {
													expr1 = r3
													goto block2
												}
												goto block1
											}
											frame.Fail()
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
	frame.Recover(checkpoint1)
	expr1 = expr0
	goto block2
block2:
	EOS(frame)
	if frame.Flow == 0 {
		ret = &Assign{Expr: expr1, Targets: []ASTExpr{r1}, Type: r2, Define: true}
		return
	}
	goto block3
block3:
	frame.Recover(checkpoint0)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == 'f' {
			frame.Consume()
			c5 = frame.Peek()
			if frame.Flow == 0 {
				if c5 == 'a' {
					frame.Consume()
					c6 = frame.Peek()
					if frame.Flow == 0 {
						if c6 == 'i' {
							frame.Consume()
							c7 = frame.Peek()
							if frame.Flow == 0 {
								if c7 == 'l' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										EOS(frame)
										if frame.Flow == 0 {
											ret = &Fail{}
											return
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
	frame.Recover(checkpoint0)
	c8 = frame.Peek()
	if frame.Flow == 0 {
		if c8 == 'r' {
			frame.Consume()
			c9 = frame.Peek()
			if frame.Flow == 0 {
				if c9 == 'e' {
					frame.Consume()
					c10 = frame.Peek()
					if frame.Flow == 0 {
						if c10 == 't' {
							frame.Consume()
							c11 = frame.Peek()
							if frame.Flow == 0 {
								if c11 == 'u' {
									frame.Consume()
									c12 = frame.Peek()
									if frame.Flow == 0 {
										if c12 == 'r' {
											frame.Consume()
											c13 = frame.Peek()
											if frame.Flow == 0 {
												if c13 == 'n' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														r4 = ParseExprList(frame)
														EOS(frame)
														if frame.Flow == 0 {
															ret = &Return{Exprs: r4}
															return
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
	frame.Recover(checkpoint0)
	r5 = ParseTargetList(frame)
	if frame.Flow == 0 {
		S(frame)
		defined0 = false
		checkpoint2 = frame.Checkpoint()
		c14 = frame.Peek()
		if frame.Flow == 0 {
			if c14 == ':' {
				frame.Consume()
				c15 = frame.Peek()
				if frame.Flow == 0 {
					if c15 == '=' {
						frame.Consume()
						defined1 = true
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
block6:
	frame.Recover(checkpoint2)
	c16 = frame.Peek()
	if frame.Flow == 0 {
		if c16 == '=' {
			frame.Consume()
			defined1 = defined0
			goto block7
		}
		frame.Fail()
		goto block8
	}
	goto block8
block7:
	S(frame)
	r6 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret = &Assign{Expr: r6, Targets: r5, Define: defined1}
			return
		}
		goto block8
	}
	goto block8
block8:
	frame.Recover(checkpoint0)
	r7 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret = r7
			return
		}
		return
	}
	return
}

func ParseCodeBlock(frame *runtime.State) (ret []ASTExpr) {
	var c0 rune
	var exprs0 []ASTExpr
	var checkpoint int
	var r0 ASTExpr
	var r1 []ASTExpr
	var exprs1 []ASTExpr
	var c1 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '{' {
			frame.Consume()
			S(frame)
			exprs0 = []ASTExpr{}
			goto block1
		}
		frame.Fail()
		return
	}
	return
block1:
	checkpoint = frame.Checkpoint()
	r0 = ParseStatement(frame)
	if frame.Flow == 0 {
		r1 = append(exprs0, r0)
		S(frame)
		exprs0 = r1
		goto block1
	}
	exprs1 = exprs0
	frame.Recover(checkpoint)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '}' {
			frame.Consume()
			ret = exprs1
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseParenthTypeList(frame *runtime.State) (ret []ASTTypeRef) {
	var c0 rune
	var r0 []ASTTypeRef
	var checkpoint0 int
	var r1 ASTTypeRef
	var r2 []ASTTypeRef
	var types0 []ASTTypeRef
	var checkpoint1 int
	var c1 rune
	var r3 ASTTypeRef
	var r4 []ASTTypeRef
	var types1 []ASTTypeRef
	var types2 []ASTTypeRef
	var types3 []ASTTypeRef
	var c2 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '(' {
			frame.Consume()
			S(frame)
			r0 = []ASTTypeRef{}
			checkpoint0 = frame.Checkpoint()
			r1 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				r2 = append(r0, r1)
				S(frame)
				types0 = r2
				goto block1
			}
			types3 = r0
			frame.Recover(checkpoint0)
			types2 = types3
			goto block3
		}
		frame.Fail()
		return
	}
	return
block1:
	checkpoint1 = frame.Checkpoint()
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == ',' {
			frame.Consume()
			S(frame)
			r3 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				r4 = append(types0, r3)
				S(frame)
				types0 = r4
				goto block1
			}
			types1 = types0
			goto block2
		}
		frame.Fail()
		types1 = types0
		goto block2
	}
	types1 = types0
	goto block2
block2:
	frame.Recover(checkpoint1)
	types2 = types1
	goto block3
block3:
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == ')' {
			frame.Consume()
			ret = types2
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseStructDecl(frame *runtime.State) (ret *StructDecl) {
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var r0 *Id
	var c_b bool
	var checkpoint0 int
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var scoped bool
	var r1 []ASTTypeRef
	var checkpoint1 int
	var c12 rune
	var c13 rune
	var c14 rune
	var c15 rune
	var c16 rune
	var c17 rune
	var c18 rune
	var c19 rune
	var r2 []ASTTypeRef
	var contains0 []ASTTypeRef
	var contains1 []ASTTypeRef
	var impl0 ASTTypeRef
	var checkpoint2 int
	var c20 rune
	var c21 rune
	var c22 rune
	var c23 rune
	var c24 rune
	var c25 rune
	var c26 rune
	var c27 rune
	var c28 rune
	var c29 rune
	var r3 ASTTypeRef
	var impl1 ASTTypeRef
	var impl2 ASTTypeRef
	var c30 rune
	var fields []*FieldDecl
	var checkpoint3 int
	var r4 *Id
	var r5 ASTTypeRef
	var c31 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 's' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 't' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'r' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'u' {
									frame.Consume()
									c4 = frame.Peek()
									if frame.Flow == 0 {
										if c4 == 'c' {
											frame.Consume()
											c5 = frame.Peek()
											if frame.Flow == 0 {
												if c5 == 't' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														r0 = Ident(frame)
														if frame.Flow == 0 {
															S(frame)
															c_b = false
															checkpoint0 = frame.Checkpoint()
															c6 = frame.Peek()
															if frame.Flow == 0 {
																if c6 == 's' {
																	frame.Consume()
																	c7 = frame.Peek()
																	if frame.Flow == 0 {
																		if c7 == 'c' {
																			frame.Consume()
																			c8 = frame.Peek()
																			if frame.Flow == 0 {
																				if c8 == 'o' {
																					frame.Consume()
																					c9 = frame.Peek()
																					if frame.Flow == 0 {
																						if c9 == 'p' {
																							frame.Consume()
																							c10 = frame.Peek()
																							if frame.Flow == 0 {
																								if c10 == 'e' {
																									frame.Consume()
																									c11 = frame.Peek()
																									if frame.Flow == 0 {
																										if c11 == 'd' {
																											frame.Consume()
																											EndKeyword(frame)
																											if frame.Flow == 0 {
																												S(frame)
																												scoped = true
																												goto block2
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
														return
													}
													return
												}
												frame.Fail()
												return
											}
											return
										}
										frame.Fail()
										return
									}
									return
								}
								frame.Fail()
								return
							}
							return
						}
						frame.Fail()
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
block1:
	frame.Recover(checkpoint0)
	scoped = c_b
	goto block2
block2:
	r1 = []ASTTypeRef{}
	checkpoint1 = frame.Checkpoint()
	c12 = frame.Peek()
	if frame.Flow == 0 {
		if c12 == 'c' {
			frame.Consume()
			c13 = frame.Peek()
			if frame.Flow == 0 {
				if c13 == 'o' {
					frame.Consume()
					c14 = frame.Peek()
					if frame.Flow == 0 {
						if c14 == 'n' {
							frame.Consume()
							c15 = frame.Peek()
							if frame.Flow == 0 {
								if c15 == 't' {
									frame.Consume()
									c16 = frame.Peek()
									if frame.Flow == 0 {
										if c16 == 'a' {
											frame.Consume()
											c17 = frame.Peek()
											if frame.Flow == 0 {
												if c17 == 'i' {
													frame.Consume()
													c18 = frame.Peek()
													if frame.Flow == 0 {
														if c18 == 'n' {
															frame.Consume()
															c19 = frame.Peek()
															if frame.Flow == 0 {
																if c19 == 's' {
																	frame.Consume()
																	EndKeyword(frame)
																	if frame.Flow == 0 {
																		S(frame)
																		r2 = ParseParenthTypeList(frame)
																		if frame.Flow == 0 {
																			S(frame)
																			contains0 = r2
																			goto block4
																		}
																		contains1 = r1
																		goto block3
																	}
																	contains1 = r1
																	goto block3
																}
																frame.Fail()
																contains1 = r1
																goto block3
															}
															contains1 = r1
															goto block3
														}
														frame.Fail()
														contains1 = r1
														goto block3
													}
													contains1 = r1
													goto block3
												}
												frame.Fail()
												contains1 = r1
												goto block3
											}
											contains1 = r1
											goto block3
										}
										frame.Fail()
										contains1 = r1
										goto block3
									}
									contains1 = r1
									goto block3
								}
								frame.Fail()
								contains1 = r1
								goto block3
							}
							contains1 = r1
							goto block3
						}
						frame.Fail()
						contains1 = r1
						goto block3
					}
					contains1 = r1
					goto block3
				}
				frame.Fail()
				contains1 = r1
				goto block3
			}
			contains1 = r1
			goto block3
		}
		frame.Fail()
		contains1 = r1
		goto block3
	}
	contains1 = r1
	goto block3
block3:
	frame.Recover(checkpoint1)
	contains0 = contains1
	goto block4
block4:
	impl0 = nil
	checkpoint2 = frame.Checkpoint()
	c20 = frame.Peek()
	if frame.Flow == 0 {
		if c20 == 'i' {
			frame.Consume()
			c21 = frame.Peek()
			if frame.Flow == 0 {
				if c21 == 'm' {
					frame.Consume()
					c22 = frame.Peek()
					if frame.Flow == 0 {
						if c22 == 'p' {
							frame.Consume()
							c23 = frame.Peek()
							if frame.Flow == 0 {
								if c23 == 'l' {
									frame.Consume()
									c24 = frame.Peek()
									if frame.Flow == 0 {
										if c24 == 'e' {
											frame.Consume()
											c25 = frame.Peek()
											if frame.Flow == 0 {
												if c25 == 'm' {
													frame.Consume()
													c26 = frame.Peek()
													if frame.Flow == 0 {
														if c26 == 'e' {
															frame.Consume()
															c27 = frame.Peek()
															if frame.Flow == 0 {
																if c27 == 'n' {
																	frame.Consume()
																	c28 = frame.Peek()
																	if frame.Flow == 0 {
																		if c28 == 't' {
																			frame.Consume()
																			c29 = frame.Peek()
																			if frame.Flow == 0 {
																				if c29 == 's' {
																					frame.Consume()
																					EndKeyword(frame)
																					if frame.Flow == 0 {
																						S(frame)
																						r3 = ParseTypeRef(frame)
																						if frame.Flow == 0 {
																							S(frame)
																							impl1 = r3
																							goto block6
																						}
																						impl2 = impl0
																						goto block5
																					}
																					impl2 = impl0
																					goto block5
																				}
																				frame.Fail()
																				impl2 = impl0
																				goto block5
																			}
																			impl2 = impl0
																			goto block5
																		}
																		frame.Fail()
																		impl2 = impl0
																		goto block5
																	}
																	impl2 = impl0
																	goto block5
																}
																frame.Fail()
																impl2 = impl0
																goto block5
															}
															impl2 = impl0
															goto block5
														}
														frame.Fail()
														impl2 = impl0
														goto block5
													}
													impl2 = impl0
													goto block5
												}
												frame.Fail()
												impl2 = impl0
												goto block5
											}
											impl2 = impl0
											goto block5
										}
										frame.Fail()
										impl2 = impl0
										goto block5
									}
									impl2 = impl0
									goto block5
								}
								frame.Fail()
								impl2 = impl0
								goto block5
							}
							impl2 = impl0
							goto block5
						}
						frame.Fail()
						impl2 = impl0
						goto block5
					}
					impl2 = impl0
					goto block5
				}
				frame.Fail()
				impl2 = impl0
				goto block5
			}
			impl2 = impl0
			goto block5
		}
		frame.Fail()
		impl2 = impl0
		goto block5
	}
	impl2 = impl0
	goto block5
block5:
	frame.Recover(checkpoint2)
	impl1 = impl2
	goto block6
block6:
	c30 = frame.Peek()
	if frame.Flow == 0 {
		if c30 == '{' {
			frame.Consume()
			S(frame)
			fields = []*FieldDecl{}
			goto block7
		}
		frame.Fail()
		return
	}
	return
block7:
	checkpoint3 = frame.Checkpoint()
	r4 = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		r5 = ParseTypeRef(frame)
		if frame.Flow == 0 {
			S(frame)
			fields = append(fields, &FieldDecl{Name: r4, Type: r5})
			goto block7
		}
		goto block8
	}
	goto block8
block8:
	frame.Recover(checkpoint3)
	c31 = frame.Peek()
	if frame.Flow == 0 {
		if c31 == '}' {
			frame.Consume()
			ret = &StructDecl{Name: r0, Implements: impl1, Fields: fields, Scoped: scoped, Contains: contains0}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseParam(frame *runtime.State) (ret *Param) {
	var r0 *NameRef
	var r1 ASTTypeRef
	r0 = ParseNameRef(frame)
	if frame.Flow == 0 {
		S(frame)
		r1 = ParseTypeRef(frame)
		if frame.Flow == 0 {
			ret = &Param{Name: r0, Type: r1}
			return
		}
		return
	}
	return
}

func ParseParamList(frame *runtime.State) (ret []*Param) {
	var r0 []*Param
	var checkpoint0 int
	var r1 *Param
	var params0 []*Param
	var checkpoint1 int
	var c rune
	var r2 *Param
	var params1 []*Param
	r0 = []*Param{}
	checkpoint0 = frame.Checkpoint()
	r1 = ParseParam(frame)
	if frame.Flow == 0 {
		params0 = append(r0, r1)
		goto block1
	}
	frame.Recover(checkpoint0)
	params1 = r0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r2 = ParseParam(frame)
			if frame.Flow == 0 {
				params0 = append(params0, r2)
				goto block1
			}
			goto block2
		}
		frame.Fail()
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint1)
	params1 = params0
	goto block3
block3:
	ret = params1
	return
}

func ParseFuncDecl(frame *runtime.State) (ret *FuncDecl) {
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var r0 *Id
	var c4 rune
	var r1 []*Param
	var c5 rune
	var r2 []ASTTypeRef
	var r3 []ASTExpr
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'f' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'u' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'n' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'c' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										r0 = Ident(frame)
										if frame.Flow == 0 {
											S(frame)
											c4 = frame.Peek()
											if frame.Flow == 0 {
												if c4 == '(' {
													frame.Consume()
													S(frame)
													r1 = ParseParamList(frame)
													S(frame)
													c5 = frame.Peek()
													if frame.Flow == 0 {
														if c5 == ')' {
															frame.Consume()
															S(frame)
															r2 = ParseReturnTypeList(frame)
															S(frame)
															r3 = ParseCodeBlock(frame)
															if frame.Flow == 0 {
																ret = &FuncDecl{Name: r0, Params: r1, ReturnTypes: r2, Block: r3, LocalInfo_Scope: &LocalInfo_Scope{}}
																return
															}
															return
														}
														frame.Fail()
														return
													}
													return
												}
												frame.Fail()
												return
											}
											return
										}
										return
									}
									return
								}
								frame.Fail()
								return
							}
							return
						}
						frame.Fail()
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseMatchState(frame *runtime.State) (ret string) {
	var checkpoint0 int
	var begin int
	var checkpoint1 int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	var slice string
	checkpoint0 = frame.Checkpoint()
	begin = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'N' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'O' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'R' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'M' {
									frame.Consume()
									c4 = frame.Peek()
									if frame.Flow == 0 {
										if c4 == 'A' {
											frame.Consume()
											c5 = frame.Peek()
											if frame.Flow == 0 {
												if c5 == 'L' {
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
	frame.Recover(checkpoint1)
	c6 = frame.Peek()
	if frame.Flow == 0 {
		if c6 == 'F' {
			frame.Consume()
			c7 = frame.Peek()
			if frame.Flow == 0 {
				if c7 == 'A' {
					frame.Consume()
					c8 = frame.Peek()
					if frame.Flow == 0 {
						if c8 == 'I' {
							frame.Consume()
							c9 = frame.Peek()
							if frame.Flow == 0 {
								if c9 == 'L' {
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
	slice = frame.Slice(begin, frame.Checkpoint())
	EndKeyword(frame)
	if frame.Flow == 0 {
		ret = slice
		return
	}
	goto block3
block3:
	frame.Recover(checkpoint0)
	ret = "NORMAL"
	return
}

func ParseTest(frame *runtime.State) (ret *Test) {
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var r0 *Id
	var r1 ASTExpr
	var r2 string
	var r3 string
	var r4 Destructure
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 't' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'e' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 's' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 't' {
									frame.Consume()
									EndKeyword(frame)
									if frame.Flow == 0 {
										S(frame)
										r0 = Ident(frame)
										if frame.Flow == 0 {
											S(frame)
											r1 = ParseExpr(frame)
											if frame.Flow == 0 {
												S(frame)
												r2 = DecodeString(frame)
												if frame.Flow == 0 {
													S(frame)
													r3 = ParseMatchState(frame)
													S(frame)
													r4 = ParseDestructure(frame)
													if frame.Flow == 0 {
														ret = &Test{Name: r0, Rule: r1, Input: r2, Flow: r3, Destructure: r4}
														return
													}
													return
												}
												return
											}
											return
										}
										return
									}
									return
								}
								frame.Fail()
								return
							}
							return
						}
						frame.Fail()
						return
					}
					return
				}
				frame.Fail()
				return
			}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseImports(frame *runtime.State) (ret []*ImportDecl) {
	var r0 []*ImportDecl
	var checkpoint0 int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var imports0 []*ImportDecl
	var checkpoint1 int
	var r1 *StringLiteral
	var r2 []*ImportDecl
	var imports1 []*ImportDecl
	var c7 rune
	var imports2 []*ImportDecl
	var imports3 []*ImportDecl
	r0 = []*ImportDecl{}
	checkpoint0 = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == 'i' {
			frame.Consume()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == 'm' {
					frame.Consume()
					c2 = frame.Peek()
					if frame.Flow == 0 {
						if c2 == 'p' {
							frame.Consume()
							c3 = frame.Peek()
							if frame.Flow == 0 {
								if c3 == 'o' {
									frame.Consume()
									c4 = frame.Peek()
									if frame.Flow == 0 {
										if c4 == 'r' {
											frame.Consume()
											c5 = frame.Peek()
											if frame.Flow == 0 {
												if c5 == 't' {
													frame.Consume()
													EndKeyword(frame)
													if frame.Flow == 0 {
														S(frame)
														c6 = frame.Peek()
														if frame.Flow == 0 {
															if c6 == '(' {
																frame.Consume()
																S(frame)
																imports0 = r0
																goto block1
															}
															frame.Fail()
															imports3 = r0
															goto block2
														}
														imports3 = r0
														goto block2
													}
													imports3 = r0
													goto block2
												}
												frame.Fail()
												imports3 = r0
												goto block2
											}
											imports3 = r0
											goto block2
										}
										frame.Fail()
										imports3 = r0
										goto block2
									}
									imports3 = r0
									goto block2
								}
								frame.Fail()
								imports3 = r0
								goto block2
							}
							imports3 = r0
							goto block2
						}
						frame.Fail()
						imports3 = r0
						goto block2
					}
					imports3 = r0
					goto block2
				}
				frame.Fail()
				imports3 = r0
				goto block2
			}
			imports3 = r0
			goto block2
		}
		frame.Fail()
		imports3 = r0
		goto block2
	}
	imports3 = r0
	goto block2
block1:
	checkpoint1 = frame.Checkpoint()
	r1 = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		r2 = append(imports0, &ImportDecl{Path: r1})
		S(frame)
		imports0 = r2
		goto block1
	}
	imports1 = imports0
	frame.Recover(checkpoint1)
	c7 = frame.Peek()
	if frame.Flow == 0 {
		if c7 == ')' {
			frame.Consume()
			imports2 = imports1
			goto block3
		}
		frame.Fail()
		imports3 = imports1
		goto block2
	}
	imports3 = imports1
	goto block2
block2:
	frame.Recover(checkpoint0)
	imports2 = imports3
	goto block3
block3:
	ret = imports2
	return
}

func ParseFile(frame *runtime.State) (ret *File) {
	var r0 []ASTDecl
	var r1 []*Test
	var r2 []*ImportDecl
	var decls0 []ASTDecl
	var tests0 []*Test
	var checkpoint0 int
	var checkpoint1 int
	var r3 *FuncDecl
	var decls1 []ASTDecl
	var tests1 []*Test
	var r4 *StructDecl
	var r5 *Test
	var decls2 []ASTDecl
	var tests2 []*Test
	var checkpoint2 int
	r0 = []ASTDecl{}
	r1 = []*Test{}
	S(frame)
	r2 = ParseImports(frame)
	S(frame)
	decls0, tests0 = r0, r1
	goto block1
block1:
	checkpoint0 = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	r3 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		decls1, tests1 = append(decls0, r3), tests0
		goto block2
	}
	frame.Recover(checkpoint1)
	r4 = ParseStructDecl(frame)
	if frame.Flow == 0 {
		decls1, tests1 = append(decls0, r4), tests0
		goto block2
	}
	frame.Recover(checkpoint1)
	r5 = ParseTest(frame)
	if frame.Flow == 0 {
		decls1, tests1 = decls0, append(tests0, r5)
		goto block2
	}
	decls2, tests2 = decls0, tests0
	frame.Recover(checkpoint0)
	checkpoint2 = frame.LookaheadBegin()
	frame.Peek()
	if frame.Flow == 0 {
		frame.Consume()
		frame.LookaheadFail(checkpoint2)
		return
	}
	frame.LookaheadNormal(checkpoint2)
	ret = &File{Imports: r2, Decls: decls2, Tests: tests2}
	return
block2:
	S(frame)
	decls0, tests0 = decls1, tests1
	goto block1
}
