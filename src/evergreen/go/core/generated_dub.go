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
	Int    *ExternalType
	UInt32 *ExternalType
	Int64  *ExternalType
	Bool   *ExternalType
	String *ExternalType
	Rune   *ExternalType
}

type Package struct {
	Path   []string
	Extern bool
}
