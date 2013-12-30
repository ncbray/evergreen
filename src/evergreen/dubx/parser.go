package dubx

import (
	"evergreen/dub"
)

type TextMatch interface {
	isTextMatch()
}
type RuneFilter struct {
	Min rune
	Max rune
}
type RuneMatch struct {
	Invert  bool
	Filters []*RuneFilter
}

func (node *RuneMatch) isTextMatch() {
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

type Token interface {
	isToken()
}
type IntTok struct {
	Text  string
	Value int
}

func (node *IntTok) isToken() {
}

type StrTok struct {
	Text  string
	Value string
}

func (node *StrTok) isToken() {
}

type RuneTok struct {
	Text  string
	Value rune
}

func (node *RuneTok) isToken() {
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
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block16
	}
block3:
	r2 = ' '
	goto block4
block4:
	r3 = r1 == r2
	goto block5
block5:
	if r3 {
		goto block1
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
		goto block1
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
		goto block1
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
		goto block1
	} else {
		goto block15
	}
block15:
	frame.Fail()
	goto block16
block16:
	frame.Recover(r0)
	goto block17
block17:
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
	r2 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block50
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
		goto block49
	}
block18:
	r13 = frame.Checkpoint()
	goto block19
block19:
	r14 = frame.Read()
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block42
	}
block20:
	r15 = 'a'
	goto block21
block21:
	r16 = r14 >= r15
	goto block22
block22:
	if r16 {
		goto block23
	} else {
		goto block26
	}
block23:
	r17 = 'z'
	goto block24
block24:
	r18 = r14 <= r17
	goto block25
block25:
	if r18 {
		goto block18
	} else {
		goto block26
	}
block26:
	r19 = 'A'
	goto block27
block27:
	r20 = r14 >= r19
	goto block28
block28:
	if r20 {
		goto block29
	} else {
		goto block32
	}
block29:
	r21 = 'Z'
	goto block30
block30:
	r22 = r14 <= r21
	goto block31
block31:
	if r22 {
		goto block18
	} else {
		goto block32
	}
block32:
	r23 = '_'
	goto block33
block33:
	r24 = r14 == r23
	goto block34
block34:
	if r24 {
		goto block18
	} else {
		goto block35
	}
block35:
	r25 = '0'
	goto block36
block36:
	r26 = r14 >= r25
	goto block37
block37:
	if r26 {
		goto block38
	} else {
		goto block41
	}
block38:
	r27 = '9'
	goto block39
block39:
	r28 = r14 <= r27
	goto block40
block40:
	if r28 {
		goto block18
	} else {
		goto block41
	}
block41:
	frame.Fail()
	goto block42
block42:
	frame.Recover(r13)
	goto block43
block43:
	r29 = frame.Slice(r1)
	goto block44
block44:
	r0 = r29
	goto block45
block45:
	S(frame)
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block50
	}
block46:
	r30 = r0
	goto block47
block47:
	ret0 = r30
	goto block48
block48:
	return
block49:
	frame.Fail()
	goto block50
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
	r3 = frame.Read()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block17
	}
block4:
	r4 = '+'
	goto block5
block5:
	r5 = r3 == r4
	goto block6
block6:
	if r5 {
		goto block45
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
		goto block45
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
		goto block45
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
		goto block45
	} else {
		goto block16
	}
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r2)
	goto block18
block18:
	r12 = frame.Read()
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block33
	}
block19:
	r13 = '<'
	goto block20
block20:
	r14 = r12 == r13
	goto block21
block21:
	if r14 {
		goto block25
	} else {
		goto block22
	}
block22:
	r15 = '>'
	goto block23
block23:
	r16 = r12 == r15
	goto block24
block24:
	if r16 {
		goto block25
	} else {
		goto block32
	}
block25:
	r17 = frame.Checkpoint()
	goto block26
block26:
	r18 = frame.Read()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block31
	}
block27:
	r19 = '='
	goto block28
block28:
	r20 = r18 == r19
	goto block29
