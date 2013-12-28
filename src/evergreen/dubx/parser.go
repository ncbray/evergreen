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
type IdTok struct {
	Text string
}

func (node *IdTok) isToken() {
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
	var r0 rune
	var r1 int
	var r2 rune
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
		goto block21
	}
block3:
	r0 = r2
	goto block4
block4:
	r3 = r0
	goto block5
block5:
	r4 = ' '
	goto block6
block6:
	r5 = r3 != r4
	goto block7
block7:
	if r5 {
		goto block8
	} else {
		goto block1
	}
block8:
	r6 = r0
	goto block9
block9:
	r7 = '\t'
	goto block10
block10:
	r8 = r6 != r7
	goto block11
block11:
	if r8 {
		goto block12
	} else {
		goto block1
	}
block12:
	r9 = r0
	goto block13
block13:
	r10 = '\r'
	goto block14
block14:
	r11 = r9 != r10
	goto block15
block15:
	if r11 {
		goto block16
	} else {
		goto block1
	}
block16:
	r12 = r0
	goto block17
block17:
	r13 = '\n'
	goto block18
block18:
	r14 = r12 != r13
	goto block19
block19:
	if r14 {
		goto block20
	} else {
		goto block1
	}
block20:
	frame.Fail()
	goto block21
block21:
	frame.Recover(r1)
	goto block22
block22:
	return
}
func Ident(frame *dub.DubState) (ret0 *IdTok) {
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
	var r31 *IdTok
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
		goto block51
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
		goto block50
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
		goto block51
	}
block46:
	r30 = r0
	goto block47
block47:
	r31 = &IdTok{Text: r30}
	goto block48
block48:
	ret0 = r31
	goto block49
block49:
	return
block50:
	frame.Fail()
	goto block51
block51:
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
func StrT(frame *dub.DubState) (ret0 *StrTok) {
	var r0 rune
	var r1 string
	var r2 int
	var r3 int
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 int
	var r8 rune
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 int
	var r16 rune
	var r17 rune
	var r18 bool
	var r19 string
	var r20 string
	var r21 *StrTok
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Checkpoint()
	goto block2
block2:
	r3 = frame.Checkpoint()
	goto block3
block3:
	r4 = frame.Read()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block37
	}
block4:
	r5 = '"'
	goto block5
block5:
	r6 = r4 == r5
	goto block6
block6:
	if r6 {
		goto block7
	} else {
		goto block36
	}
block7:
	frame.Slice(r3)
	goto block8
block8:
	r7 = frame.Checkpoint()
	goto block9
block9:
	r8 = frame.Read()
	if frame.Flow == 0 {
		goto block10
	} else {
		goto block21
	}
block10:
	r0 = r8
	goto block11
block11:
	r9 = r0
	goto block12
block12:
	r10 = '\\'
	goto block13
block13:
	r11 = r9 == r10
	goto block14
block14:
	if r11 {
		goto block15
	} else {
		goto block16
	}
block15:
	frame.Read()
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block21
	}
block16:
	r12 = r0
	goto block17
block17:
	r13 = '"'
	goto block18
block18:
	r14 = r12 == r13
	goto block19
block19:
	if r14 {
		goto block20
	} else {
		goto block8
	}
block20:
	frame.Fail()
	goto block21
block21:
	frame.Recover(r7)
	goto block22
block22:
	r15 = frame.Checkpoint()
	goto block23
block23:
	r16 = frame.Read()
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block37
	}
block24:
	r17 = '"'
	goto block25
block25:
	r18 = r16 == r17
	goto block26
block26:
	if r18 {
		goto block27
	} else {
		goto block35
	}
block27:
	frame.Slice(r15)
	goto block28
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
		goto block37
	}
block31:
	r20 = r1
	goto block32
block32:
	r21 = &StrTok{Text: r20}
	goto block33
block33:
	ret0 = r21
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
func Rune(frame *dub.DubState) (ret0 *RuneTok) {
	var r0 rune
	var r1 string
	var r2 int
	var r3 int
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 rune
	var r8 rune
	var r9 rune
	var r10 bool
	var r11 int
	var r12 rune
	var r13 rune
	var r14 bool
	var r15 string
	var r16 string
	var r17 *RuneTok
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Checkpoint()
	goto block2
block2:
	r3 = frame.Checkpoint()
	goto block3
block3:
	r4 = frame.Read()
	if frame.Flow == 0 {
		goto block4
	} else {
		goto block30
	}
block4:
	r5 = '\''
	goto block5
block5:
	r6 = r4 == r5
	goto block6
block6:
	if r6 {
		goto block7
	} else {
		goto block29
	}
block7:
	frame.Slice(r3)
	goto block8
block8:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block30
	}
block9:
	r0 = r7
	goto block10
block10:
	r8 = r0
	goto block11
block11:
	r9 = '\\'
	goto block12
block12:
	r10 = r8 == r9
	goto block13
block13:
	if r10 {
		goto block14
	} else {
		goto block15
	}
block14:
	frame.Read()
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block30
	}
block15:
	r11 = frame.Checkpoint()
	goto block16
