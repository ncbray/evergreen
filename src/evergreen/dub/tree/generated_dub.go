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
	Type ASTTypeRef
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
	Name *Id
}

func (node *NameRef) isASTExpr() {
}

type GetLocal struct {
	Info *LocalInfo
}

func (node *GetLocal) isASTExpr() {
}

type SetLocal struct {
	Info *LocalInfo
}

func (node *SetLocal) isASTExpr() {
}

type Discard struct {
}

func (node *Discard) isASTExpr() {
}

type GetFunction struct {
	Func core.Callable
}

func (node *GetFunction) isASTExpr() {
}

type GetFunctionTemplate struct {
	Template *core.IntrinsicFunctionTemplate
}

func (node *GetFunctionTemplate) isASTExpr() {
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
	Type ASTTypeRef
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
	Expr   ASTExpr
	Pos    int
	Args   []ASTExpr
	Target core.Callable
	T      core.DubType
}

func (node *Call) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
}

type Return struct {
	Pos   int
	Exprs []ASTExpr
}

func (node *Return) isASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	OpPos int
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
	Name *Id
	Type ASTTypeRef
	Info *LocalInfo
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
	var p int
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
	p = frame.Checkpoint()
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
	ret = &Id{Pos: p, Text: frame.Slice(begin, frame.Checkpoint())}
	return
}

func ParseNumericLiteral(frame *runtime.State) (ret ASTExpr) {
	var value0 int
	var c_i int
	var begin int
	var c0 rune
	var digit0 int
	var value1 int
	var checkpoint0 int
	var c1 rune
	var digit1 int
	var checkpoint1 int
	var c2 rune
	var c3 rune
	var digit2 int
	var value2 int
	var divisor0 int
	var checkpoint2 int
	var c4 rune
	var digit3 int
	var value3 int
	var divisor1 int
	var text string
	value0 = 0
	c_i = 1
	begin = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 >= '0' {
			if c0 <= '9' {
				frame.Consume()
				digit0 = int(c0) - int('0')
				value1 = value0*10 + digit0
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
				digit1 = int(c1) - int('0')
				value1 = value1*10 + digit1
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
						digit2 = int(c3) - int('0')
						value2, divisor0 = value1*10+digit2, c_i*10
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
				digit3 = int(c4) - int('0')
				value2, divisor0 = value2*10+digit3, divisor0*10
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
	text = frame.Slice(begin, frame.Checkpoint())
	if divisor1 > 1 {
		ret = &Float32Literal{Text: text, Value: float32(value3) / float32(divisor1)}
		return
	}
	ret = &IntLiteral{Text: text, Value: value3}
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
	var begin int
	var c0 rune
	var checkpoint int
	var c1 rune
	var value0 rune
	var c2 rune
	var value1 rune
	var c3 rune
	begin = frame.Checkpoint()
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
				value0 = c1
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
			value1 = EscapedChar(frame)
			if frame.Flow == 0 {
				value0 = value1
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
			ret0, ret1 = value0, frame.Slice(begin, frame.Checkpoint())
			return
		}
		frame.Fail()
		return
	}
	return
}

func DecodeBool(frame *runtime.State) (ret0 bool, ret1 string) {
	var begin int
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
	begin = frame.Checkpoint()
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
		ret0, ret1 = value, frame.Slice(begin, frame.Checkpoint())
		return
	}
	return
}

func ParseStringLiteral(frame *runtime.State) (ret *StringLiteral) {
	var begin int
	var value string
	begin = frame.Checkpoint()
	value = DecodeString(frame)
	if frame.Flow == 0 {
		ret = &StringLiteral{Pos: begin, Text: frame.Slice(begin, frame.Checkpoint()), Value: value}
		return
	}
	return
}

