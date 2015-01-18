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
	Package Package_Ref
}

func (node *ExternalType) isGoType() {
}

type TypeDefType struct {
	Name    string
	Type    GoType
	Package Package_Ref
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
	Package Package_Ref
	Methods []Function_Ref
}

func (node *StructType) isGoType() {
}

type InterfaceType struct {
	Name    string
	Fields  []*Field
	Package Package_Ref
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
}

type Function_Ref uint32

const NoFunction = ^Function_Ref(0)

type Function_Scope struct {
	objects []*Function
}

type Function struct {
	Name    string
	Package Package_Ref
	Index   Function_Ref
}

type Package_Ref uint32

const NoPackage = ^Package_Ref(0)

type Package_Scope struct {
	objects []*Package
}

type Package struct {
	Path      []string
	Extern    bool
	Functions []Function_Ref
	Index     Package_Ref
}

type CoreProgram struct {
	Package_Scope  *Package_Scope
	Function_Scope *Function_Scope
}
