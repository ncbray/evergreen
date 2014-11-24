package tree

// Is there a change that execution can flow through this block without
// jumping, returning, excepting, etc?
func NormalFlowMightExit(stmts []Stmt) bool {
	n := len(stmts)
	if n == 0 {
		return true
	}
	last := stmts[n-1]
	switch last := last.(type) {
	case *Goto, *Return:
		return false
	default:
		panic(last)
	}
}
