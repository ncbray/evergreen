package dub

type Expr interface {
  isExpr()
}

type RuneRange struct {
  Lower rune
  Upper rune
}

type RuneMatch struct {
  Invert bool
  Filters []*RuneRange
}

func (node *RuneMatch) isExpr() {
}

type Repeat struct {
	Expr Expr
	Min int
}

func (node *Repeat) isExpr() {
}

type FuncDecl struct {
	Name string
}
