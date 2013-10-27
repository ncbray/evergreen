package dub

func Dub_ReadEscaped(state *DubState) rune {
	r := state.Read()
	if state.flow == FAIL {
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
	return r;
}

func Dub_RuneMatch(state *DubState) *RuneMatch {
	r := state.Read()
	if state.flow == FAIL {
		return nil
	}
	if r != '[' {
		state.Reject()
		return nil
	}

  expr := &RuneMatch{Filters: []*RuneRange{}}

  r = state.Read()
  if state.flow == FAIL {
		return nil
	}
  if r == '^' {
		expr.Invert = true;
		goto Read
	} else {
		goto Decode
	}

Read:
  r = state.Read()
  if state.flow == FAIL {
		return nil
	}
Decode:
  if r == ']' {
		return expr
	}
  if r == '\\' {
		r = Dub_ReadEscaped(state)
		if state.flow == FAIL {
			return nil
		}
	}
	next := state.Read()
	if state.flow == FAIL {
		return nil
	}
	if next == '-' {
		other := state.Read()
		if state.flow == FAIL {
			return nil
		}
		if other == ']' {
			state.Reject()
			return nil
		}
		if other == '\\' {
			other = Dub_ReadEscaped(state)
			if state.flow == FAIL {
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
	result := Dub_RuneMatch(state);
	if state.flow == FAIL {
		return nil
	}
	pos := state.Save()
	next := state.Read()
	if state.flow == FAIL {
		goto NoPostfix
	}
	if next == '*' {
		return &Repeat{Expr: result, Min: 0}
	} else if next == '+' {
		return &Repeat{Expr: result, Min: 1}
	}

NoPostfix:
	state.Restore(pos)
	return result;
}
