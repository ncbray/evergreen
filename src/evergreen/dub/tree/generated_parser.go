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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block17
	}
block3:
	r2 = ' '
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block15
	} else {
		goto block6
	}
block6:
	r4 = '\t'
	goto block7
block7:
	r5 = r1 == r4
	goto block8
block8:
	if r5 {
		goto block15
	} else {
		goto block9
	}
block9:
	r6 = '\r'
	goto block10
block10:
	r7 = r1 == r6
	goto block11
block11:
	if r7 {
		goto block15
	} else {
		goto block12
	}
block12:
	r8 = '\n'
	goto block13
block13:
	r9 = r1 == r8
	goto block14
block14:
	if r9 {
		goto block15
	} else {
		goto block16
	}
block15:
	frame.Consume()
	goto block1
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r0)
	goto block18
block18:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.LookaheadBegin()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block27
	}
block3:
	r2 = 'a'
	goto block4
block4:
	r3 = r1 >= r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block9
	}
block6:
	r4 = 'z'
	goto block7
block7:
	r5 = r1 <= r4
	goto block8
block8:
	if r5 {
		goto block24
	} else {
		goto block9
	}
block9:
	r6 = 'A'
	goto block10
block10:
	r7 = r1 >= r6
	goto block11
block11:
	if r7 {
		goto block12
	} else {
		goto block15
	}
block12:
	r8 = 'Z'
	goto block13
block13:
	r9 = r1 <= r8
	goto block14
block14:
	if r9 {
		goto block24
	} else {
		goto block15
	}
block15:
	r10 = '_'
	goto block16
block16:
	r11 = r1 == r10
	goto block17
block17:
	if r11 {
		goto block24
	} else {
		goto block18
	}
block18:
	r12 = '0'
	goto block19
block19:
	r13 = r1 >= r12
	goto block20
block20:
	if r13 {
		goto block21
	} else {
		goto block26
	}
block21:
	r14 = '9'
	goto block22
block22:
	r15 = r1 <= r14
	goto block23
block23:
	if r15 {
		goto block24
	} else {
		goto block26
	}
block24:
	frame.Consume()
	goto block25
block25:
	frame.LookaheadFail(r0)
	goto block30
block26:
	frame.Fail()
	goto block27
block27:
	frame.LookaheadNormal(r0)
	goto block28
block28:
	S(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block30
	}
block29:
	return
block30:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block50
	}
block3:
	r2 = 'a'
	goto block4
block4:
	r3 = r1 >= r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block9
	}
block6:
	r4 = 'z'
	goto block7
block7:
	r5 = r1 <= r4
	goto block8
block8:
	if r5 {
		goto block18
	} else {
		goto block9
	}
block9:
	r6 = 'A'
	goto block10
block10:
	r7 = r1 >= r6
	goto block11
block11:
	if r7 {
		goto block12
	} else {
		goto block15
	}
block12:
	r8 = 'Z'
	goto block13
block13:
	r9 = r1 <= r8
	goto block14
block14:
	if r9 {
		goto block18
	} else {
		goto block15
	}
block15:
	r10 = '_'
	goto block16
block16:
	r11 = r1 == r10
	goto block17
block17:
	if r11 {
		goto block18
	} else {
		goto block49
	}
block18:
	frame.Consume()
	goto block19
block19:
	r12 = frame.Checkpoint()
	goto block20
block20:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block44
	}
block21:
	r14 = 'a'
	goto block22
block22:
	r15 = r13 >= r14
	goto block23
block23:
	if r15 {
		goto block24
	} else {
		goto block27
	}
block24:
	r16 = 'z'
	goto block25
block25:
	r17 = r13 <= r16
	goto block26
block26:
	if r17 {
		goto block42
	} else {
		goto block27
	}
block27:
	r18 = 'A'
	goto block28
block28:
	r19 = r13 >= r18
	goto block29
block29:
	if r19 {
		goto block30
	} else {
		goto block33
	}
block30:
	r20 = 'Z'
	goto block31
block31:
	r21 = r13 <= r20
	goto block32
block32:
	if r21 {
		goto block42
	} else {
		goto block33
	}
block33:
	r22 = '_'
	goto block34
block34:
	r23 = r13 == r22
	goto block35
block35:
	if r23 {
		goto block42
	} else {
		goto block36
	}
block36:
	r24 = '0'
	goto block37
block37:
	r25 = r13 >= r24
	goto block38
block38:
	if r25 {
		goto block39
	} else {
		goto block43
	}
block39:
	r26 = '9'
	goto block40
block40:
	r27 = r13 <= r26
	goto block41
block41:
	if r27 {
		goto block42
	} else {
		goto block43
	}
block42:
	frame.Consume()
	goto block19
block43:
	frame.Fail()
	goto block44
block44:
	frame.Recover(r12)
	goto block45
block45:
	r28 = frame.Slice(r0)
	goto block46
block46:
	S(frame)
	if frame.Flow == 0 {
		goto block47
	} else {
		goto block50
	}
block47:
	ret0 = r28
	goto block48
block48:
	return
block49:
	frame.Fail()
	goto block50
block50:
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
	goto block0
block0:
	goto block1
block1:
	r0 = 0
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block40
	}
block3:
	r2 = '0'
	goto block4
block4:
	r3 = r1 >= r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block39
	}
block6:
	r4 = '9'
	goto block7
block7:
	r5 = r1 <= r4
	goto block8
block8:
	if r5 {
		goto block9
	} else {
		goto block39
	}
block9:
	frame.Consume()
	goto block10
block10:
	r6 = int(r1)
	goto block11
block11:
	r7 = '0'
	goto block12
block12:
	r8 = int(r7)
	goto block13
block13:
	r9 = r6 - r8
	goto block14
block14:
	r10 = 10
	goto block15
block15:
	r11 = r0 * r10
	goto block16
block16:
	r12 = r11 + r9
	goto block17
block17:
	r13 = r12
	goto block18
block18:
	r14 = frame.Checkpoint()
	goto block19
block19:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block36
	}
block20:
	r16 = '0'
	goto block21
block21:
	r17 = r15 >= r16
	goto block22
block22:
	if r17 {
		goto block23
	} else {
		goto block35
	}
block23:
	r18 = '9'
	goto block24
block24:
	r19 = r15 <= r18
	goto block25
block25:
	if r19 {
		goto block26
	} else {
		goto block35
	}
block26:
	frame.Consume()
	goto block27
block27:
	r20 = int(r15)
	goto block28
block28:
	r21 = '0'
	goto block29
block29:
	r22 = int(r21)
	goto block30
block30:
	r23 = r20 - r22
	goto block31
block31:
	r24 = 10
	goto block32
block32:
	r25 = r13 * r24
	goto block33
block33:
	r26 = r25 + r23
	goto block34
block34:
	r13 = r26
	goto block18
block35:
	frame.Fail()
	goto block36
block36:
	frame.Recover(r14)
	goto block37
block37:
	ret0 = r13
	goto block38
block38:
	return
block39:
	frame.Fail()
	goto block40
block40:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block10
	}
block3:
	r2 = 'a'
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block9
	}
block6:
	frame.Consume()
	goto block7
block7:
	r4 = '\a'
	goto block8
block8:
	ret0 = r4
	goto block90
block9:
	frame.Fail()
	goto block10
block10:
	frame.Recover(r0)
	goto block11
block11:
	r5 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block19
	}
block12:
	r6 = 'b'
	goto block13
block13:
	r7 = r5 == r6
	goto block14
block14:
	if r7 {
		goto block15
	} else {
		goto block18
	}
block15:
	frame.Consume()
	goto block16
block16:
	r8 = '\b'
	goto block17
block17:
	ret0 = r8
	goto block90
block18:
	frame.Fail()
	goto block19
block19:
	frame.Recover(r0)
	goto block20
block20:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block28
	}
block21:
	r10 = 'f'
	goto block22
block22:
	r11 = r9 == r10
	goto block23
block23:
	if r11 {
		goto block24
	} else {
		goto block27
	}
block24:
	frame.Consume()
	goto block25
block25:
	r12 = '\f'
	goto block26
block26:
	ret0 = r12
	goto block90
block27:
	frame.Fail()
	goto block28
block28:
	frame.Recover(r0)
	goto block29
block29:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block30
	} else {
		goto block37
	}
block30:
	r14 = 'n'
	goto block31
block31:
	r15 = r13 == r14
	goto block32
block32:
	if r15 {
		goto block33
	} else {
		goto block36
	}
block33:
	frame.Consume()
	goto block34
block34:
	r16 = '\n'
	goto block35
block35:
	ret0 = r16
	goto block90
block36:
	frame.Fail()
	goto block37
block37:
	frame.Recover(r0)
	goto block38
block38:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		goto block39
	} else {
		goto block46
	}
block39:
	r18 = 'r'
	goto block40
block40:
	r19 = r17 == r18
	goto block41
block41:
	if r19 {
		goto block42
	} else {
		goto block45
	}
block42:
	frame.Consume()
	goto block43
block43:
	r20 = '\r'
	goto block44
block44:
	ret0 = r20
	goto block90
block45:
	frame.Fail()
	goto block46
block46:
	frame.Recover(r0)
	goto block47
block47:
	r21 = frame.Peek()
	if frame.Flow == 0 {
		goto block48
	} else {
		goto block55
	}
block48:
	r22 = 't'
	goto block49
block49:
	r23 = r21 == r22
	goto block50
block50:
	if r23 {
		goto block51
	} else {
		goto block54
	}
block51:
	frame.Consume()
	goto block52
block52:
	r24 = '\t'
	goto block53
block53:
	ret0 = r24
	goto block90
block54:
	frame.Fail()
	goto block55
block55:
	frame.Recover(r0)
	goto block56
block56:
	r25 = frame.Peek()
	if frame.Flow == 0 {
		goto block57
	} else {
		goto block64
	}
block57:
	r26 = 'v'
	goto block58
block58:
	r27 = r25 == r26
	goto block59
block59:
	if r27 {
		goto block60
	} else {
		goto block63
	}
block60:
	frame.Consume()
	goto block61
block61:
	r28 = '\v'
	goto block62
block62:
	ret0 = r28
	goto block90
block63:
	frame.Fail()
	goto block64
block64:
	frame.Recover(r0)
	goto block65
block65:
	r29 = frame.Peek()
	if frame.Flow == 0 {
		goto block66
	} else {
		goto block73
	}
block66:
	r30 = '\\'
	goto block67
block67:
	r31 = r29 == r30
	goto block68
block68:
	if r31 {
		goto block69
	} else {
		goto block72
	}
block69:
	frame.Consume()
	goto block70
block70:
	r32 = '\\'
	goto block71
block71:
	ret0 = r32
	goto block90
block72:
	frame.Fail()
	goto block73
block73:
	frame.Recover(r0)
	goto block74
block74:
	r33 = frame.Peek()
	if frame.Flow == 0 {
		goto block75
	} else {
		goto block82
	}
block75:
	r34 = '\''
	goto block76
block76:
	r35 = r33 == r34
	goto block77
block77:
	if r35 {
		goto block78
	} else {
		goto block81
	}
block78:
	frame.Consume()
	goto block79
block79:
	r36 = '\''
	goto block80
block80:
	ret0 = r36
	goto block90
block81:
	frame.Fail()
	goto block82
block82:
	frame.Recover(r0)
	goto block83
block83:
	r37 = frame.Peek()
	if frame.Flow == 0 {
		goto block84
	} else {
		goto block92
	}
block84:
	r38 = '"'
	goto block85
block85:
	r39 = r37 == r38
	goto block86
block86:
	if r39 {
		goto block87
	} else {
		goto block91
	}
block87:
	frame.Consume()
	goto block88
block88:
	r40 = '"'
	goto block89
block89:
	ret0 = r40
	goto block90
