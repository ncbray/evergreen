package dub

func Dub_ReadEscaped(state *DubState) rune {
	r := state.Read()
	if state.Flow == FAIL {
		return r
	}
	switch r {
	case 't':
		r = '\t'
	case 'r':
		r = '\r'
	case 'n':
		r = '\n'
	}
	return r
}

func Dub_RuneMatch(state *DubState) *RuneMatch {
	r := state.Read()
	if state.Flow == FAIL {
		return nil
	}
	if r != '[' {
		state.Fail()
		return nil
	}

	expr := &RuneMatch{Filters: []*RuneRange{}}

	r = state.Read()
	if state.Flow == FAIL {
		return nil
	}
	if r == '^' {
		expr.Invert = true
		goto Read
	} else {
		goto Decode
	}

Read:
	r = state.Read()
	if state.Flow == FAIL {
		return nil
	}
Decode:
	if r == ']' {
		return expr
	}
	if r == '\\' {
		r = Dub_ReadEscaped(state)
		if state.Flow == FAIL {
			return nil
		}
	}
	next := state.Read()
	if state.Flow == FAIL {
		return nil
	}
	if next == '-' {
		other := state.Read()
		if state.Flow == FAIL {
			return nil
		}
		if other == ']' {
			state.Fail()
			return nil
		}
		if other == '\\' {
			other = Dub_ReadEscaped(state)
			if state.Flow == FAIL {
				return nil
			}
		}
		expr.Filters = append(expr.Filters, &RuneRange{Lower: r, Upper: other})
		goto Read
	}
	expr.Filters = append(expr.Filters, &RuneRange{Lower: r, Upper: r})
	r = next
	goto Decode
}

func Dub_RunePostfix(state *DubState) Expr {
	result := Dub_RuneMatch(state)
	if state.Flow == FAIL {
		return nil
	}
	pos := state.Checkpoint()
	next := state.Read()
	if state.Flow == FAIL {
		goto NoPostfix
	}
	if next == '*' {
		return &Repeat{Expr: result, Min: 0}
	} else if next == '+' {
		return &Repeat{Expr: result, Min: 1}
	}

NoPostfix:
	state.Recover(pos)
	return result
}
