package tree

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"fmt"
	"strings"
)

type tupleLUT struct {
	Type *core.TupleType
	Next map[core.DubType]*tupleLUT
}

func makeTupleLUT() *tupleLUT {
	return &tupleLUT{Next: map[core.DubType]*tupleLUT{}}
}

type specialization struct {
	Func   *core.IntrinsicFunctionTemplate
	Params *core.TupleType
}

type funcTypeKey struct {
	Params *core.TupleType
	Result core.DubType
}

type typeMemoizer struct {
	Tuples      *tupleLUT
	Lists       map[core.DubType]*core.ListType
	Specialized map[specialization]core.Callable
	Funcs       map[funcTypeKey]*core.FunctionType
	Unbound     []*core.UnboundType
}

func (memo *typeMemoizer) getTuple(types []core.DubType) *core.TupleType {
	current := memo.Tuples
	for _, t := range types {
		next, ok := current.Next[t]
		if !ok {
			next = makeTupleLUT()
			current.Next[t] = next
		}
		current = next
	}
	if current.Type == nil {
		current.Type = &core.TupleType{Types: types}
	}
	return current.Type
}

func (memo *typeMemoizer) getList(t core.DubType) *core.ListType {
	lt, ok := memo.Lists[t]
	if !ok {
		lt = &core.ListType{Type: t}
		memo.Lists[t] = lt
	}
	return lt
}

// TODO remove ft parameter?
func (memo *typeMemoizer) getSpecialized(template *core.IntrinsicFunctionTemplate, params []core.DubType, ft *core.FunctionType) core.Callable {
	key := specialization{Func: template, Params: memo.getTuple(params)}
	c, ok := memo.Specialized[key]
	if !ok {
		c = &core.IntrinsicFunction{
			Name:   template.Name,
			Parent: template,
			Type:   ft,
		}
		memo.Specialized[key] = c
	}
	return c
}

func (memo *typeMemoizer) getFunctionType(params []core.DubType, result core.DubType) *core.FunctionType {
	key := funcTypeKey{Params: memo.getTuple(params), Result: result}
	ft, ok := memo.Funcs[key]
	if !ok {
		ft = &core.FunctionType{
			Params: params,
			Result: result,
		}
		memo.Funcs[key] = ft
	}
	return ft
}

func (memo *typeMemoizer) getUnbound(index int) *core.UnboundType {
	for index >= len(memo.Unbound) {
		memo.Unbound = append(memo.Unbound, &core.UnboundType{Index: len(memo.Unbound)})
	}
	return memo.Unbound[index]
}

type semanticScope struct {
	parent *semanticScope
	locals map[string]*LocalInfo
}

func (scope *semanticScope) localInfo(name string) (*LocalInfo, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return nil, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{
		parent: scope,
		locals: map[string]*LocalInfo{},
	}
}

var unresolvedType core.DubType = nil

func TypeMatches(actual core.DubType, expected core.DubType, exact bool) bool {
	if actual == unresolvedType || expected == unresolvedType {
		return true
	}

	switch actual := actual.(type) {
	case *core.StructType:
		other, ok := expected.(*core.StructType)
		if !ok {
			return false
		}
		if exact {
			return actual == other
		}
		current := actual
		for current != nil {
			if current == other {
				return true
			}
			current = current.Implements
		}
		return false
	case *core.NilType:
		_, ok := expected.(*core.StructType)
		return ok
	case *core.BuiltinType:
		other, ok := expected.(*core.BuiltinType)
		if !ok {
			return false
		}
		return actual.Name == other.Name
	case *core.ListType:
		other, ok := expected.(*core.ListType)
		if !ok {
			return false
		}
		return TypeMatches(actual.Type, other.Type, true)
	default:
		panic(actual)
	}
}

func IsDiscard(name string) bool {
	return name == "_"
}

func createLocal(ctx *semanticPassContext, decl *FuncDecl, name *Id, t core.DubType, scope *semanticScope) *LocalInfo {
	info, exists := scope.localInfo(name.Text)
	if exists {
		ctx.Status.LocationError(name.Pos, fmt.Sprintf("Tried to redefine %#v", name.Text))
		return info
	}
	info = decl.LocalInfo_Scope.Register(&LocalInfo{Name: name.Text, T: t})
	scope.locals[name.Text] = info
	return info
}

func semanticTargetPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, t core.DubType, define bool, scope *semanticScope) ASTExpr {
	switch expr := expr.(type) {
	case *NameRef:
		name := expr.Name.Text
		if IsDiscard(name) {
			return &Discard{}
		}
		switch t.(type) {
		case *core.TupleType:
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Cannot store tuple in local %#v", name))
			t = unresolvedType
		}
		var info *LocalInfo
		if define {
			info = createLocal(ctx, decl, expr.Name, t, scope)
		} else {
			var exists bool
			info, exists = scope.localInfo(name)
			if !exists {
				ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Tried to assign to unknown variable %#v", name))
				return &Discard{}
			}
			// TODO type check
		}
		return &SetLocal{Info: info}
	default:
		panic(expr)
	}
}

// TODO rewrite binary op to explicitly promote types.
func binaryOpType(ctx *semanticPassContext, lt *core.BuiltinType, op string, rt *core.BuiltinType) (*core.BuiltinType, bool) {
	builtins := ctx.Program.Index
	switch op {
	case "+", "-", "*", "/":
		switch lt {
		case builtins.Int:
			switch rt {
			case builtins.Int:
				return builtins.Int, true
			}
		case builtins.Float32:
			switch rt {
			case builtins.Float32:
				return builtins.Float32, true
			}
		case builtins.String:
			switch rt {
			case builtins.String:
				if op == "+" {
					return builtins.String, true
				}
			}
		}
	case "<", "<=", ">", ">=", "==", "!=":
		switch lt {
		case builtins.Int:
			switch rt {
			case builtins.Int:
				return builtins.Bool, true
			}
		case builtins.Float32:
			switch rt {
			case builtins.Float32:
				return builtins.Bool, true
			}
		}
	}
	return nil, false
}

func resolve(ctx *semanticPassContext, name string) (namedElement, bool) {
	result, ok := ctx.Module.Namespace[name]
	if ok {
		return result, true
	}
	result, ok = ctx.Program.Namespace[name]
	return result, ok
}

func funcType(c core.Callable) *core.FunctionType {
	switch f := c.(type) {
	case *core.Function:
		return f.Type
	case *core.IntrinsicFunction:
		return f.Type
	default:
		panic(c)
	}
}

func inferArgBindings(ctx *semanticPassContext, at core.DubType, ct core.DubType, exact bool, bindings []core.DubType) core.DubType {
	switch at := at.(type) {
	case *core.UnboundType:
		bound := bindings[at.Index]
		if bound == nil {
			bindings[at.Index] = ct
			bound = ct
		} else {
			if !TypeMatches(ct, bound, exact) {
				ctx.Status.GlobalError(fmt.Sprintf("Expected %v, but got %v", bound, ct))
			}
		}
		return bound
	case *core.ListType:
		other, ok := ct.(*core.ListType)
		if !ok {
			ctx.Status.GlobalError(fmt.Sprintf("Expected a list, but got %#v", ct))
			return unresolvedType
		}

		return ctx.Memo.getList(inferArgBindings(ctx, at.Type, other.Type, true, bindings))
	default:
		panic(at)
	}
}

func specializeBindings(ctx *semanticPassContext, at core.DubType, bindings []core.DubType) core.DubType {
	switch at := at.(type) {
	case *core.UnboundType:
		return bindings[at.Index]
	case *core.ListType:
		return ctx.Memo.getList(specializeBindings(ctx, at.Type, bindings))
	default:
		panic(at)
	}
}

func inferTypeBindings(ctx *semanticPassContext, aft *core.FunctionType, numParams int, args []core.DubType) (*core.FunctionType, []core.DubType) {
	bindings := make([]core.DubType, numParams)
	if len(aft.Params) != len(args) {
		ctx.Status.GlobalError(fmt.Sprintf("Expected %d arguments but got %d", len(aft.Params), len(args)))
		return nil, nil
	}
	specialized := make([]core.DubType, len(aft.Params))
	for i, at := range aft.Params {
		specialized[i] = inferArgBindings(ctx, at, args[i], false, bindings)
	}
	result := specializeBindings(ctx, aft.Result, bindings)

	// TODO was there an error?  Return unresolved.
	ft := ctx.Memo.getFunctionType(specialized, result)
	return ft, bindings
}

func inferTemplate(ctx *semanticPassContext, template *core.IntrinsicFunctionTemplate, args []core.DubType) (core.Callable, *core.FunctionType) {
	ft, bindings := inferTypeBindings(ctx, template.Type, len(template.Params), args)
	return ctx.Memo.getSpecialized(template, bindings, ft), ft
}