block90:
	return
block91:
	frame.Fail()
	goto block92
block92:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block42
	}
block2:
	r1 = '"'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block41
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = []rune{}
	goto block7
block7:
	r4 = r3
	goto block8
block8:
	r5 = frame.Checkpoint()
	goto block9
block9:
	r6 = frame.Checkpoint()
	goto block10
block10:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block21
	}
block11:
	r8 = '"'
	goto block12
block12:
	r9 = r7 == r8
	goto block13
block13:
	if r9 {
		goto block17
	} else {
		goto block14
	}
block14:
	r10 = '\\'
	goto block15
block15:
	r11 = r7 == r10
	goto block16
block16:
	if r11 {
		goto block17
	} else {
		goto block18
	}
block17:
	frame.Fail()
	goto block21
block18:
	frame.Consume()
	goto block19
block19:
	r12 = append(r4, r7)
	goto block20
block20:
	r4 = r12
	goto block8
block21:
	frame.Recover(r6)
	goto block22
block22:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block31
	}
block23:
	r14 = '\\'
	goto block24
block24:
	r15 = r13 == r14
	goto block25
block25:
	if r15 {
		goto block26
	} else {
		goto block30
	}
block26:
	frame.Consume()
	goto block27
block27:
	r16 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block28
	} else {
		goto block31
	}
block28:
	r17 = append(r4, r16)
	goto block29
block29:
	r4 = r17
	goto block8
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(r5)
	goto block32
block32:
	r18 = frame.Peek()
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block42
	}
block33:
	r19 = '"'
	goto block34
block34:
	r20 = r18 == r19
	goto block35
block35:
	if r20 {
		goto block36
	} else {
		goto block40
	}
block36:
	frame.Consume()
	goto block37
block37:
	r21 = string(r4)
	goto block38
block38:
	ret0 = r21
	goto block39
block39:
	return
block40:
	frame.Fail()
	goto block42
block41:
	frame.Fail()
	goto block42
block42:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block35
	}
block2:
	r1 = '\''
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block34
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = frame.Checkpoint()
	goto block7
block7:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block17
	}
block8:
	r5 = '\\'
	goto block9
block9:
	r6 = r4 == r5
	goto block10
block10:
	if r6 {
		goto block14
	} else {
		goto block11
	}
block11:
	r7 = '\''
	goto block12
block12:
	r8 = r4 == r7
	goto block13
block13:
	if r8 {
		goto block14
	} else {
		goto block15
	}
block14:
	frame.Fail()
	goto block17
block15:
	frame.Consume()
	goto block16
block16:
	r9 = r4
	goto block25
block17:
	frame.Recover(r3)
	goto block18
block18:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block35
	}
block19:
	r11 = '\\'
	goto block20
block20:
	r12 = r10 == r11
	goto block21
block21:
	if r12 {
		goto block22
	} else {
		goto block33
	}
block22:
	frame.Consume()
	goto block23
block23:
	r13 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block35
	}
block24:
	r9 = r13
	goto block25
block25:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block35
	}
block26:
	r15 = '\''
	goto block27
block27:
	r16 = r14 == r15
	goto block28
block28:
	if r16 {
		goto block29
	} else {
		goto block32
	}
block29:
	frame.Consume()
	goto block30
block30:
	ret0 = r9
	goto block31
block31:
	return
block32:
	frame.Fail()
	goto block35
block33:
	frame.Fail()
	goto block35
block34:
	frame.Fail()
	goto block35
block35:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block28
	}
block3:
	r2 = 't'
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block27
	}
block6:
	frame.Consume()
	goto block7
block7:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block28
	}
block8:
	r5 = 'r'
	goto block9
block9:
	r6 = r4 == r5
	goto block10
block10:
	if r6 {
		goto block11
	} else {
		goto block26
	}
block11:
	frame.Consume()
	goto block12
block12:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block28
	}
block13:
	r8 = 'u'
	goto block14
block14:
	r9 = r7 == r8
	goto block15
block15:
	if r9 {
		goto block16
	} else {
		goto block25
	}
block16:
	frame.Consume()
	goto block17
block17:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block28
	}
block18:
	r11 = 'e'
	goto block19
block19:
	r12 = r10 == r11
	goto block20
block20:
	if r12 {
		goto block21
	} else {
		goto block24
	}
block21:
	frame.Consume()
	goto block22
block22:
	r13 = true
	goto block23
block23:
	ret0 = r13
	goto block56
block24:
	frame.Fail()
	goto block28
block25:
	frame.Fail()
	goto block28
block26:
	frame.Fail()
	goto block28
block27:
	frame.Fail()
	goto block28
block28:
	frame.Recover(r0)
	goto block29
block29:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block30
	} else {
		goto block62
	}
block30:
	r15 = 'f'
	goto block31
block31:
	r16 = r14 == r15
	goto block32
block32:
	if r16 {
		goto block33
	} else {
		goto block61
	}
block33:
	frame.Consume()
	goto block34
block34:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		goto block35
	} else {
		goto block62
	}
block35:
	r18 = 'a'
	goto block36
block36:
	r19 = r17 == r18
	goto block37
block37:
	if r19 {
		goto block38
	} else {
		goto block60
	}
block38:
	frame.Consume()
	goto block39
block39:
	r20 = frame.Peek()
	if frame.Flow == 0 {
		goto block40
	} else {
		goto block62
	}
block40:
	r21 = 'l'
	goto block41
block41:
	r22 = r20 == r21
	goto block42
block42:
	if r22 {
		goto block43
	} else {
		goto block59
	}
block43:
	frame.Consume()
	goto block44
block44:
	r23 = frame.Peek()
	if frame.Flow == 0 {
		goto block45
	} else {
		goto block62
	}
block45:
	r24 = 's'
	goto block46
block46:
	r25 = r23 == r24
	goto block47
block47:
	if r25 {
		goto block48
	} else {
		goto block58
	}
block48:
	frame.Consume()
	goto block49
block49:
	r26 = frame.Peek()
	if frame.Flow == 0 {
		goto block50
	} else {
		goto block62
	}
block50:
	r27 = 'e'
	goto block51
block51:
	r28 = r26 == r27
	goto block52
block52:
	if r28 {
		goto block53
	} else {
		goto block57
	}
block53:
	frame.Consume()
	goto block54
block54:
	r29 = false
	goto block55
block55:
	ret0 = r29
	goto block56
block56:
	return
block57:
	frame.Fail()
	goto block62
block58:
	frame.Fail()
	goto block62
block59:
	frame.Fail()
	goto block62
block60:
	frame.Fail()
	goto block62
block61:
	frame.Fail()
	goto block62
block62:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = DecodeRune(frame)
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block8
	}
block4:
	r3 = frame.Slice(r1)
	goto block5
block5:
	S(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block8
	}
block6:
	r4 = &RuneLiteral{Text: r3, Value: r2}
	goto block7
block7:
	ret0 = r4
	goto block29
block8:
	frame.Recover(r0)
	goto block9
block9:
	r5 = frame.Checkpoint()
	goto block10
block10:
	r6 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block15
	}
block11:
	r7 = frame.Slice(r5)
	goto block12
block12:
	S(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block15
	}
block13:
	r8 = &StringLiteral{Text: r7, Value: r6}
	goto block14
block14:
	ret0 = r8
	goto block29
block15:
	frame.Recover(r0)
	goto block16
block16:
	r9 = frame.Checkpoint()
	goto block17
block17:
	r10 = DecodeInt(frame)
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block22
	}
block18:
	r11 = frame.Slice(r9)
	goto block19
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block22
	}
block20:
	r12 = &IntLiteral{Text: r11, Value: r10}
	goto block21
block21:
	ret0 = r12
	goto block29
block22:
	frame.Recover(r0)
	goto block23
block23:
	r13 = frame.Checkpoint()
	goto block24
block24:
	r14 = DecodeBool(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block30
	}
block25:
	r15 = frame.Slice(r13)
	goto block26
block26:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block30
	}
block27:
	r16 = &BoolLiteral{Text: r15, Value: r14}
	goto block28
block28:
	ret0 = r16
	goto block29
block29:
	return
block30:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block18
	}
block4:
	r3 = '+'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block16
	} else {
		goto block7
	}
block7:
	r5 = '-'
	goto block8
block8:
	r6 = r2 == r5
	goto block9
block9:
	if r6 {
		goto block16
	} else {
		goto block10
	}
block10:
	r7 = '*'
	goto block11
block11:
	r8 = r2 == r7
	goto block12
block12:
	if r8 {
		goto block16
	} else {
		goto block13
	}
block13:
	r9 = '/'
	goto block14
block14:
	r10 = r2 == r9
	goto block15
block15:
	if r10 {
		goto block16
	} else {
		goto block17
	}
block16:
	frame.Consume()
	goto block80
block17:
	frame.Fail()
	goto block18
block18:
	frame.Recover(r1)
	goto block19
block19:
	r11 = frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block36
	}
block20:
	r12 = '<'
	goto block21
block21:
	r13 = r11 == r12
	goto block22
block22:
	if r13 {
		goto block26
	} else {
		goto block23
	}
block23:
	r14 = '>'
	goto block24
block24:
	r15 = r11 == r14
	goto block25
block25:
	if r15 {
		goto block26
	} else {
		goto block35
	}
block26:
	frame.Consume()
	goto block27
block27:
	r16 = frame.Checkpoint()
	goto block28
block28:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block34
	}
block29:
	r18 = '='
	goto block30
block30:
	r19 = r17 == r18
	goto block31
block31:
	if r19 {
		goto block32
	} else {
		goto block33
	}
block32:
	frame.Consume()
	goto block80
block33:
	frame.Fail()
	goto block34
block34:
	frame.Recover(r16)
	goto block80
block35:
	frame.Fail()
	goto block36
block36:
	frame.Recover(r1)
	goto block37
block37:
	r20 = frame.Peek()
	if frame.Flow == 0 {
		goto block38
	} else {
		goto block86
	}
block38:
	r21 = '!'
	goto block39
block39:
	r22 = r20 == r21
	goto block40
block40:
	if r22 {
		goto block44
	} else {
		goto block41
	}
block41:
	r23 = '='
	goto block42
block42:
	r24 = r20 == r23
	goto block43
block43:
	if r24 {
		goto block44
	} else {
		goto block85
	}
block44:
	frame.Consume()
	goto block45
block45:
	r25 = frame.Peek()
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block86
	}
block46:
	r26 = '='
	goto block47
block47:
	r27 = r25 == r26
	goto block48
block48:
	if r27 {
		goto block49
	} else {
		goto block84
	}
block49:
	frame.Consume()
	goto block50
block50:
	r28 = frame.LookaheadBegin()
	goto block51
block51:
	r29 = frame.Peek()
	if frame.Flow == 0 {
		goto block52
	} else {
		goto block79
	}
block52:
	r30 = '+'
	goto block53
block53:
	r31 = r29 == r30
	goto block54
block54:
	if r31 {
		goto block76
	} else {
		goto block55
	}
block55:
	r32 = '-'
	goto block56
block56:
	r33 = r29 == r32
	goto block57
block57:
	if r33 {
		goto block76
	} else {
		goto block58
	}
block58:
	r34 = '*'
	goto block59
block59:
	r35 = r29 == r34
	goto block60
block60:
	if r35 {
		goto block76
	} else {
		goto block61
	}
block61:
	r36 = '/'
	goto block62
block62:
	r37 = r29 == r36
	goto block63
block63:
	if r37 {
		goto block76
	} else {
		goto block64
	}
block64:
	r38 = '<'
	goto block65
block65:
	r39 = r29 == r38
	goto block66
block66:
	if r39 {
		goto block76
	} else {
		goto block67
	}
block67:
	r40 = '>'
	goto block68
block68:
	r41 = r29 == r40
	goto block69
