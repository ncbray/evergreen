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
	var r0 string
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
	var r30 string
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Checkpoint()
	goto block2
block2:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block52
	}
block3:
	r3 = 'a'
	goto block4
block4:
	r4 = r2 >= r3
	goto block5
block5:
	if r4 {
		goto block6
	} else {
		goto block9
	}
block6:
	r5 = 'z'
	goto block7
block7:
	r6 = r2 <= r5
	goto block8
block8:
	if r6 {
		goto block18
	} else {
		goto block9
	}
block9:
	r7 = 'A'
	goto block10
block10:
	r8 = r2 >= r7
	goto block11
block11:
	if r8 {
		goto block12
	} else {
		goto block15
	}
block12:
	r9 = 'Z'
	goto block13
block13:
	r10 = r2 <= r9
	goto block14
block14:
	if r10 {
		goto block18
	} else {
		goto block15
	}
block15:
	r11 = '_'
	goto block16
block16:
	r12 = r2 == r11
	goto block17
block17:
	if r12 {
		goto block18
	} else {
		goto block51
	}
block18:
	frame.Consume()
	goto block19
block19:
	r13 = frame.Checkpoint()
	goto block20
block20:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block44
	}
block21:
	r15 = 'a'
	goto block22
block22:
	r16 = r14 >= r15
	goto block23
block23:
	if r16 {
		goto block24
	} else {
		goto block27
	}
block24:
	r17 = 'z'
	goto block25
block25:
	r18 = r14 <= r17
	goto block26
block26:
	if r18 {
		goto block42
	} else {
		goto block27
	}
block27:
	r19 = 'A'
	goto block28
block28:
	r20 = r14 >= r19
	goto block29
block29:
	if r20 {
		goto block30
	} else {
		goto block33
	}
block30:
	r21 = 'Z'
	goto block31
block31:
	r22 = r14 <= r21
	goto block32
block32:
	if r22 {
		goto block42
	} else {
		goto block33
	}
block33:
	r23 = '_'
	goto block34
block34:
	r24 = r14 == r23
	goto block35
block35:
	if r24 {
		goto block42
	} else {
		goto block36
	}
block36:
	r25 = '0'
	goto block37
block37:
	r26 = r14 >= r25
	goto block38
block38:
	if r26 {
		goto block39
	} else {
		goto block43
	}
block39:
	r27 = '9'
	goto block40
block40:
	r28 = r14 <= r27
	goto block41
block41:
	if r28 {
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
	frame.Recover(r13)
	goto block45
block45:
	r29 = frame.Slice(r1)
	goto block46
block46:
	r0 = r29
	goto block47
block47:
	S(frame)
	if frame.Flow == 0 {
		goto block48
	} else {
		goto block52
	}
block48:
	r30 = r0
	goto block49
block49:
	ret0 = r30
	goto block50
block50:
	return
block51:
	frame.Fail()
	goto block52
block52:
	return
}
func DecodeInt(frame *runtime.State) (ret0 int) {
	var r0 int
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 bool
	var r7 int
	var r8 rune
	var r9 int
	var r10 int
	var r11 int
	var r12 int
	var r13 int
	var r14 int
	var r15 int
	var r16 int
	var r17 rune
	var r18 rune
	var r19 bool
	var r20 rune
	var r21 bool
	var r22 int
	var r23 rune
	var r24 int
	var r25 int
	var r26 int
	var r27 int
	var r28 int
	var r29 int
	var r30 int
	var r31 int
	goto block0
block0:
	goto block1
block1:
	r0 = 0
	goto block2
block2:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block47
	}
block3:
	r3 = '0'
	goto block4
block4:
	r4 = r2 >= r3
	goto block5
block5:
	if r4 {
		goto block6
	} else {
		goto block46
	}
block6:
	r5 = '9'
	goto block7
block7:
	r6 = r2 <= r5
	goto block8
block8:
	if r6 {
		goto block9
	} else {
		goto block46
	}
block9:
	frame.Consume()
	goto block10
block10:
	r7 = int(r2)
	goto block11
block11:
	r8 = '0'
	goto block12
block12:
	r9 = int(r8)
	goto block13
block13:
	r10 = r7 - r9
	goto block14
block14:
	r1 = r10
	goto block15
block15:
	r11 = r0
	goto block16
block16:
	r12 = 10
	goto block17
block17:
	r13 = r11 * r12
	goto block18
block18:
	r14 = r1
	goto block19
block19:
	r15 = r13 + r14
	goto block20
block20:
	r0 = r15
	goto block21
block21:
	r16 = frame.Checkpoint()
	goto block22
block22:
	r17 = frame.Peek()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block42
	}
block23:
	r18 = '0'
	goto block24
block24:
	r19 = r17 >= r18
	goto block25
block25:
	if r19 {
		goto block26
	} else {
		goto block41
	}
block26:
	r20 = '9'
	goto block27
block27:
	r21 = r17 <= r20
	goto block28
block28:
	if r21 {
		goto block29
	} else {
		goto block41
	}
block29:
	frame.Consume()
	goto block30
block30:
	r22 = int(r17)
	goto block31
block31:
	r23 = '0'
	goto block32
block32:
	r24 = int(r23)
	goto block33
block33:
	r25 = r22 - r24
	goto block34
block34:
	r1 = r25
	goto block35
block35:
	r26 = r0
	goto block36
block36:
	r27 = 10
	goto block37
block37:
	r28 = r26 * r27
	goto block38
block38:
	r29 = r1
	goto block39
block39:
	r30 = r28 + r29
	goto block40
block40:
	r0 = r30
	goto block21
block41:
	frame.Fail()
	goto block42
block42:
	frame.Recover(r16)
	goto block43
block43:
	r31 = r0
	goto block44
block44:
	ret0 = r31
	goto block45
block45:
	return
block46:
	frame.Fail()
	goto block47
block47:
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
	var r0 []rune
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 []rune
	var r5 int
	var r6 int
	var r7 []rune
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 bool
	var r13 []rune
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 []rune
	var r18 rune
	var r19 []rune
	var r20 rune
	var r21 rune
	var r22 bool
	var r23 []rune
	var r24 string
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block45
	}
block2:
	r2 = '"'
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block44
	}
block5:
	frame.Consume()
	goto block6
block6:
	r4 = []rune{}
	goto block7
block7:
	r0 = r4
	goto block8
block8:
	r5 = frame.Checkpoint()
	goto block9
block9:
	r6 = frame.Checkpoint()
	goto block10
block10:
	r7 = r0
	goto block11
block11:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block22
	}
block12:
	r9 = '"'
	goto block13
block13:
	r10 = r8 == r9
	goto block14
block14:
	if r10 {
		goto block18
	} else {
		goto block15
	}
block15:
	r11 = '\\'
	goto block16
block16:
	r12 = r8 == r11
	goto block17
block17:
	if r12 {
		goto block18
	} else {
		goto block19
	}
block18:
	frame.Fail()
	goto block22
block19:
	frame.Consume()
	goto block20
block20:
	r13 = append(r7, r8)
	goto block21
block21:
	r0 = r13
	goto block8
block22:
	frame.Recover(r6)
	goto block23
block23:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block33
	}
block24:
	r15 = '\\'
	goto block25
block25:
	r16 = r14 == r15
	goto block26
block26:
	if r16 {
		goto block27
	} else {
		goto block32
	}
block27:
	frame.Consume()
	goto block28
block28:
	r17 = r0
	goto block29
block29:
	r18 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block30
	} else {
		goto block33
	}
block30:
	r19 = append(r17, r18)
	goto block31
block31:
	r0 = r19
	goto block8
block32:
	frame.Fail()
	goto block33
block33:
	frame.Recover(r5)
	goto block34
block34:
	r20 = frame.Peek()
	if frame.Flow == 0 {
		goto block35
	} else {
		goto block45
	}
block35:
	r21 = '"'
	goto block36
block36:
	r22 = r20 == r21
	goto block37
block37:
	if r22 {
		goto block38
	} else {
		goto block43
	}
block38:
	frame.Consume()
	goto block39
block39:
	r23 = r0
	goto block40
block40:
	r24 = string(r23)
	goto block41
block41:
	ret0 = r24
	goto block42
block42:
	return
block43:
	frame.Fail()
	goto block45
block44:
	frame.Fail()
	goto block45
block45:
	return
}
func DecodeRune(frame *runtime.State) (ret0 rune) {
	var r0 rune
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 int
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 rune
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 rune
	goto block0
block0:
	goto block1
block1:
	r0 = '\x00'
	goto block2
block2:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block37
	}
block3:
	r2 = '\''
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block6
	} else {
		goto block36
	}
block6:
	frame.Consume()
	goto block7
block7:
	r4 = frame.Checkpoint()
	goto block8
block8:
	r5 = frame.Peek()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block18
	}
block9:
	r6 = '\\'
	goto block10
block10:
	r7 = r5 == r6
	goto block11
block11:
	if r7 {
		goto block15
	} else {
		goto block12
	}
block12:
	r8 = '\''
	goto block13
block13:
	r9 = r5 == r8
	goto block14
block14:
	if r9 {
		goto block15
	} else {
		goto block16
	}
block15:
	frame.Fail()
	goto block18
block16:
	frame.Consume()
	goto block17
block17:
	r0 = r5
	goto block26
block18:
	frame.Recover(r4)
	goto block19
block19:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block37
	}
block20:
	r11 = '\\'
	goto block21
block21:
	r12 = r10 == r11
	goto block22
block22:
	if r12 {
		goto block23
	} else {
		goto block35
	}
block23:
	frame.Consume()
	goto block24
block24:
	r13 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block37
	}
block25:
	r0 = r13
	goto block26
block26:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block37
	}
block27:
	r15 = '\''
	goto block28
block28:
	r16 = r14 == r15
	goto block29
block29:
	if r16 {
		goto block30
	} else {
		goto block34
	}
block30:
	frame.Consume()
	goto block31
block31:
	r17 = r0
	goto block32
block32:
	ret0 = r17
	goto block33
block33:
	return
block34:
	frame.Fail()
	goto block37
block35:
	frame.Fail()
	goto block37
block36:
	frame.Fail()
	goto block37
block37:
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
	var r0 rune
	var r1 string
	var r2 string
	var r3 string
	var r4 int
	var r5 string
	var r6 bool
	var r7 string
	var r8 int
	var r9 int
	var r10 rune
	var r11 string
	var r12 string
	var r13 rune
	var r14 *RuneLiteral
	var r15 int
	var r16 string
	var r17 string
	var r18 string
	var r19 string
	var r20 *StringLiteral
	var r21 int
	var r22 int
	var r23 string
	var r24 string
	var r25 int
	var r26 *IntLiteral
	var r27 int
	var r28 bool
	var r29 string
	var r30 string
	var r31 bool
	var r32 *BoolLiteral
	goto block0
block0:
	goto block1
block1:
	r8 = frame.Checkpoint()
	goto block2
block2:
	r0 = '\x00'
	goto block3
block3:
	r9 = frame.Checkpoint()
	goto block4
block4:
	r10 = DecodeRune(frame)
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block13
	}
block5:
	r0 = r10
	goto block6
block6:
	r11 = frame.Slice(r9)
	goto block7
block7:
	r1 = r11
	goto block8
block8:
	S(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block13
	}
block9:
	r12 = r1
	goto block10
block10:
	r13 = r0
	goto block11
block11:
	r14 = &RuneLiteral{Text: r12, Value: r13}
	goto block12
block12:
	ret0 = r14
	goto block49
block13:
	frame.Recover(r8)
	goto block14
block14:
	r2 = ""
	goto block15
block15:
	r15 = frame.Checkpoint()
	goto block16
block16:
	r16 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block25
	}
block17:
	r2 = r16
	goto block18
block18:
	r17 = frame.Slice(r15)
	goto block19
block19:
	r3 = r17
	goto block20
block20:
	S(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block25
	}
block21:
	r18 = r3
	goto block22
block22:
	r19 = r2
	goto block23
block23:
	r20 = &StringLiteral{Text: r18, Value: r19}
	goto block24
block24:
	ret0 = r20
	goto block49
block25:
	frame.Recover(r8)
	goto block26
block26:
	r4 = 0
	goto block27
block27:
	r21 = frame.Checkpoint()
	goto block28
block28:
	r22 = DecodeInt(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block37
	}
block29:
	r4 = r22
	goto block30
block30:
	r23 = frame.Slice(r21)
	goto block31
block31:
	r5 = r23
	goto block32
block32:
	S(frame)
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block37
	}
block33:
	r24 = r5
	goto block34
block34:
	r25 = r4
	goto block35
block35:
	r26 = &IntLiteral{Text: r24, Value: r25}
	goto block36
block36:
	ret0 = r26
	goto block49
block37:
	frame.Recover(r8)
	goto block38
block38:
	r6 = false
	goto block39
block39:
	r27 = frame.Checkpoint()
	goto block40
block40:
	r28 = DecodeBool(frame)
	if frame.Flow == 0 {
		goto block41
	} else {
		goto block50
	}
block41:
	r6 = r28
	goto block42
block42:
	r29 = frame.Slice(r27)
	goto block43
block43:
	r7 = r29
	goto block44
block44:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block45
	} else {
		goto block50
	}
block45:
	r30 = r7
	goto block46
block46:
	r31 = r6
	goto block47
block47:
	r32 = &BoolLiteral{Text: r30, Value: r31}
	goto block48
block48:
	ret0 = r32
	goto block49