block29:
	if r20 {
		goto block45
	} else {
		goto block30
	}
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(r17)
	goto block45
block32:
	frame.Fail()
	goto block33
block33:
	frame.Recover(r2)
	goto block34
block34:
	r21 = frame.Read()
	if frame.Flow == 0 {
		goto block35
	} else {
		goto block53
	}
block35:
	r22 = '!'
	goto block36
block36:
	r23 = r21 == r22
	goto block37
block37:
	if r23 {
		goto block41
	} else {
		goto block38
	}
block38:
	r24 = '='
	goto block39
block39:
	r25 = r21 == r24
	goto block40
block40:
	if r25 {
		goto block41
	} else {
		goto block52
	}
block41:
	r26 = frame.Read()
	if frame.Flow == 0 {
		goto block42
	} else {
		goto block53
	}
block42:
	r27 = '='
	goto block43
block43:
	r28 = r26 == r27
	goto block44
block44:
	if r28 {
		goto block45
	} else {
		goto block51
	}
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
		goto block53
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
	goto block53
block52:
	frame.Fail()
	goto block53
block53:
	return
}
func Int(frame *dub.DubState) (ret0 *IntTok) {
	var r0 int
	var r1 int
	var r2 string
	var r3 int
	var r4 int
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 bool
	var r10 int
	var r11 rune
	var r12 int
	var r13 int
	var r14 int
	var r15 int
	var r16 int
	var r17 int
	var r18 int
	var r19 int
	var r20 rune
	var r21 rune
	var r22 bool
	var r23 rune
	var r24 bool
	var r25 int
	var r26 rune
	var r27 int
	var r28 int
	var r29 int
	var r30 int
	var r31 int
	var r32 int
	var r33 int
	var r34 string
	var r35 string
	var r36 int
	var r37 *IntTok
	goto block0
block0:
	goto block1
block1:
	r3 = frame.Checkpoint()
	goto block2
block2:
	r4 = 0
	goto block3
block3:
	r0 = r4
	goto block4
block4:
	r5 = frame.Read()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block52
	}
block5:
	r6 = '0'
	goto block6
block6:
	r7 = r5 >= r6
	goto block7
block7:
	if r7 {
		goto block8
	} else {
		goto block51
	}
block8:
	r8 = '9'
	goto block9
block9:
	r9 = r5 <= r8
	goto block10
block10:
	if r9 {
		goto block11
	} else {
		goto block51
	}
block11:
	r10 = int(r5)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block52
	}
block12:
	r11 = '0'
	goto block13
block13:
	r12 = int(r11)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block52
	}
block14:
	r13 = r10 - r12
	goto block15
block15:
	r1 = r13
	goto block16
block16:
	r14 = r0
	goto block17
block17:
	r15 = 10
	goto block18
block18:
	r16 = r14 * r15
	goto block19
block19:
	r17 = r1
	goto block20
block20:
	r18 = r16 + r17
	goto block21
block21:
	r0 = r18
	goto block22
block22:
	r19 = frame.Checkpoint()
	goto block23
block23:
	r20 = frame.Read()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block42
	}
block24:
	r21 = '0'
	goto block25
block25:
	r22 = r20 >= r21
	goto block26
block26:
	if r22 {
		goto block27
	} else {
		goto block41
	}
block27:
	r23 = '9'
	goto block28
block28:
	r24 = r20 <= r23
	goto block29
block29:
	if r24 {
		goto block30
	} else {
		goto block41
	}
block30:
	r25 = int(r20)
	if frame.Flow == 0 {
		goto block31
	} else {
		goto block42
	}
block31:
	r26 = '0'
	goto block32
block32:
	r27 = int(r26)
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block42
	}
block33:
	r28 = r25 - r27
	goto block34
block34:
	r1 = r28
	goto block35
block35:
	r29 = r0
	goto block36
block36:
	r30 = 10
	goto block37
block37:
	r31 = r29 * r30
	goto block38