block69:
	if r41 {
		goto block76
	} else {
		goto block70
	}
block70:
	r42 = '!'
	goto block71
block71:
	r43 = r29 == r42
	goto block72
block72:
	if r43 {
		goto block76
	} else {
		goto block73
	}
block73:
	r44 = '='
	goto block74
block74:
	r45 = r29 == r44
	goto block75
block75:
	if r45 {
		goto block76
	} else {
		goto block78
	}
block76:
	frame.Consume()
	goto block77
block77:
	frame.LookaheadFail(r28)
	goto block86
block78:
	frame.Fail()
	goto block79
block79:
	frame.LookaheadNormal(r28)
	goto block80
block80:
	r46 = frame.Slice(r0)
	goto block81
block81:
	S(frame)
	if frame.Flow == 0 {
		goto block82
	} else {
		goto block86
	}
block82:
	ret0 = r46
	goto block83
block83:
	return
block84:
	frame.Fail()
	goto block86
block85:
	frame.Fail()
	goto block86
block86:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block19
	}
block2:
	r1 = '/'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block18
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block19
	}
block7:
	r3 = ParseMatchChoice(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block19
	}
block8:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block19
	}
block9:
	r5 = '/'
	goto block10
block10:
	r6 = r4 == r5
	goto block11
block11:
	if r6 {
		goto block12
	} else {
		goto block17
	}
block12:
	frame.Consume()
	goto block13
block13:
	S(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block19
	}
block14:
	r7 = &StringMatch{Match: r3}
	goto block15
block15:
	ret0 = r7
	goto block16
block16:
	return
block17:
	frame.Fail()
	goto block19
block18:
	frame.Fail()
	goto block19
block19:
	return
}
func RuneMatchExpr(frame *runtime.State) (ret0 *RuneMatch) {
	var r0 rune
	var r1 rune
	var r2 bool
	var r3 *RuneRangeMatch
	var r4 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block12
	}
block2:
	r1 = '$'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block11
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block12
	}
block7:
	r3 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block12
	}
block8:
	r4 = &RuneMatch{Match: r3}
	goto block9
block9:
	ret0 = r4
	goto block10
block10:
	return
block11:
	frame.Fail()
	goto block12
block12:
	return
}
func ParseStructTypeRef(frame *runtime.State) (ret0 *TypeRef) {
	var r0 string
	var r1 *TypeRef
	goto block0
block0:
	goto block1
block1:
	r0 = Ident(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block5
	}
block2:
	r1 = &TypeRef{Name: r0}
	goto block3
block3:
	ret0 = r1
	goto block4
block4:
	return
block5:
	return
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block17
	}
block2:
	r1 = '['
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block16
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block17
	}
block7:
	r4 = ']'
	goto block8
block8:
	r5 = r3 == r4
	goto block9
block9:
	if r5 {
		goto block10
	} else {
		goto block15
	}
block10:
	frame.Consume()
	goto block11
block11:
	r6 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block17
	}
block12:
	r7 = &ListTypeRef{Type: r6}
	goto block13
block13:
	ret0 = r7
	goto block14
block14:
	return
block15:
	frame.Fail()
	goto block17
block16:
	frame.Fail()
	goto block17
block17:
	return
}
func ParseTypeRef(frame *runtime.State) (ret0 ASTTypeRef) {
	var r0 int
	var r1 *TypeRef
	var r2 *ListTypeRef
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r1
	goto block7
block4:
	frame.Recover(r0)
	goto block5
block5:
	r2 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block8
	}
block6:
	ret0 = r2
	goto block7
block7:
	return
block8:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block35
	}
block3:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block35
	}
block4:
	r3 = '{'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block34
	}
block7:
	frame.Consume()
	goto block8
block8:
	S(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block35
	}
block9:
	r5 = []*DestructureField{}
	goto block10
block10:
	r6 = r5
	goto block11
block11:
	r7 = frame.Checkpoint()
	goto block12
block12:
	r8 = Ident(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block24
	}
block13:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block24
	}
block14:
	r10 = ':'
	goto block15
block15:
	r11 = r9 == r10
	goto block16
block16:
	if r11 {
		goto block17
	} else {
		goto block23
	}
block17:
	frame.Consume()
	goto block18
block18:
	S(frame)
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block24
	}
block19:
	r12 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block24
	}
block20:
	r13 = &DestructureField{Name: r8, Destructure: r12}
	goto block21
block21:
	r14 = append(r6, r13)
	goto block22
block22:
	r6 = r14
	goto block11
block23:
	frame.Fail()
	goto block24
block24:
	frame.Recover(r7)
	goto block25
block25:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block35
	}
block26:
	r16 = '}'
	goto block27
block27:
	r17 = r15 == r16
	goto block28
block28:
	if r17 {
		goto block29
	} else {
		goto block33
	}
block29:
	frame.Consume()
	goto block30
block30:
	S(frame)
	if frame.Flow == 0 {
		goto block31
	} else {
		goto block35
	}
block31:
	r18 = &DestructureStruct{Type: r1, Args: r6}
	goto block32
block32:
	ret0 = r18
	goto block64
block33:
	frame.Fail()
	goto block35
block34:
	frame.Fail()
	goto block35
block35:
	frame.Recover(r0)
	goto block36
block36:
	r19 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block37
	} else {
		goto block60
	}
block37:
	r20 = frame.Peek()
	if frame.Flow == 0 {
		goto block38
	} else {
		goto block60
	}
block38:
	r21 = '{'
	goto block39
block39:
	r22 = r20 == r21
	goto block40
block40:
	if r22 {
		goto block41
	} else {
		goto block59
	}
block41:
	frame.Consume()
	goto block42
block42:
	S(frame)
	if frame.Flow == 0 {
		goto block43
	} else {
		goto block60
	}
block43:
	r23 = []Destructure{}
	goto block44
block44:
	r24 = r23
	goto block45
block45:
	r25 = frame.Checkpoint()
	goto block46
block46:
	r26 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block47
	} else {
		goto block49
	}
block47:
	r27 = append(r24, r26)
	goto block48
block48:
	r24 = r27
	goto block45
block49:
	frame.Recover(r25)
	goto block50
block50:
	r28 = frame.Peek()
	if frame.Flow == 0 {
		goto block51
	} else {
		goto block60
	}
block51:
	r29 = '}'
	goto block52
block52:
	r30 = r28 == r29
	goto block53
block53:
	if r30 {
		goto block54
	} else {
		goto block58
	}
block54:
	frame.Consume()
	goto block55
block55:
	S(frame)
	if frame.Flow == 0 {
		goto block56
	} else {
		goto block60
	}
block56:
	r31 = &DestructureList{Type: r19, Args: r24}
	goto block57
block57:
	ret0 = r31
	goto block64
block58:
	frame.Fail()
	goto block60
block59:
	frame.Fail()
	goto block60
block60:
	frame.Recover(r0)
	goto block61
block61:
	r32 = Literal(frame)
	if frame.Flow == 0 {
		goto block62
	} else {
		goto block65
	}
block62:
	r33 = &DestructureValue{Expr: r32}
	goto block63
block63:
	ret0 = r33
	goto block64
block64:
	return
block65:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block15
	}
block3:
	r2 = ']'
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block12
	} else {
		goto block6
	}
block6:
	r4 = '-'
	goto block7
block7:
	r5 = r1 == r4
	goto block8
block8:
	if r5 {
		goto block12
	} else {
		goto block9
	}
block9:
	r6 = '\\'
	goto block10
block10:
	r7 = r1 == r6
	goto block11
block11:
	if r7 {
		goto block12
	} else {
		goto block13
	}
block12:
	frame.Fail()
	goto block15
block13:
	frame.Consume()
	goto block14
block14:
	ret0 = r1
	goto block28
block15:
	frame.Recover(r0)
	goto block16
block16:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block30
	}
block17:
	r9 = '\\'
	goto block18
block18:
	r10 = r8 == r9
	goto block19
block19:
	if r10 {
		goto block20
	} else {
		goto block29
	}
block20:
	frame.Consume()
	goto block21
block21:
	r11 = frame.Checkpoint()
	goto block22
block22:
	r12 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block24
	}
block23:
	ret0 = r12
	goto block28
block24:
	frame.Recover(r11)
	goto block25
block25:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block30
	}
block26:
	frame.Consume()
	goto block27
block27:
	ret0 = r13
	goto block28
block28:
	return
block29:
	frame.Fail()
	goto block30
block30:
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
	goto block0
block0:
	goto block1
block1:
	r0 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block16
	}
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block11
	}
block4:
	r3 = '-'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block10
	}
block7:
	frame.Consume()
	goto block8
block8:
	r5 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block11
	}
block9:
	r6 = r5
	goto block13
block10:
	frame.Fail()
	goto block11
block11:
	frame.Recover(r1)
	goto block12
block12:
	r6 = r0
	goto block13
block13:
	r7 = &RuneFilter{Min: r0, Max: r6}
	goto block14
block14:
	ret0 = r7
	goto block15
block15:
	return
block16:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block35
	}
block2:
	r1 = '['
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block34
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = false
	goto block7
block7:
	r4 = []*RuneFilter{}
	goto block8
block8:
	r5 = frame.Checkpoint()
	goto block9
block9:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block17
	}
block10:
	r7 = '^'
	goto block11
block11:
	r8 = r6 == r7
	goto block12
block12:
	if r8 {
		goto block13
	} else {
		goto block16
	}
block13:
	frame.Consume()
	goto block14
block14:
	r9 = true
	goto block15
block15:
	r10, r11 = r9, r4
	goto block19
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r5)
	goto block18
block18:
	r10, r11 = r3, r4
	goto block19
block19:
	r12 = frame.Checkpoint()
	goto block20
block20:
	r13 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block23
	}
block21:
	r14 = append(r11, r13)
	goto block22
block22:
	r10, r11 = r10, r14
	goto block19
block23:
	frame.Recover(r12)
	goto block24
block24:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block35
	}
block25:
	r16 = ']'
	goto block26
block26:
	r17 = r15 == r16
	goto block27
block27:
	if r17 {
		goto block28
	} else {
		goto block33
	}
block28:
	frame.Consume()
	goto block29
block29:
	S(frame)
	if frame.Flow == 0 {
		goto block30
	} else {
		goto block35
	}
block30:
	r18 = &RuneRangeMatch{Invert: r10, Filters: r11}
	goto block31
block31:
	ret0 = r18
	goto block32
block32:
	return
block33:
	frame.Fail()
	goto block35
block34:
	frame.Fail()
	goto block35
block35:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r1
	goto block24
block4:
	frame.Recover(r0)
	goto block5
block5:
	r2 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block9
	}
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block9
	}
block7:
	r3 = &StringLiteralMatch{Value: r2}
	goto block8
block8:
	ret0 = r3
	goto block24
block9:
	frame.Recover(r0)
	goto block10
block10:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block27
	}
block11:
	r5 = '('
	goto block12
block12:
	r6 = r4 == r5
	goto block13
block13:
	if r6 {
		goto block14
	} else {
		goto block26
	}
block14:
	frame.Consume()
	goto block15
block15:
	S(frame)
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block27
	}
block16:
	r7 = ParseMatchChoice(frame)
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block27
	}
block17:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block27
	}
block18:
	r9 = ')'
	goto block19
block19:
	r10 = r8 == r9
	goto block20
block20:
	if r10 {
		goto block21
	} else {
		goto block25
	}
block21:
	frame.Consume()
	goto block22
block22:
	S(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block27
	}
block23:
	ret0 = r7
	goto block24
block24:
	return
block25:
	frame.Fail()
	goto block27
block26:
	frame.Fail()
	goto block27
block27:
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
	goto block0
block0:
	goto block1
block1:
	r0 = Atom(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block40
	}
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block13
	}