block49:
	return
block50:
	return
}
func BinaryOperator(frame *runtime.State) (ret0 string) {
	var r0 string
	var r1 int
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 rune
	var r16 bool
	var r17 int
	var r18 rune
	var r19 rune
	var r20 bool
	var r21 rune
	var r22 rune
	var r23 bool
	var r24 rune
	var r25 bool
	var r26 rune
	var r27 rune
	var r28 bool
	var r29 int
	var r30 rune
	var r31 rune
	var r32 bool
	var r33 rune
	var r34 bool
	var r35 rune
	var r36 bool
	var r37 rune
	var r38 bool
	var r39 rune
	var r40 bool
	var r41 rune
	var r42 bool
	var r43 rune
	var r44 bool
	var r45 rune
	var r46 bool
	var r47 string
	var r48 string
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Checkpoint()
	goto block2
block2:
	r2 = frame.Checkpoint()
	goto block3
block3:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block18
	}
block4:
	r4 = '+'
	goto block5
block5:
	r5 = r3 == r4
	goto block6
block6:
	if r5 {
		goto block16
	} else {
		goto block7
	}
block7:
	r6 = '-'
	goto block8
block8:
	r7 = r3 == r6
	goto block9
block9:
	if r7 {
		goto block16
	} else {
		goto block10
	}
block10:
	r8 = '*'
	goto block11
block11:
	r9 = r3 == r8
	goto block12
block12:
	if r9 {
		goto block16
	} else {
		goto block13
	}
block13:
	r10 = '/'
	goto block14
block14:
	r11 = r3 == r10
	goto block15
block15:
	if r11 {
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
	frame.Recover(r2)
	goto block19
block19:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block36
	}
block20:
	r13 = '<'
	goto block21
block21:
	r14 = r12 == r13
	goto block22
block22:
	if r14 {
		goto block26
	} else {
		goto block23
	}
block23:
	r15 = '>'
	goto block24
block24:
	r16 = r12 == r15
	goto block25
block25:
	if r16 {
		goto block26
	} else {
		goto block35
	}
block26:
	frame.Consume()
	goto block27
block27:
	r17 = frame.Checkpoint()
	goto block28
block28:
	r18 = frame.Peek()
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block34
	}
block29:
	r19 = '='
	goto block30
block30:
	r20 = r18 == r19
	goto block31
block31:
	if r20 {
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
	frame.Recover(r17)
	goto block80
block35:
	frame.Fail()
	goto block36
block36:
	frame.Recover(r2)
	goto block37
block37:
	r21 = frame.Peek()
	if frame.Flow == 0 {
		goto block38
	} else {
		goto block88
	}
block38:
	r22 = '!'
	goto block39
block39:
	r23 = r21 == r22
	goto block40
block40:
	if r23 {
		goto block44
	} else {
		goto block41
	}
block41:
	r24 = '='
	goto block42
block42:
	r25 = r21 == r24
	goto block43
block43:
	if r25 {
		goto block44
	} else {
		goto block87
	}
block44:
	frame.Consume()
	goto block45
block45:
	r26 = frame.Peek()
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block88
	}
block46:
	r27 = '='
	goto block47
block47:
	r28 = r26 == r27
	goto block48
block48:
	if r28 {
		goto block49
	} else {
		goto block86
	}
block49:
	frame.Consume()
	goto block50
block50:
	r29 = frame.LookaheadBegin()
	goto block51
block51:
	r30 = frame.Peek()
	if frame.Flow == 0 {
		goto block52
	} else {
		goto block79
	}
block52:
	r31 = '+'
	goto block53
block53:
	r32 = r30 == r31
	goto block54
block54:
	if r32 {
		goto block76
	} else {
		goto block55
	}
block55:
	r33 = '-'
	goto block56
block56:
	r34 = r30 == r33
	goto block57
block57:
	if r34 {
		goto block76
	} else {
		goto block58
	}
block58:
	r35 = '*'
	goto block59
block59:
	r36 = r30 == r35
	goto block60
block60:
	if r36 {
		goto block76
	} else {
		goto block61
	}
block61:
	r37 = '/'
	goto block62
block62:
	r38 = r30 == r37
	goto block63
block63:
	if r38 {
		goto block76
	} else {
		goto block64
	}
block64:
	r39 = '<'
	goto block65
block65:
	r40 = r30 == r39
	goto block66
block66:
	if r40 {
		goto block76
	} else {
		goto block67
	}
block67:
	r41 = '>'
	goto block68
block68:
	r42 = r30 == r41
	goto block69
block69:
	if r42 {
		goto block76
	} else {
		goto block70
	}
block70:
	r43 = '!'
	goto block71
block71:
	r44 = r30 == r43
	goto block72
block72:
	if r44 {
		goto block76
	} else {
		goto block73
	}
block73:
	r45 = '='
	goto block74
block74:
	r46 = r30 == r45
	goto block75
block75:
	if r46 {
		goto block76
	} else {
		goto block78
	}
block76:
	frame.Consume()
	goto block77
block77:
	frame.LookaheadFail(r29)
	goto block88
block78:
	frame.Fail()
	goto block79
block79:
	frame.LookaheadNormal(r29)
	goto block80
block80:
	r47 = frame.Slice(r1)
	goto block81
block81:
	r0 = r47
	goto block82
block82:
	S(frame)
	if frame.Flow == 0 {
		goto block83
	} else {
		goto block88
	}
block83:
	r48 = r0
	goto block84
block84:
	ret0 = r48
	goto block85
block85:
	return
block86:
	frame.Fail()
	goto block88
block87:
	frame.Fail()
	goto block88
block88:
	return
}
func StringMatchExpr(frame *runtime.State) (ret0 *StringMatch) {
	var r0 TextMatch
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 TextMatch
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 TextMatch
	var r9 *StringMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block21
	}
block2:
	r2 = '/'
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block20
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block21
	}
block7:
	r4 = ParseMatchChoice(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block21
	}
block8:
	r0 = r4
	goto block9
block9:
	r5 = frame.Peek()
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block21
	}
block10:
	r6 = '/'
	goto block11
block11:
	r7 = r5 == r6
	goto block12
block12:
	if r7 {
		goto block13
	} else {
		goto block19
	}
block13:
	frame.Consume()
	goto block14
block14:
	S(frame)
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block21
	}
block15:
	r8 = r0
	goto block16
block16:
	r9 = &StringMatch{Match: r8}
	goto block17
block17:
	ret0 = r9
	goto block18
block18:
	return
block19:
	frame.Fail()
	goto block21
block20:
	frame.Fail()
	goto block21
block21:
	return
}
func RuneMatchExpr(frame *runtime.State) (ret0 *RuneMatch) {
	var r0 *RuneRangeMatch
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 *RuneRangeMatch
	var r5 *RuneRangeMatch
	var r6 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block14
	}
block2:
	r2 = '$'
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block13
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block14
	}
block7:
	r4 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block14
	}
block8:
	r0 = r4
	goto block9
block9:
	r5 = r0
	goto block10
block10:
	r6 = &RuneMatch{Match: r5}
	goto block11
block11:
	ret0 = r6
	goto block12
block12:
	return
block13:
	frame.Fail()
	goto block14
block14:
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
	var r0 *TypeRef
	var r1 []*DestructureField
	var r2 string
	var r3 Destructure
	var r4 *ListTypeRef
	var r5 []Destructure
	var r6 int
	var r7 *TypeRef
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 []*DestructureField
	var r12 int
	var r13 string
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 Destructure
	var r18 []*DestructureField
	var r19 string
	var r20 Destructure
	var r21 *DestructureField
	var r22 []*DestructureField
	var r23 rune
	var r24 rune
	var r25 bool
	var r26 *TypeRef
	var r27 []*DestructureField
	var r28 *DestructureStruct
	var r29 *ListTypeRef
	var r30 rune
	var r31 rune
	var r32 bool
	var r33 []Destructure
	var r34 int
	var r35 []Destructure
	var r36 Destructure
	var r37 []Destructure
	var r38 rune
	var r39 rune
	var r40 bool
	var r41 *ListTypeRef
	var r42 []Destructure
	var r43 *DestructureList
	var r44 ASTExpr
	var r45 *DestructureValue
	goto block0
block0:
	goto block1
block1:
	r6 = frame.Checkpoint()
	goto block2
block2:
	r7 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block43
	}
block3:
	r0 = r7
	goto block4
block4:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block43
	}
block5:
	r9 = '{'
	goto block6
block6:
	r10 = r8 == r9
	goto block7
block7:
	if r10 {
		goto block8
	} else {
		goto block42
	}
block8:
	frame.Consume()
	goto block9
block9:
	S(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block43
	}
block10:
	r11 = []*DestructureField{}
	goto block11
block11:
	r1 = r11
	goto block12
block12:
	r12 = frame.Checkpoint()
	goto block13
block13:
	r13 = Ident(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block30
	}
block14:
	r2 = r13
	goto block15
block15:
	r14 = frame.Peek()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block30
	}
block16:
	r15 = ':'
	goto block17
block17:
	r16 = r14 == r15
	goto block18
block18:
	if r16 {
		goto block19
	} else {
		goto block29
	}
block19:
	frame.Consume()
	goto block20
block20:
	S(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block30
	}
block21:
	r17 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block30
	}
block22:
	r3 = r17
	goto block23
block23:
	r18 = r1
	goto block24
block24:
	r19 = r2
	goto block25
block25:
	r20 = r3
	goto block26
block26:
	r21 = &DestructureField{Name: r19, Destructure: r20}
	goto block27
block27:
	r22 = append(r18, r21)
	goto block28
block28:
	r1 = r22
	goto block12
block29:
	frame.Fail()
	goto block30
block30:
	frame.Recover(r12)
	goto block31
block31:
	r23 = frame.Peek()
	if frame.Flow == 0 {
		goto block32
	} else {
		goto block43
	}
block32:
	r24 = '}'
	goto block33
block33:
	r25 = r23 == r24
	goto block34
block34:
	if r25 {
		goto block35
	} else {
		goto block41
	}
block35:
	frame.Consume()
	goto block36
block36:
	S(frame)
	if frame.Flow == 0 {
		goto block37
	} else {
		goto block43
	}
block37:
	r26 = r0
	goto block38
block38:
	r27 = r1
	goto block39
block39:
	r28 = &DestructureStruct{Type: r26, Args: r27}
	goto block40
block40:
	ret0 = r28
	goto block76
block41:
	frame.Fail()
	goto block43
block42:
	frame.Fail()
	goto block43
block43:
	frame.Recover(r6)
	goto block44
block44:
	r29 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block45
	} else {
		goto block72
	}
block45:
	r4 = r29
	goto block46
block46:
	r30 = frame.Peek()
	if frame.Flow == 0 {
		goto block47
	} else {
		goto block72
	}
block47:
	r31 = '{'
	goto block48
block48:
	r32 = r30 == r31
	goto block49
block49:
	if r32 {
		goto block50
	} else {
		goto block71
	}
block50:
	frame.Consume()
	goto block51
block51:
	S(frame)
	if frame.Flow == 0 {
		goto block52
	} else {
		goto block72
	}
block52:
	r33 = []Destructure{}
	goto block53
block53:
	r5 = r33
	goto block54
block54:
	r34 = frame.Checkpoint()
	goto block55
block55:
	r35 = r5
	goto block56
block56:
	r36 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block57
	} else {
		goto block59
	}
block57:
	r37 = append(r35, r36)
	goto block58
block58:
	r5 = r37
	goto block54
block59:
	frame.Recover(r34)
	goto block60
block60:
	r38 = frame.Peek()
	if frame.Flow == 0 {
		goto block61
	} else {
		goto block72
	}
block61:
	r39 = '}'
	goto block62
block62:
	r40 = r38 == r39
	goto block63
block63:
	if r40 {
		goto block64
	} else {
		goto block70
	}
block64:
	frame.Consume()
	goto block65
block65:
	S(frame)
	if frame.Flow == 0 {
		goto block66
	} else {
		goto block72
	}
block66:
	r41 = r4
	goto block67
block67:
	r42 = r5
	goto block68
block68:
	r43 = &DestructureList{Type: r41, Args: r42}
	goto block69
block69:
	ret0 = r43
	goto block76
block70:
	frame.Fail()
	goto block72
block71:
	frame.Fail()
	goto block72
block72:
	frame.Recover(r6)
	goto block73
block73:
	r44 = Literal(frame)
	if frame.Flow == 0 {
		goto block74
	} else {
		goto block77
	}
block74:
	r45 = &DestructureValue{Expr: r44}
	goto block75
block75:
	ret0 = r45
	goto block76
block76:
	return
block77:
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
	var r1 rune
	var r2 rune
	var r3 rune
	var r4 int
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 rune
	var r11 *RuneFilter
	goto block0
block0:
	goto block1
block1:
	r2 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block20
	}
block2:
	r0 = r2
	goto block3
block3:
	r3 = r0
	goto block4
block4:
	r1 = r3
	goto block5
block5:
	r4 = frame.Checkpoint()
	goto block6
block6:
	r5 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block14
	}
block7:
	r6 = '-'
	goto block8
block8:
	r7 = r5 == r6
	goto block9
block9:
	if r7 {
		goto block10
	} else {
		goto block13
	}
block10:
	frame.Consume()
	goto block11
block11:
	r8 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block14
	}
block12:
	r1 = r8
	goto block15
block13:
	frame.Fail()
	goto block14
block14:
	frame.Recover(r4)
	goto block15
