package flow

import (
	"evergreen/dub/core"
	"evergreen/dub/tree"
	"evergreen/graph"
)

type RegisterInfo_Ref uint32

type RegisterInfo_Scope struct {
	objects []*RegisterInfo
}

type RegisterInfo struct {
	Name  string
	T     core.DubType
	Index RegisterInfo_Ref
}

type LLFunc struct {
	Name               string
	Params             []*RegisterInfo
	ReturnTypes        []core.DubType
	CFG                *graph.Graph
	Ops                []DubOp
	Edges              []int
	F                  *core.Function
	RegisterInfo_Scope *RegisterInfo_Scope
}

type DubOp interface {
	isDubOp()
}

type CoerceOp struct {
	Src *RegisterInfo
	T   core.DubType
	Dst *RegisterInfo
}

func (node *CoerceOp) isDubOp() {
}

type CopyOp struct {
	Src *RegisterInfo
	Dst *RegisterInfo
}

func (node *CopyOp) isDubOp() {
}

type ConstantNilOp struct {
	Dst *RegisterInfo
}

func (node *ConstantNilOp) isDubOp() {
}

type ConstantIntOp struct {
	Value int64
	Dst   *RegisterInfo
}

func (node *ConstantIntOp) isDubOp() {
}

type ConstantFloat32Op struct {
	Value float32
	Dst   *RegisterInfo
}

func (node *ConstantFloat32Op) isDubOp() {
}

type ConstantBoolOp struct {
	Value bool
	Dst   *RegisterInfo
}

func (node *ConstantBoolOp) isDubOp() {
}

type ConstantRuneOp struct {
	Value rune
	Dst   *RegisterInfo
}

func (node *ConstantRuneOp) isDubOp() {
}

type ConstantStringOp struct {
	Value string
	Dst   *RegisterInfo
}

func (node *ConstantStringOp) isDubOp() {
}

type BinaryOp struct {
	Left  *RegisterInfo
	Op    string
	Right *RegisterInfo
	Dst   *RegisterInfo
}

func (node *BinaryOp) isDubOp() {
}

type CallOp struct {
	Target core.Callable
	Args   []*RegisterInfo
	Dsts   []*RegisterInfo
}

func (node *CallOp) isDubOp() {
}

type KeyValue struct {
	Key   string
	Value *RegisterInfo
}

type ConstructOp struct {
	Type *core.StructType
	Args []*KeyValue
	Dst  *RegisterInfo
}

func (node *ConstructOp) isDubOp() {
}

type ConstructListOp struct {
	Type *core.ListType
	Args []*RegisterInfo
	Dst  *RegisterInfo
}

func (node *ConstructListOp) isDubOp() {
}

type Checkpoint struct {
	Dst *RegisterInfo
}

func (node *Checkpoint) isDubOp() {
}

type Recover struct {
	Src *RegisterInfo
}

func (node *Recover) isDubOp() {
}

type LookaheadBegin struct {
	Dst *RegisterInfo
}

func (node *LookaheadBegin) isDubOp() {
}

type LookaheadEnd struct {
	Failed bool
	Src    *RegisterInfo
}

func (node *LookaheadEnd) isDubOp() {
}

type Slice struct {
	Src *RegisterInfo
	Dst *RegisterInfo
}

func (node *Slice) isDubOp() {
}

type AppendOp struct {
	List  *RegisterInfo
	Value *RegisterInfo
	Dst   *RegisterInfo
}

func (node *AppendOp) isDubOp() {
}

type ReturnOp struct {
	Exprs []*RegisterInfo
}

func (node *ReturnOp) isDubOp() {
}

type Fail struct {
}

func (node *Fail) isDubOp() {
}

type Peek struct {
	Dst *RegisterInfo
}

func (node *Peek) isDubOp() {
}

type Consume struct {
}

func (node *Consume) isDubOp() {
}

type TransferOp struct {
	Srcs []*RegisterInfo
	Dsts []*RegisterInfo
}

func (node *TransferOp) isDubOp() {
}

type EntryOp struct {
}

func (node *EntryOp) isDubOp() {
}

type SwitchOp struct {
	Cond *RegisterInfo
}

func (node *SwitchOp) isDubOp() {
}

type ExitOp struct {
}

func (node *ExitOp) isDubOp() {
}

type DubPackage struct {
	Path    []string
	Structs []*core.StructType
	Funcs   []*LLFunc
	Tests   []*tree.Test
}

type DubProgram struct {
	Core     *core.CoreProgram
	Packages []*DubPackage
	LLFuncs  []*LLFunc
}