block4:
	r3 = '*'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block12
	}
block7:
	frame.Consume()
	goto block8
block8:
	S(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block13
	}
block9:
	r5 = 0
	goto block10
block10:
	r6 = &MatchRepeat{Match: r0, Min: r5}
	goto block11
block11:
	ret0 = r6
	goto block39
block12:
	frame.Fail()
	goto block13
block13:
	frame.Recover(r1)
	goto block14
block14:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block24
	}
block15:
	r8 = '+'
	goto block16
block16:
	r9 = r7 == r8
	goto block17
block17:
	if r9 {
		goto block18
	} else {
		goto block23
	}
block18:
	frame.Consume()
	goto block19
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block24
	}
block20:
	r10 = 1
	goto block21
block21:
	r11 = &MatchRepeat{Match: r0, Min: r10}
	goto block22
block22:
	ret0 = r11
	goto block39
block23:
	frame.Fail()
	goto block24
block24:
	frame.Recover(r1)
	goto block25
block25:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block37
	}
block26:
	r13 = '?'
	goto block27
block27:
	r14 = r12 == r13
	goto block28
block28:
	if r14 {
		goto block29
	} else {
		goto block36
	}
block29:
	frame.Consume()
	goto block30
block30:
	S(frame)
	if frame.Flow == 0 {
		goto block31
	} else {
		goto block37
	}
block31:
	r15 = []TextMatch{}
	goto block32
block32:
	r16 = &MatchSequence{Matches: r15}
	goto block33
block33:
	r17 = []TextMatch{r0, r16}
	goto block34
block34:
	r18 = &MatchChoice{Matches: r17}
	goto block35
block35:
	ret0 = r18
	goto block39
block36:
	frame.Fail()
	goto block37
block37:
	frame.Recover(r1)
	goto block38
block38:
	ret0 = r0
	goto block39
block39:
	return
block40:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = false
	goto block3
block3:
	r2 = frame.Checkpoint()
	goto block4
block4:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block12
	}
block5:
	r4 = '!'
	goto block6
block6:
	r5 = r3 == r4
	goto block7
block7:
	if r5 {
		goto block8
	} else {
		goto block11
	}
block8:
	frame.Consume()
	goto block9
block9:
	r6 = true
	goto block10
block10:
	r7 = r6
	goto block19
block11:
	frame.Fail()
	goto block12
block12:
	frame.Recover(r2)
	goto block13
block13:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block24
	}
block14:
	r9 = '&'
	goto block15
block15:
	r10 = r8 == r9
	goto block16
block16:
	if r10 {
		goto block17
	} else {
		goto block23
	}
block17:
	frame.Consume()
	goto block18
block18:
	r7 = r1
	goto block19
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block24
	}
block20:
	r11 = MatchPostfix(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block24
	}
block21:
	r12 = &MatchLookahead{Invert: r7, Match: r11}
	goto block22
block22:
	ret0 = r12
	goto block27
block23:
	frame.Fail()
	goto block24
block24:
	frame.Recover(r0)
	goto block25
block25:
	r13 = MatchPostfix(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block28
	}
block26:
	ret0 = r13
	goto block27
block27:
	return
block28:
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
	goto block0
block0:
	goto block1
block1:
	r0 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block17
	}
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = []TextMatch{r0}
	goto block4
block4:
	r3 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block14
	}
block5:
	r4 = append(r2, r3)
	goto block6
block6:
	r5 = r4
	goto block7
block7:
	r6 = frame.Checkpoint()
	goto block8
block8:
	r7 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block11
	}
block9:
	r8 = append(r5, r7)
	goto block10
block10:
	r5 = r8
	goto block7
block11:
	frame.Recover(r6)
	goto block12
block12:
	r9 = &MatchSequence{Matches: r5}
	goto block13
block13:
	ret0 = r9
	goto block16
block14:
	frame.Recover(r1)
	goto block15
block15:
	ret0 = r0
	goto block16
block16:
	return
block17:
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
	goto block0
block0:
	goto block1
block1:
	r0 = Sequence(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block31
	}
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = []TextMatch{r0}
	goto block4
block4:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block28
	}
block5:
	r4 = '|'
	goto block6
block6:
	r5 = r3 == r4
	goto block7
block7:
	if r5 {
		goto block8
	} else {
		goto block27
	}
block8:
	frame.Consume()
	goto block9
block9:
	S(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block28
	}
block10:
	r6 = Sequence(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block28
	}
block11:
	r7 = append(r2, r6)
	goto block12
block12:
	r8 = r7
	goto block13
block13:
	r9 = frame.Checkpoint()
	goto block14
block14:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block24
	}
block15:
	r11 = '|'
	goto block16
block16:
	r12 = r10 == r11
	goto block17
block17:
	if r12 {
		goto block18
	} else {
		goto block23
	}
block18:
	frame.Consume()
	goto block19
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block24
	}
block20:
	r13 = Sequence(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block24
	}
block21:
	r14 = append(r8, r13)
	goto block22
block22:
	r8 = r14
	goto block13
block23:
	frame.Fail()
	goto block24
block24:
	frame.Recover(r9)
	goto block25
block25:
	r15 = &MatchChoice{Matches: r8}
	goto block26
block26:
	ret0 = r15
	goto block30
block27:
	frame.Fail()
	goto block28
block28:
	frame.Recover(r1)
	goto block29
block29:
	ret0 = r0
	goto block30
block30:
	return
block31:
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
	goto block0
block0:
	goto block1
block1:
	r0 = []ASTExpr{}
	goto block2
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block19
	}
block4:
	r3 = append(r0, r2)
	goto block5
block5:
	r4 = r3
	goto block6
block6:
	r5 = frame.Checkpoint()
	goto block7
block7:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block17
	}
block8:
	r7 = ','
	goto block9
block9:
	r8 = r6 == r7
	goto block10
block10:
	if r8 {
		goto block11
	} else {
		goto block16
	}
block11:
	frame.Consume()
	goto block12
block12:
	S(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block17
	}
block13:
	r9 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block17
	}
block14:
	r10 = append(r4, r9)
	goto block15
block15:
	r4 = r10
	goto block6
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r5)
	goto block18
block18:
	r11 = r4
	goto block21
block19:
	frame.Recover(r1)
	goto block20
block20:
	r11 = r0
	goto block21
block21:
	ret0 = r11
	goto block22
block22:
	return
}
func ParseNamedExpr(frame *runtime.State) (ret0 *NamedExpr) {
	var r0 string
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 ASTExpr
	var r5 *NamedExpr
	goto block0
block0:
	goto block1
block1:
	r0 = Ident(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block13
	}
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block13
	}
block3:
	r2 = ':'
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block12
	}
block6:
	frame.Consume()
	goto block7
block7:
	S(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block13
	}
block8:
	r4 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block13
	}
block9:
	r5 = &NamedExpr{Name: r0, Expr: r4}
	goto block10
block10:
	ret0 = r5
	goto block11
block11:
	return
block12:
	frame.Fail()
	goto block13
block13:
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
	goto block0
block0:
	goto block1
block1:
	r0 = []*NamedExpr{}
	goto block2
block2:
	r1 = frame.Checkpoint()
	goto block3
block3:
	r2 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block19
	}
block4:
	r3 = append(r0, r2)
	goto block5
block5:
	r4 = r3
	goto block6
block6:
	r5 = frame.Checkpoint()
	goto block7
block7:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block17
	}
block8:
	r7 = ','
	goto block9
block9:
	r8 = r6 == r7
	goto block10
block10:
	if r8 {
		goto block11
	} else {
		goto block16
	}
block11:
	frame.Consume()
	goto block12
block12:
	S(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block17
	}
block13:
	r9 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block17
	}
block14:
	r10 = append(r4, r9)
	goto block15
block15:
	r4 = r10
	goto block6
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r5)
	goto block18
block18:
	r11 = r4
	goto block21
block19:
	frame.Recover(r1)
	goto block20
block20:
	r11 = r0
	goto block21
block21:
	ret0 = r11
	goto block22
block22:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block37
	}
block2:
	r1 = '('
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block36
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block37
	}
block7:
	r3 = []ASTTypeRef{}
	goto block8
block8:
	r4 = frame.Checkpoint()
	goto block9
block9:
	r5 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block25
	}
block10:
	r6 = append(r3, r5)
	goto block11
block11:
	r7 = r6
	goto block12
block12:
	r8 = frame.Checkpoint()
	goto block13
block13:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block23
	}
block14:
	r10 = ','
	goto block15
block15:
	r11 = r9 == r10
	goto block16
block16:
	if r11 {
		goto block17
	} else {
		goto block22
	}
block17:
	frame.Consume()
	goto block18
block18:
	S(frame)
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block23
	}
block19:
	r12 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block23
	}
block20:
	r13 = append(r7, r12)
	goto block21
block21:
	r7 = r13
	goto block12
block22:
	frame.Fail()
	goto block23
block23:
	frame.Recover(r8)
	goto block24
block24:
	r14 = r7
	goto block27
block25:
	frame.Recover(r4)
	goto block26
block26:
	r14 = r3
	goto block27
block27:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block28
	} else {
		goto block37
	}
block28:
	r16 = ')'
	goto block29
block29:
	r17 = r15 == r16
	goto block30
block30:
	if r17 {
		goto block31
	} else {
		goto block35
	}
block31:
	frame.Consume()
	goto block32
block32:
	S(frame)
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block37
	}
block33:
	ret0 = r14
	goto block34
block34:
	return
block35:
	frame.Fail()
	goto block37
block36:
	frame.Fail()
	goto block37
block37:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Checkpoint()
	goto block2
block2:
	r1 = Literal(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r1
	goto block549
block4:
	frame.Recover(r0)
	goto block5
block5:
	r2 = 0
	goto block6
block6:
	r3 = frame.Checkpoint()
	goto block7
block7:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block32
	}
block8:
	r5 = 's'
	goto block9
block9:
	r6 = r4 == r5
	goto block10
block10:
	if r6 {
		goto block11
	} else {
		goto block31
	}
block11:
	frame.Consume()
	goto block12
block12:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block32
	}
block13:
	r8 = 't'
	goto block14
block14:
	r9 = r7 == r8
	goto block15
block15:
	if r9 {
		goto block16
	} else {
		goto block30
	}
block16:
	frame.Consume()
	goto block17
block17:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block32
	}
block18:
	r11 = 'a'
	goto block19
block19:
	r12 = r10 == r11
	goto block20
block20:
	if r12 {
		goto block21
	} else {
		goto block29
	}
block21:
	frame.Consume()
	goto block22
block22:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block32
	}
block23:
	r14 = 'r'
	goto block24
block24:
	r15 = r13 == r14
	goto block25
block25:
	if r15 {
		goto block26
	} else {
		goto block28
	}
block26:
	frame.Consume()
	goto block27
block27:
	r16 = r2
	goto block55
block28:
	frame.Fail()
	goto block32
block29:
	frame.Fail()
	goto block32
block30:
	frame.Fail()
	goto block32
block31:
	frame.Fail()
	goto block32
block32:
	frame.Recover(r3)
	goto block33
block33:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		goto block34
	} else {
		goto block63
	}
block34:
	r18 = 'p'
	goto block35
block35:
	r19 = r17 == r18
	goto block36
block36:
	if r19 {
		goto block37
	} else {
		goto block62
	}
block37:
	frame.Consume()
	goto block38
block38:
	r20 = frame.Peek()
	if frame.Flow == 0 {
		goto block39
	} else {
		goto block63
	}
block39:
	r21 = 'l'
	goto block40
block40:
	r22 = r20 == r21
	goto block41
block41:
	if r22 {
		goto block42
	} else {
		goto block61
	}
block42:
	frame.Consume()
	goto block43
block43:
	r23 = frame.Peek()
	if frame.Flow == 0 {
		goto block44
	} else {
		goto block63
	}
