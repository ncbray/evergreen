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
func Int(frame *dub.DubState) (ret0 *IntTok) {
	var r0 string
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 bool
	var r7 int
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 bool
	var r13 string
	var r14 string
	var r15 *IntTok
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
		goto block27
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
		goto block26
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
		goto block26
	}
block9:
	r7 = frame.Checkpoint()
	goto block10
block10:
	r8 = frame.Read()
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block18
	}
block11:
	r9 = '0'
	goto block12
block12:
	r10 = r8 >= r9
	goto block13
block13:
	if r10 {
		goto block14
	} else {
		goto block17
	}
block14:
	r11 = '9'
	goto block15
block15:
	r12 = r8 <= r11
	goto block16
block16:
	if r12 {
		goto block9
	} else {
		goto block17
	}
block17:
	frame.Fail()
	goto block18
block18:
	frame.Recover(r7)
	goto block19
block19:
	r13 = frame.Slice(r1)
	goto block20
block20:
	r0 = r13
	goto block21
block21:
	S(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block27
	}
block22:
	r14 = r0
	goto block23
block23:
	r15 = &IntTok{Text: r14}
	goto block24
block24:
	ret0 = r15
	goto block25
block25:
	return
block26:
	frame.Fail()
	goto block27
block27:
	return
}
func EscapedChar(frame *dub.DubState) (ret0 rune) {
	var r0 rune
	var r1 rune
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 rune
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
	var r17 rune
	var r18 rune
	var r19 rune
	var r20 bool
	var r21 rune
	var r22 rune
	var r23 rune
	var r24 bool
	var r25 rune
	var r26 rune
	var r27 rune
	var r28 bool
	var r29 rune
	var r30 rune
	var r31 rune
	var r32 bool
	var r33 rune
	var r34 rune
	var r35 rune
	var r36 bool
	var r37 rune
	var r38 rune
	var r39 rune
	var r40 bool
	var r41 rune
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block65
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = r0
	goto block4
block4:
	r3 = 'a'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block9
	}
block7:
	r5 = '\a'
	goto block8
block8:
	ret0 = r5
	goto block63
block9:
	r6 = r0
	goto block10
block10:
	r7 = 'b'
	goto block11
block11:
	r8 = r6 == r7
	goto block12
block12:
	if r8 {
		goto block13
	} else {
		goto block15
	}
block13:
	r9 = '\b'
	goto block14
block14:
	ret0 = r9
	goto block63
block15:
	r10 = r0
	goto block16
block16:
	r11 = 'f'
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
	r13 = '\f'
	goto block20
block20:
	ret0 = r13
	goto block63
block21:
	r14 = r0
	goto block22
block22:
	r15 = 'n'
	goto block23
block23:
	r16 = r14 == r15
	goto block24
block24:
	if r16 {
		goto block25
	} else {
		goto block27
	}
block25:
	r17 = '\n'
	goto block26
block26:
	ret0 = r17
	goto block63
block27:
	r18 = r0
	goto block28
block28:
	r19 = 'r'
	goto block29
block29:
	r20 = r18 == r19
	goto block30
block30:
	if r20 {
		goto block31
	} else {
		goto block33
	}
block31:
	r21 = '\r'
	goto block32
block32:
	ret0 = r21
	goto block63
block33:
	r22 = r0
	goto block34
block34:
	r23 = 't'
	goto block35
block35:
	r24 = r22 == r23
	goto block36
block36:
	if r24 {
		goto block37
	} else {
		goto block39
	}
block37:
	r25 = '\t'
	goto block38
block38:
	ret0 = r25
	goto block63
block39:
	r26 = r0
	goto block40
block40:
	r27 = 'v'
	goto block41
block41:
	r28 = r26 == r27
	goto block42
block42:
	if r28 {
		goto block43
	} else {
		goto block45
	}
block43:
	r29 = '\v'
	goto block44
block44:
	ret0 = r29
	goto block63
block45:
	r30 = r0
	goto block46
block46:
	r31 = '\\'
	goto block47
block47:
	r32 = r30 == r31
	goto block48
block48:
	if r32 {
		goto block49
	} else {
		goto block51
	}
block49:
	r33 = '\\'
	goto block50
block50:
	ret0 = r33
	goto block63
block51:
	r34 = r0
	goto block52
block52:
	r35 = '\''
	goto block53
block53:
	r36 = r34 == r35
	goto block54
block54:
	if r36 {
		goto block55
	} else {
		goto block57
	}
block55:
	r37 = '\''
	goto block56
block56:
	ret0 = r37
	goto block63
block57:
	r38 = r0
	goto block58
block58:
	r39 = '"'
	goto block59
