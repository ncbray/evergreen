package dubx

import (
	"evergreen/dub"
)

type TextMatch interface {
	IsTextMatch()
}
type RuneFilter struct {
	Min rune
	Max rune
}
type RuneRangeMatch struct {
	Invert  bool
	Filters []*RuneFilter
}

func (node *RuneRangeMatch) IsTextMatch() {
}

type StringLiteralMatch struct {
	Value string
}

func (node *StringLiteralMatch) IsTextMatch() {
}

type MatchSequence struct {
	Matches []TextMatch
}

func (node *MatchSequence) IsTextMatch() {
}

type MatchChoice struct {
	Matches []TextMatch
}

func (node *MatchChoice) IsTextMatch() {
}

type MatchRepeat struct {
	Match TextMatch
	Min   int
}

func (node *MatchRepeat) IsTextMatch() {
}

type ASTExpr interface {
	IsASTExpr()
}
type RuneLiteral struct {
	Text  string
	Value rune
}

func (node *RuneLiteral) IsASTExpr() {
}

type StringLiteral struct {
	Text  string
	Value string
}

func (node *StringLiteral) IsASTExpr() {
}

type IntLiteral struct {
	Text  string
	Value int
}

func (node *IntLiteral) IsASTExpr() {
}

type BoolLiteral struct {
	Text  string
	Value bool
}

func (node *BoolLiteral) IsASTExpr() {
}

type StringMatch struct {
	Match TextMatch
}

func (node *StringMatch) IsASTExpr() {
}

type RuneMatch struct {
	Match *RuneRangeMatch
}

func (node *RuneMatch) IsASTExpr() {
}

type ASTType interface {
	IsASTType()
}
type Fake struct {
}

func (node *Fake) IsASTType() {
}

type ASTTypeRef interface {
	IsASTTypeRef()
}
type TypeRef struct {
	Name string
	T    ASTType
}

func (node *TypeRef) IsASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    ASTType
}

func (node *ListTypeRef) IsASTTypeRef() {
}

type Destructure interface {
	IsDestructure()
}
type DestructureValue struct {
	Expr ASTExpr
}

func (node *DestructureValue) IsDestructure() {
}

type DestructureField struct {
	Name        string
	Destructure Destructure
}
type DestructureStruct struct {
	Type *TypeRef
	Args []*DestructureField
}

func (node *DestructureStruct) IsDestructure() {
}

type DestructureList struct {
	Type *ListTypeRef
	Args []Destructure
}

func (node *DestructureList) IsDestructure() {
}

type If struct {
	Expr  ASTExpr
	Block []ASTExpr
}

func (node *If) IsASTExpr() {
}

type Repeat struct {
	Block []ASTExpr
	Min   int
}

func (node *Repeat) IsASTExpr() {
}

type Choice struct {
	Blocks [][]ASTExpr
}

func (node *Choice) IsASTExpr() {
}

type Optional struct {
	Block []ASTExpr
}

func (node *Optional) IsASTExpr() {
}

type Slice struct {
	Block []ASTExpr
}

func (node *Slice) IsASTExpr() {
}

type Assign struct {
	Expr   ASTExpr
	Name   string
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) IsASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) IsASTExpr() {
}

type NamedExpr struct {
	Name string
	Expr ASTExpr
}
type Construct struct {
	Type *TypeRef
	Args []*NamedExpr
}

func (node *Construct) IsASTExpr() {
}

type ConstructList struct {
	Type *ListTypeRef
	Args []ASTExpr
}

func (node *ConstructList) IsASTExpr() {
}

type Coerce struct {
	Type ASTTypeRef
	Expr ASTExpr
}

func (node *Coerce) IsASTExpr() {
}

type Call struct {
	Name string
	T    ASTType
}

func (node *Call) IsASTExpr() {
}

type Fail struct {
}

func (node *Fail) IsASTExpr() {
}

type Append struct {
	List ASTExpr
	Expr ASTExpr
	T    ASTType
}

func (node *Append) IsASTExpr() {
}

type Return struct {
	Exprs []ASTExpr
}

func (node *Return) IsASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     ASTType
}

