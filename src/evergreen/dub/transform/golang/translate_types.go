package golang

import (
	"evergreen/dub/core"
	"evergreen/dub/flow"
	dstcore "evergreen/go/core"
	ast "evergreen/go/tree"
)

const (
	STRUCT = iota
	REF
	SCOPE
)

type DubToGoLinker interface {
	SetType(s *core.StructType, subtype int, impl dstcore.GoType) dstcore.GoType
	GetType(s *core.StructType, subtype int) dstcore.GoType
	TypeRef(s *core.StructType, subtype int) *ast.NameRef
}

type linkerImpl struct {
	types []map[*core.StructType]dstcore.GoType
}

func (l *linkerImpl) SetType(s *core.StructType, subtype int, impl dstcore.GoType) dstcore.GoType {
	_, ok := l.types[subtype][s]
	if ok {
		panic(s)
	}
	l.types[subtype][s] = impl
	return impl
}

func (l *linkerImpl) GetType(s *core.StructType, subtype int) dstcore.GoType {
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
	types := []map[*core.StructType]dstcore.GoType{}
	for i := 0; i < 3; i++ {
		types = append(types, map[*core.StructType]dstcore.GoType{})
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

func builtinType(t *core.BuiltinType, ctx *DubToGoContext) dstcore.GoType {
	switch t.Name {
	case "bool":
		return ctx.index.Bool
	case "int":
		return ctx.index.Int
	case "uint32":
		return ctx.index.UInt32
	case "int64":
		return ctx.index.Int64
	case "float32":
		return ctx.index.Float32
	case "rune":
		return ctx.index.Rune
	case "string":
		return ctx.index.String
	case "graph":
		return &dstcore.PointerType{Element: ctx.graph}
	default:
		panic(t.Name)
	}
}

func goSliceType(t *core.ListType, ctx *DubToGoContext) *dstcore.SliceType {
	return &dstcore.SliceType{Element: goType(t.Type, ctx)}
}

func goType(t core.DubType, ctx *DubToGoContext) dstcore.GoType {
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
			return &dstcore.PointerType{Element: out}
		}
	default:
		panic(t)
	}
}

func goFieldType(t core.DubType, ctx *DubToGoContext) dstcore.GoType {
	switch t := t.(type) {
	case *core.StructType:
		if t.Scoped {
			return ctx.link.GetType(t, REF)
		}
	case *core.ListType:
		return &dstcore.SliceType{Element: goFieldType(t.Type, ctx)}
	}
	return goType(t, ctx)
}

func createTypeMapping(program *flow.DubProgram, coreProg *core.CoreProgram, packages []dstcore.Package_Ref, link DubToGoLinker) []dstcore.GoType {
	types := []dstcore.GoType{}
	for _, s := range coreProg.Structures {
		pIndex := coreProg.File_Scope.Get(s.File).Package
		p := packages[pIndex]

		if s.IsParent {
			if s.Scoped {
				panic(s.Name)
			}
			if len(s.Fields) != 0 {
				panic(s.Name)
			}
			types = append(types, link.SetType(s, STRUCT, &dstcore.InterfaceType{Package: p}))
		} else {
			if s.Scoped {
				types = append(types, link.SetType(s, REF, &dstcore.TypeDefType{Package: p}))
				types = append(types, link.SetType(s, SCOPE, &dstcore.StructType{Package: p}))
			}
			types = append(types, link.SetType(s, STRUCT, &dstcore.StructType{Package: p}))
		}
	}
	return types
}

func createTypes(program *flow.DubProgram, coreProg *core.CoreProgram, ctx *DubToGoContext) {
	for _, s := range coreProg.Structures {
		if s.IsParent {
			impl, _ := ctx.link.GetType(s, STRUCT).(*dstcore.InterfaceType)
			impl.Name = s.Name
			impl.Fields = []*dstcore.Field{}
			for tag := s; tag != nil; tag = tag.Implements {
				impl.Fields = append(impl.Fields, &dstcore.Field{
					Name: tagName(tag),
					Type: &dstcore.FuncType{},
				})
			}

		} else {
			impl, _ := ctx.link.GetType(s, STRUCT).(*dstcore.StructType)
			impl.Name = s.Name

			fields := []*dstcore.Field{}
			for _, f := range s.Fields {
				fields = append(fields, &dstcore.Field{
					Name: f.Name,
					Type: goFieldType(f.Type, ctx),
				})
			}
			for _, c := range s.Contains {
				if !c.Scoped {
					panic(c)
				}
				fields = append(fields, &dstcore.Field{
					Name: subtypeName(c, SCOPE),
					Type: &dstcore.PointerType{
						Element: ctx.link.GetType(c, SCOPE),
					},
				})
			}

			if s.Scoped {
				fields = append(fields, &dstcore.Field{
					Name: "Index",
					Type: ctx.link.GetType(s, REF),
				})

				ref, _ := ctx.link.GetType(s, REF).(*dstcore.TypeDefType)
				ref.Name = subtypeName(s, REF)
				ref.Type = ctx.index.UInt32

				scope, _ := ctx.link.GetType(s, SCOPE).(*dstcore.StructType)
				scope.Name = subtypeName(s, SCOPE)
				scope.Fields = []*dstcore.Field{
					&dstcore.Field{
						Name: "objects",
						Type: &dstcore.SliceType{
							Element: &dstcore.PointerType{
								Element: impl,
							},
						},
					},
				}
			}
			impl.Fields = fields
		}
	}
}