block44:
	r24 = 'u'
	goto block45
block45:
	r25 = r23 == r24
	goto block46
block46:
	if r25 {
		goto block47
	} else {
		goto block60
	}
block47:
	frame.Consume()
	goto block48
block48:
	r26 = frame.Peek()
	if frame.Flow == 0 {
		goto block49
	} else {
		goto block63
	}
block49:
	r27 = 's'
	goto block50
block50:
	r28 = r26 == r27
	goto block51
block51:
	if r28 {
		goto block52
	} else {
		goto block59
	}
block52:
	frame.Consume()
	goto block53
block53:
	r29 = 1
	goto block54
block54:
	r16 = r29
	goto block55
block55:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block56
	} else {
		goto block63
	}
block56:
	r30 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block57
	} else {
		goto block63
	}
block57:
	r31 = &Repeat{Block: r30, Min: r16}
	goto block58
block58:
	ret0 = r31
	goto block549
block59:
	frame.Fail()
	goto block63
block60:
	frame.Fail()
	goto block63
block61:
	frame.Fail()
	goto block63
block62:
	frame.Fail()
	goto block63
block63:
	frame.Recover(r0)
	goto block64
block64:
	r32 = frame.Peek()
	if frame.Flow == 0 {
		goto block65
	} else {
		goto block124
	}
block65:
	r33 = 'c'
	goto block66
block66:
	r34 = r32 == r33
	goto block67
block67:
	if r34 {
		goto block68
	} else {
		goto block123
	}
block68:
	frame.Consume()
	goto block69
block69:
	r35 = frame.Peek()
	if frame.Flow == 0 {
		goto block70
	} else {
		goto block124
	}
block70:
	r36 = 'h'
	goto block71
block71:
	r37 = r35 == r36
	goto block72
block72:
	if r37 {
		goto block73
	} else {
		goto block122
	}
block73:
	frame.Consume()
	goto block74
block74:
	r38 = frame.Peek()
	if frame.Flow == 0 {
		goto block75
	} else {
		goto block124
	}
block75:
	r39 = 'o'
	goto block76
block76:
	r40 = r38 == r39
	goto block77
block77:
	if r40 {
		goto block78
	} else {
		goto block121
	}
block78:
	frame.Consume()
	goto block79
block79:
	r41 = frame.Peek()
	if frame.Flow == 0 {
		goto block80
	} else {
		goto block124
	}
block80:
	r42 = 'o'
	goto block81
block81:
	r43 = r41 == r42
	goto block82
block82:
	if r43 {
		goto block83
	} else {
		goto block120
	}
block83:
	frame.Consume()
	goto block84
block84:
	r44 = frame.Peek()
	if frame.Flow == 0 {
		goto block85
	} else {
		goto block124
	}
block85:
	r45 = 's'
	goto block86
block86:
	r46 = r44 == r45
	goto block87
block87:
	if r46 {
		goto block88
	} else {
		goto block119
	}
block88:
	frame.Consume()
	goto block89
block89:
	r47 = frame.Peek()
	if frame.Flow == 0 {
		goto block90
	} else {
		goto block124
	}
block90:
	r48 = 'e'
	goto block91
block91:
	r49 = r47 == r48
	goto block92
block92:
	if r49 {
		goto block93
	} else {
		goto block118
	}
block93:
	frame.Consume()
	goto block94
block94:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block95
	} else {
		goto block124
	}
block95:
	r50 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block96
	} else {
		goto block124
	}
block96:
	r51 = [][]ASTExpr{r50}
	goto block97
block97:
	r52 = r51
	goto block98
block98:
	r53 = frame.Checkpoint()
	goto block99
block99:
	r54 = frame.Peek()
	if frame.Flow == 0 {
		goto block100
	} else {
		goto block115
	}
block100:
	r55 = 'o'
	goto block101
block101:
	r56 = r54 == r55
	goto block102
block102:
	if r56 {
		goto block103
	} else {
		goto block114
	}
block103:
	frame.Consume()
	goto block104
block104:
	r57 = frame.Peek()
	if frame.Flow == 0 {
		goto block105
	} else {
		goto block115
	}
block105:
	r58 = 'r'
	goto block106
block106:
	r59 = r57 == r58
	goto block107
block107:
	if r59 {
		goto block108
	} else {
		goto block113
	}
block108:
	frame.Consume()
	goto block109
block109:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block110
	} else {
		goto block115
	}
block110:
	r60 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block111
	} else {
		goto block115
	}
block111:
	r61 = append(r52, r60)
	goto block112
block112:
	r52 = r61
	goto block98
block113:
	frame.Fail()
	goto block115
block114:
	frame.Fail()
	goto block115
block115:
	frame.Recover(r53)
	goto block116
block116:
	r62 = &Choice{Blocks: r52}
	goto block117
block117:
	ret0 = r62
	goto block549
block118:
	frame.Fail()
	goto block124
block119:
	frame.Fail()
	goto block124
block120:
	frame.Fail()
	goto block124
block121:
	frame.Fail()
	goto block124
block122:
	frame.Fail()
	goto block124
block123:
	frame.Fail()
	goto block124
block124:
	frame.Recover(r0)
	goto block125
block125:
	r63 = frame.Peek()
	if frame.Flow == 0 {
		goto block126
	} else {
		goto block177
	}
block126:
	r64 = 'q'
	goto block127
block127:
	r65 = r63 == r64
	goto block128
block128:
	if r65 {
		goto block129
	} else {
		goto block176
	}
block129:
	frame.Consume()
	goto block130
block130:
	r66 = frame.Peek()
	if frame.Flow == 0 {
		goto block131
	} else {
		goto block177
	}
block131:
	r67 = 'u'
	goto block132
block132:
	r68 = r66 == r67
	goto block133
block133:
	if r68 {
		goto block134
	} else {
		goto block175
	}
block134:
	frame.Consume()
	goto block135
block135:
	r69 = frame.Peek()
	if frame.Flow == 0 {
		goto block136
	} else {
		goto block177
	}
block136:
	r70 = 'e'
	goto block137
block137:
	r71 = r69 == r70
	goto block138
block138:
	if r71 {
		goto block139
	} else {
		goto block174
	}
block139:
	frame.Consume()
	goto block140
block140:
	r72 = frame.Peek()
	if frame.Flow == 0 {
		goto block141
	} else {
		goto block177
	}
block141:
	r73 = 's'
	goto block142
block142:
	r74 = r72 == r73
	goto block143
block143:
	if r74 {
		goto block144
	} else {
		goto block173
	}
block144:
	frame.Consume()
	goto block145
block145:
	r75 = frame.Peek()
	if frame.Flow == 0 {
		goto block146
	} else {
		goto block177
	}
block146:
	r76 = 't'
	goto block147
block147:
	r77 = r75 == r76
	goto block148
block148:
	if r77 {
		goto block149
	} else {
		goto block172
	}
block149:
	frame.Consume()
	goto block150
block150:
	r78 = frame.Peek()
	if frame.Flow == 0 {
		goto block151
	} else {
		goto block177
	}
block151:
	r79 = 'i'
	goto block152
block152:
	r80 = r78 == r79
	goto block153
block153:
	if r80 {
		goto block154
	} else {
		goto block171
	}
block154:
	frame.Consume()
	goto block155
block155:
	r81 = frame.Peek()
	if frame.Flow == 0 {
		goto block156
	} else {
		goto block177
	}
block156:
	r82 = 'o'
	goto block157
block157:
	r83 = r81 == r82
	goto block158
block158:
	if r83 {
		goto block159
	} else {
		goto block170
	}
block159:
	frame.Consume()
	goto block160
block160:
	r84 = frame.Peek()
	if frame.Flow == 0 {
		goto block161
	} else {
		goto block177
	}
block161:
	r85 = 'n'
	goto block162
block162:
	r86 = r84 == r85
	goto block163
block163:
	if r86 {
		goto block164
	} else {
		goto block169
	}
block164:
	frame.Consume()
	goto block165
block165:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block166
	} else {
		goto block177
	}
block166:
	r87 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block167
	} else {
		goto block177
	}
block167:
	r88 = &Optional{Block: r87}
	goto block168
block168:
	ret0 = r88
	goto block549
block169:
	frame.Fail()
	goto block177
block170:
	frame.Fail()
	goto block177
block171:
	frame.Fail()
	goto block177
block172:
	frame.Fail()
	goto block177
block173:
	frame.Fail()
	goto block177
block174:
	frame.Fail()
	goto block177
block175:
	frame.Fail()
	goto block177
block176:
	frame.Fail()
	goto block177
block177:
	frame.Recover(r0)
	goto block178
block178:
	r89 = frame.Peek()
	if frame.Flow == 0 {
		goto block179
	} else {
		goto block212
	}
block179:
	r90 = 's'
	goto block180
block180:
	r91 = r89 == r90
	goto block181
block181:
	if r91 {
		goto block182
	} else {
		goto block211
	}
block182:
	frame.Consume()
	goto block183
block183:
	r92 = frame.Peek()
	if frame.Flow == 0 {
		goto block184
	} else {
		goto block212
	}
block184:
	r93 = 'l'
	goto block185
block185:
	r94 = r92 == r93
	goto block186
block186:
	if r94 {
		goto block187
	} else {
		goto block210
	}
block187:
	frame.Consume()
	goto block188
block188:
	r95 = frame.Peek()
	if frame.Flow == 0 {
		goto block189
	} else {
		goto block212
	}
block189:
	r96 = 'i'
	goto block190
block190:
	r97 = r95 == r96
	goto block191
block191:
	if r97 {
		goto block192
	} else {
		goto block209
	}
block192:
	frame.Consume()
	goto block193
block193:
	r98 = frame.Peek()
	if frame.Flow == 0 {
		goto block194
	} else {
		goto block212
	}
block194:
	r99 = 'c'
	goto block195
block195:
	r100 = r98 == r99
	goto block196
block196:
	if r100 {
		goto block197
	} else {
		goto block208
	}
block197:
	frame.Consume()
	goto block198
block198:
	r101 = frame.Peek()
	if frame.Flow == 0 {
		goto block199
	} else {
		goto block212
	}
block199:
	r102 = 'e'
	goto block200
block200:
	r103 = r101 == r102
	goto block201
block201:
	if r103 {
		goto block202
	} else {
		goto block207
	}
block202:
	frame.Consume()
	goto block203
block203:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block204
	} else {
		goto block212
	}
block204:
	r104 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block205
	} else {
		goto block212
	}
block205:
	r105 = &Slice{Block: r104}
	goto block206
block206:
	ret0 = r105
	goto block549
block207:
	frame.Fail()
	goto block212
block208:
	frame.Fail()
	goto block212
block209:
	frame.Fail()
	goto block212
block210:
	frame.Fail()
	goto block212
block211:
	frame.Fail()
	goto block212
block212:
	frame.Recover(r0)
	goto block213
block213:
	r106 = frame.Peek()
	if frame.Flow == 0 {
		goto block214
	} else {
		goto block230
	}
block214:
	r107 = 'i'
	goto block215
block215:
	r108 = r106 == r107
	goto block216
block216:
	if r108 {
		goto block217
	} else {
		goto block229
	}
block217:
	frame.Consume()
	goto block218
block218:
	r109 = frame.Peek()
	if frame.Flow == 0 {
		goto block219
	} else {
		goto block230
	}
block219:
	r110 = 'f'
	goto block220
block220:
	r111 = r109 == r110
	goto block221
block221:
	if r111 {
		goto block222
	} else {
		goto block228
	}
block222:
	frame.Consume()
	goto block223
block223:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block224
	} else {
		goto block230
	}
block224:
	r112 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block225
	} else {
		goto block230
	}
block225:
	r113 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block226
	} else {
		goto block230
	}