block59:
	r40 = r38 == r39
	goto block60
block60:
	if r40 {
		goto block61
	} else {
		goto block64
	}
block61:
	r41 = '"'
	goto block62
block62:
	ret0 = r41
	goto block63
block63:
	return
block64:
	frame.Fail()
	goto block65
block65:
	return
}
func StrT(frame *dub.DubState) (ret0 *StrTok) {
	var r0 rune
	var r1 string
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 int
	var r7 rune
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 rune
	var r15 rune
	var r16 bool
	var r17 string
	var r18 string
	var r19 *StrTok
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
		goto block33
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
		goto block32
	}
block6:
	r6 = frame.Checkpoint()
	goto block7
block7:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block19
	}
block8:
	r0 = r7
	goto block9
block9:
	r8 = r0
	goto block10
block10:
	r9 = '"'
	goto block11
block11:
	r10 = r8 == r9
	goto block12
block12:
	if r10 {
		goto block13
	} else {
		goto block14
	}
block13:
	frame.Fail()
	goto block19
block14:
	r11 = r0
	goto block15
block15:
	r12 = '\\'
	goto block16
block16:
	r13 = r11 == r12
	goto block17
block17:
	if r13 {
		goto block18
	} else {
		goto block6
	}
block18:
	EscapedChar(frame)
	if frame.Flow == 0 {
		goto block6
	} else {
		goto block19
	}
block19:
	frame.Recover(r6)
	goto block20
block20:
	r14 = frame.Read()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block33
	}
block21:
	r15 = '"'
	goto block22
block22:
	r16 = r14 == r15
	goto block23
block23:
	if r16 {
		goto block24
	} else {
		goto block31
	}
block24:
	r17 = frame.Slice(r2)
	goto block25
block25:
	r1 = r17
	goto block26
block26:
	S(frame)
	if frame.Flow == 0 {
		goto block27
	} else {
		goto block33
	}
block27:
	r18 = r1
	goto block28
block28:
	r19 = &StrTok{Text: r18}
	goto block29
block29:
	ret0 = r19
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
func Rune(frame *dub.DubState) (ret0 *RuneTok) {
	var r0 rune
	var r1 string
	var r2 int
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 rune
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 bool
	var r13 string
	var r14 string
	var r15 *RuneTok
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
		goto block26
	}
block3:
	r4 = '\''
	goto block4
block4:
	r5 = r3 == r4
	goto block5
block5:
	if r5 {
		goto block6
	} else {
		goto block25
	}
block6:
	r6 = frame.Read()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block26
	}
block7:
	r0 = r6
	goto block8
block8:
	r7 = r0
	goto block9
block9:
	r8 = '\\'
	goto block10
block10:
	r9 = r7 == r8
	goto block11
block11:
	if r9 {
		goto block12
	} else {
		goto block13
	}
block12:
	EscapedChar(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block26
	}
block13:
	r10 = frame.Read()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block26
	}
block14:
	r11 = '\''
	goto block15
block15:
	r12 = r10 == r11
	goto block16
block16:
	if r12 {
		goto block17
	} else {
		goto block24
	}
block17:
	r13 = frame.Slice(r2)
	goto block18
block18:
	r1 = r13
	goto block19
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block26
	}
block20:
	r14 = r1
	goto block21
block21:
	r15 = &RuneTok{Text: r14}
	goto block22
block22:
	ret0 = r15
	goto block23
block23:
	return
block24:
	frame.Fail()
	goto block26
block25:
	frame.Fail()
	goto block26
block26:
	return
}
func MatchExpr(frame *dub.DubState) (ret0 TextMatch) {
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
func ParseRuneFilterRune(frame *dub.DubState) (ret0 rune) {
	var r0 rune
	var r1 rune
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 int
	var r12 rune
	var r13 rune
	var r14 rune
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block26
	}
block2:
	r0 = r1
	goto block3
block3:
	r2 = r0
	goto block4
block4:
	r3 = ']'
	goto block5
block5:
	r4 = r2 == r3
	goto block6
block6:
	if r4 {
		goto block7
	} else {
		goto block8
	}
block7:
	frame.Fail()
	goto block26
block8:
	r5 = r0
	goto block9
block9:
	r6 = '-'
	goto block10
block10:
	r7 = r5 == r6
	goto block11
block11:
	if r7 {
		goto block12
	} else {
		goto block13
	}
block12:
	frame.Fail()
	goto block26
block13:
	r8 = r0
	goto block14
block14:
	r9 = '\\'
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
	r11 = frame.Checkpoint()
	goto block18
block18:
	r12 = EscapedChar(frame)
	if frame.Flow == 0 {
		goto block19
	} else {
		goto block20
	}
block19:
	ret0 = r12
	goto block25
block20:
	frame.Recover(r11)
	goto block21
block21:
	r13 = frame.Read()
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block26
	}
