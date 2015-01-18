package core

type GoType interface {
	isGoType()
}

type PointerType struct {
	Element GoType
}

func (node *PointerType) isGoType() {
}

type SliceType struct {
	Element GoType
}

func (node *SliceType) isGoType() {
}

type ExternalType struct {
	Name    string
	Package *Package
}

func (node *ExternalType) isGoType() {
}

type TypeDefType struct {
	Name    string
	Type    GoType
	Package *Package
}

func (node *TypeDefType) isGoType() {
}

type FuncType struct {
	Params  []GoType
	Results []GoType
}

func (node *FuncType) isGoType() {
}

type Field struct {
	Name string
	Type GoType
}

type StructType struct {
	Name    string
	Fields  []*Field
	Package *Package
	Methods []*Function
}

func (node *StructType) isGoType() {
}

type InterfaceType struct {
	Name    string
	Fields  []*Field
	Package *Package
}

func (node *InterfaceType) isGoType() {
}

type BuiltinTypeIndex struct {
	Int     *ExternalType
	UInt32  *ExternalType
	Int64   *ExternalType
	Float32 *ExternalType
	Bool    *ExternalType
	String  *ExternalType
	Rune    *ExternalType
	Append  *IntrinsicFunction
}

type Callable interface {
	isCallable()
}

type Function_Ref uint32

type Function_Scope struct {
	objects []*Function
}

type Function struct {
	Name    string
	Package *Package
	Index   Function_Ref
}

func (node *Function) isCallable() {
}

type IntrinsicFunction struct {
	Name string
}

func (node *IntrinsicFunction) isCallable() {
}

type Package_Ref uint32

type Package_Scope struct {
	objects []*Package
}

type Package struct {
	Path      []string
	Extern    bool
	Functions []*Function
	Index     Package_Ref
}

type CoreProgram struct {
	Package_Scope  *Package_Scope
	Function_Scope *Function_Scope
}