func (node *BinaryOp) IsASTExpr() {
}
func S(frame *dub.DubState) {
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
func Ident(frame *dub.DubState) (ret0 string) {
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
func DecodeInt(frame *dub.DubState) (ret0 int) {
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
func EscapedChar(frame *dub.DubState) (ret0 rune) {
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
func DecodeString(frame *dub.DubState) (ret0 string) {
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
func DecodeRune(frame *dub.DubState) (ret0 rune) {
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
func DecodeBool(frame *dub.DubState) (ret0 bool) {
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
func Literal(frame *dub.DubState) (ret0 ASTExpr) {
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
	S(frame)
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
func BinaryOperator(frame *dub.DubState) (ret0 string) {
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
	var r29 string
	var r30 string
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
	goto block50
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
	goto block50
block33:
	frame.Fail()
	goto block34
block34:
	frame.Recover(r17)
	goto block50
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
		goto block58
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
		goto block57
	}
block44:
	frame.Consume()
	goto block45
block45:
	r26 = frame.Peek()
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block58
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
		goto block56
	}
block49:
	frame.Consume()
	goto block50
block50:
	r29 = frame.Slice(r1)
	goto block51
block51:
	r0 = r29
	goto block52
block52:
	S(frame)
	if frame.Flow == 0 {
		goto block53
	} else {
		goto block58
	}
block53:
	r30 = r0
	goto block54
block54:
	ret0 = r30
	goto block55
block55:
	return
block56:
	frame.Fail()
	goto block58
block57:
	frame.Fail()
	goto block58
block58:
	return
}
func StringMatchExpr(frame *dub.DubState) (ret0 *StringMatch) {
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
func RuneMatchExpr(frame *dub.DubState) (ret0 *RuneMatch) {
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
func ParseStructTypeRef(frame *dub.DubState) (ret0 *TypeRef) {
	var r0 string
	var r1 string
	var r2 string
	var r3 *TypeRef
	goto block0
block0:
	goto block1
block1:
	r1 = Ident(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block7
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = r0
	goto block4
block4:
	r3 = &TypeRef{Name: r2}
	goto block5
block5:
	ret0 = r3
	goto block6
block6:
	return
block7:
	return
}
func ParseListTypeRef(frame *dub.DubState) (ret0 *ListTypeRef) {
	var r0 ASTTypeRef
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 ASTTypeRef
	var r8 ASTTypeRef
	var r9 *ListTypeRef
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Peek()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block19
	}
block2:
	r2 = '['
	goto block3
block3:
	r3 = r1 == r2
	goto block4
block4:
	if r3 {
		goto block5
	} else {
		goto block18
	}
block5:
	frame.Consume()
	goto block6
block6:
	r4 = frame.Peek()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block19
	}
block7:
	r5 = ']'
	goto block8
block8:
	r6 = r4 == r5
	goto block9
block9:
	if r6 {
		goto block10
	} else {
		goto block17
	}
block10:
	frame.Consume()
	goto block11
block11:
	r7 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block19
	}
block12:
	r0 = r7
	goto block13
block13:
	r8 = r0
	goto block14
block14:
	r9 = &ListTypeRef{Type: r8}
	goto block15
block15:
	ret0 = r9
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
func ParseTypeRef(frame *dub.DubState) (ret0 ASTTypeRef) {
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
func ParseDestructure(frame *dub.DubState) (ret0 Destructure) {
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
func ParseRuneFilterRune(frame *dub.DubState) (ret0 rune) {
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
func ParseRuneFilter(frame *dub.DubState) (ret0 *RuneFilter) {
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
func MatchRune(frame *dub.DubState) (ret0 *RuneRangeMatch) {
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
func Atom(frame *dub.DubState) (ret0 TextMatch) {
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
	goto block26
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
	goto block26
block11:
	frame.Recover(r2)
	goto block12
block12:
	r7 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block29
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
		goto block28
	}
block16:
	frame.Consume()
	goto block17
block17:
	r10 = ParseMatchChoice(frame)
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block29
	}
block18:
	r1 = r10
	goto block19
block19:
	r11 = frame.Peek()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block29
	}
block20:
	r12 = ')'
	goto block21
block21:
	r13 = r11 == r12
	goto block22
block22:
	if r13 {
		goto block23
	} else {
		goto block27
	}
block23:
	frame.Consume()
	goto block24
block24:
	r14 = r1
	goto block25
block25:
	ret0 = r14
	goto block26
block26:
	return
block27:
	frame.Fail()
	goto block29
block28:
	frame.Fail()
	goto block29
block29:
	return
}
func Postfix(frame *dub.DubState) (ret0 TextMatch) {
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
func Sequence(frame *dub.DubState) (ret0 TextMatch) {
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
	r2 = Postfix(frame)
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
	r7 = Postfix(frame)
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
	r11 = Postfix(frame)
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
func ParseMatchChoice(frame *dub.DubState) (ret0 TextMatch) {
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
func ParseExprList(frame *dub.DubState) (ret0 []ASTExpr) {
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
	var r13 []ASTExpr
	var r14 ASTExpr
	var r15 []ASTExpr
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 []ASTExpr
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
	r14 = ParseExpr(frame)
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
func ParseNamedExpr(frame *dub.DubState) (ret0 *NamedExpr) {
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
func ParseNamedExprList(frame *dub.DubState) (ret0 []*NamedExpr) {
	var r0 []*NamedExpr
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 []*NamedExpr
	var r5 int
	var r6 []*NamedExpr
	var r7 *NamedExpr
	var r8 []*NamedExpr
	var r9 int
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 []*NamedExpr
	var r14 *NamedExpr
	var r15 []*NamedExpr
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 []*NamedExpr
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
	r4 = []*NamedExpr{}
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
	r7 = ParseNamedExpr(frame)
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
	r14 = ParseNamedExpr(frame)
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
func ParseTypeList(frame *dub.DubState) (ret0 []ASTTypeRef) {
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
func ParseExpr(frame *dub.DubState) (ret0 ASTExpr) {
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
	var r10 bool
	var r11 string
	var r12 ASTExpr
	var r13 *TypeRef
	var r14 []*NamedExpr
	var r15 *ListTypeRef
	var r16 []ASTExpr
	var r17 string
	var r18 ASTTypeRef
	var r19 ASTExpr
	var r20 string
	var r21 ASTExpr
	var r22 string
	var r23 ASTExpr
	var r24 ASTExpr
	var r25 int
	var r26 ASTExpr
	var r27 int
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
	var r43 rune
	var r44 rune
	var r45 bool
	var r46 rune
	var r47 rune
	var r48 bool
	var r49 rune
	var r50 rune
	var r51 bool
	var r52 int
	var r53 []ASTExpr
	var r54 []ASTExpr
	var r55 int
	var r56 *Repeat
	var r57 rune
	var r58 rune
	var r59 bool
	var r60 rune
	var r61 rune
	var r62 bool
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
	var r75 []ASTExpr
	var r76 [][]ASTExpr
	var r77 int
	var r78 rune
	var r79 rune
	var r80 bool
	var r81 rune
	var r82 rune
	var r83 bool
	var r84 [][]ASTExpr
	var r85 []ASTExpr
	var r86 [][]ASTExpr
	var r87 [][]ASTExpr
	var r88 *Choice
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
	var r104 rune
	var r105 rune
	var r106 bool
	var r107 rune
	var r108 rune
	var r109 bool
	var r110 rune
	var r111 rune
	var r112 bool
	var r113 []ASTExpr
	var r114 []ASTExpr
	var r115 *Optional
	var r116 rune
	var r117 rune
	var r118 bool
	var r119 rune
	var r120 rune
	var r121 bool
	var r122 rune
	var r123 rune
	var r124 bool
	var r125 rune
	var r126 rune
	var r127 bool
	var r128 rune
	var r129 rune
	var r130 bool
	var r131 []ASTExpr
	var r132 []ASTExpr
	var r133 *Slice
	var r134 rune
	var r135 rune
	var r136 bool
	var r137 rune
	var r138 rune
	var r139 bool
	var r140 ASTExpr
	var r141 []ASTExpr
	var r142 ASTExpr
	var r143 []ASTExpr
	var r144 *If
	var r145 rune
	var r146 rune
	var r147 bool
	var r148 rune
	var r149 rune
	var r150 bool
	var r151 rune
	var r152 rune
	var r153 bool
	var r154 string
	var r155 ASTTypeRef
	var r156 int
	var r157 rune
	var r158 rune
	var r159 bool
	var r160 ASTExpr
	var r161 ASTExpr
	var r162 string
	var r163 ASTTypeRef
	var r164 bool
	var r165 *Assign
	var r166 int
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
	var r179 rune
	var r180 rune
	var r181 bool
	var r182 rune
	var r183 rune
	var r184 bool
	var r185 rune
	var r186 rune
	var r187 bool
	var r188 rune
	var r189 rune
	var r190 bool
	var r191 rune
	var r192 rune
	var r193 bool
	var r194 rune
	var r195 rune
	var r196 bool
	var r197 rune
	var r198 rune
	var r199 bool
	var r200 rune
	var r201 rune
	var r202 bool
	var r203 bool
	var r204 string
	var r205 ASTExpr
	var r206 ASTExpr
	var r207 string
	var r208 bool
	var r209 *Assign
	var r210 rune
	var r211 rune
	var r212 bool
	var r213 rune
	var r214 rune
	var r215 bool
	var r216 rune
	var r217 rune
	var r218 bool
	var r219 rune
	var r220 rune
	var r221 bool
	var r222 *TypeRef
	var r223 []*NamedExpr
	var r224 *TypeRef
	var r225 []*NamedExpr
	var r226 *Construct
	var r227 rune
	var r228 rune
	var r229 bool
	var r230 rune
	var r231 rune
	var r232 bool
	var r233 rune
	var r234 rune
	var r235 bool
	var r236 rune
	var r237 rune
	var r238 bool
	var r239 *ListTypeRef
	var r240 []ASTExpr
	var r241 *ListTypeRef
	var r242 []ASTExpr
	var r243 *ConstructList
	var r244 rune
	var r245 rune
	var r246 bool
	var r247 rune
	var r248 rune
	var r249 bool
	var r250 rune
	var r251 rune
	var r252 bool
	var r253 rune
	var r254 rune
	var r255 bool
	var r256 string
	var r257 string
	var r258 *Call
	var r259 rune
	var r260 rune
	var r261 bool
	var r262 rune
	var r263 rune
	var r264 bool
	var r265 rune
	var r266 rune
	var r267 bool
	var r268 rune
	var r269 rune
	var r270 bool
	var r271 *Fail
	var r272 rune
	var r273 rune
	var r274 bool
	var r275 rune
	var r276 rune
	var r277 bool
	var r278 rune
	var r279 rune
	var r280 bool
	var r281 rune
	var r282 rune
	var r283 bool
	var r284 rune
	var r285 rune
	var r286 bool
	var r287 rune
	var r288 rune
	var r289 bool
	var r290 ASTTypeRef
	var r291 ASTExpr
	var r292 ASTTypeRef
	var r293 ASTExpr
	var r294 *Coerce
	var r295 rune
	var r296 rune
	var r297 bool
	var r298 rune
	var r299 rune
	var r300 bool
	var r301 rune
	var r302 rune
	var r303 bool
	var r304 rune
	var r305 rune
	var r306 bool
	var r307 rune
	var r308 rune
	var r309 bool
	var r310 rune
	var r311 rune
	var r312 bool
	var r313 string
	var r314 ASTExpr
	var r315 string
	var r316 *GetName
	var r317 ASTExpr
	var r318 *Append
	var r319 string
	var r320 *Assign
	var r321 rune
	var r322 rune
	var r323 bool
	var r324 rune
	var r325 rune
	var r326 bool
	var r327 rune
	var r328 rune
	var r329 bool
	var r330 rune
	var r331 rune
	var r332 bool
	var r333 rune
	var r334 rune
	var r335 bool
	var r336 rune
	var r337 rune
	var r338 bool
	var r339 int
	var r340 []ASTExpr
	var r341 *Return
	var r342 ASTExpr
	var r343 []ASTExpr
	var r344 *Return
	var r345 []ASTExpr
	var r346 *Return
	var r347 string
	var r348 ASTExpr
	var r349 ASTExpr
	var r350 ASTExpr
	var r351 string
	var r352 ASTExpr
	var r353 *BinaryOp
	var r354 *StringMatch
	var r355 *RuneMatch
	var r356 string
	var r357 *GetName
	goto block0
block0:
	goto block1
block1:
	r25 = frame.Checkpoint()
	goto block2
block2:
	r26 = Literal(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r26
	goto block667
block4:
	frame.Recover(r25)
	goto block5
block5:
	r0 = 0
	goto block6
block6:
	r27 = frame.Checkpoint()
	goto block7
block7:
	r28 = frame.Peek()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block31
	}
block8:
	r29 = 's'
	goto block9
block9:
	r30 = r28 == r29
	goto block10
block10:
	if r30 {
		goto block11
	} else {
		goto block30
	}
block11:
	frame.Consume()
	goto block12
block12:
	r31 = frame.Peek()
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block31
	}
block13:
	r32 = 't'
	goto block14
block14:
	r33 = r31 == r32
	goto block15
block15:
	if r33 {
		goto block16
	} else {
		goto block29
	}
block16:
	frame.Consume()
	goto block17
block17:
	r34 = frame.Peek()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block31
	}
block18:
	r35 = 'a'
	goto block19
block19:
	r36 = r34 == r35
	goto block20
block20:
	if r36 {
		goto block21
	} else {
		goto block28
	}
block21:
	frame.Consume()
	goto block22
block22:
	r37 = frame.Peek()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block31
	}
block23:
	r38 = 'r'
	goto block24
block24:
	r39 = r37 == r38
	goto block25
block25:
	if r39 {
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
	frame.Recover(r27)
	goto block32
block32:
	r40 = frame.Peek()
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block65
	}
block33:
	r41 = 'p'
	goto block34
block34:
	r42 = r40 == r41
	goto block35
block35:
	if r42 {
		goto block36
	} else {
		goto block64
	}
block36:
	frame.Consume()
	goto block37
block37:
	r43 = frame.Peek()
	if frame.Flow == 0 {
		goto block38
	} else {
		goto block65
	}
block38:
	r44 = 'l'
	goto block39
block39:
	r45 = r43 == r44
	goto block40
block40:
	if r45 {
		goto block41
	} else {
		goto block63
	}
block41:
	frame.Consume()
	goto block42
block42:
	r46 = frame.Peek()
	if frame.Flow == 0 {
		goto block43
	} else {
		goto block65
	}
block43:
	r47 = 'u'
	goto block44
block44:
	r48 = r46 == r47
	goto block45
block45:
	if r48 {
		goto block46
	} else {
		goto block62
	}
block46:
	frame.Consume()
	goto block47
block47:
	r49 = frame.Peek()
	if frame.Flow == 0 {
		goto block48
	} else {
		goto block65
	}
block48:
	r50 = 's'
	goto block49
block49:
	r51 = r49 == r50
	goto block50
block50:
	if r51 {
		goto block51
	} else {
		goto block61
	}
block51:
	frame.Consume()
	goto block52
block52:
	r52 = 1
	goto block53
block53:
	r0 = r52
	goto block54
block54:
	S(frame)
	if frame.Flow == 0 {
		goto block55
	} else {
		goto block65
	}
block55:
	r53 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block56
	} else {
		goto block65
	}
block56:
	r1 = r53
	goto block57
block57:
	r54 = r1
	goto block58
block58:
	r55 = r0
	goto block59
block59:
	r56 = &Repeat{Block: r54, Min: r55}
	goto block60
block60:
	ret0 = r56
	goto block667
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
	frame.Recover(r25)
	goto block66
block66:
	r57 = frame.Peek()
	if frame.Flow == 0 {
		goto block67
	} else {
		goto block128
	}
block67:
	r58 = 'c'
	goto block68
block68:
	r59 = r57 == r58
	goto block69
block69:
	if r59 {
		goto block70
	} else {
		goto block127
	}
block70:
	frame.Consume()
	goto block71
block71:
	r60 = frame.Peek()
	if frame.Flow == 0 {
		goto block72
	} else {
		goto block128
	}
block72:
	r61 = 'h'
	goto block73
block73:
	r62 = r60 == r61
	goto block74
block74:
	if r62 {
		goto block75
	} else {
		goto block126
	}
block75:
	frame.Consume()
	goto block76
block76:
	r63 = frame.Peek()
	if frame.Flow == 0 {
		goto block77
	} else {
		goto block128
	}
block77:
	r64 = 'o'
	goto block78
block78:
	r65 = r63 == r64
	goto block79
block79:
	if r65 {
		goto block80
	} else {
		goto block125
	}
block80:
	frame.Consume()
	goto block81
block81:
	r66 = frame.Peek()
	if frame.Flow == 0 {
		goto block82
	} else {
		goto block128
	}
block82:
	r67 = 'o'
	goto block83
block83:
	r68 = r66 == r67
	goto block84
block84:
	if r68 {
		goto block85
	} else {
		goto block124
	}
block85:
	frame.Consume()
	goto block86
block86:
	r69 = frame.Peek()
	if frame.Flow == 0 {
		goto block87
	} else {
		goto block128
	}
block87:
	r70 = 's'
	goto block88
block88:
	r71 = r69 == r70
	goto block89
block89:
	if r71 {
		goto block90
	} else {
		goto block123
	}
block90:
	frame.Consume()
	goto block91
block91:
	r72 = frame.Peek()
	if frame.Flow == 0 {
		goto block92
	} else {
		goto block128
	}
block92:
	r73 = 'e'
	goto block93
block93:
	r74 = r72 == r73
	goto block94
block94:
	if r74 {
		goto block95
	} else {
		goto block122
	}
block95:
	frame.Consume()
	goto block96
block96:
	S(frame)
	if frame.Flow == 0 {
		goto block97
	} else {
		goto block128
	}
block97:
	r75 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block98
	} else {
		goto block128
	}
block98:
	r76 = [][]ASTExpr{r75}
	goto block99
block99:
	r2 = r76
	goto block100
block100:
	r77 = frame.Checkpoint()
	goto block101
block101:
	r78 = frame.Peek()
	if frame.Flow == 0 {
		goto block102
	} else {
		goto block118
	}
block102:
	r79 = 'o'
	goto block103
block103:
	r80 = r78 == r79
	goto block104
block104:
	if r80 {
		goto block105
	} else {
		goto block117
	}
block105:
	frame.Consume()
	goto block106
block106:
	r81 = frame.Peek()
	if frame.Flow == 0 {
		goto block107
	} else {
		goto block118
	}
block107:
	r82 = 'r'
	goto block108
block108:
	r83 = r81 == r82
	goto block109
block109:
	if r83 {
		goto block110
	} else {
		goto block116
	}
block110:
	frame.Consume()
	goto block111
block111:
	S(frame)
	if frame.Flow == 0 {
		goto block112
	} else {
		goto block118
	}
block112:
	r84 = r2
	goto block113
block113:
	r85 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block114
	} else {
		goto block118
	}
block114:
	r86 = append(r84, r85)
	goto block115
block115:
	r2 = r86
	goto block100
block116:
	frame.Fail()
	goto block118
block117:
	frame.Fail()
	goto block118
block118:
	frame.Recover(r77)
	goto block119
block119:
	r87 = r2
	goto block120
block120:
	r88 = &Choice{Blocks: r87}
	goto block121
block121:
	ret0 = r88
	goto block667
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
	frame.Recover(r25)
	goto block129
block129:
	r89 = frame.Peek()
	if frame.Flow == 0 {
		goto block130
	} else {
		goto block183
	}
block130:
	r90 = 'q'
	goto block131
block131:
	r91 = r89 == r90
	goto block132
block132:
	if r91 {
		goto block133
	} else {
		goto block182
	}
block133:
	frame.Consume()
	goto block134
block134:
	r92 = frame.Peek()
	if frame.Flow == 0 {
		goto block135
	} else {
		goto block183
	}
block135:
	r93 = 'u'
	goto block136
block136:
	r94 = r92 == r93
	goto block137
block137:
	if r94 {
		goto block138
	} else {
		goto block181
	}
block138:
	frame.Consume()
	goto block139
block139:
	r95 = frame.Peek()
	if frame.Flow == 0 {
		goto block140
	} else {
		goto block183
	}
block140:
	r96 = 'e'
	goto block141
block141:
	r97 = r95 == r96
	goto block142
block142:
	if r97 {
		goto block143
	} else {
		goto block180
	}
block143:
	frame.Consume()
	goto block144
block144:
	r98 = frame.Peek()
	if frame.Flow == 0 {
		goto block145
	} else {
		goto block183
	}
block145:
	r99 = 's'
	goto block146
block146:
	r100 = r98 == r99
	goto block147
block147:
	if r100 {
		goto block148
	} else {
		goto block179
	}
block148:
	frame.Consume()
	goto block149
block149:
	r101 = frame.Peek()
	if frame.Flow == 0 {
		goto block150
	} else {
		goto block183
	}
block150:
	r102 = 't'
	goto block151
block151:
	r103 = r101 == r102
	goto block152
block152:
	if r103 {
		goto block153
	} else {
		goto block178
	}
block153:
	frame.Consume()
	goto block154
block154:
	r104 = frame.Peek()
	if frame.Flow == 0 {
		goto block155
	} else {
		goto block183
	}
block155:
	r105 = 'i'
	goto block156
block156:
	r106 = r104 == r105
	goto block157
block157:
	if r106 {
		goto block158
	} else {
		goto block177
	}
block158:
	frame.Consume()
	goto block159
block159:
	r107 = frame.Peek()
	if frame.Flow == 0 {
		goto block160
	} else {
		goto block183
	}
block160:
	r108 = 'o'
	goto block161
block161:
	r109 = r107 == r108
	goto block162
block162:
	if r109 {
		goto block163
	} else {
		goto block176
	}
block163:
	frame.Consume()
	goto block164
block164:
	r110 = frame.Peek()
	if frame.Flow == 0 {
		goto block165
	} else {
		goto block183
	}
block165:
	r111 = 'n'
	goto block166
block166:
	r112 = r110 == r111
	goto block167
block167:
	if r112 {
		goto block168
	} else {
		goto block175
	}
block168:
	frame.Consume()
	goto block169
block169:
	S(frame)
	if frame.Flow == 0 {
		goto block170
	} else {
		goto block183
	}
block170:
	r113 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block171
	} else {
		goto block183
	}
block171:
	r3 = r113
	goto block172
block172:
	r114 = r3
	goto block173
block173:
	r115 = &Optional{Block: r114}
	goto block174
block174:
	ret0 = r115
	goto block667
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
	frame.Recover(r25)
	goto block184
block184:
	r116 = frame.Peek()
	if frame.Flow == 0 {
		goto block185
	} else {
		goto block220
	}
block185:
	r117 = 's'
	goto block186
block186:
	r118 = r116 == r117
	goto block187
block187:
	if r118 {
		goto block188
	} else {
		goto block219
	}
block188:
	frame.Consume()
	goto block189
block189:
	r119 = frame.Peek()
	if frame.Flow == 0 {
		goto block190
	} else {
		goto block220
	}
block190:
	r120 = 'l'
	goto block191
block191:
	r121 = r119 == r120
	goto block192
block192:
	if r121 {
		goto block193
	} else {
		goto block218
	}
block193:
	frame.Consume()
	goto block194
block194:
	r122 = frame.Peek()
	if frame.Flow == 0 {
		goto block195
	} else {
		goto block220
	}
block195:
	r123 = 'i'
	goto block196
block196:
	r124 = r122 == r123
	goto block197
block197:
	if r124 {
		goto block198
	} else {
		goto block217
	}
block198:
	frame.Consume()
	goto block199
block199:
	r125 = frame.Peek()
	if frame.Flow == 0 {
		goto block200
	} else {
		goto block220
	}
block200:
	r126 = 'c'
	goto block201
block201:
	r127 = r125 == r126
	goto block202
block202:
	if r127 {
		goto block203
	} else {
		goto block216
	}
block203:
	frame.Consume()
	goto block204
block204:
	r128 = frame.Peek()
	if frame.Flow == 0 {
		goto block205
	} else {
		goto block220
	}
block205:
	r129 = 'e'
	goto block206
block206:
	r130 = r128 == r129
	goto block207
block207:
	if r130 {
		goto block208
	} else {
		goto block215
	}
block208:
	frame.Consume()
	goto block209
block209:
	S(frame)
	if frame.Flow == 0 {
		goto block210
	} else {
		goto block220
	}
block210:
	r131 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block211
	} else {
		goto block220
	}
block211:
	r4 = r131
	goto block212
block212:
	r132 = r4
	goto block213
block213:
	r133 = &Slice{Block: r132}
	goto block214
block214:
	ret0 = r133
	goto block667
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
	frame.Recover(r25)
	goto block221
block221:
	r134 = frame.Peek()
	if frame.Flow == 0 {
		goto block222
	} else {
		goto block242
	}
block222:
	r135 = 'i'
	goto block223
block223:
	r136 = r134 == r135
	goto block224
block224:
	if r136 {
		goto block225
	} else {
		goto block241
	}
block225:
	frame.Consume()
	goto block226
block226:
	r137 = frame.Peek()
	if frame.Flow == 0 {
		goto block227
	} else {
		goto block242
	}
block227:
	r138 = 'f'
	goto block228
block228:
	r139 = r137 == r138
	goto block229
block229:
	if r139 {
		goto block230
	} else {
		goto block240
	}
block230:
	frame.Consume()
	goto block231
block231:
	S(frame)
	if frame.Flow == 0 {
		goto block232
	} else {
		goto block242
	}
block232:
	r140 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block233
	} else {
		goto block242
	}
block233:
	r5 = r140
	goto block234
block234:
	r141 = ParseCodeBlock(frame)
	if frame.Flow == 0 {
		goto block235
	} else {
		goto block242
	}
block235:
	r6 = r141
	goto block236
block236:
	r142 = r5
	goto block237
block237:
	r143 = r6
	goto block238
block238:
	r144 = &If{Expr: r142, Block: r143}
	goto block239
block239:
	ret0 = r144
	goto block667
block240:
	frame.Fail()
	goto block242
block241:
	frame.Fail()
	goto block242
block242:
	frame.Recover(r25)
	goto block243
block243:
	r145 = frame.Peek()
	if frame.Flow == 0 {
		goto block244
	} else {
		goto block284
	}
block244:
	r146 = 'v'
	goto block245
block245:
	r147 = r145 == r146
	goto block246
block246:
	if r147 {
		goto block247
	} else {
		goto block283
	}
block247:
	frame.Consume()
	goto block248
block248:
	r148 = frame.Peek()
	if frame.Flow == 0 {
		goto block249
	} else {
		goto block284
	}
block249:
	r149 = 'a'
	goto block250
block250:
	r150 = r148 == r149
	goto block251
block251:
	if r150 {
		goto block252
	} else {
		goto block282
	}
block252:
	frame.Consume()
	goto block253
block253:
	r151 = frame.Peek()
	if frame.Flow == 0 {
		goto block254
	} else {
		goto block284
	}
block254:
	r152 = 'r'
	goto block255
block255:
	r153 = r151 == r152
	goto block256
block256:
	if r153 {
		goto block257
	} else {
		goto block281
	}
block257:
	frame.Consume()
	goto block258
block258:
	S(frame)
	if frame.Flow == 0 {
		goto block259
	} else {
		goto block284
	}
block259:
	r154 = Ident(frame)
	if frame.Flow == 0 {
		goto block260
	} else {
		goto block284
	}
block260:
	r7 = r154
	goto block261
block261:
	r155 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block262
	} else {
		goto block284
	}
block262:
	r8 = r155
	goto block263
block263:
	r9 = nil
	goto block264
block264:
	r156 = frame.Checkpoint()
	goto block265
block265:
	r157 = frame.Peek()
	if frame.Flow == 0 {
		goto block266
	} else {
		goto block274
	}
block266:
	r158 = '='
	goto block267
block267:
	r159 = r157 == r158
	goto block268
block268:
	if r159 {
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
	r160 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block272
	} else {
		goto block274
	}
block272:
	r9 = r160
	goto block275
block273:
	frame.Fail()
	goto block274
block274:
	frame.Recover(r156)
	goto block275
block275:
	r161 = r9
	goto block276
block276:
	r162 = r7
	goto block277
block277:
	r163 = r8
	goto block278
block278:
	r164 = true
	goto block279
block279:
	r165 = &Assign{Expr: r161, Name: r162, Type: r163, Define: r164}
	goto block280
block280:
	ret0 = r165
	goto block667
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
	frame.Recover(r25)
	goto block285
block285:
	r10 = false
	goto block286
block286:
	r166 = frame.Checkpoint()
	goto block287
block287:
	r167 = frame.Peek()
	if frame.Flow == 0 {
		goto block288
	} else {
		goto block323
	}
block288:
	r168 = 'a'
	goto block289
block289:
	r169 = r167 == r168
	goto block290
block290:
	if r169 {
		goto block291
	} else {
		goto block322
	}
block291:
	frame.Consume()
	goto block292
block292:
	r170 = frame.Peek()
	if frame.Flow == 0 {
		goto block293
	} else {
		goto block323
	}
block293:
	r171 = 's'
	goto block294
block294:
	r172 = r170 == r171
	goto block295
block295:
	if r172 {
		goto block296
	} else {
		goto block321
	}
block296:
	frame.Consume()
	goto block297
block297:
	r173 = frame.Peek()
	if frame.Flow == 0 {
		goto block298
	} else {
		goto block323
	}
block298:
	r174 = 's'
	goto block299
block299:
	r175 = r173 == r174
	goto block300
block300:
	if r175 {
		goto block301
	} else {
		goto block320
	}
block301:
	frame.Consume()
	goto block302
block302:
	r176 = frame.Peek()
	if frame.Flow == 0 {
		goto block303
	} else {
		goto block323
	}
block303:
	r177 = 'i'
	goto block304
block304:
	r178 = r176 == r177
	goto block305
block305:
	if r178 {
		goto block306
	} else {
		goto block319
	}
block306:
	frame.Consume()
	goto block307
block307:
	r179 = frame.Peek()
	if frame.Flow == 0 {
		goto block308
	} else {
		goto block323
	}
block308:
	r180 = 'g'
	goto block309
block309:
	r181 = r179 == r180
	goto block310
block310:
	if r181 {
		goto block311
	} else {
		goto block318
	}
block311:
	frame.Consume()
	goto block312
block312:
	r182 = frame.Peek()
	if frame.Flow == 0 {
		goto block313
	} else {
		goto block323
	}
block313:
	r183 = 'n'
	goto block314
block314:
	r184 = r182 == r183
	goto block315
block315:
	if r184 {
		goto block316
	} else {
		goto block317
	}
block316:
	frame.Consume()
	goto block356
block317:
	frame.Fail()
	goto block323
block318:
	frame.Fail()
	goto block323
block319:
	frame.Fail()
	goto block323
block320:
	frame.Fail()
	goto block323
block321:
	frame.Fail()
	goto block323
block322:
	frame.Fail()
	goto block323
block323:
	frame.Recover(r166)
	goto block324
block324:
	r185 = frame.Peek()
	if frame.Flow == 0 {
		goto block325
	} else {
		goto block372
	}
block325:
	r186 = 'd'
	goto block326
block326:
	r187 = r185 == r186
	goto block327
block327:
	if r187 {
		goto block328
	} else {
		goto block371
	}
block328:
	frame.Consume()
	goto block329
block329:
	r188 = frame.Peek()
	if frame.Flow == 0 {
		goto block330
	} else {
		goto block372
	}
block330:
	r189 = 'e'
	goto block331
block331:
	r190 = r188 == r189
	goto block332
block332:
	if r190 {
		goto block333
	} else {
		goto block370
	}
block333:
	frame.Consume()
	goto block334
block334:
	r191 = frame.Peek()
	if frame.Flow == 0 {
		goto block335
	} else {
		goto block372
	}
block335:
	r192 = 'f'
	goto block336
block336:
	r193 = r191 == r192
	goto block337
block337:
	if r193 {
		goto block338
	} else {
		goto block369
	}
block338:
	frame.Consume()
	goto block339
block339:
	r194 = frame.Peek()
	if frame.Flow == 0 {
		goto block340
	} else {
		goto block372
	}
block340:
	r195 = 'i'
	goto block341
block341:
	r196 = r194 == r195
	goto block342
block342:
	if r196 {
		goto block343
	} else {
		goto block368
	}
block343:
	frame.Consume()
	goto block344
block344:
	r197 = frame.Peek()
	if frame.Flow == 0 {
		goto block345
	} else {
		goto block372
	}
block345:
	r198 = 'n'
	goto block346
block346:
	r199 = r197 == r198
	goto block347
block347:
	if r199 {
		goto block348
	} else {
		goto block367
	}
block348:
	frame.Consume()
	goto block349
block349:
	r200 = frame.Peek()
	if frame.Flow == 0 {
		goto block350
	} else {
		goto block372
	}
block350:
	r201 = 'e'
	goto block351
block351:
	r202 = r200 == r201
	goto block352
block352:
	if r202 {
		goto block353
	} else {
		goto block366
	}
block353:
	frame.Consume()
	goto block354
block354:
	r203 = true
	goto block355
block355:
	r10 = r203
	goto block356
block356:
	S(frame)
	if frame.Flow == 0 {
		goto block357
	} else {
		goto block372
	}
block357:
	r204 = Ident(frame)
	if frame.Flow == 0 {
		goto block358
	} else {
		goto block372
	}
block358:
	r11 = r204
	goto block359
block359:
	r205 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block360
	} else {
		goto block372
	}
block360:
	r12 = r205
	goto block361
block361:
	r206 = r12
	goto block362
block362:
	r207 = r11
	goto block363
block363:
	r208 = r10
	goto block364
block364:
	r209 = &Assign{Expr: r206, Name: r207, Define: r208}
	goto block365
block365:
	ret0 = r209
	goto block667
block366:
	frame.Fail()
	goto block372
block367:
	frame.Fail()
	goto block372
block368:
	frame.Fail()
	goto block372
block369:
	frame.Fail()
	goto block372
block370:
	frame.Fail()
	goto block372
block371:
	frame.Fail()
	goto block372
block372:
	frame.Recover(r25)
	goto block373
block373:
	r210 = frame.Peek()
	if frame.Flow == 0 {
		goto block374
	} else {
		goto block406
	}
block374:
	r211 = 'c'
	goto block375
block375:
	r212 = r210 == r211
	goto block376
block376:
	if r212 {
		goto block377
	} else {
		goto block405
	}
block377:
	frame.Consume()
	goto block378
block378:
	r213 = frame.Peek()
	if frame.Flow == 0 {
		goto block379
	} else {
		goto block406
	}
block379:
	r214 = 'o'
	goto block380
block380:
	r215 = r213 == r214
	goto block381
block381:
	if r215 {
		goto block382
	} else {
		goto block404
	}
block382:
	frame.Consume()
	goto block383
block383:
	r216 = frame.Peek()
	if frame.Flow == 0 {
		goto block384
	} else {
		goto block406
	}
block384:
	r217 = 'n'
	goto block385
block385:
	r218 = r216 == r217
	goto block386
block386:
	if r218 {
		goto block387
	} else {
		goto block403
	}
block387:
	frame.Consume()
	goto block388
block388:
	r219 = frame.Peek()
	if frame.Flow == 0 {
		goto block389
	} else {
		goto block406
	}
block389:
	r220 = 's'
	goto block390
block390:
	r221 = r219 == r220
	goto block391
block391:
	if r221 {
		goto block392
	} else {
		goto block402
	}
block392:
	frame.Consume()
	goto block393
block393:
	S(frame)
	if frame.Flow == 0 {
		goto block394
	} else {
		goto block406
	}
block394:
	r222 = ParseStructTypeRef(frame)
	if frame.Flow == 0 {
		goto block395
	} else {
		goto block406
	}
block395:
	r13 = r222
	goto block396
block396:
	r223 = ParseNamedExprList(frame)
	if frame.Flow == 0 {
		goto block397
	} else {
		goto block406
	}
block397:
	r14 = r223
	goto block398
block398:
	r224 = r13
	goto block399
block399:
	r225 = r14
	goto block400
block400:
	r226 = &Construct{Type: r224, Args: r225}
	goto block401
block401:
	ret0 = r226
	goto block667
block402:
	frame.Fail()
	goto block406
block403:
	frame.Fail()
	goto block406
block404:
	frame.Fail()
	goto block406
block405:
	frame.Fail()
	goto block406
block406:
	frame.Recover(r25)
	goto block407
block407:
	r227 = frame.Peek()
	if frame.Flow == 0 {
		goto block408
	} else {
		goto block440
	}
block408:
	r228 = 'c'
	goto block409
block409:
	r229 = r227 == r228
	goto block410
block410:
	if r229 {
		goto block411
	} else {
		goto block439
	}
block411:
	frame.Consume()
	goto block412
block412:
	r230 = frame.Peek()
	if frame.Flow == 0 {
		goto block413
	} else {
		goto block440
	}
block413:
	r231 = 'o'
	goto block414
block414:
	r232 = r230 == r231
	goto block415
block415:
	if r232 {
		goto block416
	} else {
		goto block438
	}
block416:
	frame.Consume()
	goto block417
block417:
	r233 = frame.Peek()
	if frame.Flow == 0 {
		goto block418
	} else {
		goto block440
	}
block418:
	r234 = 'n'
	goto block419
block419:
	r235 = r233 == r234
	goto block420
block420:
	if r235 {
		goto block421
	} else {
		goto block437
	}
block421:
	frame.Consume()
	goto block422
block422:
	r236 = frame.Peek()
	if frame.Flow == 0 {
		goto block423
	} else {
		goto block440
	}
block423:
	r237 = 'l'
	goto block424
block424:
	r238 = r236 == r237
	goto block425
block425:
	if r238 {
		goto block426
	} else {
		goto block436
	}
block426:
	frame.Consume()
	goto block427
block427:
	S(frame)
	if frame.Flow == 0 {
		goto block428
	} else {
		goto block440
	}
block428:
	r239 = ParseListTypeRef(frame)
	if frame.Flow == 0 {
		goto block429
	} else {
		goto block440
	}
block429:
	r15 = r239
	goto block430
block430:
	r240 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block431
	} else {
		goto block440
	}
block431:
	r16 = r240
	goto block432
block432:
	r241 = r15
	goto block433
block433:
	r242 = r16
	goto block434
block434:
	r243 = &ConstructList{Type: r241, Args: r242}
	goto block435
block435:
	ret0 = r243
	goto block667
block436:
	frame.Fail()
	goto block440
block437:
	frame.Fail()
	goto block440
block438:
	frame.Fail()
	goto block440
block439:
	frame.Fail()
	goto block440
block440:
	frame.Recover(r25)
	goto block441
block441:
	r244 = frame.Peek()
	if frame.Flow == 0 {
		goto block442
	} else {
		goto block471
	}
block442:
	r245 = 'c'
	goto block443
block443:
	r246 = r244 == r245
	goto block444
block444:
	if r246 {
		goto block445
	} else {
		goto block470
	}
block445:
	frame.Consume()
	goto block446
block446:
	r247 = frame.Peek()
	if frame.Flow == 0 {
		goto block447
	} else {
		goto block471
	}
block447:
	r248 = 'a'
	goto block448
block448:
	r249 = r247 == r248
	goto block449
block449:
	if r249 {
		goto block450
	} else {
		goto block469
	}
block450:
	frame.Consume()
	goto block451
block451:
	r250 = frame.Peek()
	if frame.Flow == 0 {
		goto block452
	} else {
		goto block471
	}
block452:
	r251 = 'l'
	goto block453
block453:
	r252 = r250 == r251
	goto block454
block454:
	if r252 {
		goto block455
	} else {
		goto block468
	}
block455:
	frame.Consume()
	goto block456
block456:
	r253 = frame.Peek()
	if frame.Flow == 0 {
		goto block457
	} else {
		goto block471
	}
block457:
	r254 = 'l'
	goto block458
block458:
	r255 = r253 == r254
	goto block459
block459:
	if r255 {
		goto block460
	} else {
		goto block467
	}
block460:
	frame.Consume()
	goto block461
block461:
	S(frame)
	if frame.Flow == 0 {
		goto block462
	} else {
		goto block471
	}
block462:
	r256 = Ident(frame)
	if frame.Flow == 0 {
		goto block463
	} else {
		goto block471
	}
block463:
	r17 = r256
	goto block464
block464:
	r257 = r17
	goto block465
block465:
	r258 = &Call{Name: r257}
	goto block466
block466:
	ret0 = r258
	goto block667
block467:
	frame.Fail()
	goto block471
block468:
	frame.Fail()
	goto block471
block469:
	frame.Fail()
	goto block471
block470:
	frame.Fail()
	goto block471
block471:
	frame.Recover(r25)
	goto block472
block472:
	r259 = frame.Peek()
	if frame.Flow == 0 {
		goto block473
	} else {
		goto block499
	}
block473:
	r260 = 'f'
	goto block474
block474:
	r261 = r259 == r260
	goto block475
block475:
	if r261 {
		goto block476
	} else {
		goto block498
	}
block476:
	frame.Consume()
	goto block477
block477:
	r262 = frame.Peek()
	if frame.Flow == 0 {
		goto block478
	} else {
		goto block499
	}
block478:
	r263 = 'a'
	goto block479
block479:
	r264 = r262 == r263
	goto block480
block480:
	if r264 {
		goto block481
	} else {
		goto block497
	}
block481:
	frame.Consume()
	goto block482
block482:
	r265 = frame.Peek()
	if frame.Flow == 0 {
		goto block483
	} else {
		goto block499
	}
block483:
	r266 = 'i'
	goto block484
block484:
	r267 = r265 == r266
	goto block485
block485:
	if r267 {
		goto block486
	} else {
		goto block496
	}
block486:
	frame.Consume()
	goto block487
block487:
	r268 = frame.Peek()
	if frame.Flow == 0 {
		goto block488
	} else {
		goto block499
	}
block488:
	r269 = 'l'
	goto block489
block489:
	r270 = r268 == r269
	goto block490
block490:
	if r270 {
		goto block491
	} else {
		goto block495
	}
block491:
	frame.Consume()
	goto block492
block492:
	S(frame)
	if frame.Flow == 0 {
		goto block493
	} else {
		goto block499
	}
block493:
	r271 = &Fail{}
	goto block494
block494:
	ret0 = r271
	goto block667
block495:
	frame.Fail()
	goto block499
block496:
	frame.Fail()
	goto block499
block497:
	frame.Fail()
	goto block499
block498:
	frame.Fail()
	goto block499
block499:
	frame.Recover(r25)
	goto block500
block500:
	r272 = frame.Peek()
	if frame.Flow == 0 {
		goto block501
	} else {
		goto block545
	}
block501:
	r273 = 'c'
	goto block502
block502:
	r274 = r272 == r273
	goto block503
block503:
	if r274 {
		goto block504
	} else {
		goto block544
	}
block504:
	frame.Consume()
	goto block505
block505:
	r275 = frame.Peek()
	if frame.Flow == 0 {
		goto block506
	} else {
		goto block545
	}
block506:
	r276 = 'o'
	goto block507
block507:
	r277 = r275 == r276
	goto block508
block508:
	if r277 {
		goto block509
	} else {
		goto block543
	}
block509:
	frame.Consume()
	goto block510
block510:
	r278 = frame.Peek()
	if frame.Flow == 0 {
		goto block511
	} else {
		goto block545
	}
block511:
	r279 = 'e'
	goto block512
block512:
	r280 = r278 == r279
	goto block513
block513:
	if r280 {
		goto block514
	} else {
		goto block542
	}
block514:
	frame.Consume()
	goto block515
block515:
	r281 = frame.Peek()
	if frame.Flow == 0 {
		goto block516
	} else {
		goto block545
	}
block516:
	r282 = 'r'
	goto block517
block517:
	r283 = r281 == r282
	goto block518
block518:
	if r283 {
		goto block519
	} else {
		goto block541
	}
block519:
	frame.Consume()
	goto block520
block520:
	r284 = frame.Peek()
	if frame.Flow == 0 {
		goto block521
	} else {
		goto block545
	}
block521:
	r285 = 'c'
	goto block522
block522:
	r286 = r284 == r285
	goto block523
block523:
	if r286 {
		goto block524
	} else {
		goto block540
	}
block524:
	frame.Consume()
	goto block525
block525:
	r287 = frame.Peek()
	if frame.Flow == 0 {
		goto block526
	} else {
		goto block545
	}
block526:
	r288 = 'e'
	goto block527
block527:
	r289 = r287 == r288
	goto block528
block528:
	if r289 {
		goto block529
	} else {
		goto block539
	}
block529:
	frame.Consume()
	goto block530
block530:
	S(frame)
	if frame.Flow == 0 {
		goto block531
	} else {
		goto block545
	}
block531:
	r290 = ParseTypeRef(frame)
	if frame.Flow == 0 {
		goto block532
	} else {
		goto block545
	}
block532:
	r18 = r290
	goto block533
block533:
	r291 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block534
	} else {
		goto block545
	}
block534:
	r19 = r291
	goto block535
block535:
	r292 = r18
	goto block536
block536:
	r293 = r19
	goto block537
block537:
	r294 = &Coerce{Type: r292, Expr: r293}
	goto block538
block538:
	ret0 = r294
	goto block667
block539:
	frame.Fail()
	goto block545
block540:
	frame.Fail()
	goto block545
block541:
	frame.Fail()
	goto block545
block542:
	frame.Fail()
	goto block545
block543:
	frame.Fail()
	goto block545
block544:
	frame.Fail()
	goto block545
block545:
	frame.Recover(r25)
	goto block546
block546:
	r295 = frame.Peek()
	if frame.Flow == 0 {
		goto block547
	} else {
		goto block594
	}
block547:
	r296 = 'a'
	goto block548
block548:
	r297 = r295 == r296
	goto block549
block549:
	if r297 {
		goto block550
	} else {
		goto block593
	}
block550:
	frame.Consume()
	goto block551
block551:
	r298 = frame.Peek()
	if frame.Flow == 0 {
		goto block552
	} else {
		goto block594
	}
block552:
	r299 = 'p'
	goto block553
block553:
	r300 = r298 == r299
	goto block554
block554:
	if r300 {
		goto block555
	} else {
		goto block592
	}
block555:
	frame.Consume()
	goto block556
block556:
	r301 = frame.Peek()
	if frame.Flow == 0 {
		goto block557
	} else {
		goto block594
	}
block557:
	r302 = 'p'
	goto block558
block558:
	r303 = r301 == r302
	goto block559
block559:
	if r303 {
		goto block560
	} else {
		goto block591
	}
block560:
	frame.Consume()
	goto block561
block561:
	r304 = frame.Peek()
	if frame.Flow == 0 {
		goto block562
	} else {
		goto block594
	}
block562:
	r305 = 'e'
	goto block563
block563:
	r306 = r304 == r305
	goto block564
block564:
	if r306 {
		goto block565
	} else {
		goto block590
	}
block565:
	frame.Consume()
	goto block566
block566:
	r307 = frame.Peek()
	if frame.Flow == 0 {
		goto block567
	} else {
		goto block594
	}
block567:
	r308 = 'n'
	goto block568
block568:
	r309 = r307 == r308
	goto block569
block569:
	if r309 {
		goto block570
	} else {
		goto block589
	}
block570:
	frame.Consume()
	goto block571
block571:
	r310 = frame.Peek()
	if frame.Flow == 0 {
		goto block572
	} else {
		goto block594
	}
block572:
	r311 = 'd'
	goto block573
block573:
	r312 = r310 == r311
	goto block574
block574:
	if r312 {
		goto block575
	} else {
		goto block588
	}
block575:
	frame.Consume()
	goto block576
block576:
	S(frame)
	if frame.Flow == 0 {
		goto block577
	} else {
		goto block594
	}
block577:
	r313 = Ident(frame)
	if frame.Flow == 0 {
		goto block578
	} else {
		goto block594
	}
block578:
	r20 = r313
	goto block579
block579:
	r314 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block580
	} else {
		goto block594
	}
block580:
	r21 = r314
	goto block581
block581:
	r315 = r20
	goto block582
block582:
	r316 = &GetName{Name: r315}
	goto block583
block583:
	r317 = r21
	goto block584
block584:
	r318 = &Append{List: r316, Expr: r317}
	goto block585
block585:
	r319 = r20
	goto block586
block586:
	r320 = &Assign{Expr: r318, Name: r319}
	goto block587
block587:
	ret0 = r320
	goto block667
block588:
	frame.Fail()
	goto block594
block589:
	frame.Fail()
	goto block594
block590:
	frame.Fail()
	goto block594
block591:
	frame.Fail()
	goto block594
block592:
	frame.Fail()
	goto block594
block593:
	frame.Fail()
	goto block594
block594:
	frame.Recover(r25)
	goto block595
block595:
	r321 = frame.Peek()
	if frame.Flow == 0 {
		goto block596
	} else {
		goto block645
	}
block596:
	r322 = 'r'
	goto block597
block597:
	r323 = r321 == r322
	goto block598
block598:
	if r323 {
		goto block599
	} else {
		goto block644
	}
block599:
	frame.Consume()
	goto block600
block600:
	r324 = frame.Peek()
	if frame.Flow == 0 {
		goto block601
	} else {
		goto block645
	}
block601:
	r325 = 'e'
	goto block602
block602:
	r326 = r324 == r325
	goto block603
block603:
	if r326 {
		goto block604
	} else {
		goto block643
	}
block604:
	frame.Consume()
	goto block605
block605:
	r327 = frame.Peek()
	if frame.Flow == 0 {
		goto block606
	} else {
		goto block645
	}
block606:
	r328 = 't'
	goto block607
block607:
	r329 = r327 == r328
	goto block608
block608:
	if r329 {
		goto block609
	} else {
		goto block642
	}
block609:
	frame.Consume()
	goto block610
block610:
	r330 = frame.Peek()
	if frame.Flow == 0 {
		goto block611
	} else {
		goto block645
	}
block611:
	r331 = 'u'
	goto block612
block612:
	r332 = r330 == r331
	goto block613
block613:
	if r332 {
		goto block614
	} else {
		goto block641
	}
block614:
	frame.Consume()
	goto block615
block615:
	r333 = frame.Peek()
	if frame.Flow == 0 {
		goto block616
	} else {
		goto block645
	}
block616:
	r334 = 'r'
	goto block617
block617:
	r335 = r333 == r334
	goto block618
block618:
	if r335 {
		goto block619
	} else {
		goto block640
	}
block619:
	frame.Consume()
	goto block620
block620:
	r336 = frame.Peek()
	if frame.Flow == 0 {
		goto block621
	} else {
		goto block645
	}
block621:
	r337 = 'n'
	goto block622
block622:
	r338 = r336 == r337
	goto block623
block623:
	if r338 {
		goto block624
	} else {
		goto block639
	}
block624:
	frame.Consume()
	goto block625
block625:
	S(frame)
	if frame.Flow == 0 {
		goto block626
	} else {
		goto block645
	}
block626:
	r339 = frame.Checkpoint()
	goto block627
block627:
	r340 = ParseExprList(frame)
	if frame.Flow == 0 {
		goto block628
	} else {
		goto block630
	}
block628:
	r341 = &Return{Exprs: r340}
	goto block629
block629:
	ret0 = r341
	goto block667
block630:
	frame.Recover(r339)
	goto block631
block631:
	r342 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block632
	} else {
		goto block635
	}
block632:
	r343 = []ASTExpr{r342}
	goto block633
block633:
	r344 = &Return{Exprs: r343}
	goto block634
block634:
	ret0 = r344
	goto block667
block635:
	frame.Recover(r339)
	goto block636
block636:
	r345 = []ASTExpr{}
	goto block637
block637:
	r346 = &Return{Exprs: r345}
	goto block638
block638:
	ret0 = r346
	goto block667
block639:
	frame.Fail()
	goto block645
block640:
	frame.Fail()
	goto block645
block641:
	frame.Fail()
	goto block645
block642:
	frame.Fail()
	goto block645
block643:
	frame.Fail()
	goto block645
block644:
	frame.Fail()
	goto block645
block645:
	frame.Recover(r25)
	goto block646
block646:
	r347 = BinaryOperator(frame)
	if frame.Flow == 0 {
		goto block647
	} else {
		goto block657
	}
block647:
	r22 = r347
	goto block648
block648:
	r348 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block649
	} else {
		goto block657
	}
block649:
	r23 = r348
	goto block650
block650:
	r349 = ParseExpr(frame)
	if frame.Flow == 0 {
		goto block651
	} else {
		goto block657
	}
block651:
	r24 = r349
	goto block652
block652:
	r350 = r23
	goto block653
block653:
	r351 = r22
	goto block654
block654:
	r352 = r24
	goto block655
block655:
	r353 = &BinaryOp{Left: r350, Op: r351, Right: r352}
	goto block656
block656:
	ret0 = r353
	goto block667
block657:
	frame.Recover(r25)
	goto block658
block658:
	r354 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		goto block659
	} else {
		goto block660
	}
block659:
	ret0 = r354
	goto block667
block660:
	frame.Recover(r25)
	goto block661
block661:
	r355 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		goto block662
	} else {
		goto block663
	}
block662:
	ret0 = r355
	goto block667
block663:
	frame.Recover(r25)
	goto block664
block664:
	r356 = Ident(frame)
	if frame.Flow == 0 {
		goto block665
	} else {
		goto block668
	}
block665:
	r357 = &GetName{Name: r356}
	goto block666
block666:
	ret0 = r357
	goto block667
block667:
	return
block668:
	return
}
func ParseCodeBlock(frame *dub.DubState) (ret0 []ASTExpr) {
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
