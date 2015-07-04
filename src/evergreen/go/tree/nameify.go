package tree

import (
	"evergreen/compiler"
	"evergreen/go/core"
	"fmt"
	"strings"
)

type FileInfo struct {
	Package     *core.Package
	PackageName map[*core.Package]string
	Decl        *FuncDecl // HACK
	CoreProg    *core.CoreProgram
}

func DefaultPackageName(pkg *core.Package) string {
	n := len(pkg.Path)
	return pkg.Path[n-1]
}

func (info *FileInfo) ImportedName(pkg *core.Package) string {
	name, ok := info.PackageName[pkg]
	if !ok {
		// HACK assume no import name conflicts
		name = DefaultPackageName(pkg)
		info.PackageName[pkg] = name
	}
	return name
}

func (info *FileInfo) QualifyName(pkg *core.Package, name string) string {
	if pkg != nil && pkg != info.Package {
		return fmt.Sprintf("%s.%s", info.ImportedName(pkg), name)
	}
	return name
}

func (info *FileInfo) GetCallable(f core.Callable) Expr {
	switch f := f.(type) {
	case *core.Function:
		return &GetGlobal{
			Text: info.QualifyName(f.Package, f.Name),
		}
	case *core.IntrinsicFunction:
		return &GetGlobal{
			Text: f.Name,
		}
	default:
		panic(f)
	}
}

func (info *FileInfo) LocalName(lcl *LocalInfo) string {
	return lcl.Name
}

func nameifyParam(p *Param, info *FileInfo) {
	p.Name = info.LocalName(p.Info)
	p.Type = nameifyType(p.Info.T, info)
}

func nameifyType(t TypeRef, info *FileInfo) TypeRef {
	switch t := t.(type) {
	case *NameRef:
		impl := t.T
		switch impl := impl.(type) {
		case *core.StructType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.InterfaceType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.ExternalType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.TypeDefType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		default:
			panic(impl)
		}
	case *SliceRef:
		t.Element = nameifyType(t.Element, info)
	case *PointerRef:
		t.Element = nameifyType(t.Element, info)
	case *FuncTypeRef:
		for _, p := range t.Params {
			nameifyParam(p, info)
		}
		for _, r := range t.Results {
			nameifyParam(r, info)
		}
	default:
		panic(t)
	}
	return t
}

func nameifyExpr(expr Expr, info *FileInfo) Expr {
	switch expr := expr.(type) {
	case *GetLocal, *GetName, *GetGlobal:
		// TODO
	case *GetFunction:
		return info.GetCallable(expr.Func)
	case *UnaryExpr:
		expr.Expr = nameifyExpr(expr.Expr, info)
	case *BinaryExpr:
		expr.Left = nameifyExpr(expr.Left, info)
		expr.Right = nameifyExpr(expr.Right, info)
	case *Call:
		expr.Expr = nameifyExpr(expr.Expr, info)
		for i, e := range expr.Args {
			expr.Args[i] = nameifyExpr(e, info)
		}
	case *Selector:
		expr.Expr = nameifyExpr(expr.Expr, info)
	case *Index:
		expr.Expr = nameifyExpr(expr.Expr, info)
		expr.Index = nameifyExpr(expr.Index, info)
	case *TypeCoerce:
		expr.Type = nameifyType(expr.Type, info)
		expr.Expr = nameifyExpr(expr.Expr, info)
	case *TypeAssert:
		expr.Expr = nameifyExpr(expr.Expr, info)
		expr.Type = nameifyType(expr.Type, info)
	case *StructLiteral:
		expr.Type = nameifyType(expr.Type, info)
		for i, e := range expr.Args {
			expr.Args[i].Expr = nameifyExpr(e.Expr, info)
		}
	case *ListLiteral:
		expr.Type = nameifyType(expr.Type, info)
		for i, e := range expr.Args {
			expr.Args[i] = nameifyExpr(e, info)
		}
	case *IntLiteral, *Float32Literal, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral:
		// Leaf
	default:
		panic(expr)
	}
	return expr
}

func nameifyTarget(expr Target, info *FileInfo) {
	switch expr := expr.(type) {
	case *SetName, *SetLocal:
		// TODO
	default:
		panic(expr)
	}
}