block38:
	r32 = r1
	goto block39
block39:
	r33 = r31 + r32
	goto block40
block40:
	r0 = r33
	goto block22
block41:
	frame.Fail()
	goto block42
block42:
	frame.Recover(r19)
	goto block43
block43:
	r34 = frame.Slice(r3)
	goto block44
block44:
	r2 = r34
	goto block45
block45:
	S(frame)
	if frame.Flow == 0 {
		goto block46
	} else {
		goto block52
	}
block46:
	r35 = r2
	goto block47
block47:
	r36 = r0
	goto block48
block48:
	r37 = &IntTok{Text: r35, Value: r36}
	goto block49
block49:
	ret0 = r37
	goto block50
block50:
	return
block51:
	frame.Fail()
	goto block52
block52:
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
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block9
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
		goto block8
	}
block6:
	r4 = '\a'
	goto block7
block7:
	ret0 = r4
	goto block80
block8:
	frame.Fail()
	goto block9
block9:
	frame.Recover(r0)
	goto block10
block10:
	r5 = frame.Read()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block17
	}
block11:
	r6 = 'b'
	goto block12
block12:
	r7 = r5 == r6
	goto block13
block13:
	if r7 {
		goto block14
	} else {
		goto block16
	}
block14:
	r8 = '\b'
	goto block15
block15:
	ret0 = r8
	goto block80
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r0)
	goto block18
block18:
	r9 = frame.Read()
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block25
	}
block19:
	r10 = 'f'
	goto block20
block20:
	r11 = r9 == r10
	goto block21
block21:
	if r11 {
		goto block22
	} else {
		goto block24
	}
block22:
	r12 = '\f'
	goto block23
block23:
	ret0 = r12
	goto block80
block24:
	frame.Fail()
	goto block25
block25:
	frame.Recover(r0)
	goto block26
block26:
	r13 = frame.Read()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block33
	}
block27:
	r14 = 'n'
	goto block28
block28:
	r15 = r13 == r14
	goto block29
block29:
	if r15 {
		goto block30
	} else {
		goto block32
	}
block30:
	r16 = '\n'
	goto block31
block31:
	ret0 = r16
	goto block80
block32:
	frame.Fail()
	goto block33
block33:
	frame.Recover(r0)
	goto block34
block34:
	r17 = frame.Read()
	if frame.Flow == 0 {
		goto block35
	} else {
		goto block41
	}
block35:
	r18 = 'r'
	goto block36
block36:
	r19 = r17 == r18
	goto block37
block37:
	if r19 {
		goto block38
	} else {
		goto block40
	}
block38:
	r20 = '\r'
	goto block39
block39:
	ret0 = r20
	goto block80
block40:
	frame.Fail()
	goto block41
block41:
	frame.Recover(r0)
	goto block42
block42:
	r21 = frame.Read()
	if frame.Flow == 0 {
		goto block43
	} else {
		goto block49
	}
block43:
	r22 = 't'
	goto block44
block44:
	r23 = r21 == r22
	goto block45
block45:
	if r23 {
		goto block46
	} else {
		goto block48
	}
block46:
	r24 = '\t'
	goto block47
block47:
	ret0 = r24
	goto block80
block48:
	frame.Fail()
	goto block49
block49:
	frame.Recover(r0)
	goto block50
block50:
	r25 = frame.Read()
	if frame.Flow == 0 {
		goto block51
	} else {
		goto block57
	}
block51:
	r26 = 'v'
	goto block52
block52:
	r27 = r25 == r26
	goto block53
block53:
	if r27 {
		goto block54
	} else {
		goto block56
	}
block54:
	r28 = '\v'
	goto block55
block55:
	ret0 = r28
	goto block80
block56:
	frame.Fail()
	goto block57
block57:
	frame.Recover(r0)
	goto block58
block58:
	r29 = frame.Read()
	if frame.Flow == 0 {
		goto block59
	} else {
		goto block65
	}
block59:
	r30 = '\\'
	goto block60