func specializeTemplate(ctx *semanticPassContext, template *core.IntrinsicFunctionTemplate, bindings []core.DubType) (core.Callable, *core.FunctionType) {
	if len(bindings) != len(template.Params) {
		ctx.Status.GlobalError(fmt.Sprintf("Expected %d type parameters but got %d", len(template.Params), len(bindings)))
		return nil, nil
	}

	aft := template.Type
	specialized := make([]core.DubType, len(aft.Params))
	for i, at := range aft.Params {
		specialized[i] = specializeBindings(ctx, at, bindings)
	}
	result := specializeBindings(ctx, aft.Result, bindings)
	ft := ctx.Memo.getFunctionType(specialized, result)

	return ctx.Memo.getSpecialized(template, bindings, ft), ft
}

func rewriteNamedLookup(named namedElement) (ASTExpr, core.DubType) {
	switch named := named.(type) {
	case *NamedCallable:
		return &GetFunction{Func: named.Func}, funcType(named.Func)
	case *NamedCallableTemplate:
		return &GetFunctionTemplate{Template: named.Func}, &core.FunctionTemplateType{}
	case *NamedPackage:
		return &GetPackage{Package: named.Scope.Package}, &core.PackageType{}
	default:
		panic(named)
	}
}

func semanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) (ASTExpr, core.DubType) {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(ctx, decl, expr.Block, scope)
		return expr, ctx.Void
	case *Choice:
		for _, block := range expr.Blocks {
			semanticBlockPass(ctx, decl, block, childScope(scope))
		}
		return expr, ctx.Void
	case *Optional:
		semanticBlockPass(ctx, decl, expr.Block, scope)
		return expr, ctx.Void
	case *If:
		expr.Expr, _ = semanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO check condition type
		semanticBlockPass(ctx, decl, expr.Block, childScope(scope))
		semanticBlockPass(ctx, decl, expr.Else, childScope(scope))
		return expr, ctx.Void
	case *BinaryOp:
		var l, r core.DubType
		expr.Left, l = semanticExprPass(ctx, decl, expr.Left, scope)
		expr.Right, r = semanticExprPass(ctx, decl, expr.Right, scope)
		if l == unresolvedType || r == unresolvedType {
			return expr, unresolvedType
		}
		lt, ok := l.(*core.BuiltinType)
		if !ok {
			panic(l)
		}
		rt, ok := r.(*core.BuiltinType)
		if !ok {
			panic(r)
		}
		t, ok := binaryOpType(ctx, lt, expr.Op, rt)
		if ok {
			expr.T = t
		} else {
			expr.T = unresolvedType
			sig := fmt.Sprintf("%s%s%s", lt.Name, expr.Op, rt.Name)
			ctx.Status.LocationError(expr.OpPos, fmt.Sprintf("unsupported binary op %s", sig))
		}
		return expr, expr.T
	case *NameRef:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if found {
			return &GetLocal{Info: info}, info.T
		}
		named, found := resolve(ctx, name)
		if found {
			return rewriteNamedLookup(named)
		}
		ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
		return expr, unresolvedType
	case *Assign:
		var t core.DubType
		if expr.Expr != nil {
			expr.Expr, t = semanticExprPass(ctx, decl, expr.Expr, scope)
		}
		if expr.Type != nil {
			expr.Type, t = semanticTypePass(ctx, expr.Type)
		}
		if len(expr.Targets) != 1 {
			count := 1
			switch t := t.(type) {
			case *core.TupleType:
				count = len(t.Types)
				if len(expr.Targets) == count {
					for i, target := range expr.Targets {
						expr.Targets[i] = semanticTargetPass(ctx, decl, target, t.Types[i], expr.Define, scope)
					}
					return expr, t
				}
			}
			ctx.Status.LocationError(expr.Pos, fmt.Sprintf("Expected %d values but got %d", len(expr.Targets), count))
			for i, target := range expr.Targets {
				expr.Targets[i] = semanticTargetPass(ctx, decl, target, unresolvedType, expr.Define, scope)
			}
			return expr, unresolvedType
		} else {
			expr.Targets[0] = semanticTargetPass(ctx, decl, expr.Targets[0], t, expr.Define, scope)
			return expr, t
		}
	case *StringMatch:
		return expr, ctx.Program.Index.String
	case *RuneMatch:
		return expr, ctx.Program.Index.Rune
	case *RuneLiteral:
		return expr, ctx.Program.Index.Rune
	case *StringLiteral:
		return expr, ctx.Program.Index.String
	case *IntLiteral:
		return expr, ctx.Program.Index.Int
	case *Float32Literal:
		return expr, ctx.Program.Index.Float32
	case *BoolLiteral:
		return expr, ctx.Program.Index.Bool
	case *NilLiteral:
		return expr, ctx.Program.Index.Nil
	case *Return:
		if len(decl.ReturnTypes) != len(expr.Exprs) {
			ctx.Status.LocationError(expr.Pos, fmt.Sprintf("expected %d return values, got %d", len(decl.ReturnTypes), len(expr.Exprs)))
		}
		for i, e := range expr.Exprs {
			var at core.DubType
			expr.Exprs[i], at = semanticExprPass(ctx, decl, e, scope)
			if i < len(decl.ReturnTypes) {
				et := ResolveType(decl.ReturnTypes[i])
				if !TypeMatches(at, et, false) {
					// TODO point at the exact expression.
					ctx.Status.LocationError(expr.Pos, fmt.Sprintf("return: %s vs. %s", core.TypeName(at), core.TypeName(et)))
				}
			}
		}
		return expr, ctx.Void
	case *Fail:
		return expr, ctx.Void
	case *Call:
		// Process the main expr
		var et core.DubType
		expr.Expr, et = semanticExprPass(ctx, decl, expr.Expr, scope)

		// Process args
		args := make([]core.DubType, len(expr.Args))
		for i, e := range expr.Args {
			expr.Args[i], args[i] = semanticExprPass(ctx, decl, e, scope)
		}

		// Resolve the call
		var rt core.DubType = unresolvedType
		if et != unresolvedType {
			switch ft := et.(type) {
			case *core.FunctionType:
				ref, ok := expr.Expr.(*GetFunction)
				if ok {
					expr.Target = ref.Func
					if len(args) == len(ft.Params) {
						for i, at := range args {
							et := ft.Params[i]
							if !TypeMatches(at, et, false) {
								// TODO point at the exact expression.
								ctx.Status.LocationError(expr.Pos, fmt.Sprintf("argument %d - got %s, expected %s", i, core.TypeName(at), core.TypeName(et)))
							}
						}
					} else {
						ctx.Status.LocationError(expr.Pos, fmt.Sprintf("expected %d arguments, got %d", len(ft.Params), len(args)))
					}
					rt = ft.Result
				} else {
					ctx.Status.LocationError(expr.Pos, "can only call directly referenced functions")
				}
			case *core.FunctionTemplateType:
				ref, ok := expr.Expr.(*GetFunctionTemplate)
				if ok {
					switch tmpl := ref.Template.(type) {
					case *core.IntrinsicFunctionTemplate:
						concrete, cft := inferTemplate(ctx, tmpl, args)
						expr.Target = concrete
						rt = cft.Result
					default:
						panic(tmpl)
					}
				} else {
					ctx.Status.LocationError(expr.Pos, "can only call directly referenced function templates")
				}
			default:
				ctx.Status.LocationError(expr.Pos, "expr not callable")
			}
		}

		expr.T = rt
		return expr, rt
	case *Selector:
		expr.Expr, _ = semanticExprPass(ctx, decl, expr.Expr, scope)
		switch e := expr.Expr.(type) {
		case *GetPackage:
			ns := ctx.ModuleContexts[e.Package.Index].Module.Namespace
			child, ok := ns[expr.Name.Text]
			if !ok {
				ctx.Status.LocationError(expr.Pos, fmt.Sprintf("unknown name %#v", expr.Name.Text))
				return expr, unresolvedType
			}
			return rewriteNamedLookup(child)
		default:
			panic(e)
		}
	case *SpecializeTemplate:
		var bindings []core.DubType
		expr.Expr, _ = semanticExprPass(ctx, decl, expr.Expr, scope)
		expr.Types, bindings = semanticTypeListPass(ctx, expr.Types)

		switch e := expr.Expr.(type) {
		case *GetFunctionTemplate:
			switch tmpl := e.Template.(type) {
			case *core.IntrinsicFunctionTemplate:
				concrete, cft := specializeTemplate(ctx, tmpl, bindings)
				return &GetFunction{Func: concrete}, cft
			default:
				panic(tmpl)
			}
		default:
			panic(e)
		}
	case *Construct:
		var t core.DubType
		expr.Type, t = semanticTypePass(ctx, expr.Type)
		st, ok := t.(*core.StructType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			var aft core.DubType
			arg.Expr, aft = semanticExprPass(ctx, decl, arg.Expr, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := f.Type
					if !TypeMatches(aft, eft, false) {
						ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("Expected type %s, but got %s", core.TypeName(eft), core.TypeName(aft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", core.TypeName(t), fn))
				}
			}
		}
		return expr, t
	case *ConstructList:
		var t core.DubType
		expr.Type, t = semanticTypePass(ctx, expr.Type)
		lt, ok := t.(*core.ListType)
		if t != nil && !ok {
			panic(t)
		}
		for i, arg := range expr.Args {
			var at core.DubType
			expr.Args[i], at = semanticExprPass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.GlobalError(fmt.Sprintf("%s vs. %s", core.TypeName(at), core.TypeName(lt.Type)))
				}
			}
		}
		return expr, t
	case *Coerce:
		var t core.DubType
		expr.Type, t = semanticTypePass(ctx, expr.Type)
		expr.Expr, _ = semanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO type check
		return expr, t
	default:
		panic(expr)
	}
}