func Literal(frame *runtime.State) (ret ASTExpr) {
	var checkpoint int
	var value0 rune
	var text0 string
	var r0 *StringLiteral
	var r1 ASTExpr
	var value1 bool
	var text1 string
	var c0 rune
	var c1 rune
	var c2 rune
	checkpoint = frame.Checkpoint()
	value0, text0 = DecodeRune(frame)
	if frame.Flow == 0 {
		ret = &RuneLiteral{Text: text0, Value: value0}
		return
	}
	frame.Recover(checkpoint)
	r0 = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		ret = r0
		return
	}
	frame.Recover(checkpoint)
	r1 = ParseNumericLiteral(frame)
	if frame.Flow == 0 {
		ret = r1
		return
	}
	frame.Recover(checkpoint)
	value1, text1 = DecodeBool(frame)
	if frame.Flow == 0 {
		ret = &BoolLiteral{Text: text1, Value: value1}
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
	var e TextMatch
	var c1 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '/' {
			frame.Consume()
			S(frame)
			e = ParseMatchChoice(frame)
			if frame.Flow == 0 {
				S(frame)
				c1 = frame.Peek()
				if frame.Flow == 0 {
					if c1 == '/' {
						frame.Consume()
						ret = &StringMatch{Match: e}
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
	var e *RuneRangeMatch
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == '$' {
			frame.Consume()
			S(frame)
			e = MatchRune(frame)
			if frame.Flow == 0 {
				ret = &RuneMatch{Match: e}
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
	var pkg *Id
	var c rune
	var r0 *Id
	var r1 *Id
	checkpoint = frame.Checkpoint()
	pkg = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == '.' {
				frame.Consume()
				S(frame)
				r0 = Ident(frame)
				if frame.Flow == 0 {
					ret = &QualifiedTypeRef{Package: pkg, Name: r0}
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
	r1 = Ident(frame)
	if frame.Flow == 0 {
		ret = &TypeRef{Name: r1}
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
	var t0 ASTTypeRef
	var c0 rune
	var fields0 []*DestructureField
	var checkpoint1 int
	var name *Id
	var c1 rune
	var d Destructure
	var c2 rune
	var t1 *ListTypeRef
	var c3 rune
	var fields1 []Destructure
	var checkpoint2 int
	var r0 Destructure
	var fields2 []Destructure
	var fields3 []Destructure
	var c4 rune
	var r1 ASTExpr
	checkpoint0 = frame.Checkpoint()
	t0 = ParseStructTypeRef(frame)
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
	name = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c1 = frame.Peek()
		if frame.Flow == 0 {
			if c1 == ':' {
				frame.Consume()
				S(frame)
				d = ParseDestructure(frame)
				if frame.Flow == 0 {
					S(frame)
					fields0 = append(fields0, &DestructureField{Name: name, Destructure: d})
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
			ret = &DestructureStruct{Type: t0, Args: fields0}
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(checkpoint0)
	t1 = ParseListTypeRef(frame)
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
	r0 = ParseDestructure(frame)
	if frame.Flow == 0 {
		fields2 = append(fields1, r0)
		S(frame)
		fields1 = fields2
		goto block4
	}
	fields3 = fields1
	frame.Recover(checkpoint2)
	c4 = frame.Peek()
	if frame.Flow == 0 {
		if c4 == '}' {
			frame.Consume()
			ret = &DestructureList{Type: t1, Args: fields3}
			return
		}
		frame.Fail()
		goto block5
	}
	goto block5
block5:
	frame.Recover(checkpoint0)
	r1 = Literal(frame)
	if frame.Flow == 0 {
		ret = &DestructureValue{Expr: r1}
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
	var min rune
	var checkpoint int
	var c rune
	var max0 rune
	var max1 rune
	min = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		checkpoint = frame.Checkpoint()
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == '-' {
				frame.Consume()
				max0 = ParseRuneFilterRune(frame)
				if frame.Flow == 0 {
					max1 = max0
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
	max1 = min
	goto block2
block2:
	ret = &RuneFilter{Min: min, Max: max1}
	return
}

func MatchRune(frame *runtime.State) (ret *RuneRangeMatch) {
	var c0 rune
	var c_b bool
	var filters0 []*RuneFilter
	var checkpoint0 int
	var c1 rune
	var invert bool
	var filters1 []*RuneFilter
	var checkpoint1 int
	var r *RuneFilter
	var c2 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '[' {
			frame.Consume()
			c_b = false
			filters0 = []*RuneFilter{}
			checkpoint0 = frame.Checkpoint()
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == '^' {
					frame.Consume()
					invert, filters1 = true, filters0
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
	invert, filters1 = c_b, filters0
	goto block2
block2:
	checkpoint1 = frame.Checkpoint()
	r = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		filters1 = append(filters1, r)
		goto block2
	}
	frame.Recover(checkpoint1)
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == ']' {
			frame.Consume()
			ret = &RuneRangeMatch{Invert: invert, Filters: filters1}
			return
		}
		frame.Fail()
		return
	}
	return
}

func Atom(frame *runtime.State) (ret TextMatch) {
	var checkpoint int
	var r *RuneRangeMatch
	var value string
	var c0 rune
	var e TextMatch
	var c1 rune
	checkpoint = frame.Checkpoint()
	r = MatchRune(frame)
	if frame.Flow == 0 {
		ret = r
		return
	}
	frame.Recover(checkpoint)
	value = DecodeString(frame)
	if frame.Flow == 0 {
		ret = &StringLiteralMatch{Value: value}
		return
	}
	frame.Recover(checkpoint)
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '(' {
			frame.Consume()
			S(frame)
			e = ParseMatchChoice(frame)
			if frame.Flow == 0 {
				S(frame)
				c1 = frame.Peek()
				if frame.Flow == 0 {
					if c1 == ')' {
						frame.Consume()
						ret = e
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
	var e TextMatch
	var checkpoint int
	var c0 rune
	var c1 rune
	var c2 rune
	e = Atom(frame)
	if frame.Flow == 0 {
		checkpoint = frame.Checkpoint()
		S(frame)
		c0 = frame.Peek()
		if frame.Flow == 0 {
			if c0 == '*' {
				frame.Consume()
				ret = &MatchRepeat{Match: e, Min: 0}
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
			ret = &MatchRepeat{Match: e, Min: 1}
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
			ret = &MatchChoice{Matches: []TextMatch{e, &MatchSequence{Matches: []TextMatch{}}}}
			return
		}
		frame.Fail()
		goto block3
	}
	goto block3
block3:
	frame.Recover(checkpoint)
	ret = e
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
	var e TextMatch
	var checkpoint0 int
	var l0 []TextMatch
	var r0 TextMatch
	var l1 []TextMatch
	var checkpoint1 int
	var r1 TextMatch
	e = MatchPrefix(frame)
	if frame.Flow == 0 {
		checkpoint0 = frame.Checkpoint()
		l0 = []TextMatch{e}
		S(frame)
		r0 = MatchPrefix(frame)
		if frame.Flow == 0 {
			l1 = append(l0, r0)
			goto block1
		}
		frame.Recover(checkpoint0)
		ret = e
		return
	}
	return
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	r1 = MatchPrefix(frame)
	if frame.Flow == 0 {
		l1 = append(l1, r1)
		goto block1
	}
	frame.Recover(checkpoint1)
	ret = &MatchSequence{Matches: l1}
	return
}

func ParseMatchChoice(frame *runtime.State) (ret TextMatch) {
	var e TextMatch
	var checkpoint0 int
	var l0 []TextMatch
	var c0 rune
	var r0 TextMatch
	var l1 []TextMatch
	var checkpoint1 int
	var c1 rune
	var r1 TextMatch
	e = Sequence(frame)
	if frame.Flow == 0 {
		checkpoint0 = frame.Checkpoint()
		l0 = []TextMatch{e}
		S(frame)
		c0 = frame.Peek()
		if frame.Flow == 0 {
			if c0 == '|' {
				frame.Consume()
				S(frame)
				r0 = Sequence(frame)
				if frame.Flow == 0 {
					l1 = append(l0, r0)
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
			r1 = Sequence(frame)
			if frame.Flow == 0 {
				l1 = append(l1, r1)
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
	ret = &MatchChoice{Matches: l1}
	return
block3:
	frame.Recover(checkpoint0)
	ret = e
	return
}

func ParseExprList(frame *runtime.State) (ret []ASTExpr) {
	var exprs0 []ASTExpr
	var checkpoint0 int
	var r0 ASTExpr
	var exprs1 []ASTExpr
	var checkpoint1 int
	var c rune
	var r1 ASTExpr
	var exprs2 []ASTExpr
	exprs0 = []ASTExpr{}
	checkpoint0 = frame.Checkpoint()
	r0 = ParseExpr(frame)
	if frame.Flow == 0 {
		exprs1 = append(exprs0, r0)
		goto block1
	}
	frame.Recover(checkpoint0)
	exprs2 = exprs0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r1 = ParseExpr(frame)
			if frame.Flow == 0 {
				exprs1 = append(exprs1, r1)
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
	exprs2 = exprs1
	goto block3
block3:
	ret = exprs2
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
	var name *Id
	var c rune
	var r ASTExpr
	name = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		c = frame.Peek()
		if frame.Flow == 0 {
			if c == ':' {
				frame.Consume()
				S(frame)
				r = ParseExpr(frame)
				if frame.Flow == 0 {
					ret = &NamedExpr{Name: name, Expr: r}
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
	var exprs0 []*NamedExpr
	var checkpoint0 int
	var r0 *NamedExpr
	var exprs1 []*NamedExpr
	var checkpoint1 int
	var c rune
	var r1 *NamedExpr
	var exprs2 []*NamedExpr
	exprs0 = []*NamedExpr{}
	checkpoint0 = frame.Checkpoint()
	r0 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		exprs1 = append(exprs0, r0)
		goto block1
	}
	frame.Recover(checkpoint0)
	exprs2 = exprs0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r1 = ParseNamedExpr(frame)
			if frame.Flow == 0 {
				exprs1 = append(exprs1, r1)
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
	exprs2 = exprs1
	goto block3
block3:
	ret = exprs2
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
	var t0 ASTTypeRef
	var c7 rune
	var e0 ASTExpr
	var c8 rune
	var t1 ASTTypeRef
	var c9 rune
	var args0 []*NamedExpr
	var c10 rune
	var t2 *ListTypeRef
	var c11 rune
	var args1 []ASTExpr
	var c12 rune
	var r1 *StringMatch
	var r2 *RuneMatch
	var c13 rune
	var e1 ASTExpr
	var c14 rune
	var r3 *NameRef
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
																t0 = ParseTypeRef(frame)
																if frame.Flow == 0 {
																	S(frame)
																	c7 = frame.Peek()
																	if frame.Flow == 0 {
																		if c7 == ',' {
																			frame.Consume()
																			S(frame)
																			e0 = ParseExpr(frame)
																			if frame.Flow == 0 {
																				S(frame)
																				c8 = frame.Peek()
																				if frame.Flow == 0 {
																					if c8 == ')' {
																						frame.Consume()
																						ret = &Coerce{Type: t0, Expr: e0}
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
	t1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c9 = frame.Peek()
		if frame.Flow == 0 {
			if c9 == '{' {
				frame.Consume()
				S(frame)
				args0 = ParseNamedExprList(frame)
				S(frame)
				c10 = frame.Peek()
				if frame.Flow == 0 {
					if c10 == '}' {
						frame.Consume()
						ret = &Construct{Type: t1, Args: args0}
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
	t2 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		S(frame)
		c11 = frame.Peek()
		if frame.Flow == 0 {
			if c11 == '{' {
				frame.Consume()
				S(frame)
				args1 = ParseExprList(frame)
				S(frame)
				c12 = frame.Peek()
				if frame.Flow == 0 {
					if c12 == '}' {
						frame.Consume()
						ret = &ConstructList{Type: t2, Args: args1}
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
	r1 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		ret = r1
		return
	}
	frame.Recover(checkpoint)
	r2 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		ret = r2
		return
	}
	frame.Recover(checkpoint)
	c13 = frame.Peek()
	if frame.Flow == 0 {
		if c13 == '(' {
			frame.Consume()
			S(frame)
			e1 = ParseExpr(frame)
			if frame.Flow == 0 {
				S(frame)
				c14 = frame.Peek()
				if frame.Flow == 0 {
					if c14 == ')' {
						frame.Consume()
						ret = e1
						return
					}
					frame.Fail()
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
block4:
	frame.Recover(checkpoint)
	r3 = ParseNameRef(frame)
	if frame.Flow == 0 {
		ret = r3
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

func PrimaryExprPostfix(frame *runtime.State) (ret ASTExpr) {
	var e0 ASTExpr
	var e1 ASTExpr
	var checkpoint int
	var pos int
	var c0 rune
	var args []ASTExpr
	var c1 rune
	e0 = PrimaryExpr(frame)
	if frame.Flow == 0 {
		e1 = e0
		goto block1
	}
	return
block1:
	checkpoint = frame.Checkpoint()
	S(frame)
	pos = frame.Checkpoint()
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '(' {
			frame.Consume()
			S(frame)
			args = ParseExprList(frame)
			S(frame)
			c1 = frame.Peek()
			if frame.Flow == 0 {
				if c1 == ')' {
					frame.Consume()
					e1 = &Call{Expr: e1, Pos: pos, Args: args}
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
block2:
	frame.Recover(checkpoint)
	ret = e1
	return
}

func ParseBinaryOp(frame *runtime.State, min_prec int) (ret ASTExpr) {
	var e0 ASTExpr
	var e1 ASTExpr
	var checkpoint int
	var opPos int
	var op string
	var prec int
	var r ASTExpr
	e0 = PrimaryExprPostfix(frame)
	if frame.Flow == 0 {
		e1 = e0
		goto block1
	}
	return
block1:
	checkpoint = frame.Checkpoint()
	S(frame)
	opPos = frame.Checkpoint()
	op, prec = BinaryOperator(frame)
	if frame.Flow == 0 {
		if prec < min_prec {
			frame.Fail()
			goto block2
		}
		S(frame)
		r = ParseBinaryOp(frame, prec+1)
		if frame.Flow == 0 {
			e1 = &BinaryOp{Left: e1, Op: op, OpPos: opPos, Right: r}
			goto block1
		}
		goto block2
	}
	goto block2
block2:
	frame.Recover(checkpoint)
	ret = e1
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
	var block0 []ASTExpr
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var block1 []ASTExpr
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var c12 rune
	var c13 rune
	var r0 []ASTExpr
	var blocks0 [][]ASTExpr
	var c14 rune
	var c15 rune
	var r1 []ASTExpr
	var blocks1 [][]ASTExpr
	var checkpoint1 int
	var c16 rune
	var c17 rune
	var r2 []ASTExpr
	var c18 rune
	var c19 rune
	var c20 rune
	var c21 rune
	var c22 rune
	var c23 rune
	var c24 rune
	var c25 rune
	var block2 []ASTExpr
	var c26 rune
	var c27 rune
	var expr ASTExpr
	var block3 []ASTExpr
	var else_0 []ASTExpr
	var checkpoint2 int
	var c28 rune
	var c29 rune
	var c30 rune
	var c31 rune
	var else_1 []ASTExpr
	var else_2 []ASTExpr
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
										block0 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											ret = &Repeat{Block: block0, Min: 0}
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
										block1 = ParseCodeBlock(frame)
										if frame.Flow == 0 {
											ret = &Repeat{Block: block1, Min: 1}
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
														r0 = ParseCodeBlock(frame)
														if frame.Flow == 0 {
															blocks0 = [][]ASTExpr{r0}
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
																				r1 = ParseCodeBlock(frame)
																				if frame.Flow == 0 {
																					blocks1 = append(blocks0, r1)
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
						r2 = ParseCodeBlock(frame)
						if frame.Flow == 0 {
							blocks1 = append(blocks1, r2)
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
	ret = &Choice{Blocks: blocks1}
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
																		block2 = ParseCodeBlock(frame)
																		if frame.Flow == 0 {
																			ret = &Optional{Block: block2}
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
						expr = ParseExpr(frame)
						if frame.Flow == 0 {
							S(frame)
							block3 = ParseCodeBlock(frame)
							if frame.Flow == 0 {
								else_0 = []ASTExpr{}
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
																	else_1 = ParseCodeBlock(frame)
																	if frame.Flow == 0 {
																		else_2 = else_1
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
	else_2 = else_0
	goto block8
block8:
	ret = &If{Expr: expr, Block: block3, Else: else_2}
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
	var r ASTExpr
	var c0 rune
	var c1 rune
	var c2 rune
	var name *NameRef
	var t ASTTypeRef
	var expr0 ASTExpr
	var checkpoint1 int
	var c3 rune
	var expr1 ASTExpr
	var expr2 ASTExpr
	var c4 rune
	var c5 rune
	var c6 rune
	var c7 rune
	var pos int
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var c12 rune
	var c13 rune
	var exprs []ASTExpr
	var names []ASTExpr
	var defined0 bool
	var checkpoint2 int
	var c14 rune
	var c15 rune
	var defined1 bool
	var c16 rune
	var expr3 ASTExpr
	var e ASTExpr
	checkpoint0 = frame.Checkpoint()
	r = ParseCompoundStatement(frame)
	if frame.Flow == 0 {
		ret = r
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
								name = ParseNameRef(frame)
								if frame.Flow == 0 {
									S(frame)
									t = ParseTypeRef(frame)
									if frame.Flow == 0 {
										expr0 = nil
										checkpoint1 = frame.Checkpoint()
										S(frame)
										c3 = frame.Peek()
										if frame.Flow == 0 {
											if c3 == '=' {
												frame.Consume()
												S(frame)
												expr1 = ParseExpr(frame)
												if frame.Flow == 0 {
													expr2 = expr1
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
	expr2 = expr0
	goto block2
block2:
	EOS(frame)
	if frame.Flow == 0 {
		ret = &Assign{Expr: expr2, Targets: []ASTExpr{name}, Type: t, Define: true}
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
	pos = frame.Checkpoint()
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
														exprs = ParseExprList(frame)
														EOS(frame)
														if frame.Flow == 0 {
															ret = &Return{Pos: pos, Exprs: exprs}
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
	names = ParseTargetList(frame)
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
	expr3 = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret = &Assign{Expr: expr3, Targets: names, Define: defined1}
			return
		}
		goto block8
	}
	goto block8
block8:
	frame.Recover(checkpoint0)
	e = ParseExpr(frame)
	if frame.Flow == 0 {
		EOS(frame)
		if frame.Flow == 0 {
			ret = e
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
	var r ASTExpr
	var exprs1 []ASTExpr
	var exprs2 []ASTExpr
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
	r = ParseStatement(frame)
	if frame.Flow == 0 {
		exprs1 = append(exprs0, r)
		S(frame)
		exprs0 = exprs1
		goto block1
	}
	exprs2 = exprs0
	frame.Recover(checkpoint)
	c1 = frame.Peek()
	if frame.Flow == 0 {
		if c1 == '}' {
			frame.Consume()
			ret = exprs2
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseParenthTypeList(frame *runtime.State) (ret []ASTTypeRef) {
	var c0 rune
	var types0 []ASTTypeRef
	var checkpoint0 int
	var r0 ASTTypeRef
	var types1 []ASTTypeRef
	var types2 []ASTTypeRef
	var checkpoint1 int
	var c1 rune
	var r1 ASTTypeRef
	var types3 []ASTTypeRef
	var types4 []ASTTypeRef
	var types5 []ASTTypeRef
	var types6 []ASTTypeRef
	var c2 rune
	c0 = frame.Peek()
	if frame.Flow == 0 {
		if c0 == '(' {
			frame.Consume()
			S(frame)
			types0 = []ASTTypeRef{}
			checkpoint0 = frame.Checkpoint()
			r0 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				types1 = append(types0, r0)
				S(frame)
				types2 = types1
				goto block1
			}
			types6 = types0
			frame.Recover(checkpoint0)
			types5 = types6
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
			r1 = ParseTypeRef(frame)
			if frame.Flow == 0 {
				types3 = append(types2, r1)
				S(frame)
				types2 = types3
				goto block1
			}
			types4 = types2
			goto block2
		}
		frame.Fail()
		types4 = types2
		goto block2
	}
	types4 = types2
	goto block2
block2:
	frame.Recover(checkpoint1)
	types5 = types4
	goto block3
block3:
	c2 = frame.Peek()
	if frame.Flow == 0 {
		if c2 == ')' {
			frame.Consume()
			ret = types5
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
	var name *Id
	var c_b bool
	var checkpoint0 int
	var c6 rune
	var c7 rune
	var c8 rune
	var c9 rune
	var c10 rune
	var c11 rune
	var scoped bool
	var contains0 []ASTTypeRef
	var checkpoint1 int
	var c12 rune
	var c13 rune
	var c14 rune
	var c15 rune
	var c16 rune
	var c17 rune
	var c18 rune
	var c19 rune
	var contains1 []ASTTypeRef
	var contains2 []ASTTypeRef
	var contains3 []ASTTypeRef
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
	var impl1 ASTTypeRef
	var impl2 ASTTypeRef
	var impl3 ASTTypeRef
	var c30 rune
	var fields []*FieldDecl
	var checkpoint3 int
	var fn *Id
	var ft ASTTypeRef
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
														name = Ident(frame)
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
	contains0 = []ASTTypeRef{}
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
																		contains1 = ParseParenthTypeList(frame)
																		if frame.Flow == 0 {
																			S(frame)
																			contains2 = contains1
																			goto block4
																		}
																		contains3 = contains0
																		goto block3
																	}
																	contains3 = contains0
																	goto block3
																}
																frame.Fail()
																contains3 = contains0
																goto block3
															}
															contains3 = contains0
															goto block3
														}
														frame.Fail()
														contains3 = contains0
														goto block3
													}
													contains3 = contains0
													goto block3
												}
												frame.Fail()
												contains3 = contains0
												goto block3
											}
											contains3 = contains0
											goto block3
										}
										frame.Fail()
										contains3 = contains0
										goto block3
									}
									contains3 = contains0
									goto block3
								}
								frame.Fail()
								contains3 = contains0
								goto block3
							}
							contains3 = contains0
							goto block3
						}
						frame.Fail()
						contains3 = contains0
						goto block3
					}
					contains3 = contains0
					goto block3
				}
				frame.Fail()
				contains3 = contains0
				goto block3
			}
			contains3 = contains0
			goto block3
		}
		frame.Fail()
		contains3 = contains0
		goto block3
	}
	contains3 = contains0
	goto block3
block3:
	frame.Recover(checkpoint1)
	contains2 = contains3
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
																						impl1 = ParseTypeRef(frame)
																						if frame.Flow == 0 {
																							S(frame)
																							impl2 = impl1
																							goto block6
																						}
																						impl3 = impl0
																						goto block5
																					}
																					impl3 = impl0
																					goto block5
																				}
																				frame.Fail()
																				impl3 = impl0
																				goto block5
																			}
																			impl3 = impl0
																			goto block5
																		}
																		frame.Fail()
																		impl3 = impl0
																		goto block5
																	}
																	impl3 = impl0
																	goto block5
																}
																frame.Fail()
																impl3 = impl0
																goto block5
															}
															impl3 = impl0
															goto block5
														}
														frame.Fail()
														impl3 = impl0
														goto block5
													}
													impl3 = impl0
													goto block5
												}
												frame.Fail()
												impl3 = impl0
												goto block5
											}
											impl3 = impl0
											goto block5
										}
										frame.Fail()
										impl3 = impl0
										goto block5
									}
									impl3 = impl0
									goto block5
								}
								frame.Fail()
								impl3 = impl0
								goto block5
							}
							impl3 = impl0
							goto block5
						}
						frame.Fail()
						impl3 = impl0
						goto block5
					}
					impl3 = impl0
					goto block5
				}
				frame.Fail()
				impl3 = impl0
				goto block5
			}
			impl3 = impl0
			goto block5
		}
		frame.Fail()
		impl3 = impl0
		goto block5
	}
	impl3 = impl0
	goto block5
block5:
	frame.Recover(checkpoint2)
	impl2 = impl3
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
	fn = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		ft = ParseTypeRef(frame)
		if frame.Flow == 0 {
			S(frame)
			fields = append(fields, &FieldDecl{Name: fn, Type: ft})
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
			ret = &StructDecl{Name: name, Implements: impl2, Fields: fields, Scoped: scoped, Contains: contains2}
			return
		}
		frame.Fail()
		return
	}
	return
}

func ParseParam(frame *runtime.State) (ret *Param) {
	var name *Id
	var type0 ASTTypeRef
	name = Ident(frame)
	if frame.Flow == 0 {
		S(frame)
		type0 = ParseTypeRef(frame)
		if frame.Flow == 0 {
			ret = &Param{Name: name, Type: type0}
			return
		}
		return
	}
	return
}

func ParseParamList(frame *runtime.State) (ret []*Param) {
	var params0 []*Param
	var checkpoint0 int
	var r0 *Param
	var params1 []*Param
	var checkpoint1 int
	var c rune
	var r1 *Param
	var params2 []*Param
	params0 = []*Param{}
	checkpoint0 = frame.Checkpoint()
	r0 = ParseParam(frame)
	if frame.Flow == 0 {
		params1 = append(params0, r0)
		goto block1
	}
	frame.Recover(checkpoint0)
	params2 = params0
	goto block3
block1:
	checkpoint1 = frame.Checkpoint()
	S(frame)
	c = frame.Peek()
	if frame.Flow == 0 {
		if c == ',' {
			frame.Consume()
			S(frame)
			r1 = ParseParam(frame)
			if frame.Flow == 0 {
				params1 = append(params1, r1)
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
	params2 = params1
	goto block3
block3:
	ret = params2
	return
}

func ParseFuncDecl(frame *runtime.State) (ret *FuncDecl) {
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var name *Id
	var c4 rune
	var params []*Param
	var c5 rune
	var retTypes []ASTTypeRef
	var block []ASTExpr
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
										name = Ident(frame)
										if frame.Flow == 0 {
											S(frame)
											c4 = frame.Peek()
											if frame.Flow == 0 {
												if c4 == '(' {
													frame.Consume()
													S(frame)
													params = ParseParamList(frame)
													S(frame)
													c5 = frame.Peek()
													if frame.Flow == 0 {
														if c5 == ')' {
															frame.Consume()
															S(frame)
															retTypes = ParseReturnTypeList(frame)
															S(frame)
															block = ParseCodeBlock(frame)
															if frame.Flow == 0 {
																ret = &FuncDecl{Name: name, Params: params, ReturnTypes: retTypes, Block: block, LocalInfo_Scope: &LocalInfo_Scope{}}
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
	var name *Id
	var rule ASTExpr
	var input string
	var flow string
	var d Destructure
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
										name = Ident(frame)
										if frame.Flow == 0 {
											S(frame)
											rule = ParseExpr(frame)
											if frame.Flow == 0 {
												S(frame)
												input = DecodeString(frame)
												if frame.Flow == 0 {
													S(frame)
													flow = ParseMatchState(frame)
													S(frame)
													d = ParseDestructure(frame)
													if frame.Flow == 0 {
														ret = &Test{Name: name, Rule: rule, Input: input, Flow: flow, Destructure: d}
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
	var imports0 []*ImportDecl
	var checkpoint0 int
	var c0 rune
	var c1 rune
	var c2 rune
	var c3 rune
	var c4 rune
	var c5 rune
	var c6 rune
	var imports1 []*ImportDecl
	var checkpoint1 int
	var r *StringLiteral
	var imports2 []*ImportDecl
	var imports3 []*ImportDecl
	var c7 rune
	var imports4 []*ImportDecl
	var imports5 []*ImportDecl
	imports0 = []*ImportDecl{}
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
																imports1 = imports0
																goto block1
															}
															frame.Fail()
															imports5 = imports0
															goto block2
														}
														imports5 = imports0
														goto block2
													}
													imports5 = imports0
													goto block2
												}
												frame.Fail()
												imports5 = imports0
												goto block2
											}
											imports5 = imports0
											goto block2
										}
										frame.Fail()
										imports5 = imports0
										goto block2
									}
									imports5 = imports0
									goto block2
								}
								frame.Fail()
								imports5 = imports0
								goto block2
							}
							imports5 = imports0
							goto block2
						}
						frame.Fail()
						imports5 = imports0
						goto block2
					}
					imports5 = imports0
					goto block2
				}
				frame.Fail()
				imports5 = imports0
				goto block2
			}
			imports5 = imports0
			goto block2
		}
		frame.Fail()
		imports5 = imports0
		goto block2
	}
	imports5 = imports0
	goto block2
block1:
	checkpoint1 = frame.Checkpoint()
	r = ParseStringLiteral(frame)
	if frame.Flow == 0 {
		imports2 = append(imports1, &ImportDecl{Path: r})
		S(frame)
		imports1 = imports2
		goto block1
	}
	imports3 = imports1
	frame.Recover(checkpoint1)
	c7 = frame.Peek()
	if frame.Flow == 0 {
		if c7 == ')' {
			frame.Consume()
			imports4 = imports3
			goto block3
		}
		frame.Fail()
		imports5 = imports3
		goto block2
	}
	imports5 = imports3
	goto block2
block2:
	frame.Recover(checkpoint0)
	imports4 = imports5
	goto block3
block3:
	ret = imports4
	return
}

func ParseFile(frame *runtime.State) (ret *File) {
	var decls0 []ASTDecl
	var tests0 []*Test
	var imports []*ImportDecl
	var decls1 []ASTDecl
	var tests1 []*Test
	var checkpoint0 int
	var checkpoint1 int
	var r0 *FuncDecl
	var decls2 []ASTDecl
	var tests2 []*Test
	var r1 *StructDecl
	var r2 *Test
	var decls3 []ASTDecl
	var tests3 []*Test
	var checkpoint2 int
	decls0 = []ASTDecl{}
	tests0 = []*Test{}
	S(frame)
	imports = ParseImports(frame)
	S(frame)
	decls1, tests1 = decls0, tests0
	goto block1
block1:
	checkpoint0 = frame.Checkpoint()
	checkpoint1 = frame.Checkpoint()
	r0 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		decls2, tests2 = append(decls1, r0), tests1
		goto block2
	}
	frame.Recover(checkpoint1)
	r1 = ParseStructDecl(frame)
	if frame.Flow == 0 {
		decls2, tests2 = append(decls1, r1), tests1
		goto block2
	}
	frame.Recover(checkpoint1)
	r2 = ParseTest(frame)
	if frame.Flow == 0 {
		decls2, tests2 = decls1, append(tests1, r2)
		goto block2
	}
	decls3, tests3 = decls1, tests1
	frame.Recover(checkpoint0)
	checkpoint2 = frame.LookaheadBegin()
	frame.Peek()
	if frame.Flow == 0 {
		frame.Consume()
		frame.LookaheadFail(checkpoint2)
		return
	}
	frame.LookaheadNormal(checkpoint2)
	ret = &File{Imports: imports, Decls: decls3, Tests: tests3}
	return
block2:
	S(frame)
	decls1, tests1 = decls2, tests2
	goto block1
}