block22:
	ret0 = r13
	goto block25
block23:
	r14 = r0
	goto block24
block24:
	ret0 = r14
	goto block25
block25:
	return
block26:
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
	var r0 []*RuneFilter
	var r1 rune
	var r2 rune
	var r3 bool
	var r4 []*RuneFilter
	var r5 int
	var r6 []*RuneFilter
	var r7 *RuneFilter
	var r8 []*RuneFilter
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 []*RuneFilter
	var r13 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block24
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
		goto block23
	}
block5:
	r4 = []*RuneFilter{}
	goto block6
block6:
	r0 = r4
	goto block7
block7:
	r5 = frame.Checkpoint()
	goto block8
block8:
	r6 = r0
	goto block9
block9:
	r7 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block12
	}
block10:
	r8 = append(r6, r7)
	goto block11
block11:
	r0 = r8
	goto block7
block12:
	frame.Recover(r5)
	goto block13
block13:
	r9 = frame.Read()
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block24
	}
block14:
	r10 = ']'
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
	S(frame)
	if frame.Flow == 0 {
		goto block18
	} else {
		goto block24
	}
block18:
	r12 = r0
	goto block19
block19:
	r13 = &RuneMatch{Filters: r12}
	goto block20
block20:
	ret0 = r13
	goto block21
block21:
	return
block22:
	frame.Fail()
	goto block24
block23:
	frame.Fail()
	goto block24
block24:
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
	var r1 rune
	var r2 TextMatch
	var r3 int
	var r4 rune
	var r5 rune
	var r6 rune
	var r7 bool
	var r8 TextMatch
	var r9 int
	var r10 *MatchRepeat
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 TextMatch
	var r15 int
	var r16 *MatchRepeat
	var r17 rune
	var r18 rune
	var r19 bool
	var r20 TextMatch
	var r21 []TextMatch
	var r22 *MatchSequence
	var r23 []TextMatch
	var r24 *MatchChoice
	var r25 TextMatch
	goto block0
block0:
	goto block1
block1:
	r2 = Atom(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block40
	}
block2:
	r0 = r2
	goto block3
block3:
	r3 = frame.Checkpoint()
	goto block4
block4:
	r4 = frame.Read()
	if frame.Flow == 0 {
		goto block5
	} else {
		goto block36
	}
block5:
	r1 = r4
	goto block6
block6:
	r5 = r1
	goto block7
block7:
	r6 = '*'
	goto block8
block8:
	r7 = r5 == r6
	goto block9
block9:
	if r7 {
		goto block10
	} else {
		goto block15
	}
block10:
	S(frame)
	if frame.Flow == 0 {
		goto block11
	} else {
		goto block36
	}
block11:
	r8 = r0
	goto block12
block12:
	r9 = 0
	goto block13
block13:
	r10 = &MatchRepeat{Match: r8, Min: r9}
	goto block14
block14:
	ret0 = r10
	goto block39
block15:
	r11 = r1
	goto block16
block16:
	r12 = '+'
	goto block17
block17:
	r13 = r11 == r12
	goto block18
block18:
	if r13 {
		goto block19
	} else {
		goto block24
	}
block19:
	S(frame)
	if frame.Flow == 0 {
		goto block20
	} else {
		goto block36
	}
block20:
	r14 = r0
	goto block21
block21:
	r15 = 1
	goto block22
block22:
	r16 = &MatchRepeat{Match: r14, Min: r15}
	goto block23
block23:
	ret0 = r16
	goto block39
block24:
	r17 = r1
	goto block25
block25:
	r18 = '?'
	goto block26
block26:
	r19 = r17 == r18
	goto block27
block27:
	if r19 {
		goto block28
	} else {
		goto block35
	}
block28:
	S(frame)
	if frame.Flow == 0 {
		goto block29
	} else {
		goto block36
	}
block29:
	r20 = r0
	goto block30
block30:
	r21 = []TextMatch{}
	goto block31
block31:
	r22 = &MatchSequence{Matches: r21}
	goto block32
block32:
	r23 = []TextMatch{r20, r22}
	goto block33
block33:
	r24 = &MatchChoice{Matches: r23}
	goto block34
block34:
	ret0 = r24
	goto block39
block35:
	frame.Fail()
	goto block36
block36:
	frame.Recover(r3)
	goto block37
block37:
	r25 = r0
	goto block38
block38:
	ret0 = r25
	goto block39
block39:
	return
block40:
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
	r0 = r14
	goto block21
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
	r0 = r20
	goto block33
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