func semanticTypeListPass(ctx *semanticPassContext, nodes []ASTTypeRef) ([]ASTTypeRef, []core.DubType) {
	types := make([]core.DubType, len(nodes))
	for i, node := range nodes {
		nodes[i], types[i] = semanticTypePass(ctx, node)
	}
	return nodes, types
}

func semanticTypePass(ctx *semanticPassContext, node ASTTypeRef) (ASTTypeRef, core.DubType) {
	switch node := node.(type) {
	case *TypeRef:
		name := node.Name.Text
		d, ok := resolve(ctx, name)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return node, unresolvedType
		}
		t, ok := AsType(d)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
			return node, unresolvedType
		}
		return &GetType{Type: t}, t
	case *QualifiedTypeRef:
		mname := node.Package.Text
		pkg, ok := ctx.Module.Namespace[mname]
		if !ok {
			ctx.Status.LocationError(node.Package.Pos, fmt.Sprintf("Could not resolve name %#v", mname))
			return node, unresolvedType
		}
		scope, ok := AsPackage(pkg)
		if !ok {
			ctx.Status.LocationError(node.Package.Pos, fmt.Sprintf("%#v is not a package", mname))
			return node, unresolvedType
		}
		name := node.Name.Text
		d, ok := scope.Namespace[name]
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return node, unresolvedType
		}
		t, ok := AsType(d)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
			return node, unresolvedType
		}
		return &GetType{Type: t}, t
	case *ListTypeRef:
		var t core.DubType
		node.Type, t = semanticTypePass(ctx, node.Type)
		if t == unresolvedType {
			return node, unresolvedType
		}
		t = ctx.Memo.getList(t)
		return &GetType{Type: t}, t
	default:
		panic(node)
	}
}

