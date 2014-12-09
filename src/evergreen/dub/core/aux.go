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
