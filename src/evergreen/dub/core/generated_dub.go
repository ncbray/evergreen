package core

type DubType interface {
	isDubType()
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) isDubType() {
}

type NilType struct {
}

func (node *NilType) isDubType() {
}

type ListType struct {
	Type DubType
}

func (node *ListType) isDubType() {
}

type TupleType struct {
	Types []DubType
}

func (node *TupleType) isDubType() {
}

type FunctionType struct {
	Params []DubType
	Result DubType
}

func (node *FunctionType) isDubType() {
}

type UnboundType struct {
	Index int
}

func (node *UnboundType) isDubType() {
}

type FunctionTemplateType struct {
}

func (node *FunctionTemplateType) isDubType() {
}

type PackageType struct {
}

func (node *PackageType) isDubType() {
}

type FieldType struct {
	Name string
	Type DubType
}

type StructType struct {
	Name       string
	Implements *StructType
	Fields     []*FieldType
	Scoped     bool
	Contains   []*StructType
	IsParent   bool
	File       *File
}

func (node *StructType) isDubType() {
}

type Package_Ref uint32

type Package_Scope struct {
	objects []*Package
}

type Package struct {
	Path  []string
	Files []*File
	Index Package_Ref
}

type File_Ref uint32

type File_Scope struct {
	objects []*File
}

type File struct {
	Name    string
	Package *Package
	Index   File_Ref
}

type Callable interface {
	isCallable()
}

type Function_Ref uint32

type Function_Scope struct {
	objects []*Function
}

type Function struct {
	Name  string
	Type  *FunctionType
	File  *File
	Index Function_Ref
}

func (node *Function) isCallable() {
}

type IntrinsicFunction struct {
	Name   string
	Parent *IntrinsicFunctionTemplate
	Type   *FunctionType
}

func (node *IntrinsicFunction) isCallable() {
}

type TemplateParam struct {
	Name string
}

type CallableTemplate interface {
	isCallableTemplate()
}

type FunctionTemplate struct {
	Name string
}

func (node *FunctionTemplate) isCallableTemplate() {
}

type IntrinsicFunctionTemplate struct {
	Name   string
	Params []*TemplateParam
	Type   *FunctionType
}

func (node *IntrinsicFunctionTemplate) isCallableTemplate() {
}

type BuiltinTypeIndex struct {
	String   *BuiltinType
	Rune     *BuiltinType
	Int      *BuiltinType
	Int64    *BuiltinType
	Float32  *BuiltinType
	Bool     *BuiltinType
	Graph    *BuiltinType
	Nil      *NilType
	Append   *IntrinsicFunctionTemplate
	Position *IntrinsicFunction
	Slice    *IntrinsicFunction
}

type CoreProgram struct {
	Builtins       *BuiltinTypeIndex
	Structures     []*StructType
	Package_Scope  *Package_Scope
	File_Scope     *File_Scope
	Function_Scope *Function_Scope
}
