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