block60:
	r31 = r29 == r30
	goto block61
block61:
	if r31 {
		goto block62
	} else {
		goto block64
	}
block62:
	r32 = '\\'
	goto block63
block63:
	ret0 = r32
	goto block80
block64:
	frame.Fail()
	goto block65
block65:
	frame.Recover(r0)
	goto block66
block66:
	r33 = frame.Read()
	if frame.Flow == 0 {
		goto block67
	} else {
		goto block73
	}
block67:
	r34 = '\''
	goto block68
block68:
	r35 = r33 == r34
	goto block69
block69:
	if r35 {
		goto block70
	} else {
		goto block72
	}
block70:
	r36 = '\''
	goto block71
block71:
	ret0 = r36
	goto block80
block72:
	frame.Fail()
	goto block73
block73:
	frame.Recover(r0)
	goto block74
block74:
	r37 = frame.Read()
	if frame.Flow == 0 {
		goto block75
	} else {
		goto block82
	}
block75:
	r38 = '"'
	goto block76
block76:
	r39 = r37 == r38
	goto block77
block77:
	if r39 {
		goto block78
	} else {
		goto block81
	}
block78:
	r40 = '"'
	goto block79
block79:
	ret0 = r40
	goto block80
block80:
	return
block81:
	frame.Fail()
	goto block82
block82:
	return
}
func StrT(frame *dub.DubState) (ret0 *StrTok) {
	var r0 []rune
	var r1 string
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 []rune
	var r7 int
	var r8 int
	var r9 []rune
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 rune
	var r14 bool
	var r15 []rune
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 []rune
	var r20 rune
	var r21 []rune
	var r22 rune
	var r23 rune
	var r24 bool
	var r25 string
	var r26 string
	var r27 []rune
	var r28 string
	var r29 *StrTok
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Checkpoint()
	goto block2
block2:
	r3 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block47
	}
block3:
	r4 = '"'
	goto block4
block4:
	r5 = r3 == r4
	goto block5
block5:
	if r5 {
		goto block6
	} else {
		goto block46
	}
block6:
	r6 = []rune{}
	goto block7
block7:
	r0 = r6
	goto block8
block8:
	r7 = frame.Checkpoint()
	goto block9
block9:
	r8 = frame.Checkpoint()
	goto block10
block10:
	r9 = r0
	goto block11
block11:
	r10 = frame.Read()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block21
	}
block12:
	r11 = '"'
	goto block13
block13:
	r12 = r10 == r11
	goto block14
block14:
	if r12 {
		goto block18
	} else {
		goto block15
	}
block15:
	r13 = '\\'
	goto block16
block16:
	r14 = r10 == r13
	goto block17
block17:
	if r14 {
		goto block18
	} else {
		goto block19
	}
block18:
	frame.Fail()
	goto block21
block19:
	r15 = append(r9, r10)
	goto block20
block20:
	r0 = r15
	goto block8
block21:
	frame.Recover(r8)
	goto block22
block22:
	r16 = frame.Read()
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block31
	}
block23:
	r17 = '\\'
	goto block24
block24:
	r18 = r16 == r17
	goto block25
block25:
	if r18 {
		goto block26
	} else {
		goto block30
	}
block26:
	r19 = r0
	goto block27
block27:
	r20 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block28
	} else {
		goto block31
	}
block28:
	r21 = append(r19, r20)
	goto block29
block29:
	r0 = r21
	goto block8
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(r7)
	goto block32
block32:
	r22 = frame.Read()
	if frame.Flow == 0 {
		goto block33
	} else {
		goto block47
	}
block33:
	r23 = '"'
	goto block34
block34:
	r24 = r22 == r23
	goto block35
block35:
	if r24 {
		goto block36
	} else {
		goto block45
	}
block36:
	r25 = frame.Slice(r2)
	goto block37
block37:
	r1 = r25
	goto block38
block38:
	S(frame)
	if frame.Flow == 0 {
		goto block39
	} else {
		goto block47
	}