block15:
	r9 = r0
	goto block16
block16:
	r10 = r1
	goto block17
block17:
	r11 = &RuneFilter{Min: r9, Max: r10}
	goto block18
block18:
	ret0 = r11
	goto block19
block19:
	return
block20:
	return
}
func MatchRune(frame *runtime.State) (ret0 *RuneRangeMatch) {
	var r0 bool
	var r1 []*RuneFilter
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 bool
	var r6 []*RuneFilter
	var r7 int
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 bool
	var r12 int
	var r13 []*RuneFilter
	var r14 *RuneFilter
	var r15 []*RuneFilter
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 bool
	var r20 []*RuneFilter
	var r21 *RuneRangeMatch
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block39
	}
block2:
	r3 = '['
	goto block3
block3:
	r4 = r2 == r3
	goto block4
block4:
	if r4 {
		goto block5
	} else {
		goto block38
	}
block5:
	frame.Consume()
	goto block6
block6:
	r5 = false
	goto block7
block7:
	r0 = r5
	goto block8
block8:
	r6 = []*RuneFilter{}
	goto block9
block9:
	r1 = r6
	goto block10
block10:
	r7 = frame.Checkpoint()
	goto block11
block11:
	r8 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block19
	}
block12:
	r9 = '^'
	goto block13
block13:
	r10 = r8 == r9
	goto block14
block14:
	if r10 {
		goto block15
	} else {
		goto block18
	}
block15:
	frame.Consume()
	goto block16
block16:
	r11 = true
	goto block17
block17:
	r0 = r11
	goto block20
block18:
	frame.Fail()
	goto block19
block19:
	frame.Recover(r7)
	goto block20
block20:
	r12 = frame.Checkpoint()
	goto block21
block21:
	r13 = r1
	goto block22
block22:
	r14 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block25
	}
block23:
	r15 = append(r13, r14)
	goto block24
block24:
	r1 = r15
	goto block20
block25:
	frame.Recover(r12)
	goto block26
block26:
	r16 = frame.Peek()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block39
	}
block27:
	r17 = ']'
	goto block28
block28:
	r18 = r16 == r17
	goto block29
block29:
	if r18 {
		goto block30
	} else {
		goto block37
	}
block30:
	frame.Consume()
	goto block31
block31:
	S(frame)
	if frame.Flow == 0 {
		goto block32
	} else {
		goto block39
	}
block32:
	r19 = r0
	goto block33
block33:
	r20 = r1
	goto block34
block34:
	r21 = &RuneRangeMatch{Invert: r19, Filters: r20}
	goto block35
block35:
	ret0 = r21
	goto block36
block36:
	return
block37:
	frame.Fail()
	goto block39
block38:
	frame.Fail()
	goto block39
block39:
	return
}
func Atom(frame *runtime.State) (ret0 TextMatch) {
	var r0 string
	var r1 TextMatch
	var r2 int
	var r3 *RuneRangeMatch
	var r4 string
	var r5 string
	var r6 *StringLiteralMatch
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 TextMatch
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 TextMatch
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Checkpoint()
	goto block2
block2:
	r3 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r3
	goto block28
block4:
	frame.Recover(r2)
	goto block5
block5:
	r4 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block11
	}
block6:
	r0 = r4
	goto block7
block7:
	S(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block11
	}
block8:
	r5 = r0
	goto block9
block9:
	r6 = &StringLiteralMatch{Value: r5}
	goto block10
block10:
	ret0 = r6
	goto block28
block11:
	frame.Recover(r2)
	goto block12
block12:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block31
	}
block13:
	r8 = '('
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
	S(frame)
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block31
	}
block18:
	r10 = ParseMatchChoice(frame)
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block31
	}
block19:
	r1 = r10
	goto block20
block20:
	r11 = frame.Peek()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block31
	}
block21:
	r12 = ')'
	goto block22
block22:
	r13 = r11 == r12
	goto block23
block23:
	if r13 {
		goto block24
	} else {
		goto block29
	}
block24:
	frame.Consume()
	goto block25
block25:
	S(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block31
	}
block26:
	r14 = r1
	goto block27
block27:
	ret0 = r14
	goto block28
block28:
	return
block29:
	frame.Fail()
	goto block31
block30:
	frame.Fail()
	goto block31
block31:
	return
}
func MatchPostfix(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 TextMatch
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 TextMatch
	var r7 int
	var r8 *MatchRepeat
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 TextMatch
	var r13 int
	var r14 *MatchRepeat
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 TextMatch
	var r19 []TextMatch
	var r20 *MatchSequence
	var r21 []TextMatch
	var r22 *MatchChoice
	var r23 TextMatch
	goto block0
block0:
	goto block1
block1:
	r1 = Atom(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block45
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = frame.Checkpoint()
	goto block4
block4:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block15
	}
block5:
	r4 = '*'
	goto block6
block6:
	r5 = r3 == r4
	goto block7
block7:
	if r5 {
		goto block8
	} else {
		goto block14
	}
block8:
	frame.Consume()
	goto block9
block9:
	S(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block15
	}
block10:
	r6 = r0
	goto block11
block11:
	r7 = 0
	goto block12
block12:
	r8 = &MatchRepeat{Match: r6, Min: r7}
	goto block13
block13:
	ret0 = r8
	goto block44
block14:
	frame.Fail()
	goto block15
block15:
	frame.Recover(r2)
	goto block16
block16:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block27
	}
block17:
	r10 = '+'
	goto block18
block18:
	r11 = r9 == r10
	goto block19
block19:
	if r11 {
		goto block20
	} else {
		goto block26
	}
block20:
	frame.Consume()
	goto block21
block21:
	S(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block27
	}
block22:
	r12 = r0
	goto block23
block23:
	r13 = 1
	goto block24
block24:
	r14 = &MatchRepeat{Match: r12, Min: r13}
	goto block25
block25:
	ret0 = r14
	goto block44
block26:
	frame.Fail()
	goto block27
block27:
	frame.Recover(r2)
	goto block28
block28:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block41
	}
block29:
	r16 = '?'
	goto block30
block30:
	r17 = r15 == r16
	goto block31
block31:
	if r17 {
		goto block32
	} else {
		goto block40
	}
block32:
	frame.Consume()
	goto block33
block33:
	S(frame)
	if frame.Flow == 0 {
		goto block34
	} else {
		goto block41
	}
block34:
	r18 = r0
	goto block35
block35:
	r19 = []TextMatch{}
	goto block36
block36:
	r20 = &MatchSequence{Matches: r19}
	goto block37
block37:
	r21 = []TextMatch{r18, r20}
	goto block38
block38:
	r22 = &MatchChoice{Matches: r21}
	goto block39
block39:
	ret0 = r22
	goto block44
block40:
	frame.Fail()
	goto block41
block41:
	frame.Recover(r2)
	goto block42
block42:
	r23 = r0
	goto block43
block43:
	ret0 = r23
	goto block44
block44:
	return
block45:
	return
}
func MatchPrefix(frame *runtime.State) (ret0 TextMatch) {
	var r0 bool
	var r1 int
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 bool
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 bool
	var r11 TextMatch
	var r12 *MatchLookahead
	var r13 TextMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Checkpoint()
	goto block2
block2:
	r0 = false
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
	r0 = r6
	goto block18
block11:
	frame.Fail()
	goto block12
block12:
	frame.Recover(r2)
	goto block13
block13:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block24
	}
block14:
	r8 = '&'
	goto block15
block15:
	r9 = r7 == r8
	goto block16
block16:
	if r9 {
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
	r10 = r0
	goto block20
block20:
	r11 = MatchPostfix(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block24
	}
block21:
	r12 = &MatchLookahead{Invert: r10, Match: r11}
	goto block22
block22:
	ret0 = r12
	goto block27
block23:
	frame.Fail()
	goto block24
block24:
	frame.Recover(r1)
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
	var r1 []TextMatch
	var r2 TextMatch
	var r3 int
	var r4 TextMatch
	var r5 []TextMatch
	var r6 []TextMatch
	var r7 TextMatch
	var r8 []TextMatch
	var r9 int
	var r10 []TextMatch
	var r11 TextMatch
	var r12 []TextMatch
	var r13 []TextMatch
	var r14 *MatchSequence
	var r15 TextMatch
	goto block0
block0:
	goto block1
block1:
	r2 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block24
	}
block2:
	r0 = r2
	goto block3
block3:
	r3 = frame.Checkpoint()
	goto block4
block4:
	r4 = r0
	goto block5
block5:
	r5 = []TextMatch{r4}
	goto block6
block6:
	r1 = r5
	goto block7
block7:
	r6 = r1
	goto block8
block8:
	r7 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block20
	}
block9:
	r8 = append(r6, r7)
	goto block10
block10:
	r1 = r8
	goto block11
block11:
	r9 = frame.Checkpoint()
	goto block12
block12:
	r10 = r1
	goto block13
block13:
	r11 = MatchPrefix(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block16
	}
block14:
	r12 = append(r10, r11)
	goto block15
block15:
	r1 = r12
	goto block11
block16:
	frame.Recover(r9)
	goto block17
block17:
	r13 = r1
	goto block18
block18:
	r14 = &MatchSequence{Matches: r13}
	goto block19
block19:
	ret0 = r14
	goto block23
block20:
	frame.Recover(r3)
	goto block21
block21:
	r15 = r0
	goto block22
block22:
	ret0 = r15
	goto block23
block23:
	return
block24:
	return
}
func ParseMatchChoice(frame *runtime.State) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 []TextMatch
	var r2 TextMatch
	var r3 int
	var r4 TextMatch
	var r5 []TextMatch
	var r6 rune
	var r7 rune
	var r8 bool
	var r9 []TextMatch
	var r10 TextMatch
	var r11 []TextMatch
	var r12 int
	var r13 rune
	var r14 rune
	var r15 bool
	var r16 []TextMatch
	var r17 TextMatch
	var r18 []TextMatch
	var r19 []TextMatch
	var r20 *MatchChoice
	var r21 TextMatch
	goto block0
block0:
	goto block1
block1:
	r2 = Sequence(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block38
	}
block2:
	r0 = r2
	goto block3
block3:
	r3 = frame.Checkpoint()
	goto block4
block4:
	r4 = r0
	goto block5
block5:
	r5 = []TextMatch{r4}
	goto block6
block6:
	r1 = r5
	goto block7
block7:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block34
	}
block8:
	r7 = '|'
	goto block9
block9:
	r8 = r6 == r7
	goto block10
block10:
	if r8 {
		goto block11
	} else {
		goto block33
	}
block11:
	frame.Consume()
	goto block12
block12:
	S(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block34
	}
block13:
	r9 = r1
	goto block14
block14:
	r10 = Sequence(frame)
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block34
	}
block15:
	r11 = append(r9, r10)
	goto block16
block16:
	r1 = r11
	goto block17
block17:
	r12 = frame.Checkpoint()
	goto block18
block18:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block29
	}
block19:
	r14 = '|'
	goto block20
block20:
	r15 = r13 == r14
	goto block21
block21:
	if r15 {
		goto block22
	} else {
		goto block28
	}
block22:
	frame.Consume()
	goto block23
block23:
	S(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block29
	}
block24:
	r16 = r1
	goto block25
block25:
	r17 = Sequence(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block29
	}
block26:
	r18 = append(r16, r17)
	goto block27
block27:
	r1 = r18
	goto block17
block28:
	frame.Fail()
	goto block29
block29:
	frame.Recover(r12)
	goto block30
block30:
	r19 = r1
	goto block31
block31:
	r20 = &MatchChoice{Matches: r19}
	goto block32
block32:
	ret0 = r20
	goto block37
block33:
	frame.Fail()
	goto block34
block34:
	frame.Recover(r3)
	goto block35
block35:
	r21 = r0
	goto block36
block36:
	ret0 = r21
	goto block37
block37:
	return
block38:
	return
}
func ParseExprList(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 []ASTExpr
	var r1 []ASTExpr
	var r2 int
	var r3 []ASTExpr
	var r4 ASTExpr
	var r5 []ASTExpr
	var r6 int
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 []ASTExpr
	var r11 ASTExpr
	var r12 []ASTExpr
	var r13 []ASTExpr
	goto block0
block0:
	goto block1
block1:
	r1 = []ASTExpr{}
	goto block2
block2:
	r0 = r1
	goto block3
block3:
	r2 = frame.Checkpoint()
	goto block4
block4:
	r3 = r0
	goto block5
block5:
	r4 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block21
	}
block6:
	r5 = append(r3, r4)
	goto block7
block7:
	r0 = r5
	goto block8
block8:
	r6 = frame.Checkpoint()
	goto block9
block9:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block20
	}
block10:
	r8 = ','
	goto block11
block11:
	r9 = r7 == r8
	goto block12
block12:
	if r9 {
		goto block13
	} else {
		goto block19
	}
block13:
	frame.Consume()
	goto block14
block14:
	S(frame)
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block20
	}
block15:
	r10 = r0
	goto block16
block16:
	r11 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block20
	}
block17:
	r12 = append(r10, r11)
	goto block18
block18:
	r0 = r12
	goto block8
block19:
	frame.Fail()
	goto block20
block20:
	frame.Recover(r6)
	goto block22
block21:
	frame.Recover(r2)
	goto block22
block22:
	r13 = r0
	goto block23
block23:
	ret0 = r13
	goto block24