func refLocation(node ASTTypeRef) int {
	switch node := node.(type) {
	case *TypeRef:
		return node.Name.Pos
	case *ListTypeRef:
		// TODO more precise location
		return refLocation(node.Type)
	default:
		panic(node)
	}
}

func semanticStructTypePass(ctx *semanticPassContext, node ASTTypeRef) *core.StructType {
	node, t := semanticTypePass(ctx, node)
	if t == unresolvedType {
		return nil
	}
	st, ok := t.(*core.StructType)
	if !ok {
		ctx.Status.LocationError(refLocation(node), "Not a structure type.")
		return nil
	}
	return st
}

func semanticBlockPass(ctx *semanticPassContext, decl *FuncDecl, block []ASTExpr, scope *semanticScope) {
	for i, expr := range block {
		block[i], _ = semanticExprPass(ctx, decl, expr, scope)
	}
}

func semanticFuncSignaturePass(ctx *semanticPassContext, decl *FuncDecl) {
	args := make([]core.DubType, len(decl.Params))
	for i, p := range decl.Params {
		p.Type, args[i] = semanticTypePass(ctx, p.Type)
	}
	var result core.DubType
	if len(decl.ReturnTypes) != 1 {
		results := make([]core.DubType, len(decl.ReturnTypes))
		for i, t := range decl.ReturnTypes {
			decl.ReturnTypes[i], results[i] = semanticTypePass(ctx, t)
		}
		result = ctx.Memo.getTuple(results)
	} else {
		decl.ReturnTypes[0], result = semanticTypePass(ctx, decl.ReturnTypes[0])
	}
	decl.F.Type = ctx.Memo.getFunctionType(args, result)
}