block226:
	r114 = &If{Expr: r112, Block: r113}
	goto block227
block227:
	ret0 = r114
	goto block549
block228:
	frame.Fail()
	goto block230
block229:
	frame.Fail()
	goto block230
block230:
	frame.Recover(r0)
	goto block231
block231:
	r115 = frame.Peek()
	if frame.Flow == 0 {
		goto block232
	} else {
		goto block268
	}
block232:
	r116 = 'v'
	goto block233
block233:
	r117 = r115 == r116
	goto block234
block234:
	if r117 {
		goto block235
	} else {
		goto block267
	}
block235:
	frame.Consume()
	goto block236
block236:
	r118 = frame.Peek()
	if frame.Flow == 0 {
		goto block237
	} else {
		goto block268
	}
block237:
	r119 = 'a'
	goto block238
block238:
	r120 = r118 == r119
	goto block239
block239:
	if r120 {
		goto block240
	} else {
		goto block266
	}
block240:
	frame.Consume()
	goto block241
block241:
	r121 = frame.Peek()
	if frame.Flow == 0 {
		goto block242
	} else {
		goto block268
	}
block242:
	r122 = 'r'
	goto block243
block243:
	r123 = r121 == r122
	goto block244
block244:
	if r123 {
		goto block245
	} else {
		goto block265
	}
block245:
	frame.Consume()
	goto block246
block246:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block247
	} else {
		goto block268
	}
block247:
	r124 = Ident(frame)
	if frame.Flow == 0 {
		goto block248
	} else {
		goto block268
	}
block248:
	r125 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block249
	} else {
		goto block268
	}
block249:
	r126 = nil
	goto block250
block250:
	r127 = frame.Checkpoint()
	goto block251
block251:
	r128 = frame.Peek()
	if frame.Flow == 0 {
		goto block252
	} else {
		goto block260
	}
block252:
	r129 = '='
	goto block253
block253:
	r130 = r128 == r129
	goto block254
block254:
	if r130 {
		goto block255
	} else {
		goto block259
	}
block255:
	frame.Consume()
	goto block256
block256:
	S(frame)
	if frame.Flow == 0 {
		goto block257
	} else {
		goto block260
	}
block257:
	r131 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block258
	} else {
		goto block260
	}
block258:
	r132 = r131
	goto block262
block259:
	frame.Fail()
	goto block260
block260:
	frame.Recover(r127)
	goto block261
block261:
	r132 = r126
	goto block262
block262:
	r133 = true
	goto block263
block263:
	r134 = &Assign{Expr: r132, Name: r124, Type: r125, Define: r133}
	goto block264
block264:
	ret0 = r134
	goto block549
block265:
	frame.Fail()
	goto block268
block266:
	frame.Fail()
	goto block268
block267:
	frame.Fail()
	goto block268
block268:
	frame.Recover(r0)
	goto block269
block269:
	r135 = frame.Peek()
	if frame.Flow == 0 {
		goto block270
	} else {
		goto block296
	}
block270:
	r136 = 'f'
	goto block271
block271:
	r137 = r135 == r136
	goto block272
block272:
	if r137 {
		goto block273
	} else {
		goto block295
	}
block273:
	frame.Consume()
	goto block274
block274:
	r138 = frame.Peek()
	if frame.Flow == 0 {
		goto block275
	} else {
		goto block296
	}
block275:
	r139 = 'a'
	goto block276
block276:
	r140 = r138 == r139
	goto block277
block277:
	if r140 {
		goto block278
	} else {
		goto block294
	}
block278:
	frame.Consume()
	goto block279
block279:
	r141 = frame.Peek()
	if frame.Flow == 0 {
		goto block280
	} else {
		goto block296
	}
block280:
	r142 = 'i'
	goto block281
block281:
	r143 = r141 == r142
	goto block282
block282:
	if r143 {
		goto block283
	} else {
		goto block293
	}
block283:
	frame.Consume()
	goto block284
block284:
	r144 = frame.Peek()
	if frame.Flow == 0 {
		goto block285
	} else {
		goto block296
	}
block285:
	r145 = 'l'
	goto block286
block286:
	r146 = r144 == r145
	goto block287
block287:
	if r146 {
		goto block288
	} else {
		goto block292
	}
block288:
	frame.Consume()
	goto block289
block289:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block290
	} else {
		goto block296
	}
block290:
	r147 = &Fail{}
	goto block291
block291:
	ret0 = r147
	goto block549
block292:
	frame.Fail()
	goto block296
block293:
	frame.Fail()
	goto block296
block294:
	frame.Fail()
	goto block296
block295:
	frame.Fail()
	goto block296
block296:
	frame.Recover(r0)
	goto block297
block297:
	r148 = frame.Peek()
	if frame.Flow == 0 {
		goto block298
	} else {
		goto block338
	}
block298:
	r149 = 'c'
	goto block299
block299:
	r150 = r148 == r149
	goto block300
block300:
	if r150 {
		goto block301
	} else {
		goto block337
	}
block301:
	frame.Consume()
	goto block302
block302:
	r151 = frame.Peek()
	if frame.Flow == 0 {
		goto block303
	} else {
		goto block338
	}
block303:
	r152 = 'o'
	goto block304
block304:
	r153 = r151 == r152
	goto block305
block305:
	if r153 {
		goto block306
	} else {
		goto block336
	}
block306:
	frame.Consume()
	goto block307
block307:
	r154 = frame.Peek()
	if frame.Flow == 0 {
		goto block308
	} else {
		goto block338
	}
block308:
	r155 = 'e'
	goto block309
block309:
	r156 = r154 == r155
	goto block310
block310:
	if r156 {
		goto block311
	} else {
		goto block335
	}
block311:
	frame.Consume()
	goto block312
block312:
	r157 = frame.Peek()
	if frame.Flow == 0 {
		goto block313
	} else {
		goto block338
	}
block313:
	r158 = 'r'
	goto block314
block314:
	r159 = r157 == r158
	goto block315
block315:
	if r159 {
		goto block316
	} else {
		goto block334
	}
block316:
	frame.Consume()
	goto block317
block317:
	r160 = frame.Peek()
	if frame.Flow == 0 {
		goto block318
	} else {
		goto block338
	}
block318:
	r161 = 'c'
	goto block319
block319:
	r162 = r160 == r161
	goto block320
block320:
	if r162 {
		goto block321
	} else {
		goto block333
	}
block321:
	frame.Consume()
	goto block322
block322:
	r163 = frame.Peek()
	if frame.Flow == 0 {
		goto block323
	} else {
		goto block338
	}
block323:
	r164 = 'e'
	goto block324
block324:
	r165 = r163 == r164
	goto block325
block325:
	if r165 {
		goto block326
	} else {
		goto block332
	}
block326:
	frame.Consume()
	goto block327
block327:
	S(frame)
	if frame.Flow == 0 {
		goto block328
	} else {
		goto block338
	}
block328:
	r166 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block329
	} else {
		goto block338
	}
block329:
	r167 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block330
	} else {
		goto block338
	}
block330:
	r168 = &Coerce{Type: r166, Expr: r167}
	goto block331
block331:
	ret0 = r168
	goto block549
block332:
	frame.Fail()
	goto block338
block333:
	frame.Fail()
	goto block338
block334:
	frame.Fail()
	goto block338
block335:
	frame.Fail()
	goto block338
block336:
	frame.Fail()
	goto block338
block337:
	frame.Fail()
	goto block338
block338:
	frame.Recover(r0)
	goto block339
block339:
	r169 = frame.Peek()
	if frame.Flow == 0 {
		goto block340
	} else {
		goto block382
	}
block340:
	r170 = 'a'
	goto block341
block341:
	r171 = r169 == r170
	goto block342
block342:
	if r171 {
		goto block343
	} else {
		goto block381
	}
block343:
	frame.Consume()
	goto block344
block344:
	r172 = frame.Peek()
	if frame.Flow == 0 {
		goto block345
	} else {
		goto block382
	}
block345:
	r173 = 'p'
	goto block346
block346:
	r174 = r172 == r173
	goto block347
block347:
	if r174 {
		goto block348
	} else {
		goto block380
	}
block348:
	frame.Consume()
	goto block349
block349:
	r175 = frame.Peek()
	if frame.Flow == 0 {
		goto block350
	} else {
		goto block382
	}
block350:
	r176 = 'p'
	goto block351
block351:
	r177 = r175 == r176
	goto block352
block352:
	if r177 {
		goto block353
	} else {
		goto block379
	}
block353:
	frame.Consume()
	goto block354
block354:
	r178 = frame.Peek()
	if frame.Flow == 0 {
		goto block355
	} else {
		goto block382
	}
block355:
	r179 = 'e'
	goto block356
block356:
	r180 = r178 == r179
	goto block357
block357:
	if r180 {
		goto block358
	} else {
		goto block378
	}
block358:
	frame.Consume()
	goto block359
block359:
	r181 = frame.Peek()
	if frame.Flow == 0 {
		goto block360
	} else {
		goto block382
	}
block360:
	r182 = 'n'
	goto block361
block361:
	r183 = r181 == r182
	goto block362
block362:
	if r183 {
		goto block363
	} else {
		goto block377
	}
block363:
	frame.Consume()
	goto block364
block364:
	r184 = frame.Peek()
	if frame.Flow == 0 {
		goto block365
	} else {
		goto block382
	}
block365:
	r185 = 'd'
	goto block366
block366:
	r186 = r184 == r185
	goto block367
block367:
	if r186 {
		goto block368
	} else {
		goto block376
	}
block368:
	frame.Consume()
	goto block369
block369:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block370
	} else {
		goto block382
	}
block370:
	r187 = Ident(frame)
	if frame.Flow == 0 {
		goto block371
	} else {
		goto block382
	}
block371:
	r188 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block372
	} else {
		goto block382
	}
block372:
	r189 = &GetName{Name: r187}
	goto block373
block373:
	r190 = &Append{List: r189, Expr: r188}
	goto block374
block374:
	r191 = &Assign{Expr: r190, Name: r187}
	goto block375
block375:
	ret0 = r191
	goto block549
block376:
	frame.Fail()
	goto block382
block377:
	frame.Fail()
	goto block382
block378:
	frame.Fail()
	goto block382
block379:
	frame.Fail()
	goto block382
block380:
	frame.Fail()
	goto block382
block381:
	frame.Fail()
	goto block382
block382:
	frame.Recover(r0)
	goto block383
block383:
	r192 = frame.Peek()
	if frame.Flow == 0 {
		goto block384
	} else {
		goto block447
	}
block384:
	r193 = 'r'
	goto block385
block385:
	r194 = r192 == r193
	goto block386
block386:
	if r194 {
		goto block387
	} else {
		goto block446
	}
block387:
	frame.Consume()
	goto block388
block388:
	r195 = frame.Peek()
	if frame.Flow == 0 {
		goto block389
	} else {
		goto block447
	}
block389:
	r196 = 'e'
	goto block390
block390:
	r197 = r195 == r196
	goto block391
block391:
	if r197 {
		goto block392
	} else {
		goto block445
	}
block392:
	frame.Consume()
	goto block393
block393:
	r198 = frame.Peek()
	if frame.Flow == 0 {
		goto block394
	} else {
		goto block447
	}
block394:
	r199 = 't'
	goto block395
block395:
	r200 = r198 == r199
	goto block396
block396:
	if r200 {
		goto block397
	} else {
		goto block444
	}
block397:
	frame.Consume()
	goto block398
block398:
	r201 = frame.Peek()
	if frame.Flow == 0 {
		goto block399
	} else {
		goto block447
	}
block399:
	r202 = 'u'
	goto block400
block400:
	r203 = r201 == r202
	goto block401
block401:
	if r203 {
		goto block402
	} else {
		goto block443
	}
block402:
	frame.Consume()
	goto block403