block16:
	r12 = frame.Read()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block30
	}
block17:
	r13 = '\''
	goto block18
block18:
	r14 = r12 == r13
	goto block19
block19:
	if r14 {
		goto block20
	} else {
		goto block28
	}
block20:
	frame.Slice(r11)
	goto block21
block21:
	r15 = frame.Slice(r2)
	goto block22
block22:
	r1 = r15
	goto block23
block23:
	S(frame)
	if frame.Flow == 0 {
		goto block24
	} else {
		goto block30
	}
block24:
	r16 = r1
	goto block25
block25:
	r17 = &RuneTok{Text: r16}
	goto block26
block26:
	ret0 = r17
	goto block27
block27:
	return
block28:
	frame.Fail()
	goto block30
block29:
	frame.Fail()
	goto block30
block30:
	return
}
func MatchExpr(frame *dub.DubState) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 rune
	var r3 rune
	var r4 bool
	var r5 TextMatch
	var r6 int
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
	r2 = frame.Read()
	if frame.Flow == 0 {
		goto block3
	} else {
		goto block22
	}
block3:
	r3 = '/'
	goto block4
block4:
	r4 = r2 == r3
	goto block5
block5:
	if r4 {
		goto block6
	} else {
		goto block21
	}
block6:
	frame.Slice(r1)
	goto block7
block7:
	S(frame)
	if frame.Flow == 0 {
		goto block8
	} else {
		goto block22
	}
block8:
	r5 = Choice(frame)
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block22
	}
block9:
	r0 = r5
	goto block10
block10:
	r6 = frame.Checkpoint()
	goto block11
block11:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block22
	}
block12:
	r8 = '/'
	goto block13
block13:
	r9 = r7 == r8
	goto block14
block14:
	if r9 {
		goto block15
	} else {
		goto block20
	}
block15:
	frame.Slice(r6)
	goto block16
block16:
	S(frame)
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block22
	}
block17:
	r10 = r0
	goto block18
block18:
	ret0 = r10
	goto block19
block19:
	return
block20:
	frame.Fail()
	goto block22
block21:
	frame.Fail()
	goto block22
block22:
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
	goto block0
block0:
	goto block1
block1:
	r1 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block16
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
	goto block16
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
	goto block16
block13:
	r8 = r0
	goto block14
block14:
	ret0 = r8
	goto block15
block15:
	return
block16:
	return
}
func ParseRuneFilter(frame *dub.DubState) (ret0 *RuneFilter) {
	var r0 rune
	var r1 rune
	var r2 rune
	var r3 rune
	var r4 rune
	var r5 int
	var r6 rune
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 rune
	var r11 rune
	var r12 rune
	var r13 *RuneFilter
	goto block0
block0:
	goto block1
block1:
	r3 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block21
	}
block2:
	r0 = r3
	goto block3
block3:
	r4 = r0
	goto block4
block4:
	r1 = r4
	goto block5
block5:
	r5 = frame.Checkpoint()
	goto block6
block6:
	r6 = frame.Read()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block15
	}
block7:
	r2 = r6
	goto block8
block8:
	r7 = r2
	goto block9
block9:
	r8 = '-'
	goto block10
block10:
	r9 = r7 != r8
	goto block11
block11:
	if r9 {
		goto block12
	} else {
		goto block13
	}
block12:
	frame.Fail()
	goto block15
block13:
	r10 = ParseRuneFilterRune(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block15
	}
block14:
	r1 = r10
	goto block16
block15:
	frame.Recover(r5)
	goto block16
block16:
	r11 = r0
	goto block17
block17:
	r12 = r1
	goto block18
block18:
	r13 = &RuneFilter{Min: r11, Max: r12}
	goto block19
block19:
	ret0 = r13
	goto block20
block20:
	return