func semanticFuncBodyPass(ctx *semanticPassContext, decl *FuncDecl) {
	scope := childScope(nil)
	for _, p := range decl.Params {
		p.Info = createLocal(ctx, decl, p.Name, ResolveType(p.Type), scope)
	}
	semanticBlockPass(ctx, decl, decl.Block, scope)
}

func semanticStructPass(ctx *semanticPassContext, decl *StructDecl) {
	t := decl.T
	t.Name = decl.Name.Text
	t.Scoped = decl.Scoped
	t.Contains = make([]*core.StructType, len(decl.Contains))
	for i, c := range decl.Contains {
		t.Contains[i] = semanticStructTypePass(ctx, c)
	}
	if decl.Implements != nil {
		t.Implements = semanticStructTypePass(ctx, decl.Implements)
	}
	t.Fields = make([]*core.FieldType, len(decl.Fields))
	for i, f := range decl.Fields {
		var ft core.DubType
		f.Type, ft = semanticTypePass(ctx, f.Type)
		t.Fields[i] = &core.FieldType{
			Name: f.Name.Text,
			Type: ft,
		}
	}
}

func semanticDestructurePass(ctx *semanticPassContext, decl *FuncDecl, d Destructure, scope *semanticScope) core.DubType {
	switch d := d.(type) {
	case *DestructureStruct:
		var t core.DubType
		d.Type, t = semanticTypePass(ctx, d.Type)
		st, ok := t.(*core.StructType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			aft := semanticDestructurePass(ctx, decl, arg.Destructure, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := f.Type
					if !TypeMatches(aft, eft, false) {
						ctx.Status.GlobalError(fmt.Sprintf("%s.%s: %s vs. %s", core.TypeName(t), fn, core.TypeName(aft), core.TypeName(eft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", core.TypeName(t), fn))
				}
			}
		}
		return t
	case *DestructureList:
		var t core.DubType
		d.Type, t = semanticTypePass(ctx, d.Type)
		lt, ok := t.(*core.ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			at := semanticDestructurePass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.GlobalError(fmt.Sprintf("%s vs. %s", core.TypeName(at), core.TypeName(lt.Type)))
				}
			}
		}
		return t
	case *DestructureValue:
		var t core.DubType
		d.Expr, t = semanticExprPass(ctx, decl, d.Expr, scope)
		return t
	default:
		panic(d)
	}
}

func semanticTestPass(ctx *semanticPassContext, tst *Test) {
	// HACK no real context
	scope := childScope(nil)
	tst.Rule, tst.Type = semanticExprPass(ctx, nil, tst.Rule, scope)

	at := semanticDestructurePass(ctx, nil, tst.Destructure, scope)
	if !TypeMatches(at, tst.Type, false) {
		ctx.Status.GlobalError(fmt.Sprintf("destructure %s vs. %s", core.TypeName(at), core.TypeName(tst.Type)))
	}
}

type ProgramScope struct {
	Index     *core.BuiltinTypeIndex
	Namespace map[string]namedElement
}

type namedElement interface {
	isNamedElement()
}

type NamedType struct {
	T core.DubType
}

func (element *NamedType) isNamedElement() {
}

func AsType(node namedElement) (core.DubType, bool) {
	switch node := node.(type) {
	case *NamedType:
		return node.T, true
	default:
		return nil, false
	}
}

type NamedCallable struct {
	Func core.Callable
}

func (element *NamedCallable) isNamedElement() {
}

type NamedCallableTemplate struct {
	Func core.CallableTemplate
}

func (element *NamedCallableTemplate) isNamedElement() {
}

type NamedPackage struct {
	Scope *ModuleScope
}

func (element *NamedPackage) isNamedElement() {
}

func AsPackage(node namedElement) (*ModuleScope, bool) {
	switch node := node.(type) {
	case *NamedPackage:
		return node.Scope, true
	default:
		return nil, false
	}
}

type ModuleScope struct {
	Package   *core.Package
	Path      []string
	Namespace map[string]namedElement
}

func GetField(node *core.StructType, name string) *core.FieldType {
	for _, decl := range node.Fields {
		if decl.Name == name {
			return decl
		}
	}
	return nil
}

func ResolveType(ref ASTTypeRef) core.DubType {
	switch ref := ref.(type) {
	case *GetType:
		return ref.Type
	case *TypeRef, *QualifiedTypeRef, *ListTypeRef:
		return unresolvedType
	default:
		panic(ref)
	}
}

func ReturnType(ctx *semanticPassContext, node core.Callable, args []core.DubType) core.DubType {
	builtins := ctx.Program.Index

	// TODO check argument types
	switch node := node.(type) {
	case *core.Function:
		f := ctx.Functions[node.Index]
		if len(f.ReturnTypes) != 1 {
			types := make([]core.DubType, len(f.ReturnTypes))
			for i, t := range f.ReturnTypes {
				types[i] = ResolveType(t)
			}
			return ctx.Memo.getTuple(types)
		} else {
			return ResolveType(f.ReturnTypes[0])
		}
	case *core.IntrinsicFunction:
		switch node {
		case ctx.Program.Index.Position:
			if len(args) != 0 {
				panic(args)
			}
			return builtins.Int
		case ctx.Program.Index.Slice:
			if len(args) != 2 {
				panic(args)
			}
			return builtins.String
		default:
			panic(node)
		}
	default:
		panic(node)
	}
}

func makeBuiltinTypeIndex(memo *typeMemoizer) *core.BuiltinTypeIndex {
	index := &core.BuiltinTypeIndex{
		String:  &core.BuiltinType{Name: "string"},
		Rune:    &core.BuiltinType{Name: "rune"},
		Int:     &core.BuiltinType{Name: "int"},
		Int64:   &core.BuiltinType{Name: "int64"},
		Float32: &core.BuiltinType{Name: "float32"},
		Bool:    &core.BuiltinType{Name: "bool"},
		Graph:   &core.BuiltinType{Name: "graph"},
		Nil:     &core.NilType{},
	}

	index.Append = &core.IntrinsicFunctionTemplate{
		Name: "append",
		Params: []*core.TemplateParam{
			&core.TemplateParam{
				Name: "T",
			},
		},
		Type: memo.getFunctionType(
			[]core.DubType{
				memo.getList(memo.getUnbound(0)),
				memo.getUnbound(0),
			},
			memo.getList(memo.getUnbound(0)),
		),
	}
	index.Position = &core.IntrinsicFunction{
		Name: "position",
		Type: memo.getFunctionType(
			[]core.DubType{},
			index.Int,
		),
	}
	index.Slice = &core.IntrinsicFunction{
		Name: "slice",
		// TODO memoize
		Type: memo.getFunctionType(
			[]core.DubType{
				index.Int,
				index.Int,
			},
			index.String,
		),
	}
	return index
}

func addIntrinsicFunction(f *core.IntrinsicFunction, namespace map[string]namedElement) {
	namespace[f.Name] = &NamedCallable{
		Func: f,
	}
}

func addIntrinsicFunctionTemplate(f *core.IntrinsicFunctionTemplate, namespace map[string]namedElement) {
	namespace[f.Name] = &NamedCallableTemplate{
		Func: f,
	}
}

func addBuiltinType(t *core.BuiltinType, namespace map[string]namedElement) {
	namespace[t.Name] = &NamedType{
		T: t,
	}
}

func MakeProgramScope(program *Program) *ProgramScope {
	ns := map[string]namedElement{}
	programScope := &ProgramScope{
		Index:     program.Builtins,
		Namespace: ns,
	}
	builtins := program.Builtins

	addBuiltinType(builtins.String, ns)
	addBuiltinType(builtins.Rune, ns)
	addBuiltinType(builtins.Int, ns)
	addBuiltinType(builtins.Int64, ns)
	addBuiltinType(builtins.Float32, ns)
	addBuiltinType(builtins.Bool, ns)
	addBuiltinType(builtins.Graph, ns)

	addIntrinsicFunctionTemplate(builtins.Append, ns)
	addIntrinsicFunction(builtins.Position, ns)
	addIntrinsicFunction(builtins.Slice, ns)

	return programScope
}

type semanticPassContext struct {
	Program        *ProgramScope
	Module         *ModuleScope
	ModuleContexts []*semanticPassContext
	Status         compiler.PassStatus
	Core           *core.CoreProgram
	Functions      []*FuncDecl
	Memo           *typeMemoizer
	Void           *core.TupleType
}

func resolveImport(ctx *semanticPassContext, imp *ImportDecl) {
	pos := imp.Path.Pos
	path := imp.Path.Value
	parts := strings.Split(path, "/")

	// HACK O(n^2)
	for _, other := range ctx.ModuleContexts {
		otherPath := other.Module.Path
		found := strings.Join(otherPath, "/") == path
		if found {
			name := parts[len(parts)-1]
			// HACK should use file-local namespace.
			_, exists := ctx.Module.Namespace[name]
			if exists {
				ctx.Status.LocationError(pos, fmt.Sprintf("Tried to redefine %#v", name))
			} else {
				ctx.Module.Namespace[name] = &NamedPackage{Scope: other.Module}
			}
			return
		}
	}
	ctx.Status.LocationError(pos, fmt.Sprintf("cannot find module %#v", path))
}

func indexModule(ctx *semanticPassContext, pkg *Package) {
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			resolveImport(ctx, imp)
		}

		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Namespace[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					f := &core.Function{
						Name: name,
						File: file.F,
					}

					if len(decl.TemplateParams) == 0 {
						decl.F = ctx.Core.Function_Scope.Register(f)
						ctx.Functions = append(ctx.Functions, decl)

						ctx.Module.Namespace[name] = &NamedCallable{
							Func: decl.F,
						}
					} else {
						f := &core.FunctionTemplate{
							Name: name,
						}
						ctx.Module.Namespace[name] = &NamedCallableTemplate{
							Func: f,
						}
					}
				}
			case *StructDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Namespace[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					st := &core.StructType{
						File: file.F,
					}
					decl.T = st
					ctx.Core.Structures = append(ctx.Core.Structures, st)
					ctx.Module.Namespace[name] = &NamedType{T: st}
				}
			default:
				panic(decl)
			}
		}
	}
}

