package core

func MakeBuiltinTypeIndex() *BuiltinTypeIndex {
	return &BuiltinTypeIndex{
		Int:     &ExternalType{Name: "int", Package: NoPackage},
		UInt32:  &ExternalType{Name: "uint32", Package: NoPackage},
		Int64:   &ExternalType{Name: "int64", Package: NoPackage},
		Float32: &ExternalType{Name: "float32", Package: NoPackage},
		Bool:    &ExternalType{Name: "bool", Package: NoPackage},
		String:  &ExternalType{Name: "string", Package: NoPackage},
		Rune:    &ExternalType{Name: "rune", Package: NoPackage},
	}
}

func (scope *Package_Scope) Get(ref Package_Ref) *Package {
	return scope.objects[ref]
}

func (scope *Package_Scope) Register(info *Package) Package_Ref {
	index := Package_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *Package_Scope) Len() int {
	return len(scope.objects)
}

func (scope *Package_Scope) Iter() *packageIterator {
	return &packageIterator{scope: scope, current: -1}
}

type packageIterator struct {
	scope   *Package_Scope
	current int
}

func (iter *packageIterator) Next() bool {
	iter.current += 1
	return iter.current < len(iter.scope.objects)
}

func (iter *packageIterator) Value() (Package_Ref, *Package) {
	return Package_Ref(iter.current), iter.scope.objects[iter.current]
}

func (scope *Function_Scope) Get(ref Function_Ref) *Function {
	return scope.objects[ref]
}

func (scope *Function_Scope) Register(info *Function) Function_Ref {
	index := Function_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *Function_Scope) Len() int {
	return len(scope.objects)
}

func InsertFunctionIntoPackage(coreProg *CoreProgram, pRef Package_Ref, fRef Function_Ref) {
	p := coreProg.Package_Scope.Get(pRef)
	f := coreProg.Function_Scope.Get(fRef)

	f.Package = pRef
	p.Functions = append(p.Functions, fRef)
}