block21:
	return
}
func MatchRune(frame *dub.DubState) (ret0 *RuneMatch) {
	var r0 rune
	var r1 []*RuneFilter
	var r2 rune
	var r3 rune
	var r4 rune
	var r5 bool
	var r6 []*RuneFilter
	var r7 int
	var r8 []*RuneFilter
	var r9 *RuneFilter
	var r10 []*RuneFilter
	var r11 rune
	var r12 rune
	var r13 bool
	var r14 []*RuneFilter
	var r15 *RuneMatch
	goto block0
block0:
	goto block1
block1:
	r2 = frame.Read()
	if frame.Flow == 0 {
		goto block2
	} else {
		goto block26
	}
block2:
	r0 = r2
	goto block3
block3:
	r3 = r0
	goto block4
block4:
	r4 = '['
	goto block5
block5:
	r5 = r3 != r4
	goto block6
block6:
	if r5 {
		goto block7
	} else {
		goto block8
	}
block7:
	frame.Fail()
	goto block26
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
	r8 = r1
	goto block12
block12:
	r9 = ParseRuneFilter(frame)
	if frame.Flow == 0 {
		goto block13
	} else {
		goto block15
	}
block13:
	r10 = append(r8, r9)
	goto block14
block14:
	r1 = r10
	goto block10
block15:
	frame.Recover(r7)
	goto block16
block16:
	r11 = frame.Read()
	if frame.Flow == 0 {
		goto block17
	} else {
		goto block26
	}
block17:
	r12 = ']'
	goto block18
block18:
	r13 = r11 != r12
	goto block19
block19:
	if r13 {
		goto block20
	} else {
		goto block21
	}
block20:
	frame.Fail()
	goto block26
block21:
	S(frame)
	if frame.Flow == 0 {
		goto block22
	} else {
		goto block26
	}
block22:
	r14 = r1
	goto block23
block23:
	r15 = &RuneMatch{Filters: r14}
	goto block24
block24:
	ret0 = r15
	goto block25
block25:
	return
block26:
	return
}
func Atom(frame *dub.DubState) (ret0 TextMatch) {
	var r0 TextMatch
	var r1 int
	var r2 *RuneMatch
	var r3 int
	var r4 rune
	var r5 rune
	var r6 bool
	var r7 TextMatch
	var r8 int
	var r9 rune
	var r10 rune
	var r11 bool
	var r12 TextMatch
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
	goto block21
block4:
	frame.Recover(r1)
	goto block5
block5:
	r3 = frame.Checkpoint()
	goto block6
block6:
	r4 = frame.Read()
	if frame.Flow == 0 {
		goto block7
	} else {
		goto block24
	}
block7:
	r5 = '('
	goto block8
block8:
	r6 = r4 == r5
	goto block9
block9:
	if r6 {
		goto block10
	} else {
		goto block23
	}
block10:
	frame.Slice(r3)
	goto block11
block11:
	r7 = Choice(frame)
	if frame.Flow == 0 {
		goto block12
	} else {
		goto block24
	}
block12:
	r0 = r7
	goto block13
block13:
	r8 = frame.Checkpoint()
	goto block14
block14:
	r9 = frame.Read()
	if frame.Flow == 0 {
		goto block15
	} else {
		goto block24
	}
block15:
	r10 = ')'
	goto block16
block16:
	r11 = r9 == r10
	goto block17
block17:
	if r11 {
		goto block18
	} else {
		goto block22
	}
block18:
	frame.Slice(r8)
	goto block19
block19:
	r12 = r0
	goto block20
block20:
	ret0 = r12
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
	var r6 int
	var r7 rune
	var r8 rune
	var r9 bool
	var r10 []TextMatch
	var r11 TextMatch
	var r12 []TextMatch
	var r13 int
	var r14 int
	var r15 rune
	var r16 rune
	var r17 bool
	var r18 []TextMatch
	var r19 TextMatch
	var r20 []TextMatch
	var r21 []TextMatch
	var r22 *MatchChoice
	var r23 TextMatch
	goto block0
block0:
	goto block1
block1:
	r2 = Sequence(frame)
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
	r4 = r0
	goto block5
block5:
	r5 = []TextMatch{r4}
	goto block6
block6:
	r1 = r5
	goto block7
block7:
	r6 = frame.Checkpoint()
	goto block8
block8:
	r7 = frame.Read()
	if frame.Flow == 0 {
		goto block9
	} else {
		goto block36
	}
block9:
	r8 = '|'
	goto block10
block10:
	r9 = r7 == r8
	goto block11
block11:
	if r9 {
		goto block12
	} else {
		goto block35
	}
block12:
	frame.Slice(r6)
	goto block13
block13:
	S(frame)
	if frame.Flow == 0 {
		goto block14
	} else {
		goto block36
	}
block14:
	r10 = r1
	goto block15
block15:
	r11 = Sequence(frame)
	if frame.Flow == 0 {
		goto block16
	} else {
		goto block36
	}
block16:
	r12 = append(r10, r11)
	goto block17
block17:
	r1 = r12
	goto block18
block18:
	r13 = frame.Checkpoint()
	goto block19
block19:
	r14 = frame.Checkpoint()
	goto block20
block20:
	r15 = frame.Read()
	if frame.Flow == 0 {
		goto block21
	} else {
		goto block31
	}
block21:
	r16 = '|'
	goto block22
block22:
	r17 = r15 == r16
	goto block23
block23:
	if r17 {
		goto block24
	} else {
		goto block30
	}
block24:
	frame.Slice(r14)
	goto block25
block25:
	S(frame)
	if frame.Flow == 0 {
		goto block26
	} else {
		goto block31
	}
block26:
	r18 = r1
	goto block27
block27:
	r19 = Sequence(frame)
	if frame.Flow == 0 {
		goto block28
	} else {
		goto block31
	}
block28:
	r20 = append(r18, r19)
	goto block29
block29:
	r1 = r20
	goto block18
block30:
	frame.Fail()
	goto block31
block31:
	frame.Recover(r13)
	goto block32
block32:
	r21 = r1
	goto block33
block33:
	r22 = &MatchChoice{Matches: r21}
	goto block34
block34:
	r0 = r22
	goto block37
block35:
	frame.Fail()
	goto block36
block36:
	frame.Recover(r3)
	goto block37
block37:
	r23 = r0
	goto block38
block38:
	ret0 = r23
	goto block39
block39:
	return
block40:
	return
}
