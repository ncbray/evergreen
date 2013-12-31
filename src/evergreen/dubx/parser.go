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
	r4 = Choice(frame)
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
	r10 = Choice(frame)
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
func Choice(frame *dub.DubState) (ret0 TextMatch) {
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
func ParseExpr(frame *dub.DubState) (ret0 ASTExpr) {
	var r0 string
	var r1 int
	var r2 ASTExpr
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
	var r16 string
	var r17 *Call
	var r18 rune
	var r19 rune
	var r20 bool
	var r21 rune
	var r22 rune
	var r23 bool
	var r24 rune
	var r25 rune
	var r26 bool
	var r27 rune
	var r28 rune
	var r29 bool
	var r30 *Fail
	var r31 *StringMatch
	var r32 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Checkpoint()
	goto block2
block2:
	r2 = Literal(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r2
	goto block69
block4:
	frame.Recover(r1)
	goto block5
block5:
	r3 = frame.Peek()
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block35
	}
block6:
	r4 = 'c'
	goto block7
block7:
	r5 = r3 == r4
	goto block8
block8:
	if r5 {
		goto block9
	} else {
		goto block34
	}
block9:
	frame.Consume()
	goto block10
block10:
	r6 = frame.Peek()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block35
	}
block11:
	r7 = 'a'
	goto block12
block12:
	r8 = r6 == r7
	goto block13
block13:
	if r8 {
		goto block14
	} else {
		goto block33
	}
block14:
	frame.Consume()
	goto block15
block15:
	r9 = frame.Peek()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block35
	}
block16:
	r10 = 'l'
	goto block17
block17:
	r11 = r9 == r10
	goto block18
block18:
	if r11 {
		goto block19
	} else {
		goto block32
	}
block19:
	frame.Consume()
	goto block20
block20:
	r12 = frame.Peek()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block35
	}
block21:
	r13 = 'l'
	goto block22
block22:
	r14 = r12 == r13
	goto block23
block23:
	if r14 {
		goto block24
	} else {
		goto block31
	}
block24:
	frame.Consume()
	goto block25
block25:
	S(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block35
	}
block26:
	r15 = Ident(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block35
	}
block27:
	r0 = r15
	goto block28
block28:
	r16 = r0
	goto block29
block29:
	r17 = &Call{Name: r16}
	goto block30
block30:
	ret0 = r17
	goto block69
block31:
	frame.Fail()
	goto block35
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
	frame.Recover(r1)
	goto block36
block36:
	r18 = frame.Peek()
	if frame.Flow == 0 {
		goto block37
	} else {
		goto block63
	}
block37:
	r19 = 'f'
	goto block38
block38:
	r20 = r18 == r19
	goto block39
block39:
	if r20 {
		goto block40
	} else {
		goto block62
	}
block40:
	frame.Consume()
	goto block41
block41:
	r21 = frame.Peek()
	if frame.Flow == 0 {
		goto block42
	} else {
		goto block63
	}
block42:
	r22 = 'a'
	goto block43
block43:
	r23 = r21 == r22
	goto block44
block44:
	if r23 {
		goto block45
	} else {
		goto block61
	}
block45:
	frame.Consume()
	goto block46
block46:
	r24 = frame.Peek()
	if frame.Flow == 0 {
		goto block47
	} else {
		goto block63
	}
block47:
	r25 = 'i'
	goto block48
block48:
	r26 = r24 == r25
	goto block49
block49:
	if r26 {
		goto block50
	} else {
		goto block60
	}
block50:
	frame.Consume()
	goto block51
block51:
	r27 = frame.Peek()
	if frame.Flow == 0 {
		goto block52
	} else {
		goto block63
	}
block52:
	r28 = 'l'
	goto block53
block53:
	r29 = r27 == r28
	goto block54
block54:
	if r29 {
		goto block55
	} else {
		goto block59
	}
block55:
	frame.Consume()
	goto block56
block56:
	S(frame)
	if frame.Flow == 0 {
		goto block57
	} else {
		goto block63
	}
block57:
	r30 = &Fail{}
	goto block58
block58:
	ret0 = r30
	goto block69
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
	frame.Recover(r1)
	goto block64
block64:
	r31 = StringMatchExpr(frame)
	if frame.Flow == 0 {
		goto block65
	} else {
		goto block66
	}
block65:
	ret0 = r31
	goto block69
block66:
	frame.Recover(r1)
	goto block67
block67:
	r32 = RuneMatchExpr(frame)
	if frame.Flow == 0 {
		goto block68
	} else {
		goto block70
	}
block68:
	ret0 = r32
	goto block69
block69:
	return
block70:
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