block24:
	return
}
func ParseNamedExpr(frame *runtime.State) (ret0 *NamedExpr) {
	var r0 string
	var r1 string
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 string
	var r6 ASTExpr
	var r7 *NamedExpr
	goto block0
block0:
	goto block1
block1:
	r1 = Ident(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block15
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = frame.Peek()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block15
	}
block4:
	r3 = ':'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block14
	}
block7:
	frame.Consume()
	goto block8
block8:
	S(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block15
	}
block9:
	r5 = r0
	goto block10
block10:
	r6 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block15
	}
block11:
	r7 = &NamedExpr{Name: r5, Expr: r6}
	goto block12
block12:
	ret0 = r7
	goto block13
block13:
	return
block14:
	frame.Fail()
	goto block15
block15:
	return
}
func ParseNamedExprList(frame *runtime.State) (ret0 []*NamedExpr) {
	var r0 []*NamedExpr
	var r1 []*NamedExpr
	var r2 int
	var r3 []*NamedExpr
	var r4 *NamedExpr
	var r5 []*NamedExpr
	var r6 int
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 []*NamedExpr
	var r11 *NamedExpr
	var r12 []*NamedExpr
	var r13 []*NamedExpr
	goto block0
block0:
	goto block1
block1:
	r1 = []*NamedExpr{}
	goto block2
block2:
	r0 = r1
	goto block3
block3:
	r2 = frame.Checkpoint()
	goto block4
block4:
	r3 = r0
	goto block5
block5:
	r4 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block21
	}
block6:
	r5 = append(r3, r4)
	goto block7
block7:
	r0 = r5
	goto block8
block8:
	r6 = frame.Checkpoint()
	goto block9
block9:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block20
	}
block10:
	r8 = ','
	goto block11
block11:
	r9 = r7 == r8
	goto block12
block12:
	if r9 {
		goto block13
	} else {
		goto block19
	}
block13:
	frame.Consume()
	goto block14
block14:
	S(frame)
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block20
	}
block15:
	r10 = r0
	goto block16
block16:
	r11 = ParseNamedExpr(frame)
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block20
	}
block17:
	r12 = append(r10, r11)
	goto block18
block18:
	r0 = r12
	goto block8
block19:
	frame.Fail()
	goto block20
block20:
	frame.Recover(r6)
	goto block22
block21:
	frame.Recover(r2)
	goto block22
block22:
	r13 = r0
	goto block23
block23:
	ret0 = r13
	goto block24
block24:
	return
}
func ParseTypeList(frame *runtime.State) (ret0 []ASTTypeRef) {
	var r0 []ASTTypeRef
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 []ASTTypeRef
	var r5 int
	var r6 []ASTTypeRef
	var r7 ASTTypeRef
	var r8 []ASTTypeRef
	var r9 int
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 []ASTTypeRef
	var r14 ASTTypeRef
	var r15 []ASTTypeRef
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 []ASTTypeRef
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block39
	}
block2:
	r2 = '('
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block38
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block39
	}
block7:
	r4 = []ASTTypeRef{}
	goto block8
block8:
	r0 = r4
	goto block9
block9:
	r5 = frame.Checkpoint()
	goto block10
block10:
	r6 = r0
	goto block11
block11:
	r7 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block27
	}
block12:
	r8 = append(r6, r7)
	goto block13
block13:
	r0 = r8
	goto block14
block14:
	r9 = frame.Checkpoint()
	goto block15
block15:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block26
	}
block16:
	r11 = ','
	goto block17
block17:
	r12 = r10 == r11
	goto block18
block18:
	if r12 {
		goto block19
	} else {
		goto block25
	}
block19:
	frame.Consume()
	goto block20
block20:
	S(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block26
	}
block21:
	r13 = r0
	goto block22
block22:
	r14 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block26
	}
block23:
	r15 = append(r13, r14)
	goto block24
block24:
	r0 = r15
	goto block14
block25:
	frame.Fail()
	goto block26
block26:
	frame.Recover(r9)
	goto block28
block27:
	frame.Recover(r5)
	goto block28
block28:
	r16 = frame.Peek()
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block39
	}
block29:
	r17 = ')'
	goto block30
block30:
	r18 = r16 == r17
	goto block31
block31:
	if r18 {
		goto block32
	} else {
		goto block37
	}
block32:
	frame.Consume()
	goto block33
block33:
	S(frame)
	if frame.Flow == 0 {
		goto block34
	} else {
		goto block39
	}
block34:
	r19 = r0
	goto block35
block35:
	ret0 = r19
	goto block36
block36:
	return
block37:
	frame.Fail()
	goto block39
block38:
	frame.Fail()
	goto block39
block39:
	return
}
func ParseExpr(frame *runtime.State) (ret0 ASTExpr) {
	var r0 int
	var r1 []ASTExpr
	var r2 [][]ASTExpr
	var r3 []ASTExpr
	var r4 []ASTExpr
	var r5 ASTExpr
	var r6 []ASTExpr
	var r7 string
	var r8 ASTTypeRef
	var r9 ASTExpr
	var r10 ASTTypeRef
	var r11 ASTExpr
	var r12 string
	var r13 ASTExpr
	var r14 []ASTExpr
	var r15 string
	var r16 string
	var r17 ASTExpr
	var r18 ASTExpr
	var r19 *TypeRef
	var r20 []*NamedExpr
	var r21 *ListTypeRef
	var r22 []ASTExpr
	var r23 string
	var r24 bool
	var r25 ASTExpr
	var r26 int
	var r27 ASTExpr
	var r28 int
	var r29 rune
	var r30 rune
	var r31 bool
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
	var r50 rune
	var r51 rune
	var r52 bool
	var r53 int
	var r54 []ASTExpr
	var r55 []ASTExpr
	var r56 int
	var r57 *Repeat
	var r58 rune
	var r59 rune
	var r60 bool
	var r61 rune
	var r62 rune
	var r63 bool
	var r64 rune
	var r65 rune
	var r66 bool
	var r67 rune
	var r68 rune
	var r69 bool
	var r70 rune
	var r71 rune
	var r72 bool
	var r73 rune
	var r74 rune
	var r75 bool
	var r76 []ASTExpr
	var r77 [][]ASTExpr
	var r78 int
	var r79 rune
	var r80 rune
	var r81 bool
	var r82 rune
	var r83 rune
	var r84 bool
	var r85 [][]ASTExpr
	var r86 []ASTExpr
	var r87 [][]ASTExpr
	var r88 [][]ASTExpr
	var r89 *Choice
	var r90 rune
	var r91 rune
	var r92 bool
	var r93 rune
	var r94 rune
	var r95 bool
	var r96 rune
	var r97 rune
	var r98 bool
	var r99 rune
	var r100 rune
	var r101 bool
	var r102 rune
	var r103 rune
	var r104 bool
	var r105 rune
	var r106 rune
	var r107 bool
	var r108 rune
	var r109 rune
	var r110 bool
	var r111 rune
	var r112 rune
	var r113 bool
	var r114 []ASTExpr
	var r115 []ASTExpr
	var r116 *Optional
	var r117 rune
	var r118 rune
	var r119 bool
	var r120 rune
	var r121 rune
	var r122 bool
	var r123 rune
	var r124 rune
	var r125 bool
	var r126 rune
	var r127 rune
	var r128 bool
	var r129 rune
	var r130 rune
	var r131 bool
	var r132 []ASTExpr
	var r133 []ASTExpr
	var r134 *Slice
	var r135 rune
	var r136 rune
	var r137 bool
	var r138 rune
	var r139 rune
	var r140 bool
	var r141 ASTExpr
	var r142 []ASTExpr
	var r143 ASTExpr
	var r144 []ASTExpr
	var r145 *If
	var r146 rune
	var r147 rune
	var r148 bool
	var r149 rune
	var r150 rune
	var r151 bool
	var r152 rune
	var r153 rune
	var r154 bool
	var r155 string
	var r156 ASTTypeRef
	var r157 int
	var r158 rune
	var r159 rune
	var r160 bool
	var r161 ASTExpr
	var r162 ASTExpr
	var r163 string
	var r164 ASTTypeRef
	var r165 bool
	var r166 *Assign
	var r167 rune
	var r168 rune
	var r169 bool
	var r170 rune
	var r171 rune
	var r172 bool
	var r173 rune
	var r174 rune
	var r175 bool
	var r176 rune
	var r177 rune
	var r178 bool
	var r179 *Fail
	var r180 rune
	var r181 rune
	var r182 bool
	var r183 rune
	var r184 rune
	var r185 bool
	var r186 rune
	var r187 rune
	var r188 bool
	var r189 rune
	var r190 rune
	var r191 bool
	var r192 rune
	var r193 rune
	var r194 bool
	var r195 rune
	var r196 rune
	var r197 bool
	var r198 ASTTypeRef
	var r199 ASTExpr
	var r200 ASTTypeRef
	var r201 ASTExpr
	var r202 *Coerce
	var r203 rune
	var r204 rune
	var r205 bool
	var r206 rune
	var r207 rune
	var r208 bool
	var r209 rune
	var r210 rune
	var r211 bool
	var r212 rune
	var r213 rune
	var r214 bool
	var r215 rune
	var r216 rune
	var r217 bool
	var r218 rune
	var r219 rune
	var r220 bool
	var r221 string
	var r222 ASTExpr
	var r223 string
	var r224 *GetName
	var r225 ASTExpr
	var r226 *Append
	var r227 string
	var r228 *Assign
	var r229 rune
	var r230 rune
	var r231 bool
	var r232 rune
	var r233 rune
	var r234 bool
	var r235 rune
	var r236 rune
	var r237 bool
	var r238 rune
	var r239 rune
	var r240 bool
	var r241 rune
	var r242 rune
	var r243 bool
	var r244 rune
	var r245 rune
	var r246 bool
	var r247 int
	var r248 rune
	var r249 rune
	var r250 bool
	var r251 []ASTExpr
	var r252 rune
	var r253 rune
	var r254 bool
	var r255 []ASTExpr
	var r256 *Return
	var r257 ASTExpr
	var r258 []ASTExpr
	var r259 *Return
	var r260 []ASTExpr
	var r261 *Return
	var r262 string
	var r263 rune
	var r264 rune
	var r265 bool
	var r266 rune
	var r267 rune
	var r268 bool
	var r269 string
	var r270 *Call
	var r271 string
	var r272 ASTExpr
	var r273 ASTExpr
	var r274 ASTExpr
	var r275 string
	var r276 ASTExpr
	var r277 *BinaryOp
	var r278 *TypeRef
	var r279 rune
	var r280 rune
	var r281 bool
	var r282 []*NamedExpr
	var r283 rune
	var r284 rune
	var r285 bool
	var r286 *TypeRef
	var r287 []*NamedExpr
	var r288 *Construct
	var r289 *ListTypeRef
	var r290 rune
	var r291 rune
	var r292 bool
	var r293 []ASTExpr
	var r294 rune
	var r295 rune
	var r296 bool
	var r297 *ListTypeRef
	var r298 []ASTExpr
	var r299 *ConstructList
	var r300 *StringMatch
	var r301 *RuneMatch
	var r302 string
	var r303 int
	var r304 int
	var r305 rune
	var r306 rune
	var r307 bool
	var r308 rune
	var r309 rune
	var r310 bool
	var r311 bool
	var r312 rune
	var r313 rune
	var r314 bool
	var r315 ASTExpr
	var r316 ASTExpr
	var r317 string
	var r318 bool
	var r319 *Assign
	var r320 string
	var r321 *GetName
	goto block0
block0:
	goto block1
block1:
	r26 = frame.Checkpoint()
	goto block2
block2:
	r27 = Literal(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r27
	goto block597
block4:
	frame.Recover(r26)
	goto block5
block5:
	r0 = 0
	goto block6
block6:
	r28 = frame.Checkpoint()
	goto block7
block7:
	r29 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block31
	}
block8:
	r30 = 's'
	goto block9
block9:
	r31 = r29 == r30
	goto block10
block10:
	if r31 {
		goto block11
	} else {
		goto block30
	}
block11:
	frame.Consume()
	goto block12
block12:
	r32 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block31
	}
block13:
	r33 = 't'
	goto block14
block14:
	r34 = r32 == r33
	goto block15
block15:
	if r34 {
		goto block16
	} else {
		goto block29
	}
block16:
	frame.Consume()
	goto block17
block17:
	r35 = frame.Peek()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block31
	}
block18:
	r36 = 'a'
	goto block19
block19:
	r37 = r35 == r36
	goto block20
block20:
	if r37 {
		goto block21
	} else {
		goto block28
	}
block21:
	frame.Consume()
	goto block22
block22:
	r38 = frame.Peek()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block31
	}
block23:
	r39 = 'r'
	goto block24
block24:
	r40 = r38 == r39
	goto block25
block25:
	if r40 {
		goto block26
	} else {
		goto block27
	}
block26:
	frame.Consume()
	goto block54
block27:
	frame.Fail()
	goto block31
block28:
	frame.Fail()
	goto block31
block29:
	frame.Fail()
	goto block31
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(r28)
	goto block32
block32:
	r41 = frame.Peek()
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block65
	}
block33:
	r42 = 'p'
	goto block34
block34:
	r43 = r41 == r42
	goto block35
block35:
	if r43 {
		goto block36
	} else {
		goto block64
	}
block36:
	frame.Consume()
	goto block37
block37:
	r44 = frame.Peek()
	if frame.Flow == 0 {
		goto block38
	} else {
		goto block65
	}
block38:
	r45 = 'l'
	goto block39