block39:
	r26 = r1
	goto block40
block40:
	r27 = r0
	goto block41
block41:
	r28 = string(r27)
	if frame.Flow == 0 {
		goto block42
	} else {
		goto block47
	}
block42:
	r29 = &StrTok{Text: r26, Value: r28}
	goto block43
block43:
	ret0 = r29
	goto block44
block44:
	return
block45:
	frame.Fail()
	goto block47
block46:
	frame.Fail()
	goto block47
block47:
	return
}
func Rune(frame *dub.DubState) (ret0 *RuneTok) {
	var r0 rune
	var r1 string
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 int
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 rune
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 string
	var r20 string
	var r21 rune
	var r22 *RuneTok
	goto block0
block0:
	goto block1
block1:
	r0 = '\x00'
	goto block2
block2:
	r2 = frame.Checkpoint()
	goto block3
block3:
	r3 = frame.Read()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block39
	}
block4:
	r4 = '\''
	goto block5
block5:
	r5 = r3 == r4
	goto block6
block6:
	if r5 {
		goto block7
	} else {
		goto block38
	}
block7:
	r6 = frame.Checkpoint()
	goto block8
block8:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block17
	}
block9:
	r8 = '\\'
	goto block10
block10:
	r9 = r7 == r8
	goto block11
block11:
	if r9 {
		goto block15
	} else {
		goto block12
	}
block12:
	r10 = '\''
	goto block13
block13:
	r11 = r7 == r10
	goto block14
block14:
	if r11 {
		goto block15
	} else {
		goto block16
	}
block15:
	frame.Fail()
	goto block17
block16:
	r0 = r7
	goto block24
block17:
	frame.Recover(r6)
	goto block18
block18:
	r12 = frame.Read()
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block39
	}
block19:
	r13 = '\\'
	goto block20
block20:
	r14 = r12 == r13
	goto block21
block21:
	if r14 {
		goto block22
	} else {
		goto block37
	}
block22:
	r15 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block23
	} else {
		goto block39
	}
block23:
	r0 = r15
	goto block24
block24:
	r16 = frame.Read()
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block39
	}
block25:
	r17 = '\''
	goto block26
block26:
	r18 = r16 == r17
	goto block27
block27:
	if r18 {
		goto block28
	} else {
		goto block36
	}
block28:
	r19 = frame.Slice(r2)
	goto block29
block29:
	r1 = r19
	goto block30
block30:
	S(frame)
	if frame.Flow == 0 {
		goto block31
	} else {
		goto block39
	}
block31:
	r20 = r1
	goto block32
block32:
	r21 = r0
	goto block33
block33:
	r22 = &RuneTok{Text: r20, Value: r21}
	goto block34
block34:
	ret0 = r22
	goto block35
block35:
	return
block36:
	frame.Fail()
	goto block39
block37:
	frame.Fail()
	goto block39
block38:
	frame.Fail()
	goto block39
block39:
	return
}
func StringMatchExpr(frame *dub.DubState) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 TextMatch
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 TextMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block18
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
		goto block17
	}
block5:
	S(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block18
	}
block6:
	r4 = Choice(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block18
	}
block7:
	r0 = r4
	goto block8
block8:
	r5 = frame.Read()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block18
	}
block9:
	r6 = '/'
	goto block10
block10:
	r7 = r5 == r6
	goto block11
block11:
	if r7 {
		goto block12
	} else {
		goto block16
	}
block12:
	S(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block18
	}
block13:
	r8 = r0
	goto block14
block14:
	ret0 = r8
	goto block15
block15:
	return
block16:
	frame.Fail()
	goto block18
block17:
	frame.Fail()
	goto block18
block18:
	return
}
func RuneMatchExpr(frame *dub.DubState) (ret0 *RuneMatch) {
	var r0 *RuneMatch
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 *RuneMatch
	var r5 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block12
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
		goto block11
	}
block5:
	S(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block12
	}
