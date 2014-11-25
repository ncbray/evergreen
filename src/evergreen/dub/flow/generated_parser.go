package flow

import (
	"evergreen/base"
	"evergreen/dub/tree"
)

type RegisterInfo_Ref uint32

const NoRegisterInfo = ^RegisterInfo_Ref(0)

type RegisterInfo_Scope struct {
	objects []*RegisterInfo
}

type RegisterInfo struct {
	T tree.DubType
}

type LLFunc struct {
	Name               string
	Params             []RegisterInfo_Ref
	ReturnTypes        []tree.DubType
	CFG                *base.Graph
	Ops                []DubOp
	RegisterInfo_Scope *RegisterInfo_Scope
}

type DubOp interface {
	isDubOp()
}

type CoerceOp struct {
	Src RegisterInfo_Ref
	T   tree.DubType
	Dst RegisterInfo_Ref
}

func (node *CoerceOp) isDubOp() {
}

type CopyOp struct {
	Src RegisterInfo_Ref
	Dst RegisterInfo_Ref
}

func (node *CopyOp) isDubOp() {
}

type ConstantNilOp struct {
	Dst RegisterInfo_Ref
}

func (node *ConstantNilOp) isDubOp() {
}

type ConstantIntOp struct {
	Value int64
	Dst   RegisterInfo_Ref
}

func (node *ConstantIntOp) isDubOp() {
}

type ConstantBoolOp struct {
	Value bool
	Dst   RegisterInfo_Ref
}

func (node *ConstantBoolOp) isDubOp() {
}

type ConstantRuneOp struct {
	Value rune
	Dst   RegisterInfo_Ref
}

func (node *ConstantRuneOp) isDubOp() {
}

type ConstantStringOp struct {
	Value string
	Dst   RegisterInfo_Ref
}

func (node *ConstantStringOp) isDubOp() {
}

type BinaryOp struct {
	Left  RegisterInfo_Ref
	Op    string
	Right RegisterInfo_Ref
	Dst   RegisterInfo_Ref
}

func (node *BinaryOp) isDubOp() {
}

type CallOp struct {
	Name   string
	Target *LLFunc
	Args   []RegisterInfo_Ref
	Dsts   []RegisterInfo_Ref
}

func (node *CallOp) isDubOp() {
}

type KeyValue struct {
	Key   string
	Value RegisterInfo_Ref
}

type ConstructOp struct {
	Type *tree.StructType
	Args []*KeyValue
	Dst  RegisterInfo_Ref
}

func (node *ConstructOp) isDubOp() {
}

type ConstructListOp struct {
	Type *tree.ListType
	Args []RegisterInfo_Ref
	Dst  RegisterInfo_Ref
}

func (node *ConstructListOp) isDubOp() {
}

type Checkpoint struct {
	Dst RegisterInfo_Ref
}

func (node *Checkpoint) isDubOp() {
}

type Recover struct {
	Src RegisterInfo_Ref
}

func (node *Recover) isDubOp() {
}

type LookaheadBegin struct {
	Dst RegisterInfo_Ref
}

func (node *LookaheadBegin) isDubOp() {
}

type LookaheadEnd struct {
	Failed bool
	Src    RegisterInfo_Ref
}

func (node *LookaheadEnd) isDubOp() {
}

type Slice struct {
	Src RegisterInfo_Ref
	Dst RegisterInfo_Ref
}

func (node *Slice) isDubOp() {
}

type AppendOp struct {
	List  RegisterInfo_Ref
	Value RegisterInfo_Ref
	Dst   RegisterInfo_Ref
}

func (node *AppendOp) isDubOp() {
}

type ReturnOp struct {
	Exprs []RegisterInfo_Ref
}

func (node *ReturnOp) isDubOp() {
}

type Fail struct {
}

func (node *Fail) isDubOp() {
}

type Peek struct {
	Dst RegisterInfo_Ref
}

func (node *Peek) isDubOp() {
}

type Consume struct {
}

func (node *Consume) isDubOp() {
}

type TransferOp struct {
	Srcs []RegisterInfo_Ref
	Dsts []RegisterInfo_Ref
}

func (node *TransferOp) isDubOp() {
}

type EntryOp struct {
}

func (node *EntryOp) isDubOp() {
}

type SwitchOp struct {
	Cond RegisterInfo_Ref
}

func (node *SwitchOp) isDubOp() {
}

type FlowExitOp struct {
	Flow int
}

func (node *FlowExitOp) isDubOp() {
}

type ExitOp struct {
}

func (node *ExitOp) isDubOp() {
}
