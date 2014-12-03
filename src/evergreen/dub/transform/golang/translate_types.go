package golang

import (
	"evergreen/dub/core"
	"evergreen/dub/flow"
	ast "evergreen/go/tree"
)

const (
	STRUCT = iota
	REF
	SCOPE
)

type DubToGoLinker interface {
	SetType(s *core.StructType, subtype int, impl ast.GoType)
	GetType(s *core.StructType, subtype int) ast.GoType
	TypeRef(s *core.StructType, subtype int) *ast.NameRef
}

type linkerImpl struct {
	types []map[*core.StructType]ast.GoType
}

func (l *linkerImpl) SetType(s *core.StructType, subtype int, impl ast.GoType) {
	_, ok := l.types[subtype][s]
	if ok {
		panic(s)
	}
	l.types[subtype][s] = impl
}

func (l *linkerImpl) GetType(s *core.StructType, subtype int) ast.GoType {
	e, ok := l.types[subtype][s]
	if !ok {
		panic(s)
	}
	return e
}

func (l *linkerImpl) TypeRef(s *core.StructType, subtype int) *ast.NameRef {
	return &ast.NameRef{
		Name: subtypeName(s, subtype),
		T:    l.GetType(s, subtype),
	}
}

func makeLinker() DubToGoLinker {
	types := []map[*core.StructType]ast.GoType{}
	for i := 0; i < 3; i++ {
		types = append(types, map[*core.StructType]ast.GoType{})
	}
	return &linkerImpl{
		types: types,
	}
}

func subtypeName(s *core.StructType, subtype int) string {
	name := s.Name
	switch subtype {
	case STRUCT:
		// Nothing
	case REF:
		name += "_Ref"
	case SCOPE:
		name += "_Scope"
	default:
		panic(subtype)
	}
	return name
}

func tagName(s *core.StructType) string {
	return "is" + s.Name
}

func makeBuiltinTypes() *ast.BuiltinTypeIndex {
	return &ast.BuiltinTypeIndex{
		Int:    &ast.ExternalType{Name: "int"},
		UInt32: &ast.ExternalType{Name: "uint32"},
		Int64:  &ast.ExternalType{Name: "int64"},
		Bool:   &ast.ExternalType{Name: "bool"},
		String: &ast.ExternalType{Name: "string"},
		Rune:   &ast.ExternalType{Name: "rune"},
	}
}

func builtinType(t *core.BuiltinType, ctx *DubToGoContext) ast.GoType {
	switch t.Name {
	case "bool":
		return ctx.index.Bool
	case "int":
		return ctx.index.Int
	case "uint32":
		return ctx.index.UInt32
	case "int64":
		return ctx.index.Int64
	case "rune":
		return ctx.index.Rune
	case "string":
		return ctx.index.String
	case "graph":
		return &ast.PointerType{Element: ctx.graph}
	default:
		panic(t.Name)
	}
}

func goSliceType(t *core.ListType, ctx *DubToGoContext) *ast.SliceType {
	return &ast.SliceType{Element: goType(t.Type, ctx)}
}

func goType(t core.DubType, ctx *DubToGoContext) ast.GoType {
	switch t := t.(type) {
	case *core.BuiltinType:
		return builtinType(t, ctx)
	case *core.ListType:
		return goSliceType(t, ctx)
	case *core.StructType:
		out := ctx.link.GetType(t, STRUCT)
		if t.IsParent {
			return out
		} else {
			return &ast.PointerType{Element: out}
		}
	default:
		panic(t)
	}
}

func goFieldType(t core.DubType, ctx *DubToGoContext) ast.GoType {
	switch t := t.(type) {
	case *core.StructType:
		if t.Scoped {
			return ctx.link.GetType(t, REF)
		}
	case *core.ListType:
		return &ast.SliceType{Element: goFieldType(t.Type, ctx)}
	}
	return goType(t, ctx)
}

func createTypeMapping(program []*flow.DubPackage, link DubToGoLinker) {
	for _, dubPkg := range program {
		for _, s := range dubPkg.Structs {
			if s.IsParent {
				link.SetType(s, STRUCT, &ast.InterfaceType{})
			} else {
				if s.Scoped {
					link.SetType(s, REF, &ast.TypeDefType{})
					link.SetType(s, SCOPE, &ast.StructType{})
				}
				link.SetType(s, STRUCT, &ast.StructType{})
			}
		}
	}
}

func createTypes(program []*flow.DubPackage, ctx *DubToGoContext) {
	for _, dubPkg := range program {
		for _, s := range dubPkg.Structs {
			if s.IsParent {
				impl, _ := ctx.link.GetType(s, STRUCT).(*ast.InterfaceType)
				impl.Name = s.Name
				impl.Fields = []*ast.Field{}
				for tag := s; tag != nil; tag = tag.Implements {
					impl.Fields = append(impl.Fields, &ast.Field{
						Name: tagName(tag),
						Type: &ast.FuncType{},
					})
				}

			} else {
				impl, _ := ctx.link.GetType(s, STRUCT).(*ast.StructType)
				impl.Name = s.Name

				fields := []*ast.Field{}
				for _, f := range s.Fields {
					fields = append(fields, &ast.Field{
						Name: f.Name,
						Type: goFieldType(f.Type, ctx),
					})
				}
				for _, c := range s.Contains {
					if !c.Scoped {
						panic(c)
					}
					fields = append(fields, &ast.Field{
						Name: subtypeName(c, SCOPE),
						Type: &ast.PointerType{
							Element: ctx.link.GetType(c, SCOPE),
						},
					})
				}
				impl.Fields = fields

				if s.Scoped {
					ref, _ := ctx.link.GetType(s, REF).(*ast.TypeDefType)
					ref.Name = subtypeName(s, REF)
					ref.Type = ctx.index.UInt32

					scope, _ := ctx.link.GetType(s, SCOPE).(*ast.StructType)
					scope.Name = subtypeName(s, SCOPE)
					scope.Fields = []*ast.Field{
						&ast.Field{
							Name: "objects",
							Type: &ast.SliceType{
								Element: &ast.PointerType{
									Element: impl,
								},
							},
						},
					}
				}
			}

		}
	}
}