block6:
	r4 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block12
	}
block7:
	r0 = r4
	goto block8
block8:
	r5 = r0
	goto block9
block9:
	ret0 = r5
	goto block10
block10:
	return
block11:
	frame.Fail()
	goto block12
block12:
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
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block14
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
	goto block14
block13:
	ret0 = r1
	goto block25
block14:
	frame.Recover(r0)
	goto block15
block15:
	r8 = frame.Read()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block27
	}
block16:
	r9 = '\\'
	goto block17
block17:
	r10 = r8 == r9
	goto block18
block18:
	if r10 {
		goto block19
	} else {
		goto block26
	}
block19:
	r11 = frame.Checkpoint()
	goto block20
block20:
	r12 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block22
	}
block21:
	ret0 = r12
	goto block25
block22:
	frame.Recover(r11)
	goto block23
block23:
	r13 = frame.Read()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block27
	}
block24:
	ret0 = r13
	goto block25
block25:
	return
block26:
	frame.Fail()
	goto block27
block27:
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
		goto block19
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
	r5 = frame.Read()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block13
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
		goto block12
	}
block10:
	r8 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block13
	}
block11:
	r1 = r8
	goto block14
block12:
	frame.Fail()
	goto block13
block13:
	frame.Recover(r4)
	goto block14
block14:
	r9 = r0
	goto block15
block15:
	r10 = r1
	goto block16
block16:
	r11 = &RuneFilter{Min: r9, Max: r10}
	goto block17
block17:
	ret0 = r11
	goto block18
block18:
	return
block19:
	return
}
func MatchRune(frame *dub.DubState) (ret0 *RuneMatch) {
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
	var r21 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block36
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
		goto block35
	}
block5:
	r5 = false
	goto block6
block6:
	r0 = r5
	goto block7
block7:
	r6 = []*RuneFilter{}
	goto block8
block8:
	r1 = r6
	goto block9
block9:
	r7 = frame.Checkpoint()
	goto block10
block10:
	r8 = frame.Read()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block17
	}
block11:
	r9 = '^'
	goto block12
block12:
	r10 = r8 == r9
	goto block13
block13:
	if r10 {
		goto block14
	} else {
		goto block16
	}
block14:
	r11 = true
	goto block15
block15:
	r0 = r11
	goto block18
block16:
	frame.Fail()
	goto block17
block17:
	frame.Recover(r7)
	goto block18
block18:
	r12 = frame.Checkpoint()
	goto block19
block19:
	r13 = r1
	goto block20
block20:
	r14 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block23
	}
block21:
	r15 = append(r13, r14)
	goto block22
block22:
	r1 = r15
	goto block18
block23:
	frame.Recover(r12)
	goto block24
block24:
	r16 = frame.Read()
	if frame.Flow == 0 {
		goto block25
	} else {
		goto block36
	}
block25:
	r17 = ']'
	goto block26
block26:
	r18 = r16 == r17
	goto block27
block27:
	if r18 {
		goto block28
	} else {
		goto block34
	}
block28:
	S(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block36
	}
block29:
	r19 = r0
	goto block30
block30:
	r20 = r1
	goto block31
block31:
	r21 = &RuneMatch{Invert: r19, Filters: r20}
	goto block32
block32:
	ret0 = r21
	goto block33
block33:
	return
block34:
	frame.Fail()
	goto block36
block35:
	frame.Fail()
	goto block36
block36:
	return
}
func Atom(frame *dub.DubState) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 *RuneMatch
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 TextMatch
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 TextMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Checkpoint()
	goto block2
block2:
	r2 = MatchRune(frame)
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block4
	}
block3:
	ret0 = r2
	goto block17
block4:
	frame.Recover(r1)
	goto block5
block5:
	r3 = frame.Read()
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block20
	}
block6:
	r4 = '('
	goto block7
block7:
	r5 = r3 == r4
	goto block8
block8:
	if r5 {
		goto block9
	} else {
		goto block19
	}