func nameifyStmt(stmt Stmt, info *FileInfo) Stmt {
	switch stmt := stmt.(type) {
	case *Var:
		stmt.Name = info.LocalName(stmt.Info)
		stmt.Type = nameifyType(stmt.Type, info)
		if stmt.Expr != nil {
			stmt.Expr = nameifyExpr(stmt.Expr, info)
		}
	case *Assign:
		for i, e := range stmt.Sources {
			stmt.Sources[i] = nameifyExpr(e, info)
		}
		for _, e := range stmt.Targets {
			nameifyTarget(e, info)
		}
	case *If:
		stmt.Cond = nameifyExpr(stmt.Cond, info)
		nameifyBody(stmt.T, info)
		if stmt.F != nil {
			nameifyBody(stmt.F, info)
		}
	case *For:
		nameifyBody(stmt.Block, info)
	case *BlockStmt:
		nameifyBody(stmt.Block, info)
	case *Return:
		for i, e := range stmt.Args {
			stmt.Args[i] = nameifyExpr(e, info)
		}
	case *Goto:
		// TODO
	case *Label:
		// TODO
	default:
		// TODO unhack.
		expr, ok := stmt.(Expr)
		if ok {
			return nameifyExpr(expr, info)
		} else {
			panic(stmt)
		}
	}
	return stmt
}

func nameifyBody(block *Block, info *FileInfo) {
	for i, stmt := range block.Body {
		block.Body[i] = nameifyStmt(stmt, info)
	}
}

func baseName(lcl *LocalInfo) string {
	letters := []rune(lcl.Name)
	end := len(letters)
	for end > 0 && letters[end-1] >= '0' && letters[end-1] <= '9' {
		end -= 1
	}
	name := string(letters[:end])
	if name == "" {
		name = "r"
	}
	return name
}

func isGoKeyword(s string) bool {
	switch s {
	case "break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "got", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var":
		return true
	default:
		return false
	}
}

func nameifyFunc(decl *FuncDecl, info *FileInfo) {
	info.Decl = decl

	CompactFunc(decl)

	num_vars := decl.LocalInfo_Scope.Len()

	// Normalize variable names and count identical names.
	count := make(map[string]int, num_vars)
	iter := decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		name := baseName(lcl)
		lcl.Name = name
		count[name] += 1
	}

	// Disambiguate identical names.
	uids := make(map[string]int, num_vars)
	iter = decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		name := baseName(lcl)
		if count[name] > 1 || isGoKeyword(name) {
			uid, _ := uids[name]
			uids[name] += 1
			lcl.Name = fmt.Sprintf("%s%d", name, uid)
		}
	}

	iter = decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		lcl.T = nameifyType(lcl.T, info)
	}
	if decl.Recv != nil {
		nameifyParam(decl.Recv, info)
	}
	// HACK not writing back to decl.Type due to type widening.
	// Function types should not be rewritten, however.
	nameifyType(decl.Type, info)
	nameifyBody(decl.Block, info)
	InsertVarDecls(decl)
}

func nameifyDecl(decl Decl, info *FileInfo) {
	switch decl := decl.(type) {
	case *InterfaceDecl:
		for _, field := range decl.Fields {
			field.Type = nameifyType(field.Type, info)
		}
	case *StructDecl:
		for _, field := range decl.Fields {
			field.Type = nameifyType(field.Type, info)
		}
	case *TypeDefDecl:
		decl.Type = nameifyType(decl.Type, info)
	case *FuncDecl:
		nameifyFunc(decl, info)
	case *VarDecl:
		decl.Type = nameifyType(decl.Type, info)
		decl.Expr = nameifyExpr(decl.Expr, info)
	default:
		panic(decl)
	}
}

func nameifyFile(pkg *PackageAST, file *FileAST, coreProg *core.CoreProgram) {
	p := pkg.P
	file.Package = p.Path[len(p.Path)-1]

	info := &FileInfo{
		Package:     p,
		PackageName: map[*core.Package]string{},
		CoreProg:    coreProg,
	}

	for _, decl := range file.Decls {
		nameifyDecl(decl, info)
	}

	if len(file.Imports) != 0 {
		panic(file.Imports)
	}
	for pkg, name := range info.PackageName {
		path := strings.Join(pkg.Path, "/")
		file.Imports = append(file.Imports, &Import{Path: path, Name: name})
	}
}

func Nameify(status compiler.PassStatus, prog *ProgramAST, coreProg *core.CoreProgram) {
	status.Begin()
	defer status.End()

	nameifyPrepass(prog)
	for _, pkg := range prog.Packages {
		p := pkg.P
		pkgName := ""
		if len(p.Path) > 0 {
			pkgName = p.Path[len(p.Path)-1]
		}
		for _, file := range pkg.Files {
			file.Package = pkgName
			nameifyFile(pkg, file, coreProg)
		}
	}
}

func nameifyPrepass(prog *ProgramAST) {
	for _, pkg := range prog.Packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *InterfaceDecl:
					decl.T.Package = pkg.P
				case *StructDecl:
					decl.T.Package = pkg.P
				case *FuncDecl:
					decl.Package = pkg.P
				case *OpaqueDecl:
					decl.T.Package = pkg.P
				case *TypeDefDecl:
					decl.T.Package = pkg.P
				case *VarDecl:
					// Leaf
				default:
					panic(decl)
				}
			}
		}
	}
}
