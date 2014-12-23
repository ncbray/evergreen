package core

func MakeBuiltinTypeIndex() *BuiltinTypeIndex {
	return &BuiltinTypeIndex{
		Int:    &ExternalType{Name: "int", Package: NoPackage},
		UInt32: &ExternalType{Name: "uint32", Package: NoPackage},
		Int64:  &ExternalType{Name: "int64", Package: NoPackage},
		Bool:   &ExternalType{Name: "bool", Package: NoPackage},
		String: &ExternalType{Name: "string", Package: NoPackage},
		Rune:   &ExternalType{Name: "rune", Package: NoPackage},
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