func resolveSignatures(ctx *semanticPassContext, pkg *Package) {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				if len(decl.TemplateParams) == 0 {
					// Needed for resolving calls in the next step.
					semanticFuncSignaturePass(ctx, decl)
				}
			case *StructDecl:
				// Needed for resolving field reference types.
				semanticStructPass(ctx, decl)
			default:
				panic(decl)
			}
		}
	}
}

func semanticModulePass(ctx *semanticPassContext, pkg *Package) {
	// Resolve the declaration contents.
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				semanticFuncBodyPass(ctx, decl)
			case *StructDecl:
			default:
				panic(decl)
			}
		}
		for _, tst := range file.Tests {
			semanticTestPass(ctx, tst)
		}
	}
}

func SemanticPass(program *Program, status compiler.PassStatus) *core.CoreProgram {
	status.Begin()
	defer status.End()

	memo := &typeMemoizer{
		Tuples:      makeTupleLUT(),
		Lists:       map[core.DubType]*core.ListType{},
		Specialized: map[specialization]core.Callable{},
		Funcs:       map[funcTypeKey]*core.FunctionType{},
	}
	voidType := memo.getTuple(nil)

	program.Builtins = makeBuiltinTypeIndex(memo)

	programScope := MakeProgramScope(program)
	coreProg := &core.CoreProgram{
		Builtins:       program.Builtins,
		Package_Scope:  &core.Package_Scope{},
		File_Scope:     &core.File_Scope{},
		Function_Scope: &core.Function_Scope{},
	}

	for _, pkg := range program.Packages {
		corePkg := &core.Package{Path: pkg.Path}
		packageRef := coreProg.Package_Scope.Register(corePkg)
		pkg.P = packageRef

		for _, file := range pkg.Files {
			coreFile := &core.File{Name: file.Name, Package: pkg.P}
			fileRef := coreProg.File_Scope.Register(coreFile)
			file.F = fileRef
			corePkg.Files = append(corePkg.Files, fileRef)
		}
	}

	ctxs := make([]*semanticPassContext, len(program.Packages))
	for i, pkg := range program.Packages {
		moduleScope := &ModuleScope{
			Package:   pkg.P,
			Path:      pkg.Path,
			Namespace: map[string]namedElement{},
		}
		ctxs[i] = &semanticPassContext{
			Program:        programScope,
			Module:         moduleScope,
			ModuleContexts: ctxs,
			Status:         status,
			Core:           coreProg,
			Memo:           memo,
			Void:           voidType,
		}
	}

	for i, pkg := range program.Packages {
		indexModule(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return nil
	}
	for i, pkg := range program.Packages {
		resolveSignatures(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return nil
	}
	for i, pkg := range program.Packages {
		semanticModulePass(ctxs[i], pkg)
	}
	return coreProg
}

func (scope *LocalInfo_Scope) Get(ref LocalInfo_Ref) *LocalInfo {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *LocalInfo_Scope) Register(info *LocalInfo) *LocalInfo {
	info.Index = LocalInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
}

func (scope *LocalInfo_Scope) Len() int {
	return len(scope.objects)
}