block39:
	r46 = r44 == r45
	goto block40
block40:
	if r46 {
		goto block41
	} else {
		goto block63
	}
block41:
	frame.Consume()
	goto block42
block42:
	r47 = frame.Peek()
	if frame.Flow == 0 {
		goto block43
	} else {
		goto block65
	}
block43:
	r48 = 'u'
	goto block44
block44:
	r49 = r47 == r48
	goto block45
block45:
	if r49 {
		goto block46
	} else {
		goto block62
	}
block46:
	frame.Consume()
	goto block47
block47:
	r50 = frame.Peek()
	if frame.Flow == 0 {
		goto block48
	} else {
		goto block65
	}
block48:
	r51 = 's'
	goto block49
block49:
	r52 = r50 == r51
	goto block50
block50:
	if r52 {
		goto block51
	} else {
		goto block61
	}
block51:
	frame.Consume()
	goto block52
block52:
	r53 = 1
	goto block53
block53:
	r0 = r53
	goto block54
block54:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block55
	} else {
		goto block65
	}
block55:
	r54 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block56
	} else {
		goto block65
	}
block56:
	r1 = r54
	goto block57
block57:
	r55 = r1
	goto block58
block58:
	r56 = r0
	goto block59
block59:
	r57 = &Repeat{Block: r55, Min: r56}
	goto block60
block60:
	ret0 = r57
	goto block597
block61:
	frame.Fail()
	goto block65
block62:
	frame.Fail()
	goto block65
block63:
	frame.Fail()
	goto block65
block64:
	frame.Fail()
	goto block65
block65:
	frame.Recover(r26)
	goto block66
block66:
	r58 = frame.Peek()
	if frame.Flow == 0 {
		goto block67
	} else {
		goto block128
	}
block67:
	r59 = 'c'
	goto block68
block68:
	r60 = r58 == r59
	goto block69
block69:
	if r60 {
		goto block70
	} else {
		goto block127
	}
block70:
	frame.Consume()
	goto block71
block71:
	r61 = frame.Peek()
	if frame.Flow == 0 {
		goto block72
	} else {
		goto block128
	}
block72:
	r62 = 'h'
	goto block73
block73:
	r63 = r61 == r62
	goto block74
block74:
	if r63 {
		goto block75
	} else {
		goto block126
	}
block75:
	frame.Consume()
	goto block76
block76:
	r64 = frame.Peek()
	if frame.Flow == 0 {
		goto block77
	} else {
		goto block128
	}
block77:
	r65 = 'o'
	goto block78
block78:
	r66 = r64 == r65
	goto block79
block79:
	if r66 {
		goto block80
	} else {
		goto block125
	}
block80:
	frame.Consume()
	goto block81
block81:
	r67 = frame.Peek()
	if frame.Flow == 0 {
		goto block82
	} else {
		goto block128
	}
block82:
	r68 = 'o'
	goto block83
block83:
	r69 = r67 == r68
	goto block84
block84:
	if r69 {
		goto block85
	} else {
		goto block124
	}
block85:
	frame.Consume()
	goto block86
block86:
	r70 = frame.Peek()
	if frame.Flow == 0 {
		goto block87
	} else {
		goto block128
	}
block87:
	r71 = 's'
	goto block88
block88:
	r72 = r70 == r71
	goto block89
block89:
	if r72 {
		goto block90
	} else {
		goto block123
	}
block90:
	frame.Consume()
	goto block91
block91:
	r73 = frame.Peek()
	if frame.Flow == 0 {
		goto block92
	} else {
		goto block128
	}
block92:
	r74 = 'e'
	goto block93
block93:
	r75 = r73 == r74
	goto block94
block94:
	if r75 {
		goto block95
	} else {
		goto block122
	}
block95:
	frame.Consume()
	goto block96
block96:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block97
	} else {
		goto block128
	}
block97:
	r76 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block98
	} else {
		goto block128
	}
block98:
	r77 = [][]ASTExpr{r76}
	goto block99
block99:
	r2 = r77
	goto block100
block100:
	r78 = frame.Checkpoint()
	goto block101
block101:
	r79 = frame.Peek()
	if frame.Flow == 0 {
		goto block102
	} else {
		goto block118
	}
block102:
	r80 = 'o'
	goto block103
block103:
	r81 = r79 == r80
	goto block104
block104:
	if r81 {
		goto block105
	} else {
		goto block117
	}
block105:
	frame.Consume()
	goto block106
block106:
	r82 = frame.Peek()
	if frame.Flow == 0 {
		goto block107
	} else {
		goto block118
	}
block107:
	r83 = 'r'
	goto block108
block108:
	r84 = r82 == r83
	goto block109
block109:
	if r84 {
		goto block110
	} else {
		goto block116
	}
block110:
	frame.Consume()
	goto block111
block111:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block112
	} else {
		goto block118
	}
block112:
	r85 = r2
	goto block113
block113:
	r86 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block114
	} else {
		goto block118
	}
block114:
	r87 = append(r85, r86)
	goto block115
block115:
	r2 = r87
	goto block100
block116:
	frame.Fail()
	goto block118
block117:
	frame.Fail()
	goto block118
block118:
	frame.Recover(r78)
	goto block119
block119:
	r88 = r2
	goto block120
block120:
	r89 = &Choice{Blocks: r88}
	goto block121
block121:
	ret0 = r89
	goto block597
block122:
	frame.Fail()
	goto block128
block123:
	frame.Fail()
	goto block128
block124:
	frame.Fail()
	goto block128
block125:
	frame.Fail()
	goto block128
block126:
	frame.Fail()
	goto block128
block127:
	frame.Fail()
	goto block128
block128:
	frame.Recover(r26)
	goto block129
block129:
	r90 = frame.Peek()
	if frame.Flow == 0 {
		goto block130
	} else {
		goto block183
	}
block130:
	r91 = 'q'
	goto block131
block131:
	r92 = r90 == r91
	goto block132
block132:
	if r92 {
		goto block133
	} else {
		goto block182
	}
block133:
	frame.Consume()
	goto block134
block134:
	r93 = frame.Peek()
	if frame.Flow == 0 {
		goto block135
	} else {
		goto block183
	}
block135:
	r94 = 'u'
	goto block136
block136:
	r95 = r93 == r94
	goto block137
block137:
	if r95 {
		goto block138
	} else {
		goto block181
	}
block138:
	frame.Consume()
	goto block139
block139:
	r96 = frame.Peek()
	if frame.Flow == 0 {
		goto block140
	} else {
		goto block183
	}
block140:
	r97 = 'e'
	goto block141
block141:
	r98 = r96 == r97
	goto block142
block142:
	if r98 {
		goto block143
	} else {
		goto block180
	}
block143:
	frame.Consume()
	goto block144
block144:
	r99 = frame.Peek()
	if frame.Flow == 0 {
		goto block145
	} else {
		goto block183
	}
block145:
	r100 = 's'
	goto block146
block146:
	r101 = r99 == r100
	goto block147
block147:
	if r101 {
		goto block148
	} else {
		goto block179
	}
block148:
	frame.Consume()
	goto block149
block149:
	r102 = frame.Peek()
	if frame.Flow == 0 {
		goto block150
	} else {
		goto block183
	}
block150:
	r103 = 't'
	goto block151
block151:
	r104 = r102 == r103
	goto block152
block152:
	if r104 {
		goto block153
	} else {
		goto block178
	}
block153:
	frame.Consume()
	goto block154
block154:
	r105 = frame.Peek()
	if frame.Flow == 0 {
		goto block155
	} else {
		goto block183
	}
block155:
	r106 = 'i'
	goto block156
block156:
	r107 = r105 == r106
	goto block157
block157:
	if r107 {
		goto block158
	} else {
		goto block177
	}
block158:
	frame.Consume()
	goto block159
block159:
	r108 = frame.Peek()
	if frame.Flow == 0 {
		goto block160
	} else {
		goto block183
	}
block160:
	r109 = 'o'
	goto block161
block161:
	r110 = r108 == r109
	goto block162
block162:
	if r110 {
		goto block163
	} else {
		goto block176
	}
block163:
	frame.Consume()
	goto block164
block164:
	r111 = frame.Peek()
	if frame.Flow == 0 {
		goto block165
	} else {
		goto block183
	}
block165:
	r112 = 'n'
	goto block166
block166:
	r113 = r111 == r112
	goto block167
block167:
	if r113 {
		goto block168
	} else {
		goto block175
	}
block168:
	frame.Consume()
	goto block169
block169:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block170
	} else {
		goto block183
	}
block170:
	r114 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block171
	} else {
		goto block183
	}
block171:
	r3 = r114
	goto block172
block172:
	r115 = r3
	goto block173
block173:
	r116 = &Optional{Block: r115}
	goto block174
block174:
	ret0 = r116
	goto block597
block175:
	frame.Fail()
	goto block183
block176:
	frame.Fail()
	goto block183
block177:
	frame.Fail()
	goto block183
block178:
	frame.Fail()
	goto block183
block179:
	frame.Fail()
	goto block183
block180:
	frame.Fail()
	goto block183
block181:
	frame.Fail()
	goto block183
block182:
	frame.Fail()
	goto block183
block183:
	frame.Recover(r26)
	goto block184
block184:
	r117 = frame.Peek()
	if frame.Flow == 0 {
		goto block185
	} else {
		goto block220
	}
block185:
	r118 = 's'
	goto block186
block186:
	r119 = r117 == r118
	goto block187
block187:
	if r119 {
		goto block188
	} else {
		goto block219
	}
block188:
	frame.Consume()
	goto block189
block189:
	r120 = frame.Peek()
	if frame.Flow == 0 {
		goto block190
	} else {
		goto block220
	}
block190:
	r121 = 'l'
	goto block191
block191:
	r122 = r120 == r121
	goto block192
block192:
	if r122 {
		goto block193
	} else {
		goto block218
	}
block193:
	frame.Consume()
	goto block194
block194:
	r123 = frame.Peek()
	if frame.Flow == 0 {
		goto block195
	} else {
		goto block220
	}
block195:
	r124 = 'i'
	goto block196
block196:
	r125 = r123 == r124
	goto block197
block197:
	if r125 {
		goto block198
	} else {
		goto block217
	}
block198:
	frame.Consume()
	goto block199
block199:
	r126 = frame.Peek()
	if frame.Flow == 0 {
		goto block200
	} else {
		goto block220
	}
block200:
	r127 = 'c'
	goto block201
block201:
	r128 = r126 == r127
	goto block202
block202:
	if r128 {
		goto block203
	} else {
		goto block216
	}
block203:
	frame.Consume()
	goto block204
block204:
	r129 = frame.Peek()
	if frame.Flow == 0 {
		goto block205
	} else {
		goto block220
	}
block205:
	r130 = 'e'
	goto block206
block206:
	r131 = r129 == r130
	goto block207
block207:
	if r131 {
		goto block208
	} else {
		goto block215
	}
block208:
	frame.Consume()
	goto block209
block209:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block210
	} else {
		goto block220
	}
block210:
	r132 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block211
	} else {
		goto block220
	}
block211:
	r4 = r132
	goto block212
block212:
	r133 = r4
	goto block213
block213:
	r134 = &Slice{Block: r133}
	goto block214
block214:
	ret0 = r134
	goto block597
block215:
	frame.Fail()
	goto block220
block216:
	frame.Fail()
	goto block220
block217:
	frame.Fail()
	goto block220
block218:
	frame.Fail()
	goto block220
block219:
	frame.Fail()
	goto block220
block220:
	frame.Recover(r26)
	goto block221
block221:
	r135 = frame.Peek()
	if frame.Flow == 0 {
		goto block222
	} else {
		goto block242
	}
block222:
	r136 = 'i'
	goto block223
block223:
	r137 = r135 == r136
	goto block224
block224:
	if r137 {
		goto block225
	} else {
		goto block241
	}
block225:
	frame.Consume()
	goto block226
block226:
	r138 = frame.Peek()
	if frame.Flow == 0 {
		goto block227
	} else {
		goto block242
	}
block227:
	r139 = 'f'
	goto block228
block228:
	r140 = r138 == r139
	goto block229
block229:
	if r140 {
		goto block230
	} else {
		goto block240
	}
block230:
	frame.Consume()
	goto block231
block231:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block232
	} else {
		goto block242
	}
block232:
	r141 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block233
	} else {
		goto block242
	}
block233:
	r5 = r141
	goto block234
block234:
	r142 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block235
	} else {
		goto block242
	}
block235:
	r6 = r142
	goto block236
block236:
	r143 = r5
	goto block237
block237:
	r144 = r6
	goto block238
block238:
	r145 = &If{Expr: r143, Block: r144}
	goto block239
block239:
	ret0 = r145
	goto block597
block240:
	frame.Fail()
	goto block242
block241:
	frame.Fail()
	goto block242
block242:
	frame.Recover(r26)
	goto block243
block243:
	r146 = frame.Peek()
	if frame.Flow == 0 {
		goto block244
	} else {
		goto block284
	}
block244:
	r147 = 'v'
	goto block245
block245:
	r148 = r146 == r147
	goto block246
block246:
	if r148 {
		goto block247
	} else {
		goto block283
	}
block247:
	frame.Consume()
	goto block248
block248:
	r149 = frame.Peek()
	if frame.Flow == 0 {
		goto block249
	} else {
		goto block284
	}
block249:
	r150 = 'a'
	goto block250
block250:
	r151 = r149 == r150
	goto block251
block251:
	if r151 {
		goto block252
	} else {
		goto block282
	}