block9:
	r6 = Choice(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block20
	}
block10:
	r0 = r6
	goto block11
block11:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block20
	}
block12:
	r8 = ')'
	goto block13
block13:
	r9 = r7 == r8
	goto block14
block14:
	if r9 {
		goto block15
	} else {
		goto block18
	}
block15:
	r10 = r0
	goto block16
block16:
	ret0 = r10
	goto block17
block17:
	return
block18:
	frame.Fail()
	goto block20
block19:
	frame.Fail()
	goto block20
block20:
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
		goto block42
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = frame.Checkpoint()
	goto block4
block4:
	r3 = frame.Read()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block14
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
		goto block13
	}
block8:
	S(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block14
	}
block9:
	r6 = r0
	goto block10
block10:
	r7 = 0
	goto block11
block11:
	r8 = &MatchRepeat{Match: r6, Min: r7}
	goto block12
block12:
	ret0 = r8
	goto block41
block13:
	frame.Fail()
	goto block14
block14:
	frame.Recover(r2)
	goto block15
block15:
	r9 = frame.Read()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block25
	}
block16:
	r10 = '+'
	goto block17
block17:
	r11 = r9 == r10
	goto block18
block18:
	if r11 {
		goto block19
	} else {
		goto block24
	}
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block25
	}
block20:
	r12 = r0
	goto block21
block21:
	r13 = 1
	goto block22
block22:
	r14 = &MatchRepeat{Match: r12, Min: r13}
	goto block23
block23:
	ret0 = r14
	goto block41
block24:
	frame.Fail()
	goto block25
block25:
	frame.Recover(r2)
	goto block26
block26:
	r15 = frame.Read()
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block38
	}
block27:
	r16 = '?'
	goto block28
block28:
	r17 = r15 == r16
	goto block29
block29:
	if r17 {
		goto block30
	} else {
		goto block37
	}
block30:
	S(frame)
	if frame.Flow == 0 {
		goto block31
	} else {
		goto block38
	}
block31:
	r18 = r0
	goto block32
block32:
	r19 = []TextMatch{}
	goto block33
block33:
	r20 = &MatchSequence{Matches: r19}
	goto block34
block34:
	r21 = []TextMatch{r18, r20}
	goto block35
block35:
	r22 = &MatchChoice{Matches: r21}
	goto block36
block36:
	ret0 = r22
	goto block41
block37:
	frame.Fail()
	goto block38
block38:
	frame.Recover(r2)
	goto block39
block39:
	r23 = r0
	goto block40
block40:
	ret0 = r23
	goto block41
block41:
	return
block42:
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
		goto block36
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
	r6 = frame.Read()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block32
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
		goto block31
	}
block11:
	S(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block32
	}
block12:
	r9 = r1
	goto block13
block13:
	r10 = Sequence(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block32
	}
block14:
	r11 = append(r9, r10)
	goto block15
block15:
	r1 = r11
	goto block16
block16:
	r12 = frame.Checkpoint()
	goto block17
block17:
	r13 = frame.Read()
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block27
	}
block18:
	r14 = '|'
	goto block19
block19:
	r15 = r13 == r14
	goto block20
block20:
	if r15 {
		goto block21
	} else {
		goto block26
	}
block21:
	S(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block27
	}
block22:
	r16 = r1
	goto block23
block23:
	r17 = Sequence(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block27
	}
block24:
	r18 = append(r16, r17)
	goto block25
block25:
	r1 = r18
	goto block16
block26:
	frame.Fail()
	goto block27
block27:
	frame.Recover(r12)
	goto block28
block28:
	r19 = r1
	goto block29
block29:
	r20 = &MatchChoice{Matches: r19}
	goto block30
block30:
	ret0 = r20
	goto block35
block31:
	frame.Fail()
	goto block32
block32:
	frame.Recover(r3)
	goto block33
block33:
	r21 = r0
	goto block34
block34:
	ret0 = r21
	goto block35
block35:
	return
block36:
	return
}
