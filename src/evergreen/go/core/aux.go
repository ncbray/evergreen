package core

func MakeBuiltinTypeIndex() *BuiltinTypeIndex {
	return &BuiltinTypeIndex{
		Int:     &ExternalType{Name: "int"},
		UInt32:  &ExternalType{Name: "uint32"},
		Int64:   &ExternalType{Name: "int64"},
		Float32: &ExternalType{Name: "float32"},
		Bool:    &ExternalType{Name: "bool"},
		String:  &ExternalType{Name: "string"},
		Rune:    &ExternalType{Name: "rune"},

		Append: &IntrinsicFunction{Name: "append"},
	}
}

func (scope *Package_Scope) Get(ref Package_Ref) *Package {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *Package_Scope) Register(info *Package) *Package {
	info.Index = Package_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
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
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *Function_Scope) Register(info *Function) *Function {
	info.Index = Function_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
}

func (scope *Function_Scope) Len() int {
	return len(scope.objects)
}

func InsertFunctionIntoPackage(coreProg *CoreProgram, p *Package, f *Function) {
	f.Package = p
	p.Functions = append(p.Functions, f)
}