block252:
	frame.Consume()
	goto block253
block253:
	r152 = frame.Peek()
	if frame.Flow == 0 {
		goto block254
	} else {
		goto block284
	}
block254:
	r153 = 'r'
	goto block255
block255:
	r154 = r152 == r153
	goto block256
block256:
	if r154 {
		goto block257
	} else {
		goto block281
	}
block257:
	frame.Consume()
	goto block258
block258:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block259
	} else {
		goto block284
	}
block259:
	r155 = Ident(frame)
	if frame.Flow == 0 {
		goto block260
	} else {
		goto block284
	}
block260:
	r7 = r155
	goto block261
block261:
	r156 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block262
	} else {
		goto block284
	}
block262:
	r8 = r156
	goto block263
block263:
	r9 = nil
	goto block264
block264:
	r157 = frame.Checkpoint()
	goto block265
block265:
	r158 = frame.Peek()
	if frame.Flow == 0 {
		goto block266
	} else {
		goto block274
	}
block266:
	r159 = '='
	goto block267
block267:
	r160 = r158 == r159
	goto block268
block268:
	if r160 {
		goto block269
	} else {
		goto block273
	}
block269:
	frame.Consume()
	goto block270
block270:
	S(frame)
	if frame.Flow == 0 {
		goto block271
	} else {
		goto block274
	}
block271:
	r161 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block272
	} else {
		goto block274
	}
block272:
	r9 = r161
	goto block275
block273:
	frame.Fail()
	goto block274
block274:
	frame.Recover(r157)
	goto block275
block275:
	r162 = r9
	goto block276
block276:
	r163 = r7
	goto block277
block277:
	r164 = r8
	goto block278
block278:
	r165 = true
	goto block279
block279:
	r166 = &Assign{Expr: r162, Name: r163, Type: r164, Define: r165}
	goto block280
block280:
	ret0 = r166
	goto block597
block281:
	frame.Fail()
	goto block284
block282:
	frame.Fail()
	goto block284
block283:
	frame.Fail()
	goto block284
block284:
	frame.Recover(r26)
	goto block285
block285:
	r167 = frame.Peek()
	if frame.Flow == 0 {
		goto block286
	} else {
		goto block312
	}
block286:
	r168 = 'f'
	goto block287
block287:
	r169 = r167 == r168
	goto block288
block288:
	if r169 {
		goto block289
	} else {
		goto block311
	}
block289:
	frame.Consume()
	goto block290
block290:
	r170 = frame.Peek()
	if frame.Flow == 0 {
		goto block291
	} else {
		goto block312
	}
block291:
	r171 = 'a'
	goto block292
block292:
	r172 = r170 == r171
	goto block293
block293:
	if r172 {
		goto block294
	} else {
		goto block310
	}
block294:
	frame.Consume()
	goto block295
block295:
	r173 = frame.Peek()
	if frame.Flow == 0 {
		goto block296
	} else {
		goto block312
	}
block296:
	r174 = 'i'
	goto block297
block297:
	r175 = r173 == r174
	goto block298
block298:
	if r175 {
		goto block299
	} else {
		goto block309
	}
block299:
	frame.Consume()
	goto block300
block300:
	r176 = frame.Peek()
	if frame.Flow == 0 {
		goto block301
	} else {
		goto block312
	}
block301:
	r177 = 'l'
	goto block302
block302:
	r178 = r176 == r177
	goto block303
block303:
	if r178 {
		goto block304
	} else {
		goto block308
	}
block304:
	frame.Consume()
	goto block305
block305:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block306
	} else {
		goto block312
	}
block306:
	r179 = &Fail{}
	goto block307
block307:
	ret0 = r179
	goto block597
block308:
	frame.Fail()
	goto block312
block309:
	frame.Fail()
	goto block312
block310:
	frame.Fail()
	goto block312
block311:
	frame.Fail()
	goto block312
block312:
	frame.Recover(r26)
	goto block313
block313:
	r180 = frame.Peek()
	if frame.Flow == 0 {
		goto block314
	} else {
		goto block358
	}
block314:
	r181 = 'c'
	goto block315
block315:
	r182 = r180 == r181
	goto block316
block316:
	if r182 {
		goto block317
	} else {
		goto block357
	}
block317:
	frame.Consume()
	goto block318
block318:
	r183 = frame.Peek()
	if frame.Flow == 0 {
		goto block319
	} else {
		goto block358
	}
block319:
	r184 = 'o'
	goto block320
block320:
	r185 = r183 == r184
	goto block321
block321:
	if r185 {
		goto block322
	} else {
		goto block356
	}
block322:
	frame.Consume()
	goto block323
block323:
	r186 = frame.Peek()
	if frame.Flow == 0 {
		goto block324
	} else {
		goto block358
	}
block324:
	r187 = 'e'
	goto block325
block325:
	r188 = r186 == r187
	goto block326
block326:
	if r188 {
		goto block327
	} else {
		goto block355
	}
block327:
	frame.Consume()
	goto block328
block328:
	r189 = frame.Peek()
	if frame.Flow == 0 {
		goto block329
	} else {
		goto block358
	}
block329:
	r190 = 'r'
	goto block330
block330:
	r191 = r189 == r190
	goto block331
block331:
	if r191 {
		goto block332
	} else {
		goto block354
	}
block332:
	frame.Consume()
	goto block333
block333:
	r192 = frame.Peek()
	if frame.Flow == 0 {
		goto block334
	} else {
		goto block358
	}
block334:
	r193 = 'c'
	goto block335
block335:
	r194 = r192 == r193
	goto block336
block336:
	if r194 {
		goto block337
	} else {
		goto block353
	}
block337:
	frame.Consume()
	goto block338
block338:
	r195 = frame.Peek()
	if frame.Flow == 0 {
		goto block339
	} else {
		goto block358
	}
block339:
	r196 = 'e'
	goto block340
block340:
	r197 = r195 == r196
	goto block341
block341:
	if r197 {
		goto block342
	} else {
		goto block352
	}
block342:
	frame.Consume()
	goto block343
block343:
	S(frame)
	if frame.Flow == 0 {
		goto block344
	} else {
		goto block358
	}
block344:
	r198 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block345
	} else {
		goto block358
	}
block345:
	r10 = r198
	goto block346
block346:
	r199 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block347
	} else {
		goto block358
	}
block347:
	r11 = r199
	goto block348
block348:
	r200 = r10
	goto block349
block349:
	r201 = r11
	goto block350
block350:
	r202 = &Coerce{Type: r200, Expr: r201}
	goto block351
block351:
	ret0 = r202
	goto block597
block352:
	frame.Fail()
	goto block358
block353:
	frame.Fail()
	goto block358
block354:
	frame.Fail()
	goto block358
block355:
	frame.Fail()
	goto block358
block356:
	frame.Fail()
	goto block358
block357:
	frame.Fail()
	goto block358
block358:
	frame.Recover(r26)
	goto block359
block359:
	r203 = frame.Peek()
	if frame.Flow == 0 {
		goto block360
	} else {
		goto block407
	}
block360:
	r204 = 'a'
	goto block361
block361:
	r205 = r203 == r204
	goto block362
block362:
	if r205 {
		goto block363
	} else {
		goto block406
	}
block363:
	frame.Consume()
	goto block364
block364:
	r206 = frame.Peek()
	if frame.Flow == 0 {
		goto block365
	} else {
		goto block407
	}
block365:
	r207 = 'p'
	goto block366
block366:
	r208 = r206 == r207
	goto block367
block367:
	if r208 {
		goto block368
	} else {
		goto block405
	}
block368:
	frame.Consume()
	goto block369
block369:
	r209 = frame.Peek()
	if frame.Flow == 0 {
		goto block370
	} else {
		goto block407
	}
block370:
	r210 = 'p'
	goto block371
block371:
	r211 = r209 == r210
	goto block372
block372:
	if r211 {
		goto block373
	} else {
		goto block404
	}
block373:
	frame.Consume()
	goto block374
block374:
	r212 = frame.Peek()
	if frame.Flow == 0 {
		goto block375
	} else {
		goto block407
	}
block375:
	r213 = 'e'
	goto block376
block376:
	r214 = r212 == r213
	goto block377
block377:
	if r214 {
		goto block378
	} else {
		goto block403
	}
block378:
	frame.Consume()
	goto block379
block379:
	r215 = frame.Peek()
	if frame.Flow == 0 {
		goto block380
	} else {
		goto block407
	}
block380:
	r216 = 'n'
	goto block381
block381:
	r217 = r215 == r216
	goto block382
block382:
	if r217 {
		goto block383
	} else {
		goto block402
	}
block383:
	frame.Consume()
	goto block384
block384:
	r218 = frame.Peek()
	if frame.Flow == 0 {
		goto block385
	} else {
		goto block407
	}
block385:
	r219 = 'd'
	goto block386
block386:
	r220 = r218 == r219
	goto block387
block387:
	if r220 {
		goto block388
	} else {
		goto block401
	}
block388:
	frame.Consume()
	goto block389
block389:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block390
	} else {
		goto block407
	}
block390:
	r221 = Ident(frame)
	if frame.Flow == 0 {
		goto block391
	} else {
		goto block407
	}
block391:
	r12 = r221
	goto block392
block392:
	r222 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block393
	} else {
		goto block407
	}
block393:
	r13 = r222
	goto block394
block394:
	r223 = r12
	goto block395
block395:
	r224 = &GetName{Name: r223}
	goto block396
block396:
	r225 = r13
	goto block397
block397:
	r226 = &Append{List: r224, Expr: r225}
	goto block398
block398:
	r227 = r12
	goto block399
block399:
	r228 = &Assign{Expr: r226, Name: r227}
	goto block400
block400:
	ret0 = r228
	goto block597
block401:
	frame.Fail()
	goto block407
block402:
	frame.Fail()
	goto block407
block403:
	frame.Fail()
	goto block407
block404:
	frame.Fail()
	goto block407
block405:
	frame.Fail()
	goto block407
block406:
	frame.Fail()
	goto block407
block407:
	frame.Recover(r26)
	goto block408
block408:
	r229 = frame.Peek()
	if frame.Flow == 0 {
		goto block409
	} else {
		goto block474
	}
block409:
	r230 = 'r'
	goto block410
block410:
	r231 = r229 == r230
	goto block411
block411:
	if r231 {
		goto block412
	} else {
		goto block473
	}
block412:
	frame.Consume()
	goto block413
block413:
	r232 = frame.Peek()
	if frame.Flow == 0 {
		goto block414
	} else {
		goto block474
	}
block414:
	r233 = 'e'
	goto block415
block415:
	r234 = r232 == r233
	goto block416
block416:
	if r234 {
		goto block417
	} else {
		goto block472
	}
block417:
	frame.Consume()
	goto block418
block418:
	r235 = frame.Peek()
	if frame.Flow == 0 {
		goto block419
	} else {
		goto block474
	}
block419:
	r236 = 't'
	goto block420
block420:
	r237 = r235 == r236
	goto block421
block421:
	if r237 {
		goto block422
	} else {
		goto block471
	}
block422:
	frame.Consume()
	goto block423
block423:
	r238 = frame.Peek()
	if frame.Flow == 0 {
		goto block424
	} else {
		goto block474
	}
block424:
	r239 = 'u'
	goto block425
block425:
	r240 = r238 == r239
	goto block426
block426:
	if r240 {
		goto block427
	} else {
		goto block470
	}
block427:
	frame.Consume()
	goto block428
block428:
	r241 = frame.Peek()
	if frame.Flow == 0 {
		goto block429
	} else {
		goto block474
	}
block429:
	r242 = 'r'
	goto block430
block430:
	r243 = r241 == r242
	goto block431
block431:
	if r243 {
		goto block432
	} else {
		goto block469
	}
block432:
	frame.Consume()
	goto block433
block433:
	r244 = frame.Peek()
	if frame.Flow == 0 {
		goto block434
	} else {
		goto block474
	}
block434:
	r245 = 'n'
	goto block435
block435:
	r246 = r244 == r245
	goto block436
block436:
	if r246 {
		goto block437
	} else {
		goto block468
	}
block437:
	frame.Consume()
	goto block438
block438:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block439
	} else {
		goto block474
	}
block439:
	r247 = frame.Checkpoint()
	goto block440
block440:
	r248 = frame.Peek()
	if frame.Flow == 0 {
		goto block441
	} else {
		goto block459
	}
block441:
	r249 = '('
	goto block442
block442:
	r250 = r248 == r249
	goto block443
block443:
	if r250 {
		goto block444
	} else {
		goto block458
	}
block444:
	frame.Consume()
	goto block445
block445:
	S(frame)
	if frame.Flow == 0 {
		goto block446
	} else {
		goto block459
	}
block446:
	r251 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block447
	} else {
		goto block459
	}
block447:
	r14 = r251
	goto block448
block448:
	r252 = frame.Peek()
	if frame.Flow == 0 {
		goto block449
	} else {
		goto block459
	}
block449:
	r253 = ')'
	goto block450
block450:
	r254 = r252 == r253
	goto block451
block451:
	if r254 {
		goto block452
	} else {
		goto block457
	}
block452:
	frame.Consume()
	goto block453
block453:
	S(frame)
	if frame.Flow == 0 {
		goto block454
	} else {
		goto block459
	}
block454:
	r255 = r14
	goto block455
block455:
	r256 = &Return{Exprs: r255}
	goto block456
block456:
	ret0 = r256
	goto block597
