package core

func TypeName(t DubType) string {
	switch t := t.(type) {
	case *StructType:
		return t.Name
	case *ListType:
		return "[]" + TypeName(t.Type)
	case *BuiltinType:
		return t.Name
	default:
		panic(t)
	}
}

func (scope *Function_Scope) Get(ref Function_Ref) *Function {
	if scope.objects[ref].Index != ref {
		panic(ref)
	}
	return scope.objects[ref]
}

func (scope *Function_Scope) Register(info *Function) Function_Ref {
	info.Index = Function_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info.Index
}

func (scope *Function_Scope) Len() int {
	return len(scope.objects)
}

func (scope *Function_Scope) Iter() *functionIterator {
	return &functionIterator{scope: scope, current: -1}
}

type functionIterator struct {
	scope   *Function_Scope
	current int
}

func (iter *functionIterator) Next() bool {
	iter.current += 1
	return iter.current < len(iter.scope.objects)
}

func (iter *functionIterator) Value() (Function_Ref, *Function) {
	return Function_Ref(iter.current), iter.scope.objects[iter.current]
}

func (scope *Package_Scope) Get(ref Package_Ref) *Package {
	if scope.objects[ref].Index != ref {
		panic(ref)
	}
	return scope.objects[ref]
}

func (scope *Package_Scope) Register(info *Package) Package_Ref {
	info.Index = Package_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info.Index
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

func (scope *File_Scope) Get(ref File_Ref) *File {
	if scope.objects[ref].Index != ref {
		panic(ref)
	}
	return scope.objects[ref]
}

func (scope *File_Scope) Register(info *File) File_Ref {
	info.Index = File_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info.Index
}

func (scope *File_Scope) Len() int {
	return len(scope.objects)
}

func (scope *File_Scope) Iter() *fileIterator {
	return &fileIterator{scope: scope, current: -1}
}

type fileIterator struct {
	scope   *File_Scope
	current int
}

func (iter *fileIterator) Next() bool {
	iter.current += 1
	return iter.current < len(iter.scope.objects)
}

func (iter *fileIterator) Value() (File_Ref, *File) {
	return File_Ref(iter.current), iter.scope.objects[iter.current]
}