block403:
	r204 = frame.Peek()
	if frame.Flow == 0 {
		goto block404
	} else {
		goto block447
	}
block404:
	r205 = 'r'
	goto block405
block405:
	r206 = r204 == r205
	goto block406
block406:
	if r206 {
		goto block407
	} else {
		goto block442
	}
block407:
	frame.Consume()
	goto block408
block408:
	r207 = frame.Peek()
	if frame.Flow == 0 {
		goto block409
	} else {
		goto block447
	}
block409:
	r208 = 'n'
	goto block410
block410:
	r209 = r207 == r208
	goto block411
block411:
	if r209 {
		goto block412
	} else {
		goto block441
	}
block412:
	frame.Consume()
	goto block413
block413:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block414
	} else {
		goto block447
	}
block414:
	r210 = frame.Checkpoint()
	goto block415
block415:
	r211 = frame.Peek()
	if frame.Flow == 0 {
		goto block416
	} else {
		goto block432
	}
block416:
	r212 = '('
	goto block417
block417:
	r213 = r211 == r212
	goto block418
block418:
	if r213 {
		goto block419
	} else {
		goto block431
	}
block419:
	frame.Consume()
	goto block420
block420:
	S(frame)
	if frame.Flow == 0 {
		goto block421
	} else {
		goto block432
	}
block421:
	r214 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block422
	} else {
		goto block432
	}
block422:
	r215 = frame.Peek()
	if frame.Flow == 0 {
		goto block423
	} else {
		goto block432
	}
block423:
	r216 = ')'
	goto block424
block424:
	r217 = r215 == r216
	goto block425
block425:
	if r217 {
		goto block426
	} else {
		goto block430
	}
block426:
	frame.Consume()
	goto block427
block427:
	S(frame)
	if frame.Flow == 0 {
		goto block428
	} else {
		goto block432
	}
block428:
	r218 = &Return{Exprs: r214}
	goto block429
block429:
	ret0 = r218
	goto block549
block430:
	frame.Fail()
	goto block432
block431:
	frame.Fail()
	goto block432
block432:
	frame.Recover(r210)
	goto block433
block433:
	r219 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block434
	} else {
		goto block437
	}
block434:
	r220 = []ASTExpr{r219}
	goto block435
block435:
	r221 = &Return{Exprs: r220}
	goto block436
block436:
	ret0 = r221
	goto block549
block437:
	frame.Recover(r210)
	goto block438
block438:
	r222 = []ASTExpr{}
	goto block439
block439:
	r223 = &Return{Exprs: r222}
	goto block440
block440:
	ret0 = r223
	goto block549
block441:
	frame.Fail()
	goto block447
block442:
	frame.Fail()
	goto block447
block443:
	frame.Fail()
	goto block447
block444:
	frame.Fail()
	goto block447
block445:
	frame.Fail()
	goto block447
block446:
	frame.Fail()
	goto block447
block447:
	frame.Recover(r0)
	goto block448
block448:
	r224 = Ident(frame)
	if frame.Flow == 0 {
		goto block449
	} else {
		goto block465
	}
block449:
	r225 = frame.Peek()
	if frame.Flow == 0 {
		goto block450
	} else {
		goto block465
	}
block450:
	r226 = '('
	goto block451
block451:
	r227 = r225 == r226
	goto block452
block452:
	if r227 {
		goto block453
	} else {
		goto block464
	}
block453:
	frame.Consume()
	goto block454
block454:
	S(frame)
	if frame.Flow == 0 {
		goto block455
	} else {
		goto block465
	}
block455:
	r228 = frame.Peek()
	if frame.Flow == 0 {
		goto block456
	} else {
		goto block465
	}
block456:
	r229 = ')'
	goto block457
block457:
	r230 = r228 == r229
	goto block458
block458:
	if r230 {
		goto block459
	} else {
		goto block463
	}
block459:
	frame.Consume()
	goto block460
block460:
	S(frame)
	if frame.Flow == 0 {
		goto block461
	} else {
		goto block465
	}
block461:
	r231 = &Call{Name: r224}
	goto block462
block462:
	ret0 = r231
	goto block549
block463:
	frame.Fail()
	goto block465
block464:
	frame.Fail()
	goto block465
block465:
	frame.Recover(r0)
	goto block466
block466:
	r232 = BinaryOperator(frame)
	if frame.Flow == 0 {
		goto block467
	} else {
		goto block471
	}
block467:
	r233 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block468
	} else {
		goto block471
	}
block468:
	r234 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block469
	} else {
		goto block471
	}
block469:
	r235 = &BinaryOp{Left: r233, Op: r232, Right: r234}
	goto block470
block470:
	ret0 = r235
	goto block549
block471:
	frame.Recover(r0)
	goto block472
block472:
	r236 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block473
	} else {
		goto block490
	}
block473:
	r237 = frame.Peek()
	if frame.Flow == 0 {
		goto block474
	} else {
		goto block490
	}
block474:
	r238 = '{'
	goto block475
block475:
	r239 = r237 == r238
	goto block476
block476:
	if r239 {
		goto block477
	} else {
		goto block489
	}
block477:
	frame.Consume()
	goto block478
block478:
	S(frame)
	if frame.Flow == 0 {
		goto block479
	} else {
		goto block490
	}
block479:
	r240 = ParseNamedExprList(frame)
	if frame.Flow == 0 {
		goto block480
	} else {
		goto block490
	}
block480:
	r241 = frame.Peek()
	if frame.Flow == 0 {
		goto block481
	} else {
		goto block490
	}
block481:
	r242 = '}'
	goto block482
block482:
	r243 = r241 == r242
	goto block483
block483:
	if r243 {
		goto block484
	} else {
		goto block488
	}
block484:
	frame.Consume()
	goto block485
block485:
	S(frame)
	if frame.Flow == 0 {
		goto block486
	} else {
		goto block490
	}
block486:
	r244 = &Construct{Type: r236, Args: r240}
	goto block487
block487:
	ret0 = r244
	goto block549
block488:
	frame.Fail()
	goto block490
block489:
	frame.Fail()
	goto block490
block490:
	frame.Recover(r0)
	goto block491
block491:
	r245 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block492
	} else {
		goto block509
	}
block492:
	r246 = frame.Peek()
	if frame.Flow == 0 {
		goto block493
	} else {
		goto block509
	}
block493:
	r247 = '{'
	goto block494
block494:
	r248 = r246 == r247
	goto block495
block495:
	if r248 {
		goto block496
	} else {
		goto block508
	}
block496:
	frame.Consume()
	goto block497
block497:
	S(frame)
	if frame.Flow == 0 {
		goto block498
	} else {
		goto block509
	}
block498:
	r249 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block499
	} else {
		goto block509
	}
block499:
	r250 = frame.Peek()
	if frame.Flow == 0 {
		goto block500
	} else {
		goto block509
	}
block500:
	r251 = '}'
	goto block501
block501:
	r252 = r250 == r251
	goto block502
block502:
	if r252 {
		goto block503
	} else {
		goto block507
	}
block503:
	frame.Consume()
	goto block504
block504:
	S(frame)
	if frame.Flow == 0 {
		goto block505
	} else {
		goto block509
	}
block505:
	r253 = &ConstructList{Type: r245, Args: r249}
	goto block506
block506:
	ret0 = r253
	goto block549
block507:
	frame.Fail()
	goto block509
block508:
	frame.Fail()
	goto block509
block509:
	frame.Recover(r0)
	goto block510
block510:
	r254 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		goto block511
	} else {
		goto block512
	}
block511:
	ret0 = r254
	goto block549
block512:
	frame.Recover(r0)
	goto block513
block513:
	r255 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		goto block514
	} else {
		goto block515
	}
block514:
	ret0 = r255
	goto block549
block515:
	frame.Recover(r0)
	goto block516
block516:
	r256 = Ident(frame)
	if frame.Flow == 0 {
		goto block517
	} else {
		goto block550
	}
block517:
	r257 = frame.Checkpoint()
	goto block518
block518:
	r258 = false
	goto block519
block519:
	r259 = frame.Checkpoint()
	goto block520
block520:
	r260 = frame.Peek()
	if frame.Flow == 0 {
		goto block521
	} else {
		goto block534
	}
block521:
	r261 = ':'
	goto block522
block522:
	r262 = r260 == r261
	goto block523
block523:
	if r262 {
		goto block524
	} else {
		goto block533
	}
block524:
	frame.Consume()
	goto block525
block525:
	r263 = frame.Peek()
	if frame.Flow == 0 {
		goto block526
	} else {
		goto block534
	}
block526:
	r264 = '='
	goto block527
block527:
	r265 = r263 == r264
	goto block528
block528:
	if r265 {
		goto block529
	} else {
		goto block532
	}
block529:
	frame.Consume()
	goto block530
block530:
	r266 = true
	goto block531
block531:
	r267 = r266
	goto block541
block532:
	frame.Fail()
	goto block534
block533:
	frame.Fail()
	goto block534
block534:
	frame.Recover(r259)
	goto block535
block535:
	r268 = frame.Peek()
	if frame.Flow == 0 {
		goto block536
	} else {
		goto block546
	}
block536:
	r269 = '='
	goto block537
block537:
	r270 = r268 == r269
	goto block538
block538:
	if r270 {
		goto block539
	} else {
		goto block545
	}
block539:
	frame.Consume()
	goto block540
block540:
	r267 = r258
	goto block541
block541:
	S(frame)
	if frame.Flow == 0 {
		goto block542
	} else {
		goto block546
	}
block542:
	r271 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block543
	} else {
		goto block546
	}
block543:
	r272 = &Assign{Expr: r271, Name: r256, Define: r267}
	goto block544
block544:
	ret0 = r272
	goto block549
block545:
	frame.Fail()
	goto block546
block546:
	frame.Recover(r257)
	goto block547
block547:
	r273 = &GetName{Name: r256}
	goto block548
block548:
	ret0 = r273
	goto block549
block549:
	return
block550:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block33
	}
block2:
	r1 = '{'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block32
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block33
	}
block7:
	r3 = []ASTExpr{}
	goto block8
block8:
	r4 = r3
	goto block9
block9:
	r5 = frame.Checkpoint()
	goto block10
block10:
	r6 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block22
	}
block11:
	r7 = append(r4, r6)
	goto block12
block12:
	r8 = frame.Checkpoint()
	goto block13
block13:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block20
	}
block14:
	r10 = ';'
	goto block15
block15:
	r11 = r9 == r10
	goto block16
block16:
	if r11 {
		goto block17
	} else {
		goto block19
	}
block17:
	frame.Consume()
	goto block18
block18:
	S(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block20
	}
block19:
	frame.Fail()
	goto block20
block20:
	frame.Recover(r8)
	goto block21
block21:
	r4 = r7
	goto block9
block22:
	frame.Recover(r5)
	goto block23
block23:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block33
	}
block24:
	r13 = '}'
	goto block25
block25:
	r14 = r12 == r13
	goto block26
block26:
	if r14 {
		goto block27
	} else {
		goto block31
	}
block27:
	frame.Consume()
	goto block28
block28:
	S(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block33
	}
block29:
	ret0 = r4
	goto block30
block30:
	return
block31:
	frame.Fail()
	goto block33
block32:
	frame.Fail()
	goto block33
block33:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block132
	}
block2:
	r1 = 's'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block131
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block132
	}
block7:
	r4 = 't'
	goto block8
block8:
	r5 = r3 == r4
	goto block9
block9:
	if r5 {
		goto block10
	} else {
		goto block130
	}
block10:
	frame.Consume()
	goto block11
block11:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block132
	}
block12:
	r7 = 'r'
	goto block13
block13:
	r8 = r6 == r7
	goto block14
block14:
	if r8 {
		goto block15
	} else {
		goto block129
	}