block457:
	frame.Fail()
	goto block459
block458:
	frame.Fail()
	goto block459
block459:
	frame.Recover(r247)
	goto block460
block460:
	r257 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block461
	} else {
		goto block464
	}
block461:
	r258 = []ASTExpr{r257}
	goto block462
block462:
	r259 = &Return{Exprs: r258}
	goto block463
block463:
	ret0 = r259
	goto block597
block464:
	frame.Recover(r247)
	goto block465
block465:
	r260 = []ASTExpr{}
	goto block466
block466:
	r261 = &Return{Exprs: r260}
	goto block467
block467:
	ret0 = r261
	goto block597
block468:
	frame.Fail()
	goto block474
block469:
	frame.Fail()
	goto block474
block470:
	frame.Fail()
	goto block474
block471:
	frame.Fail()
	goto block474
block472:
	frame.Fail()
	goto block474
block473:
	frame.Fail()
	goto block474
block474:
	frame.Recover(r26)
	goto block475
block475:
	r262 = Ident(frame)
	if frame.Flow == 0 {
		goto block476
	} else {
		goto block494
	}
block476:
	r15 = r262
	goto block477
block477:
	r263 = frame.Peek()
	if frame.Flow == 0 {
		goto block478
	} else {
		goto block494
	}
block478:
	r264 = '('
	goto block479
block479:
	r265 = r263 == r264
	goto block480
block480:
	if r265 {
		goto block481
	} else {
		goto block493
	}
block481:
	frame.Consume()
	goto block482
block482:
	S(frame)
	if frame.Flow == 0 {
		goto block483
	} else {
		goto block494
	}
block483:
	r266 = frame.Peek()
	if frame.Flow == 0 {
		goto block484
	} else {
		goto block494
	}
block484:
	r267 = ')'
	goto block485
block485:
	r268 = r266 == r267
	goto block486
block486:
	if r268 {
		goto block487
	} else {
		goto block492
	}
block487:
	frame.Consume()
	goto block488
block488:
	S(frame)
	if frame.Flow == 0 {
		goto block489
	} else {
		goto block494
	}
block489:
	r269 = r15
	goto block490
block490:
	r270 = &Call{Name: r269}
	goto block491
block491:
	ret0 = r270
	goto block597
block492:
	frame.Fail()
	goto block494
block493:
	frame.Fail()
	goto block494
block494:
	frame.Recover(r26)
	goto block495
block495:
	r271 = BinaryOperator(frame)
	if frame.Flow == 0 {
		goto block496
	} else {
		goto block506
	}
block496:
	r16 = r271
	goto block497
block497:
	r272 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block498
	} else {
		goto block506
	}
block498:
	r17 = r272
	goto block499
block499:
	r273 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block500
	} else {
		goto block506
	}
block500:
	r18 = r273
	goto block501
block501:
	r274 = r17
	goto block502
block502:
	r275 = r16
	goto block503
block503:
	r276 = r18
	goto block504
block504:
	r277 = &BinaryOp{Left: r274, Op: r275, Right: r276}
	goto block505
block505:
	ret0 = r277
	goto block597
block506:
	frame.Recover(r26)
	goto block507
block507:
	r278 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block508
	} else {
		goto block529
	}
block508:
	r19 = r278
	goto block509
block509:
	r279 = frame.Peek()
	if frame.Flow == 0 {
		goto block510
	} else {
		goto block529
	}
block510:
	r280 = '{'
	goto block511
block511:
	r281 = r279 == r280
	goto block512
block512:
	if r281 {
		goto block513
	} else {
		goto block528
	}
block513:
	frame.Consume()
	goto block514
block514:
	S(frame)
	if frame.Flow == 0 {
		goto block515
	} else {
		goto block529
	}
block515:
	r282 = ParseNamedExprList(frame)
	if frame.Flow == 0 {
		goto block516
	} else {
		goto block529
	}
block516:
	r20 = r282
	goto block517
block517:
	r283 = frame.Peek()
	if frame.Flow == 0 {
		goto block518
	} else {
		goto block529
	}
block518:
	r284 = '}'
	goto block519
block519:
	r285 = r283 == r284
	goto block520
block520:
	if r285 {
		goto block521
	} else {
		goto block527
	}
block521:
	frame.Consume()
	goto block522
block522:
	S(frame)
	if frame.Flow == 0 {
		goto block523
	} else {
		goto block529
	}
block523:
	r286 = r19
	goto block524
block524:
	r287 = r20
	goto block525
block525:
	r288 = &Construct{Type: r286, Args: r287}
	goto block526
block526:
	ret0 = r288
	goto block597
block527:
	frame.Fail()
	goto block529
block528:
	frame.Fail()
	goto block529
block529:
	frame.Recover(r26)
	goto block530
block530:
	r289 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block531
	} else {
		goto block552
	}
block531:
	r21 = r289
	goto block532
block532:
	r290 = frame.Peek()
	if frame.Flow == 0 {
		goto block533
	} else {
		goto block552
	}
block533:
	r291 = '{'
	goto block534
block534:
	r292 = r290 == r291
	goto block535
block535:
	if r292 {
		goto block536
	} else {
		goto block551
	}
block536:
	frame.Consume()
	goto block537
block537:
	S(frame)
	if frame.Flow == 0 {
		goto block538
	} else {
		goto block552
	}
block538:
	r293 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block539
	} else {
		goto block552
	}
block539:
	r22 = r293
	goto block540
block540:
	r294 = frame.Peek()
	if frame.Flow == 0 {
		goto block541
	} else {
		goto block552
	}
block541:
	r295 = '}'
	goto block542
block542:
	r296 = r294 == r295
	goto block543
block543:
	if r296 {
		goto block544
	} else {
		goto block550
	}
block544:
	frame.Consume()
	goto block545
block545:
	S(frame)
	if frame.Flow == 0 {
		goto block546
	} else {
		goto block552
	}
block546:
	r297 = r21
	goto block547
block547:
	r298 = r22
	goto block548
block548:
	r299 = &ConstructList{Type: r297, Args: r298}
	goto block549
block549:
	ret0 = r299
	goto block597
block550:
	frame.Fail()
	goto block552
block551:
	frame.Fail()
	goto block552
block552:
	frame.Recover(r26)
	goto block553
block553:
	r300 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		goto block554
	} else {
		goto block555
	}
block554:
	ret0 = r300
	goto block597
block555:
	frame.Recover(r26)
	goto block556
block556:
	r301 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		goto block557
	} else {
		goto block558
	}
block557:
	ret0 = r301
	goto block597
block558:
	frame.Recover(r26)
	goto block559
block559:
	r302 = Ident(frame)
	if frame.Flow == 0 {
		goto block560
	} else {
		goto block598
	}
block560:
	r23 = r302
	goto block561
block561:
	r303 = frame.Checkpoint()
	goto block562
block562:
	r24 = false
	goto block563
block563:
	r304 = frame.Checkpoint()
	goto block564
block564:
	r305 = frame.Peek()
	if frame.Flow == 0 {
		goto block565
	} else {
		goto block578
	}
block565:
	r306 = ':'
	goto block566
block566:
	r307 = r305 == r306
	goto block567
block567:
	if r307 {
		goto block568
	} else {
		goto block577
	}
block568:
	frame.Consume()
	goto block569
block569:
	r308 = frame.Peek()
	if frame.Flow == 0 {
		goto block570
	} else {
		goto block578
	}
block570:
	r309 = '='
	goto block571
block571:
	r310 = r308 == r309
	goto block572
block572:
	if r310 {
		goto block573
	} else {
		goto block576
	}
block573:
	frame.Consume()
	goto block574
block574:
	r311 = true
	goto block575
block575:
	r24 = r311
	goto block584
block576:
	frame.Fail()
	goto block578
block577:
	frame.Fail()
	goto block578
block578:
	frame.Recover(r304)
	goto block579
block579:
	r312 = frame.Peek()
	if frame.Flow == 0 {
		goto block580
	} else {
		goto block593
	}
block580:
	r313 = '='
	goto block581
block581:
	r314 = r312 == r313
	goto block582
block582:
	if r314 {
		goto block583
	} else {
		goto block592
	}
block583:
	frame.Consume()
	goto block584
block584:
	S(frame)
	if frame.Flow == 0 {
		goto block585
	} else {
		goto block593
	}
block585:
	r315 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block586
	} else {
		goto block593
	}
block586:
	r25 = r315
	goto block587
block587:
	r316 = r25
	goto block588
block588:
	r317 = r23
	goto block589
block589:
	r318 = r24
	goto block590
block590:
	r319 = &Assign{Expr: r316, Name: r317, Define: r318}
	goto block591
block591:
	ret0 = r319
	goto block597
block592:
	frame.Fail()
	goto block593
block593:
	frame.Recover(r303)
	goto block594
block594:
	r320 = r23
	goto block595
block595:
	r321 = &GetName{Name: r320}
	goto block596
block596:
	ret0 = r321
	goto block597
block597:
	return
block598:
	return
}
func ParseCodeBlock(frame *runtime.State) (ret0 []ASTExpr) {
	var r0 []ASTExpr
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 []ASTExpr
	var r5 int
	var r6 []ASTExpr
	var r7 ASTExpr
	var r8 []ASTExpr
	var r9 int
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 rune
	var r14 rune
	var r15 bool
	var r16 []ASTExpr
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block35
	}
block2:
	r2 = '{'
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block34
	}
block5:
	frame.Consume()
	goto block6
block6:
	S(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block35
	}
block7:
	r4 = []ASTExpr{}
	goto block8
block8:
	r0 = r4
	goto block9
block9:
	r5 = frame.Checkpoint()
	goto block10
block10:
	r6 = r0
	goto block11
block11:
	r7 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block23
	}
block12:
	r8 = append(r6, r7)
	goto block13
block13:
	r0 = r8
	goto block14
block14:
	r9 = frame.Checkpoint()
	goto block15
block15:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block22
	}
block16:
	r11 = ';'
	goto block17
block17:
	r12 = r10 == r11
	goto block18
block18:
	if r12 {
		goto block19
	} else {
		goto block21
	}
block19:
	frame.Consume()
	goto block20
block20:
	S(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block22
	}
block21:
	frame.Fail()
	goto block22
block22:
	frame.Recover(r9)
	goto block9
block23:
	frame.Recover(r5)
	goto block24
block24:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block35
	}
block25:
	r14 = '}'
	goto block26
block26:
	r15 = r13 == r14
	goto block27
block27:
	if r15 {
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
	r16 = r0
	goto block31
block31:
	ret0 = r16
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
func ParseStructDecl(frame *runtime.State) (ret0 *StructDecl) {
	var r0 string
	var r1 ASTTypeRef
	var r2 []*FieldDecl
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
	var r18 rune
	var r19 rune
	var r20 bool
	var r21 string
	var r22 int
	var r23 rune
	var r24 rune
	var r25 bool
	var r26 rune
	var r27 rune
	var r28 bool
	var r29 rune
	var r30 rune
	var r31 bool
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
	var r50 rune
	var r51 rune
	var r52 bool
	var r53 ASTTypeRef
	var r54 rune
	var r55 rune
	var r56 bool
	var r57 []*FieldDecl
	var r58 int
	var r59 []*FieldDecl
	var r60 string
	var r61 ASTTypeRef
	var r62 *FieldDecl
	var r63 []*FieldDecl
	var r64 rune
	var r65 rune
	var r66 bool
	var r67 string
	var r68 ASTTypeRef
	var r69 []*FieldDecl
	var r70 *StructDecl
	goto block0
block0:
	goto block1
block1:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block136
	}
block2:
	r4 = 's'
	goto block3
block3:
	r5 = r3 == r4
	goto block4
block4:
	if r5 {
		goto block5
	} else {
		goto block135
	}
block5:
	frame.Consume()
	goto block6
block6:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block136
	}
block7:
	r7 = 't'
	goto block8
block8:
	r8 = r6 == r7
	goto block9
block9:
	if r8 {
		goto block10
	} else {
		goto block134
	}
block10:
	frame.Consume()
	goto block11
block11:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block136
	}
block12:
	r10 = 'r'
	goto block13
block13:
	r11 = r9 == r10
	goto block14
block14:
	if r11 {
		goto block15
	} else {
		goto block133
	}
block15:
	frame.Consume()
	goto block16
block16:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block136
	}
block17:
	r13 = 'u'
	goto block18
block18:
	r14 = r12 == r13
	goto block19
block19:
	if r14 {
		goto block20
	} else {
		goto block132
	}
block20:
	frame.Consume()
	goto block21
block21:
	r15 = frame.Peek()
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block136
	}
block22:
	r16 = 'c'
	goto block23
block23:
	r17 = r15 == r16
	goto block24
block24:
	if r17 {
		goto block25
	} else {
		goto block131
	}
block25:
	frame.Consume()
	goto block26
block26:
	r18 = frame.Peek()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block136
	}
block27:
	r19 = 't'
	goto block28
block28:
	r20 = r18 == r19
	goto block29
block29:
	if r20 {
		goto block30
	} else {
		goto block130
	}
block30:
	frame.Consume()
	goto block31
block31:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block32
	} else {
		goto block136
	}
block32:
	r21 = Ident(frame)
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block136
	}
block33:
	r0 = r21
	goto block34
block34:
	r1 = nil
	goto block35
block35:
	r22 = frame.Checkpoint()
	goto block36
block36:
	r23 = frame.Peek()
	if frame.Flow == 0 {
		goto block37
	} else {
		goto block99
	}
block37:
	r24 = 'i'
	goto block38
block38:
	r25 = r23 == r24
	goto block39
block39:
	if r25 {
		goto block40
	} else {
		goto block98
	}
block40:
	frame.Consume()
	goto block41
block41:
	r26 = frame.Peek()
	if frame.Flow == 0 {
		goto block42
	} else {
		goto block99
	}
block42:
	r27 = 'm'
	goto block43
block43:
	r28 = r26 == r27
	goto block44
block44:
	if r28 {
		goto block45
	} else {
		goto block97
	}
block45:
	frame.Consume()
	goto block46
block46:
	r29 = frame.Peek()
	if frame.Flow == 0 {
		goto block47
	} else {
		goto block99
	}
block47:
	r30 = 'p'
	goto block48
block48:
	r31 = r29 == r30
	goto block49
block49:
	if r31 {
		goto block50
	} else {
		goto block96
	}
block50:
	frame.Consume()
	goto block51
block51:
	r32 = frame.Peek()
	if frame.Flow == 0 {
		goto block52
	} else {
		goto block99
	}
block52:
	r33 = 'l'
	goto block53
block53:
	r34 = r32 == r33
	goto block54
block54:
	if r34 {
		goto block55
	} else {
		goto block95
	}
block55:
	frame.Consume()
	goto block56
block56:
	r35 = frame.Peek()
	if frame.Flow == 0 {
		goto block57
	} else {
		goto block99
	}
block57:
	r36 = 'e'
	goto block58
block58:
	r37 = r35 == r36
	goto block59
block59:
	if r37 {
		goto block60
	} else {
		goto block94
	}
block60:
	frame.Consume()
	goto block61
block61:
	r38 = frame.Peek()
	if frame.Flow == 0 {
		goto block62
	} else {
		goto block99
	}
block62:
	r39 = 'm'
	goto block63
block63:
	r40 = r38 == r39
	goto block64
block64:
	if r40 {
		goto block65
	} else {
		goto block93
	}
block65:
	frame.Consume()
	goto block66
block66:
	r41 = frame.Peek()
	if frame.Flow == 0 {
		goto block67
	} else {
		goto block99
	}
block67:
	r42 = 'e'
	goto block68
block68:
	r43 = r41 == r42
	goto block69
block69:
	if r43 {
		goto block70
	} else {
		goto block92
	}
block70:
	frame.Consume()
	goto block71
block71:
	r44 = frame.Peek()
	if frame.Flow == 0 {
		goto block72
	} else {
		goto block99
	}
block72:
	r45 = 'n'
	goto block73
block73:
	r46 = r44 == r45
	goto block74
block74:
	if r46 {
		goto block75
	} else {
		goto block91
	}
block75:
	frame.Consume()
	goto block76
block76:
	r47 = frame.Peek()
	if frame.Flow == 0 {
		goto block77
	} else {
		goto block99
	}
block77:
	r48 = 't'
	goto block78
block78:
	r49 = r47 == r48
	goto block79
block79:
	if r49 {
		goto block80
	} else {
		goto block90
	}
block80:
	frame.Consume()
	goto block81
block81:
	r50 = frame.Peek()
	if frame.Flow == 0 {
		goto block82
	} else {
		goto block99
	}
block82:
	r51 = 's'
	goto block83
block83:
	r52 = r50 == r51
	goto block84
block84:
	if r52 {
		goto block85
	} else {
		goto block89
	}
block85:
	frame.Consume()
	goto block86
block86:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block87
	} else {
		goto block99
	}
block87:
	r53 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block88
	} else {
		goto block99
	}
block88:
	r1 = r53
	goto block100
block89:
	frame.Fail()
	goto block99
block90:
	frame.Fail()
	goto block99
block91:
	frame.Fail()
	goto block99
block92:
	frame.Fail()
	goto block99
block93:
	frame.Fail()
	goto block99
block94:
	frame.Fail()
	goto block99
block95:
	frame.Fail()
	goto block99
block96:
	frame.Fail()
	goto block99
block97:
	frame.Fail()
	goto block99
block98:
	frame.Fail()
	goto block99
block99:
	frame.Recover(r22)
	goto block100
block100:
	r54 = frame.Peek()
	if frame.Flow == 0 {
		goto block101
	} else {
		goto block136
	}
block101:
	r55 = '{'
	goto block102
block102:
	r56 = r54 == r55
	goto block103
block103:
	if r56 {
		goto block104
	} else {
		goto block129
	}
block104:
	frame.Consume()
	goto block105
block105:
	S(frame)
	if frame.Flow == 0 {
		goto block106
	} else {
		goto block136
	}
block106:
	r57 = []*FieldDecl{}
	goto block107
block107:
	r2 = r57
	goto block108
block108:
	r58 = frame.Checkpoint()
	goto block109
block109:
	r59 = r2
	goto block110
block110:
	r60 = Ident(frame)
	if frame.Flow == 0 {
		goto block111
	} else {
		goto block115
	}
block111:
	r61 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block112
	} else {
		goto block115
	}
block112:
	r62 = &FieldDecl{Name: r60, Type: r61}
	goto block113
block113:
	r63 = append(r59, r62)
	goto block114
block114:
	r2 = r63
	goto block108
block115:
	frame.Recover(r58)
	goto block116
block116:
	r64 = frame.Peek()
	if frame.Flow == 0 {
		goto block117
	} else {
		goto block136
	}
block117:
	r65 = '}'
	goto block118
block118:
	r66 = r64 == r65
	goto block119
block119:
	if r66 {
		goto block120
	} else {
		goto block128
	}
block120:
	frame.Consume()
	goto block121
block121:
	S(frame)
	if frame.Flow == 0 {
		goto block122
	} else {
		goto block136
	}
block122:
	r67 = r0
	goto block123
block123:
	r68 = r1
	goto block124
block124:
	r69 = r2
	goto block125
block125:
	r70 = &StructDecl{Name: r67, Implements: r68, Fields: r69}
	goto block126
block126:
	ret0 = r70
	goto block127
block127:
	return
block128:
	frame.Fail()
	goto block136
block129:
	frame.Fail()
	goto block136
block130:
	frame.Fail()
	goto block136
block131:
	frame.Fail()
	goto block136
block132:
	frame.Fail()
	goto block136
block133:
	frame.Fail()
	goto block136
block134:
	frame.Fail()
	goto block136
block135:
	frame.Fail()
	goto block136
block136:
	return
}
func ParseFuncDecl(frame *runtime.State) (ret0 *FuncDecl) {
	var r0 string
	var r1 []ASTTypeRef
	var r2 []ASTExpr
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
	var r15 string
	var r16 []ASTTypeRef
	var r17 []ASTExpr
	var r18 string
	var r19 []ASTTypeRef
	var r20 []ASTExpr
	var r21 *FuncDecl
	goto block0
block0:
	goto block1
block1:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block38
	}
block2:
	r4 = 'f'
	goto block3
block3:
	r5 = r3 == r4
	goto block4
block4:
	if r5 {
		goto block5
	} else {
		goto block37
	}
block5:
	frame.Consume()
	goto block6
block6:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block38
	}
block7:
	r7 = 'u'
	goto block8
block8:
	r8 = r6 == r7
	goto block9
block9:
	if r8 {
		goto block10
	} else {
		goto block36
	}
block10:
	frame.Consume()
	goto block11
block11:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block38
	}
block12:
	r10 = 'n'
	goto block13
block13:
	r11 = r9 == r10
	goto block14
block14:
	if r11 {
		goto block15
	} else {
		goto block35
	}
block15:
	frame.Consume()
	goto block16
block16:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block38
	}
block17:
	r13 = 'c'
	goto block18
block18:
	r14 = r12 == r13
	goto block19
block19:
	if r14 {
		goto block20
	} else {
		goto block34
	}
block20:
	frame.Consume()
	goto block21
block21:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block38
	}
block22:
	r15 = Ident(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block38
	}
block23:
	r0 = r15
	goto block24
block24:
	r16 = ParseTypeList(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block38
	}
block25:
	r1 = r16
	goto block26
block26:
	r17 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block38
	}
block27:
	r2 = r17
	goto block28
block28:
	r18 = r0
	goto block29
block29:
	r19 = r1
	goto block30
block30:
	r20 = r2
	goto block31
block31:
	r21 = &FuncDecl{Name: r18, ReturnTypes: r19, Block: r20}
	goto block32
block32:
	ret0 = r21
	goto block33
block33:
	return
block34:
	frame.Fail()
	goto block38
block35:
	frame.Fail()
	goto block38
block36:
	frame.Fail()
	goto block38
block37:
	frame.Fail()
	goto block38
block38:
	return
}
func ParseTest(frame *runtime.State) (ret0 *Test) {
	var r0 string
	var r1 string
	var r2 string
	var r3 Destructure
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
	var r16 string
	var r17 string
	var r18 string
	var r19 Destructure
	var r20 string
	var r21 string
	var r22 string
	var r23 Destructure
	var r24 *Test
	goto block0
block0:
	goto block1
block1:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block42
	}
block2:
	r5 = 't'
	goto block3
block3:
	r6 = r4 == r5
	goto block4
block4:
	if r6 {
		goto block5
	} else {
		goto block41
	}
block5:
	frame.Consume()
	goto block6
block6:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block42
	}
block7:
	r8 = 'e'
	goto block8
block8:
	r9 = r7 == r8
	goto block9
block9:
	if r9 {
		goto block10
	} else {
		goto block40
	}
block10:
	frame.Consume()
	goto block11
block11:
	r10 = frame.Peek()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block42
	}
block12:
	r11 = 's'
	goto block13
block13:
	r12 = r10 == r11
	goto block14
block14:
	if r12 {
		goto block15
	} else {
		goto block39
	}
block15:
	frame.Consume()
	goto block16
block16:
	r13 = frame.Peek()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block42
	}
block17:
	r14 = 't'
	goto block18
block18:
	r15 = r13 == r14
	goto block19
block19:
	if r15 {
		goto block20
	} else {
		goto block38
	}
block20:
	frame.Consume()
	goto block21
block21:
	EndKeyword(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block42
	}
block22:
	r16 = Ident(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block42
	}
block23:
	r0 = r16
	goto block24
block24:
	r17 = Ident(frame)
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block42
	}
block25:
	r1 = r17
	goto block26
block26:
	r18 = DecodeString(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block42
	}
block27:
	r2 = r18
	goto block28
block28:
	S(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block42
	}
block29:
	r19 = ParseDestructure(frame)
	if frame.Flow == 0 {
		goto block30
	} else {
		goto block42
	}
block30:
	r3 = r19
	goto block31
block31:
	r20 = r0
	goto block32
block32:
	r21 = r1
	goto block33
block33:
	r22 = r2
	goto block34
block34:
	r23 = r3
	goto block35
block35:
	r24 = &Test{Rule: r20, Name: r21, Input: r22, Destructure: r23}
	goto block36
block36:
	ret0 = r24
	goto block37
block37:
	return
block38:
	frame.Fail()
	goto block42
block39:
	frame.Fail()
	goto block42
block40:
	frame.Fail()
	goto block42
block41:
	frame.Fail()
	goto block42
block42:
	return
}
func ParseFile(frame *runtime.State) (ret0 *File) {
	var r0 []ASTDecl
	var r1 []*Test
	var r2 []ASTDecl
	var r3 []*Test
	var r4 int
	var r5 int
	var r6 []ASTDecl
	var r7 *FuncDecl
	var r8 []ASTDecl
	var r9 []ASTDecl
	var r10 *StructDecl
	var r11 []ASTDecl
	var r12 []*Test
	var r13 *Test
	var r14 []*Test
	var r15 int
	var r16 []ASTDecl
	var r17 []*Test
	var r18 *File
	goto block0
block0:
	goto block1
block1:
	r2 = []ASTDecl{}
	goto block2
block2:
	r0 = r2
	goto block3
block3:
	r3 = []*Test{}
	goto block4
block4:
	r1 = r3
	goto block5
block5:
	r4 = frame.Checkpoint()
	goto block6
block6:
	r5 = frame.Checkpoint()
	goto block7
block7:
	r6 = r0
	goto block8
block8:
	r7 = ParseFuncDecl(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block11
	}
block9:
	r8 = append(r6, r7)
	goto block10
block10:
	r0 = r8
	goto block5
block11:
	frame.Recover(r5)
	goto block12
block12:
	r9 = r0
	goto block13
block13:
	r10 = ParseStructDecl(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block16
	}
block14:
	r11 = append(r9, r10)
	goto block15
block15:
	r0 = r11
	goto block5
block16:
	frame.Recover(r5)
	goto block17
block17:
	r12 = r1
	goto block18
block18:
	r13 = ParseTest(frame)
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block21
	}
block19:
	r14 = append(r12, r13)
	goto block20
block20:
	r1 = r14
	goto block5
block21:
	frame.Recover(r4)
	goto block22
block22:
	r15 = frame.LookaheadBegin()
	goto block23
block23:
	frame.Peek()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block27
	}
block24:
	frame.Consume()
	goto block25
block25:
	frame.LookaheadFail(r15)
	goto block26
block26:
	return
block27:
	frame.LookaheadNormal(r15)
	goto block28
block28:
	r16 = r0
	goto block29
block29:
	r17 = r1
	goto block30
block30:
	r18 = &File{Decls: r16, Tests: r17}
	goto block31
block31:
	ret0 = r18
	goto block32
block32:
	return
}