block15:
	frame.Consume()
	goto block16
block16:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block132
	}
block17:
	r10 = 'u'
	goto block18
block18:
	r11 = r9 == r10
	goto block19
block19:
	if r11 {
		goto block20
	} else {
		goto block128
	}
block20:
	frame.Consume()
	goto block21
block21:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block132
	}
block22:
	r13 = 'c'
	goto block23
block23:
	r14 = r12 == r13
	goto block24
block24:
	if r14 {
		goto block25
	} else {
		goto block127
	}
block25:
	frame.Consume()
	goto block26
block26:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block132
	}
block27:
	r16 = 't'
	goto block28
block28:
	r17 = r15 == r16
	goto block29
block29:
	if r17 {
		goto block30
	} else {
		goto block126
	}
block30:
	frame.Consume()
	goto block31
block31:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block32
	} else {
		goto block132
	}
block32:
	r18 = Ident(frame)
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block132
	}
block33:
	r19 = nil
	goto block34
block34:
	r20 = frame.Checkpoint()
	goto block35
block35:
	r21 = frame.Peek()
	if frame.Flow == 0 {
		goto block36
	} else {
		goto block98
	}
block36:
	r22 = 'i'
	goto block37
block37:
	r23 = r21 == r22
	goto block38
block38:
	if r23 {
		goto block39
	} else {
		goto block97
	}
block39:
	frame.Consume()
	goto block40
block40:
	r24 = frame.Peek()
	if frame.Flow == 0 {
		goto block41
	} else {
		goto block98
	}
block41:
	r25 = 'm'
	goto block42
block42:
	r26 = r24 == r25
	goto block43
block43:
	if r26 {
		goto block44
	} else {
		goto block96
	}
block44:
	frame.Consume()
	goto block45
block45:
	r27 = frame.Peek()
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block98
	}
block46:
	r28 = 'p'
	goto block47
block47:
	r29 = r27 == r28
	goto block48
block48:
	if r29 {
		goto block49
	} else {
		goto block95
	}
block49:
	frame.Consume()
	goto block50
block50:
	r30 = frame.Peek()
	if frame.Flow == 0 {
		goto block51
	} else {
		goto block98
	}
block51:
	r31 = 'l'
	goto block52
block52:
	r32 = r30 == r31
	goto block53
block53:
	if r32 {
		goto block54
	} else {
		goto block94
	}
block54:
	frame.Consume()
	goto block55
block55:
	r33 = frame.Peek()
	if frame.Flow == 0 {
		goto block56
	} else {
		goto block98
	}
block56:
	r34 = 'e'
	goto block57
block57:
	r35 = r33 == r34
	goto block58
block58:
	if r35 {
		goto block59
	} else {
		goto block93
	}
block59:
	frame.Consume()
	goto block60
block60:
	r36 = frame.Peek()
	if frame.Flow == 0 {
		goto block61
	} else {
		goto block98
	}
block61:
	r37 = 'm'
	goto block62
block62:
	r38 = r36 == r37
	goto block63
block63:
	if r38 {
		goto block64
	} else {
		goto block92
	}
block64:
	frame.Consume()
	goto block65
block65:
	r39 = frame.Peek()
	if frame.Flow == 0 {
		goto block66
	} else {
		goto block98
	}
block66:
	r40 = 'e'
	goto block67
block67:
	r41 = r39 == r40
	goto block68
block68:
	if r41 {
		goto block69
	} else {
		goto block91
	}
block69:
	frame.Consume()
	goto block70
block70:
	r42 = frame.Peek()
	if frame.Flow == 0 {
		goto block71
	} else {
		goto block98
	}
block71:
	r43 = 'n'
	goto block72
block72:
	r44 = r42 == r43
	goto block73
block73:
	if r44 {
		goto block74
	} else {
		goto block90
	}
block74:
	frame.Consume()
	goto block75
block75:
	r45 = frame.Peek()
	if frame.Flow == 0 {
		goto block76
	} else {
		goto block98
	}
block76:
	r46 = 't'
	goto block77
block77:
	r47 = r45 == r46
	goto block78
block78:
	if r47 {
		goto block79
	} else {
		goto block89
	}
block79:
	frame.Consume()
	goto block80
block80:
	r48 = frame.Peek()
	if frame.Flow == 0 {
		goto block81
	} else {
		goto block98
	}
block81:
	r49 = 's'
	goto block82
block82:
	r50 = r48 == r49
	goto block83
block83:
	if r50 {
		goto block84
	} else {
		goto block88
	}
block84:
	frame.Consume()
	goto block85
block85:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block86
	} else {
		goto block98
	}
block86:
	r51 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block87
	} else {
		goto block98
	}
block87:
	r52 = r51
	goto block100
block88:
	frame.Fail()
	goto block98
block89:
	frame.Fail()
	goto block98
block90:
	frame.Fail()
	goto block98
block91:
	frame.Fail()
	goto block98
block92:
	frame.Fail()
	goto block98
block93:
	frame.Fail()
	goto block98
block94:
	frame.Fail()
	goto block98
block95:
	frame.Fail()
	goto block98
block96:
	frame.Fail()
	goto block98
block97:
	frame.Fail()
	goto block98
block98:
	frame.Recover(r20)
	goto block99
block99:
	r52 = r19
	goto block100
block100:
	r53 = frame.Peek()
	if frame.Flow == 0 {
		goto block101
	} else {
		goto block132
	}
block101:
	r54 = '{'
	goto block102
block102:
	r55 = r53 == r54
	goto block103
block103:
	if r55 {
		goto block104
	} else {
		goto block125
	}
block104:
	frame.Consume()
	goto block105
block105:
	S(frame)
	if frame.Flow == 0 {
		goto block106
	} else {
		goto block132
	}
block106:
	r56 = []*FieldDecl{}
	goto block107
block107:
	r57 = r56
	goto block108
block108:
	r58 = frame.Checkpoint()
	goto block109
block109:
	r59 = Ident(frame)
	if frame.Flow == 0 {
		goto block110
	} else {
		goto block114
	}
block110:
	r60 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block111
	} else {
		goto block114
	}
block111:
	r61 = &FieldDecl{Name: r59, Type: r60}
	goto block112
block112:
	r62 = append(r57, r61)
	goto block113
block113:
	r57 = r62
	goto block108
block114:
	frame.Recover(r58)
	goto block115
block115:
	r63 = frame.Peek()
	if frame.Flow == 0 {
		goto block116
	} else {
		goto block132
	}
block116:
	r64 = '}'
	goto block117
block117:
	r65 = r63 == r64
	goto block118
block118:
	if r65 {
		goto block119
	} else {
		goto block124
	}
block119:
	frame.Consume()
	goto block120
block120:
	S(frame)
	if frame.Flow == 0 {
		goto block121
	} else {
		goto block132
	}
block121:
	r66 = &StructDecl{Name: r18, Implements: r52, Fields: r57}
	goto block122
block122:
	ret0 = r66
	goto block123
block123:
	return
block124:
	frame.Fail()
	goto block132
block125:
	frame.Fail()
	goto block132
block126:
	frame.Fail()
	goto block132
block127:
	frame.Fail()
	goto block132
block128:
	frame.Fail()
	goto block132
block129:
	frame.Fail()
	goto block132
block130:
	frame.Fail()
	goto block132
block131:
	frame.Fail()
	goto block132
block132:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block32
	}
block2:
	r1 = 'f'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block31
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block32
	}
block7:
	r4 = 'u'
	goto block8
block8:
	r5 = r3 == r4
	goto block9
block9:
	if r5 {
		goto block10
	} else {
		goto block30
	}
block10:
	frame.Consume()
	goto block11
block11:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block32
	}
block12:
	r7 = 'n'
	goto block13
block13:
	r8 = r6 == r7
	goto block14
block14:
	if r8 {
		goto block15
	} else {
		goto block29
	}
block15:
	frame.Consume()
	goto block16
block16:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block32
	}
block17:
	r10 = 'c'
	goto block18
block18:
	r11 = r9 == r10
	goto block19
block19:
	if r11 {
		goto block20
	} else {
		goto block28
	}
block20:
	frame.Consume()
	goto block21
block21:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block32
	}
block22:
	r12 = Ident(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block32
	}
block23:
	r13 = ParseTypeList(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block32
	}
block24:
	r14 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block32
	}
block25:
	r15 = &FuncDecl{Name: r12, ReturnTypes: r13, Block: r14}
	goto block26
block26:
	ret0 = r15
	goto block27
block27:
	return
block28:
	frame.Fail()
	goto block32
block29:
	frame.Fail()
	goto block32
block30:
	frame.Fail()
	goto block32
block31:
	frame.Fail()
	goto block32
block32:
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
	goto block0
block0:
	goto block1
block1:
	r0 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block34
	}
block2:
	r1 = 't'
	goto block3
block3:
	r2 = r0 == r1
	goto block4
block4:
	if r2 {
		goto block5
	} else {
		goto block33
	}
block5:
	frame.Consume()
	goto block6
block6:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block34
	}
block7:
	r4 = 'e'
	goto block8
block8:
	r5 = r3 == r4
	goto block9
block9:
	if r5 {
		goto block10
	} else {
		goto block32
	}
block10:
	frame.Consume()
	goto block11
block11:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block34
	}
block12:
	r7 = 's'
	goto block13
block13:
	r8 = r6 == r7
	goto block14
block14:
	if r8 {
		goto block15
	} else {
		goto block31
	}
block15:
	frame.Consume()
	goto block16
block16:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block34
	}
block17:
	r10 = 't'
	goto block18
block18:
	r11 = r9 == r10
	goto block19
block19:
	if r11 {
		goto block20
	} else {
		goto block30
	}
block20:
	frame.Consume()
	goto block21
block21:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block34
	}
block22:
	r12 = Ident(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block34
	}
block23:
	r13 = Ident(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block34
	}
block24:
	r14 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block34
	}
block25:
	S(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block34
	}
block26:
	r15 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block34
	}
block27:
	r16 = &Test{Rule: r12, Name: r13, Input: r14, Destructure: r15}
	goto block28
block28:
	ret0 = r16
	goto block29
block29:
	return
block30:
	frame.Fail()
	goto block34
block31:
	frame.Fail()
	goto block34
block32:
	frame.Fail()
	goto block34
block33:
	frame.Fail()
	goto block34
block34:
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
	goto block0
block0:
	goto block1
block1:
	r0 = []ASTDecl{}
	goto block2
block2:
	r1 = []*Test{}
	goto block3
block3:
	r2, r3 = r0, r1
	goto block4
block4:
	r4 = frame.Checkpoint()
	goto block5
block5:
	r5 = frame.Checkpoint()
	goto block6
block6:
	r6 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block9
	}
block7:
	r7 = append(r2, r6)
	goto block8
block8:
	r2, r3 = r7, r3
	goto block4
block9:
	frame.Recover(r5)
	goto block10
block10:
	r8 = ParseStructDecl(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block13
	}
block11:
	r9 = append(r2, r8)
	goto block12
block12:
	r2, r3 = r9, r3
	goto block4
block13:
	frame.Recover(r5)
	goto block14
block14:
	r10 = ParseTest(frame)
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block17
	}
block15:
	r11 = append(r3, r10)
	goto block16
block16:
	r2, r3 = r2, r11
	goto block4
block17:
	frame.Recover(r4)
	goto block18
block18:
	r12 = frame.LookaheadBegin()
	goto block19
block19:
	frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block23
	}
block20:
	frame.Consume()
	goto block21
block21:
	frame.LookaheadFail(r12)
	goto block22
block22:
	return
block23:
	frame.LookaheadNormal(r12)
	goto block24
block24:
	r13 = &File{Decls: r2, Tests: r3}
	goto block25
block25:
	ret0 = r13
	goto block26
block26:
	return
}
